"""Model Registry MLflow Tracking Store Implementation."""

from __future__ import annotations

import os
from typing import TYPE_CHECKING, List, Optional, Sequence, Any

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
        Param,
        Run,
        RunInfo,
        RunStatus,
        RunTag,
        ViewType,
    )
    from mlflow.store.entities.paged_list import PagedList
    from mlflow.models.model import Model

from .api_client import ModelRegistryAPIClient
from .operations.experiment import ExperimentOperations
from .operations.metric import MetricOperations
from .operations.model import ModelOperations
from .operations.run import RunOperations
from .operations.search import SearchOperations
from .utils import parse_tracking_uri


class ModelRegistryTrackingStore:
    """MLflow tracking store that uses Kubeflow Model Registry as the backend."""

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

        # Initialize API client
        self.api_client = ModelRegistryAPIClient(self.base_url)

        # Initialize operation classes
        self.experiments = ExperimentOperations(self.api_client, self.artifact_uri)
        self.runs = RunOperations(self.api_client, self.artifact_uri)
        self.metrics = MetricOperations(self.api_client)
        self.models = ModelOperations(self.api_client, self.artifact_uri)
        self.search = SearchOperations(self.api_client, self.artifact_uri)

    # Experiment operations
    def create_experiment(
        self,
        name: str,
        artifact_location: Optional[str] = None,
        tags: Optional[List[ExperimentTag]] = None,
    ) -> str:
        """
        Create a new experiment.
        If an experiment with the given name already exists, throws exception.

        Args:
            name: Desired name for an experiment.
            artifact_location: Base location for artifacts in runs. May be None.
            tags: Experiment tags to set upon experiment creation

        Returns:
            experiment_id (string) for the newly created experiment if successful, else None.

        """
        return self.experiments.create_experiment(name, artifact_location, tags)

    def get_experiment(self, experiment_id: str) -> Experiment:
        """
        Fetch the experiment by ID from the backend store.

        Args:
            experiment_id: String id for the experiment

        Returns:
            A single :py:class:`mlflow.entities.Experiment` object if it exists,
            otherwise raises an exception.
        """
        return self.experiments.get_experiment(experiment_id)

    def get_experiment_by_name(self, experiment_name: str) -> Optional[Experiment]:
        """
        Fetch the experiment by name from the backend store.

        Args:
            experiment_name: Name of experiment

        Returns:
            A single :py:class:`mlflow.entities.Experiment` object if it exists.
        """
        return self.experiments.get_experiment_by_name(experiment_name)

    def delete_experiment(self, experiment_id: str) -> None:
        """
        Delete the experiment from the backend store. Deleted experiments can be restored until
        permanently deleted.

        Args:
            experiment_id: String id for the experiment.
        """
        self.experiments.delete_experiment(experiment_id)

    def restore_experiment(self, experiment_id: str) -> None:
        """
        Restore deleted experiment unless it is permanently deleted.

        Args:
            experiment_id: String id for the experiment.
        """
        self.experiments.restore_experiment(experiment_id)

    def rename_experiment(self, experiment_id: str, new_name: str) -> None:
        """
        Update an experiment's name. The new name must be unique.

        Args:
            experiment_id: String id for the experiment.
            new_name: New name for the experiment.
        """
        self.experiments.rename_experiment(experiment_id, new_name)

    def search_experiments(
        self,
        view_type: Optional[ViewType] = None,
        max_results: int = 1000,
        filter_string: Optional[str] = None,
        order_by: Optional[List[str]] = None,
        page_token: Optional[str] = None,
    ) -> PagedList[Experiment]:
        """
        Search for experiments that match the specified search query.

        Args:
            view_type: One of enum values ``ACTIVE_ONLY``, ``DELETED_ONLY``, or ``ALL``
                defined in :py:class:`mlflow.entities.ViewType`.
            max_results: Maximum number of experiments desired. Certain server backend may apply
                its own limit.
            filter_string: Filter query string (e.g., ``"name = 'my_experiment'"``), defaults to
                searching for all experiments. The following identifiers, comparators, and logical
                operators are supported.

                Identifiers
                  - ``name``: Experiment name
                  - ``creation_time``: Experiment creation time
                  - ``last_update_time``: Experiment last update time
                  - ``tags.<tag_key>``: Experiment tag. If ``tag_key`` contains
                    spaces, it must be wrapped with backticks (e.g., ``"tags.`extra key`"``).

                Comparators for string attributes and tags
                  - ``=``: Equal to
                  - ``!=``: Not equal to
                  - ``LIKE``: Case-sensitive pattern match
                  - ``ILIKE``: Case-insensitive pattern match

                Comparators for numeric attributes
                  - ``=``: Equal to
                  - ``!=``: Not equal to
                  - ``<``: Less than
                  - ``<=``: Less than or equal to
                  - ``>``: Greater than
                  - ``>=``: Greater than or equal to

                Logical operators
                  - ``AND``: Combines two sub-queries and returns True if both of them are True.

            order_by: List of columns to order by. The ``order_by`` column can contain an optional
                ``DESC`` or ``ASC`` value (e.g., ``"name DESC"``). The default ordering is ``ASC``,
                so ``"name"`` is equivalent to ``"name ASC"``. If unspecified, defaults to
                ``["last_update_time DESC"]``, which lists experiments updated most recently first.
                The following fields are supported:

                - ``experiment_id``: Experiment ID
                - ``name``: Experiment name
                - ``creation_time``: Experiment creation time
                - ``last_update_time``: Experiment last update time

            page_token: Token specifying the next page of results. It should be obtained from
                a ``search_experiments`` call.

        Returns:
            A :py:class:`PagedList <mlflow.store.entities.PagedList>` of
            :py:class:`Experiment <mlflow.entities.Experiment>` objects. The pagination token
            for the next page can be obtained via the ``token`` attribute of the object.

        """
        return self.experiments.search_experiments(
            view_type, max_results, filter_string, order_by, page_token
        )

    def set_experiment_tag(self, experiment_id: str, tag: ExperimentTag) -> None:
        """
        Set a tag for the specified experiment

        Args:
            experiment_id: String id for the experiment.
            tag: :py:class:`mlflow.entities.ExperimentTag` instance to set.
        """
        self.experiments.set_experiment_tag(experiment_id, tag)

    # Run operations
    def create_run(
        self,
        experiment_id: str,
        user_id: Optional[str] = None,
        start_time: Optional[int] = None,
        tags: Optional[List[RunTag]] = None,
        run_name: Optional[str] = None,
    ) -> Run:
        """
        Create a run under the specified experiment ID, setting the run's status to "RUNNING"
        and the start time to the current time.

        Args:
            experiment_id: String id of the experiment for this run.
            user_id: ID of the user launching this run.
            start_time: Start time of the run.
            tags: A dictionary of string keys and string values.
            run_name: Name of the run.

        Returns:
            The created Run object
        """
        return self.runs.create_run(experiment_id, user_id, start_time, tags, run_name)

    def get_run(self, run_id: str) -> Run:
        """
        Fetch the run from backend store. The resulting :py:class:`Run <mlflow.entities.Run>`
        contains a collection of run metadata - :py:class:`RunInfo <mlflow.entities.RunInfo>`,
        as well as a collection of run parameters, tags, and metrics -
        :py:class:`RunData <mlflow.entities.RunData>`. In the case where multiple metrics with the
        same key are logged for the run, the :py:class:`RunData <mlflow.entities.RunData>` contains
        the value at the latest timestamp for each metric. If there are multiple values with the
        latest timestamp for a given metric, the maximum of these values is returned.

        Args:
            run_id: Unique identifier for the run.

        Returns:
            A single :py:class:`mlflow.entities.Run` object, if the run exists. Otherwise,
            raises an exception.
        """
        return self.runs.get_run(run_id)

    def update_run_info(
        self,
        run_id: str,
        run_status: Optional[RunStatus] = None,
        end_time: Optional[int] = None,
        run_name: Optional[str] = None,
    ) -> RunInfo:
        """
        Update the metadata of the specified run.

        Returns:
            mlflow.entities.RunInfo: Describing the updated run.
        """
        return self.runs.update_run_info(run_id, run_status, end_time, run_name)

    def delete_run(self, run_id: str) -> None:
        """
        Delete a run.

        Args:
            run_id: The ID of the run to delete.

        """
        self.runs.delete_run(run_id)

    def restore_run(self, run_id: str) -> None:
        """
        Restore a run.

        Args:
            run_id: The ID of the run to restore.

        """
        self.runs.restore_run(run_id)

    def log_metric(self, run_id: str, metric: Metric) -> None:
        """
        Log a metric for the specified run

        Args:
            run_id: String id for the run
            metric: `mlflow.entities.Metric` instance to log
        """
        self.runs.log_metric(run_id, metric)

    def log_param(self, run_id: str, param: Param) -> None:
        """
        Log a param for the specified run

        Args:
            run_id: String id for the run
            param: :py:class:`mlflow.entities.Param` instance to log
        """
        self.runs.log_param(run_id, param)

    def log_batch(
        self,
        run_id: str,
        metrics: Sequence[Metric] = (),
        params: Sequence[Param] = (),
        tags: Sequence[RunTag] = (),
    ) -> None:
        """
        Log multiple metrics, params, and tags for the specified run

        Args:
            run_id: String id for the run
            metrics: List of :py:class:`mlflow.entities.Metric` instances to log
            params: List of :py:class:`mlflow.entities.Param` instances to log
            tags: List of :py:class:`mlflow.entities.RunTag` instances to log

        Returns:
            None.
        """
        self.runs.log_batch(run_id, metrics, params, tags)

    # Async logging methods
    def log_batch_async(
        self,
        run_id: str,
        metrics: List[Metric],
        params: List[Param],
        tags: List[RunTag],
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

    def log_inputs(
        self,
        run_id: str,
        datasets: Optional[List[DatasetInput]] = None,
        models: Optional[List[LoggedModelInput]] = None,
    ) -> None:
        """
        Log inputs, such as datasets, to the specified run.

        Args:
            run_id: String id for the run
            datasets: List of :py:class:`mlflow.entities.DatasetInput` instances to log
                as inputs to the run.
            models: List of :py:class:`mlflow.entities.LoggedModelInput` instances to log
                as inputs to the run.

        Returns:
            None.
        """
        self.runs.log_inputs(run_id, datasets, models)

    def log_outputs(self, run_id: str, models: List[LoggedModelOutput]) -> None:
        """
        Log outputs, such as models, to the specified run.

        Args:
            run_id: String id for the run
            models: List of :py:class:`mlflow.entities.LoggedModelOutput` instances to log
                as outputs of the run.

        Returns:
            None.
        """
        self.runs.log_outputs(run_id, models)

    def set_tag(self, run_id: str, tag: RunTag) -> None:
        """
        Set a tag for the specified run

        Args:
            run_id: String id for the run.
            tag: :py:class:`mlflow.entities.RunTag` instance to set.
        """
        self.runs.set_tag(run_id, tag)

    def delete_tag(self, run_id: str, key: str) -> None:
        """
        Delete a tag from the specified run

        Args:
            run_id: String id for the run.
            key: Key of the tag to delete.
        """
        self.runs.delete_tag(run_id, key)

    # Metric operations
    def get_metric_history(
        self,
        run_id: str,
        metric_key: str,
        max_results: Optional[int] = None,
        page_token: Optional[str] = None,
    ) -> List[Metric]:
        """
        Return a list of metric objects corresponding to all values logged for a given metric
        within a run.

        Args:
            run_id: Unique identifier for run.
            metric_key: Metric name within the run.
            max_results: Maximum number of metric history events (steps) to return per paged
                query.
            page_token: A Token specifying the next paginated set of results of metric history.
                This value is obtained as a return value from a paginated call to GetMetricHistory.

        Returns:
            A list of :py:class:`mlflow.entities.Metric` entities if logged, else empty list.
        """
        return self.metrics.get_metric_history(
            run_id, metric_key, max_results, page_token
        )

    def search_runs(
        self,
        experiment_ids: Optional[List[str]] = None,
        filter_string: Optional[str] = None,
        run_view_type: Optional[int] = None,
        max_results: int = 1000,
        order_by: Optional[List[str]] = None,
        page_token: Optional[str] = None,
    ) -> PagedList[Run]:
        """
        Return runs that match the given list of search expressions within the experiments.

        Args:
            experiment_ids: List of experiment ids to scope the search.
            filter_string: A search filter string.
            run_view_type: ACTIVE_ONLY, DELETED_ONLY, or ALL runs.
            max_results: Maximum number of runs desired.
            order_by: List of order_by clauses.
            page_token: Token specifying the next page of results. It should be obtained from
                a ``search_runs`` call.

        Returns:
            A :py:class:`PagedList <mlflow.store.entities.PagedList>` of
            :py:class:`Run <mlflow.entities.Run>` objects that satisfy the search expressions.
            If the underlying tracking store supports pagination, the token for the next page may
            be obtained via the ``token`` attribute of the returned object; however, some store
            implementations may not support pagination and thus the returned token would not be
            meaningful in such cases.
        """
        return self.search.search_runs(
            experiment_ids,
            filter_string,
            run_view_type,
            max_results,
            order_by,
            page_token,
        )

    def get_metric_history_bulk_interval_from_steps(
        self,
        run_id: str,
        metric_key: str,
        steps: list[int],
        max_results: int | None = None,
    ):
        """
        Return a list of metric objects corresponding to all values logged
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
        return self.metrics.get_metric_history_bulk_interval_from_steps(
            run_id, metric_key, steps, max_results
        )

    def get_metric_history_bulk(
        self,
        run_ids: list[str],
        metric_key: str,
        max_results: int | None = None,
    ):
        """
        Return a list of metric objects corresponding to all values logged
        for a given metric within multiple runs.

        Args:
            run_ids: List of unique identifiers for runs.
            metric_key: Metric name within the runs.
            max_results: Maximum number of metric history events (steps) to return per run.

        Returns:
            A list of MetricWithRunId objects:
                - key: Metric name within the run.
                - value: Metric value.
                - timestamp: Metric timestamp.
                - step: Metric step.
                - run_id: Unique identifier for run.
        """
        return self.metrics.get_metric_history_bulk(run_ids, metric_key, max_results)

    # Model operations
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
        """
        Create a new logged model.

        Args:
            experiment_id: ID of the experiment to which the model belongs.
            name: Name of the model. If not specified, a random name will be generated.
            source_run_id: ID of the run that produced the model.
            tags: Tags to set on the model.
            params: Parameters to set on the model.
            model_type: Type of the model.

        Returns:
            The created model.
        """
        return self.models.create_logged_model(
            name,
            source_run_id,
            experiment_id,
            model_type,
            artifact_location,
            tags,
            params,
        )

    def get_logged_model(self, model_id: str) -> LoggedModel:
        """
        Fetch the logged model with the specified ID.

        Args:
            model_id: ID of the model to fetch.

        Returns:
            The fetched model.
        """
        return self.models.get_logged_model(model_id)

    def delete_logged_model(self, model_id: str) -> None:
        """
        Delete the logged model with the specified ID.

        Args:
            model_id: ID of the model to delete.
        """
        self.models.delete_logged_model(model_id)

    def delete_logged_model_tag(self, model_id: str, key: str) -> None:
        """
        Delete a tag from the specified logged model.

        Args:
            model_id: ID of the model.
            key: Key of the tag to delete.
        """
        self.models.delete_logged_model_tag(model_id, key)

    def search_logged_models(
        self,
        experiment_ids: List[str],
        filter_string: Optional[str] = None,
        datasets: Optional[List[dict[str, Any]]] = None,
        max_results: Optional[int] = None,
        order_by: Optional[List[dict[str, Any]]] = None,
        page_token: Optional[str] = None,
    ) -> PagedList[LoggedModel]:
        """
        Search for logged models that match the specified search criteria.

        Args:
            experiment_ids: List of experiment ids to scope the search.
            filter_string: A search filter string.
            datasets: List of dictionaries to specify datasets on which to apply metrics filters.
                The following fields are supported:

                name (str): Required. Name of the dataset.
                digest (str): Optional. Digest of the dataset.
            max_results: Maximum number of logged models desired.
            order_by: List of dictionaries to specify the ordering of the search results.
                The following fields are supported:

                field_name (str): Required. Name of the field to order by, e.g. "metrics.accuracy".
                ascending: (bool): Optional. Whether the order is ascending or not.
                dataset_name: (str): Optional. If ``field_name`` refers to a metric, this field
                    specifies the name of the dataset associated with the metric. Only metrics
                    associated with the specified dataset name will be considered for ordering.
                    This field may only be set if ``field_name`` refers to a metric.
                dataset_digest (str): Optional. If ``field_name`` refers to a metric, this field
                    specifies the digest of the dataset associated with the metric. Only metrics
                    associated with the specified dataset name and digest will be considered for
                    ordering. This field may only be set if ``dataset_name`` is also set.
            page_token: Token specifying the next page of results.

        Returns:
            A :py:class:`PagedList <mlflow.store.entities.PagedList>` of
            :py:class:`LoggedModel <mlflow.entities.LoggedModel>` objects.
        """
        return self.models.search_logged_models(
            experiment_ids, filter_string, datasets, max_results, order_by, page_token
        )

    def record_logged_model(self, run_id: str, mlflow_model: Model) -> None:
        """Record a logged model for a run."""
        return self.models.record_logged_model(run_id, mlflow_model)

    def finalize_logged_model(self, model_id: str, status: Any) -> LoggedModel:
        """
        Finalize a model by updating its status.

        Args:
            model_id: ID of the model to finalize.
            status: Final status to set on the model.

        Returns:
            The updated model.
        """
        return self.models.finalize_logged_model(model_id, status)

    def set_logged_model_tags(self, model_id: str, tags):
        """
        Set tags on the specified logged model.

        Args:
            model_id: ID of the model.
            tags: Tags to set on the model.

        Returns:
            None
        """
        return self.models.set_logged_model_tags(model_id, tags)

    def log_logged_model_params(self, model_id: str, params):
        """
        Log parameters for the specified logged model.

        Args:
            model_id: ID of the model.
            params: Parameters to log on the model.

        Returns:
            None
        """
        return self.models.log_logged_model_params(model_id, params)

    def _search_datasets(self, experiment_ids):
        """
        Search for datasets across experiments.

        Args:
            experiment_ids: List of experiment IDs to search.

        Returns:
            List of dataset summaries.
        """
        return self.search._search_datasets(experiment_ids)
