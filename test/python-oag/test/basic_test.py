import asyncio
import pytest
import subprocess
import time
import os
import sys
import requests
from uuid import uuid4

import mr_openapi
from mr_openapi import RegisteredModelCreate
from mr_openapi import ModelVersion
from mr_openapi import ModelArtifact
from mr_openapi import DocArtifact
from mr_openapi import Artifact
from mr_openapi.api import ModelRegistryServiceApi

REGISTRY_HOST = "localhost"
REGISTRY_PORT = 8080
REGISTRY_URL = f"http://{REGISTRY_HOST}:{REGISTRY_PORT}"
MAX_POLL_TIME = 1200 # the first build is extremely slow if using docker-compose-*local*.yaml for bootstrap of builder image
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
    print(f" Starting Docker Compose in folder {model_registry_root_dir}")
    subprocess.call(f"{DOCKER} compose -f docker-compose.yaml build", shell=True, cwd=model_registry_root_dir)
    p = subprocess.Popen(f"{DOCKER} compose -f docker-compose.yaml up", shell=True, cwd=model_registry_root_dir)
    request.addfinalizer(p.kill)
    def teardown():
        print(f" Closing Docker Compose in folder {model_registry_root_dir}")
        subprocess.call(f"{DOCKER} compose -f docker-compose.yaml down", shell=True, cwd=model_registry_root_dir)

    request.addfinalizer(teardown)
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


def get_client():
    c = mr_openapi.Configuration(host=REGISTRY_URL)
    api_client = mr_openapi.ApiClient(configuration=c)
    client = ModelRegistryServiceApi(api_client=api_client)
    return client


def test_registered_model_create_and_retrieve():
    client = get_client()

    payload = RegisteredModelCreate()
    payload.name = f"FOO{uuid4()}"
    payload.description = "a foo"

    create_registered_model = client.create_registered_model(payload)
    assert create_registered_model is not None
    print(create_registered_model)
    print(create_registered_model.id)

    by_find = client.find_registered_model(name=create_registered_model.name)
    print(by_find)
    print(by_find.id)
    assert by_find.id == create_registered_model.id


def test_model_version_create_and_retrieve():
    client = get_client()

    rm = RegisteredModelCreate()
    rm.name = f"BAR{uuid4()}"
    rm.description = "a bar"

    create_registered_model = client.create_registered_model(rm)
    assert create_registered_model is not None
    print(create_registered_model)
    print(create_registered_model.id)

    payload = ModelVersion(registered_model_id=create_registered_model.id) # necessary here for required attrs, or else pydantic failure
    payload.author = "me"
    payload.name = "v1"
    payload.description = "a v1 for bar"

    create_model_version = client.create_registered_model_version(create_registered_model.id, payload)
    assert create_model_version is not None

    return_model_version = client.get_model_version(create_model_version.id)
    assert return_model_version is not None
    print(return_model_version)
    assert return_model_version.name == payload.name


def test_model_artifact_create_and_retrieve():
    client = get_client()

    rm = RegisteredModelCreate()
    rm.name = f"BAZ{uuid4()}"
    rm.description = "a baz"

    create_registered_model = client.create_registered_model(rm)
    assert create_registered_model is not None
    print(create_registered_model)
    print(create_registered_model.id)

    mv = ModelVersion(registered_model_id=create_registered_model.id) # necessary here for required attrs, or else pydantic failure
    mv.author = "me"
    mv.name = "v1"
    mv.description = "a v1 for baz"

    create_model_version = client.create_registered_model_version(create_registered_model.id, mv)
    assert create_model_version is not None

    payload = Artifact(ModelArtifact(artifact_type="model-artifact")) # necessary here for required attrs, or else pydantic failure
    payload.actual_instance.name = "mnist"
    payload.actual_instance.uri = "https://acme.org/mnist.onnx"

    create_model_artifact = client.create_model_version_artifact(create_model_version.id, payload)
    assert create_model_artifact is not None
    create_model_artifact = create_model_artifact.actual_instance
    assert create_model_artifact is not None
    print(create_model_artifact)
    print(create_model_artifact.id)
    assert isinstance(create_model_artifact, ModelArtifact)

    return_model_artifact = client.get_model_artifact(create_model_artifact.id)
    assert return_model_artifact is not None
    print(return_model_artifact)
    assert return_model_artifact.id == create_model_artifact.id

    payload = Artifact(DocArtifact(artifact_type="doc-artifact")) # necessary here for required attrs, or else pydantic failure
    payload.actual_instance.uri = "https://acme.org/README.md"

    create_doc_artifact = client.create_model_version_artifact(create_model_version.id, payload)
    assert create_doc_artifact is not None
    create_doc_artifact = create_doc_artifact.actual_instance
    assert create_doc_artifact is not None
    print(create_doc_artifact)
    assert create_doc_artifact.id != create_model_artifact.id
    
    list_artifacts = client.get_model_version_artifacts(create_model_version.id)
    assert list_artifacts is not None
    print(list_artifacts)
    assert list_artifacts.size == 2
