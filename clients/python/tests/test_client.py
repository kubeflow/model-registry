import os
from itertools import islice

import pytest
import requests

from model_registry import ModelRegistry, utils
from model_registry.exceptions import StoreError
from model_registry.types import ModelArtifact


def test_secure_client():
    os.environ["CERT"] = ""
    os.environ["KF_PIPELINES_SA_TOKEN_PATH"] = ""
    with pytest.raises(StoreError) as e:
        ModelRegistry("anything", author="test_author")

    assert "user token" in str(e.value).lower()


@pytest.mark.e2e
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


@pytest.mark.e2e
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


@pytest.mark.e2e
def test_register_existing_version(client: ModelRegistry):
    params = {
        "name": "test_model",
        "uri": "s3",
        "model_format_name": "test_format",
        "model_format_version": "test_version",
        "version": "1.0.0",
    }
    client.register_model(**params, metadata=None)

    with pytest.raises(StoreError):
        client.register_model(**params, metadata=None)


@pytest.mark.e2e
async def test_update_models(client: ModelRegistry):
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

    new_description = "updated description"
    rm.description = new_description
    mv.description = new_description
    ma.description = new_description
    assert client.update(rm).description == new_description
    assert client.update(mv).description == new_description
    assert client.update(ma).description == new_description


@pytest.mark.e2e
async def test_update_logical_model_with_labels(client: ModelRegistry):
    """As a MLOps engineer I would like to store some labels

    A custom property of type string, with empty string value, shall be considered a Label; this is also semantically compatible for properties having empty string values in general.
    """
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
    mv = client.get_model_version(name, version)
    assert mv
    assert mv.id
    ma = client.get_model_artifact(name, version)
    assert ma
    assert ma.id

    rm_labels = {
        "my-label1": "",
        "my-label2": "",
    }
    rm.custom_properties = rm_labels
    client.update(rm)

    mv_labels = {
        "my-label3": "",
        "my-label4": "",
    }
    mv.custom_properties = mv_labels
    client.update(mv)

    ma_labels = {
        "my-label5": "",
        "my-label6": "",
    }
    ma.custom_properties = ma_labels
    client.update(ma)

    rm = client.get_registered_model(name)
    assert rm
    assert rm.custom_properties == rm_labels
    mv = client.get_model_version(name, version)
    assert mv
    assert mv.custom_properties == mv_labels
    ma = client.get_model_artifact(name, version)
    assert ma
    assert ma.custom_properties == ma_labels


@pytest.mark.e2e
async def test_patch_model_artifacts_artifact_type(client: ModelRegistry):
    """Patching ModelArtifact requires `artifactType` value which was previously not required

    reported with https://issues.redhat.com/browse/RHOAIENG-15326
    """
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
    mv = client.get_model_version(name, version)
    assert mv
    assert mv.id
    ma = client.get_model_artifact(name, version)
    assert ma
    assert ma.id

    payload = { "modelFormatName": "foo" }
    from .conftest import REGISTRY_HOST, REGISTRY_PORT
    response = requests.patch(url=f"{REGISTRY_HOST}:{REGISTRY_PORT}/api/model_registry/v1alpha3/model_artifacts/{ma.id}", json=payload, timeout=10, headers={"Content-Type": "application/json"})
    assert response.status_code == 200
    ma = client.get_model_artifact(name, version)
    assert ma
    assert ma.id
    assert ma.model_format_name == "foo"


@pytest.mark.e2e
async def test_update_preserves_model_info(client: ModelRegistry):
    name = "test_model"
    version = "1.0.0"
    uri = "s3"
    model_fmt_name = "test_format"
    model_fmt_version = "test_version"
    rm = client.register_model(
        name,
        uri,
        model_format_name=model_fmt_name,
        model_format_version=model_fmt_version,
        version=version,
    )
    assert rm.id

    mr_api = client._api
    mv = await mr_api.get_model_version_by_params(rm.id, version)
    assert mv
    assert mv.id
    ma = await mr_api.get_model_artifact_by_params(name, mv.id)
    assert ma

    new_description = "updated description"
    ma = ModelArtifact(id=ma.id, uri=uri, description=new_description)

    updated_ma = client.update(ma)
    assert updated_ma.description == new_description
    assert updated_ma.uri == uri
    assert updated_ma.id == ma.id
    assert updated_ma.model_format_name == model_fmt_name
    assert updated_ma.model_format_version == model_fmt_version


