"""Pager for iterating over items."""

from __future__ import annotations

import asyncio
from collections.abc import AsyncIterator, Awaitable, Iterator
from dataclasses import dataclass, field
from typing import Callable, Generic, TypeVar, cast

from .base import BaseModel
from .options import ListOptions, OrderByField

T = TypeVar("T", bound=BaseModel)


@dataclass
class Pager(Generic[T], Iterator[T], AsyncIterator[T]):
    """Pager for iterating over items.

    Assumes that page_fn is a paged function that takes ListOptions and returns a list of items.
    """

    page_fn: (
        Callable[[ListOptions], list[T]] | Callable[[ListOptions], Awaitable[list[T]]]
    )
    options: ListOptions = field(default_factory=ListOptions)

    def __post_init__(self):
        self.restart()
        if asyncio.iscoroutinefunction(self.page_fn):
            self.__next__ = NotImplemented
            self.next_page = self._anext_page
            self.next_item = self._anext_item
        else:
            self.__anext__ = NotImplemented
            self.next_page = self._next_page
            self.next_item = self._next_item

    def restart(self) -> Pager[T]:
        """Reset the pager.

        This keeps the current options and page function, but resets the internal state.
        """
        # as MLMD loops over pages, we need to keep track of the first page or we'll loop forever
        self._start = None
        self._current_page = None
        # tracks the next item on the current page
        self._i = 0
        self.options.next_page_token = None
        return self

    def order_by_creation_time(self) -> Pager[T]:
        """Order items by creation time.

        This resets the pager.
        """
        self.options.order_by = OrderByField.CREATE_TIME
        return self.restart()

    def order_by_update_time(self) -> Pager[T]:
        """Order items by update time.

        This resets the pager.
        """
        self.options.order_by = OrderByField.LAST_UPDATE_TIME
        return self.restart()

    def order_by_id(self) -> Pager[T]:
        """Order items by ID.

        This resets the pager.
        """
        self.options.order_by = OrderByField.ID
        return self.restart()

    def limit(self, limit: int) -> Pager[T]:
        """Limit the number of items to return.

        This resets the pager.
        """
        self.options.limit = limit
        return self.restart()

    def ascending(self) -> Pager[T]:
        """Order items in ascending order.

        This resets the pager.
        """
        self.options.is_asc = True
        return self.restart()

    def descending(self) -> Pager[T]:
        """Order items in descending order.

        This resets the pager.
        """
        self.options.is_asc = False
        return self.restart()

    def _next_page(self) -> list[T]:
        """Get the next page of items.

        This will automatically loop over pages.
        """
        return cast(list[T], self.page_fn(self.options))

    async def _anext_page(self) -> list[T]:
        """Get the next page of items.

        This will automatically loop over pages.
        """
        return await cast(Awaitable[list[T]], self.page_fn(self.options))

    def _needs_fetch(self) -> bool:
        return not self._current_page or self._i >= len(self._current_page)

    def _next_item(self) -> T:
        """Get the next item in the pager.

        This variant won't check for looping, so it's useful for manual iteration/scripting.

        NOTE: This won't check for looping, so use with caution.
        If you want to check for looping, use the pythonic `next()`.
        """
        if self._needs_fetch():
            self._current_page = self._next_page()
            self._i = 0
        assert self._current_page

        item = self._current_page[self._i]
        self._i += 1
        return item

    async def _anext_item(self) -> T:
        """Get the next item in the pager.

        This variant won't check for looping, so it's useful for manual iteration/scripting.

        NOTE: This won't check for looping, so use with caution.
        If you want to check for looping, use the pythonic `next()`.
        """
        if self._needs_fetch():
            self._current_page = await self._anext_page()
            self._i = 0
        assert self._current_page

        item = self._current_page[self._i]
        self._i += 1
        return item

    def __next__(self) -> T:
        check_looping = self._needs_fetch()

        item = self._next_item()

        if not self._start:
            self._start = self.options.next_page_token
        elif check_looping and self.options.next_page_token == self._start:
            raise StopIteration

        return item

    async def __anext__(self) -> T:
        check_looping = self._needs_fetch()

        item = await self._anext_item()

        if not self._start:
            self._start = self.options.next_page_token
        elif check_looping and self.options.next_page_token == self._start:
            raise StopAsyncIteration

        return item

    def __iter__(self) -> Iterator[T]:
        return self

    def __aiter__(self) -> AsyncIterator[T]:
        return self
