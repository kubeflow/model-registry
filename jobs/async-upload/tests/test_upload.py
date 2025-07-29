from dataclasses import asdict
import pytest
from unittest.mock import Mock, patch
from model_registry.utils import S3Params, OCIParams
from job.upload import _get_upload_params, perform_upload
from job.models import (
    AsyncUploadConfig,
    S3StorageConfig,
    OCIStorageConfig,
    ModelConfig,
    StorageConfig,
    RegistryConfig
)

class TestGetUploadParams:
    """Test cases for _get_upload_params function"""

    def test_get_upload_params_s3_config(self):
        """Test _get_upload_params with S3 configuration returns S3Params"""
        config = AsyncUploadConfig(
            source=S3StorageConfig(
                bucket="source-bucket",
                key="source-key",
                access_key_id="source-key-id",
                secret_access_key="source-secret",
                region="us-west-1"
            ),
            destination=S3StorageConfig(
                bucket="test-bucket",
                key="test-key",
                endpoint_url="https://s3.amazonaws.com",
                access_key_id="test-access-key",
                secret_access_key="test-secret-key",
                region="us-east-1"
            ),
            model=ModelConfig(
                id="test-model",
                version_id="test-version", 
                artifact_id="test-artifact"
            ),
            storage=StorageConfig(path="/tmp/test"),
            registry=RegistryConfig(server_address="test-server")
        )

        result = _get_upload_params(config)

        assert isinstance(result, S3Params)
        assert result.bucket_name == "test-bucket"
        assert result.s3_prefix == "test-key"
        assert result.endpoint_url == "https://s3.amazonaws.com"
        assert result.access_key_id == "test-access-key"
        assert result.secret_access_key == "test-secret-key"
        assert result.region == "us-east-1"

    def test_get_upload_params_oci_config(self):
        """Test _get_upload_params with OCI configuration returns OCIParams"""
        config = AsyncUploadConfig(
            source=S3StorageConfig(
                bucket="source-bucket",
                key="source-key",
                access_key_id="source-key-id",
                secret_access_key="source-secret",
                region="us-west-1"
            ),
            destination=OCIStorageConfig(
                uri="quay.io/example/test:latest",
                registry="quay.io",
                username="test-user",
                password="test-password",
                base_image="foo-bar:latest",
                enable_tls_verify=False,
                credentials_path="/tmp/test-creds"
            ),
            model=ModelConfig(
                id="abc",
                version_id="def",
                artifact_id="123"
            ),
            storage=StorageConfig(path="/tmp/test-model"),
            registry=RegistryConfig(server_address="test-server")
        )

        result = _get_upload_params(config)

        assert isinstance(result, OCIParams)
        assert result.base_image == "foo-bar:latest"
        assert result.oci_ref == "quay.io/example/test:latest"
        assert result.dest_dir == "/tmp/test-model"
        assert result.oci_username == "test-user"
        assert result.oci_password == "test-password"
        
    def test_get_upload_params_unsupported_type(self):
        """Test _get_upload_params with unsupported destination type raises ValueError"""
        # Create a mock config with an unsupported destination type
        config = Mock(spec=AsyncUploadConfig)
        config.destination = Mock()
        config.destination.__class__.__name__ = "UnsupportedStorageConfig"
        
        with pytest.raises(ValueError, match="Unsupported destination type"):
            _get_upload_params(config)

    def test_get_upload_params_oci_with_none_values(self):
        """Test _get_upload_params with OCI config containing None values"""
        config = AsyncUploadConfig(
            source=S3StorageConfig(
                bucket="source-bucket",
                key="source-key",
                access_key_id="source-key-id",
                secret_access_key="source-secret",
                region="us-west-1"
            ),
            destination=OCIStorageConfig(
                uri="quay.io/example/test:latest",
                registry="quay.io",
                username=None,
                password=None,
                base_image="foo-bar:latest",
                enable_tls_verify=False,
                credentials_path=None
            ),
            model=ModelConfig(
                id="abc",
                version_id="def",
                artifact_id="123"
            ),
            storage=StorageConfig(path="/tmp/test-model"),
            registry=RegistryConfig(server_address="test-server")
        )

        result = _get_upload_params(config)

        assert isinstance(result, OCIParams)
        assert result.base_image == "foo-bar:latest"
        assert result.oci_ref == "quay.io/example/test:latest"
        assert result.dest_dir == "/tmp/test-model"
        assert result.oci_username is None
        assert result.oci_password is None


