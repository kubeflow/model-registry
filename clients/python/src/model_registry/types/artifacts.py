"""Artifact types for model registry.

Artifacts represent pieces of data.
This could be datasets, models, metrics, or any other piece of data produced or consumed by an
execution, such as an experiment run.

Those types are used to map between proto types based on artifacts and Python objects.

Todo:
    * Move part of the description to API Reference docs (#120).
"""

from __future__ import annotations

from abc import ABC, abstractmethod
from typing import Any, TypeVar

from typing_extensions import override

from mr_openapi import (
    Artifact as ArtifactBaseModel,
)
from mr_openapi import (
    ArtifactState,
    ModelArtifactCreate,
    ModelArtifactUpdate,
)
from mr_openapi import (
    DocArtifact as DocArtifactBaseModel,
)
from mr_openapi import (
    ModelArtifact as ModelArtifactBaseModel,
)

from .base import BaseResourceModel

A = TypeVar("A", bound="Artifact")


class Artifact(BaseResourceModel, ABC):
    """Base class for all artifacts.

    Attributes:
        name: Name of the artifact.
        uri: URI of the artifact.
        state: State of the artifact.
    """

    name: str | None = None
    uri: str
    state: ArtifactState = ArtifactState.UNKNOWN

    @classmethod
    def from_artifact(cls: type[A], source: ArtifactBaseModel) -> A:
        """Convert a base artifact."""
        model = source.actual_instance
        assert model
        return cls.from_basemodel(model)

    @staticmethod
    def validate_artifact(source: ArtifactBaseModel) -> DocArtifact | ModelArtifact:
        """Validate an artifact."""
        model = source.actual_instance
        assert model
        if isinstance(model, DocArtifactBaseModel):
            return DocArtifact.from_basemodel(model)
        return ModelArtifact.from_basemodel(model)

    @abstractmethod
    def as_basemodel(self) -> Any:
        """Wrap the object in a BaseModel object."""

    def wrap(self) -> ArtifactBaseModel:
        """Wrap the object in a ArtifactBaseModel object."""
        return ArtifactBaseModel(self.as_basemodel())


class DocArtifact(Artifact):
    """Represents a Document Artifact.

    Attributes:
        name: Name of the document.
        uri: URI of the document.
        description: Description of the object.
        external_id: Customizable ID. Has to be unique among instances of the same type.
    """

    @override
    def create(self, **kwargs) -> Any:
        raise NotImplementedError

    @override
    def update(self, **kwargs) -> Any:
        raise NotImplementedError

    @override
    def as_basemodel(self) -> DocArtifactBaseModel:
        return DocArtifactBaseModel(
            customProperties=self._map_custom_properties(),
            **self._props_as_dict(exclude=("custom_properties")),
            artifactType="doc-artifact",
        )

    @classmethod
    @override
    def from_basemodel(cls, source: DocArtifactBaseModel) -> DocArtifact:
        assert source.name
        assert source.uri
        assert source.state
        return cls(
            id=source.id,
            name=source.name,
            description=source.description,
            external_id=source.external_id,
            create_time_since_epoch=source.create_time_since_epoch,
            last_update_time_since_epoch=source.last_update_time_since_epoch,
            uri=source.uri,
            state=source.state,
            custom_properties=cls._unmap_custom_properties(source.custom_properties)
            if source.custom_properties
            else None,
        )


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

    _model_version_id: str | None = None

    @override
    def create(self, **kwargs) -> ModelArtifactCreate:
        """Create a new ModelArtifactCreate object."""
        return ModelArtifactCreate(
            customProperties=self._map_custom_properties(),
            **self._props_as_dict(exclude=("id", "custom_properties")),
            **kwargs,
        )

    @override
    def update(self, **kwargs) -> ModelArtifactUpdate:
        """Create a new ModelArtifactUpdate object."""
        return ModelArtifactUpdate(
            customProperties=self._map_custom_properties(),
            **self._props_as_dict(exclude=("id", "name", "custom_properties")),
            **kwargs,
        )

    @override
    def as_basemodel(self) -> ModelArtifactBaseModel:
        return ModelArtifactBaseModel(
            customProperties=self._map_custom_properties(),
            **self._props_as_dict(exclude=("custom_properties")),
            artifactType="model-artifact",
        )

    @classmethod
    @override
    def from_basemodel(cls, source: ModelArtifactBaseModel) -> ModelArtifact:
        """Create a new ModelArtifact object from a BaseModel object."""
        assert source.name
        assert source.uri
        assert source.state
        return cls(
            id=source.id,
            name=source.name,
            description=source.description,
            external_id=source.external_id,
            create_time_since_epoch=source.create_time_since_epoch,
            last_update_time_since_epoch=source.last_update_time_since_epoch,
            uri=source.uri,
            model_format_name=source.model_format_name,
            model_format_version=source.model_format_version,
            storage_key=source.storage_key,
            storage_path=source.storage_path,
            service_account_name=source.service_account_name,
            state=source.state,
            custom_properties=cls._unmap_custom_properties(source.custom_properties)
            if source.custom_properties
            else None,
        )
