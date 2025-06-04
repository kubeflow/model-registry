from __future__ import annotations
import base64
import json
import configargparse as cap
from typing import Any, Dict, Mapping
from pathlib import Path


def _parser() -> cap.ArgumentParser:
    """Parse command line arguments and config files"""
    p = cap.ArgumentParser(
        default_config_files=[],
        auto_env_var_prefix="MODEL_SYNC_",
        description="Synchronise AI models between OCI registries and/or S3 buckets",
    )

    # --- source ---
    # s3
    # TODO: We should be able to infer the type from the credentials provided, therefore no default needed
    p.add("--source-type", choices=["s3", "oci"], default="s3")
    p.add("--source-aws-bucket")
    p.add("--source-aws-region")
    p.add("--source-aws-access-key-id")
    p.add("--source-aws-secret-access-key")
    p.add("--source-aws-endpoint")
    # OCI registry
    p.add("--source-oci-uri")
    p.add("--source-oci-username")
    p.add("--source-oci-password")

    # --- destination ---
    # s3
    # TODO: We should be able to infer the type from the credentials provided, therefore no default needed
    p.add("--destination-type", choices=["s3", "oci"], default="oci")
    p.add("--destination-aws-bucket")
    p.add("--destination-aws-region")
    p.add("--destination-aws-access-key-id")
    p.add("--destination-aws-secret-access-key")
    p.add("--destination-aws-endpoint")
    # OCI registry
    p.add("--destination-oci-uri")
    p.add("--destination-oci-username")
    p.add("--destination-oci-password")

    # --- model-registry --- TODO: use IDs https://github.com/kubeflow/model-registry/issues/1108#issuecomment-2880448765
    p.add("--model-name")
    p.add("--model-version")
    p.add("--model-format")
    p.add("--registry-url")

    p.add(
        "--source-s3-credentials-path",
        metavar="PATH",
    )
    p.add(
        "--destination-s3-credentials-path",
        metavar="PATH",
    )
    p.add(
        "--source-oci-credentials-path",
        metavar="PATH",
    )
    p.add(
        "--destination-oci-credentials-path",
        metavar="PATH",
    )

    return p


def _load_s3_credentials(path: str | Path, store: Mapping[str, Any], side: str) -> None:
    """
    The path must be a folder, with a number of files that match a typical AWS config, ie - a Secret mounted to a pod with keys:
    - AWS_ACCESS_KEY_ID
    - AWS_SECRET_ACCESS_KEY
    - AWS_BUCKET
    - AWS_REGION
    - AWS_ENDPOINT_URL

    These would be mounted as files with the names above and whose contents are the secret values.

    These keys are loaded into the config store
    """

    # Validate the path is a directory
    p = Path(path).expanduser()
    if not p.is_dir():
        raise FileNotFoundError(f"{side}-credentials folder not found: {p}")

    # Load the credentials from the files
    for file in p.glob("*"):
        if file.is_file():
            if file.name.startswith("AWS_"):
                # Converts a file with name AWS_ACCESS_KEY_ID to aws_access_key_id
                key_name = file.name.lower().replace("aws_", "")
                store["s3"][key_name] = file.read_text()
            else:
                raise ValueError(f"Invalid credential file name: {file.name}")


def _load_oci_credentials(
    path: str | Path, store: Mapping[str, Any], side: str
) -> None:
    """
    The path must point to a file which is a `.dockerconfigjson` file

    A typical file looks like this:
    ```json
    {
        "auths": {
            "registry.example.com": {
                "auth": "base64(username:password)",
                "email": "user@example.com"
            }
        },
        // ...
    }
    ```
    """

    # Validate the path is a file
    p = Path(path).expanduser()
    if not p.is_file():
        raise FileNotFoundError(f"{side}-credentials file not found: {p}")

    # Load the credentials from the file
    docker_config = json.loads(p.read_text())

    # Validate the docker config is valid
    if "auths" not in docker_config:
        raise ValueError("Invalid docker config file")

    # Load the credentials from the docker config, the URI passed in via config is used as a key to find the correct credentials
    auth = docker_config["auths"][store["uri"]]["auth"]
    username, password = base64.b64decode(auth).decode("utf-8").split(":")
    store["username"] = username
    store["password"] = password
    store["email"] = docker_config["auths"][store["uri"]]["email"]


def _validate_oci_config(cfg: Dict[str, Any]) -> None:
    """Validates the OCI config is valid"""
    # if the username is set but the password is not (and vice-versa) throw an error
    if cfg["oci"]["username"] is not None and cfg["oci"]["password"] is None:
        raise ValueError("OCI password must be set")
    if cfg["oci"]["username"] is None and cfg["oci"]["password"] is not None:
        raise ValueError("OCI username must be set")

    if cfg["oci"]["uri"] is None:
        raise ValueError("OCI URI must be set")


def _validate_s3_config(cfg: Dict[str, Any]) -> None:
    """Validates the S3 config is valid"""
    if cfg["s3"]["access_key_id"] is None or cfg["s3"]["secret_access_key"] is None:
        raise ValueError("S3 credentials must be set")
    if cfg["s3"]["bucket"] is None:
        raise ValueError("S3 bucket must be set")


