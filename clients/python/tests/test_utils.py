import os

import pytest

from model_registry.exceptions import MissingMetadata
from model_registry.utils import (
    get_files_from_path,
    s3_uri_from,
    save_to_oci_registry,
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
    base_image = "busybox:latest"
    dest_dir, _ = get_temp_dir_with_models
    oci_ref = "localhost:5001/foo/bar:latest"

    backend = "skopeo"

    save_to_oci_registry(base_image, oci_ref, dest_dir, get_temp_dir, backend)


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
    get_mock_custom_oci_backend["is_available"].assert_called_once()
    get_mock_custom_oci_backend["pull"].assert_called_once()
    get_mock_custom_oci_backend["push"].assert_called_once()
    assert uri == f"oci://{oci_ref}"


def test_save_to_oci_registry_backend_not_found():
    backend = "non-existent"
    with pytest.raises(
        ValueError, match=f"'{backend}' is not an available backend to use."
    ) as e:
        save_to_oci_registry("", "", [], "", backend)

    assert f"'{backend}' is not an available backend to use." in str(e.value)


def test_get_files_from_path_no_path():
    path = "/in/val/id/pa/th"
    with pytest.raises(ValueError, match="Please ensure path is correct.") as e:
        get_files_from_path(path)
    assert e


def test_get_files_from_path_single_file(get_model_file):
    file = get_files_from_path(get_model_file)
    # It returns only 1 file in the list, and it is a tuple of (absolute_path, filename)
    assert len(file) == 1
    assert file[0] == (get_model_file, os.path.basename(get_model_file))


def test_get_files_from_path_multiple_files(get_temp_dir_with_models):
    path, generated_files = get_temp_dir_with_models
    files = get_files_from_path(path)
    # It returns the same number of files as were generated, and it is a list tuple of (absolute_path, filename)
    assert len(files) == len(generated_files)
    for abs, filename in files:
        assert abs == os.path.join(path, filename)
        assert filename == os.path.relpath(abs, path)
