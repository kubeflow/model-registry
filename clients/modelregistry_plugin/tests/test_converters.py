"""Tests for MLflowEntityConverter."""

import json
from unittest.mock import Mock

import pytest
from mlflow.entities import (
    Dataset,
    DatasetInput,
    Experiment,
    ExperimentTag,
    LoggedModel,
    LoggedModelParameter,
    LoggedModelTag,
    Metric,
    Param,
    Run,
    RunData,
    RunInfo,
    RunInputs,
    RunOutputs,
    RunStatus,
    RunTag,
)
from mlflow.exceptions import MlflowException
from mlflow.models import Model

from modelregistry_plugin.converters import MLflowEntityConverter
from modelregistry_plugin.utils import ModelIOType


class TestMLflowEntityConverter:
    def test_to_mlflow_metric(self):
        """Test converting Model Registry metric to MLflow Metric."""
        metric_data = {
            "name": "accuracy",
            "value": 0.95,
            "timestamp": 1234567890,
            "step": 1,
        }

        metric = MLflowEntityConverter.to_mlflow_metric(metric_data)

        assert isinstance(metric, Metric)
        assert metric.key == "accuracy"
        assert metric.value == 0.95
        assert metric.timestamp == 1234567890
        assert metric.step == 1

    def test_to_mlflow_metric_with_create_time_fallback(self):
        """Test converting metric with createTimeSinceEpoch fallback."""
        metric_data = {
            "name": "accuracy",
            "value": 0.95,
            "createTimeSinceEpoch": 1234567890,
            "step": 1,
        }

        metric = MLflowEntityConverter.to_mlflow_metric(metric_data)

        assert metric.timestamp == 1234567890
        assert metric.step == 1

    def test_to_mlflow_metric_default_step(self):
        """Test converting metric with default step value."""
        metric_data = {
            "name": "accuracy",
            "value": 0.95,
            "timestamp": 1234567890,
        }

        metric = MLflowEntityConverter.to_mlflow_metric(metric_data)

        assert metric.step == 0

    def test_to_mlflow_param(self):
        """Test converting Model Registry parameter to MLflow Param."""
        param_data = {
            "name": "learning_rate",
            "value": "0.01",
        }

        param = MLflowEntityConverter.to_mlflow_param(param_data)

        assert isinstance(param, Param)
        assert param.key == "learning_rate"
        assert param.value == "0.01"

    def test_to_mlflow_experiment(self):
        """Test converting Model Registry experiment to MLflow Experiment."""
        experiment_data = {
            "id": "exp-123",
            "name": "test-experiment",
            "externalId": "s3://bucket/artifacts/exp-123",
            "state": "LIVE",
            "customProperties": {
                "key1": "value1"
            },
        }

        experiment = MLflowEntityConverter.to_mlflow_experiment(
            experiment_data, "s3://bucket/artifacts"
        )

        assert isinstance(experiment, Experiment)
        assert experiment.experiment_id == "exp-123"
        assert experiment.name == "test-experiment"
        assert experiment.artifact_location == "s3://bucket/artifacts/exp-123"
        assert len(experiment.tags) == 1
        assert experiment.tags["key1"] == "value1"

    def test_to_mlflow_experiment_with_default_artifact_uri(self):
        """Test converting experiment with default artifact URI."""
        experiment_data = {
            "id": "exp-123",
            "name": "test-experiment",
            "state": "LIVE",
            "customProperties": {},
        }

        experiment = MLflowEntityConverter.to_mlflow_experiment(
            experiment_data, "s3://bucket/artifacts"
        )

        assert experiment.artifact_location == "s3://bucket/artifacts"

    def test_to_mlflow_run_info(self):
        """Test converting Model Registry run to MLflow RunInfo."""
        run_data = {
            "id": "run-123",
            "experimentId": "exp-123",
            "name": "test-run",
            "status": "RUNNING",
            "owner": "user123",
            "startTimeSinceEpoch": 1234567890,
            "externalId": "s3://bucket/artifacts/exp-123/run-123",
            "customProperties": {},
        }

        run_info = MLflowEntityConverter.to_mlflow_run_info(
            run_data, "s3://bucket/artifacts/exp-123/run-123"
        )

        assert isinstance(run_info, RunInfo)
        assert run_info.run_id == "run-123"
        assert run_info.experiment_id == "exp-123"
        assert run_info.run_name == "test-run"
        assert run_info.user_id == "user123"
        assert run_info.status == RunStatus.RUNNING
        assert run_info.start_time == 1234567890
        assert run_info.artifact_uri == "s3://bucket/artifacts/exp-123/run-123"

    def test_to_mlflow_run_info_with_end_time(self):
        """Test converting run with end time."""
        run_data = {
            "id": "run-123",
            "experimentId": "exp-123",
            "status": "FINISHED",
            "state": "TERMINATED",
            "startTimeSinceEpoch": 1234567890,
            "endTimeSinceEpoch": 1234567999,
            "externalId": "s3://bucket/artifacts/exp-123/run-123",
        }

        run_info = MLflowEntityConverter.to_mlflow_run_info(
            run_data, "s3://bucket/artifacts/exp-123/run-123"
        )

        assert run_info.end_time == 1234567999

    def test_to_mlflow_run_info_with_create_time_fallback(self):
        """Test converting run with createTimeSinceEpoch fallback."""
        run_data = {
            "id": "run-123",
            "experimentId": "exp-123",
            "status": "RUNNING",
            "createTimeSinceEpoch": 1234567890,
            "externalId": "s3://bucket/artifacts/exp-123/run-123",
        }

        run_info = MLflowEntityConverter.to_mlflow_run_info(
            run_data, "s3://bucket/artifacts/exp-123/run-123"
        )

        assert run_info.start_time == 1234567890

    def test_to_mlflow_dataset_input(self):
        """Test converting Model Registry dataset to MLflow DatasetInput."""
        dataset_data = {
            "name": "test-dataset",
            "digest": "digest123",
            "sourceType": "csv",
            "source": "s3://bucket/data.csv",
            "schema": "schema",
            "profile": "profile",
            "customProperties": {
                "key1": "value1"
            },
        }

        dataset_input = MLflowEntityConverter.to_mlflow_dataset_input(dataset_data)

        assert isinstance(dataset_input, DatasetInput)
        assert dataset_input.dataset.name == "test-dataset"
        assert dataset_input.dataset.digest == "digest123"
        assert dataset_input.dataset.source_type == "csv"
        assert dataset_input.dataset.source == "s3://bucket/data.csv"
        assert dataset_input.dataset.schema == "schema"
        assert dataset_input.dataset.profile == "profile"
        assert len(dataset_input.tags) == 1
        assert dataset_input.tags[0].key == "key1"
        assert dataset_input.tags[0].value == "value1"

    def test_to_mlflow_logged_model_with_serialized_model(self):
        """Test converting logged model with serialized MLflow model data."""
        model_dict = {
            "artifact_path": "model",
            "run_id": "run-123",
            "flavors": {"python_function": {}},
            "model_uuid": "uuid-123",
            "utc_time_created": "2023-01-01T00:00:00Z",
            "mlflow_version": "2.0.0",
        }

        model_data = {
            "id": "model-123",
            "name": "test-model",
            "uri": "s3://bucket/model",
            "createTimeSinceEpoch": 1234567890,
            "lastUpdateTimeSinceEpoch": 1234567890,
            "customProperties": {
                "mlflow__logged_model": json.dumps(model_dict),
            },
        }

        logged_model = MLflowEntityConverter.to_mlflow_logged_model(model_data)

        assert isinstance(logged_model, Model)
        assert logged_model.model_uuid == "uuid-123"
        assert logged_model.artifact_path == "model"

    def test_to_mlflow_logged_model_with_invalid_serialized_data(self):
        """Test converting logged model with invalid serialized data."""
        model_data = {
            "id": "model-123",
            "name": "test-model",
            "uri": "s3://bucket/model",
            "customProperties": {
                "mlflow__logged_model": "invalid json",
            },
        }

        with pytest.raises(MlflowException) as exc_info:
            MLflowEntityConverter.to_mlflow_logged_model(model_data)

        assert "Failed to deserialize stored MLflow model" in str(exc_info.value)

    def test_to_mlflow_logged_model_with_tags_and_params(self):
        """Test converting logged model with tags and parameters."""
        model_data = {
            "id": "model-123",
            "name": "test-model",
            "uri": "s3://bucket/model",
            "createTimeSinceEpoch": 1234567890,
            "lastUpdateTimeSinceEpoch": 1234567890,
            "customProperties": {
                "mlflow__experiment_id": "exp-123",
                "mlflow__source_run_id": "run-123",
                "mlflow__model_type": "sklearn",
                "mlflow__utc_time_created": "2023-01-01T00:00:00Z",
                "tag1": "value1",
                "param_lr": "0.01",
            },
        }

        logged_model = MLflowEntityConverter.to_mlflow_logged_model(model_data)

        assert isinstance(logged_model, LoggedModel)
        assert logged_model.model_id == "model-123"
        assert logged_model.name == "test-model"
        assert logged_model.experiment_id == "exp-123"
        assert logged_model.source_run_id == "run-123"
        assert logged_model.model_type == "sklearn"
        assert logged_model.artifact_location == "s3://bucket/model"
        assert len(logged_model.tags) == 1
        assert logged_model.tags["tag1"] == "value1"
        assert len(logged_model.params) == 1
        assert logged_model.params["lr"] == "0.01"

    def test_to_mlflow_run(self):
        """Test converting Model Registry run with artifacts to MLflow Run."""
        run_data = {
            "id": "run-123",
            "experimentId": "exp-123",
            "name": "test-run",
            "status": "RUNNING",
            "owner": "user123",
            "startTimeSinceEpoch": 1234567890,
            "externalId": "s3://bucket/artifacts/exp-123/run-123",
            "customProperties": {
                "tag1": "value1"
            },
        }

        artifacts = [
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
            {
                "artifactType": "model-artifact",
                "id": "model-123",
                "customProperties": {
                    "mlflow__model_io_type": ModelIOType.OUTPUT.value,
                    "mlflow__step": "1",
                },
            },
        ]

        run = MLflowEntityConverter.to_mlflow_run(
            run_data, artifacts, "s3://bucket/artifacts/exp-123/run-123"
        )

        assert isinstance(run, Run)
        assert run.info.run_id == "run-123"
        assert run.info.experiment_id == "exp-123"
        assert run.info.run_name == "test-run"
        assert run.info.user_id == "user123"
        assert run.info.status == RunStatus.RUNNING
        assert run.info.artifact_uri == "s3://bucket/artifacts/exp-123/run-123"

        # Check metrics
        assert len(run.data.metrics) == 1
        assert "accuracy" in run.data.metrics
        assert run.data.metrics["accuracy"] == 0.95

        # Check parameters
        assert len(run.data.params) == 1
        assert run.data.params["learning_rate"] == "0.01"

        # Check tags
        assert len(run.data.tags) == 1
        assert run.data.tags["tag1"] == "value1"

        # Check dataset inputs
        assert len(run.inputs.dataset_inputs) == 1
        assert run.inputs.dataset_inputs[0].dataset.name == "test-dataset"

        # Check model outputs
        assert len(run.outputs.model_outputs) == 1
        assert run.outputs.model_outputs[0].model_id == "model-123"
        assert run.outputs.model_outputs[0].step == 1

    def test_to_mlflow_run_with_input_model(self):
        """Test converting run with input model."""
        run_data = {
            "id": "run-123",
            "experimentId": "exp-123",
            "status": "RUNNING",
            "startTimeSinceEpoch": 1234567890,
            "externalId": "s3://bucket/artifacts/exp-123/run-123",
            "customProperties": {},
        }

        artifacts = [
            {
                "artifactType": "model-artifact",
                "id": "model-123",
                "customProperties": {
                    "mlflow__model_io_type": ModelIOType.INPUT.value,
                },
            },
        ]

        run = MLflowEntityConverter.to_mlflow_run(
            run_data, artifacts, "s3://bucket/artifacts/exp-123/run-123"
        )

        assert len(run.inputs.model_inputs) == 1
        assert run.inputs.model_inputs[0].model_id == "model-123"
        assert len(run.outputs.model_outputs) == 0

    def test_to_mlflow_run_with_unknown_artifact_type(self):
        """Test converting run with unknown artifact type."""
        run_data = {
            "id": "run-123",
            "experimentId": "exp-123",
            "status": "RUNNING",
            "startTimeSinceEpoch": 1234567890,
            "externalId": "s3://bucket/artifacts/exp-123/run-123",
            "customProperties": {},
        }

        artifacts = [
            {
                "artifactType": "unknown-type",
                "name": "unknown",
            },
        ]

        run = MLflowEntityConverter.to_mlflow_run(
            run_data, artifacts, "s3://bucket/artifacts/exp-123/run-123"
        )

        # Unknown artifact types should be ignored
        assert len(run.data.metrics) == 0
        assert len(run.data.params) == 0
        assert len(run.inputs.dataset_inputs) == 0
        assert len(run.inputs.model_inputs) == 0
        assert len(run.outputs.model_outputs) == 0 