from typing import Any, Dict

from model_registry import ModelRegistry
from mr_openapi import ArtifactState


def validate_and_get_model_registry_client(config: Dict[str, Any]) -> ModelRegistry:
    """
    Validates the model registry client configuration and returns a ModelRegistry client.
    """
    return ModelRegistry(**config["registry"])


def set_artifact_pending(
    client: ModelRegistry, config: Dict[str, Any]
) -> None:
    """
    Sets the model artifact to pending.
    """
    artifact = client.get_model_artifact(
        config["model"]["name"], config["model"]["version_name"]
    )

    if artifact is None:
        raise ValueError(f"Artifact {config['model']['name']}/{config['model']['version_name']} not found")
    
    artifact.update(state=ArtifactState.PENDING)
    client.update(artifact)



def update_model_artifact_uri(
    uri: str, client: ModelRegistry, config: Dict[str, Any]
) -> None:
    artifact = client.get_model_artifact(
        config["model"]["name"], config["model"]["version_name"]
    )
    if artifact is None:
        raise ValueError(f"Artifact {config['model']['name']}/{config['model']['version_name']} not found, was it deleted since starting this job?")
    
    # Set the state of the artifact to LIVE and set the URI
    artifact.update(state=ArtifactState.LIVE, uri=uri)
    client.update(artifact)
