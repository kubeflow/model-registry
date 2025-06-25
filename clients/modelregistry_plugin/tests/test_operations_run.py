"""Tests for RunOperations."""

import os
from unittest.mock import Mock, patch

import pytest
from mlflow.entities import (
    DatasetInput,
    LoggedModelInput,
    LoggedModelOutput,
    Metric,
    Param,
    Run,
    RunInfo,
    RunStatus,
    RunTag,
)

from modelregistry_plugin.operations.run import RunOperations


class TestRunOperations:
    @pytest.fixture
    def api_client(self):
        """Create a mock API client."""
        return Mock()

    @pytest.fixture
    def run_ops(self, api_client):
        """Create a RunOperations instance for testing."""
        return RunOperations(api_client, "s3://bucket/artifacts")

    def test_init(self, run_ops, api_client):
        """Test RunOperations initialization."""
        assert run_ops.api_client == api_client
        assert run_ops.artifact_uri == "s3://bucket/artifacts"
        assert run_ops.DEFAULT_ARTIFACT_PAGE_SIZE == 1000

    def test_create_run(self, run_ops, api_client):
        """Test creating a run."""
        # Mock responses
        api_client.post.return_value = {
            "id": "run-123",
            "experimentId": "exp-123",
            "createTimeSinceEpoch": "1234567890",
        }
        api_client.get.return_value = {
            "id": "exp-123",
            "name": "test-experiment",
            "externalId": "s3://bucket/artifacts/experiments/exp-123",
        }
        api_client.patch.return_value = {}

        run = run_ops.create_run("exp-123", start_time=1234567890)

        assert isinstance(run, Run)
        assert run.info.run_id == "run-123"
        assert run.info.experiment_id == "exp-123"
        assert run.info.status == RunStatus.RUNNING
        assert (
            run.info.artifact_uri == "s3://bucket/artifacts/experiments/exp-123/run-123"
        )

        # Check POST call
        post_call = api_client.post.call_args
        assert post_call[0][0] == "/experiment_runs"
        json_data = post_call[1]["json"]
        assert json_data["experimentId"] == "exp-123"
        assert json_data["status"] == "RUNNING"
        assert json_data["startTimeSinceEpoch"] == "1234567890"

        # Check GET call
        api_client.get.assert_called_once_with("/experiments/exp-123")

        # Check PATCH call
        patch_call = api_client.patch.call_args
        assert patch_call[0][0] == "/experiment_runs/run-123"
        json_data = patch_call[1]["json"]
        assert (
            json_data["externalId"]
            == "s3://bucket/artifacts/experiments/exp-123/run-123"
        )

    def test_create_run_with_user_and_tags(self, run_ops, api_client):
        """Test creating a run with user ID and tags."""
        # Mock responses
        api_client.post.return_value = {
            "id": "run-123",
            "experimentId": "exp-123",
            "owner": "user123",
            "createTimeSinceEpoch": "1234567890",
            "customProperties": {
                "key1": "value1",
                "key2": "value2",
            },
        }
        api_client.get.return_value = {
            "id": "exp-123",
            "name": "test-experiment",
            "externalId": "s3://bucket/artifacts/experiments/exp-123",
        }
        api_client.patch.return_value = {}

        tags = [RunTag("key1", "value1"), RunTag("key2", "value2")]
        run = run_ops.create_run("exp-123", user_id="user123", tags=tags)

        # Check POST call
        post_call = api_client.post.call_args
        json_data = post_call[1]["json"]
        assert json_data["owner"] == "user123"
        custom_props = json_data["customProperties"]
        assert custom_props["key1"] == "value1"
        assert custom_props["key2"] == "value2"

        # Check returned run has expected tag keys and values
        assert run.data.tags.get("key1") == "value1"
        assert run.data.tags.get("key2") == "value2"
        assert run.info.user_id == "user123"

    def test_get_run(self, run_ops, api_client):
        """Test getting a run by ID."""
        # Mock run data
        run_data = {
            "id": "run-123",
            "experimentId": "exp-123",
            "name": "test-run",
            "state": "RUNNING",
            "owner": "user123",
            "startTimeSinceEpoch": "1234567890",
            "externalId": "s3://bucket/artifacts/experiments/exp-123/run-123",
            "customProperties": {"key1": "value1"},
        }
        api_client.get.side_effect = [
            run_data,  # GET run
            {"items": []},  # GET artifacts
        ]

        run = run_ops.get_run("run-123")

        assert isinstance(run, Run)
        assert run.info.run_id == "run-123"
        assert run.info.experiment_id == "exp-123"
        assert run.info.run_name == "test-run"
        assert run.info.user_id == "user123"
        assert (
            run.info.artifact_uri == "s3://bucket/artifacts/experiments/exp-123/run-123"
        )
        assert len(run.data.tags) == 1
        assert run.data.tags["key1"] == "value1"

        # Check API calls
        assert api_client.get.call_count == 2
        api_client.get.assert_any_call("/experiment_runs/run-123")
        api_client.get.assert_any_call(
            "/experiment_runs/run-123/artifacts", params={"pageSize": 1000}
        )

    def test_update_run_info(self, run_ops, api_client):
        """Test updating run information."""
        run_data = {
            "id": "run-123",
            "experimentId": "exp-123",
            "status": "FINISHED",
            "endTimeSinceEpoch": "1234567899",
            "startTimeSinceEpoch": "1234567890",
            "owner": "user123",
            "externalId": "s3://bucket/artifacts/experiments/exp-123/run-123",
        }
        api_client.patch.return_value = run_data

        run_info = run_ops.update_run_info(
            "run-123", RunStatus.FINISHED, end_time=1234567899
        )

        assert isinstance(run_info, RunInfo)
        assert run_info.run_id == "run-123"
        assert run_info.status == RunStatus.FINISHED
        assert (
            run_info.artifact_uri == "s3://bucket/artifacts/experiments/exp-123/run-123"
        )

        patch_call = api_client.patch.call_args
        assert patch_call[0][0] == "/experiment_runs/run-123"
        json_data = patch_call[1]["json"]
        assert json_data["status"] == "FINISHED"
        assert json_data["endTimeSinceEpoch"] == "1234567899"

    def test_delete_run(self, run_ops, api_client):
        """Test deleting a run."""
        api_client.patch.return_value = {}

        run_ops.delete_run("run-123")

        api_client.patch.assert_called_once_with(
            "/experiment_runs/run-123", json={"state": "ARCHIVED"}
        )

    def test_restore_run(self, run_ops, api_client):
        """Test restoring a run."""
        api_client.patch.return_value = {}

        run_ops.restore_run("run-123")

        api_client.patch.assert_called_once_with(
            "/experiment_runs/run-123", json={"state": "LIVE"}
        )

    def test_log_metric(self, run_ops, api_client):
        """Test logging a metric."""
        api_client.post.return_value = {}

        metric = Metric("accuracy", 0.95, 1234567890, 1)
        run_ops.log_metric("run-123", metric)

        post_call = api_client.post.call_args
        assert post_call[0][0] == "/experiment_runs/run-123/artifacts"
        json_data = post_call[1]["json"]
        assert json_data["artifactType"] == "metric"
        assert json_data["name"] == "accuracy"
        assert json_data["value"] == 0.95
        assert json_data["step"] == 1
        assert json_data["timestamp"] == "1234567890"

    def test_log_param(self, run_ops, api_client):
        """Test logging a parameter."""
        api_client.post.return_value = {}

        param = Param("learning_rate", "0.01")
        run_ops.log_param("run-123", param)

        post_call = api_client.post.call_args
        assert post_call[0][0] == "/experiment_runs/run-123/artifacts"
        json_data = post_call[1]["json"]
        assert json_data["artifactType"] == "parameter"
        assert json_data["name"] == "learning_rate"
        assert json_data["value"] == "0.01"
        assert json_data["parameterType"] == "string"

    def test_log_batch(self, run_ops, api_client):
        """Test batch logging."""
        run_data = {
            "id": "run-123",
            "experimentId": "exp-123",
            "customProperties": {"existing": "tag"},
        }
        api_client.get.return_value = run_data
        api_client.patch.return_value = {}
        api_client.post.return_value = {}

        metrics = [Metric("acc", 0.9, 1234567890, 1)]
        params = [Param("lr", "0.01")]
        tags = [RunTag("new", "tag")]

        run_ops.log_batch("run-123", metrics, params, tags)

        # Check GET call
        api_client.get.assert_called_once_with("/experiment_runs/run-123")

        # Check PATCH call for tags
        patch_call = api_client.patch.call_args
        assert patch_call[0][0] == "/experiment_runs/run-123"
        json_data = patch_call[1]["json"]
        custom_props = json_data["customProperties"]
        assert custom_props["new"] == "tag"

        # Check POST calls for metrics and params
        assert api_client.post.call_count == 2

    def test_log_inputs_datasets(self, run_ops, api_client):
        """Test logging dataset inputs."""
        api_client.post.return_value = {}

        dataset = Mock()
        dataset.dataset.name = "test-dataset"
        dataset.dataset.digest = "digest123"
        dataset.dataset.source_type = "csv"
        dataset.dataset.source = "s3://bucket/data.csv"
        dataset.dataset.schema = "schema"
        dataset.dataset.profile = "profile"
        dataset.tags = [RunTag("tag1", "value1")]

        dataset_input = DatasetInput(dataset=dataset.dataset, tags=dataset.tags)
        run_ops.log_inputs("run-123", datasets=[dataset_input])

        post_call = api_client.post.call_args
        assert post_call[0][0] == "/experiment_runs/run-123/artifacts"
        json_data = post_call[1]["json"]
        assert json_data["artifactType"] == "dataset-artifact"
        assert json_data["name"] == "test-dataset"
        assert json_data["digest"] == "digest123"
        assert json_data["sourceType"] == "csv"
        assert json_data["source"] == "s3://bucket/data.csv"
        assert json_data["schema"] == "schema"
        assert json_data["profile"] == "profile"
        assert json_data["customProperties"]["tag1"] == "value1"

    def test_log_inputs_models(self, run_ops, api_client):
        """Test logging model inputs."""
        model_data = {"customProperties": {"existing": "value"}}
        api_client.get.return_value = model_data
        api_client.post.return_value = {}

        model_input = LoggedModelInput(model_id="model-123")
        run_ops.log_inputs("run-123", models=[model_input])

        # Check GET call
        api_client.get.assert_called_once_with("/artifacts/model-123")

        # Check POST call
        post_call = api_client.post.call_args
        assert post_call[0][0] == "/experiment_runs/run-123/artifacts"
        json_data = post_call[1]["json"]
        custom_props = json_data["customProperties"]
        assert custom_props["mlflow__model_io_type"] == "input"

    def test_log_outputs(self, run_ops, api_client):
        """Test logging model outputs."""
        model_data = {"customProperties": {"existing": "value"}}
        api_client.get.return_value = model_data
        api_client.post.return_value = {}

        model_output = LoggedModelOutput(model_id="model-123", step=1)
        run_ops.log_outputs("run-123", [model_output])

        # Check GET call
        api_client.get.assert_called_once_with("/artifacts/model-123")

        # Check POST call
        post_call = api_client.post.call_args
        assert post_call[0][0] == "/experiment_runs/run-123/artifacts"
        json_data = post_call[1]["json"]
        custom_props = json_data["customProperties"]
        assert custom_props["mlflow__model_io_type"] == "output"

    def test_set_tag(self, run_ops, api_client):
        """Test setting a run tag."""
        run_data = {
            "id": "run-123",
            "experimentId": "exp-123",
            "customProperties": {"existing": "value"},
        }
        api_client.get.return_value = run_data
        api_client.patch.return_value = {}

        tag = RunTag("key1", "value1")
        run_ops.set_tag("run-123", tag)

        # Check GET call
        api_client.get.assert_called_once_with("/experiment_runs/run-123")

        # Check PATCH call
        patch_call = api_client.patch.call_args
        assert patch_call[0][0] == "/experiment_runs/run-123"
        json_data = patch_call[1]["json"]
        custom_props = json_data["customProperties"]
        assert custom_props["key1"] == "value1"

    def test_get_all_run_artifacts(self, run_ops, api_client):
        """Test getting all artifacts for a run with pagination."""
        # Mock responses for pagination
        response1 = {
            "items": [{"id": "artifact1"}, {"id": "artifact2"}],
            "nextPageToken": "token123",
        }
        response2 = {
            "items": [{"id": "artifact3"}],
            "nextPageToken": "",
        }
        api_client.get.side_effect = [response1, response2]

        artifacts = run_ops._get_all_run_artifacts("run-123")

        assert len(artifacts) == 3
        assert artifacts[0]["id"] == "artifact1"
        assert artifacts[1]["id"] == "artifact2"
        assert artifacts[2]["id"] == "artifact3"

        # Check API calls
        assert api_client.get.call_count == 2
        api_client.get.assert_any_call(
            "/experiment_runs/run-123/artifacts", params={"pageSize": 1000}
        )
        api_client.get.assert_any_call(
            "/experiment_runs/run-123/artifacts",
            params={"pageSize": 1000, "pageToken": "token123"},
        )

    @patch.dict(os.environ, {"MODEL_REGISTRY_ARTIFACT_PAGE_SIZE": "500"})
    def test_get_all_run_artifacts_with_custom_page_size(self, run_ops, api_client):
        """Test getting artifacts with custom page size from environment."""
        api_client.get.return_value = {"items": [], "nextPageToken": ""}

        run_ops._get_all_run_artifacts("run-123")

        api_client.get.assert_called_once_with(
            "/experiment_runs/run-123/artifacts", params={"pageSize": 500}
        )

    def test_get_all_run_artifacts_single_page(self, run_ops, api_client):
        """Test getting artifacts with single page response."""
        api_client.get.return_value = {
            "items": [{"id": "artifact1"}],
            "nextPageToken": "",
        }

        artifacts = run_ops._get_all_run_artifacts("run-123")

        assert len(artifacts) == 1
        assert artifacts[0]["id"] == "artifact1"

        api_client.get.assert_called_once()

    def test_get_all_run_artifacts_empty_response(self, run_ops, api_client):
        """Test getting artifacts with empty response."""
        api_client.get.return_value = {"items": [], "nextPageToken": ""}

        artifacts = run_ops._get_all_run_artifacts("run-123")

        assert len(artifacts) == 0
        api_client.get.assert_called_once()
