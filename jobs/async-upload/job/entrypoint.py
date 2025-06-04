import logging
from config import get_config

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

        # TODO: Validate credentials and connection to Model Registry

        logger.info(
            f"Source: {config.source.type.upper()} storage at {config.source.endpoint or 'default endpoint'}"
        )
        logger.info(
            f"Destination: {config.destination.type.upper()} storage at {config.destination.endpoint or 'default endpoint'}"
        )

        # TODO: Implement the main upload logic here

    except ValueError as e:
        logger.error(f"Configuration error: {str(e)}")
        raise
    except Exception as e:
        logger.error(f"Unexpected error: {str(e)}")
        raise


if __name__ == "__main__":  # pragma: no cover
    main()
