from typing import Any, Dict
import os
import logging

from model_registry import ModelRegistry
from model_registry.utils import _connect_to_s3

logger = logging.getLogger(__name__)

HF_URI_PREFIX = "hf://"


def download_from_s3(client: ModelRegistry, config: Dict[str, Any]):
    logger.debug("🔍 Downloading model from S3...")
    logger.debug("🔍 Source config: %s", config["source"]["s3"])
    source_config = config["source"]["s3"]
    s3_client, _ = _connect_to_s3(
        source_config["endpoint_url"],
        source_config["access_key_id"],
        source_config["secret_access_key"],
        source_config["region"],
        multipart_threshold=1024 * 1024,
        multipart_chunksize=1024 * 1024,
        max_pool_connections=10,
    )

    bucket_name = source_config["bucket"]
    prefix = source_config["key"]

    # TODO: It might make sense to check if the provided key points to a single file first before assuming the directory needs to be traversed

    paginator = s3_client.get_paginator("list_objects_v2")
    for page in paginator.paginate(Bucket=bucket_name, Prefix=prefix):
        for obj in page.get("Contents", []):
            key = obj["Key"]
            if key.endswith("/"):
                continue
            relative = os.path.relpath(key, prefix)
            local_path = os.path.join(config["storage"]["path"], relative)
            os.makedirs(os.path.dirname(local_path), exist_ok=True)
            logger.info(f"⏳ Downloading s3://{bucket_name}/{key} → {local_path}")
            s3_client.download_file(bucket_name, key, local_path)
            logger.info(f"☑️ Downloaded s3://{bucket_name}/{key} → {local_path}")
    logger.debug("✅ Model files downloaded from S3")


def download_from_hf(uri: str, dest_dir: str) -> str:
    """
    adapted from kserve:
    https://github.com/kserve/kserve/blob/4edbb36c520c2e880842229bfc56b7f11d766822/python/storage/kserve_storage/kserve_storage.py#L292-L322
    """
    from huggingface_hub import snapshot_download

    if not uri.startswith(HF_URI_PREFIX):
        raise ValueError(f"Expected URI beginning with {HF_URI_PREFIX}")

    components = uri[len(HF_URI_PREFIX) :].split("/")
    if len(components) != 2:
        raise ValueError(
            "URI must contain exactly one '/' separating the repo and model name"
        )

    repo, model_id = components
    if not repo:
        raise ValueError("Repository name cannot be empty")

    model_name, _, hash_value = model_id.partition(":")
    if not model_name:
        raise ValueError("Model name cannot be empty")

    revision = hash_value if hash_value else None
    return snapshot_download(
        repo_id=f"{repo}/{model_name}", revision=revision, local_dir=dest_dir
    )


def perform_download(client: ModelRegistry, config: Dict[str, Any]):
    logger.info("📥 Downloading model from source...")
    # Download the model from the defined source
    source_type = config["source"]["type"]
    if source_type == "s3":
        download_from_s3(client, config)
    elif source_type == "hf" or config["source"]["uri"].startswith(HF_URI_PREFIX):
        download_from_hf(config["source"]["uri"], config["storage"]["path"])
    elif source_type == "oci":
        # TODO: Implement the OCI download logic here
        raise ValueError("OCI source is not supported yet")
    else:
        raise ValueError(f"Unsupported source type: {config['source']['type']}")
    logger.info("✅ Model downloaded from source")
