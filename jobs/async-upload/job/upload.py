from typing import Any, Dict
from model_registry import ModelRegistry
from model_registry.utils import OCIParams, S3Params


def _get_upload_params(config: Dict[str, Any]) -> S3Params | OCIParams:
    """
    Returns the upload params for the destination type
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
            base_image=config["model"]["name"],
            oci_ref=destination_config["oci"]["uri"],
            dest_dir=config["storage"]["path"],
            oci_username=destination_config["oci"]["username"],
            oci_password=destination_config["oci"]["password"],
        )
    else:
        raise ValueError(f"Unsupported destination type: {destination_config['type']}")


def perform_upload(client: ModelRegistry, config: Dict[str, Any]):
    """
    Performs the upload of the model to the destination
    """
    upload_params = _get_upload_params(config)

    client.upload_artifact_and_register_model(
        model_files_path=config["storage"]["path"],
        name=config["model"]["name"],
        version=config["model"]["version"],
        model_format_name=config["model"]["format"],
        model_format_version=config["model"]["format_version"],
        upload_params=upload_params,
    )
