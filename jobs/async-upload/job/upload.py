from dataclasses import asdict
import os
import shutil
from pathlib import Path
from typing import Any, Dict
from model_registry.utils import OCIParams, S3Params, save_to_oci_registry


def _prepare_modelcar_structure(config: Dict[str, Any], model_files_path: str) -> str:
    """
    Prepare model files in KServe Modelcars-compatible structure.

    KServe Modelcars expects models to be in /models directory within the OCI image.
    This function creates the proper directory structure.

    Args:
        config: The configuration dictionary containing model metadata
        model_files_path: Path to the downloaded model files

    Returns:
        Path to the prepared modelcar directory structure
    """
    model_path = Path(model_files_path)

    # Create modelcar structure directory
    modelcar_base = model_path.parent / "modelcar"
    models_dir = modelcar_base / "models"

    # Clean up any existing structure
    if modelcar_base.exists():
        shutil.rmtree(modelcar_base)

    models_dir.mkdir(parents=True, exist_ok=True)

    # Copy model files to /models directory structure
    if model_path.is_file():
        # Single file - copy to models directory
        shutil.copy2(model_path, models_dir / model_path.name)
    elif model_path.is_dir():
        # Directory - copy contents to models directory
        for item in model_path.iterdir():
            if item.is_file():
                shutil.copy2(item, models_dir / item.name)
            else:
                shutil.copytree(item, models_dir / item.name)
    # If path doesn't exist, just create empty modelcar structure (already done above)

    return str(modelcar_base)


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
            base_image=config["model"]["name"],
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

    # For OCI uploads, prepare KServe Modelcars-compatible structure
    if config["destination"]["type"] == "oci" and os.path.exists(model_files_path):
        # Prepare modelcar directory structure (/models)
        modelcar_path = _prepare_modelcar_structure(config, model_files_path)

        # Update the model files path to use the modelcar structure
        upload_path = modelcar_path
    else:
        # For S3 or if model files don't exist, use original path
        upload_path = model_files_path

    upload_params = _get_upload_params(config)

    if isinstance(upload_params, S3Params):
        raise ValueError("S3 upload destination is not supported")
    elif isinstance(upload_params, OCIParams):
        uri = save_to_oci_registry(
            **asdict(upload_params),
            model_files_path=upload_path,
        )
    else:
        raise ValueError("Unsupported destination type")

    return uri
