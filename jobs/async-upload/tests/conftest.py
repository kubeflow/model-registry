import logging
import pytest

logging.basicConfig(level=logging.INFO)

def pytest_collection_modifyitems(config, items):
    for item in items: 
        skip_e2e = pytest.mark.skip(
            reason="this is an end-to-end test, requires explicit opt-in --e2e option to run."
        )
        skip_not_e2e = pytest.mark.skip(
            reason="skipping non-e2e tests; opt-out of --e2e -like options to run."
        )
        if "e2e" in item.keywords:
            if not config.getoption("--e2e"):
                item.add_marker(skip_e2e)
        elif config.getoption("--e2e"):
            item.add_marker(skip_not_e2e)


def pytest_addoption(parser):
    parser.addoption(
        "--e2e",
        action="store_true",
        default=False,
        help="opt-in to run tests marked with e2e",
    )
