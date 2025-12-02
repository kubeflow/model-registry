import json
import os
from contextlib import contextmanager
from pathlib import Path

import pytest

from model_registry import utils
from model_registry.exceptions import MissingMetadata
from model_registry.utils import (
    _get_files_from_path,
    s3_uri_from,
    save_to_oci_registry,
    temp_auth_file,
)


def test_s3_uri_builder():
    s3_uri = s3_uri_from(
        "test-path",
        "test-bucket",
        endpoint="test-endpoint",
        region="test-region",
    )
    assert (
        s3_uri
        == "s3://test-bucket/test-path?endpoint=test-endpoint&defaultRegion=test-region"
    )


def test_s3_uri_builder_without_env():
    os.environ.pop("AWS_S3_BUCKET", None)
    os.environ.pop("AWS_S3_ENDPOINT", None)
    os.environ.pop("AWS_DEFAULT_REGION", None)
    with pytest.raises(MissingMetadata) as e:
        s3_uri_from(
            "test-path",
        )
    assert "custom environment" in str(e.value).lower()

    with pytest.raises(MissingMetadata) as e:
        s3_uri_from(
            "test-path",
            "test-bucket",
        )
    assert "non-default bucket" in str(e.value).lower()


def test_s3_uri_builder_with_only_default_bucket_env():
    os.environ["AWS_S3_BUCKET"] = "test-bucket"
    os.environ.pop("AWS_S3_ENDPOINT", None)
    os.environ.pop("AWS_DEFAULT_REGION", None)
    with pytest.raises(MissingMetadata) as e:
        s3_uri_from(
            "test-path",
        )
    assert "missing environment variable" in str(e.value).lower()


def test_s3_uri_builder_with_other_default_variables():
    os.environ.pop("AWS_S3_BUCKET", None)
    os.environ["AWS_S3_ENDPOINT"] = "test-endpoint"
    os.environ["AWS_DEFAULT_REGION"] = "test-region"
    with pytest.raises(MissingMetadata) as e:
        s3_uri_from(
            "test-path",
        )
    assert "custom environment" in str(e.value).lower()

    with pytest.raises(MissingMetadata) as e:
        s3_uri_from(
            "test-path",
            "test-bucket",
        )
    assert "non-default bucket" in str(e.value).lower()


def test_s3_uri_builder_with_complete_env():
    os.environ["AWS_S3_BUCKET"] = "test-bucket"
    os.environ["AWS_S3_ENDPOINT"] = "test-endpoint"
    os.environ["AWS_DEFAULT_REGION"] = "test-region"
    assert s3_uri_from("test-path") == s3_uri_from("test-path", "test-bucket")


@pytest.mark.e2e(type="oci")
def test_save_to_oci_registry_with_skopeo(get_temp_dir_with_models, get_temp_dir):
    base_image = "quay.io/mmortari/hello-world-wait:latest"
    dest_dir, _ = get_temp_dir_with_models
    oci_ref = "localhost:5001/foo/bar:latest"

    save_to_oci_registry(
        base_image=base_image,
        oci_ref=oci_ref,
        model_files_path=dest_dir,
        dest_dir=get_temp_dir,
        custom_oci_backend=utils._get_skopeo_backend(
            push_args=[
                "--dest-tls-verify=false",
                "--src-username=non_user",
                "--src-password=nonpassword",
            ],
        ),
    )


def test_save_to_oci_registry_with_custom_backend(
    get_temp_dir_with_models, get_temp_dir, get_mock_custom_oci_backend
):
    backend = "something_custom"
    # similar to other test
    base_image = "busybox:latest"
    dest_dir, _ = get_temp_dir_with_models
    oci_ref = "localhost:5001/foo/bar:latest"

    uri = save_to_oci_registry(
        base_image=base_image,
        oci_ref=oci_ref,
        model_files_path=dest_dir,
        dest_dir=get_temp_dir,
        backend=backend,
        modelcard=None,
        custom_oci_backend=get_mock_custom_oci_backend,
    )
    # Ensure our mocked backend was called
    get_mock_custom_oci_backend.is_available.assert_called_once()
    get_mock_custom_oci_backend.pull.assert_called_once()
    get_mock_custom_oci_backend.push.assert_called_once()
    assert uri == f"oci://{oci_ref}"


