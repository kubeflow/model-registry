"""Tests creation and retrieval of base models."""

from collections.abc import AsyncIterator

import pytest

import mr_openapi
from mr_openapi import (
    Artifact,
    DocArtifact,
    MetadataValue,
    ModelArtifact,
    ModelRegistryServiceApi,
    ModelVersionCreate,
    RegisteredModelCreate,
)

from .conftest import REGISTRY_URL, cleanup


@pytest.fixture
@cleanup
async def client(user_token: str, verify_ssl: bool) -> AsyncIterator[ModelRegistryServiceApi]:
    params = {"verify_ssl": verify_ssl, "access_token": user_token}
    config = mr_openapi.Configuration(REGISTRY_URL, **params)  # type: ignore[arg-type]
    api_client = mr_openapi.ApiClient(config)
    client = mr_openapi.ModelRegistryServiceApi(api_client)
    yield client
    await api_client.close()


@pytest.fixture
def rm_create() -> RegisteredModelCreate:
    return RegisteredModelCreate(name="registered", description="a registered model")


@pytest.mark.e2e
async def test_registered_model(
    client: ModelRegistryServiceApi, rm_create: RegisteredModelCreate
):
    rm_create.custom_properties = {
        "key1": MetadataValue.from_dict(
            {"string_value": "value1", "metadataType": "MetadataStringValue"},
        )
    }

    new_rm = await client.create_registered_model(rm_create)
    print("created RM", new_rm, "with ID", new_rm.id)
    assert rm_create.name == new_rm.name
    assert rm_create.description == new_rm.description
    assert new_rm.custom_properties == rm_create.custom_properties

    by_find = await client.find_registered_model(name=new_rm.name)
    print("found RM", by_find, "with ID", by_find.id)
    assert by_find == new_rm
    assert by_find.id == new_rm.id
    assert new_rm.name == by_find.name
    assert new_rm.description == by_find.description


@pytest.fixture
async def mv_create(
    client: ModelRegistryServiceApi, rm_create: RegisteredModelCreate
) -> ModelVersionCreate:
    # HACK: create an RM first because we need an ID for the instance
    rm = await client.create_registered_model(rm_create)
    assert rm is not None
    return ModelVersionCreate(
        name="version",
        author="author",
        registeredModelId=str(rm.id),
        description="a model version",
    )


@pytest.mark.e2e
async def test_model_version(
    client: ModelRegistryServiceApi, mv_create: ModelVersionCreate
):
    mv_create.custom_properties = {
        "key1": MetadataValue.from_dict(
            {"string_value": "value1", "metadataType": "MetadataStringValue"},
        )
    }

    new_mv = await client.create_model_version(mv_create)
    print("created MV", new_mv, "with ID", new_mv.id)
    assert mv_create.name == new_mv.name
    assert mv_create.author == new_mv.author
    assert mv_create.description == new_mv.description
    assert mv_create.custom_properties == new_mv.custom_properties

    by_find = await client.get_model_version(str(new_mv.id))
    print("found MV", by_find)
    assert new_mv.id == by_find.id
    assert new_mv.name == by_find.name
    assert new_mv.author == by_find.author
    assert new_mv.description == by_find.description
    assert new_mv.custom_properties == by_find.custom_properties


@pytest.mark.e2e
async def test_model_artifact(
    client: ModelRegistryServiceApi, mv_create: ModelVersionCreate
):
    mv = await client.create_model_version(mv_create)
    assert mv is not None

    ma_create = ModelArtifact(
        name="model",
        uri="uri",
        artifactType="model-artifact",
        description="a model artifact",
        customProperties={
            "key1": MetadataValue.from_dict(
                {"string_value": "value1", "metadataType": "MetadataStringValue"},
            )
        },
    )

    new_ma = (
        await client.upsert_model_version_artifact(str(mv.id), Artifact(ma_create))
    ).actual_instance
    assert new_ma is not None
    print("created MA", new_ma, "with ID", new_ma.id)
    assert isinstance(new_ma, ModelArtifact)
    assert ma_create.name == new_ma.name
    assert ma_create.uri == new_ma.uri
    assert ma_create.description == new_ma.description
    assert ma_create.custom_properties == new_ma.custom_properties

    by_find = await client.get_model_artifact(str(new_ma.id))
    assert by_find is not None
    print("found MA", by_find)
    assert new_ma.id == by_find.id
    assert new_ma.name == by_find.name
    assert new_ma.uri == by_find.uri
    assert new_ma.description == by_find.description
    assert new_ma.custom_properties == by_find.custom_properties

    doc_art = DocArtifact(
        artifactType="doc-artifact",
        uri="https://acme.org/README.md",
        customProperties={
            "key1": MetadataValue.from_dict(
                {"string_value": "value1", "metadataType": "MetadataStringValue"},
            )
        },
    )

    new_da = (
        await client.upsert_model_version_artifact(str(mv.id), Artifact(doc_art))
    ).actual_instance
    assert new_da is not None
    print("created DA", new_da, "with ID", new_da.id)
    assert isinstance(new_da, DocArtifact)
    assert new_da.id != new_ma.id
    assert new_da.uri == doc_art.uri

    list_artifacts = await client.get_model_version_artifacts(str(mv.id))
    assert list_artifacts is not None
    print("list artifacts", list_artifacts)
    assert list_artifacts.size == 2
