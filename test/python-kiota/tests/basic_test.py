import asyncio
from dataclasses import dataclass
from typing import Optional
from httpx import QueryParams
import pytest
import subprocess
import time
import os
import sys
import requests
import json
from kiota_abstractions.headers_collection import HeadersCollection
from kiota_abstractions.base_request_configuration import RequestConfiguration
from kiota_abstractions.authentication.anonymous_authentication_provider import (
    AnonymousAuthenticationProvider,
)
from kiota_http.httpx_request_adapter import HttpxRequestAdapter
from apisdk.client.registry_client import RegistryClient
from apisdk.client.models.registered_model_create import RegisteredModelCreate
from apisdk.client.models.registered_model_state import RegisteredModelState
from apisdk.client.api.model_registry.v1alpha3.registered_model.registered_model_request_builder import Registered_modelRequestBuilder
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
    root_folder = os.path.join(sys.path[0], "..", "..", "..")
    print(f" Starting Docker Compose in folder {root_folder}")
    subprocess.call(f"{DOCKER} compose -f docker-compose-local.yaml build", shell=True, cwd=root_folder)
    p = subprocess.Popen(f"{DOCKER} compose -f docker-compose-local.yaml up", shell=True, cwd=root_folder)
    request.addfinalizer(p.kill)
    request.addfinalizer(cleanup)
    poll_for_ready()
    
def cleanup():
    root_folder = os.path.join(sys.path[0], "..", "..", "..")
    print(f" Closing Docker Compose in folder {root_folder}")
    subprocess.Popen(f"{DOCKER} compose -f docker-compose-local.yaml down", shell=False, cwd=root_folder)

# workaround: https://stackoverflow.com/a/72104554
@pytest.fixture(scope="session", autouse=True)
def event_loop():
    try:
        loop = asyncio.get_running_loop()
    except RuntimeError:
        loop = asyncio.new_event_loop()
    yield loop
    loop.close()

# registered Model
# registered version
# model artifact

@pytest.mark.asyncio
async def test_registered_model_create_and_retrieve():
    auth_provider = AnonymousAuthenticationProvider()
    request_adapter = HttpxRequestAdapter(auth_provider)
    request_adapter.base_url = REGISTRY_URL
    client = RegistryClient(request_adapter)

    payload = RegisteredModelCreate()
    payload.name = "FOO"
    payload.description = "a foo"

    # TODO: doesn't work it infer type_id = 10 for some reasons
    create_registered_model = await client.api.model_registry.v1alpha3.registered_models.post(payload)
    assert create_registered_model is not None

    query_params = Registered_modelRequestBuilder.Registered_modelRequestBuilderGetQueryParameters(
        name= create_registered_model.name
    )
    return_model_artifact = await client.api.model_registry.v1alpha3.registered_model.get(RequestConfiguration(query_params=query_params))
    print(return_model_artifact)
