"""
Tests for authentication utilities
"""

import os
import pytest
from unittest.mock import mock_open, patch

from model_registry_mlflow.auth import get_auth_headers, _get_k8s_token, _token_cache


class TestAuth:
    def setup_method(self):
        """Reset token cache before each test."""
        # Reset the cache state for testing
        _token_cache._token = None
        _token_cache._token_mtime = None

    def test_get_auth_headers_env_token_success(self):
        """Test successful authentication with environment token."""
        with patch.dict(os.environ, {"MODEL_REGISTRY_TOKEN": "test-env-token"}):
            headers = {"Content-Type": "application/json"}
            get_auth_headers(headers)

            assert headers["Authorization"] == "Bearer test-env-token"
            assert headers["Content-Type"] == "application/json"

    @patch("os.path.getmtime", return_value=1000.0)
    @patch("builtins.open", new_callable=mock_open, read_data="k8s-token")
    def test_get_auth_headers_k8s_token_success(self, mock_file, mock_getmtime):
        """Test successful authentication with Kubernetes token."""
        with patch.dict(os.environ, {}, clear=True):
            headers = {"Content-Type": "application/json"}
            get_auth_headers(headers)

            assert headers["Authorization"] == "Bearer k8s-token"
            mock_file.assert_called_once_with(
                "/var/run/secrets/kubernetes.io/serviceaccount/token", "r"
            )

    def test_get_auth_headers_env_token_priority(self):
        """Test that environment token takes priority over Kubernetes token."""
        with patch.dict(os.environ, {"MODEL_REGISTRY_TOKEN": "env-token"}):
            with patch("model_registry_mlflow.auth._get_k8s_token") as mock_k8s:
                mock_k8s.return_value = "k8s-token"

                headers = {}
                get_auth_headers(headers)

                assert headers["Authorization"] == "Bearer env-token"
                # K8s token function should not be called when env token exists
                mock_k8s.assert_not_called()

    def test_get_auth_headers_no_auth_raises_error(self):
        """Test RuntimeError when no authentication is available."""
        with patch.dict(os.environ, {}, clear=True):
            with patch("model_registry_mlflow.auth._get_k8s_token", return_value=None):
                headers = {}

                with pytest.raises(
                    RuntimeError, match="No authentication token available"
                ):
                    get_auth_headers(headers)

    @patch("builtins.open", side_effect=FileNotFoundError("No such file or directory"))
    def test_get_k8s_token_file_not_found(self, mock_file):
        """Test RuntimeError when token file doesn't exist."""
        with pytest.raises(
            RuntimeError,
            match="Error accessing Kubernetes token file.*No such file or directory",
        ):
            _get_k8s_token()

    @patch("os.path.getmtime", side_effect=PermissionError("Permission denied"))
    def test_get_k8s_token_permission_error(self, mock_getmtime):
        """Test RuntimeError when permission denied accessing token file."""
        with pytest.raises(
            RuntimeError,
            match="Error accessing Kubernetes token file.*Permission denied",
        ):
            _get_k8s_token()

    @patch("os.path.getmtime", return_value=1000.0)
    @patch("builtins.open", new_callable=mock_open)
    def test_get_k8s_token_unicode_decode_error(self, mock_file, mock_getmtime):
        """Test RuntimeError when token file contains invalid UTF-8."""
        mock_file.return_value.read.side_effect = UnicodeDecodeError(
            "utf-8", b"\x80\x81", 0, 1, "invalid start byte"
        )

        with pytest.raises(
            RuntimeError,
            match="Error accessing Kubernetes token file.*invalid start byte",
        ):
            _get_k8s_token()

    @patch("os.path.getmtime", side_effect=IOError("I/O error"))
    def test_get_k8s_token_io_error(self, mock_getmtime):
        """Test RuntimeError when I/O error occurs reading token file."""
        with pytest.raises(
            RuntimeError, match="Error accessing Kubernetes token file.*I/O error"
        ):
            _get_k8s_token()

    @patch("os.path.getmtime", return_value=1000.0)
    @patch("builtins.open", new_callable=mock_open, read_data="cached-token")
    def test_k8s_token_caching_behavior(self, mock_file, mock_getmtime):
        """Test that Kubernetes token is cached efficiently."""
        # First call should read the file
        token1 = _get_k8s_token()
        assert token1 == "cached-token"
        assert mock_file.call_count == 1

        # Second call should use cached token (no additional file reads)
        token2 = _get_k8s_token()
        assert token2 == "cached-token"
        assert mock_file.call_count == 1  # File should not be read again
        assert mock_getmtime.call_count == 2  # mtime checked both times

    def test_get_auth_headers_modifies_existing_headers(self):
        """Test that get_auth_headers modifies existing headers dictionary."""
        with patch.dict(os.environ, {"MODEL_REGISTRY_TOKEN": "test-token"}):
            headers = {
                "Content-Type": "application/json",
                "Accept": "application/json",
                "Custom-Header": "custom-value",
            }

            get_auth_headers(headers)

            assert headers["Authorization"] == "Bearer test-token"
            assert (
                headers["Custom-Header"] == "custom-value"
            )  # Existing headers preserved
            assert len(headers) == 4  # Authorization added to existing headers

    @patch("os.path.getmtime")
    @patch("builtins.open", new_callable=mock_open, read_data="new-token")
    def test_k8s_token_mtime_based_caching(self, mock_file, mock_getmtime):
        """Test that token caching works correctly based on file modification time."""
        # First call - file mtime is 1000.0
        mock_getmtime.return_value = 1000.0

        token1 = _get_k8s_token()
        assert token1 == "new-token"
        assert mock_file.call_count == 1
        assert mock_getmtime.call_count == 1  # Only checked once before read

        # Second call - same mtime, should use cache
        mock_getmtime.reset_mock()
        mock_file.reset_mock()
        mock_getmtime.return_value = 1000.0

        token2 = _get_k8s_token()
        assert token2 == "new-token"
        assert mock_file.call_count == 0  # File not read
        assert mock_getmtime.call_count == 1  # Only checked mtime once

        # Third call - different mtime, should refresh
        mock_file.return_value.read.return_value = "updated-token"
        mock_getmtime.return_value = 2000.0

        token3 = _get_k8s_token()
        assert token3 == "updated-token"
        assert mock_file.call_count == 1  # File read again
        assert mock_getmtime.call_count == 2  # Checked once more

    @patch("os.path.getmtime", side_effect=OSError("File not found"))
    def test_k8s_token_mtime_error_clears_cache(self, mock_getmtime):
        """Test that mtime errors clear the cache properly."""
        # Set up initial cache state
        _token_cache._token = "old-token"
        _token_cache._token_mtime = 1000.0

        # Call should clear cache and raise error
        with pytest.raises(
            RuntimeError, match="Error accessing Kubernetes token file.*File not found"
        ):
            _get_k8s_token()

        # Verify cache was cleared
        assert _token_cache._token is None
        assert _token_cache._token_mtime is None

    @patch("os.path.getmtime")
    @patch(
        "builtins.open",
        side_effect=UnicodeDecodeError("utf-8", b"\x80", 0, 1, "invalid"),
    )
    def test_k8s_token_read_error_clears_cache(self, mock_file, mock_getmtime):
        """Test that file read errors clear the cache properly."""
        mock_getmtime.return_value = 1000.0

        # Set up initial cache state
        _token_cache._token = "old-token"
        _token_cache._token_mtime = 500.0

        # Call should clear cache and raise error
        with pytest.raises(
            RuntimeError, match="Error accessing Kubernetes token file.*invalid"
        ):
            _get_k8s_token()

        # Verify cache was cleared
        assert _token_cache._token is None
        assert _token_cache._token_mtime is None

    @patch("os.path.getmtime")
    @patch("builtins.open", new_callable=mock_open)
    def test_k8s_token_simple_caching_behavior(self, mock_file, mock_getmtime):
        """Test that token caching stores the initial mtime value."""
        # Simple case - single mtime check
        mock_getmtime.return_value = 1000.0
        mock_file.return_value.read.return_value = "token-content"

        token = _get_k8s_token()
        assert token == "token-content"

        # Cache should store the initial mtime
        assert _token_cache._token_mtime == 1000.0
        assert _token_cache._token == "token-content"

        # Verify only one mtime check was made
        assert mock_getmtime.call_count == 1
