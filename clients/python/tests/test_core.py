"""Tests for user facing model registry APIs."""

import pytest
from model_registry.core import ModelRegistryAPIClient
from model_registry.types import ModelArtifact, ModelVersion, RegisteredModel


@pytest.fixture()
def client():
    pass


def test_insert_registered_model(client: ModelRegistryAPIClient):
    registered_model = RegisteredModel(name="test rm")
    rm = client.upsert_registered_model(registered_model)
    assert rm.id
    assert rm.name == registered_model.name
    assert rm.external_id is None
    assert rm.description is None
    assert rm.create_time_since_epoch
    assert rm.last_update_time_since_epoch


def test_update_registered_model(client: ModelRegistryAPIClient):
    registered_model = RegisteredModel(name="updated rm")
    rm = client.upsert_registered_model(registered_model)
    last_update = rm.last_update_time_since_epoch
    rm.description = "lorem ipsum"
    rm = client.upsert_registered_model(rm)

    assert rm.description == "lorem ipsum"
    assert rm.last_update_time_since_epoch != last_update


@pytest.fixture()
def registered_model(client: ModelRegistryAPIClient) -> RegisteredModel:
    return client.upsert_registered_model(
        RegisteredModel(name="registered", external_id="mr id")
    )


def test_get_registered_model_by_id(
    client: ModelRegistryAPIClient,
    registered_model: RegisteredModel,
):
    assert (rm := client.get_registered_model_by_id(str(registered_model.id)))
    assert rm == registered_model


def test_get_registered_model_by_name(
    client: ModelRegistryAPIClient,
    registered_model: RegisteredModel,
):
    assert (rm := client.get_registered_model_by_params(name=registered_model.name))
    assert rm == registered_model


def test_get_registered_model_by_external_id(
    client: ModelRegistryAPIClient,
    registered_model: RegisteredModel,
):
    assert (
        rm := client.get_registered_model_by_params(
            external_id=registered_model.external_id
        )
    )
    assert rm == registered_model


def test_get_registered_models(
    client: ModelRegistryAPIClient, registered_model: RegisteredModel
):
    rm2 = client.upsert_registered_model(RegisteredModel(name="rm2"))

    rms = client.get_registered_models()
    assert [registered_model, rm2] == rms


def test_insert_model_version(
    client: ModelRegistryAPIClient,
    registered_model: RegisteredModel,
):
    model_version = ModelVersion(name="test version", author="test author")
    mv = client.upsert_model_version(model_version, str(registered_model.id))
    assert mv.id
    assert mv.name == model_version.name
    assert mv.external_id is None
    assert mv.description is None
    assert mv.create_time_since_epoch
    assert mv.last_update_time_since_epoch
    assert mv.author == model_version.author


def test_update_model_version(
    client: ModelRegistryAPIClient, registered_model: RegisteredModel
):
    model_version = ModelVersion(name="updated mv", author="test author")
    mv = client.upsert_model_version(model_version, str(registered_model.id))
    last_update = mv.last_update_time_since_epoch
    mv.description = "lorem ipsum"
    mv = client.upsert_model_version(mv, str(registered_model.id))

    assert mv.description == "lorem ipsum"
    assert mv.last_update_time_since_epoch != last_update


@pytest.fixture()
def model_version(
    client: ModelRegistryAPIClient, registered_model: RegisteredModel
) -> ModelVersion:
    return client.upsert_model_version(
        ModelVersion(name="version", author="author", external_id="mv id"),
        str(registered_model.id),
    )


def test_get_model_version_by_id(
    client: ModelRegistryAPIClient, model_version: ModelVersion
):
    assert (mv := client.get_model_version_by_id(str(model_version.id)))
    assert mv == model_version


def test_get_model_version_by_name(
    client: ModelRegistryAPIClient,
    registered_model: RegisteredModel,
    model_version: ModelVersion,
):
    assert (
        mv := client.get_model_version_by_params(
            registered_model_id=str(registered_model.id), name=model_version.name
        )
    )
    assert mv == model_version


def test_get_model_version_by_external_id(
    client: ModelRegistryAPIClient, model_version: ModelVersion
):
    assert (
        mv := client.get_model_version_by_params(external_id=model_version.external_id)
    )
    assert mv == model_version


def test_get_model_versions(
    client: ModelRegistryAPIClient,
    registered_model: RegisteredModel,
    model_version: ModelVersion,
):
    mv2 = client.upsert_model_version(
        ModelVersion(name="mv2", author="author"), str(registered_model.id)
    )

    mvs = client.get_model_versions(str(registered_model.id))
    assert [model_version, mv2] == mvs


def test_insert_model_artifact(
    client: ModelRegistryAPIClient,
    model_version: ModelVersion,
):
    props = {
        "name": "test model",
        "uri": "test uri",
        "model_format_name": "test format",
        "model_format_version": "test version",
        "storage_key": "test key",
        "storage_path": "test path",
        "service_account_name": "test service account",
    }
    ma = client.upsert_model_artifact(ModelArtifact(**props), str(model_version.id))
    assert ma.id
    assert ma.name == "test model"
    assert ma.uri
    assert ma.description is None
    assert ma.external_id is None
    assert ma.create_time_since_epoch
    assert ma.last_update_time_since_epoch
    assert ma.model_format_name
    assert ma.model_format_version
    assert ma.storage_key
    assert ma.storage_path
    assert ma.service_account_name


def test_update_model_artifact(
    client: ModelRegistryAPIClient, model_version: ModelVersion
):
    model = ModelArtifact(name="updated model", uri="uri")
    ma = client.upsert_model_artifact(model, str(model_version.id))
    last_update = ma.last_update_time_since_epoch
    ma.description = "lorem ipsum"
    ma = client.upsert_model_artifact(ma, str(model_version.id))

    assert ma.description == "lorem ipsum"
    assert ma.last_update_time_since_epoch != last_update


@pytest.fixture()
def model(
    client: ModelRegistryAPIClient,
    model_version: ModelVersion,
) -> ModelArtifact:
    return client.upsert_model_artifact(
        ModelArtifact(name="model", uri="uri", external_id="ma id"),
        str(model_version.id),
    )


def test_get_model_artifact_by_id(client: ModelRegistryAPIClient, model: ModelArtifact):
    assert (ma := client.get_model_artifact_by_id(str(model.id)))
    assert ma == model


def test_get_model_artifact_by_name(
    client: ModelRegistryAPIClient, model_version: ModelVersion, model: ModelArtifact
):
    assert (
        ma := client.get_model_artifact_by_params(
            name=model.name, model_version_id=str(model_version.id)
        )
    )
    assert ma == model


def test_get_model_artifact_by_external_id(
    client: ModelRegistryAPIClient, model: ModelArtifact
):
    assert (ma := client.get_model_artifact_by_params(external_id=model.external_id))
    assert ma == model


def test_get_all_model_artifacts(
    client: ModelRegistryAPIClient, model_version: ModelVersion, model: ModelArtifact
):
    ma2 = client.upsert_model_artifact(
        ModelArtifact(name="ma2", uri="uri"), str(model_version.id)
    )

    mas = client.get_model_artifacts()
    assert [model, ma2] == mas


def test_get_model_artifacts_by_mv_id(
    client: ModelRegistryAPIClient, model_version: ModelVersion, model: ModelArtifact
):
    ma2 = client.upsert_model_artifact(
        ModelArtifact(name="ma2", uri="uri"), str(model_version.id)
    )

    mas = client.get_model_artifacts(str(model_version.id))
    assert [model, ma2] == mas
