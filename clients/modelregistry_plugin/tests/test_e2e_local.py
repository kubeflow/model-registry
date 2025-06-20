"""
End-to-end tests for ModelRegistryStore with local Model Registry server.

This test suite starts a local Model Registry server with MLMD backend using SQLite
and runs comprehensive tests against it. This provides a self-contained testing
environment that doesn't require external dependencies.
"""

import os
import pytest
import time
import uuid
import subprocess
import tempfile
import shutil
import signal
import threading
from pathlib import Path
from unittest.mock import patch

import mlflow
from mlflow.entities import (
    Experiment, Run, RunInfo, RunData, RunStatus, RunTag, Param, Metric,
    ViewType, LifecycleStage, ExperimentTag
)

from modelregistry_plugin.store import ModelRegistryStore


class LocalModelRegistryServer:
    """Manages a local Model Registry server with MLMD backend."""
    
    def __init__(self, temp_dir: Path):
        self.temp_dir = temp_dir
        self.mlmd_server_process = None
        self.model_registry_process = None
        self.mlmd_port = 9090
        self.model_registry_port = 8080
        self.mlmd_db_path = temp_dir / "metadata.sqlite.db"
        self.conn_config_path = temp_dir / "conn_config.pb"
        
    def setup_mlmd_config(self):
        """Create MLMD connection configuration for SQLite."""
        # Ensure the temp directory exists and is world-writable
        self.temp_dir.mkdir(parents=True, exist_ok=True)
        os.chmod(self.temp_dir, 0o777)
        
        db_path = self.temp_dir / "metadata.sqlite.db"
        conn_config_content = f"""connection_config {{
  sqlite {{
    filename_uri: '{db_path}'
    connection_mode: READWRITE_OPENCREATE
  }}
}}
"""
        with open(self.conn_config_path, 'w') as f:
            f.write(conn_config_content)
        
        # Ensure the config file has proper permissions
        os.chmod(self.conn_config_path, 0o644)
        
        print(f"📁 Created MLMD config at: {self.conn_config_path}")
        print(f"📁 Database will be created at: {db_path}")
    
    def start_mlmd_server(self):
        """Start the MLMD server using Docker."""
        try:
            # Use absolute path for Docker volume mount
            temp_dir_abs = str(self.temp_dir.absolute())
            
            print(f"🐳 Starting MLMD server with volume mount: {temp_dir_abs}:/tmp/shared")
            
            # Start MLMD server
            self.mlmd_server_process = subprocess.Popen([
                'docker', 'run', '--rm',
                '-p', f'{self.mlmd_port}:8080',
                '-v', f'{temp_dir_abs}:/tmp/shared',
                '-e', 'METADATA_STORE_SERVER_CONFIG_FILE=/tmp/shared/conn_config.pb',
                'gcr.io/tfx-oss-public/ml_metadata_store_server:1.14.0'
            ], stdout=subprocess.PIPE, stderr=subprocess.PIPE)
            
            # Wait for MLMD server to start
            time.sleep(5)
            
            # Check if MLMD server is running
            if self.mlmd_server_process.poll() is not None:
                stdout, stderr = self.mlmd_server_process.communicate()
                print(f"MLMD stdout: {stdout.decode()}")
                print(f"MLMD stderr: {stderr.decode()}")
                raise RuntimeError(f"MLMD server failed to start: {stderr.decode()}")
                
            print(f"✅ MLMD server started on port {self.mlmd_port}")
            
        except Exception as e:
            print(f"❌ Failed to start MLMD server: {e}")
            raise
    
    def start_model_registry_server(self):
        """Start the Model Registry server."""
        try:
            # Get the project root (two levels up from the plugin directory)
            project_root = Path(__file__).parent.parent.parent.parent
            
            # Build the Model Registry server
            print("🔨 Building Model Registry server...")
            build_result = subprocess.run(
                ['make', 'build'],
                cwd=project_root,
                capture_output=True,
                text=True
            )
            
            if build_result.returncode != 0:
                print(f"Build stdout: {build_result.stdout}")
                print(f"Build stderr: {build_result.stderr}")
                raise RuntimeError(f"Failed to build Model Registry: {build_result.stderr}")
            
            # Start Model Registry server
            model_registry_bin = project_root / "model-registry"
            if not model_registry_bin.exists():
                raise RuntimeError("Model Registry binary not found after build")
            
            print(f"🚀 Starting Model Registry server, connecting to MLMD on localhost:{self.mlmd_port}")
            
            self.model_registry_process = subprocess.Popen([
                str(model_registry_bin),
                'proxy',
                '--hostname', '0.0.0.0',
                '--port', str(self.model_registry_port),
                '--mlmd-hostname', 'localhost',
                '--mlmd-port', str(self.mlmd_port),  # Use the MLMD port (9090)
                '--datastore-type', 'mlmd'
            ], stdout=subprocess.PIPE, stderr=subprocess.PIPE)
            
            # Wait for Model Registry server to start
            time.sleep(10)
            
            # Check if Model Registry server is running
            if self.model_registry_process.poll() is not None:
                stdout, stderr = self.model_registry_process.communicate()
                print(f"Model Registry stdout: {stdout.decode()}")
                print(f"Model Registry stderr: {stderr.decode()}")
                raise RuntimeError(f"Model Registry server failed to start: {stderr.decode()}")
                
            print(f"✅ Model Registry server started on port {self.model_registry_port}")
            
        except Exception as e:
            print(f"❌ Failed to start Model Registry server: {e}")
            raise
    
    def start(self):
        """Start both MLMD and Model Registry servers."""
        print("🚀 Starting local Model Registry test environment...")
        
        # Setup MLMD configuration
        self.setup_mlmd_config()
        
        # Start servers
        self.start_mlmd_server()
        self.start_model_registry_server()
        
        # Wait a bit more for everything to be ready
        time.sleep(5)
        
        print("✅ Local Model Registry test environment ready!")
    
    def stop(self):
        """Stop both servers and clean up."""
        print("🛑 Stopping local Model Registry test environment...")
        
        if self.model_registry_process:
            self.model_registry_process.terminate()
            try:
                self.model_registry_process.wait(timeout=10)
            except subprocess.TimeoutExpired:
                self.model_registry_process.kill()
        
        if self.mlmd_server_process:
            self.mlmd_server_process.terminate()
            try:
                self.mlmd_server_process.wait(timeout=10)
            except subprocess.TimeoutExpired:
                self.mlmd_server_process.kill()
        
        # Clean up Docker containers
        try:
            subprocess.run(['docker', 'ps', '-q', '--filter', 'ancestor=gcr.io/tfx-oss-public/ml_metadata_store_server:1.14.0'], 
                         capture_output=True, text=True, check=True)
            subprocess.run(['docker', 'stop', '$(docker ps -q --filter ancestor=gcr.io/tfx-oss-public/ml_metadata_store_server:1.14.0)'], 
                         shell=True, capture_output=True)
        except:
            pass  # Ignore cleanup errors
        
        print("✅ Local Model Registry test environment stopped.")


