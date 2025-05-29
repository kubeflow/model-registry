import os
import random
import tempfile
import shutil
from pathlib import Path
import pytest
from job.config import get_config

MR_PREFIX = "MODEL_SYNC"
MR_SOURCE_PREFIX = "MODEL_SYNC_SOURCE"
MR_DEST_PREFIX = "MODEL_SYNC_DESTINATION"


@pytest.fixture
def env_vars():
    """Fixture to set up environment variables for testing"""
    # Save original environment
    original_env = dict(os.environ)

    os.environ["MODEL_SYNC_SOURCE_TYPE"] = "oci"
    os.environ["MODEL_SYNC_SOURCE_OCI_URI"] = "quay.io/example/models"

    # Set up test environment variables
    os.environ["MODEL_SYNC_DESTINATION_TYPE"] = "s3"
    os.environ["MODEL_SYNC_DESTINATION_AWS_ACCESS_KEY_ID"] = "destination_key_env"
    os.environ["MODEL_SYNC_DESTINATION_AWS_SECRET_ACCESS_KEY"] = (
        "destination_secret_env"
    )
    os.environ["MODEL_SYNC_DESTINATION_AWS_REGION"] = "us-east-1"
    os.environ["MODEL_SYNC_DESTINATION_AWS_BUCKET"] = "destination-bucket-env"

    os.environ["MODEL_SYNC_MODEL_NAME"] = "my-model"
    os.environ["MODEL_SYNC_MODEL_VERSION"] = "1.0.0"
    os.environ["MODEL_SYNC_MODEL_FORMAT"] = "onnx"
    os.environ["MODEL_SYNC_REGISTRY_URL"] = "https://registry.example.com"

    yield

    # Restore original environment
    os.environ.clear()
    os.environ.update(original_env)


@pytest.fixture
def s3_credentials_folder():
    """
    Fixture to create temporary config files for testing
    """
    # Create source config file
    with tempfile.TemporaryDirectory() as temp_dir:
        # Create the credentials files and store their contents
        credentials = {
            "access_key_id": f"file_key_{random.randint(1000, 9999)}",
            "secret_access_key": f"file_secret_{random.randint(1000, 9999)}",
            "region": f"eu-west-1_{random.randint(1000, 9999)}",
            "bucket": f"file-bucket_{random.randint(1000, 9999)}",
        }

        # Write the credentials to files
        with open(os.path.join(temp_dir, "AWS_ACCESS_KEY_ID"), "w") as f:
            f.write(credentials["access_key_id"])
        with open(os.path.join(temp_dir, "AWS_SECRET_ACCESS_KEY"), "w") as f:
            f.write(credentials["secret_access_key"])
        with open(os.path.join(temp_dir, "AWS_REGION"), "w") as f:
            f.write(credentials["region"])
        with open(os.path.join(temp_dir, "AWS_BUCKET"), "w") as f:
            f.write(credentials["bucket"])

        yield Path(temp_dir), credentials

        # Clean up temporary folder
        shutil.rmtree(temp_dir)


def test_credentials_dest_folder_based_config(s3_credentials_folder):
    """Test configuration using credentials folder"""
    folder_location, expected_credentials = s3_credentials_folder

    # The same folder/files are used for both source and destination
    config = get_config(
        [
            "--destination-s3-credentials-path",
            str(folder_location),
            "--source-s3-credentials-path",
            str(folder_location),
            "--model-name",
            "my-model",
            "--model-version",
            "1.0.0",
            "--model-format",
            "onnx",
            "--registry-url",
            "https://registry.example.com",
        ]
    )

    # Source credentials were not set, so they are None
    assert config["source"]["type"] == "s3"
    assert (
        config["source"]["s3"]["access_key_id"] == expected_credentials["access_key_id"]
    )
    assert (
        config["source"]["s3"]["secret_access_key"]
        == expected_credentials["secret_access_key"]
    )
    assert config["source"]["s3"]["region"] == expected_credentials["region"]
    assert config["source"]["s3"]["bucket"] == expected_credentials["bucket"]

    assert config["destination"]["type"] == "s3"
    assert (
        config["destination"]["s3"]["access_key_id"]
        == expected_credentials["access_key_id"]
    )
    assert (
        config["destination"]["s3"]["secret_access_key"]
        == expected_credentials["secret_access_key"]
    )
    assert config["destination"]["s3"]["region"] == expected_credentials["region"]
    assert config["destination"]["s3"]["bucket"] == expected_credentials["bucket"]


def test_env_based_config(env_vars):
    """Test configuration using environment variables"""
    config = get_config([])

    assert config["source"]["type"] == "oci"
    assert config["source"]["oci"]["uri"] == "quay.io/example/models"
    assert config["destination"]["type"] == "s3"
    assert config["destination"]["s3"]["access_key_id"] == "destination_key_env"
    assert config["destination"]["s3"]["secret_access_key"] == "destination_secret_env"
    assert config["destination"]["s3"]["region"] == "us-east-1"
    assert config["destination"]["s3"]["bucket"] == "destination-bucket-env"


def test_params_based_config():
    """Test configuration using parameters"""
    config = get_config(
        [
            "--source-type",
            "oci",
            "--source-oci-uri",
            "quay.io/example/params",
            "--destination-type",
            "s3",
            "--destination-aws-bucket",
            "destination-bucket-params",
            "--destination-aws-region",
            "eu-central-1",
            "--destination-aws-access-key-id",
            "destination_key_params",
            "--destination-aws-secret-access-key",
            "destination_secret_params",
            "--model-name",
            "my-model",
            "--model-version",
            "1.0.0",
            "--model-format",
            "onnx",
            "--registry-url",
            "https://registry.example.com",
        ]
    )

    assert config["source"]["type"] == "oci"
    assert config["source"]["oci"]["uri"] == "quay.io/example/params"
    assert config["destination"]["type"] == "s3"
    assert config["destination"]["s3"]["access_key_id"] == "destination_key_params"
    assert (
        config["destination"]["s3"]["secret_access_key"] == "destination_secret_params"
    )
    assert config["destination"]["s3"]["region"] == "eu-central-1"
    assert config["destination"]["s3"]["bucket"] == "destination-bucket-params"


def test_params_will_override_env_config(env_vars):
    """Test configuration using parameters and environment variables

    source is from env, destination is from params
    """
    config = get_config(
        [
            "--destination-type",
            "s3",
            "--destination-aws-bucket",
            "destination-bucket-params",
            "--destination-aws-region",
            "eu-central-1",
            "--destination-aws-access-key-id",
            "destination_key_params",
            "--destination-aws-secret-access-key",
            "destination_secret_params",
            "--model-name",
            "my-model",
            "--model-version",
            "1.0.0",
            "--model-format",
            "onnx",
            "--registry-url",
            "https://registry.example.com",
        ]
    )

    assert config["source"]["type"] == "oci"
    assert config["source"]["oci"]["uri"] == "quay.io/example/models"
    assert config["destination"]["type"] == "s3"
    assert config["destination"]["s3"]["access_key_id"] == "destination_key_params"
    assert (
        config["destination"]["s3"]["secret_access_key"] == "destination_secret_params"
    )
    assert config["destination"]["s3"]["region"] == "eu-central-1"
    assert config["destination"]["s3"]["bucket"] == "destination-bucket-params"
