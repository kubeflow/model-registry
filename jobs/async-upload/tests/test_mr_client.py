import os
import aiohttp
from model_registry import ModelRegistry
from model_registry.exceptions import StoreError
import pytest
from job.mr_client import validate_and_get_model_registry_client
from job.config import get_config


@pytest.fixture
def minimal_env_source_dest_vars():
    original_env = dict(os.environ)

    vars = {
        "type": "oci",
        "oci_uri": "quay.io/example/oci",
        "oci_registry": "quay.io",
        "oci_username": "oci_username_env",
        "oci_password": "oci_password_env",
    }

    # Set up test environment variables
    for key, value in vars.items():
        os.environ[f"MODEL_SYNC_DESTINATION_{key.upper()}"] = value
    for key, value in vars.items():
        os.environ[f"MODEL_SYNC_SOURCE_{key.upper()}"] = value

    vars = {
        "model_upload_intent": "update_artifact",
        "model_artifact_id": "123",
        "registry_server_address": "http://localhost",
        "registry_port": "8080",
        "registry_author": "author",
    }

    for key, value in vars.items():
        os.environ[f"MODEL_SYNC_{key.upper()}"] = value

    yield vars

    # Restore original environment
    os.environ.clear()
    os.environ.update(original_env)


@pytest.mark.e2e
def test_model_registry_config_throws_error_on_missing_user_token(
    minimal_env_source_dest_vars,
):
    """Test that the model registry config throws an error on missing user token because it's a secure connection"""
    sample_config = get_config([])
    with pytest.raises(StoreError) as e:
        validate_and_get_model_registry_client(sample_config.registry)
    assert "user token must be provided for secure connection" in str(e.value)


@pytest.mark.e2e
def test_model_registry_config_correct(minimal_env_source_dest_vars):
    """Test that the model registry config is correct"""
    sample_config = get_config(["--registry-is-secure", False])
    # Note: Instantiating the client will ping the GET /.../registered_models endpoint, validating the connection
    client = validate_and_get_model_registry_client(sample_config.registry)
    assert isinstance(client, ModelRegistry)


@pytest.mark.e2e
def test_model_registry_config_throws_when_mr_is_unreachable(
    minimal_env_source_dest_vars,
):
    """Test that the model registry config is correct"""
    sample_config = get_config(
        [
            "--registry-is-secure",
            "false",
            "--registry-port",
            "1337",  # Note: the E2E test should expose 8080, this is a purposely invalid port
        ]
    )

    with pytest.raises(aiohttp.client_exceptions.ClientConnectorError) as e:
        validate_and_get_model_registry_client(sample_config.registry)
    assert "Cannot connect to host" in str(e.value)