def test_save_to_oci_registry_with_username_password(mocker, tmp_path):
    model_files_path = tmp_path / "model-files"
    model_files_path.mkdir()
    (model_files_path / "model.bin").touch()
    dest_dir = tmp_path / "dest"

    temp_auth_file_info = {}

    @contextmanager
    def temp_auth_file_wrapper(auth):
        with temp_auth_file(auth) as f:
            temp_auth_file_info["path"] = f.name
            temp_auth_file_info["contents"] = Path(f.name).read_text()
            yield f

    mock_skopeo_pull = mocker.patch("olot.backend.skopeo.skopeo_pull")
    mock_skopeo_push = mocker.patch("olot.backend.skopeo.skopeo_push")
    mocker.patch("olot.basics.oci_layers_on_top")
    mocker.patch("model_registry.utils.temp_auth_file", side_effect=temp_auth_file_wrapper)

    save_to_oci_registry(
        base_image="busybox",
        oci_ref="quay.io/example/example:latest",
        model_files_path=model_files_path,
        dest_dir=dest_dir,
        backend="skopeo",
        oci_username="user32",
        oci_password="zi3327",  # noqa: S106
    )

    assert mock_skopeo_pull.call_args.args == ("busybox", dest_dir, ["--src-authfile", mocker.ANY])
    assert mock_skopeo_pull.call_args.kwargs == {}
    assert mock_skopeo_push.call_args.args == (dest_dir, "quay.io/example/example:latest", ["--dest-authfile", mocker.ANY])
    assert mock_skopeo_push.call_args.kwargs == {}
    assert json.loads(temp_auth_file_info["contents"]) == {"auths": {"quay.io/example/example": {"auth": "dXNlcjMyOnppMzMyNw=="}}}
    assert not Path(temp_auth_file_info["path"]).exists()


@pytest.mark.e2e(type="oci")
def test_save_to_oci_registry_auth_params(
    get_temp_dir_with_models,
    get_temp_dir,
    get_mock_skopeo_backend_for_auth,
):
    # similar to other test
    base_image = "busybox:latest"
    dest_dir, _ = get_temp_dir_with_models
    oci_ref = "localhost:5001/foo/bar:latest"
    backend, skopeo_pull_mock, skopeo_push_mock, generic_params = (
        get_mock_skopeo_backend_for_auth
    )

    assert os.getenv(".dockerconfigjson") is not None  # noqa: SIM112 (no capitalization)

    save_to_oci_registry(
        base_image=base_image,
        oci_ref=oci_ref,
        model_files_path=dest_dir,
        dest_dir=get_temp_dir,
        modelcard=None,
        custom_oci_backend=backend,
    )
    skopeo_pull_mock.assert_called_once()
    skopeo_push_mock.assert_called_once()
    args, _ = skopeo_pull_mock.call_args
    params = args[2]
    assert generic_params[0] in params
    assert generic_params[-1] in params


def test_save_to_oci_registry_backend_not_found():
    backend = "non-existent"
    with pytest.raises(
        ValueError, match=f"'{backend}' is not an available backend to use."
    ) as e:
        save_to_oci_registry("", "", [], "", backend)  # type: ignore[arg-type]

    assert f"'{backend}' is not an available backend to use." in str(e.value)


def test_get_files_from_path_no_path():
    path = "/in/val/id/pa/th"
    with pytest.raises(ValueError, match="Please ensure path is correct.") as e:
        _get_files_from_path(path)
    assert e


def test_get_files_from_path_single_file(get_model_file):
    file = _get_files_from_path(get_model_file)
    # It returns only 1 file in the list, and it is a tuple of (absolute_path, filename)
    assert len(file) == 1
    assert file[0] == (get_model_file, os.path.basename(get_model_file))


def test_get_files_from_path_multiple_files(get_temp_dir_with_models):
    path, generated_files = get_temp_dir_with_models
    files = _get_files_from_path(path)
    # It returns the same number of files as were generated, and it is a list tuple of (absolute_path, filename)
    assert len(files) == len(generated_files)
    for abs, filename in files:
        assert abs == os.path.join(path, filename)
        assert filename == os.path.relpath(abs, path)
