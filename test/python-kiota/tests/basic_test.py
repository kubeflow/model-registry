import asyncio
from apisdk.client.models.model_artifact import ModelArtifact
from apisdk.client.models.doc_artifact import DocArtifact
from apisdk.client.models.model_version import ModelVersion
import pytest
import subprocess
import time
import os
import sys
import requests
from uuid import uuid4
from kiota_abstractions.base_request_configuration import RequestConfiguration
from kiota_abstractions.authentication.anonymous_authentication_provider import (
    AnonymousAuthenticationProvider,
)
from kiota_http.httpx_request_adapter import HttpxRequestAdapter
from apisdk.client.registry_client import RegistryClient
from apisdk.client.models.registered_model_create import RegisteredModelCreate
from apisdk.client.models.model_version_create import ModelVersionCreate
from apisdk.client.api.model_registry.v1alpha3.registered_model.registered_model_request_builder import Registered_modelRequestBuilder
from apisdk.client.api.model_registry.v1alpha3.model_version.model_version_request_builder import Model_versionRequestBuilder
from apisdk.client.api.model_registry.v1alpha3.model_versions.model_versions_request_builder import Model_versionsRequestBuilder
# from apisdk.client.models.model_artifact_create import ModelArtifactCreate
# from apisdk.client.api.model_registry.v1alpha3.model_artifact.model_artifact_request_builder import Model_artifactRequestBuilder


REGISTRY_HOST = "localhost"
REGISTRY_PORT = 8080
REGISTRY_URL = f"http://{REGISTRY_HOST}:{REGISTRY_PORT}"
MAX_POLL_TIME = 1200 # the first build is extremely slow
POLL_INTERVAL = 1
DOCKER = os.getenv("DOCKER", "docker")
start_time = time.time()


def poll_for_ready():
    while True:
        elapsed_time = time.time() - start_time
        if elapsed_time >= MAX_POLL_TIME:
            print("Polling timed out.")
            break

        print("Attempt to connect")
        try:
            response = requests.get(REGISTRY_URL)
            if response.status_code == 404:
                print("Server is up!")
                break
        except requests.exceptions.ConnectionError:
            pass

        # Wait for the specified poll interval before trying again
        time.sleep(POLL_INTERVAL)


@pytest.fixture(scope="session", autouse=True)
def registry_server(request):
    model_registry_root_dir = model_registry_root(request)
    print(
        "Assuming this is the Model Registry root directory:", model_registry_root_dir
    )
    shared_volume = model_registry_root_dir / "test/config/ml-metadata"
    sqlite_db_file = shared_volume / "metadata.sqlite.db"
    if sqlite_db_file.exists():
        msg = f"The file {sqlite_db_file} already exists; make sure to cancel it before running these tests."
        raise FileExistsError(msg)
    root_folder = os.path.join(sys.path[0], "..", "..", "..")
    print(f" Starting Docker Compose in folder {root_folder}")
    subprocess.call(f"{DOCKER} compose -f docker-compose-local.yaml build", shell=True, cwd=root_folder)
    p = subprocess.Popen(f"{DOCKER} compose -f docker-compose-local.yaml up", shell=True, cwd=root_folder)
    request.addfinalizer(p.kill)
    request.addfinalizer(cleanup)
    poll_for_ready()


def model_registry_root(request):
    return (request.config.rootpath / "../..").resolve()  # resolves to absolute path


@pytest.fixture(scope="session", autouse=True)
def plain_wrapper(request):
    sqlite_db_file = (
        model_registry_root(request) / "test/config/ml-metadata/metadata.sqlite.db"
    )

    def teardown():
        try:
            os.remove(sqlite_db_file)
            print(f"Removed {sqlite_db_file} successfully.")
        except Exception as e:
            print(f"An error occurred while removing {sqlite_db_file}: {e}")
        print("plain_wrapper_after_each done.")

    request.addfinalizer(teardown)


def cleanup():
    root_folder = os.path.join(sys.path[0], "..", "..", "..")
    print(f" Closing Docker Compose in folder {root_folder}")
    subprocess.call(f"{DOCKER} compose -f docker-compose-local.yaml down", shell=True, cwd=root_folder)

