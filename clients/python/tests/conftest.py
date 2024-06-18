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
def mlmd_port(request) -> int:
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
    port = int(container.get_exposed_port(8080))
    print("port:", port)

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
    _throwaway_store = metadata_store.MetadataStore(
        MetadataStoreClientConfig(host="localhost", port=port)
    )
    wait_for_grpc(container, _throwaway_store)

    return port


def model_registry_root(request):
    return (request.config.rootpath / "../..").resolve()  # resolves to absolute path


@pytest.fixture()
def plain_wrapper(request, mlmd_port: int) -> MLMDStore:
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

    to_return = MLMDStore.from_config("localhost", mlmd_port)
    sanity_check_mlmd_connection_to_db(to_return)
    return to_return


def set_type_attrs(mlmd_obj: ProtoTypeType, name: str, props: list[str]):
    mlmd_obj.name = name
    for key in props:
        mlmd_obj.properties[key] = metadata_store_pb2.STRING
    return mlmd_obj


def sanity_check_mlmd_connection_to_db(overview: MLMDStore):
    # sanity check before each test: connect to MLMD directly, and dry-run any of the gRPC (read) operations;
    # on newer Podman might delay in recognising volume mount files for sqlite3 db,
    # hence in case of error "Cannot connect sqlite3 database: unable to open database file" make some retries.
    retry_count = 0
    while retry_count < 3:
        retry_count += 1
        try:
            overview.store.get_artifact_types()
            return
        except Exception as e:
            if (
                str(e)
                == "Cannot connect sqlite3 database: unable to open database file"
            ):
                time.sleep(1)
            else:
                msg = "Failed to sanity check before each test, another type of error detected."
                raise RuntimeError(msg) from e
    msg = "Failed to sanity check before each test."
    raise RuntimeError(msg)


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

    plain_wrapper.store.put_artifact_type(ma_type)

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

    plain_wrapper.store.put_context_type(mv_type)

    rm_type = set_type_attrs(
        ContextType(),
        RegisteredModel.get_proto_type_name(),
        [
            "description",
            "state",
            "owner",
        ],
    )

    plain_wrapper.store.put_context_type(rm_type)

    return plain_wrapper


@pytest.fixture()
def mr_api(store_wrapper: MLMDStore) -> ModelRegistryAPIClient:
    mr = object.__new__(ModelRegistryAPIClient)
    mr.store = store_wrapper
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
            msg = f"wait_for_grpc not ready {timeout:.3f} seconds"
            raise TimeoutError(msg)
        time.sleep(interval)
