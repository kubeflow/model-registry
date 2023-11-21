"""Tests behavior of the MLMD wrapper.

Tests whether the wrapper is properly handling misuses of the MLMD store, as common use cases
are already covered by the Registry client.
"""

from ml_metadata.proto import (
    Artifact,
    ArtifactType,
    Context,
    ContextType,
)
from model_registry.exceptions import (
    DuplicateException,
    StoreException,
    TypeNotFoundException,
)
from model_registry.store import MLMDStore
from pytest import fixture, raises


@fixture
def artifact(store_wrapper: MLMDStore) -> Artifact:
    art_type = ArtifactType()
    art_type.name = "test_artifact"

    art = Artifact()
    art.name = "test_artifact"
    art.type_id = store_wrapper._mlmd_store.put_artifact_type(art_type)

    return art


@fixture
def context(store_wrapper: MLMDStore) -> Context:
    ctx_type = ContextType()
    ctx_type.name = "test_context"

    ctx = Context()
    ctx.name = "test_context"
    ctx.type_id = store_wrapper._mlmd_store.put_context_type(ctx_type)

    return ctx


def test_get_undefined_artifact_type_id(store_wrapper: MLMDStore):
    with raises(TypeNotFoundException):
        store_wrapper.get_type_id(Artifact, "undefined")


def test_get_undefined_context_type_id(store_wrapper: MLMDStore):
    with raises(TypeNotFoundException):
        store_wrapper.get_type_id(Context, "undefined")


def test_put_invalid_artifact(store_wrapper: MLMDStore, artifact: Artifact):
    artifact.properties["null"].int_value = 0

    with raises(StoreException):
        store_wrapper.put_artifact(artifact)


def test_put_duplicate_artifact(store_wrapper: MLMDStore, artifact: Artifact):
    store_wrapper._mlmd_store.put_artifacts([artifact])
    with raises(DuplicateException):
        store_wrapper.put_artifact(artifact)


def test_put_invalid_context(store_wrapper: MLMDStore, context: Context):
    context.properties["null"].int_value = 0

    with raises(StoreException):
        store_wrapper.put_context(context)


def test_put_duplicate_context(store_wrapper: MLMDStore, context: Context):
    store_wrapper._mlmd_store.put_contexts([context])

    with raises(DuplicateException):
        store_wrapper.put_context(context)


def test_put_attribution_with_invalid_context(
    store_wrapper: MLMDStore, artifact: Artifact
):
    art_id = store_wrapper._mlmd_store.put_artifacts([artifact])[0]

    with raises(StoreException) as store_error:
        store_wrapper.put_attribution(0, art_id)

    assert "context" in str(store_error.value).lower()


def test_put_attribution_with_invalid_artifact(
    store_wrapper: MLMDStore, context: Context
):
    ctx_id = store_wrapper._mlmd_store.put_contexts([context])[0]

    with raises(StoreException) as store_error:
        store_wrapper.put_attribution(ctx_id, 0)

    assert "artifact" in str(store_error.value).lower()


def test_get_undefined_artifact_by_id(store_wrapper: MLMDStore):
    with raises(StoreException):
        store_wrapper.get_artifact("dup", 0)


def test_get_undefined_context_by_id(store_wrapper: MLMDStore):
    with raises(StoreException):
        store_wrapper.get_context("dup", 0)
