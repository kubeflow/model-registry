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


def validate_custom_property_sorted(items: list[dict], property_field: str, sort_order: str) -> bool:
    """Verify items are sorted by a custom property value.

    Expected behavior:
    1. Items WITH the property appear first, sorted by value (ASC/DESC)
    2. Items WITHOUT the property appear after, sorted by ID ASC (fallback)

    Args:
        items: List of artifact items from API response
        property_field: Property field path (e.g., "accuracy.double_value")
        sort_order: Sort order (ASC or DESC)

    Returns:
        True if sorted correctly, False otherwise
    """
    if len(items) < 2:
        pytest.fail("List has fewer than 2 items, double check the data you are passing to this function")

    property_name, value_type = property_field.rsplit(".", 1)

    items_with = []
    items_without = []

    for item in items:
        custom_props = item.get("customProperties", {})
        if property_name in custom_props:
            value = custom_props[property_name].get(value_type)
            if value is not None:
                items_with.append((item, value))
            else:
                items_without.append(item)
        else:
            items_without.append(item)

    # Verify items with property are sorted by value
    if len(items_with) >= 2:
        values = [v for _, v in items_with]
        if sort_order == "ASC":
            if not all(values[i] <= values[i + 1] for i in range(len(values) - 1)):
                return False
        elif sort_order == "DESC":
            if not all(values[i] >= values[i + 1] for i in range(len(values) - 1)):
                return False

    # Verify items without property are sorted by ID ASC (fallback)
    if len(items_without) >= 2:
        ids = [int(item["id"]) for item in items_without]
        if not all(ids[i] <= ids[i + 1] for i in range(len(ids) - 1)):
            return False

    # Verify items with property come before items without
    if items_with and items_without:
        last_with_idx = items.index(items_with[-1][0])
        first_without_idx = items.index(items_without[0])
        if last_with_idx > first_without_idx:
            return False

    return True
