"""Utility functions for Model Registry plugin."""

from __future__ import annotations

import os
from enum import Enum
from typing import TYPE_CHECKING, Any, Optional

if TYPE_CHECKING:
    from mlflow.entities import LoggedModelStatus
from urllib.parse import urlparse


class ModelIOType(Enum):
    INPUT = "input"
    OUTPUT = "output"
    UNKNOWN = "unknown"


def parse_tracking_uri(tracking_uri: str) -> tuple[str, int, bool]:
    """Parse the tracking URI to extract connection details.

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


def convert_timestamp(timestamp_input: Any) -> Optional[int]:
    """Convert timestamp input to milliseconds since epoch.

    Args:
        timestamp_input: Timestamp from Model Registry (string or int)

    Returns:
        Timestamp in milliseconds since epoch, or None if conversion fails
    """
    if timestamp_input is None:
        return None

    try:
        # If it's already an integer, return it
        if isinstance(timestamp_input, int):
            return timestamp_input

        # Convert to string and handle
        timestamp_str = str(timestamp_input)
        if not timestamp_str:
            return None

        # Assuming timestamp is in ISO format or milliseconds
        if timestamp_str.isdigit():
            return int(timestamp_str)
        # Parse ISO format timestamp
        from datetime import datetime

        dt = datetime.fromisoformat(timestamp_str.replace("Z", "+00:00"))
        return int(dt.timestamp() * 1000)
    except (ValueError, TypeError, AttributeError, OSError):
        return None


def convert_modelregistry_state(response: dict):
    """Convert Model Registry state to MLflow lifecycle stage.

    Args:
        response: Model Registry Response with state (LIVE, ARCHIVED, etc.)

    Returns:
        Corresponding MLflow lifecycle stage
    """
    from mlflow.entities import LifecycleStage

    return (
        LifecycleStage.ACTIVE
        if response.get("state") != "ARCHIVED"
        else LifecycleStage.DELETED
    )


def convert_to_model_artifact_state(state: Optional[LoggedModelStatus]) -> str:
    """Convert MLflow LoggedModelStatus to Model Artifact State."""
    from mlflow.entities import LoggedModelStatus

    if state is None:
        return "UNKNOWN"

    # map to values from openapi/model_artifact_state.go
    if state == LoggedModelStatus.UNSPECIFIED:
        return "UNKNOWN"
    if state == LoggedModelStatus.PENDING:
        return "PENDING"
    if state == LoggedModelStatus.READY:
        return "LIVE"
    if state == LoggedModelStatus.FAILED:
        return "ABANDONED"
    return None


def convert_to_mlflow_logged_model_status(state: Optional[str]) -> LoggedModelStatus:
    """Convert Model Artifact State to MLflow LoggedModelStatus."""
    from mlflow.entities import LoggedModelStatus

    if state is None:
        return LoggedModelStatus.UNSPECIFIED

    # map to values from openapi/model_artifact_state.go
    if state == "PENDING":
        return LoggedModelStatus.PENDING
    if state == "LIVE":
        return LoggedModelStatus.READY
    if state == "ABANDONED":
        return LoggedModelStatus.FAILED
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