class TestModelRegistryStoreE2ELocal:
    """End-to-end tests for ModelRegistryStore with local Model Registry server."""
    
    @pytest.fixture(scope="class")
    def local_server(self):
        """Start and manage a local Model Registry server."""
        # Create temporary directory for test data under /tmp for Docker compatibility
        temp_dir = Path(tempfile.mkdtemp(dir='/tmp', prefix="model_registry_e2e_"))
        
        server = LocalModelRegistryServer(temp_dir)
        
        try:
            server.start()
            yield server
        finally:
            server.stop()
            # Clean up temporary directory
            shutil.rmtree(temp_dir, ignore_errors=True)
    
    @pytest.fixture(scope="class")
    def store(self, local_server):
        """Create a ModelRegistryStore instance connected to the local server."""
        store_uri = f"modelregistry://localhost:{local_server.model_registry_port}"
        return ModelRegistryStore(store_uri=store_uri)
    
    @pytest.fixture(scope="class")
    def experiment_name(self):
        """Generate a unique experiment name for testing."""
        return f"e2e-local-test-{uuid.uuid4().hex[:8]}"
    
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
    
    def test_local_server_connection(self, store):
        """Test that we can connect to the local Model Registry server."""
        # Try to list experiments to verify connection
        experiments = store.list_experiments()
        assert isinstance(experiments, list)
        print(f"✅ Successfully connected to local Model Registry. Found {len(experiments)} experiments.")
    
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


def test_mlflow_integration_local():
    """Test that the ModelRegistryStore works with MLflow's tracking API using local server."""
    # This test will be run separately since it needs its own server instance
    pytest.skip("This test requires a separate server instance")


if __name__ == "__main__":
    # Allow running the tests directly
    pytest.main([__file__, "-v"]) 