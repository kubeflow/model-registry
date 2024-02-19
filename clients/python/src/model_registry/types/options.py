"""Options for listing objects.

Provides a thin wrappers around the options classes defined in the MLMD Py lib.
"""

from __future__ import annotations

from enum import Enum

from attrs import define, field
from ml_metadata.metadata_store import ListOptions as MLMDListOptions
from ml_metadata.metadata_store import OrderByField as MLMDOrderByField


class OrderByField(Enum):
    """Fields to order by."""

    CREATE_TIME = MLMDOrderByField.CREATE_TIME
    UPDATE_TIME = MLMDOrderByField.UPDATE_TIME
    ID = MLMDOrderByField.ID


@define
class ListOptions:
    """Options for listing objects.

    Attributes:
        limit: Maximum number of objects to return.
        order_by: Field to order by.
        is_asc: Whether to order in ascending order. Defaults to True.
    """

    limit: int | None = field(default=None)
    order_by: OrderByField | None = field(default=None)
    is_asc: bool = field(default=True)

    def as_mlmd_list_options(self) -> MLMDListOptions:
        """Convert to MLMD ListOptions.

        Returns:
            MLMD ListOptions.
        """
        return MLMDListOptions(
            limit=self.limit,
            order_by=OrderByField(self.order_by).value if self.order_by else None,
            is_asc=self.is_asc,
        )
