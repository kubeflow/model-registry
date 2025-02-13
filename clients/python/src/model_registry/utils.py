"""Utilities for the model registry."""

from __future__ import annotations

import os
import pathlib

from typing_extensions import overload

from ._utils import required_args
from .exceptions import MissingMetadata, StoreError


@overload
def s3_uri_from(
    path: str,
) -> str: ...


@overload
def s3_uri_from(
    path: str,
    bucket: str,
) -> str: ...


@overload
def s3_uri_from(
    path: str,
    bucket: str,
    *,
    endpoint: str,
    region: str,
) -> str: ...


@required_args(
    (),
    (  # pre-configured env
        "bucket",
    ),
    (  # custom env or non-default bucket
        "bucket",
        "endpoint",
        "region",
    ),
)
def s3_uri_from(
    path: str,
    bucket: str | None = None,
    *,
    endpoint: str | None = None,
    region: str | None = None,
) -> str:
    """Build an S3 URI.

    This helper function builds an S3 URI from a path and a bucket name, assuming you have a configured environment
    with a default bucket, endpoint, and region set.
    If you don't, you must provide all three optional arguments.
    That is also the case for custom environments, where the default bucket is not set, or if you want to use a
    different bucket.

    Args:
        path: Storage path.
        bucket: Name of the S3 bucket. Defaults to AWS_S3_BUCKET.
        endpoint: Endpoint of the S3 bucket. Defaults to AWS_S3_ENDPOINT.
        region: Region of the S3 bucket. Defaults to AWS_DEFAULT_REGION.

    Returns:
        S3 URI.
    """
    default_bucket = os.environ.get("AWS_S3_BUCKET")
    if not bucket:
        if not default_bucket:
            msg = "Custom environment requires all arguments"
            raise MissingMetadata(msg)
        bucket = default_bucket
    elif (not default_bucket or default_bucket != bucket) and not endpoint:
        msg = (
            "bucket_endpoint and bucket_region must be provided for non-default bucket"
        )
        raise MissingMetadata(msg)

    endpoint = endpoint or os.getenv("AWS_S3_ENDPOINT")
    region = region or os.getenv("AWS_DEFAULT_REGION")

    if not (endpoint and region):
        msg = "Missing environment variables: bucket_endpoint and bucket_region are required"
        raise MissingMetadata(msg)

    # https://alexwlchan.net/2020/s3-keys-are-not-file-paths/ nor do they resolve to valid URls
    # FIXME: is this safe?
    return f"s3://{bucket}/{path}?endpoint={endpoint}&defaultRegion={region}"

def save_to_oci_registry(
        base_image: str,
        dest_dir: str | os.PathLike,
        oci_ref: str,
        model_files: list[os.PathLike],
        backend: str = "skopeo",
        modelcard: os.PathLike | None = None,
):
    """Appends a list of files to an OCI-based image.

    Args:
        base_image: The image to append model files to. This image will be downloaded to the location at `dest_dir`
        dest_dir: The location to save the downloaded and extracted base image to.
        oci_ref: Destination of where to push the newly layered image to
        model_files: List of files to add to the base_image as layers
        backend: The CLI tool to use to perform the oci image pull/push. One of: "skopeo", "oras"
        modelcard: Optional, path to the modelcard to additionally include as a layer

    Raises:
        ValueError: If the chosen backend is not installed on the host
        StoreError: If the chosen backend is an invalid option
        StoreError: If `olot` is not installed as a python package
    Returns:
        None.
    """
    try:
        from olot.basics import oci_layers_on_top
    except ImportError as e:
        msg = """Package `olot` is not installed.
To save models to OCI compatible storage, start by installing the `olot` package, either directly or as an
extra (available as `model-registry[olot]`), e.g.:
```sh
!pip install --pre model-registry[olot]
```
or
```sh
!pip install olot
```
        """
        raise StoreError(msg) from e

    local_image_path = pathlib.Path(dest_dir)

    if backend == "skopeo":
        from olot.backend.skopeo import is_skopeo, skopeo_pull, skopeo_push

        if not is_skopeo():
            msg = "skopeo is selected, but it is not present on the machine. Please validate the skopeo cli is installed and available in the PATH"
            raise ValueError(msg)

        skopeo_pull(base_image, local_image_path)
        oci_layers_on_top(local_image_path, model_files, modelcard)
        skopeo_push(dest_dir, oci_ref)

    elif backend == "oras":
        from olot.backend.oras_cp import is_oras, oras_pull, oras_push
        if not is_oras():
            msg = "oras is selected, but it is not present on the machine. Please validate the oras cli is installed and available in the PATH"
            raise ValueError(msg)

        oras_pull(base_image, local_image_path)
        oci_layers_on_top(local_image_path, model_files, modelcard)
        oras_push(local_image_path, oci_ref)

    else:
        msg = f"Invalid backend chosen: '{backend}'"
        raise StoreError(msg)
