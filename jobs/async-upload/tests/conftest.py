import logging
import pytest

logging.basicConfig(level=logging.INFO)

def pytest_collection_modifyitems(config, items):
    for item in items: 
        skip_e2e = pytest.mark.skip(
            reason="this is an end-to-end test, requires explicit opt-in --e2e option to run."
        )
        skip_integration = pytest.mark.skip(
            reason="this is an integration test, requires explicit opt-in --integration option to run."
        )
        skip_not_e2e = pytest.mark.skip(
            reason="skipping non-e2e tests; opt-out of --e2e -like options to run."
        )
        skip_not_integration = pytest.mark.skip(
            reason="skipping non-integration tests; opt-out of --integration -like options to run."
        )
        
        e2e_option = config.getoption("--e2e")
        integration_option = config.getoption("--integration")
        
        if "e2e" in item.keywords:
            if not e2e_option:
                item.add_marker(skip_e2e)
        elif "integration" in item.keywords:
            if not integration_option:
                item.add_marker(skip_integration)
        elif e2e_option:
            item.add_marker(skip_not_e2e)
        elif integration_option:
            item.add_marker(skip_not_integration)


def pytest_addoption(parser):
    parser.addoption(
        "--e2e",
        action="store_true",
        default=False,
        help="opt-in to run tests marked with e2e",
    )
    parser.addoption(
        "--integration",
        action="store_true",
        default=False,
        help="opt-in to run tests marked with integration",
    )
