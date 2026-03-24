"""E2E tests for MCP server named query functionality.

Tests namedQuery parameter for pre-defined filter templates.

To run these tests:
1. Start the catalog service with MCP test data and namedQueries loaded
2. Set CATALOG_URL environment variable (default: http://localhost:8081)
3. Run: pytest --e2e tests/mcp_servers/test_named_queries.py
"""

import pytest

from model_catalog import CatalogAPIClient, CatalogValidationError


class TestMCPServerNamedQueries:
    """Test suite for MCP server named query functionality."""

    @pytest.mark.parametrize(
        "named_query, expected_security_indicators",
        [
            pytest.param(
                "production_ready",
                {"verifiedSource": True},
                id="production_ready",
            ),
            pytest.param(
                "security_focused",
                {"sast": True, "readOnlyTools": True},
                id="security_focused",
            ),
        ],
    )
    def test_named_query_execution(
        self,
        api_client: CatalogAPIClient,
        suppress_ssl_warnings: None,
        named_query: str,
        expected_security_indicators: dict[str, bool],
    ):
        """TC-API-011: Test executing a named query filters servers by security indicators."""
        response = api_client.get_mcp_servers(named_query=named_query)
        items = response.get("items", [])
        assert len(items) == 1, f"Expected 1 server matching '{named_query}', got {len(items)}"
        assert items[0]["name"] == "calculator"

        security_indicators = items[0].get("securityIndicators", {})
        for prop_name, expected_value in expected_security_indicators.items():
            assert security_indicators.get(prop_name) is expected_value, (
                f"Expected {prop_name}={expected_value}, got {security_indicators.get(prop_name)}"
            )

    @pytest.mark.parametrize(
        "filter_query, expected_count, expected_names",
        [
            pytest.param(
                "provider='Math Community'",
                1,
                {"calculator"},
                id="matching_overlap",
            ),
            pytest.param(
                "provider='Weather Community'",
                0,
                set(),
                id="no_overlap",
            ),
        ],
    )
    def test_named_query_combined_with_filter_query(
        self,
        api_client: CatalogAPIClient,
        suppress_ssl_warnings: None,
        filter_query: str,
        expected_count: int,
        expected_names: set[str],
    ):
        """TC-API-013: Test combining namedQuery with filterQuery."""
        response = api_client.get_mcp_servers(
            named_query="production_ready",
            filter_query=filter_query,
        )
        items = response.get("items", [])
        assert len(items) == expected_count, (
            f"Expected {expected_count} server(s) for namedQuery + '{filter_query}', got {len(items)}"
        )
        assert {s["name"] for s in items} == expected_names

    def test_unknown_named_query_returns_error(
        self,
        api_client: CatalogAPIClient,
        suppress_ssl_warnings: None,
    ):
        """TC-ERROR-005: Test that an unknown named query returns 400."""
        with pytest.raises(CatalogValidationError, match="unknown named query") as exc_info:
            api_client.get_mcp_servers(named_query="nonExistentQuery")
        assert exc_info.value.status_code == 400
