"""
End-to-end tests for ModelRegistryStore with real Model Registry server.

These tests use MLflow's standard APIs and rely on the tracking URI configuration
to connect to the Model Registry server. The tracking URI should be configured
via environment variables or the Makefile.

Required environment variables:
- MLFLOW_TRACKING_URI: Set to "modelregistry://host:port" format
- MODEL_REGISTRY_TOKEN: Authentication token for the Model Registry server
"""

import os
import pytest
import uuid
import tempfile
import pandas as pd

import mlflow
from mlflow.entities import (
    Experiment,
    RunStatus,
    ViewType,
    LifecycleStage,
)
from mlflow.exceptions import MlflowException
from sklearn.ensemble import RandomForestClassifier
from sklearn.datasets import make_classification
from sklearn.model_selection import train_test_split


@pytest.mark.skipif(
    not os.getenv("MLFLOW_TRACKING_URI"),
    reason="MLFLOW_TRACKING_URI environment variable not set",
)
class TestModelRegistryStoreE2E:
    """End-to-end tests for ModelRegistryStore using MLflow APIs."""

    @pytest.fixture(scope="function")
    def experiment_name(self):
        """Generate a unique experiment name for testing."""
        return f"e2e-test-{uuid.uuid4().hex[:8]}"

    @pytest.fixture(scope="function")
    def experiment_id(self, experiment_name):
        """Create a test experiment and return its ID."""
        try:
            experiment_id = mlflow.create_experiment(experiment_name)
            yield experiment_id
        finally:
            # Cleanup: delete the experiment
            try:
                mlflow.delete_experiment(experiment_id)
            except MlflowException:
                pass  # Ignore cleanup errors

    @pytest.fixture
    def sample_dataset(self):
        """Create a lightweight sample dataset."""
        # Generate a small synthetic dataset
        X, y = make_classification(
            n_samples=100,
            n_features=5,
            n_informative=3,
            n_redundant=1,
            n_clusters_per_class=1,
            random_state=42,
        )

        # Create DataFrame
        feature_names = [f"feature_{i}" for i in range(X.shape[1])]
        df = pd.DataFrame(X, columns=feature_names)
        df["target"] = y

        return df

    @pytest.fixture
    def sample_model(self, sample_dataset):
        """Create a lightweight sample model."""
        # Prepare data
        X = sample_dataset.drop("target", axis=1)
        y = sample_dataset["target"]

        # Split data
        X_train, X_test, y_train, y_test = train_test_split(
            X, y, test_size=0.2, random_state=42
        )

        # Train a small model
        model = RandomForestClassifier(n_estimators=10, random_state=42)
        model.fit(X_train, y_train)

        return model, X_test, y_test

    def test_mlflow_connection(self):
        """Test that we can connect to the tracking server via MLflow."""
        # Try to list experiments to verify connection
        experiments = mlflow.search_experiments()
        assert isinstance(experiments, list)
        print(
            f"Successfully connected to tracking server via MLflow. Found {len(experiments)} experiments."
        )

    def test_create_and_get_experiment(self, experiment_name):
        """Test creating and retrieving an experiment using MLflow."""
        # Create experiment
        experiment_id = mlflow.create_experiment(experiment_name)
        assert experiment_id is not None
        assert isinstance(experiment_id, str)

        # Get experiment by ID
        experiment = mlflow.get_experiment(experiment_id)
        assert isinstance(experiment, Experiment)
        assert experiment.experiment_id == experiment_id
        assert experiment.name == experiment_name
        assert experiment.lifecycle_stage == LifecycleStage.ACTIVE

        # Get experiment by name
        experiment_by_name = mlflow.get_experiment_by_name(experiment_name)
        assert experiment_by_name is not None
        assert experiment_by_name.experiment_id == experiment_id

        # Get experiment by name that doesn't exist
        experiment_by_name = mlflow.get_experiment_by_name("nonexistent")
        assert experiment_by_name is None

        # Cleanup
        mlflow.delete_experiment(experiment_id)

    def test_create_and_log_run(self, experiment_id, sample_dataset, sample_model):
        """Test creating and logging to a run using MLflow."""
        model, X_test, y_test = sample_model

        with mlflow.start_run(experiment_id=experiment_id) as run:
            # Log parameters
            mlflow.log_param("learning_rate", "0.01")
            mlflow.log_param("epochs", "100")
            mlflow.log_param("model_type", "random_forest")

            # Log metrics
            mlflow.log_metric("accuracy", 0.95, step=1)
            mlflow.log_metric("loss", 0.05, step=1)

            # Log multiple metric values
            for i in range(2, 6):
                mlflow.log_metric("accuracy", 0.95 + (i * 0.01), step=i)
                mlflow.log_metric("loss", 0.05 - (i * 0.01), step=i)

            # Log tags
            mlflow.set_tag("test_tag", "test_value")
            mlflow.set_tag("model_version", "1.0")

            # Log the model
            mlflow.sklearn.log_model(model, "model")

            # Log dataset as artifact
            dataset_path = tempfile.mktemp(suffix=".csv")
            sample_dataset.to_csv(dataset_path, index=False)
            mlflow.log_artifact(dataset_path, "dataset")

            run_id = run.info.run_id

        # Verify the run was created and logged
        retrieved_run = mlflow.get_run(run_id)
        assert retrieved_run.info.run_id == run_id
        assert retrieved_run.info.experiment_id == experiment_id
        assert retrieved_run.info.status == RunStatus.FINISHED
        assert retrieved_run.info.lifecycle_stage == LifecycleStage.ACTIVE

        # Verify parameters
        assert retrieved_run.data.params["learning_rate"] == "0.01"
        assert retrieved_run.data.params["epochs"] == "100"
        assert retrieved_run.data.params["model_type"] == "random_forest"

        # Verify metrics (should show the latest value for each metric key)
        assert retrieved_run.data.metrics["accuracy"] == 1.0  # Latest accuracy value
        assert retrieved_run.data.metrics["loss"] == 0.01  # Latest loss value

        # Verify tags
        assert retrieved_run.data.tags["test_tag"] == "test_value"
        assert retrieved_run.data.tags["model_version"] == "1.0"

        # Test metric history
        accuracy_history = mlflow.get_metric_history(run_id, "accuracy")
        assert len(accuracy_history) == 5  # Should have 5 accuracy values

        # Cleanup
        os.unlink(dataset_path)

    def test_log_dataset_inputs(self, experiment_id, sample_dataset):
        """Test logging dataset inputs using MLflow."""
        with mlflow.start_run(experiment_id=experiment_id) as run:
            # Create a temporary dataset file
            dataset_path = tempfile.mktemp(suffix=".csv")
            sample_dataset.to_csv(dataset_path, index=False)

            # Log dataset as artifact
            mlflow.log_artifact(dataset_path, "input_dataset")

            # Log dataset metadata as parameters
            mlflow.log_param("dataset_rows", str(len(sample_dataset)))
            mlflow.log_param("dataset_columns", str(len(sample_dataset.columns)))
            mlflow.log_param(
                "dataset_features", str(len(sample_dataset.columns) - 1)
            )  # excluding target

            # Log dataset statistics as metrics
            mlflow.log_metric(
                "dataset_size_mb", len(sample_dataset.to_csv()) / (1024 * 1024)
            )

            run_id = run.info.run_id

        # Verify dataset was logged
        retrieved_run = mlflow.get_run(run_id)
        assert retrieved_run.data.params["dataset_rows"] == str(len(sample_dataset))
        assert retrieved_run.data.params["dataset_columns"] == str(
            len(sample_dataset.columns)
        )

        # Cleanup
        os.unlink(dataset_path)

    def test_log_multiple_datasets(self, experiment_id):
        """Test logging multiple datasets using MLflow."""
        with mlflow.start_run(experiment_id=experiment_id) as run:
            # Create multiple small datasets
            datasets = []
            for i in range(3):
                # Create small synthetic dataset
                X, y = make_classification(
                    n_samples=50, n_features=3, random_state=42 + i
                )
                df = pd.DataFrame(X, columns=[f"feature_{j}" for j in range(3)])
                df["target"] = y
                datasets.append(df)

                # Save to temporary file
                dataset_path = tempfile.mktemp(suffix=f"_dataset_{i}.csv")
                df.to_csv(dataset_path, index=False)

                # Log dataset
                mlflow.log_artifact(dataset_path, f"dataset_{i}")

                # Log dataset metadata
                mlflow.log_param(f"dataset_{i}_rows", str(len(df)))
                mlflow.log_param(f"dataset_{i}_columns", str(len(df.columns)))

                # Cleanup temp file
                os.unlink(dataset_path)

            run_id = run.info.run_id

        # Verify datasets were logged
        retrieved_run = mlflow.get_run(run_id)
        assert retrieved_run.data.params["dataset_0_rows"] == "50"
        assert retrieved_run.data.params["dataset_1_rows"] == "50"
        assert retrieved_run.data.params["dataset_2_rows"] == "50"

    def test_log_model_inputs_and_outputs(self, experiment_id, sample_model):
        """Test logging model inputs and outputs using MLflow."""
        model, X_test, y_test = sample_model

        # First run: create and log a model
        with mlflow.start_run(experiment_id=experiment_id) as input_run:
            # Log the model
            mlflow.sklearn.log_model(model, "input_model")

            # Log model metadata
            mlflow.log_param("model_type", "random_forest")
            mlflow.log_param("n_estimators", "10")
            mlflow.log_metric("input_model_score", model.score(X_test, y_test))

            input_run_id = input_run.info.run_id

        # Second run: use the model as input and create a new model as output
        with mlflow.start_run(experiment_id=experiment_id) as output_run:
            # Load the input model
            loaded_model = mlflow.sklearn.load_model(
                f"runs:/{input_run_id}/input_model"
            )
            assert loaded_model is not None

            # Create a new model (fine-tuned version)
            new_model = RandomForestClassifier(n_estimators=15, random_state=42)
            new_model.fit(X_test, y_test)  # Using test data for simplicity

            # Log the new model
            mlflow.sklearn.log_model(new_model, "output_model")

            # Log model metadata
            mlflow.log_param("model_type", "random_forest_finetuned")
            mlflow.log_param("n_estimators", "15")
            mlflow.log_metric("output_model_score", new_model.score(X_test, y_test))

            # Log reference to input model
            mlflow.log_param("input_model_run_id", input_run_id)

            output_run_id = output_run.info.run_id

        # Verify both runs
        input_retrieved = mlflow.get_run(input_run_id)
        output_retrieved = mlflow.get_run(output_run_id)

        assert input_retrieved.data.params["model_type"] == "random_forest"
        assert output_retrieved.data.params["model_type"] == "random_forest_finetuned"
        assert output_retrieved.data.params["input_model_run_id"] == input_run_id

    def test_log_model_with_steps(self, experiment_id, sample_dataset):
        """Test logging model outputs with different steps using MLflow."""
        X = sample_dataset.drop("target", axis=1)
        y = sample_dataset["target"]

        with mlflow.start_run(experiment_id=experiment_id) as run:
            # Simulate training with multiple steps
            for step in [0, 10, 20, 50, 100]:
                # Create a model for this step (simulating training progress)
                model = RandomForestClassifier(
                    n_estimators=5 + step // 20, random_state=42
                )
                model.fit(X, y)

                # Log model for this step
                mlflow.sklearn.log_model(model, f"model_step_{step}")

                # Log metrics for this step
                score = model.score(X, y)
                mlflow.log_metric("accuracy", score, step=step)
                mlflow.log_metric("model_complexity", model.n_estimators, step=step)

                # Log parameters for this step
                mlflow.log_param(f"n_estimators_step_{step}", str(model.n_estimators))

            run_id = run.info.run_id

        # Verify step-wise logging
        retrieved_run = mlflow.get_run(run_id)
        assert retrieved_run.data.metrics["accuracy"] == 1.0  # Latest accuracy
        assert retrieved_run.data.metrics["model_complexity"] == 10  # Latest complexity

        # Test metric history
        accuracy_history = mlflow.get_metric_history(run_id, "accuracy")
        assert len(accuracy_history) == 5  # Should have 5 accuracy values

    def test_log_inputs_outputs_with_metrics_params(
        self, experiment_id, sample_dataset, sample_model
    ):
        """Test logging inputs/outputs along with metrics and parameters using MLflow."""
        model, X_test, y_test = sample_model

        with mlflow.start_run(experiment_id=experiment_id) as run:
            # Log parameters
            mlflow.log_param("model_type", "random_forest")
            mlflow.log_param("n_estimators", "10")
            mlflow.log_param("random_state", "42")

            # Log metrics
            mlflow.log_metric("training_accuracy", 0.95, step=1)
            mlflow.log_metric("validation_accuracy", 0.92, step=1)

            # Log dataset as input
            dataset_path = tempfile.mktemp(suffix=".csv")
            sample_dataset.to_csv(dataset_path, index=False)
            mlflow.log_artifact(dataset_path, "input_dataset")

            # Log model as output
            mlflow.sklearn.log_model(model, "output_model")

            # Log additional metrics
            mlflow.log_metric("test_accuracy", model.score(X_test, y_test))
            mlflow.log_metric("model_size_mb", 0.1)  # Approximate size

            # Log tags
            mlflow.set_tag("pipeline_stage", "training")
            mlflow.set_tag("data_version", "1.0")

            run_id = run.info.run_id

        # Verify run data
        retrieved_run = mlflow.get_run(run_id)
        assert len(retrieved_run.data.params) >= 3
        assert len(retrieved_run.data.metrics) >= 4
        assert retrieved_run.data.params["model_type"] == "random_forest"
        assert retrieved_run.data.metrics["training_accuracy"] == 0.95
        assert retrieved_run.data.tags["pipeline_stage"] == "training"

        # Cleanup
        os.unlink(dataset_path)

    def test_experiment_lifecycle(self, experiment_name):
        """Test experiment lifecycle operations using MLflow."""
        # Create experiment
        experiment_id = mlflow.create_experiment(experiment_name)

        try:
            # Verify experiment is active
            experiment = mlflow.get_experiment(experiment_id)
            assert experiment.lifecycle_stage == LifecycleStage.ACTIVE

            # Delete experiment
            mlflow.delete_experiment(experiment_id)

            # Verify experiment is deleted
            experiment = mlflow.get_experiment(experiment_id)
            assert experiment.lifecycle_stage == LifecycleStage.DELETED

            # Restore experiment using the tracking store
            from mlflow.tracking import _get_store

            store = _get_store()
            store.restore_experiment(experiment_id)

            # Verify experiment is active again
            experiment = mlflow.get_experiment(experiment_id)
            assert experiment.lifecycle_stage == LifecycleStage.ACTIVE

        finally:
            # Cleanup
            try:
                mlflow.delete_experiment(experiment_id)
            except MlflowException:
                pass

    def test_run_lifecycle(self, experiment_id):
        """Test run lifecycle operations using MLflow."""
        with mlflow.start_run(experiment_id=experiment_id) as run:
            mlflow.log_param("test_param", "test_value")
            run_id = run.info.run_id

        # Verify run is active
        retrieved_run = mlflow.get_run(run_id)
        assert retrieved_run.info.lifecycle_stage == LifecycleStage.ACTIVE

        # Delete run
        mlflow.delete_run(run_id)

        # Verify run is deleted
        retrieved_run = mlflow.get_run(run_id)
        assert retrieved_run.info.lifecycle_stage == LifecycleStage.DELETED

        # Restore run using the tracking store
        from mlflow.tracking import _get_store

        store = _get_store()
        store.restore_run(run_id)

        # Verify run is active again
        retrieved_run = mlflow.get_run(run_id)
        assert retrieved_run.info.lifecycle_stage == LifecycleStage.ACTIVE

    def test_search_experiments(self, experiment_id):
        """Test searching experiments using MLflow."""
        # Search all experiments
        experiments = mlflow.search_experiments(view_type=ViewType.ALL)
        assert isinstance(experiments, list)
        assert len(experiments) > 0

        # Search active experiments only
        active_experiments = mlflow.search_experiments(view_type=ViewType.ACTIVE_ONLY)
        assert isinstance(active_experiments, list)

        # Verify our test experiment is in the results
        experiment_ids = [exp.experiment_id for exp in active_experiments]
        assert experiment_id in experiment_ids

    def test_search_runs(self, experiment_id):
        """Test searching runs using MLflow."""
        # Create a test run
        with mlflow.start_run(experiment_id=experiment_id) as run:
            mlflow.log_param("search_test", "value")
            run_id = run.info.run_id

        try:
            # Search runs in the experiment
            runs = mlflow.search_runs(
                experiment_ids=[experiment_id], run_view_type=ViewType.ACTIVE_ONLY
            )
            assert isinstance(runs, pd.DataFrame)
            assert len(runs) > 0

            # Verify our test run is in the results
            run_ids = runs["run_id"].tolist()
            assert run_id in run_ids

        finally:
            # Cleanup
            mlflow.delete_run(run_id)

    def test_experiment_tags(self, experiment_id):
        """Test setting and managing experiment tags using MLflow."""
        # set active experiment to avoid errors with mlflow.tracking.default_experiment
        mlflow.set_experiment(experiment_id=experiment_id)

        # Set experiment tag
        mlflow.set_experiment_tag("test_tag", "test_value")

        # Get experiment and verify tag
        experiment = mlflow.get_experiment(experiment_id)
        assert experiment.tags["test_tag"] == "test_value"

    def test_run_tags(self, experiment_id):
        """Test setting and managing run tags using MLflow."""
        with mlflow.start_run(experiment_id=experiment_id) as run:
            # Set run tag
            mlflow.set_tag("run_test_tag", "run_test_value")

            # Set another tag
            mlflow.set_tag("run_test_tag2", "run_test_value2")

            run_id = run.info.run_id

        try:
            # Get run and verify tags
            retrieved_run = mlflow.get_run(run_id)
            assert retrieved_run.data.tags["run_test_tag"] == "run_test_value"
            assert retrieved_run.data.tags["run_test_tag2"] == "run_test_value2"

        finally:
            # Cleanup
            mlflow.delete_run(run_id)

    def test_batch_logging(self, experiment_id):
        """Test batch logging of metrics, parameters, and tags using MLflow."""
        with mlflow.start_run(experiment_id=experiment_id) as run:
            # Log multiple parameters
            params = {
                "batch_param1": "value1",
                "batch_param2": "value2",
                "batch_param3": "value3",
            }
            for key, value in params.items():
                mlflow.log_param(key, value)

            # Log multiple metrics
            metrics = {"batch_metric1": 1.0, "batch_metric2": 2.0, "batch_metric3": 3.0}
            for key, value in metrics.items():
                mlflow.log_metric(key, value, step=1)

            # Log multiple tags
            tags = {
                "batch_tag1": "tag_value1",
                "batch_tag2": "tag_value2",
                "batch_tag3": "tag_value3",
            }
            for key, value in tags.items():
                mlflow.set_tag(key, value)

            run_id = run.info.run_id

        # Verify batch data
        retrieved_run = mlflow.get_run(run_id)
        assert len(retrieved_run.data.params) >= 3
        assert len(retrieved_run.data.metrics) >= 3
        assert len(retrieved_run.data.tags) >= 3

        # Verify specific values
        assert retrieved_run.data.params["batch_param1"] == "value1"
        assert retrieved_run.data.metrics["batch_metric1"] == 1.0
        assert retrieved_run.data.tags["batch_tag1"] == "tag_value1"

    def test_metric_history(self, experiment_id):
        """Test metric history functionality using MLflow."""
        with mlflow.start_run(experiment_id=experiment_id) as run:
            # Log metrics with different steps
            for step in range(1, 11):
                mlflow.log_metric("test_metric", step * 0.1, step=step)

            run_id = run.info.run_id

        # Test metric history
        metric_history = mlflow.get_metric_history(run_id, "test_metric")
        assert len(metric_history) == 10

        # Verify values are in order
        for i, metric in enumerate(metric_history):
            assert metric.value == (i + 1) * 0.1
            assert metric.step == i + 1

    def test_artifact_logging(self, experiment_id):
        """Test artifact logging using MLflow."""
        with mlflow.start_run(experiment_id=experiment_id) as run:
            # Create a simple text file
            text_content = (
                "This is a test artifact file.\nIt contains multiple lines.\n"
            )
            text_path = tempfile.mktemp(suffix=".txt")
            with open(text_path, "w") as f:
                f.write(text_content)

            # Log the text file
            mlflow.log_artifact(text_path, "text_artifacts")

            # Create a JSON file
            json_content = '{"key": "value", "number": 42, "list": [1, 2, 3]}'
            json_path = tempfile.mktemp(suffix=".json")
            with open(json_path, "w") as f:
                f.write(json_content)

            # Log the JSON file
            mlflow.log_artifact(json_path, "json_artifacts")

            run_id = run.info.run_id

        # Cleanup temp files
        os.unlink(text_path)
        os.unlink(json_path)

        # Verify artifacts were logged (we can't easily verify artifact content via MLflow API)
        # but we can verify the run was created successfully
        retrieved_run = mlflow.get_run(run_id)
        assert retrieved_run.info.run_id == run_id


def test_mlflow_integration():
    """Test that the ModelRegistryStore works with MLflow's tracking API."""
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
        except MlflowException:
            pass


if __name__ == "__main__":
    # Allow running the tests directly
    pytest.main([__file__, "-v"])
