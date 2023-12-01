"""Tests for user facing model registry APIs."""

import pytest
from attrs import evolve
from ml_metadata.proto import (
    Artifact,
    Attribution,
    Context,
    ParentContext,
)
from model_registry.core import ModelRegistryAPIClient
from model_registry.exceptions import StoreException
from model_registry.store import MLMDStore
from model_registry.types import ModelArtifact, ModelVersion, RegisteredModel

from . import Mapped


@pytest.fixture()
def model(store_wrapper: MLMDStore) -> Mapped:
    art = Artifact()
    # we can't test the name directly as it's prefixed
    art.name = "model"
    art.type_id = store_wrapper.get_type_id(
        Artifact, ModelArtifact.get_proto_type_name()
    )

    art.uri = "uri"

    return Mapped(art, ModelArtifact("model", "uri"))


@pytest.fixture()
def model_version(store_wrapper: MLMDStore, model: Mapped) -> Mapped:
    ctx = Context()
    # we can't test the name directly as it's prefixed
    ctx.name = "version"
    ctx.type_id = store_wrapper.get_type_id(Context, ModelVersion.get_proto_type_name())
    ctx.properties["author"].string_value = "author"
    ctx.properties["model_name"].string_value = model.py.name
    ctx.properties["state"].string_value = "LIVE"

    return Mapped(ctx, ModelVersion(model.py.name, "version", "author"))


@pytest.fixture()
def registered_model(store_wrapper: MLMDStore, model: Mapped) -> Mapped:
    ctx = Context()
    ctx.name = model.py.name
    ctx.type_id = store_wrapper.get_type_id(
        Context, RegisteredModel.get_proto_type_name()
    )
    ctx.properties["state"].string_value = "LIVE"

    return Mapped(ctx, RegisteredModel(model.py.name))


# TODO: should we test insert/update separately?
def test_upsert_registered_model(
    mr_api: ModelRegistryAPIClient, registered_model: Mapped
):
    mr_api.upsert_registered_model(registered_model.py)

    rm_proto = mr_api._store._mlmd_store.get_context_by_type_and_name(
        RegisteredModel.get_proto_type_name(), registered_model.proto.name
    )
    assert rm_proto is not None
    assert registered_model.py.id == str(rm_proto.id)
    assert registered_model.py.name == rm_proto.name


def test_get_registered_model_by_id(
    mr_api: ModelRegistryAPIClient,
    registered_model: Mapped,
):
    rm_id = mr_api._store._mlmd_store.put_contexts([registered_model.proto])[0]

    mlmd_rm = mr_api.get_registered_model_by_id(str(rm_id))
    assert mlmd_rm.id == str(rm_id)
    assert mlmd_rm.name == registered_model.py.name
    assert mlmd_rm.name == registered_model.proto.name


def test_get_registered_model_by_name(
    mr_api: ModelRegistryAPIClient,
    registered_model: Mapped,
):
    rm_id = mr_api._store._mlmd_store.put_contexts([registered_model.proto])[0]

    mlmd_rm = mr_api.get_registered_model_by_params(name=registered_model.py.name)
    assert mlmd_rm.id == str(rm_id)
    assert mlmd_rm.name == registered_model.py.name
    assert mlmd_rm.name == registered_model.proto.name


def test_get_registered_model_by_external_id(
    mr_api: ModelRegistryAPIClient,
    registered_model: Mapped,
):
    registered_model.py.external_id = "external_id"
    registered_model.proto.external_id = "external_id"
    rm_id = mr_api._store._mlmd_store.put_contexts([registered_model.proto])[0]

    mlmd_rm = mr_api.get_registered_model_by_params(
        external_id=registered_model.py.external_id
    )
    assert mlmd_rm.id == str(rm_id)
    assert mlmd_rm.name == registered_model.py.name
    assert mlmd_rm.name == registered_model.proto.name


def test_get_registered_models(
    mr_api: ModelRegistryAPIClient, registered_model: Mapped
):
    rm1_id = mr_api._store._mlmd_store.put_contexts([registered_model.proto])[0]
    registered_model.proto.name = "model2"
    rm2_id = mr_api._store._mlmd_store.put_contexts([registered_model.proto])[0]

    mlmd_rms = mr_api.get_registered_models()
    assert len(mlmd_rms) == 2
    assert mlmd_rms[0].id in [str(rm1_id), str(rm2_id)]


