"""Experiment operations for Model Registry store."""

from __future__ import annotations

from typing import TYPE_CHECKING, List, Optional

if TYPE_CHECKING:
    from mlflow.entities import Experiment, ExperimentTag, ViewType
    from mlflow.store.entities.paged_list import PagedList

from ..api_client import ModelRegistryAPIClient
from ..converters import MLflowEntityConverter


class ExperimentOperations:
    """Handles all experiment-related operations."""

    def __init__(self, api_client: ModelRegistryAPIClient, artifact_uri: str):
        """Initialize experiment operations.

        Args:
            api_client: API client for making requests
            artifact_uri: Default artifact URI
        """
        self.api_client = api_client
        self.artifact_uri = artifact_uri

    def create_experiment(
        self,
        name: str,
        artifact_location: Optional[str] = None,
        tags: Optional[List[ExperimentTag]] = None,
    ) -> str:
        """Create a new experiment in Model Registry.

        Args:
            name: Name of the experiment
            artifact_location: Artifact location for the experiment
            tags: Tags to set on the experiment

        Returns:
            Experiment ID
        """
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
            payload["externalId"] = None

        if tags:
            for tag in tags:
                payload["customProperties"][tag.key] = tag.value

        experiment_data = self.api_client.post("/experiments", json=payload)
        experiment_id = str(experiment_data["id"])

        # If no artifact_location was provided but we have artifact_uri, update with the default pattern
        if not artifact_location and self.artifact_uri:
            default_artifact_location = (
                f"{self.artifact_uri}/experiments/{experiment_id}"
            )
            update_payload = {"externalId": default_artifact_location}
            self.api_client.patch(f"/experiments/{experiment_id}", json=update_payload)

        return experiment_id

    def get_experiment(self, experiment_id: str) -> Experiment:
        """Get experiment by ID.

        Args:
            experiment_id: ID of the experiment

        Returns:
            Experiment entity
        """
        experiment_data = self.api_client.get(f"/experiments/{experiment_id}")
        return MLflowEntityConverter.to_mlflow_experiment(
            experiment_data, self.artifact_uri
        )

    def get_experiment_by_name(self, experiment_name: str) -> Optional[Experiment]:
        """Get experiment by name.

        Args:
            experiment_name: Name of the experiment

        Returns:
            Experiment entity or None if not found
        """
        from mlflow.exceptions import MlflowException

        try:
            exp_data = self.api_client.get(
                "/experiment", params={"name": experiment_name}
            )
            return MLflowEntityConverter.to_mlflow_experiment(
                exp_data, self.artifact_uri
            )
        except MlflowException as e:
            if e.get_http_status_code() == 404 and "not found" in e.message:
                return None
            raise e

    def delete_experiment(self, experiment_id: str) -> None:
        """Delete an experiment.

        Args:
            experiment_id: ID of the experiment to delete
        """
        payload = {"state": "ARCHIVED"}
        self.api_client.patch(f"/experiments/{experiment_id}", json=payload)

    def restore_experiment(self, experiment_id: str) -> None:
        """Restore a deleted experiment.

        Args:
            experiment_id: ID of the experiment to restore
        """
        payload = {"state": "LIVE"}
        self.api_client.patch(f"/experiments/{experiment_id}", json=payload)

    def rename_experiment(self, experiment_id: str, new_name: str) -> None:
        """Rename an experiment.

        Args:
            experiment_id: ID of the experiment
            new_name: New name for the experiment
        """
        payload = {"name": new_name}
        self.api_client.patch(f"/experiments/{experiment_id}", json=payload)

    def search_experiments(
        self,
        view_type: Optional[ViewType] = None,
        max_results: int = 1000,
        filter_string: Optional[str] = None,
        order_by: Optional[List[str]] = None,
        page_token: Optional[str] = None,
    ) -> PagedList[Experiment]:
        """Search for experiments.

        Args:
            view_type: Type of experiments to search
            max_results: Maximum number of results
            filter_string: Filter string (not supported yet)
            order_by: Order by fields (not supported yet)
            page_token: Token for pagination

        Returns:
            PagedList of experiments
        """
        from mlflow.entities import LifecycleStage, ViewType
        from mlflow.store.entities.paged_list import PagedList

        if view_type is None:
            view_type = ViewType.ACTIVE_ONLY

        # TODO: Add support for filter_string and order_by in Model Registry API
        params = {}
        if max_results:
            params["pageSize"] = max_results
        if page_token:
            params["pageToken"] = page_token

        response_data = self.api_client.get("/experiments", params=params)
        items = response_data.get("items", [])

        experiments = []
        for exp_data in items:
            lifecycle_stage = MLflowEntityConverter.to_mlflow_experiment(
                exp_data, self.artifact_uri
            ).lifecycle_stage

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
                MLflowEntityConverter.to_mlflow_experiment(exp_data, self.artifact_uri)
            )

        nextPageToken = response_data.get("nextPageToken")
        return PagedList(experiments, nextPageToken if nextPageToken != "" else None)

    def set_experiment_tag(self, experiment_id: str, tag: ExperimentTag) -> None:
        """Set a tag on an experiment.

        Args:
            experiment_id: ID of the experiment
            tag: Tag to set
        """
        experiment = self.api_client.get(f"/experiments/{experiment_id}")
        custom_props = experiment.get("customProperties", {})
        custom_props[tag.key] = tag.value

        payload = {"customProperties": custom_props}
        self.api_client.patch(f"/experiments/{experiment_id}", json=payload)
