"""Tests for ModelRegistryAPIClient."""

import os
from unittest.mock import Mock, patch

import pytest
import requests
from mlflow.exceptions import MlflowException

from modelregistry_plugin.api_client import ModelRegistryAPIClient


class TestModelRegistryAPIClient:
    @pytest.fixture
    def api_client(self):
        """Create a ModelRegistryAPIClient instance for testing."""
        return ModelRegistryAPIClient("http://localhost:8080")

    @pytest.fixture
    def mock_response(self):
        """Create a mock response object."""
        response = Mock(spec=requests.Response)
        response.ok = True
        response.json.return_value = {}
        return response

    def test_init(self, api_client):
        """Test API client initialization."""
        assert api_client.base_url == "http://localhost:8080"
        assert isinstance(api_client.session, requests.Session)

    def test_init_with_trailing_slash(self):
        """Test API client initialization with trailing slash."""
        client = ModelRegistryAPIClient("http://localhost:8080/")
        assert client.base_url == "http://localhost:8080"

    def test_init_without_protocol(self):
        """Test API client initialization without protocol."""
        client = ModelRegistryAPIClient("localhost:8080")
        assert client.base_url == "localhost:8080"

    @patch("modelregistry_plugin.api_client.get_auth_headers")
    @patch("modelregistry_plugin.api_client.requests.Session.request")
    def test_request_success(
        self, mock_session_request, mock_auth_headers, api_client, mock_response
    ):
        """Test successful API request."""
        mock_auth_headers.return_value = {"Authorization": "Bearer token"}
        mock_session_request.return_value = mock_response

        response_data = api_client.request("GET", "/test")

        mock_session_request.assert_called_once()
        call_args = mock_session_request.call_args
        assert call_args[0][0] == "GET"  # method
        assert call_args[0][1] == "http://localhost:8080/test"  # url
        assert call_args[1]["headers"]["Authorization"] == "Bearer token"
        assert response_data == mock_response.json.return_value

    @patch("modelregistry_plugin.api_client.get_auth_headers")
    @patch("modelregistry_plugin.api_client.requests.Session.request")
    def test_request_with_json_data(
        self, mock_session_request, mock_auth_headers, api_client, mock_response
    ):
        """Test API request with JSON data."""
        mock_auth_headers.return_value = {"Authorization": "Bearer token"}
        mock_session_request.return_value = mock_response

        json_data = {"customProperties": {"key1": "value1"}}
        api_client.request("POST", "/test", json=json_data)

        call_args = mock_session_request.call_args
        assert (
            call_args[1]["json"]["customProperties"]["key1"]["string_value"] == "value1"
        )
        assert (
            call_args[1]["json"]["customProperties"]["key1"]["metadataType"]
            == "MetadataStringValue"
        )

    @patch("modelregistry_plugin.api_client.get_auth_headers")
    @patch("modelregistry_plugin.api_client.requests.Session.request")
    def test_request_with_response_items(
        self, mock_session_request, mock_auth_headers, api_client
    ):
        """Test API request with response containing items."""
        mock_auth_headers.return_value = {"Authorization": "Bearer token"}

        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {
            "items": [
                {
                    "customProperties": {
                        "key1": {
                            "string_value": "value1",
                            "metadataType": "MetadataStringValue",
                        }
                    }
                }
            ]
        }
        mock_session_request.return_value = mock_response

        response_data = api_client.request("GET", "/test")

        # Check that customProperties were converted back to MLflow format
        assert response_data["items"][0]["customProperties"]["key1"] == "value1"

    @patch("modelregistry_plugin.api_client.get_auth_headers")
    @patch("modelregistry_plugin.api_client.requests.Session.request")
    def test_request_with_single_response(
        self, mock_session_request, mock_auth_headers, api_client
    ):
        """Test API request with single response object."""
        mock_auth_headers.return_value = {"Authorization": "Bearer token"}

        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {
            "customProperties": {
                "key1": {
                    "string_value": "value1",
                    "metadataType": "MetadataStringValue",
                }
            }
        }
        mock_session_request.return_value = mock_response

        response_data = api_client.request("GET", "/test")

        # Check that customProperties were converted back to MLflow format
        assert response_data["customProperties"]["key1"] == "value1"

    @patch("modelregistry_plugin.api_client.get_auth_headers")
    @patch("modelregistry_plugin.api_client.requests.Session.request")
    def test_request_http_error(
        self, mock_session_request, mock_auth_headers, api_client
    ):
        """Test API request with HTTP error."""
        mock_auth_headers.return_value = {"Authorization": "Bearer token"}

        mock_response = Mock(spec=requests.Response)
        mock_response.ok = False
        mock_response.status_code = 404
        mock_response.json.return_value = {"message": "Not found"}
        mock_response.text = "Not found"

        mock_session_request.side_effect = requests.exceptions.HTTPError(
            response=mock_response
        )

        with pytest.raises(MlflowException) as exc_info:
            api_client.request("GET", "/test")

        assert "Model Registry API error: Not found" in str(exc_info.value)

    @patch("modelregistry_plugin.api_client.get_auth_headers")
    @patch("modelregistry_plugin.api_client.requests.Session.request")
    def test_request_network_error(
        self, mock_session_request, mock_auth_headers, api_client
    ):
        """Test API request with network error."""
        mock_auth_headers.return_value = {"Authorization": "Bearer token"}

        mock_session_request.side_effect = requests.exceptions.ConnectionError(
            "Connection failed"
        )

        with pytest.raises(MlflowException) as exc_info:
            api_client.request("GET", "/test")

        assert "Network error connecting to Model Registry: Connection failed" in str(
            exc_info.value
        )

    @patch("modelregistry_plugin.api_client.get_auth_headers")
    @patch("modelregistry_plugin.api_client.requests.Session.request")
    def test_request_invalid_json_response(
        self, mock_session_request, mock_auth_headers, api_client
    ):
        """Test API request with invalid JSON response."""
        mock_auth_headers.return_value = {"Authorization": "Bearer token"}

        mock_response = Mock(spec=requests.Response)
        mock_response.ok = False
        mock_response.status_code = 400
        mock_response.json.side_effect = ValueError("Invalid JSON")
        mock_response.text = "Invalid JSON"

        mock_session_request.side_effect = requests.exceptions.HTTPError(
            response=mock_response
        )

        with pytest.raises(MlflowException) as exc_info:
            api_client.request("GET", "/test")

        assert "Model Registry API error: Invalid JSON" in str(exc_info.value)

    @patch("modelregistry_plugin.api_client.get_auth_headers")
    @patch("modelregistry_plugin.api_client.requests.Session.request")
    def test_get_method(
        self, mock_session_request, mock_auth_headers, api_client, mock_response
    ):
        """Test GET method."""
        mock_auth_headers.return_value = {"Authorization": "Bearer token"}
        mock_session_request.return_value = mock_response

        api_client.get("/test", params={"key": "value"})

        call_args = mock_session_request.call_args
        assert call_args[0][0] == "GET"
        assert call_args[1]["params"]["key"] == "value"

    @patch("modelregistry_plugin.api_client.get_auth_headers")
    @patch("modelregistry_plugin.api_client.requests.Session.request")
    def test_post_method(
        self, mock_session_request, mock_auth_headers, api_client, mock_response
    ):
        """Test POST method."""
        mock_auth_headers.return_value = {"Authorization": "Bearer token"}
        mock_session_request.return_value = mock_response

        api_client.post("/test", json={"data": "value"})

        call_args = mock_session_request.call_args
        assert call_args[0][0] == "POST"
        assert call_args[1]["json"]["data"] == "value"

    @patch("modelregistry_plugin.api_client.get_auth_headers")
    @patch("modelregistry_plugin.api_client.requests.Session.request")
    def test_patch_method(
        self, mock_session_request, mock_auth_headers, api_client, mock_response
    ):
        """Test PATCH method."""
        mock_auth_headers.return_value = {"Authorization": "Bearer token"}
        mock_session_request.return_value = mock_response

        api_client.patch("/test", json={"data": "value"})

        call_args = mock_session_request.call_args
        assert call_args[0][0] == "PATCH"
        assert call_args[1]["json"]["data"] == "value"

    @patch("modelregistry_plugin.api_client.get_auth_headers")
    @patch("modelregistry_plugin.api_client.requests.Session.request")
    def test_delete_method(
        self, mock_session_request, mock_auth_headers, api_client, mock_response
    ):
        """Test DELETE method."""
        mock_auth_headers.return_value = {"Authorization": "Bearer token"}
        mock_session_request.return_value = mock_response

        api_client.delete("/test")

        call_args = mock_session_request.call_args
        assert call_args[0][0] == "DELETE"

    @patch("modelregistry_plugin.api_client.get_auth_headers")
    @patch("modelregistry_plugin.api_client.requests.Session.request")
    def test_request_with_custom_headers(
        self, mock_session_request, mock_auth_headers, api_client, mock_response
    ):
        """Test API request with custom headers."""
        mock_auth_headers.return_value = {"Authorization": "Bearer token"}
        mock_session_request.return_value = mock_response

        api_client.request("GET", "/test", headers={"Custom-Header": "value"})

        call_args = mock_session_request.call_args
        headers = call_args[1]["headers"]
        assert headers["Authorization"] == "Bearer token"
        assert headers["Custom-Header"] == "value"

    @patch("modelregistry_plugin.api_client.get_auth_headers")
    @patch("modelregistry_plugin.api_client.requests.Session.request")
    def test_request_with_params(
        self, mock_session_request, mock_auth_headers, api_client, mock_response
    ):
        """Test API request with query parameters."""
        mock_auth_headers.return_value = {"Authorization": "Bearer token"}
        mock_session_request.return_value = mock_response

        api_client.request("GET", "/test", params={"key": "value"})

        call_args = mock_session_request.call_args
        assert call_args[1]["params"]["key"] == "value"

    def test_retry_strategy_configured(self, api_client):
        """Test that retry strategy is properly configured."""
        # Check that HTTPAdapter is mounted with retry strategy
        adapter = api_client.session.get_adapter("http://")
        assert adapter.max_retries.total == 3
        assert adapter.max_retries.backoff_factor == 1
        assert adapter.max_retries.status_forcelist == [429, 500, 502, 503, 504]

        adapter = api_client.session.get_adapter("https://")
        assert adapter.max_retries.total == 3
        assert adapter.max_retries.backoff_factor == 1
        assert adapter.max_retries.status_forcelist == [429, 500, 502, 503, 504]

    def test_ca_cert_configuration_http(self):
        """Test that CA certificate configuration is skipped for HTTP URLs."""
        with patch("modelregistry_plugin.api_client.logger") as mock_logger:
            client = ModelRegistryAPIClient("http://localhost:8080")
            assert client.session.verify is True  # Default requests behavior
            mock_logger.debug.assert_called_with(
                "Using HTTP connection, skipping CA certificate configuration"
            )

    def test_ca_cert_configuration_explicit_path(self):
        """Test CA certificate configuration with explicit path."""
        with patch("pathlib.Path.exists", return_value=True):
            with patch("modelregistry_plugin.api_client.logger") as mock_logger:
                client = ModelRegistryAPIClient(
                    "https://localhost:8080", ca_cert_path="/path/to/cert.pem"
                )
                assert client.session.verify == "/path/to/cert.pem"
                mock_logger.info.assert_called_with(
                    "Configuring custom CA certificate: /path/to/cert.pem"
                )

    def test_ca_cert_configuration_explicit_path_not_found(self):
        """Test CA certificate configuration with explicit path that doesn't exist."""
        with patch("pathlib.Path.exists", return_value=False):
            with patch("modelregistry_plugin.api_client.logger") as mock_logger:
                client = ModelRegistryAPIClient(
                    "https://localhost:8080", ca_cert_path="/path/to/nonexistent.pem"
                )
                assert client.session.verify is True  # Should fall back to default
                mock_logger.warning.assert_called_with(
                    "CA certificate file not found: /path/to/nonexistent.pem, falling back to system CA"
                )

    @patch.dict(os.environ, {"MODELREGISTRY_CA_CERT_PATH": "/env/path/to/cert.pem"})
    def test_ca_cert_configuration_env_var(self):
        """Test CA certificate configuration via environment variable."""
        with patch("pathlib.Path.exists", return_value=True):
            with patch("modelregistry_plugin.api_client.logger") as mock_logger:
                client = ModelRegistryAPIClient("https://localhost:8080")
                assert client.session.verify == "/env/path/to/cert.pem"
                # Check that the environment variable message was logged (first call)
                mock_logger.info.assert_any_call(
                    "Using CA certificate from environment variable MODELREGISTRY_CA_CERT_PATH: /env/path/to/cert.pem"
                )

    def test_ca_cert_configuration_kubernetes_ca_file_exists(self):
        """Test CA certificate configuration when Kubernetes CA file exists."""
        with patch("modelregistry_plugin.api_client.Path") as mock_path:
            # Mock the Path constructor and its exists method
            mock_path_instance = Mock()
            mock_path_instance.exists.return_value = True
            mock_path.return_value = mock_path_instance

            with patch("modelregistry_plugin.api_client.logger") as mock_logger:
                client = ModelRegistryAPIClient("https://localhost:8080")
                assert (
                    client.session.verify
                    == "/run/secrets/kubernetes.io/serviceaccount/ca.crt"
                )
                mock_logger.info.assert_any_call(
                    "Using Kubernetes default CA certificate: /run/secrets/kubernetes.io/serviceaccount/ca.crt"
                )

    def test_ca_cert_configuration_default_system_ca(self):
        """Test CA certificate configuration falls back to system CA."""
        with patch("pathlib.Path.exists", return_value=False):
            with patch("modelregistry_plugin.api_client.logger") as mock_logger:
                client = ModelRegistryAPIClient("https://localhost:8080")
                assert client.session.verify is True  # System default
                mock_logger.debug.assert_called_with("Using system default CA bundle")

    @patch("modelregistry_plugin.api_client.get_auth_headers")
    @patch("modelregistry_plugin.api_client.requests.Session.request")
    def test_request_ssl_error(
        self, mock_session_request, mock_auth_headers, api_client
    ):
        """Test API request with SSL error."""
        mock_auth_headers.return_value = {"Authorization": "Bearer token"}
        mock_session_request.side_effect = requests.exceptions.SSLError("SSL Error")

        with pytest.raises(MlflowException) as exc_info:
            api_client.request("GET", "/test")

        assert (
            "SSL certificate verification failed connecting to Model Registry: SSL Error"
            in str(exc_info.value)
        )

    def test_api_client_with_ca_cert_parameter(self):
        """Test API client initialization with ca_cert_path parameter."""
        with patch("pathlib.Path.exists", return_value=True):
            client = ModelRegistryAPIClient(
                "https://localhost:8080", ca_cert_path="/custom/ca.pem"
            )
            assert client.session.verify == "/custom/ca.pem"

    def test_api_client_constructor_backward_compatibility(self):
        """Test that the original constructor signature still works."""
        client = ModelRegistryAPIClient("http://localhost:8080")
        assert client.base_url == "http://localhost:8080"
        assert isinstance(client.session, requests.Session)
