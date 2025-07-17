from dataclasses import asdict, fields
import logging
from typing import Any, Dict
from model_registry import utils
from model_registry.utils import OCIParams, S3Params, save_to_oci_registry

logger = logging.getLogger(__name__)

def _get_upload_params(config: Dict[str, Any]) -> S3Params | OCIParams:
    """
    Returns the upload params for the destination type

    Args:
        config: Configuration dictionary
    """
    logger.debug("🔍 Getting upload params for destination type: %s", config["destination"]["type"])
    destination_config = config["destination"]
    if destination_config["type"] == "s3":
        return S3Params(
            bucket_name=destination_config["s3"]["bucket"],
            s3_prefix=destination_config["s3"]["key"],
            endpoint_url=destination_config["s3"]["endpoint_url"],
            access_key_id=destination_config["s3"]["access_key_id"],
            secret_access_key=destination_config["s3"]["secret_access_key"],
            region=destination_config["s3"]["region"],
        )
    elif destination_config["type"] == "oci":
        push_args = []
        # Note: These are all skopeo args, see: https://github.com/containers/skopeo/blob/main/docs/skopeo-copy.1.md
        if not destination_config["oci"]["enable_tls_verify"]:
            push_args.append("--dest-tls-verify=false")
        if destination_config["credentials_path"]:
            push_args.append("--authfile")
            push_args.append(destination_config["credentials_path"])

        return OCIParams(
            base_image=destination_config["oci"]["base_image"],
            oci_ref=destination_config["oci"]["uri"],
            dest_dir=config["storage"]["path"],
            oci_username=destination_config["oci"]["username"],
            oci_password=destination_config["oci"]["password"],
            # Same as the default backend, but with additional args included
            custom_oci_backend=utils._get_skopeo_backend(
                push_args=push_args
            ),
        )
    else:
        raise ValueError(f"Unsupported destination type: {destination_config['type']}")


def perform_upload(config: Dict[str, Any]) -> str:
    """
    Performs the upload of the model to the destination with KServe Modelcars compatibility

    Returns:
        The URI of the uploaded model
    """
    model_files_path = config["storage"]["path"]

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
