"""Artifact types for model registry.

Artifacts represent pieces of data.
This could be datasets, models, metrics, or any other piece of data produced or consumed by an
execution, such as an experiment run.

Those types are used to map between proto types based on artifacts and Python objects.

Todo:
    * Move part of the description to API Reference docs (#120).
"""

from __future__ import annotations

from enum import Enum, unique
from uuid import uuid4

from attrs import define, field
from ml_metadata.proto import Artifact
from typing_extensions import override

from .base import Prefixable, ProtoBase


@unique
class ArtifactState(Enum):
    """State of an artifact."""

    UNKNOWN = Artifact.UNKNOWN
    PENDING = Artifact.PENDING
    LIVE = Artifact.LIVE
    MARKED_FOR_DELETION = Artifact.MARKED_FOR_DELETION
    DELETED = Artifact.DELETED
    ABANDONED = Artifact.ABANDONED
    REFERENCE = Artifact.REFERENCE


@define(slots=False)
class BaseArtifact(ProtoBase):
    """Abstract base class for all artifacts.

    Attributes:
        name: Name of the artifact.
        uri: URI of the artifact.
        state: State of the artifact.
    """

    name: str
    uri: str
    state: ArtifactState = field(init=False, default=ArtifactState.UNKNOWN)

    @classmethod
    @override
    def get_proto_type(cls) -> type[Artifact]:
        return Artifact

    @override
    def map(self, type_id: int) -> Artifact:
        mlmd_obj = super().map(type_id)
        mlmd_obj.uri = self.uri
        mlmd_obj.state = self.state.value
        return mlmd_obj

    @classmethod
    @override
    def unmap(cls, mlmd_obj: Artifact) -> BaseArtifact:
        py_obj = super().unmap(mlmd_obj)
        assert isinstance(
            py_obj, BaseArtifact
        ), f"Expected BaseArtifact, got {type(py_obj)}"
        py_obj.uri = mlmd_obj.uri
        py_obj.state = ArtifactState(mlmd_obj.state)
        return py_obj


@define(slots=False, auto_attribs=True)
class ModelArtifact(BaseArtifact, Prefixable):
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
    model_format_name: str | None = field(kw_only=True, default=None)
    model_format_version: str | None = field(kw_only=True, default=None)
    storage_key: str | None = field(kw_only=True, default=None)
    storage_path: str | None = field(kw_only=True, default=None)
    service_account_name: str | None = field(kw_only=True, default=None)

    _model_version_id: str | None = field(init=False, default=None)

    @property
    @override
    def mlmd_name_prefix(self) -> str:
        return self._model_version_id if self._model_version_id else uuid4().hex

    @override
    def map(self, type_id: int) -> Artifact:
        mlmd_obj = super().map(type_id)
        props = {
            "model_format_name": self.model_format_name,
            "model_format_version": self.model_format_version,
            "storage_key": self.storage_key,
            "storage_path": self.storage_path,
            "service_account_name": self.service_account_name,
        }
        self._map_props(props, mlmd_obj.properties)
        return mlmd_obj

    @override
    @classmethod
    def unmap(cls, mlmd_obj: Artifact) -> ModelArtifact:
        py_obj = super().unmap(mlmd_obj)
        assert isinstance(
            py_obj, ModelArtifact
        ), f"Expected ModelArtifact, got {type(py_obj)}"
        py_obj.model_format_name = mlmd_obj.properties["model_format_name"].string_value
        py_obj.model_format_version = mlmd_obj.properties[
            "model_format_version"
        ].string_value
        py_obj.storage_key = mlmd_obj.properties["storage_key"].string_value
        py_obj.storage_path = mlmd_obj.properties["storage_path"].string_value
        py_obj.service_account_name = mlmd_obj.properties[
            "service_account_name"
        ].string_value
        return py_obj
