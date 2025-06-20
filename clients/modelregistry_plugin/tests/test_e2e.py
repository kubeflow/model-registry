"""
End-to-end tests for ModelRegistryStore with real Model Registry server.

These tests require a running Model Registry server and proper authentication.
Set the following environment variables to configure the test:

- MODEL_REGISTRY_HOST: Hostname of the Model Registry server
- MODEL_REGISTRY_PORT: Port of the Model Registry server (default: 8080)
- MODEL_REGISTRY_SECURE: Whether to use HTTPS (default: false)
- MODEL_REGISTRY_TOKEN: Authentication token for the Model Registry server
"""

import os
import pytest
import time
import uuid
from unittest.mock import patch

import mlflow
from mlflow.entities import (
    Experiment, Run, RunInfo, RunData, RunStatus, RunTag, Param, Metric,
    ViewType, LifecycleStage, ExperimentTag
)

from modelregistry_plugin.store import ModelRegistryStore


class TestModelRegistryStoreE2E:
    """End-to-end tests for ModelRegistryStore with real Model Registry server."""
    
    @pytest.fixture(scope="class")
    def store(self):
        """Create a ModelRegistryStore instance with real server connection."""
        host = os.getenv("MODEL_REGISTRY_HOST")
        port = os.getenv("MODEL_REGISTRY_PORT", "8080")
        secure = os.getenv("MODEL_REGISTRY_SECURE", "false").lower() == "true"
        token = os.getenv("MODEL_REGISTRY_TOKEN")
        
        if not host:
            pytest.skip("MODEL_REGISTRY_HOST environment variable not set")
        
        if not token:
            pytest.skip("MODEL_REGISTRY_TOKEN environment variable not set")
        
        # Construct the store URI
        scheme = "https" if secure else "http"
        store_uri = f"modelregistry://{host}:{port}"
        
        # Set the token in environment for auth
        os.environ["MODEL_REGISTRY_TOKEN"] = token
        
        return ModelRegistryStore(store_uri=store_uri)
    
    @pytest.fixture(scope="class")
    def experiment_name(self):
        """Generate a unique experiment name for testing."""
        return f"e2e-test-{uuid.uuid4().hex[:8]}"
    
    @pytest.fixture(scope="class")
    def experiment_id(self, store, experiment_name):
        """Create a test experiment and return its ID."""
        try:
            experiment_id = store.create_experiment(experiment_name)
            yield experiment_id
        finally:
            # Cleanup: delete the experiment
            try:
                store.delete_experiment(experiment_id)
            except:
                pass  # Ignore cleanup errors
    
    def test_store_connection(self, store):
        """Test that we can connect to the Model Registry server."""
        # Try to list experiments to verify connection
        experiments = store.list_experiments()
        assert isinstance(experiments, list)
        print(f"Successfully connected to Model Registry. Found {len(experiments)} experiments.")
    
    def test_create_and_get_experiment(self, store, experiment_name):
        """Test creating and retrieving an experiment."""
        # Create experiment
        experiment_id = store.create_experiment(experiment_name)
        assert experiment_id is not None
        assert isinstance(experiment_id, str)
        
        # Get experiment by ID
        experiment = store.get_experiment(experiment_id)
        assert isinstance(experiment, Experiment)
        assert experiment.experiment_id == experiment_id
        assert experiment.name == experiment_name
        assert experiment.lifecycle_stage == LifecycleStage.ACTIVE
        
        # Get experiment by name
        experiment_by_name = store.get_experiment_by_name(experiment_name)
        assert experiment_by_name is not None
        assert experiment_by_name.experiment_id == experiment_id
        
        # Cleanup
        store.delete_experiment(experiment_id)
    
    def test_create_and_get_run(self, store, experiment_id):
        """Test creating and retrieving a run."""
        # Create run
        run_name = f"test-run-{uuid.uuid4().hex[:8]}"
        run = store.create_run(
            experiment_id=experiment_id,
            user_id="test-user",
            run_name=run_name
        )
        
        assert isinstance(run, Run)
        assert run.info.experiment_id == experiment_id
        assert run.info.user_id == "test-user"
        assert run.info.status == RunStatus.RUNNING
        assert run.info.lifecycle_stage == LifecycleStage.ACTIVE
        
        run_id = run.info.run_id
        
        # Get run by ID
        retrieved_run = store.get_run(run_id)
        assert isinstance(retrieved_run, Run)
        assert retrieved_run.info.run_id == run_id
        assert retrieved_run.info.experiment_id == experiment_id
        
        # Update run status
        updated_info = store.update_run_info(
            run_id=run_id,
            run_status=RunStatus.FINISHED,
            end_time=int(time.time() * 1000)
        )
        assert updated_info.status == RunStatus.FINISHED
        
        # Cleanup
        store.delete_run(run_id)
    
    def test_log_metrics_and_params(self, store, experiment_id):
        """Test logging metrics and parameters."""
        # Create run
        run = store.create_run(experiment_id=experiment_id)
        run_id = run.info.run_id
        
        try:
            # Log parameters
            param1 = Param(key="learning_rate", value="0.01")
            param2 = Param(key="epochs", value="100")
            store.log_param(run_id, param1)
            store.log_param(run_id, param2)
            
            # Log metrics
            metric1 = Metric(key="accuracy", value=0.95, timestamp=int(time.time() * 1000), step=1)
            metric2 = Metric(key="loss", value=0.05, timestamp=int(time.time() * 1000), step=1)
            store.log_metric(run_id, metric1)
            store.log_metric(run_id, metric2)
            
            # Get run and verify data
            retrieved_run = store.get_run(run_id)
            assert len(retrieved_run.data.params) == 2
            assert len(retrieved_run.data.metrics) == 2
            
            # Verify parameters
            param_dict = {p.key: p.value for p in retrieved_run.data.params}
            assert param_dict["learning_rate"] == "0.01"
            assert param_dict["epochs"] == "100"
            
            # Verify metrics
            metric_dict = {m.key: m.value for m in retrieved_run.data.metrics}
            assert metric_dict["accuracy"] == 0.95
            assert metric_dict["loss"] == 0.05
            
            # Test metric history
            accuracy_history = store.get_metric_history(run_id, "accuracy")
            assert len(accuracy_history) == 1
            assert accuracy_history[0].key == "accuracy"
            assert accuracy_history[0].value == 0.95
            
        finally:
            # Cleanup
            store.delete_run(run_id)
    
    def test_log_batch(self, store, experiment_id):
        """Test batch logging of metrics, parameters, and tags."""
        # Create run
        run = store.create_run(experiment_id=experiment_id)
        run_id = run.info.run_id
        
        try:
            # Prepare batch data
            metrics = [
                Metric(key="batch_metric1", value=1.0, timestamp=int(time.time() * 1000), step=1),
                Metric(key="batch_metric2", value=2.0, timestamp=int(time.time() * 1000), step=1)
            ]
            params = [
                Param(key="batch_param1", value="value1"),
                Param(key="batch_param2", value="value2")
            ]
            tags = [
                RunTag(key="batch_tag1", value="tag_value1"),
                RunTag(key="batch_tag2", value="tag_value2")
            ]
            
            # Log batch
            store.log_batch(run_id, metrics=metrics, params=params, tags=tags)
            
            # Verify batch data
            retrieved_run = store.get_run(run_id)
            assert len(retrieved_run.data.metrics) == 2
            assert len(retrieved_run.data.params) == 2
            assert len(retrieved_run.data.tags) == 2
            
        finally:
            # Cleanup
            store.delete_run(run_id)
    
    def test_search_experiments(self, store, experiment_id):
        """Test searching experiments."""
        # Search all experiments
        experiments = store.search_experiments(view_type=ViewType.ALL)
        assert isinstance(experiments, list)
        assert len(experiments) > 0
        
        # Search active experiments only
        active_experiments = store.search_experiments(view_type=ViewType.ACTIVE_ONLY)
        assert isinstance(active_experiments, list)
        
        # Verify our test experiment is in the results
        experiment_ids = [exp.experiment_id for exp in active_experiments]
        assert experiment_id in experiment_ids
    
    def test_search_runs(self, store, experiment_id):
        """Test searching runs."""
        # Create a test run
        run = store.create_run(experiment_id=experiment_id)
        run_id = run.info.run_id
        
        try:
            # Search runs in the experiment
            runs = store.search_runs(
                experiment_ids=[experiment_id],
                run_view_type=ViewType.ACTIVE_ONLY
            )
            assert isinstance(runs, list)
            assert len(runs) > 0
            
            # Verify our test run is in the results
            run_ids = [r.info.run_id for r in runs]
            assert run_id in run_ids
            
        finally:
            # Cleanup
            store.delete_run(run_id)
    
    def test_experiment_tags(self, store, experiment_id):
        """Test setting and managing experiment tags."""
        # Set experiment tag
        tag = ExperimentTag(key="test_tag", value="test_value")
        store.set_experiment_tag(experiment_id, tag)
        
        # Get experiment and verify tag
        experiment = store.get_experiment(experiment_id)
        tag_dict = {t.key: t.value for t in experiment.tags}
        assert "test_tag" in tag_dict
        assert tag_dict["test_tag"] == "test_value"
    
    def test_run_tags(self, store, experiment_id):
        """Test setting and managing run tags."""
        # Create run
        run = store.create_run(experiment_id=experiment_id)
        run_id = run.info.run_id
        
        try:
            # Set run tag
            tag = RunTag(key="run_test_tag", value="run_test_value")
            store.set_tag(run_id, tag)
            
            # Get run and verify tag
            retrieved_run = store.get_run(run_id)
            tag_dict = {t.key: t.value for t in retrieved_run.data.tags}
            assert "run_test_tag" in tag_dict
            assert tag_dict["run_test_tag"] == "run_test_value"
            
            # Delete tag
            store.delete_tag(run_id, "run_test_tag")
            
            # Verify tag is deleted
            retrieved_run = store.get_run(run_id)
            tag_dict = {t.key: t.value for t in retrieved_run.data.tags}
            assert "run_test_tag" not in tag_dict
            
        finally:
            # Cleanup
            store.delete_run(run_id)
    
    def test_experiment_lifecycle(self, store, experiment_name):
        """Test experiment lifecycle operations."""
        # Create experiment
        experiment_id = store.create_experiment(experiment_name)
        
        try:
            # Verify experiment is active
            experiment = store.get_experiment(experiment_id)
            assert experiment.lifecycle_stage == LifecycleStage.ACTIVE
            
            # Delete experiment
            store.delete_experiment(experiment_id)
            
            # Verify experiment is deleted
            experiment = store.get_experiment(experiment_id)
            assert experiment.lifecycle_stage == LifecycleStage.DELETED
            
            # Restore experiment
            store.restore_experiment(experiment_id)
            
            # Verify experiment is active again
            experiment = store.get_experiment(experiment_id)
            assert experiment.lifecycle_stage == LifecycleStage.ACTIVE
            
        finally:
            # Cleanup
            try:
                store.delete_experiment(experiment_id)
            except:
                pass
    
    def test_run_lifecycle(self, store, experiment_id):
        """Test run lifecycle operations."""
        # Create run
        run = store.create_run(experiment_id=experiment_id)
        run_id = run.info.run_id
        
        try:
            # Verify run is active
            run = store.get_run(run_id)
            assert run.info.lifecycle_stage == LifecycleStage.ACTIVE
            
            # Delete run
            store.delete_run(run_id)
            
            # Verify run is deleted
            run = store.get_run(run_id)
            assert run.info.lifecycle_stage == LifecycleStage.DELETED
            
            # Restore run
            store.restore_run(run_id)
            
            # Verify run is active again
            run = store.get_run(run_id)
            assert run.info.lifecycle_stage == LifecycleStage.ACTIVE
            
        finally:
            # Cleanup
            try:
                store.delete_run(run_id)
            except:
                pass


