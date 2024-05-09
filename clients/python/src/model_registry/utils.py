"""Utilities for the model registry."""

from __future__ import annotations

import os
from collections import namedtuple
from typing import Callable

import grpc
from attr import dataclass
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


# https://github.com/grpc/grpc/blob/master/examples/python/interceptors/headers/generic_client_interceptor.py
@dataclass
class GenericClientInterceptor(  # noqa: D101
    grpc.UnaryUnaryClientInterceptor,
    grpc.UnaryStreamClientInterceptor,
    grpc.StreamUnaryClientInterceptor,
    grpc.StreamStreamClientInterceptor,
):
    fn: Callable

    def intercept_unary_unary(self, continuation, client_call_details, request):  # noqa: D102
        new_details, new_request_iterator, postprocess = self.fn(
            client_call_details, iter((request,)), False, False
        )
        response = continuation(new_details, next(new_request_iterator))
        return postprocess(response) if postprocess else response

    def intercept_unary_stream(self, continuation, client_call_details, request):  # noqa: D102
        new_details, new_request_iterator, postprocess = self.fn(
            client_call_details, iter((request,)), False, True
        )
        response_it = continuation(new_details, next(new_request_iterator))
        return postprocess(response_it) if postprocess else response_it

    def intercept_stream_unary(  # noqa: D102
        self, continuation, client_call_details, request_iterator
    ):
        new_details, new_request_iterator, postprocess = self.fn(
            client_call_details, request_iterator, True, False
        )
        response = continuation(new_details, new_request_iterator)
        return postprocess(response) if postprocess else response

    def intercept_stream_stream(  # noqa: D102
        self, continuation, client_call_details, request_iterator
    ):
        new_details, new_request_iterator, postprocess = self.fn(
            client_call_details, request_iterator, True, True
        )
        response_it = continuation(new_details, new_request_iterator)
        return postprocess(response_it) if postprocess else response_it


# https://github.com/grpc/grpc/blob/master/examples/python/interceptors/headers/header_manipulator_client_interceptor.py
# we need to subclass ClientCallDetails to add a constructor (it's ABC)
class ClientCallDetails(  # noqa: D101
    namedtuple("ClientCallDetails", ("method", "timeout", "metadata", "credentials")),
    grpc.ClientCallDetails,
):
    pass


def header_adder_interceptor(header, value):
    """Create a client interceptor that adds a header to requests."""

    def intercept_call(
        client_call_details,
        request_iterator,
        request_streaming,
        response_streaming,
    ):
        metadata = list(client_call_details.metadata or [])
        metadata.append(
            (
                header,
                value,
            )
        )
        return (
            ClientCallDetails(
                client_call_details.method,
                client_call_details.timeout,
                metadata,
                client_call_details.credentials,
            ),
            request_iterator,
            None,
        )

    return GenericClientInterceptor(intercept_call)
