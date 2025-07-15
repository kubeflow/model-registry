import asyncio
import logging
import os
import sys

from job.upload import perform_upload
from .config import get_config
from .mr_client import (
    validate_and_get_model_registry_client,
    set_artifact_pending,
    update_model_artifact_uri,
)
from .download import perform_download

# Configure logging
log_level = os.getenv('LOGLEVEL', logging.INFO)
logging.basicConfig(
    level=log_level, format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
    force=True
)
logger = logging.getLogger(__name__)

# Test logging configuration immediately
logger.info("📝 Logging configuration initialized successfully")


async def main() -> None:
    """
    Main entrypoint for the async upload job.
    Validates source and destination credentials before proceeding.
    """
    logger.info("🚀 Starting async upload job...")
    try:
        # Get complete configuration
        config = get_config()

        client = validate_and_get_model_registry_client(config)

        # Queue up model registration
        await set_artifact_pending(client, config)

        # Download the model from the defined source
        perform_download(client, config)


        # Upload the model to the destination
        uri = perform_upload(config)

        await update_model_artifact_uri(uri, client, config)

    except ValueError as e:
        logger.error(f"Configuration error: {str(e)}")
        raise
    except Exception as e:
        logger.error(f"Unexpected error: {str(e)}")
        raise
    logger.info("🏁 Job completed successfully")


if __name__ == "__main__":  # pragma: no cover
    asyncio.run(main())
