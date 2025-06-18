"""
Model Registry MLflow Tracking Store Implementation
"""

import os
import json

import mlflow
import requests
from typing import List, Optional, Dict, Any, Sequence

from mlflow.store.tracking.abstract_store import AbstractStore
from mlflow.entities import (
    Experiment, Run, RunInfo, RunData, RunStatus, RunTag, Param, Metric,
    ViewType, LifecycleStage, ExperimentTag, DatasetInput, LoggedModelInput, LoggedModelOutput
)
from mlflow.exceptions import MlflowException
from mlflow.protos.databricks_pb2 import INVALID_PARAMETER_VALUE, RESOURCE_DOES_NOT_EXIST
from mlflow.utils import time

from .auth import get_auth_headers
from .utils import parse_tracking_uri, convert_timestamp, convert_modelregistry_state

class ModelRegistryStore(AbstractStore):
    """
    MLflow tracking store that uses Kubeflow Model Registry as the backend.
    """
    
    def __init__(self, store_uri: str = None, artifact_uri: str = None):
        """
        Initialize the Model Registry store.
        
        Args:
            store_uri: URI for the Model Registry (e.g., "modelregistry://localhost:8080")
            artifact_uri: URI for artifact storage (optional)
        """
        super().__init__()
        
        if store_uri:
            self.store_uri = store_uri
        else:
            self.store_uri = os.getenv("MLFLOW_TRACKING_URI", "modelregistry://localhost:8080")
            
        # Parse the tracking URI to get connection details
        self.host, self.port, self.secure = parse_tracking_uri(self.store_uri)
        self.base_url = f"{'https' if self.secure else 'http'}://{self.host}:{self.port}/api/model_registry/v1alpha3"
        
        self.artifact_uri = artifact_uri

    def _get_artifact_location(self, response: Dict) -> str:
        return response["externalId"] or self.artifact_uri

    def _request(self, method: str, endpoint: str, **kwargs) -> requests.Response:
        """Make authenticated request to Model Registry API."""
        url = f"{self.base_url}/{endpoint.lstrip('/')}"
        headers = get_auth_headers()
        headers.update(kwargs.pop("headers", {}))
        
        response = requests.request(method, url, headers=headers, **kwargs)
        
        if not response.ok:
            try:
                error_detail = response.json().get("message", response.text)
            except:
                error_detail = response.text
            raise MlflowException(
                f"Model Registry API error: {error_detail}",
                error_code=INVALID_PARAMETER_VALUE # TODO: map HTTP status code to MLflow error code
            )
        
        return response
    
    # Experiment operations
    def create_experiment(self, name: str, artifact_location: str = None, tags: List[ExperimentTag] = None) -> str:
        """Create a new experiment in Model Registry."""
        payload = {
            "name": name,
            "description": f"MLflow experiment: {name}",
            "externalId": artifact_location or self.artifact_uri,
            "state": "LIVE",
            "customProperties": {}
        }
        
        if tags:
            for tag in tags:
                payload["customProperties"][tag.key] = tag.value
                
        response = self._request("POST", "/experiments", json=payload)
        experiment_data = response.json()
        return str(experiment_data["id"])
    
    def get_experiment(self, experiment_id: str) -> Experiment:
        """Get experiment by ID."""
        response = self._request("GET", f"/experiments/{experiment_id}")
        experiment_data = response.json()
        
        return Experiment(
            experiment_id=str(experiment_data["id"]),
            name=experiment_data["name"],
            artifact_location=self._get_artifact_location(experiment_data),
            lifecycle_stage=convert_modelregistry_state(experiment_data),
            tags=[ExperimentTag(k, v) for k, v in experiment_data.get("customProperties", {}).items()]
        )
    
    def get_experiment_by_name(self, experiment_name: str) -> Optional[Experiment]:
        """Get experiment by name."""
        response = self._request("GET", "/experiments", params={"name": experiment_name})
        experiments = response.json().get("experiments", [])
        
        for exp_data in experiments:
            if exp_data["name"] == experiment_name:
                return Experiment(
                    experiment_id=str(exp_data["id"]),
                    name=exp_data["name"],
                    artifact_location=self._get_artifact_location(exp_data),
                    lifecycle_stage=convert_modelregistry_state(exp_data),
                    tags=[ExperimentTag(k, v) for k, v in exp_data.get("customProperties", {}).items()]
                )
        return None
    
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
    
    def list_experiments(self, view_type: ViewType = ViewType.ACTIVE_ONLY, max_results: int = None, page_token: str = None) -> List[Experiment]:
        """List experiments."""
        params = {}
        if max_results:
            params["pageSize"] = max_results
        if page_token:
            params["pageToken"] = page_token

        response = self._request("GET", "/experiments", params=params)
        items = response.json().get("items", [])
        
        experiments = []
        for exp_data in items:
            lifecycle_stage = convert_modelregistry_state(exp_data)

            # TODO add filtering in model registry server
            if view_type == ViewType.ACTIVE_ONLY and lifecycle_stage == LifecycleStage.DELETED:
                continue
            elif view_type == ViewType.DELETED_ONLY and lifecycle_stage == LifecycleStage.ACTIVE:
                continue
                
            experiments.append(Experiment(
                experiment_id=str(exp_data["id"]),
                name=exp_data["name"],
                artifact_location=self._get_artifact_location(exp_data),
                lifecycle_stage=lifecycle_stage,
                tags=[ExperimentTag(k, v) for k, v in exp_data.get("customProperties", {}).items()]
            ))
        
        return experiments

    # Run operations
    def create_run(self, experiment_id: str, user_id: str = None, start_time: int = None, tags: List[RunTag] = None, run_name: str = None) -> Run:
        """Create a new run."""
        payload = {
            "experimentId": experiment_id,
            "name": run_name or f"run-{start_time or 0}",
            "description": f"MLflow run in experiment {experiment_id}",
            "startTimeSinceEpoch": str(start_time or time.get_current_time_millis()),
            "status": "RUNNING",
            "customProperties": {}
        }

        if user_id:
            payload["owner"] = user_id

        if tags:
            for tag in tags:
                payload["customProperties"][tag.key] = tag.value
                
        response = self._request("POST", "/experiment_runs", json=payload)
        run_data = response.json()
        
        run_info = RunInfo(
            run_id=str(run_data["id"]),
            experiment_id=experiment_id,
            user_id=user_id or "unknown",
            status=RunStatus.RUNNING,
            start_time=start_time or convert_timestamp(run_data.get("createTimeSinceEpoch")),
            end_time=None,
            lifecycle_stage=LifecycleStage.ACTIVE,
            artifact_uri=self._get_artifact_location(run_data),
            run_name=run_name
        )
        
        return Run(run_info=run_info, run_data=RunData())
    
    def get_run(self, run_id: str) -> Run:
        """Get run by ID."""
        response = self._request("GET", f"/experiment_runs/{run_id}")
        run_data = response.json()
        
        # Get metrics, parameters, and tags
        return self._getMLflowRun(run_data)
    
    def update_run_info(self, run_id: str, run_status: RunStatus, end_time: int = None, run_name: str = None) -> RunInfo:
        """Update run information."""
        payload = {}
        if run_status:
            payload["status"] = RunStatus.to_string(run_status)
        if end_time:
            payload["endTimeSinceEpoch"] = str(end_time)
        if run_name:
            payload["name"] = run_name
            
        response = self._request("PATCH", f"/experiment_runs/{run_id}", json=payload)
        run_data = response.json()
        
        return RunInfo(
            run_id=str(run_data["id"]),
            experiment_id=str(run_data["experimentId"]),
            user_id=run_data["owner"] or "unknown",
            status=RunStatus.from_string(run_data.get("status", "RUNNING")),
            start_time=convert_timestamp(run_data.get("startTimeSinceEpoch") or run_data.get("createTimeSinceEpoch")),
            end_time=convert_timestamp(run_data.get("endTimeSinceEpoch")),
            lifecycle_stage=convert_modelregistry_state(run_data),
            artifact_uri=self._get_artifact_location(run_data),
            run_name=run_data.get("name")
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
        payload = {
            "artifactType": "metric",
            "name": metric.key,
            "value": metric.value,
            "step": str(metric.step or 0),
            "timestamp": str(metric.timestamp or time.get_current_time_millis()),
            "customProperties": {},
        }
        self._request("POST", f"/experiment_runs/{run_id}/artifacts", json=payload)
    
    def _get_run_metrics(self, run_id: str) -> List[Metric]:
        """Get all metrics for a run."""
        response = self._request("GET", f"/experiment_runs/{run_id}/artifacts",
                                 params={"artifactType": "metric"})
        items = response.json().get("items", [])
        
        metrics = []
        for metric_data in items:
            metrics.append(Metric(
                key=metric_data["name"],
                value=float(metric_data["value"]),
                timestamp=int(metric_data.get("timestamp") or metric_data.get("createTimeSinceEpoch")),
                step=int(metric_data.get("step", "0")),
            ))
        return metrics
    
    # Parameter operations  
    def log_param(self, run_id: str, param: Param) -> None:
        """Log a parameter for a run."""
        payload = {
            "artifactType": "parameter",
            "name": param.key,
            "value": param.value,
            "parameterType": "string", # since MLflow doesn't provide the type, default to string
        }
        self._request("POST", f"/experiment_runs/{run_id}/artifacts", json=payload)
    
    def _get_run_params(self, run_id: str) -> List[Param]:
        """Get all parameters for a run."""
        response = self._request("GET", f"/experiment_runs/{run_id}/artifacts",
                                 params={"artifactType": "parameter"})
        items = response.json().get("items", [])
        
        params = []
        for param_data in items:
            params.append(Param(
                key=param_data["name"],
                value=str(param_data["value"])
            ))
        return params
    
    # Tag operations
    def set_experiment_tag(self, experiment_id: str, tag: ExperimentTag) -> None:
        """Set a tag on an experiment."""
        # Get current experiment to preserve other properties
        experiment = self.get_experiment(experiment_id)
        custom_props = {t.key: t.value for t in experiment.tags}
        custom_props[tag.key] = tag.value
        
        payload = {"customProperties": custom_props}
        self._request("PATCH", f"/experiments/{experiment_id}", json=payload)
    
    def set_tag(self, run_id: str, tag: RunTag) -> None:
        """Set a tag on a run."""
        # Get current run to preserve other properties
        run = self.get_run(run_id)
        custom_props = {t.key: t.value for t in run.data.tags}
        custom_props[tag.key] = tag.value
        
        payload = {"customProperties": custom_props}
        self._request("PATCH", f"/experiment_runs/{run_id}", json=payload)
    
    def delete_tag(self, run_id: str, key: str) -> None:
        """Delete a tag from a run."""
        run = self.get_run(run_id)
        custom_props = {t.key: t.value for t in run.data.tags if t.key != key}
        
        payload = {"customProperties": custom_props}
        self._request("PATCH", f"/experiment_runs/{run_id}", json=payload)
    
    # Search operations (simplified implementation)
    def search_runs(self, experiment_ids: List[str], filter_string: str = "", 
                   run_view_type: ViewType = ViewType.ACTIVE_ONLY, max_results: int = 1000, 
                   order_by: List[str] = None, page_token: str = None) -> List[Run]:
        """Search for runs."""
        all_runs = []

        # TODO add support for filter_string in ModelRegistry API
        for experiment_id in experiment_ids:
            params = {"parentResourceId": experiment_id}
            if max_results:
                params["pageSize"] = str(min(max_results, 1000))
                
            response = self._request("GET", "/experiment_runs", params=params)
            runs_data = response.json().get("items", [])
            
            for run_data in runs_data:
                run = self._getMLflowRun(run_data)
                all_runs.append(run)
                
        return all_runs[:max_results] if max_results else all_runs

    def _getMLflowRun(self, run_data):
        run_id = run_data["id"]
        metrics = self._get_run_metrics(run_id)
        params = self._get_run_params(run_id)
        tags = [RunTag(k, v) for k, v in run_data.get("customProperties", {}).items()]
        run_info = RunInfo(
            run_id=str(run_data["id"]),
            experiment_id=str(run_data["experimentId"]),
            user_id=run_data["owner"] or "unknown",
            status=RunStatus.from_string(run_data.get("state", "RUNNING")),
            start_time=convert_timestamp(run_data.get("startTimeSinceEpoch") or run_data.get("createTimeSinceEpoch")),
            end_time=convert_timestamp(run_data.get("endTimeSinceEpoch")) if run_data.get(
                "state") == "TERMINATED" else None,
            lifecycle_stage=convert_modelregistry_state(run_data),
            artifact_uri=self._get_artifact_location(run_data),
            run_name=run_data.get("name")
        )
        run_data_obj = RunData(metrics=metrics, params=params, tags=tags)
        run = Run(run_info=run_info, run_data=run_data_obj)
        return run

    def log_batch(self, run_id: str, metrics: Sequence[Metric] = (), 
                 params: Sequence[Param] = (), tags: Sequence[RunTag] = ()) -> None:
        """Log a batch of metrics, parameters, and tags."""
        # Get current run to preserve other properties
        run = self.get_run(run_id)
        custom_props = {t.key: t.value for t in run.data.tags}
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

    def log_inputs(self, run_id: str, datasets: Optional[list[DatasetInput]] = None,
                   models: Optional[list[LoggedModelInput]] = None,) -> None:
        """Log inputs for a run.
        
        Args:
            run_id: The ID of the run to log inputs for
            inputs: Dictionary of input names to their values
        """
        for input_name, input_value in inputs.items():
            payload = {
                "artifactType": "input",
                "name": input_name,
                "value": str(input_value),
                "customProperties": {
                    "input_type": type(input_value).__name__
                }
            }
            self._request("POST", f"/experiment_runs/{run_id}/artifacts", json=payload)

    def log_outputs(self, run_id: str, models: list[LoggedModelOutput]) -> None:
        """Log outputs for a run.
        
        Args:
            run_id: The ID of the run to log outputs for
            outputs: Dictionary of output names to their values
        """
        for output_name, output_value in outputs.items():
            payload = {
                "artifactType": "output",
                "name": output_name,
                "value": str(output_value),
                "customProperties": {
                    "output_type": type(output_value).__name__
                }
            }
            self._request("POST", f"/experiment_runs/{run_id}/artifacts", json=payload)

    def record_logged_model(self, run_id: str, mlflow_model) -> None:
        """Record a logged model."""
        # This would integrate with Model Registry's model versioning
        # For now, we'll store it as a tag
        model_tag = RunTag("mlflow.logged_model", json.dumps(mlflow_model.to_dict()))
        self.set_tag(run_id, model_tag)
