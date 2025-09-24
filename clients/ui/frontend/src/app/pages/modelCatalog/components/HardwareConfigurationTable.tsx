import * as React from 'react';
import { Table, DashboardEmptyTableView } from 'mod-arch-shared';
import {
  HardwareConfiguration,
  HardwareConfigurationTableProps,
} from '~/app/pages/modelCatalog/types/hardwareConfiguration';
import { hardwareConfigColumns } from './HardwareConfigurationTableColumns';
import HardwareConfigurationTableRow from './HardwareConfigurationTableRow';

const HardwareConfigurationTable = ({
  configurations,
  isLoading = false,
}: HardwareConfigurationTableProps): React.JSX.Element => {
  const clearFilters = React.useCallback(() => {
    // No filters for now, but keeping the interface consistent
  }, []);

  return (
    <Table
      data-testid="hardware-configuration-table"
      data={configurations}
      columns={hardwareConfigColumns}
      defaultSortColumn={0}
      enablePagination
      emptyTableView={<DashboardEmptyTableView onClearFilters={clearFilters} />}
      rowRenderer={(config: HardwareConfiguration) => (
        <HardwareConfigurationTableRow key={config.id} configuration={config} />
      )}
      loading={isLoading}
    />
  );
};

export default HardwareConfigurationTable;