def test_upsert_model_version(
    mr_api: ModelRegistryAPIClient,
    model_version: Mapped,
    registered_model: Mapped,
):
    rm_id = mr_api._store._mlmd_store.put_contexts([registered_model.proto])[0]
    rm_id = str(rm_id)

    mr_api.upsert_model_version(model_version.py, rm_id)

    mv_proto = mr_api._store._mlmd_store.get_context_by_type_and_name(
        ModelVersion.get_proto_type_name(), f"{rm_id}:{model_version.proto.name}"
    )
    assert mv_proto is not None
    assert model_version.py.id == str(mv_proto.id)
    assert model_version.py.version != mv_proto.name


def test_get_model_version_by_id(mr_api: ModelRegistryAPIClient, model_version: Mapped):
    model_version.proto.name = f"1:{model_version.proto.name}"
    ctx_id = mr_api._store._mlmd_store.put_contexts([model_version.proto])[0]

    id = str(ctx_id)
    mlmd_mv = mr_api.get_model_version_by_id(id)
    assert mlmd_mv.id == id
    assert mlmd_mv.name == model_version.py.name
    assert mlmd_mv.version != model_version.proto.name


def test_get_model_version_by_name(
    mr_api: ModelRegistryAPIClient, model_version: Mapped
):
    model_version.proto.name = f"1:{model_version.proto.name}"
    mv_id = mr_api._store._mlmd_store.put_contexts([model_version.proto])[0]

    mlmd_mv = mr_api.get_model_version_by_params(
        registered_model_id="1", version=model_version.py.name
    )
    assert mlmd_mv.id == str(mv_id)
    assert mlmd_mv.name == model_version.py.name
    assert mlmd_mv.name != model_version.proto.name


def test_get_model_version_by_external_id(
    mr_api: ModelRegistryAPIClient, model_version: Mapped
):
    model_version.proto.name = f"1:{model_version.proto.name}"
    model_version.proto.external_id = "external_id"
    model_version.py.external_id = "external_id"
    mv_id = mr_api._store._mlmd_store.put_contexts([model_version.proto])[0]

    mlmd_mv = mr_api.get_model_version_by_params(
        external_id=model_version.py.external_id
    )
    assert mlmd_mv.id == str(mv_id)
    assert mlmd_mv.name == model_version.py.name
    assert mlmd_mv.name != model_version.proto.name


def test_get_model_versions(
    mr_api: ModelRegistryAPIClient,
    model_version: Mapped,
    registered_model: Mapped,
):
    rm_id = mr_api._store._mlmd_store.put_contexts([registered_model.proto])[0]

    model_version.proto.name = f"{rm_id}:version"
    mv1_id = mr_api._store._mlmd_store.put_contexts([model_version.proto])[0]
    model_version.proto.name = f"{rm_id}:version2"
    mv2_id = mr_api._store._mlmd_store.put_contexts([model_version.proto])[0]

    mr_api._store._mlmd_store.put_parent_contexts(
        [
            ParentContext(parent_id=rm_id, child_id=mv1_id),
            ParentContext(parent_id=rm_id, child_id=mv2_id),
        ]
    )

    mlmd_mvs = mr_api.get_model_versions(str(rm_id))
    assert len(mlmd_mvs) == 2
    assert mlmd_mvs[0].id in [str(mv1_id), str(mv2_id)]


def test_upsert_model_artifact(
    monkeypatch,
    mr_api: ModelRegistryAPIClient,
    model: Mapped,
    model_version: Mapped,
):
    monkeypatch.setattr(ModelArtifact, "mlmd_name_prefix", "test_prefix")

    mv_id = mr_api._store._mlmd_store.put_contexts([model_version.proto])[0]
    mv_id = str(mv_id)

    mr_api.upsert_model_artifact(model.py, mv_id)

    ma_proto = mr_api._store._mlmd_store.get_artifact_by_type_and_name(
        ModelArtifact.get_proto_type_name(), f"test_prefix:{model.proto.name}"
    )
    assert ma_proto is not None
    assert model.py.id == str(ma_proto.id)
    assert model.py.name != ma_proto.name


def test_upsert_duplicate_model_artifact_with_different_version(
    mr_api: ModelRegistryAPIClient, model: Mapped, model_version: Mapped
):
    mv1_id = mr_api._store._mlmd_store.put_contexts([model_version.proto])[0]
    mv1_id = str(mv1_id)

    model_version.proto.name = "version2"
    mv2_id = mr_api._store._mlmd_store.put_contexts([model_version.proto])[0]
    mv2_id = str(mv2_id)

    ma1 = evolve(model.py)
    mr_api.upsert_model_artifact(ma1, mv1_id)
    ma2 = evolve(model.py)
    mr_api.upsert_model_artifact(ma2, mv2_id)

    ma_protos = mr_api._store._mlmd_store.get_artifacts_by_id(
        [int(ma1.id), int(ma2.id)]
    )
    assert ma1.name == ma2.name
    assert ma1.name != str(ma_protos[0].name)
    assert ma2.name != str(ma_protos[1].name)


