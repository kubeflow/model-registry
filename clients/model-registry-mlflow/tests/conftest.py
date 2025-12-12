"""
Shared fixtures for model_registry_mlflow tests.
"""

import uuid
from pathlib import Path

import pytest


def get_testdata_dir() -> Path:
    """Get the testdata directory, creating it if it doesn't exist."""
    testdata_dir = Path(__file__).parent / "testdata"
    testdata_dir.mkdir(exist_ok=True)
    return testdata_dir


def create_temp_file(suffix: str = "", prefix: str = "test_") -> str:
    """Create a temporary file in the testdata directory."""
    testdata_dir = get_testdata_dir()
    temp_filename = f"{prefix}{uuid.uuid4().hex[:8]}{suffix}"
    temp_path = testdata_dir / temp_filename
    return str(temp_path)


def create_temp_dir(prefix: str = "test_") -> Path:
    """Create a temporary directory in the testdata directory."""
    testdata_dir = get_testdata_dir()
    temp_dirname = f"{prefix}{uuid.uuid4().hex[:8]}"
    temp_path = testdata_dir / temp_dirname
    temp_path.mkdir(parents=True, exist_ok=True)
    return temp_path


@pytest.fixture(scope="session", autouse=True)
def setup_testdata_directory():
    """Ensure testdata directory exists and is cleaned up after tests."""
    testdata_dir = get_testdata_dir()
    yield testdata_dir
    # Cleanup after all tests are done
    import shutil

    if testdata_dir.exists():
        shutil.rmtree(testdata_dir, ignore_errors=True)
