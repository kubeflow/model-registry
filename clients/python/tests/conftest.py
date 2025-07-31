import base64
import contextlib
import inspect
import json
import os
import pathlib
import shutil
import subprocess
import tempfile
import time
from collections.abc import Generator
from contextlib import asynccontextmanager
from pathlib import Path
from typing import Any
from unittest.mock import Mock, patch
from urllib.parse import urlparse

import pytest
import requests
import schemathesis
from schemathesis import Case, Response
from schemathesis.generation.stateful.state_machine import APIStateMachine
from schemathesis.specs.openapi.schemas import BaseOpenAPISchema

from model_registry import ModelRegistry
from model_registry.utils import BackendDefinition, _get_skopeo_backend

from .constants import DEFAULT_API_TIMEOUT


def pytest_addoption(parser):
    parser.addoption("--e2e", action="store_true", help="run end-to-end tests")
    parser.addoption("--fuzz", action="store_true", help="run fuzzing tests")


def pytest_collection_modifyitems(config, items):
    skip_reasons = {
        "e2e": pytest.mark.skip(reason="this is an end-to-end test, requires explicit opt-in --e2e option to run."),
        "fuzz": pytest.mark.skip(reason="this is a fuzzing test, requires explicit opt-in --fuzz option to run."),
        "skip": pytest.mark.skip(reason="skipping non-e2e and non-fuzz tests"),
    }
    e2e = config.getoption("--e2e")
    fuzz = config.getoption("--fuzz")

    for item in items:
        if e2e:
            if "e2e" not in item.keywords:
                item.add_marker(skip_reasons["skip"])
        elif fuzz:
            if "fuzz" not in item.keywords:
                item.add_marker(skip_reasons["skip"])
        else:
            if "e2e" in item.keywords:
                item.add_marker(skip_reasons["e2e"])
            if "fuzz" in item.keywords:
                item.add_marker(skip_reasons["fuzz"])


def pytest_report_teststatus(report, config):
    test_name = report.head_line
    if report.passed:
        if report.when == "call":
            print(f"\nTEST: {test_name} STATUS: \033[0;32mPASSED\033[0m")

    elif report.skipped:
        print(f"\nTEST: {test_name} STATUS: \033[1;33mSKIPPED\033[0m")

    elif report.failed:
        if report.when != "call":
            print(f"\nTEST: {test_name} [{report.when}] STATUS: \033[0;31mERROR\033[0m")
        else:
            print(f"\nTEST: {test_name} STATUS: \033[0;31mFAILED\033[0m")


REGISTRY_URL = os.environ.get("MR_URL", "http://localhost:8080")
parsed = urlparse(REGISTRY_URL)
host, port = parsed.netloc.split(":")
REGISTRY_HOST = f"{parsed.scheme}://{host}"
REGISTRY_PORT = int(port)

MAX_POLL_TIME = 10
POLL_INTERVAL = 1
start_time = time.time()


@pytest.fixture(scope="session")
def root(request) -> Path:
    return (request.config.rootpath / "../..").resolve()  # resolves to absolute path


def poll_for_ready():
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


def cleanup(client):
    async def yield_and_restart(root):
        poll_for_ready()
        if inspect.iscoroutinefunction(client) or inspect.isasyncgenfunction(client):
            async with asynccontextmanager(client)() as async_client:
                yield async_client
        else:
            yield client()

        print("Cleaning DB...")
        subprocess.call(  # noqa: S602
            "./scripts/cleanup.sh",
            shell=True,
            cwd=root,
        )

    return yield_and_restart




@pytest.fixture
@cleanup
def client() -> ModelRegistry:
    return ModelRegistry(REGISTRY_HOST, REGISTRY_PORT, author="author", is_secure=False)


@pytest.fixture
@cleanup
def client_attrs() -> dict[str, any]:
    return {
        "host": REGISTRY_HOST,
        "port": REGISTRY_PORT,
        "author": "author",
        "ssl": False,
    }


@pytest.fixture(scope="module")
def setup_env_user_token():
    with tempfile.NamedTemporaryFile(delete=False) as token_file:
        token_file.write(b"Token")
    old_token_path = os.getenv("KF_PIPELINES_SA_TOKEN_PATH")
    os.environ["KF_PIPELINES_SA_TOKEN_PATH"] = token_file.name

    yield token_file.name

    if old_token_path is None:
        del os.environ["KF_PIPELINES_SA_TOKEN_PATH"]
    else:
        os.environ["KF_PIPELINES_SA_TOKEN_PATH"] = old_token_path
    os.remove(token_file.name)