def test_upsert_duplicate_model_artifact_with_same_version(
    mr_api: ModelRegistryAPIClient, model: Mapped, model_version: Mapped
):
    mv_id = mr_api._store._mlmd_store.put_contexts([model_version.proto])[0]
    mv_id = str(mv_id)

    ma1 = evolve(model.py)
    mr_api.upsert_model_artifact(ma1, mv_id)
    ma2 = evolve(model.py)
    with pytest.raises(StoreException):
        mr_api.upsert_model_artifact(ma2, mv_id)


def test_get_model_artifact_by_id(mr_api: ModelRegistryAPIClient, model: Mapped):
    model.proto.name = f"test_prefix:{model.proto.name}"
    id = mr_api._store._mlmd_store.put_artifacts([model.proto])[0]
    id = str(id)

    mlmd_ma = mr_api.get_model_artifact_by_id(id)

    assert mlmd_ma.id == id
    assert mlmd_ma.name == model.py.name
    assert mlmd_ma.name != model.proto.name


def test_get_model_artifact_by_model_version_id(
    mr_api: ModelRegistryAPIClient, model: Mapped, model_version: Mapped
):
    mv_id = mr_api._store._mlmd_store.put_contexts([model_version.proto])[0]

    model.proto.name = f"test_prefix:{model.proto.name}"
    ma_id = mr_api._store._mlmd_store.put_artifacts([model.proto])[0]

    mr_api._store._mlmd_store.put_attributions_and_associations(
        [Attribution(context_id=mv_id, artifact_id=ma_id)], []
    )

    mlmd_ma = mr_api.get_model_artifact_by_params(model_version_id=str(mv_id))

    assert mlmd_ma.id == str(ma_id)
    assert mlmd_ma.name == model.py.name
    assert mlmd_ma.name != model.proto.name


def test_get_model_artifact_by_external_id(
    mr_api: ModelRegistryAPIClient, model: Mapped
):
    model.proto.name = f"test_prefix:{model.proto.name}"
    model.proto.external_id = "external_id"
    model.py.external_id = "external_id"

    id = mr_api._store._mlmd_store.put_artifacts([model.proto])[0]
    id = str(id)

    mlmd_ma = mr_api.get_model_artifact_by_params(external_id=model.py.external_id)

    assert mlmd_ma.id == id
    assert mlmd_ma.name == model.py.name
    assert mlmd_ma.name != model.proto.name


def test_get_all_model_artifacts(mr_api: ModelRegistryAPIClient, model: Mapped):
    model.proto.name = "test_prefix:model1"
    ma1_id = mr_api._store._mlmd_store.put_artifacts([model.proto])[0]
    model.proto.name = "test_prefix:model2"
    ma2_id = mr_api._store._mlmd_store.put_artifacts([model.proto])[0]

    mlmd_mas = mr_api.get_model_artifacts()
    assert len(mlmd_mas) == 2
    assert mlmd_mas[0].id in [str(ma1_id), str(ma2_id)]


def test_get_model_artifacts_by_mv_id(
    mr_api: ModelRegistryAPIClient, model: Mapped, model_version: Mapped
):
    mv1_id = mr_api._store._mlmd_store.put_contexts([model_version.proto])[0]

    model_version.proto.name = "version2"
    mv2_id = mr_api._store._mlmd_store.put_contexts([model_version.proto])[0]

    model.proto.name = "test_prefix:model1"
    ma1_id = mr_api._store._mlmd_store.put_artifacts([model.proto])[0]
    model.proto.name = "test_prefix:model2"
    ma2_id = mr_api._store._mlmd_store.put_artifacts([model.proto])[0]

    mr_api._store._mlmd_store.put_attributions_and_associations(
        [
            Attribution(context_id=mv1_id, artifact_id=ma1_id),
            Attribution(context_id=mv2_id, artifact_id=ma2_id),
        ],
        [],
    )

    mlmd_mas = mr_api.get_model_artifacts(str(mv1_id))
    assert len(mlmd_mas) == 1
    assert mlmd_mas[0].id == str(ma1_id)
