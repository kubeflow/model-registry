"""
Unit tests for models.py - specifically testing model validation logic.
"""
import pytest
from pydantic import ValidationError

from job.models import ModelConfig, UploadIntent


class TestModelConfigValidateModelIds:
    """Test cases for ModelConfig.validate_model_ids method"""

    def test_create_model_intent_with_all_ids_succeeds(self):
        """Test that create_model intent succeeds when all required IDs are provided"""
        config = ModelConfig(
            upload_intent=UploadIntent.create_model,
            id="test-model-id",
            version_id="test-version-id",
            artifact_id="test-artifact-id"
        )
        
        # Should not raise any validation errors
        assert config.upload_intent == UploadIntent.create_model
        assert config.id == "test-model-id"
        assert config.version_id == "test-version-id"
        assert config.artifact_id == "test-artifact-id"

    def test_create_model_intent_no_validation(self):
        """Test that create_model intent has no validation requirements"""
        # Should succeed even with empty strings since create_model has no validation
        config = ModelConfig(
            upload_intent=UploadIntent.create_model,
            id="",  # Empty string - should be fine
            version_id="",  # Empty string - should be fine
            artifact_id=""  # Empty string - should be fine
        )
        
        # Should not raise any validation errors
        assert config.upload_intent == UploadIntent.create_model
        assert config.id == ""
        assert config.version_id == ""
        assert config.artifact_id == ""

    def test_create_model_intent_with_partial_data_succeeds(self):
        """Test that create_model intent succeeds with partial data"""
        config = ModelConfig(
            upload_intent=UploadIntent.create_model,
            id="test-model-id",
            version_id="",  # Empty string - should be fine
            artifact_id="test-artifact-id"
        )
        
        # Should not raise any validation errors
        assert config.upload_intent == UploadIntent.create_model
        assert config.id == "test-model-id"
        assert config.version_id == ""
        assert config.artifact_id == "test-artifact-id"

    def test_create_model_intent_with_mixed_data_succeeds(self):
        """Test that create_model intent succeeds with mixed data"""
        config = ModelConfig(
            upload_intent=UploadIntent.create_model,
            id="test-model-id",
            version_id="test-version-id",
            artifact_id=""  # Empty string - should be fine
        )
        
        # Should not raise any validation errors
        assert config.upload_intent == UploadIntent.create_model
        assert config.id == "test-model-id"
        assert config.version_id == "test-version-id"
        assert config.artifact_id == ""

    def test_create_model_intent_allows_all_empty_ids(self):
        """Test that create_model intent allows all empty IDs"""
        config = ModelConfig(
            upload_intent=UploadIntent.create_model,
            id="",
            version_id="",
            artifact_id=""
        )
        
        # Should not raise any validation errors
        assert config.upload_intent == UploadIntent.create_model
        assert config.id == ""
        assert config.version_id == ""
        assert config.artifact_id == ""

    def test_create_version_intent_with_required_ids_succeeds(self):
        """Test that create_version intent succeeds when model ID and version ID are provided"""
        config = ModelConfig(
            upload_intent=UploadIntent.create_version,
            id="test-model-id",
            version_id="test-version-id",
            artifact_id="test-artifact-id"  # Can be provided but not required
        )
        
        # Should not raise any validation errors
        assert config.upload_intent == UploadIntent.create_version
        assert config.id == "test-model-id"
        assert config.version_id == "test-version-id"
        assert config.artifact_id == "test-artifact-id"

    def test_create_version_intent_without_artifact_id_succeeds(self):
        """Test that create_version intent succeeds without artifact ID"""
        config = ModelConfig(
            upload_intent=UploadIntent.create_version,
            id="test-model-id",
            version_id="test-version-id",
            artifact_id="any-value"  # Not validated for create_version
        )
        
        # Should not raise any validation errors
        assert config.upload_intent == UploadIntent.create_version
        assert config.id == "test-model-id"
        assert config.version_id == "test-version-id"

    def test_create_version_intent_missing_model_id_fails(self):
        """Test that create_version intent fails when model ID is missing"""
        with pytest.raises(ValidationError) as exc_info:
            ModelConfig(
                upload_intent=UploadIntent.create_version,
                id="",  # Empty string
                version_id="test-version-id",
                artifact_id="test-artifact-id"
            )
        
        error_message = str(exc_info.value)
        assert "Model ID must be set when intent is create_version" in error_message

    def test_create_version_intent_without_version_id_succeeds(self):
        """Test that create_version intent succeeds without version ID (only model ID required)"""
        config = ModelConfig(
            upload_intent=UploadIntent.create_version,
            id="test-model-id",
            version_id="",  # Empty string - should be fine
            artifact_id="test-artifact-id"
        )
        
        # Should not raise any validation errors since only model ID is required
        assert config.upload_intent == UploadIntent.create_version
        assert config.id == "test-model-id"
        assert config.version_id == ""
        assert config.artifact_id == "test-artifact-id"

    def test_create_version_intent_missing_model_id_only_fails(self):
        """Test that create_version intent fails only when model ID is missing"""
        with pytest.raises(ValidationError) as exc_info:
            ModelConfig(
                upload_intent=UploadIntent.create_version,
                id="",  # Empty string - this should fail
                version_id="",  # Empty string - this is fine
                artifact_id="test-artifact-id"
            )
        
        error_message = str(exc_info.value)
        assert "Model ID must be set when intent is create_version" in error_message

    def test_update_artifact_intent_with_artifact_id_succeeds(self):
        """Test that update_artifact intent succeeds when artifact ID is provided"""
        config = ModelConfig(
            upload_intent=UploadIntent.update_artifact,
            id="test-model-id",  # Can be provided but not required
            version_id="test-version-id",  # Can be provided but not required
            artifact_id="test-artifact-id"
        )
        
        # Should not raise any validation errors
        assert config.upload_intent == UploadIntent.update_artifact
        assert config.id == "test-model-id"
        assert config.version_id == "test-version-id"
        assert config.artifact_id == "test-artifact-id"

    def test_update_artifact_intent_without_model_and_version_ids_succeeds(self):
        """Test that update_artifact intent succeeds without model ID and version ID"""
        config = ModelConfig(
            upload_intent=UploadIntent.update_artifact,
            id="any-value",  # Not validated for update_artifact
            version_id="any-value",  # Not validated for update_artifact
            artifact_id="test-artifact-id"
        )
        
        # Should not raise any validation errors
        assert config.upload_intent == UploadIntent.update_artifact
        assert config.artifact_id == "test-artifact-id"

    def test_update_artifact_intent_missing_artifact_id_fails(self):
        """Test that update_artifact intent fails when artifact ID is missing"""
        with pytest.raises(ValidationError) as exc_info:
            ModelConfig(
                upload_intent=UploadIntent.update_artifact,
                id="test-model-id",
                version_id="test-version-id",
                artifact_id=""  # Empty string
            )
        
        error_message = str(exc_info.value)
        assert "Model Artifact ID must be set when intent is update_artifact" in error_message

    def test_default_intent_is_update_artifact(self):
        """Test that the default upload intent is update_artifact"""
        config = ModelConfig(
            id="test-model-id",
            version_id="test-version-id",
            artifact_id="test-artifact-id"
        )
        
        # Default should be update_artifact
        assert config.upload_intent == UploadIntent.update_artifact

    def test_default_intent_validation_with_valid_artifact_id(self):
        """Test that default intent (update_artifact) validates correctly with artifact ID"""
        config = ModelConfig(
            id="test-model-id",
            version_id="test-version-id",
            artifact_id="test-artifact-id"
        )
        
        # Should not raise any validation errors with default intent
        assert config.upload_intent == UploadIntent.update_artifact
        assert config.artifact_id == "test-artifact-id"

    def test_default_intent_validation_fails_without_artifact_id(self):
        """Test that default intent (update_artifact) validation fails without artifact ID"""
        with pytest.raises(ValidationError) as exc_info:
            ModelConfig(
                id="test-model-id",
                version_id="test-version-id",
                artifact_id=""  # Empty string
            )
        
        error_message = str(exc_info.value)
        assert "Model Artifact ID must be set when intent is update_artifact" in error_message

    def test_pydantic_field_validation_prevents_none_values(self):
        """Test that Pydantic field validation prevents None values for required fields"""
        # Pydantic should prevent None values due to Field(...) constraints
        with pytest.raises(ValidationError):
            ModelConfig(
                upload_intent=UploadIntent.update_artifact,
                id="test-model-id",
                version_id="test-version-id",
                artifact_id=None  # This should be caught by Pydantic field validation
            )

    def test_whitespace_only_strings_behavior(self):
        """Test behavior with whitespace-only strings for different intents"""
        # create_model: No validation, so whitespace is fine
        config1 = ModelConfig(
            upload_intent=UploadIntent.create_model,
            id="   ",  # Whitespace only - fine for create_model
            version_id="   ",
            artifact_id="   "
        )
        assert config1.upload_intent == UploadIntent.create_model
        
        # create_version: Only validates model ID, whitespace-only string is truthy
        config2 = ModelConfig(
            upload_intent=UploadIntent.create_version,
            id="   ",  # Whitespace only - truthy, so validation passes
            version_id="",
            artifact_id=""
        )
        assert config2.upload_intent == UploadIntent.create_version
        
        # update_artifact: Only validates artifact ID, whitespace-only string is truthy
        config3 = ModelConfig(
            upload_intent=UploadIntent.update_artifact,
            id="",
            version_id="",
            artifact_id="   "  # Whitespace only - truthy, so validation passes
        )
        assert config3.upload_intent == UploadIntent.update_artifact
