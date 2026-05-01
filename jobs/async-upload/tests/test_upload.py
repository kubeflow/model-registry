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
    RegistryConfig,
    UpdateArtifactIntent
)

class TestGetUploadParamsOCIPullArgs:
    """Test cases for base image pull_args passed to _get_skopeo_backend"""

    @patch("job.upload.utils._get_skopeo_backend")
    def test_pull_args_with_tls_disabled_and_credentials(self, mock_get_backend):
        mock_get_backend.return_value = Mock()

        config = AsyncUploadConfig(
            source=S3StorageConfig(
                bucket="src-bucket", key="src-key",
                access_key_id="id", secret_access_key="secret", region="us-east-1"
            ),
            destination=OCIStorageConfig(
                uri="registry.internal/model:v1", registry="registry.internal",
                username="user", password="pass",
                base_image="registry.internal/busybox:latest",
                base_image_tls_verify=False,
                base_image_credentials_path="/etc/pull-secret/.dockerconfigjson",
                enable_tls_verify=False,
                credentials_path="/tmp/push-creds"
            ),
            model=ModelConfig(intent=UpdateArtifactIntent(artifact_id="123")),
            storage=StorageConfig(path="/tmp/test"),
            registry=RegistryConfig(server_address="test-server")
        )

        _get_upload_params(config)

        mock_get_backend.assert_called_once()
        call_kwargs = mock_get_backend.call_args
        pull_args = call_kwargs.kwargs.get("pull_args") or call_kwargs[1].get("pull_args")
        push_args = call_kwargs.kwargs.get("push_args") or call_kwargs[1].get("push_args")

        assert "--src-tls-verify=false" in pull_args
        assert "--authfile" in pull_args
        assert "/etc/pull-secret/.dockerconfigjson" in pull_args
        assert "--dest-tls-verify=false" in push_args
        assert "--authfile" in push_args
        assert "/tmp/push-creds" in push_args

    @patch("job.upload.utils._get_skopeo_backend")
    def test_pull_args_defaults_are_empty(self, mock_get_backend):
        mock_get_backend.return_value = Mock()

        config = AsyncUploadConfig(
            source=S3StorageConfig(
                bucket="src-bucket", key="src-key",
                access_key_id="id", secret_access_key="secret", region="us-east-1"
            ),
            destination=OCIStorageConfig(
                uri="quay.io/org/model:v1", registry="quay.io",
                username="user", password="pass",
                base_image="quay.io/quay/busybox:latest",
            ),
            model=ModelConfig(intent=UpdateArtifactIntent(artifact_id="123")),
            storage=StorageConfig(path="/tmp/test"),
            registry=RegistryConfig(server_address="test-server")
        )

        _get_upload_params(config)

        mock_get_backend.assert_called_once()
        call_kwargs = mock_get_backend.call_args
        pull_args = call_kwargs.kwargs.get("pull_args") or call_kwargs[1].get("pull_args")
        push_args = call_kwargs.kwargs.get("push_args") or call_kwargs[1].get("push_args")

        assert pull_args == []
        assert push_args == []

    @patch("job.upload.utils._get_skopeo_backend")
    def test_pull_args_independent_from_push_args(self, mock_get_backend):
        """Base image pull settings are independent from destination push settings"""
        mock_get_backend.return_value = Mock()

        config = AsyncUploadConfig(
            source=S3StorageConfig(
                bucket="src-bucket", key="src-key",
                access_key_id="id", secret_access_key="secret", region="us-east-1"
            ),
            destination=OCIStorageConfig(
                uri="registry.internal/model:v1", registry="registry.internal",
                username="user", password="pass",
                base_image="registry.internal/busybox:latest",
                base_image_tls_verify=False,
                base_image_credentials_path="/etc/pull-secret/.dockerconfigjson",
                enable_tls_verify=True,
                credentials_path=None,
            ),
            model=ModelConfig(intent=UpdateArtifactIntent(artifact_id="123")),
            storage=StorageConfig(path="/tmp/test"),
            registry=RegistryConfig(server_address="test-server")
        )

        _get_upload_params(config)

        call_kwargs = mock_get_backend.call_args
        pull_args = call_kwargs.kwargs.get("pull_args") or call_kwargs[1].get("pull_args")
        push_args = call_kwargs.kwargs.get("push_args") or call_kwargs[1].get("push_args")

        assert "--src-tls-verify=false" in pull_args
        assert "--authfile" in pull_args
        assert push_args == []


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
                endpoint="https://s3.amazonaws.com",
                access_key_id="test-access-key",
                secret_access_key="test-secret-key",
                region="us-east-1"
            ),
            model=ModelConfig(
                intent=UpdateArtifactIntent(artifact_id="test-artifact")
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
                intent=UpdateArtifactIntent(artifact_id="123")
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
                intent=UpdateArtifactIntent(artifact_id="123")
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
                intent=UpdateArtifactIntent(artifact_id="123")
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
                intent=UpdateArtifactIntent(artifact_id="test-artifact")
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
                intent=UpdateArtifactIntent(artifact_id="123")
            ),
            storage=StorageConfig(path="/tmp/test-model"),
            registry=RegistryConfig(server_address="test-server")
        )

        # Execute and verify exception is propagated
        with pytest.raises(Exception, match="Upload failed"):
            perform_upload(config)

        # Verify client method was called
        mock_save_to_oci_registry.assert_called_once()
