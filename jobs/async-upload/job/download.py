from typing import Any, Dict

from model_registry import ModelRegistry
from model_registry.utils import _connect_to_s3


def download_from_s3(client: ModelRegistry, config: Dict[str, Any]):
    source_config = config["source"]["s3"]
    s3_client, _ = _connect_to_s3(
        endpoint_url=source_config["endpoint"],
        access_key_id=source_config["access_key_id"],
        secret_access_key=source_config["secret_access_key"],
        region=source_config["region"],
    )
    s3_client.download_file(
        source_config["bucket"], source_config["key"], config["storage"]["path"]
    )


def perform_download(client: ModelRegistry, config: Dict[str, Any]):
    # Download the model from the defined source
    if config.source.type == "s3":
        download_from_s3(client, config)
    elif config.source.type == "oci":
        # TODO: Implement the OCI download logic here
        raise ValueError("OCI source is not supported yet")
    else:
        raise ValueError(f"Unsupported source type: {config.source.type}")
