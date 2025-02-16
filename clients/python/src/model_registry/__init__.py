"""Main package for the Kubeflow model registry."""

__version__ = "0.2.14"

from ._client import ModelRegistry
from ._async_task_runner_base import AsyncTaskRunnerBase

__all__ = [
    "ModelRegistry",
    "AsyncTaskRunnerBase",
]