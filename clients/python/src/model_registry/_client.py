"""Standard client for the model registry."""
from __future__ import annotations

from .core import ModelRegistryAPIClient
from .exceptions import StoreException
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
            server_address (str): Server address.
            port (int): Server port.
            author (str): Name of the author.
            client_key (str, optional): The PEM-encoded private key as a byte string.
            server_cert (str, optional): The PEM-encoded certificate as a byte string.
            custom_ca (str, optional): The PEM-encoded root certificates as a byte string.
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
        self, rm: RegisteredModel, version: str, /, description: str | None
    ) -> ModelVersion:
        assert rm.id is not None, "Registered model must have an ID"
        if self._api.get_model_version_by_params(rm.id, version):
            msg = f"Version {version} already exists"
            raise StoreException(msg)

        mv = ModelVersion(rm.name, version, self._author, description=description)
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
        description: str | None = None,
        storage_key: str | None = None,
        storage_path: str | None = None,
        service_account_name: str | None = None,
    ) -> RegisteredModel:
        """Register a model.

        Version has to be unique for the model.
        Either `storage_key` and `storage_path`, or `service_account_name` must be provided.

        Args:
            name (str): Name of the model.
            uri (str): URI of the model.

        Keyword Args:
            model_format_name (str): Name of the model format.
            model_format_version (str): Version of the model format.
            version (str): Version of the model.
            description (str, optional): Description of the model.
            storage_key (str, optional): Storage key.
            storage_path (str, optional): Storage path.
            service_account_name (str, optional): Service account name.

        Returns:
            RegisteredModel: Registered model.
        """
        rm = self._register_model(name)
        mv = self._register_new_version(rm, version, description=description)
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

    def get_registered_model(self, name: str) -> RegisteredModel:
        """Get a registered model.

        Args:
            name (str): Name of the model.

        Returns:
            RegisteredModel: Registered model.
        """
        return self._api.get_registered_model_by_params(name)

    def get_model_version(self, name: str, version: str) -> ModelVersion:
        """Get a model version.

        Args:
            name (str): Name of the model.
            version (str): Version of the model.

        Returns:
            ModelVersion: Model version.
        """
        rm = self._api.get_registered_model_by_params(name)
        return self._api.get_model_version_by_params(rm.id, version)

    def get_model_artifact(self, name: str, version: str) -> ModelArtifact:
        """Get a model artifact.

        Args:
            name (str): Name of the model.
            version (str): Version of the model.

        Returns:
            ModelArtifact: Model artifact.
        """
        mv = self.get_model_version(name, version)
        return self._api.get_model_artifact_by_params(mv.id)
