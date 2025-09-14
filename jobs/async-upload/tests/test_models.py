"""
Unit tests for models.py - specifically testing model validation logic.
"""
import pytest
from pydantic import ValidationError

from job.models import (
    ModelConfig,
    UploadIntent,
    CreateModelIntent,
    CreateVersionIntent,
    UpdateArtifactIntent,
    ConfigMapMetadata,
    RegisteredModelMetadata,
    ModelVersionMetadata,
    ModelArtifactMetadata,
)


class TestModelConfigIntentTypes:
    """Test cases for ModelConfig intent types using discriminated union"""

    def test_create_model_intent_with_no_ids_succeeds(self):
        """Test that create_model intent succeeds when no ids are provided"""
        intent = CreateModelIntent()
        config = ModelConfig(intent=intent)

        # Should not raise any validation errors
        assert config.intent.intent_type == UploadIntent.create_model
        assert isinstance(config.intent, CreateModelIntent)

    def test_create_model_intent_structure(self):
        """Test that CreateModelIntent only has the intent_type field"""
        intent = CreateModelIntent()
        config = ModelConfig(intent=intent)

        # CreateModelIntent should only have intent_type field
        assert config.intent.intent_type == UploadIntent.create_model
        assert not hasattr(config.intent, 'model_id')
        assert not hasattr(config.intent, 'artifact_id')

    def test_create_version_intent_with_required_ids_succeeds(self):
        """Test that create_version intent succeeds when model ID is provided"""
        intent = CreateVersionIntent(model_id="test-model-id")
        config = ModelConfig(intent=intent)

        # Should not raise any validation errors
        assert config.intent.intent_type == UploadIntent.create_version
        assert isinstance(config.intent, CreateVersionIntent)
        assert config.intent.model_id == "test-model-id"

    def test_create_version_intent_missing_model_id_fails(self):
        """Test that create_version intent fails when model ID is missing"""
        with pytest.raises(ValidationError) as exc_info:
            CreateVersionIntent()
            # model_id is a required field

        error_message = str(exc_info.value)
        assert "Field required" in error_message

    def test_update_artifact_intent_with_artifact_id_succeeds(self):
        """Test that update_artifact intent succeeds when artifact ID is provided"""
        intent = UpdateArtifactIntent(artifact_id="test-artifact-id")
        config = ModelConfig(intent=intent)

        # Should not raise any validation errors
        assert config.intent.intent_type == UploadIntent.update_artifact
        assert isinstance(config.intent, UpdateArtifactIntent)
        assert config.intent.artifact_id == "test-artifact-id"

    def test_update_artifact_intent_structure(self):
        """Test that UpdateArtifactIntent only has the required fields"""
        intent = UpdateArtifactIntent(artifact_id="test-artifact-id")
        config = ModelConfig(intent=intent)

        # UpdateArtifactIntent should only have intent_type and artifact_id fields
        assert config.intent.intent_type == UploadIntent.update_artifact
        assert config.intent.artifact_id == "test-artifact-id"
        assert not hasattr(config.intent, 'model_id')

    def test_update_artifact_intent_missing_artifact_id_fails(self):
        """Test that update_artifact intent fails when artifact ID is missing"""
        with pytest.raises(ValidationError) as exc_info:
            UpdateArtifactIntent()
            # artifact_id is a required field

        error_message = str(exc_info.value)
        assert "Field required" in error_message

    def test_discriminated_union_serialization(self):
        """Test that discriminated union serialization works correctly"""
        # Test CreateModelIntent
        create_intent = CreateModelIntent()
        config = ModelConfig(intent=create_intent)
        serialized = config.model_dump()
        assert serialized['intent']['intent_type'] == 'create_model'

        # Test CreateVersionIntent
        version_intent = CreateVersionIntent(model_id="test-id")
        config = ModelConfig(intent=version_intent)
        serialized = config.model_dump()
        assert serialized['intent']['intent_type'] == 'create_version'
        assert serialized['intent']['model_id'] == 'test-id'

        # Test UpdateArtifactIntent
        artifact_intent = UpdateArtifactIntent(artifact_id="artifact-id")
        config = ModelConfig(intent=artifact_intent)
        serialized = config.model_dump()
        assert serialized['intent']['intent_type'] == 'update_artifact'
        assert serialized['intent']['artifact_id'] == 'artifact-id'


class TestMetadataModels:

    def test_configmap_metadata_structure(self):
        metadata = ConfigMapMetadata(
            registered_model=RegisteredModelMetadata(
                name="test-model",
                description="A test model",
                owner="test-user"
            ),
            model_version=ModelVersionMetadata(
                name="1.0.0",
                description="Initial version",
                author="test-user"
            ),
            model_artifact=ModelArtifactMetadata(
                name="test-model-artifact",
                model_format_name="tensorflow",
                model_format_version="2.8"
            ),
        )

        assert metadata.registered_model.name == "test-model"
        assert metadata.model_version.name == "1.0.0"
        assert metadata.model_artifact.model_format_name == "tensorflow"

    def test_registered_model_metadata_requires_name_or_id(self):
        with pytest.raises(ValidationError, match="Must provide either name or id"):
            RegisteredModelMetadata()

    def test_registered_model_metadata_cannot_have_both_name_and_id(self):
        with pytest.raises(ValidationError, match="Cannot provide both name and id"):
            RegisteredModelMetadata(name="test", id="123")

    def test_registered_model_metadata_accepts_name_only(self):
        rm = RegisteredModelMetadata(name="test")
        assert rm.name == "test"
        assert rm.id is None

    def test_registered_model_metadata_accepts_id_only(self):
        rm = RegisteredModelMetadata(id="123")
        assert rm.id == "123"
        assert rm.name is None
