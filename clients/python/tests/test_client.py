"""Tests for user facing model registry APIs."""

from collections import namedtuple

from attrs import evolve
from ml_metadata.proto import (
    ArtifactType,
    Artifact,
    Attribution,
    ContextType,
    Context,
    ParentContext,
    metadata_store_pb2,
)
from model_registry import ModelRegistry
from model_registry.exceptions import StoreException
from model_registry.store import MLMDStore
from model_registry.types import ModelArtifact, ModelVersion, RegisteredModel
from pytest import fixture, raises


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
    # we can't test the name directly as it's prefixed
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
    # we can't test the name directly as it's prefixed
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
    model_registry: ModelRegistry,
    registered_model: Mapped,
    model_version: Mapped,
    model: Mapped,
):
    model.proto.name = f"test_prefix:{model.proto.name}"
    ma_id = model_registry._store._mlmd_store.put_artifacts([model.proto])[0]

    model_version.proto.name = f"1:{model_version.proto.name}"
    mv_id = model_registry._store._mlmd_store.put_contexts([model_version.proto])[0]

    model_registry._store._mlmd_store.put_attributions_and_associations(
        [Attribution(context_id=mv_id, artifact_id=ma_id)], []
    )

    rm_id = model_registry._store._mlmd_store.put_contexts([registered_model.proto])[0]

    model_registry._store._mlmd_store.put_parent_contexts(
        [ParentContext(parent_id=rm_id, child_id=mv_id)]
    )

    mlmd_rm = model_registry.get_registered_model_by_id(str(rm_id))
    assert mlmd_rm.id == str(rm_id)
    assert mlmd_rm.name == registered_model.py.name
    assert mlmd_rm.name == registered_model.proto.name


def test_get_registered_model_by_name(
    model_registry: ModelRegistry,
    registered_model: Mapped,
    model_version: Mapped,
    model: Mapped,
):
    model.proto.name = f"test_prefix:{model.proto.name}"
    ma_id = model_registry._store._mlmd_store.put_artifacts([model.proto])[0]

    model_version.proto.name = f"1:{model_version.proto.name}"
    mv_id = model_registry._store._mlmd_store.put_contexts([model_version.proto])[0]

    model_registry._store._mlmd_store.put_attributions_and_associations(
        [Attribution(context_id=mv_id, artifact_id=ma_id)], []
    )

    rm_id = model_registry._store._mlmd_store.put_contexts([registered_model.proto])[0]

    model_registry._store._mlmd_store.put_parent_contexts(
        [ParentContext(parent_id=rm_id, child_id=mv_id)]
    )

    mlmd_rm = model_registry.get_registered_model_by_params(
        name=registered_model.py.name
    )
    assert mlmd_rm.id == str(rm_id)
    assert mlmd_rm.name == registered_model.py.name
    assert mlmd_rm.name == registered_model.proto.name


def test_get_registered_model_by_external_id(
    model_registry: ModelRegistry,
    registered_model: Mapped,
    model_version: Mapped,
    model: Mapped,
):
    model.proto.name = f"test_prefix:{model.proto.name}"
    ma_id = model_registry._store._mlmd_store.put_artifacts([model.proto])[0]

    model_version.proto.name = f"1:{model_version.proto.name}"
    mv_id = model_registry._store._mlmd_store.put_contexts([model_version.proto])[0]

    model_registry._store._mlmd_store.put_attributions_and_associations(
        [Attribution(context_id=mv_id, artifact_id=ma_id)], []
    )

    registered_model.py.external_id = "external_id"
    registered_model.proto.external_id = "external_id"

    rm_id = model_registry._store._mlmd_store.put_contexts([registered_model.proto])[0]

    model_registry._store._mlmd_store.put_parent_contexts(
        [ParentContext(parent_id=rm_id, child_id=mv_id)]
    )

    mlmd_rm = model_registry.get_registered_model_by_params(
        external_id=registered_model.py.external_id
    )
    assert mlmd_rm.id == str(rm_id)
    assert mlmd_rm.name == registered_model.py.name
    assert mlmd_rm.name == registered_model.proto.name


def test_get_registered_models(model_registry: ModelRegistry, registered_model: Mapped):
    rm1_id = model_registry._store._mlmd_store.put_contexts([registered_model.proto])[0]
    registered_model.proto.name = "model2"
    rm2_id = model_registry._store._mlmd_store.put_contexts([registered_model.proto])[0]

    mlmd_rms = model_registry.get_registered_models()
    assert len(mlmd_rms) == 2
    assert mlmd_rms[0].id in [str(rm1_id), str(rm2_id)]


