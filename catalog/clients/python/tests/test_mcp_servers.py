"""E2E tests for MCP server catalog functionality.

Tests MCP server listing, retrieval by ID, tools, and custom properties.

To run these tests:
1. Start the catalog service with MCP test data loaded
2. Set CATALOG_URL environment variable (default: http://localhost:8081)
3. Run: pytest --e2e tests/test_mcp_servers.py
"""

import random

from model_catalog import CatalogAPIClient


class TestMCPServerBasics:
    """Test suite for basic MCP server listing and retrieval."""

    def test_mcp_servers_loaded(
        self,
        api_client: CatalogAPIClient,
        suppress_ssl_warnings: None,
        test_mcp_catalog_data: dict,
    ):
        """Test that all expected MCP servers are loaded from YAML."""
        response = api_client.get_mcp_servers()
        actual_names = {s["name"] for s in response.get("items", [])}
        expected_names = {s["name"] for s in test_mcp_catalog_data["mcp_servers"]}
        assert actual_names == expected_names

    def test_mcp_server_providers(
        self,
        api_client: CatalogAPIClient,
        suppress_ssl_warnings: None,
        test_mcp_catalog_data: dict,
    ):
        """Test that MCP server providers match expected values."""
        response = api_client.get_mcp_servers()
        actual_providers = {s["name"]: s["provider"] for s in response["items"]}
        expected_providers = {s["name"]: s["provider"] for s in test_mcp_catalog_data["mcp_servers"]}
        assert actual_providers == expected_providers

    def test_mcp_server_get_by_id(
        self,
        api_client: CatalogAPIClient,
        suppress_ssl_warnings: None,
    ):
        """Test that an MCP server can be retrieved by ID."""
        response = api_client.get_mcp_servers()
        assert response.get("items"), "No MCP servers found"
        server = random.choice(response["items"])
        single = api_client.get_mcp_server(server_id=server["id"])
        assert single["name"] == server["name"]

    def test_mcp_server_pagination_fields(
        self,
        api_client: CatalogAPIClient,
        suppress_ssl_warnings: None,
    ):
        """Test that MCP servers response includes pagination fields."""
        response = api_client.get_mcp_servers()
        assert "size" in response
        assert "pageSize" in response

    def test_filter_by_name(
        self,
        api_client: CatalogAPIClient,
        suppress_ssl_warnings: None,
    ):
        """Test filtering MCP servers by exact name."""
        target_name = "calculator"
        response = api_client.get_mcp_servers(name=target_name)
        items = response.get("items", [])
        assert len(items) == 1, f"Expected 1 server named '{target_name}', got {len(items)}"
        assert items[0]["name"] == target_name

class TestMCPServerTools:
    """Test suite for MCP server tools functionality."""

    def test_tool_count_without_include(
        self,
        api_client: CatalogAPIClient,
        suppress_ssl_warnings: None,
        test_mcp_catalog_data: dict,
    ):
        """Test that toolCount reflects actual tools even without includeTools."""
        response = api_client.get_mcp_servers()
        expected_counts = {
            s["name"]: len(s.get("tools", [])) for s in test_mcp_catalog_data["mcp_servers"]
        }
        for server in response.get("items", []):
            name = server["name"]
            if name in expected_counts:
                assert server.get("toolCount", 0) == expected_counts[name], (
                    f"Server '{name}': expected toolCount {expected_counts[name]}, got {server.get('toolCount', 0)}"
                )

    def test_tools_included_when_requested(
        self,
        api_client: CatalogAPIClient,
        suppress_ssl_warnings: None,
        test_mcp_catalog_data: dict,
    ):
        """Test that tools are returned when include_tools=True."""
        response = api_client.get_mcp_servers(include_tools=True)
        expected_tools = {
            s["name"]: [t["name"] for t in s.get("tools", [])] for s in test_mcp_catalog_data["mcp_servers"]
        }
        for server in response.get("items", []):
            name = server["name"]
            if name in expected_tools:
                actual_tool_names = [t["name"] for t in server.get("tools", [])]
                assert sorted(actual_tool_names) == sorted(expected_tools[name])

class TestMCPServerCustomProperties:
    """Test suite for MCP server custom properties."""

    def test_custom_properties_loaded(
        self,
        api_client: CatalogAPIClient,
        suppress_ssl_warnings: None,
        test_mcp_catalog_data: dict,
    ):
        """Test that customProperties are correctly loaded from YAML."""
        response = api_client.get_mcp_servers()
        servers_by_name = {s["name"]: s for s in response.get("items", [])}

        for yaml_server in test_mcp_catalog_data["mcp_servers"]:
            name = yaml_server["name"]
            expected_props = yaml_server.get("customProperties")
            if expected_props:
                assert name in servers_by_name
                assert servers_by_name[name].get("customProperties") == expected_props
