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
import sys
import psutil
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
        # Track process IDs for cleanup
        self.mlmd_pid = None
        self.model_registry_pid = None
        
    def setup_mlmd_config(self):
        """Create MLMD connection configuration for SQLite."""
        # Ensure the temp directory exists and is world-writable
        self.temp_dir.mkdir(parents=True, exist_ok=True)
        os.chmod(self.temp_dir, 0o777)
        
        db_path = self.temp_dir / "metadata.sqlite.db"
        conn_config_content = f"""connection_config {{
  sqlite {{
    filename_uri: '/tmp/shared/metadata.sqlite.db'
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
            
            # Store the process ID
            self.mlmd_pid = self.mlmd_server_process.pid
            
            # Wait for MLMD server to start
            time.sleep(5)
            
            # Check if MLMD server is running
            if self.mlmd_server_process.poll() is not None:
                stdout, stderr = self.mlmd_server_process.communicate()
                print(f"MLMD stdout: {stdout.decode()}")
                print(f"MLMD stderr: {stderr.decode()}")
                raise RuntimeError(f"MLMD server failed to start: {stderr.decode()}")
                
            print(f"✅ MLMD server started on port {self.mlmd_port} (PID: {self.mlmd_pid})")
            
        except Exception as e:
            print(f"❌ Failed to start MLMD server: {e}")
            raise
    
    def start_model_registry_server(self):
        """Start the Model Registry server."""
        try:
            # Get the project root (two levels up from the plugin directory)
            project_root = Path(__file__).parent.parent.parent.parent
            
            print(f"🚀 Starting Model Registry server with Go, connecting to MLMD on localhost:{self.mlmd_port}")
            
            # Start Model Registry server using Go command
            self.model_registry_process = subprocess.Popen([
                'go', 'run', 'main.go',
                'proxy', '0.0.0.0',
                '--hostname', '0.0.0.0',
                '--port', str(self.model_registry_port),
                '--mlmd-hostname', 'localhost',
                '--mlmd-port', str(self.mlmd_port),  # Use the MLMD port (9090)
                '--datastore-type', 'mlmd'
            ], cwd=project_root, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
            
            # Store the process ID
            self.model_registry_pid = self.model_registry_process.pid
            
            # Wait for Model Registry server to start
            time.sleep(10)
            
            # Check if Model Registry server is running
            if self.model_registry_process.poll() is not None:
                stdout, stderr = self.model_registry_process.communicate()
                print(f"Model Registry stdout: {stdout.decode()}")
                print(f"Model Registry stderr: {stderr.decode()}")
                raise RuntimeError(f"Model Registry server failed to start: {stderr.decode()}")
                
            print(f"✅ Model Registry server started on port {self.model_registry_port} (PID: {self.model_registry_pid})")
            
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
        
        # Stop Model Registry server
        if self.model_registry_process and self.model_registry_pid:
            try:
                print(f"  🛑 Stopping Model Registry server (PID: {self.model_registry_pid})...")
                # Terminate child process if it exists
                try:
                    parent = psutil.Process(self.model_registry_pid)
                    children = parent.children()
                    if children:
                        child = children[0]  # Get the first (and likely only) child
                        print(f"    🛑 Terminating child process: {child.pid} ({child.name()})")
                        child.terminate()
                        try:
                            child.wait(timeout=5)
                            print(f"    ✅ Child process {child.pid} terminated gracefully")
                        except psutil.TimeoutExpired:
                            print(f"    ⚠️  Child process {child.pid} didn't terminate gracefully, killing...")
                            child.kill()
                            child.wait(timeout=2)
                            print(f"    ✅ Child process {child.pid} killed")
                        except psutil.NoSuchProcess:
                            print(f"    ✅ Child process {child.pid} already terminated")
                except psutil.NoSuchProcess:
                    print(f"    ℹ️  Parent process {self.model_registry_pid} not found")
                except Exception as e:
                    print(f"    ⚠️  Error terminating child process: {e}")
                
                # Terminate parent process
                self.model_registry_process.terminate()
                try:
                    self.model_registry_process.wait(timeout=10)
                    print("  ✅ Model Registry server stopped gracefully")
                except subprocess.TimeoutExpired:
                    print("  ⚠️  Model Registry server didn't stop gracefully, killing...")
                    self.model_registry_process.kill()
                    try:
                        self.model_registry_process.wait(timeout=5)
                        print("  ✅ Model Registry server killed")
                    except subprocess.TimeoutExpired:
                        print("  ❌ Failed to kill Model Registry server")
            except Exception as e:
                print(f"  ❌ Error stopping Model Registry server: {e}")
        
        # Stop MLMD server
        if self.mlmd_server_process and self.mlmd_pid:
            try:
                print(f"  🛑 Stopping MLMD server (PID: {self.mlmd_pid})...")
                self.mlmd_server_process.terminate()
                try:
                    self.mlmd_server_process.wait(timeout=10)
                    print("  ✅ MLMD server stopped gracefully")
                except subprocess.TimeoutExpired:
                    print("  ⚠️  MLMD server didn't stop gracefully, killing...")
                    self.mlmd_server_process.kill()
                    try:
                        self.mlmd_server_process.wait(timeout=5)
                        print("  ✅ MLMD server killed")
                    except subprocess.TimeoutExpired:
                        print("  ❌ Failed to kill MLMD server")
            except Exception as e:
                print(f"  ❌ Error stopping MLMD server: {e}")
        
        # Clean up Docker containers (MLMD server)
        try:
            print("  🐳 Cleaning up Docker containers...")
            # Find running MLMD containers
            result = subprocess.run(
                ['docker', 'ps', '-q', '--filter', 'ancestor=gcr.io/tfx-oss-public/ml_metadata_store_server:1.14.0'], 
                capture_output=True, text=True, check=False
            )
            
            if result.returncode == 0 and result.stdout.strip():
                container_ids = result.stdout.strip().split('\n')
                for container_id in container_ids:
                    if container_id.strip():
                        print(f"    🛑 Stopping Docker container: {container_id}")
                        subprocess.run(
                            ['docker', 'stop', container_id], 
                            capture_output=True, text=True, check=False
                        )
                        print(f"    🗑️  Removing Docker container: {container_id}")
                        subprocess.run(
                            ['docker', 'rm', container_id], 
                            capture_output=True, text=True, check=False
                        )
                print("  ✅ Docker containers cleaned up")
            else:
                print("  ℹ️  No running MLMD Docker containers found")
                
        except Exception as e:
            print(f"  ❌ Error cleaning up Docker containers: {e}")
        
        # Reset process references and PIDs
        self.model_registry_process = None
        self.mlmd_server_process = None
        self.model_registry_pid = None
        self.mlmd_pid = None
        
        print("✅ Local Model Registry test environment stopped.")


class TestModelRegistryStoreE2ELocal:
    """End-to-end tests for ModelRegistryStore with local Model Registry server."""
    
    # Class-level storage for cleanup
    _local_server = None
    _class_store = None
    
    @classmethod
    def setup_class(cls):
        """Setup signal handlers for graceful cleanup."""
        def signal_handler(signum, frame):
            print(f"\n🛑 Received signal {signum}, cleaning up...")
            if cls._local_server:
                try:
                    cls._local_server.stop()
                except Exception as e:
                    print(f"❌ Error during signal cleanup: {e}")
            sys.exit(1)
        
        # Register signal handlers for graceful cleanup
        signal.signal(signal.SIGINT, signal_handler)
        signal.signal(signal.SIGTERM, signal_handler)
    
    @classmethod
    def teardown_class(cls):
        """Ensure cleanup happens at class teardown."""
        print("🧹 Cleaning up class-level resources...")
        
        # Clean up local server
        if cls._local_server:
            try:
                print(f"🛑 Stopping local server with tracked PIDs...")
                if cls._local_server.model_registry_pid:
                    print(f"  📋 Model Registry PID: {cls._local_server.model_registry_pid}")
                if cls._local_server.mlmd_pid:
                    print(f"  📋 MLMD PID: {cls._local_server.mlmd_pid}")
                
                cls._local_server.stop()
                print("  ✅ Successfully stopped local server")
            except Exception as e:
                print(f"❌ Error during class teardown cleanup: {e}")
        else:
            print("  ℹ️  No local server to clean up")
        
        # Clear class reference
        cls._local_server = None
        
        print("✅ Class-level cleanup completed.")
    
    @pytest.fixture(scope="class")
    def local_server(self):
        """Start and manage a local Model Registry server."""
        # Create temporary directory for test data under /tmp for Docker compatibility
        # TODO switch to a temp directory in the module directory
        temp_dir = Path(tempfile.mkdtemp(dir='/tmp', prefix="model_registry_e2e_"))
        
        server = LocalModelRegistryServer(temp_dir)
        
        try:
            server.start()
            # Store reference for class-level cleanup
            TestModelRegistryStoreE2ELocal._local_server = server
            yield server
        except Exception as e:
            print(f"❌ Failed to start local server: {e}")
            try:
                server.stop()
            except Exception as cleanup_error:
                print(f"❌ Error during cleanup after startup failure: {cleanup_error}")
            raise
        finally:
            try:
                server.stop()
            except Exception as e:
                print(f"❌ Error during final cleanup: {e}")
            try:
                shutil.rmtree(temp_dir, ignore_errors=True)
            except Exception as e:
                print(f"❌ Error cleaning up temp directory: {e}")
            TestModelRegistryStoreE2ELocal._local_server = None
    
    @pytest.fixture(scope="class")
    def store(self, local_server):
        """Create a ModelRegistryStore instance connected to the local server."""
        store_uri = f"modelregistry://localhost:{local_server.model_registry_port}"
        return ModelRegistryStore(store_uri=store_uri)
    
    @pytest.fixture
    def experiment_name(self):
        """Generate a unique experiment name for testing."""
        return f"e2e-local-test-{uuid.uuid4().hex[:8]}"
    
    @pytest.fixture
    def experiment_id(self, store, experiment_name):
        """Create a test experiment and return its ID.
        
        This experiment will be cleaned up after each test that uses it.
        """
        experiment_id = None
        try:
            experiment_id = store.create_experiment(experiment_name)
            print(f"📁 Created test experiment: {experiment_id} ({experiment_name})")
            yield experiment_id
        except Exception as e:
            print(f"❌ Failed to create experiment '{experiment_name}': {e}")
            raise
        finally:
            # Cleanup: delete the experiment
            if experiment_id:
                try:
                    store.delete_experiment(experiment_id)
                    print(f"✅ Cleaned up test experiment: {experiment_id}")
                except Exception as e:
                    print(f"❌ Error deleting test experiment {experiment_id}: {e}")
                    # Fail the test if cleanup fails - this could indicate resource leaks
                    pytest.fail(f"Failed to clean up experiment {experiment_id}: {e}")
            else:
                print("⚠️  No test experiment to clean up (creation failed)")
    
    @pytest.fixture
    def run_id(self, store, experiment_id):
        """Create a test run and return its ID.
        
        This run will be cleaned up after each test that uses it.
        """
        run_id = None
        try:
            run = store.create_run(
                experiment_id=experiment_id,
                user_id="test-user",
                run_name="test-run"
            )
            run_id = run.info.run_id
            print(f"🏃 Created test run: {run_id}")
            yield run_id
        except Exception as e:
            print(f"❌ Failed to create run in experiment {experiment_id}: {e}")
            raise
        finally:
            # Cleanup: delete the run
            if run_id:
                try:
                    store.delete_run(run_id)
                    print(f"✅ Cleaned up test run: {run_id}")
                except Exception as e:
                    print(f"❌ Error deleting test run {run_id}: {e}")
                    # Fail the test if cleanup fails - this could indicate resource leaks
                    pytest.fail(f"Failed to clean up run {run_id}: {e}")
            else:
                print("⚠️  No test run to clean up (creation failed)")
    
    def test_local_server_connection(self, store):
        """Test that we can connect to the local Model Registry server."""
        experiments = store.list_experiments()
        assert isinstance(experiments, list)
        print(f"✅ Successfully connected to local Model Registry. Found {len(experiments)} experiments.")
    
    def test_experiment_exists(self, store, experiment_id):
        """Test that the experiment exists and can be retrieved."""
        experiment = store.get_experiment(experiment_id)
        assert isinstance(experiment, Experiment)
        assert experiment.experiment_id == experiment_id
        assert experiment.lifecycle_stage == LifecycleStage.ACTIVE
        print(f"✅ Experiment exists: {experiment.name}")
    
    def test_run_exists(self, store, run_id):
        """Test that the run exists and can be retrieved."""
        retrieved_run = store.get_run(run_id)
        assert isinstance(retrieved_run, Run)
        assert retrieved_run.info.run_id == run_id
        assert retrieved_run.info.status == RunStatus.RUNNING
        assert retrieved_run.info.lifecycle_stage == LifecycleStage.ACTIVE
        print(f"✅ Run exists: {retrieved_run.info.run_name}")
    
    def test_run_logging(self, store, run_id):
        """Test logging to the run."""
        # Log parameters
        param1 = Param(key="learning_rate", value="0.001")
        param2 = Param(key="epochs", value="50")
        store.log_param(run_id, param1)
        store.log_param(run_id, param2)
        
        # Log metrics
        metric1 = Metric(key="accuracy", value=0.98, timestamp=int(time.time() * 1000), step=1)
        metric2 = Metric(key="loss", value=0.02, timestamp=int(time.time() * 1000), step=1)
        store.log_metric(run_id, metric1)
        store.log_metric(run_id, metric2)
        
        # Verify logged data
        retrieved_run = store.get_run(run_id)
        assert len(retrieved_run.data.params) >= 2
        assert len(retrieved_run.data.metrics) >= 2
        
        # Verify parameters
        assert retrieved_run.data.params["learning_rate"] == "0.001"
        assert retrieved_run.data.params["epochs"] == "50"
        
        # Verify metrics
        assert retrieved_run.data.metrics["accuracy"] == 0.98
        assert retrieved_run.data.metrics["loss"] == 0.02
        
        print(f"✅ Successfully logged data to run: {run_id}")
    
    def test_run_batch_logging(self, store, run_id):
        """Test batch logging to the run."""
        # Prepare batch data
        metrics = [
            Metric(key="batch_metric1", value=1.5, timestamp=int(time.time() * 1000), step=2),
            Metric(key="batch_metric2", value=2.5, timestamp=int(time.time() * 1000), step=2)
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
        assert len(retrieved_run.data.metrics) >= 2
        assert len(retrieved_run.data.params) >= 2
        assert len(retrieved_run.data.tags) >= 2
        
        print(f"✅ Successfully batch logged data to run: {run_id}")
    
    def test_run_tags(self, store, run_id):
        """Test setting and managing tags on the run."""
        # Set run tag
        tag = RunTag(key="test_tag", value="test_value")
        store.set_tag(run_id, tag)
        
        # Get run and verify tag
        retrieved_run = store.get_run(run_id)
        assert retrieved_run.data.tags.get("test_tag") == "test_value"
        
        # Delete tag
        store.delete_tag(run_id, "test_tag")
        
        # Verify tag is deleted
        retrieved_run = store.get_run(run_id)
        assert retrieved_run.data.tags.get("test_tag") is None
        
        print(f"✅ Successfully managed tags on run: {run_id}")
    
    def test_run_lifecycle(self, store, run_id):
        """Test run lifecycle operations."""
        # Verify run is active
        run = store.get_run(run_id)
        assert run.info.lifecycle_stage == LifecycleStage.ACTIVE
        
        # Update run status to finished
        updated_info = store.update_run_info(
            run_id=run_id,
            run_status=RunStatus.FINISHED,
            end_time=int(time.time() * 1000)
        )
        assert updated_info.status == RunStatus.FINISHED
        
        # Verify run is still active (not deleted)
        run = store.get_run(run_id)
        assert run.info.lifecycle_stage == LifecycleStage.ACTIVE
        
        print(f"✅ Successfully tested lifecycle operations on run: {run_id}")
    
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
            assert retrieved_run.data.params["learning_rate"] == "0.01"
            assert retrieved_run.data.params["epochs"] == "100"
            
            # Verify metrics
            assert retrieved_run.data.metrics["accuracy"] == 0.95
            assert retrieved_run.data.metrics["loss"] == 0.05
            
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
        assert experiment.tags.get("test_tag") == "test_value"
    
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
            assert retrieved_run.data.tags.get("run_test_tag") == "run_test_value"
            
            # Delete tag
            store.delete_tag(run_id, "run_test_tag")
            
            # Verify tag is deleted
            retrieved_run = store.get_run(run_id)
            assert retrieved_run.data.tags.get("run_test_tag") is None
            
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
                print(f"✅ Cleaned up experiment from lifecycle test: {experiment_id}")
            except Exception as e:
                print(f"❌ Error deleting experiment {experiment_id} in lifecycle test: {e}")
                # Fail the test if cleanup fails - this could indicate resource leaks
                pytest.fail(f"Failed to clean up experiment {experiment_id} in lifecycle test: {e}")
    
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
                print(f"✅ Cleaned up run from lifecycle test: {run_id}")
            except Exception as e:
                print(f"❌ Error deleting run {run_id} in lifecycle test: {e}")
                # Fail the test if cleanup fails - this could indicate resource leaks
                pytest.fail(f"Failed to clean up run {run_id} in lifecycle test: {e}")


if __name__ == "__main__":
    # Allow running the tests directly
    pytest.main([__file__, "-v"]) 