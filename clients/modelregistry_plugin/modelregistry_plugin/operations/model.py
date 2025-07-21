"""Model operations for Model Registry store."""

from __future__ import annotations

import json
from typing import TYPE_CHECKING, List, Optional

if TYPE_CHECKING:
    from mlflow.entities import (
        LoggedModel,
        LoggedModelParameter,
        LoggedModelTag,
        LoggedModelStatus,
    )
    from mlflow.store.entities.paged_list import PagedList
    from mlflow.models.model import Model

from ..api_client import ModelRegistryAPIClient
from ..converters import MLflowEntityConverter


class ModelOperations:
    """Handles all model-related operations."""

    def __init__(self, api_client: ModelRegistryAPIClient, artifact_uri: str):
        """Initialize model operations.

        Args:
            api_client: API client for making requests
            artifact_uri: Base artifact URI for the store
        """
        self.api_client = api_client
        self.artifact_uri = artifact_uri

    def record_logged_model(self, run_id: str, model: Model) -> LoggedModel:
        """Record a logged model from a run.

        Args:
            run_id: ID of the run
            model: MLflow Model to record

        Returns:
            LoggedModel entity
        """
        model_info = model.get_model_info()

        # Get experiment ID from the run
        run_data = self.api_client.get(f"/experiment_runs/{run_id}")
        experiment_id = run_data.get("experimentId")

        # Convert model to dict for serialization
        model_dict = model.to_dict()

        payload = {
            "artifactType": "model-artifact",
            "name": model_dict.get("name", model.model_uuid or model.name),
            "description": model_dict.get(
                "description", f"MLflow logged model: {model.name}"
            ),
            "uri": model_dict.get("artifact_location", model_info.model_uri),
            "customProperties": {
                # Standard MLflow custom properties - try to get from model_dict first
                "mlflow__name": model_dict.get("name", model.name),
                "mlflow__version": model_dict.get("version", "1"),
                "mlflow__user_id": model_dict.get("user_id", "unknown"),
                "mlflow__status_message": model_dict.get("status_message", ""),
                "mlflow__artifact_location": model_dict.get(
                    "artifact_location", model_info.model_uri
                ),
                "mlflow__description": model_dict.get(
                    "description", f"MLflow logged model: {model.name}"
                ),
                "mlflow__experiment_id": model_dict.get("experiment_id", experiment_id),
                "mlflow__source_run_id": model_dict.get("source_run_id", run_id),
                "mlflow__model_type": model_dict.get("model_type", "unknown"),
                # MLflow model-specific properties
                "mlflow__artifactPath": model_info.artifact_path,
                "mlflow__model_uuid": model.model_uuid,
                "mlflow__utc_time_created": model_info.utc_time_created,
                "mlflow__mlflow_version": model_info.mlflow_version,
                "mlflow__model_io_type": "output",
                # Store the full model dict for backward compatibility
                "mlflow__logged_model": json.dumps(model_dict),
            },
        }

        # Add flavors as custom properties
        for flavor_name, flavor_config in model_info.flavors.items():
            payload["customProperties"][f"mlflow__flavor_{flavor_name}"] = json.dumps(
                flavor_config
            )

        model_data = self.api_client.post(
            f"/experiment_runs/{run_id}/artifacts", json=payload
        )
        return MLflowEntityConverter.to_mlflow_logged_model(model_data)

    def create_logged_model(
        self,
        name: str,
        source_run_id: Optional[str] = None,
        experiment_id: Optional[str] = None,
        model_type: Optional[str] = None,
        artifact_location: Optional[str] = None,
        tags: Optional[List[LoggedModelTag]] = None,
        params: Optional[List[LoggedModelParameter]] = None,
    ) -> LoggedModel:
        """Create a new logged model.

        Args:
            name: Name of the model
            source_run_id: ID of the run that produced the model
            experiment_id: ID of the experiment to which the model belongs
            model_type: Type of the model
            artifact_location: Artifact location for the model
            tags: Tags to set on the model
            params: Parameters to set on the model

        Returns:
            Created LoggedModel entity
        """
        # Check if experiment exists and get artifact location
        if experiment_id:
            experiment_data = self.api_client.get(f"/experiments/{experiment_id}")
            if not artifact_location:
                artifact_location = (
                    experiment_data.get("externalId") or self.artifact_uri
                )

        # Set artifact location
        if artifact_location:
            if source_run_id:
                final_artifact_location = f"{artifact_location}/{source_run_id}/{name}"
            else:
                final_artifact_location = f"{artifact_location}/{name}"
        else:
            final_artifact_location = ""

        payload = {
            "artifactType": "model-artifact",
            "name": name,
            "description": f"MLflow logged model: {name}",
            "customProperties": {
                "mlflow__model_type": model_type or "unknown",
                "mlflow__experiment_id": experiment_id,
                "mlflow__source_run_id": source_run_id,
                "mlflow__name": name,
                "mlflow__version": "1",
                "mlflow__user_id": "unknown",  # TODO: Get from context
                "mlflow__status_message": "",
                "mlflow__artifact_location": final_artifact_location,
                "mlflow__description": f"MLflow logged model: {name}",
            },
        }

        # Set artifact location
        if artifact_location:
            if source_run_id:
                payload["uri"] = f"{artifact_location}/{source_run_id}/{name}"
            else:
                payload["uri"] = f"{artifact_location}/{name}"

        if tags:
            for tag in tags:
                payload["customProperties"][tag.key] = tag.value

        if params:
            for param in params:
                payload["customProperties"][f"param_{param.key}"] = param.value

        model_data = self.api_client.post(
            f"/experiment_runs/{source_run_id}/artifacts", json=payload
        )
        return MLflowEntityConverter.to_mlflow_logged_model(model_data)

    def get_logged_model(self, model_id: str) -> LoggedModel:
        """Get a logged model by ID.

        Args:
            model_id: ID of the model

        Returns:
            LoggedModel entity
        """
        model_data = self.api_client.get(f"/artifacts/{model_id}")
        return MLflowEntityConverter.to_mlflow_logged_model(model_data)

    def delete_logged_model(self, model_id: str) -> None:
        """Delete a logged model.

        Args:
            model_id: ID of the model to delete
        """
        # Get current model to preserve other properties
        model_data = self.api_client.get(f"/artifacts/{model_id}")
        custom_props = model_data.get("customProperties", {})
        custom_props["state"] = "MARKED_FOR_DELETION"

        payload = {"artifactType": "model-artifact", "customProperties": custom_props}
        self.api_client.patch(f"/artifacts/{model_id}", json=payload)

    def delete_logged_model_tag(self, model_id: str, key: str) -> None:
        """Delete a tag from a logged model.

        Args:
            model_id: ID of the model
            key: Key of the tag to delete
        """
        model = self.api_client.get(f"/artifacts/{model_id}")
        custom_props = model.get("customProperties", {})
        if key in custom_props:
            del custom_props[key]

        payload = {"customProperties": custom_props}
        self.api_client.patch(f"/artifacts/{model_id}", json=payload)

    def search_logged_models(
        self,
        experiment_ids: List[str],
        filter_string: Optional[str] = None,
        datasets: Optional[List[dict]] = None,
        max_results: Optional[int] = None,
        order_by: Optional[List[dict]] = None,
        page_token: Optional[str] = None,
    ) -> PagedList[LoggedModel]:
        """Search for logged models across experiments.

        Args:
            experiment_ids: List of experiment IDs to search
            filter_string: Filter string (not supported yet)
            datasets: List of datasets for metrics filtering (not supported yet)
            max_results: Maximum number of results
            order_by: List of order specifications (not supported yet)
            page_token: Token for pagination

        Returns:
            PagedList of LoggedModel entities
        """
        from mlflow.store.entities.paged_list import PagedList

        # TODO: Add support for filter_string in ModelRegistry API
        # TODO: Add support for datasets filtering in ModelRegistry API
        # TODO: Add support for order_by mapping to ModelRegistry API
        # TODO: Add support for pagination in ModelRegistry API across list of experiments

        all_models = []

        # Iterate over experiment_ids and get all runs, and get all model-artifacts for each run
        for experiment_id in experiment_ids:
            # Get runs from experiment
            runs_response = self.api_client.get(
                f"/experiments/{experiment_id}/experiment_runs"
            )
            runs = runs_response.get("items", [])

            for run in runs:
                run_id = run["id"]
                # Get artifacts from run
                artifacts_response = self.api_client.get(
                    f"/experiment_runs/{run_id}/artifacts",
                    params={"artifactType": "model-artifact"},
                )
                artifacts = artifacts_response.get("items", [])

                for artifact in artifacts:
                    all_models.append(
                        MLflowEntityConverter.to_mlflow_logged_model(artifact)
                    )

        # Apply max_results limit if specified
        if max_results and len(all_models) > max_results:
            all_models = all_models[:max_results]

        # Return PagedList with no paging across experiments
        return PagedList(items=all_models, token=None)

    def finalize_logged_model(
        self, model_id: str, status: LoggedModelStatus
    ) -> LoggedModel:
        """Finalize a logged model with status.

        Args:
            model_id: ID of the model
            status: Status to set

        Returns:
            Updated LoggedModel entity
        """
        from ..utils import convert_to_model_artifact_state

        # Convert MLflow status to Model Registry state
        model_state = convert_to_model_artifact_state(status)

        # Use /artifacts endpoint with state and artifactType discriminator properties
        payload = {"artifactType": "model-artifact", "state": model_state}

        model_data = self.api_client.patch(f"/artifacts/{model_id}", json=payload)
        return MLflowEntityConverter.to_mlflow_logged_model(model_data)

    def set_logged_model_tags(self, model_id: str, tags: List[LoggedModelTag]) -> None:
        """Set multiple tags on a logged model.

        Args:
            model_id: ID of the model
            tags: List of tags to set
        """
        model = self.api_client.get(f"/artifacts/{model_id}")
        custom_props = model.get("customProperties", {})

        for tag in tags:
            custom_props[tag.key] = tag.value

        self.api_client.patch(
            f"/artifacts/{model_id}", json={"customProperties": custom_props}
        )

    def log_logged_model_params(
        self, model_id: str, params: List[LoggedModelParameter]
    ) -> None:
        """Log parameters for a logged model.

        Args:
            model_id: ID of the model
            params: List of parameters to log
        """
        model = self.api_client.get(f"/artifacts/{model_id}")
        custom_props = model.get("customProperties", {})

        for param in params:
            custom_props[f"param_{param.key}"] = param.value

        self.api_client.patch(
            f"/artifacts/{model_id}", json={"customProperties": custom_props}
        )
