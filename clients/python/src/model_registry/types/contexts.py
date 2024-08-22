"""Context types for model registry.

Contexts group related Artifacts together.
They provide a way to organize and categorize components in a workflow.

Those types are used to map between proto types based on contexts and Python objects.

Todo:
    * Move part of the description to API Reference docs (#120).
"""

from __future__ import annotations

from typing_extensions import override

from mr_openapi import (
    ModelVersion as ModelVersionBaseModel,
)
from mr_openapi import (
    ModelVersionCreate,
    ModelVersionState,
    ModelVersionUpdate,
    RegisteredModelCreate,
    RegisteredModelState,
    RegisteredModelUpdate,
)
from mr_openapi import (
    RegisteredModel as RegisteredModelBaseModel,
)

from .base import BaseResourceModel


class ModelVersion(BaseResourceModel):
    """Represents a model version.

    Attributes:
        name: Name of this version.
        author: Author of the model version.
        description: Description of the object.
        external_id: Customizable ID. Has to be unique among instances of the same type.
        artifacts: Artifacts associated with this version.
    """

    name: str
    author: str | None = None
    state: ModelVersionState = ModelVersionState.LIVE

    @override
    def create(self, *, registered_model_id: str, **kwargs) -> ModelVersionCreate:  # type: ignore[override]
        return ModelVersionCreate(
            registeredModelId=registered_model_id,
            customProperties=self._map_custom_properties(),
            **self._props_as_dict(exclude=("id", "custom_properties")),
            **kwargs,
        )

    @override
    def update(self, **kwargs) -> ModelVersionUpdate:
        return ModelVersionUpdate(
            customProperties=self._map_custom_properties(),
            **self._props_as_dict(exclude=("id", "name", "custom_properties")),
            **kwargs,
        )

    @classmethod
    @override
    def from_basemodel(cls, source: ModelVersionBaseModel) -> ModelVersion:
        assert source.name
        assert source.state
        return cls(
            id=source.id,
            name=source.name,
            state=source.state,
            author=source.author,
            description=source.description,
            external_id=source.external_id,
            create_time_since_epoch=source.create_time_since_epoch,
            last_update_time_since_epoch=source.last_update_time_since_epoch,
            custom_properties=cls._unmap_custom_properties(source.custom_properties)
            if source.custom_properties
            else None,
        )


class RegisteredModel(BaseResourceModel):
    """Represents a registered model.

    Attributes:
        name: Registered model name.
        owner: Owner of this Registered Model.
        description: Description of the object.
        external_id: Customizable ID. Has to be unique among instances of the same type.
    """

    name: str
    owner: str | None = None
    state: RegisteredModelState = RegisteredModelState.LIVE

    @override
    def create(self, **kwargs) -> RegisteredModelCreate:
        return RegisteredModelCreate(
            customProperties=self._map_custom_properties(),
            **self._props_as_dict(exclude=("id", "custom_properties")),
            **kwargs,
        )

    @override
    def update(self, **kwargs) -> RegisteredModelUpdate:
        return RegisteredModelUpdate(
            customProperties=self._map_custom_properties(),
            **self._props_as_dict(exclude=("id", "name", "custom_properties")),
            **kwargs,
        )

    @classmethod
    @override
    def from_basemodel(cls, source: RegisteredModelBaseModel) -> RegisteredModel:
        assert source.name
        assert source.state
        return cls(
            id=source.id,
            name=source.name,
            owner=source.owner,
            state=source.state,
            description=source.description,
            external_id=source.external_id,
            create_time_since_epoch=source.create_time_since_epoch,
            last_update_time_since_epoch=source.last_update_time_since_epoch,
            custom_properties=cls._unmap_custom_properties(source.custom_properties)
            if source.custom_properties
            else None,
        )
