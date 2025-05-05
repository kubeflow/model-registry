import asyncio
import base64
import inspect
import json
import os
import pathlib
import shutil
import subprocess
import tempfile
import time
from contextlib import asynccontextmanager
from pathlib import Path
from unittest.mock import Mock, patch
from urllib.parse import urlparse

import pytest
import requests
import uvloop

from model_registry import ModelRegistry
from model_registry.utils import BackendDefinition, _get_skopeo_backend


def pytest_addoption(parser):
    parser.addoption("--e2e", action="store_true", help="run end-to-end tests")


def pytest_collection_modifyitems(config, items):
    for item in items:
        if "e2e" in item.keywords:
            skip_e2e = pytest.mark.skip(
                reason="this is an end-to-end test, requires explicit opt-in --e2e option to run."
            )
            if not config.getoption("--e2e"):
                item.add_marker(skip_e2e)
            continue


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


# workaround: https://github.com/pytest-dev/pytest-asyncio/issues/706#issuecomment-2147044022
@pytest.fixture(scope="session", autouse=True)
def event_loop():
    loop = asyncio.get_event_loop_policy().get_event_loop()
    yield loop
    loop.close()


@pytest.fixture
def uv_event_loop():
    old_policy = asyncio.get_event_loop_policy()
    policy = uvloop.EventLoopPolicy()
    asyncio.set_event_loop_policy(policy)
    loop = uvloop.new_event_loop()
    asyncio.set_event_loop(loop)
    yield loop
    loop.close()
    asyncio.set_event_loop_policy(old_policy)


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
