import * as React from 'react';
import { CatalogPerformanceMetricsArtifact } from '~/app/modelCatalogTypes';
import {
  useManageColumns,
  UseManageColumnsResult,
} from '~/app/shared/components/manageColumns/useManageColumns';
import {
  LatencyMetric,
  LatencyMetricFieldName,
  getLatencyPropertyKey,
  parseLatencyFilterKey,
} from '~/concepts/modelCatalog/const';
import {
  hardwareConfigColumns,
  HardwareConfigColumn,
  HardwareConfigColumnField,
  STICKY_COLUMN_FIELDS,
  DEFAULT_VISIBLE_COLUMN_FIELDS,
  LATENCY_COLUMN_FIELDS,
  TPS_COLUMN_FIELDS,
  HARDWARE_CONFIG_COLUMNS_STORAGE_KEY,
} from './HardwareConfigurationTableColumns';

interface UseHardwareConfigColumnsResult {
  /** Final columns to render in the table (sticky + visible managed columns) */
  columns: HardwareConfigColumn[];
  /** Result from useManageColumns hook, to be passed directly to ManageColumnsModal */
  manageColumnsResult: UseManageColumnsResult<CatalogPerformanceMetricsArtifact>;
}

/**
 * Type guard to check if a string is a HardwareConfigColumnField
 */
const isHardwareConfigColumnField = (id: string): id is HardwareConfigColumnField =>
  hardwareConfigColumns.some((col) => col.field === id);

/**
 * Check if a column field is a latency column (TTFT, E2E, ITL - not TPS)
 */
const isLatencyColumnField = (field: string): boolean =>
  LATENCY_COLUMN_FIELDS.some((f) => f === field);

/**
 * Check if a column field is a TPS column
 */
const isTpsColumnField = (field: string): boolean => TPS_COLUMN_FIELDS.some((f) => f === field);

/**
 * Custom hook that combines useManageColumns with the latency filter effect logic.
 *
 * When the latency filter changes:
 * - The filtered latency column becomes selected
 * - Other latency columns become deselected
 * - The corresponding TPS column becomes selected
 * - Other TPS columns become deselected
 *
 * When the filter is cleared, the current state is preserved.
 */
export const useHardwareConfigColumns = (
  activeLatencyField: LatencyMetricFieldName | undefined,
): UseHardwareConfigColumnsResult => {
  // Track the previous latency filter to detect changes
  const prevLatencyFieldRef = React.useRef<LatencyMetricFieldName | undefined>(undefined);

  // Separate sticky columns (always visible) from manageable columns
  const { stickyColumns, manageableColumns } = React.useMemo(() => {
    const sticky = hardwareConfigColumns.filter((col) => STICKY_COLUMN_FIELDS.includes(col.field));
    const manageable = hardwareConfigColumns.filter(
      (col) => !STICKY_COLUMN_FIELDS.includes(col.field),
    );
    return { stickyColumns: sticky, manageableColumns: manageable };
  }, []);

  // Use the manage columns hook for manageable columns only
  const manageColumnsResult = useManageColumns<CatalogPerformanceMetricsArtifact>({
    allColumns: manageableColumns,
    storageKey: HARDWARE_CONFIG_COLUMNS_STORAGE_KEY,
    defaultVisibleColumnIds: DEFAULT_VISIBLE_COLUMN_FIELDS,
  });

  // Effect to update visible columns when latency filter changes
  React.useEffect(() => {
    // Only react to changes, not initial mount
    if (prevLatencyFieldRef.current === activeLatencyField) {
      return;
    }
    prevLatencyFieldRef.current = activeLatencyField;

    // If filter is cleared, keep current state
    if (!activeLatencyField) {
      return;
    }

    // Parse the active filter to get metric and percentile
    const parsed = parseLatencyFilterKey(activeLatencyField);
    const activePropertyKey = parsed.propertyKey;
    const matchingTpsPropertyKey = getLatencyPropertyKey(LatencyMetric.TPS, parsed.percentile);

    // Build new visible column IDs:
    // - Keep all non-latency/non-TPS columns
    // - Remove all latency and TPS columns
    // - Add the active latency column and matching TPS column
    const newVisibleIds = manageColumnsResult.visibleColumnIds.filter((id) => {
      const isLatencyColumn = isLatencyColumnField(id);
      const isTpsColumn = isTpsColumnField(id);
      return !isLatencyColumn && !isTpsColumn;
    });

    // Add the active latency column (if not already present)
    if (!newVisibleIds.includes(activePropertyKey)) {
      newVisibleIds.push(activePropertyKey);
    }

    // Add the matching TPS column (if not already present)
    if (!newVisibleIds.includes(matchingTpsPropertyKey)) {
      newVisibleIds.push(matchingTpsPropertyKey);
    }

    manageColumnsResult.setVisibleColumnIds(newVisibleIds);
  }, [activeLatencyField, manageColumnsResult]);

  // Combine sticky + visible managed columns
  // Map SortableData back to HardwareConfigColumn by finding the original column
  const columns = React.useMemo((): HardwareConfigColumn[] => {
    const visibleManaged = manageColumnsResult.visibleColumns
      .map((sortableCol) => {
        // Find the original HardwareConfigColumn by field
        if (isHardwareConfigColumnField(sortableCol.field)) {
          return manageableColumns.find((col) => col.field === sortableCol.field);
        }
        return undefined;
      })
      .filter((col): col is HardwareConfigColumn => col !== undefined);

    return [...stickyColumns, ...visibleManaged];
  }, [stickyColumns, manageColumnsResult.visibleColumns, manageableColumns]);

  return {
    columns,
    manageColumnsResult,
  };
};
