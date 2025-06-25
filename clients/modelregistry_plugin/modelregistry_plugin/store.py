"""Model Registry MLflow Tracking Store Implementation."""

from __future__ import annotations

import json
import os
import uuid
from collections.abc import Sequence
from typing import TYPE_CHECKING, Any

import requests

from modelregistry_plugin.auth import get_auth_headers
from modelregistry_plugin.utils import (
    ModelIOType,
    convert_modelregistry_state,
    convert_timestamp,
    convert_to_mlflow_logged_model_status,
    convert_to_model_artifact_state,
    fromModelRegistryCustomProperties,
    parse_tracking_uri,
    toModelRegistryCustomProperties,
)

if TYPE_CHECKING:
    from mlflow.entities import (
        DatasetInput,
        Experiment,
        ExperimentTag,
        LoggedModel,
        LoggedModelInput,
        LoggedModelOutput,
        LoggedModelParameter,
        LoggedModelTag,
        Metric,
        PagedList,
        Param,
        Run,
        RunInfo,
        RunInputs,
        RunOutputs,
        RunStatus,
        RunTag,
        ViewType,
    )
    from mlflow.entities.metric import MetricWithRunId
    from mlflow.models.model import Model


class ModelRegistryStore:
    """MLflow tracking store that uses Kubeflow Model Registry as the backend."""

    # Default page size for artifact fetching
    DEFAULT_ARTIFACT_PAGE_SIZE = 1000

    def __init__(
        self, store_uri: str | None = None, artifact_uri: str | None = None
    ) -> None:
        """Initialize the Model Registry store.

        Args:
            store_uri: URI for the Model Registry (e.g., "modelregistry://localhost:8080")
            artifact_uri: URI for artifact storage (optional)
        """
        # Import MLflow modules here to avoid circular imports
        from mlflow.store.tracking.abstract_store import AbstractStore
        from mlflow.store.tracking.file_store import _default_root_dir
        from mlflow.utils.file_utils import (
            local_file_uri_to_path,
            path_to_local_file_uri,
        )

        # Initialize as AbstractStore
        AbstractStore.__init__(self)

        if store_uri:
            self.store_uri = store_uri
        else:
            self.store_uri = os.getenv(
                "MLFLOW_TRACKING_URI", "modelregistry://localhost:8080"
            )

        # Parse the tracking URI to get connection details
        self.host, self.port, self.secure = parse_tracking_uri(self.store_uri)
        self.base_url = f"{'https' if self.secure else 'http'}://{self.host}:{self.port}/api/model_registry/v1alpha3"

        if artifact_uri:
            self.artifact_uri = artifact_uri
        else:
            # FIXME use local file store by default?
            file_store_path = local_file_uri_to_path(_default_root_dir())
            self.artifact_uri = path_to_local_file_uri(file_store_path)

    def _get_artifact_location(self, response: dict) -> str:
        return response.get("externalId") or self.artifact_uri

    def _getMLflowMetric(self, metric_data: dict) -> Metric:
        """Create an MLflow Metric entity from Model Registry metric data."""
        # Import MLflow modules here to avoid circular imports
        from mlflow.entities import Metric

        return Metric(
            key=metric_data["name"],
            value=float(metric_data["value"]),
            timestamp=int(
                metric_data.get("timestamp") or metric_data.get("createTimeSinceEpoch")
            ),
            step=metric_data.get("step") or 0,
        )

    def _getMLflowParam(self, param_data: dict) -> Param:
        """Create an MLflow Param entity from Model Registry parameter data."""
        # Import MLflow modules here to avoid circular imports
        from mlflow.entities import Param

        return Param(key=param_data["name"], value=str(param_data["value"]))

    def _get_all_run_artifacts(self, run_id: str) -> list[dict]:
        """Get all artifacts for a run with pagination support.

        Args:
            run_id: The ID of the run to fetch artifacts for

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

            response = self._request(
                "GET", f"/experiment_runs/{run_id}/artifacts", params=params
            )

            items = response.get("items", [])
            all_artifacts.extend(items)

            # Check for next page
            next_page_token = response.get("nextPageToken")
            if not next_page_token or next_page_token == "":
                break
            page_token = next_page_token

        return all_artifacts

    def _request(self, method: str, endpoint: str, **kwargs) -> dict:
        """Make authenticated request to Model Registry API."""
        # Import MLflow modules here to avoid circular imports
        from mlflow.exceptions import MlflowException, get_error_code

        url = f"{self.base_url}/{endpoint.lstrip('/')}"
        headers = get_auth_headers()
        headers.update(kwargs.pop("headers", {}))

        # convert customProperties to ModelRegistry customProperties format
        json_data = kwargs.get("json")
        if json_data is not None:
            toModelRegistryCustomProperties(json_data)

        response = requests.request(method, url, headers=headers, **kwargs)

        response_json = response.json()
        if not response.ok:
            try:
                error_detail = response_json.get("message", response.text)
            except (AttributeError, TypeError, KeyError):
                error_detail = response.text
            msg = f"Model Registry API error: {error_detail}"
            raise MlflowException(
                msg,
                error_code=get_error_code(response.status_code),
            )

        # convert ModelRegistry customProperties format back to MLflow customProperties format
        if response_json.get("items"):
            for item in response_json.get("items"):
                fromModelRegistryCustomProperties(item)
        else:
            fromModelRegistryCustomProperties(response_json)
        return response_json

    # Async logging methods copied from mlflow.store.tracking.abstract_store.py
    def log_batch_async(
        self,
        run_id: str,
        metrics: list[Metric],
        params: list[Param],
        tags: list[RunTag],
    ) -> Any:
        """Log multiple metrics, params, and tags for the specified run in async fashion.
        This API does not offer immediate consistency of the data. When API returns,
        data is accepted but not persisted/processed by back end. Data would be processed
        in near real time fashion.

        Args:
            run_id: String id for the run.
            metrics: List of :py:class:`mlflow.entities.Metric` instances to log.
            params: List of :py:class:`mlflow.entities.Param` instances to log.
            tags: List of :py:class:`mlflow.entities.RunTag` instances to log.

        Returns:
            An :py:class:`mlflow.utils.async_logging.run_operations.RunOperations` instance
            that represents future for logging operation.
        """
        if not self._async_logging_queue.is_active():
            self._async_logging_queue.activate()

        return self._async_logging_queue.log_batch_async(
            run_id=run_id, metrics=metrics, params=params, tags=tags
        )

    def end_async_logging(self) -> None:
        """Ends the async logging queue. This method is a no-op if the queue is not active. This is
        different from flush as it just stops the async logging queue from accepting
        new data (moving the queue state TEAR_DOWN state), but flush will ensure all data
        is processed before returning (moving the queue to IDLE state).
        """
        if self._async_logging_queue.is_active():
            self._async_logging_queue.end_async_logging()

    def flush_async_logging(self) -> None:
        """Flushes the async logging queue. This method is a no-op if the queue is already
        at IDLE state. This methods also shutdown the logging worker threads.
        After flushing, logging thread is setup again.
        """
        if not self._async_logging_queue.is_idle():
            self._async_logging_queue.flush()

    def shut_down_async_logging(self) -> None:
        """Shuts down the async logging queue. This method is a no-op if the queue is already
        at IDLE state. This methods also shutdown the logging worker threads.
        """
        if not self._async_logging_queue.is_idle():
            self._async_logging_queue.shut_down_async_logging()

    # Experiment operations
    def create_experiment(
        self,
        name: str,
        artifact_location: str | None = None,
        tags: list[ExperimentTag] | None = None,
    ) -> str:
        """Create a new experiment in Model Registry."""
        payload = {
            "name": name,
            "description": f"MLflow experiment: {name}",
            "state": "LIVE",
            "customProperties": {},
        }

        # Set externalId based on artifact_location or default pattern
        if artifact_location:
            payload["externalId"] = artifact_location
        elif self.artifact_uri:
            # We'll set this after getting the experiment ID
            pass
        else:
            # No artifact_uri available, set to None
            payload["externalId"] = None

        if tags:
            for tag in tags:
                payload["customProperties"][tag.key] = tag.value

        experiment_data = self._request("POST", "/experiments", json=payload)
        experiment_id = str(experiment_data["id"])

        # If no artifact_location was provided but we have artifact_uri, update with the default pattern
        if not artifact_location and self.artifact_uri:
            default_artifact_location = (
                f"{self.artifact_uri}/experiments/{experiment_id}"
            )
            update_payload = {"externalId": default_artifact_location}
            self._request("PATCH", f"/experiments/{experiment_id}", json=update_payload)

        return experiment_id

    def get_experiment(self, experiment_id: str) -> Experiment:
        """Get experiment by ID."""
        # Import MLflow modules here to avoid circular imports
        from mlflow.entities import Experiment, ExperimentTag

        experiment_data = self._request("GET", f"/experiments/{experiment_id}")

        return Experiment(
            experiment_id=str(experiment_data["id"]),
            name=experiment_data["name"],
            artifact_location=self._get_artifact_location(experiment_data),
            lifecycle_stage=convert_modelregistry_state(experiment_data),
            tags=[
                ExperimentTag(k, v)
                for k, v in experiment_data.get("customProperties", {}).items()
            ],
        )

    def get_experiment_by_name(self, experiment_name: str) -> Experiment | None:
        """Get experiment by name."""
        # Import MLflow modules here to avoid circular imports
        from mlflow.entities import Experiment, ExperimentTag
        from mlflow.exceptions import MlflowException

        try:
            exp_data = self._request(
                "GET", "/experiment", params={"name": experiment_name}
            )
            return Experiment(
                experiment_id=str(exp_data["id"]),
                name=exp_data["name"],
                artifact_location=self._get_artifact_location(exp_data),
                lifecycle_stage=convert_modelregistry_state(exp_data),
                tags=[
                    ExperimentTag(k, v)
                    for k, v in exp_data.get("customProperties", {}).items()
                ],
            )
        except MlflowException as e:
            if e.get_http_status_code() == 404 and "not found" in e.message:
                return None
            raise e

    def delete_experiment(self, experiment_id: str) -> None:
        """Delete an experiment."""
        # Model Registry doesn't support deletion, so we mark as archived
        payload = {"state": "ARCHIVED"}
        self._request("PATCH", f"/experiments/{experiment_id}", json=payload)

    def restore_experiment(self, experiment_id: str) -> None:
        """Restore a deleted experiment."""
        payload = {"state": "LIVE"}
        self._request("PATCH", f"/experiments/{experiment_id}", json=payload)

    # TODO this won't work until ModelRegistry supports mutable resource names
    def rename_experiment(self, experiment_id: str, new_name: str) -> None:
        """Rename an experiment."""
        payload = {"name": new_name}
        self._request("PATCH", f"/experiments/{experiment_id}", json=payload)

    def list_experiments(
        self,
        view_type: ViewType | None = None,
        max_results: int | None = None,
        page_token: str | None = None,
    ) -> PagedList[Experiment]:
        """List experiments."""
        # Import MLflow modules here to avoid circular imports
        from mlflow.entities import Experiment, ExperimentTag, LifecycleStage, ViewType
        from mlflow.store.entities.paged_list import PagedList

        if view_type is None:
            view_type = ViewType.ACTIVE_ONLY

        params = {}
        if max_results:
            params["pageSize"] = max_results
        if page_token:
            params["pageToken"] = page_token

        response_data = self._request("GET", "/experiments", params=params)
        items = response_data.get("items", [])

        experiments = []
        for exp_data in items:
            lifecycle_stage = convert_modelregistry_state(exp_data)

            # TODO add filtering in model registry server
            if (
                view_type == ViewType.ACTIVE_ONLY
                and lifecycle_stage == LifecycleStage.DELETED
            ) or (
                view_type == ViewType.DELETED_ONLY
                and lifecycle_stage == LifecycleStage.ACTIVE
            ):
                continue

            experiments.append(
                Experiment(
                    experiment_id=str(exp_data["id"]),
                    name=exp_data["name"],
                    artifact_location=self._get_artifact_location(exp_data),
                    lifecycle_stage=lifecycle_stage,
                    tags=[
                        ExperimentTag(k, v)
                        for k, v in exp_data.get("customProperties", {}).items()
                    ],
                )
            )

        nextPageToken = response_data.get("nextPageToken")
        return PagedList(experiments, nextPageToken if nextPageToken != "" else None)

    def search_experiments(
        self,
        view_type: ViewType | None = None,
        max_results: int = 1000,  # TODO: Import SEARCH_MAX_RESULTS_DEFAULT
        filter_string: str | None = None,
        order_by: list[str] | None = None,
        page_token: str | None = None,
    ) -> PagedList[Experiment]:
        """Search for experiments that match the specified search query.

        Args:
            view_type: One of enum values ACTIVE_ONLY, DELETED_ONLY, or ALL
            max_results: Maximum number of experiments desired
            filter_string: Filter query string (not supported in Model Registry yet)
            order_by: List of columns to order by (not supported in Model Registry yet)
            page_token: Token specifying the next page of results

        Returns:
            A PagedList of Experiment objects
        """
        # Import MLflow modules here to avoid circular imports
        from mlflow.entities import Experiment, ExperimentTag, LifecycleStage, ViewType
        from mlflow.store.entities.paged_list import PagedList

        if view_type is None:
            view_type = ViewType.ACTIVE_ONLY

        # TODO: Add support for filter_string and order_by in Model Registry API
        if filter_string:
            # For now, we'll ignore filter_string as Model Registry doesn't support it yet
            pass

        if order_by:
            # For now, we'll ignore order_by as Model Registry doesn't support it yet
            pass

        params = {}
        if max_results:
            params["pageSize"] = max_results
        if page_token:
            params["pageToken"] = page_token

        response_data = self._request("GET", "/experiments", params=params)
        items = response_data.get("items", [])

        experiments = []
        for exp_data in items:
            lifecycle_stage = convert_modelregistry_state(exp_data)

            # Filter by view_type
            if (
                view_type == ViewType.ACTIVE_ONLY
                and lifecycle_stage == LifecycleStage.DELETED
            ) or (
                view_type == ViewType.DELETED_ONLY
                and lifecycle_stage == LifecycleStage.ACTIVE
            ):
                continue

            experiments.append(
                Experiment(
                    experiment_id=str(exp_data["id"]),
                    name=exp_data["name"],
                    artifact_location=self._get_artifact_location(exp_data),
                    lifecycle_stage=lifecycle_stage,
                    tags=[
                        ExperimentTag(k, v)
                        for k, v in exp_data.get("customProperties", {}).items()
                    ],
                )
            )

        nextPageToken = response_data.get("nextPageToken")
        return PagedList(experiments, nextPageToken if nextPageToken != "" else None)

    # Run operations
    def create_run(
        self,
        experiment_id: str,
        user_id: str | None = None,
        start_time: int | None = None,
        tags: list[RunTag] | None = None,
        run_name: str | None = None,
    ) -> Run:
        """Create a new run."""
        # Import MLflow modules here to avoid circular imports
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

        run_data = self._request("POST", "/experiment_runs", json=payload)
        run_id = str(run_data["id"])

        # Get the experiment to determine its externalId
        experiment_data = self._request("GET", f"/experiments/{experiment_id}")

        # Set the artifact location for the run using experiment's externalId as prefix
        artifact_location = self._get_artifact_location(experiment_data)
        if artifact_location:
            artifact_location = f"{artifact_location}/{run_id}"
            update_payload = {"externalId": artifact_location}
            self._request("PATCH", f"/experiment_runs/{run_id}", json=update_payload)

        run_info = RunInfo(
            run_id=run_id,
            experiment_id=experiment_id,
            user_id=user_id or "unknown",
            status=RunStatus.RUNNING,
            start_time=start_time
            or convert_timestamp(run_data.get("createTimeSinceEpoch")),
            end_time=None,
            lifecycle_stage=LifecycleStage.ACTIVE,
            artifact_uri=artifact_location,
            run_name=run_name,
        )

        # get tags from run_data
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
        """Get run by ID."""
        run_data = self._request("GET", f"/experiment_runs/{run_id}")

        # Get metrics, parameters, tags, inputs, outputs
        return self._getMLflowRun(run_data)

    def update_run_info(
        self,
        run_id: str,
        run_status: RunStatus | None = None,
        end_time: int | None = None,
        run_name: str | None = None,
    ) -> RunInfo:
        """Update run information."""
        # Import MLflow modules here to avoid circular imports
        from mlflow.entities import RunInfo, RunStatus

        payload = {}
        if run_status:
            payload["status"] = RunStatus.to_string(run_status)
        if end_time:
            payload["endTimeSinceEpoch"] = str(end_time)
        if run_name:
            payload["name"] = run_name

        run_data = self._request("PATCH", f"/experiment_runs/{run_id}", json=payload)

        return RunInfo(
            run_id=str(run_data["id"]),
            experiment_id=str(run_data["experimentId"]),
            user_id=run_data.get("owner") or "unknown",
            status=RunStatus.from_string(run_data.get("status", "RUNNING")),
            start_time=convert_timestamp(
                run_data.get("startTimeSinceEpoch")
                or run_data.get("createTimeSinceEpoch")
            ),
            end_time=convert_timestamp(run_data.get("endTimeSinceEpoch")),
            lifecycle_stage=convert_modelregistry_state(run_data),
            artifact_uri=self._get_artifact_location(run_data),
            run_name=run_data.get("name"),
        )

    def delete_run(self, run_id: str) -> None:
        """Delete a run."""
        payload = {"state": "ARCHIVED"}
        self._request("PATCH", f"/experiment_runs/{run_id}", json=payload)

    def restore_run(self, run_id: str) -> None:
        """Restore a deleted run."""
        payload = {"state": "LIVE"}
        self._request("PATCH", f"/experiment_runs/{run_id}", json=payload)

    # Metric operations
    def log_metric(self, run_id: str, metric: Metric) -> None:
        """Log a metric for a run."""
        # Import MLflow modules here to avoid circular imports
        from mlflow.utils.time import get_current_time_millis

        payload = {
            "artifactType": "metric",
            "name": metric.key,
            "value": metric.value,
            "step": metric.step or 0,
            "timestamp": str(metric.timestamp or get_current_time_millis()),
            "customProperties": {},
        }
        self._request("POST", f"/experiment_runs/{run_id}/artifacts", json=payload)

    def _get_run_metrics(self, run_id: str) -> list[Metric]:
        """Get all metrics for a run."""
        # Import MLflow modules here to avoid circular imports

        items = self._request(
            "GET",
            f"/experiment_runs/{run_id}/artifacts",
            params={"artifactType": "metric"},
        ).get("items", [])

        metrics = []
        for metric_data in items:
            metrics.append(self._getMLflowMetric(metric_data))
        return metrics

    def get_metric_history(
        self,
        run_id: str,
        metric_key: str,
        max_results: int | None = None,
        page_token: str | None = None,
    ) -> PagedList[Metric]:
        """Return a list of metric objects corresponding to all values logged for a given metric
        within a run.

        Args:
            run_id: Unique identifier for run
            metric_key: Metric name within the run
            max_results: Maximum number of metric history events to return per paged query
            page_token: Token specifying the next paginated set of results

        Returns:
            A PagedList of Metric entities if logged, else empty PagedList
        """
        # Import MLflow modules here to avoid circular imports
        from mlflow.store.entities.paged_list import PagedList

        params = {"name": metric_key}
        if max_results:
            params["pageSize"] = max_results
        if page_token:
            params["pageToken"] = page_token

        # Get all metrics for the run
        response_data = self._request(
            "GET", f"/experiment_runs/{run_id}/metric_history", params=params
        )
        next_page_token = response_data.get("nextPageToken")
        items = response_data.get("items", [])
        metrics = []
        for metric_data in items:
            metrics.append(self._getMLflowMetric(metric_data))
        return PagedList(metrics, next_page_token if next_page_token != "" else None)

    # NOTE: Copied from mlflow.store.tracking.abstract_store.py
    def get_metric_history_bulk_interval_from_steps(
        self, run_id: str, metric_key: str, steps: list[int], max_results: int | None
    ) -> PagedList[MetricWithRunId]:
        """Return a list of metric objects corresponding to all values logged
        for a given metric within a run for the specified steps.

        Args:
            run_id: Unique identifier for run.
            metric_key: Metric name within the run.
            steps: List of steps for which to return metrics.
            max_results: Maximum number of metric history events (steps) to return.

        Returns:
            A list of MetricWithRunId objects:
                - key: Metric name within the run.
                - value: Metric value.
                - timestamp: Metric timestamp.
                - step: Metric step.
                - run_id: Unique identifier for run.
        """
        # Import MLflow modules here to avoid circular imports
        from mlflow.entities.metric import MetricWithRunId
        from mlflow.store.entities.paged_list import PagedList

        params = {"name": metric_key}
        if max_results:
            params["pageSize"] = max_results
        if steps:
            params["stepIds"] = ",".join(str(step) for step in steps)

        response_data = self._request(
            "GET", f"/experiment_runs/{run_id}/metric_history", params=params
        )
        next_page_token = response_data.get("nextPageToken")
        items = response_data.get("items", [])
        metrics = []
        for metric_data in items:
            metrics.append(
                MetricWithRunId(
                    metric=self._getMLflowMetric(metric_data),
                    run_id=run_id,
                )
            )
        return PagedList(metrics, next_page_token if next_page_token != "" else None)

    # Parameter operations
    def log_param(self, run_id: str, param: Param) -> None:
        """Log a parameter for a run."""
        # Import MLflow modules here to avoid circular imports

        payload = {
            "artifactType": "parameter",
            "name": param.key,
            "value": param.value,
            "parameterType": "string",  # since MLflow doesn't provide the type, default to string
        }
        self._request("POST", f"/experiment_runs/{run_id}/artifacts", json=payload)

    def _get_run_params(self, run_id: str) -> list[Param]:
        """Get all parameters for a run."""
        # Import MLflow modules here to avoid circular imports

        items = self._request(
            "GET",
            f"/experiment_runs/{run_id}/artifacts",
            params={"artifactType": "parameter"},
        ).get("items", [])

        params = []
        for param_data in items:
            params.append(self._getMLflowParam(param_data))
        return params

    def _get_run_inputs_outputs(self, run_id: str) -> tuple[RunInputs, RunOutputs]:
        """Get all inputs and outputs for a run (datasets and models)."""
        # Import MLflow modules here to avoid circular imports
        from mlflow.entities import (
            LoggedModelInput,
            LoggedModelOutput,
            RunInputs,
            RunOutputs,
        )

        # Get all artifacts and filter by type
        all_artifacts = self._get_all_run_artifacts(run_id)

        dataset_inputs = []
        input_models = []
        output_models = []

        for artifact_data in all_artifacts:
            artifact_type = artifact_data.get("artifactType")

            if artifact_type == "dataset-artifact":
                # Create DatasetInput entity using helper method
                dataset_input = self._getMLflowDatasetInput(artifact_data)
                dataset_inputs.append(dataset_input)
            elif artifact_type == "model-artifact":
                # Check the model's io_type to determine if it's input or output
                custom_props = artifact_data.get("customProperties", {})
                model_id = artifact_data["id"]
                io_type = custom_props.get("mlflow__model_io_type")
                step = int(
                    custom_props.get("mlflow__step") or 0
                )  # FIXME set this when logged model is created

                if io_type == ModelIOType.INPUT.value:
                    # Create LoggedModelInput entity
                    model_input = LoggedModelInput(model_id=model_id)
                    input_models.append(model_input)
                else:  # default to output
                    # Create LoggedModelOutput entity
                    model_output = LoggedModelOutput(model_id=model_id, step=step)
                    output_models.append(model_output)

        return RunInputs(
            dataset_inputs=dataset_inputs, model_inputs=input_models
        ), RunOutputs(model_outputs=output_models)

    def _getMLflowDatasetInput(self, dataset_data: dict) -> DatasetInput:
        """Create an MLflow DatasetInput entity from Model Registry dataset data."""
        # Import MLflow modules here to avoid circular imports
        from mlflow.entities import Dataset, DatasetInput, LoggedModelTag

        tags = []
        for key, value in dataset_data.get("customProperties", {}).items():
            tags.append(LoggedModelTag(key=key, value=value))

        # Create Dataset entity
        dataset = Dataset(
            name=dataset_data["name"],
            digest=dataset_data.get("digest", ""),
            source_type=dataset_data.get("sourceType", ""),
            source=dataset_data.get("source", ""),
            schema=dataset_data.get("schema", ""),
            profile=dataset_data.get("profile", ""),
        )

        # Create and return DatasetInput entity
        return DatasetInput(dataset=dataset, tags=tags)

    def _getMLflowLoggedModel(self, model_data: dict) -> LoggedModel:
        """Create an MLflow LoggedModel entity from Model Registry model data."""
        # Import MLflow modules here to avoid circular imports
        from mlflow.entities import LoggedModel, LoggedModelParameter, LoggedModelTag
        from mlflow.exceptions import MlflowException
        from mlflow.models.model import Model

        custom_props = model_data.get("customProperties", {})

        # Check if the model has the serialized MLflow model data
        if "mlflow__logged_model" in custom_props:
            try:
                # Deserialize the stored model dictionary
                model_dict = json.loads(custom_props["mlflow__logged_model"])

                # Create and return the MLflow Model object
                return Model.from_dict(model_dict)
            except (json.JSONDecodeError, KeyError, TypeError) as e:
                msg = f"Failed to deserialize stored MLflow model: {e}"
                raise MlflowException(msg) from e

        # Extract tags and params from customProperties
        tags = []
        params = []
        for key, value in custom_props.items():
            if key.startswith("mlflow__"):
                # Skip mlflow__* keys as they are used for internal tracking
                continue
            if key.startswith("param_"):
                params.append(LoggedModelParameter(key=key[6:], value=value))
            else:
                tags.append(LoggedModelTag(key=key, value=value))

        return LoggedModel(
            model_id=str(model_data["id"]),
            experiment_id=custom_props.get("mlflow__experiment_id"),
            name=model_data["name"],
            source_run_id=custom_props.get("mlflow__source_run_id"),
            artifact_location=model_data.get("uri", ""),
            creation_timestamp=convert_timestamp(
                model_data.get("mlflow__utc_time_created")
                or model_data.get("createTimeSinceEpoch")
            ),
            last_updated_timestamp=convert_timestamp(
                model_data.get("lastUpdateTimeSinceEpoch")
            ),
            model_type=custom_props.get("mlflow__model_type"),
            status=convert_to_mlflow_logged_model_status(model_data.get("state")),
            tags=tags,
            params=params,
        )

    def log_batch(
        self,
        run_id: str,
        metrics: Sequence[Metric] = (),
        params: Sequence[Param] = (),
        tags: Sequence[RunTag] = (),
    ) -> None:
        """Log a batch of metrics, parameters, and tags."""
        # Get current run to preserve other properties
        run_data = self._request("GET", f"/experiment_runs/{run_id}")
        custom_props = run_data.get("customProperties", {}) or {}
        for tag in tags:
            custom_props[tag.key] = tag.value
        payload = {"customProperties": custom_props}
        self._request("PATCH", f"/experiment_runs/{run_id}", json=payload)
        # iterate and log metrics and params
        # TODO add support for batch logging in Model Registry REST API
        for metric in metrics:
            self.log_metric(run_id, metric)
        for param in params:
            self.log_param(run_id, param)

    def log_inputs(
        self,
        run_id: str,
        datasets: list[DatasetInput] | None = None,
        models: list[LoggedModelInput] | None = None,
    ) -> None:
        """Log inputs for a run.

        Args:
            run_id: The ID of the run to log inputs for
            datasets: List of dataset inputs
            models: List of logged model inputs
        """
        if datasets:
            for datasetInput in datasets:
                payload = {
                    "artifactType": "dataset-artifact",
                    "name": datasetInput.dataset.name,
                    "digest": datasetInput.dataset.digest,
                    "sourceType": datasetInput.dataset.source_type,
                    "source": datasetInput.dataset.source,
                    "schema": datasetInput.dataset.schema,
                    "profile": datasetInput.dataset.profile,
                    "customProperties": {},
                }
                if datasetInput.tags:
                    for tag in datasetInput.tags:
                        payload["customProperties"][tag.key] = tag.value
                self._request(
                    "POST", f"/experiment_runs/{run_id}/artifacts", json=payload
                )

        if models:
            for model in models:
                # Get current model to preserve other properties
                model_data = self._request("GET", f"/artifacts/{model.model_id}")
                custom_props = model_data.get("customProperties", {})
                custom_props["mlflow__model_io_type"] = ModelIOType.INPUT.value

                payload = {
                    "artifactType": "model-artifact",
                    "id": model.model_id,
                    "customProperties": custom_props,
                }
                self._request(
                    "POST", f"/experiment_runs/{run_id}/artifacts", json=payload
                )

    def log_outputs(self, run_id: str, models: list[LoggedModelOutput]) -> None:
        """Log outputs for a run.

        Args:
            run_id: The ID of the run to log outputs for
            models: List of logged model outputs
        """
        for model in models:
            # Get current model to preserve other properties
            model_data = self._request("GET", f"/artifacts/{model.model_id}")
            custom_props = model_data.get("customProperties", {})
            custom_props["mlflow__model_io_type"] = ModelIOType.OUTPUT.value

            payload = {
                "artifactType": "model-artifact",
                "id": model.model_id,
                "customProperties": custom_props,
            }
            self._request("POST", f"/experiment_runs/{run_id}/artifacts", json=payload)

    def record_logged_model(self, run_id: str, mlflow_model: Model) -> None:
        """Record a logged model.

        Args:
            run_id: The ID of the run to record the model for
            mlflow_model: The MLflow model to record
        """
        # Import MLflow modules here to avoid circular imports

        model: Model = mlflow_model
        model_info = model.get_model_info()
        model_id = model.model_id
        model_uuid = model.model_uuid
        model_dict = model.to_dict()
        model_name = (
            model_dict.get("name") or model_uuid or str(uuid.uuid4())
        )  # TODO: check the naming convention

        # Create a model artifact in Model Registry
        payload = {
            "artifactType": "model-artifact",
            "name": model_name,
            "uri": model_info.model_uri,  # TODO: also set the externalId
            "customProperties": {
                "mlflow__artifactPath": model_info.artifact_path,
                "mlflow__model_uuid": model_uuid,
                "mlflow__utc_time_created": model_info.utc_time_created,
                "mlflow__mlflow_version": model_info.mlflow_version,
                "mlflow__flavor": str(model_info.flavors),
                "mlflow__source_run_id": run_id or model_info.run_id,
            },
        }
        if model_id:
            payload["id"] = model_id

        # Store the full model dict as a tag for backward compatibility
        payload["customProperties"]["mlflow__logged_model"] = json.dumps(model_dict)

        # Create the model artifact
        self._request("POST", f"/experiment_runs/{run_id}/artifacts", json=payload)

    def create_logged_model(
        self,
        experiment_id: str,
        name: str | None = None,
        source_run_id: str | None = None,
        tags: list[LoggedModelTag] | None = None,
        params: list[LoggedModelParameter] | None = None,
        model_type: str | None = None,
    ) -> LoggedModel:
        """Create a new logged model.

        Args:
            experiment_id: ID of the experiment to which the model belongs
            name: Name of the model. If not specified, a random name will be generated
            source_run_id: ID of the run that produced the model
            tags: Tags to set on the model
            params: Parameters to set on the model
            model_type: Type of the model

        Returns:
            The created LoggedModel object
        """
        experiment_data = self._request("GET", f"/experiments/{experiment_id}")
        artifact_location = self._get_artifact_location(experiment_data)

        payload = {
            "artifactType": "model-artifact",
            "name": name or str(uuid.uuid4()),
            "customProperties": {
                "mlflow__model_type": model_type or "unknown",
                "mlflow__experiment_id": experiment_id,
                "mlflow__source_run_id": source_run_id,
            },
        }
        # TODO: check whether this is correct for mlflow
        if artifact_location:
            if source_run_id:
                artifact_location = f"{artifact_location}/{source_run_id}/{name}"
            else:
                artifact_location = f"{artifact_location}/{name}"
            payload["uri"] = artifact_location

        if tags:
            for tag in tags:
                payload["customProperties"][tag.key] = tag.value

        if params:
            for param in params:
                payload["customProperties"][f"param_{param.key}"] = param.value

        # TODO source_run_id is optional, but we need to handle it
        model_data = self._request(
            "POST", f"/experiment_runs/{source_run_id}/artifacts", json=payload
        )

        # Use the helper method to create LoggedModel entity
        return self._getMLflowLoggedModel(model_data)

    def search_logged_models(
        self,
        experiment_ids: list[str],
        filter_string: str | None = None,
        datasets: list[dict[str, Any]] | None = None,
        max_results: int | None = None,
        order_by: list[dict[str, Any]] | None = None,
        page_token: str | None = None,
    ) -> PagedList[LoggedModel]:
        """Search for logged models that match the specified search criteria.

        Args:
            experiment_ids: List of experiment ids to scope the search
            filter_string: A search filter string
            datasets: List of dictionaries specifying datasets for metric filters
            max_results: Maximum number of logged models desired
            order_by: List of dictionaries specifying result ordering
            page_token: Token specifying the next page of results

        Returns:
            A PagedList of LoggedModel objects
        """
        # Import MLflow modules here to avoid circular imports
        from mlflow.store.entities.paged_list import PagedList

        params = {"artifactType": "model-artifact", "experimentIds": experiment_ids}

        # TODO add support for filter_string in ModelRegistry API
        # TODO add support for datasets filtering in ModelRegistry API
        # TODO add support for mlflow order_by mapping to ModelRegistry API
        # TODO add support for pagination in ModelRegistry API across list of experiments
        if max_results:
            params["pageSize"] = str(max_results)
        if page_token:
            params["pageToken"] = page_token

        # iterate over experiment_ids and get all runs, and get all model-artifacts for each run
        models = []
        for experiment_id in experiment_ids:
            response = self._request(
                "GET", f"/experiments/{experiment_id}/experiment_runs", params=params
            )
            runs_data = response.get("items", [])
            for run_data in runs_data:
                items = self._request(
                    "GET",
                    f"/experiment_runs/{run_data['id']}/artifacts",
                    params={"artifactType": "model-artifact"},
                ).get("items", [])
                for item in items:
                    models.append(self.get_logged_model(item["id"]))

        return PagedList(models, None)  # no paging across experiments

    def finalize_logged_model(self, model_id: str, status: Any) -> LoggedModel:
        """Finalize a model by updating its status.

        Args:
            model_id: ID of the model to finalize
            status: Final status to set on the model

        Returns:
            The updated LoggedModel
        """
        payload = {"state": convert_to_model_artifact_state(status)}

        model_data = self._request("PATCH", f"/artifacts/{model_id}", json=payload)
        # FIXME Model Regisistry is not returning artifact state in the response
        # so we hack it for now
        response = self._getMLflowLoggedModel(model_data)
        response.status = status

        # Use the helper method to create LoggedModel entity
        return response

    def set_logged_model_tags(self, model_id: str, tags: list[LoggedModelTag]) -> None:
        """Set tags on the specified logged model.

        Args:
            model_id: ID of the model
            tags: Tags to set on the model
        """
        # Get current model to preserve other properties
        model_data = self._request("GET", f"/artifacts/{model_id}")
        custom_props = model_data.get("customProperties", {})

        # Update custom properties with new tags
        for tag in tags:
            custom_props[tag.key] = tag.value

        payload = {"artifactType": "model-artifact", "customProperties": custom_props}
        self._request("PATCH", f"/artifacts/{model_id}", json=payload)

    def delete_logged_model_tag(self, model_id: str, key: str) -> None:
        """Delete a tag from the specified logged model.

        Args:
            model_id: ID of the model
            key: Key of the tag to delete
        """
        # Get current model to preserve other properties
        model_data = self._request("GET", f"/artifacts/{model_id}")
        custom_props = model_data.get("customProperties", {})

        # Remove the specified tag
        if key in custom_props:
            del custom_props[key]

        payload = {"artifactType": "model-artifact", "customProperties": custom_props}
        self._request("PATCH", f"/artifacts/{model_id}", json=payload)

    def get_logged_model(self, model_id: str) -> LoggedModel:
        """Fetch the logged model with the specified ID.

        Args:
            model_id: ID of the model to fetch

        Returns:
            The fetched LoggedModel
        """
        model_data = self._request("GET", f"/artifacts/{model_id}")

        # Use the helper method to create LoggedModel entity
        return self._getMLflowLoggedModel(model_data)

    def delete_logged_model(self, model_id: str) -> None:
        """Delete the logged model with the specified ID.

        Args:
            model_id: ID of the model to delete
        """
        # Model Registry doesn't support deletion, so we mark as archived
        payload = {
            "artifactType": "model-artifact",
            "customProperties": {
                "state": "MARKED_FOR_DELETION"  # TODO: handle ModelArtifactState.MARKED_FOR_DELETION in ModelRegistry
            },
        }
        self._request("PATCH", f"/artifacts/{model_id}", json=payload)

    # Tag operations
    def set_experiment_tag(self, experiment_id: str, tag: ExperimentTag) -> None:
        """Set a tag on an experiment."""
        # Get current experiment to preserve other properties
        experiment = self._request("GET", f"/experiments/{experiment_id}")
        custom_props = experiment.get("customProperties", {})
        custom_props[tag.key] = tag.value

        payload = {"customProperties": custom_props}
        self._request("PATCH", f"/experiments/{experiment_id}", json=payload)

    def set_tag(self, run_id: str, tag: RunTag) -> None:
        """Set a tag on a run."""
        # Get current run to preserve other properties
        run = self._request("GET", f"/experiment_runs/{run_id}")
        custom_props = run.get("customProperties", {})
        custom_props[tag.key] = tag.value

        payload = {"customProperties": custom_props}
        self._request("PATCH", f"/experiment_runs/{run_id}", json=payload)

    def delete_tag(self, run_id: str, key: str) -> None:
        """Delete a tag from a run."""
        run = self._request("GET", f"/experiment_runs/{run_id}")
        custom_props = run.get("customProperties", {})
        if key in custom_props:
            del custom_props[key]

        payload = {"customProperties": custom_props}
        self._request("PATCH", f"/experiment_runs/{run_id}", json=payload)

    # Search operations (simplified implementation)
    def search_runs(
        self,
        experiment_ids: list[str],
        filter_string: str = "",
        run_view_type: ViewType | None = None,
        max_results: int = 1000,
        order_by: list[str] | None = None,
        page_token: str | None = None,
    ) -> PagedList[Run]:
        """Search for runs."""
        # Import MLflow modules here to avoid circular imports
        from mlflow.entities import ViewType
        from mlflow.store.entities.paged_list import PagedList

        if run_view_type is None:
            run_view_type = ViewType.ACTIVE_ONLY

        all_runs = []

        # TODO add support for filter_string in ModelRegistry API
        for experiment_id in experiment_ids:
            response = self._request(
                "GET",
                f"/experiments/{experiment_id}/experiment_runs",
                params={"pageSize": str(min(max_results, 1000))},
            )
            runs_data = response.get("items", [])

            for run_data in runs_data:
                # avoid nested calls for run data by filtering out runs that are not active or deleted
                if (
                    run_view_type == ViewType.ACTIVE_ONLY
                    and run_data.get("state") == "ARCHIVED"
                ) or (
                    run_view_type == ViewType.DELETED_ONLY
                    and run_data.get("state") == "ACTIVE"
                ):
                    continue
                run = self._getMLflowRun(run_data)

                all_runs.append(run)

        return PagedList(all_runs, None)  # no paging across experiments

    def _getMLflowRun(self, run_data: dict) -> Run:
        # Import MLflow modules here to avoid circular imports
        from mlflow.entities import (
            LoggedModelInput,
            LoggedModelOutput,
            Run,
            RunData,
            RunInfo,
            RunInputs,
            RunOutputs,
            RunStatus,
            RunTag,
        )

        run_id = run_data["id"]

        # Get all artifacts in a single API call
        all_artifacts = self._get_all_run_artifacts(run_id)

        # Process artifacts based on their artifactType discriminator
        metrics = []
        params = []
        dataset_inputs = []
        input_models = []
        output_models = []

        for artifact_data in all_artifacts:
            artifact_type = artifact_data.get("artifactType")

            if artifact_type == "metric":
                metrics.append(self._getMLflowMetric(artifact_data))
            elif artifact_type == "parameter":
                params.append(self._getMLflowParam(artifact_data))
            elif artifact_type == "dataset-artifact":
                dataset_input = self._getMLflowDatasetInput(artifact_data)
                dataset_inputs.append(dataset_input)
            elif artifact_type == "model-artifact":
                # Check the model's io_type to determine if it's input or output
                custom_props = artifact_data.get("customProperties", {})
                model_id = artifact_data["id"]
                io_type = custom_props.get("mlflow__model_io_type")
                step = int(
                    custom_props.get("mlflow__step") or 0
                )  # FIXME set this when logged model is created

                if io_type == ModelIOType.INPUT.value:
                    # Create LoggedModelInput entity
                    model_input = LoggedModelInput(model_id=model_id)
                    input_models.append(model_input)
                else:  # default to output
                    # Create LoggedModelOutput entity
                    model_output = LoggedModelOutput(model_id=model_id, step=step)
                    output_models.append(model_output)

        # Create RunInputs and RunOutputs
        run_inputs = RunInputs(dataset_inputs=dataset_inputs, model_inputs=input_models)
        run_outputs = RunOutputs(model_outputs=output_models)

        tags = [RunTag(k, v) for k, v in run_data.get("customProperties", {}).items()]
        run_info = RunInfo(
            run_id=str(run_data["id"]),
            experiment_id=str(run_data["experimentId"]),
            user_id=run_data.get("owner") or "unknown",
            status=RunStatus.from_string(run_data.get("status", "RUNNING")),
            start_time=convert_timestamp(
                run_data.get("startTimeSinceEpoch")
                or run_data.get("createTimeSinceEpoch")
            ),
            end_time=convert_timestamp(run_data.get("endTimeSinceEpoch"))
            if run_data.get("state") == "TERMINATED"
            else None,
            lifecycle_stage=convert_modelregistry_state(run_data),
            artifact_uri=self._get_artifact_location(run_data),
            run_name=run_data.get("name"),
        )
        run_data_obj = RunData(metrics=metrics, params=params, tags=tags)
        return Run(
            run_info=run_info,
            run_inputs=run_inputs,
            run_outputs=run_outputs,
            run_data=run_data_obj,
        )
