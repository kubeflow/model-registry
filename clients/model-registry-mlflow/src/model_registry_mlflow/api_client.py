"""Model Registry API Client for handling HTTP communication."""

from __future__ import annotations

import os
import logging
from pathlib import Path
from typing import Any, Dict, Optional

import requests

from model_registry_mlflow.auth import get_auth_headers
from model_registry_mlflow.utils import (
    fromModelRegistryCustomProperties,
    toModelRegistryCustomProperties,
)

logger = logging.getLogger(__name__)

# Default paths for CA certificates
DEFAULT_K8S_CA_CERT_PATH = "/run/secrets/kubernetes.io/serviceaccount/ca.crt"
DEFAULT_CA_CERT_ENV_VAR = "MODELREGISTRY_CA_CERT_PATH"


class ModelRegistryAPIClient:
    """Handles all HTTP communication with Model Registry."""

    def __init__(self, base_url: str, ca_cert_path: Optional[str] = None):
        """Initialize the API client.

        Args:
            base_url: Base URL for the Model Registry API
            ca_cert_path: Path to CA certificate file. If None, will attempt to auto-detect.
        """
        self.base_url = base_url.rstrip("/")
        self.session = requests.Session()

        # Configure CA certificate
        self._configure_ca_cert(ca_cert_path)

        # Configure retry strategy
        retry_strategy = requests.adapters.Retry(
            total=3,
            backoff_factor=1,
            status_forcelist=[429, 500, 502, 503, 504],
        )
        adapter = requests.adapters.HTTPAdapter(max_retries=retry_strategy)
        self.session.mount("http://", adapter)
        self.session.mount("https://", adapter)

    def _configure_ca_cert(self, ca_cert_path: Optional[str] = None) -> None:
        """Configure CA certificate for SSL verification.

        Priority order:
        1. Explicitly provided ca_cert_path parameter
        2. Environment variable MODELREGISTRY_CA_CERT_PATH
        3. Kubernetes default CA cert (if file exists)
        4. System default CA bundle (requests default)

        Args:
            ca_cert_path: Path to CA certificate file
        """
        # Only configure CA for HTTPS URLs
        if not self.base_url.startswith("https://"):
            logger.debug("Using HTTP connection, skipping CA certificate configuration")
            return

        final_ca_path = ca_cert_path

        # Check environment variable if no explicit path provided
        if not final_ca_path:
            env_ca_path = os.getenv(DEFAULT_CA_CERT_ENV_VAR)
            if env_ca_path:
                logger.info(
                    f"Using CA certificate from environment variable {DEFAULT_CA_CERT_ENV_VAR}: {env_ca_path}"
                )
                final_ca_path = env_ca_path

        # Check for Kubernetes default CA if still no path
        if not final_ca_path and Path(DEFAULT_K8S_CA_CERT_PATH).exists():
            logger.info(
                f"Using Kubernetes default CA certificate: {DEFAULT_K8S_CA_CERT_PATH}"
            )
            final_ca_path = DEFAULT_K8S_CA_CERT_PATH

        # Configure the session with CA certificate
        if final_ca_path:
            if Path(final_ca_path).exists():
                logger.info(f"Configuring custom CA certificate: {final_ca_path}")
                self.session.verify = final_ca_path
            else:
                logger.warning(
                    f"CA certificate file not found: {final_ca_path}, falling back to system CA"
                )
                # Let requests use system CA bundle (default behavior)
        else:
            logger.debug("Using system default CA bundle")
            # Let requests use system CA bundle (default behavior)

    def _get_request_headers(self) -> Dict[str, str]:
        """Get request headers."""

        from mlflow.exceptions import MlflowException

        headers = {"Content-Type": "application/json", "Accept": "application/json"}

        if self.base_url.startswith("https://"):
            try:
                get_auth_headers(headers)
            except RuntimeError as e:
                logger.error(f"Error getting authentication headers: {e}")
                raise MlflowException(
                    f"Error getting authentication headers: {e}"
                ) from e
        return headers

    def request(self, method: str, endpoint: str, **kwargs) -> Dict[str, Any]:
        """Make authenticated request to Model Registry API.

        Args:
            method: HTTP method (GET, POST, PATCH, DELETE)
            endpoint: API endpoint path
            **kwargs: Additional request parameters

        Returns:
            Response data as dictionary

        Raises:
            MlflowException: If the request fails
        """
        from mlflow.exceptions import MlflowException, get_error_code

        url = f"{self.base_url}/{endpoint.lstrip('/')}"
        headers = self._get_request_headers()
        headers.update(kwargs.pop("headers", {}))

        # Convert customProperties to ModelRegistry format for outgoing requests
        json_data = kwargs.get("json")
        if json_data is not None:
            toModelRegistryCustomProperties(json_data)

        try:
            response = self.session.request(method, url, headers=headers, **kwargs)
            response.raise_for_status()
            response_json = response.json()
        except requests.exceptions.SSLError as e:
            # Handle TLS certificate errors specifically
            msg = (
                f"TLS certificate verification failed connecting to Model Registry: {e}"
            )
            logger.error(msg)
            raise MlflowException(msg) from e
        except requests.exceptions.RequestException as e:
            # Handle HTTP errors
            if hasattr(e, "response") and e.response is not None:
                try:
                    error_detail = e.response.json().get("message", e.response.text)
                except (ValueError, KeyError):
                    error_detail = e.response.text
                msg = f"Model Registry API error: {error_detail}"
                raise MlflowException(
                    msg,
                    error_code=get_error_code(e.response.status_code),
                ) from e
            else:
                # Handle network errors
                msg = f"Network error connecting to Model Registry: {e}"
                raise MlflowException(msg) from e

        # Convert ModelRegistry customProperties format back to MLflow format
        if response_json.get("items"):
            for item in response_json.get("items"):
                fromModelRegistryCustomProperties(item)
        else:
            fromModelRegistryCustomProperties(response_json)

        return response_json

    def get(self, endpoint: str, **kwargs) -> Dict[str, Any]:
        """Make GET request."""
        return self.request("GET", endpoint, **kwargs)

    def post(self, endpoint: str, **kwargs) -> Dict[str, Any]:
        """Make POST request."""
        return self.request("POST", endpoint, **kwargs)

    def patch(self, endpoint: str, **kwargs) -> Dict[str, Any]:
        """Make PATCH request."""
        return self.request("PATCH", endpoint, **kwargs)

    def delete(self, endpoint: str, **kwargs) -> Dict[str, Any]:
        """Make DELETE request."""
        return self.request("DELETE", endpoint, **kwargs)
