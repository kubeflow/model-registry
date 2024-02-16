"""Standard client for the model registry."""
from __future__ import annotations

import os
from typing import get_args
from warnings import warn

from .core import ModelRegistryAPIClient
from .exceptions import StoreException
from .store import ScalarType
from .types import ModelArtifact, ModelVersion, RegisteredModel


class ModelRegistry:
    """Model registry client."""

    def __init__(
        self,
        server_address: str,
        port: int,
        author: str,
        client_key: str | None = None,
        server_cert: str | None = None,
        custom_ca: str | None = None,
    ):
        """Constructor.

        Args:
            server_address: Server address.
            port: Server port.
            author: Name of the author.
            client_key: The PEM-encoded private key as a byte string.
            server_cert: The PEM-encoded certificate as a byte string.
            custom_ca: The PEM-encoded root certificates as a byte string.
        """
        # TODO: get args from env
        self._author = author
        self._api = ModelRegistryAPIClient(
            server_address, port, client_key, server_cert, custom_ca
        )

    def _register_model(self, name: str) -> RegisteredModel:
        if rm := self._api.get_registered_model_by_params(name):
            return rm

        rm = RegisteredModel(name)
        self._api.upsert_registered_model(rm)
        return rm

    def _register_new_version(
        self, rm: RegisteredModel, version: str, author: str, /, **kwargs
    ) -> ModelVersion:
        assert rm.id is not None, "Registered model must have an ID"
        if self._api.get_model_version_by_params(rm.id, version):
            msg = f"Version {version} already exists"
            raise StoreException(msg)

        mv = ModelVersion(rm.name, version, author, **kwargs)
        self._api.upsert_model_version(mv, rm.id)
        return mv

    def _register_model_artifact(
        self, mv: ModelVersion, uri: str, /, **kwargs
    ) -> ModelArtifact:
        assert mv.id is not None, "Model version must have an ID"
        ma = ModelArtifact(mv.model_name, uri, **kwargs)
        self._api.upsert_model_artifact(ma, mv.id)
        return ma

    def register_model(
        self,
        name: str,
        uri: str,
        *,
        model_format_name: str,
        model_format_version: str,
        version: str,
        author: str | None = None,
        description: str | None = None,
        storage_key: str | None = None,
        storage_path: str | None = None,
        service_account_name: str | None = None,
        metadata: dict[str, ScalarType] | None = None,
    ) -> RegisteredModel:
        """Register a model.

        Either `storage_key` and `storage_path`, or `service_account_name` must be provided.

        Args:
            name: Name of the model.
            uri: URI of the model.

        Keyword Args:
            version: Version of the model. Has to be unique.
            model_format_name: Name of the model format.
            model_format_version: Version of the model format.
            description: Description of the model.
            author: Author of the model. Defaults to the client author.
            storage_key: Storage key.
            storage_path: Storage path.
            service_account_name: Service account name.
            metadata: Additional version metadata. Defaults to values returned by `default_metadata()`.

        Returns:
            Registered model.
        """
        rm = self._register_model(name)
        mv = self._register_new_version(
            rm,
            version,
            author or self._author,
            description=description,
            metadata=metadata or self.default_metadata(),
        )
        self._register_model_artifact(
            mv,
            uri,
            model_format_name=model_format_name,
            model_format_version=model_format_version,
            storage_key=storage_key,
            storage_path=storage_path,
            service_account_name=service_account_name,
        )

        return rm

    def default_metadata(self) -> dict[str, ScalarType]:
        """Default metadata valorisations.

        When not explicitly supplied by the end users, these valorisations will be used
        by default.

        Returns:
            default metadata valorisations.
        """
        return {
            key: os.environ[key] for key in ["AWS_S3_ENDPOINT", "AWS_S3_BUCKET", "AWS_DEFAULT_REGION"] if key in os.environ
        }

    def register_hf_model(
        self,
        repo: str,
        path: str,
        *,
        version: str,
        model_format_name: str,
        model_format_version: str,
        author: str | None = None,
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
            raise StoreException(msg) from e

        api = HfApi()
        try:
            model_info = api.model_info(repo, revision=git_ref)
        except utils.RepositoryNotFoundError as e:
            msg = f"Repository {repo} does not exist"
            raise StoreException(msg) from e
        except utils.RevisionNotFoundError as e:
            # TODO: as all hf-hub client calls default to using main, should we provide a tip?
            msg = f"Revision {git_ref} does not exist"
            raise StoreException(msg) from e

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
            **self.default_metadata(),
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
                    if isinstance(v, get_args(ScalarType))
                }
            )
        return self.register_model(
            model_name or model_info.id,
            source_uri,
            author=author or model_author,
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
        return self._api.get_registered_model_by_params(name)

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
        if not (rm := self._api.get_registered_model_by_params(name)):
            msg = f"Model {name} does not exist"
            raise StoreException(msg)
        return self._api.get_model_version_by_params(rm.id, version)

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
            raise StoreException(msg)
        return self._api.get_model_artifact_by_params(mv.id)
