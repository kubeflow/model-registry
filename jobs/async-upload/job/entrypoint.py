import asyncio
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


async def main() -> None:
    """
    Main entrypoint for the async upload job.
    Validates source and destination credentials before proceeding.
    """
    logger.info("🚀 Starting async upload job...")
    try:
        # Get complete configuration
        try:
            config = get_config()
        except Exception as e:
            raise RuntimeError("Failed to get config") from e

        client = validate_and_get_model_registry_client(config.registry)

        # Handle different intents
        intent = config.model.intent
        
        if isinstance(intent, UpdateArtifactIntent):
            # Original "Option 2" flow - update existing artifact
            logger.info("📋 Processing update_artifact intent")
            
            # Queue up model registration
            await set_artifact_pending(client, intent.artifact_id)

            # Download the model from the defined source
            perform_download(config)

            # Upload the model to the destination
            uri = perform_upload(config)

            await update_model_artifact_uri(client, intent.artifact_id, uri)
            
        elif isinstance(intent, CreateModelIntent):
            # "Option 1" flow - create new model, version, and artifact
            logger.info("📋 Processing create_model intent")
            
            if not config.metadata:
                raise ValueError("create_model intent requires ConfigMap metadata")

            # Download the model from the defined source
            perform_download(config)

            # Upload the model to the destination
            uri = perform_upload(config)

            # Create the complete model registry entry
            await create_model_and_artifact(client, config.metadata, uri)
            
        elif isinstance(intent, CreateVersionIntent):
            # "Option 1" flow - create new version and artifact under existing model
            logger.info("📋 Processing create_version intent")
            
            if not config.metadata:
                raise ValueError("create_version intent requires ConfigMap metadata")

            # Download the model from the defined source
            perform_download(config)

            # Upload the model to the destination
            uri = perform_upload(config)

            # Create the version and artifact under existing model
            await create_version_and_artifact(client, intent.model_id, config.metadata, uri)
        
        else:
            raise ValueError(f"Unknown intent type: {type(intent)}")

    except BaseException as e:
        record_error(e)
        raise
    else:
        logger.info("🏁 Job completed successfully")


if __name__ == "__main__":  # pragma: no cover
    asyncio.run(main())
