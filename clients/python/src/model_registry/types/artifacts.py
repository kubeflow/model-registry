"""Artifact types for model registry.

Artifacts represent pieces of data.
This could be datasets, models, metrics, or any other piece of data produced or consumed by an
execution, such as an experiment run.

Those types are used to map between proto types based on artifacts and Python objects.

Todo:
    * Move part of the description to API Reference docs (#120).
"""

from __future__ import annotations

from abc import ABC
from typing import TypeVar

from .base import BaseResourceModel

A = TypeVar("A", bound="Artifact")


class Artifact(BaseResourceModel, ABC):
    """Base class for all artifacts.

    Attributes:
        name: Name of the artifact.
        uri: URI of the artifact.
        state: State of the artifact.
    """

    uri: str
    # state: ArtifactState = ArtifactState.UNKNOWN


class DocArtifact(Artifact):
    """Represents a Document Artifact.

    Attributes:
        name: Name of the document.
        uri: URI of the document.
        description: Description of the object.
        external_id: Customizable ID. Has to be unique among instances of the same type.
    """


class ModelArtifact(Artifact):
    """Represents a Model.

    Attributes:
        name: Name of the model.
        uri: URI of the model.
        description: Description of the object.
        external_id: Customizable ID. Has to be unique among instances of the same type.
        model_format_name: Name of the model format.
        model_format_version: Version of the model format.
        storage_key: Storage secret name.
        storage_path: Storage path of the model.
        service_account_name: Name of the service account with storage secret.
    """

    # TODO: this could be an enum of valid formats
    model_format_name: str | None = None
    model_format_version: str | None = None
    storage_key: str | None = None
    storage_path: str | None = None
    service_account_name: str | None = None
