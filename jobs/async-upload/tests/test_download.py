import os
import pytest
from unittest.mock import Mock, patch
from job.download import download_from_s3
from job.config import get_config
from job.mr_client import validate_and_get_model_registry_client


@pytest.fixture
def minimal_env_source_dest_vars():
    original_env = dict(os.environ)

    # Destination variables
    dest_vars = {
        "type": "oci",
        "oci_uri": "quay.io/example/oci",
        "oci_username": "oci_username_env",
        "oci_password": "oci_password_env",
    }

    # Source variables - using the correct format from existing tests
    source_vars = {
        "type": "s3",
        "aws_bucket": "test-bucket",
        "aws_key": "test-key",
        "aws_access_key_id": "test-access-key-id",
        "aws_secret_access_key": "test-secret-access-key",
        "aws_endpoint": "http://localhost:9000",
    }

    # Set up test environment variables
    for key, value in dest_vars.items():
        os.environ[f"MODEL_SYNC_DESTINATION_{key.upper()}"] = value
    for key, value in source_vars.items():
        os.environ[f"MODEL_SYNC_SOURCE_{key.upper()}"] = value

    # Model and registry variables
    model_vars = {
        "model_name": "my-model",
        "model_version": "1.0.0",
        "model_format": "onnx",
        "registry_server_address": "http://localhost",
        "registry_port": "8080",
        "registry_author": "author",
        "storage_path": "/tmp/model-sync",
    }

    for key, value in model_vars.items():
        os.environ[f"MODEL_SYNC_{key.upper()}"] = value

    yield model_vars

    # Restore original environment
    os.environ.clear()
    os.environ.update(original_env)


def test_download_from_s3(minimal_env_source_dest_vars):
    """Test download_from_s3 function with proper configuration"""

    # Get configuration from environment variables
    config = get_config([])

    # Verify the configuration is set up correctly
    assert config["source"]["type"] == "s3"
    assert config["source"]["s3"]["bucket"] == "test-bucket"
    assert config["source"]["s3"]["key"] == "test-key"
    assert config["source"]["s3"]["access_key_id"] == "test-access-key-id"
    assert config["source"]["s3"]["secret_access_key"] == "test-secret-access-key"
    assert config["source"]["s3"]["endpoint"] == "http://localhost:9000"
    assert config["storage"]["path"] == "/tmp/model-sync"

    # Create mock ModelRegistry client
    with patch("job.mr_client.ModelRegistry") as mock_registry_class:
        mock_client = Mock()
        mock_registry_class.return_value = mock_client
        client = validate_and_get_model_registry_client(config)

        # Mock the S3 client and _connect_to_s3 function
        with patch("job.download._connect_to_s3") as mock_connect:
            mock_s3_client = Mock()
            mock_transfer_config = Mock()
            mock_connect.return_value = (mock_s3_client, mock_transfer_config)

            # Call the function under test
            download_from_s3(client, config)

            # Verify _connect_to_s3 was called with correct parameters
            mock_connect.assert_called_once_with(
                endpoint_url="http://localhost:9000",
                access_key_id="test-access-key-id",
                secret_access_key="test-secret-access-key",
                region=None,  # Not set in the test fixture
            )

            # Verify download_file was called with correct parameters
            mock_s3_client.download_file.assert_called_once_with(
                "test-bucket",  # bucket
                "test-key",  # key
                "/tmp/model-sync",  # local path
            )


def test_download_from_s3_with_region(minimal_env_source_dest_vars):
    """Test download_from_s3 function with region specified"""

    # Set region in environment
    os.environ["MODEL_SYNC_SOURCE_AWS_REGION"] = "us-west-2"

    config = get_config([])

    # Create mock ModelRegistry client
    with patch("job.mr_client.ModelRegistry") as mock_registry_class:
        mock_client = Mock()
        mock_registry_class.return_value = mock_client
        client = validate_and_get_model_registry_client(config)

        # Mock the S3 client and _connect_to_s3 function
        with patch("job.download._connect_to_s3") as mock_connect:
            mock_s3_client = Mock()
            mock_transfer_config = Mock()
            mock_connect.return_value = (mock_s3_client, mock_transfer_config)

            # Call the function under test
            download_from_s3(client, config)

            # Verify _connect_to_s3 was called with correct parameters including region
            mock_connect.assert_called_once_with(
                endpoint_url="http://localhost:9000",
                access_key_id="test-access-key-id",
                secret_access_key="test-secret-access-key",
                region="us-west-2",
            )


def test_download_from_s3_connection_error(minimal_env_source_dest_vars):
    """Test download_from_s3 function when S3 connection fails"""

    config = get_config([])

    # Create mock ModelRegistry client
    with patch("job.mr_client.ModelRegistry") as mock_registry_class:
        mock_client = Mock()
        mock_registry_class.return_value = mock_client
        client = validate_and_get_model_registry_client(config)

        # Mock _connect_to_s3 to raise an exception
        with patch("job.download._connect_to_s3") as mock_connect:
            mock_connect.side_effect = Exception("Connection failed")

            # Test that the exception is propagated
            with pytest.raises(Exception, match="Connection failed"):
                download_from_s3(client, config)


def test_download_from_s3_download_error(minimal_env_source_dest_vars):
    """Test download_from_s3 function when file download fails"""

    config = get_config([])

    # Create mock ModelRegistry client
    with patch("job.mr_client.ModelRegistry") as mock_registry_class:
        mock_client = Mock()
        mock_registry_class.return_value = mock_client
        client = validate_and_get_model_registry_client(config)

        # Mock the S3 client and _connect_to_s3 function
        with patch("job.download._connect_to_s3") as mock_connect:
            mock_s3_client = Mock()
            mock_transfer_config = Mock()
            mock_connect.return_value = (mock_s3_client, mock_transfer_config)

            # Mock download_file to raise an exception
            mock_s3_client.download_file.side_effect = Exception("Download failed")

            # Test that the exception is propagated
            with pytest.raises(Exception, match="Download failed"):
                download_from_s3(client, config)
