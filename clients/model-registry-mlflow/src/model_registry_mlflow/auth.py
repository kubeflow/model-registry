"""Authentication utilities for Model Registry API."""

from __future__ import annotations
from typing import Final, Optional
import os
import threading


# Kubernetes service account token path
_TOKEN_PATH: Final[str] = "/var/run/secrets/kubernetes.io/serviceaccount/token"


class K8sTokenCache:
    """Thread-safe caching for Kubernetes service account tokens."""

    def __init__(self):
        self._token: Optional[str] = None
        self._token_mtime: Optional[float] = None
        self._lock = threading.Lock()

    def get_token(self) -> Optional[str]:
        """Get the Kubernetes token with caching."""
        with self._lock:
            try:
                mtime = os.path.getmtime(_TOKEN_PATH)
                if self._token_mtime == mtime and self._token is not None:
                    return self._token

                # Need to refresh token - read file
                with open(_TOKEN_PATH, "r") as f:
                    self._token = f.read().strip()
                self._token_mtime = mtime
                return self._token

            except (OSError, UnicodeDecodeError) as e:
                # Clear cache on any error to prevent stale data
                self._token = None
                self._token_mtime = None
                raise RuntimeError(
                    f"Error accessing Kubernetes token file {_TOKEN_PATH}: {e}"
                ) from e


# Global token cache instance
_token_cache = K8sTokenCache()


def _get_k8s_token() -> Optional[str]:
    """Get Kubernetes service account token with caching."""
    global _token_cache
    return _token_cache.get_token()


def get_auth_headers(headers: dict[str, str] = None) -> None:
    """Get authentication headers for Model Registry API requests.

    Args:
        headers: Dictionary to add authorization headers to

    Raises:
        RuntimeError: If no authentication token can be obtained
    """

    # Token-based authentication
    token = os.getenv("MODEL_REGISTRY_TOKEN")
    if token:
        headers["Authorization"] = f"Bearer {token}"
        return

    # Kubernetes service account token
    k8s_token = _get_k8s_token()
    if k8s_token:
        headers["Authorization"] = f"Bearer {k8s_token}"
        return

    # No authentication token available
    raise RuntimeError(
        "No authentication token available. Set MODEL_REGISTRY_TOKEN environment variable "
        "or ensure Kubernetes service account token is available."
    )