# workaround: https://stackoverflow.com/a/72104554
@pytest.fixture(scope="session", autouse=True)
def event_loop():
    try:
        loop = asyncio.get_running_loop()
    except RuntimeError:
        loop = asyncio.new_event_loop()
    yield loop
    loop.close()


def get_client():
    auth_provider = AnonymousAuthenticationProvider()
    request_adapter = HttpxRequestAdapter(auth_provider)
    request_adapter.base_url = REGISTRY_URL
    client = RegistryClient(request_adapter)
    return client


@pytest.mark.asyncio
async def test_registered_model_create_and_retrieve():
    client = get_client()

    payload = RegisteredModelCreate()
    payload.name = f"FOO{uuid4()}"
    payload.description = "a foo"

    create_registered_model = await client.api.model_registry.v1alpha3.registered_models.post(payload)
    assert create_registered_model is not None
    print(create_registered_model)
    print(create_registered_model.id)

    query_params = Registered_modelRequestBuilder.Registered_modelRequestBuilderGetQueryParameters(
        name= create_registered_model.name
    )
    return_model_artifact = await client.api.model_registry.v1alpha3.registered_model.get(RequestConfiguration(query_parameters=query_params))
    print(return_model_artifact)
    print(return_model_artifact.id)


@pytest.mark.asyncio
async def test_model_version_create_and_retrieve():
    client = get_client()

    rm = RegisteredModelCreate()
    rm.name = f"BAR{uuid4()}"
    rm.description = "a bar"

    create_registered_model = await client.api.model_registry.v1alpha3.registered_models.post(rm)
    assert create_registered_model is not None
    print(create_registered_model)
    print(create_registered_model.id)

    payload = ModelVersion()
    payload.author = "me"
    payload.name = "v1"
    payload.description = "a v1 for bar"

    create_model_version = await client.api.model_registry.v1alpha3.registered_models.by_registeredmodel_id(create_registered_model.id).versions.post(payload)
    assert create_model_version is not None

    return_model_version = await client.api.model_registry.v1alpha3.model_versions.by_modelversion_id(create_model_version.id).get()
    assert return_model_version is not None
    print(return_model_version)


@pytest.mark.asyncio
async def test_model_artifact_create_and_retrieve():
    client = get_client()

    rm = RegisteredModelCreate()
    rm.name = f"BAZ{uuid4()}"
    rm.description = "a baz"

    create_registered_model = await client.api.model_registry.v1alpha3.registered_models.post(rm)
    assert create_registered_model is not None
    print(create_registered_model)
    print(create_registered_model.id)

    mv = ModelVersion()
    mv.author = "me"
    mv.name = "v1"
    mv.description = "a v1 for baz"

    create_model_version = await client.api.model_registry.v1alpha3.registered_models.by_registeredmodel_id(create_registered_model.id).versions.post(mv)
    assert create_model_version is not None

    payload = ModelArtifact()
    payload.name = "mnist"
    payload.uri = "https://acme.org/mnist.onnx"

    create_model_artifact = await client.api.model_registry.v1alpha3.model_versions.by_modelversion_id(create_model_version.id).artifacts.post(payload)
    assert create_model_artifact is not None
    create_model_artifact = create_model_artifact.model_artifact
    print(create_model_artifact)
    print(create_model_artifact.id)

    return_model_artifact = await client.api.model_registry.v1alpha3.model_artifacts.by_modelartifact_id(create_model_artifact.id).get()
    assert return_model_artifact is not None
    print(return_model_artifact)

    payload = DocArtifact()
    payload.uri = "https://acme.org/mnist.onnx"

    create_doc_artifact = await client.api.model_registry.v1alpha3.model_versions.by_modelversion_id(create_model_version.id).artifacts.post(payload)
    assert create_doc_artifact is not None
    create_doc_artifact = create_doc_artifact.doc_artifact
    assert create_doc_artifact is not None
    print(create_doc_artifact)
    
    # How to retrieve a doc_artifact?
