import os
import aiohttp
from model_registry import ModelRegistry
from model_registry.exceptions import StoreError
import pytest
from job.mr_client import (
    validate_and_get_model_registry_client,
    validate_create_model_intent,
    validate_create_version_intent,
)
from job.config import get_config
from job.models import (
    ConfigMapMetadata,
    RegisteredModelMetadata,
    ModelVersionMetadata,
    ModelArtifactMetadata,
)


@pytest.fixture
def sample_metadata():
    """Fixture providing sample ConfigMapMetadata for testing."""
    import time
    timestamp = int(time.time())
    return ConfigMapMetadata(
        registered_model=RegisteredModelMetadata(name=f"test-model-e2e-{timestamp}"),
        model_version=ModelVersionMetadata(name="v1.0.0"),
        model_artifact=ModelArtifactMetadata(name="test-artifact")
    )


@pytest.fixture
def version_only_metadata():
    """Fixture providing metadata for create_version intent (no registered_model)."""
    return ConfigMapMetadata(
        model_version=ModelVersionMetadata(name="v2.0.0"),
        model_artifact=ModelArtifactMetadata(name="test-artifact")
    )


@pytest.fixture
def mr_client(minimal_env_source_dest_vars):
    """Fixture providing a real ModelRegistry client for e2e tests."""
    sample_config = get_config(["--registry-is-secure", "false"])
    return validate_and_get_model_registry_client(sample_config.registry)


@pytest.fixture
def minimal_env_source_dest_vars():
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
    for key, value in vars.items():
        os.environ[f"MODEL_SYNC_SOURCE_{key.upper()}"] = value

    vars = {
        "model_upload_intent": "update_artifact",
        "model_artifact_id": "123",
        "registry_server_address": "http://localhost",
        "registry_port": "8080",
        "registry_author": "author",
    }

    for key, value in vars.items():
        os.environ[f"MODEL_SYNC_{key.upper()}"] = value

    yield vars

    # Restore original environment
    os.environ.clear()
    os.environ.update(original_env)


@pytest.mark.e2e
def test_model_registry_config_throws_error_on_missing_user_token(
    minimal_env_source_dest_vars,
):
    """Test that the model registry config throws an error on missing user token because it's a secure connection"""
    sample_config = get_config([])
    with pytest.raises(StoreError) as e:
        validate_and_get_model_registry_client(sample_config.registry)
    assert "user token must be provided for secure connection" in str(e.value)


@pytest.mark.e2e
def test_model_registry_config_correct(minimal_env_source_dest_vars):
    """Test that the model registry config is correct"""
    sample_config = get_config(["--registry-is-secure", False])
    # Note: Instantiating the client will ping the GET /.../registered_models endpoint, validating the connection
    client = validate_and_get_model_registry_client(sample_config.registry)
    assert isinstance(client, ModelRegistry)


@pytest.mark.e2e
def test_model_registry_config_throws_when_mr_is_unreachable(
    minimal_env_source_dest_vars,
):
    """Test that the model registry config is correct"""
    sample_config = get_config(
        [
            "--registry-is-secure",
            "false",
            "--registry-port",
            "1337",  # Note: the E2E test should expose 8080, this is a purposely invalid port
        ]
    )

    with pytest.raises(aiohttp.client_exceptions.ClientConnectorError) as e:
        validate_and_get_model_registry_client(sample_config.registry)
    assert "Cannot connect to host" in str(e.value)


@pytest.mark.e2e
@pytest.mark.asyncio
async def test_validate_create_model_intent_success(mr_client, sample_metadata):
    """Test successful validation for create_model intent with non-existent model."""
    # Ensure the test model doesn't exist by trying to clean it up first
    try:
        existing_model = await mr_client._api.get_registered_model_by_params(sample_metadata.registered_model.name)
        if existing_model:
            # Clean up any existing test model to ensure clean state
            await mr_client._api.delete_registered_model(existing_model.id)
    except Exception:
        pass  # Model doesn't exist, which is what we want
    
    # Should not raise any exception when model doesn't exist
    await validate_create_model_intent(mr_client, sample_metadata)


@pytest.mark.e2e
@pytest.mark.asyncio
async def test_validate_create_model_intent_model_already_exists(mr_client, sample_metadata):
    """Test validation failure when model already exists."""
    # First, create a model to ensure it exists
    existing_model = await mr_client._register_model(
        name=sample_metadata.registered_model.name,
        description="Test model for fast-fail validation",
        owner="test-user"
    )
    
    try:
        # Should raise ValueError with friendly message
        with pytest.raises(ValueError) as exc_info:
            await validate_create_model_intent(mr_client, sample_metadata)
        
        assert f"Cannot create model: RegisteredModel with name '{sample_metadata.registered_model.name}' already exists" in str(exc_info.value)
        assert "Use 'create_version' intent to add a new version" in str(exc_info.value)
    finally:
        # Clean up the test model
        try:
            await mr_client._api.delete_registered_model(existing_model.id)
        except Exception:
            pass


