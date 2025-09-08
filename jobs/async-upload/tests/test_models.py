"""
Unit tests for models.py - specifically testing model validation logic.
"""
import pytest
from pydantic import ValidationError

from job.models import ModelConfig, UploadIntent


class TestModelConfigValidateModelIds:
    """Test cases for ModelConfig.validate_model_ids method"""

    def test_create_model_intent_with_no_ids_succeeds(self):
        """Test that create_model intent succeeds when no ids are provided"""
        config = ModelConfig(
            upload_intent=UploadIntent.create_model,
        )
        
        # Should not raise any validation errors
        assert config.upload_intent == UploadIntent.create_model

    def test_create_model_intent_with_any_id_fails(self):
        """Test that create_model intent fails when any id is provided"""
        expected_error_message = "Model ID, Model Version ID and Model Artifact ID cannot be set when intent is create_model"
        
        with pytest.raises(ValidationError) as exc_info:
            ModelConfig(
                upload_intent=UploadIntent.create_model,
                id="fail",
            )
        assert expected_error_message in str(exc_info.value)

        with pytest.raises(ValidationError) as exc_info:
            ModelConfig(
                upload_intent=UploadIntent.create_model,
                version_id="fail",
            )
        assert expected_error_message in str(exc_info.value)

        with pytest.raises(ValidationError) as exc_info:
            ModelConfig(
                upload_intent=UploadIntent.create_model,
                artifact_id="fail",
            )
        assert expected_error_message in str(exc_info.value)

    def test_create_version_intent_with_required_ids_succeeds(self):
        """Test that create_version intent succeeds when model ID and version ID are provided"""
        config = ModelConfig(
            upload_intent=UploadIntent.create_version,
            id="test-model-id",
        )
        
        # Should not raise any validation errors
        assert config.upload_intent == UploadIntent.create_version
        assert config.id == "test-model-id"

    def test_create_version_intent_missing_model_id_fails(self):
        """Test that create_version intent fails when model ID is missing"""
        with pytest.raises(ValidationError) as exc_info:
            ModelConfig(
                upload_intent=UploadIntent.create_version,
                # model id is missing
            )
        
        error_message = str(exc_info.value)
        assert "Model ID must be set when intent is create_version" in error_message

    def test_update_artifact_intent_with_artifact_id_succeeds(self):
        """Test that update_artifact intent succeeds when artifact ID is provided"""
        config = ModelConfig(
            upload_intent=UploadIntent.update_artifact,
            artifact_id="test-artifact-id"
        )
        
        # Should not raise any validation errors
        assert config.upload_intent == UploadIntent.update_artifact
        assert config.artifact_id == "test-artifact-id"

    def test_update_artifact_intent_without_model_and_version_ids_succeeds(self):
        """Test that update_artifact intent succeeds without model ID and version ID"""
        config = ModelConfig(
            upload_intent=UploadIntent.update_artifact,
            # no model id or version id
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
            )
        
        error_message = str(exc_info.value)
        assert "Model Artifact ID must be set when intent is update_artifact" in error_message

    def test_default_intent_is_update_artifact(self):
        """Test that the default upload intent is update_artifact"""
        config = ModelConfig(
            # intent is not set
            artifact_id="test-artifact-id"
        )
        
        # Default should be update_artifact
        assert config.upload_intent == UploadIntent.update_artifact
