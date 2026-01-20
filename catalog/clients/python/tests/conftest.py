"""Pytest configuration and fixtures for catalog tests.

This module follows the model-registry Python client pattern:
- Assumes catalog service is already running (K8s, local, etc.)
- Uses environment variables for configuration
"""

import logging
import os
import time
from collections.abc import Generator
from pathlib import Path

import yaml
import pytest
import requests

from model_catalog import CatalogAPIClient

from .constants import (
    API_BASE_PATH,
    CATALOG_URL,
    CLIENT_TIMEOUT,
    MAX_BACKOFF,
    MAX_POLL_TIME,
    POLL_INTERVAL,
    get_verify_ssl,
)

# Configure logging
logging.basicConfig(
    format="%(asctime)s.%(msecs)03d - %(name)s:%(levelname)s: %(message)s",
    datefmt="%H:%M:%S",
    level=logging.WARNING,
)

logger = logging.getLogger("model-catalog")


def pytest_addoption(parser):
    """Add custom command line options."""
    parser.addoption("--e2e", action="store_true", help="run end-to-end tests")
    parser.addoption("--fuzz", action="store_true", help="run fuzzing tests")


def pytest_configure(config):
    """Register custom markers."""
    config.addinivalue_line("markers", "e2e: mark test as end-to-end test")
    config.addinivalue_line("markers", "fuzz: mark test as fuzzing test")
    config.addinivalue_line("markers", "huggingface: mark test as requiring HuggingFace API")


def _auto_mark_test(item) -> None:
    """Auto-mark tests based on their location."""
    path = str(item.fspath)
    if "fuzz_api" in path:
        item.add_marker(pytest.mark.fuzz)
    elif "tests" in path:
        item.add_marker(pytest.mark.e2e)


def _apply_skip_markers(item, *, e2e: bool, fuzz: bool) -> None:
    """Apply skip markers based on CLI flags."""
    skip_e2e = pytest.mark.skip(reason="need --e2e option to run E2E tests")
    skip_fuzz = pytest.mark.skip(reason="need --fuzz option to run fuzzing tests")
    skip_other = pytest.mark.skip(reason="skipping non-selected tests")

    if e2e:
        if "e2e" not in item.keywords:
            item.add_marker(skip_other)
    elif fuzz:
        if "fuzz" not in item.keywords:
            item.add_marker(skip_other)
    else:
        # No flag specified - skip both e2e and fuzz tests
        if "e2e" in item.keywords:
            item.add_marker(skip_e2e)
        if "fuzz" in item.keywords:
            item.add_marker(skip_fuzz)


def pytest_collection_modifyitems(config, items):
    """Modify test collection based on markers and options."""
    e2e = config.getoption("--e2e")
    fuzz = config.getoption("--fuzz")

    for item in items:
        _auto_mark_test(item)
        _apply_skip_markers(item, e2e=e2e, fuzz=fuzz)


def pytest_report_teststatus(report, config):
    """Custom test status reporting."""
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


# Maximum directory levels to traverse when searching for repo root
_MAX_PARENT_LEVELS = 10


@pytest.fixture(scope="session")
def root(request) -> Path:
    """Get repository root directory.

    Navigates up from catalog/clients/python to find the repo root.
    The repo root is identified by the presence of a .git directory.

    Raises:
        RuntimeError: If the repository root cannot be found.
    """
    current = request.config.rootpath
    # Walk up looking for .git directory (repo root marker)
    for _ in range(_MAX_PARENT_LEVELS):
        if (current / ".git").exists():
            return current
        current = current.parent
    # Fail explicitly if repo root not found
    msg = (
        f"Could not find repository root (.git directory) starting from "
        f"{request.config.rootpath}. Searched {_MAX_PARENT_LEVELS} levels up."
    )
    raise RuntimeError(msg)


@pytest.fixture(scope="session")
def user_token() -> str | None:
    """Get user token from environment."""
    return os.getenv("AUTH_TOKEN")


@pytest.fixture(scope="session")
def request_headers(user_token: str | None) -> dict[str, str]:
    """Get request headers including authorization if token is set."""
    headers = {"Content-Type": "application/json"}
    if user_token:
        headers["Authorization"] = f"Bearer {user_token}"
    return headers


