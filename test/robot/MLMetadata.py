from ml_metadata.metadata_store import metadata_store
from ml_metadata.metadata_store.metadata_store import ListOptions
from ml_metadata.proto import metadata_store_pb2
from ml_metadata import proto
from typing import List, Optional

class MLMetadata(metadata_store.MetadataStore):
    def __init__(self, host: str = 'localhost', port: int = 9090):
        client_connection_config = metadata_store_pb2.MetadataStoreClientConfig()
        client_connection_config.host = host
        client_connection_config.port = port
        print(client_connection_config)
        super().__init__(client_connection_config)

    def get_contexts_by_single_id(self, context_id: int) -> List[proto.Context]:
        return self.get_contexts_by_id([context_id])

    def get_artifacts_by_single_id(self, artifact_id: int) -> List[proto.Artifact]:
        return self.get_artifacts_by_id([artifact_id])