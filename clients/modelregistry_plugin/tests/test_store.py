"""
Tests for ModelRegistryStore
"""

import pytest
from unittest.mock import Mock, patch, MagicMock
import requests
import json
import uuid

from modelregistry_plugin.store import ModelRegistryStore
from mlflow.entities import (
    Experiment, Run, RunInfo, RunData, RunStatus, ExperimentTag, RunTag, Param, Metric,
    ViewType, LifecycleStage, DatasetInput, LoggedModelInput, LoggedModelOutput,
    LoggedModel, LoggedModelTag, LoggedModelParameter, LoggedModelStatus
)
from mlflow.store.entities.paged_list import PagedList
from mlflow.models import Model
from mlflow.exceptions import MlflowException


class TestModelRegistryStore:
    
    @pytest.fixture
    def store(self):
        """Create a ModelRegistryStore instance for testing."""
        return ModelRegistryStore("modelregistry://localhost:8080", "s3://bucket/artifacts")
    
    @pytest.fixture
    def mock_response(self):
        """Create a mock response object."""
        response = Mock(spec=requests.Response)
        response.ok = True
        response.json.return_value = {}
        return response
    
    @pytest.fixture
    def mock_model(self):
        """Create a mock MLflow Model object."""
        model = Mock(spec=Model)
        model.to_dict.return_value = {
            "artifact_path": "model",
            "run_id": "run-123",
            "flavors": {"python_function": {}},
            "model_uuid": "uuid-123",
            "utc_time_created": "2023-01-01T00:00:00Z",
            "mlflow_version": "2.0.0"
        }
        model_info = Mock()
        model_info.model_uri = "runs:/run-123/model"
        model_info.artifact_path = "model"
        model_info.model_uuid = "uuid-123"
        model_info.utc_time_created = "2023-01-01T00:00:00Z"
        model_info.mlflow_version = "2.0.0"
        model_info.flavors = {"python_function": {}}
        model.get_model_info.return_value = model_info
        return model

    # Initialization tests
    def test_init(self, store):
        """Test store initialization."""
        assert store.host == "localhost"
        assert store.port == 8080
        assert store.secure is False
        assert store.base_url == "http://localhost:8080/api/model_registry/v1alpha3"
    
    def test_init_secure(self):
        """Test store initialization with secure connection."""
        with patch.dict('os.environ', {'MODEL_REGISTRY_SECURE': 'true'}):
            store = ModelRegistryStore("modelregistry://localhost:8080")
            assert store.secure is True
    
    def test_init_with_artifact_uri(self):
        """Test store initialization with artifact URI."""
        store = ModelRegistryStore("modelregistry://localhost:8080", "s3://bucket/artifacts")
        assert store.artifact_uri == "s3://bucket/artifacts"
    
    def test_init_from_env(self):
        """Test store initialization from environment variable."""
        with patch.dict('os.environ', {'MLFLOW_TRACKING_URI': 'modelregistry://test:9090'}):
            store = ModelRegistryStore()
            assert store.host == "test"
            assert store.port == 9090

    # Request method tests
    @patch('modelregistry_plugin.store.requests.request')
    def test_request_success(self, mock_request, store, mock_response):
        """Test successful API request."""
        mock_request.return_value = mock_response
        
        response_data = store._request("GET", "/test")
        
        mock_request.assert_called_once()
        assert response_data == mock_response.json.return_value
    
    @patch('modelregistry_plugin.store.requests.request')
    def test_request_failure(self, mock_request, store):
        """Test failed API request."""
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = False
        mock_response.json.return_value = {"message": "Test error"}
        mock_response.status_code = 404
        mock_request.return_value = mock_response
        
        with pytest.raises(MlflowException) as exc_info:
            store._request("GET", "/test")

        assert "Model Registry API error: Test error" in str(exc_info.value)

    
    @patch('modelregistry_plugin.store.requests.request')
    def test_request_with_custom_properties(self, mock_request, store):
        """Test request with custom properties conversion."""
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {
            "customProperties": {
                "key1": {"string_value": "value1", "metadataType": "MetadataStringValue"}
            }
        }
        mock_request.return_value = mock_response
        
        response_data = store._request("POST", "/test", json={
            "customProperties": {"key1": "value1"}
        })
        
        # Check that custom properties were converted in request and response
        call_args = mock_request.call_args
        assert call_args[1]['json']['customProperties']['key1']['string_value'] == "value1"
        assert call_args[1]['json']['customProperties']['key1']['metadataType'] == "MetadataStringValue"
        assert response_data['customProperties']['key1'] == "value1"

    # Experiment operations tests
    @patch('modelregistry_plugin.store.requests.request')
    def test_create_experiment(self, mock_request, store):
        """Test experiment creation."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {"id": "exp-123"}
        mock_request.return_value = mock_response
        
        experiment_id = store.create_experiment("test-experiment")
        
        assert experiment_id == "exp-123"
        # Should make 2 calls: POST to create, then PATCH to set default artifact location
        assert mock_request.call_count == 2
        
        # Check first call (POST to create experiment)
        call_args = mock_request.call_args_list[0]
        assert call_args[0][0] == "POST"  # method
        assert "/experiments" in call_args[0][1]  # endpoint
        json_data = call_args[1]['json']
        assert json_data["name"] == "test-experiment"
        assert json_data["description"] == "MLflow experiment: test-experiment"
        assert json_data["state"] == "LIVE"
        assert json_data["customProperties"] == {}
        # externalId should not be in the initial payload when artifact_uri is available
        assert "externalId" not in json_data
        
        # Check second call (PATCH to set default artifact location)
        call_args = mock_request.call_args_list[1]
        assert call_args[0][0] == "PATCH"  # method
        assert "/experiments/exp-123" in call_args[0][1]  # endpoint
        json_data = call_args[1]['json']
        assert json_data["externalId"] == "s3://bucket/artifacts/experiments/exp-123"
    
    @patch('modelregistry_plugin.store.requests.request')
    def test_create_experiment_with_artifact_location(self, mock_request, store):
        """Test experiment creation with explicit artifact_location."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {"id": "exp-123"}
        mock_request.return_value = mock_response
        
        experiment_id = store.create_experiment("test-experiment", artifact_location="s3://custom/location")
        
        assert experiment_id == "exp-123"
        # Should make only 1 call since artifact_location was provided
        assert mock_request.call_count == 1
        
        call_args = mock_request.call_args
        assert call_args[0][0] == "POST"  # method
        assert "/experiments" in call_args[0][1]  # endpoint
        json_data = call_args[1]['json']
        assert json_data["name"] == "test-experiment"
        assert json_data["externalId"] == "s3://custom/location"
    
    @patch('modelregistry_plugin.store.requests.request')
    def test_create_experiment_without_artifact_uri(self, mock_request):
        """Test experiment creation when store has no artifact_uri."""
        # Create store without artifact_uri
        store = ModelRegistryStore("modelregistry://localhost:8080")
        
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {"id": "exp-123"}
        mock_request.return_value = mock_response
        
        experiment_id = store.create_experiment("test-experiment")
        
        assert experiment_id == "exp-123"
        # Should make only 1 call since no artifact_uri available
        assert mock_request.call_count == 1
        
        call_args = mock_request.call_args
        assert call_args[0][0] == "POST"  # method
        assert "/experiments" in call_args[0][1]  # endpoint
        json_data = call_args[1]['json']
        assert json_data["name"] == "test-experiment"
        assert json_data.get("externalId") is None
    
    @patch('modelregistry_plugin.store.requests.request')
    def test_create_experiment_with_tags(self, mock_request, store):
        """Test experiment creation with tags."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {"id": "exp-123"}
        mock_request.return_value = mock_response
        
        tags = [ExperimentTag("key1", "value1"), ExperimentTag("key2", "value2")]
        
        experiment_id = store.create_experiment("test-experiment", tags=tags)
        
        mock_request.assert_called_once()
        call_args = mock_request.call_args
        json_data = call_args[1]['json']
        custom_props = json_data['customProperties']
        assert custom_props['key1']['string_value'] == "value1"
        assert custom_props['key2']['string_value'] == "value2"
    
    @patch('modelregistry_plugin.store.requests.request')
    def test_get_experiment(self, mock_request, store):
        """Test getting experiment by ID."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {
            "id": "exp-123",
            "name": "test-experiment",
            "externalId": "s3://bucket/artifacts/experiments/exp-123",
            "customProperties": {
                "key1": {"string_value": "value1", "metadataType": "MetadataStringValue"}
            }
        }
        mock_request.return_value = mock_response
        
        experiment = store.get_experiment("exp-123")
        
        assert isinstance(experiment, Experiment)
        assert experiment.experiment_id == "exp-123"
        assert experiment.name == "test-experiment"
        assert experiment.artifact_location == "s3://bucket/artifacts/experiments/exp-123"
        assert len(experiment.tags) == 1
        assert experiment.tags["key1"] == "value1"
    
    @patch('modelregistry_plugin.store.requests.request')
    def test_get_experiment_by_name(self, mock_request, store):
        """Test getting experiment by name."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {
            "id": "exp-123",
            "name": "test-experiment",
            "externalId": "s3://bucket/artifacts/experiments/exp-123",
            "customProperties": {
                "key1": {"string_value": "value1", "metadataType": "MetadataStringValue"}
            }
        }
        mock_request.return_value = mock_response
        
        experiment = store.get_experiment_by_name("test-experiment")

        mock_request.assert_called_once()
        call_args = mock_request.call_args
        assert call_args[0][0] == "GET"  # method
        assert "/experiment" in call_args[0][1]  # endpoint
        params = call_args[1]['params']
        assert params["name"] == "test-experiment"

        assert isinstance(experiment, Experiment)
        assert experiment.experiment_id == "exp-123"
        assert experiment.name == "test-experiment"
        assert experiment.artifact_location == "s3://bucket/artifacts/experiments/exp-123"
        assert len(experiment.tags) == 1
        assert experiment.tags["key1"] == "value1"
    
    @patch('modelregistry_plugin.store.requests.request')
    def test_get_experiment_by_name_not_found(self, mock_request, store):
        """Test getting experiment by name when not found."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {"experiments": []}
        mock_request.return_value = mock_response
        
        experiment = store.get_experiment_by_name("nonexistent")
        
        assert experiment is None
    
    @patch('modelregistry_plugin.store.requests.request')
    def test_delete_experiment(self, mock_request, store):
        """Test deleting an experiment."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {}
        mock_request.return_value = mock_response
        
        store.delete_experiment("exp-123")
        
        mock_request.assert_called_once()
        call_args = mock_request.call_args
        assert call_args[0][0] == "PATCH"  # method
        assert "/experiments/exp-123" in call_args[0][1]  # endpoint
        json_data = call_args[1]['json']
        assert json_data["state"] == "ARCHIVED"
    
    @patch('modelregistry_plugin.store.requests.request')
    def test_restore_experiment(self, mock_request, store):
        """Test restoring an experiment."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {}
        mock_request.return_value = mock_response
        
        store.restore_experiment("exp-123")
        
        mock_request.assert_called_once()
        call_args = mock_request.call_args
        assert call_args[0][0] == "PATCH"  # method
        assert "/experiments/exp-123" in call_args[0][1]  # endpoint
        json_data = call_args[1]['json']
        assert json_data["state"] == "LIVE"
    
    @patch('modelregistry_plugin.store.requests.request')
    def test_rename_experiment(self, mock_request, store):
        """Test renaming an experiment."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {}
        mock_request.return_value = mock_response
        
        store.rename_experiment("exp-123", "new-name")
        
        mock_request.assert_called_once()
        call_args = mock_request.call_args
        assert call_args[0][0] == "PATCH"  # method
        assert "/experiments/exp-123" in call_args[0][1]  # endpoint
        json_data = call_args[1]['json']
        assert json_data["name"] == "new-name"
    
    @patch('modelregistry_plugin.store.requests.request')
    def test_list_experiments(self, mock_request, store):
        """Test listing experiments."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {
            "items": [
                {
                    "id": "1", 
                    "name": "exp1", 
                    "state": "LIVE",
                    "externalId": "s3://bucket/artifacts/experiments/1",
                    "customProperties": {
                        "tag1": {"string_value": "value1", "metadataType": "MetadataStringValue"}
                    }
                },
                {
                    "id": "2", 
                    "name": "exp2", 
                    "state": "ARCHIVED",
                    "externalId": "s3://bucket/artifacts/experiments/2",
                    "customProperties": {
                        "tag2": {"string_value": "value2", "metadataType": "MetadataStringValue"}
                    }
                }
            ]
        }
        mock_request.return_value = mock_response
        
        experiments = store.list_experiments(view_type=ViewType.ALL)
        
        assert len(experiments) == 2
        assert experiments[0].experiment_id == "1"
        assert experiments[0].name == "exp1"
        assert experiments[0].artifact_location == "s3://bucket/artifacts/experiments/1"
        assert len(experiments[0].tags) == 1
        assert experiments[0].tags["tag1"] == "value1"
        assert experiments[1].experiment_id == "2"
        assert experiments[1].name == "exp2"
        assert experiments[1].artifact_location == "s3://bucket/artifacts/experiments/2"
        assert len(experiments[1].tags) == 1
        assert experiments[1].tags["tag2"] == "value2"
    
    @patch('modelregistry_plugin.store.requests.request')
    def test_list_experiments_with_pagination(self, mock_request, store):
        """Test listing experiments with pagination."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {"items": []}
        mock_request.return_value = mock_response
        
        store.list_experiments(max_results=10, page_token="token123")
        
        mock_request.assert_called_once()
        call_args = mock_request.call_args
        params = call_args[1]['params']
        assert params['pageSize'] == 10
        assert params['pageToken'] == "token123"

    @patch('modelregistry_plugin.store.requests.request')
    def test_search_experiments(self, mock_request, store):
        """Test searching experiments."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {
            "items": [
                {
                    "id": "1", 
                    "name": "exp1", 
                    "state": "LIVE",
                    "externalId": "s3://bucket/artifacts/experiments/1",
                    "customProperties": {
                        "tag1": {"string_value": "value1", "metadataType": "MetadataStringValue"}
                    }
                },
                {
                    "id": "2", 
                    "name": "exp2", 
                    "state": "ARCHIVED",
                    "externalId": "s3://bucket/artifacts/experiments/2",
                    "customProperties": {
                        "tag2": {"string_value": "value2", "metadataType": "MetadataStringValue"}
                    }
                }
            ],
            "nextPageToken": "token123"
        }
        mock_request.return_value = mock_response
        
        result = store.search_experiments(
            view_type=ViewType.ALL,
            max_results=10,
            filter_string="name='exp1'",
            order_by=["name"],
            page_token="token123"
        )
        
        assert isinstance(result, PagedList)
        assert len(result) == 2
        assert result[0].experiment_id == "1"
        assert result[0].name == "exp1"
        assert result[0].artifact_location == "s3://bucket/artifacts/experiments/1"
        assert len(result[0].tags) == 1
        assert result[0].tags["tag1"] == "value1"
        assert result[1].experiment_id == "2"
        assert result[1].name == "exp2"
        assert result[1].artifact_location == "s3://bucket/artifacts/experiments/2"
        assert len(result[1].tags) == 1
        assert result[1].tags["tag2"] == "value2"
        assert result.token == "token123"
        
        # Verify API call parameters
        mock_request.assert_called_once()
        call_args = mock_request.call_args
        assert call_args[0][0] == "GET"  # method
        assert "/experiments" in call_args[0][1]  # endpoint
        params = call_args[1]['params']
        assert params['pageSize'] == 10
        assert params['pageToken'] == "token123"

    @patch('modelregistry_plugin.store.requests.request')
    def test_search_experiments_active_only(self, mock_request, store):
        """Test searching experiments with ACTIVE_ONLY view type."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {
            "items": [
                {
                    "id": "1", 
                    "name": "exp1", 
                    "state": "LIVE",
                    "externalId": "s3://bucket/artifacts/experiments/1",
                    "customProperties": {}
                },
                {
                    "id": "2", 
                    "name": "exp2", 
                    "state": "ARCHIVED",
                    "externalId": "s3://bucket/artifacts/experiments/2",
                    "customProperties": {}
                }
            ]
        }
        mock_request.return_value = mock_response
        
        result = store.search_experiments(view_type=ViewType.ACTIVE_ONLY)
        
        # Should only return active experiments
        assert len(result) == 1
        assert result[0].experiment_id == "1"
        assert result[0].name == "exp1"
        assert result[0].artifact_location == "s3://bucket/artifacts/experiments/1"

    @patch('modelregistry_plugin.store.requests.request')
    def test_search_experiments_deleted_only(self, mock_request, store):
        """Test searching experiments with DELETED_ONLY view type."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {
            "items": [
                {
                    "id": "1", 
                    "name": "exp1", 
                    "state": "LIVE",
                    "externalId": "s3://bucket/artifacts/experiments/1",
                    "customProperties": {}
                },
                {
                    "id": "2", 
                    "name": "exp2", 
                    "state": "ARCHIVED",
                    "externalId": "s3://bucket/artifacts/experiments/2",
                    "customProperties": {}
                }
            ]
        }
        mock_request.return_value = mock_response
        
        result = store.search_experiments(view_type=ViewType.DELETED_ONLY)
        
        # Should only return deleted experiments
        assert len(result) == 1
        assert result[0].experiment_id == "2"
        assert result[0].name == "exp2"
        assert result[0].artifact_location == "s3://bucket/artifacts/experiments/2"

    @patch('modelregistry_plugin.store.requests.request')
    def test_search_experiments_with_filter_string(self, mock_request, store):
        """Test searching experiments with filter string (should be ignored for now)."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {"items": []}
        mock_request.return_value = mock_response
        
        # This should not raise an error even though filter_string is not supported yet
        result = store.search_experiments(filter_string="name='test'")
        
        assert isinstance(result, PagedList)
        assert len(result) == 0

    @patch('modelregistry_plugin.store.requests.request')
    def test_search_experiments_with_order_by(self, mock_request, store):
        """Test searching experiments with order_by (should be ignored for now)."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {"items": []}
        mock_request.return_value = mock_response
        
        # This should not raise an error even though order_by is not supported yet
        result = store.search_experiments(order_by=["name"])
        
        assert isinstance(result, PagedList)
        assert len(result) == 0

    # Run operations tests
    @patch('modelregistry_plugin.store.requests.request')
    def test_create_run(self, mock_request, store):
        """Test run creation."""
        # Mock the raw response from Model Registry API for creating run
        mock_response_create = Mock(spec=requests.Response)
        mock_response_create.ok = True
        mock_response_create.json.return_value = {
            "id": "run-123",
            "experimentId": "exp-123",
            "createTimeSinceEpoch": "1234567890"
        }
        
        # Mock the raw response from Model Registry API for getting experiment
        mock_response_experiment = Mock(spec=requests.Response)
        mock_response_experiment.ok = True
        mock_response_experiment.json.return_value = {
            "id": "exp-123",
            "name": "test-experiment",
            "externalId": "s3://bucket/artifacts/experiments/exp-123",
            "customProperties": {}
        }
        
        # Mock the raw response from Model Registry API for updating run with artifact location
        mock_response_update = Mock(spec=requests.Response)
        mock_response_update.ok = True
        mock_response_update.json.return_value = {}
        
        mock_request.side_effect = [mock_response_create, mock_response_experiment, mock_response_update]
        
        run = store.create_run("exp-123", start_time=1234567890)
        
        assert isinstance(run, Run)
        assert run.info.run_id == "run-123"
        assert run.info.experiment_id == "exp-123"
        assert run.info.status == RunStatus.RUNNING
        
        # Should make 3 calls: POST to create run, GET to get experiment, then PATCH to set artifact location
        assert mock_request.call_count == 3
        
        # Check first call (POST to create run)
        call_args = mock_request.call_args_list[0]
        assert call_args[0][0] == "POST"  # method
        assert "/experiment_runs" in call_args[0][1]  # endpoint
        json_data = call_args[1]['json']
        assert json_data["experimentId"] == "exp-123"
        assert json_data["status"] == "RUNNING"
        
        # Check second call (GET experiment)
        call_args = mock_request.call_args_list[1]
        assert call_args[0][0] == "GET"  # method
        assert "/experiments/exp-123" in call_args[0][1]  # endpoint
        
        # Check third call (PATCH to set artifact location)
        call_args = mock_request.call_args_list[2]
        assert call_args[0][0] == "PATCH"  # method
        assert "/experiment_runs/run-123" in call_args[0][1]  # endpoint
        json_data = call_args[1]['json']
        assert json_data["externalId"] == "s3://bucket/artifacts/experiments/exp-123/run-123"
    
    @patch('modelregistry_plugin.store.requests.request')
    def test_create_run_with_user_and_tags(self, mock_request, store):
        """Test run creation with user ID and tags."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {
            "id": "run-123",
            "experimentId": "exp-123",
            "createTimeSinceEpoch": "1234567890"
        }
        mock_request.return_value = mock_response
        
        tags = [RunTag("key1", "value1"), RunTag("key2", "value2")]
        
        run = store.create_run("exp-123", user_id="user123", tags=tags)
        
        mock_request.assert_called_once()
        call_args = mock_request.call_args
        json_data = call_args[1]['json']
        custom_props = json_data['customProperties']
        assert custom_props['key1']['string_value'] == "value1"
        assert custom_props['key2']['string_value'] == "value2"
    
    @patch('modelregistry_plugin.store.requests.request')
    def test_get_run(self, mock_request, store):
        """Test getting run by ID."""
        # First call: get run data
        mock_response_run = Mock(spec=requests.Response)
        mock_response_run.ok = True
        mock_response_run.json.return_value = {
            "id": "run-123",
            "experimentId": "exp-123",
            "name": "test-run",
            "state": "RUNNING",
            "owner": "user123",
            "startTimeSinceEpoch": "1234567890",
            "externalId": "s3://bucket/artifacts/experiments/exp-123/run-123",
            "customProperties": {
                "key1": {"string_value": "value1", "metadataType": "MetadataStringValue"}
            }
        }
        # Second call: get metrics (empty)
        mock_response_metrics = Mock(spec=requests.Response)
        mock_response_metrics.ok = True
        mock_response_metrics.json.return_value = {"items": []}
        # Third call: get params (empty)
        mock_response_params = Mock(spec=requests.Response)
        mock_response_params.ok = True
        mock_response_params.json.return_value = {"items": []}

        mock_request.side_effect = [mock_response_run, mock_response_metrics, mock_response_params]

        run = store.get_run("run-123")

        assert isinstance(run, Run)
        assert run.info.run_id == "run-123"
        assert run.info.experiment_id == "exp-123"
        assert run.info.run_name == "test-run"
        assert run.info.user_id == "user123"
        assert run.info.artifact_uri == "s3://bucket/artifacts/experiments/exp-123/run-123"
        assert len(run.data.tags) == 1
        assert run.data.tags["key1"] == "value1"
    
    @patch('modelregistry_plugin.store.requests.request')
    def test_update_run_info(self, mock_request, store):
        """Test updating run info."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {
            "id": "run-123",
            "experimentId": "exp-123",
            "status": "FINISHED",
            "endTimeSinceEpoch": "1234567899",
            "startTimeSinceEpoch": "1234567890",
            "owner": "user123",
            "externalId": "s3://bucket/artifacts/experiments/exp-123/run-123"
        }
        mock_request.return_value = mock_response
        
        run_info = store.update_run_info("run-123", RunStatus.FINISHED, end_time=1234567899)
        
        assert isinstance(run_info, RunInfo)
        assert run_info.run_id == "run-123"
        assert run_info.status == RunStatus.FINISHED
        assert run_info.artifact_uri == "s3://bucket/artifacts/experiments/exp-123/run-123"
    
    @patch('modelregistry_plugin.store.requests.request')
    def test_delete_run(self, mock_request, store):
        """Test deleting a run."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {}
        mock_request.return_value = mock_response
        
        store.delete_run("run-123")
        
        mock_request.assert_called_once()
        call_args = mock_request.call_args
        assert call_args[0][0] == "PATCH"  # method
        assert "/experiment_runs/run-123" in call_args[0][1]  # endpoint
        json_data = call_args[1]['json']
        assert json_data["state"] == "ARCHIVED"
    
    @patch('modelregistry_plugin.store.requests.request')
    def test_restore_run(self, mock_request, store):
        """Test restoring a run."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {}
        mock_request.return_value = mock_response
        
        store.restore_run("run-123")
        
        mock_request.assert_called_once()
        call_args = mock_request.call_args
        assert call_args[0][0] == "PATCH"  # method
        assert "/experiment_runs/run-123" in call_args[0][1]  # endpoint
        json_data = call_args[1]['json']
        assert json_data["state"] == "LIVE"

    # Metric and parameter tests
    @patch('modelregistry_plugin.store.requests.request')
    def test_log_metric(self, mock_request, store):
        """Test logging a metric."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = { 
            "id": "metric-123", 
            "name": "accuracy", 
            "value": 0.95, 
            "step": 1, 
            "timestamp": 1234567890, 
            "createTimeSinceEpoch": 1234567890, 
            "lastModifiedTimeSinceEpoch": 1234567890 
        }
        mock_request.return_value = mock_response
        
        metric = Metric("accuracy", 0.95, 1234567890, 1)
        
        store.log_metric("run-123", metric)
        
        mock_request.assert_called_once()
        call_args = mock_request.call_args
        assert call_args[0][0] == "POST"  # method
        assert "/experiment_runs/run-123/artifacts" in call_args[0][1]  # endpoint
        json_data = call_args[1]['json']
        assert json_data["artifactType"] == "metric"
        assert json_data["name"] == "accuracy"
        assert json_data["value"] == 0.95
        assert json_data["step"] == 1
        assert json_data["timestamp"] == "1234567890"
        assert json_data["customProperties"] == {}
    
    @patch('modelregistry_plugin.store.requests.request')
    def test_get_run_metrics(self, mock_request, store):
        """Test getting run metrics."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {
            "items": [
                {
                    "name": "accuracy",
                    "value": 0.95,
                    "timestamp": "1234567890",
                    "step": 1,
                    "createTimeSinceEpoch": "1234567890"
                }
            ]
        }
        mock_request.return_value = mock_response
        
        metrics = store._get_run_metrics("run-123")
        
        assert len(metrics) == 1
        assert metrics[0].key == "accuracy"
        assert metrics[0].value == 0.95
        assert metrics[0].step == 1

    @patch('modelregistry_plugin.store.requests.request')
    def test_get_metric_history(self, mock_request, store):
        """Test getting metric history for a specific metric key."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {
            "items": [
                {
                    "name": "accuracy",
                    "value": 0.95,
                    "timestamp": "1234567890",
                    "step": 1,
                    "createTimeSinceEpoch": "1234567890"
                },
                {
                    "name": "accuracy",
                    "value": 0.97,
                    "timestamp": "1234567891",
                    "step": 2,
                    "createTimeSinceEpoch": "1234567891"
                },
                {
                    "name": "loss",
                    "value": 0.1,
                    "timestamp": "1234567892",
                    "step": 1,
                    "createTimeSinceEpoch": "1234567892"
                }
            ]
        }
        mock_request.return_value = mock_response
        
        # Get metric history for "accuracy"
        metrics = store.get_metric_history("run-123", "accuracy")
        
        assert len(metrics) == 2
        assert all(metric.key == "accuracy" for metric in metrics)
        assert metrics[0].value == 0.95
        assert metrics[0].step == 1
        assert metrics[1].value == 0.97
        assert metrics[1].step == 2
        
        # Verify metrics are sorted by timestamp and step
        assert metrics[0].timestamp <= metrics[1].timestamp

    @patch('modelregistry_plugin.store.requests.request')
    def test_get_metric_history_empty(self, mock_request, store):
        """Test getting metric history when no metrics exist."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {"items": []}
        mock_request.return_value = mock_response
        
        metrics = store.get_metric_history("run-123", "nonexistent")
        
        assert len(metrics) == 0

    @patch('modelregistry_plugin.store.requests.request')
    def test_get_metric_history_with_max_results(self, mock_request, store):
        """Test getting metric history with max_results limit."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {
            "items": [
                {
                    "name": "accuracy",
                    "value": 0.95,
                    "timestamp": "1234567890",
                    "step": 1,
                    "createTimeSinceEpoch": "1234567890"
                },
                {
                    "name": "accuracy",
                    "value": 0.97,
                    "timestamp": "1234567891",
                    "step": 2,
                    "createTimeSinceEpoch": "1234567891"
                },
                {
                    "name": "accuracy",
                    "value": 0.98,
                    "timestamp": "1234567892",
                    "step": 3,
                    "createTimeSinceEpoch": "1234567892"
                }
            ]
        }
        mock_request.return_value = mock_response
        
        # Get metric history with max_results=2
        metrics = store.get_metric_history("run-123", "accuracy", max_results=2)
        
        assert len(metrics) == 2
        assert all(metric.key == "accuracy" for metric in metrics)
        # Should return the first 2 metrics (sorted by timestamp and step)
        assert metrics[0].step == 1
        assert metrics[1].step == 2

    @patch('modelregistry_plugin.store.requests.request')
    def test_get_metric_history_with_page_token(self, mock_request, store):
        """Test getting metric history with page_token (should be ignored for now)."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {
            "items": [
                {
                    "name": "accuracy",
                    "value": 0.95,
                    "timestamp": "1234567890",
                    "step": 1,
                    "createTimeSinceEpoch": "1234567890"
                }
            ]
        }
        mock_request.return_value = mock_response
        
        # This should not raise an error even though page_token is not fully implemented yet
        metrics = store.get_metric_history("run-123", "accuracy", page_token="token123")
        
        assert len(metrics) == 1
        assert metrics[0].key == "accuracy"

    @patch('modelregistry_plugin.store.requests.request')
    def test_get_metric_history_sorting(self, mock_request, store):
        """Test that metric history is properly sorted by timestamp and step."""
        # Mock the raw response from Model Registry API with unsorted data
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {
            "items": [
                {
                    "name": "accuracy",
                    "value": 0.97,
                    "timestamp": "1234567891",
                    "step": 2,
                    "createTimeSinceEpoch": "1234567891"
                },
                {
                    "name": "accuracy",
                    "value": 0.95,
                    "timestamp": "1234567890",
                    "step": 1,
                    "createTimeSinceEpoch": "1234567890"
                },
                {
                    "name": "accuracy",
                    "value": 0.98,
                    "timestamp": "1234567890",
                    "step": 3,
                    "createTimeSinceEpoch": "1234567890"
                }
            ]
        }
        mock_request.return_value = mock_response
        
        metrics = store.get_metric_history("run-123", "accuracy")
        
        assert len(metrics) == 3
        # Should be sorted by timestamp first, then step
        assert metrics[0].timestamp == 1234567890 and metrics[0].step == 1
        assert metrics[1].timestamp == 1234567890 and metrics[1].step == 3
        assert metrics[2].timestamp == 1234567891 and metrics[2].step == 2

    @patch('modelregistry_plugin.store.requests.request')
    def test_log_param(self, mock_request, store):
        """Test logging a parameter."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {
            "id": "param-123", 
            "name": "learning_rate", 
            "value": "0.01", 
            "parameterType": "string", 
            "createTimeSinceEpoch": 1234567890, 
            "lastModifiedTimeSinceEpoch": 1234567890 
        }
        mock_request.return_value = mock_response
        
        param = Param("learning_rate", "0.01")
        
        store.log_param("run-123", param)
        
        mock_request.assert_called_once()
        call_args = mock_request.call_args
        assert call_args[0][0] == "POST"  # method
        assert "/experiment_runs/run-123/artifacts" in call_args[0][1]  # endpoint
        json_data = call_args[1]['json']
        assert json_data["artifactType"] == "parameter"
        assert json_data["name"] == "learning_rate"
        assert json_data["value"] == "0.01"
        assert json_data["parameterType"] == "string"
    
    @patch('modelregistry_plugin.store.requests.request')
    def test_get_run_params(self, mock_request, store):
        """Test getting run parameters."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {
            "items": [
                {
                    "name": "learning_rate",
                    "value": "0.01",
                    "parameterType": "string"
                }
            ]
        }
        mock_request.return_value = mock_response
        
        params = store._get_run_params("run-123")
        
        assert len(params) == 1
        assert params[0].key == "learning_rate"
        assert params[0].value == "0.01"

    # Tag management tests
    @patch('modelregistry_plugin.store.requests.request')
    def test_set_experiment_tag(self, mock_request, store):
        """Test setting experiment tag."""
        # Mock the raw response from Model Registry API for GET request
        mock_response_get = Mock(spec=requests.Response)
        mock_response_get.ok = True
        mock_response_get.json.return_value = {
            "id": "exp-123",
            "name": "test-experiment",
            "customProperties": {
                "existing": {"string_value": "value", "metadataType": "MetadataStringValue"}
            }
        }
        # Mock the raw response from Model Registry API for PATCH request
        mock_response_patch = Mock(spec=requests.Response)
        mock_response_patch.ok = True
        mock_response_patch.json.return_value = {}
        
        mock_request.side_effect = [mock_response_get, mock_response_patch]
        
        tag = ExperimentTag("key1", "value1")
        
        store.set_experiment_tag("exp-123", tag)
        
        # Should call PATCH to update tags
        call_args = mock_request.call_args_list
        assert len(call_args) == 2  # GET + PATCH
        patch_call = call_args[1]
        json_data = patch_call[1]['json']
        custom_props = json_data['customProperties']
        assert custom_props['key1']['string_value'] == "value1"
    
    @patch('modelregistry_plugin.store.requests.request')
    def test_set_tag(self, mock_request, store):
        """Test setting run tag."""
        # Mock the raw response from Model Registry API for GET run request
        mock_response_get_run = Mock(spec=requests.Response)
        mock_response_get_run.ok = True
        mock_response_get_run.json.return_value = {
            "id": "run-123",
            "experimentId": "exp-123",
            "name": "test-run",
            "state": "RUNNING",
            "owner": "user123",
            "startTimeSinceEpoch": "1234567890",
            "customProperties": {
                "existing": {"string_value": "value", "metadataType": "MetadataStringValue"}
            }
        }
        # Mock the raw response from Model Registry API for PATCH request
        mock_response_patch = Mock(spec=requests.Response)
        mock_response_patch.ok = True
        mock_response_patch.json.return_value = {
            "id": "run-123",
            "experimentId": "exp-123",
            "name": "test-run",
            "state": "RUNNING",
            "owner": "user123",
            "startTimeSinceEpoch": "1234567890",
            "customProperties": {
                "existing": {"string_value": "value", "metadataType": "MetadataStringValue"},
                "key1": {"string_value": "value1", "metadataType": "MetadataStringValue"}
            }
        }
        
        mock_request.side_effect = [
            mock_response_get_run,  # GET run
            mock_response_patch     # PATCH
        ]
        
        tag = RunTag("key1", "value1")
        
        store.set_tag("run-123", tag)
        
        # Should call PATCH to update tags
        call_args = mock_request.call_args_list
        assert len(call_args) == 2  # GET run + PATCH
        patch_call = call_args[1]
        json_data = patch_call[1]['json']
        custom_props = json_data['customProperties']
        assert custom_props['key1']['string_value'] == "value1"

    @patch('modelregistry_plugin.store.requests.request')
    def test_delete_tag(self, mock_request, store):
        """Test deleting a run tag."""
        # Mock the raw response from Model Registry API for GET run request
        mock_response_get_run = Mock(spec=requests.Response)
        mock_response_get_run.ok = True
        mock_response_get_run.json.return_value = {
            "id": "run-123",
            "experimentId": "exp-123",
            "name": "test-run",
            "state": "RUNNING",
            "owner": "user123",
            "startTimeSinceEpoch": "1234567890",
            "customProperties": {
                "key1": {"string_value": "value1", "metadataType": "MetadataStringValue"},
                "key2": {"string_value": "value2", "metadataType": "MetadataStringValue"}
            }
        }
        # Mock the raw response from Model Registry API for PATCH request
        mock_response_patch = Mock(spec=requests.Response)
        mock_response_patch.ok = True
        mock_response_patch.json.return_value = {
            "id": "run-123",
            "experimentId": "exp-123",
            "name": "test-run",
            "state": "RUNNING",
            "owner": "user123",
            "startTimeSinceEpoch": "1234567890",
            "customProperties": {
                "key2": {"string_value": "value2", "metadataType": "MetadataStringValue"}
            }
        }
        
        mock_request.side_effect = [
            mock_response_get_run,  # GET run
            mock_response_patch     # PATCH
        ]
        
        store.delete_tag("run-123", "key1")
        
        # Should call PATCH to update tags without key1
        call_args = mock_request.call_args_list
        assert len(call_args) == 2  # GET run + GET metrics + GET params + PATCH
        patch_call = call_args[1]
        json_data = patch_call[1]['json']
        custom_props = json_data['customProperties']
        assert "key1" not in custom_props
        assert "key2" in custom_props

    # Batch operations tests
    @patch('modelregistry_plugin.store.requests.request')
    def test_log_batch(self, mock_request, store):
        """Test batch logging."""
        # Mock the raw response from Model Registry API for GET run request
        mock_response_get_run = Mock(spec=requests.Response)
        mock_response_get_run.ok = True
        mock_response_get_run.json.return_value = {
            "id": "run-123",
            "experimentId": "exp-123",
            "name": "test-run",
            "state": "RUNNING",
            "owner": "user123",
            "startTimeSinceEpoch": "1234567890",
            "customProperties": {
                "existing": {"string_value": "tag", "metadataType": "MetadataStringValue"}
            }
        }
        # Mock the raw response from Model Registry API for GET metrics request
        mock_response_metrics = Mock(spec=requests.Response)
        mock_response_metrics.ok = True
        mock_response_metrics.json.return_value = {"items": []}
        # Mock the raw response from Model Registry API for GET params request
        mock_response_params = Mock(spec=requests.Response)
        mock_response_params.ok = True
        mock_response_params.json.return_value = {"items": []}
        # Mock the raw response from Model Registry API for PATCH request
        mock_response_patch = Mock(spec=requests.Response)
        mock_response_patch.ok = True
        mock_response_patch.json.return_value = {}
        # Mock the raw response from Model Registry API for metric POST request
        mock_response_metric = Mock(spec=requests.Response)
        mock_response_metric.ok = True
        mock_response_metric.json.return_value = {}
        # Mock the raw response from Model Registry API for param POST request
        mock_response_param = Mock(spec=requests.Response)
        mock_response_param.ok = True
        mock_response_param.json.return_value = {}
        
        mock_request.side_effect = [
            mock_response_get_run,  # get_run
            mock_response_metrics,  # GET metrics
            mock_response_params,   # GET params
            mock_response_patch,    # PATCH for tags
            mock_response_metric,   # POST for metric
            mock_response_param     # POST for param
        ]
        
        metrics = [Metric("acc", 0.9, 1234567890, 1)]
        params = [Param("lr", "0.01")]
        tags = [RunTag("new", "tag")]
        
        store.log_batch("run-123", metrics, params, tags)
        
        # Check that the correct number of calls were made
        assert mock_request.call_count == 4  # GET run + PATCH + metric POST + param POST

    # Input/Output logging tests
    @patch('modelregistry_plugin.store.requests.request')
    def test_log_inputs_datasets(self, mock_request, store):
        """Test logging dataset inputs."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {}
        mock_request.return_value = mock_response
        
        dataset = Mock()
        dataset.dataset.name = "test-dataset"
        dataset.dataset.digest = "digest123"
        dataset.dataset.source_type = "csv"
        dataset.dataset.source = "s3://bucket/data.csv"
        dataset.dataset.schema = "schema"
        dataset.dataset.profile = "profile"
        dataset.tags = [RunTag("tag1", "value1")]
        
        dataset_input = DatasetInput(dataset=dataset.dataset, tags=dataset.tags)
        
        store.log_inputs("run-123", datasets=[dataset_input])
        
        mock_request.assert_called_once()
        call_args = mock_request.call_args
        assert call_args[0][0] == "POST"  # method
        assert "/experiment_runs/run-123/artifacts" in call_args[0][1]  # endpoint
        json_data = call_args[1]['json']
        assert json_data["artifactType"] == "dataset-artifact"
        assert json_data["name"] == "test-dataset"
        assert json_data["digest"] == "digest123"
        assert json_data["sourceType"] == "csv"
        assert json_data["source"] == "s3://bucket/data.csv"
        assert json_data["schema"] == "schema"
        assert json_data["profile"] == "profile"
        assert json_data["customProperties"]["tag1"]["string_value"] == "value1"
    
    @patch('modelregistry_plugin.store.requests.request')
    def test_log_inputs_models(self, mock_request, store):
        """Test logging model inputs."""
        # Mock the raw response from Model Registry API for GET request
        mock_response_get = Mock(spec=requests.Response)
        mock_response_get.ok = True
        mock_response_get.json.return_value = {
            "customProperties": {"existing": {"string_value": "value", "metadataType": "MetadataStringValue"}}
        }
        # Mock the raw response from Model Registry API for POST request
        mock_response_post = Mock(spec=requests.Response)
        mock_response_post.ok = True
        mock_response_post.json.return_value = {}
        
        mock_request.side_effect = [mock_response_get, mock_response_post]
        
        model_input = LoggedModelInput(model_id="model-123")
        
        store.log_inputs("run-123", models=[model_input])
        
        call_args = mock_request.call_args_list
        assert len(call_args) == 2  # GET + POST
        post_call = call_args[1]
        json_data = post_call[1]['json']
        custom_props = json_data['customProperties']
        assert custom_props['mlflow.model_io_type']['string_value'] == "input"
    
    @patch('modelregistry_plugin.store.requests.request')
    def test_log_outputs(self, mock_request, store):
        """Test logging model outputs."""
        # Mock the raw response from Model Registry API for GET request
        mock_response_get = Mock(spec=requests.Response)
        mock_response_get.ok = True
        mock_response_get.json.return_value = {
            "customProperties": {"existing": {"string_value": "value", "metadataType": "MetadataStringValue"}}
        }
        # Mock the raw response from Model Registry API for POST request
        mock_response_post = Mock(spec=requests.Response)
        mock_response_post.ok = True
        mock_response_post.json.return_value = {}
        
        mock_request.side_effect = [mock_response_get, mock_response_post]
        
        model_output = LoggedModelOutput(model_id="model-123", step=1)
        
        store.log_outputs("run-123", [model_output])
        
        call_args = mock_request.call_args_list
        assert len(call_args) == 2  # GET + POST
        post_call = call_args[1]
        json_data = post_call[1]['json']
        custom_props = json_data['customProperties']
        assert custom_props['mlflow.model_io_type']['string_value'] == "output"

    # Logged model tests
    @patch('modelregistry_plugin.store.requests.request')
    def test_record_logged_model(self, mock_request, store, mock_model):
        """Test recording a logged model."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {}
        mock_request.return_value = mock_response
        
        store.record_logged_model("run-123", mock_model)
        
        mock_request.assert_called_once()
        call_args = mock_request.call_args
        assert call_args[0][0] == "POST"  # method
        assert "/experiment_runs/run-123/artifacts" in call_args[0][1]  # endpoint
        json_data = call_args[1]['json']
        assert json_data["artifactType"] == "model-artifact"
        assert json_data["name"] == "uuid-123"
        assert json_data["uri"] == "runs:/run-123/model"
        assert json_data["customProperties"]["artifactPath"]["string_value"] == "model"
        assert json_data["customProperties"]["model_uuid"]["string_value"] == "uuid-123"
    
    @patch('modelregistry_plugin.store.requests.request')
    def test_create_logged_model(self, mock_request, store):
        """Test creating a logged model."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {
            "id": "model-123",
            "name": "test-model",
            "experimentId": "exp-123",
            "uri": "s3://bucket/model",
            "createTimeSinceEpoch": "1234567890",
            "lastUpdateTimeSinceEpoch": "1234567890"
        }
        mock_request.return_value = mock_response
        
        tags = [LoggedModelTag("key1", "value1")]
        params = [LoggedModelParameter("param1", "value1")]
        
        logged_model = store.create_logged_model(
            "exp-123", 
            name="test-model",
            source_run_id="run-123",
            tags=tags,
            params=params,
            model_type="sklearn"
        )
        
        assert isinstance(logged_model, LoggedModel)
        assert logged_model.model_id == "model-123"
        assert logged_model.experiment_id == "exp-123"
        assert logged_model.name == "test-model"
        assert logged_model.source_run_id == "run-123"
        assert logged_model.model_type == "sklearn"
    
    @patch('modelregistry_plugin.store.requests.request')
    def test_search_logged_models(self, mock_request, store):
        """Test searching logged models."""
        # Mock the raw response from Model Registry API for getting runs request
        mock_response_runs = Mock(spec=requests.Response)
        mock_response_runs.ok = True
        mock_response_runs.json.return_value = {
            "items": [
                {
                    "id": "run-123",
                    "experimentId": "exp-123",
                    "name": "test-run"
                }
            ]
        }
        
        # Mock the raw response from Model Registry API for getting artifacts request
        mock_response_artifacts = Mock(spec=requests.Response)
        mock_response_artifacts.ok = True
        mock_response_artifacts.json.return_value = {
            "items": [
                {
                    "id": "model-123",
                    "name": "test-model"
                }
            ]
        }
        
        # Mock the raw response from Model Registry API for getting individual logged model request
        mock_response_model = Mock(spec=requests.Response)
        mock_response_model.ok = True
        mock_response_model.json.return_value = {
            "id": "model-123",
            "name": "test-model",
            "experimentId": "exp-123",
            "uri": "s3://bucket/model",
            "createTimeSinceEpoch": "1234567890",
            "lastUpdateTimeSinceEpoch": "1234567890",
            "customProperties": {
                "model_type": {"string_value": "sklearn", "metadataType": "MetadataStringValue"},
                "source_run_id": {"string_value": "run-123", "metadataType": "MetadataStringValue"},
                "experiment_id": {"string_value": "exp-123", "metadataType": "MetadataStringValue"}
            }
        }
        
        mock_request.side_effect = [
            mock_response_runs,      # GET /experiments/exp-123/experiment_runs
            mock_response_artifacts, # GET /experiment_runs/run-123/artifacts
            mock_response_model      # GET /artifacts/model-123
        ]
        
        result = store.search_logged_models(
            ["exp-123"],
            filter_string="model_type='sklearn'",
            max_results=10,
            page_token="token123"
        )
        
        assert isinstance(result, PagedList)
        assert len(result) == 1
        assert result[0].model_id == "model-123"
        assert result[0].name == "test-model"
        assert result[0].experiment_id == "exp-123"
        assert result[0].model_type == "sklearn"
        assert result[0].source_run_id == "run-123"
        
        # Verify API calls were made correctly
        assert mock_request.call_count == 3
        
        # Check first call: GET /experiments/exp-123/experiment_runs
        call_args = mock_request.call_args_list[0]
        assert call_args[0][0] == "GET"  # method
        assert "/experiments/exp-123/experiment_runs" in call_args[0][1]  # endpoint
        params = call_args[1]['params']
        assert params['artifactType'] == "model-artifact"
        assert params['experimentIds'] == ["exp-123"]
        assert params['pageSize'] == "10"
        assert params['pageToken'] == "token123"
        
        # Check second call: GET /experiment_runs/run-123/artifacts
        call_args = mock_request.call_args_list[1]
        assert call_args[0][0] == "GET"  # method
        assert "/experiment_runs/run-123/artifacts" in call_args[0][1]  # endpoint
        params = call_args[1]['params']
        assert params['artifactType'] == "model-artifact"
        
        # Check third call: GET /artifacts/model-123
        call_args = mock_request.call_args_list[2]
        assert call_args[0][0] == "GET"  # method
        assert "/artifacts/model-123" in call_args[0][1]  # endpoint

    @patch('modelregistry_plugin.store.requests.request')
    def test_finalize_logged_model(self, mock_request, store):
        """Test finalizing a logged model."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {
            "id": "model-123",
            "name": "test-model",
            "experimentId": "exp-123",
            "uri": "s3://bucket/model",
            "createTimeSinceEpoch": "1234567890",
            "lastUpdateTimeSinceEpoch": "1234567890",
            "customProperties": {
                "status": {"string_value": "READY", "metadataType": "MetadataStringValue"},
                "model_type": {"string_value": "sklearn", "metadataType": "MetadataStringValue"}
            }
        }
        mock_request.return_value = mock_response
        
        logged_model = store.finalize_logged_model("model-123", LoggedModelStatus.READY)
        
        assert isinstance(logged_model, LoggedModel)
        assert logged_model.model_id == "model-123"
        assert logged_model.model_type == "sklearn"
    
    @patch('modelregistry_plugin.store.requests.request')
    def test_set_logged_model_tags(self, mock_request, store):
        """Test setting logged model tags."""
        # Mock the raw response from Model Registry API for GET request
        mock_response_get = Mock(spec=requests.Response)
        mock_response_get.ok = True
        mock_response_get.json.return_value = {
            "id": "model-123",
            "name": "test-model",
            "customProperties": {
                "existing": {"string_value": "value", "metadataType": "MetadataStringValue"}
            }
        }
        # Mock the raw response from Model Registry API for PATCH request
        mock_response_patch = Mock(spec=requests.Response)
        mock_response_patch.ok = True
        mock_response_patch.json.return_value = {}
        
        mock_request.side_effect = [mock_response_get, mock_response_patch]
        
        tags = [LoggedModelTag("key1", "value1"), LoggedModelTag("key2", "value2")]
        
        store.set_logged_model_tags("model-123", tags)
        
        # Should call PATCH to update tags
        call_args = mock_request.call_args_list
        assert len(call_args) == 2  # GET + PATCH
        patch_call = call_args[1]
        json_data = patch_call[1]['json']
        custom_props = json_data['customProperties']
        assert custom_props['key1']['string_value'] == "value1"
        assert custom_props['key2']['string_value'] == "value2"
    
    @patch('modelregistry_plugin.store.requests.request')
    def test_delete_logged_model_tag(self, mock_request, store):
        """Test deleting a logged model tag."""
        # Mock the raw response from Model Registry API for GET request
        mock_response_get = Mock(spec=requests.Response)
        mock_response_get.ok = True
        mock_response_get.json.return_value = {
            "id": "model-123",
            "name": "test-model",
            "customProperties": {
                "key1": {"string_value": "value1", "metadataType": "MetadataStringValue"},
                "key2": {"string_value": "value2", "metadataType": "MetadataStringValue"}
            }
        }
        # Mock the raw response from Model Registry API for PATCH request
        mock_response_patch = Mock(spec=requests.Response)
        mock_response_patch.ok = True
        mock_response_patch.json.return_value = {}
        
        mock_request.side_effect = [mock_response_get, mock_response_patch]
        
        store.delete_logged_model_tag("model-123", "key1")
        
        # Should call PATCH to update tags without key1
        call_args = mock_request.call_args_list
        assert len(call_args) == 2  # GET + PATCH
        patch_call = call_args[1]
        json_data = patch_call[1]['json']
        custom_props = json_data['customProperties']
        assert "key1" not in custom_props
        assert "key2" in custom_props
    
    @patch('modelregistry_plugin.store.requests.request')
    def test_get_logged_model(self, mock_request, store):
        """Test getting a logged model."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {
            "id": "model-123",
            "name": "test-model",
            "uri": "s3://bucket/model",
            "createTimeSinceEpoch": "1234567890",
            "lastUpdateTimeSinceEpoch": "1234567890",
            "customProperties": {
                "experiment_id": {"string_value": "exp-123", "metadataType": "MetadataStringValue"},
                "model_type": {"string_value": "sklearn", "metadataType": "MetadataStringValue"},
                "source_run_id": {"string_value": "run-123", "metadataType": "MetadataStringValue"},
                "param_lr": {"string_value": "0.01", "metadataType": "MetadataStringValue"}
            }
        }
        mock_request.return_value = mock_response
        
        logged_model = store.get_logged_model("model-123")
        
        assert isinstance(logged_model, LoggedModel)
        assert logged_model.model_id == "model-123"
        assert logged_model.name == "test-model"
        assert logged_model.experiment_id == "exp-123"
        assert logged_model.model_type == "sklearn"
        assert logged_model.source_run_id == "run-123"
        assert len(logged_model.params) == 1
        assert logged_model.params["lr"] == "0.01"
    
    @patch('modelregistry_plugin.store.requests.request')
    def test_delete_logged_model(self, mock_request, store):
        """Test deleting a logged model."""
        # Mock the raw response from Model Registry API for GET request
        mock_response_get = Mock(spec=requests.Response)
        mock_response_get.ok = True
        mock_response_get.json.return_value = {
            "id": "model-123",
            "name": "test-model",
            "customProperties": {
                "existing": {"string_value": "value", "metadataType": "MetadataStringValue"}
            }
        }
        # Mock the raw response from Model Registry API for PATCH request
        mock_response_patch = Mock(spec=requests.Response)
        mock_response_patch.ok = True
        mock_response_patch.json.return_value = {}
        
        mock_request.side_effect = [mock_response_get, mock_response_patch]
        
        store.delete_logged_model("model-123")
        
        # Should call PATCH to mark as archived
        call_args = mock_request.call_args_list
        assert len(call_args) == 1  # PATCH only
        patch_call = call_args[0]
        json_data = patch_call[1]['json']
        assert json_data["artifactType"] == "model-artifact"
        custom_props = json_data['customProperties']
        assert custom_props['state']['string_value'] == "MARKED_FOR_DELETION"

    # Search runs test
    @patch('modelregistry_plugin.store.requests.request')
    def test_search_runs_all(self, mock_request, store):
        """Test searching runs."""
        # Mock the raw response from Model Registry API for search request
        mock_response_search = Mock(spec=requests.Response)
        mock_response_search.ok = True
        mock_response_search.json.return_value = {
            "items": [
                {
                    "id": "run-123",
                    "experimentId": "exp-123",
                    "name": "test-run",
                    "state": "RUNNING",
                    "owner": "user123",
                    "startTimeSinceEpoch": "1234567890",
                    "externalId": "s3://bucket/artifacts/experiments/exp-123/run-123",
                    "customProperties": {
                        "tag1": {"string_value": "value1", "metadataType": "MetadataStringValue"}
                    }
                }
            ],
            "nextPageToken": "token123"
        }
        # Mock the raw response from Model Registry API for metrics request
        mock_response_metrics = Mock(spec=requests.Response)
        mock_response_metrics.ok = True
        mock_response_metrics.json.return_value = {"items": []}
        # Mock the raw response from Model Registry API for params request
        mock_response_params = Mock(spec=requests.Response)
        mock_response_params.ok = True
        mock_response_params.json.return_value = {"items": []}
        
        mock_request.side_effect = [
            mock_response_search,  # search request
            mock_response_metrics,  # metrics request
            mock_response_params    # params request
        ]
        
        result = store.search_runs(
            ["exp-123"],
            filter_string="status='RUNNING'",
            max_results=10,
            page_token="token123"
        )
        
        assert isinstance(result, PagedList)
        assert len(result) == 1
        assert result[0].info.run_id == "run-123"
        assert result[0].info.artifact_uri == "s3://bucket/artifacts/experiments/exp-123/run-123"
        assert result.token == "token123"

    @patch('modelregistry_plugin.store.requests.request')
    def test_search_runs_active_only(self, mock_request, store):
        """Test searching runs with ACTIVE_ONLY view type."""
        # Mock the raw response from Model Registry API for search request
        mock_response_search = Mock(spec=requests.Response)
        mock_response_search.ok = True
        mock_response_search.json.return_value = {
            "items": [
                {
                    "id": "run-123",
                    "experimentId": "exp-123",
                    "name": "active-run",
                    "status": "RUNNING",  # MLflow RunStatus
                    "state": "LIVE",      # ModelRegistry lifecycle state
                    "owner": "user123",
                    "startTimeSinceEpoch": "1234567890",
                    "externalId": "s3://bucket/artifacts/experiments/exp-123/run-123",
                    "customProperties": {}
                },
                {
                    "id": "run-456",
                    "experimentId": "exp-123",
                    "name": "deleted-run",
                    "status": "FINISHED",  # MLflow RunStatus
                    "state": "ARCHIVED",   # ModelRegistry lifecycle state
                    "owner": "user123",
                    "startTimeSinceEpoch": "1234567890",
                    "externalId": "s3://bucket/artifacts/experiments/exp-123/run-456",
                    "customProperties": {}
                }
            ],
            "nextPageToken": "token123"
        }
        # Mock the raw response from Model Registry API for metrics request (for active run)
        mock_response_metrics_active = Mock(spec=requests.Response)
        mock_response_metrics_active.ok = True
        mock_response_metrics_active.json.return_value = {"items": []}
        # Mock the raw response from Model Registry API for params request (for active run)
        mock_response_params_active = Mock(spec=requests.Response)
        mock_response_params_active.ok = True
        mock_response_params_active.json.return_value = {"items": []}
        # Mock the raw response from Model Registry API for metrics request (for deleted run)
        mock_response_metrics_deleted = Mock(spec=requests.Response)
        mock_response_metrics_deleted.ok = True
        mock_response_metrics_deleted.json.return_value = {"items": []}
        # Mock the raw response from Model Registry API for params request (for deleted run)
        mock_response_params_deleted = Mock(spec=requests.Response)
        mock_response_params_deleted.ok = True
        mock_response_params_deleted.json.return_value = {"items": []}
        
        mock_request.side_effect = [
            mock_response_search,      # search request
            mock_response_metrics_active,   # metrics request for active run
            mock_response_params_active,    # params request for active run
            mock_response_metrics_deleted,  # metrics request for deleted run
            mock_response_params_deleted    # params request for deleted run
        ]
        
        result = store.search_runs(
            ["exp-123"],
            run_view_type=ViewType.ACTIVE_ONLY
        )
        
        # Should return all runs since filtering is not implemented yet
        assert isinstance(result, PagedList)
        assert len(result) == 1
        assert result[0].info.run_id == "run-123"
        assert result[0].info.artifact_uri == "s3://bucket/artifacts/experiments/exp-123/run-123"

    @patch('modelregistry_plugin.store.requests.request')
    def test_search_runs_deleted_only(self, mock_request, store):
        """Test searching runs with DELETED_ONLY view type."""
        # Mock the raw response from Model Registry API for search request
        mock_response_search = Mock(spec=requests.Response)
        mock_response_search.ok = True
        mock_response_search.json.return_value = {
            "items": [
                {
                    "id": "run-123",
                    "experimentId": "exp-123",
                    "name": "active-run",
                    "status": "RUNNING",  # MLflow RunStatus
                    "state": "LIVE",      # ModelRegistry lifecycle state
                    "owner": "user123",
                    "startTimeSinceEpoch": "1234567890",
                    "externalId": "s3://bucket/artifacts/experiments/exp-123/run-123",
                    "customProperties": {}
                },
                {
                    "id": "run-456",
                    "experimentId": "exp-123",
                    "name": "deleted-run",
                    "status": "FINISHED",  # MLflow RunStatus
                    "state": "ARCHIVED",   # ModelRegistry lifecycle state
                    "owner": "user123",
                    "startTimeSinceEpoch": "1234567890",
                    "externalId": "s3://bucket/artifacts/experiments/exp-123/run-456",
                    "customProperties": {}
                }
            ],
            "nextPageToken": "token123"
        }
        # Mock the raw response from Model Registry API for metrics request (for active run)
        mock_response_metrics_active = Mock(spec=requests.Response)
        mock_response_metrics_active.ok = True
        mock_response_metrics_active.json.return_value = {"items": []}
        # Mock the raw response from Model Registry API for params request (for active run)
        mock_response_params_active = Mock(spec=requests.Response)
        mock_response_params_active.ok = True
        mock_response_params_active.json.return_value = {"items": []}
        # Mock the raw response from Model Registry API for metrics request (for deleted run)
        mock_response_metrics_deleted = Mock(spec=requests.Response)
        mock_response_metrics_deleted.ok = True
        mock_response_metrics_deleted.json.return_value = {"items": []}
        # Mock the raw response from Model Registry API for params request (for deleted run)
        mock_response_params_deleted = Mock(spec=requests.Response)
        mock_response_params_deleted.ok = True
        mock_response_params_deleted.json.return_value = {"items": []}
        
        mock_request.side_effect = [
            mock_response_search,      # search request
            mock_response_metrics_active,   # metrics request for active run
            mock_response_params_active,    # params request for active run
            mock_response_metrics_deleted,  # metrics request for deleted run
            mock_response_params_deleted    # params request for deleted run
        ]
        
        result = store.search_runs(
            ["exp-123"],
            run_view_type=ViewType.DELETED_ONLY
        )
        
        # Should return all runs since filtering is not implemented yet
        assert isinstance(result, PagedList)
        assert len(result) == 1
        assert result[0].info.run_id == "run-456"
        assert result[0].info.artifact_uri == "s3://bucket/artifacts/experiments/exp-123/run-456"

    # Error handling tests
    @patch('modelregistry_plugin.store.requests.request')
    def test_request_error_handling(self, mock_request, store):
        """Test error handling in requests."""
        # Mock the raw response from Model Registry API with error
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = False
        mock_response.text = "Server error"
        mock_response.json.return_value = {"message": "Invalid JSON"}
        mock_response.status_code = 400
        mock_request.return_value = mock_response
        
        with pytest.raises(MlflowException) as exc_info:
            store._request("GET", "/test")
        
        assert "Model Registry API error: Invalid JSON" in str(exc_info.value)
    
    def test_get_artifact_location(self, store):
        """Test getting artifact location."""
        response_with_external_id = {"externalId": "s3://bucket/artifacts"}
        response_without_external_id = {}
        
        assert store._get_artifact_location(response_with_external_id) == "s3://bucket/artifacts"
        assert store._get_artifact_location(response_without_external_id) is None
        
        # Test with artifact_uri set
        store.artifact_uri = "s3://default/artifacts"
        assert store._get_artifact_location(response_without_external_id) == "s3://default/artifacts"

    @patch('modelregistry_plugin.store.requests.request')
    def test_create_run_without_artifact_uri(self, mock_request):
        """Test run creation when experiment has no artifact_location."""
        # Create store without artifact_uri
        store = ModelRegistryStore("modelregistry://localhost:8080")
        
        # Mock the raw response from Model Registry API for creating run
        mock_response_create = Mock(spec=requests.Response)
        mock_response_create.ok = True
        mock_response_create.json.return_value = {
            "id": "run-123",
            "experimentId": "exp-123",
            "createTimeSinceEpoch": "1234567890"
        }
        
        # Mock the raw response from Model Registry API for getting experiment (no artifact_location)
        mock_response_experiment = Mock(spec=requests.Response)
        mock_response_experiment.ok = True
        mock_response_experiment.json.return_value = {
            "id": "exp-123",
            "name": "test-experiment",
            "customProperties": {}
        }
        
        mock_request.side_effect = [mock_response_create, mock_response_experiment]
        
        run = store.create_run("exp-123", start_time=1234567890)
        
        assert isinstance(run, Run)
        assert run.info.run_id == "run-123"
        assert run.info.experiment_id == "exp-123"
        assert run.info.status == RunStatus.RUNNING
        assert run.info.artifact_uri is None
        
        # Should make only 2 calls since experiment has no artifact_location
        assert mock_request.call_count == 2
        
        # Check first call (POST to create run)
        call_args = mock_request.call_args_list[0]
        assert call_args[0][0] == "POST"  # method
        assert "/experiment_runs" in call_args[0][1]  # endpoint
        json_data = call_args[1]['json']
        assert json_data["experimentId"] == "exp-123"
        assert json_data["status"] == "RUNNING"
        assert json_data.get("externalId") is None

    @patch('modelregistry_plugin.store.requests.request')
    def test_get_metric_history_basic(self, mock_request, store):
        """Test basic metric history retrieval."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {
            "items": [
                {
                    "name": "accuracy",
                    "value": 0.95,
                    "timestamp": "1234567890",
                    "step": 1,
                    "createTimeSinceEpoch": "1234567890"
                },
                {
                    "name": "accuracy",
                    "value": 0.97,
                    "timestamp": "1234567891",
                    "step": 2,
                    "createTimeSinceEpoch": "1234567891"
                }
            ]
        }
        mock_request.return_value = mock_response
        
        metrics = store.get_metric_history("run-123", "accuracy")
        
        assert len(metrics) == 2
        assert all(metric.key == "accuracy" for metric in metrics)
        assert metrics[0].value == 0.95
        assert metrics[0].step == 1
        assert metrics[0].timestamp == 1234567890
        assert metrics[1].value == 0.97
        assert metrics[1].step == 2
        assert metrics[1].timestamp == 1234567891
        
        # Verify API call
        mock_request.assert_called_once()
        call_args = mock_request.call_args
        assert call_args[0][0] == "GET"  # method
        assert "/experiment_runs/run-123/metric_history" in call_args[0][1]  # endpoint
        params = call_args[1]['params']
        assert params["name"] == "accuracy"

    @patch('modelregistry_plugin.store.requests.request')
    def test_get_metric_history_with_max_results(self, mock_request, store):
        """Test metric history with max_results limit."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {
            "items": [
                {
                    "name": "accuracy",
                    "value": 0.95,
                    "timestamp": "1234567890",
                    "step": 1,
                    "createTimeSinceEpoch": "1234567890"
                }
            ]
        }
        mock_request.return_value = mock_response
        
        metrics = store.get_metric_history("run-123", "accuracy", max_results=1)
        
        assert len(metrics) == 1
        assert metrics[0].key == "accuracy"
        assert metrics[0].value == 0.95
        
        # Verify API call with max_results
        mock_request.assert_called_once()
        call_args = mock_request.call_args
        params = call_args[1]['params']
        assert params["name"] == "accuracy"
        assert params["pageSize"] == 1

    @patch('modelregistry_plugin.store.requests.request')
    def test_get_metric_history_with_page_token(self, mock_request, store):
        """Test metric history with page_token."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {
            "items": [
                {
                    "name": "accuracy",
                    "value": 0.98,
                    "timestamp": "1234567892",
                    "step": 3,
                    "createTimeSinceEpoch": "1234567892"
                }
            ]
        }
        mock_request.return_value = mock_response
        
        metrics = store.get_metric_history("run-123", "accuracy", page_token="token123")
        
        assert len(metrics) == 1
        assert metrics[0].key == "accuracy"
        assert metrics[0].value == 0.98
        
        # Verify API call with page_token
        mock_request.assert_called_once()
        call_args = mock_request.call_args
        params = call_args[1]['params']
        assert params["name"] == "accuracy"
        assert params["pageToken"] == "token123"

    @patch('modelregistry_plugin.store.requests.request')
    def test_get_metric_history_empty(self, mock_request, store):
        """Test metric history when no metrics exist."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {"items": []}
        mock_request.return_value = mock_response
        
        metrics = store.get_metric_history("run-123", "nonexistent")
        
        assert len(metrics) == 0
        
        # Verify API call
        mock_request.assert_called_once()
        call_args = mock_request.call_args
        params = call_args[1]['params']
        assert params["name"] == "nonexistent"

    @patch('modelregistry_plugin.store.requests.request')
    def test_get_metric_history_uses_timestamp_fallback(self, mock_request, store):
        """Test that metric history uses createTimeSinceEpoch when timestamp is not available."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {
            "items": [
                {
                    "name": "accuracy",
                    "value": 0.95,
                    "step": 1,
                    "createTimeSinceEpoch": "1234567890"
                    # Note: no "timestamp" field
                }
            ]
        }
        mock_request.return_value = mock_response
        
        metrics = store.get_metric_history("run-123", "accuracy")
        
        assert len(metrics) == 1
        assert metrics[0].timestamp == 1234567890  # Should use createTimeSinceEpoch

    @patch('modelregistry_plugin.store.requests.request')
    def test_get_metric_history_uses_timestamp_over_create_time(self, mock_request, store):
        """Test that metric history prefers timestamp over createTimeSinceEpoch."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {
            "items": [
                {
                    "name": "accuracy",
                    "value": 0.95,
                    "timestamp": "1234567890",
                    "step": 1,
                    "createTimeSinceEpoch": "9999999999"  # Different value
                }
            ]
        }
        mock_request.return_value = mock_response
        
        metrics = store.get_metric_history("run-123", "accuracy")
        
        assert len(metrics) == 1
        assert metrics[0].timestamp == 1234567890  # Should use timestamp, not createTimeSinceEpoch

    @patch('modelregistry_plugin.store.requests.request')
    def test_get_metric_history_bulk_interval_from_steps_basic(self, mock_request, store):
        """Test basic bulk metric history retrieval for specific steps."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {
            "items": [
                {
                    "name": "accuracy",
                    "value": 0.95,
                    "timestamp": "1234567890",
                    "step": 1,
                    "createTimeSinceEpoch": "1234567890"
                },
                {
                    "name": "accuracy",
                    "value": 0.97,
                    "timestamp": "1234567891",
                    "step": 2,
                    "createTimeSinceEpoch": "1234567891"
                },
                {
                    "name": "accuracy",
                    "value": 0.98,
                    "timestamp": "1234567892",
                    "step": 3,
                    "createTimeSinceEpoch": "1234567892"
                },
                {
                    "name": "accuracy",
                    "value": 0.99,
                    "timestamp": "1234567893",
                    "step": 4,
                    "createTimeSinceEpoch": "1234567893"
                }
            ]
        }
        mock_request.return_value = mock_response
        
        # Get metrics for specific steps
        metrics = store.get_metric_history_bulk_interval_from_steps(
            "run-123", "accuracy", steps=[1, 3], max_results=2
        )
        
        assert len(metrics) == 2
        assert all(m.run_id == "run-123" for m in metrics)
        assert all(m.key == "accuracy" for m in metrics)
        
        # Should return metrics for steps 1 and 3, sorted by step then timestamp
        assert metrics[0].step == 1
        assert metrics[0].value == 0.95
        assert metrics[1].step == 3
        assert metrics[1].value == 0.98

    @patch('modelregistry_plugin.store.requests.request')
    def test_get_metric_history_bulk_interval_from_steps_filtering(self, mock_request, store):
        """Test that bulk metric history correctly filters by steps."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {
            "items": [
                {
                    "name": "accuracy",
                    "value": 0.95,
                    "timestamp": "1234567890",
                    "step": 1,
                    "createTimeSinceEpoch": "1234567890"
                },
                {
                    "name": "accuracy",
                    "value": 0.97,
                    "timestamp": "1234567891",
                    "step": 2,
                    "createTimeSinceEpoch": "1234567891"
                },
                {
                    "name": "accuracy",
                    "value": 0.98,
                    "timestamp": "1234567892",
                    "step": 3,
                    "createTimeSinceEpoch": "1234567892"
                }
            ]
        }
        mock_request.return_value = mock_response
        
        # Request steps that don't exist
        metrics = store.get_metric_history_bulk_interval_from_steps(
            "run-123", "accuracy", steps=[5, 6], max_results=10
        )
        
        assert len(metrics) == 0

    @patch('modelregistry_plugin.store.requests.request')
    def test_get_metric_history_bulk_interval_from_steps_max_results(self, mock_request, store):
        """Test that bulk metric history respects max_results limit."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {
            "items": [
                {
                    "name": "accuracy",
                    "value": 0.95,
                    "timestamp": "1234567890",
                    "step": 1,
                    "createTimeSinceEpoch": "1234567890"
                },
                {
                    "name": "accuracy",
                    "value": 0.97,
                    "timestamp": "1234567891",
                    "step": 2,
                    "createTimeSinceEpoch": "1234567891"
                },
                {
                    "name": "accuracy",
                    "value": 0.98,
                    "timestamp": "1234567892",
                    "step": 3,
                    "createTimeSinceEpoch": "1234567892"
                }
            ]
        }
        mock_request.return_value = mock_response
        
        # Request more steps than max_results
        metrics = store.get_metric_history_bulk_interval_from_steps(
            "run-123", "accuracy", steps=[1, 2, 3], max_results=2
        )
        
        assert len(metrics) == 2  # Should be limited by max_results

    @patch('modelregistry_plugin.store.requests.request')
    def test_get_metric_history_bulk_interval_from_steps_sorting(self, mock_request, store):
        """Test that bulk metric history sorts correctly by step then timestamp."""
        # Mock the raw response from Model Registry API with unsorted data
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {
            "items": [
                {
                    "name": "accuracy",
                    "value": 0.97,
                    "timestamp": "1234567891",
                    "step": 2,
                    "createTimeSinceEpoch": "1234567891"
                },
                {
                    "name": "accuracy",
                    "value": 0.95,
                    "timestamp": "1234567890",
                    "step": 1,
                    "createTimeSinceEpoch": "1234567890"
                },
                {
                    "name": "accuracy",
                    "value": 0.96,
                    "timestamp": "1234567895",
                    "step": 1,  # Same step, later timestamp
                    "createTimeSinceEpoch": "1234567895"
                }
            ]
        }
        mock_request.return_value = mock_response
        
        metrics = store.get_metric_history_bulk_interval_from_steps(
            "run-123", "accuracy", steps=[1, 2], max_results=10
        )
        
        assert len(metrics) == 3
        
        # Should be sorted by step first, then timestamp
        assert metrics[0].step == 1 and metrics[0].timestamp == 1234567890
        assert metrics[1].step == 1 and metrics[1].timestamp == 1234567895
        assert metrics[2].step == 2 and metrics[2].timestamp == 1234567891

    @patch('modelregistry_plugin.store.requests.request')
    def test_get_metric_history_bulk_interval_from_steps_empty_steps(self, mock_request, store):
        """Test bulk metric history with empty steps list."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {
            "items": [
                {
                    "name": "accuracy",
                    "value": 0.95,
                    "timestamp": "1234567890",
                    "step": 1,
                    "createTimeSinceEpoch": "1234567890"
                }
            ]
        }
        mock_request.return_value = mock_response
        
        metrics = store.get_metric_history_bulk_interval_from_steps(
            "run-123", "accuracy", steps=[], max_results=10
        )
        
        assert len(metrics) == 0

    @patch('modelregistry_plugin.store.requests.request')
    def test_get_metric_history_bulk_interval_from_steps_zero_max_results(self, mock_request, store):
        """Test bulk metric history with zero max_results."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {
            "items": [
                {
                    "name": "accuracy",
                    "value": 0.95,
                    "timestamp": "1234567890",
                    "step": 1,
                    "createTimeSinceEpoch": "1234567890"
                }
            ]
        }
        mock_request.return_value = mock_response
        
        metrics = store.get_metric_history_bulk_interval_from_steps(
            "run-123", "accuracy", steps=[1], max_results=0
        )
        
        assert len(metrics) == 0

    @patch('modelregistry_plugin.store.requests.request')
    def test_get_metric_history_bulk_interval_from_steps_metric_with_run_id_structure(self, mock_request, store):
        """Test that MetricWithRunId objects have the correct structure."""
        # Mock the raw response from Model Registry API
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = True
        mock_response.json.return_value = {
            "items": [
                {
                    "name": "accuracy",
                    "value": 0.95,
                    "timestamp": "1234567890",
                    "step": 1,
                    "createTimeSinceEpoch": "1234567890"
                }
            ]
        }
        mock_request.return_value = mock_response
        
        metrics = store.get_metric_history_bulk_interval_from_steps(
            "run-123", "accuracy", steps=[1], max_results=10
        )
        
        assert len(metrics) == 1
        metric = metrics[0]
        
        # Check expected values
        assert metric.run_id == "run-123"
        assert metric.key == "accuracy"
        assert metric.value == 0.95
        assert metric.step == 1
        assert metric.timestamp == 1234567890