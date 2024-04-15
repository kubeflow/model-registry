"""Tests behavior of the MLMD wrapper.

Tests whether the wrapper is properly handling misuses of the MLMD store, as common use cases
are already covered by the Registry client.
"""

import pytest
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
from model_registry.types.options import MLMDListOptions


@pytest.fixture()
def artifact(plain_wrapper: MLMDStore) -> Artifact:
    art_type = ArtifactType()
    art_type.name = "test_artifact"

    art = Artifact()
    art.name = "test_artifact"
    art.type_id = plain_wrapper._mlmd_store.put_artifact_type(art_type)

    return art


@pytest.fixture()
def context(plain_wrapper: MLMDStore) -> Context:
    ctx_type = ContextType()
    ctx_type.name = "test_context"

    ctx = Context()
    ctx.name = "test_context"
    ctx.type_id = plain_wrapper._mlmd_store.put_context_type(ctx_type)

    return ctx


def test_get_undefined_artifact_type_id(plain_wrapper: MLMDStore):
    with pytest.raises(TypeNotFoundException):
        plain_wrapper.get_type_id(Artifact, "undefined")


def test_get_undefined_context_type_id(plain_wrapper: MLMDStore):
    with pytest.raises(TypeNotFoundException):
        plain_wrapper.get_type_id(Context, "undefined")


@pytest.mark.usefixtures("artifact")
def test_get_no_artifacts(plain_wrapper: MLMDStore):
    arts = plain_wrapper.get_artifacts("test_artifact", MLMDListOptions())
    assert arts == []


def test_get_undefined_artifacts(plain_wrapper: MLMDStore):
    with pytest.raises(TypeNotFoundException):
        plain_wrapper.get_artifacts("undefined", MLMDListOptions())


@pytest.mark.usefixtures("context")
def test_get_no_contexts(plain_wrapper: MLMDStore):
    ctxs = plain_wrapper.get_contexts("test_context", MLMDListOptions())
    assert ctxs == []


def test_get_undefined_contexts(plain_wrapper: MLMDStore):
    with pytest.raises(TypeNotFoundException):
        plain_wrapper.get_contexts("undefined", MLMDListOptions())


def test_put_invalid_artifact(plain_wrapper: MLMDStore, artifact: Artifact):
    artifact.properties["null"].int_value = 0

    with pytest.raises(StoreException):
        plain_wrapper.put_artifact(artifact)


def test_put_duplicate_artifact(plain_wrapper: MLMDStore, artifact: Artifact):
    plain_wrapper._mlmd_store.put_artifacts([artifact])
    with pytest.raises(DuplicateException):
        plain_wrapper.put_artifact(artifact)


def test_put_invalid_context(plain_wrapper: MLMDStore, context: Context):
    context.properties["null"].int_value = 0

    with pytest.raises(StoreException):
        plain_wrapper.put_context(context)


def test_put_duplicate_context(plain_wrapper: MLMDStore, context: Context):
    plain_wrapper._mlmd_store.put_contexts([context])

    with pytest.raises(DuplicateException):
        plain_wrapper.put_context(context)


def test_put_attribution_with_invalid_context(
    plain_wrapper: MLMDStore, artifact: Artifact
):
    art_id = plain_wrapper._mlmd_store.put_artifacts([artifact])[0]

    with pytest.raises(StoreException) as store_error:
        plain_wrapper.put_attribution(0, art_id)

    assert "context" in str(store_error.value).lower()


def test_put_attribution_with_invalid_artifact(
    plain_wrapper: MLMDStore, context: Context
):
    ctx_id = plain_wrapper._mlmd_store.put_contexts([context])[0]

    with pytest.raises(StoreException) as store_error:
        plain_wrapper.put_attribution(ctx_id, 0)

    assert "artifact" in str(store_error.value).lower()


def test_get_undefined_artifact_by_id(plain_wrapper: MLMDStore):
    assert plain_wrapper.get_artifact("dup", 0) is None


def test_get_undefined_context_by_id(plain_wrapper: MLMDStore):
    assert plain_wrapper.get_context("dup", 0) is None
