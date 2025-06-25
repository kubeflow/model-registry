"""
Minimal tests for ModelRegistryTrackingStore

These tests focus only on store-specific functionality:
- Initialization and configuration
- Basic delegation to operation classes
- Store-specific error handling

Business logic is tested in test_operations_*.py files.
Integration is tested in test_e2e_*.py files.
"""

from unittest.mock import patch

import pytest

from modelregistry_plugin.store import ModelRegistryTrackingStore


class TestModelRegistryTrackingStore:
    def test_init_default(self):
        """Test store initialization with defaults."""
        # Clear environment to test true defaults
        with patch.dict("os.environ", {"MLFLOW_TRACKING_URI": ""}, clear=False):
            store = ModelRegistryTrackingStore()
            assert store.host == "localhost"
            assert store.port == 8080
            assert store.secure is False
            assert store.base_url == "http://localhost:8080/api/model_registry/v1alpha3"
            assert store.artifact_uri is not None

    def test_init_with_custom_uri(self):
        """Test store initialization with custom URI."""
        store = ModelRegistryTrackingStore(
            "modelregistry://example.com:9090", "s3://bucket/artifacts"
        )
        assert store.host == "example.com"
        assert store.port == 9090
        assert store.secure is False
        assert store.artifact_uri == "s3://bucket/artifacts"

    def test_init_secure(self):
        """Test store initialization with secure connection."""
        with patch.dict("os.environ", {"MODEL_REGISTRY_SECURE": "true"}):
            store = ModelRegistryTrackingStore("modelregistry://localhost:8080")
            assert store.secure is True
            assert "https://" in store.base_url

    def test_init_from_env(self):
        """Test store initialization from environment variable."""
        with patch.dict(
            "os.environ", {"MLFLOW_TRACKING_URI": "modelregistry://test:9090"}
        ):
            store = ModelRegistryTrackingStore()
            assert store.host == "test"
            assert store.port == 9090

    def test_operation_classes_initialized(self):
        """Test that all operation classes are properly initialized."""
        store = ModelRegistryTrackingStore("modelregistry://localhost:8080")

        # Verify operation classes exist
        assert hasattr(store, "experiments")
        assert hasattr(store, "runs")
        assert hasattr(store, "metrics")
        assert hasattr(store, "models")
        assert hasattr(store, "search")

        # Verify they have the same API client
        assert store.experiments.api_client is store.api_client
        assert store.runs.api_client is store.api_client
        assert store.metrics.api_client is store.api_client
        assert store.models.api_client is store.api_client
        assert store.search.api_client is store.api_client

    def test_basic_delegation_smoke_test(self):
        """Smoke test that store methods delegate to operation classes."""
        store = ModelRegistryTrackingStore("modelregistry://localhost:8080")

        # Test a few key methods exist and are callable
        # We don't test the actual implementation - that's in operation tests
        assert callable(store.create_experiment)
        assert callable(store.get_experiment)
        assert callable(store.create_run)
        assert callable(store.get_run)
        assert callable(store.log_metric)
        assert callable(store.log_param)
        assert callable(store.search_runs)
        assert callable(store.search_experiments)

    def test_api_client_configuration(self):
        """Test that API client is properly configured."""
        store = ModelRegistryTrackingStore("modelregistry://localhost:8080")

        assert store.api_client.base_url == store.base_url
        assert hasattr(store.api_client, "session")
        assert hasattr(store.api_client, "request")


if __name__ == "__main__":
    # Allow running the tests directly
    pytest.main([__file__, "-v"])
