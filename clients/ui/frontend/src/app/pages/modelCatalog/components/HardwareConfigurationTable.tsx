import * as React from 'react';
import { DashboardEmptyTableView, Table } from 'mod-arch-shared';
import { Spinner } from '@patternfly/react-core';
import { OuterScrollContainer } from '@patternfly/react-table';
import { CatalogPerformanceMetricsArtifact } from '~/app/modelCatalogTypes';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import {
  clearAllFilters,
  parseLatencyFieldName,
} from '~/app/pages/modelCatalog/utils/hardwareConfigurationFilterUtils';
import { getActiveLatencyFieldName } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import {
  ALL_LATENCY_FIELD_NAMES,
  getLatencyFieldName,
  LatencyMetric,
} from '~/concepts/modelCatalog/const';
import { hardwareConfigColumns, HardwareConfigColumn } from './HardwareConfigurationTableColumns';
import HardwareConfigurationTableRow from './HardwareConfigurationTableRow';
import HardwareConfigurationFilterToolbar from './HardwareConfigurationFilterToolbar';

type HardwareConfigurationTableProps = {
  performanceArtifacts: CatalogPerformanceMetricsArtifact[];
  isLoading?: boolean;
};

const HardwareConfigurationTable: React.FC<HardwareConfigurationTableProps> = ({
  performanceArtifacts,
  isLoading = false,
}) => {
  const { setFilterData, filterData } = React.useContext(ModelCatalogContext);

  // Note: Filtering is now done server-side via the /performance_artifacts endpoint.
  // The performanceArtifacts prop contains pre-filtered data from the server.

  // Get the active latency filter field name (if any)
  const activeLatencyField = getActiveLatencyFieldName(filterData);

  // When a latency filter is selected, show only that column and hide other latency columns
  // Also show the TPS column with the matching percentile (e.g., TTFT P90 filter shows TPS P90)
  const filteredColumns = React.useMemo((): HardwareConfigColumn[] => {
    if (!activeLatencyField) {
      // No latency filter selected, show all columns
      return hardwareConfigColumns;
    }

    // Parse the active filter field name to extract metric and percentile
    const parsed = parseLatencyFieldName(activeLatencyField);
    if (!parsed) {
      return hardwareConfigColumns;
    }

    // Build the matching TPS field name using the same percentile (e.g., TTFT P90 filter shows TPS P90)
    const matchingTpsField = getLatencyFieldName(LatencyMetric.TPS, parsed.percentile);

    // Filter out latency columns that don't match the active filter
    return hardwareConfigColumns.filter((column) => {
      // Check if this column is a latency column
      const isLatencyColumn = ALL_LATENCY_FIELD_NAMES.some(
        (fieldName) => fieldName === column.field,
      );

      // If it's not a latency column, keep it
      if (!isLatencyColumn) {
        return true;
      }

      // Show TPS column with matching percentile (they measure throughput, not latency delay)
      if (column.field === matchingTpsField) {
        return true;
      }

      // If it's a latency column (not TPS), only keep it if it matches the active filter
      return column.field === activeLatencyField;
    });
  }, [activeLatencyField]);

  if (isLoading) {
    return <Spinner size="lg" />;
  }

  const toolbarContent = (
    <HardwareConfigurationFilterToolbar performanceArtifacts={performanceArtifacts} />
  );
  const handleClearFilters = () => {
    clearAllFilters(setFilterData);
  };

  return (
    <OuterScrollContainer>
      <Table
        data-testid="hardware-configuration-table"
        variant="compact"
        isStickyHeader
        hasStickyColumns
        data={performanceArtifacts}
        columns={filteredColumns}
        toolbarContent={toolbarContent}
        onClearFilters={handleClearFilters}
        defaultSortColumn={0}
        emptyTableView={<DashboardEmptyTableView onClearFilters={handleClearFilters} />}
        rowRenderer={(artifact) => (
          <HardwareConfigurationTableRow
            key={artifact.customProperties?.config_id?.string_value}
            performanceArtifact={artifact}
            columns={filteredColumns}
          />
        )}
      />
    </OuterScrollContainer>
  );
};

export default HardwareConfigurationTable;
