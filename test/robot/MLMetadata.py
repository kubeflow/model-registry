from ml_metadata.metadata_store import metadata_store
from ml_metadata.proto import metadata_store_pb2

class MLMetadata(metadata_store.MetadataStore):
    def __init__(self, host: str = 'localhost', port: int = 8081):
        client_connection_config = metadata_store_pb2.MetadataStoreClientConfig()
        client_connection_config.host = host
        client_connection_config.port = port
        print(client_connection_config)
        super().__init__(client_connection_config)
