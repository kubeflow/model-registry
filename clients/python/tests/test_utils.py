import os

import pytest

from model_registry.exceptions import MissingMetadata
from model_registry.utils import s3_uri_from, is_s3_uri, is_oci_uri


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

