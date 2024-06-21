import os

import pytest
from model_registry import ModelRegistry, utils
from model_registry.exceptions import StoreException


@pytest.fixture()
def client() -> ModelRegistry:
    pass


def test_secure_client():
    os.environ["CERT"] = ""
    os.environ["KF_PIPELINES_SA_TOKEN_PATH"] = ""
    with pytest.raises(StoreException) as e:
        ModelRegistry("anything", author="test_author")

    assert "user token" in str(e.value).lower()


def test_register_new(client: ModelRegistry):
    name = "test_model"
    version = "1.0.0"
    rm = client.register_model(
        name,
        "s3",
        model_format_name="test_format",
        model_format_version="test_version",
        version=version,
    )
    assert rm.id

    mr_api = client._api
    mv = mr_api.get_model_version_by_params(rm.id, version)
    assert mv
    assert mv.id
    ma = mr_api.get_model_artifact_by_params(name, mv.id)
    assert ma


def test_register_new_using_s3_uri_builder(client: ModelRegistry):
    name = "test_model"
    version = "1.0.0"
    uri = utils.s3_uri_from(
        "storage/path", "my-bucket", endpoint="my-endpoint", region="my-region"
    )
    rm = client.register_model(
        name,
        uri,
        model_format_name="test_format",
        model_format_version="test_version",
        version=version,
    )
    assert rm.id is not None

    mr_api = client._api
    assert (mv := mr_api.get_model_version_by_params(rm.id, version))
    assert mv.id
    assert (ma := mr_api.get_model_artifact_by_params(mv.id))
    assert ma.uri == uri


def test_register_existing_version(client: ModelRegistry):
    params = {
        "name": "test_model",
        "uri": "s3",
        "model_format_name": "test_format",
        "model_format_version": "test_version",
        "version": "1.0.0",
    }
    client.register_model(**params)

    with pytest.raises(StoreException):
        client.register_model(**params)


def test_get(client: ModelRegistry):
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

    assert (_rm := client.get_registered_model(name))
    assert rm.id == _rm.id

    mr_api = client._api
    assert (mv := mr_api.get_model_version_by_params(rm.id, version))
    assert mv.id
    assert (ma := mr_api.get_model_artifact_by_params(name, mv.id))

    assert (_mv := client.get_model_version(name, version))
    assert mv.id == _mv.id
    assert mv.custom_properties == metadata
    assert (_ma := client.get_model_artifact(name, version))
    assert ma.id == _ma.id


def test_hf_import(client: ModelRegistry):
    pytest.importorskip("huggingface_hub")
    name = "openai-community/gpt2"
    version = "1.2.3"
    author = "test author"

    assert client.register_hf_model(
        name,
        "onnx/decoder_model.onnx",
        author=author,
        version=version,
        model_format_name="test format",
        model_format_version="test version",
    )
    assert (mv := client.get_model_version(name, version))
    assert mv.author == author
    assert mv.custom_properties
    assert mv.custom_properties["model_author"] == author
    assert mv.custom_properties["model_origin"] == "huggingface_hub"
    assert (
        mv.custom_properties["source_uri"]
        == "https://huggingface.co/openai-community/gpt2/resolve/main/onnx/decoder_model.onnx"
    )
    assert mv.custom_properties["repo"] == name
    assert client.get_model_artifact(name, version)


def test_hf_import_default_env(client: ModelRegistry):
    """Test setting environment variables, hence triggering defaults, does _not_ interfere with HF metadata"""
    pytest.importorskip("huggingface_hub")
    name = "openai-community/gpt2"
    version = "1.2.3"
    author = "test author"
    env_values = {
        "AWS_S3_ENDPOINT": "value1",
        "AWS_S3_BUCKET": "value2",
        "AWS_DEFAULT_REGION": "value3",
    }
    for k, v in env_values.items():
        os.environ[k] = v

    assert client.register_hf_model(
        name,
        "onnx/decoder_model.onnx",
        author=author,
        version=version,
        model_format_name="test format",
        model_format_version="test version",
    )
    assert (mv := client.get_model_version(name, version))
    assert mv.custom_properties
    assert mv.custom_properties["model_author"] == author
    assert mv.custom_properties["model_origin"] == "huggingface_hub"
    assert (
        mv.custom_properties["source_uri"]
        == "https://huggingface.co/openai-community/gpt2/resolve/main/onnx/decoder_model.onnx"
    )
    assert mv.custom_properties["repo"] == name
    assert client.get_model_artifact(name, version)

    for k in env_values:
        os.environ.pop(k)
