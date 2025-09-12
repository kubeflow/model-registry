import logging

from job.models import ModelConfig, RegistryConfig, UpdateArtifactIntent
from model_registry import ModelRegistry
from mr_openapi import ArtifactState

logger = logging.getLogger(__name__)

def validate_and_get_model_registry_client(config: RegistryConfig) -> ModelRegistry:
    """
    Validates the model registry client configuration and returns a ModelRegistry client.
    """
    logger.debug(f"🔍 Creating ModelRegistry client with config: {config}")
    return ModelRegistry(**config.model_dump())


async def set_artifact_pending(
    client: ModelRegistry, config: ModelConfig
) -> None:
    """
    Sets the model artifact to pending.
    """
    if not isinstance(config.intent, UpdateArtifactIntent):
        raise ValueError("set_artifact_pending can only be used with UpdateArtifactIntent")
    
    artifact_id = config.intent.artifact_id
    logger.debug("🔍 Setting artifact to pending: %s", artifact_id)
    artifact = await client._api.get_model_artifact_by_id(artifact_id)

    if artifact is None:
        raise ValueError(f"Artifact {artifact_id} not found")
    
    artifact.state = ArtifactState.PENDING
    await client._api.upsert_model_artifact(artifact)
    logger.debug("✅ Artifact set to pending: %s", artifact_id)



async def update_model_artifact_uri(
    uri: str, client: ModelRegistry, config: ModelConfig
) -> None:
    if not isinstance(config.intent, UpdateArtifactIntent):
        raise ValueError("update_model_artifact_uri can only be used with UpdateArtifactIntent")
    
    artifact_id = config.intent.artifact_id
    logger.debug("🔍 Updating model artifact URI: %s", uri)
    artifact = await client._api.get_model_artifact_by_id(artifact_id)

    if artifact is None:
        raise ValueError(f"Artifact {artifact_id} not found")
    

    # Set the state of the artifact to LIVE and set the URI
    artifact.state=ArtifactState.LIVE
    artifact.uri=uri
    await client._api.upsert_model_artifact(artifact)
    logger.debug("✅ Model artifact URI updated: %s", uri)