class TestPerformUpload:
    """Test cases for perform_upload function"""

    @patch("job.upload.save_to_oci_registry")
    def test_perform_upload_oci(
        self, mock_save_to_oci_registry
    ):
        """Test perform_upload with OCI destination"""

        mock_save_to_oci_registry.return_value = 'quay.io/example/oci/abc:def'

        config = AsyncUploadConfig(
            source=S3StorageConfig(
                bucket="source-bucket",
                key="source-key",
                access_key_id="source-key-id",
                secret_access_key="source-secret",
                region="us-west-1"
            ),
            destination=OCIStorageConfig(
                uri="quay.io/example/oci",
                registry="quay.io",
                username="oci_user",
                password="oci_pass",
                base_image="foo-bar:latest",
                enable_tls_verify=False,
                credentials_path="/tmp/test-creds"
            ),
            model=ModelConfig(
                id="abc",
                version_id="def",
                artifact_id="123"
            ),
            storage=StorageConfig(path="/tmp/test-model"),
            registry=RegistryConfig(server_address="test-server")
        )

        # Act
        result_uri = perform_upload(config)

        # - And the returned URI is forwarded
        assert result_uri == "quay.io/example/oci/abc:def"

    @patch("job.upload._get_upload_params")
    def test_perform_upload_propagates_exceptions_from_get_upload_params(
        self, mock_get_upload_params
    ):
        """Test perform_upload propagates exceptions from _get_upload_params"""
        # Setup
        mock_client = Mock()
        mock_get_upload_params.side_effect = ValueError("Invalid config")

        config = AsyncUploadConfig(
            source=S3StorageConfig(
                bucket="source-bucket",
                key="source-key",
                access_key_id="source-key-id",
                secret_access_key="source-secret",
                region="us-west-1"
            ), 
            destination=S3StorageConfig(
                bucket="test-bucket",
                key="test-key",
                access_key_id="test-access-key",
                secret_access_key="test-secret-key",
                region="us-east-1"
            ),
            model=ModelConfig(
                id="test-model",
                version_id="1.0.0",
                artifact_id="test-artifact"
            ),
            storage=StorageConfig(path="/tmp/test-model"),
            registry=RegistryConfig(server_address="test-server")
        )

        # Execute and verify exception is propagated
        with pytest.raises(ValueError, match="Invalid config"):
            perform_upload(config)

        # Verify client method was not called
        mock_client.upload_artifact_and_register_model.assert_not_called()

    @patch("job.upload.save_to_oci_registry")
    def test_perform_upload_propagates_exceptions_from_client(
        self, mock_save_to_oci_registry
    ):
        """Test perform_upload propagates exceptions from client method"""
        # Setup
        mock_save_to_oci_registry.side_effect = Exception(
            "Upload failed"
        )

        config = AsyncUploadConfig(
            source=S3StorageConfig(
                bucket="source-bucket",
                key="source-key",
                access_key_id="source-key-id",
                secret_access_key="source-secret",
                region="us-west-1"
            ),
            destination=OCIStorageConfig(
                uri="quay.io/example/oci",
                registry="quay.io",
                username="oci_user",
                password="oci_pass",
                base_image="foo-bar:latest",
                enable_tls_verify=False,
                credentials_path="/tmp/test-creds"
            ),
            model=ModelConfig(
                id="abc",
                version_id="def",
                artifact_id="123"
            ),
            storage=StorageConfig(path="/tmp/test-model"),
            registry=RegistryConfig(server_address="test-server")
        )

        # Execute and verify exception is propagated
        with pytest.raises(Exception, match="Upload failed"):
            perform_upload(config)

        # Verify client method was called
        mock_save_to_oci_registry.assert_called_once()
