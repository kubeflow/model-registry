"""
Tests for authentication utilities
"""

from unittest.mock import patch, mock_open
import os

from modelregistry_plugin.auth import get_auth_headers


class TestAuth:
    def test_get_auth_headers_no_auth(self):
        """Test getting headers with no authentication."""
        with patch.dict(os.environ, {}, clear=True):
            headers = get_auth_headers()

            expected = {
                "Content-Type": "application/json",
                "Accept": "application/json",
            }
            assert headers == expected

    def test_get_auth_headers_token(self):
        """Test getting headers with token authentication."""
        with patch.dict(os.environ, {"MODEL_REGISTRY_TOKEN": "test-token"}):
            headers = get_auth_headers()

            assert headers["Authorization"] == "Bearer test-token"
            assert headers["Content-Type"] == "application/json"

    @patch("os.path.exists")
    @patch("builtins.open", new_callable=mock_open, read_data="k8s-token")
    def test_get_auth_headers_k8s_token(self, mock_file, mock_exists):
        """Test getting headers with Kubernetes service account token."""
        mock_exists.return_value = True

        with patch.dict(os.environ, {}, clear=True):
            headers = get_auth_headers()

            assert headers["Authorization"] == "Bearer k8s-token"
            mock_exists.assert_called_once_with(
                "/var/run/secrets/kubernetes.io/serviceaccount/token"
            )
            mock_file.assert_called_once_with(
                "/var/run/secrets/kubernetes.io/serviceaccount/token", "r"
            )
