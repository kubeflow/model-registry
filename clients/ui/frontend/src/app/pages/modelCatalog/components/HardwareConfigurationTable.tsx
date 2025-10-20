import * as React from 'react';
import { DashboardEmptyTableView, Table } from 'mod-arch-shared';
import { Spinner } from '@patternfly/react-core';
import { OuterScrollContainer } from '@patternfly/react-table';
import { CatalogPerformanceMetricsArtifact } from '~/app/modelCatalogTypes';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import {
  filterHardwareConfigurationArtifacts,
  clearAllFilters,
} from '~/app/pages/modelCatalog/utils/hardwareConfigurationFilterUtils';
import { getLatencyFilterConfig } from '~/app/pages/modelCatalog/utils/latencyFilterState';
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
  const { filterData, setFilterData } = React.useContext(ModelCatalogContext);

  // Apply filters to the artifacts
  const filteredArtifacts = React.useMemo(() => {
    const latencyConfig = getLatencyFilterConfig();
    return filterHardwareConfigurationArtifacts(performanceArtifacts, filterData, latencyConfig);
  }, [performanceArtifacts, filterData]);

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
        data={filteredArtifacts}
        columns={hardwareConfigColumns}
        toolbarContent={toolbarContent}
        onClearFilters={handleClearFilters}
        defaultSortColumn={0}
        emptyTableView={<DashboardEmptyTableView onClearFilters={handleClearFilters} />}
        rowRenderer={(artifact) => (
          <HardwareConfigurationTableRow
            key={artifact.customProperties.config_id?.string_value}
            performanceArtifact={artifact}
          />
        )}
      />
    </OuterScrollContainer>
  );
};

export default HardwareConfigurationTable;
