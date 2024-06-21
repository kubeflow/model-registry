"""Tests creation and retrieval of base models."""

import asyncio
import os
import subprocess
import time
from time import sleep

import mr_openapi
import pytest
import requests
from mr_openapi import (
    Artifact,
    DocArtifact,
    MetadataValue,
    ModelArtifact,
    ModelVersionCreate,
    RegisteredModelCreate,
)

REGISTRY_HOST = "http://localhost"
REGISTRY_PORT = 8080
REGISTRY_URL = f"{REGISTRY_HOST}:{REGISTRY_PORT}"
COMPOSE_FILE = "docker-compose.yaml"
MAX_POLL_TIME = 1200  # the first build is extremely slow if using docker-compose-*local*.yaml for bootstrap of builder image
POLL_INTERVAL = 1
DOCKER = os.getenv("DOCKER", "docker")


def poll_for_ready():
    start_time = time.time()
    while True:
        elapsed_time = time.time() - start_time
        if elapsed_time >= MAX_POLL_TIME:
            print("Polling timed out.")
            break

        print("Attempt to connect")
        try:
            response = requests.get(REGISTRY_URL, timeout=MAX_POLL_TIME)
            if response.status_code == 404:
                print("Server is up!")
                break
        except requests.exceptions.ConnectionError:
            pass

        # Wait for the specified poll interval before trying again
        time.sleep(POLL_INTERVAL)


@pytest.fixture(scope="session", autouse=True)
def _compose_mr(root):
    print("Assuming this is the Model Registry root directory:", root)
    shared_volume = root / "test/config/ml-metadata"
    sqlite_db_file = shared_volume / "metadata.sqlite.db"
    if sqlite_db_file.exists():
        msg = f"The file {sqlite_db_file} already exists; make sure to cancel it before running these tests."
        raise FileExistsError(msg)
    print(f" Starting Docker Compose in folder {root}")
    p = subprocess.Popen(
        f"{DOCKER} compose -f {COMPOSE_FILE} up --build",
        shell=True,  # noqa: S602
        cwd=root,
    )
    yield

    p.kill()
    print(f" Closing Docker Compose in folder {root}")
    subprocess.call(
        f"{DOCKER} compose -f {COMPOSE_FILE} down",
        shell=True,  # noqa: S602
        cwd=root,
    )
    try:
        os.remove(sqlite_db_file)
        print(f"Removed {sqlite_db_file} successfully.")
    except Exception as e:
        print(f"An error occurred while removing {sqlite_db_file}: {e}")


# workaround: https://github.com/pytest-dev/pytest-asyncio/issues/706#issuecomment-2147044022
@pytest.fixture(scope="session", autouse=True)
def event_loop():
    loop = asyncio.get_event_loop_policy().get_event_loop()
    yield loop
    loop.close()


@pytest.fixture()
async def client(root):
    poll_for_ready()

    config = mr_openapi.Configuration(REGISTRY_URL)
    api_client = mr_openapi.ApiClient(config)
    client = mr_openapi.ModelRegistryServiceApi(api_client)
    yield client
    await api_client.close()

    sqlite_db_file = root / "test/config/ml-metadata/metadata.sqlite.db"
    try:
        os.remove(sqlite_db_file)
        print(f"Removed {sqlite_db_file} successfully.")
    except Exception as e:
        print(f"An error occurred while removing {sqlite_db_file}: {e}")
    # we have to wait to make sure the server restarts after the file is gone
    sleep(1)

    print("Restarting model-registry...")
    subprocess.call(
        f"{DOCKER} compose -f {COMPOSE_FILE} restart model-registry",
        shell=True,  # noqa: S602
        cwd=root,
    )


@pytest.fixture()
def rm_create() -> RegisteredModelCreate:
    return RegisteredModelCreate(name="registered", description="a registered model")


@pytest.fixture()
async def mv_create(client, rm_create) -> ModelVersionCreate:
    # HACK: create an RM first because we need an ID for the instance
    rm = await client.create_registered_model(rm_create)
    assert rm is not None
    return ModelVersionCreate(
        name="version",
        author="author",
        registeredModelId=rm.id,
        description="a model version",
    )


async def test_registered_model(client, rm_create):
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


async def test_model_version(client, mv_create):
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

    by_find = await client.get_model_version(new_mv.id)
    print("found MV", by_find)
    assert new_mv.id == by_find.id
    assert new_mv.name == by_find.name
    assert new_mv.author == by_find.author
    assert new_mv.description == by_find.description
    assert new_mv.custom_properties == by_find.custom_properties


async def test_model_artifact(client, mv_create):
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
        await client.create_model_version_artifact(mv.id, Artifact(ma_create))
    ).actual_instance
    assert new_ma is not None
    print("created MA", new_ma, "with ID", new_ma.id)
    assert isinstance(new_ma, ModelArtifact)
    assert ma_create.name == new_ma.name
    assert ma_create.uri == new_ma.uri
    assert ma_create.description == new_ma.description
    assert ma_create.custom_properties == new_ma.custom_properties

    by_find = await client.get_model_artifact(new_ma.id)
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
        await client.create_model_version_artifact(mv.id, Artifact(doc_art))
    ).actual_instance
    assert new_da is not None
    print("created DA", new_da, "with ID", new_da.id)
    assert isinstance(new_da, DocArtifact)
    assert new_da.id != new_ma.id
    assert new_da.uri == doc_art.uri

    list_artifacts = await client.get_model_version_artifacts(mv.id)
    assert list_artifacts is not None
    print("list artifacts", list_artifacts)
    assert list_artifacts.size == 2
