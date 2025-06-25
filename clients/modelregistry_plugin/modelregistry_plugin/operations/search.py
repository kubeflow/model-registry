"""Search operations for Model Registry store."""

from __future__ import annotations

from typing import TYPE_CHECKING, List, Optional

if TYPE_CHECKING:
    from mlflow.entities import Run
    from mlflow.store.entities.paged_list import PagedList

from ..api_client import ModelRegistryAPIClient
from ..converters import MLflowEntityConverter


class SearchOperations:
    """Handles all search-related operations."""

    def __init__(self, api_client: ModelRegistryAPIClient, artifact_uri: str):
        """Initialize search operations.

        Args:
            api_client: API client for making requests
            artifact_uri: Default artifact URI
        """
        self.api_client = api_client
        self.artifact_uri = artifact_uri

    def search_runs(
        self,
        experiment_ids: Optional[List[str]] = None,
        filter_string: Optional[str] = None,
        run_view_type: Optional[int] = None,
        max_results: int = 1000,
        order_by: Optional[List[str]] = None,
        page_token: Optional[str] = None,
    ) -> PagedList[Run]:
        """Search for runs.

        Args:
            experiment_ids: List of experiment IDs to search
            filter_string: Filter string (not supported yet)
            run_view_type: Run view type
            max_results: Maximum number of results
            order_by: Order by fields (not supported yet)
            page_token: Token for pagination

        Returns:
            PagedList of runs
        """
        from mlflow.entities import LifecycleStage, ViewType
        from mlflow.store.entities.paged_list import PagedList

        if run_view_type is None:
            run_view_type = ViewType.ACTIVE_ONLY

        # TODO: Add support for filter_string and order_by in Model Registry API
        all_runs = []

        if experiment_ids:
            for experiment_id in experiment_ids:
                params = {}
                if max_results:
                    params["pageSize"] = str(min(max_results, 1000))
                if page_token:
                    params["pageToken"] = page_token

                response = self.api_client.get(
                    f"/experiments/{experiment_id}/experiment_runs", params=params
                )
                runs = response.get("items", [])

                for run_data in runs:
                    # Filter by view_type
                    lifecycle_stage = MLflowEntityConverter.to_mlflow_run_info(
                        run_data, self.artifact_uri
                    ).lifecycle_stage

                    if (
                        run_view_type == ViewType.ACTIVE_ONLY
                        and lifecycle_stage == LifecycleStage.DELETED
                    ) or (
                        run_view_type == ViewType.DELETED_ONLY
                        and lifecycle_stage == LifecycleStage.ACTIVE
                    ):
                        continue

                    # Get artifacts for the run
                    artifacts_response = self.api_client.get(
                        f"/experiment_runs/{run_data['id']}/artifacts"
                    )
                    artifacts = artifacts_response.get("items", [])

                    artifact_location = run_data.get("externalId") or self.artifact_uri
                    run = MLflowEntityConverter.to_mlflow_run(
                        run_data, artifacts, artifact_location
                    )
                    all_runs.append(run)

        return PagedList(all_runs, None)  # no paging across experiments
