"""
Utility functions for Model Registry plugin
"""

import os
from typing import Optional, Tuple, Dict, Any
from urllib.parse import urlparse
from enum import Enum
from mlflow.entities import LoggedModelStatus
from mlflow.entities import LifecycleStage


class ModelIOType(Enum):
    INPUT = "input"
    OUTPUT = "output"
    UNKNOWN = "unknown"


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

            dt = datetime.fromisoformat(timestamp_str.replace("Z", "+00:00"))
            return int(dt.timestamp() * 1000)
    except (ValueError, TypeError, AttributeError, OSError):
        return None


def convert_modelregistry_state(response: Dict):
    """
    Convert Model Registry state to MLflow lifecycle stage.

    Args:
        response: Model Registry Response with state (LIVE, ARCHIVED, etc.)

    Returns:
        Corresponding MLflow lifecycle stage
    """
    return (
        LifecycleStage.ACTIVE
        if response.get("state") != "ARCHIVED"
        else LifecycleStage.DELETED
    )


def convert_to_model_artifact_state(state: Optional[LoggedModelStatus]) -> str:
    """
    Convert MLflow LoggedModelStatus to Model Artifact State.
    """
    if state is None:
        return "UNKNOWN"

    # map to values from openapi/model_artifact_state.go
    if state == LoggedModelStatus.UNSPECIFIED:
        return "UNKNOWN"
    elif state == LoggedModelStatus.PENDING:
        return "PENDING"
    elif state == LoggedModelStatus.READY:
        return "LIVE"
    elif state == LoggedModelStatus.FAILED:
        return "ABANDONED"


def convert_to_mlflow_logged_model_status(state: Optional[str]) -> LoggedModelStatus:
    """
    Convert Model Artifact State to MLflow LoggedModelStatus.
    """
    if state is None:
        return LoggedModelStatus.UNSPECIFIED

    # map to values from openapi/model_artifact_state.go
    if state == "PENDING":
        return LoggedModelStatus.PENDING
    elif state == "LIVE":
        return LoggedModelStatus.READY
    elif state == "ABANDONED":
        return LoggedModelStatus.FAILED
    else:
        return LoggedModelStatus.UNSPECIFIED


def toModelRegistryCustomProperties(json):
    if json.get("customProperties"):
        # wrap each custom property with MetadataStringVal
        custom_props = {}
        for k, v in json["customProperties"].items():
            if v is not None:
                # TODO: handle other types of custom properties
                custom_props[k] = {
                    "string_value": str(v),
                    "metadataType": "MetadataStringValue",
                }
        json["customProperties"] = custom_props


def fromModelRegistryCustomProperties(response_json):
    if response_json.get("customProperties"):
        custom_props = {}
        for k, v in response_json["customProperties"].items():
            if v["metadataType"] == "MetadataStringValue":
                custom_props[k] = v["string_value"]
            elif v["metadataType"] == "MetadataIntValue":
                custom_props[k] = int(v["int_value"])
            elif v["metadataType"] == "MetadataFloatValue":
                custom_props[k] = float(v["float_value"])
            elif v["metadataType"] == "MetadataBoolValue":
                custom_props[k] = bool(v["bool_value"])
            else:
                custom_props[k] = v["string_value"]
        response_json["customProperties"] = custom_props
