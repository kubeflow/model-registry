// TODO this hook was copied from odh-dashboard temporarily and should be abstracted out into mod-arch-shared.

import React from 'react';
import { useBrowserStorage } from 'mod-arch-core';
import { SortableData, CHECKBOX_FIELD_ID, KEBAB_FIELD_ID, EXPAND_FIELD_ID } from 'mod-arch-shared';

// Fields that are never manageable by users (UI chrome columns)
const NON_MANAGEABLE_FIELDS = [CHECKBOX_FIELD_ID, KEBAB_FIELD_ID, EXPAND_FIELD_ID];

/**
 * Represents a column that can be managed (shown/hidden, reordered)
 */
export interface ManagedColumn {
  /** Unique identifier for the column (typically matches SortableData.field) */
  id: string;
  /** Display label for the column */
  label: string;
  /** Whether the column is currently visible */
  isVisible: boolean;
}

/**
 * Configuration for the useManageColumns hook
 */
export interface UseManageColumnsConfig<T> {
  /** All possible columns (the full column definition) */
  allColumns: SortableData<T>[];
  /** Unique key for localStorage persistence */
  storageKey: string;
  /** Default visible column fields when no localStorage value exists */
  defaultVisibleFields?: string[];
  /** Maximum number of manageable columns that can be visible */
  maxVisibleColumns?: number;
}

/**
 * Return type for the useManageColumns hook
 */
export interface UseManageColumnsResult<T> {
  /** The columns to render in the table, filtered and ordered by visibility settings */
  visibleColumns: SortableData<T>[];
  /** All manageable columns with their current visibility state (for the modal) */
  managedColumns: ManagedColumn[];
  /** Callback to update which columns are visible (called from modal) */
  setVisibleColumnIds: (columnIds: string[]) => void;
  /** The currently visible column IDs (for display purposes like "X of Y selected") */
  visibleColumnIds: string[];
}

export const useManageColumns = <T>({
  allColumns,
  storageKey,
  defaultVisibleFields,
  maxVisibleColumns,
}: UseManageColumnsConfig<T>): UseManageColumnsResult<T> => {
  // Get manageable columns (those that can be shown/hidden)
  const manageableColumns = React.useMemo(
    () => allColumns.filter((col) => !NON_MANAGEABLE_FIELDS.includes(col.field)),
    [allColumns],
  );

  // Calculate default visible fields if not provided
  const effectiveDefaultVisibleFields = React.useMemo(() => {
    if (defaultVisibleFields) {
      return defaultVisibleFields;
    }
    // Default: show first maxVisibleColumns columns, or first 2 if not specified
    const manageableFields = manageableColumns.map((col) => col.field);
    const defaultCount = maxVisibleColumns ?? 2;
    return manageableFields.slice(0, defaultCount);
  }, [defaultVisibleFields, manageableColumns, maxVisibleColumns]);

  // Persist visible column IDs to localStorage
  const [storedVisibleIds, setStoredVisibleIds] = useBrowserStorage<string[]>(
    storageKey,
    effectiveDefaultVisibleFields,
    true, // jsonify
  );

  // Build the managed columns list for the modal
  // Preserves order from storage, adds any new columns at the end
  const managedColumns: ManagedColumn[] = React.useMemo(() => {
    const orderedIds = [...storedVisibleIds];
    const allManageableIds = manageableColumns.map((col) => col.field);

    // Add columns that are not in storage (new columns)
    allManageableIds.forEach((id) => {
      if (!orderedIds.includes(id)) {
        orderedIds.push(id);
      }
    });

    return orderedIds
      .filter((id) => allManageableIds.includes(id)) // Remove deleted columns
      .map((id) => {
        const col = manageableColumns.find((c) => c.field === id);
        if (!col) {
          return null;
        }
        return {
          id: col.field,
          label: col.label,
          isVisible: storedVisibleIds.includes(col.field),
        };
      })
      .filter((col): col is ManagedColumn => col !== null);
  }, [manageableColumns, storedVisibleIds]);

  // Build the final visible columns for the table
  const visibleColumns: SortableData<T>[] = React.useMemo(() => {
    const result: SortableData<T>[] = [];

    // First, add columns that come before manageable columns (like checkbox)
    allColumns.forEach((col) => {
      if (col.field === CHECKBOX_FIELD_ID || col.field === EXPAND_FIELD_ID) {
        result.push(col);
      }
    });

    // Add visible manageable columns in their stored order
    const visibleManageableIds = storedVisibleIds.filter((id) =>
      manageableColumns.some((col) => col.field === id),
    );

    visibleManageableIds.forEach((id) => {
      const col = manageableColumns.find((c) => c.field === id);
      if (col) {
        result.push(col);
      }
    });

    // Add kebab column at the end if present
    const kebabCol = allColumns.find((col) => col.field === KEBAB_FIELD_ID);
    if (kebabCol) {
      result.push(kebabCol);
    }

    return result;
  }, [allColumns, storedVisibleIds, manageableColumns]);

  const setVisibleColumnIds = React.useCallback(
    (columnIds: string[]) => {
      setStoredVisibleIds(columnIds);
    },
    [setStoredVisibleIds],
  );

  return {
    visibleColumns,
    managedColumns,
    setVisibleColumnIds,
    visibleColumnIds: storedVisibleIds,
  };
};
