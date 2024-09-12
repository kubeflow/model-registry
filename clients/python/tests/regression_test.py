import pytest

from model_registry import ModelRegistry
from model_registry.types.artifacts import ModelArtifact


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


@pytest.mark.e2e
def test_get_few_registered_models(client: ModelRegistry):
    """Test regression for paging without next page token.

    Reported on: https://github.com/kubeflow/model-registry/issues/348
    """
    models = 9

    for name in [f"test_model{i}" for i in range(models)]:
        client.register_model(
            name,
            "s3",
            model_format_name="test_format",
            model_format_version="test_version",
            version="1.0.0",
        )

    i = 0
    for rm in client.get_registered_models():
        print(f"found {rm}")
        i += 1
        assert i < models + 1

    assert i == models


@pytest.mark.e2e
async def test_create_standalone_model_artifact(client: ModelRegistry):
    """Test regression for creating standalone model artifact.

    Reported on: https://github.com/kubeflow/model-registry/issues/231"""
    ma = ModelArtifact(uri="s3")
    async with client._api.get_client() as api:
        new_raw_ma = await api.create_model_artifact(ma.create())
        new_ma = ModelArtifact.from_basemodel(new_raw_ma)
        assert new_ma.id
        assert new_ma.uri == "s3"

    rm = client.register_model(
        "test_model",
        "s3",
        version="1.0.0",
        model_format_name="x",
        model_format_version="y",
    )
    assert rm.id
    mv = client.get_model_version("test_model", "1.0.0")
    assert mv
    assert mv.id
    mv_ma = await client._api.upsert_model_version_artifact(new_ma, mv.id)
    assert mv_ma.id == new_ma.id
