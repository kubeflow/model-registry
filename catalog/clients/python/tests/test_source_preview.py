"""E2E tests for source preview functionality.

Tests the source preview API endpoint.

To run these tests:
1. Start the catalog service with test data loaded
2. Set CATALOG_URL environment variable (default: http://localhost:8081)
3. Run: pytest --e2e tests/test_source_preview.py
"""

import pytest

from model_catalog import CatalogAPIClient, CatalogAPIError, CatalogValidationError


class TestSourcePreview:
    """Test suite for source preview functionality.

    Tests the preview endpoint which allows testing source configurations
    without modifying the running service state.
    """

    def test_preview_valid_source_config(self, api_client: CatalogAPIClient, suppress_ssl_warnings: None):
        """Test previewing a valid source configuration."""
        # Preview expects a flat config with inline catalog data
        config_content = """
type: yaml
"""
        catalog_data = """
source: Preview Test Source
models:
  - name: valid-preview-model
    description: A model for valid config test
"""
        response = api_client.preview_source(
            config_content=config_content,
            catalog_data=catalog_data,
        )
        assert isinstance(response, dict)
        assert "items" in response

    @pytest.mark.huggingface
    def test_preview_huggingface_source(self, api_client: CatalogAPIClient, suppress_ssl_warnings: None):
        """Test previewing a HuggingFace source configuration.

        HuggingFace preview works without credentials for public models.
        Uses includedModels to specify which models to fetch from HF API.
        """
        # HuggingFace preview requires 'type' and 'includedModels'
        # Note: This is a single source config, not wrapped in 'catalogs'
        config_content = """type: huggingface
includedModels:
  - 'openai-community/gpt2'
"""
        response = api_client.preview_source(config_content=config_content)
        assert isinstance(response, dict)

        # Preview should return items (models from HuggingFace)
        assert "items" in response
        items = response.get("items", [])
        assert isinstance(items, list)
        assert len(items) >= 1

        # Verify the model name
        model_names = [item.get("name") for item in items]
        assert "openai-community/gpt2" in model_names

    def test_preview_yaml_source(self, api_client: CatalogAPIClient, suppress_ssl_warnings: None):
        """Test previewing a YAML catalog source with inline data."""
        # Preview expects a flat config, not the full sources.yaml structure
        config_content = """
type: yaml
"""
        catalog_data = """
source: Inline YAML Preview
models:
  - name: preview-model
    description: Model for preview test
"""
        response = api_client.preview_source(
            config_content=config_content,
            catalog_data=catalog_data,
        )
        assert isinstance(response, dict)

        # Preview should return items
        assert "items" in response
        items = response.get("items", [])
        model_names = [m.get("name") for m in items]
        assert "preview-model" in model_names, f"Expected 'preview-model' in {model_names}"

    def test_preview_invalid_config_returns_error(self, api_client: CatalogAPIClient, suppress_ssl_warnings: None):
        """Test that previewing invalid config returns an error."""
        invalid_config = "not: valid: yaml: config: [[[]]"
        with pytest.raises((CatalogValidationError, CatalogAPIError)):
            api_client.preview_source(config_content=invalid_config)

    def test_preview_empty_config(self, api_client: CatalogAPIClient, suppress_ssl_warnings: None):
        """Test previewing with empty/minimal config.

        The server may either:
        1. Accept empty catalogs and return an empty items list
        2. Reject empty catalogs with a validation error

        Both behaviors are acceptable for this edge case.
        """
        config_content = """
catalogs: []
"""
        try:
            response = api_client.preview_source(config_content=config_content)
            assert isinstance(response, dict)
            # Empty catalogs should return empty items
            items = response.get("items", [])
            assert isinstance(items, list)
        except CatalogValidationError:
            # Server may reject empty catalogs as invalid - this is acceptable
            pass

    def test_preview_response_structure(self, api_client: CatalogAPIClient, suppress_ssl_warnings: None):
        """Test that preview response has expected structure."""
        # Preview expects a flat config, not the full sources.yaml structure
        config_content = """
type: yaml
"""
        catalog_data = """
source: Structure Test
models:
  - name: structure-test-model
    description: Testing response structure
"""
        response = api_client.preview_source(
            config_content=config_content,
            catalog_data=catalog_data,
        )

        # Preview should return model-like structure
        assert isinstance(response, dict)
        assert "items" in response, "Response missing 'items' field"
        assert isinstance(response["items"], list)
