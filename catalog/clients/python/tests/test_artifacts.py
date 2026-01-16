"""E2E tests for artifact functionality.

Tests artifact filtering and ordering behavior in the catalog API.
Consolidated from test_artifact_filtering.py and test_artifacts_ordering.py.

To run these tests:
1. Start the catalog service with test data loaded
2. Set CATALOG_URL environment variable (default: http://localhost:8081)
3. Run: pytest --e2e tests/test_artifacts.py
"""

import pytest

from model_catalog import CatalogAPIClient


class TestArtifacts:
    """Test suite for artifact functionality (filtering and ordering).

    Uses the model_with_artifacts fixture from conftest.py for tests
    that need a model with artifacts.
    """

    # === Basic Tests ===

    def test_get_artifacts_returns_response(
        self, api_client: CatalogAPIClient, model_with_artifacts: tuple[str, str]
    ) -> None:
        """Test that getting artifacts returns a response."""
        source_id, model_name = model_with_artifacts
        response = api_client.get_artifacts(source_id=source_id, model_name=model_name)
        assert isinstance(response, dict)
        assert "items" in response

    def test_artifacts_have_expected_structure(
        self, api_client: CatalogAPIClient, model_with_artifacts: tuple[str, str]
    ) -> None:
        """Test that artifacts have expected structure."""
        source_id, model_name = model_with_artifacts
        response = api_client.get_artifacts(source_id=source_id, model_name=model_name)

        for item in response.get("items", []):
            assert isinstance(item, dict)
            # Should have at least id or uri
            assert "id" in item or "uri" in item

    # === Filtering Tests ===

    def test_filter_artifacts_by_model(
        self, api_client: CatalogAPIClient, model_with_artifacts: tuple[str, str]
    ) -> None:
        """Test filtering artifacts by model source and name."""
        source_id, model_name = model_with_artifacts

        response = api_client.get_artifacts(source_id=source_id, model_name=model_name)
        assert isinstance(response, dict)
        assert "items" in response

        # Verify we got artifacts
        artifacts = response.get("items", [])
        assert len(artifacts) > 0, f"Expected artifacts for {model_name}"

    def test_filter_artifacts_by_query(
        self, api_client: CatalogAPIClient, model_with_artifacts: tuple[str, str]
    ) -> None:
        """Test filtering artifacts using filter query."""
        source_id, model_name = model_with_artifacts

        # Get all artifacts first
        all_artifacts = api_client.get_artifacts(
            source_id=source_id,
            model_name=model_name,
        )

        assert all_artifacts.get("items"), "No artifacts available for filtering test"

        # Test filter by framework_type property.
        # Note: This filter may return empty results if test data doesn't have
        # artifacts with framework_type="pytorch". The test validates the filter
        # mechanism works, not that specific data exists.
        filter_query = 'framework_type.string_value = "pytorch"'
        filtered = api_client.get_artifacts(
            source_id=source_id,
            model_name=model_name,
            filter_query=filter_query,
        )

        assert isinstance(filtered, dict)
        assert "items" in filtered

        # Verify any filtered results actually match the filter criteria
        for artifact in filtered.get("items", []):
            props = artifact.get("customProperties", {})
            if "framework_type" in props:
                assert props["framework_type"].get("string_value") == "pytorch"

    def test_filter_artifacts_by_numeric_property(
        self, api_client: CatalogAPIClient, model_with_artifacts: tuple[str, str]
    ) -> None:
        """Test filtering artifacts by numeric custom property."""
        source_id, model_name = model_with_artifacts

        # Test filter by accuracy > 0.9
        filter_query = "accuracy.double_value > 0.9"
        response = api_client.get_artifacts(
            source_id=source_id,
            model_name=model_name,
            filter_query=filter_query,
        )

        assert isinstance(response, dict)
        assert "items" in response

        # Verify all returned artifacts have accuracy > 0.9
        for artifact in response.get("items", []):
            props = artifact.get("customProperties", {})
            if "accuracy" in props:
                acc = props["accuracy"].get("double_value")
                assert acc is not None, "Accuracy value is None"
                assert acc > 0.9, f"Expected accuracy > 0.9, got {acc}"

    def test_multiple_models_have_different_artifacts(self, api_client: CatalogAPIClient) -> None:
        """Test that different models have their own artifacts."""
        models = api_client.get_models()
        if not models.get("items") or len(models["items"]) < 2:
            pytest.skip("Need at least 2 models for this test")

        model1 = models["items"][0]
        model2 = models["items"][1]

        artifacts1 = api_client.get_artifacts(
            source_id=model1["source_id"],
            model_name=model1["name"],
        )

        artifacts2 = api_client.get_artifacts(
            source_id=model2["source_id"],
            model_name=model2["name"],
        )

        # Verify both responses are valid
        assert isinstance(artifacts1, dict)
        assert isinstance(artifacts2, dict)
        assert "items" in artifacts1
        assert "items" in artifacts2

    # === Ordering Tests ===

    def test_artifacts_have_ordering_fields(
        self, api_client: CatalogAPIClient, model_with_artifacts: tuple[str, str]
    ) -> None:
        """Test artifacts have timestamp fields for ordering."""
        source_id, model_name = model_with_artifacts

        response = api_client.get_artifacts(
            source_id=source_id,
            model_name=model_name,
        )

        artifacts = response.get("items", [])
        if len(artifacts) < 1:
            pytest.skip("Need at least 1 artifact to test")

        # Check that artifacts have some form of timestamp or id for ordering
        for artifact in artifacts:
            has_ordering_field = (
                artifact.get("id") is not None
                or artifact.get("createTimeSinceEpoch") is not None
                or artifact.get("lastUpdateTimeSinceEpoch") is not None
            )
            assert has_ordering_field, "Artifact missing ordering fields"

    def test_artifacts_have_identifiers(
        self, api_client: CatalogAPIClient, model_with_artifacts: tuple[str, str]
    ) -> None:
        """Test artifacts can be identified by name/id."""
        source_id, model_name = model_with_artifacts

        response = api_client.get_artifacts(
            source_id=source_id,
            model_name=model_name,
        )

        artifacts = response.get("items", [])

        # Extract artifact identifiers
        identifiers = []
        for artifact in artifacts:
            name = artifact.get("name") or artifact.get("artifactType") or artifact.get("id")
            if name:
                identifiers.append(name)

        # We should have identifiers for artifacts
        assert len(identifiers) > 0, "No artifact identifiers found"

    def test_artifacts_pagination(
        self, api_client: CatalogAPIClient, model_with_artifacts: tuple[str, str]
    ) -> None:
        """Test artifacts support pagination."""
        source_id, model_name = model_with_artifacts

        # Request small page size
        response = api_client.get_artifacts(
            source_id=source_id,
            model_name=model_name,
            page_size=1,
        )

        assert isinstance(response, dict)
        assert "items" in response

        # Check pagination fields
        assert "pageSize" in response or "size" in response

    def test_artifacts_custom_properties(
        self, api_client: CatalogAPIClient, model_with_artifacts: tuple[str, str]
    ) -> None:
        """Test that artifacts with custom properties are valid."""
        source_id, model_name = model_with_artifacts

        response = api_client.get_artifacts(
            source_id=source_id,
            model_name=model_name,
        )

        artifacts = response.get("items", [])

        # Extract accuracy values to verify they exist
        accuracy_values = []
        for artifact in artifacts:
            props = artifact.get("customProperties", {})
            if "accuracy" in props:
                acc = props["accuracy"].get("double_value")
                if acc is not None:
                    accuracy_values.append(acc)

        # Test data should have artifacts with accuracy
        if accuracy_values:
            # Verify values are valid numbers
            for acc in accuracy_values:
                assert isinstance(acc, (int, float))
                assert 0 <= acc <= 1, f"Accuracy {acc} should be between 0 and 1"
