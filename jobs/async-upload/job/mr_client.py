import logging
from dataclasses import dataclass

from model_registry import ModelRegistry
from model_registry.types import ArtifactState
from job.models import (
    ModelConfig,
    RegistryConfig,
    UpdateArtifactIntent,
    CreateModelIntent,
    CreateVersionIntent,
    ConfigMapMetadata,
)

logger = logging.getLogger(__name__)


@dataclass
class CreatedEntityIds:
    """IDs of created/updated entities in the model registry."""
    registered_model_id: str | None = None
    model_version_id: str | None = None
    model_artifact_id: str | None = None


def validate_and_get_model_registry_client(config: RegistryConfig) -> ModelRegistry:
    """
    Validates the model registry client configuration and returns a ModelRegistry client.
    """
    logger.debug(f"🔍 Creating ModelRegistry client with config: {config}")
    return ModelRegistry(**config.model_dump())


async def set_artifact_pending(client: ModelRegistry, artifact_id: str) -> None:
    """
    Sets the model artifact to pending.
    """
    logger.debug("🔍 Setting artifact to pending: %s", artifact_id)
    artifact = await client._api.get_model_artifact_by_id(artifact_id)

    if artifact is None:
        raise ValueError(f"Artifact {artifact_id} not found")

    artifact.state = ArtifactState.PENDING
    await client._api.upsert_model_artifact(artifact)
    logger.debug("✅ Artifact set to pending: %s", artifact_id)


async def validate_create_model_intent(client: ModelRegistry, metadata: ConfigMapMetadata) -> None:
    """Fast-fail validation for create_model intent."""
    if not metadata or not metadata.registered_model or not metadata.model_version or not metadata.model_artifact:
        raise ValueError(
            "create_model intent requires complete metadata for registered_model, model_version, and model_artifact"
        )

    if not metadata.registered_model.name:
        raise ValueError("RegisteredModel name is required for create_model intent")

    # Fast-fail check: ensure RegisteredModel doesn't already exist
    existing_rm = await client._api.get_registered_model_by_params(metadata.registered_model.name)
    if existing_rm:
        raise ValueError(
            f"Cannot create model: RegisteredModel with name '{metadata.registered_model.name}' already exists. "
            f"Use 'create_version' intent to add a new version to this existing model."
        )

    logger.debug("✅ create_model intent validation passed")


async def validate_create_version_intent(client: ModelRegistry, model_id: str, metadata: ConfigMapMetadata) -> None:
    """Fast-fail validation for create_version intent."""
    if not metadata or not metadata.model_version or not metadata.model_artifact:
        raise ValueError("create_version intent requires metadata for model_version and model_artifact")

    if not metadata.model_version.name:
        raise ValueError("ModelVersion name is required for create_version intent")

    # Fast-fail check: ensure RegisteredModel exists
    rm = await client._api.get_registered_model_by_id(model_id)
    if not rm:
        raise ValueError(
            f"Cannot create version: RegisteredModel with ID '{model_id}' not found. "
            f"Use 'create_model' intent to create a new model first."
        )

    # Fast-fail check: ensure ModelVersion doesn't already exist
    try:
        existing_mv = await client._api.get_model_version_by_params(
            registered_model_id=model_id, 
            name=metadata.model_version.name
        )
        if existing_mv:
            raise ValueError(
                f"Cannot create version: ModelVersion with name '{metadata.model_version.name}' already exists "
                f"under RegisteredModel '{rm.name}' (ID: {model_id}). "
                f"Use 'update_artifact' intent to update an existing version's artifact."
            )
    except Exception as e:
        # If the API call fails for reasons other than "not found", we should still fail
        if "not found" not in str(e).lower():
            raise ValueError(f"Failed to check if ModelVersion exists: {e}") from e

    logger.debug("✅ create_version intent validation passed")


async def create_model_and_artifact(client: ModelRegistry, metadata: ConfigMapMetadata, uri: str) -> CreatedEntityIds:
    """Creates a new registered model, model version, and model artifact.
    
    Returns:
        CreatedEntityIds: IDs of the created registered model, model version, and model artifact.
    """
    logger.debug("🔍 Creating new registered model, version, and artifact")
    rm = await _create_registered_model(client, metadata.registered_model)
    mv, artifact = await _create_version_and_artifact_for_model(client, rm, uri, metadata)
    return CreatedEntityIds(
        registered_model_id=str(rm.id) if rm.id else None,
        model_version_id=str(mv.id) if mv.id else None,
        model_artifact_id=str(artifact.id) if artifact.id else None,
    )


