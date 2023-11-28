from typing import Union

import pytest
from ml_metadata.proto import (
    ArtifactType,
    ConnectionConfig,
    ContextType,
    metadata_store_pb2,
)
from model_registry.store.wrapper import MLMDStore
from model_registry.types import ModelArtifact, ModelVersion, RegisteredModel


@pytest.fixture()
def plain_wrapper() -> MLMDStore:
    config = ConnectionConfig()
    config.fake_database.SetInParent()
    return MLMDStore(config)


def set_type_attrs(
    mlmd_obj: Union[ArtifactType, ContextType], name: str, props: list[str]
):
    mlmd_obj.name = name
    for key in props:
        mlmd_obj.properties[key] = metadata_store_pb2.STRING
    return mlmd_obj


@pytest.fixture()
def store_wrapper(plain_wrapper: MLMDStore) -> MLMDStore:
    ma_type = set_type_attrs(
        ArtifactType(),
        ModelArtifact.get_proto_type_name(),
        [
            "description",
            "model_format_name",
            "model_format_version",
            "storage_key",
            "storage_path",
            "service_account_name",
        ],
    )

    plain_wrapper._mlmd_store.put_artifact_type(ma_type)

    mv_type = set_type_attrs(
        ContextType(),
        ModelVersion.get_proto_type_name(),
        [
            "author",
            "description",
            "model_name",
            "tags",
        ],
    )

    plain_wrapper._mlmd_store.put_context_type(mv_type)

    rm_type = set_type_attrs(
        ContextType(),
        RegisteredModel.get_proto_type_name(),
        [
            "description",
        ],
    )

    plain_wrapper._mlmd_store.put_context_type(rm_type)

    return plain_wrapper