@pytest.mark.e2e
async def test_update_existing_model_artifact(client: ModelRegistry):
    """Updating uri (or other properties) by re-using and call to update

    reported via slack
    """
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
    mv = client.get_model_version(name, version)
    assert mv
    assert mv.id
    ma = client.get_model_artifact(name, version)
    assert ma
    assert ma.id

    something_else = "https://something.else/model.onnx"
    ma.uri = something_else
    response = client.update(ma)
    assert response
    assert response.uri == something_else

    ma = client.get_model_artifact(name, version)
    assert ma.uri == something_else


@pytest.mark.e2e
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


@pytest.mark.e2e
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

    rm_iter = client.get_registered_models().page_size(10)
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


@pytest.mark.e2e
def test_get_registered_models_order_by(client: ModelRegistry):
    models = 5

    rms = []
    for name in [f"test_model{i}" for i in range(models)]:
        rms.append(
            client.register_model(
                name,
                "s3",
                model_format_name="test_format",
                model_format_version="test_version",
                version="1.0.0",
            )
        )

    # id ordering should match creation order
    i = 0
    for rm, by_id in zip(
        rms,
        client.get_registered_models().order_by_id(),
    ):
        assert rm.id == by_id.id
        i += 1

    assert i == models

    # and obviously, creation ordering should match creation ordering
    i = 0
    for rm, by_creation in zip(
        rms,
        client.get_registered_models().order_by_creation_time(),
    ):
        assert rm.id == by_creation.id
        i += 1

    assert i == models

    # update order should match creation ordering by default
    i = 0
    for rm, by_update in zip(
        rms,
        client.get_registered_models().order_by_update_time(),
    ):
        assert rm.id == by_update.id
        i += 1

    assert i == models

    # now update the models in reverse order
    for rm in reversed(rms):
        rm.description = "updated"
        client.update(rm)

    # and they should match in reverse
    i = 0
    for rm, by_update in zip(
        reversed(rms),
        client.get_registered_models().order_by_update_time(),
    ):
        assert rm.id == by_update.id
        i += 1

    assert i == models

    # or if descending is explicitly set
    i = 0
    for rm, by_update in zip(
        rms,
        client.get_registered_models().order_by_update_time().descending(),
    ):
        assert rm.id == by_update.id
        i += 1

    assert i == models


@pytest.mark.e2e
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

    rm_iter = client.get_registered_models().page_size(model_count - 1)
    models = []
    for rm in islice(rm_iter, page):
        models.append(rm)
    assert len(models) == page
    rm_iter.restart()
    complete = list(rm_iter)
    assert len(complete) == model_count
    assert complete[:page] == models


@pytest.mark.e2e
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

    mv_iter = client.get_model_versions(name).page_size(10)
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


@pytest.mark.e2e
@pytest.mark.xfail(
    reason="MLMD issue tracked on: https://github.com/kubeflow/model-registry/issues/358"
)
def test_get_model_versions_order_by(client: ModelRegistry):
    name = "test_model"
    models = 5
    mvs = []
    for v in [f"1.0.{i}" for i in range(models)]:
        client.register_model(
            name,
            "s3",
            model_format_name="test_format",
            model_format_version="test_version",
            version=v,
        )
        mvs.append(client.get_model_version(name, v))

    i = 0
    for mv, by_id in zip(
        mvs,
        client.get_model_versions(name).order_by_id(),
    ):
        assert mv.id == by_id.id
        i += 1

    assert i == models

    i = 0
    for mv, by_creation in zip(
        mvs,
        client.get_model_versions(name).order_by_creation_time(),
    ):
        assert mv.id == by_creation.id
        i += 1

    assert i == models

    i = 0
    for mv, by_update in zip(
        mvs,
        client.get_model_versions(name).order_by_update_time(),
    ):
        assert mv.id == by_update.id
        i += 1

    assert i == models

    for mv in reversed(mvs):
        mv.description = "updated"
        client.update(mv)

    i = 0
    for mv, by_update in zip(
        reversed(mvs),
        client.get_model_versions(name).order_by_update_time(),
    ):
        assert mv.id == by_update.id
        i += 1

    assert i == models

    i = 0
    for mv, by_update in zip(
        mvs,
        client.get_model_versions(name).order_by_update_time().descending(),
    ):
        assert mv.id == by_update.id
        i += 1

    assert i == models


@pytest.mark.e2e
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

    mv_iter = client.get_model_versions(name).page_size(model_count - 1)
    models = []
    for rm in islice(mv_iter, page):
        models.append(rm)
    assert len(models) == page
    mv_iter.restart()
    complete = list(mv_iter)
    assert len(complete) == model_count
    assert complete[:page] == models


@pytest.mark.e2e
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


@pytest.mark.e2e
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
