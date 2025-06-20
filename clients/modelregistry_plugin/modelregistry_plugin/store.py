"""
Model Registry MLflow Tracking Store Implementation
"""

import json
import os
from typing import List, Optional, Dict, Sequence, Any
import uuid

import requests

from .auth import get_auth_headers
from .utils import convert_to_mlflow_logged_model_status, convert_to_model_artifact_state, parse_tracking_uri, \
    convert_timestamp, convert_modelregistry_state, ModelIOType, toModelRegistryCustomProperties, \
    fromModelRegistryCustomProperties


class ModelRegistryStore:
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
        # Import MLflow modules here to avoid circular imports
        from mlflow.store.tracking.abstract_store import AbstractStore
        
        # Initialize as AbstractStore
        AbstractStore.__init__(self)
        
        if store_uri:
            self.store_uri = store_uri
        else:
            self.store_uri = os.getenv("MLFLOW_TRACKING_URI", "modelregistry://localhost:8080")
            
        # Parse the tracking URI to get connection details
        self.host, self.port, self.secure = parse_tracking_uri(self.store_uri)
        self.base_url = f"{'https' if self.secure else 'http'}://{self.host}:{self.port}/api/model_registry/v1alpha3"
        
        self.artifact_uri = artifact_uri

    def _get_artifact_location(self, response: Dict) -> str:
        return response.get("externalId") or self.artifact_uri

    def _request(self, method: str, endpoint: str, **kwargs) -> requests.Response:
        """Make authenticated request to Model Registry API."""
        # Import MLflow exceptions locally to avoid circular imports
        from mlflow.exceptions import MlflowException, get_error_code
        
        url = f"{self.base_url}/{endpoint.lstrip('/')}"
        headers = get_auth_headers()
        headers.update(kwargs.pop("headers", {}))

        # convert customProperties to ModelRegistry customProperties format
        json = kwargs.get("json", None)
        if json is not None:
            toModelRegistryCustomProperties(json)

        response = requests.request(method, url, headers=headers, **kwargs)

        response_json = response.json()
        if not response.ok:
            try:
                error_detail = response_json.get("message", response.text)
            except:
                error_detail = response.text
            raise MlflowException(
                f"Model Registry API error: {error_detail}",
                error_code=get_error_code(response.status_code)
            )
        
        # convert ModelRegistry customProperties format back to MLflow customProperties format
        if response_json.get("items"):
            for item in response_json.get("items"):
                fromModelRegistryCustomProperties(item)
        else:
            fromModelRegistryCustomProperties(response_json)
        return response

    # Experiment operations
    def create_experiment(self, name: str, artifact_location: str = None, tags: List = None) -> str:
        """Create a new experiment in Model Registry."""
        # Import MLflow entities locally to avoid circular imports
        from mlflow.entities import ExperimentTag
        
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
    
    def get_experiment(self, experiment_id: str):
        """Get experiment by ID."""
        # Import MLflow entities locally to avoid circular imports
        from mlflow.entities import Experiment, ExperimentTag
        
        response = self._request("GET", f"/experiments/{experiment_id}")
        experiment_data = response.json()
        
        return Experiment(
            experiment_id=str(experiment_data["id"]),
            name=experiment_data["name"],
            artifact_location=self._get_artifact_location(experiment_data),
            lifecycle_stage=convert_modelregistry_state(experiment_data),
            tags=[ExperimentTag(k, v) for k, v in experiment_data.get("customProperties", {}).items()]
        )
    
    def get_experiment_by_name(self, experiment_name: str):
        """Get experiment by name."""
        # Import MLflow entities locally to avoid circular imports
        from mlflow.entities import Experiment, ExperimentTag
        
        try:
            response = self._request("GET", "/experiment", params={"name": experiment_name})
            exp_data = response.json()
            return Experiment(
                experiment_id=str(exp_data["id"]),
                name=exp_data["name"],
                artifact_location=self._get_artifact_location(exp_data),
                lifecycle_stage=convert_modelregistry_state(exp_data),
                tags=[ExperimentTag(k, v) for k, v in exp_data.get("customProperties", {}).items()]
            )
        except:
            # TODO look for not found error
            pass
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
    
    def list_experiments(self, view_type=None, max_results: int = None, page_token: str = None):
        """List experiments."""
        # Import MLflow entities locally to avoid circular imports
        from mlflow.entities import Experiment, ExperimentTag, ViewType, LifecycleStage
        
        if view_type is None:
            view_type = ViewType.ACTIVE_ONLY
            
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

    def search_experiments(
        self,
        view_type=None,
        max_results=1000,  # TODO: Import SEARCH_MAX_RESULTS_DEFAULT
        filter_string=None,
        order_by=None,
        page_token=None,
    ):
        """
        Search for experiments that match the specified search query.
        
        Args:
            view_type: One of enum values ACTIVE_ONLY, DELETED_ONLY, or ALL
            max_results: Maximum number of experiments desired
            filter_string: Filter query string (not supported in Model Registry yet)
            order_by: List of columns to order by (not supported in Model Registry yet)
            page_token: Token specifying the next page of results
            
        Returns:
            A PagedList of Experiment objects
        """
        # Import MLflow entities locally to avoid circular imports
        from mlflow.entities import Experiment, ExperimentTag, ViewType, LifecycleStage
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

        response = self._request("GET", "/experiments", params=params)
        response_data = response.json()
        items = response_data.get("items", [])
        
        experiments = []
        for exp_data in items:
            lifecycle_stage = convert_modelregistry_state(exp_data)

            # Filter by view_type
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
        
        return PagedList(experiments, response_data.get("nextPageToken"))

    # Run operations
    def create_run(self, experiment_id: str, user_id: str = None, start_time: int = None, tags: List = None, run_name: str = None):
        """Create a new run."""
        # Import MLflow entities locally to avoid circular imports
        from mlflow.entities import Run, RunInfo, RunData, RunStatus, RunTag, LifecycleStage
        from mlflow.utils import time
        
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
    
    def get_run(self, run_id: str):
        """Get run by ID."""
        response = self._request("GET", f"/experiment_runs/{run_id}")
        run_data = response.json()
        
        # Get metrics, parameters, and tags
        return self._getMLflowRun(run_data)
    
    def update_run_info(self, run_id: str, run_status=None, end_time: int = None, run_name: str = None):
        """Update run information."""
        # Import MLflow entities locally to avoid circular imports
        from mlflow.entities import RunInfo, RunStatus
        
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
    def log_metric(self, run_id: str, metric) -> None:
        """Log a metric for a run."""
        payload = {
            "artifactType": "metric",
            "name": metric.key,
            "value": metric.value,
            "step": metric.step or 0,
            "timestamp": str(metric.timestamp or time.get_current_time_millis()),
            "customProperties": {},
        }
        self._request("POST", f"/experiment_runs/{run_id}/artifacts", json=payload)
    
    def _get_run_metrics(self, run_id: str):
        """Get all metrics for a run."""
        # Import MLflow entities locally to avoid circular imports
        from mlflow.entities import Metric
        
        response = self._request("GET", f"/experiment_runs/{run_id}/artifacts",
                                 params={"artifactType": "metric"})
        items = response.json().get("items", [])
        
        metrics = []
        for metric_data in items:
            metrics.append(Metric(
                key=metric_data["name"],
                value=float(metric_data["value"]),
                timestamp=int(metric_data.get("timestamp") or metric_data.get("createTimeSinceEpoch")),
                step=metric_data.get("step") or 0,
            ))
        return metrics
    
    def get_metric_history(self, run_id, metric_key, max_results=None, page_token=None):
        """
        Return a list of metric objects corresponding to all values logged for a given metric
        within a run.
        
        Args:
            run_id: Unique identifier for run
            metric_key: Metric name within the run
            max_results: Maximum number of metric history events to return per paged query
            page_token: Token specifying the next paginated set of results
            
        Returns:
            A list of Metric entities if logged, else empty list
        """
        # Get all metrics for the run
        # TODO use metric history API instead
        all_metrics = self._get_run_metrics(run_id)
        
        # Filter by metric key
        filtered_metrics = [metric for metric in all_metrics if metric.key == metric_key]
        
        # Sort by timestamp and step for consistent ordering
        filtered_metrics.sort(key=lambda m: (m.timestamp, m.step))
        
        # Apply pagination if max_results is specified
        if max_results is not None:
            # TODO: Implement proper pagination with page_token
            # For now, just limit the results
            filtered_metrics = filtered_metrics[:max_results]
        
        return filtered_metrics
    
    # Parameter operations  
    def log_param(self, run_id: str, param) -> None:
        """Log a parameter for a run."""
        payload = {
            "artifactType": "parameter",
            "name": param.key,
            "value": param.value,
            "parameterType": "string", # since MLflow doesn't provide the type, default to string
        }
        self._request("POST", f"/experiment_runs/{run_id}/artifacts", json=payload)
    
    def _get_run_params(self, run_id: str):
        """Get all parameters for a run."""
        # Import MLflow entities locally to avoid circular imports
        from mlflow.entities import Param
        
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
    def set_experiment_tag(self, experiment_id: str, tag) -> None:
        """Set a tag on an experiment."""
        # Get current experiment to preserve other properties
        experiment = self.get_experiment(experiment_id)
        custom_props = {k: v for k, v in experiment.tags.items()}
        custom_props[tag.key] = tag.value
        
        payload = {"customProperties": custom_props}
        self._request("PATCH", f"/experiments/{experiment_id}", json=payload)
    
    def set_tag(self, run_id: str, tag) -> None:
        """Set a tag on a run."""
        # Get current run to preserve other properties
        run = self.get_run(run_id)
        custom_props = {k: v for k, v in run.data.tags.items()}
        custom_props[tag.key] = tag.value
        
        payload = {"customProperties": custom_props}
        self._request("PATCH", f"/experiment_runs/{run_id}", json=payload)
    
    def delete_tag(self, run_id: str, key: str) -> None:
        """Delete a tag from a run."""
        run = self.get_run(run_id)
        custom_props = {k: v for k, v in run.data.tags.items()}
        if key in custom_props:
            del custom_props[key]
        
        payload = {"customProperties": custom_props}
        self._request("PATCH", f"/experiment_runs/{run_id}", json=payload)
    
    # Search operations (simplified implementation)
    def search_runs(self, experiment_ids: List[str], filter_string: str = "",
                   run_view_type=None, max_results: int = 1000, 
                   order_by: List[str] = None, page_token: str = None):
        """Search for runs."""
        # Import MLflow entities locally to avoid circular imports
        from mlflow.entities import ViewType, LifecycleStage
        from mlflow.store.entities.paged_list import PagedList
        
        if run_view_type is None:
            run_view_type = ViewType.ACTIVE_ONLY
            
        all_runs = []

        # TODO add support for filter_string in ModelRegistry API
        for experiment_id in experiment_ids:
            response = self._request("GET", f"/experiments/{experiment_id}/experiment_runs", 
                                     params={"pageSize": str(min(max_results, 1000))})
            runs_data = response.json().get("items", [])
            
            for run_data in runs_data:
                run = self._getMLflowRun(run_data)
                # compare run.info.lifecycle_stage with run_view_type
                if run_view_type == ViewType.ACTIVE_ONLY and run.info.lifecycle_stage == LifecycleStage.DELETED:
                    continue
                elif run_view_type == ViewType.DELETED_ONLY and run.info.lifecycle_stage == LifecycleStage.ACTIVE:
                    continue
                all_runs.append(run)
                
        return PagedList(all_runs, response.json().get("nextPageToken"))

    def _getMLflowRun(self, run_data):
        # Import MLflow entities locally to avoid circular imports
        from mlflow.entities import Run, RunInfo, RunData, RunStatus, RunTag, LifecycleStage
        
        run_id = run_data["id"]
        metrics = self._get_run_metrics(run_id)
        params = self._get_run_params(run_id)
        tags = [RunTag(k, v) for k, v in run_data.get("customProperties", {}).items()]
        run_info = RunInfo(
            run_id=str(run_data["id"]),
            experiment_id=str(run_data["experimentId"]),
            user_id=run_data["owner"] or "unknown",
            status=RunStatus.from_string(run_data.get("status", "RUNNING")),
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

    def log_batch(self, run_id: str, metrics: Sequence = (), 
                 params: Sequence = (), tags: Sequence = ()) -> None:
        """Log a batch of metrics, parameters, and tags."""
        # Get current run to preserve other properties
        response = self._request("GET", f"/experiment_runs/{run_id}")
        run_data = response.json()
        custom_props = run_data["customProperties"] or {}
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

    def log_inputs(self, run_id: str, datasets: Optional[list] = None,
                   models: Optional[list] = None) -> None:
        """Log inputs for a run.
        
        Args:
            run_id: The ID of the run to log inputs for
            datasets: List of dataset inputs
            models: List of logged model inputs
        """
        # Import MLflow entities locally to avoid circular imports
        from mlflow.entities import DatasetInput, LoggedModelInput
        
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
                    "customProperties": {}
                }
                if datasetInput.tags:
                    for tag in datasetInput.tags:
                        payload["customProperties"][tag.key] = tag.value
                self._request("POST", f"/experiment_runs/{run_id}/artifacts", json=payload)

        if models:
            for model in models:
                # Get current model to preserve other properties
                response = self._request("GET", f"/artifacts/{model.model_id}")
                model_data = response.json()
                custom_props = model_data.get("customProperties", {})
                custom_props["mlflow.model_io_type"] = ModelIOType.INPUT.value
                
                payload = {
                    "artifactType": "model-artifact",
                    "id": model.model_id,
                    "customProperties": custom_props
                }
                self._request("POST", f"/experiment_runs/{run_id}/artifacts", json=payload)

    def log_outputs(self, run_id: str, models: list) -> None:
        """Log outputs for a run.
        
        Args:
            run_id: The ID of the run to log outputs for
            models: List of logged model outputs
        """
        # Import MLflow entities locally to avoid circular imports
        from mlflow.entities import LoggedModelOutput
        
        for model in models:
            # Get current model to preserve other properties
            response = self._request("GET", f"/artifacts/{model.model_id}")
            model_data = response.json()
            custom_props = model_data.get("customProperties", {})
            custom_props["mlflow.model_io_type"] = ModelIOType.OUTPUT.value
            
            payload = {
                "artifactType": "model-artifact",
                "id": model.model_id,
                "customProperties": custom_props
            }
            self._request("POST", f"/experiment_runs/{run_id}/artifacts", json=payload)

    def record_logged_model(self, run_id: str, mlflow_model) -> None:
        """Record a logged model.
        
        Args:
            run_id: The ID of the run to record the model for
            mlflow_model: The MLflow model to record
        """
        model_dict = mlflow_model.to_dict()
        model_info = mlflow_model.get_model_info()
        
        # Create a model artifact in Model Registry
        payload = {
            "artifactType": "model-artifact",
            "name": model_dict.get("model_uuid") or str(uuid.uuid4()), # generate a unique name if not provided
            "uri": model_info.model_uri,
            "customProperties": {
                "artifactPath": model_info.artifact_path,
                "model_uuid": model_info.model_uuid,
                "utc_time_created": model_info.utc_time_created,
                "mlflow_version": model_info.mlflow_version,
                "flavor": str(model_info.flavors),
                "source_run_id": run_id,
            }
        }
        
        # Store the full model dict as a tag for backward compatibility
        payload["customProperties"]["mlflow.logged_model"] = json.dumps(model_dict)
        
        # Create the model artifact
        self._request("POST", f"/experiment_runs/{run_id}/artifacts", json=payload)

    def create_logged_model(
        self,
        experiment_id: str,
        name: Optional[str] = None,
        source_run_id: Optional[str] = None,
        tags: Optional[list] = None,
        params: Optional[list] = None,
        model_type: Optional[str] = None,
    ):
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
        # Import MLflow entities locally to avoid circular imports
        from mlflow.entities import LoggedModel, LoggedModelTag, LoggedModelParameter, LoggedModelStatus
        
        payload = {
            "artifactType": "model-artifact",
            "name": name or str(uuid.uuid4()),
            "customProperties": {
                "model_type": model_type or "unknown",
                "experiment_id": experiment_id,
                "source_run_id": source_run_id,
            }
        }
        
        if tags:
            for tag in tags:
                payload["customProperties"][tag.key] = tag.value
            
        if params:
            for param in params:
                payload["customProperties"][f"param_{param.key}"] = param.value

        # TODO source_run_id is optional, but we need to handle it
        response = self._request("POST", f"/experiment_runs/{source_run_id}/artifacts", json=payload)
        model_data = response.json()
        
        return LoggedModel(
            model_id=str(model_data["id"]),
            experiment_id=experiment_id,
            name=model_data["name"],
            source_run_id=source_run_id,
            artifact_location=model_data["uri"],
            creation_timestamp=convert_timestamp(model_data["createTimeSinceEpoch"]),
            last_updated_timestamp=convert_timestamp(model_data["updateTimeSinceEpoch"]),
            model_type=model_type,
            status=LoggedModelStatus.READY,
            tags=tags or [],
            params=params or []
        )

    def search_logged_models(
        self,
        experiment_ids: list[str],
        filter_string: Optional[str] = None,
        datasets: Optional[list[dict[str, Any]]] = None,
        max_results: Optional[int] = None,
        order_by: Optional[list[dict[str, Any]]] = None,
        page_token: Optional[str] = None,
    ):
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
        # Import MLflow entities locally to avoid circular imports
        from mlflow.entities import LoggedModel, LoggedModelTag, LoggedModelParameter, LoggedModelStatus
        from mlflow.store.entities.paged_list import PagedList
        
        params = {
            "artifactType": "model-artifact",
            "experimentIds": experiment_ids
        }
        
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
            response = self._request("GET", f"/experiments/{experiment_id}/experiment_runs", params=params)
            runs_data = response.json().get("items", [])
            for run_data in runs_data:
                response = self._request("GET", f"/experiment_runs/{run_data['id']}/artifacts", 
                                         params={"artifactType": "model-artifact"})
                items = response.json().get("items", [])
                for item in items:
                    models.append(self.get_logged_model(item["id"]))
        
        return PagedList(models, None)

    def finalize_logged_model(self, model_id: str, status):
        """Finalize a model by updating its status.
        
        Args:
            model_id: ID of the model to finalize
            status: Final status to set on the model
            
        Returns:
            The updated LoggedModel
        """
        # Import MLflow entities locally to avoid circular imports
        from mlflow.entities import LoggedModel, LoggedModelTag, LoggedModelParameter, LoggedModelStatus
        
        payload = {
            "state": convert_to_model_artifact_state(status)
        }
        
        response = self._request("PATCH", f"/artifacts/{model_id}", json=payload)
        model_data = response.json()
        
        # Extract tags and params from customProperties
        custom_props = model_data.get("customProperties", {})
        tags = []
        params = []
        
        for key, value in custom_props.items():
            if key.startswith("param_"):
                params.append(LoggedModelParameter(key=key[6:], value=value))
            else:
                tags.append(LoggedModelTag(key=key, value=value))
            
        return LoggedModel(
            model_id=str(model_data["id"]),
            experiment_id=custom_props.get("experiment_id"),
            name=model_data["name"],
            source_run_id=custom_props.get("source_run_id"),
            artifact_location=model_data["uri"],
            creation_timestamp=convert_timestamp(model_data["createTimeSinceEpoch"]),
            last_updated_timestamp=convert_timestamp(model_data["updateTimeSinceEpoch"]),
            model_type=custom_props.get("model_type"),
            status=status,
            tags=tags,
            params=params
        )

    def set_logged_model_tags(self, model_id: str, tags: list) -> None:
        """Set tags on the specified logged model.
        
        Args:
            model_id: ID of the model
            tags: Tags to set on the model
        """
        # Get current model to preserve other properties
        response = self._request("GET", f"/artifacts/{model_id}")
        model_data = response.json()
        custom_props = model_data.get("customProperties", {})
        
        # Update custom properties with new tags
        for tag in tags:
            custom_props[tag.key] = tag.value
        
        payload = {
            "artifactType": "model-artifact",
            "customProperties": custom_props
        }
        self._request("PATCH", f"/artifacts/{model_id}", json=payload)

    def delete_logged_model_tag(self, model_id: str, key: str) -> None:
        """Delete a tag from the specified logged model.
        
        Args:
            model_id: ID of the model
            key: Key of the tag to delete
        """
        # Get current model to preserve other properties
        response = self._request("GET", f"/artifacts/{model_id}")
        model_data = response.json()
        custom_props = model_data.get("customProperties", {})
        
        # Remove the specified tag
        if key in custom_props:
            del custom_props[key]
        
        payload = {
            "artifactType": "model-artifact",
            "customProperties": custom_props
        }
        self._request("PATCH", f"/artifacts/{model_id}", json=payload)

    def get_logged_model(self, model_id: str):
        """Fetch the logged model with the specified ID.
        
        Args:
            model_id: ID of the model to fetch
            
        Returns:
            The fetched LoggedModel
        """
        # Import MLflow entities locally to avoid circular imports
        from mlflow.entities import LoggedModel, LoggedModelTag, LoggedModelParameter, LoggedModelStatus
        
        response = self._request("GET", f"/artifacts/{model_id}")
        model_data = response.json()
        
        # Extract tags and params from customProperties
        custom_props = model_data.get("customProperties", {})
        tags = []
        params = []
        
        for key, value in custom_props.items():
            if key.startswith("param_"):
                params.append(LoggedModelParameter(key=key[6:], value=value))
            else:
                tags.append(LoggedModelTag(key=key, value=value))
            
        return LoggedModel(
            model_id=str(model_data["id"]),
            experiment_id=custom_props.get("experiment_id"),
            name=model_data["name"],
            source_run_id=custom_props.get("source_run_id"),
            artifact_location=model_data["uri"],
            creation_timestamp=convert_timestamp(model_data["createTimeSinceEpoch"]),
            last_updated_timestamp=convert_timestamp(model_data["updateTimeSinceEpoch"]),
            model_type=custom_props.get("model_type"),
            status=convert_to_mlflow_logged_model_status(custom_props.get("state")),
            tags=tags,
            params=params
        )

    def delete_logged_model(self, model_id: str) -> None:
        """Delete the logged model with the specified ID.
        
        Args:
            model_id: ID of the model to delete
        """
        # Model Registry doesn't support deletion, so we mark as archived
        payload = {
            "artifactType": "model-artifact",
            "customProperties": {
                "state": "MARKED_FOR_DELETION" # TODO: handle ModelArtifactState.MARKED_FOR_DELETION in ModelRegistry
            }
        }
        self._request("PATCH", f"/artifacts/{model_id}", json=payload)