async def create_version_and_artifact(
    client: ModelRegistry, model_id: str, metadata: ConfigMapMetadata, uri: str
) -> CreatedEntityIds:
    """Creates a new model version and model artifact under an existing registered model.
    
    Returns:
        CreatedEntityIds: IDs of the existing registered model, and created model version and model artifact.
    """
    logger.debug("🔍 Creating new version and artifact for model ID: %s", model_id)

    rm = await client._api.get_registered_model_by_id(model_id)
    if not rm:
        raise ValueError(f"RegisteredModel with ID '{model_id}' not found")

    mv, artifact = await _create_version_and_artifact_for_model(client, rm, uri, metadata)
    return CreatedEntityIds(
        registered_model_id=str(rm.id) if rm.id else None,
        model_version_id=str(mv.id) if mv.id else None,
        model_artifact_id=str(artifact.id) if artifact.id else None,
    )


async def update_model_artifact_uri(
    client: ModelRegistry,
    artifact_id: str,
    uri: str,
    registered_model_id: str | None = None,
    model_version_id: str | None = None,
) -> CreatedEntityIds:
    """Updates the model artifact URI and sets state to LIVE.
    
    Args:
        client: Model registry client.
        artifact_id: ID of the artifact to update.
        uri: New URI for the artifact.
        registered_model_id: Optional registered model ID to pass through to output.
        model_version_id: Optional model version ID to pass through to output.
    
    Returns:
        CreatedEntityIds: IDs passed through to output. For update_artifact intent,
        model and version IDs will be in the output if and only if they are passed in.
    """
    logger.debug("🔍 Updating model artifact URI: %s", uri)
    artifact = await client._api.get_model_artifact_by_id(artifact_id)

    if artifact is None:
        raise ValueError(f"Artifact {artifact_id} not found")

    # Set the state of the artifact to LIVE and set the URI
    artifact.state = ArtifactState.LIVE
    artifact.uri = uri
    await client._api.upsert_model_artifact(artifact)
    logger.debug("✅ Model artifact URI updated: %s", uri)
    
    # Pass through the IDs that were provided
    return CreatedEntityIds(
        registered_model_id=registered_model_id,
        model_version_id=model_version_id,
        model_artifact_id=artifact_id,
    )


async def _create_registered_model(client: ModelRegistry, rm_metadata):
    """Creates a new registered model and returns it."""
    rm = await client._register_model(
        name=rm_metadata.name,
        owner=rm_metadata.owner,
        description=rm_metadata.description,
        custom_properties=rm_metadata.custom_properties or {},
    )
    logger.debug("✅ Created RegisteredModel: %s (ID: %s)", rm.name, rm.id)
    return rm


async def _create_version_and_artifact_for_model(
    client: ModelRegistry, rm, uri: str, metadata: ConfigMapMetadata
) -> tuple:
    """Creates a model version and artifact under the given registered model.
    
    Returns:
        tuple: (ModelVersion, ModelArtifact) - The created model version and artifact objects.
    """
    mv_metadata = metadata.model_version
    version_name = mv_metadata.name or "1.0.0"

    mv = await client._register_new_version(
        rm,
        version_name,
        mv_metadata.author or client._author,
        description=mv_metadata.description,
        custom_properties=mv_metadata.custom_properties or {},
    )
    logger.debug("✅ Created ModelVersion: %s (ID: %s)", mv.name, mv.id)

    # Create model artifact
    ma_metadata = metadata.model_artifact
    artifact_name = ma_metadata.name or rm.name

    artifact = await client._register_model_artifact(
        mv,
        artifact_name,
        uri,
        model_format_name=ma_metadata.model_format_name,
        model_format_version=ma_metadata.model_format_version,
        storage_key=ma_metadata.storage_key,
        storage_path=ma_metadata.storage_path,
        service_account_name=ma_metadata.service_account_name,
        model_source_kind=ma_metadata.model_source_kind,
        model_source_class=ma_metadata.model_source_class,
        model_source_group=ma_metadata.model_source_group,
        model_source_id=ma_metadata.model_source_id,
        model_source_name=ma_metadata.model_source_name,
        custom_properties=ma_metadata.custom_properties or {},
    )
    logger.debug("✅ Created ModelArtifact: %s (ID: %s) with URI: %s", artifact.name, artifact.id, uri)
    # Set the artifact state to LIVE since it has a valid URI
    artifact.state = ArtifactState.LIVE
    await client._api.upsert_model_artifact(artifact)
    logger.debug("✅ Updated ModelArtifact state to LIVE: %s", artifact.id)
    
    return (mv, artifact)