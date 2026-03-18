"""E2E tests for labels endpoint assetType filtering.

Tests that the /labels endpoint correctly filters labels by assetType
query parameter (models, mcp_servers, default, invalid).

To run these tests:
1. Start the catalog service with test data loaded
2. Set CATALOG_URL environment variable (default: http://localhost:8081)
3. Run: pytest --e2e tests/test_labels.py
"""

import pytest

from model_catalog import CatalogAPIClient, CatalogValidationError


class TestLabelAssetTypeFilter:
    """Test suite for /labels endpoint assetType query parameter filtering."""

    @pytest.mark.parametrize(
        "asset_type",
        [None, "models", "mcp_servers"],
        ids=["default-models", "explicit-models", "mcp-servers"],
    )
    def test_asset_type_filters_labels(
        self,
        asset_type: str | None,
        api_client: CatalogAPIClient,
        suppress_ssl_warnings: None,
    ):
        """Test that assetType filter returns the correct labels."""
        response = api_client.get_labels(asset_type=asset_type)
        assert isinstance(response.get("items"), list)

        actual_names = {label.get("name") for label in response["items"]}
        default_names = {label.get("name") for label in api_client.get_labels()["items"]}

        if asset_type == "mcp_servers":
            assert len(actual_names) > 0
            assert actual_names.isdisjoint(default_names)
        else:
            # None and "models" should return the same result
            assert actual_names == default_names

    def test_invalid_asset_type_returns_400(
        self,
        api_client: CatalogAPIClient,
        suppress_ssl_warnings: None,
    ):
        """Test that an invalid assetType value returns a 400 error."""
        with pytest.raises(CatalogValidationError, match="invalid value"):
            api_client.get_labels(asset_type="invalid_value")
