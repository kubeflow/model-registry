"""Entity converters for Model Registry to MLflow conversion."""

from __future__ import annotations

import json
from typing import TYPE_CHECKING, Any, Dict, List

if TYPE_CHECKING:
    from mlflow.entities import (
        DatasetInput,
        Experiment,
        LoggedModel,
        Metric,
        Param,
        Run,
        RunInfo,
    )

from model_registry_mlflow.utils import (
    ModelIOType,
    convert_modelregistry_state,
    convert_timestamp,
    convert_to_mlflow_logged_model_status,
)


class MLflowEntityConverter:
    """Handles conversion between Model Registry and MLflow entities."""

    @staticmethod
    def to_mlflow_metric(metric_data: Dict[str, Any]) -> Metric:
        """Convert Model Registry metric to MLflow Metric.

        Args:
            metric_data: Metric data from Model Registry

        Returns:
            MLflow Metric entity
        """
        from mlflow.entities import Metric

        return Metric(
            key=metric_data["name"],
            value=float(metric_data["value"]),
            timestamp=convert_timestamp(
                metric_data.get("timestamp") or metric_data.get("createTimeSinceEpoch")
            ),
            step=metric_data.get("step", 0),
        )

    @staticmethod
    def to_mlflow_param(param_data: Dict[str, Any]) -> Param:
        """Convert Model Registry parameter to MLflow Param.

        Args:
            param_data: Parameter data from Model Registry

        Returns:
            MLflow Param entity
        """
        from mlflow.entities import Param

        return Param(
            key=param_data["name"],
            value=str(param_data["value"]),
        )

    @staticmethod
    def to_mlflow_experiment(
        experiment_data: Dict[str, Any], artifact_uri: str
    ) -> Experiment:
        """Convert Model Registry experiment to MLflow Experiment.

        Args:
            experiment_data: Experiment data from Model Registry
            artifact_uri: Default artifact URI

        Returns:
            MLflow Experiment entity
        """
        from mlflow.entities import Experiment, ExperimentTag

        return Experiment(
            experiment_id=str(experiment_data["id"]),
            name=experiment_data["name"],
            artifact_location=experiment_data.get("externalId") or artifact_uri,
            lifecycle_stage=convert_modelregistry_state(experiment_data),
            tags=[
                ExperimentTag(k, v)
                for k, v in experiment_data.get("customProperties", {}).items()
            ],
        )

    @staticmethod
    def to_mlflow_run_info(run_data: Dict[str, Any], artifact_uri: str) -> RunInfo:
        """Convert Model Registry run to MLflow RunInfo.

        Args:
            run_data: Run data from Model Registry
            artifact_uri: Artifact URI for the run

        Returns:
            MLflow RunInfo entity
        """
        from mlflow.entities import RunInfo

        start_time = convert_timestamp(
            run_data.get("startTimeSinceEpoch") or run_data.get("createTimeSinceEpoch")
        )
        if start_time is None:
            # Provide a default timestamp if none is available
            import time

            start_time = int(time.time() * 1000)

        return RunInfo(
            run_id=str(run_data["id"]),
            experiment_id=str(run_data["experimentId"]),
            user_id=run_data.get("owner") or "unknown",
            status=run_data.get("status", "RUNNING"),
            start_time=start_time,
            end_time=convert_timestamp(run_data.get("endTimeSinceEpoch"))
            if run_data.get("state") == "TERMINATED"
            else None,
            lifecycle_stage=convert_modelregistry_state(run_data),
            artifact_uri=artifact_uri,
            run_name=run_data.get("name"),
        )

    @staticmethod
    def to_mlflow_dataset_input(dataset_data: Dict[str, Any]) -> DatasetInput:
        """Convert Model Registry dataset to MLflow DatasetInput.

        Args:
            dataset_data: Dataset data from Model Registry

        Returns:
            MLflow DatasetInput entity
        """
        from mlflow.entities import Dataset, DatasetInput, InputTag

        tags = []
        for key, value in dataset_data.get("customProperties", {}).items():
            tags.append(InputTag(key=key, value=value))

        dataset = Dataset(
            name=dataset_data["name"],
            digest=dataset_data.get("digest", ""),
            source_type=dataset_data.get("sourceType", ""),
            source=dataset_data.get("source", ""),
            schema=dataset_data.get("schema", ""),
            profile=dataset_data.get("profile", ""),
        )

        return DatasetInput(dataset=dataset, tags=tags)

    @staticmethod
    def to_mlflow_logged_model(model_data: Dict[str, Any]) -> LoggedModel:
        """Convert Model Registry model to MLflow LoggedModel.

        Args:
            model_data: Model data from Model Registry

        Returns:
            MLflow LoggedModel entity
        """
        from mlflow.entities import LoggedModel, LoggedModelParameter, LoggedModelTag
        from mlflow.exceptions import MlflowException
        from mlflow.models.model import Model

        custom_props = model_data.get("customProperties", {})

        # Check if the model has the serialized MLflow model data
        if "mlflow__logged_model" in custom_props:
            try:
                model_dict = json.loads(custom_props["mlflow__logged_model"])
                return Model.from_dict(model_dict)
            except (json.JSONDecodeError, KeyError, TypeError) as e:
                msg = f"Failed to deserialize stored MLflow model: {e}"
                raise MlflowException(msg) from e

        # Extract tags and params from customProperties
        tags = []
        params = []
        for key, value in custom_props.items():
            if key.startswith("mlflow__"):
                continue  # Ignore all mlflow__* keys
            if key.startswith("param_"):
                params.append(LoggedModelParameter(key=key[6:], value=value))
            else:
                tags.append(LoggedModelTag(key=key, value=value))

        return LoggedModel(
            model_id=str(model_data["id"]),
            experiment_id=custom_props.get("mlflow__experiment_id"),
            name=custom_props.get("mlflow__name", model_data["name"]),
            source_run_id=custom_props.get("mlflow__source_run_id"),
            artifact_location=custom_props.get(
                "mlflow__artifact_location", model_data.get("uri", "")
            ),
            creation_timestamp=convert_timestamp(
                model_data.get("createTimeSinceEpoch")
            ),
            last_updated_timestamp=convert_timestamp(
                model_data.get("lastUpdatedTimeSinceEpoch")
            ),
            model_type=custom_props.get("mlflow__model_type"),
            status=convert_to_mlflow_logged_model_status(model_data.get("state")),
            tags=tags,
            params=params,
        )

    @staticmethod
    def to_mlflow_run(
        run_data: Dict[str, Any],
        artifacts: List[Dict[str, Any]],
        artifact_uri: str,
    ) -> Run:
        """Convert Model Registry run to MLflow Run.

        Args:
            run_data: Run data from Model Registry
            artifacts: List of artifacts for the run
            artifact_uri: Artifact URI for the run

        Returns:
            MLflow Run entity
        """
        from mlflow.entities import (
            LoggedModelInput,
            LoggedModelOutput,
            Run,
            RunData,
            RunInputs,
            RunOutputs,
            RunTag,
        )

        # Convert run info
        run_info = MLflowEntityConverter.to_mlflow_run_info(run_data, artifact_uri)

        # Extract metrics, params, and tags from artifacts
        metrics = []
        params = []
        tags = []
        dataset_inputs = []
        model_inputs = []
        model_outputs = []

        for artifact in artifacts:
            artifact_type = artifact.get("artifactType", "")
            if artifact_type == "metric":
                metrics.append(MLflowEntityConverter.to_mlflow_metric(artifact))
            elif artifact_type == "parameter":
                params.append(MLflowEntityConverter.to_mlflow_param(artifact))
            elif artifact_type == "dataset-artifact":
                dataset_inputs.append(
                    MLflowEntityConverter.to_mlflow_dataset_input(artifact)
                )
            elif (
                artifact_type == "model-artifact"
                and artifact.get("customProperties", {}).get("mlflow__model_io_type")
                == ModelIOType.INPUT.value
            ):
                model_inputs.append(LoggedModelInput(model_id=artifact["id"]))
            elif artifact_type == "model-artifact":  # default to output
                custom_props = artifact.get("customProperties", {})
                step = int(custom_props.get("mlflow__step", 0))
                model_outputs.append(
                    LoggedModelOutput(model_id=artifact["id"], step=step)
                )

        # Add tags from customProperties
        for key, value in run_data.get("customProperties", {}).items():
            tags.append(RunTag(key=key, value=value))

        run_data_obj = RunData(
            metrics=metrics,
            params=params,
            tags=tags,
        )

        run_inputs = RunInputs(
            dataset_inputs=dataset_inputs,
            model_inputs=model_inputs,
        )

        run_outputs = RunOutputs(
            model_outputs=model_outputs,
        )

        return Run(
            run_info=run_info,
            run_data=run_data_obj,
            run_inputs=run_inputs,
            run_outputs=run_outputs,
        )
