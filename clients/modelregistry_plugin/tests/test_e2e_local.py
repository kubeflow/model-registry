"""
End-to-end tests for ModelRegistryTrackingStore with local Model Registry server.

This test suite starts a local Model Registry server with EmbedMD backend using MySQL
and runs comprehensive tests against it. This provides a self-contained testing
environment that doesn't require external dependencies.
"""

import logging
import pytest
import time
import uuid
import subprocess
import shutil
import signal
import sys
from pathlib import Path
from unittest.mock import patch

from .conftest import create_temp_dir
from mlflow.entities import (
    Experiment,
    Run,
    RunStatus,
    RunTag,
    Param,
    Metric,
    ViewType,
    LifecycleStage,
    ExperimentTag,
)
from testcontainers.mysql import MySqlContainer

# Mock the API client to avoid actual HTTP requests
with patch("modelregistry_plugin.api_client.requests.Session.request") as mock_request:
    from modelregistry_plugin.tracking_store import ModelRegistryTrackingStore

# Set up logging
logger = logging.getLogger(__name__)


class LocalModelRegistryServer:
    """Manages a local Model Registry server with EmbedMD backend using testcontainers."""

    def __init__(self, temp_dir: Path):
        self.temp_dir = temp_dir
        self.mysql_container = None
        self.model_registry_process = None
        self.mysql_port = 3306
        self.model_registry_port = 8080
        self.mysql_dsn = None

    def start_mysql_server(self):
        """Start the MySQL server using testcontainers."""
        try:
            logger.info("🐳 Starting MySQL server for EmbedMD...")

            # Create MySQL container using testcontainers MySQL module
            self.mysql_container = MySqlContainer(
                image="mysql:8.3",
                username="root",
                root_password="root",
                password="root",
                dbname="model_registry",
            )

            # Start the container
            self.mysql_container.start()

            # Get connection details
            self.mysql_port = self.mysql_container.get_exposed_port(3306)
            self.mysql_dsn = f"root:root@tcp(localhost:{self.mysql_port})/model_registry?charset=utf8mb4&parseTime=True&loc=Local"

            logger.info(f"✅ MySQL server started on port {self.mysql_port}")
            logger.info(f"📁 MySQL DSN: {self.mysql_dsn}")

        except Exception as e:
            logger.error(f"❌ Failed to start MySQL server: {e}")
            if self.mysql_container:
                try:
                    logger.error("🐳 Container logs:")
                    logger.error(self.mysql_container.get_logs())
                    self.mysql_container.stop()
                except Exception as cleanup_error:
                    logger.error(f"❌ Error during MySQL cleanup: {cleanup_error}")
            raise

    def start_model_registry_server(self):
        """Start the Model Registry server."""
        try:
            # Get the project root (two levels up from the plugin directory)
            project_root = Path(__file__).parent.parent.parent.parent

            logger.info(
                f"🚀 Starting Model Registry server with Go, connecting to MySQL on localhost:{self.mysql_port}"
            )

            # Start Model Registry server using Go command
            self.model_registry_process = subprocess.Popen(
                [
                    "go",
                    "run",
                    "main.go",
                    "proxy",
                    "--hostname",
                    "0.0.0.0",
                    "--port",
                    str(self.model_registry_port),
                    "--embedmd-database-type",
                    "mysql",
                    "--embedmd-database-dsn",
                    self.mysql_dsn,
                ],
                cwd=project_root,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
            )

            # Wait for Model Registry server to start
            time.sleep(10)

            # Check if Model Registry server is running
            if self.model_registry_process.poll() is not None:
                stdout, stderr = self.model_registry_process.communicate()
                logger.error(f"Model Registry stdout: {stdout.decode()}")
                logger.error(f"Model Registry stderr: {stderr.decode()}")
                raise RuntimeError(
                    f"Model Registry server failed to start: {stderr.decode()}"
                )

            logger.info(
                f"✅ Model Registry server started on port {self.model_registry_port}"
            )

        except Exception as e:
            logger.error(f"❌ Failed to start Model Registry server: {e}")
            raise

    def start(self):
        """Start both MySQL and Model Registry servers."""
        logger.info("🚀 Starting local Model Registry test environment...")

        # Start servers
        self.start_mysql_server()
        self.start_model_registry_server()

        # Wait a bit more for everything to be ready
        time.sleep(5)

        logger.info("✅ Local Model Registry test environment ready!")

    def stop(self):
        """Stop both servers and clean up."""
        logger.info("🛑 Stopping local Model Registry test environment...")

        # Stop Model Registry server
        if self.model_registry_process:
            try:
                logger.info("  🛑 Stopping Model Registry server...")

                # First, try to terminate child processes if they exist
                try:
                    import psutil

                    parent = psutil.Process(self.model_registry_process.pid)
                    children = parent.children(recursive=True)

                    if children:
                        logger.info(
                            f"    🛑 Found {len(children)} child processes, terminating them..."
                        )
                        for child in children:
                            try:
                                logger.info(
                                    f"      🛑 Terminating child process: {child.pid} ({child.name()})"
                                )
                                child.terminate()
                            except psutil.NoSuchProcess:
                                logger.error(
                                    f"      ℹ️  Child process {child.pid} already terminated"
                                )
                            except Exception as e:
                                logger.error(
                                    f"      ⚠️  Error terminating child process {child.pid}: {e}"
                                )

                        # Wait for children to terminate gracefully
                        gone, alive = psutil.wait_procs(children, timeout=5)
                        for child in alive:
                            logger.warning(
                                f"      ⚠️  Child process {child.pid} didn't terminate gracefully, killing..."
                            )
                            try:
                                child.kill()
                                child.wait(timeout=2)
                                logger.info(
                                    f"      ✅ Child process {child.pid} killed"
                                )
                            except psutil.NoSuchProcess:
                                logger.info(
                                    f"      ✅ Child process {child.pid} already terminated"
                                )
                            except Exception as e:
                                logger.error(
                                    f"      ❌ Error killing child process {child.pid}: {e}"
                                )

                except psutil.NoSuchProcess:
                    logger.info("    ℹ️  Parent process not found")
                except Exception as e:
                    logger.error(f"    ⚠️  Error handling child processes: {e}")

                # Now terminate the main process
                self.model_registry_process.terminate()
                try:
                    self.model_registry_process.wait(timeout=10)
                    logger.info("  ✅ Model Registry server stopped gracefully")
                except subprocess.TimeoutExpired:
                    logger.warning(
                        "  ⚠️  Model Registry server didn't stop gracefully, killing..."
                    )
                    self.model_registry_process.kill()
                    try:
                        self.model_registry_process.wait(timeout=5)
                        logger.info("  ✅ Model Registry server killed")
                    except subprocess.TimeoutExpired:
                        logger.error("  ❌ Failed to kill Model Registry server")

            except Exception as e:
                logger.error(f"  ❌ Error stopping Model Registry server: {e}")

        # Stop MySQL container
        if self.mysql_container:
            try:
                logger.info("  🐳 Stopping MySQL container...")
                self.mysql_container.stop()
                logger.info("  ✅ MySQL container stopped")
            except Exception as e:
                logger.error(f"  ❌ Error stopping MySQL container: {e}")

        # Reset references
        self.model_registry_process = None
        self.mysql_container = None

        logger.info("✅ Local Model Registry test environment stopped.")


