"""E2E tests for remote MCP server support (RHOAIENG-46621).

Tests remote server fields, transport derivation, and deployment mode filtering.

To run these tests:
1. Start the catalog service with MCP test data loaded
2. Set CATALOG_URL environment variable (default: http://localhost:8081)
3. Run: pytest --e2e tests/test_mcp_remote_servers.py
"""

from typing import Self

import pytest

from model_catalog import CatalogAPIClient


class TestRemoteServerFieldsAndTransports:
    """TC-API-018, TC-API-025: Verify remote server fields and transport derivation."""

    @pytest.mark.parametrize(
        "provider_filter, expected_mode, expected_transports",
        [
            pytest.param("provider='CloudMCP'", "remote", {"http"}, id="remote_http"),
            pytest.param("provider='StreamCorp'", "remote", {"sse"}, id="remote_sse"),
            pytest.param("provider='HybridCorp'", "local", {"stdio", "http"}, id="hybrid"),
        ],
    )
    def test_server_fields_and_transports(
        self: Self,
        api_client: CatalogAPIClient,
        suppress_ssl_warnings: None,
        provider_filter: str,
        expected_mode: str,
        expected_transports: set[str],
    ):
        """Verify deploymentMode, endpoints, and transports for each server type."""
        response = api_client.get_mcp_servers(filter_query=provider_filter)
        items = response.get("items", [])
        assert len(items) == 1, f"Expected 1 server for '{provider_filter}', got {len(items)}"

        server = items[0]
        assert server["deploymentMode"] == expected_mode
        assert set(server["transports"]) == expected_transports
        assert server.get("endpoints"), "Server should have endpoints"


class TestRemoteServerFiltering:
    """TC-API-007, TC-API-008: Verify filtering by deploymentMode and transports."""

    @pytest.mark.parametrize(
        "filter_query, expected_names, excluded_names",
        [
            pytest.param(
                "deploymentMode='remote'",
                {"remote-http-server", "remote-sse-server"},
                {"weather-api", "hybrid-server"},
                id="deployment_mode_remote",
            ),
            pytest.param(
                "deploymentMode='local'",
                {"weather-api", "hybrid-server"},
                {"remote-http-server", "remote-sse-server"},
                id="deployment_mode_local",
            ),
            pytest.param(
                "transports='http'",
                {"remote-http-server", "hybrid-server"},
                {"remote-sse-server"},
                id="transport_http",
            ),
            pytest.param(
                "transports='sse'",
                {"remote-sse-server"},
                {"remote-http-server"},
                id="transport_sse",
            ),
        ],
    )
    def test_filter(
        self: Self,
        api_client: CatalogAPIClient,
        suppress_ssl_warnings: None,
        filter_query: str,
        expected_names: set[str],
        excluded_names: set[str],
    ):
        """Verify filter returns correct servers."""
        response = api_client.get_mcp_servers(filter_query=filter_query)
        items = response.get("items", [])
        names = {server["name"] for server in items}
        assert expected_names <= names, f"Missing expected servers: {expected_names - names}"
        assert not (excluded_names & names), f"Unexpected servers returned: {excluded_names & names}"
