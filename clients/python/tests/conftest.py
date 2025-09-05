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

import pytest
import requests

from model_registry import ModelRegistry
from model_registry.utils import BackendDefinition, _get_skopeo_backend

from .constants import (
    MAX_POLL_TIME,
    POLL_INTERVAL,
    REGISTRY_HOST,
    REGISTRY_PORT,
    REGISTRY_URL,
)


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
    if config.getoption("--quiet", default=False):
        return

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


start_time = time.time()


@pytest.fixture(scope="session")
def root(request) -> Path:
    return (request.config.rootpath / "../..").resolve()  # resolves to absolute path


@pytest.fixture(scope="session")
def user_token() -> str:
    return os.getenv("AUTH_TOKEN", None)


@pytest.fixture(scope="session")
def request_headers(user_token: str) -> dict[str, str]:
    headers = {"Content-Type": "application/json"}
    if user_token:
        headers["Authorization"] = f"Bearer {user_token}"
    return headers


@pytest.fixture(scope="session")
def verify_ssl() -> bool:
    verify_ssl_env = os.environ.get("VERIFY_SSL")
    if verify_ssl_env is None:
        return None
    return verify_ssl_env.lower() == "true"


def poll_for_ready(user_token, verify_ssl):
    params = {
        "url": REGISTRY_URL,
        "headers": {"Authorization": f"Bearer {user_token}", } if user_token else None,
        "verify": verify_ssl
    }
    while True:
        elapsed_time = time.time() - start_time
        if elapsed_time >= MAX_POLL_TIME:
            print("Polling timed out.")
            break
        print(f"Attempt to connect to server {REGISTRY_URL}")
        try:
            response = requests.get(**params, timeout=MAX_POLL_TIME)
            if response.status_code == 404:
                print("Server is up!")
                break
        except requests.exceptions.ConnectionError:
            pass

        # Wait for the specified poll interval before trying again
        time.sleep(POLL_INTERVAL)


def cleanup(fixture_func):
    async def yield_and_restart(root, request):
        # Access fixture values through request
        try:
            user_token = request.getfixturevalue("user_token")
            verify_ssl = request.getfixturevalue("verify_ssl")
        except pytest.FixtureLookupError:
            user_token = None
            verify_ssl = None

        poll_for_ready(user_token=user_token, verify_ssl=verify_ssl)

        if inspect.iscoroutinefunction(fixture_func) or inspect.isasyncgenfunction(fixture_func):
            async with asynccontextmanager(fixture_func)(user_token=user_token, verify_ssl=verify_ssl) as async_client:
                yield async_client
        else:
            # Check if fixture function expects parameters
            sig = inspect.signature(fixture_func)
            if "user_token" in sig.parameters:
                yield fixture_func(user_token=user_token)
            else:
                # For fixtures that don't take parameters (like client_attrs)
                yield fixture_func()

        print("Cleaning DB...")
        subprocess.call(  # noqa: S602
            "./scripts/cleanup.sh",
            shell=True,
            cwd=root,
        )

    return yield_and_restart


@pytest.fixture
@cleanup
def client(user_token: str) -> ModelRegistry:
    return ModelRegistry(REGISTRY_HOST, REGISTRY_PORT, author="author", is_secure=False, user_token=user_token)


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