@pytest.fixture
def get_model_file():
    with tempfile.NamedTemporaryFile(delete=False, suffix=".onnx") as model_file:
        pass

    yield model_file.name

    os.remove(model_file.name)


@pytest.fixture
def get_temp_dir_with_models():
    temp_dir = tempfile.mkdtemp()
    file_paths = []
    for _ in range(3):
        tmp_file = tempfile.NamedTemporaryFile(  # noqa: SIM115
            delete=False, dir=temp_dir, suffix=".onnx"
        )
        file_paths.append(tmp_file.name)
        tmp_file.close()

    yield temp_dir, file_paths

    for file in file_paths:
        if os.path.exists(file):
            os.remove(file)
    os.rmdir(temp_dir)


@pytest.fixture
def get_temp_dir():
    temp_dir = tempfile.mkdtemp()

    yield temp_dir

    shutil.rmtree(temp_dir)


@pytest.fixture
def get_temp_dir_with_nested_models():
    temp_dir = tempfile.mkdtemp()
    nested_dir = tempfile.mkdtemp(dir=temp_dir)

    file_paths = []
    for _ in range(3):
        tmp_file = tempfile.NamedTemporaryFile(  # noqa: SIM115
            delete=False, dir=nested_dir, suffix=".onnx"
        )
        file_paths.append(tmp_file.name)
        tmp_file.close()

    yield temp_dir, file_paths

    for file in file_paths:
        if os.path.exists(file):
            os.remove(file)
    os.rmdir(nested_dir)
    os.rmdir(temp_dir)


@pytest.fixture
def get_large_model_dir():
    """Creates a directory containing a large model file (300-500MB) for testing file size extremes."""
    temp_dir = tempfile.mkdtemp()

    # Create the large model file
    model_file = os.path.join(temp_dir, "large_model.onnx")
    with open(model_file, "wb") as f:
        # Write random data in chunks to create a large file
        chunk_size = 1024 * 1024  # 1MB chunks
        target_size = 400 * 1024 * 1024  # 400MB target size
        bytes_written = 0

        while bytes_written < target_size:
            # Generate random bytes for this chunk
            chunk = os.urandom(min(chunk_size, target_size - bytes_written))
            f.write(chunk)
            bytes_written += len(chunk)

    yield temp_dir

    # Cleanup
    shutil.rmtree(temp_dir)


@pytest.fixture
def patch_s3_env(monkeypatch: pytest.MonkeyPatch):
    s3_endpoint = os.getenv("KF_MR_TEST_S3_ENDPOINT")
    access_key_id = os.getenv("KF_MR_TEST_ACCESS_KEY_ID")
    secret_access_key = os.getenv("KF_MR_TEST_SECRET_ACCESS_KEY")
    bucket = os.getenv("KF_MR_TEST_BUCKET_NAME") or "default"
    region = "east"

    monkeypatch.setenv("AWS_S3_ENDPOINT", s3_endpoint)
    monkeypatch.setenv("AWS_ACCESS_KEY_ID", access_key_id)
    monkeypatch.setenv("AWS_SECRET_ACCESS_KEY", secret_access_key)
    monkeypatch.setenv("AWS_DEFAULT_REGION", region)
    monkeypatch.setenv("AWS_S3_BUCKET", bucket)

    return (bucket, s3_endpoint, access_key_id, secret_access_key, region)


# These are trimmed down versions of whats found in the example specs found here: https://github.com/opencontainers/image-spec/blob/main/image-layout.md#oci-layout-file
index_json_contents = """{
  "schemaVersion": 2,
  "mediaType": "application/vnd.oci.image.index.v1+json",
  "manifests": [],
  "annotations": {
    "com.example.index.revision": "r124356"
  }
}"""
oci_layout_contents = """{"imageLayoutVersion": "1.0.0"}"""


@pytest.fixture
def get_mock_custom_oci_backend():
    is_available_mock = Mock()
    is_available_mock.return_value = True
    pull_mock = Mock()
    push_mock = Mock()

    def pull_mock_imple(base_image, dest_dir, **kwargs):
        pathlib.Path(dest_dir).joinpath("oci-layout").write_text(oci_layout_contents)
        pathlib.Path(dest_dir).joinpath("index.json").write_text(index_json_contents)

    pull_mock.side_effect = pull_mock_imple
    return BackendDefinition(
        is_available=is_available_mock, pull=pull_mock, push=push_mock
    )


