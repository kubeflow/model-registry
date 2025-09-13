import logging

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


async def create_model_and_artifact(client: ModelRegistry, metadata: ConfigMapMetadata, uri: str) -> None:
    """Creates a new registered model, model version, and model artifact."""
    if not metadata or not metadata.registered_model or not metadata.model_version or not metadata.model_artifact:
        raise ValueError(
            "create_model intent requires complete metadata for registered_model, model_version, and model_artifact"
        )

    if not metadata.registered_model.name:
        raise ValueError("RegisteredModel name is required for create_model intent")

    logger.debug("🔍 Creating new registered model, version, and artifact")
    rm = await _create_registered_model(client, metadata.registered_model)
    await _create_version_and_artifact_for_model(client, rm, uri, metadata)


async def create_version_and_artifact(
    client: ModelRegistry, model_id: str, metadata: ConfigMapMetadata, uri: str
) -> None:
    """Creates a new model version and model artifact under an existing registered model."""
    if not metadata or not metadata.model_version or not metadata.model_artifact:
        raise ValueError("create_version intent requires metadata for model_version and model_artifact")

    logger.debug("🔍 Creating new version and artifact for model ID: %s", model_id)

    rm = await client._api.get_registered_model_by_id(model_id)
    if not rm:
        raise ValueError(f"RegisteredModel with ID '{model_id}' not found")

    await _create_version_and_artifact_for_model(client, rm, uri, metadata)


async def update_model_artifact_uri(client: ModelRegistry, artifact_id: str, uri: str) -> None:
    logger.debug("🔍 Updating model artifact URI: %s", uri)
    artifact = await client._api.get_model_artifact_by_id(artifact_id)

    if artifact is None:
        raise ValueError(f"Artifact {artifact_id} not found")

    # Set the state of the artifact to LIVE and set the URI
    artifact.state = ArtifactState.LIVE
    artifact.uri = uri
    await client._api.upsert_model_artifact(artifact)
    logger.debug("✅ Model artifact URI updated: %s", uri)


async def _create_registered_model(client: ModelRegistry, rm_metadata):
    """Creates a new registered model and returns it."""
    # Check if model already exists
    existing_rm = await client._api.get_registered_model_by_params(rm_metadata.name)
    if existing_rm:
        raise ValueError(f"RegisteredModel with name '{rm_metadata.name}' already exists")

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
) -> None:
    """Creates a model version and artifact under the given registered model."""
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