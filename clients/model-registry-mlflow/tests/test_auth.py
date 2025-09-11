"""
Tests for authentication utilities
"""

import os
import time
from pathlib import Path
import pytest
from unittest.mock import patch

from model_registry_mlflow.auth import get_auth_headers, _get_k8s_token, K8sTokenCache
from .conftest import create_temp_file


class TestAuth:
    @pytest.fixture
    def temp_token_file(self):
        """Create a temporary token file and mock _TOKEN_PATH."""
        # Create temp file path using conftest helper
        temp_path = create_temp_file(suffix=".token", prefix="k8s_token_")

        # Create the file with initial content
        Path(temp_path).write_text("initial-token")

        # Mock the _TOKEN_PATH to point to our temp file
        with patch("model_registry_mlflow.auth._TOKEN_PATH", temp_path):
            yield temp_path

    def setup_method(self):
        """Reset token cache before each test."""
        # Reset the cache state for testing
        K8sTokenCache._token = None
        K8sTokenCache._token_mtime = None

    def test_get_auth_headers_env_token_success(self):
        """Test successful authentication with environment token."""
        with patch.dict(os.environ, {"MODEL_REGISTRY_TOKEN": "test-env-token"}):
            headers = {"Content-Type": "application/json"}
            get_auth_headers(headers)

            assert headers["Authorization"] == "Bearer test-env-token"
            assert headers["Content-Type"] == "application/json"

    def test_get_auth_headers_k8s_token_success(self, temp_token_file):
        """Test successful authentication with Kubernetes token."""
        # Write token to temp file
        Path(temp_token_file).write_text("k8s-token")

        with patch.dict(os.environ, {}, clear=True):
            headers = {"Content-Type": "application/json"}
            get_auth_headers(headers)

            assert headers["Authorization"] == "Bearer k8s-token"

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

    def test_get_k8s_token_file_not_found(self):
        """Test RuntimeError when token file doesn't exist."""
        # Mock _TOKEN_PATH to point to non-existent file
        with patch("model_registry_mlflow.auth._TOKEN_PATH", "/nonexistent/path/token"):
            with pytest.raises(
                RuntimeError,
                match="Error accessing Kubernetes token file.*No such file or directory",
            ):
                _get_k8s_token()

    def test_get_k8s_token_permission_error(self):
        """Test RuntimeError when permission denied accessing token file."""
        # Mock _TOKEN_PATH to point to file with no permissions
        with patch("model_registry_mlflow.auth._TOKEN_PATH", "/root/restricted_file"):
            with pytest.raises(
                RuntimeError,
                match="Error accessing Kubernetes token file.*Permission denied",
            ):
                _get_k8s_token()

    def test_get_k8s_token_unicode_decode_error(self, temp_token_file):
        """Test RuntimeError when token file contains invalid UTF-8."""
        # Write invalid UTF-8 bytes to temp file
        Path(temp_token_file).write_bytes(b"\x80\x81")

        with pytest.raises(
            RuntimeError,
            match="Error accessing Kubernetes token file.*invalid",
        ):
            _get_k8s_token()

    @patch("os.path.getmtime", side_effect=IOError("I/O error"))
    def test_get_k8s_token_io_error(self, mock_getmtime):
        """Test RuntimeError when I/O error occurs reading token file."""
        with pytest.raises(
            RuntimeError, match="Error accessing Kubernetes token file.*I/O error"
        ):
            _get_k8s_token()

    def test_k8s_token_caching_behavior(self, temp_token_file):
        """Test that Kubernetes token is cached efficiently."""
        # Write token to temp file
        Path(temp_token_file).write_text("cached-token")
        initial_mtime = os.path.getmtime(temp_token_file)

        # First call should read the file
        token1 = _get_k8s_token()
        assert token1 == "cached-token"

        # Second call should use cached token (no additional file reads needed)
        token2 = _get_k8s_token()
        assert token2 == "cached-token"

        # Verify mtime hasn't changed (cache should be used)
        assert os.path.getmtime(temp_token_file) == initial_mtime

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

    def test_k8s_token_mtime_based_caching(self, temp_token_file):
        """Test that token caching works correctly based on file modification time."""
        # First call - write initial token
        Path(temp_token_file).write_text("new-token")

        token1 = _get_k8s_token()
        assert token1 == "new-token"

        # Second call - same file, should use cache
        token2 = _get_k8s_token()
        assert token2 == "new-token"

        # Third call - modify file (change mtime), should refresh
        time.sleep(0.01)  # Ensure different mtime
        Path(temp_token_file).write_text("updated-token")

        token3 = _get_k8s_token()
        assert token3 == "updated-token"

    @patch("os.path.getmtime", side_effect=OSError("File not found"))
    def test_k8s_token_mtime_error_clears_cache(self, mock_getmtime):
        """Test that mtime errors clear the cache properly."""
        # Set up initial cache state
        K8sTokenCache._token = "old-token"
        K8sTokenCache._token_mtime = 1000.0

        # Call should clear cache and raise error
        with pytest.raises(
            RuntimeError, match="Error accessing Kubernetes token file.*File not found"
        ):
            _get_k8s_token()

        # Verify cache was cleared
        assert K8sTokenCache._token is None
        assert K8sTokenCache._token_mtime is None

    def test_k8s_token_read_error_clears_cache(self, temp_token_file):
        """Test that file read errors clear the cache properly."""
        # Set up initial cache state
        K8sTokenCache._token = "old-token"
        K8sTokenCache._token_mtime = 500.0

        # Write invalid UTF-8 to cause read error
        Path(temp_token_file).write_bytes(b"\x80\x81")

        # Call should clear cache and raise error
        with pytest.raises(
            RuntimeError, match="Error accessing Kubernetes token file.*invalid"
        ):
            _get_k8s_token()

        # Verify cache was cleared
        assert K8sTokenCache._token is None
        assert K8sTokenCache._token_mtime is None

    def test_k8s_token_simple_caching_behavior(self, temp_token_file):
        """Test that token caching stores the initial mtime value."""
        # Write token to temp file
        Path(temp_token_file).write_text("token-content")
        expected_mtime = os.path.getmtime(temp_token_file)

        token = _get_k8s_token()
        assert token == "token-content"

        # Cache should store the actual file mtime
        assert K8sTokenCache._token_mtime == expected_mtime
        assert K8sTokenCache._token == "token-content"
