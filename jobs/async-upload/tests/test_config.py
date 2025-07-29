import os
import random
import tempfile
import shutil
from pathlib import Path
import pytest
from job.config import get_config

from job.models import S3StorageConfig, OCIStorageConfig

MR_PREFIX = "MODEL_SYNC"
MR_SOURCE_PREFIX = "MODEL_SYNC_SOURCE"
MR_DEST_PREFIX = "MODEL_SYNC_DESTINATION"


@pytest.fixture
def source_s3_env_vars():
    """Fixture to set up s3-specific source environment variables"""
    # Save original environment
    original_env = dict(os.environ)

    vars = {
        "type": "s3",
        "aws_access_key_id": "source_key_env",
        "aws_secret_access_key": "source_secret_env",
        "aws_region": "us-east-1",
        "aws_bucket": "source-bucket-env",
        "aws_key": "source-key-env",
    }

    for key, value in vars.items():
        os.environ[f"MODEL_SYNC_SOURCE_{key.upper()}"] = value

    yield vars

    # Restore original environment
    os.environ.clear()
    os.environ.update(original_env)


@pytest.fixture
def destination_oci_env_vars():
    """Fixture to set up oci-specific destination environment variables"""
    # Save original environment
    original_env = dict(os.environ)

    vars = {
        "type": "oci",
        "oci_uri": "quay.io/example/oci",
        "oci_registry": "quay.io",
        "oci_username": "oci_username_env",
        "oci_password": "oci_password_env",
    }

    # Set up test environment variables
    for key, value in vars.items():
        os.environ[f"MODEL_SYNC_DESTINATION_{key.upper()}"] = value

    yield vars

    # Restore original environment
    os.environ.clear()
    os.environ.update(original_env)


@pytest.fixture
def model_env_vars():
    """Fixture to set up environment variables for testing"""
    # Save original environment
    original_env = dict(os.environ)

    vars = {
        "model_id": "1234",
        "model_version_id": "0987",
        "model_artifact_id": "5678",
        "registry_server_address": "https://registry.example.com",
    }

    for key, value in vars.items():
        os.environ[f"MODEL_SYNC_{key.upper()}"] = value

    yield vars

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
            "key": f"file-key_{random.randint(1000, 9999)}",
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
        with open(os.path.join(temp_dir, "AWS_KEY"), "w") as f:
            f.write(credentials["key"])

        yield Path(temp_dir), credentials

        # Clean up temporary folder
        shutil.rmtree(temp_dir)


def test_s3_file_to_oci_env_config(
    s3_credentials_folder, destination_oci_env_vars, model_env_vars
):
    """Tests a configuration where the source is S3, using a credentials path, to a destination with OCI env vars"""
    folder_location, expected_credentials = s3_credentials_folder

    # The same folder/files are used for both source and destination
    config = get_config(
        [
            "--source-type",
            "s3",
            "--source-s3-credentials-path",
            str(folder_location),
            "--registry-server-address",
            "https://registry.example.com",
        ]
    )

    # Source credentials were not set, so they are None
    assert isinstance(config.source, S3StorageConfig)
    assert config.source.access_key_id == expected_credentials["access_key_id"]
    assert config.source.secret_access_key == expected_credentials["secret_access_key"]
    assert config.source.region == expected_credentials["region"]
    assert config.source.bucket == expected_credentials["bucket"]

    assert isinstance(config.destination, OCIStorageConfig)
    assert config.destination.uri == destination_oci_env_vars["oci_uri"]
    assert config.destination.username == destination_oci_env_vars["oci_username"]
    assert config.destination.password == destination_oci_env_vars["oci_password"]


def test_env_based_s3_to_oci_config(
    model_env_vars, source_s3_env_vars, destination_oci_env_vars
):
    """Test configuration using environment variables"""
    config = get_config([])

    assert isinstance(config.source, S3StorageConfig)
    assert config.source.access_key_id == source_s3_env_vars["aws_access_key_id"]
    assert (
        config.source.secret_access_key
        == source_s3_env_vars["aws_secret_access_key"]
    )
    assert config.source.region == source_s3_env_vars["aws_region"]
    assert config.source.bucket == source_s3_env_vars["aws_bucket"]

    assert isinstance(config.destination, OCIStorageConfig)
    assert config.destination.uri == destination_oci_env_vars["oci_uri"]
    assert config.destination.username == destination_oci_env_vars["oci_username"]
    assert config.destination.password == destination_oci_env_vars["oci_password"]

    assert config.model.id == model_env_vars["model_id"]
    assert config.model.version_id == model_env_vars["model_version_id"]
    assert config.model.artifact_id == model_env_vars["model_artifact_id"]
    assert config.registry.server_address == model_env_vars["registry_server_address"]


def test_params_based_config():
    """Test configuration using parameters"""
    config = get_config(
        [
            "--source-type",
            "oci",
            "--source-oci-uri",
            "quay.io/example/params",
            "--source-oci-registry",
            "quay.io",
            "--destination-type",
            "s3",
            "--destination-aws-bucket",
            "destination-bucket-params",
            "--destination-aws-key",
            "destination-key-params",
            "--destination-aws-region",
            "eu-central-1",
            "--destination-aws-access-key-id",
            "destination_key_params",
            "--destination-aws-secret-access-key",
            "destination_secret_params",
            "--model-id",
            "1234",
            "--model-version-id",
            "0987",
            "--model-artifact-id",
            "5678",
            "--registry-server-address",
            "https://registry.example.com",
        ]
    )

    assert isinstance(config.source, OCIStorageConfig)
    assert config.source.uri == "quay.io/example/params"
    assert isinstance(config.destination, S3StorageConfig)
    assert config.destination.access_key_id == "destination_key_params"
    assert config.destination.secret_access_key == "destination_secret_params"
    assert config.destination.region == "eu-central-1"
    assert config.destination.bucket == "destination-bucket-params"


def test_params_will_override_env_config(
    model_env_vars, source_s3_env_vars, destination_oci_env_vars
):
    """Test a configuration in which ENV vars are set, but override params are provided to the CLI"""

    override_vars = {
        "aws_bucket": f"source-bucket-{random.randint(1000, 9999)}",
        "aws_region": f"eu-central-{random.randint(1000, 9999)}",
        "aws_access_key_id": f"source-key-{random.randint(1000, 9999)}",
        "aws_secret_access_key": f"source-secret-{random.randint(1000, 9999)}",
    }
    config = get_config(
        [
            "--source-type",
            "s3",
            "--source-aws-bucket",
            override_vars["aws_bucket"],
            "--source-aws-region",
            override_vars["aws_region"],
            "--source-aws-access-key-id",
            override_vars["aws_access_key_id"],
            "--source-aws-secret-access-key",
            override_vars["aws_secret_access_key"],
            "--model-id",
            "1234",
            "--model-version-id",
            "0987",
            "--model-artifact-id",
            "5678",
            "--registry-server-address",
            "https://registry.example.com",
        ]
    )

    assert isinstance(config.source, S3StorageConfig)
    assert config.source.access_key_id == override_vars["aws_access_key_id"]
    assert config.source.secret_access_key == override_vars["aws_secret_access_key"]
    assert config.source.region == override_vars["aws_region"]
    assert config.source.bucket == override_vars["aws_bucket"]
