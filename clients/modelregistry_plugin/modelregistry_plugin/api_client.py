"""Model Registry API Client for handling HTTP communication."""

from __future__ import annotations

from typing import Any, Dict

import requests

from modelregistry_plugin.auth import get_auth_headers
from modelregistry_plugin.utils import (
    fromModelRegistryCustomProperties,
    toModelRegistryCustomProperties,
)


class ModelRegistryAPIClient:
    """Handles all HTTP communication with Model Registry."""

    def __init__(self, base_url: str):
        """Initialize the API client.

        Args:
            base_url: Base URL for the Model Registry API
        """
        self.base_url = base_url.rstrip("/")
        self.session = requests.Session()

        # Configure retry strategy
        retry_strategy = requests.adapters.Retry(
            total=3,
            backoff_factor=1,
            status_forcelist=[429, 500, 502, 503, 504],
        )
        adapter = requests.adapters.HTTPAdapter(max_retries=retry_strategy)
        self.session.mount("http://", adapter)
        self.session.mount("https://", adapter)

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
        headers = get_auth_headers()
        headers.update(kwargs.pop("headers", {}))

        # Convert customProperties to ModelRegistry format for outgoing requests
        json_data = kwargs.get("json")
        if json_data is not None:
            toModelRegistryCustomProperties(json_data)

        try:
            response = self.session.request(method, url, headers=headers, **kwargs)
            response.raise_for_status()
            response_json = response.json()
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
