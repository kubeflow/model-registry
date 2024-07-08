"""Standard client for the model registry."""

from __future__ import annotations

import os
from pathlib import Path
from typing import get_args
from warnings import warn

from .core import ModelRegistryAPIClient
from .exceptions import StoreException
from .store import ScalarType
from .types import ModelArtifact, ModelVersion, RegisteredModel
from .integrator import ModelInfoManager


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
        custom_ca: bytes | None = None,
    ):
        """Constructor.

        Args:
            server_address: Server address.
            port: Server port. Defaults to 443.

        Keyword Args:
            author: Name of the author.
            is_secure: Whether to use a secure connection. Defaults to True.
            user_token: The PEM-encoded user token as a byte string. Defaults to content of path on envvar KF_PIPELINES_SA_TOKEN_PATH.
            custom_ca: The PEM-encoded root certificates as a byte string. Defaults to contents of path on envvar CERT.
        """
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
                if ca_path := os.getenv("CERT"):
                    root_ca = Path(ca_path).read_bytes()
                    # client might have a default CA setup
            else:
                root_ca = custom_ca

            self._api = ModelRegistryAPIClient.secure_connection(
                server_address, port, user_token, root_ca
            )
        elif custom_ca:
            msg = "Custom CA provided without secure connection"
            raise StoreException(msg)
        else:
            self._api = ModelRegistryAPIClient.insecure_connection(
                server_address, port, user_token
            )

    def _register_model(self, name: str, **kwargs) -> RegisteredModel:
        if rm := self._api.get_registered_model_by_params(name):
            return rm

        rm = RegisteredModel(name, **kwargs)
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
        storage_key: str | None = None,
        storage_path: str | None = None,
        service_account_name: str | None = None,
        author: str | None = None,
        owner: str | None = None,
        description: str | None = None,
        metadata: dict[str, ScalarType] | None = None,
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
        rm = self._register_model(name, owner=owner or self._author)
        mv = self._register_new_version(
            rm,
            version,
            author or self._author,
            description=description,
            metadata=metadata or {},
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
        params = locals()
        params.pop('self', None)
        model_info = ModelInfoManager.get_model_info("HuggingFace", params)
        
        return self.register_model(
            model_name or model_info["model_name"],
            model_info["source_uri"],
            author=author or model_info["author"],
            version=version,
            model_format_name=model_format_name,
            model_format_version=model_format_version,
            description=description,
            storage_path=model_info["storage_path"],
            metadata=model_info["metadata"],
        )
        
        
    def register_Mlflow_model(
        self,
        tracking_uri: str,
        registered_name: str,
        *,
        version: str,
        model_format_name: str,
        model_format_version: str,
        author: str | None = None,
        model_name: str | None = None,
        description: str | None = None,
        registered_version: int | None = None,
    ) -> RegisteredModel:
        """Register a MlFlow Registered model.

        This imports a model from MlFlowg and registers it in the model registry.
        Note that the model is not downloaded.

        Args:
            tracking_uri: URI for the MlFlow Server
            registered_name: registered model name in MlFlow MR.

        Keyword Args:
            version: Version of the model. Has to be unique.
            model_format_name: Name of the model format.
            model_format_version: Version of the model format.
            author: Author of the model. Defaults to repo owner.
            model_name: Name of the model. Defaults to the model name in MlFlow Registry.
            description: Description of the model.
            registered_version: registered model version in MlFlow MR, if None will pick the latest.

        Returns:
            Registered model.
        """
        params = locals()
        params.pop('self', None)
        model_info = ModelInfoManager.get_model_info("MlFlow", params)
        
        return self.register_model(
            model_name or model_info["name"],
            model_info["uri"],
            author=author or model_info["author"],
            version= version or model_info["version"],
            model_format_name=model_format_name,
            model_format_version=model_format_version,
            description=description,
            storage_path=model_info["uri"],
            metadata=model_info["metadata"],
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