def _validate_model_config(cfg: Dict[str, Any]) -> None:
    """Validates the model config is valid"""
    if cfg["name"] is None or cfg["version"] is None or cfg["format"] is None:
        raise ValueError("Model must be set")


def _validate_registry_config(cfg: Dict[str, Any]) -> None:
    """Validates the registry config is valid"""
    if cfg["url"] is None:
        raise ValueError("Registry URL must be set")


def _validate_store(cfg: Dict[str, Any]) -> None:
    """Validates the store is valid"""
    if cfg["type"] == "s3":
        _validate_s3_config(cfg)
    elif cfg["type"] == "oci":
        _validate_oci_config(cfg)
    else:
        raise ValueError("Source type must be set")


def _validate_config(cfg: Dict[str, Any]) -> None:
    """Validates the config is has the credentials and locations needed to perform an async upload"""

    # Ensure the source is valid
    _validate_store(cfg["source"])
    _validate_store(cfg["destination"])

    # Ensure the model is valid
    _validate_model_config(cfg["model"])

    # Ensure the registry is valid
    _validate_registry_config(cfg["registry"])


def get_config(argv: list[str] | None = None) -> Dict[str, Any]:
    """
    Return a plain nested dict suitable for boto3 / skopeo / register_model.

    Priority of the config is:
    1. Command-line arguments
    2. Environment variables
    3. Credentials files
    4. Default values
    """
    args = _parser().parse_args(argv)

    # Initialize config with command-line arguments and defaults
    cfg = {
        "source": {
            "type": args.source_type,
            "s3": {
                "bucket": None,
                "region": None,
                "access_key_id": None,
                "secret_access_key": None,
                "endpoint": None,
            },
            "oci": {
                "uri": args.source_oci_uri,
                "username": None,
                "password": None,
                "email": None,
            },
        },
        "destination": {
            "type": args.destination_type,
            "s3": {
                "bucket": None,
                "region": None,
                "access_key_id": None,
                "secret_access_key": None,
                "endpoint": None,
            },
            "oci": {
                "uri": args.destination_oci_uri,
                "username": None,
                "password": None,
                "email": None,
            },
        },
        "model": {
            "name": args.model_name,
            "version": args.model_version,
            "format": args.model_format,
            # TODO: Add the rest of the needed values
        },
        "registry": {
            "url": args.registry_url,
            # TODO: Add the rest of the needed values
        },
    }

    # Load credentials from files, if provided
    if args.source_s3_credentials_path:
        _load_s3_credentials(args.source_s3_credentials_path, cfg["source"], "source")
    elif args.source_oci_credentials_path:
        _load_oci_credentials(args.source_oci_credentials_path, cfg["source"], "source")

    if args.destination_s3_credentials_path:
        _load_s3_credentials(
            args.destination_s3_credentials_path, cfg["destination"], "destination"
        )
    elif args.destination_oci_credentials_path:
        _load_oci_credentials(
            args.destination_oci_credentials_path, cfg["destination"], "destination"
        )

    # TODO: Maybe clean this up, its a little manual
    # Override with command-line arguments if provided. configargparse will prioritize CLI > ENV
    if args.source_aws_bucket:
        cfg["source"]["s3"]["bucket"] = args.source_aws_bucket
    if args.source_aws_region:
        cfg["source"]["s3"]["region"] = args.source_aws_region
    if args.source_aws_access_key_id:
        cfg["source"]["s3"]["access_key_id"] = args.source_aws_access_key_id
    if args.source_aws_secret_access_key:
        cfg["source"]["s3"]["secret_access_key"] = args.source_aws_secret_access_key
    if args.source_aws_endpoint:
        cfg["source"]["s3"]["endpoint"] = args.source_aws_endpoint

    if args.destination_aws_bucket:
        cfg["destination"]["s3"]["bucket"] = args.destination_aws_bucket
    if args.destination_aws_region:
        cfg["destination"]["s3"]["region"] = args.destination_aws_region
    if args.destination_aws_access_key_id:
        cfg["destination"]["s3"]["access_key_id"] = args.destination_aws_access_key_id
    if args.destination_aws_secret_access_key:
        cfg["destination"]["s3"]["secret_access_key"] = (
            args.destination_aws_secret_access_key
        )
    if args.destination_aws_endpoint:
        cfg["destination"]["s3"]["endpoint"] = args.destination_aws_endpoint

    if args.source_oci_uri:
        cfg["source"]["oci"]["uri"] = args.source_oci_uri
    if args.source_oci_username:
        cfg["source"]["oci"]["username"] = args.source_oci_username
    if args.source_oci_password:
        cfg["source"]["oci"]["password"] = args.source_oci_password

    if args.destination_oci_uri:
        cfg["destination"]["oci"]["uri"] = args.destination_oci_uri
    if args.destination_oci_username:
        cfg["destination"]["oci"]["username"] = args.destination_oci_username
    if args.destination_oci_password:
        cfg["destination"]["oci"]["password"] = args.destination_oci_password

    _validate_config(cfg)

    return cfg
