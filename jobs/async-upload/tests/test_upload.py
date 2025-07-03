import pytest
from unittest.mock import Mock, patch, MagicMock
from pathlib import Path
import shutil
import tempfile
import os
from model_registry.utils import S3Params, OCIParams
from job.upload import _get_upload_params, perform_upload, _prepare_modelcar_structure


class TestPrepareModelcarStructure:
    """Test cases for _prepare_modelcar_structure function"""

    def test_prepare_modelcar_structure_single_file(self):
        """Test _prepare_modelcar_structure with a single model file"""
        with tempfile.TemporaryDirectory() as temp_dir:
            # Create a test model file
            model_file = Path(temp_dir) / "model.onnx"
            model_file.write_text("dummy model content")

            config = {"model": {"name": "test-model"}}

            # Execute
            result = _prepare_modelcar_structure(config, str(model_file))

            # Verify
            result_path = Path(result)
            assert result_path.exists()
            assert result_path.name == "modelcar"

            models_dir = result_path / "models"
            assert models_dir.exists()
            assert models_dir.is_dir()

            copied_file = models_dir / "model.onnx"
            assert copied_file.exists()
            assert copied_file.read_text() == "dummy model content"

    def test_prepare_modelcar_structure_directory(self):
        """Test _prepare_modelcar_structure with a directory of model files"""
        with tempfile.TemporaryDirectory() as temp_dir:
            # Create a test model directory with multiple files
            model_dir = Path(temp_dir) / "model_files"
            model_dir.mkdir()

            # Create test files
            (model_dir / "model.onnx").write_text("model content")
            (model_dir / "config.json").write_text('{"version": "1.0"}')

            # Create subdirectory
            subdir = model_dir / "assets"
            subdir.mkdir()
            (subdir / "weights.bin").write_text("weights content")

            config = {"model": {"name": "test-model"}}

            # Execute
            result = _prepare_modelcar_structure(config, str(model_dir))

            # Verify
            result_path = Path(result)
            assert result_path.exists()
            assert result_path.name == "modelcar"

            models_dir = result_path / "models"
            assert models_dir.exists()
            assert models_dir.is_dir()

            # Check files were copied
            assert (models_dir / "model.onnx").exists()
            assert (models_dir / "model.onnx").read_text() == "model content"
            assert (models_dir / "config.json").exists()
            assert (models_dir / "config.json").read_text() == '{"version": "1.0"}'

            # Check subdirectory was copied
            assert (models_dir / "assets").exists()
            assert (models_dir / "assets").is_dir()
            assert (models_dir / "assets" / "weights.bin").exists()
            assert (
                models_dir / "assets" / "weights.bin"
            ).read_text() == "weights content"

    def test_prepare_modelcar_structure_cleanup_existing(self):
        """Test _prepare_modelcar_structure cleans up existing modelcar directory"""
        with tempfile.TemporaryDirectory() as temp_dir:
            # Create a test model file
            model_file = Path(temp_dir) / "model.onnx"
            model_file.write_text("new model content")

            # Create existing modelcar directory
            existing_modelcar = Path(temp_dir) / "modelcar"
            existing_modelcar.mkdir()
            (existing_modelcar / "old_file.txt").write_text("old content")

            config = {"model": {"name": "test-model"}}

            # Execute
            result = _prepare_modelcar_structure(config, str(model_file))

            # Verify old content is gone and new content is there
            result_path = Path(result)
            assert result_path.exists()
            assert not (result_path / "old_file.txt").exists()

            models_dir = result_path / "models"
            assert (models_dir / "model.onnx").exists()
            assert (models_dir / "model.onnx").read_text() == "new model content"

    def test_prepare_modelcar_structure_empty_directory(self):
        """Test _prepare_modelcar_structure with empty directory"""
        with tempfile.TemporaryDirectory() as temp_dir:
            # Create empty model directory
            model_dir = Path(temp_dir) / "empty_model"
            model_dir.mkdir()

            config = {"model": {"name": "test-model"}}

            # Execute
            result = _prepare_modelcar_structure(config, str(model_dir))

            # Verify structure is created but models directory is empty
            result_path = Path(result)
            assert result_path.exists()

            models_dir = result_path / "models"
            assert models_dir.exists()
            assert models_dir.is_dir()
            assert len(list(models_dir.iterdir())) == 0

    def test_prepare_modelcar_structure_nonexistent_path(self):
        """Test _prepare_modelcar_structure with nonexistent path"""
        with tempfile.TemporaryDirectory() as temp_dir:
            nonexistent_path = Path(temp_dir) / "nonexistent"
            config = {"model": {"name": "test-model"}}

            # Execute
            result = _prepare_modelcar_structure(config, str(nonexistent_path))

            # Verify structure is created even if source doesn't exist
            result_path = Path(result)
            assert result_path.exists()

            models_dir = result_path / "models"
            assert models_dir.exists()
            assert models_dir.is_dir()
            assert len(list(models_dir.iterdir())) == 0

    def test_prepare_modelcar_structure_path_objects(self):
        """Test _prepare_modelcar_structure handles Path objects correctly"""
        with tempfile.TemporaryDirectory() as temp_dir:
            # Create a test model file
            model_file = Path(temp_dir) / "model.onnx"
            model_file.write_text("dummy model content")

            config = {"model": {"name": "test-model"}}

            # Execute with Path object
            result = _prepare_modelcar_structure(config, str(model_file))

            # Verify
            result_path = Path(result)
            assert result_path.exists()
            assert isinstance(result, str)  # Should return string path

            models_dir = result_path / "models"
            assert (models_dir / "model.onnx").exists()


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
            "model": {"name": "test-model"},
            "storage": {"path": "/tmp/test-model"},
        }

        result = _get_upload_params(config)

        assert isinstance(result, OCIParams)
        assert result.base_image == "test-model"
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
            "model": {"name": "test-model"},
            "storage": {"path": "/tmp/test-model"},
        }

        result = _get_upload_params(config)

        assert isinstance(result, OCIParams)
        assert result.base_image == "test-model"
        assert result.oci_ref == "quay.io/example/test:latest"
        assert result.dest_dir == "/tmp/test-model"
        assert result.oci_username is None
        assert result.oci_password is None


