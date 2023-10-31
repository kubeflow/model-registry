from collections.abc import Sequence

from ml_metadata import MetadataStore, errors
from ml_metadata.proto import (
    Artifact,
    Attribution,
    Context,
    MetadataStoreClientConfig,
    ParentContext,
)

from .base import ProtoType
from model_registry.exceptions import (
    DuplicateException,
    ServerException,
    StoreException,
    TypeNotFoundException,
    UnsupportedTypeException,
)


class MLMDStore:
    """MLMD storage backend."""

    # cache for MLMD type IDs
    _type_ids: dict[str, int] = {}

    def __init__(self, config: MetadataStoreClientConfig):
        """Constructor.

        Args:
            config (MetadataStoreClientConfig): MLMD config.
        """
        self._mlmd_store = MetadataStore(config)

    def get_type_id(self, mlmd_pt: ProtoType, type_name: str) -> int:
        """Get backend ID for a type.

        Args:
            mlmd_pt (ProtoType): Proto type.
            type_name (str): Name of the type.

        Returns:
            int: Backend ID.

        Raises:
            TypeNotFoundException: If the type doesn't exist.
            ServerException: If there was an error getting the type.
            UnsupportedTypeException: If the type is not supported.
        """
        if type_name in self._type_ids:
            return self._type_ids[type_name]

        if isinstance(mlmd_pt, Artifact):
            mlmd_pt_name = "artifact"
            get_type = self._mlmd_store.get_artifact_type
        elif isinstance(mlmd_pt, Context):
            mlmd_pt_name = "context"
            get_type = self._mlmd_store.get_context_type
        else:
            raise UnsupportedTypeException(f"Unsupported type: {mlmd_pt}")

        try:
            _type = get_type(type_name)
        except errors.NotFoundError as e:
            raise TypeNotFoundException(
                f"{mlmd_pt_name} type {type_name} does not exist"
            ) from e
        except errors.InternalError as e:
            raise ServerException(
                f"Couldn't get {mlmd_pt_name} type {type_name} from MLMD store"
            ) from e

        self._type_ids[type_name] = _type.id

        return _type.id

    def put_artifact(self, artifact: Artifact) -> int:
        """Put an artifact in the store.

        Args:
            artifact (Artifact): Artifact to put.

        Returns:
            int: ID of the artifact.

        Raises:
            DuplicateException: If an artifact with the same name or external id already exists.
            TypeNotFoundException: If the type doesn't exist.
            StoreException: If the artifact isn't properly formed.
        """
        try:
            return self._mlmd_store.put_artifacts([artifact])[0]
        except errors.AlreadyExistsError as e:
            raise DuplicateException(f"Artifact {artifact.name} already exists") from e
        except errors.InvalidArgumentError as e:
            raise StoreException("Artifact has invalid properties") from e
        except errors.NotFoundError as e:
            raise TypeNotFoundException(
                f"Artifact type {artifact.type} does not exist"
            ) from e

    def put_context(self, context: Context) -> int:
        """Put a context in the store.

        Args:
            context (Context): Context to put.

        Returns:
            int: ID of the context.

        Raises:
            DuplicateException: If a context with the same name or external id already exists.
            TypeNotFoundException: If the type doesn't exist.
            StoreException: If the context isn't propertly formed.
        """
        try:
            return self._mlmd_store.put_contexts([context])[0]
        except errors.AlreadyExistsError as e:
            raise DuplicateException(f"Context {context.name} already exists") from e
        except errors.InvalidArgumentError as e:
            raise StoreException("Context has invalid properties") from e
        except errors.NotFoundError as e:
            raise TypeNotFoundException(
                f"Context type {context.type} does not exist"
            ) from e

    def _filter_type(
        self, type_name: str, protos: Sequence[ProtoType]
    ) -> Sequence[ProtoType]:
        return [proto for proto in protos if proto.type == type_name]

    def get_context(self, ctx_type_name: str, id: int) -> Context:
        """Get a context from the store.

        Args:
            ctx_type_name (str): Name of the context type.
            id (int): ID of the context.

        Returns:
            Context: Context.

        Raises:
            StoreException: If the context doesn't exist.
        """
        contexts = self._mlmd_store.get_contexts_by_id([id])

        contexts = self._filter_type(ctx_type_name, contexts)
        if contexts:
            return contexts[0]

        raise StoreException(f"Context with ID {id} does not exist")

    def put_context_parent(self, parent_id: int, child_id: int):
        """Put a parent-child relationship between two contexts.

        Args:
            parent_id (int): ID of the parent context.
            child_id (int): ID of the child context.

        Raises:
            StoreException: If the parent context doesn't exist.
            ServerException: If there was an error putting the parent context.
        """
        try:
            self._mlmd_store.put_parent_contexts(
                [ParentContext(parent_id=parent_id, child_id=child_id)]
            )
        except errors.AlreadyExistsError as e:
            raise StoreException(
                f"Parent context {parent_id} already exists for context {child_id}"
            ) from e
        except errors.InternalError as e:
            raise ServerException(
                f"Couldn't put parent context {parent_id} for context {child_id}"
            ) from e

    def put_attribution(self, context_id: int, artifact_id: int):
        """Put an attribution relationship between a context and an artifact.

        Args:
            context_id (int): ID of the context.
            artifact_id (int): ID of the artifact.

        Raises:
            StoreException: Invalid argument.
        """
        attribution = Attribution(context_id=context_id, artifact_id=artifact_id)
        try:
            self._mlmd_store.put_attributions_and_associations([attribution], [])
        except errors.InvalidArgumentError as e:
            if "artifact" in str(e).lower():
                raise StoreException(
                    f"Artifact with ID {artifact_id} does not exist"
                ) from e
            elif "context" in str(e).lower():
                raise StoreException(
                    f"Context with ID {context_id} does not exist"
                ) from e
            else:
                raise StoreException(f"Invalid argument: {e}") from e

    def get_artifact(self, art_type_name: str, id: int) -> Artifact:
        """Get an artifact from the store.

        Args:
            art_type_name (str): Name of the artifact type.
            id (int): ID of the artifact.

        Returns:
            Artifact: Artifact.

        Raises:
            StoreException: If the context doesn't exist.
        """
        artifacts = self._mlmd_store.get_artifacts_by_id([id])

        artifacts = self._filter_type(art_type_name, artifacts)
        if artifacts:
            return artifacts[0]

        raise StoreException(f"Artifact with ID {id} does not exist")
