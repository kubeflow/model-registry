"""
Tests for utility functions
"""

from unittest.mock import patch
import os

from mlflow.entities import LifecycleStage, LoggedModelStatus

from modelregistry_plugin.utils import (
    parse_tracking_uri,
    convert_timestamp,
    convert_modelregistry_state,
    convert_to_model_artifact_state,
    convert_to_mlflow_logged_model_status,
    toModelRegistryCustomProperties,
    fromModelRegistryCustomProperties,
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

    def test_parse_tracking_uri_unknown_scheme(self):
        """Test parsing URI with unknown scheme."""
        host, port, secure = parse_tracking_uri("http://example.com:8080")

        assert host == "example.com"
        assert port == 8080
        assert secure is False

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

    def test_convert_to_model_artifact_state(self):
        """Test converting MLflow LoggedModelStatus to Model Artifact State."""
        # Test UNSPECIFIED status
        result = convert_to_model_artifact_state(LoggedModelStatus.UNSPECIFIED)
        assert result == "UNKNOWN"

        # Test PENDING status
        result = convert_to_model_artifact_state(LoggedModelStatus.PENDING)
        assert result == "PENDING"

        # Test READY status
        result = convert_to_model_artifact_state(LoggedModelStatus.READY)
        assert result == "LIVE"

        # Test FAILED status
        result = convert_to_model_artifact_state(LoggedModelStatus.FAILED)
        assert result == "ABANDONED"

        # Test None status
        result = convert_to_model_artifact_state(None)
        assert result == "UNKNOWN"

    def test_convert_to_mlflow_logged_model_status(self):
        """Test converting Model Artifact State to MLflow LoggedModelStatus."""
        # Test PENDING state
        result = convert_to_mlflow_logged_model_status("PENDING")
        assert result == LoggedModelStatus.PENDING

        # Test READY state
        result = convert_to_mlflow_logged_model_status("READY")
        assert result == LoggedModelStatus.READY

        # Test ABANDONED state
        result = convert_to_mlflow_logged_model_status("ABANDONED")
        assert result == LoggedModelStatus.FAILED

        # Test unknown state
        result = convert_to_mlflow_logged_model_status("UNKNOWN_STATE")
        assert result == LoggedModelStatus.UNSPECIFIED

        # Test None state
        result = convert_to_mlflow_logged_model_status(None)
        assert result == LoggedModelStatus.UNSPECIFIED

    def test_toModelRegistryCustomProperties(self):
        """Test converting custom properties to model registry format."""
        # Test with custom properties
        json_data = {
            "customProperties": {
                "string_prop": "test_value",
                "int_prop": 42,
                "float_prop": 3.14,
                "bool_prop": True,
                "none_prop": None,
            }
        }

        toModelRegistryCustomProperties(json_data)

        # Verify string property is wrapped correctly
        assert json_data["customProperties"]["string_prop"]["string_value"] == "test_value"
        assert json_data["customProperties"]["string_prop"]["metadataType"] == "MetadataStringValue"

        # Verify int property is wrapped correctly
        assert json_data["customProperties"]["int_prop"]["string_value"] == "42"
        assert json_data["customProperties"]["int_prop"]["metadataType"] == "MetadataStringValue"

        # Verify float property is wrapped correctly
        assert json_data["customProperties"]["float_prop"]["string_value"] == "3.14"
        assert json_data["customProperties"]["float_prop"]["metadataType"] == "MetadataStringValue"

        # Verify bool property is wrapped correctly
        assert json_data["customProperties"]["bool_prop"]["string_value"] == "True"
        assert json_data["customProperties"]["bool_prop"]["metadataType"] == "MetadataStringValue"

        # Verify None property is not included
        assert "none_prop" not in json_data["customProperties"]

        # Test with no custom properties
        json_data_no_props = {"other_field": "value"}
        toModelRegistryCustomProperties(json_data_no_props)
        assert "customProperties" not in json_data_no_props

    def test_fromModelRegistryCustomProperties(self):
        """Test converting custom properties from model registry format."""
        # Test with different metadata types
        response_json = {
            "customProperties": {
                "string_prop": {
                    "string_value": "test_value",
                    "metadataType": "MetadataStringValue"
                },
                "int_prop": {
                    "int_value": 42,
                    "metadataType": "MetadataIntValue"
                },
                "float_prop": {
                    "float_value": 3.14,
                    "metadataType": "MetadataFloatValue"
                },
                "bool_prop": {
                    "bool_value": True,
                    "metadataType": "MetadataBoolValue"
                },
                "unknown_prop": {
                    "string_value": "unknown_value",
                    "metadataType": "MetadataUnknownType"
                }
            }
        }

        fromModelRegistryCustomProperties(response_json)

        # Verify string property is converted correctly
        assert response_json["customProperties"]["string_prop"] == "test_value"

        # Verify int property is converted correctly
        assert response_json["customProperties"]["int_prop"] == 42

        # Verify float property is converted correctly
        assert response_json["customProperties"]["float_prop"] == 3.14

        # Verify bool property is converted correctly
        assert response_json["customProperties"]["bool_prop"] is True

        # Verify unknown type defaults to string value
        assert response_json["customProperties"]["unknown_prop"] == "unknown_value"

        # Test with no custom properties
        response_json_no_props = {"other_field": "value"}
        fromModelRegistryCustomProperties(response_json_no_props)
        assert "customProperties" not in response_json_no_props
