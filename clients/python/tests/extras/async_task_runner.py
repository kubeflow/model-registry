# from https://gist.github.com/blink1073/969aeba85f32c285235750626f2eadd8
# Slightly reformatted for linting purposes March 3, 2025, no logic changed.
# License below. Not shipped for use within the library code, only used for test examples.

"""
Copyright (c) 2022 Steven Silvester
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its
   contributors may be used to endorse or promote products derived from
   this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
"""

import asyncio
import atexit
from collections.abc import Coroutine
from threading import Lock, Thread
from typing import Any, Optional


class AsyncTaskRunner:
    """A singleton task runner that runs an asyncio event loop on a background thread."""

    __instance = None

    @staticmethod
    def get_instance():
        """Get an AsyncTaskRunner (singleton)."""
        if AsyncTaskRunner.__instance is None:
            AsyncTaskRunner()
        assert AsyncTaskRunner.__instance is not None
        return AsyncTaskRunner.__instance

    def __init__(self):
        """Initialize."""
        # make sure it is a singleton
        if AsyncTaskRunner.__instance is not None:
            msg = "This class is a singleton!"
            raise Exception(msg)
        AsyncTaskRunner.__instance = self
        # initialize variables
        self.__io_loop: Optional[asyncio.AbstractEventLoop] = None
        self.__runner_thread: Optional[Thread] = None
        self.__lock = Lock()
        # register exit handler
        atexit.register(self._close)

    def _close(self):
        """Clean up. Stop the loop if running."""
        if self.__io_loop:
            self.__io_loop.stop()

    def _runner(self) -> None:
        """Function to run in a thread."""
        loop = self.__io_loop
        assert loop is not None
        try:
            loop.run_forever()
        finally:
            loop.close()

    def run(self, coro: Coroutine) -> Any:
        """Synchronously run a coroutine on a background thread."""
        with self.__lock:
            if self.__io_loop is None:
                # If the asyncio loop does not exist
                self.__io_loop = asyncio.new_event_loop()
                self.__runner_thread = Thread(target=self._runner, daemon=True)
                self.__runner_thread.start()
        # run coroutine thread safe inside a thread. This return concurrent future
        fut = asyncio.run_coroutine_threadsafe(coro, self.__io_loop)
        # get concurrent future result
        return fut.result()
