"""Standard client for the model registry."""

from __future__ import annotations

import os
from pathlib import Path
from typing import Any, get_args
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


class ModelRegistry:
    """Model registry client."""

    def __init__(
        self,
        server_address: str,
        port: int = 443,
        *,
        author: str,
        is_secure: bool = True,
        user_token: bytes | None = None,
        custom_ca: str | None = None,
    ):
        """Constructor.

        Args:
            server_address: Server address.
            port: Server port. Defaults to 443.

        Keyword Args:
            author: Name of the author.
            is_secure: Whether to use a secure connection. Defaults to True.
            user_token: The PEM-encoded user token as a byte string. Defaults to content of path on envvar KF_PIPELINES_SA_TOKEN_PATH.
            custom_ca: Path to the PEM-encoded root certificates as a byte string. Defaults to path on envvar CERT.
        """
        import nest_asyncio

        nest_asyncio.apply()

        # TODO: get remaining args from env
        self._author = author

        if not user_token:
            # /var/run/secrets/kubernetes.io/serviceaccount/token
            sa_token = os.environ.get("KF_PIPELINES_SA_TOKEN_PATH")
            if sa_token:
                user_token = Path(sa_token).read_bytes()
            else:
                warn("User access token is missing", stacklevel=2)

        if is_secure:
            root_ca = None
            if not custom_ca:
                if cert := os.getenv("CERT"):
                    root_ca = cert
                    # client might have a default CA setup
            else:
                root_ca = custom_ca

            if not user_token:
                msg = "user token must be provided for secure connection"
                raise StoreError(msg)

            self._api = ModelRegistryAPIClient.secure_connection(
                server_address, port, user_token=user_token, custom_ca=root_ca
            )
        elif custom_ca:
            msg = "Custom CA provided without secure connection, conflicting options"
            raise StoreError(msg)
        else:
            self._api = ModelRegistryAPIClient.insecure_connection(
                server_address, port, user_token
            )

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
        return await self._api.upsert_model_artifact(
            ModelArtifact(name=name, uri=uri, **kwargs), mv.id
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
        metadata: dict[str, SupportedTypes] | None = None,
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
            msg = "huggingface_hub is not installed"
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
