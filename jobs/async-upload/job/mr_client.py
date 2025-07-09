from typing import Any, Dict

from model_registry import ModelRegistry
from model_registry.types import RegisteredModel


def validate_and_get_model_registry_client(config: Dict[str, Any]) -> ModelRegistry:
    """
    Validates the model registry client configuration and returns a ModelRegistry client.
    """
    return ModelRegistry(**config["registry"])


def perform_model_registration(
    client: ModelRegistry, config: Dict[str, Any]
) -> RegisteredModel:
    """
    Performs the model registration.
    """
    return client.register_model(
        config["model"]["name"],
        # URI will be populated after the upload completes
        "",
        model_format_name=config["model"]["model_format_name"],
        model_format_version=config["model"]["model_format_version"],
        version=config["model"]["version"],
        storage_key=config["model"]["storage_key"],
        storage_path=config["model"]["storage_path"],
    )


def update_model_registration(
    uri: str, client: ModelRegistry, registered_model: RegisteredModel
) -> None:
    updated_model = registered_model.update(uri=uri)
    client.update(updated_model)
