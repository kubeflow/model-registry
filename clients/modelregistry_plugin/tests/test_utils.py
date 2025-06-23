"""
Tests for utility functions
"""

from unittest.mock import patch
import os

from mlflow.entities import LifecycleStage

from modelregistry_plugin.utils import (
    parse_tracking_uri,
    convert_timestamp,
    convert_modelregistry_state,
)


class TestUtils:
    def test_parse_tracking_uri_basic(self):
        """Test parsing basic tracking URI."""
        host, port, secure = parse_tracking_uri("modelregistry://localhost:8080")

        assert host == "localhost"
        assert port == 8080
        assert secure is False

    def test_parse_tracking_uri_https(self):
        """Test parsing HTTPS tracking URI."""
        host, port, secure = parse_tracking_uri(
            "modelregistry+https://example.com:9000"
        )

        assert host == "example.com"
        assert port == 9000
        assert secure is True

    def test_parse_tracking_uri_http(self):
        """Test parsing HTTP tracking URI."""
        host, port, secure = parse_tracking_uri("modelregistry+http://example.com:8080")

        assert host == "example.com"
        assert port == 8080
        assert secure is False

    def test_parse_tracking_uri_with_env_vars(self):
        """Test parsing URI with environment variable overrides."""
        with patch.dict(
            os.environ,
            {
                "MODEL_REGISTRY_HOST": "env-host",
                "MODEL_REGISTRY_PORT": "9090",
                "MODEL_REGISTRY_SECURE": "true",
            },
        ):
            host, port, secure = parse_tracking_uri("modelregistry://")

            assert host == "env-host"
            assert port == 9090
            assert secure is True

    def test_convert_timestamp_numeric(self):
        """Test converting numeric timestamp."""
        result = convert_timestamp("1234567890")
        assert result == 1234567890

    def test_convert_timestamp_iso(self):
        """Test converting ISO format timestamp."""
        # This test might need adjustment based on exact ISO format handling
        result = convert_timestamp("2023-01-01T00:00:00Z")
        assert isinstance(result, int)
        assert result > 0

    def test_convert_timestamp_invalid(self):
        """Test converting invalid timestamp."""
        result = convert_timestamp("invalid")
        assert result is None

    def test_convert_timestamp_empty(self):
        """Test converting empty timestamp."""
        result = convert_timestamp("")
        assert result is None

        result = convert_timestamp(None)
        assert result is None

    def test_convert_modelregistry_state(self):
        """Test converting model registry state to MLflow lifecycle stage."""
        result = convert_modelregistry_state({"state": "LIVE"})
        assert result == LifecycleStage.ACTIVE

    def test_convert_modelregistry_state_archived(self):
        """Test converting model registry state to MLflow lifecycle stage."""
        result = convert_modelregistry_state({"state": "ARCHIVED"})
        assert result == LifecycleStage.DELETED

    def test_convert_modelregistry_state_unknown(self):
        """Test converting model registry state to MLflow lifecycle stage."""
        result = convert_modelregistry_state({})
        assert result == LifecycleStage.ACTIVE
