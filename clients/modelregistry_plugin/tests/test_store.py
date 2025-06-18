"""
Tests for ModelRegistryStore
"""

import pytest
from unittest.mock import Mock, patch, MagicMock
import requests

from modelregistry_plugin.store import ModelRegistryStore
from mlflow.entities import Experiment, Run, RunInfo, RunData, RunStatus, ExperimentTag, RunTag, Param, Metric


class TestModelRegistryStore:
    
    @pytest.fixture
    def store(self):
        """Create a ModelRegistryStore instance for testing."""
        return ModelRegistryStore("modelregistry://localhost:8080")
    
    @pytest.fixture
    def mock_response(self):
        """Create a mock response object."""
        response = Mock(spec=requests.Response)
        response.ok = True
        response.json.return_value = {}
        return response
    
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
    
    @patch('modelregistry_plugin.store.requests.request')
    def test_request_success(self, mock_request, store, mock_response):
        """Test successful API request."""
        mock_request.return_value = mock_response
        
        response = store._request("GET", "/test")
        
        mock_request.assert_called_once()
        assert response == mock_response
    
    @patch('modelregistry_plugin.store.requests.request')
    def test_request_failure(self, mock_request, store):
        """Test failed API request."""
        mock_response = Mock(spec=requests.Response)
        mock_response.ok = False
        mock_response.json.return_value = {"message": "Test error"}
        mock_request.return_value = mock_response
        
        with pytest.raises(Exception):  # MlflowException
            store._request("GET", "/test")
    
    @patch('modelregistry_plugin.store.ModelRegistryStore._request')
    def test_create_experiment(self, mock_request, store):
        """Test experiment creation."""
        mock_request.return_value.json.return_value = {"id": "123"}
        
        experiment_id = store.create_experiment("test-experiment")
        
        assert experiment_id == "123"
        mock_request.assert_called_once_with(
            "POST", "/experiments", 
            json={
                "name": "test-experiment",
                "description": "MLflow experiment: test-experiment",
                "customProperties": {}
            }
        )
    
    @patch('modelregistry_plugin.store.ModelRegistryStore._request')
    def test_get_experiment(self, mock_request, store):
        """Test getting experiment by ID."""
        mock_request.return_value.json.return_value = {
            "id": "123",
            "name": "test-experiment",
            "customProperties": {"tag1": "value1"}
        }
        
        experiment = store.get_experiment("123")
        
        assert isinstance(experiment, Experiment)
        assert experiment.experiment_id == "123"
        assert experiment.name == "test-experiment"
        assert len(experiment.tags) == 1
        assert experiment.tags[0].key == "tag1"
        assert experiment.tags[0].value == "value1"
    
    @patch('modelregistry_plugin.store.ModelRegistryStore._request')
    def test_create_run(self, mock_request, store):
        """Test run creation."""
        mock_request.return_value.json.return_value = {
            "id": "run-123",
            "experimentId": "exp-123"
        }
        
        run = store.create_run("exp-123", start_time=1234567890)
        
        assert isinstance(run, Run)
        assert run.info.run_id == "run-123"
        assert run.info.experiment_id == "exp-123"
        assert run.info.status == RunStatus.RUNNING
    
    @patch('modelregistry_plugin.store.ModelRegistryStore._request')
    @patch('modelregistry_plugin.store.ModelRegistryStore._get_run_metrics')
    @patch('modelregistry_plugin.store.ModelRegistryStore._get_run_params')
    def test_get_run(self, mock_get_params, mock_get_metrics, mock_request, store):
        """Test getting run by ID."""
        mock_request.return_value.json.return_value = {
            "id": "run-123",
            "experimentId": "exp-123",
            "name": "test-run",
            "state": "RUNNING",
            "customProperties": {"tag1": "value1"}
        }
        mock_get_metrics.return_value = []
        mock_get_params.return_value = []
        
        run = store.get_run("run-123")
        
        assert isinstance(run, Run)
        assert run.info.run_id == "run-123"
        assert run.info.experiment_id == "exp-123"
        assert run.info.run_name == "test-run"
        assert len(run.data.tags) == 1
    
    @patch('modelregistry_plugin.store.ModelRegistryStore._request')
    def test_log_metric(self, mock_request, store):
        """Test logging a metric."""
        metric = Metric("accuracy", 0.95, 1234567890, 1)
        
        store.log_metric("run-123", metric)
        
        mock_request.assert_called_once_with(
            "POST", "/metrics",
            json={
                "name": "accuracy",
                "experimentRunId": "run-123",
                "value": 0.95,
                "customProperties": {
                    "timestamp": "1234567890",
                    "step": "1"
                }
            }
        )
    
    @patch('modelregistry_plugin.store.ModelRegistryStore._request')
    def test_log_param(self, mock_request, store):
        """Test logging a parameter."""
        param = Param("learning_rate", "0.01")
        
        store.log_param("run-123", param)
        
        mock_request.assert_called_once_with(
            "POST", "/parameters",
            json={
                "name": "learning_rate",
                "experimentRunId": "run-123",
                "value": "0.01"
            }
        )
    
    @patch('modelregistry_plugin.store.ModelRegistryStore._request')
    def test_list_experiments(self, mock_request, store):
        """Test listing experiments."""
        mock_request.return_value.json.return_value = {
            "experiments": [
                {"id": "1", "name": "exp1", "state": "LIVE"},
                {"id": "2", "name": "exp2", "state": "ARCHIVED"}
            ]
        }
        
        experiments = store.list_experiments()
        
        assert len(experiments) == 2
        assert experiments[0].experiment_id == "1"
        assert experiments[0].name == "exp1"