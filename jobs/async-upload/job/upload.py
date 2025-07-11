from dataclasses import asdict
from typing import Any, Dict
from model_registry.utils import OCIParams, S3Params, save_to_oci_registry


def _get_upload_params(config: Dict[str, Any]) -> S3Params | OCIParams:
    """
    Returns the upload params for the destination type

    Args:
        config: Configuration dictionary
    """
    destination_config = config["destination"]
    if destination_config["type"] == "s3":
        return S3Params(
            bucket_name=destination_config["s3"]["bucket"],
            s3_prefix=destination_config["s3"]["key"],
            endpoint_url=destination_config["s3"]["endpoint"],
            access_key_id=destination_config["s3"]["access_key_id"],
            secret_access_key=destination_config["s3"]["secret_access_key"],
            region=destination_config["s3"]["region"],
        )
    elif destination_config["type"] == "oci":
        return OCIParams(
            base_image=config["model"]["artifact_id"],
            oci_ref=destination_config["oci"]["uri"],
            dest_dir=config["storage"]["path"],
            oci_username=destination_config["oci"]["username"],
            oci_password=destination_config["oci"]["password"],
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

    if isinstance(upload_params, S3Params):
        raise ValueError("S3 upload destination is not supported")
    elif isinstance(upload_params, OCIParams):
        uri = save_to_oci_registry(
            **asdict(upload_params),
            model_files_path=model_files_path,
        )
    else:
        raise ValueError("Unsupported destination type")

    return uri
