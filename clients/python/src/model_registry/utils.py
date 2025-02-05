"""Utilities for the model registry."""

from __future__ import annotations

import os
import re

from typing_extensions import overload

from ._utils import required_args
from .exceptions import MissingMetadata


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

s3_prefix = "s3://"

def is_s3_uri(uri: str):
    """Checks whether a string is a valid S3 URI
    
    This helper function checks whether the string starts with the correct s3 prefix (s3://) and
    whether the string contains both a bucket and a key.
    
    Args:
        uri: The URI to check
        
    Returns:
        Boolean indicating whether it is a valid S3 URI
    """
    if not uri.startswith(s3_prefix):
        return False
    # Slice the uri from prefix onward, then check if there are 2 components when splitting on "/"
    path = uri[len(s3_prefix) :]
    if len(path.split("/", 1)) != 2:
        return False
    return True

oci_pattern = r'^oci://(?P<host>[^/]+)/(?P<repository>[A-Za-z0-9_\-/]+)(:(?P<tag>[A-Za-z0-9_.-]+))?$'

def is_oci_uri(uri: str):
    """Checks whether a string is a valid OCI URI
    
    The expected format is:
        oci://<host>/<repository>[:<tag>]
        
    Examples of valid URIs:
        oci://registry.example.com/my-namespace/my-repo:latest
        oci://localhost:5000/my-repo

    Args:
        uri: The URI to check
        
    Returns:
        Boolean indicating whether it is a valid OCI URI
    """
    return re.match(oci_pattern, uri) is not None
