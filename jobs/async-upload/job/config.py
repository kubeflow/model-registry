from __future__ import annotations
import json
import logging
import configargparse as cap
from typing import Any, Dict
from pathlib import Path

from .models import (
    AsyncUploadConfig,
    OCIStorageConfig,
    S3StorageConfig, 
    SourceConfig, 
    DestinationConfig, 
    ModelConfig, 
    StorageConfig, 
    RegistryConfig,
    S3Config,
    OCIConfig,
    URISourceConfig,
    SourceType,
    DestinationType,
    URISourceStorageConfig,
    UploadIntent
)

logger = logging.getLogger(__name__)


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
    p.add_argument("--source-type", choices=["s3", "oci", "uri"], default="s3")
    p.add_argument("--source-uri")
    p.add_argument("--source-aws-bucket")
    p.add_argument("--source-aws-key")
    p.add_argument("--source-aws-region")
    p.add_argument("--source-aws-access-key-id")
    p.add_argument("--source-aws-secret-access-key")
    p.add_argument("--source-aws-endpoint")
    # OCI registry
    p.add_argument("--source-oci-uri")
    p.add_argument("--source-oci-registry")
    p.add_argument("--source-oci-username")
    p.add_argument("--source-oci-password")

    # --- destination ---
    # s3
    # TODO: We should be able to infer the type from the credentials provided, therefore no default needed
    p.add_argument("--destination-type", choices=["s3", "oci"], default="oci")
    p.add_argument("--destination-aws-bucket")
    p.add_argument("--destination-aws-key")
    p.add_argument("--destination-aws-region")
    p.add_argument("--destination-aws-access-key-id")
    p.add_argument("--destination-aws-secret-access-key")
    p.add_argument("--destination-aws-endpoint")
    # OCI registry
    p.add_argument("--destination-oci-uri")
    p.add_argument("--destination-oci-registry")
    p.add_argument("--destination-oci-username")
    p.add_argument("--destination-oci-password")
    p.add_argument("--destination-oci-base-image", default="busybox:latest")
    # The `type` converter is needed here to support env-based booleans
    # See: https://github.com/bw2/ConfigArgParse/tree/master?tab=readme-ov-file#special-values
    p.add_argument("--destination-oci-enable-tls-verify", default=True, type=str2bool)

    # --- model-registry model data ---
    # This intent determines the action to take once the model has been uploaded to the destination.
    p.add_argument(
        "--model-upload-intent", 
        type=UploadIntent, 
        choices=tuple(UploadIntent), 
        default=UploadIntent.update_artifact
    )
    p.add_argument("--model-id")
    p.add_argument("--model-version-id")
    p.add_argument("--model-artifact-id")

    # --- model-storage configuration ---
    p.add_argument("--storage-path", default="/tmp/model-sync")

    # --- model-registry client ---
    p.add_argument("--registry-server-address")
    p.add_argument("--registry-port", default=443)
    p.add_argument("--registry-is-secure", default=True)
    p.add_argument("--registry-author")
    p.add_argument("--registry-user-token", default=None)
    p.add_argument("--registry-user-token-envvar", default=None)
    p.add_argument("--registry-custom-ca", default=None)
    p.add_argument("--registry-custom-ca-envvar", default=None)
    p.add_argument("--registry-log-level", default=logging.WARNING)

    # TODO: The type of credential should be inferrable from the `type` specified in the source/destination
    p.add_argument(
        "--source-s3-credentials-path",
        metavar="PATH",
    )
    p.add_argument(
        "--destination-s3-credentials-path",
        metavar="PATH",
    )
    p.add_argument(
        "--source-oci-credentials-path",
        metavar="PATH",
    )
    p.add_argument(
        "--destination-oci-credentials-path",
        metavar="PATH",
    )
    p.add_argument(
        "--source-uri-credentials-path",
        metavar="PATH",
    )

    return p