@pytest.fixture
def get_mock_skopeo_backend_for_auth(monkeypatch):
    user_auth = b"myuser:passwordhere"
    upc_encoded = base64.b64encode(user_auth)
    auth_data = {
        "auths": {"localhost:5001": {"auth": upc_encoded.decode(), "email": ""}}
    }
    auth_json = json.dumps(auth_data)
    monkeypatch.setenv(".dockerconfigjson", auth_json)
    generic_auth_vars = ["--username", "myuser", "--password", "mypasswordhere"]

    with (
        patch("olot.backend.skopeo.skopeo_pull") as skopeo_pull_mock,
        patch("olot.backend.skopeo.skopeo_push") as skopeo_push_mock,
        patch("olot.basics.oci_layers_on_top"),
    ):
        backend = _get_skopeo_backend(
            pull_args=generic_auth_vars, push_args=generic_auth_vars
        )

        def mock_override(base_image, dest_dir, params):
            return params

        skopeo_pull_mock.side_effect = mock_override
        skopeo_push_mock.side_effect = mock_override
        yield backend, skopeo_pull_mock, skopeo_push_mock, generic_auth_vars

@pytest.fixture(scope="session")
def generated_schema(pytestconfig: pytest.Config ) -> BaseOpenAPISchema:
    """Generate the schema for the API"""

    os.environ["API_HOST"] = REGISTRY_URL
    config = schemathesis.config.SchemathesisConfig.from_path(f"{pytestconfig.rootpath}/schemathesis.toml")
    local_schema_path = f"{pytestconfig.rootpath}/../../api/openapi/model-registry.yaml"
    schema = schemathesis.openapi.from_path(
        path=local_schema_path,
        config=config,
    )
    schema.config.output.sanitization.update(enabled=False)
    return schema

@pytest.fixture
def cleanup_artifacts(request: pytest.FixtureRequest, auth_headers: dict):
    """Cleanup artifacts created during the test."""
    created_ids = []
    def register(artifact_id):
        created_ids.append(artifact_id)

    yield register

    for artifact_id in created_ids:
        del_url = f"{REGISTRY_URL}/api/model_registry/v1alpha3/artifacts/{artifact_id}"
        try:
            requests.delete(del_url, headers=auth_headers, timeout=DEFAULT_API_TIMEOUT)
        except Exception as e:
            print(f"Failed to delete artifact {artifact_id}: {e}")

@pytest.fixture
def artifact_resource():
    """Create an artifact resource for the test."""
    @contextlib.contextmanager
    def _artifact_resource(auth_headers: dict, payload: dict) -> Generator[str, None, None]:
        create_endpoint = f"{REGISTRY_URL}/api/model_registry/v1alpha3/artifacts"
        resp = requests.post(create_endpoint, headers=auth_headers, json=payload, timeout=DEFAULT_API_TIMEOUT)
        resp.raise_for_status()
        artifact_id = resp.json()["id"]
        try:
            yield artifact_id
        finally:
            del_url = f"{REGISTRY_URL}/api/model_registry/v1alpha3/artifacts/{artifact_id}"
            try:
                requests.delete(del_url, headers=auth_headers, timeout=DEFAULT_API_TIMEOUT)
            except Exception as e:
                print(f"Failed to delete artifact {artifact_id}: {e}")
    return _artifact_resource

@pytest.fixture
def auth_headers(setup_env_user_token):
    """Provides authorization headers for API requests."""
    return {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {setup_env_user_token}"
    }

@pytest.fixture
def state_machine(generated_schema: BaseOpenAPISchema, auth_headers: str) -> APIStateMachine:
    BaseAPIWorkflow = generated_schema.as_state_machine()

    class APIWorkflow(BaseAPIWorkflow):  # type: ignore
        headers: dict[str, str]

        def setup(self) -> None:
            print("Cleaning up database")
            subprocess.run(
                ["../../scripts/cleanup.sh"],
                capture_output=True,
                check=True
            )
            self.headers = auth_headers

        def before_call(self, case: Case) -> None:
            print(f"Checking: {case.method} {case.path}")
        def get_call_kwargs(self, case: Case) -> dict[str, Any]:
            return {"verify": False, "headers": self.headers}

        def after_call(self, response: Response, case: Case) -> None:
            print(f"{case.method} {case.path} -> {response.status_code},")
    return APIWorkflow
