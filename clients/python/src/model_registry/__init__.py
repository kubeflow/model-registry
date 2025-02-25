"""Main package for the Kubeflow model registry."""

__version__ = "0.2.15"

from ._async_task_runner_base import AsyncTaskRunnerBase
from ._client import ModelRegistry

__all__ = [
    "ModelRegistry",
    "AsyncTaskRunnerBase",
]
