import pytest

from model_registry import ModelRegistry


@pytest.mark.e2e
def test_create_tagged_version(client: ModelRegistry):
    """Test regression for creating tagged versions.

    Reported on: https://github.com/kubeflow/model-registry/issues/255
    """
    name = "test_model"
    version = "model:latest"
    rm = client.register_model(
        name,
        "s3",
        model_format_name="test_format",
        model_format_version="test_version",
        version=version,
    )
    assert rm.id
    mv = client.get_model_version(name, version)
    assert mv
    assert mv.id

@pytest.mark.e2e
def test_get_model_without_user_token(setup_env_user_token, client):
    """Test regression for using client methods without an user_token in the init arguments.

    Reported on: https://github.com/kubeflow/model-registry/issues/340
    """
    assert setup_env_user_token != ""
    name = "test_model"
    version = "1.0.0"
    metadata = {"a": 1, "b": "2"}
    rm = client.register_model(
        name,
        "s3",
        model_format_name="test_format",
        model_format_version="test_version",
        version=version,
        metadata=metadata,
    )
    assert rm.id
    assert (_rm := client.get_registered_model(name))
    assert rm.id == _rm.id
