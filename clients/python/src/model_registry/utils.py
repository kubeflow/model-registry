"""Utilities for the model registry."""

from __future__ import annotations

import base64
import json
import os
import tempfile
from dataclasses import dataclass
from pathlib import Path
from subprocess import CalledProcessError
from typing import Callable, Protocol, TypeVar

from typing_extensions import Literal, overload

from ._utils import required_args
from .exceptions import MissingMetadata, StoreError

# Generic return type
T = TypeVar("T")


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
        msg = "endpoint must be provided for non-default bucket"
        raise MissingMetadata(msg)

    endpoint = endpoint or os.getenv("AWS_S3_ENDPOINT")
    region = region or os.getenv("AWS_DEFAULT_REGION")
    if not endpoint:
        msg = "Missing environment variable: `bucket` is required"
        raise MissingMetadata(msg)

    # https://alexwlchan.net/2020/s3-keys-are-not-file-paths/ nor do they resolve to valid URls
    # FIXME: is this safe?
    if not region:
        return f"s3://{bucket}/{path}?endpoint={endpoint}"
    return f"s3://{bucket}/{path}?endpoint={endpoint}&defaultRegion={region}"


class PullFn(Protocol):
    """Pull function definition."""

    def __call__(self, base: str, dest: Path, **kwargs) -> T: ...  # noqa: D102


class PushFn(Protocol):
    """Push function definition."""

    def __call__(self, src: Path, oci_ref: str, **kwargs) -> T: ...  # noqa: D102


@dataclass
class BackendDefinition:
    """Holds the 3 core callables for a backend.

    - is_available() -> bool
    - pull(base_image: str, dest_dir: Path, **kwargs) -> T
    - push(local_image_path: Path, oci_ref: str, **kwargs) -> T.
    """

    is_available: Callable[[], bool]
    pull: PullFn
    push: PushFn


def _kwargs_to_params(kwargs: dict[str, str]) -> list[str]:
    """Convert kwargs to list of params.

    Args:
        kwargs: The keyword args dict.
    """
    args = []
    for k, v in kwargs.items():
        args.append(f"{k}")
        args.append(str(v))
    return args


def _get_skopeo_backend(
    pull_args: list[str] | None = None, push_args: list[str] | None = None
) -> BackendDefinition:
    try:
        from olot.backend.skopeo import is_skopeo, skopeo_pull, skopeo_push
    except ImportError as e:
        msg = "Could not import 'olot.backend.skopeo'. Ensure that 'olot' is installed if you want to use the 'skopeo' backend."
        raise ImportError(msg) from e

    def wrapped_pull(base_image: str, dest: Path, **kwargs) -> T:
        kwargs = _backend_specific_params("skopeo", "pull", **kwargs)
        params = _kwargs_to_params(kwargs)
        params.extend(pull_args or [])

        return _scrub_errors(lambda: skopeo_pull(base_image, dest, params))

    def wrapped_push(src: Path, oci_ref: str, **kwargs) -> T:
        kwargs = _backend_specific_params("skopeo", "push", **kwargs)
        params = _kwargs_to_params(kwargs)
        params.extend(push_args or [])

        return _scrub_errors(lambda: skopeo_push(src, oci_ref, params))

    return BackendDefinition(
        is_available=is_skopeo, pull=wrapped_pull, push=wrapped_push
    )


def _get_oras_backend(
    pull_args: list[str] | None = None, push_args: list[str] | None = None
) -> BackendDefinition:
    try:
        from olot.backend.oras_cp import is_oras, oras_pull, oras_push
    except ImportError as e:
        msg = "Could not import 'olot.backend.oras_cp'. Ensure that 'olot' is installed if you want to use the 'oras_cp' backend."
        raise ImportError(msg) from e

    def wrapped_pull(base_image: str, dest: Path, **kwargs) -> T:
        kwargs = _backend_specific_params("oras", "pull", **kwargs)
        params = _kwargs_to_params(kwargs)
        params.extend(pull_args or [])

        return _scrub_errors(lambda: oras_pull(base_image, dest, params))

    def wrapped_push(src: Path, oci_ref: str, **kwargs) -> T:
        kwargs = _backend_specific_params("oras", "push", **kwargs)
        params = _kwargs_to_params(kwargs)
        params.extend(push_args or [])

        return _scrub_errors(lambda: oras_push(src, oci_ref, params))

    return BackendDefinition(
        is_available=is_oras,
        pull=wrapped_pull,
        push=wrapped_push,
    )


