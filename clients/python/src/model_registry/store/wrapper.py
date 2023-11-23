"""MLMD storage backend wrapper."""

from __future__ import annotations

from collections.abc import Sequence
from typing import ClassVar

from ml_metadata import errors
from ml_metadata.metadata_store import ListOptions, MetadataStore
from ml_metadata.proto import (
    Artifact,
    Attribution,
    Context,
    MetadataStoreClientConfig,
    ParentContext,
)

from model_registry.exceptions import (
    DuplicateException,
    ServerException,
    StoreException,
    TypeNotFoundException,
    UnsupportedTypeException,
)

from .base import ProtoType


class MLMDStore:
    """MLMD storage backend."""

    # cache for MLMD type IDs
    _type_ids: ClassVar[dict[str, int]] = {}

    def __init__(self, config: MetadataStoreClientConfig):
        """Constructor.

        Args:
            config: MLMD config.
        """
        self._mlmd_store = MetadataStore(config)

    def get_type_id(self, mlmd_pt: type[ProtoType], type_name: str) -> int:
        """Get backend ID for a type.

        Args:
            mlmd_pt: Proto type.
            type_name: Name of the type.

        Returns:
            Backend ID.

        Raises:
            TypeNotFoundException: If the type doesn't exist.
            ServerException: If there was an error getting the type.
            UnsupportedTypeException: If the type is not supported.
        """
        if type_name in self._type_ids:
            return self._type_ids[type_name]

        if mlmd_pt is Artifact:
            mlmd_pt_name = "artifact"
            get_type = self._mlmd_store.get_artifact_type
        elif mlmd_pt is Context:
            mlmd_pt_name = "context"
            get_type = self._mlmd_store.get_context_type
        else:
            msg = f"Unsupported type: {mlmd_pt}"
            raise UnsupportedTypeException(msg)

        try:
            _type = get_type(type_name)
        except errors.NotFoundError as e:
            msg = f"{mlmd_pt_name} type {type_name} does not exist"
            raise TypeNotFoundException(msg) from e
        except errors.InternalError as e:
            msg = f"Couldn't get {mlmd_pt_name} type {type_name} from MLMD store"
            raise ServerException(msg) from e

        self._type_ids[type_name] = _type.id

        return _type.id

    def put_artifact(self, artifact: Artifact) -> int:
        """Put an artifact in the store.

        Args:
            artifact: Artifact to put.

        Returns:
            ID of the artifact.

        Raises:
            DuplicateException: If an artifact with the same name or external id already exists.
            TypeNotFoundException: If the type doesn't exist.
            StoreException: If the artifact isn't properly formed.
        """
        try:
            return self._mlmd_store.put_artifacts([artifact])[0]
        except errors.AlreadyExistsError as e:
            msg = f"Artifact {artifact.name} already exists"
            raise DuplicateException(msg) from e
        except errors.InvalidArgumentError as e:
            msg = "Artifact has invalid properties"
            raise StoreException(msg) from e
        except errors.NotFoundError as e:
            msg = f"Artifact type {artifact.type} does not exist"
            raise TypeNotFoundException(msg) from e

    def put_context(self, context: Context) -> int:
        """Put a context in the store.

        Args:
            context: Context to put.

        Returns:
            ID of the context.

        Raises:
            DuplicateException: If a context with the same name or external id already exists.
            TypeNotFoundException: If the type doesn't exist.
            StoreException: If the context isn't propertly formed.
        """
        try:
            return self._mlmd_store.put_contexts([context])[0]
        except errors.AlreadyExistsError as e:
            msg = f"Context {context.name} already exists"
            raise DuplicateException(msg) from e
        except errors.InvalidArgumentError as e:
            msg = "Context has invalid properties"
            raise StoreException(msg) from e
        except errors.NotFoundError as e:
            msg = f"Context type {context.type} does not exist"
            raise TypeNotFoundException(msg) from e

    def _filter_type(
        self, type_name: str, protos: Sequence[ProtoType]
    ) -> Sequence[ProtoType]:
        return [proto for proto in protos if proto.type == type_name]

    def get_context(
        self,
        ctx_type_name: str,
        id: int | None = None,
        name: str | None = None,
        external_id: str | None = None,
    ) -> Context:
        """Get a context from the store.

        This gets a context either by ID, name or external ID.
        If multiple arguments are provided, the simplest query will be performed.

        Args:
            ctx_type_name: Name of the context type.
            id: ID of the context.
            name: Name of the context.
            external_id: External ID of the context.

        Returns:
            Context.

        Raises:
            StoreException: If the context doesn't exist.
        """
        if name is not None:
            return self._mlmd_store.get_context_by_type_and_name(ctx_type_name, name)

        if id is not None:
            contexts = self._mlmd_store.get_contexts_by_id([id])
        elif external_id is not None:
            contexts = self._mlmd_store.get_contexts_by_external_ids([external_id])
        else:
            msg = "Either id, name or external_id must be provided"
            raise StoreException(msg)

        contexts = self._filter_type(ctx_type_name, contexts)
        if contexts:
            return contexts[0]

        msg = f"Context with ID {id} does not exist"
        raise StoreException(msg)

    def get_contexts(
        self, ctx_type_name: str, options: ListOptions
    ) -> Sequence[Context]:
        """Get contexts from the store.

        Args:
            ctx_type_name: Name of the context type.
            options: List options.

        Returns:
            Contexts.
        """
        # TODO: should we make options optional?
        # if options is not None:
        try:
            contexts = self._mlmd_store.get_contexts(options)
        except errors.InvalidArgumentError as e:
            msg = f"Invalid arguments for get_contexts: {e}"
            raise StoreException(msg) from e
        except errors.InternalError as e:
            msg = "Couldn't get contexts from MLMD store"
            raise ServerException(msg) from e

        contexts = self._filter_type(ctx_type_name, contexts)
        # else:
        #     contexts = self._mlmd_store.get_contexts_by_type(ctx_type_name)

        if not contexts:
            msg = f"Context type {ctx_type_name} does not exist"
            raise StoreException(msg)

        return contexts

    def put_context_parent(self, parent_id: int, child_id: int):
        """Put a parent-child relationship between two contexts.

        Args:
            parent_id: ID of the parent context.
            child_id: ID of the child context.

        Raises:
            StoreException: If the parent context doesn't exist.
            ServerException: If there was an error putting the parent context.
        """
        try:
            self._mlmd_store.put_parent_contexts(
                [ParentContext(parent_id=parent_id, child_id=child_id)]
            )
        except errors.AlreadyExistsError as e:
            msg = f"Parent context {parent_id} already exists for context {child_id}"
            raise StoreException(msg) from e
        except errors.InternalError as e:
            msg = f"Couldn't put parent context {parent_id} for context {child_id}"
            raise ServerException(msg) from e

    def put_attribution(self, context_id: int, artifact_id: int):
        """Put an attribution relationship between a context and an artifact.

        Args:
            context_id: ID of the context.
            artifact_id: ID of the artifact.

        Raises:
            StoreException: Invalid argument.
        """
        attribution = Attribution(context_id=context_id, artifact_id=artifact_id)
        try:
            self._mlmd_store.put_attributions_and_associations([attribution], [])
        except errors.InvalidArgumentError as e:
            if "artifact" in str(e).lower():
                msg = f"Artifact with ID {artifact_id} does not exist"
                raise StoreException(msg) from e

            if "context" in str(e).lower():
                msg = f"Context with ID {context_id} does not exist"
                raise StoreException(msg) from e

            msg = f"Invalid argument: {e}"
            raise StoreException(msg) from e

    def get_artifact(
        self,
        art_type_name: str,
        id: int | None = None,
        name: str | None = None,
        external_id: str | None = None,
    ) -> Artifact:
        """Get an artifact from the store.

        Gets an artifact either by ID, name or external ID.

        Args:
            art_type_name: Name of the artifact type.
            id: ID of the artifact.
            name: Name of the artifact.
            external_id: External ID of the artifact.

        Returns:
            Artifact.

        Raises:
            StoreException: If the context doesn't exist.
        """
        if name is not None:
            return self._mlmd_store.get_artifact_by_type_and_name(art_type_name, name)

        if id is not None:
            artifacts = self._mlmd_store.get_artifacts_by_id([id])
        elif external_id is not None:
            artifacts = self._mlmd_store.get_artifacts_by_external_ids([external_id])
        else:
            msg = "Either id, name or external_id must be provided"
            raise StoreException(msg)

        artifacts = self._filter_type(art_type_name, artifacts)
        if artifacts:
            return artifacts[0]

        msg = f"Artifact with ID {id} does not exist"
        raise StoreException(msg)

    def get_attributed_artifact(self, art_type_name: str, ctx_id: int) -> Artifact:
        """Get an artifact from the store by its attributed context.

        Args:
            art_type_name: Name of the artifact type.
            ctx_id: ID of the context.

        Returns:
            Artifact.
        """
        try:
            artifacts = self._mlmd_store.get_artifacts_by_context(ctx_id)
        except errors.InternalError as e:
            msg = f"Couldn't get artifacts by context {ctx_id}"
            raise ServerException(msg) from e
        artifacts = self._filter_type(art_type_name, artifacts)
        if artifacts:
            return artifacts[0]
        msg = "No artifacts found"
        raise StoreException(msg)

    def get_artifacts(
        self, art_type_name: str, options: ListOptions
    ) -> Sequence[Artifact]:
        """Get artifacts from the store.

        Args:
            art_type_name: Name of the artifact type.
            options: List options.

        Returns:
            Artifacts.
        """
        try:
            artifacts = self._mlmd_store.get_artifacts(options)
        except errors.InvalidArgumentError as e:
            msg = f"Invalid arguments for get_artifacts: {e}"
            raise StoreException(msg) from e
        except errors.InternalError as e:
            msg = "Couldn't get artifacts from MLMD store"
            raise ServerException(msg) from e

        artifacts = self._filter_type(art_type_name, artifacts)
        if not artifacts:
            msg = f"Artifact type {art_type_name} does not exist"
            raise StoreException(msg)

        return artifacts