@pytest.fixture(scope="session")
def verify_ssl() -> bool:
    """Get SSL verification setting from environment."""
    return get_verify_ssl(logger)


def poll_for_ready(user_token: str | None, verify_ssl: bool) -> None:
    """Wait for catalog service to be ready using exponential backoff.

    Args:
        user_token: Optional auth token.
        verify_ssl: Whether to verify SSL certificates.
    """
    url = f"{CATALOG_URL}{API_BASE_PATH}/sources"
    headers = {"Authorization": f"Bearer {user_token}"} if user_token else None

    # Exponential backoff: start at POLL_INTERVAL, double each time, cap at MAX_BACKOFF
    backoff = POLL_INTERVAL
    poll_start = time.time()

    while True:
        elapsed_time = time.time() - poll_start
        if elapsed_time >= MAX_POLL_TIME:
            msg = f"Catalog service not ready after {int(elapsed_time)}s at {url}"
            logger.error(msg)
            raise TimeoutError(msg)
        logger.info("Attempting to connect to server %s", url)
        try:
            response = requests.get(url, headers=headers, verify=verify_ssl, timeout=MAX_BACKOFF)
            if response.status_code < 500:  # Accept any non-5xx response
                logger.info("Server is up!")
                return
        except requests.exceptions.ConnectionError:
            pass

        time.sleep(backoff)
        backoff = min(backoff * 2, MAX_BACKOFF)  # Exponential backoff with cap


@pytest.fixture(scope="session")
def api_client(user_token: str | None, verify_ssl: bool) -> Generator[CatalogAPIClient, None, None]:
    """Create API client for the catalog service.

    This is a session-scoped fixture that connects to the already-running
    catalog service specified by CATALOG_URL environment variable.

    Timeout is configurable via CATALOG_CLIENT_TIMEOUT env var (default 30s).
    """
    poll_for_ready(user_token=user_token, verify_ssl=verify_ssl)
    with CatalogAPIClient(CATALOG_URL, timeout=CLIENT_TIMEOUT, verify_ssl=verify_ssl) as client:
        yield client


@pytest.fixture(scope="session")
def model_with_artifacts(api_client: CatalogAPIClient) -> tuple[str, str]:
    """Get a model that has artifacts for testing.

    Searches available models to find one with artifacts.
    Fails if no models or no models with artifacts are found.

    Returns:
        Tuple of (source_id, model_name) for a model with artifacts.

    Raises:
        pytest.fail: If no models are available or no model has artifacts.
    """
    models = api_client.get_models()
    if not models.get("items"):
        pytest.fail("No models available - test data may not be loaded")

    # Find a model that has artifacts
    for model in models["items"]:
        source_id = model.get("source_id")
        model_name = model.get("name")
        if not source_id or not model_name:
            continue

        # Check if this model has artifacts
        artifacts = api_client.get_artifacts(source_id=source_id, model_name=model_name)
        if artifacts.get("items"):
            return source_id, model_name

    # Fallback to first model with required fields
    model = models["items"][0]
    source_id = model.get("source_id")
    model_name = model.get("name")

    if not source_id or not model_name:
        pytest.fail("Model missing source_id or name - test data may be malformed")

    return source_id, model_name


@pytest.fixture(scope="session")
def testdata_dir(root) -> Path:
    """Get path to testdata directory."""
    return root / "test" / "testdata"


@pytest.fixture(scope="session")
def local_testdata_dir() -> Path:
    """Get path to local testdata directory (in tests/)."""
    return Path(__file__).parent / "testdata"


@pytest.fixture(scope="session")
def test_catalog_data(root: Path) -> dict:
    """Load test catalog data used by E2E tests.

    Returns:
        Dictionary containing the test catalog YAML data.
    """
    test_catalog_path = (
        root
        / "manifests"
        / "kustomize"
        / "options"
        / "catalog"
        / "overlays"
        / "e2e"
        / "test-catalog.yaml"
    )
    with open(test_catalog_path) as f:
        return yaml.safe_load(f)