def _backend_specific_params(
    backend: Literal["skopeo", "oras"], type: Literal["push", "pull"], **kwargs
) -> dict:
    """Generate params based on the backend and action.

    Args:
        backend: The backend to use supported in olot.
        type: The action to perform.
        kwargs: Additional args defined below.

    Keyword Args:
        username: the usrername of the registry.
        password: the password of the registry.
    """
    # Determine backend
    if backend == "skopeo":
        prefix = "--src" if type == "pull" else "--dest"
    elif backend == "oras":
        prefix = "--from" if type == "pull" else "--to"

    # Actual param specifications
    if username := kwargs.pop("username", None):
        kwargs[f"{prefix}-username"] = username

    if password := kwargs.pop("password", None):
        kwargs[f"{prefix}-password"] = password

    return kwargs


def _scrub_errors(func: Callable[[], T]) -> T:
    """Scrub errors of any subprocess command with sensitive data.

    Args:
        func: A partial or lambda function that has not been yet executed.
    """
    try:
        return func()
    except (CalledProcessError, Exception) as e:
        msg = """Problem with command"""
        raise RuntimeError(msg, e.returncode, e.stderr) from None


@dataclass
class OCIParams:
    """Parameters for the OCI client to perform the upload.

    Allows for some customization of how to perform the upload step when uploading via OCI
    """

    base_image: str
    oci_ref: str
    dest_dir: str | os.PathLike = None
    backend: str = "skopeo"
    modelcard: os.PathLike | None = None
    custom_oci_backend: BackendDefinition = None
    oci_auth_env_var: str | None = None
    oci_username: str | None = None
    oci_password: str | None = None


@dataclass
class S3Params:
    """Parameters for the S3 Client (boto3) to perform the upload.

    Allows for some amount of customization when performing an upload, such as providing a custom endpoint url, access keys, etc.
    """

    bucket_name: str
    s3_prefix: str
    endpoint_url: str | None = None
    access_key_id: str | None = None
    secret_access_key: str | None = None
    region: str | None = None


# A dict mapping backend names to their definitions
BackendDict = dict[str, Callable[[], BackendDefinition]]

DEFAULT_BACKENDS: BackendDict = {
    "skopeo": _get_skopeo_backend,
    "oras": _get_oras_backend,
}


