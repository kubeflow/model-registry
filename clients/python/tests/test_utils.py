import os
import pathlib

import pytest

from model_registry.exceptions import MissingMetadata
from model_registry.utils import s3_uri_from, save_to_oci_registry


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
    assert "missing environment variables" in str(e.value).lower()


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
def test_save_to_oci_registry_with_skopeo():
    # TODO: We need a good source registry which is oci-compliant and very small in size
    base_image = "quay.io/mmortari/hello-world-wait:latest"
    dest_dir = "tests/data"
    oci_ref = "localhost:5001/foo/bar:latest"

    # Create a sample file named README.md to be added to the registry
    pathlib.Path(dest_dir).mkdir(parents=True, exist_ok=True)
    readme_file_path = os.path.join(dest_dir, "README.md")
    with open(readme_file_path, "w") as f:
        f.write("")

    model_files = [readme_file_path]
    backend = "skopeo"

    save_to_oci_registry(base_image, dest_dir, oci_ref, model_files, backend)
