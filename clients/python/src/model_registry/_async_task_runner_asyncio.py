import asyncio
from collections.abc import Coroutine
from typing import Any

from ._async_task_runner_base import AsyncTaskRunnerBase

SINGLETON = "This class is a singleton!"


class AsyncTaskRunnerAcyncio(AsyncTaskRunnerBase):
    """A singleton task runner that runs an asyncio event loop on a background thread."""

    __instance = None

    @staticmethod
    def get_instance():
        """Get an AsyncTaskRunner (singleton)."""
        if AsyncTaskRunnerAcyncio.__instance is None:
            return AsyncTaskRunnerAcyncio()
        return AsyncTaskRunnerAcyncio.__instance

    def __init__(self):
        """Initialize."""
        if AsyncTaskRunnerAcyncio.__instance is not None:
            raise Exception(SINGLETON)
        import nest_asyncio
        nest_asyncio.apply()
        AsyncTaskRunnerAcyncio.__instance = self

    def run(self, coro: Coroutine) -> Any:
        """Synchronously run a coroutine on a background thread."""
        try:
            loop = asyncio.get_event_loop()
        except RuntimeError:
            loop = asyncio.new_event_loop()
            asyncio.set_event_loop(loop)
        return loop.run_until_complete(coro)
