"""E2E tests for model search functionality.

Tests the 'q' search parameter for text-based model search.

To run these tests:
1. Start the catalog service with test data loaded
2. Set CATALOG_URL environment variable (default: http://localhost:8081)
3. Run: pytest --e2e tests/test_search.py
"""

import pytest

from model_catalog import CatalogAPIClient
from tests.search_utils import (
    validate_model_contains_search_term,
    validate_search_results_against_test_data,
)


class TestModelSearch:
    """Test suite for model text search using 'q' parameter."""

    @pytest.mark.parametrize(
        "search_term",
        [
            "alpha",
            "test",
            "ordering",
            "pagination",
        ],
    )
    def test_q_parameter_basic_search(
        self,
        search_term: str,
        api_client: CatalogAPIClient,
        test_catalog_data: dict,
        suppress_ssl_warnings: None,
        kind_cluster: bool,
    ):
        """Test basic search functionality with q parameter using test data validation.

        1. Execute API search
        2. Validate results against test data
        3. Ensure each returned model contains the search term
        """
        response = api_client.get_models(q=search_term)

        assert "items" in response
        models = response.get("items", [])

        # Validate API results against test data (only for Kind clusters)
        if kind_cluster:
            is_valid, errors = validate_search_results_against_test_data(
                api_response=response,
                search_term=search_term,
                catalog_data=test_catalog_data,
            )

            assert is_valid, f"API search results do not match test data for '{search_term}': {errors}"

        # Additional validation: ensure returned models actually contain the search term
        for model in models:
            assert validate_model_contains_search_term(model, search_term), (
                f"Model '{model.get('name')}' doesn't contain search term '{search_term}' in any searchable field"
            )

    @pytest.mark.parametrize(
        "search_term,case_variant",
        [
            ("alpha", "ALPHA"),
            ("test", "TEST"),
        ],
    )
    def test_q_parameter_case_insensitive(
        self,
        search_term: str,
        case_variant: str,
        api_client: CatalogAPIClient,
        test_catalog_data: dict,
        suppress_ssl_warnings: None,
    ):
        """Test that search is case insensitive.

        Validates that searching for the same term in different cases
        returns identical results.
        """
        response1 = api_client.get_models(q=search_term)
        models1 = response1.get("items", [])
        response2 = api_client.get_models(q=case_variant)
        models2 = response2.get("items", [])

        assert len(models1) == len(models2), f"Got {len(models1)} models vs {len(models2)}"

        sorted_m1 = sorted(models1, key=lambda x: x["id"])
        sorted_m2 = sorted(models2, key=lambda x: x.get("id"))

        assert sorted_m1 == sorted_m2

    @pytest.mark.parametrize(
        "search_term",
        [
            "nonexistent_search_term_12345_abcdef",
        ],
    )
    def test_q_parameter_no_results(
        self,
        api_client: CatalogAPIClient,
        test_catalog_data: dict,
        search_term: str,
        suppress_ssl_warnings: None,
    ):
        """Test search with term that should return no results."""
        response = api_client.get_models(q=search_term)
        assert "items" in response

        is_valid, errors = validate_search_results_against_test_data(
            api_response=response,
            search_term=search_term,
            catalog_data=test_catalog_data,
        )

        assert is_valid, f"API search results do not match test data for '{search_term}': {errors}"
        models = response.get("items", [])
        assert len(models) == 0

    @pytest.mark.parametrize("search_term", ["", None])
    def test_q_parameter_empty_query(
        self,
        search_term: str | None,
        api_client: CatalogAPIClient,
        test_catalog_data: dict,
        suppress_ssl_warnings: None,
        kind_cluster: bool,
    ):
        """Test behavior with empty or None q parameter.

        Empty/None search should return all models (no filtering).
        """
        response = api_client.get_models(q=search_term)
        assert "items" in response
        models_q = response.get("items", [])
        response = api_client.get_models()
        models_actual = response.get("items", [])

        # Validate against test data only for Kind clusters
        if kind_cluster:
            models_yaml = test_catalog_data.get("models", [])
            assert len(models_q) == len(models_actual) == len(models_yaml)
        else:
            # For non-Kind clusters, just verify empty/None query returns same as no query
            assert len(models_q) == len(models_actual)
