from __future__ import annotations

"""Authentication utilities for Model Registry API."""

import os


def get_auth_headers() -> dict[str, str]:
    """Get authentication headers for Model Registry API requests.

    Returns:
        Dictionary of headers for authentication
    """
    headers = {"Content-Type": "application/json", "Accept": "application/json"}

    # Token-based authentication
    token = os.getenv("MODEL_REGISTRY_TOKEN")
    if token:
        headers["Authorization"] = f"Bearer {token}"
        return headers

    # Kubernetes service account token
    token_path = "/var/run/secrets/kubernetes.io/serviceaccount/token"
    if os.path.exists(token_path):
        with open(token_path) as f:
            token = f.read().strip()
        headers["Authorization"] = f"Bearer {token}"
        return headers

    return headers
