from dataclasses import asdict, fields
import logging
from typing import Any, Dict
from model_registry import utils
from model_registry.utils import OCIParams, S3Params, save_to_oci_registry

from .models import AsyncUploadConfig, DestinationConfig, OCIStorageConfig, S3StorageConfig

logger = logging.getLogger(__name__)

def _get_upload_params(config: AsyncUploadConfig) -> S3Params | OCIParams:
    """
    Returns the upload params for the destination type

    Args:
        config: Configuration dictionary
    """
    destination_config = config.destination
    logger.debug("🔍 Getting upload params for destination type: %s", destination_config)
    if isinstance(config.destination, S3StorageConfig):
        return S3Params(
            bucket_name=config.destination.bucket,
            s3_prefix=config.destination.key,
            endpoint_url=config.destination.endpoint_url,
            access_key_id=config.destination.access_key_id,
            secret_access_key=config.destination.secret_access_key,
            region=config.destination.region,
        )
    elif isinstance(destination_config, OCIStorageConfig):
        push_args = []
        # Note: These are all skopeo args, see: https://github.com/containers/skopeo/blob/main/docs/skopeo-copy.1.md
        if not destination_config.enable_tls_verify:
            push_args.append("--dest-tls-verify=false")
        if destination_config.credentials_path:
            push_args.append("--authfile")
            push_args.append(destination_config.credentials_path)

        return OCIParams(
            base_image=destination_config.base_image,
            oci_ref=destination_config.uri,
            dest_dir=config.storage.path,
            oci_username=destination_config.username,
            oci_password=destination_config.password,
            # Same as the default backend, but with additional args included
            custom_oci_backend=utils._get_skopeo_backend(
                push_args=push_args
            ),
        )
    else:
        raise ValueError(f"Unsupported destination type")


def perform_upload(config: AsyncUploadConfig) -> str:
    """
    Performs the upload of the model to the destination with KServe Modelcars compatibility

    Returns:
        The URI of the uploaded model
    """
    model_files_path = config.storage.path

    upload_params = _get_upload_params(config)
    logger.debug("🔍 Upload params: %s", upload_params)

    logger.info("📤 Uploading model to destination...")
    if isinstance(upload_params, S3Params):
        raise ValueError("S3 upload destination is not supported")
    elif isinstance(upload_params, OCIParams):
        uri = save_to_oci_registry(
            **{field.name: getattr(upload_params, field.name) for field in fields(upload_params)},
            model_files_path=model_files_path
        )
    else:
        raise ValueError("Unsupported destination type")

    logger.info("✅ Model uploaded to destination: %s", uri)
    return uri
