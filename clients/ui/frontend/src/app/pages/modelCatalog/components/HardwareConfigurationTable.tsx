import * as React from 'react';
import { DashboardEmptyTableView, Table } from 'mod-arch-shared';
import { Button, Spinner } from '@patternfly/react-core';
import { ColumnsIcon } from '@patternfly/react-icons';
import { OuterScrollContainer } from '@patternfly/react-table';
import { CatalogPerformanceMetricsArtifact } from '~/app/modelCatalogTypes';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import { getActiveLatencyFieldName } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import { parseLatencyFilterKey } from '~/concepts/modelCatalog/const';
import { ManageColumnsModal } from '~/app/shared/components/manageColumns/ManageColumnsModal';
import HardwareConfigurationTableRow from './HardwareConfigurationTableRow';
import HardwareConfigurationFilterToolbar from './HardwareConfigurationFilterToolbar';
import { useHardwareConfigColumns } from './useHardwareConfigColumns';

type HardwareConfigurationTableProps = {
  performanceArtifacts: CatalogPerformanceMetricsArtifact[];
  isLoading?: boolean;
};

const HardwareConfigurationTable: React.FC<HardwareConfigurationTableProps> = ({
  performanceArtifacts,
  isLoading = false,
}) => {
  const { filterData, resetPerformanceFiltersToDefaults } = React.useContext(ModelCatalogContext);

  // Note: Filtering is now done server-side via the /performance_artifacts endpoint.
  // The performanceArtifacts prop contains pre-filtered data from the server.

  // Get the active latency filter field name (if any)
  const activeLatencyField = getActiveLatencyFieldName(filterData);

  // Use the custom hook that combines manage columns with latency filter effects
  const { columns, manageColumnsResult } = useHardwareConfigColumns(activeLatencyField);

  // When a latency filter is active sort by the corresponding latency column in ASC
  const { sortColumn, sortIndex } = React.useMemo(() => {
    if (!activeLatencyField) {
      return { sortColumn: null, sortIndex: 0 };
    }

    const { propertyKey } = parseLatencyFilterKey(activeLatencyField);
    const index = columns.findIndex((col) => col.field === propertyKey);

    return {
      sortColumn: index !== -1 ? columns[index] : null,
      sortIndex: index !== -1 ? index : 0,
    };
  }, [activeLatencyField, columns]);

  // Pre-sort the data by the active latency column as a fallback
  const sortedData = React.useMemo(() => {
    if (!sortColumn) {
      return performanceArtifacts;
    }
    const sortFn = sortColumn.sortable;
    if (typeof sortFn !== 'function') {
      return performanceArtifacts;
    }
    return [...performanceArtifacts].toSorted((a, b) => sortFn(a, b, sortColumn.field));
  }, [sortColumn, performanceArtifacts]);

  if (isLoading) {
    return <Spinner size="lg" />;
  }

  const handleClearFilters = () => {
    // On details page, reset performance filters to defaults (not basic filters from landing page)
    resetPerformanceFiltersToDefaults();
  };

  const manageColumnsButton = (
    <Button
      variant="link"
      icon={<ColumnsIcon />}
      onClick={manageColumnsResult.openModal}
      data-testid="manage-columns-button"
    >
      Customize columns
    </Button>
  );

  const toolbarContent = (
    <HardwareConfigurationFilterToolbar
      onResetAllFilters={handleClearFilters}
      includePerformanceFilters
      toolbarActions={manageColumnsButton}
    />
  );

  return (
    <>
      <OuterScrollContainer>
        <Table
          data-testid="hardware-configuration-table"
          variant="compact"
          isStickyHeader
          hasStickyColumns
          data={sortedData}
          columns={columns}
          toolbarContent={toolbarContent}
          onClearFilters={handleClearFilters}
          defaultSortColumn={sortIndex}
          emptyTableView={<DashboardEmptyTableView onClearFilters={handleClearFilters} />}
          rowRenderer={(artifact) => (
            <HardwareConfigurationTableRow
              key={artifact.customProperties?.config_id?.string_value}
              performanceArtifact={artifact}
              columns={columns}
            />
          )}
        />
      </OuterScrollContainer>
      <ManageColumnsModal
        manageColumnsResult={manageColumnsResult}
        description="Manage the columns that appear in the hardware configuration table."
        dataTestId="hardware-config-manage-columns"
      />
    </>
  );
};

export default HardwareConfigurationTable;
