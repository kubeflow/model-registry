import asyncio
import json
import logging
import os

from job.upload import perform_upload
from .config import get_config
from .mr_client import (
    validate_and_get_model_registry_client,
    set_artifact_pending,
    update_model_artifact_uri,
    create_model_and_artifact,
    create_version_and_artifact,
    validate_create_model_intent,
    validate_create_version_intent,
    CreatedEntityIds,
)
from .models import CreateModelIntent, CreateVersionIntent, UpdateArtifactIntent
from .download import perform_download

# Configure logging
log_level = os.getenv("LOGLEVEL", logging.INFO)
logging.basicConfig(
    level=log_level,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
    force=True,
)
logger = logging.getLogger(__name__)

# Test logging configuration immediately
logger.info("📝 Logging configuration initialized successfully")

termination_message_path = os.environ.get(
    "TERMINATION_MESSAGE_PATH", "/dev/termination-log"
)


def write_to_termination_message_path(message):
    with open(termination_message_path, "w") as f:
        f.write(message)


def record_error(exc):
    message = f"Unexpected error: {exc}"
    write_to_termination_message_path(message)
    logger.error(message)


def write_success_result(entity_ids: CreatedEntityIds, intent_type: str) -> None:
    """Write success result with entity IDs to termination message path.
    
    Args:
        entity_ids: CreatedEntityIds object containing the IDs of created/updated entities.
        intent_type: The intent type string (e.g., "update_artifact", "create_model", "create_version").
    """
    result_dict = {"intent": intent_type}
    
    if entity_ids.registered_model_id:
        result_dict["RegisteredModel"] = {"id": entity_ids.registered_model_id}
    
    if entity_ids.model_version_id:
        result_dict["ModelVersion"] = {"id": entity_ids.model_version_id}
    
    if entity_ids.model_artifact_id:
        result_dict["ModelArtifact"] = {"id": entity_ids.model_artifact_id}
    
    result_json = json.dumps(result_dict, indent=2)
    write_to_termination_message_path(result_json)
    logger.info(f"✅ Success result written to termination message path: {result_json}")


async def main() -> None:
    """
    Main entrypoint for the async upload job.
    Validates source and destination credentials before proceeding.
    """
    logger.info("🚀 Starting async upload job...")
    try:
        try:
            config = get_config()
        except Exception as e:
            raise RuntimeError("Failed to get config") from e

        client = validate_and_get_model_registry_client(config.registry)

        intent = config.model.intent
        entity_ids: CreatedEntityIds | None = None
        
        if isinstance(intent, UpdateArtifactIntent):
            logger.info("📋 Processing update_artifact intent")
            await set_artifact_pending(client, intent.artifact_id)
            perform_download(config)
            uri = perform_upload(config)
            # Pass through optional model_id and version_id if provided
            entity_ids = await update_model_artifact_uri(
                client,
                intent.artifact_id,
                uri,
                registered_model_id=intent.model_id,
                model_version_id=intent.version_id,
            )
        elif isinstance(intent, CreateModelIntent):
            logger.info("📋 Processing create_model intent")
            if not config.metadata:
                raise ValueError("create_model intent requires ConfigMap metadata")
            # Fast-fail validation before any expensive operations
            logger.info("🔍 Validating create_model intent...")
            await validate_create_model_intent(client, config.metadata)
            perform_download(config)
            uri = perform_upload(config)
            entity_ids = await create_model_and_artifact(client, config.metadata, uri)
        elif isinstance(intent, CreateVersionIntent):
            logger.info("📋 Processing create_version intent")
            if not config.metadata:
                raise ValueError("create_version intent requires ConfigMap metadata")
            # Fast-fail validation before any expensive operations
            logger.info("🔍 Validating create_version intent...")
            await validate_create_version_intent(client, intent.model_id, config.metadata)
            perform_download(config)
            uri = perform_upload(config)
            entity_ids = await create_version_and_artifact(client, intent.model_id, config.metadata, uri)
        else:
            raise ValueError(f"Unknown intent type: {type(intent)}")
        
        # Write success result to termination message path
        if entity_ids:
            write_success_result(entity_ids, intent.intent_type)
    except BaseException as e:
        record_error(e)
        raise
    else:
        logger.info("🏁 Job completed successfully")


if __name__ == "__main__":  # pragma: no cover
    asyncio.run(main())