def test_mlflow_integration():
    """Test that the ModelRegistryStore works with MLflow's tracking API."""
    host = os.getenv("MODEL_REGISTRY_HOST")
    port = os.getenv("MODEL_REGISTRY_PORT", "8080")
    secure = os.getenv("MODEL_REGISTRY_SECURE", "false").lower() == "true"
    token = os.getenv("MODEL_REGISTRY_TOKEN")
    
    if not host or not token:
        pytest.skip("MODEL_REGISTRY_HOST or MODEL_REGISTRY_TOKEN not set")
    
    # Set the token in environment for auth
    os.environ["MODEL_REGISTRY_TOKEN"] = token
    
    # Set tracking URI
    scheme = "https" if secure else "http"
    tracking_uri = f"modelregistry://{host}:{port}"
    mlflow.set_tracking_uri(tracking_uri)
    
    # Test MLflow integration
    experiment_name = f"mlflow-e2e-test-{uuid.uuid4().hex[:8]}"
    
    try:
        # Create experiment using MLflow
        experiment_id = mlflow.create_experiment(experiment_name)
        assert experiment_id is not None
        
        # Start a run using MLflow
        with mlflow.start_run(experiment_id=experiment_id) as run:
            # Log parameters
            mlflow.log_param("mlflow_param", "mlflow_value")
            
            # Log metrics
            mlflow.log_metric("mlflow_metric", 0.99)
            
            # Log tags
            mlflow.set_tag("mlflow_tag", "mlflow_tag_value")
            
            run_id = run.info.run_id
        
        # Verify the run was created
        retrieved_run = mlflow.get_run(run_id)
        assert retrieved_run.info.run_id == run_id
        assert retrieved_run.data.params["mlflow_param"] == "mlflow_value"
        assert retrieved_run.data.metrics["mlflow_metric"] == 0.99
        assert retrieved_run.data.tags["mlflow_tag"] == "mlflow_tag_value"
        
    finally:
        # Cleanup
        try:
            mlflow.delete_experiment(experiment_id)
        except:
            pass


if __name__ == "__main__":
    # Allow running the tests directly
    pytest.main([__file__, "-v"]) 