from model_registry.core import ModelRegistryAPIClient


def test_connection_args():
    """Test connection arguments parsing."""
    # old-style: explicit port
    server_address = "http://localhost"
    port = 8080
    host = f"{server_address}:{port}"

    client = ModelRegistryAPIClient.insecure_connection(server_address, port)
    assert client.config.host == host

    # new-style: port in url
    server_address_with_port = "http://localhost:9090"
    client = ModelRegistryAPIClient.insecure_connection(server_address_with_port, port)
    assert client.config.host == server_address_with_port

    # Secure connection tests
    client = ModelRegistryAPIClient.secure_connection(
        server_address, port, user_token="token"  # noqa: S106
    )
    assert client.config.host == host

    client = ModelRegistryAPIClient.secure_connection(
        server_address_with_port, port, user_token="token"  # noqa: S106
    )
    assert client.config.host == server_address_with_port
