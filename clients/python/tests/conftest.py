import os
from pathlib import Path

import pytest


@pytest.fixture(scope="session")
def root(request) -> Path:
    return (request.config.rootpath / "../..").resolve()  # resolves to absolute path


@pytest.fixture(scope="session")
def _compose_mr(root):
    print("Assuming this is the Model Registry root directory:", root)
    shared_volume = root / "test/config/ml-metadata"
    sqlite_db_file = shared_volume / "metadata.sqlite.db"
    if sqlite_db_file.exists():
        msg = f"The file {sqlite_db_file} already exists; make sure to cancel it before running these tests."
        raise FileExistsError(msg)

    yield

    try:
        os.remove(sqlite_db_file)
        print(f"Removed {sqlite_db_file} successfully.")
    except Exception as e:
        print(f"An error occurred while removing {sqlite_db_file}: {e}")