class TestPerformUpload:
    """Test cases for perform_upload function"""

    @patch("job.upload._get_upload_params")
    def test_perform_upload_s3_calls_client_with_correct_params(
        self, mock_get_upload_params
    ):
        """Test perform_upload calls client.upload_artifact_and_register_model with correct parameters for S3"""
        # Setup
        mock_client = Mock()
        mock_upload_params = Mock()
        mock_get_upload_params.return_value = mock_upload_params

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

        # Execute
        perform_upload(mock_client, config)

        # Verify _get_upload_params was called with config
        mock_get_upload_params.assert_called_once_with(config)

        # Verify client method was called with correct parameters
        mock_client.upload_artifact_and_register_model.assert_called_once_with(
            model_files_path="/tmp/test-model",
            name="test-model",
            version="1.0.0",
            model_format_name="onnx",
            model_format_version="1.16",
            upload_params=mock_upload_params,
        )

    @patch("job.upload._get_upload_params")
    @patch("job.upload._prepare_modelcar_structure")
    @patch("os.path.exists")
    def test_perform_upload_oci_with_modelcar_preparation(
        self, mock_exists, mock_prepare_modelcar, mock_get_upload_params
    ):
        """Test perform_upload with OCI destination prepares modelcar structure"""
        # Setup
        mock_client = Mock()
        mock_upload_params = Mock()
        mock_get_upload_params.return_value = mock_upload_params
        mock_exists.return_value = True
        mock_prepare_modelcar.return_value = "/tmp/modelcar"

        config = {
            "destination": {"type": "oci"},
            "storage": {"path": "/tmp/test-model"},
            "model": {
                "name": "test-model",
                "version": "1.0.0",
                "format": "onnx",
                "format_version": "1.16",
            },
        }

        # Execute
        perform_upload(mock_client, config)

        # Verify modelcar preparation was called
        mock_prepare_modelcar.assert_called_once_with(config, "/tmp/test-model")
        mock_exists.assert_called_once_with("/tmp/test-model")

        # Verify client method was called with modelcar path
        mock_client.upload_artifact_and_register_model.assert_called_once_with(
            model_files_path="/tmp/modelcar",
            name="test-model",
            version="1.0.0",
            model_format_name="onnx",
            model_format_version="1.16",
            upload_params=mock_upload_params,
        )

    @patch("job.upload._get_upload_params")
    @patch("job.upload._prepare_modelcar_structure")
    @patch("os.path.exists")
    def test_perform_upload_oci_without_existing_files(
        self, mock_exists, mock_prepare_modelcar, mock_get_upload_params
    ):
        """Test perform_upload with OCI destination when model files don't exist"""
        # Setup
        mock_client = Mock()
        mock_upload_params = Mock()
        mock_get_upload_params.return_value = mock_upload_params
        mock_exists.return_value = False

        config = {
            "destination": {"type": "oci"},
            "storage": {"path": "/tmp/test-model"},
            "model": {
                "name": "test-model",
                "version": "1.0.0",
                "format": "onnx",
                "format_version": "1.16",
            },
        }

        # Execute
        perform_upload(mock_client, config)

        # Verify modelcar preparation was NOT called
        mock_prepare_modelcar.assert_not_called()
        mock_exists.assert_called_once_with("/tmp/test-model")

        # Verify client method was called with original path
        mock_client.upload_artifact_and_register_model.assert_called_once_with(
            model_files_path="/tmp/test-model",
            name="test-model",
            version="1.0.0",
            model_format_name="onnx",
            model_format_version="1.16",
            upload_params=mock_upload_params,
        )

    @patch("job.upload._get_upload_params")
    def test_perform_upload_with_s3_params(self, mock_get_upload_params):
        """Test perform_upload integration with S3 parameters"""
        # Setup
        mock_client = Mock()
        s3_params = S3Params(
            bucket_name="test-bucket",
            s3_prefix="test-key",
            endpoint_url="https://s3.amazonaws.com",
            access_key_id="test-access-key",
            secret_access_key="test-secret-key",
            region="us-east-1",
        )
        mock_get_upload_params.return_value = s3_params

        config = {
            "destination": {"type": "s3"},
            "storage": {"path": "/tmp/test-model"},
            "model": {
                "name": "test-model",
                "version": "2.0.0",
                "format": "pytorch",
                "format_version": "1.12",
            },
        }

        # Execute
        perform_upload(mock_client, config)

        # Verify
        mock_client.upload_artifact_and_register_model.assert_called_once_with(
            model_files_path="/tmp/test-model",
            name="test-model",
            version="2.0.0",
            model_format_name="pytorch",
            model_format_version="1.12",
            upload_params=s3_params,
        )

    @patch("job.upload._get_upload_params")
    @patch("job.upload._prepare_modelcar_structure")
    @patch("os.path.exists")
    def test_perform_upload_with_oci_params(
        self, mock_exists, mock_prepare_modelcar, mock_get_upload_params
    ):
        """Test perform_upload integration with OCI parameters"""
        # Setup
        mock_client = Mock()
        oci_params = OCIParams(
            base_image="test-model",
            oci_ref="quay.io/example/test:latest",
            dest_dir="/tmp/test-model",
            oci_username="test-user",
            oci_password="test-password",
        )
        mock_get_upload_params.return_value = oci_params
        mock_exists.return_value = True
        mock_prepare_modelcar.return_value = "/tmp/modelcar"

        config = {
            "destination": {"type": "oci"},
            "storage": {"path": "/tmp/test-model"},
            "model": {
                "name": "test-model",
                "version": "3.0.0",
                "format": "tensorflow",
                "format_version": "2.8",
            },
        }

        # Execute
        perform_upload(mock_client, config)

        # Verify
        mock_client.upload_artifact_and_register_model.assert_called_once_with(
            model_files_path="/tmp/modelcar",
            name="test-model",
            version="3.0.0",
            model_format_name="tensorflow",
            model_format_version="2.8",
            upload_params=oci_params,
        )

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
            perform_upload(mock_client, config)

        # Verify client method was not called
        mock_client.upload_artifact_and_register_model.assert_not_called()

    @patch("job.upload._get_upload_params")
    def test_perform_upload_propagates_exceptions_from_client(
        self, mock_get_upload_params
    ):
        """Test perform_upload propagates exceptions from client method"""
        # Setup
        mock_client = Mock()
        mock_upload_params = Mock()
        mock_get_upload_params.return_value = mock_upload_params
        mock_client.upload_artifact_and_register_model.side_effect = Exception(
            "Upload failed"
        )

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
        with pytest.raises(Exception, match="Upload failed"):
            perform_upload(mock_client, config)

        # Verify client method was called
        mock_client.upload_artifact_and_register_model.assert_called_once()

    @patch("job.upload._get_upload_params")
    @patch("job.upload._prepare_modelcar_structure")
    @patch("os.path.exists")
    def test_perform_upload_propagates_exceptions_from_prepare_modelcar(
        self, mock_exists, mock_prepare_modelcar, mock_get_upload_params
    ):
        """Test perform_upload propagates exceptions from _prepare_modelcar_structure"""
        # Setup
        mock_client = Mock()
        mock_upload_params = Mock()
        mock_get_upload_params.return_value = mock_upload_params
        mock_exists.return_value = True
        mock_prepare_modelcar.side_effect = OSError("Permission denied")

        config = {
            "destination": {"type": "oci"},
            "storage": {"path": "/tmp/test-model"},
            "model": {
                "name": "test-model",
                "version": "1.0.0",
                "format": "onnx",
                "format_version": "1.16",
            },
        }

        # Execute and verify exception is propagated
        with pytest.raises(OSError, match="Permission denied"):
            perform_upload(mock_client, config)

        # Verify client method was not called
        mock_client.upload_artifact_and_register_model.assert_not_called()
