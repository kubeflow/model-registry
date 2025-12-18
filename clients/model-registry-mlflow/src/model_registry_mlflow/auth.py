"""Authentication utilities for Model Registry API."""

from __future__ import annotations
from typing import Final, Optional
import os
import threading


# Kubernetes service account token path
_TOKEN_PATH: Final[str] = "/var/run/secrets/kubernetes.io/serviceaccount/token"


class K8sTokenCache:
    """Thread-safe caching for Kubernetes service account tokens."""

    _token: Optional[str] = None
    _token_mtime: Optional[float] = None
    _lock = threading.Lock()

    @classmethod
    def get_token(cls) -> Optional[str]:
        """Get the Kubernetes token with caching."""
        with cls._lock:
            try:
                mtime = os.path.getmtime(_TOKEN_PATH)
                if cls._token_mtime == mtime and cls._token is not None:
                    return cls._token

                # Need to refresh token - read file with mtime consistency check
                cls._token_mtime = mtime
                while True:
                    try:
                        with open(_TOKEN_PATH, "r") as f:
                            cls._token = f.read().strip()
                        # Check if file was modified during read
                        if os.path.getmtime(_TOKEN_PATH) == cls._token_mtime:
                            return cls._token
                        # File was modified during read, try again
                        cls._token_mtime = os.path.getmtime(_TOKEN_PATH)
                    except OSError:
                        # File temporarily unavailable (e.g., during atomic replacement)
                        # Update mtime and retry
                        try:
                            cls._token_mtime = os.path.getmtime(_TOKEN_PATH)
                        except OSError:
                            # File doesn't exist, let outer exception handler deal with it
                            raise

            except (OSError, UnicodeDecodeError) as e:
                # Clear cache on any error to prevent stale data
                cls._token = None
                cls._token_mtime = None
                raise RuntimeError(
                    f"Error accessing Kubernetes token file {_TOKEN_PATH}: {e}"
                ) from e


def _get_k8s_token() -> Optional[str]:
    """Get Kubernetes service account token with caching."""
    return K8sTokenCache.get_token()


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
