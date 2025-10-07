"""Run operations for Model Registry store."""

from __future__ import annotations

import os
from typing import TYPE_CHECKING, List, Optional, Sequence

if TYPE_CHECKING:
    from mlflow.entities import (
        DatasetInput,
        LoggedModelInput,
        LoggedModelOutput,
        Metric,
        Param,
        Run,
        RunInfo,
        RunStatus,
        RunTag,
    )

from ..api_client import ModelRegistryAPIClient
from ..converters import MLflowEntityConverter
from ..utils import ModelIOType


class RunOperations:
    """Handles all run-related operations."""

    # Default page size for artifact fetching
    DEFAULT_ARTIFACT_PAGE_SIZE = 1000

    def __init__(self, api_client: ModelRegistryAPIClient, artifact_uri: str):
        """Initialize run operations.

        Args:
            api_client: API client for making requests
            artifact_uri: Default artifact URI
        """
        self.api_client = api_client
        self.artifact_uri = artifact_uri

    def create_run(
        self,
        experiment_id: str,
        user_id: Optional[str] = None,
        start_time: Optional[int] = None,
        tags: Optional[List[RunTag]] = None,
        run_name: Optional[str] = None,
    ) -> Run:
        """Create a new run.

        Args:
            experiment_id: ID of the experiment
            user_id: User ID
            start_time: Start time
            tags: Tags to set on the run
            run_name: Name of the run

        Returns:
            Run entity
        """
        from mlflow.entities import (
            LifecycleStage,
            Run,
            RunData,
            RunInfo,
            RunInputs,
            RunOutputs,
            RunStatus,
            RunTag,
        )
        from mlflow.utils.time import get_current_time_millis

        payload = {
            "experimentId": experiment_id,
            "name": run_name or f"run-{start_time or 0}",
            "description": f"MLflow run in experiment {experiment_id}",
            "startTimeSinceEpoch": str(start_time or get_current_time_millis()),
            "status": "RUNNING",
            "customProperties": {},
        }

        if user_id:
            payload["owner"] = user_id

        if tags:
            for tag in tags:
                payload["customProperties"][tag.key] = tag.value

        run_data = self.api_client.post("/experiment_runs", json=payload)
        run_id = str(run_data["id"])

        # Get the experiment to determine its externalId
        experiment_data = self.api_client.get(f"/experiments/{experiment_id}")

        # Set the artifact location for the run using experiment's externalId as prefix
        artifact_location = experiment_data.get("externalId") or self.artifact_uri
        if artifact_location:
            artifact_location = f"{artifact_location}/{run_id}"
            update_payload = {"externalId": artifact_location}
            self.api_client.patch(f"/experiment_runs/{run_id}", json=update_payload)

        run_info = RunInfo(
            run_id=run_id,
            experiment_id=experiment_id,
            user_id=user_id or "unknown",
            status=RunStatus.to_string(RunStatus.RUNNING),
            start_time=start_time or run_data.get("createTimeSinceEpoch"),
            end_time=None,
            lifecycle_stage=LifecycleStage.ACTIVE,
            artifact_uri=artifact_location,
            run_name=run_name,
        )

        # Get tags from run_data
        run_tags = [
            RunTag(k, v) for k, v in run_data.get("customProperties", {}).items()
        ]

        return Run(
            run_info=run_info,
            run_inputs=RunInputs(dataset_inputs=[], model_inputs=[]),
            run_outputs=RunOutputs(model_outputs=[]),
            run_data=RunData(tags=run_tags),
        )

    def get_run(self, run_id: str) -> Run:
        """Get run by ID.

        Args:
            run_id: ID of the run

        Returns:
            Run entity
        """
        run_data = self.api_client.get(f"/experiment_runs/{run_id}")
        all_artifacts = self._get_all_run_artifacts(run_id)
        artifact_location = run_data.get("externalId") or self.artifact_uri
        return MLflowEntityConverter.to_mlflow_run(
            run_data, all_artifacts, artifact_location
        )

    def update_run_info(
        self,
        run_id: str,
        run_status: Optional[RunStatus] = None,
        end_time: Optional[int] = None,
        run_name: Optional[str] = None,
    ) -> RunInfo:
        """Update run information.

        Args:
            run_id: ID of the run
            run_status: New run status
            end_time: End time
            run_name: New run name

        Returns:
            Updated RunInfo
        """
        from mlflow.entities import RunStatus

        payload = {}
        if run_status:
            payload["status"] = RunStatus.to_string(run_status)
        if end_time:
            payload["endTimeSinceEpoch"] = str(end_time)
        if run_name:
            payload["name"] = run_name

        run_data = self.api_client.patch(f"/experiment_runs/{run_id}", json=payload)
        artifact_location = run_data.get("externalId") or self.artifact_uri

        return MLflowEntityConverter.to_mlflow_run_info(run_data, artifact_location)

    def delete_run(self, run_id: str) -> None:
        """Delete a run.

        Args:
            run_id: ID of the run to delete
        """
        payload = {"state": "ARCHIVED"}
        self.api_client.patch(f"/experiment_runs/{run_id}", json=payload)

    def restore_run(self, run_id: str) -> None:
        """Restore a deleted run.

        Args:
            run_id: ID of the run to restore
        """
        payload = {"state": "LIVE"}
        self.api_client.patch(f"/experiment_runs/{run_id}", json=payload)

    def log_metric(self, run_id: str, metric: Metric) -> None:
        """Log a metric for a run.

        Args:
            run_id: ID of the run
            metric: Metric to log
        """
        from mlflow.utils.time import get_current_time_millis

        payload = {
            "artifactType": "metric",
            "name": metric.key,
            "value": metric.value,
            "step": metric.step or 0,
            "timestamp": str(metric.timestamp or get_current_time_millis()),
            "customProperties": {},
        }
        self.api_client.post(f"/experiment_runs/{run_id}/artifacts", json=payload)

    def log_param(self, run_id: str, param: Param) -> None:
        """Log a parameter for a run.

        Args:
            run_id: ID of the run
            param: Parameter to log
        """
        payload = {
            "artifactType": "parameter",
            "name": param.key,
            "value": param.value,
            "parameterType": "string",
        }
        self.api_client.post(f"/experiment_runs/{run_id}/artifacts", json=payload)

    def log_batch(
        self,
        run_id: str,
        metrics: Sequence[Metric] = (),
        params: Sequence[Param] = (),
        tags: Sequence[RunTag] = (),
    ) -> None:
        """Log a batch of metrics, parameters, and tags.

        Args:
            run_id: ID of the run
            metrics: Metrics to log
            params: Parameters to log
            tags: Tags to log
        """
        # Get current run to preserve other properties
        run_data = self.api_client.get(f"/experiment_runs/{run_id}")
        custom_props = run_data.get("customProperties", {}) or {}
        for tag in tags:
            custom_props[tag.key] = tag.value
        payload = {"customProperties": custom_props}
        self.api_client.patch(f"/experiment_runs/{run_id}", json=payload)

        # Log metrics and params individually
        # TODO: Add support for batch logging in Model Registry REST API
        for metric in metrics:
            self.log_metric(run_id, metric)
        for param in params:
            self.log_param(run_id, param)

    def log_inputs(
        self,
        run_id: str,
        datasets: Optional[List[DatasetInput]] = None,
        models: Optional[List[LoggedModelInput]] = None,
    ) -> None:
        """Log inputs for a run.

        Args:
            run_id: ID of the run
            datasets: Dataset inputs to log
            models: Model inputs to log
        """
        if datasets:
            for dataset_input in datasets:
                payload = {
                    "artifactType": "dataset-artifact",
                    "name": dataset_input.dataset.name,
                    "digest": dataset_input.dataset.digest,
                    "sourceType": dataset_input.dataset.source_type,
                    "source": dataset_input.dataset.source,
                    "schema": dataset_input.dataset.schema,
                    "profile": dataset_input.dataset.profile,
                    "customProperties": {},
                }
                if dataset_input.tags:
                    for tag in dataset_input.tags:
                        payload["customProperties"][tag.key] = tag.value
                self.api_client.post(
                    f"/experiment_runs/{run_id}/artifacts", json=payload
                )

        if models:
            for model in models:
                # Get current model to preserve other properties
                model_data = self.api_client.get(f"/artifacts/{model.model_id}")
                custom_props = model_data.get("customProperties", {})
                custom_props["mlflow__model_io_type"] = ModelIOType.INPUT.value

                payload = {
                    "artifactType": "model-artifact",
                    "id": model.model_id,
                    "customProperties": custom_props,
                }
                self.api_client.post(
                    f"/experiment_runs/{run_id}/artifacts", json=payload
                )

    def log_outputs(self, run_id: str, models: List[LoggedModelOutput]) -> None:
        """Log outputs for a run.

        Args:
            run_id: ID of the run
            models: Model outputs to log
        """
        for model in models:
            # Get current model to preserve other properties
            model_data = self.api_client.get(f"/artifacts/{model.model_id}")
            custom_props = model_data.get("customProperties", {})
            custom_props["mlflow__model_io_type"] = ModelIOType.OUTPUT.value

            payload = {
                "artifactType": "model-artifact",
                "id": model.model_id,
                "customProperties": custom_props,
            }
            self.api_client.post(f"/experiment_runs/{run_id}/artifacts", json=payload)

    def set_tag(self, run_id: str, tag: RunTag) -> None:
        """Set a tag on a run.

        Args:
            run_id: ID of the run
            tag: Tag to set
        """
        run = self.api_client.get(f"/experiment_runs/{run_id}")
        custom_props = run.get("customProperties", {})
        custom_props[tag.key] = tag.value

        payload = {"customProperties": custom_props}
        self.api_client.patch(f"/experiment_runs/{run_id}", json=payload)

    def delete_tag(self, run_id: str, key: str) -> None:
        """Delete a tag from a run.

        Args:
            run_id: ID of the run
            key: Key of the tag to delete
        """
        from mlflow.exceptions import MlflowException
        from mlflow.protos.databricks_pb2 import RESOURCE_DOES_NOT_EXIST

        run = self.api_client.get(f"/experiment_runs/{run_id}")
        custom_props = run.get("customProperties", {})

        if key not in custom_props:
            raise MlflowException(
                f"No tag with name: {key} in run with id {run_id}",
                error_code=RESOURCE_DOES_NOT_EXIST,
            )

        del custom_props[key]
        payload = {"customProperties": custom_props}
        self.api_client.patch(f"/experiment_runs/{run_id}", json=payload)

    def _get_all_run_artifacts(self, run_id: str) -> List[dict]:
        """Get all artifacts for a run with pagination support.

        Args:
            run_id: ID of the run

        Returns:
            List of all artifact data dictionaries
        """
        # Get page size from environment variable or use default
        page_size = int(
            os.getenv(
                "MODEL_REGISTRY_ARTIFACT_PAGE_SIZE",
                str(self.DEFAULT_ARTIFACT_PAGE_SIZE),
            )
        )

        all_artifacts = []
        page_token = None

        while True:
            params = {"pageSize": page_size}
            if page_token:
                params["pageToken"] = page_token

            response = self.api_client.get(
                f"/experiment_runs/{run_id}/artifacts", params=params
            )

            items = response.get("items", [])
            all_artifacts.extend(items)

            # Check for next page
            next_page_token = response.get("nextPageToken")
            if not next_page_token or next_page_token == "":
                break
            page_token = next_page_token

        return all_artifacts
