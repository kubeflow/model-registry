"""Tests for user facing model registry APIs."""

from collections import namedtuple

from ml_metadata.proto import (
    ArtifactType,
    Artifact,
    ContextType,
    Context,
    metadata_store_pb2,
)
from model_registry import ModelRegistry
from model_registry.store import MLMDStore
from model_registry.types import ModelArtifact, ModelVersion, RegisteredModel
from pytest import fixture


@fixture
def model_registry(store_wrapper: MLMDStore) -> ModelRegistry:
    mr = object.__new__(ModelRegistry)
    mr._store = store_wrapper
    return mr


Mapped = namedtuple("Mapped", ["proto", "py"])


@fixture
def model(store_wrapper: MLMDStore) -> Mapped:
    art_type = ArtifactType()
    art_type.name = ModelArtifact.get_proto_type_name()
    props = [
        "description",
        "modelFormatName",
        "modelFormatVersion",
        "runtime",
        "storageKey",
        "storagePath",
        "serviceAccountName",
    ]
    for key in props:
        art_type.properties[key] = metadata_store_pb2.STRING

    art = Artifact()
    art.name = "model"
    art.type_id = store_wrapper._mlmd_store.put_artifact_type(art_type)
    art.uri = "uri"

    return Mapped(art, ModelArtifact("model", "uri"))


@fixture
def model_version(store_wrapper: MLMDStore, model: Mapped) -> Mapped:
    ctx_type = ContextType()
    ctx_type.name = ModelVersion.get_proto_type_name()
    props = [
        "author",
        "description",
        "model_name",
        "tags",
    ]
    for key in props:
        ctx_type.properties[key] = metadata_store_pb2.STRING

    ctx = Context()
    ctx.name = "version"
    ctx.type_id = store_wrapper._mlmd_store.put_context_type(ctx_type)
    ctx.properties["author"].string_value = "author"
    ctx.properties["model_name"].string_value = model.py.name

    return Mapped(ctx, ModelVersion(model.py, "version", "author"))


@fixture
def registered_model(store_wrapper: MLMDStore, model: Mapped) -> Mapped:
    ctx_type = ContextType()
    ctx_type.name = RegisteredModel.get_proto_type_name()
    ctx_type.properties["description"] = metadata_store_pb2.STRING

    ctx = Context()
    ctx.name = model.py.name
    ctx.type_id = store_wrapper._mlmd_store.put_context_type(ctx_type)

    return Mapped(ctx, RegisteredModel(model.py.name))


# TODO: should we test insert/update separately?
def test_upsert_registered_model(
    model_registry: ModelRegistry, registered_model: Mapped
):
    model_registry.upsert_registered_model(registered_model.py)

    rm_proto = model_registry._store._mlmd_store.get_context_by_type_and_name(
        RegisteredModel.get_proto_type_name(), registered_model.proto.name
    )
    assert rm_proto is not None
    assert registered_model.py.id == str(rm_proto.id)
    assert registered_model.py.name == rm_proto.name


def test_get_registered_model_by_id(
    model_registry: ModelRegistry, registered_model: Mapped
):
    id = model_registry._store._mlmd_store.put_contexts([registered_model.proto])[0]
    id = str(id)

    mlmd_rm = model_registry.get_registered_model_by_id(id)
    assert mlmd_rm.id == id
    assert mlmd_rm.name == registered_model.proto.name


def test_upsert_model_version(
    model_registry: ModelRegistry, model_version: Mapped, registered_model: Mapped
):
    rm_id = model_registry._store._mlmd_store.put_contexts([registered_model.proto])[0]
    rm_id = str(rm_id)

    model_registry.upsert_model_version(model_version.py, rm_id)

    mv_proto = model_registry._store._mlmd_store.get_context_by_type_and_name(
        ModelVersion.get_proto_type_name(), model_version.proto.name
    )
    assert mv_proto is not None
    assert model_version.py.id == str(mv_proto.id)
    assert model_version.py.version == mv_proto.name


def test_get_model_version_by_id(model_registry: ModelRegistry, model_version: Mapped):
    id = model_registry._store._mlmd_store.put_contexts([model_version.proto])[0]
    id = str(id)

    mlmd_mv = model_registry.get_model_version_by_id(id)
    assert mlmd_mv.id == id
    assert mlmd_mv.version == model_version.proto.name


def test_upsert_model_artifact(
    model_registry: ModelRegistry, model: Mapped, model_version: Mapped
):
    mv_id = model_registry._store._mlmd_store.put_contexts([model_version.proto])[0]
    mv_id = str(mv_id)

    model_registry.upsert_model_artifact(model.py, mv_id)

    ma_proto = model_registry._store._mlmd_store.get_artifact_by_type_and_name(
        ModelArtifact.get_proto_type_name(), model.proto.name
    )
    assert ma_proto is not None
    assert model.py.id == str(ma_proto.id)
    assert model.py.name == ma_proto.name


def test_get_model_artifact_by_id(model_registry: ModelRegistry, model: Mapped):
    id = model_registry._store._mlmd_store.put_artifacts([model.proto])[0]
    id = str(id)

    mlmd_ma = model_registry.get_model_artifact_by_id(id)

    assert mlmd_ma.id == id
    assert mlmd_ma.name == model.proto.name