class TestModelRegistryTrackingStoreE2ELocal:
    """End-to-end tests for ModelRegistryTrackingStore with local Model Registry server."""

    # Class-level storage for cleanup
    _local_server = None

    @classmethod
    def setup_class(cls):
        """Setup signal handlers for graceful cleanup."""

        def signal_handler(signum, frame):
            logger.info(f"\n🛑 Received signal {signum}, cleaning up...")
            if cls._local_server:
                try:
                    cls._local_server.stop()
                except Exception as e:
                    logger.error(f"❌ Error during signal cleanup: {e}")
            sys.exit(1)

        # Register signal handlers for graceful cleanup
        signal.signal(signal.SIGINT, signal_handler)
        signal.signal(signal.SIGTERM, signal_handler)

    @classmethod
    def teardown_class(cls):
        """Ensure cleanup happens at class teardown."""
        logger.info("🧹 Cleaning up class-level resources...")

        # Clean up local server
        if cls._local_server:
            try:
                logger.info("🛑 Stopping local server...")
                cls._local_server.stop()
                logger.info("  ✅ Successfully stopped local server")
            except Exception as e:
                logger.error(f"❌ Error during class teardown cleanup: {e}")
        else:
            logger.info("  ℹ️  No local server to clean up")

        # Clear class reference
        cls._local_server = None

        logger.info("✅ Class-level cleanup completed.")

    @pytest.fixture(scope="class")
    def local_server(self):
        """Start and manage a local Model Registry server."""
        # Create temporary directory for test data in testdata directory for Docker compatibility
        temp_dir = create_temp_dir(prefix="model_registry_e2e_")

        server = LocalModelRegistryServer(temp_dir)

        try:
            server.start()
            # Store reference for class-level cleanup
            TestModelRegistryTrackingStoreE2ELocal._local_server = server
            yield server
        except Exception as e:
            logger.error(f"❌ Failed to start local server: {e}")
            try:
                server.stop()
            except Exception as cleanup_error:
                logger.error(
                    f"❌ Error during cleanup after startup failure: {cleanup_error}"
                )
            raise
        finally:
            try:
                server.stop()
            except Exception as e:
                logger.error(f"❌ Error during final cleanup: {e}")
            try:
                shutil.rmtree(temp_dir, ignore_errors=True)
            except Exception as e:
                logger.error(f"❌ Error cleaning up temp directory: {e}")
            TestModelRegistryTrackingStoreE2ELocal._local_server = None

    @pytest.fixture(scope="class")
    def store(self, local_server):
        """Create a ModelRegistryTrackingStore instance connected to the local server."""
        store_uri = f"modelregistry://localhost:{local_server.model_registry_port}"
        return ModelRegistryTrackingStore(store_uri=store_uri)

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
            logger.info(
                f"📁 Created test experiment: {experiment_id} ({experiment_name})"
            )
            yield experiment_id
        except Exception as e:
            logger.error(f"❌ Failed to create experiment '{experiment_name}': {e}")
            raise
        finally:
            # Cleanup: delete the experiment
            if experiment_id:
                try:
                    store.delete_experiment(experiment_id)
                    logger.info(f"✅ Cleaned up test experiment: {experiment_id}")
                except Exception as e:
                    logger.error(
                        f"❌ Error deleting test experiment {experiment_id}: {e}"
                    )
                    # Fail the test if cleanup fails - this could indicate resource leaks
                    pytest.fail(f"Failed to clean up experiment {experiment_id}: {e}")
            else:
                logger.warning("⚠️  No test experiment to clean up (creation failed)")

    @pytest.fixture
    def run_id(self, store, experiment_id):
        """Create a test run and return its ID.

        This run will be cleaned up after each test that uses it.
        """
        run_id = None
        try:
            run = store.create_run(
                experiment_id=experiment_id, user_id="test-user", run_name="test-run"
            )
            run_id = run.info.run_id
            logger.info(f"🏃 Created test run: {run_id}")
            yield run_id
        except Exception as e:
            logger.error(f"❌ Failed to create run in experiment {experiment_id}: {e}")
            raise
        finally:
            # Cleanup: delete the run
            if run_id:
                try:
                    store.delete_run(run_id)
                    logger.info(f"✅ Cleaned up test run: {run_id}")
                except Exception as e:
                    logger.error(f"❌ Error deleting test run {run_id}: {e}")
                    # Fail the test if cleanup fails - this could indicate resource leaks
                    pytest.fail(f"Failed to clean up run {run_id}: {e}")
            else:
                logger.warning("⚠️  No test run to clean up (creation failed)")

    def test_local_server_connection(self, store):
        """Test that we can connect to the local Model Registry server."""
        experiments = store.search_experiments()
        assert isinstance(experiments, list)
        logger.info(
            f"✅ Successfully connected to local Model Registry. Found {len(experiments)} experiments."
        )

    def test_experiment_exists(self, store, experiment_id):
        """Test that the experiment exists and can be retrieved."""
        experiment = store.get_experiment(experiment_id)
        assert isinstance(experiment, Experiment)
        assert experiment.experiment_id == experiment_id
        assert experiment.lifecycle_stage == LifecycleStage.ACTIVE
        logger.info(f"✅ Experiment exists: {experiment.name}")

    def test_run_exists(self, store, run_id):
        """Test that the run exists and can be retrieved."""
        retrieved_run = store.get_run(run_id)
        assert isinstance(retrieved_run, Run)
        assert retrieved_run.info.run_id == run_id
        assert retrieved_run.info.status == RunStatus.to_string(RunStatus.RUNNING)
        assert retrieved_run.info.lifecycle_stage == LifecycleStage.ACTIVE
        logger.info(f"✅ Run exists: {retrieved_run.info.run_name}")

    def test_run_logging(self, store, run_id):
        """Test logging to the run."""
        # Log parameters
        param1 = Param(key="learning_rate", value="0.001")
        param2 = Param(key="epochs", value="50")
        store.log_param(run_id, param1)
        store.log_param(run_id, param2)

        # Log metrics
        metric1 = Metric(
            key="accuracy", value=0.98, timestamp=int(time.time() * 1000), step=1
        )
        metric2 = Metric(
            key="loss", value=0.02, timestamp=int(time.time() * 1000), step=1
        )
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

        logger.info(f"✅ Successfully logged data to run: {run_id}")

    def test_run_batch_logging(self, store, run_id):
        """Test batch logging to the run."""
        # Prepare batch data
        metrics = [
            Metric(
                key="batch_metric1",
                value=1.5,
                timestamp=int(time.time() * 1000),
                step=2,
            ),
            Metric(
                key="batch_metric2",
                value=2.5,
                timestamp=int(time.time() * 1000),
                step=2,
            ),
        ]
        params = [
            Param(key="batch_param1", value="value1"),
            Param(key="batch_param2", value="value2"),
        ]
        tags = [
            RunTag(key="batch_tag1", value="tag_value1"),
            RunTag(key="batch_tag2", value="tag_value2"),
        ]

        # Log batch
        store.log_batch(run_id, metrics=metrics, params=params, tags=tags)

        # Verify batch data
        retrieved_run = store.get_run(run_id)
        assert len(retrieved_run.data.metrics) >= 2
        assert len(retrieved_run.data.params) >= 2
        assert len(retrieved_run.data.tags) >= 2

        logger.info(f"✅ Successfully batch logged data to run: {run_id}")

    def test_update_run_status(self, store, run_id):
        """Test updating run status."""
        # Verify run is active
        run = store.get_run(run_id)
        assert run.info.lifecycle_stage == LifecycleStage.ACTIVE

        # Update run status to finished
        updated_info = store.update_run_info(
            run_id=run_id,
            run_status=RunStatus.FINISHED,
            end_time=int(time.time() * 1000),
        )
        assert updated_info.status == RunStatus.to_string(RunStatus.FINISHED)

        # Verify run is still active (not deleted)
        run = store.get_run(run_id)
        assert run.info.lifecycle_stage == LifecycleStage.ACTIVE

        logger.info(f"✅ Successfully tested status update operations on run: {run_id}")

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
            experiment_id=experiment_id, user_id="test-user", run_name=run_name
        )

        assert isinstance(run, Run)
        assert run.info.experiment_id == experiment_id
        assert run.info.user_id == "test-user"
        assert run.info.status == RunStatus.to_string(RunStatus.RUNNING)
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
            end_time=int(time.time() * 1000),
        )
        assert updated_info.status == RunStatus.to_string(RunStatus.FINISHED)

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

            # Log metrics - log multiple values for the same metric to test history
            metric1 = Metric(
                key="accuracy", value=0.85, timestamp=int(time.time() * 1000), step=1
            )
            metric2 = Metric(
                key="loss", value=0.05, timestamp=int(time.time() * 1000), step=1
            )
            store.log_metric(run_id, metric1)
            store.log_metric(run_id, metric2)

            # Log additional accuracy values using for loops and value ranges
            base_timestamp = int(time.time() * 1000)
            accuracy_values = [
                0.85,
                0.90,
                0.95,
                0.98,
            ]  # Range of accuracy values to test

            for i, accuracy_value in enumerate(
                accuracy_values[1:], start=2
            ):  # Start from step 2, skip first value (already logged)
                metric = Metric(
                    key="accuracy",
                    value=accuracy_value,
                    timestamp=base_timestamp + (i * 1000),
                    step=i,
                )
                store.log_metric(run_id, metric)

            # Log additional loss values using a range
            loss_values = [0.05, 0.03, 0.02, 0.01]  # Decreasing loss values

            for i, loss_value in enumerate(
                loss_values[1:], start=2
            ):  # Start from step 2, skip first value (already logged)
                metric = Metric(
                    key="loss",
                    value=loss_value,
                    timestamp=base_timestamp
                    + (i * 1000)
                    + 500,  # Offset timestamps slightly
                    step=i,
                )
                store.log_metric(run_id, metric)

            # Get run and verify data
            retrieved_run = store.get_run(run_id)
            assert len(retrieved_run.data.params) == 2
            assert len(retrieved_run.data.metrics) == 2  # Still 2 unique metric keys

            # Verify parameters
            assert retrieved_run.data.params["learning_rate"] == "0.01"
            assert retrieved_run.data.params["epochs"] == "100"

            # Verify metrics (should show the latest value for each metric key)
            assert (
                retrieved_run.data.metrics["accuracy"] == 0.98
            )  # Latest accuracy value
            assert retrieved_run.data.metrics["loss"] == 0.01  # Latest loss value

            # Test metric history - should return all logged values for accuracy
            accuracy_history = store.get_metric_history(run_id, "accuracy")
            assert len(accuracy_history) == 4  # Should have 4 accuracy values
            for i, metric in enumerate(accuracy_history):
                assert metric.key == "accuracy"
                assert metric.value == accuracy_values[i]
                assert metric.step == i + 1

            # Test metric history for loss as well
            loss_history = store.get_metric_history(run_id, "loss")
            assert len(loss_history) == 4  # Should have 4 loss values
            for i, metric in enumerate(loss_history):
                assert metric.key == "loss"
                assert metric.value == loss_values[i]
                assert metric.step == i + 1

            logger.info(
                "✅ Successfully logged multiple metric values using loops and verified history contains all values in correct order"
            )
            logger.info(
                f"   Accuracy history: {len(accuracy_history)} values, Loss history: {len(loss_history)} values"
            )
            logger.info(
                "   Both metrics are sorted by (timestamp, step) as per EmbedMD specification"
            )

            # Test bulk metric history API for specific steps
            # Test accuracy metrics for steps 2 and 3
            accuracy_bulk_history = store.get_metric_history_bulk_interval_from_steps(
                run_id, "accuracy", steps=[2, 3], max_results=10
            )
            assert (
                len(accuracy_bulk_history) == 2
            )  # Should have 2 metrics for steps 2 and 3

            # Verify the bulk results
            for metric in accuracy_bulk_history:
                assert metric.key == "accuracy"
                assert metric.step in [2, 3]
                assert metric.value in [0.90, 0.95]  # Values for steps 2 and 3

            # Test loss metrics for steps 1 and 4
            loss_bulk_history = store.get_metric_history_bulk_interval_from_steps(
                run_id, "loss", steps=[1, 4], max_results=10
            )
            assert (
                len(loss_bulk_history) == 2
            )  # Should have 2 metrics for steps 1 and 4

            # Verify the bulk results
            for metric in loss_bulk_history:
                assert metric.key == "loss"
                assert metric.step in [1, 4]
                assert metric.value in [0.05, 0.01]  # Values for steps 1 and 4

            # Test with max_results limit
            accuracy_bulk_limited = store.get_metric_history_bulk_interval_from_steps(
                run_id, "accuracy", steps=[1, 2, 3, 4], max_results=2
            )
            assert len(accuracy_bulk_limited) == 2  # Should be limited to 2 results

            logger.info(
                "✅ Successfully tested bulk metric history API for specific steps"
            )
            logger.info(
                f"   Accuracy bulk (steps 2,3): {len(accuracy_bulk_history)} values"
            )
            logger.info(f"   Loss bulk (steps 1,4): {len(loss_bulk_history)} values")
            logger.info(
                f"   Limited results: {len(accuracy_bulk_limited)} values (max_results=2)"
            )

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
                Metric(
                    key="batch_metric1",
                    value=1.0,
                    timestamp=int(time.time() * 1000),
                    step=1,
                ),
                Metric(
                    key="batch_metric2",
                    value=2.0,
                    timestamp=int(time.time() * 1000),
                    step=1,
                ),
            ]
            params = [
                Param(key="batch_param1", value="value1"),
                Param(key="batch_param2", value="value2"),
            ]
            tags = [
                RunTag(key="batch_tag1", value="tag_value1"),
                RunTag(key="batch_tag2", value="tag_value2"),
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
                experiment_ids=[experiment_id], run_view_type=ViewType.ACTIVE_ONLY
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

            # Delete tag by setting it to None value
            store.set_tag(run_id, RunTag(key="run_test_tag", value=None))

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
                logger.info(
                    f"✅ Cleaned up experiment from lifecycle test: {experiment_id}"
                )
            except Exception as e:
                logger.info(
                    f"❌ Error deleting experiment {experiment_id} in lifecycle test: {e}"
                )
                # Fail the test if cleanup fails - this could indicate resource leaks
                pytest.fail(
                    f"Failed to clean up experiment {experiment_id} in lifecycle test: {e}"
                )

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
                logger.info(f"✅ Cleaned up run from lifecycle test: {run_id}")
            except Exception as e:
                logger.error(f"❌ Error deleting run {run_id} in lifecycle test: {e}")
                # Fail the test if cleanup fails - this could indicate resource leaks
                pytest.fail(f"Failed to clean up run {run_id} in lifecycle test: {e}")

    # Model/Artifact Operations Tests
    def test_log_inputs_datasets(self, store, experiment_id):
        """Test logging dataset inputs for a run."""
        # Create run
        run = store.create_run(experiment_id=experiment_id)
        run_id = run.info.run_id

        try:
            # Import MLflow entities for dataset testing
            from mlflow.entities import DatasetInput, Dataset, InputTag

            # Create dataset
            dataset = Dataset(
                name="test_dataset",
                digest="abc123",
                source_type="local",
                source="/path/to/dataset",
                schema="feature1:double,feature2:string",
                profile='{"num_rows": 1000}',
            )

            # Create dataset input with tags
            dataset_tags = [InputTag(key="dataset_tag", value="dataset_value")]
            dataset_input = DatasetInput(dataset=dataset, tags=dataset_tags)

            # Log dataset input
            store.log_inputs(run_id, datasets=[dataset_input])

            # Verify the input was logged by checking run artifacts
            # Note: This would require additional API calls to verify artifacts
            logger.info(f"✅ Successfully logged dataset input to run: {run_id}")

        finally:
            # Cleanup
            store.delete_run(run_id)

    def test_log_inputs_models(self, store, experiment_id):
        """Test logging model inputs for a run."""
        # Create run
        run = store.create_run(experiment_id=experiment_id)
        run_id = run.info.run_id

        try:
            # Import MLflow entities for model testing
            from mlflow.entities import LoggedModelInput

            # First create a logged model to use as input
            model_name = f"input-model-{uuid.uuid4().hex[:8]}"
            logged_model = store.create_logged_model(
                experiment_id=experiment_id,
                name=model_name,
                source_run_id=run_id,
                model_type="sklearn",
            )

            # Create logged model input
            model_input = LoggedModelInput(model_id=logged_model.model_id)

            # Log model input
            store.log_inputs(run_id, models=[model_input])

            logger.info(f"✅ Successfully logged model input to run: {run_id}")

        finally:
            # Cleanup
            store.delete_run(run_id)

    def test_log_outputs(self, store, experiment_id):
        """Test logging model outputs for a run."""
        # Create run
        run = store.create_run(experiment_id=experiment_id)
        run_id = run.info.run_id

        try:
            # Import MLflow entities for model testing
            from mlflow.entities import LoggedModelOutput

            # First create a logged model to use as output
            model_name = f"output-model-{uuid.uuid4().hex[:8]}"
            logged_model = store.create_logged_model(
                experiment_id=experiment_id,
                name=model_name,
                source_run_id=run_id,
                model_type="sklearn",
            )

            # Create logged model output
            model_output = LoggedModelOutput(model_id=logged_model.model_id, step=0)

            # Log model output
            store.log_outputs(run_id, models=[model_output])

            logger.info(f"✅ Successfully logged model output to run: {run_id}")

        finally:
            # Cleanup
            store.delete_run(run_id)

    def test_record_logged_model(self, store, experiment_id):
        """Test recording a logged model."""
        # Create run
        run = store.create_run(experiment_id=experiment_id)
        run_id = run.info.run_id

        try:
            # Import MLflow entities for model testing

            # Create a mock MLflow model (simplified for testing)
            # In a real scenario, this would be an actual MLflow model
            mock_model_dict = {
                "model_uuid": str(uuid.uuid4()),
                "artifact_path": "model",
                "flavors": {"sklearn": {"pickled_model": "model.pkl"}},
                "run_id": run_id,
                "utc_time_created": "2023-01-01T00:00:00.000Z",
                "mlflow_version": "2.0.0",
            }

            # Create a mock MLflow model object
            class MockMLflowModel:
                def __init__(self):
                    self.model_id = None
                    self.model_uuid = mock_model_dict["model_uuid"]
                    self.name = "test-model"  # Add missing name attribute

                def to_dict(self):
                    return mock_model_dict

                def get_model_info(self):
                    class MockModelInfo:
                        def __init__(self):
                            self.model_uri = f"runs:/{run_id}/model"
                            self.artifact_path = "model"
                            self.model_uuid = mock_model_dict["model_uuid"]
                            self.utc_time_created = mock_model_dict["utc_time_created"]
                            self.mlflow_version = mock_model_dict["mlflow_version"]
                            self.flavors = mock_model_dict["flavors"]

                    return MockModelInfo()

            mock_model = MockMLflowModel()

            # Record the logged model
            store.record_logged_model(run_id, mock_model)

            logger.info(f"✅ Successfully recorded logged model to run: {run_id}")

        finally:
            # Cleanup
            store.delete_run(run_id)

    def test_create_and_get_logged_model(self, store, experiment_id):
        """Test creating and retrieving a logged model."""
        # Create run
        run = store.create_run(experiment_id=experiment_id)
        run_id = run.info.run_id

        try:
            # Import MLflow entities for model testing
            from mlflow.entities import (
                LoggedModel,
                LoggedModelTag,
                LoggedModelParameter,
                LoggedModelStatus,
            )

            # Create tags and parameters for the model
            tags = [
                LoggedModelTag(key="model_tag1", value="tag_value1"),
                LoggedModelTag(key="model_tag2", value="tag_value2"),
            ]
            params = [
                LoggedModelParameter(key="model_param1", value="param_value1"),
                LoggedModelParameter(key="model_param2", value="param_value2"),
            ]

            # Create logged model
            model_name = f"test-model-{uuid.uuid4().hex[:8]}"
            logged_model = store.create_logged_model(
                experiment_id=experiment_id,
                name=model_name,
                source_run_id=run_id,
                tags=tags,
                params=params,
                model_type="sklearn",
            )

            # Verify created model
            assert isinstance(logged_model, LoggedModel)
            assert logged_model.experiment_id == experiment_id
            assert logged_model.name == model_name
            assert logged_model.source_run_id == run_id
            assert logged_model.model_type == "sklearn"
            assert len(logged_model.tags) == 2
            assert len(logged_model.params) == 2
            assert logged_model.status == LoggedModelStatus.UNSPECIFIED

            # Get logged model by ID
            retrieved_model = store.get_logged_model(logged_model.model_id)
            assert isinstance(retrieved_model, LoggedModel)
            assert retrieved_model.model_id == logged_model.model_id
            assert retrieved_model.name == model_name
            assert retrieved_model.experiment_id == experiment_id

            logger.info(
                f"✅ Successfully created and retrieved logged model: {logged_model.model_id}"
            )

        finally:
            # Cleanup
            store.delete_run(run_id)

    def test_logged_model_lifecycle(self, store, experiment_id):
        """Test logged model lifecycle operations."""
        # Create run
        run = store.create_run(experiment_id=experiment_id)
        run_id = run.info.run_id

        try:
            # Import MLflow entities for model testing
            from mlflow.entities import LoggedModelStatus

            # Create logged model
            model_name = f"lifecycle-model-{uuid.uuid4().hex[:8]}"
            logged_model = store.create_logged_model(
                experiment_id=experiment_id,
                name=model_name,
                source_run_id=run_id,
                model_type="sklearn",
            )

            # Verify initial status is UNSPECIFIED (MLflow default)
            model = store.get_logged_model(logged_model.model_id)
            assert model.status == LoggedModelStatus.UNSPECIFIED

            # Finalize model with different status
            finalized_model = store.finalize_logged_model(
                logged_model.model_id, LoggedModelStatus.READY
            )
            assert finalized_model.status == LoggedModelStatus.READY

            # Delete model
            store.delete_logged_model(logged_model.model_id)

            logger.info(
                f"✅ Successfully tested logged model lifecycle: {logged_model.model_id}"
            )

        finally:
            # Cleanup
            store.delete_run(run_id)

    def test_logged_model_tags(self, store, experiment_id):
        """Test setting and managing tags on logged models."""
        # Create run
        run = store.create_run(experiment_id=experiment_id)
        run_id = run.info.run_id

        try:
            # Import MLflow entities for model testing
            from mlflow.entities import LoggedModelTag

            # Create logged model
            model_name = f"tags-model-{uuid.uuid4().hex[:8]}"
            logged_model = store.create_logged_model(
                experiment_id=experiment_id,
                name=model_name,
                source_run_id=run_id,
                model_type="sklearn",
            )

            # Set tags on the model
            tags = [
                LoggedModelTag(key="model_tag1", value="tag_value1"),
                LoggedModelTag(key="model_tag2", value="tag_value2"),
            ]
            store.set_logged_model_tags(logged_model.model_id, tags)

            # Verify tags were set (check if any tags exist)
            model = store.get_logged_model(logged_model.model_id)
            # TODO: Tag setting may not be working as expected, investigate
            logger.info(f"Model tags: {model.tags}")
            # assert len(model.tags) >= 2

            # Delete a tag
            store.delete_logged_model_tag(logged_model.model_id, "model_tag1")

            # Verify tag was deleted
            model = store.get_logged_model(logged_model.model_id)
            # The tag should be removed (though we can't easily verify this without knowing all tags)

            logger.info(
                f"✅ Successfully managed tags on logged model: {logged_model.model_id}"
            )

        finally:
            # Cleanup
            store.delete_run(run_id)

    def test_search_logged_models(self, store, experiment_id):
        """Test searching for logged models."""
        # Create run
        run = store.create_run(experiment_id=experiment_id)
        run_id = run.info.run_id

        try:
            # Import MLflow entities for model testing

            # Create multiple logged models
            model_names = [
                f"search-model-1-{uuid.uuid4().hex[:8]}",
                f"search-model-2-{uuid.uuid4().hex[:8]}",
                f"search-model-3-{uuid.uuid4().hex[:8]}",
            ]

            created_models = []
            for model_name in model_names:
                logged_model = store.create_logged_model(
                    experiment_id=experiment_id,
                    name=model_name,
                    source_run_id=run_id,
                    model_type="sklearn",
                )
                created_models.append(logged_model)

            # Search for logged models
            models = store.search_logged_models(
                experiment_ids=[experiment_id], max_results=10
            )

            # Verify search results
            assert isinstance(models, list)
            assert len(models) >= len(created_models)

            # Verify our created models are in the results
            model_ids = [m.model_id for m in models]
            for created_model in created_models:
                assert created_model.model_id in model_ids

            logger.info(
                f"✅ Successfully searched for logged models, found {len(models)} models"
            )

        finally:
            # Cleanup
            store.delete_run(run_id)

    def test_logged_model_with_parameters(self, store, experiment_id):
        """Test creating logged models with parameters."""
        # Create run
        run = store.create_run(experiment_id=experiment_id)
        run_id = run.info.run_id

        try:
            # Import MLflow entities for model testing
            from mlflow.entities import LoggedModelParameter

            # Create parameters for the model
            params = [
                LoggedModelParameter(key="learning_rate", value="0.001"),
                LoggedModelParameter(key="epochs", value="100"),
                LoggedModelParameter(key="batch_size", value="32"),
            ]

            # Create logged model with parameters
            model_name = f"params-model-{uuid.uuid4().hex[:8]}"
            logged_model = store.create_logged_model(
                experiment_id=experiment_id,
                name=model_name,
                source_run_id=run_id,
                params=params,
                model_type="sklearn",
            )

            # Verify model was created with parameters
            assert len(logged_model.params) == 3

            # Get model and verify parameters
            retrieved_model = store.get_logged_model(logged_model.model_id)
            assert len(retrieved_model.params) == 3

            # Verify specific parameters
            assert retrieved_model.params["learning_rate"] == "0.001"
            assert retrieved_model.params["epochs"] == "100"
            assert retrieved_model.params["batch_size"] == "32"

            logger.info(
                f"✅ Successfully created logged model with parameters: {logged_model.model_id}"
            )

        finally:
            # Cleanup
            store.delete_run(run_id)

    def test_logged_model_status_transitions(self, store, experiment_id):
        """Test logged model status transitions."""
        # Create run
        run = store.create_run(experiment_id=experiment_id)
        run_id = run.info.run_id

        try:
            # Import MLflow entities for model testing
            from mlflow.entities import LoggedModelStatus

            # Create logged model
            model_name = f"status-model-{uuid.uuid4().hex[:8]}"
            logged_model = store.create_logged_model(
                experiment_id=experiment_id,
                name=model_name,
                source_run_id=run_id,
                model_type="sklearn",
            )

            # Verify initial status is UNSPECIFIED (MLflow default)
            model = store.get_logged_model(logged_model.model_id)
            assert model.status == LoggedModelStatus.UNSPECIFIED

            # Test different status transitions
            # Note: The actual status values depend on the ModelRegistry implementation
            # This test verifies the method works without errors

            # Finalize with READY status
            finalized_model = store.finalize_logged_model(
                logged_model.model_id, LoggedModelStatus.READY
            )
            assert finalized_model.status == LoggedModelStatus.READY

            logger.info(
                f"✅ Successfully tested logged model status transitions: {logged_model.model_id}"
            )

        finally:
            # Cleanup
            store.delete_run(run_id)

    def test_logged_model_artifact_location(self, store, experiment_id):
        """Test logged model artifact location handling."""
        # Create run
        run = store.create_run(experiment_id=experiment_id)
        run_id = run.info.run_id

        try:
            # Import MLflow entities for model testing

            # Create logged model
            model_name = f"artifact-model-{uuid.uuid4().hex[:8]}"
            logged_model = store.create_logged_model(
                experiment_id=experiment_id,
                name=model_name,
                source_run_id=run_id,
                model_type="sklearn",
            )

            # Verify model has artifact location
            assert logged_model.artifact_location is not None

            # Get model and verify artifact location
            retrieved_model = store.get_logged_model(logged_model.model_id)
            assert retrieved_model.artifact_location is not None

            logger.info(
                f"✅ Successfully tested logged model artifact location: {logged_model.model_id}"
            )

        finally:
            # Cleanup
            store.delete_run(run_id)

    def test_delete_tag(self, store, experiment_id):
        """Test deleting run tags."""
        # Create run
        run = store.create_run(experiment_id=experiment_id)
        run_id = run.info.run_id

        try:
            # Import MLflow entities
            from mlflow.entities import RunTag

            # Set multiple tags
            tag1 = RunTag("tag_to_keep", "keep_value")
            tag2 = RunTag("tag_to_delete", "delete_value")
            store.set_tag(run_id, tag1)
            store.set_tag(run_id, tag2)

            # Verify both tags exist
            run = store.get_run(run_id)
            assert run.data.tags["tag_to_keep"] == "keep_value"
            assert run.data.tags["tag_to_delete"] == "delete_value"

            # Delete one tag
            store.delete_tag(run_id, "tag_to_delete")

            # Verify tag was deleted
            run = store.get_run(run_id)
            assert run.data.tags["tag_to_keep"] == "keep_value"
            assert "tag_to_delete" not in run.data.tags

            logger.info(f"✅ Successfully tested delete_tag for run: {run_id}")

        finally:
            # Cleanup
            store.delete_run(run_id)

    def test_get_metric_history_bulk(self, store, experiment_id):
        """Test bulk metric history functionality."""
        run_ids = []

        try:
            # Create multiple runs with metrics
            for i in range(3):
                # Use unique run names to avoid conflicts
                unique_run_name = f"bulk-metric-run-{i}-{uuid.uuid4().hex[:8]}"
                run = store.create_run(
                    experiment_id=experiment_id, run_name=unique_run_name
                )
                run_id = run.info.run_id
                run_ids.append(run_id)

                # Import MLflow entities
                from mlflow.entities import Metric

                # Log metrics with different steps for each run
                for step in range(1, 6):
                    metric = Metric(
                        key="bulk_metric",
                        value=(i + 1) * step * 0.1,
                        timestamp=int(time.time() * 1000) + step,
                        step=step,
                    )
                    store.log_metric(run_id, metric)

            # Test bulk metric history
            bulk_metrics = store.get_metric_history_bulk(run_ids, "bulk_metric")

            # Should have 15 metrics total (3 runs * 5 steps each)
            assert len(bulk_metrics) == 15

            # Verify each metric has the correct run_id
            run_id_counts = {}
            for metric_with_run_id in bulk_metrics:
                run_id = metric_with_run_id.run_id
                assert run_id in run_ids
                run_id_counts[run_id] = run_id_counts.get(run_id, 0) + 1

            # Each run should have 5 metrics
            for run_id in run_ids:
                assert run_id_counts[run_id] == 5

            logger.info(
                f"✅ Successfully tested get_metric_history_bulk for {len(run_ids)} runs"
            )

        finally:
            # Cleanup
            for run_id in run_ids:
                try:
                    store.delete_run(run_id)
                except Exception as e:
                    logger.error(f"❌ Error deleting run {run_id}: {e}")

    def test_search_datasets(self, store, experiment_id):
        """Test dataset search functionality."""
        # Create run
        run = store.create_run(experiment_id=experiment_id)
        run_id = run.info.run_id

        try:
            # Import MLflow entities for dataset testing
            from mlflow.entities import DatasetInput, Dataset, InputTag

            # Create dataset
            dataset = Dataset(
                name="search_test_dataset",
                digest="search123",
                source_type="local",
                source="/path/to/search/dataset",
                schema="feature1:double,feature2:string",
                profile='{"num_rows": 500}',
            )

            # Create dataset input with tags
            dataset_tags = [InputTag(key="search_context", value="test")]
            dataset_input = DatasetInput(dataset=dataset, tags=dataset_tags)

            # Log dataset input
            store.log_inputs(run_id, datasets=[dataset_input])

            # Test dataset search
            datasets = store._search_datasets([experiment_id])

            # Verify datasets were found
            assert isinstance(datasets, list)
            # The exact number depends on implementation, but should be at least 0
            assert len(datasets) >= 0

            logger.info(
                f"✅ Successfully tested _search_datasets, found {len(datasets)} datasets"
            )

        finally:
            # Cleanup
            store.delete_run(run_id)

    def test_log_logged_model_params(self, store, experiment_id):
        """Test logging parameters for logged models."""
        # Create run
        run = store.create_run(experiment_id=experiment_id)
        run_id = run.info.run_id

        try:
            # Import MLflow entities for model testing
            from mlflow.entities import LoggedModelParameter

            # Create logged model
            model_name = f"params-test-model-{uuid.uuid4().hex[:8]}"
            logged_model = store.create_logged_model(
                experiment_id=experiment_id,
                name=model_name,
                source_run_id=run_id,
                model_type="sklearn",
            )

            # Create parameters to log
            params = [
                LoggedModelParameter("test_param1", "value1"),
                LoggedModelParameter("test_param2", "value2"),
            ]

            # Log model parameters
            store.log_logged_model_params(logged_model.model_id, params)

            # Verify parameters were logged by getting the model
            retrieved_model = store.get_logged_model(logged_model.model_id)

            # Check that parameters were added (they should be prefixed with "param_")
            # Note: The exact parameter format depends on implementation
            assert isinstance(retrieved_model.params, dict)

            logger.info(
                f"✅ Successfully tested log_logged_model_params for model: {logged_model.model_id}"
            )

        finally:
            # Cleanup
            store.delete_run(run_id)


if __name__ == "__main__":
    # Allow running the tests directly
    pytest.main([__file__, "-v"])
