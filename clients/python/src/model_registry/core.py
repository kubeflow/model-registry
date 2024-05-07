"""Client for the model registry."""

from __future__ import annotations

from dataclasses import dataclass

import grpc

from .exceptions import StoreException
from .store import MLMDStore, ProtoType
from .types import ListOptions, ModelArtifact, ModelVersion, RegisteredModel
from .types.base import ProtoBase
from .types.options import MLMDListOptions
from .utils import header_adder_interceptor


@dataclass
class ModelRegistryAPIClient:
    """Model registry API."""

    store: MLMDStore

    @classmethod
    def secure_connection(
        cls,
        server_address: str,
        port: int = 443,
        user_token: bytes | None = None,
        custom_ca: bytes | None = None,
    ) -> ModelRegistryAPIClient:
        """Constructor.

        Args:
            server_address: Server address.
            port: Server port. Defaults to 443.
            user_token: The PEM-encoded user token as a byte string.
            custom_ca: The PEM-encoded root certificates as a byte string. Defaults to GRPC_DEFAULT_SSL_ROOTS_FILE_PATH, then system default.
        """
        if not user_token:
            msg = "user token must be provided for secure connection"
            raise StoreException(msg)

        chan = grpc.secure_channel(
            f"{server_address}:{port}",
            grpc.composite_channel_credentials(
                # custom_ca = None will get the default root certificates
                grpc.ssl_channel_credentials(custom_ca),
                grpc.access_token_call_credentials(user_token),
            ),
        )

        return cls(MLMDStore.from_channel(chan))

    @classmethod
    def insecure_connection(
        cls,
        server_address: str,
        port: int,
        user_token: bytes | None = None,
    ) -> ModelRegistryAPIClient:
        """Constructor.

        Args:
            server_address: Server address.
            port: Server port.
            user_token: The PEM-encoded user token as a byte string.
        """
        if user_token:
            chan = grpc.intercept_channel(
                grpc.insecure_channel(f"{server_address}:{port}"),
                # header key has to be lowercase
                header_adder_interceptor("authorization", f"Bearer {user_token}"),
            )
        else:
            chan = grpc.insecure_channel(f"{server_address}:{port}")

        return cls(MLMDStore.from_channel(chan))

    def _map(self, py_obj: ProtoBase) -> ProtoType:
        """Map a Python object to a proto object.

        Helper around the `map` method of the Python object.

        Args:
            py_obj: Python object.

        Returns:
            Proto object.
        """
        type_id = self.store.get_type_id(
            py_obj.get_proto_type(), py_obj.get_proto_type_name()
        )
        return py_obj.map(type_id)

    def upsert_registered_model(self, registered_model: RegisteredModel) -> str:
        """Upsert a registered model.

        Updates or creates a registered model on the server.
        This updates the registered_model instance passed in with new data from the servers.

        Args:
            registered_model: Registered model.

        Returns:
            ID of the registered model.
        """
        id = self.store.put_context(self._map(registered_model))
        new_py_rm = RegisteredModel.unmap(
            self.store.get_context(RegisteredModel.get_proto_type_name(), id)
        )
        id = str(id)
        registered_model.id = id
        registered_model.create_time_since_epoch = new_py_rm.create_time_since_epoch
        registered_model.last_update_time_since_epoch = (
            new_py_rm.last_update_time_since_epoch
        )
        return id

    def get_registered_model_by_id(self, id: str) -> RegisteredModel | None:
        """Fetch a registered model by its ID.

        Args:
            id: Registered model ID.

        Returns:
            Registered model.
        """
        proto_rm = self.store.get_context(
            RegisteredModel.get_proto_type_name(), id=int(id)
        )
        if proto_rm is not None:
            return RegisteredModel.unmap(proto_rm)

        return None

    def get_registered_model_by_params(
        self, name: str | None = None, external_id: str | None = None
    ) -> RegisteredModel | None:
        """Fetch a registered model by its name or external ID.

        Args:
            name: Registered model name.
            external_id: Registered model external ID.

        Returns:
            Registered model.

        Raises:
            StoreException: If neither name nor external ID is provided.
        """
        if name is None and external_id is None:
            msg = "Either name or external_id must be provided"
            raise StoreException(msg)
        proto_rm = self.store.get_context(
            RegisteredModel.get_proto_type_name(),
            name=name,
            external_id=external_id,
        )
        if proto_rm is not None:
            return RegisteredModel.unmap(proto_rm)

        return None

    def get_registered_models(
        self, options: ListOptions | None = None
    ) -> list[RegisteredModel]:
        """Fetch registered models.

        Args:
            options: Options for listing registered models.

        Returns:
            Registered models.
        """
        mlmd_options = options.as_mlmd_list_options() if options else MLMDListOptions()
        proto_rms = self.store.get_contexts(
            RegisteredModel.get_proto_type_name(), mlmd_options
        )
        return [RegisteredModel.unmap(proto_rm) for proto_rm in proto_rms]

    def upsert_model_version(
        self, model_version: ModelVersion, registered_model_id: str
    ) -> str:
        """Upsert a model version.

        Updates or creates a model version on the server.
        This updates the model_version instance passed in with new data from the servers.

        Args:
            model_version: Model version to upsert.
            registered_model_id: ID of the registered model this version will be associated to.

        Returns:
            ID of the model version.
        """
        # this is not ideal but we need this info for the prefix
        model_version._registered_model_id = registered_model_id
        id = self.store.put_context(self._map(model_version))
        self.store.put_context_parent(int(registered_model_id), id)
        new_py_mv = ModelVersion.unmap(
            self.store.get_context(ModelVersion.get_proto_type_name(), id)
        )
        id = str(id)
        model_version.id = id
        model_version.create_time_since_epoch = new_py_mv.create_time_since_epoch
        model_version.last_update_time_since_epoch = (
            new_py_mv.last_update_time_since_epoch
        )
        return id

    def get_model_version_by_id(self, model_version_id: str) -> ModelVersion | None:
        """Fetch a model version by its ID.

        Args:
            model_version_id: Model version ID.

        Returns:
            Model version.
        """
        proto_mv = self.store.get_context(
            ModelVersion.get_proto_type_name(), id=int(model_version_id)
        )
        if proto_mv is not None:
            return ModelVersion.unmap(proto_mv)

        return None

    def get_model_versions(
        self, registered_model_id: str, options: ListOptions | None = None
    ) -> list[ModelVersion]:
        """Fetch model versions by registered model ID.

        Args:
            registered_model_id: Registered model ID.
            options: Options for listing model versions.

        Returns:
            Model versions.
        """
        mlmd_options = options.as_mlmd_list_options() if options else MLMDListOptions()
        mlmd_options.filter_query = f"parent_contexts_a.id = {registered_model_id}"
        return [
            ModelVersion.unmap(proto_mv)
            for proto_mv in self.store.get_contexts(
                ModelVersion.get_proto_type_name(), mlmd_options
            )
        ]

    def get_model_version_by_params(
        self,
        registered_model_id: str | None = None,
        version: str | None = None,
        external_id: str | None = None,
    ) -> ModelVersion | None:
        """Fetch a model version by associated parameters.

        Either fetches by using external ID or by using registered model ID and version.

        Args:
            registered_model_id: Registered model ID.
            version: Model version.
            external_id: Model version external ID.

        Returns:
            Model version.

        Raises:
            StoreException: If neither external ID nor registered model ID and version is provided.
        """
        if external_id is not None:
            proto_mv = self.store.get_context(
                ModelVersion.get_proto_type_name(), external_id=external_id
            )
        elif registered_model_id is None or version is None:
            msg = (
                "Either registered_model_id and version or external_id must be provided"
            )
            raise StoreException(msg)
        else:
            proto_mv = self.store.get_context(
                ModelVersion.get_proto_type_name(),
                name=f"{registered_model_id}:{version}",
            )
        if proto_mv is not None:
            return ModelVersion.unmap(proto_mv)

        return None

    def upsert_model_artifact(
        self, model_artifact: ModelArtifact, model_version_id: str
    ) -> str:
        """Upsert a model artifact.

        Updates or creates a model artifact on the server.
        This updates the model_artifact instance passed in with new data from the servers.

        Args:
            model_artifact: Model artifact to upsert.
            model_version_id: ID of the model version this artifact will be associated to.

        Returns:
            ID of the model artifact.

        Raises:
            StoreException: If the model version already has a model artifact.
        """
        mv_id = int(model_version_id)
        if self.store.get_attributed_artifact(
            ModelArtifact.get_proto_type_name(), mv_id
        ):
            msg = f"Model version with ID {mv_id} already has a model artifact"
            raise StoreException(msg)

        model_artifact._model_version_id = model_version_id
        id = self.store.put_artifact(self._map(model_artifact))
        self.store.put_attribution(mv_id, id)
        new_py_ma = ModelArtifact.unmap(
            self.store.get_artifact(ModelArtifact.get_proto_type_name(), id)
        )
        id = str(id)
        model_artifact.id = id
        model_artifact.create_time_since_epoch = new_py_ma.create_time_since_epoch
        model_artifact.last_update_time_since_epoch = (
            new_py_ma.last_update_time_since_epoch
        )
        return id

    def get_model_artifact_by_id(self, id: str) -> ModelArtifact | None:
        """Fetch a model artifact by its ID.

        Args:
            id: Model artifact ID.

        Returns:
            Model artifact.
        """
        proto_ma = self.store.get_artifact(ModelArtifact.get_proto_type_name(), int(id))
        if proto_ma is not None:
            return ModelArtifact.unmap(proto_ma)

        return None

    def get_model_artifact_by_params(
        self, model_version_id: str | None = None, external_id: str | None = None
    ) -> ModelArtifact | None:
        """Fetch a model artifact either by external ID or by the ID of its associated model version.

        Args:
            model_version_id: ID of the associated model version.
            external_id: Model artifact external ID.

        Returns:
            Model artifact.

        Raises:
            StoreException: If neither external ID nor model version ID is provided.
        """
        if external_id:
            proto_ma = self.store.get_artifact(
                ModelArtifact.get_proto_type_name(), external_id=external_id
            )
        elif not model_version_id:
            msg = "Either model_version_id or external_id must be provided"
            raise StoreException(msg)
        else:
            proto_ma = self.store.get_attributed_artifact(
                ModelArtifact.get_proto_type_name(), int(model_version_id)
            )
        if proto_ma is not None:
            return ModelArtifact.unmap(proto_ma)

        return None

    def get_model_artifacts(
        self,
        model_version_id: str | None = None,
        options: ListOptions | None = None,
    ) -> list[ModelArtifact]:
        """Fetches model artifacts.

        Args:
            model_version_id: ID of the associated model version.
            options: Options for listing model artifacts.

        Returns:
            Model artifacts.
        """
        mlmd_options = options.as_mlmd_list_options() if options else MLMDListOptions()
        if model_version_id is not None:
            mlmd_options.filter_query = f"contexts_a.id = {model_version_id}"

        proto_mas = self.store.get_artifacts(
            ModelArtifact.get_proto_type_name(), mlmd_options
        )
        return [ModelArtifact.unmap(proto_ma) for proto_ma in proto_mas]
