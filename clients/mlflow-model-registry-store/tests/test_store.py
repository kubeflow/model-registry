"""Tests for ModelRegistryStore."""

import pytest
from unittest.mock import Mock, patch

from mlflow.exceptions import MlflowException
from mlflow.entities.model_registry import RegisteredModel, ModelVersion

from mlflow_model_registry_store.store import ModelRegistryStore


class TestModelRegistryStore:
    """Test cases for ModelRegistryStore."""

    def test_init_with_valid_uri(self):
        """Test store initialization with valid URI."""
        uri = "modelregistry://localhost:8080?author=testuser&is-secure=false"

        with patch('mlflow_model_registry_store.store.ModelRegistry') as mock_registry:
            store = ModelRegistryStore(store_uri=uri)

            assert store.server_address == "http://localhost"
            assert store.port == 8080
            assert store.author == "testuser"
            assert store.is_secure is False

            # Verify ModelRegistry was called with correct parameters
            mock_registry.assert_called_once_with(
                server_address="http://localhost",
                port=8080,
                author="testuser",
                is_secure=False,
                user_token=None,
                custom_ca=None,
            )

    def test_init_with_invalid_uri(self):
        """Test store initialization with invalid URI."""
        invalid_uri = "invalid://localhost:8080"

        with pytest.raises(MlflowException) as exc_info:
            ModelRegistryStore(store_uri=invalid_uri)

        assert "Invalid store URI" in str(exc_info.value)

    def test_init_with_missing_uri(self):
        """Test store initialization with missing URI."""
        with pytest.raises(MlflowException) as exc_info:
            ModelRegistryStore(store_uri=None)

        assert "Store URI is required" in str(exc_info.value)

    @patch('mlflow_model_registry_store.store.ModelRegistry')
    def test_convert_mr_to_mlflow_registered_model(self, mock_registry):
        """Test conversion from ModelRegistry RegisteredModel to MLflow RegisteredModel."""
        uri = "modelregistry://localhost:8080?author=testuser&is-secure=false"
        store = ModelRegistryStore(store_uri=uri)

        # Mock Model Registry RegisteredModel
        mr_model = Mock()
        mr_model.name = "test-model"
        mr_model.description = "Test model description"
        mr_model.create_time_since_epoch = 1234567890
        mr_model.last_update_time_since_epoch = 1234567900
        mr_model.custom_properties = {"key1": "value1", "key2": "value2"}

        # Convert to MLflow format
        mlflow_model = store._convert_mr_to_mlflow_registered_model(mr_model)

        assert isinstance(mlflow_model, RegisteredModel)
        assert mlflow_model.name == "test-model"
        assert mlflow_model.description == "Test model description"
        assert mlflow_model.creation_timestamp == 1234567890
        assert mlflow_model.last_updated_timestamp == 1234567900
        assert len(mlflow_model.tags) == 2

    @patch('mlflow_model_registry_store.store.ModelRegistry')
    def test_convert_mr_to_mlflow_model_version(self, mock_registry):
        """Test conversion from ModelRegistry ModelVersion to MLflow ModelVersion."""
        uri = "modelregistry://localhost:8080?author=testuser&is-secure=false"
        store = ModelRegistryStore(store_uri=uri)

        # Mock Model Registry ModelVersion
        mr_version = Mock()
        mr_version.name = "v1"
        mr_version.description = "Version 1"
        mr_version.create_time_since_epoch = 1234567890
        mr_version.last_update_time_since_epoch = 1234567900
        mr_version.custom_properties = {"stage": "production"}

        # Mock artifact
        mr_artifact = Mock()
        mr_artifact.uri = "s3://bucket/model"

        # Convert to MLflow format
        mlflow_version = store._convert_mr_to_mlflow_model_version(mr_version, mr_artifact)

        assert isinstance(mlflow_version, ModelVersion)
        assert mlflow_version.version == "v1"
        assert mlflow_version.description == "Version 1"
        assert mlflow_version.source == "s3://bucket/model"
        assert mlflow_version.status == "READY"

    @patch('mlflow_model_registry_store.store.ModelRegistry')
    def test_create_registered_model_success(self, mock_registry):
        """Test successful registered model creation."""
        uri = "modelregistry://localhost:8080?author=testuser&is-secure=false"
        store = ModelRegistryStore(store_uri=uri)

        # Mock the client methods
        store._client.get_registered_model.return_value = None  # Model doesn't exist

        mock_mr_model = Mock()
        mock_mr_model.name = "test-model"
        mock_mr_model.description = "Test description"
        mock_mr_model.create_time_since_epoch = 1234567890
        mock_mr_model.last_update_time_since_epoch = 1234567890
        mock_mr_model.custom_properties = {}

        store._client.async_runner.return_value = mock_mr_model

        # Create registered model
        result = store.create_registered_model("test-model", description="Test description")

        assert isinstance(result, RegisteredModel)
        assert result.name == "test-model"

    @patch('mlflow_model_registry_store.store.ModelRegistry')
    def test_create_registered_model_already_exists(self, mock_registry):
        """Test registered model creation when model already exists."""
        uri = "modelregistry://localhost:8080?author=testuser&is-secure=false"
        store = ModelRegistryStore(store_uri=uri)

        # Mock existing model
        existing_model = Mock()
        store._client.get_registered_model.return_value = existing_model

        # Should raise exception
        with pytest.raises(MlflowException) as exc_info:
            store.create_registered_model("test-model")

        assert "already exists" in str(exc_info.value)

    @patch('mlflow_model_registry_store.store.ModelRegistry')
    def test_get_registered_model(self, mock_registry):
        """Test getting a registered model."""
        uri = "modelregistry://localhost:8080?author=testuser&is-secure=false"
        store = ModelRegistryStore(store_uri=uri)

        # Mock the client method
        mock_mr_model = Mock()
        mock_mr_model.name = "test-model"
        mock_mr_model.description = "Test description"
        mock_mr_model.create_time_since_epoch = 1234567890
        mock_mr_model.last_update_time_since_epoch = 1234567890
        mock_mr_model.custom_properties = {}

        store._client.get_registered_model.return_value = mock_mr_model

        # Get registered model
        result = store.get_registered_model("test-model")

        assert isinstance(result, RegisteredModel)
        assert result.name == "test-model"
        store._client.get_registered_model.assert_called_once_with("test-model")

    @patch('mlflow_model_registry_store.store.ModelRegistry')
    def test_get_registered_model_not_found(self, mock_registry):
        """Test getting a registered model that doesn't exist."""
        uri = "modelregistry://localhost:8080?author=testuser&is-secure=false"
        store = ModelRegistryStore(store_uri=uri)

        # Mock the client method to return None
        store._client.get_registered_model.return_value = None

        # Get registered model should raise exception
        with pytest.raises(MlflowException) as exc_info:
            store.get_registered_model("nonexistent-model")

        assert "not found" in str(exc_info.value)

    @patch('mlflow_model_registry_store.store.ModelRegistry')
    def test_rename_registered_model_not_supported(self, mock_registry):
        """Test that renaming registered models is not supported."""
        uri = "modelregistry://localhost:8080?author=testuser&is-secure=false"
        store = ModelRegistryStore(store_uri=uri)

        with pytest.raises(MlflowException) as exc_info:
            store.rename_registered_model("old-name", "new-name")

        assert "not supported" in str(exc_info.value)

    @patch('mlflow_model_registry_store.store.ModelRegistry')
    def test_delete_registered_model_archives_model_and_versions(self, mock_registry):
        """Test that deleting registered models archives the model and its versions."""
        from model_registry.types import RegisteredModelState, ModelVersionState

        uri = "modelregistry://localhost:8080?author=testuser&is-secure=false"
        store = ModelRegistryStore(store_uri=uri)

        # Mock existing model
        mock_model = Mock()
        mock_model.name = "test-model"
        mock_model.state = RegisteredModelState.LIVE
        store._client.get_registered_model.return_value = mock_model

        # Mock model versions
        mock_version1 = Mock()
        mock_version1.name = "1"
        mock_version1.state = ModelVersionState.LIVE

        mock_version2 = Mock()
        mock_version2.name = "2"
        mock_version2.state = ModelVersionState.LIVE

        # Mock the pager to return versions
        mock_pager = Mock()
        mock_pager.__iter__ = Mock(return_value=iter([mock_version1, mock_version2]))
        store._client.get_model_versions.return_value = mock_pager

        # Delete the model
        store.delete_registered_model("test-model")

        # Verify versions were archived
        assert mock_version1.state == ModelVersionState.ARCHIVED
        assert mock_version2.state == ModelVersionState.ARCHIVED

        # Verify model was archived
        assert mock_model.state == RegisteredModelState.ARCHIVED

        # Verify update was called for all objects (2 versions + 1 model)
        assert store._client.update.call_count == 3

    @patch('mlflow_model_registry_store.store.ModelRegistry')
    def test_delete_model_version_archives_version(self, mock_registry):
        """Test that deleting model version archives the specific version."""
        from model_registry.types import ModelVersionState

        uri = "modelregistry://localhost:8080?author=testuser&is-secure=false"
        store = ModelRegistryStore(store_uri=uri)

        # Mock model version
        mock_version = Mock()
        mock_version.name = "1"
        mock_version.state = ModelVersionState.LIVE
        store._client.get_model_version.return_value = mock_version

        # Delete the version
        store.delete_model_version("test-model", "1")

        # Verify version was archived
        assert mock_version.state == ModelVersionState.ARCHIVED

        # Verify update was called
        store._client.update.assert_called_once_with(mock_version)