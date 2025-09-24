"""Tests for user facing model registry APIs."""

import pytest

from model_registry.core import ModelRegistryAPIClient
from model_registry.types import (
    DocArtifact,
    ModelArtifact,
    ModelVersion,
    Pager,
    RegisteredModel,
)

from .conftest import REGISTRY_HOST, REGISTRY_PORT, cleanup


@pytest.fixture
@cleanup
def client(user_token) -> ModelRegistryAPIClient:
    return ModelRegistryAPIClient.insecure_connection(REGISTRY_HOST, REGISTRY_PORT, user_token=user_token)


@pytest.mark.e2e
async def test_insert_registered_model(client: ModelRegistryAPIClient):
    registered_model = RegisteredModel(name="test rm")
    rm = await client.upsert_registered_model(registered_model)
    assert rm.id
    assert rm.name == registered_model.name
    assert rm.external_id is None
    assert rm.description is None
    assert rm.create_time_since_epoch
    assert rm.last_update_time_since_epoch


@pytest.mark.e2e
async def test_update_registered_model(client: ModelRegistryAPIClient):
    registered_model = RegisteredModel(name="updated rm")
    rm = await client.upsert_registered_model(registered_model)
    last_update = rm.last_update_time_since_epoch
    rm.description = "lorem ipsum"
    rm = await client.upsert_registered_model(rm)

    assert rm.description == "lorem ipsum"
    assert rm.last_update_time_since_epoch != last_update


@pytest.fixture
async def registered_model(client: ModelRegistryAPIClient) -> RegisteredModel:
    return await client.upsert_registered_model(
        RegisteredModel(name="registered", external_id="mr id")
    )


@pytest.mark.e2e
async def test_get_registered_model_by_id(
        client: ModelRegistryAPIClient,
        registered_model: RegisteredModel,
):
    assert (rm := await client.get_registered_model_by_id(str(registered_model.id)))
    assert rm == registered_model


@pytest.mark.e2e
async def test_get_registered_model_by_name(
        client: ModelRegistryAPIClient,
        registered_model: RegisteredModel,
):
    assert (
        rm := await client.get_registered_model_by_params(name=registered_model.name)
    )
    assert rm == registered_model


@pytest.mark.e2e
async def test_get_registered_model_by_external_id(
        client: ModelRegistryAPIClient,
        registered_model: RegisteredModel,
):
    assert registered_model.external_id
    assert (
        rm := await client.get_registered_model_by_params(
            external_id=registered_model.external_id
        )
    )
    assert rm == registered_model


@pytest.mark.e2e
async def test_get_registered_models(
        client: ModelRegistryAPIClient, registered_model: RegisteredModel
):
    rm2 = await client.upsert_registered_model(RegisteredModel(name="rm2"))

    rms = await client.get_registered_models()
    assert [registered_model, rm2] == rms


@pytest.mark.e2e
async def test_page_through_registered_models(client: ModelRegistryAPIClient):
    models = 6
    for i in range(models):
        await client.upsert_registered_model(RegisteredModel(name=f"rm{i}"))
    pager = Pager(client.get_registered_models).page_size(5)
    total = 0
    async for _ in pager:
        total += 1
    assert total == models


@pytest.mark.e2e
async def test_insert_model_version(
        client: ModelRegistryAPIClient,
        registered_model: RegisteredModel,
):
    model_version = ModelVersion(name="test version", author="test author")
    mv = await client.upsert_model_version(model_version, str(registered_model.id))
    assert mv.id
    assert mv.name == model_version.name
    assert mv.external_id is None
    assert mv.description is None
    assert mv.create_time_since_epoch
    assert mv.last_update_time_since_epoch
    assert mv.author == model_version.author


@pytest.mark.e2e
async def test_update_model_version(
        client: ModelRegistryAPIClient, registered_model: RegisteredModel
):
    model_version = ModelVersion(name="updated mv", author="test author")
    mv = await client.upsert_model_version(model_version, str(registered_model.id))
    last_update = mv.last_update_time_since_epoch
    mv.description = "lorem ipsum"
    mv = await client.upsert_model_version(mv, str(registered_model.id))

    assert mv.description == "lorem ipsum"
    assert mv.last_update_time_since_epoch != last_update


@pytest.fixture
async def model_version(
        client: ModelRegistryAPIClient, registered_model: RegisteredModel
) -> ModelVersion:
    return await client.upsert_model_version(
        ModelVersion(name="version", author="author", external_id="mv id"),
        str(registered_model.id),
    )


@pytest.mark.e2e
async def test_get_model_version_by_id(
        client: ModelRegistryAPIClient, model_version: ModelVersion
):
    assert (mv := await client.get_model_version_by_id(str(model_version.id)))
    assert mv == model_version


@pytest.mark.e2e
async def test_get_model_version_by_name(
        client: ModelRegistryAPIClient,
        registered_model: RegisteredModel,
        model_version: ModelVersion,
):
    assert (
        mv := await client.get_model_version_by_params(
            registered_model_id=str(registered_model.id), name=model_version.name
        )
    )
    assert mv == model_version


@pytest.mark.e2e
async def test_get_model_version_by_external_id(
        client: ModelRegistryAPIClient, model_version: ModelVersion
):
    assert (
        mv := await client.get_model_version_by_params(
            external_id=str(model_version.external_id)
        )
    )
    assert mv == model_version


