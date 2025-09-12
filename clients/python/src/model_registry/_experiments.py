from __future__ import annotations  # noqa: I001

import time
from contextlib import AbstractContextManager
from dataclasses import dataclass
from typing import Any, Callable, Literal

from model_registry.core import ModelRegistryAPIClient
from model_registry.exceptions import StoreError
from model_registry.types.artifacts import (
    ArtifactState,
    DataSet,
    ExperimentRunArtifact,
    ExperimentRunArtifactTypes,
    Metric,
    Parameter,
    ParameterType,
)
from model_registry.types.experiments import ExperimentRun

from .utils import S3Params, ThreadSafeVariable, upload_to_s3

LogType = Literal["params", "metrics", "datasets"]


@dataclass
class RunContext:
    id: str | None = None
    name: str | None = None
    run_id: str | None = None
    active: bool = False


@dataclass
class RunInfo:
    experiment_id: str
    id: str
    name: str


class ActiveExperimentRun(AbstractContextManager):
    def __init__(
        self,
        thread_safe_ctx: ThreadSafeVariable,
        experiment_run: ExperimentRun,
        api: ModelRegistryAPIClient,
        async_runner: Callable,
    ):
        self._thread_safe_ctx = thread_safe_ctx
        self._exp_run = experiment_run
        self.info = RunInfo(
            id=experiment_run.id,
            name=experiment_run.name,
            experiment_id=experiment_run.experiment_id,
        )
        self.__api = api
        self.__async_runner = async_runner
        self._logs: ExperimentRunArtifactTypes = ExperimentRunArtifactTypes()

    def __enter__(self):
        return self

    def __exit__(self, exc_type, exc_value, traceback):
        """Exit the context manager and upsert the logs to the experiment run."""
        temp_artifacts: ExperimentRunArtifactTypes = ExperimentRunArtifactTypes()
        for log in self.get_logs():
            server_log = self.__async_runner(
                self.__api.upsert_experiment_run_artifact(
                    experiment_run_id=self.info.id, artifact=log
                )
            )
            log_type = type(server_log)
            if log_type is Parameter:
                temp_artifacts.params[log.name] = server_log
            elif log_type is Metric:
                temp_artifacts.metrics[log.name] = server_log
            elif log_type is DataSet:
                temp_artifacts.datasets[log.name] = server_log
        self._logs = temp_artifacts
        self._thread_safe_ctx.set(RunContext(active=False))

    def log_param(self, key: str, value: Any, *, description: str | None = None):
        """Log a parameter to the experiment run.

        The parameter type is inferred from the value type.

        Args:
            key: Name of the parameter.
            value: Value of the parameter.

        Keyword Args:
            description: Description of the parameter.
        """
        param_type = ParameterType.STRING  # consistent with param_type default init
        if isinstance(value, bool):
            param_type = ParameterType.BOOLEAN
        elif isinstance(value, (int, float)):
            # TODO: ensure for numpy and other numeric types
            param_type = ParameterType.NUMBER
        elif isinstance(value, dict):
            param_type = ParameterType.OBJECT
        else:
            param_type = ParameterType.STRING
        self._logs.params[key] = Parameter(
            name=key,
            value=value,
            description=description,
            parameter_type=param_type,
        )

    def log_metric(
        self,
        key: str,
        value: Any,
        step: int = 0,
        timestamp: int | None = None,
        *,
        description: str | None = None,
    ):
        """Log a metric to the experiment run.

        Args:
            key: Name of the metric.
            value: Value of the metric.
            step: Step number for multi-step metrics (e.g., training epochs).
            timestamp: Unix timestamp in milliseconds when the metric was recorded.

        Keyword Args:
            description: Description of the metric.
        """
        self._logs.metrics[key] = Metric(
            name=key,
            value=value,
            step=step,
            state=ArtifactState.LIVE,
            timestamp=timestamp or str(int(time.time() * 1000)),
            description=description,
        )

    def log_dataset(
        self,
        name: str | None = None,
        uri: str | None = None,
        source_type: str | None = None,
        source: str | None = None,
        schema: dict | None = None,
        profile: dict | None = None,
        *,
        description: str | None = None,
        file_path: str | None = None,
        s3_auth: S3Params | None = None,
    ):
        """Log a dataset to the experiment run.

        Args:
            name: The name of the dataset.
            uri: The uri of the dataset if already uploaded to S3.
            source_type: Type of the source for the dataset.
            source: Location or connection string for the dataset source.
            schema: JSON schema or description of the dataset structure.
            profile: Statistical profile or summary of the dataset.

        Keyword Args:
            description: Description of the dataset.
            file_path: The path to the dataset file.
            s3_auth: S3Params for uploading the dataset to S3.
        """
        if not uri and not file_path:
            msg = "Either `uri` or `file_path` must be provided."
            raise ValueError(msg)
        try:
            uri = (
                upload_to_s3(
                    s3_auth=s3_auth,
                    path=file_path,
                )
                if file_path
                else uri
            )
        except Exception as e:
            msg = f"Failed to upload dataset to S3: {e}"
            raise StoreError(msg) from e
        self._logs.datasets[name] = DataSet(
            name=name,
            uri=uri,
            source_type=source_type,
            source=source,
            schema=schema,
            profile=profile,
            description=description,
        )

    def get_log(self, type: LogType, key: str) -> ExperimentRunArtifact:
        """Get a log from the experiment run.

        Args:
            type: Type of the log (params, metrics, or datasets).
            key: Key of the log.
        """
        return self._logs.__getattribute__(type)[key]

    def get_logs(self) -> list[ExperimentRunArtifact]:
        """Return every recorded artifact (params + metrics) in one flat list."""
        params = self._logs.params.values()
        metrics = self._logs.metrics.values()
        datasets = self._logs.datasets.values()

        return list(params) + list(metrics) + list(datasets)
