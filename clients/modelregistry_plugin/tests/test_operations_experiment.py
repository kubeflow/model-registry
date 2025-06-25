"""Tests for ExperimentOperations."""

from unittest.mock import Mock

import pytest
from mlflow.entities import Experiment, ExperimentTag, ViewType
from mlflow.exceptions import MlflowException
from mlflow.store.entities.paged_list import PagedList

from modelregistry_plugin.operations.experiment import ExperimentOperations


class TestExperimentOperations:
    @pytest.fixture
    def api_client(self):
        """Create a mock API client."""
        return Mock()

    @pytest.fixture
    def experiment_ops(self, api_client):
        """Create an ExperimentOperations instance for testing."""
        return ExperimentOperations(api_client, "s3://bucket/artifacts")

    def test_init(self, experiment_ops, api_client):
        """Test ExperimentOperations initialization."""
        assert experiment_ops.api_client == api_client
        assert experiment_ops.artifact_uri == "s3://bucket/artifacts"

    def test_create_experiment(self, experiment_ops, api_client):
        """Test creating an experiment."""
        api_client.post.return_value = {"id": "exp-123"}
        api_client.patch.return_value = {}

        experiment_id = experiment_ops.create_experiment("test-experiment")

        assert experiment_id == "exp-123"
        assert api_client.post.call_count == 1
        assert api_client.patch.call_count == 1

        # Check POST call
        post_call = api_client.post.call_args
        assert post_call[0][0] == "/experiments"
        json_data = post_call[1]["json"]
        assert json_data["name"] == "test-experiment"
        assert json_data["description"] == "MLflow experiment: test-experiment"
        assert json_data["state"] == "LIVE"
        assert json_data["customProperties"] == {}

        # Check PATCH call
        patch_call = api_client.patch.call_args
        assert patch_call[0][0] == "/experiments/exp-123"
        json_data = patch_call[1]["json"]
        assert json_data["externalId"] == "s3://bucket/artifacts/experiments/exp-123"

    def test_create_experiment_with_artifact_location(self, experiment_ops, api_client):
        """Test creating an experiment with explicit artifact location."""
        api_client.post.return_value = {"id": "exp-123"}

        experiment_id = experiment_ops.create_experiment(
            "test-experiment", artifact_location="s3://custom/location"
        )

        assert experiment_id == "exp-123"
        assert api_client.post.call_count == 1
        assert api_client.patch.call_count == 0

        post_call = api_client.post.call_args
        json_data = post_call[1]["json"]
        assert json_data["externalId"] == "s3://custom/location"

    def test_create_experiment_with_tags(self, experiment_ops, api_client):
        """Test creating an experiment with tags."""
        api_client.post.return_value = {"id": "exp-123"}
        api_client.patch.return_value = {}

        tags = [ExperimentTag("key1", "value1"), ExperimentTag("key2", "value2")]
        experiment_id = experiment_ops.create_experiment("test-experiment", tags=tags)

        assert experiment_id == "exp-123"

        post_call = api_client.post.call_args
        json_data = post_call[1]["json"]
        custom_props = json_data["customProperties"]
        assert custom_props["key1"] == "value1"
        assert custom_props["key2"] == "value2"

    def test_create_experiment_without_artifact_uri(self, api_client):
        """Test creating an experiment without artifact URI."""
        experiment_ops = ExperimentOperations(api_client, None)
        api_client.post.return_value = {"id": "exp-123"}

        experiment_id = experiment_ops.create_experiment("test-experiment")

        assert experiment_id == "exp-123"

        post_call = api_client.post.call_args
        json_data = post_call[1]["json"]
        assert json_data["externalId"] is None

    def test_get_experiment(self, experiment_ops, api_client):
        """Test getting an experiment by ID."""
        experiment_data = {
            "id": "exp-123",
            "name": "test-experiment",
            "externalId": "s3://bucket/artifacts/exp-123",
            "state": "LIVE",
            "customProperties": {"key1": "value1"},
        }
        api_client.get.return_value = experiment_data

        experiment = experiment_ops.get_experiment("exp-123")

        assert isinstance(experiment, Experiment)
        assert experiment.experiment_id == "exp-123"
        assert experiment.name == "test-experiment"
        assert experiment.artifact_location == "s3://bucket/artifacts/exp-123"
        assert len(experiment.tags) == 1
        assert experiment.tags["key1"] == "value1"

        api_client.get.assert_called_once_with("/experiments/exp-123")

    def test_get_experiment_by_name_found(self, experiment_ops, api_client):
        """Test getting an experiment by name when found."""
        experiment_data = {
            "id": "exp-123",
            "name": "test-experiment",
            "externalId": "s3://bucket/artifacts/exp-123",
            "state": "LIVE",
            "customProperties": {},
        }
        api_client.get.return_value = experiment_data

        experiment = experiment_ops.get_experiment_by_name("test-experiment")

        assert isinstance(experiment, Experiment)
        assert experiment.experiment_id == "exp-123"
        assert experiment.name == "test-experiment"

        api_client.get.assert_called_once_with(
            "/experiment", params={"name": "test-experiment"}
        )

    def test_get_experiment_by_name_not_found(self, experiment_ops, api_client):
        """Test getting an experiment by name when not found."""
        from mlflow.exceptions import MlflowException

        # Mock the API client to raise a proper MlflowException with 404 status
        # The API client creates exceptions with error_code, so we need to mock that
        exception = MlflowException("Model Registry API error: not found")
        # Set the attributes that the implementation checks
        exception._http_status_code = 404
        # The message should be the full error message, not just "not found"
        exception.message = "Model Registry API error: not found"
        # Set the error_code that the API client would set
        exception.error_code = "RESOURCE_DOES_NOT_EXIST"

        # Mock the API client to raise this exception
        api_client.get.side_effect = exception

        experiment = experiment_ops.get_experiment_by_name("nonexistent")

        assert experiment is None
        api_client.get.assert_called_once_with(
            "/experiment", params={"name": "nonexistent"}
        )

    def test_get_experiment_by_name_other_error(self, experiment_ops, api_client):
        """Test getting an experiment by name with other error."""
        api_client.get.side_effect = MlflowException("server error", error_code=500)

        with pytest.raises(MlflowException) as exc_info:
            experiment_ops.get_experiment_by_name("test-experiment")

        assert "server error" in str(exc_info.value)

    def test_delete_experiment(self, experiment_ops, api_client):
        """Test deleting an experiment."""
        api_client.patch.return_value = {}

        experiment_ops.delete_experiment("exp-123")

        api_client.patch.assert_called_once_with(
            "/experiments/exp-123", json={"state": "ARCHIVED"}
        )

    def test_restore_experiment(self, experiment_ops, api_client):
        """Test restoring an experiment."""
        api_client.patch.return_value = {}

        experiment_ops.restore_experiment("exp-123")

        api_client.patch.assert_called_once_with(
            "/experiments/exp-123", json={"state": "LIVE"}
        )

    def test_rename_experiment(self, experiment_ops, api_client):
        """Test renaming an experiment."""
        api_client.patch.return_value = {}

        experiment_ops.rename_experiment("exp-123", "new-name")

        api_client.patch.assert_called_once_with(
            "/experiments/exp-123", json={"name": "new-name"}
        )

    def test_search_experiments(self, experiment_ops, api_client):
        """Test searching experiments."""
        response_data = {
            "items": [
                {
                    "id": "1",
                    "name": "exp1",
                    "state": "LIVE",
                    "externalId": "s3://bucket/artifacts/exp1",
                    "customProperties": {},
                },
                {
                    "id": "2",
                    "name": "exp2",
                    "state": "ARCHIVED",
                    "externalId": "s3://bucket/artifacts/exp2",
                    "customProperties": {},
                },
            ],
            "nextPageToken": "token123",
        }
        api_client.get.return_value = response_data

        result = experiment_ops.search_experiments(
            view_type=ViewType.ALL,
            max_results=10,
            filter_string="name='exp1'",
            order_by=["name"],
            page_token="token123",
        )

        assert isinstance(result, PagedList)
        assert len(result) == 2
        assert result[0].experiment_id == "1"
        assert result[0].name == "exp1"
        assert result[1].experiment_id == "2"
        assert result[1].name == "exp2"
        assert result.token == "token123"

        api_client.get.assert_called_once_with(
            "/experiments", params={"pageSize": 10, "pageToken": "token123"}
        )

    def test_search_experiments_active_only(self, experiment_ops, api_client):
        """Test searching experiments with ACTIVE_ONLY view type."""
        response_data = {
            "items": [
                {
                    "id": "1",
                    "name": "exp1",
                    "state": "LIVE",
                    "externalId": "s3://bucket/artifacts/exp1",
                    "customProperties": {},
                },
                {
                    "id": "2",
                    "name": "exp2",
                    "state": "ARCHIVED",
                    "externalId": "s3://bucket/artifacts/exp2",
                    "customProperties": {},
                },
            ]
        }
        api_client.get.return_value = response_data

        result = experiment_ops.search_experiments(view_type=ViewType.ACTIVE_ONLY)

        # Should only return active experiments
        assert len(result) == 1
        assert result[0].experiment_id == "1"
        assert result[0].name == "exp1"

    def test_search_experiments_deleted_only(self, experiment_ops, api_client):
        """Test searching experiments with DELETED_ONLY view type."""
        response_data = {
            "items": [
                {
                    "id": "1",
                    "name": "exp1",
                    "state": "LIVE",
                    "externalId": "s3://bucket/artifacts/exp1",
                    "customProperties": {},
                },
                {
                    "id": "2",
                    "name": "exp2",
                    "state": "ARCHIVED",
                    "externalId": "s3://bucket/artifacts/exp2",
                    "customProperties": {},
                },
            ]
        }
        api_client.get.return_value = response_data

        result = experiment_ops.search_experiments(view_type=ViewType.DELETED_ONLY)

        # Should only return deleted experiments
        assert len(result) == 1
        assert result[0].experiment_id == "2"
        assert result[0].name == "exp2"

    def test_set_experiment_tag(self, experiment_ops, api_client):
        """Test setting an experiment tag."""
        experiment_data = {
            "id": "exp-123",
            "name": "test-experiment",
            "customProperties": {"existing": "value"},
        }
        api_client.get.return_value = experiment_data
        api_client.patch.return_value = {}

        tag = ExperimentTag("key1", "value1")
        experiment_ops.set_experiment_tag("exp-123", tag)

        # Check GET call
        api_client.get.assert_called_once_with("/experiments/exp-123")

        # Check PATCH call
        patch_call = api_client.patch.call_args
        assert patch_call[0][0] == "/experiments/exp-123"
        json_data = patch_call[1]["json"]
        custom_props = json_data["customProperties"]
        assert custom_props["key1"] == "value1"
        assert custom_props["existing"] == "value"
