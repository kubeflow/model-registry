from collections.abc import Coroutine
from typing import Any

NOT_IMPLEMENTED = "Must be implemented by subclass"


class AsyncTaskRunnerBase:
    """A base task runner that runs an asyncio event loop on a background thread.

    A user can add his own representation of this class
    """
    def run(self, coro: Coroutine) -> Any:
        """Synchronously run a coroutine on a background thread."""
        raise ValueError(NOT_IMPLEMENTED)
