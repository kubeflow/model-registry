"""
Utility functions for Model Registry plugin
"""

import os
from typing import Tuple, Dict, Any
from urllib.parse import urlparse

from mlflow.entities import LifecycleStage


def parse_tracking_uri(tracking_uri: str) -> Tuple[str, int, bool]:
    """
    Parse the tracking URI to extract connection details.
    
    Args:
        tracking_uri: URI like "modelregistry://localhost:8080" or "modelregistry+https://host:port"
    
    Returns:
        Tuple of (host, port, secure)
    """
    parsed = urlparse(tracking_uri)
    
    # Handle different URI schemes
    if parsed.scheme == "modelregistry":
        secure = os.getenv("MODEL_REGISTRY_SECURE", "false").lower() == "true"
    elif parsed.scheme == "modelregistry+https":
        secure = True
    elif parsed.scheme == "modelregistry+http":
        secure = False
    else:
        secure = False
    
    host = parsed.hostname or os.getenv("MODEL_REGISTRY_HOST", "localhost")
    port = parsed.port or int(os.getenv("MODEL_REGISTRY_PORT", "8080"))
    
    return host, port, secure


def convert_timestamp(timestamp_str: str) -> Any:
    """
    Convert timestamp string to milliseconds since epoch.
    
    Args:
        timestamp_str: Timestamp string from Model Registry
        
    Returns:
        Timestamp in milliseconds since epoch
    """
    if not timestamp_str:
        return None
        
    try:
        # Assuming timestamp is in ISO format or milliseconds
        if timestamp_str.isdigit():
            return int(timestamp_str)
        else:
            # Parse ISO format timestamp
            from datetime import datetime
            dt = datetime.fromisoformat(timestamp_str.replace('Z', '+00:00'))
            return int(dt.timestamp() * 1000)
    except:
        return None


def convert_modelregistry_state(response: Dict):
    """
    Convert Model Registry state to MLflow lifecycle stage.

    Args:
        response: Model Registry Response with state (LIVE, ARCHIVED, etc.)

    Returns:
        Corresponding MLflow lifecycle stage
    """
    return LifecycleStage.ACTIVE if response.get("state") != "ARCHIVED" else LifecycleStage.DELETED
