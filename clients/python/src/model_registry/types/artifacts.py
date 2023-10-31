"""Artifact types for model registry.

Artifacts represent pieces of data.
This could be datasets, models, metrics, or any other piece of data produced or consumed by an
execution, such as an experiment run.

Those types are used to map between proto types based on artifacts and Python objects.

TODO:
    * Move part of the description to API Reference docs (#120).
"""

from enum import Enum, unique
from typing import Optional

from attrs import define, field

from .base import ProtoBase


@unique
class ArtifactState(Enum):
    """State of an artifact."""

    UNKNOWN = 0
    PENDING = 1
    LIVE = 2
    MARKED_FOR_DELETION = 3
    DELETED = 4
    ABANDONED = 5
    REFERENCE = 6


@define(slots=False)
class BaseArtifact(ProtoBase):
    """Abstract base class for all artifacts.

    Attributes:
        name (str): Name of the artifact.
        uri (str): URI of the artifact.
        state (ArtifactState): State of the artifact.
    """

    name: str
    uri: str
    state: ArtifactState = field(init=False, default=ArtifactState.UNKNOWN)


@define(slots=False)
class ModelArtifact(BaseArtifact):
    """Represents a Model.

    Attributes:
        name (str): Name of the model.
        uri (str): URI of the model.
        description (str, optional): Description of the object.
        external_id (str, optional): Customizable ID. Has to be unique among instances of the same type.
        model_format_name (str, optional): Name of the model format.
        model_format_version (str, optional): Version of the model format.
        runtime (str, optional): Runtime of the model.
        storage_key (str, optional): Storage key of the model.
        storage_path (str, optional): Storage path of the model.
        service_account_name (str, optional): Service account name of the model.
    """

    # TODO: this could be an enum of valid formats
    model_format_name: Optional[str] = field(kw_only=True, default=None)
    model_format_version: Optional[str] = field(kw_only=True, default=None)
    runtime: Optional[str] = field(kw_only=True, default=None)
    storage_key: Optional[str] = field(kw_only=True, default=None)
    storage_path: Optional[str] = field(kw_only=True, default=None)
    service_account_name: Optional[str] = field(kw_only=True, default=None)
