"""Context types for model registry.

Contexts group related Artifacts together.
They provide a way to organize and categorize components in a workflow.

Those types are used to map between proto types based on contexts and Python objects.

Todo:
    * Move part of the description to API Reference docs (#120).
"""

from __future__ import annotations

from abc import ABC
from enum import Enum, unique

from attrs import define, field
from ml_metadata.proto import Context
from typing_extensions import override

from model_registry.store import ScalarType

from .artifacts import BaseArtifact
from .base import Prefixable, ProtoBase


@unique
class ContextState(Enum):
    """State of the context.

    LIVE: The context is live and can be used.
    ARCHIVED: The context is archived and can't be used.
    """

    LIVE = "LIVE"
    ARCHIVED = "ARCHIVED"


@define(slots=False, init=False)
class BaseContext(ProtoBase, ABC):
    """Abstract base class for all contexts."""

    state: ContextState = field(init=False, default=ContextState.LIVE)

    @override
    def map(self, type_id: int) -> Context:
        mlmd_obj = super().map(type_id)
        mlmd_obj.properties["state"].string_value = self.state.value
        return mlmd_obj

    @classmethod
    @override
    def unmap(cls, mlmd_obj: Context) -> BaseContext:
        py_obj = super().unmap(mlmd_obj)
        assert isinstance(
            py_obj, BaseContext
        ), f"Expected BaseContext, got {type(py_obj)}"
        py_obj.state = ContextState(mlmd_obj.properties["state"].string_value)
        return py_obj

    @classmethod
    @override
    def get_proto_type(cls) -> type[Context]:
        return Context


@define(slots=False)
class ModelVersion(BaseContext, Prefixable):
    """Represents a model version.

    Attributes:
        model_name: Name of the model associated with this version.
        version: Version of the model.
        author: Author of the model version.
        description: Description of the object.
        external_id: Customizable ID. Has to be unique among instances of the same type.
        artifacts: Artifacts associated with this version.
        tags: Tags associated with this version.
        metadata: Metadata associated with this version.
    """

    model_name: str
    version: str
    author: str
    artifacts: list[BaseArtifact] = field(init=False, factory=list)
    tags: list[str] = field(init=False, factory=list)
    metadata: dict[str, ScalarType] = field(init=False, factory=dict)

    _registered_model_id: int | None = field(init=False, default=None)

    def __attrs_post_init__(self) -> None:
        self.name = self.version

    @property
    @override
    def mlmd_name_prefix(self):
        assert (
            self._registered_model_id is not None
        ), "There's no registered model associated with this version"
        return self._registered_model_id

    @override
    def map(self, type_id: int) -> Context:
        mlmd_obj = super().map(type_id)
        # this should match the name of the registered model
        mlmd_obj.properties["model_name"].string_value = self.model_name
        mlmd_obj.properties["author"].string_value = self.author
        if self.tags:
            mlmd_obj.properties["tags"].string_value = ",".join(self.tags)
        self._map_props(self.metadata, mlmd_obj.custom_properties)
        return mlmd_obj

    @classmethod
    @override
    def unmap(cls, mlmd_obj: Context) -> ModelVersion:
        py_obj = super().unmap(mlmd_obj)
        assert isinstance(
            py_obj, ModelVersion
        ), f"Expected ModelVersion, got {type(py_obj)}"
        py_obj.version = py_obj.name
        py_obj.model_name = mlmd_obj.properties["model_name"].string_value
        py_obj.author = mlmd_obj.properties["author"].string_value
        tags = mlmd_obj.properties["tags"].string_value
        if tags:
            tags = tags.split(",")
        py_obj.tags = tags or []
        py_obj.metadata = cls._unmap_props(mlmd_obj.custom_properties)
        return py_obj


@define(slots=False)
class RegisteredModel(BaseContext):
    """Represents a registered model.

    Attributes:
        name: Registered model name.
        description: Description of the object.
        external_id: Customizable ID. Has to be unique among instances of the same type.
    """

    name: str

    @classmethod
    @override
    def unmap(cls, mlmd_obj: Context) -> RegisteredModel:
        py_obj = super().unmap(mlmd_obj)
        assert isinstance(
            py_obj, RegisteredModel
        ), f"Expected RegisteredModel, got {type(py_obj)}"
        return py_obj
