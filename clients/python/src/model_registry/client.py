"""Client for the model registry.
"""
from collections.abc import Sequence
from typing import Optional

from ml_metadata.proto import MetadataStoreClientConfig

from .exceptions import StoreException
from .store import ProtoType, MLMDStore
from .types import ModelArtifact, ModelVersion, RegisteredModel, ListOptions
from .types.base import ProtoBase
from .types.options import MLMDListOptions


class ModelRegistry:
    """Model registry client."""

    def __init__(
        self,
        server_address: str,
        port: int,
        client_key: Optional[str] = None,
        server_cert: Optional[str] = None,
        custom_ca: Optional[str] = None,
    ):
        """Constructor.

        Args:
            server_address (str): Server address.
            port (int): Server port.
            client_key (str, optional): The PEM-encoded private key as a byte string.
            server_cert (str, optional): The PEM-encoded certificate as a byte string.
            custom_ca (str, optional): The PEM-encoded root certificates as a byte string.
        """
        config = MetadataStoreClientConfig()
        config.host = server_address
        config.port = port
        if client_key is not None:
            config.ssl_config.client_key = client_key
        if server_cert is not None:
            config.ssl_config.server_cert = server_cert
        if custom_ca is not None:
            config.ssl_config.custom_ca = custom_ca
        self._store = MLMDStore(config)

    def _map(self, py_obj: ProtoBase) -> ProtoType:
        """Map a Python object to a proto object.

        Helper around the `map` method of the Python object.

        Args:
            py_obj (Mappable): Python object.

        Returns:
            ProtoType: Proto object.
        """
        type_id = self._store.get_type_id(
            py_obj.get_proto_type(), py_obj.get_proto_type_name()
        )
        return py_obj.map(type_id)

    def upsert_registered_model(self, registered_model: RegisteredModel) -> str:
        """Upsert a registered model.

        Updates or creates a registered model on the server.
        This updates the registered_model instance passed in with new data from the servers.

        Args:
            registered_model (RegisteredModel): Registered model.

        Returns:
            str: ID of the registered model.
        """
        id = self._store.put_context(self._map(registered_model))
        new_py_rm = RegisteredModel.unmap(
            self._store.get_context(RegisteredModel.get_proto_type_name(), id)
        )
        id = str(id)
        registered_model.id = id
        registered_model.create_time_since_epoch = new_py_rm.create_time_since_epoch
        registered_model.last_update_time_since_epoch = (
            new_py_rm.last_update_time_since_epoch
        )
        return id

    def get_registered_model_by_id(self, id: str) -> RegisteredModel:
        """Fetch a registered model by its ID.

        Args:
            id (str): Registered model ID.

        Returns:
            RegisteredModel: Registered model.
        """
        py_rm = RegisteredModel.unmap(
            self._store.get_context(RegisteredModel.get_proto_type_name(), id=int(id))
        )
        versions = self.get_model_versions(id)
        assert isinstance(versions, list), "Expected a list"
        py_rm.versions = versions
        return py_rm

    def get_registered_model_by_params(
        self, name: Optional[str] = None, external_id: Optional[str] = None
    ) -> RegisteredModel:
        """Fetch a registered model by its name or external ID.

        Args:
            name (str, optional): Registered model name.
            external_id (str, optional): Registered model external ID.

        Returns:
            RegisteredModel: Registered model.
        """
        if name is None and external_id is None:
            raise StoreException("Either name or external_id must be provided")
        py_rm = RegisteredModel.unmap(
            self._store.get_context(
                RegisteredModel.get_proto_type_name(),
                name=name,
                external_id=external_id,
            )
        )
        assert py_rm.id is not None
        versions = self.get_model_versions(py_rm.id)
        assert isinstance(versions, list), "Expected a list"
        py_rm.versions = versions
        return py_rm

    def get_registered_models(
        self, options: Optional[ListOptions] = None
    ) -> Sequence[RegisteredModel]:
        """Fetch registered models.

        Args:
            options (ListOptions, optional): Options for listing registered models.

        Returns:
            Sequence[RegisteredModel]: Registered models.
        """
        mlmd_options = options.as_mlmd_list_options() if options else None
        proto_rms = self._store.get_contexts(
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
            model_version (ModelVersion): Model version to upsert.
            registered_model_id (str): ID of the registered model this version will be associated to.

        Returns:
            str: ID of the model version.
        """
        rm_id = int(registered_model_id)
        # this is not ideal but we need this info for the prefix
        model_version._registered_model_id = rm_id
        id = self._store.put_context(self._map(model_version))
        self._store.put_context_parent(rm_id, id)
        new_py_mv = ModelVersion.unmap(
            self._store.get_context(ModelVersion.get_proto_type_name(), id)
        )
        id = str(id)
        model_version.id = id
        model_version.create_time_since_epoch = new_py_mv.create_time_since_epoch
        model_version.last_update_time_since_epoch = (
            new_py_mv.last_update_time_since_epoch
        )
        return id

    def get_model_version_by_id(self, model_version_id: str) -> ModelVersion:
        """Fetch a model version by its ID.

        Args:
            model_version_id (str): Model version ID.

        Returns:
            ModelVersion: Model version.
        """
        py_mv = ModelVersion.unmap(
            self._store.get_context(
                ModelVersion.get_proto_type_name(), id=int(model_version_id)
            )
        )
        py_mv.model = self.get_model_artifact_by_params(
            model_version_id=model_version_id
        )
        return py_mv

    def get_model_versions(
        self, registered_model_id: str, options: Optional[ListOptions] = None
    ) -> Sequence[ModelVersion]:
        """Fetch model versions by registered model ID.

        Args:
            registered_model_id (str): Registered model ID.
            options (ListOptions, optional): Options for listing model versions.

        Returns:
            Sequence[ModelVersion]: Model versions.
        """
        mlmd_options = options.as_mlmd_list_options() if options else MLMDListOptions()
        mlmd_options.filter_query = f"parent_contexts_a.id = {registered_model_id}"
        proto_mvs = self._store.get_contexts(
            ModelVersion.get_proto_type_name(), mlmd_options
        )
        py_mvs: list[ModelVersion] = []
        for proto_mv in proto_mvs:
            py_mv = ModelVersion.unmap(proto_mv)
            assert py_mv.id is not None, "Model version ID is None"
            py_mv.model = self.get_model_artifact_by_params(model_version_id=py_mv.id)
            py_mvs.append(py_mv)
        return py_mvs

    def get_model_version_by_params(
        self,
        registered_model_id: Optional[str] = None,
        version: Optional[str] = None,
        external_id: Optional[str] = None,
    ) -> ModelVersion:
        """Fetch a model version by associated parameters.

        Either fetches by using external ID or by using registered model ID and version.

        Args:
            registered_model_id (str, optional): Registered model ID.
            version (str, optional): Model version.
            external_id (str, optional): Model version external ID.

        Returns:
            ModelVersion: Model version.
        """
        if external_id is not None:
            proto_mv = self._store.get_context(
                ModelVersion.get_proto_type_name(), external_id=external_id
            )
        elif registered_model_id is None or version is None:
            raise StoreException(
                "Either registered_model_id and version or external_id must be provided"
            )
        else:
            proto_mv = self._store.get_context(
                ModelVersion.get_proto_type_name(),
                name=f"{registered_model_id}:{version}",
            )
        py_mv = ModelVersion.unmap(proto_mv)
        py_mv.model = self.get_model_artifact_by_params(model_version_id=py_mv.id)
        return py_mv

    def upsert_model_artifact(
        self, model_artifact: ModelArtifact, model_version_id: str
    ) -> str:
        """Upsert a model artifact.

        Updates or creates a model artifact on the server.
        This updates the model_artifact instance passed in with new data from the servers.

        Args:
            model_artifact (ModelArtifact): Model artifact to upsert.
            model_version_id (str): ID of the model version this artifact will be associated to.

        Returns:
            str: ID of the model artifact.
        """
        mv_id = int(model_version_id)
        try:
            self._store.get_attributed_artifact(
                ModelArtifact.get_proto_type_name(), mv_id
            )
            raise StoreException(
                f"Model version with ID {mv_id} already has a model artifact"
            )
        except StoreException as e:
            if "found" not in str(e).lower():
                raise
        model_artifact._model_version_id = model_version_id
        id = self._store.put_artifact(self._map(model_artifact))
        self._store.put_attribution(mv_id, id)
        new_py_ma = ModelArtifact.unmap(
            self._store.get_artifact(ModelArtifact.get_proto_type_name(), id)
        )
        id = str(id)
        model_artifact.id = id
        model_artifact.create_time_since_epoch = new_py_ma.create_time_since_epoch
        model_artifact.last_update_time_since_epoch = (
            new_py_ma.last_update_time_since_epoch
        )
        return id

    def get_model_artifact_by_id(self, id: str) -> ModelArtifact:
        """Fetch a model artifact by its ID.

        Args:
            id (str): Model artifact ID.

        Returns:
            ModelArtifact: Model artifact.
        """
        proto_ma = self._store.get_artifact(
            ModelArtifact.get_proto_type_name(), int(id)
        )
        return ModelArtifact.unmap(proto_ma)

    def get_model_artifact_by_params(
        self, model_version_id: Optional[str] = None, external_id: Optional[str] = None
    ) -> ModelArtifact:
        """Fetch a model artifact either by external ID or by the ID of its associated model version.

        Args:
            model_version_id (str, optional): ID of the associated model version.
            external_id (str, optional): Model artifact external ID.

        Returns:
            ModelArtifact: Model artifact.
        """
        if external_id:
            proto_ma = self._store.get_artifact(
                ModelArtifact.get_proto_type_name(), external_id=external_id
            )
        elif not model_version_id:
            raise StoreException(
                "Either model_version_id or external_id must be provided"
            )
        else:
            proto_ma = self._store.get_attributed_artifact(
                ModelArtifact.get_proto_type_name(), int(model_version_id)
            )
        return ModelArtifact.unmap(proto_ma)

    def get_model_artifacts(
        self,
        model_version_id: Optional[str] = None,
        options: Optional[ListOptions] = None,
    ) -> Sequence[ModelArtifact]:
        """Fetches model artifacts.

        Args:
            model_version_id (str, optional): ID of the associated model version.
            options (ListOptions, optional): Options for listing model artifacts.

        Returns:
            Sequence[ModelArtifact]: Model artifacts.
        """
        mlmd_options = options.as_mlmd_list_options() if options else MLMDListOptions()
        if model_version_id is not None:
            mlmd_options.filter_query = f"contexts_a.id = {model_version_id}"

        proto_mas = self._store.get_artifacts(
            ModelArtifact.get_proto_type_name(), mlmd_options
        )
        return [ModelArtifact.unmap(proto_ma) for proto_ma in proto_mas]
