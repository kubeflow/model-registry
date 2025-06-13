from typing import Any, Dict

from model_registry import ModelRegistry


def validate_and_get_model_registry_client(config: Dict[str, Any]) -> ModelRegistry:
    """
    Validates the model registry client configuration and returns a ModelRegistry client.
    """
    return ModelRegistry(**config["registry"])
