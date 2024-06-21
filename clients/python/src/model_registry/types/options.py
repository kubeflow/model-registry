"""Options for listing objects.

Provides a thin wrappers around the options classes defined in the MLMD Py lib.
"""

from __future__ import annotations

from dataclasses import dataclass


@dataclass
class ListOptions:
    """Options for listing objects.

    Attributes:
        limit: Maximum number of objects to return.
        order_by: Field to order by.
        is_asc: Whether to order in ascending order. Defaults to True.
    """

    limit: int | None = None
    # order_by: OrderByField | None = None
    is_asc: bool = True
