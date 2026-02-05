"""Utility functions for sorting tests."""

import pytest
from typing import Any


def get_field_value(item: dict[str, Any], field: str) -> Any:
    """Extract field value from an item for sorting comparison.

    Args:
        item: Model/item dictionary from API response
        field: Field name (ID, CREATE_TIME, LAST_UPDATE_TIME)

    Returns:
        The field value, converted appropriately for comparison

    Raises:
        ValueError: If field is invalid or missing from item
    """
    field_mapping = {
        "ID": "id",
        "CREATE_TIME": "createTimeSinceEpoch",
        "LAST_UPDATE_TIME": "lastUpdateTimeSinceEpoch",
    }

    if field not in field_mapping:
        raise ValueError(f"Invalid field: {field}")

    api_field = field_mapping[field]
    value = item.get(api_field)

    if value is None:
        raise ValueError(f"Field {field} ({api_field}) is missing from item: {item}")

    # Convert ID for proper numeric comparison
    if field == "ID":
        try:
            return int(value)
        except ValueError:
            return str(value)

    return value


def validate_items_sorted_correctly(items: list[dict], field: str, order: str) -> bool:
    """Verify items are sorted correctly by the specified field.

    Args:
        items: List of items to validate
        field: Field name to check sorting on (ID, CREATE_TIME, LAST_UPDATE_TIME)
        order: Sort order (ASC or DESC)

    Returns:
        True if items are sorted correctly, False otherwise
    """
    if len(items) < 2:
        pytest.fail("List has fewer than 2 items, double check the data you are passing to this function")

    values = [get_field_value(item, field) for item in items]

    if order == "ASC":
        return all(values[i] <= values[i + 1] for i in range(len(values) - 1))
    elif order == "DESC":
        return all(values[i] >= values[i + 1] for i in range(len(values) - 1))
    else:
        raise ValueError(f"Invalid sort order: {order}")
