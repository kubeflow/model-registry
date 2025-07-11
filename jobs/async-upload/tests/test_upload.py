from dataclasses import asdict
import pytest
from unittest.mock import Mock, patch
from model_registry.utils import S3Params, OCIParams
from job.upload import _get_upload_params, perform_upload

class TestGetUploadParams:
    """Test cases for _get_upload_params function"""

    def test_get_upload_params_s3_config(self):
        """Test _get_upload_params with S3 configuration returns S3Params"""
        config = {
            "destination": {
                "type": "s3",
                "s3": {
                    "bucket": "test-bucket",
                    "key": "test-key",
                    "endpoint": "https://s3.amazonaws.com",
                    "access_key_id": "test-access-key",
                    "secret_access_key": "test-secret-key",
                    "region": "us-east-1",
                },
            }
        }

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
        config = {
            "destination": {
                "type": "oci",
                "oci": {
                    "uri": "quay.io/example/test:latest",
                    "username": "test-user",
                    "password": "test-password",
                },
            },
            "model": {"id": "abc", "version_id": "def", "artifact_id": "123"},
            "storage": {"path": "/tmp/test-model"},
        }

        result = _get_upload_params(config)

        assert isinstance(result, OCIParams)
        assert result.base_image == "123"
        assert result.oci_ref == "quay.io/example/test:latest"
        assert result.dest_dir == "/tmp/test-model"
        assert result.oci_username == "test-user"
        assert result.oci_password == "test-password"

    def test_get_upload_params_unsupported_type(self):
        """Test _get_upload_params with unsupported destination type raises ValueError"""
        config = {"destination": {"type": "ftp"}}

        with pytest.raises(ValueError, match="Unsupported destination type: ftp"):
            _get_upload_params(config)

    def test_get_upload_params_s3_with_none_values(self):
        """Test _get_upload_params with S3 config containing None values"""
        config = {
            "destination": {
                "type": "s3",
                "s3": {
                    "bucket": "test-bucket",
                    "key": "test-key",
                    "endpoint": None,
                    "access_key_id": None,
                    "secret_access_key": None,
                    "region": None,
                },
            }
        }

        result = _get_upload_params(config)

        assert isinstance(result, S3Params)
        assert result.bucket_name == "test-bucket"
        assert result.s3_prefix == "test-key"
        assert result.endpoint_url is None
        assert result.access_key_id is None
        assert result.secret_access_key is None
        assert result.region is None

    def test_get_upload_params_oci_with_none_values(self):
        """Test _get_upload_params with OCI config containing None values"""
        config = {
            "destination": {
                "type": "oci",
                "oci": {
                    "uri": "quay.io/example/test:latest",
                    "username": None,
                    "password": None,
                },
            },
            "model": {"id": "abc", "version_id": "def", "artifact_id": "123"},
            "storage": {"path": "/tmp/test-model"},
        }

        result = _get_upload_params(config)

        assert isinstance(result, OCIParams)
        assert result.base_image == "123"
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

        config = {
            "destination": {"type": "oci", "oci": {"uri": "quay.io/example/oci", "username": "oci_user", "password": "oci_pass"}},
            "storage": {"path": "/tmp/test-model"},
            "model": {
                "id": "abc",
                "version_id": "def",
                "artifact_id": "123",
                "format_version": "1.16",
            },
        }

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

        config = {
            "destination": {"type": "s3"},
            "storage": {"path": "/tmp/test-model"},
            "model": {
                "name": "test-model",
                "version": "1.0.0",
                "format": "onnx",
                "format_version": "1.16",
            },
        }

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

        config = {
            "destination": {"type": "oci", "oci": {"uri": "quay.io/example/oci", "username": "oci_user", "password": "oci_pass"}},
            "storage": {"path": "/tmp/test-model"},
            "model": {
                "id": "abc",
                "version_id": "def",
                "artifact_id": "123",
                "format_version": "1.16",
            },
        }

        # Execute and verify exception is propagated
        with pytest.raises(Exception, match="Upload failed"):
            perform_upload(config)

        # Verify client method was called
        mock_save_to_oci_registry.assert_called_once()
