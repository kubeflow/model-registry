"""
Authentication utilities for Model Registry API
"""

import os
from typing import Dict


def get_auth_headers() -> Dict[str, str]:
    """
    Get authentication headers for Model Registry API requests.

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
        with open(token_path, "r") as f:
            token = f.read().strip()
        headers["Authorization"] = f"Bearer {token}"
        return headers

    return headers
