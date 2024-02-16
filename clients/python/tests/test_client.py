import os

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
    metadata = {"a": 1, "b": "2"}

    rm = mr_client.register_model(
        name,
        "s3",
        model_format_name="test_format",
        model_format_version="test_version",
        version=version,
        metadata=metadata
    )

    assert (_rm := mr_client.get_registered_model(name))
    assert rm.id == _rm.id

    mr_api = mr_client._api
    assert (mv := mr_api.get_model_version_by_params(rm.id, version))
    assert (ma := mr_api.get_model_artifact_by_params(mv.id))

    assert (_mv := mr_client.get_model_version(name, version))
    assert mv.id == _mv.id
    assert mv.metadata == metadata
    assert (_ma := mr_client.get_model_artifact(name, version))
    assert ma.id == _ma.id


def test_default_md(mr_client: ModelRegistry):
    name = "test_model"
    version = "1.0.0"
    env_values = {"AWS_S3_ENDPOINT": "value1", "AWS_S3_BUCKET": "value2", "AWS_DEFAULT_REGION": "value3"}
    for k, v in env_values.items():
        os.environ[k] = v

    assert mr_client.register_model(
        name,
        "s3",
        model_format_name="test_format",
        model_format_version="test_version",
        version=version,
        # ensure leave empty metadata
    )
    assert (mv := mr_client.get_model_version(name, version))
    assert mv.metadata == env_values

    for k in env_values:
        os.environ.pop(k)


def test_hf_import(mr_client: ModelRegistry):
    pytest.importorskip("huggingface_hub")
    name = "openai-community/gpt2"
    version = "1.2.3"
    author = "test author"

    assert mr_client.register_hf_model(
        name,
        "onnx/decoder_model.onnx",
        author=author,
        version=version,
        model_format_name="test format",
        model_format_version="test version",
    )
    assert (mv := mr_client.get_model_version(name, version))
    assert mv.author == author
    assert mv.metadata["model_author"] == author
    assert mv.metadata["model_origin"] == "huggingface_hub"
    assert mv.metadata["source_uri"] == "https://huggingface.co/openai-community/gpt2/resolve/main/onnx/decoder_model.onnx"
    assert mv.metadata["repo"] == name
    assert mr_client.get_model_artifact(name, version)


def test_hf_import_default_env(mr_client: ModelRegistry):
    """Test setting environment variables, hence triggering defaults, does _not_ interfere with HF metadata
    """
    pytest.importorskip("huggingface_hub")
    name = "openai-community/gpt2"
    version = "1.2.3"
    author = "test author"
    env_values = {"AWS_S3_ENDPOINT": "value1", "AWS_S3_BUCKET": "value2", "AWS_DEFAULT_REGION": "value3"}
    for k, v in env_values.items():
        os.environ[k] = v

    assert mr_client.register_hf_model(
        name,
        "onnx/decoder_model.onnx",
        author=author,
        version=version,
        model_format_name="test format",
        model_format_version="test version",
    )
    assert (mv := mr_client.get_model_version(name, version))
    assert mv.metadata["model_author"] == author
    assert mv.metadata["model_origin"] == "huggingface_hub"
    assert mv.metadata["source_uri"] == "https://huggingface.co/openai-community/gpt2/resolve/main/onnx/decoder_model.onnx"
    assert mv.metadata["repo"] == name
    assert mr_client.get_model_artifact(name, version)

    for k in env_values:
        os.environ.pop(k)
