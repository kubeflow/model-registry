import * as React from 'react';
import { DashboardEmptyTableView, Table } from 'mod-arch-shared';
import { Button, Spinner } from '@patternfly/react-core';
import { ColumnsIcon } from '@patternfly/react-icons';
import { OuterScrollContainer } from '@patternfly/react-table';
import { CatalogPerformanceMetricsArtifact } from '~/app/modelCatalogTypes';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import { getActiveLatencyFieldName } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import { ManageColumnsModal } from '~/app/shared/components/manageColumns/ManageColumnsModal';
import HardwareConfigurationTableRow from './HardwareConfigurationTableRow';
import HardwareConfigurationFilterToolbar from './HardwareConfigurationFilterToolbar';
import { useHardwareConfigColumns, ControlledTableSortProps } from './useHardwareConfigColumns';

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

  // Use the custom hook that combines manage columns with the latency filter + sort logic
  const {
    columns,
    manageColumnsResult,
    sortState: { sortIndex, sortDirection, onSortIndexChange, onSortDirectionChange },
  } = useHardwareConfigColumns(activeLatencyField);

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

  // Controlled sort props exist at runtime but not in mod-arch-shared Table typings yet
  const controlledSortProps: ControlledTableSortProps = {
    sortIndex,
    sortDirection,
    onSortIndexChange,
    onSortDirectionChange,
  };

  return (
    <>
      <OuterScrollContainer>
        <Table
          data-testid="hardware-configuration-table"
          variant="compact"
          isStickyHeader
          hasStickyColumns
          data={performanceArtifacts}
          columns={columns}
          toolbarContent={toolbarContent}
          onClearFilters={handleClearFilters}
          defaultSortColumn={sortIndex}
          {...controlledSortProps}
          emptyTableView={<DashboardEmptyTableView onClearFilters={handleClearFilters} />}
          rowRenderer={(artifact: CatalogPerformanceMetricsArtifact) => (
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
