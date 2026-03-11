"""E2E tests for MCP server filtering and error handling.

Tests filterQuery and negative cases.

To run these tests:
1. Start the catalog service with MCP test data loaded
2. Set CATALOG_URL environment variable (default: http://localhost:8081)
3. Run: pytest --e2e tests/test_mcp_servers_filtering.py
"""

import pytest

from model_catalog import CatalogAPIClient, CatalogAPIError, CatalogNotFoundError


class TestMCPServerFiltering:
    """Test suite for MCP server filterQuery functionality."""

    def test_filter_by_provider(
        self,
        api_client: CatalogAPIClient,
        suppress_ssl_warnings: None,
        test_mcp_catalog_data: dict,
        kind_cluster: bool,
    ):
        """Test filtering MCP servers by provider field."""
        target_provider = "Math Community"
        response = api_client.get_mcp_servers(filter_query=f"provider='{target_provider}'")
        items = response.get("items", [])
        assert len(items) == 1, f"Expected 1 server from '{target_provider}', got {len(items)}"
        assert items[0]["name"] == "calculator"
        assert items[0]["provider"] == target_provider

        if kind_cluster:
            yaml_matches = [s["name"] for s in test_mcp_catalog_data["mcp_servers"] if s["provider"] == target_provider]
            assert yaml_matches == ["calculator"]

    def test_filter_by_multiple_providers_in(
        self,
        api_client: CatalogAPIClient,
        suppress_ssl_warnings: None,
        test_mcp_catalog_data: dict,
        kind_cluster: bool,
    ):
        """Test filtering MCP servers using IN operator for multiple providers."""
        target_providers = {"Weather Community", "Math Community"}
        filter_expr = "provider IN ('Weather Community','Math Community')"
        response = api_client.get_mcp_servers(filter_query=filter_expr)
        items = response.get("items", [])
        assert len(items) == 2, f"Expected 2 servers for IN filter, got {len(items)}"
        assert all(s["provider"] in target_providers for s in items)

        if kind_cluster:
            yaml_matches = {s["name"] for s in test_mcp_catalog_data["mcp_servers"] if s["provider"] in target_providers}
            actual_names = {s["name"] for s in items}
            assert actual_names == yaml_matches

    def test_filter_by_boolean(
        self,
        api_client: CatalogAPIClient,
        suppress_ssl_warnings: None,
    ):
        """Test filtering MCP servers by boolean security indicators."""
        response = api_client.get_mcp_servers(filter_query="verifiedSource=true AND sast=true")
        items = response.get("items", [])
        assert len(items) == 1, f"Expected 1 server with verifiedSource+sast, got {len(items)}"
        assert items[0]["name"] == "calculator"

    def test_complex_filter_and_or(
        self,
        api_client: CatalogAPIClient,
        suppress_ssl_warnings: None,
    ):
        """Test complex filter with AND/OR boolean logic."""
        filter_expr = "license='MIT' AND (provider='Math Community' OR provider='Weather Community')"
        response = api_client.get_mcp_servers(filter_query=filter_expr)
        items = response.get("items", [])
        assert len(items) == 2, f"Expected 2 servers for complex filter, got {len(items)}"
        actual_names = {s["name"] for s in items}
        assert actual_names == {"weather-api", "calculator"}

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


class TestMCPServerNegative:
    """Test suite for MCP server error handling and negative cases."""

    def test_get_nonexistent_server_404(
        self,
        api_client: CatalogAPIClient,
        suppress_ssl_warnings: None,
    ):
        """Test that requesting a non-existent MCP server returns 404."""
        with pytest.raises(CatalogNotFoundError):
            api_client.get_mcp_server(server_id="99999")

    @pytest.mark.parametrize(
        "invalid_filter_query",
        [
            pytest.param(
                "provider IN ('test','value'))",
                id="malformed_syntax_unbalanced_parentheses",
            ),
            pytest.param(
                "provider = 'test' AND",
                id="malformed_syntax_incomplete_expression",
            ),
        ],
    )
    def test_invalid_filter_syntax(
        self,
        api_client: CatalogAPIClient,
        suppress_ssl_warnings: None,
        invalid_filter_query: str,
    ):
        """Test that invalid filterQuery syntax raises an error."""
        with pytest.raises(CatalogAPIError, match="invalid filter query"):
            api_client.get_mcp_servers(filter_query=invalid_filter_query)

