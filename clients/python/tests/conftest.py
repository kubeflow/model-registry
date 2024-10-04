import asyncio
import inspect
import os
import subprocess
import tempfile
import time
from contextlib import asynccontextmanager
from pathlib import Path
from time import sleep

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


REGISTRY_HOST = "http://localhost"
REGISTRY_PORT = 8080
REGISTRY_URL = f"{REGISTRY_HOST}:{REGISTRY_PORT}"
COMPOSE_FILE = "docker-compose.yaml"
MAX_POLL_TIME = 1200  # the first build is extremely slow if using docker-compose-*local*.yaml for bootstrap of builder image
POLL_INTERVAL = 1
DOCKER = os.getenv("DOCKER", "docker")
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


@pytest.fixture(scope="session")
def _compose_mr(root):
    print("Assuming this is the Model Registry root directory:", root)
    shared_volume = root / "test/config/ml-metadata"
    sqlite_db_file = shared_volume / "metadata.sqlite.db"
    if sqlite_db_file.exists():
        msg = f"The file {sqlite_db_file} already exists; make sure to cancel it before running these tests."
        raise FileExistsError(msg)
    print(f" Starting Docker Compose in folder {root}")
    p = subprocess.Popen(  # noqa: S602
        f"{DOCKER} compose -f {COMPOSE_FILE} up",
        shell=True,
        cwd=root,
    )
    yield

    p.kill()
    print(f" Closing Docker Compose in folder {root}")
    subprocess.call(  # noqa: S602
        f"{DOCKER} compose -f {COMPOSE_FILE} down",
        shell=True,
        cwd=root,
    )
    try:
        os.remove(sqlite_db_file)
        print(f"Removed {sqlite_db_file} successfully.")
    except Exception as e:
        print(f"An error occurred while removing {sqlite_db_file}: {e}")


def cleanup(client):
    async def yield_and_restart(_compose_mr, root):
        poll_for_ready()
        if inspect.iscoroutinefunction(client) or inspect.isasyncgenfunction(client):
            async with asynccontextmanager(client)() as async_client:
                yield async_client
        else:
            yield client()

        sqlite_db_file = root / "test/config/ml-metadata/metadata.sqlite.db"
        try:
            os.remove(sqlite_db_file)
            print(f"Removed {sqlite_db_file} successfully.")
        except Exception as e:
            print(f"An error occurred while removing {sqlite_db_file}: {e}")
        # we have to wait to make sure the server restarts after the file is gone
        sleep(1)

        print("Restarting model-registry...")
        subprocess.call(  # noqa: S602
            f"{DOCKER} compose -f {COMPOSE_FILE} restart model-registry",
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
