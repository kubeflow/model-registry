import logging

from job.models import ModelConfig, RegistryConfig
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
    logger.debug("🔍 Setting artifact to pending: %s", config.artifact_id)
    artifact = await client._api.get_model_artifact_by_id(config.artifact_id)

    if artifact is None:
        raise ValueError(f"Artifact {config.artifact_id} not found")
    
    artifact.state = ArtifactState.PENDING
    await client._api.upsert_model_artifact(artifact)
    logger.debug("✅ Artifact set to pending: %s", config.artifact_id)



async def update_model_artifact_uri(
    uri: str, client: ModelRegistry, config: ModelConfig
) -> None:
    logger.debug("🔍 Updating model artifact URI: %s", uri)
    artifact = await client._api.get_model_artifact_by_id(config.artifact_id)

    if artifact is None:
        raise ValueError(f"Artifact {config.artifact_id} not found")
    

    # Set the state of the artifact to LIVE and set the URI
    artifact.state=ArtifactState.LIVE
    artifact.uri=uri
    await client._api.upsert_model_artifact(artifact)
    logger.debug("✅ Model artifact URI updated: %s", uri)