import pytest
from model_registry import ModelRegistry
from model_registry.core import ModelRegistryAPIClient
from model_registry.exceptions import StoreException


@pytest.fixture()
def mr_client(mr_api: ModelRegistryAPIClient) -> ModelRegistry:
    mr = ModelRegistry.__new__(ModelRegistry)
    mr._api = mr_api
    mr._author = "test_author"
    return mr


def test_register_new(mr_client: ModelRegistry):
    name = "test_model"
    version = "1.0.0"
    rm = mr_client.register_model(
        name,
        "s3",
        model_format_name="test_format",
        model_format_version="test_version",
        version=version,
    )
    assert rm.id is not None

    mr_api = mr_client._api
    assert (mv := mr_api.get_model_version_by_params(rm.id, version)) is not None
    assert mr_api.get_model_artifact_by_params(mv.id) is not None


def test_register_existing_version(mr_client: ModelRegistry):
    params = {
        "name": "test_model",
        "uri": "s3",
        "model_format_name": "test_format",
        "model_format_version": "test_version",
        "version": "1.0.0",
    }
    mr_client.register_model(**params)

    with pytest.raises(StoreException):
        mr_client.register_model(**params)


def test_get(mr_client: ModelRegistry):
    name = "test_model"
    version = "1.0.0"

    rm = mr_client.register_model(
        name,
        "s3",
        model_format_name="test_format",
        model_format_version="test_version",
        version=version,
    )

    assert rm.id == mr_client.get_registered_model(name).id

    mr_api = mr_client._api
    mv = mr_api.get_model_version_by_params(rm.id, version)
    ma = mr_api.get_model_artifact_by_params(mv.id)

    assert mv.id == mr_client.get_model_version(name, version).id
    assert ma.id == mr_client.get_model_artifact(name, version).id
