import asyncio
import logging
import os

from job.upload import perform_upload
from .config import get_config
from .mr_client import (
    validate_and_get_model_registry_client,
    set_artifact_pending,
    update_model_artifact_uri,
)
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

        # Queue up model registration
        await set_artifact_pending(client, config.model)

        # Download the model from the defined source
        perform_download(config)

        # Upload the model to the destination
        uri = perform_upload(config)

        await update_model_artifact_uri(uri, client, config.model)

    except BaseException as e:
        record_error(e)
        raise
    else:
        logger.info("🏁 Job completed successfully")


if __name__ == "__main__":  # pragma: no cover
    asyncio.run(main())
