"""E2E tests for model filterQuery functionality.

Tests model-level filterQuery search with IN, ILIKE, AND/OR operators.
Migrated from opendatahub-tests (RHOAIENG-45971).

To run these tests:
1. Start the catalog service with test data loaded
2. Set CATALOG_URL environment variable (default: http://localhost:8081)
3. Run: pytest --e2e tests/test_filter_query.py
"""

import pytest

from model_catalog import CatalogAPIClient, CatalogAPIError
from tests.search_utils import get_expected_models_for_filter_query


class TestFilterQuery:
    """Test suite for model-level filterQuery search (RHOAIENG-33658)."""

    def test_search_models_by_filter_query(
        self,
        api_client: CatalogAPIClient,
        test_catalog_data: dict,
        suppress_ssl_warnings: None,
        kind_cluster: bool,
    ) -> None:
        """Test that filterQuery returns models matching combined filter criteria.

        Validates IN operator on custom properties, ILIKE for case-insensitive
        pattern matching, equality on JSON array property (tasks),
        and AND logic combining all conditions.
        """
        sizes = ("7B", "13B")
        filter_query = "size IN ('7B','13B') AND name ILIKE '%test%' AND tasks = 'text-generation'"

        response = api_client.get_models(filter_query=filter_query)

        assert "items" in response
        models = response.get("items", [])
        assert len(models) > 0, f"Expected models matching filter query: {filter_query}"

        for model in models:
            model_name = model["name"]
            assert "test" in model_name.lower()
            model_tasks = model["tasks"]
            assert "text-generation" in model_tasks
            size_value = model["customProperties"]["size"]["string_value"]
            assert size_value in sizes

        if kind_cluster:
            expected_names = get_expected_models_for_filter_query(
                catalog_data=test_catalog_data,
                name_pattern="test",
                task="text-generation",
                sizes=sizes,
            )
            actual_names = {model["name"] for model in models}
            assert actual_names == expected_names


class TestFilterQueryNegative:
    """Test suite for negative filterQuery cases (RHOAIENG-36938)."""

    @pytest.mark.parametrize(
        "invalid_filter_query",
        [
            pytest.param(
                "fake IN ('test','value'))",
                id="malformed_syntax_unbalanced_parentheses",
            ),
            pytest.param(
                "name = 'test' AND",
                id="malformed_syntax_incomplete_expression",
            ),
        ],
    )
    def test_search_models_by_invalid_filter_query(
        self,
        api_client: CatalogAPIClient,
        invalid_filter_query: str,
        suppress_ssl_warnings: None,
    ) -> None:
        """Test that invalid filterQuery raises a validation error."""
        with pytest.raises(CatalogAPIError, match="invalid filter query"):
            api_client.get_models(filter_query=invalid_filter_query)

    def test_filter_query_valid_but_no_matches(
        self,
        api_client: CatalogAPIClient,
        suppress_ssl_warnings: None,
    ) -> None:
        """Test that a valid filterQuery with no matching data returns zero results."""
        filter_query = "name = 'nonexistent_model_xyz_12345'"

        response = api_client.get_models(filter_query=filter_query)

        assert response["items"] == []
        assert response["size"] == 0