@pytest.mark.e2e
@pytest.mark.asyncio
async def test_validate_create_model_intent_missing_metadata(mr_client):
    """Test validation failure when required metadata is missing."""
    # Test with None metadata
    with pytest.raises(ValueError) as exc_info:
        await validate_create_model_intent(mr_client, None)
    assert "create_model intent requires complete metadata" in str(exc_info.value)
    
    # Test with missing registered_model
    metadata = ConfigMapMetadata(
        model_version=ModelVersionMetadata(name="v1.0.0"),
        model_artifact=ModelArtifactMetadata(name="test-artifact")
    )
    with pytest.raises(ValueError) as exc_info:
        await validate_create_model_intent(mr_client, metadata)
    assert "create_model intent requires complete metadata" in str(exc_info.value)
    
    # Test with missing model name - this will fail at Pydantic validation level
    with pytest.raises((ValueError, Exception)) as exc_info:
        metadata = ConfigMapMetadata(
            registered_model=RegisteredModelMetadata(name=None),
            model_version=ModelVersionMetadata(name="v1.0.0"),
            model_artifact=ModelArtifactMetadata(name="test-artifact")
        )
        await validate_create_model_intent(mr_client, metadata)
    # Either Pydantic validation error or our custom validation error is acceptable
    assert ("Must provide either name or id" in str(exc_info.value) or 
            "RegisteredModel name is required" in str(exc_info.value))


@pytest.mark.e2e
@pytest.mark.asyncio
async def test_validate_create_version_intent_success(mr_client, version_only_metadata):
    """Test successful validation for create_version intent."""
    # First, create a model to use as parent
    import time
    timestamp = int(time.time())
    parent_model = await mr_client._register_model(
        name=f"test-parent-model-e2e-{timestamp}",
        description="Parent model for version validation test",
        owner="test-user"
    )
    
    try:
        # Should not raise any exception when model exists but version doesn't
        await validate_create_version_intent(mr_client, str(parent_model.id), version_only_metadata)
    finally:
        # Clean up the test model
        try:
            await mr_client._api.delete_registered_model(parent_model.id)
        except Exception:
            pass


@pytest.mark.e2e
@pytest.mark.asyncio
async def test_validate_create_version_intent_model_not_found(mr_client, version_only_metadata):
    """Test validation failure when parent model doesn't exist."""
    # Use a non-existent model ID
    non_existent_model_id = "99999"
    
    # Should raise ValueError with friendly message
    with pytest.raises(ValueError) as exc_info:
        await validate_create_version_intent(mr_client, non_existent_model_id, version_only_metadata)
    
    assert f"Cannot create version: RegisteredModel with ID '{non_existent_model_id}' not found" in str(exc_info.value)
    assert "Use 'create_model' intent to create a new model first" in str(exc_info.value)


@pytest.mark.e2e
@pytest.mark.asyncio
async def test_validate_create_version_intent_version_already_exists(mr_client, version_only_metadata):
    """Test validation failure when version already exists."""
    # First, create a model and version to ensure they exist
    import time
    timestamp = int(time.time())
    parent_model = await mr_client._register_model(
        name=f"test-parent-model-with-version-e2e-{timestamp}",
        description="Parent model for version conflict test",
        owner="test-user"
    )
    
    existing_version = await mr_client._register_new_version(
        parent_model,
        version_only_metadata.model_version.name,
        "test-user",
        description="Existing version for conflict test"
    )
    
    try:
        # Should raise ValueError with friendly message
        with pytest.raises(ValueError) as exc_info:
            await validate_create_version_intent(mr_client, str(parent_model.id), version_only_metadata)
        
        assert f"Cannot create version: ModelVersion with name '{version_only_metadata.model_version.name}' already exists" in str(exc_info.value)
        assert f"under RegisteredModel '{parent_model.name}'" in str(exc_info.value)
        assert "Use 'update_artifact' intent to update an existing version's artifact" in str(exc_info.value)
    finally:
        # Clean up the test model (this will also clean up the version)
        try:
            await mr_client._api.delete_registered_model(parent_model.id)
        except Exception:
            pass


@pytest.mark.e2e
@pytest.mark.asyncio
async def test_validate_create_version_intent_missing_metadata(mr_client):
    """Test validation failure when required metadata is missing."""
    # Test with None metadata
    with pytest.raises(ValueError) as exc_info:
        await validate_create_version_intent(mr_client, "model-123", None)
    assert "create_version intent requires metadata for model_version and model_artifact" in str(exc_info.value)
    
    # Test with missing model_version
    metadata = ConfigMapMetadata(
        model_artifact=ModelArtifactMetadata(name="test-artifact")
    )
    with pytest.raises(ValueError) as exc_info:
        await validate_create_version_intent(mr_client, "model-123", metadata)
    assert "create_version intent requires metadata for model_version and model_artifact" in str(exc_info.value)
    
    # Test with missing version name - this will fail at Pydantic validation level
    with pytest.raises((ValueError, Exception)) as exc_info:
        metadata = ConfigMapMetadata(
            model_version=ModelVersionMetadata(name=None),
            model_artifact=ModelArtifactMetadata(name="test-artifact")
        )
        await validate_create_version_intent(mr_client, "model-123", metadata)
    # Either Pydantic validation error or our custom validation error is acceptable
    assert ("name" in str(exc_info.value).lower() or 
            "ModelVersion name is required" in str(exc_info.value))
