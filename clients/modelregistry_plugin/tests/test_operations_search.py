"""Tests for SearchOperations."""

from unittest.mock import Mock

import pytest
from mlflow.entities import ViewType
from mlflow.store.entities.paged_list import PagedList

from modelregistry_plugin.operations.search import SearchOperations


class TestSearchOperations:
    @pytest.fixture
    def api_client(self):
        """Create a mock API client."""
        return Mock()

    @pytest.fixture
    def search_ops(self, api_client):
        """Create a SearchOperations instance for testing."""
        return SearchOperations(api_client, "s3://bucket/artifacts")

    def test_init(self, search_ops, api_client):
        """Test SearchOperations initialization."""
        assert search_ops.api_client == api_client
        assert search_ops.artifact_uri == "s3://bucket/artifacts"

    def test_search_runs_all(self, search_ops, api_client):
        """Test searching runs with all parameters."""
        search_data = {
            "items": [
                {
                    "id": "run-123",
                    "experimentId": "exp-123",
                    "name": "test-run",
                    "status": "RUNNING",
                    "owner": "user123",
                    "startTimeSinceEpoch": "1234567890",
                    "externalId": "s3://bucket/artifacts/experiments/exp-123/run-123",
                    "customProperties": {"tag1": "value1"},
                }
            ],
            "nextPageToken": "token123",
        }
        artifacts_data = {"items": []}

        # Only one run, so only one artifact call
        api_client.get.side_effect = [search_data, artifacts_data]

        result = search_ops.search_runs(
            ["exp-123"], max_results=10, page_token="token123"
        )

        assert isinstance(result, PagedList)
        assert len(result) == 1
        assert result[0].info.run_id == "run-123"
        assert (
            result[0].info.artifact_uri
            == "s3://bucket/artifacts/experiments/exp-123/run-123"
        )
        assert result.token is None  # Token is not returned in current implementation

        # Verify API calls
        assert api_client.get.call_count == 2
        api_client.get.assert_any_call(
            "/experiments/exp-123/experiment_runs",
            params={"pageSize": "10", "pageToken": "token123"},
        )
        api_client.get.assert_any_call("/experiment_runs/run-123/artifacts")

    def test_search_runs_active_only(self, search_ops, api_client):
        """Test searching runs with ACTIVE_ONLY view type."""
        search_data = {
            "items": [
                {
                    "id": "run-123",
                    "experimentId": "exp-123",
                    "name": "active-run",
                    "status": "RUNNING",
                    "state": "LIVE",
                    "owner": "user123",
                    "startTimeSinceEpoch": "1234567890",
                    "externalId": "s3://bucket/artifacts/experiments/exp-123/run-123",
                    "customProperties": {},
                },
                {
                    "id": "run-456",
                    "experimentId": "exp-123",
                    "name": "deleted-run",
                    "status": "FINISHED",
                    "state": "ARCHIVED",
                    "owner": "user123",
                    "startTimeSinceEpoch": "1234567890",
                    "externalId": "s3://bucket/artifacts/experiments/exp-123/run-456",
                    "customProperties": {},
                },
            ],
            "nextPageToken": "token123",
        }
        artifacts_data = {"items": []}

        # Two runs in data, but only one returned after filtering, so only one artifact call
        api_client.get.side_effect = [search_data, artifacts_data]

        result = search_ops.search_runs(["exp-123"], run_view_type=ViewType.ACTIVE_ONLY)

        # Should return only the active run
        assert isinstance(result, PagedList)
        assert len(result) == 1
        assert result[0].info.run_id == "run-123"
        assert (
            result[0].info.artifact_uri
            == "s3://bucket/artifacts/experiments/exp-123/run-123"
        )

        # Verify API calls
        assert api_client.get.call_count == 2
        api_client.get.assert_any_call(
            "/experiments/exp-123/experiment_runs",
            params={"pageSize": "1000"},
        )
        api_client.get.assert_any_call("/experiment_runs/run-123/artifacts")

    def test_search_runs_deleted_only(self, search_ops, api_client):
        """Test searching runs with DELETED_ONLY view type."""
        search_data = {
            "items": [
                {
                    "id": "run-123",
                    "experimentId": "exp-123",
                    "name": "active-run",
                    "status": "RUNNING",
                    "state": "LIVE",
                    "owner": "user123",
                    "startTimeSinceEpoch": "1234567890",
                    "externalId": "s3://bucket/artifacts/experiments/exp-123/run-123",
                    "customProperties": {},
                },
                {
                    "id": "run-456",
                    "experimentId": "exp-123",
                    "name": "deleted-run",
                    "status": "FINISHED",
                    "state": "ARCHIVED",
                    "owner": "user123",
                    "startTimeSinceEpoch": "1234567890",
                    "externalId": "s3://bucket/artifacts/experiments/exp-123/run-456",
                    "customProperties": {},
                },
            ],
            "nextPageToken": "token123",
        }
        artifacts_data = {"items": []}

        # Two runs in data, but only one returned after filtering, so only one artifact call
        api_client.get.side_effect = [search_data, artifacts_data]

        result = search_ops.search_runs(
            ["exp-123"], run_view_type=ViewType.DELETED_ONLY
        )

        # Should return only the deleted run
        assert isinstance(result, PagedList)
        assert len(result) == 1
        assert result[0].info.run_id == "run-456"
        assert (
            result[0].info.artifact_uri
            == "s3://bucket/artifacts/experiments/exp-123/run-456"
        )

        # Verify API calls
        assert api_client.get.call_count == 2
        api_client.get.assert_any_call(
            "/experiments/exp-123/experiment_runs",
            params={"pageSize": "1000"},
        )
        api_client.get.assert_any_call("/experiment_runs/run-456/artifacts")

    def test_search_runs_with_filter_string(self, search_ops, api_client):
        """Test searching runs with filter string (should be ignored for now)."""
        search_data = {"items": []}
        api_client.get.return_value = search_data

        # This should not raise an error even though filter_string is not supported yet
        result = search_ops.search_runs(["exp-123"], filter_string="status='RUNNING'")

        assert isinstance(result, PagedList)
        assert len(result) == 0

        # Verify API call
        api_client.get.assert_called_once_with(
            "/experiments/exp-123/experiment_runs",
            params={"pageSize": "1000"},
        )

    def test_search_runs_with_order_by(self, search_ops, api_client):
        """Test searching runs with order_by (should be ignored for now)."""
        search_data = {"items": []}
        api_client.get.return_value = search_data

        # This should not raise an error even though order_by is not supported yet
        result = search_ops.search_runs(["exp-123"], order_by=["name"])

        assert isinstance(result, PagedList)
        assert len(result) == 0

        # Verify API call
        api_client.get.assert_called_once_with(
            "/experiments/exp-123/experiment_runs",
            params={"pageSize": "1000"},
        )

    def test_search_runs_multiple_experiments(self, search_ops, api_client):
        """Test searching runs across multiple experiments."""
        search_data1 = {
            "items": [
                {
                    "id": "run-1",
                    "experimentId": "exp-1",
                    "name": "run1",
                    "state": "RUNNING",
                    "owner": "user123",
                    "startTimeSinceEpoch": "1234567890",
                    "externalId": "s3://bucket/artifacts/experiments/exp-1/run-1",
                    "customProperties": {},
                }
            ]
        }
        search_data2 = {
            "items": [
                {
                    "id": "run-2",
                    "experimentId": "exp-2",
                    "name": "run2",
                    "state": "RUNNING",
                    "owner": "user123",
                    "startTimeSinceEpoch": "1234567890",
                    "externalId": "s3://bucket/artifacts/experiments/exp-2/run-2",
                    "customProperties": {},
                }
            ]
        }
        artifacts_data = {"items": []}

        api_client.get.side_effect = [
            search_data1,
            artifacts_data,
            search_data2,
            artifacts_data,
        ]

        result = search_ops.search_runs(["exp-1", "exp-2"])

        assert len(result) == 2
        assert result[0].info.run_id == "run-1"
        assert result[0].info.experiment_id == "exp-1"
        assert result[1].info.run_id == "run-2"
        assert result[1].info.experiment_id == "exp-2"

        # Verify API calls
        assert api_client.get.call_count == 4
        api_client.get.assert_any_call(
            "/experiments/exp-1/experiment_runs",
            params={"pageSize": "1000"},
        )
        api_client.get.assert_any_call("/experiment_runs/run-1/artifacts")
        api_client.get.assert_any_call(
            "/experiments/exp-2/experiment_runs",
            params={"pageSize": "1000"},
        )
        api_client.get.assert_any_call("/experiment_runs/run-2/artifacts")

    def test_search_runs_empty_response(self, search_ops, api_client):
        """Test searching runs with empty response."""
        search_data = {"items": []}
        api_client.get.return_value = search_data

        result = search_ops.search_runs(["exp-123"])

        assert isinstance(result, PagedList)
        assert len(result) == 0

        api_client.get.assert_called_once_with(
            "/experiments/exp-123/experiment_runs",
            params={"pageSize": "1000"},
        )

    def test_search_runs_with_pagination(self, search_ops, api_client):
        """Test searching runs with pagination parameters."""
        search_data = {
            "items": [
                {
                    "id": "run-123",
                    "experimentId": "exp-123",
                    "name": "test-run",
                    "status": "RUNNING",
                    "owner": "user123",
                    "startTimeSinceEpoch": "1234567890",
                    "externalId": "s3://bucket/artifacts/experiments/exp-123/run-123",
                    "customProperties": {},
                }
            ],
            "nextPageToken": "token123",
        }
        artifacts_data = {"items": []}

        # Only one run, so only one artifact call
        api_client.get.side_effect = [search_data, artifacts_data]

        result = search_ops.search_runs(
            ["exp-123"], max_results=5, page_token="token123"
        )

        assert isinstance(result, PagedList)
        assert len(result) == 1
        assert result[0].info.run_id == "run-123"

        # Verify API calls
        assert api_client.get.call_count == 2
        api_client.get.assert_any_call(
            "/experiments/exp-123/experiment_runs",
            params={"pageSize": "5", "pageToken": "token123"},
        )
        api_client.get.assert_any_call("/experiment_runs/run-123/artifacts")

    def test_search_runs_with_artifacts(self, search_ops, api_client):
        """Test searching runs with artifacts."""
        search_data = {
            "items": [
                {
                    "id": "run-123",
                    "experimentId": "exp-123",
                    "name": "test-run",
                    "status": "RUNNING",
                    "owner": "user123",
                    "startTimeSinceEpoch": "1234567890",
                    "externalId": "s3://bucket/artifacts/experiments/exp-123/run-123",
                    "customProperties": {},
                }
            ]
        }
        artifacts_data = {
            "items": [
                {
                    "artifactType": "metric",
                    "name": "accuracy",
                    "value": 0.95,
                    "timestamp": 1234567890,
                    "step": 1,
                },
                {
                    "artifactType": "parameter",
                    "name": "learning_rate",
                    "value": "0.01",
                },
                {
                    "artifactType": "dataset-artifact",
                    "name": "test-dataset",
                    "digest": "digest123",
                    "sourceType": "csv",
                    "source": "s3://bucket/data.csv",
                    "schema": "schema",
                    "profile": "profile",
                    "customProperties": {},
                },
            ]
        }

        # Only one run, so only one artifact call
        api_client.get.side_effect = [search_data, artifacts_data]

        result = search_ops.search_runs(["exp-123"])

        assert isinstance(result, PagedList)
        assert len(result) == 1
        run = result[0]
        assert run.info.run_id == "run-123"
        assert (
            run.info.artifact_uri == "s3://bucket/artifacts/experiments/exp-123/run-123"
        )

        # Check that artifacts were processed
        assert len(run.data.metrics) == 1
        assert "accuracy" in run.data.metrics
        assert run.data.metrics["accuracy"] == 0.95

        assert len(run.data.params) == 1
        assert run.data.params["learning_rate"] == "0.01"

        assert len(run.inputs.dataset_inputs) == 1
        assert run.inputs.dataset_inputs[0].dataset.name == "test-dataset"

        # Verify API calls
        assert api_client.get.call_count == 2
        api_client.get.assert_any_call(
            "/experiments/exp-123/experiment_runs",
            params={"pageSize": "1000"},
        )
        api_client.get.assert_any_call("/experiment_runs/run-123/artifacts")

    def test_search_runs_with_metrics_and_params(self, search_ops, api_client):
        """Test searching runs with metrics and parameters."""
        search_data = {
            "items": [
                {
                    "id": "run-123",
                    "experimentId": "exp-123",
                    "name": "test-run",
                    "status": "RUNNING",
                    "owner": "user123",
                    "startTimeSinceEpoch": "1234567890",
                    "externalId": "s3://bucket/artifacts/experiments/exp-123/run-123",
                    "customProperties": {},
                }
            ]
        }
        artifacts_data = {
            "items": [
                {
                    "artifactType": "metric",
                    "name": "accuracy",
                    "value": 0.95,
                    "timestamp": 1234567890,
                    "step": 1,
                },
                {
                    "artifactType": "parameter",
                    "name": "learning_rate",
                    "value": "0.01",
                },
            ]
        }

        # Only one run, so only one artifact call
        api_client.get.side_effect = [search_data, artifacts_data]

        result = search_ops.search_runs(["exp-123"])

        assert isinstance(result, PagedList)
        assert len(result) == 1
        run = result[0]
        assert run.info.run_id == "run-123"

        # Check metrics and params
        assert len(run.data.metrics) == 1
        assert "accuracy" in run.data.metrics
        assert run.data.metrics["accuracy"] == 0.95

        assert len(run.data.params) == 1
        assert run.data.params["learning_rate"] == "0.01"

        # Verify API calls
        assert api_client.get.call_count == 2
        api_client.get.assert_any_call(
            "/experiments/exp-123/experiment_runs",
            params={"pageSize": "1000"},
        )
        api_client.get.assert_any_call("/experiment_runs/run-123/artifacts")

    def test_search_runs_default_view_type(self, search_ops, api_client):
        """Test that search_runs defaults to ViewType.ALL when no run_view_type is specified."""
        search_data = {
            "items": [
                {
                    "id": "run-123",
                    "experimentId": "exp-123",
                    "name": "active-run",
                    "status": "RUNNING",
                    "state": "LIVE",
                    "owner": "user123",
                    "startTimeSinceEpoch": "1234567890",
                    "externalId": "s3://bucket/artifacts/experiments/exp-123/run-123",
                    "customProperties": {},
                },
                {
                    "id": "run-456",
                    "experimentId": "exp-123",
                    "name": "deleted-run",
                    "status": "FINISHED",
                    "state": "ARCHIVED",
                    "owner": "user123",
                    "startTimeSinceEpoch": "1234567890",
                    "externalId": "s3://bucket/artifacts/experiments/exp-123/run-456",
                    "customProperties": {},
                },
            ]
        }
        artifacts_data = {"items": []}

        api_client.get.side_effect = [search_data, artifacts_data, artifacts_data]

        # Call without specifying run_view_type - should default to ViewType.ALL
        result = search_ops.search_runs(["exp-123"])

        # Should return both active and deleted runs (ALL view type)
        assert isinstance(result, PagedList)
        assert len(result) == 2
        assert result[0].info.run_id == "run-123"
        assert result[1].info.run_id == "run-456"

        # Verify API calls
        assert api_client.get.call_count == 3
        api_client.get.assert_any_call(
            "/experiments/exp-123/experiment_runs",
            params={"pageSize": "1000"},
        )
        api_client.get.assert_any_call("/experiment_runs/run-123/artifacts")

    def test_search_datasets(self, search_ops, api_client):
        """Test searching for datasets across experiments."""
        runs_data1 = {"items": [{"id": "run-1", "experimentId": "exp-1"}]}
        runs_data2 = {"items": [{"id": "run-2", "experimentId": "exp-2"}]}
        artifacts_data1 = {
            "items": [
                {
                    "name": "dataset1",
                    "digest": "digest1",
                    "artifactType": "dataset-artifact",
                    "customProperties": {"mlflow.dataset.context": "training"},
                }
            ]
        }
        artifacts_data2 = {
            "items": [
                {
                    "name": "dataset2",
                    "digest": "digest2",
                    "artifactType": "dataset-artifact",
                    "customProperties": {"mlflow.dataset.context": "validation"},
                }
            ]
        }

        api_client.get.side_effect = [
            runs_data1,  # GET /experiments/exp-1/experiment_runs
            artifacts_data1,  # GET /experiment_runs/run-1/artifacts
            runs_data2,  # GET /experiments/exp-2/experiment_runs
            artifacts_data2,  # GET /experiment_runs/run-2/artifacts
        ]

        result = search_ops._search_datasets(["exp-1", "exp-2"])

        assert len(result) == 2
        assert result[0].experiment_id == "exp-1"
        assert result[0].name == "dataset1"
        assert result[0].digest == "digest1"
        assert result[0].context == "training"
        assert result[1].experiment_id == "exp-2"
        assert result[1].name == "dataset2"
        assert result[1].digest == "digest2"
        assert result[1].context == "validation"

        # Verify API calls
        assert api_client.get.call_count == 4
        api_client.get.assert_any_call("/experiments/exp-1/experiment_runs")
        api_client.get.assert_any_call(
            "/experiment_runs/run-1/artifacts",
            params={"artifactType": "dataset-artifact"},
        )
        api_client.get.assert_any_call("/experiments/exp-2/experiment_runs")
        api_client.get.assert_any_call(
            "/experiment_runs/run-2/artifacts",
            params={"artifactType": "dataset-artifact"},
        )

    def test_search_datasets_deduplication(self, search_ops, api_client):
        """Test that _search_datasets deduplicates datasets."""
        runs_data = {
            "items": [
                {"id": "run-1", "experimentId": "exp-1"},
                {"id": "run-2", "experimentId": "exp-1"},
            ]
        }
        # Both runs have the same dataset
        artifacts_data = {
            "items": [
                {
                    "name": "dataset1",
                    "digest": "digest1",
                    "artifactType": "dataset-artifact",
                    "customProperties": {},
                }
            ]
        }

        api_client.get.side_effect = [
            runs_data,  # GET /experiments/exp-1/experiment_runs
            artifacts_data,  # GET /experiment_runs/run-1/artifacts
            artifacts_data,  # GET /experiment_runs/run-2/artifacts
        ]

        result = search_ops._search_datasets(["exp-1"])

        # Should only return one dataset despite being in two runs
        assert len(result) == 1
        assert result[0].experiment_id == "exp-1"
        assert result[0].name == "dataset1"
        assert result[0].digest == "digest1"

        # Verify API calls
        assert api_client.get.call_count == 3

    def test_search_datasets_empty_response(self, search_ops, api_client):
        """Test _search_datasets with no datasets found."""
        runs_data = {"items": [{"id": "run-1", "experimentId": "exp-1"}]}
        artifacts_data = {"items": []}  # No dataset artifacts

        api_client.get.side_effect = [runs_data, artifacts_data]

        result = search_ops._search_datasets(["exp-1"])

        assert len(result) == 0

    def test_search_datasets_no_runs(self, search_ops, api_client):
        """Test _search_datasets with no runs in experiment."""
        runs_data = {"items": []}  # No runs

        api_client.get.return_value = runs_data

        result = search_ops._search_datasets(["exp-1"])

        assert len(result) == 0
        api_client.get.assert_called_once_with("/experiments/exp-1/experiment_runs")