def test_upsert_model_version(
    model_registry: ModelRegistry, model_version: Mapped, registered_model: Mapped
):
    rm_id = model_registry._store._mlmd_store.put_contexts([registered_model.proto])[0]
    rm_id = str(rm_id)

    model_registry.upsert_model_version(model_version.py, rm_id)

    mv_proto = model_registry._store._mlmd_store.get_context_by_type_and_name(
        ModelVersion.get_proto_type_name(), f"{rm_id}:{model_version.proto.name}"
    )
    assert mv_proto is not None
    assert model_version.py.id == str(mv_proto.id)
    assert model_version.py.version != mv_proto.name


def test_get_model_version_by_id(
    model_registry: ModelRegistry, model: Mapped, model_version: Mapped
):
    model.proto.name = f"test_prefix:{model.proto.name}"
    art_id = model_registry._store._mlmd_store.put_artifacts([model.proto])[0]

    model_version.proto.name = f"1:{model_version.proto.name}"
    ctx_id = model_registry._store._mlmd_store.put_contexts([model_version.proto])[0]

    model_registry._store._mlmd_store.put_attributions_and_associations(
        [Attribution(context_id=ctx_id, artifact_id=art_id)], []
    )

    id = str(ctx_id)
    mlmd_mv = model_registry.get_model_version_by_id(id)
    assert mlmd_mv.id == id
    assert mlmd_mv.name == model_version.py.name
    assert mlmd_mv.version != model_version.proto.name


def test_get_model_version_by_name(
    model_registry: ModelRegistry, model_version: Mapped
):
    model_version.proto.name = f"1:{model_version.proto.name}"

    id = model_registry._store._mlmd_store.put_contexts([model_version.proto])[0]
    id = str(id)

    mlmd_mv = model_registry.get_model_version_by_params(
        registered_model_id="1", version=model_version.py.name
    )
    assert mlmd_mv.id == id
    assert mlmd_mv.name == model_version.py.name
    assert mlmd_mv.name != model_version.proto.name


def test_get_model_version_by_external_id(
    model_registry: ModelRegistry, model_version: Mapped
):
    model_version.proto.name = f"1:{model_version.proto.name}"
    model_version.proto.external_id = "external_id"
    model_version.py.external_id = "external_id"

    id = model_registry._store._mlmd_store.put_contexts([model_version.proto])[0]
    id = str(id)

    mlmd_mv = model_registry.get_model_version_by_params(
        external_id=model_version.py.external_id
    )
    assert mlmd_mv.id == id
    assert mlmd_mv.name == model_version.py.name
    assert mlmd_mv.name != model_version.proto.name


def test_get_model_versions(
    model_registry: ModelRegistry,
    model_version: Mapped,
    registered_model: Mapped,
    model: Mapped,
):
    model.proto.name = f"test_prefix:{model.proto.name}"
    ma_id = model_registry._store._mlmd_store.put_artifacts([model.proto])[0]

    rm_id = model_registry._store._mlmd_store.put_contexts([registered_model.proto])[0]

    model_version.proto.name = f"{rm_id}:version"
    mv1_id = model_registry._store._mlmd_store.put_contexts([model_version.proto])[0]
    model_version.proto.name = f"{rm_id}:version2"
    mv2_id = model_registry._store._mlmd_store.put_contexts([model_version.proto])[0]

    model_registry._store._mlmd_store.put_attributions_and_associations(
        [
            Attribution(context_id=mv1_id, artifact_id=ma_id),
            Attribution(context_id=mv2_id, artifact_id=ma_id),
        ],
        [],
    )

    model_registry._store._mlmd_store.put_parent_contexts(
        [
            ParentContext(parent_id=rm_id, child_id=mv1_id),
            ParentContext(parent_id=rm_id, child_id=mv2_id),
        ]
    )

    mlmd_mvs = model_registry.get_model_versions(str(rm_id))
    assert len(mlmd_mvs) == 2
    assert mlmd_mvs[0].id in [str(mv1_id), str(mv2_id)]


def test_upsert_model_artifact(
    monkeypatch, model_registry: ModelRegistry, model: Mapped, model_version: Mapped
):
    monkeypatch.setattr(ModelArtifact, "mlmd_name_prefix", "test_prefix")

    mv_id = model_registry._store._mlmd_store.put_contexts([model_version.proto])[0]
    mv_id = str(mv_id)

    model_registry.upsert_model_artifact(model.py, mv_id)

    ma_proto = model_registry._store._mlmd_store.get_artifact_by_type_and_name(
        ModelArtifact.get_proto_type_name(), f"test_prefix:{model.proto.name}"
    )
    assert ma_proto is not None
    assert model.py.id == str(ma_proto.id)
    assert model.py.name != ma_proto.name


