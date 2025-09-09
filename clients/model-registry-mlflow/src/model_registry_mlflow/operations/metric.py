"""Metric operations for Model Registry store."""

from __future__ import annotations

from typing import TYPE_CHECKING, List, Optional

if TYPE_CHECKING:
    from mlflow.entities import Metric
    from mlflow.entities.metric import MetricWithRunId
    from mlflow.store.entities.paged_list import PagedList

from ..api_client import ModelRegistryAPIClient
from ..converters import MLflowEntityConverter


class MetricOperations:
    """Handles all metric-related operations."""

    def __init__(self, api_client: ModelRegistryAPIClient):
        """Initialize metric operations.

        Args:
            api_client: API client for making requests
        """
        self.api_client = api_client

    def get_metric_history(
        self,
        run_id: str,
        metric_key: str,
        max_results: Optional[int] = None,
        page_token: Optional[str] = None,
    ) -> PagedList[Metric]:
        """Get metric history for a run.

        Args:
            run_id: ID of the run
            metric_key: Key of the metric
            max_results: Maximum number of results
            page_token: Token for pagination

        Returns:
            PagedList of metrics
        """
        from mlflow.store.entities.paged_list import PagedList

        params = {"name": metric_key}
        if max_results:
            params["pageSize"] = max_results
        if page_token:
            params["pageToken"] = page_token

        response = self.api_client.get(
            f"/experiment_runs/{run_id}/metric_history", params=params
        )
        next_page_token = response.get("nextPageToken")
        items = response.get("items", [])

        metrics = []
        for item in items:
            metrics.append(MLflowEntityConverter.to_mlflow_metric(item))

        return PagedList(metrics, next_page_token if next_page_token != "" else None)

    def get_metric_history_bulk_interval_from_steps(
        self,
        run_id: str,
        metric_key: str,
        steps: Optional[List[int]] = None,
        max_results: Optional[int] = None,
        page_token: Optional[str] = None,
    ) -> PagedList[MetricWithRunId]:
        """Get metric history for a run with step filtering.

        Args:
            run_id: ID of the run
            metric_key: Key of the metric
            steps: List of steps to filter by
            max_results: Maximum number of results
            page_token: Token for pagination

        Returns:
            PagedList of metrics with run IDs
        """
        from mlflow.entities.metric import MetricWithRunId
        from mlflow.store.entities.paged_list import PagedList

        params = {"name": metric_key}
        if max_results:
            params["pageSize"] = max_results
        if page_token:
            params["pageToken"] = page_token
        if steps:
            params["stepIds"] = ",".join(map(str, steps))

        response = self.api_client.get(
            f"/experiment_runs/{run_id}/metric_history", params=params
        )
        next_page_token = response.get("nextPageToken")
        items = response.get("items", [])

        metrics = []
        for item in items:
            metric = MLflowEntityConverter.to_mlflow_metric(item)
            metrics.append(
                MetricWithRunId(
                    metric,
                    run_id,
                )
            )

        return PagedList(metrics, next_page_token if next_page_token != "" else None)

    def get_metric_history_bulk(
        self,
        run_ids: List[str],
        metric_key: str,
        max_results: Optional[int] = None,
    ) -> List[MetricWithRunId]:
        """Get metric history for multiple runs.

        Args:
            run_ids: List of run IDs
            metric_key: Key of the metric
            max_results: Maximum number of results per run

        Returns:
            List of metrics with run IDs
        """
        from mlflow.entities.metric import MetricWithRunId

        all_metrics = []

        for run_id in run_ids:
            params = {"name": metric_key}
            if max_results:
                params["pageSize"] = max_results

            response = self.api_client.get(
                f"/experiment_runs/{run_id}/metric_history", params=params
            )
            items = response.get("items", [])

            for item in items:
                metric = MLflowEntityConverter.to_mlflow_metric(item)
                all_metrics.append(
                    MetricWithRunId(
                        metric,
                        run_id,
                    )
                )

        return all_metrics
