"""Model Registry REST API.

REST API for Model Registry to create and manage ML model metadata

The version of the OpenAPI document: v1alpha3
Generated by OpenAPI Generator (https://openapi-generator.tech)

Do not edit the class manually.
"""  # noqa: E501

from __future__ import annotations

import json
from enum import Enum

from typing_extensions import Self


class ModelVersionState(str, Enum):
    """- LIVE: A state indicating that the `ModelVersion` exists - ARCHIVED: A state indicating that the `ModelVersion` has been archived."""

    """
    allowed enum values
    """
    LIVE = "LIVE"
    ARCHIVED = "ARCHIVED"

    @classmethod
    def from_json(cls, json_str: str) -> Self:
        """Create an instance of ModelVersionState from a JSON string."""
        return cls(json.loads(json_str))