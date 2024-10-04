import asyncio
import inspect
import os
import subprocess
import tempfile
import time
from contextlib import asynccontextmanager
from pathlib import Path
from urllib.parse import urlparse

import pytest
import requests

from model_registry import ModelRegistry


def pytest_addoption(parser):
    parser.addoption("--e2e", action="store_true", help="run end-to-end tests")


def pytest_collection_modifyitems(config, items):
    for item in items:
        skip_e2e = pytest.mark.skip(
            reason="this is an end-to-end test, requires explicit opt-in --e2e option to run."
        )
        if "e2e" in item.keywords:
            if not config.getoption("--e2e"):
                item.add_marker(skip_e2e)
            continue


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
@cleanup
def client() -> ModelRegistry:
    return ModelRegistry(REGISTRY_HOST, REGISTRY_PORT, author="author", is_secure=False)


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