def test_upsert_duplicate_model_artifact_with_different_version(
    model_registry: ModelRegistry, model: Mapped, model_version: Mapped
):
    mv1_id = model_registry._store._mlmd_store.put_contexts([model_version.proto])[0]
    mv1_id = str(mv1_id)

    model_version.proto.name = "version2"
    mv2_id = model_registry._store._mlmd_store.put_contexts([model_version.proto])[0]
    mv2_id = str(mv2_id)

    ma1 = evolve(model.py)
    model_registry.upsert_model_artifact(ma1, mv1_id)
    ma2 = evolve(model.py)
    model_registry.upsert_model_artifact(ma2, mv2_id)

    ma_protos = model_registry._store._mlmd_store.get_artifacts_by_id(
        [int(ma1.id), int(ma2.id)]
    )
    assert ma1.name == ma2.name
    assert ma1.name != str(ma_protos[0].name)
    assert ma2.name != str(ma_protos[1].name)


def test_upsert_duplicate_model_artifact_with_same_version(
    model_registry: ModelRegistry, model: Mapped, model_version: Mapped
):
    mv_id = model_registry._store._mlmd_store.put_contexts([model_version.proto])[0]
    mv_id = str(mv_id)

    ma1 = evolve(model.py)
    model_registry.upsert_model_artifact(ma1, mv_id)
    ma2 = evolve(model.py)
    with raises(StoreException):
        model_registry.upsert_model_artifact(ma2, mv_id)


def test_get_model_artifact_by_id(model_registry: ModelRegistry, model: Mapped):
    model.proto.name = f"test_prefix:{model.proto.name}"
    id = model_registry._store._mlmd_store.put_artifacts([model.proto])[0]
    id = str(id)

    mlmd_ma = model_registry.get_model_artifact_by_id(id)

    assert mlmd_ma.id == id
    assert mlmd_ma.name == model.py.name
    assert mlmd_ma.name != model.proto.name


def test_get_model_artifact_by_model_version_id(
    model_registry: ModelRegistry, model: Mapped, model_version: Mapped
):
    mv_id = model_registry._store._mlmd_store.put_contexts([model_version.proto])[0]

    model.proto.name = f"test_prefix:{model.proto.name}"
    ma_id = model_registry._store._mlmd_store.put_artifacts([model.proto])[0]

    model_registry._store._mlmd_store.put_attributions_and_associations(
        [Attribution(context_id=mv_id, artifact_id=ma_id)], []
    )

    mlmd_ma = model_registry.get_model_artifact_by_params(model_version_id=str(mv_id))

    assert mlmd_ma.id == str(ma_id)
    assert mlmd_ma.name == model.py.name
    assert mlmd_ma.name != model.proto.name


def test_get_model_artifact_by_external_id(
    model_registry: ModelRegistry, model: Mapped
):
    model.proto.name = f"test_prefix:{model.proto.name}"
    model.proto.external_id = "external_id"
    model.py.external_id = "external_id"

    id = model_registry._store._mlmd_store.put_artifacts([model.proto])[0]
    id = str(id)

    mlmd_ma = model_registry.get_model_artifact_by_params(
        external_id=model.py.external_id
    )

    assert mlmd_ma.id == id
    assert mlmd_ma.name == model.py.name
    assert mlmd_ma.name != model.proto.name


def test_get_model_artifacts(
    model_registry: ModelRegistry, model: Mapped, model_version: Mapped
):
    mv1_id = model_registry._store._mlmd_store.put_contexts([model_version.proto])[0]

    model_version.proto.name = "version2"
    mv2_id = model_registry._store._mlmd_store.put_contexts([model_version.proto])[0]

    model.proto.name = "test_prefix:model1"
    ma1_id = model_registry._store._mlmd_store.put_artifacts([model.proto])[0]
    model.proto.name = "test_prefix:model2"
    ma2_id = model_registry._store._mlmd_store.put_artifacts([model.proto])[0]

    model_registry._store._mlmd_store.put_attributions_and_associations(
        [
            Attribution(context_id=mv1_id, artifact_id=ma1_id),
            Attribution(context_id=mv2_id, artifact_id=ma2_id),
        ],
        [],
    )

    mlmd_mas = model_registry.get_model_artifacts()
    assert len(mlmd_mas) == 2
    assert mlmd_mas[0].id in [str(ma1_id), str(ma2_id)]
