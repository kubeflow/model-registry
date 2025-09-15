"""Tests for MetricOperations."""

from unittest.mock import Mock

import pytest

from model_registry_mlflow.operations.metric import MetricOperations


class TestMetricOperations:
    @pytest.fixture
    def api_client(self):
        """Create a mock API client."""
        return Mock()

    @pytest.fixture
    def metric_ops(self, api_client):
        """Create a MetricOperations instance for testing."""
        return MetricOperations(api_client)

    def test_init(self, metric_ops, api_client):
        """Test MetricOperations initialization."""
        assert metric_ops.api_client == api_client

    def test_get_metric_history(self, metric_ops, api_client):
        """Test getting metric history for a specific metric key."""
        response_data = {
            "items": [
                {
                    "name": "accuracy",
                    "value": 0.95,
                    "timestamp": "1234567890",
                    "step": 1,
                    "createTimeSinceEpoch": "1234567890",
                    "artifactType": "metric",
                },
                {
                    "name": "accuracy",
                    "value": 0.97,
                    "timestamp": "1234567891",
                    "step": 2,
                    "createTimeSinceEpoch": "1234567891",
                    "artifactType": "metric",
                },
            ],
            "nextPageToken": "token202",
        }
        api_client.get.return_value = response_data

        result = metric_ops.get_metric_history("run-123", "accuracy")

        assert len(result) == 2
        assert result[0].value == 0.95
        assert result[0].step == 1
        assert result[0].timestamp == 1234567890
        assert result[1].value == 0.97
        assert result[1].step == 2
        assert result[1].timestamp == 1234567891
        assert all(metric.key == "accuracy" for metric in result)

        # Verify API call
        api_client.get.assert_called_once_with(
            "/experiment_runs/run-123/metric_history", params={"name": "accuracy"}
        )

        # Verify nextPageToken is handled correctly
        assert result.token == "token202"

    def test_get_metric_history_empty(self, metric_ops, api_client):
        """Test getting metric history when no metrics exist."""
        api_client.get.return_value = {"items": []}

        metrics = metric_ops.get_metric_history("run-123", "nonexistent")

        assert len(metrics) == 0

        # Verify API call
        api_client.get.assert_called_once_with(
            "/experiment_runs/run-123/metric_history", params={"name": "nonexistent"}
        )

    def test_get_metric_history_uses_timestamp_fallback(self, metric_ops, api_client):
        """Test that metric history uses createTimeSinceEpoch when timestamp is not available."""
        response_data = {
            "items": [
                {
                    "name": "accuracy",
                    "value": 0.95,
                    "step": 1,
                    "createTimeSinceEpoch": "1234567890",
                    "artifactType": "metric",
                    # Note: no "timestamp" field
                }
            ]
        }
        api_client.get.return_value = response_data

        metrics = metric_ops.get_metric_history("run-123", "accuracy")

        assert len(metrics) == 1
        assert metrics[0].timestamp == 1234567890  # Should use createTimeSinceEpoch

    def test_get_metric_history_uses_timestamp_over_create_time(
        self, metric_ops, api_client
    ):
        """Test that metric history prefers timestamp over createTimeSinceEpoch."""
        response_data = {
            "items": [
                {
                    "name": "accuracy",
                    "value": 0.95,
                    "timestamp": "1234567890",
                    "step": 1,
                    "createTimeSinceEpoch": "9999999999",  # Different value
                    "artifactType": "metric",
                }
            ]
        }
        api_client.get.return_value = response_data

        metrics = metric_ops.get_metric_history("run-123", "accuracy")

        assert len(metrics) == 1
        assert (
            metrics[0].timestamp == 1234567890
        )  # Should use timestamp, not createTimeSinceEpoch

    def test_get_metric_history_bulk_interval_from_steps(self, metric_ops, api_client):
        """Test basic bulk metric history retrieval for specific steps."""
        response_data = {
            "items": [
                {
                    "name": "accuracy",
                    "value": 0.95,
                    "timestamp": "1234567890",
                    "step": 1,
                    "createTimeSinceEpoch": "1234567890",
                    "artifactType": "metric",
                },
                {
                    "name": "accuracy",
                    "value": 0.98,
                    "timestamp": "1234567892",
                    "step": 3,
                    "createTimeSinceEpoch": "1234567892",
                    "artifactType": "metric",
                },
            ],
            "nextPageToken": "token505",
        }
        api_client.get.return_value = response_data

        result = metric_ops.get_metric_history_bulk_interval_from_steps(
            "run-123", "accuracy", steps=[1, 3], max_results=2
        )

        assert len(result) == 2
        assert result[0].value == 0.95
        assert result[0].step == 1
        assert result[0].key == "accuracy"
        assert result[1].value == 0.98
        assert result[1].step == 3
        assert result[1].key == "accuracy"
        assert result.token == "token505"

        api_client.get.assert_called_once_with(
            "/experiment_runs/run-123/metric_history",
            params={"name": "accuracy", "pageSize": 2, "stepIds": "1,3"},
        )

    def test_get_metric_history_bulk_interval_from_steps_sorts_correctly(
        self, metric_ops, api_client
    ):
        """Test that bulk metric history sorts correctly by step and timestamp."""
        response_data = {
            "items": [
                {
                    "name": "accuracy",
                    "value": 0.95,
                    "timestamp": "1234567890",
                    "step": 1,
                    "createTimeSinceEpoch": "1234567890",
                    "artifactType": "metric",
                },
                {
                    "name": "accuracy",
                    "value": 0.97,
                    "timestamp": "1234567891",
                    "step": 1,
                    "createTimeSinceEpoch": "1234567891",
                    "artifactType": "metric",
                },
                {
                    "name": "accuracy",
                    "value": 0.98,
                    "timestamp": "1234567892",
                    "step": 3,
                    "createTimeSinceEpoch": "1234567892",
                    "artifactType": "metric",
                },
            ]
        }
        api_client.get.return_value = response_data

        result = metric_ops.get_metric_history_bulk_interval_from_steps(
            "run-123", "accuracy", steps=[1, 3]
        )

        assert len(result) == 3
        assert result[0].step == 1
        assert result[0].value == 0.95
        assert result[0].key == "accuracy"
        assert result[1].step == 1
        assert result[1].value == 0.97
        assert result[1].key == "accuracy"
        assert result[2].step == 3
        assert result[2].value == 0.98
        assert result[2].key == "accuracy"

        api_client.get.assert_called_once_with(
            "/experiment_runs/run-123/metric_history",
            params={"name": "accuracy", "stepIds": "1,3"},
        )

    def test_get_metric_history_bulk_interval_from_steps_with_pagination(
        self, metric_ops, api_client
    ):
        """Test bulk metric history retrieval with pagination parameters."""
        response_data = {
            "items": [
                {
                    "name": "accuracy",
                    "value": 0.95,
                    "timestamp": "1234567890",
                    "step": 1,
                    "createTimeSinceEpoch": "1234567890",
                    "artifactType": "metric",
                }
            ],
            "nextPageToken": "token123",
        }
        api_client.get.return_value = response_data

        result = metric_ops.get_metric_history_bulk_interval_from_steps(
            "run-123", "accuracy", steps=[1], max_results=10, page_token="token123"
        )

        assert len(result) == 1
        assert result[0].step == 1
        assert result[0].value == 0.95
        assert result[0].key == "accuracy"
        assert result.token == "token123"

        api_client.get.assert_called_once_with(
            "/experiment_runs/run-123/metric_history",
            params={
                "name": "accuracy",
                "pageSize": 10,
                "pageToken": "token123",
                "stepIds": "1",
            },
        )

    def test_get_metric_history_bulk_interval_from_steps_filters_by_steps(
        self, metric_ops, api_client
    ):
        """Test that bulk metric history filters by specified steps."""
        response_data = {
            "items": [
                {
                    "name": "accuracy",
                    "value": 0.95,
                    "timestamp": "1234567890",
                    "step": 1,
                    "createTimeSinceEpoch": "1234567890",
                    "artifactType": "metric",
                },
                {
                    "name": "accuracy",
                    "value": 0.98,
                    "timestamp": "1234567892",
                    "step": 3,
                    "createTimeSinceEpoch": "1234567892",
                    "artifactType": "metric",
                },
            ]
        }
        api_client.get.return_value = response_data

        result = metric_ops.get_metric_history_bulk_interval_from_steps(
            "run-123", "accuracy", steps=[1, 3]
        )

        assert len(result) == 2
        assert result[0].step == 1
        assert result[0].value == 0.95
        assert result[0].key == "accuracy"
        assert result[1].step == 3
        assert result[1].value == 0.98
        assert result[1].key == "accuracy"

        api_client.get.assert_called_once_with(
            "/experiment_runs/run-123/metric_history",
            params={"name": "accuracy", "stepIds": "1,3"},
        )

    def test_get_metric_history_bulk_interval_from_steps_with_max_results(
        self, metric_ops, api_client
    ):
        """Test bulk metric history retrieval with max_results parameter."""
        response_data = {
            "items": [
                {
                    "name": "accuracy",
                    "value": 0.95,
                    "timestamp": "1234567890",
                    "step": 1,
                    "createTimeSinceEpoch": "1234567890",
                    "artifactType": "metric",
                },
                {
                    "name": "accuracy",
                    "value": 0.97,
                    "timestamp": "1234567891",
                    "step": 2,
                    "createTimeSinceEpoch": "1234567891",
                    "artifactType": "metric",
                },
            ]
        }
        api_client.get.return_value = response_data

        result = metric_ops.get_metric_history_bulk_interval_from_steps(
            "run-123", "accuracy", steps=[1, 2, 3], max_results=2
        )

        assert len(result) == 2  # Both metrics match the steps filter
        assert result[0].step == 1
        assert result[0].value == 0.95
        assert result[0].key == "accuracy"
        assert result[1].step == 2
        assert result[1].value == 0.97
        assert result[1].key == "accuracy"

        api_client.get.assert_called_once_with(
            "/experiment_runs/run-123/metric_history",
            params={"name": "accuracy", "pageSize": 2, "stepIds": "1,2,3"},
        )

    def test_get_metric_history_bulk_interval_from_steps_no_steps_specified(
        self, metric_ops, api_client
    ):
        """Test bulk metric history retrieval when no steps are specified."""
        response_data = {
            "items": [
                {
                    "name": "accuracy",
                    "value": 0.95,
                    "timestamp": "1234567890",
                    "step": 1,
                    "createTimeSinceEpoch": "1234567890",
                    "artifactType": "metric",
                },
                {
                    "name": "accuracy",
                    "value": 0.97,
                    "timestamp": "1234567891",
                    "step": 2,
                    "createTimeSinceEpoch": "1234567891",
                    "artifactType": "metric",
                },
            ]
        }
        api_client.get.return_value = response_data

        result = metric_ops.get_metric_history_bulk_interval_from_steps(
            "run-123", "accuracy"
        )

        assert len(result) == 2
        assert result[0].step == 1
        assert result[0].value == 0.95
        assert result[0].key == "accuracy"
        assert result[1].step == 2
        assert result[1].value == 0.97
        assert result[1].key == "accuracy"

        api_client.get.assert_called_once_with(
            "/experiment_runs/run-123/metric_history", params={"name": "accuracy"}
        )

    def test_get_metric_history_with_page_token(self, metric_ops, api_client):
        """Test getting metric history with page_token parameter."""
        response_data = {
            "items": [
                {
                    "name": "accuracy",
                    "value": 0.95,
                    "timestamp": "1234567890",
                    "step": 1,
                    "createTimeSinceEpoch": "1234567890",
                    "artifactType": "metric",
                }
            ],
            "nextPageToken": "next_token_456",
        }
        api_client.get.return_value = response_data

        result = metric_ops.get_metric_history(
            "run-123", "accuracy", page_token="token_123"
        )

        assert len(result) == 1
        assert result[0].value == 0.95
        assert result[0].step == 1
        assert result[0].key == "accuracy"
        assert result.token == "next_token_456"

        # Verify API call includes page_token
        api_client.get.assert_called_once_with(
            "/experiment_runs/run-123/metric_history",
            params={"name": "accuracy", "pageToken": "token_123"},
        )

    def test_get_metric_history_with_max_results_and_page_token(
        self, metric_ops, api_client
    ):
        """Test getting metric history with both max_results and page_token parameters."""
        response_data = {
            "items": [
                {
                    "name": "accuracy",
                    "value": 0.95,
                    "timestamp": "1234567890",
                    "step": 1,
                    "createTimeSinceEpoch": "1234567890",
                    "artifactType": "metric",
                }
            ],
            "nextPageToken": "next_token_789",
        }
        api_client.get.return_value = response_data

        result = metric_ops.get_metric_history(
            "run-123", "accuracy", max_results=10, page_token="token_456"
        )

        assert len(result) == 1
        assert result[0].value == 0.95
        assert result[0].step == 1
        assert result[0].key == "accuracy"
        assert result.token == "next_token_789"

        # Verify API call includes both parameters
        api_client.get.assert_called_once_with(
            "/experiment_runs/run-123/metric_history",
            params={"name": "accuracy", "pageSize": 10, "pageToken": "token_456"},
        )

    def test_get_metric_history_with_max_results(self, metric_ops, api_client):
        """Test getting metric history with max_results parameter."""
        response_data = {
            "items": [
                {
                    "name": "accuracy",
                    "value": 0.95,
                    "timestamp": "1234567890",
                    "step": 1,
                    "createTimeSinceEpoch": "1234567890",
                    "artifactType": "metric",
                }
            ],
            "nextPageToken": "",
        }
        api_client.get.return_value = response_data

        result = metric_ops.get_metric_history("run-123", "accuracy", max_results=5)

        assert len(result) == 1
        assert result[0].value == 0.95
        assert result[0].step == 1
        assert result[0].key == "accuracy"
        assert result.token is None  # Empty nextPageToken should become None

        # Verify API call includes max_results as pageSize
        api_client.get.assert_called_once_with(
            "/experiment_runs/run-123/metric_history",
            params={"name": "accuracy", "pageSize": 5},
        )

    def test_get_metric_history_bulk(self, metric_ops, api_client):
        """Test getting metric history for multiple runs."""
        response_data_run1 = {
            "items": [
                {
                    "name": "accuracy",
                    "value": 0.95,
                    "timestamp": "1234567890",
                    "step": 1,
                    "createTimeSinceEpoch": "1234567890",
                    "artifactType": "metric",
                }
            ]
        }
        response_data_run2 = {
            "items": [
                {
                    "name": "accuracy",
                    "value": 0.97,
                    "timestamp": "1234567891",
                    "step": 1,
                    "createTimeSinceEpoch": "1234567891",
                    "artifactType": "metric",
                }
            ]
        }
        api_client.get.side_effect = [response_data_run1, response_data_run2]

        result = metric_ops.get_metric_history_bulk(["run-123", "run-456"], "accuracy")

        assert len(result) == 2
        assert result[0].value == 0.95
        assert result[0].step == 1
        assert result[0].key == "accuracy"
        assert result[0].run_id == "run-123"
        assert result[1].value == 0.97
        assert result[1].step == 1
        assert result[1].key == "accuracy"
        assert result[1].run_id == "run-456"

        # Verify API calls
        assert api_client.get.call_count == 2
        api_client.get.assert_any_call(
            "/experiment_runs/run-123/metric_history", params={"name": "accuracy"}
        )
        api_client.get.assert_any_call(
            "/experiment_runs/run-456/metric_history", params={"name": "accuracy"}
        )

    def test_get_metric_history_bulk_with_max_results(self, metric_ops, api_client):
        """Test getting metric history for multiple runs with max_results."""
        response_data = {
            "items": [
                {
                    "name": "accuracy",
                    "value": 0.95,
                    "timestamp": "1234567890",
                    "step": 1,
                    "createTimeSinceEpoch": "1234567890",
                    "artifactType": "metric",
                }
            ]
        }
        api_client.get.return_value = response_data

        result = metric_ops.get_metric_history_bulk(
            ["run-123"], "accuracy", max_results=10
        )

        assert len(result) == 1
        assert result[0].value == 0.95
        assert result[0].run_id == "run-123"

        # Verify API call includes max_results
        api_client.get.assert_called_once_with(
            "/experiment_runs/run-123/metric_history",
            params={"name": "accuracy", "pageSize": 10},
        )

    def test_get_metric_history_bulk_empty_runs(self, metric_ops, api_client):
        """Test getting metric history for empty run list."""
        result = metric_ops.get_metric_history_bulk([], "accuracy")

        assert len(result) == 0
        api_client.get.assert_not_called()