@pytest.mark.e2e
async def test_get_model_versions(
        client: ModelRegistryAPIClient,
        registered_model: RegisteredModel,
        model_version: ModelVersion,
):
    mv2 = await client.upsert_model_version(
        ModelVersion(name="mv2", author="author"), str(registered_model.id)
    )

    mvs = await client.get_model_versions(str(registered_model.id))
    assert [model_version, mv2] == mvs


@pytest.mark.e2e
async def test_page_through_model_versions(
        client: ModelRegistryAPIClient, registered_model: RegisteredModel
):
    models = 6
    for i in range(models):
        await client.upsert_model_version(
            ModelVersion(name=f"mv{i}"), str(registered_model.id)
        )
    pager = Pager(
        lambda o: client.get_model_versions(str(registered_model.id), o)
    ).page_size(5)
    total = 0
    async for _ in pager:
        total += 1
    assert total == models


@pytest.mark.e2e
async def test_insert_model_version_artifact(
        client: ModelRegistryAPIClient, model_version: ModelVersion
):
    model = DocArtifact(
        name="test model",
        uri="test uri",
    )
    assert model_version.id
    da = await client.upsert_model_version_artifact(model, model_version.id)
    assert da.id
    assert da.name == "test model"
    assert da.uri
    assert da.description is None
    assert da.external_id is None


@pytest.mark.e2e
async def test_update_model_version_artifact(
        client: ModelRegistryAPIClient, model_version: ModelVersion
):
    model = DocArtifact(name="updated model", uri="uri")
    da = await client.upsert_model_version_artifact(model, str(model_version.id))
    assert da.id
    last_update = da.last_update_time_since_epoch
    da.description = "lorem ipsum"
    da = await client.upsert_model_version_artifact(da, str(model_version.id))

    assert da.description == "lorem ipsum"
    assert da.last_update_time_since_epoch != last_update


@pytest.mark.e2e
async def test_insert_model_artifact(
        client: ModelRegistryAPIClient,
):
    model = ModelArtifact(
        name="test model",
        uri="test uri",
        model_format_name="test format",
        model_format_version="test version",
        storage_key="test key",
        storage_path="test path",
        service_account_name="test service account",
        model_source_kind="test source kind",
        model_source_class="test source class",
        model_source_group="test source group",
        model_source_id="test source id",
        model_source_name="test source name",
    )
    ma = await client.upsert_model_artifact(model)
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
    assert ma.model_source_kind
    assert ma.model_source_class
    assert ma.model_source_group
    assert ma.model_source_id
    assert ma.model_source_name


@pytest.mark.e2e
async def test_update_model_artifact(client: ModelRegistryAPIClient):
    model = ModelArtifact(name="updated model", uri="uri")
    ma = await client.upsert_model_artifact(model)
    last_update = ma.last_update_time_since_epoch
    ma.description = "lorem ipsum"
    ma = await client.upsert_model_artifact(ma)

    assert ma.description == "lorem ipsum"
    assert ma.last_update_time_since_epoch != last_update


@pytest.fixture
async def model(
        client: ModelRegistryAPIClient,
        model_version: ModelVersion,
) -> ModelArtifact:
    return await client.upsert_model_version_artifact(
        ModelArtifact(name="model", uri="uri", external_id="ma id"),
        str(model_version.id),
    )


@pytest.mark.e2e
async def test_get_model_artifact_by_id(
        client: ModelRegistryAPIClient, model: ModelArtifact
):
    assert (ma := await client.get_model_artifact_by_id(str(model.id)))
    assert ma == model


@pytest.mark.e2e
async def test_get_model_artifact_by_name(
        client: ModelRegistryAPIClient, model_version: ModelVersion, model: ModelArtifact
):
    assert (
        ma := await client.get_model_artifact_by_params(
            name=str(model.name), model_version_id=str(model_version.id)
        )
    )
    assert ma == model


@pytest.mark.e2e
async def test_get_model_artifact_by_external_id(
        client: ModelRegistryAPIClient, model: ModelArtifact
):
    assert (
        ma := await client.get_model_artifact_by_params(
            external_id=str(model.external_id)
        )
    )
    assert ma == model


@pytest.mark.e2e
async def test_get_all_model_artifacts(
        client: ModelRegistryAPIClient, model: ModelArtifact
):
    ma2 = await client.upsert_model_artifact(ModelArtifact(name="ma2", uri="uri"))

    mas = await client.get_model_artifacts()
    assert [model, ma2] == mas


@pytest.mark.e2e
async def test_get_model_version_artifacts_by_mv_id(
        client: ModelRegistryAPIClient, model_version: ModelVersion, model: ModelArtifact
):
    ma2 = await client.upsert_model_version_artifact(
        ModelArtifact(name="ma2", uri="uri"), str(model_version.id)
    )

    mas = await client.get_model_artifacts(str(model_version.id))
    assert [model, ma2] == mas


@pytest.mark.e2e
async def test_page_through_model_version_artifacts(
        client: ModelRegistryAPIClient,
        registered_model: RegisteredModel,
        model_version: ModelVersion,
):
    _ = registered_model
    models = 6
    for i in range(models):
        art = ModelArtifact(name=f"ma{i}", uri="uri") if i % 2 == 0 else DocArtifact(name=f"ma{i}", uri="uri")
        await client.upsert_model_version_artifact(art, str(model_version.id))
    pager = Pager(
        lambda o: client.get_model_version_artifacts(str(model_version.id), o)
    ).page_size(5)
    total = 0
    async for _ in pager:
        total += 1
    assert total == models
