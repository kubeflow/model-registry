import * as React from 'react';
import { DashboardEmptyTableView, Table } from 'mod-arch-shared';
import { Spinner } from '@patternfly/react-core';
import { OuterScrollContainer } from '@patternfly/react-table';
import { CatalogPerformanceMetricsArtifact } from '~/app/modelCatalogTypes';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import { clearAllFilters } from '~/app/pages/modelCatalog/utils/hardwareConfigurationFilterUtils';
import { hardwareConfigColumns } from './HardwareConfigurationTableColumns';
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
  const { setFilterData } = React.useContext(ModelCatalogContext);

  // Note: Filtering is now done server-side via the /performance_artifacts endpoint.
  // The performanceArtifacts prop contains pre-filtered data from the server.

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
        columns={hardwareConfigColumns}
        toolbarContent={toolbarContent}
        onClearFilters={handleClearFilters}
        defaultSortColumn={0}
        emptyTableView={<DashboardEmptyTableView onClearFilters={handleClearFilters} />}
        rowRenderer={(artifact) => (
          <HardwareConfigurationTableRow
            key={artifact.customProperties?.config_id?.string_value}
            performanceArtifact={artifact}
          />
        )}
      />
    </OuterScrollContainer>
  );
};

export default HardwareConfigurationTable;
