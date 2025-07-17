from typing import Any, Dict
import logging

from model_registry import ModelRegistry
from mr_openapi import ArtifactState

logger = logging.getLogger(__name__)

def validate_and_get_model_registry_client(config: Dict[str, Any]) -> ModelRegistry:
    """
    Validates the model registry client configuration and returns a ModelRegistry client.
    """
    client_config = config["registry"]
    logger.debug(f"🔍 Creating ModelRegistry client with config: {client_config}")
    return ModelRegistry(
        server_address=client_config["server_address"],
        port=client_config["port"],
        author=client_config["author"],
        is_secure=client_config["is_secure"],
        user_token=client_config["user_token"],
        user_token_envvar=client_config["user_token_envvar"],
        custom_ca=client_config["custom_ca"],
        custom_ca_envvar=client_config["custom_ca_envvar"],
        log_level=client_config["log_level"],
    )


async def set_artifact_pending(
    client: ModelRegistry, config: Dict[str, Any]
) -> None:
    """
    Sets the model artifact to pending.
    """
    logger.debug("🔍 Setting artifact to pending: %s", config['model']['artifact_id'])
    artifact = await client._api.get_model_artifact_by_id(config['model']['artifact_id'])

    if artifact is None:
        raise ValueError(f"Artifact {config['model']['artifact_id']} not found")
    
    artifact.state = ArtifactState.PENDING
    await client._api.upsert_model_artifact(artifact)
    logger.debug("✅ Artifact set to pending: %s", config['model']['artifact_id'])



async def update_model_artifact_uri(
    uri: str, client: ModelRegistry, config: Dict[str, Any]
) -> None:
    logger.debug("🔍 Updating model artifact URI: %s", uri)
    artifact = await client._api.get_model_artifact_by_id(config['model']['artifact_id'])

    if artifact is None:
        raise ValueError(f"Artifact {config['model']['artifact_id']} not found")
    

    # Set the state of the artifact to LIVE and set the URI
    artifact.state=ArtifactState.LIVE
    artifact.uri=uri
    await client._api.upsert_model_artifact(artifact)
    logger.debug("✅ Model artifact URI updated: %s", uri)