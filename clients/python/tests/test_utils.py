import os
import pathlib
from unittest.mock import Mock

import pytest

from model_registry.exceptions import MissingMetadata
<<<<<<< HEAD
from model_registry.utils import s3_uri_from, save_to_oci_registry
=======
from model_registry.utils import s3_uri_from, is_s3_uri, is_oci_uri
>>>>>>> 6ebe337 (chore: add oci and s3 helper methods)


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

    save_to_oci_registry(base_image, oci_ref, model_files, dest_dir, backend)


# These are trimmed down versions of whats found in the example specs found here: https://github.com/opencontainers/image-spec/blob/main/image-layout.md#oci-layout-file
index_json_contents = """{
  "schemaVersion": 2,
  "mediaType": "application/vnd.oci.image.index.v1+json",
  "manifests": [],
  "annotations": {
    "com.example.index.revision": "r124356"
  }
}"""
oci_layout_contents = """{"imageLayoutVersion": "1.0.0"}"""

def test_save_to_oci_registry_with_custom_backend():
    is_available_mock = Mock()
    is_available_mock.return_value = True
    pull_mock = Mock()
    push_mock = Mock()
    def pull_mock_imple(base_image, dest_dir):
        pathlib.Path(dest_dir).joinpath("oci-layout").write_text(oci_layout_contents)
        pathlib.Path(dest_dir).joinpath("index.json").write_text(index_json_contents)

    pull_mock.side_effect = pull_mock_imple


    backend = "something_custom"
    backend_registry = {
        "something_custom": lambda: {
            "is_available": is_available_mock,
            "pull": pull_mock,
            "push": push_mock,
        }
    }

    # similar to other test
    base_image = "quay.io/mmortari/hello-world-wait:latest"
    dest_dir = "tests/data"
    oci_ref = "localhost:5001/foo/bar:latest"

    # Create a sample file named README.md to be added to the registry
    pathlib.Path(dest_dir).mkdir(parents=True, exist_ok=True)
    readme_file_path = os.path.join(dest_dir, "README.md")
    with open(readme_file_path, "w") as f:
        f.write("")

    model_files = [readme_file_path]

    uri = save_to_oci_registry(base_image, oci_ref, model_files, dest_dir, backend, None, backend_registry)
    # Ensure our mocked backend was called
    is_available_mock.assert_called_once()
    pull_mock.assert_called_once()
    push_mock.assert_called_once()
    assert uri == f"oci://{oci_ref}"

def test_save_to_oci_registry_with_custom_backend_unavailable():
    is_available_mock = Mock()
    is_available_mock.return_value = False # Backend is unavailable, expect an error
    pull_mock = Mock()
    push_mock = Mock()

    backend = "something_custom"
    backend_registry = {
        "something_custom": lambda: {
            "is_available": is_available_mock,
            "pull": pull_mock,
            "push": push_mock,
        }
    }


    with pytest.raises(ValueError, match=f"Backend '{backend}' is selected, but not available on the system. Ensure the dependencies for '{backend}' are installed in your environment.") as e:
        save_to_oci_registry("", "", [], "", backend, backend_registry=backend_registry)

    assert f"Backend '{backend}' is selected, but not available on the system." in str(e.value)

def test_save_to_oci_registry_backend_not_found():
    backend = "non-existent"
    with pytest.raises(ValueError, match=f"'{backend}' is not an available backend to use.") as e:
        save_to_oci_registry("", "", [], "", backend)

    assert f"'{backend}' is not an available backend to use." in str(e.value)

def test_is_s3_uri_with_valid_uris():
    test_cases = [
        "s3://my-bucket/my-file.txt",
        "s3://my-bucket/my-folder/my-file.conf",
        "s3://my-bucket/my-folder/my-sub-folder/my-file.sh",
    ]
    for test in test_cases:
        assert is_s3_uri(test) == True

def test_is_s3_uri_with_invalid_uris():
    test_cases = [
        "",
        "s3://",
        "s3://my-file.txt",
        "my-bucket/my-file.sh",
    ]
    for test in test_cases:
        assert is_s3_uri(test) == False

def test_is_oci_uri_with_valid_uris():
    test_cases = [
        "oci://registry.example.com/my-namespace/my-repo:latest",
        "oci://localhost:5000/my-repo",
        "oci://registry.example.com/my-repo",
        "oci://registry.example.com/my-repo:1.0.0",
    ]

    for test in test_cases:
        assert is_oci_uri(test) == True

def test_is_oci_uri_with_invalid_uris():
    test_cases = [
        "",
        "oci://",
        "oci://registry.example.com"
        "oci://localhost:5000"
        "oci://registry.example.com/"
        "oci://registry.example.com/my-repo/"
        "oci://registry.example.com/my-repo/something-wrong"
    ]

    for test in test_cases:
        assert is_oci_uri(test) == False

