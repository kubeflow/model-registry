from ._async_task_runner_thread import AsyncTaskRunnerThread
from ._async_task_runner_base import AsyncTaskRunnerBase


class AsyncTaskRunnerFactory:
    """A factory to create an AsyncTaskRunner.

    A user can overwrite it to use his own AsyncTaskRunner implementation
    """
    @staticmethod
    def get_instance() -> AsyncTaskRunnerBase:
        return AsyncTaskRunnerThread.get_instance()

