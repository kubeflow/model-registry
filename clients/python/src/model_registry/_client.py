"""Standard client for the model registry."""

from __future__ import annotations

import logging
import os
from collections.abc import Mapping
from dataclasses import asdict
from pathlib import Path
from typing import TYPE_CHECKING, Any, TypeVar, Union, get_args
from warnings import warn

from .core import ModelRegistryAPIClient
from .exceptions import StoreError
from .types import (
    ListOptions,
    ModelArtifact,
    ModelVersion,
    Pager,
    RegisteredModel,
    SupportedTypes,
)
from .utils import OCIParams, S3Params, save_to_oci_registry, s3_uri_from

ModelTypes = Union[RegisteredModel, ModelVersion, ModelArtifact]
TModel = TypeVar("TModel", bound=ModelTypes)

logging.basicConfig(
    format="%(asctime)s.%(msecs)03d - %(name)s:%(levelname)s: %(message)s",
    datefmt="%H:%M:%S",
    level=logging.WARNING,  # the default loglevel
    handlers=[
        # logging.FileHandler(
        #     LOGS
        #     / "log-{}-{}.log".format(
        #         datetime.now(tz=datetime.now().astimezone().tzinfo).strftime(
        #             "%Y-%m-%d-%H-%M-%S"
        #         ),
        #         os.getpid(),
        #     ),
        #     encoding="utf-8",
        #     delay=False,
        # ),
        logging.StreamHandler(),
    ],
)

logger = logging.getLogger("model-registry")

DEFAULT_USER_TOKEN_ENVVAR = "KF_PIPELINES_SA_TOKEN_PATH"  # noqa: S105

# If we want to forward reference
if TYPE_CHECKING:
    from botocore.client import BaseClient