def _load_s3_credentials(path: str | Path, store: S3Config) -> None:
    """
    The path must be a folder, with a number of files that match a typical AWS config, ie - a Secret mounted to a pod with keys:
    - AWS_ACCESS_KEY_ID
    - AWS_SECRET_ACCESS_KEY
    - AWS_REGION, or AWS_DEFAULT_REGION, see: https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html#envvars-list:~:text=command%20line%20parameter.-,AWS_DEFAULT_REGION,-The%20Default%20region
    - AWS_ENDPOINT_URL, or AWS_S3_ENDPOINT
    - AWS_BUCKET, or AWS_S3_BUCKET

    These would be mounted as files with the names above and whose contents are the secret values.

    These keys are loaded into the config store

    Mutates the `store` in-place
    """

    logger.info(f"🔍 Loading S3 credentials from {path}")

    # Validate the path is a directory
    p = Path(path).expanduser()
    if not p.is_dir():
        raise FileNotFoundError(f"credentials folder not found: {p}")

    # Load the credentials from the files
    aws_access_key_file = p / "AWS_ACCESS_KEY_ID"
    if aws_access_key_file.exists():
        store.access_key_id = aws_access_key_file.read_text()
    else:
        logger.warning("AWS_ACCESS_KEY_ID not found in %s", p)

    aws_secret_key_file = p / "AWS_SECRET_ACCESS_KEY"
    if aws_secret_key_file.exists():
        store.secret_access_key = aws_secret_key_file.read_text()
    else:
        logger.warning("AWS_SECRET_ACCESS_KEY not found in %s", p)

    aws_region_file = p / "AWS_REGION"
    if aws_region_file.exists():
        store.region = aws_region_file.read_text()
    else:
        aws_region_file = p / "AWS_DEFAULT_REGION"
        if aws_region_file.exists():
            store.region = aws_region_file.read_text()
        else:
            logger.warning("AWS_REGION, or AWS_DEFAULT_REGION not found in %s", p)

    aws_endpoint_file = p / "AWS_ENDPOINT_URL"
    if aws_endpoint_file.exists():
        store.endpoint = aws_endpoint_file.read_text()
    else:
        aws_endpoint_file = p / "AWS_S3_ENDPOINT"
        if aws_endpoint_file.exists():
            store.endpoint = aws_endpoint_file.read_text()
        else:
            logger.warning("AWS_ENDPOINT_URL, or AWS_S3_ENDPOINT not found in %s", p)

    aws_bucket_file = p / "AWS_BUCKET"
    if aws_bucket_file.exists():
        store.bucket = aws_bucket_file.read_text()
    else:
        aws_bucket_file = p / "AWS_S3_BUCKET"
        if aws_bucket_file.exists():
            store.bucket = aws_bucket_file.read_text()
        else:
            logger.warning("AWS_BUCKET, or AWS_S3_BUCKET not found in %s", p)


