"""Base types for model registry."""

from __future__ import annotations

from abc import ABC
from typing import Union

from pydantic import BaseModel

SupportedTypes = Union[bool, int, float, str]


class BaseResourceModel(BaseModel, ABC):
    """Abstract base type for protos.

    This is a type defining common functionality for all types representing Model Registry resources,
    such as Artifacts, Contexts, and Executions.

    Attributes:
        id: Object ID. Auto-assigned when put on the server.
        name: Name of the object.
        description: Description of the object.
        external_id: Customizable ID. Has to be unique among instances of the same type.
        create_time_since_epoch: Seconds elapsed since object creation time, measured against epoch.
        last_update_time_since_epoch: Seconds elapsed since object last update time, measured against epoch.
    """

    name: str
    id: str | None = None
    description: str | None = None
    external_id: str | None = None
    create_time_since_epoch: str | None = None
    last_update_time_since_epoch: str | None = None
    custom_properties: dict[str, SupportedTypes] | None = None