def save_to_oci_registry(  # noqa: C901 ( complex args >8 )
    base_image: str,
    oci_ref: str,
    model_files_path: str | os.PathLike,
    dest_dir: str | os.PathLike = None,
    backend: str = "skopeo",
    modelcard: os.PathLike | None = None,
    custom_oci_backend: BackendDefinition | None = None,
    oci_auth_env_var: str | None = None,
    oci_username: str | None = None,
    oci_password: str | None = None,
) -> str:
    """Appends a list of files to an OCI-based image.

    Args:
        base_image: The image to append model files to. This image will be downloaded to the location at `dest_dir`
        dest_dir: The location to save the downloaded and extracted base image to.
        oci_ref: Destination of where to push the newly layered image to. eg, "quay.io/my-org/my-registry:1.0.0"
        model_files_path: Path to the files to add to the base_image as layers
        backend: The CLI tool to use to perform the oci image pull/push. One of: "skopeo", "oras"
        modelcard: [Optional] Path to the modelcard to additionally include as a layer
        custom_oci_backend: [Optional] If you would like to use your own OCI Backend layer, you can provide it here
        oci_auth_env_var: [Optional] The environment variable that holds the auth/config JSON for OCI registry auth.
        oci_username: [Optional] The username to the OCI registry.
        oci_password: [Optional] (Must be used with OCI username) The password to the OCI registry.

    Raises:
        ValueError: If the chosen backend is not installed on the host
        ValueError: If the chosen backend is an invalid option
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

    # Check for OCI Auth Env and a default
    auth: str = None
    if oci_auth_env_var:
        env_value = _validate_env_var(oci_auth_env_var)
        auth = _extract_auth_json(env_value)
    else:
        try:
            env_value = _validate_env_var(".dockerconfigjson")
            auth = _extract_auth_json(env_value)
        except ValueError:
            pass

    # If a custom backend is provided, use it, else fetch the backend out of the registry
    if custom_oci_backend:
        backend_def = custom_oci_backend
    elif backend in DEFAULT_BACKENDS:
        # Fetching the backend definition can throw an error, but it should bubble up as it has the appropriate messaging
        backend_def = DEFAULT_BACKENDS[backend]()
    else:
        msg = f"'{backend}' is not an available backend to use. Available backends: {DEFAULT_BACKENDS.keys()}"
        raise ValueError(msg)
    if not backend_def.is_available():
        msg = f"Backend '{backend}' is selected, but not available on the system. Ensure the dependencies for '{backend}' are installed in your environment."
        raise ValueError(msg)

    if dest_dir is None:
        dest_dir = tempfile.mkdtemp()
    local_image_path = Path(dest_dir)

    # Set params
    params = {}

    # User/pass
    if auth:
        usr_pass = auth.split(":")
        params["username"] = usr_pass[0]
        params["password"] = usr_pass[-1]
    elif oci_username and oci_password:
        params["username"] = oci_username
        params["password"] = oci_password

    backend_def.pull(base_image, local_image_path, **params)
    # Extract the absolute path from the files found in the path
    files = [file[0] for file in get_files_from_path(model_files_path)]
    oci_layers_on_top(local_image_path, files, modelcard)
    backend_def.push(local_image_path, oci_ref, **params)

    # Return the OCI URI

    return f"oci://{oci_ref}"


def get_files_from_path(path: str) -> list[tuple[str, str]]:
    """Given a path, get the list of files.

    If the path points to a single file, that file's absolute_path and filename will be the only entry returned

    If the path points to a directory, the directory will be walked to fetch all the absolute and relative filepaths for each file

    Args:
        path: Location (directory or file) to extract filenames from

    Returns:
        A list of 2-entry tuples containing (absolute_path, relative_path) from the path provided

    Raises:
        ValueError: If the path provided does not already exist.
    """
    if not os.path.exists(path):
        msg = f"Path '{path}' does not exist. Please ensure path is correct."
        raise ValueError(msg)

    files = []

    is_file = os.path.isfile(path)
    if is_file:
        # When just a single file, return it
        filename = os.path.basename(path)
        file = (path, filename)
        files.append(file)
        return files

    for root, _, filenames in os.walk(path):
        for filename in filenames:
            absolute_path = os.path.join(root, filename)
            relative_path = os.path.relpath(absolute_path, path)
            files.append((absolute_path, relative_path))

    return files


def _validate_env_var(var: str) -> str:
    """Validate that an env var exists.

    Args:
        var: The env var to lookup.
    """
    if not (env_var := os.getenv(var)):
        msg = f"Cannot find environment variable '{var}'."
        raise ValueError(msg)
    return env_var


def _extract_auth_json(auth_data: str) -> str:
    """Extract the auth JSON from a string value.

    Args:
        auth_data: The Auth JSON string.
    """
    try:
        auth_json = json.loads(auth_data)
        if type(auth_json) is not dict:
            msg = ""
            raise TypeError(msg)
        registries = auth_json["auths"]
        reg_keys = list(registries.keys())
        if len(reg_keys) > 1:
            msg = f"Auth JSON has multiple registries ({', '.join(reg_keys)}). This is not supported."
            raise ValueError(msg)

        key = registries[reg_keys[0]]["auth"]
        auth = base64.b64decode(key)
        return auth.decode()

    except (AttributeError, KeyError) as e:
        msg = "This is an invalid Auth JSON."
        raise ValueError(msg) from e
    except json.JSONDecodeError as e:
        invalid_json_msg = "Auth data does not contain valid JSON."
        raise ValueError(invalid_json_msg) from e
