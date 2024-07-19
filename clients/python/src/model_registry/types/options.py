"""Options for listing objects.

Provides a thin wrappers around the options classes defined in the MLMD Py lib.
"""

from __future__ import annotations

from dataclasses import dataclass
from typing import Any

from mr_openapi import OrderByField, SortOrder


@dataclass
class ListOptions:
    """Options for listing objects.

    Attributes:
        limit: Maximum number of objects to return.
        order_by: Field to order by.
        is_asc: Whether to order in ascending order. Defaults to True.
        next_page_token: Token to use to retrieve next page of results.
    """

    limit: int | None = None
    order_by: OrderByField | None = None
    is_asc: bool = True
    next_page_token: str | None = None

    @classmethod
    def order_by_creation_time(cls, **kwargs) -> ListOptions:
        """Return options to order by creation time."""
        return cls(order_by=OrderByField.CREATE_TIME, **kwargs)

    @classmethod
    def order_by_update_time(cls, **kwargs) -> ListOptions:
        """Return options to order by update time."""
        return cls(order_by=OrderByField.LAST_UPDATE_TIME, **kwargs)

    @classmethod
    def order_by_id(cls, **kwargs) -> ListOptions:
        """Return options to order by ID."""
        return cls(order_by=OrderByField.ID, **kwargs)

    def as_options(self) -> dict[str, Any]:
        """Convert to options dictionary."""
        options = {}
        if self.limit is not None:
            options["page_size"] = str(self.limit)
        if self.order_by is not None:
            options["order_by"] = self.order_by
        if self.is_asc is not None:
            options["sort_order"] = SortOrder.ASC if self.is_asc else SortOrder.DESC
        if self.next_page_token is not None:
            options["next_page_token"] = self.next_page_token
        return options
