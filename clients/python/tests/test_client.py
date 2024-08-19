import os
from itertools import islice

import pytest

from model_registry import ModelRegistry, utils
from model_registry.exceptions import StoreError

from .conftest import REGISTRY_HOST, REGISTRY_PORT, cleanup


@pytest.fixture
@cleanup
def client() -> ModelRegistry:
    return ModelRegistry(REGISTRY_HOST, REGISTRY_PORT, author="author", is_secure=False)


def test_secure_client():
    os.environ["CERT"] = ""
    os.environ["KF_PIPELINES_SA_TOKEN_PATH"] = ""
    with pytest.raises(StoreError) as e:
        ModelRegistry("anything", author="test_author")

    assert "user token" in str(e.value).lower()


async def test_register_new(client: ModelRegistry):
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
    mv = await mr_api.get_model_version_by_params(rm.id, version)
    assert mv
    assert mv.id
    ma = await mr_api.get_model_artifact_by_params(name, mv.id)
    assert ma


async def test_register_new_using_s3_uri_builder(client: ModelRegistry):
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
    mv = await mr_api.get_model_version_by_params(rm.id, version)
    assert mv
    assert mv.id
    ma = await mr_api.get_model_artifact_by_params(name, mv.id)
    assert ma
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

    with pytest.raises(StoreError):
        client.register_model(**params)


async def test_get(client: ModelRegistry):
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

    mr_api = client._api
    assert (mv := await mr_api.get_model_version_by_params(rm.id, version))
    assert mv.id
    assert (ma := await mr_api.get_model_artifact_by_params(name, mv.id))

    assert (_mv := client.get_model_version(name, version))
    assert mv.id == _mv.id
    assert mv.custom_properties == metadata
    assert (_ma := client.get_model_artifact(name, version))
    assert ma.id == _ma.id


def test_get_registered_models(client: ModelRegistry):
    models = 21

    for name in [f"test_model{i}" for i in range(models)]:
        client.register_model(
            name,
            "s3",
            model_format_name="test_format",
            model_format_version="test_version",
            version="1.0.0",
        )

    rm_iter = client.get_registered_models().limit(10)
    i = 0
    prev_tok = None
    changes = 0
    with pytest.raises(StopIteration):  # noqa: PT012
        while i < 50 and next(rm_iter):
            if rm_iter.options.next_page_token != prev_tok:
                print(
                    f"Token changed from {prev_tok} to {rm_iter.options.next_page_token} at {i}"
                )
                prev_tok = rm_iter.options.next_page_token
                changes += 1
            i += 1

    assert changes == 3
    assert i == models


def test_get_registered_models_and_reset(client: ModelRegistry):
    model_count = 6
    page = model_count // 2

    for name in [f"test_model{i}" for i in range(model_count)]:
        client.register_model(
            name,
            "s3",
            model_format_name="test_format",
            model_format_version="test_version",
            version="1.0.0",
        )

    rm_iter = client.get_registered_models().limit(model_count - 1)
    models = []
    for rm in islice(rm_iter, page):
        models.append(rm)
    assert len(models) == page
    rm_iter.restart()
    complete = list(rm_iter)
    assert len(complete) == model_count
    assert complete[:page] == models


def test_get_model_versions(client: ModelRegistry):
    name = "test_model"
    models = 21

    for v in [f"1.0.{i}" for i in range(models)]:
        client.register_model(
            name,
            "s3",
            model_format_name="test_format",
            model_format_version="test_version",
            version=v,
        )

    mv_iter = client.get_model_versions(name).limit(10)
    i = 0
    prev_tok = None
    changes = 0
    with pytest.raises(StopIteration):  # noqa: PT012
        while i < 50 and next(mv_iter):
            if mv_iter.options.next_page_token != prev_tok:
                print(
                    f"Token changed from {prev_tok} to {mv_iter.options.next_page_token} at {i}"
                )
                prev_tok = mv_iter.options.next_page_token
                changes += 1
            i += 1

    assert changes == 3
    assert i == models


def test_get_model_versions_and_reset(client: ModelRegistry):
    name = "test_model"

    model_count = 6
    page = model_count // 2

    for v in [f"1.0.{i}" for i in range(model_count)]:
        client.register_model(
            name,
            "s3",
            model_format_name="test_format",
            model_format_version="test_version",
            version=v,
        )

    mv_iter = client.get_model_versions(name).limit(model_count - 1)
    models = []
    for rm in islice(mv_iter, page):
        models.append(rm)
    assert len(models) == page
    mv_iter.restart()
    complete = list(mv_iter)
    assert len(complete) == model_count
    assert complete[:page] == models


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
