"""Tests for ModelOperations."""

from unittest.mock import Mock

import pytest
from mlflow.entities import (
    LoggedModel,
    LoggedModelParameter,
    LoggedModelStatus,
    LoggedModelTag,
)
from mlflow.models import Model
from mlflow.store.entities.paged_list import PagedList

from modelregistry_plugin.operations.model import ModelOperations


class TestModelOperations:
    @pytest.fixture
    def api_client(self):
        """Create a mock API client."""
        return Mock()

    @pytest.fixture
    def model_ops(self, api_client):
        """Create a ModelOperations instance for testing."""
        return ModelOperations(api_client, "s3://bucket/artifacts")

    def test_init(self, model_ops, api_client):
        """Test ModelOperations initialization."""
        assert model_ops.api_client == api_client
        assert model_ops.artifact_uri == "s3://bucket/artifacts"

    def test_record_logged_model(self, model_ops, api_client):
        """Test recording a logged model."""
        # Mock MLflow Model
        model = Mock(spec=Model)
        model.to_dict.return_value = {
            "artifact_path": "model",
            "run_id": "run-123",
            "flavors": {"python_function": {}},
            "model_uuid": "uuid-123",
            "utc_time_created": "2023-01-01T00:00:00Z",
            "mlflow_version": "2.0.0",
        }
        model_info = Mock()
        model_info.model_uri = "runs:/run-123/model"
        model_info.artifact_path = "model"
        model_info.model_uuid = "uuid-123"
        model_info.utc_time_created = "2023-01-01T00:00:00Z"
        model_info.mlflow_version = "2.0.0"
        model_info.flavors = {"python_function": {}}
        model.get_model_info.return_value = model_info
        model.model_id = None
        model.model_uuid = "uuid-123"
        model.name = "test-model"
        model.source_run_id = "run-123"

        api_client.post.return_value = {
            "id": "model-123",
            "name": "test-model",
            "uri": "s3://bucket/artifacts/experiments/exp-123/run-123/test-model",
            "createTimeSinceEpoch": "1234567890",
            "lastUpdateTimeSinceEpoch": "1234567890",
            "artifactType": "model-artifact",
            "state": "LIVE",
            "customProperties": {
                "mlflow__model_io_type": {
                    "string_value": "output",
                    "metadataType": "MetadataStringValue",
                },
                "mlflow__source_run_id": {
                    "string_value": "run-123",
                    "metadataType": "MetadataStringValue",
                },
                "mlflow__model_uuid": {
                    "string_value": "uuid-123",
                    "metadataType": "MetadataStringValue",
                },
            },
        }

        model_ops.record_logged_model("run-123", model)

        post_call = api_client.post.call_args
        assert post_call[0][0] == "/experiment_runs/run-123/artifacts"
        json_data = post_call[1]["json"]
        assert json_data["artifactType"] == "model-artifact"
        assert json_data["name"] == "uuid-123"
        assert json_data["uri"] == "runs:/run-123/model"
        assert json_data["customProperties"]["mlflow__artifactPath"] == "model"
        assert json_data["customProperties"]["mlflow__model_uuid"] == "uuid-123"

    def test_create_logged_model(self, model_ops, api_client):
        """Test creating a logged model."""
        # Mock experiment data
        experiment_data = {
            "id": "exp-123",
            "name": "test-experiment",
            "externalId": "s3://bucket/artifacts/experiments/exp-123",
            "customProperties": {},
        }
        api_client.get.return_value = experiment_data

        # Mock model artifact response
        model_data = {
            "id": "model-123",
            "name": "test-model",
            "uri": "s3://bucket/artifacts/experiments/exp-123/run-123/test-model",
            "createTimeSinceEpoch": "1234567890",
            "lastUpdateTimeSinceEpoch": "1234567890",
            "artifactType": "model-artifact",
            "state": "LIVE",
            "customProperties": {
                "mlflow__model_type": "sklearn",
                "mlflow__source_run_id": "run-123",
                "mlflow__experiment_id": "exp-123",
                "key1": "value1",
                "param_param1": "value1",
            },
        }
        api_client.post.return_value = model_data

        tags = [LoggedModelTag("key1", "value1")]
        params = [LoggedModelParameter("param1", "value1")]

        logged_model = model_ops.create_logged_model(
            name="test-model",
            source_run_id="run-123",
            experiment_id="exp-123",
            tags=tags,
            params=params,
            model_type="sklearn",
        )

        assert isinstance(logged_model, LoggedModel)
        assert logged_model.model_id == "model-123"
        assert logged_model.experiment_id == "exp-123"
        assert logged_model.name == "test-model"
        assert logged_model.source_run_id == "run-123"
        assert logged_model.model_type == "sklearn"
        assert logged_model.params["param1"] == "value1"
        assert logged_model.tags["key1"] == "value1"
        assert logged_model.status == LoggedModelStatus.READY

        # Check API calls
        assert api_client.get.call_count == 1
        api_client.get.assert_called_once_with("/experiments/exp-123")

        assert api_client.post.call_count == 1
        post_call = api_client.post.call_args
        assert post_call[0][0] == "/experiment_runs/run-123/artifacts"
        json_data = post_call[1]["json"]
        assert json_data["artifactType"] == "model-artifact"
        assert json_data["name"] == "test-model"
        assert json_data["customProperties"]["mlflow__model_type"] == "sklearn"
        assert json_data["customProperties"]["mlflow__experiment_id"] == "exp-123"
        assert json_data["customProperties"]["mlflow__source_run_id"] == "run-123"
        assert json_data["customProperties"]["key1"] == "value1"
        assert json_data["customProperties"]["param_param1"] == "value1"
        assert (
            json_data["uri"]
            == "s3://bucket/artifacts/experiments/exp-123/run-123/test-model"
        )

    def test_search_logged_models(self, model_ops, api_client):
        """Test searching logged models."""
        # Mock runs response
        runs_data = {
            "items": [{"id": "run-123", "experimentId": "exp-123", "name": "test-run"}]
        }
        api_client.get.side_effect = [
            runs_data,  # GET /experiments/exp-123/experiment_runs
            {  # GET /experiment_runs/run-123/artifacts
                "items": [
                    {
                        "id": "model-123",
                        "name": "test-model",
                        "experimentId": "exp-123",
                        "uri": "s3://bucket/model",
                        "createTimeSinceEpoch": "1234567890",
                        "lastUpdateTimeSinceEpoch": "1234567890",
                        "artifactType": "model-artifact",
                        "customProperties": {
                            "mlflow__model_type": "sklearn",
                            "mlflow__source_run_id": "run-123",
                            "mlflow__experiment_id": "exp-123",
                        },
                    }
                ]
            },
        ]

        result = model_ops.search_logged_models(
            ["exp-123"],
            filter_string="model_type='sklearn'",
            max_results=10,
            page_token="token123",
        )

        assert isinstance(result, PagedList)
        assert len(result) == 1
        assert result[0].model_id == "model-123"
        assert result[0].name == "test-model"
        assert result[0].experiment_id == "exp-123"
        assert result[0].model_type == "sklearn"
        assert result[0].source_run_id == "run-123"

        # Check API calls
        assert api_client.get.call_count == 2

        # Check first call: GET /experiments/exp-123/experiment_runs
        call_args = api_client.get.call_args_list[0]
        assert call_args[0][0] == "/experiments/exp-123/experiment_runs"
        # First call should not have params since we're getting all runs from the experiment

        # Check second call: GET /experiment_runs/run-123/artifacts
        call_args = api_client.get.call_args_list[1]
        assert call_args[0][0] == "/experiment_runs/run-123/artifacts"
        params = call_args[1]["params"]
        assert params["artifactType"] == "model-artifact"

    def test_finalize_logged_model(self, model_ops, api_client):
        """Test finalizing a logged model."""
        model_data = {
            "id": "model-123",
            "name": "test-model",
            "experimentId": "exp-123",
            "uri": "s3://bucket/model",
            "createTimeSinceEpoch": "1234567890",
            "lastUpdateTimeSinceEpoch": "1234567890",
            "state": "LIVE",
            "customProperties": {
                "status": "READY",
                "mlflow__model_type": "sklearn",
            },
        }
        api_client.patch.return_value = model_data

        logged_model = model_ops.finalize_logged_model(
            "model-123", LoggedModelStatus.READY
        )

        assert isinstance(logged_model, LoggedModel)
        assert logged_model.model_id == "model-123"
        assert logged_model.model_type == "sklearn"
        assert logged_model.status == LoggedModelStatus.READY

        patch_call = api_client.patch.call_args
        assert patch_call[0][0] == "/artifacts/model-123"
        json_data = patch_call[1]["json"]
        assert json_data["state"] == "LIVE"
        assert json_data["artifactType"] == "model-artifact"

    def test_set_logged_model_tags(self, model_ops, api_client):
        """Test setting logged model tags."""
        model_data = {
            "id": "model-123",
            "name": "test-model",
            "customProperties": {"existing": "value"},
        }
        api_client.get.return_value = model_data
        api_client.patch.return_value = {}

        tags = [LoggedModelTag("key1", "value1"), LoggedModelTag("key2", "value2")]
        model_ops.set_logged_model_tags("model-123", tags)

        # Check GET call
        api_client.get.assert_called_once_with("/artifacts/model-123")

        # Check PATCH call
        patch_call = api_client.patch.call_args
        assert patch_call[0][0] == "/artifacts/model-123"
        json_data = patch_call[1]["json"]
        custom_props = json_data["customProperties"]
        assert custom_props["key1"] == "value1"
        assert custom_props["key2"] == "value2"

    def test_delete_logged_model_tag(self, model_ops, api_client):
        """Test deleting a logged model tag."""
        model_data = {
            "id": "model-123",
            "name": "test-model",
            "customProperties": {
                "key1": "value1",
                "key2": "value2",
            },
        }
        api_client.get.return_value = model_data
        api_client.patch.return_value = {}

        model_ops.delete_logged_model_tag("model-123", "key1")

        # Check GET call
        api_client.get.assert_called_once_with("/artifacts/model-123")

        # Check PATCH call
        patch_call = api_client.patch.call_args
        assert patch_call[0][0] == "/artifacts/model-123"
        json_data = patch_call[1]["json"]
        custom_props = json_data["customProperties"]
        assert "key1" not in custom_props
        assert "key2" in custom_props

    def test_get_logged_model(self, model_ops, api_client):
        """Test getting a logged model."""
        model_data = {
            "id": "model-123",
            "name": "test-model",
            "uri": "s3://bucket/model",
            "createTimeSinceEpoch": "1234567890",
            "lastUpdateTimeSinceEpoch": "1234567890",
            "customProperties": {
                "mlflow__experiment_id": "exp-123",
                "mlflow__model_type": "sklearn",
                "mlflow__source_run_id": "run-123",
                "param_lr": "0.01",
            },
        }
        api_client.get.return_value = model_data

        logged_model = model_ops.get_logged_model("model-123")

        assert isinstance(logged_model, LoggedModel)
        assert logged_model.model_id == "model-123"
        assert logged_model.name == "test-model"
        assert logged_model.experiment_id == "exp-123"
        assert logged_model.model_type == "sklearn"
        assert logged_model.source_run_id == "run-123"
        assert len(logged_model.params) == 1
        assert logged_model.params["lr"] == "0.01"

        api_client.get.assert_called_once_with("/artifacts/model-123")

    def test_delete_logged_model(self, model_ops, api_client):
        """Test deleting a logged model."""
        model_data = {
            "id": "model-123",
            "name": "test-model",
            "customProperties": {"existing": "value"},
        }
        api_client.get.return_value = model_data
        api_client.patch.return_value = {}

        model_ops.delete_logged_model("model-123")

        # Check GET call
        api_client.get.assert_called_once_with("/artifacts/model-123")

        # Check PATCH call
        patch_call = api_client.patch.call_args
        assert patch_call[0][0] == "/artifacts/model-123"
        json_data = patch_call[1]["json"]
        assert json_data["customProperties"]["state"] == "MARKED_FOR_DELETION"

    def test_search_logged_models_empty_response(self, model_ops, api_client):
        """Test searching logged models with empty response."""
        api_client.get.return_value = {"items": []}

        result = model_ops.search_logged_models(["exp-123"])

        assert isinstance(result, PagedList)
        assert len(result) == 0

        api_client.get.assert_called_once_with(
            "/experiments/exp-123/experiment_runs",
        )

    def test_search_logged_models_multiple_experiments(self, model_ops, api_client):
        """Test searching logged models across multiple experiments."""
        # Mock responses for multiple experiments
        runs_data1 = {"items": [{"id": "run-1"}]}
        runs_data2 = {"items": [{"id": "run-2"}]}
        artifacts_data1 = {
            "items": [
                {
                    "id": "model-1",
                    "name": "model1",
                    "experimentId": "exp-1",
                    "uri": "s3://bucket/model1",
                    "createTimeSinceEpoch": "1234567890",
                    "lastUpdateTimeSinceEpoch": "1234567890",
                    "artifactType": "model-artifact",
                    "customProperties": {
                        "mlflow__model_type": "sklearn",
                        "mlflow__source_run_id": "run-1",
                        "mlflow__experiment_id": "exp-1",
                    },
                }
            ]
        }
        artifacts_data2 = {
            "items": [
                {
                    "id": "model-2",
                    "name": "model2",
                    "experimentId": "exp-2",
                    "uri": "s3://bucket/model2",
                    "createTimeSinceEpoch": "1234567890",
                    "lastUpdateTimeSinceEpoch": "1234567890",
                    "artifactType": "model-artifact",
                    "customProperties": {
                        "mlflow__model_type": "sklearn",
                        "mlflow__source_run_id": "run-2",
                        "mlflow__experiment_id": "exp-2",
                    },
                }
            ]
        }

        api_client.get.side_effect = [
            runs_data1,  # GET /experiments/exp-1/experiment_runs
            artifacts_data1,  # GET /experiment_runs/run-1/artifacts
            runs_data2,  # GET /experiments/exp-2/experiment_runs
            artifacts_data2,  # GET /experiment_runs/run-2/artifacts
        ]

        result = model_ops.search_logged_models(["exp-1", "exp-2"])

        assert len(result) == 2
        assert result[0].model_id == "model-1"
        assert result[0].experiment_id == "exp-1"
        assert result[1].model_id == "model-2"
        assert result[1].experiment_id == "exp-2"

        # Verify API calls
        assert api_client.get.call_count == 4
        api_client.get.assert_any_call("/experiments/exp-1/experiment_runs")
        api_client.get.assert_any_call(
            "/experiment_runs/run-1/artifacts",
            params={"artifactType": "model-artifact"},
        )
        api_client.get.assert_any_call("/experiments/exp-2/experiment_runs")
        api_client.get.assert_any_call(
            "/experiment_runs/run-2/artifacts",
            params={"artifactType": "model-artifact"},
        )