def _load_oci_credentials(
    path: str | Path, store: OCIConfig
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
    logger.info(f"🔍 Loading OCI credentials from {path}")
    # Validate the path is a file
    p = Path(path).expanduser()
    if not p.is_file():
        raise FileNotFoundError(f"credentials file not found: {p}")

    # Load the credentials from the file
    docker_config = json.loads(p.read_text())

    # Validate the docker config is valid
    if "auths" not in docker_config:
        raise ValueError("Invalid docker config file")

    # Load the credentials from the docker config, the URI passed in via config is used as a key to find the correct credentials
    registry = store.registry
    auth = docker_config["auths"][registry]["auth"]
    # TODO: This might not be the correct way to parse this
    username, password = auth.split(":") if auth else (None, None)
    store.username = username
    store.password = password
    store.email = docker_config["auths"][registry]["email"]


def _load_uri_credentials(path: str | Path, store: URISourceConfig) -> None:
    """
    The path must be a folder containing a Secret mounted to a pod with a key "URI".

    For example, a Secret like:
    ```yaml
    kind: Secret
    apiVersion: v1
    metadata:
      name: my-uri-credential
    stringData:
      URI: hf:/some/repo
    ```

    This would be mounted as a file with name "URI" and whose contents are the URI value.

    Mutates the `store` in-place
    """
    logger.info(f"🔍 Loading URI credentials from {path}")

    # Validate the path is a directory
    p = Path(path).expanduser()
    if not p.is_dir():
        raise FileNotFoundError(f"credentials folder not found: {p}")

    # Load the URI from the file
    uri_file = p / "URI"
    if not uri_file.is_file():
        raise FileNotFoundError(f"URI credential file not found: {uri_file}")

    store.uri = uri_file.read_text().strip()


def str2bool(x):
    """Convert a config string to boolean. This is needed because configargparse doesn't support boolean optional action as env vars"""
    if isinstance(x, bool):
        return x
    val = x.lower()
    if val in ("yes", "y", "true", "t", "1"):
        return True
    if val in ("no", "n", "false", "f", "0"):
        return False
    raise ValueError(f"Invalid boolean value: {x!r}")


def get_config(argv: list[str] | None = None) -> AsyncUploadConfig:
    """
    Return a validated AsyncUploadConfig instance.

    Priority of the config is:
    1. Command-line arguments
    2. Environment variables
    3. Credentials files
    4. Default values
    """
    args = _parser().parse_args(argv)
    logger.debug("🔍 Command-line arguments: %s", args)
 
    # Create source config based on type
    if args.source_type == "s3":
        s3_config = S3Config(
            bucket=args.source_aws_bucket,
            key=args.source_aws_key,
            region=args.source_aws_region,
            access_key_id=args.source_aws_access_key_id,
            secret_access_key=args.source_aws_secret_access_key,
            endpoint=args.source_aws_endpoint,
        )
        # Load credentials from files, if provided
        if args.source_s3_credentials_path:
            _load_s3_credentials(args.source_s3_credentials_path, s3_config)
        source_config = S3StorageConfig(
            **s3_config.model_dump(),
            credentials_path=args.source_s3_credentials_path
        )
    elif args.source_type == "oci":
        source_config = OCIStorageConfig(
            uri=args.source_oci_uri,
            registry=args.source_oci_registry,
            username=args.source_oci_username,
            password=args.source_oci_password,
            email=None,
            credentials_path=args.source_oci_credentials_path
        )
        # Load credentials from files, if provided
        if args.source_oci_credentials_path:
            _load_oci_credentials(args.source_oci_credentials_path, source_config)

    elif args.source_type == "uri":
        uri_config = URISourceConfig(uri=args.source_uri)
        # Load credentials from files, if provided
        if args.source_uri_credentials_path:
            _load_uri_credentials(args.source_uri_credentials_path, uri_config)
        source_config = URISourceStorageConfig(
            **uri_config.model_dump(),
            credentials_path=args.source_uri_credentials_path
        )
    else:
        raise ValueError(f"Unsupported source type: {args.source_type}")

    # Create destination config based on type
    if args.destination_type == "s3":
        destination_config = S3StorageConfig(
            bucket=args.destination_aws_bucket,
            key=args.destination_aws_key,
            region=args.destination_aws_region,
            access_key_id=args.destination_aws_access_key_id,
            secret_access_key=args.destination_aws_secret_access_key,
            endpoint=args.destination_aws_endpoint,
            credentials_path=args.destination_s3_credentials_path
        )
        # Load credentials from files, if provided
        if args.destination_s3_credentials_path:
            _load_s3_credentials(args.destination_s3_credentials_path, destination_config)
    elif args.destination_type == "oci":
        destination_config = OCIStorageConfig(
            uri=args.destination_oci_uri,
            registry=args.destination_oci_registry,
            username=args.destination_oci_username,
            password=args.destination_oci_password,
            email=None,
            base_image=args.destination_oci_base_image,
            enable_tls_verify=args.destination_oci_enable_tls_verify,
            credentials_path=args.destination_oci_credentials_path
        )
        # Load credentials from files, if provided
        if args.destination_oci_credentials_path:
            _load_oci_credentials(args.destination_oci_credentials_path, destination_config)

    else:
        raise ValueError(f"Unsupported destination type: {args.destination_type}")

    # Create model instances
    try:
        config = AsyncUploadConfig(
            source=source_config,
            destination=destination_config,
            model=ModelConfig(
                upload_intent=args.model_upload_intent,
                id=args.model_id,
                version_id=args.model_version_id,
                artifact_id=args.model_artifact_id,
            ),
            storage=StorageConfig(
                path=args.storage_path,
            ),
            registry=RegistryConfig(
                server_address=args.registry_server_address,
                port=args.registry_port,
                is_secure=args.registry_is_secure,
                author=args.registry_author,
                user_token=args.registry_user_token,
                user_token_envvar=args.registry_user_token_envvar,
                custom_ca=args.registry_custom_ca,
                custom_ca_envvar=args.registry_custom_ca_envvar,
                log_level=args.registry_log_level,
            ),
        )
    except Exception as e:
        logger.error("❌ Configuration validation failed: %s", e)
        raise

    logger.info("📦 Configuration loaded successfully")

    # Log the configuration (with sensitive values sanitized)
    sanitized_cfg = _sanitize_config_for_logging(config.model_dump())
    logger.debug("🔍 Configuration loaded: %s", json.dumps(sanitized_cfg, indent=2))

    return config


def _sanitize_config_for_logging(cfg: Dict[str, Any]) -> Dict[str, Any]:
    """
    Create a sanitized copy of the config for logging purposes, masking sensitive values.
    """
    import copy

    sanitized = copy.deepcopy(cfg)

    # Mask sensitive S3 credentials in source
    if "source" in sanitized:
        source = sanitized["source"]
        if "secret_access_key" in source and source["secret_access_key"]:
            source["secret_access_key"] = "***"
        if "access_key_id" in source and source["access_key_id"]:
            source["access_key_id"] = "***"
        if "password" in source and source["password"]:
            source["password"] = "***"

    # Mask sensitive credentials in destination  
    if "destination" in sanitized:
        destination = sanitized["destination"]
        if "secret_access_key" in destination and destination["secret_access_key"]:
            destination["secret_access_key"] = "***"
        if "access_key_id" in destination and destination["access_key_id"]:
            destination["access_key_id"] = "***"
        if "password" in destination and destination["password"]:
            destination["password"] = "***"

    # Mask sensitive registry credentials
    if "registry" in sanitized:
        registry = sanitized["registry"]
        if "user_token" in registry and registry["user_token"]:
            registry["user_token"] = "***"
        if "custom_ca" in registry and registry["custom_ca"]:
            registry["custom_ca"] = "***"

    return sanitized