class ModelRegistry:
    """Model registry client."""

    def __init__(
        self,
        server_address: str,
        port: int = 443,
        *,
        author: str,
        is_secure: bool = True,
        user_token: str | None = None,
        user_token_envvar: str = DEFAULT_USER_TOKEN_ENVVAR,
        custom_ca: str | None = None,
        custom_ca_envvar: str | None = None,
        log_level: int = logging.WARNING,
    ):
        """Constructor.

        Args:
            server_address: Server address.
            port: Server port. Defaults to 443.

        Keyword Args:
            author: Name of the author.
            is_secure: Whether to use a secure connection. Defaults to True.
            user_token: The PEM-encoded user token as a string.
            user_token_envvar: Environment variable to read the user token from if it's not passed as an arg. Defaults to KF_PIPELINES_SA_TOKEN_PATH.
            custom_ca: Path to the PEM-encoded root certificates as a string.
            custom_ca_envvar: Environment variable to read the custom CA from if it's not passed as an arg.
            log_level: Log level. Defaults to logging.WARNING.
        """
        logger.setLevel(log_level)

        import nest_asyncio

        logger.debug("Setting up reentrant async event loop")
        nest_asyncio.apply()

        # TODO: get remaining args from env
        self._author = author

        if not user_token and user_token_envvar:
            logger.info("Reading user token from %s", user_token_envvar)
            # /var/run/secrets/kubernetes.io/serviceaccount/token
            if sa_token := os.environ.get(user_token_envvar):
                if user_token_envvar == DEFAULT_USER_TOKEN_ENVVAR:
                    logger.warning(
                        f"Sourcing user token from default envvar: {DEFAULT_USER_TOKEN_ENVVAR}"
                    )
                user_token = Path(sa_token).read_text()
            else:
                warn("User access token is missing", stacklevel=2)

        if is_secure:
            if (
                not custom_ca
                and custom_ca_envvar
                and (cert := os.getenv(custom_ca_envvar))
            ):
                logger.info(
                    "Using custom CA envvar %s",
                    custom_ca_envvar,
                )
                custom_ca = cert
                # client might have a default CA setup

            if not user_token:
                msg = "user token must be provided for secure connection"
                raise StoreError(msg)

            self._api = ModelRegistryAPIClient.secure_connection(
                server_address, port, user_token=user_token, custom_ca=custom_ca
            )
        else:
            self._api = ModelRegistryAPIClient.insecure_connection(
                server_address, port, user_token
            )
        self.get_registered_models().page_size(1)._next_page()

    def async_runner(self, coro: Any) -> Any:
        import asyncio

        try:
            loop = asyncio.get_event_loop()
        except RuntimeError:
            loop = asyncio.new_event_loop()
            asyncio.set_event_loop(loop)
        return loop.run_until_complete(coro)

    async def _register_model(self, name: str, **kwargs) -> RegisteredModel:
        if rm := await self._api.get_registered_model_by_params(name):
            return rm

        return await self._api.upsert_registered_model(
            RegisteredModel(name=name, **kwargs)
        )

    async def _register_new_version(
        self, rm: RegisteredModel, version: str, author: str, /, **kwargs
    ) -> ModelVersion:
        assert rm.id is not None, "Registered model must have an ID"
        if await self._api.get_model_version_by_params(rm.id, version):
            msg = f"Version {version} already exists"
            raise StoreError(msg)

        return await self._api.upsert_model_version(
            ModelVersion(name=version, author=author, **kwargs), rm.id
        )

    async def _register_model_artifact(
        self, mv: ModelVersion, name: str, uri: str, /, **kwargs
    ) -> ModelArtifact:
        assert mv.id is not None, "Model version must have an ID"
        return await self._api.upsert_model_version_artifact(
            ModelArtifact(name=name, uri=uri, **kwargs), mv.id
        )

    def _upload_to_s3(
        self,
        artifact_local_path: str,
        destination: str,
        region_name: str | None = None,
    ) -> str:
        """Uploads a file to an S3 bucket."""
        msg = "S3 Upload is not supported"
        raise StoreError(msg)

    def upload_artifact_and_register_model(
        self,
        name: str,
        model_files: list[os.PathLike],
        *,
        # Artifact/Model Params
        version: str,
        model_format_name: str,
        model_format_version: str,
        storage_path: str | None = None,
        storage_key: str | None = None,
        service_account_name: str | None = None,
        author: str | None = None,
        owner: str | None = None,
        description: str | None = None,
        metadata: Mapping[str, SupportedTypes] | None = None,
        # Upload/client Params
        upload_params: OCIParams | S3Params,
    ) -> RegisteredModel:
        if isinstance(upload_params, S3Params):
            # TODO: Dont use a mock function here
            destination_uri = self._upload_to_s3(model_files, destination_uri)
        elif isinstance(upload_params, OCIParams):
            print(asdict(upload_params))
            destination_uri = save_to_oci_registry(**asdict(upload_params), model_files=model_files)
        else:
            msg = 'Param "upload_params" is required to perform an upload. Please ensure the value provided is valid'
            raise ValueError(msg)

        return self.register_model(
            name,
            destination_uri,
            model_format_name=model_format_name,
            model_format_version=model_format_version,
            version=version,
            storage_key=storage_key,
            storage_path=storage_path,
            service_account_name=service_account_name,
            author=author,
            owner=owner,
            description=description,
            metadata=metadata,
        )


    def register_model(
        self,
        name: str,
        uri: str,
        *,
        model_format_name: str,
        model_format_version: str,
        version: str,
        storage_key: str | None = None,
        storage_path: str | None = None,
        service_account_name: str | None = None,
        author: str | None = None,
        owner: str | None = None,
        description: str | None = None,
        metadata: Mapping[str, SupportedTypes] | None = None,
    ) -> RegisteredModel:
        """Register a model.

        This registers a model in the model registry. The model is not downloaded, and has to be stored prior to
        registration.

        Most models can be registered using their URI, along with optional connection-specific parameters, `storage_key`
        and `storage_path` or, simply a `service_account_name`.
        URI builder utilities are recommended when referring to specialized storage; for example `utils.s3_uri_from`
        helper when using S3 object storage data connections.

        Args:
            name: Name of the model.
            uri: URI of the model.

        Keyword Args:
            version: Version of the model. Has to be unique.
            model_format_name: Name of the model format.
            model_format_version: Version of the model format.
            description: Description of the model.
            author: Author of the model. Defaults to the client author.
            owner: Owner of the model. Defaults to the client author.
            storage_key: Storage key.
            storage_path: Storage path.
            service_account_name: Service account name.
            metadata: Additional version metadata. Defaults to values returned by `default_metadata()`.

        Returns:
            Registered model.
        """
        rm = self.async_runner(self._register_model(name, owner=owner or self._author))
        mv = self.async_runner(
            self._register_new_version(
                rm,
                version,
                author or self._author,
                description=description,
                custom_properties=metadata or {},
            )
        )
        self.async_runner(
            self._register_model_artifact(
                mv,
                name,
                uri,
                model_format_name=model_format_name,
                model_format_version=model_format_version,
                storage_key=storage_key,
                storage_path=storage_path,
                service_account_name=service_account_name,
            )
        )

        return rm

    def update(self, model: TModel) -> TModel:
        """Update a model."""
        if not model.id:
            msg = "Model must have an ID"
            raise StoreError(msg)
        if not isinstance(model, get_args(ModelTypes)):
            msg = f"Model must be one of {get_args(ModelTypes)}"
            raise StoreError(msg)
        if isinstance(model, RegisteredModel):
            return self.async_runner(self._api.upsert_registered_model(model))
        if isinstance(model, ModelVersion):
            return self.async_runner(self._api.upsert_model_version(model, None))
        return self.async_runner(self._api.upsert_model_artifact(model))

    def register_hf_model(
        self,
        repo: str,
        path: str,
        *,
        version: str,
        model_format_name: str,
        model_format_version: str,
        author: str | None = None,
        owner: str | None = None,
        model_name: str | None = None,
        description: str | None = None,
        git_ref: str = "main",
    ) -> RegisteredModel:
        """Register a Hugging Face model.

        This imports a model from Hugging Face hub and registers it in the model registry.
        Note that the model is not downloaded.

        Args:
            repo: Name of the repository from Hugging Face hub.
            path: URI of the model.

        Keyword Args:
            version: Version of the model. Has to be unique.
            model_format_name: Name of the model format.
            model_format_version: Version of the model format.
            author: Author of the model. Defaults to repo owner.
            owner: Owner of the model. Defaults to the client author.
            model_name: Name of the model. Defaults to the repo name.
            description: Description of the model.
            git_ref: Git reference to use. Defaults to `main`.

        Returns:
            Registered model.
        """
        try:
            from huggingface_hub import HfApi, hf_hub_url, utils
        except ImportError as e:
            msg = """package `huggingface-hub` is not installed.
            To import models from Hugging Face Hub, start by installing the `huggingface-hub` package, either directly or as an
            extra (available as `model-registry[hf]`), e.g.:
            ```sh
            !pip install --pre model-registry[hf]
            ```
            or
            ```sh
            !pip install huggingface-hub
            ```
            """
            raise StoreError(msg) from e

        api = HfApi()
        try:
            model_info = api.model_info(repo, revision=git_ref)
        except utils.RepositoryNotFoundError as e:
            msg = f"Repository {repo} does not exist"
            raise StoreError(msg) from e
        except utils.RevisionNotFoundError as e:
            # TODO: as all hf-hub client calls default to using main, should we provide a tip?
            msg = f"Revision {git_ref} does not exist"
            raise StoreError(msg) from e

        if not author:
            # model author can be None if the repo is in a "global" namespace (i.e. no / in repo).
            if model_info.author is None:
                model_author = "unknown"
                warn(
                    "Model author is unknown. This is likely because the model is in a global namespace.",
                    stacklevel=2,
                )
            else:
                model_author = model_info.author
        else:
            model_author = author
        source_uri = hf_hub_url(repo, path, revision=git_ref)
        metadata = {
            "repo": repo,
            "source_uri": source_uri,
            "model_origin": "huggingface_hub",
            "model_author": model_author,
        }
        # card_data is the new field, but let's use the old one for backwards compatibility.
        if card_data := model_info.cardData:
            metadata.update(
                {
                    k: v
                    for k, v in card_data.to_dict().items()
                    # TODO: (#151) preserve tags, possibly other complex metadata
                    if isinstance(v, get_args(SupportedTypes))
                }
            )
        return self.register_model(
            model_name or model_info.id,
            source_uri,
            author=author or model_author,
            owner=owner or self._author,
            version=version,
            model_format_name=model_format_name,
            model_format_version=model_format_version,
            description=description,
            storage_path=path,
            metadata=metadata,
        )

    def get_registered_model(self, name: str) -> RegisteredModel | None:
        """Get a registered model.

        Args:
            name: Name of the model.

        Returns:
            Registered model.
        """
        return self.async_runner(self._api.get_registered_model_by_params(name))

    def get_model_version(self, name: str, version: str) -> ModelVersion | None:
        """Get a model version.

        Args:
            name: Name of the model.
            version: Version of the model.

        Returns:
            Model version.

        Raises:
            StoreException: If the model does not exist.
        """
        if not (rm := self.get_registered_model(name)):
            msg = f"Model {name} does not exist"
            raise StoreError(msg)
        assert rm.id
        return self.async_runner(self._api.get_model_version_by_params(rm.id, version))

    def get_model_artifact(self, name: str, version: str) -> ModelArtifact | None:
        """Get a model artifact.

        Args:
            name: Name of the model.
            version: Version of the model.

        Returns:
            Model artifact.

        Raises:
            StoreException: If either the model or the version don't exist.
        """
        if not (mv := self.get_model_version(name, version)):
            msg = f"Version {version} does not exist"
            raise StoreError(msg)
        assert mv.id
        return self.async_runner(self._api.get_model_artifact_by_params(name, mv.id))

    def get_registered_models(self) -> Pager[RegisteredModel]:
        """Get a pager for registered models.

        Returns:
            Iterable pager for registered models.
        """

        def rm_list(options: ListOptions) -> list[RegisteredModel]:
            return self.async_runner(self._api.get_registered_models(options))

        return Pager[RegisteredModel](rm_list)

    def get_model_versions(self, name: str) -> Pager[ModelVersion]:
        """Get a pager for model versions.

        Args:
            name: Name of the model.

        Returns:
            Iterable pager for model versions.

        Raises:
            StoreException: If the model does not exist.
        """
        if not (rm := self.get_registered_model(name)):
            msg = f"Model {name} does not exist"
            raise StoreError(msg)

        def rm_versions(options: ListOptions) -> list[ModelVersion]:
            # type checkers can't restrict the type inside a nested function: https://mypy.readthedocs.io/en/stable/common_issues.html#narrowing-and-inner-functions
            assert rm.id
            return self.async_runner(self._api.get_model_versions(rm.id, options))

        return Pager[ModelVersion](rm_versions)

    def save_to_s3(
        self,
        path: str,
        bucket_name: str,
        s3_prefix: str,
        *,
        endpoint_url: str | None = None,
        access_key_id: str | None = None,
        secret_access_key: str | None = None,
        region: str | None = None,
    ) -> str:
        """Saves a model to an S3 compatible storage.

        Args:
            path: Location to where the model(s) or artifact(s) are located. \
                Can recursively upload nested files in folders.
            bucket_name: The bucket to use for the S3 compatible object storage.
            s3_prefix: The path to prefix under root of bucket.

        Keyword Args:
            endpoint_url: The endpoint URL for the S3 comaptible storage if not using AWS S3.
            access_key_id: The S3 compatible object storage access ID.
            secret_access_key: The S3 compatible object storage secret access key.
            region: The region name for the S3 object storage.

        Returns:
            The S3 URI to the uploaded files.

        Raises:
            StoreError: If there was an issue uploading to S3.
        """
        # Get mixed credentials
        endpoint_url, access_key_id, secret_access_key, region = self.__s3_creds(
            endpoint_url, access_key_id, secret_access_key, region
        )

        s3 = self.__connect_to_s3(
            endpoint_url=endpoint_url,
            access_key_id=access_key_id,
            secret_access_key=secret_access_key,
        )
        try:
            return self.__upload_to_s3(
                path=path,
                path_prefix=s3_prefix,
                bucket=bucket_name,
                s3=s3,
                endpoint_url=endpoint_url,
                region=region,
            )
        except Exception as e:
            raise e

    def __upload_to_s3(  # noqa: C901
        self,
        path: str,
        bucket: str,
        s3: BaseClient,
        path_prefix: str,
        *,
        endpoint_url: str | None = None,
        region: str | None = None,
    ) -> str:
        """Internal method for recursively uploading all files to S3.

        Args:
            path: The path to where the models or artifacts are.
            bucket: The name of the S3 bucket.
            s3: The S3 Client object.
            path_prefix: The folder prefix to store under the root of the bucket.

        Keyword Args:
            endpoint_url: The endpoint url for the S3 bucket.
            region: The region name for the S3 bucket.

        Returns:  The S3 URI path of the uploaded files.

        Raises:
            StoreError if `path` does not exist.
            StoreError if `path_prefix` is not set.
        """
        is_file = os.path.isfile(path)

        if not path_prefix:
            msg = "`path_prefix` must be set."
            raise StoreError(msg)

        if not os.path.exists(path):
            msg = f"Path '{path}' does not exist. Please ensure path is correct."
            raise StoreError(msg)

        if path_prefix.endswith("/"):
            path_prefix = path_prefix[:-1]

        uri = s3_uri_from(
            path=path_prefix,
            bucket=bucket,
            endpoint=endpoint_url,
            region=region,
        )

        if is_file:
            filename = os.path.basename(path)
            if not filename:
                msg = "An error occured determining if the supplied path was file."
                raise StoreError(msg)

            s3_key = os.path.join(path_prefix, filename)
            s3.upload_file(path, bucket, s3_key)

            return uri

        for root, _, files in os.walk(path):
            for filename in files:
                file_path = os.path.join(root, filename)
                relative_path = os.path.relpath(file_path, path)
                s3_key = os.path.join(path_prefix, relative_path)
                s3.upload_file(file_path, bucket, s3_key)

        return uri

    def __s3_creds(
        self,
        endpoint_url: str | None = None,
        access_key_id: str | None = None,
        secret_access_key: str | None = None,
        region: str | None = None,
    ):
        """Internal method to return mix and matched S3 credentials based on presence.

        Args:
            endpoint_url: The S3 compatible object storage endpoint.
            access_key_id: The S3 compatible object storage access key ID.
            secret_access_key: The S3 compatible object storage secret access key.
            region: The region name for the S3 object storage.

        Raises:
            ValueError if the required values are None.

        Returns:
            tuple(endpoint, access_key_id, secret_access_key, region)
        """
        aws_s3_endpoint = os.getenv("AWS_S3_ENDPOINT")
        aws_access_key_id = os.getenv("AWS_ACCESS_KEY_ID")
        aws_secret_access_key = os.getenv("AWS_SECRET_ACCESS_KEY")
        aws_default_region = os.getenv("AWS_DEFAULT_REGION")

        # Set values to parameter values or environment values
        endpoint_url = endpoint_url if endpoint_url else aws_s3_endpoint
        access_key_id = access_key_id if access_key_id else aws_access_key_id
        secret_access_key = (
            secret_access_key if secret_access_key else aws_secret_access_key
        )
        region = region if region else aws_default_region

        if not any((aws_s3_endpoint, endpoint_url)):
            msg = """Please set either `AWS_S3_ENDPOINT` as environment variable
            or specify `endpoint_url` as the parameter.
            """
            raise ValueError(msg)

        if not any((aws_access_key_id, aws_secret_access_key)) and not any(
            (access_key_id, secret_access_key)
        ):
            msg = """Envrionment variables `AWS_ACCESS_KEY_ID` or `AWS_SECRET_ACCESS_KEY` were not set.
            Please either set these environment variables or pass them in as parameters using
            `aws_access_key_id` or `aws_access_secret_id`.
            """
            raise ValueError(msg)

        return endpoint_url, access_key_id, secret_access_key, region

    def __connect_to_s3(
        self,
        endpoint_url: str | None = None,
        access_key_id: str | None = None,
        secret_access_key: str | None = None,
        region: str | None = None,
    ) -> None:
        """Internal method to connect to Boto3 Client.

        Args:
            endpoint_url: The S3 compatible object storage endpoint.
            access_key_id: The S3 compatible object storage access key ID.
            secret_access_key: The S3 compatible object storage secret access key.
            region: The region name for the S3 object storage.

        Raises:
            StoreError: If Boto3 is not installed.
            ValueError: If the appropriate values are not supplied.
        """
        try:
            from boto3 import client  # type: ignore
        except ImportError as e:
            msg = """package `boto3` is not installed.
            To save models to an S3 compatible storage, start by installing the `boto3` package, either directly or as an
            extra (available as `model-registry[boto3]`), e.g.:
            ```sh
            !pip install --pre model-registry[boto3]
            ```
            or
            ```sh
            !pip install boto3
            ```
            """
            raise StoreError(msg) from e

        return client(
            "s3",
            endpoint_url=endpoint_url,
            aws_access_key_id=access_key_id,
            aws_secret_access_key=secret_access_key,
            region_name=region,
        )
