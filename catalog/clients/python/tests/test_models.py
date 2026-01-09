"""E2E tests for model functionality.

Tests model filtering, pagination, and basic operations.

To run these tests:
1. Start the catalog service with test data loaded
2. Set CATALOG_URL environment variable (default: http://localhost:8081)
3. Run: pytest --e2e tests/test_models.py
"""

import pytest

from model_catalog import CatalogAPIClient


class TestModels:
    """Test suite for model functionality."""

    def test_get_models_returns_items(self, api_client: CatalogAPIClient):
        """Test that getting all models returns a list of items."""
        response = api_client.get_models()
        assert "items" in response
        assert isinstance(response["items"], list)

    def test_models_have_required_fields(self, api_client: CatalogAPIClient):
        """Test that models have required fields."""
        response = api_client.get_models()
        assert response.get("items"), "No models found"

        for model in response["items"]:
            assert "name" in model, "Model missing name"
            assert "source_id" in model, f"Model {model.get('name')} missing source_id"

    def test_filter_models_by_source(self, api_client: CatalogAPIClient):
        """Test filtering models by source."""
        sources = api_client.get_sources()
        if sources.get("items"):
            source_id = sources["items"][0].get("id") or sources["items"][0].get("name")
            if source_id:
                response = api_client.get_models(source=source_id)
                assert "items" in response

                # All returned models should be from the specified source
                for model in response.get("items", []):
                    assert model.get("source_id") == source_id

    def test_get_models_with_pagination(self, api_client: CatalogAPIClient):
        """Test model pagination."""
        response = api_client.get_models(page_size=5)
        assert "items" in response
        assert len(response["items"]) <= 5

    def test_pagination_next_page(self, api_client: CatalogAPIClient):
        """Test getting next page of models."""
        response = api_client.get_models(page_size=3)

        if response.get("nextPageToken"):
            response2 = api_client.get_models(
                page_size=3,
                next_page_token=response["nextPageToken"],
            )
            assert "items" in response2

            # Different pages should have different models
            names1 = {m["name"] for m in response["items"]}
            names2 = {m["name"] for m in response2["items"]}
            # No overlap between pages
            assert not names1.intersection(names2)

    def test_models_reference_valid_sources(self, api_client: CatalogAPIClient):
        """Test that all models reference valid enabled sources."""
        sources = api_client.get_sources(page_size=100)
        enabled_source_ids = {
            s.get("id") or s.get("name")
            for s in sources.get("items", [])
            if s.get("enabled") is True or s.get("status") == "available"
        }

        models = api_client.get_models()

        for model in models.get("items", []):
            source_id = model.get("source_id")
            assert source_id in enabled_source_ids, f"Model {model.get('name')} has invalid source_id: {source_id}"

    def test_models_from_disabled_source_excluded(self, api_client: CatalogAPIClient):
        """Test that models from disabled sources don't appear in results."""
        sources = api_client.get_sources()
        assert sources.get("items"), "No sources found"

        disabled_source_ids = [
            s.get("id") or s.get("name")
            for s in sources["items"]
            if s.get("enabled") is False or s.get("status") == "disabled"
        ]

        if not disabled_source_ids:
            pytest.skip("No disabled sources found in test data")

        models = api_client.get_models()
        model_source_ids = [m.get("source_id") for m in models.get("items", [])]

        for disabled_id in disabled_source_ids:
            assert disabled_id not in model_source_ids, f"Model from disabled source {disabled_id} found in results"

    def test_model_count_consistency(self, api_client: CatalogAPIClient):
        """Test that model counts per source are consistent."""
        sources = api_client.get_sources()
        models = api_client.get_models()

        model_counts: dict[str, int] = {}
        for model in models.get("items", []):
            source_id = model.get("source_id")
            model_counts[source_id] = model_counts.get(source_id, 0) + 1

        for source_id in model_counts:
            source = next(
                (s for s in sources["items"] if s.get("id") == source_id),
                None,
            )
            if source:
                assert source.get("enabled") is True or source.get("status") == "available", (
                    f"Source {source_id} has {model_counts[source_id]} models but is not enabled"
                )
