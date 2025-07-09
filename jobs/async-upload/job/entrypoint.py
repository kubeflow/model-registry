import logging

from job.upload import perform_upload
from .config import get_config
from .mr_client import (
    validate_and_get_model_registry_client,
    perform_model_registration,
    update_model_registration,
)
from .download import perform_download

# Configure logging
logging.basicConfig(
    level=logging.INFO, format="%(asctime)s - %(name)s - %(levelname)s - %(message)s"
)
logger = logging.getLogger(__name__)


def main() -> None:
    """
    Main entrypoint for the async upload job.
    Validates source and destination credentials before proceeding.
    """
    try:
        # Get complete configuration
        config = get_config()

        client = validate_and_get_model_registry_client(config)

        logger.info(
            f"Source: {config['source']['type'].upper()} storage at {config['source'].get('endpoint') or 'default endpoint'}"
        )
        logger.info(
            f"Destination: {config['destination']['type'].upper()} storage at {config['destination'].get('endpoint') or 'default endpoint'}"
        )

        # Download the model from the defined source
        perform_download(client, config)

        # Queue up model registration
        registered_model = perform_model_registration(client, config)

        # KServe Modelcars compatibility is handled within perform_upload()
        # For OCI destinations, it automatically creates the required /models directory structure
        # and prepares the model files in a KServe Modelcars-compatible format

        uri = perform_upload(config)

        update_model_registration(uri, client, registered_model)

    except ValueError as e:
        logger.error(f"Configuration error: {str(e)}")
        raise
    except Exception as e:
        logger.error(f"Unexpected error: {str(e)}")
        raise


if __name__ == "__main__":  # pragma: no cover
    main()
