import os
import time
from typing import Union

import pytest
from ml_metadata import errors, metadata_store
from ml_metadata.proto import (
    ArtifactType,
    ContextType,
    metadata_store_pb2,
)
from ml_metadata.proto.metadata_store_pb2 import MetadataStoreClientConfig
from model_registry.core import ModelRegistryAPIClient
from model_registry.store.wrapper import MLMDStore
from model_registry.types import ModelArtifact, ModelVersion, RegisteredModel
from testcontainers.core.container import DockerContainer
from testcontainers.core.waiting_utils import wait_for_logs

ProtoTypeType = Union[ArtifactType, ContextType]


# ruff: noqa: PT021 supported
@pytest.fixture(scope="session")
def mlmd_conn(request) -> MetadataStoreClientConfig:
    model_registry_root_dir = model_registry_root(request)
    print(
        "Assuming this is the Model Registry root directory:", model_registry_root_dir
    )
    shared_volume = model_registry_root_dir / "test/config/ml-metadata"
    sqlite_db_file = shared_volume / "metadata.sqlite.db"
    if sqlite_db_file.exists():
        msg = f"The file {sqlite_db_file} already exists; make sure to cancel it before running these tests."
        raise FileExistsError(msg)
    container = DockerContainer("gcr.io/tfx-oss-public/ml_metadata_store_server:1.14.0")
    container.with_exposed_ports(8080)
    container.with_volume_mapping(
        shared_volume,
        "/tmp/shared",  # noqa: S108
        "rw",
    )
    container.with_env(
        "METADATA_STORE_SERVER_CONFIG_FILE",
        "/tmp/shared/conn_config.pb",  # noqa: S108
    )
    container.start()
    wait_for_logs(container, "Server listening on")
    os.system('docker container ls --format "table {{.ID}}\t{{.Names}}\t{{.Ports}}" -a')  # noqa governed test
    print("waited for logs and port")
    cfg = MetadataStoreClientConfig(
        host="localhost", port=int(container.get_exposed_port(8080))
    )
    print(cfg)

    # this callback is needed in order to perform the container.stop()
    # removing this callback might result in mlmd container shutting down before the tests had chance to fully run,
    # and resulting in grpc connection resets.
    def teardown():
        container.stop()
        print("teardown of plain_wrapper completed.")

    request.addfinalizer(teardown)

    time.sleep(
        3
    )  # allowing some time for mlmd grpc to fully stabilize (is "spent" once per pytest session anyway)
    _throwaway_store = metadata_store.MetadataStore(cfg)
    wait_for_grpc(container, _throwaway_store)

    return cfg


def model_registry_root(request):
    return (request.config.rootpath / "../..").resolve()  # resolves to absolute path


@pytest.fixture()
def plain_wrapper(request, mlmd_conn: MetadataStoreClientConfig) -> MLMDStore:
    sqlite_db_file = (
        model_registry_root(request) / "test/config/ml-metadata/metadata.sqlite.db"
    )

    def teardown():
        try:
            os.remove(sqlite_db_file)
            print(f"Removed {sqlite_db_file} successfully.")
        except Exception as e:
            print(f"An error occurred while removing {sqlite_db_file}: {e}")
        print("plain_wrapper_after_each done.")

    request.addfinalizer(teardown)

    return MLMDStore(mlmd_conn)


def set_type_attrs(mlmd_obj: ProtoTypeType, name: str, props: list[str]):
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
            "state",
        ],
    )

    plain_wrapper._mlmd_store.put_context_type(mv_type)

    rm_type = set_type_attrs(
        ContextType(),
        RegisteredModel.get_proto_type_name(),
        [
            "description",
            "state",
        ],
    )

    plain_wrapper._mlmd_store.put_context_type(rm_type)

    return plain_wrapper


@pytest.fixture()
def mr_api(store_wrapper: MLMDStore) -> ModelRegistryAPIClient:
    mr = object.__new__(ModelRegistryAPIClient)
    mr._store = store_wrapper
    return mr


def wait_for_grpc(
    container: DockerContainer,
    store: metadata_store.MetadataStore,
    timeout=6,
    interval=2,
):
    start = time.time()
    while True:
        duration = time.time() - start
        results = None
        try:
            results = store.get_contexts()
        except errors.UnavailableError as e:
            print(e)
            print("Container logs:\n", container.get_logs())
            print("Container not ready. Retrying...")
        if results is not None:
            return duration
        if timeout and duration > timeout:
            raise TimeoutError("wait_for_grpc not ready %.3f seconds" % timeout)
        time.sleep(interval)
