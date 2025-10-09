/* eslint-disable @typescript-eslint/consistent-type-assertions */
import * as React from 'react';
import { DashboardEmptyTableView, Table } from 'mod-arch-shared';
import { Spinner } from '@patternfly/react-core';
import { OuterScrollContainer } from '@patternfly/react-table';
import { CatalogPerformanceMetricsArtifact } from '~/app/modelCatalogTypes';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import { hardwareConfigColumns } from './HardwareConfigurationTableColumns';
import HardwareConfigurationTableRow from './HardwareConfigurationTableRow';
import HardwareConfigurationFilterToolbar from './HardwareConfigurationFilterToolbar';
import { useHardwareTypeFilterState } from '../utils/hardwareTypeFilterState';
import {
  filterHardwareConfigurationArtifacts,
  clearAllFilters,
} from '../utils/hardwareConfigurationFilterUtils';

type HardwareConfigurationTableProps = {
  performanceArtifacts: CatalogPerformanceMetricsArtifact[];
  isLoading?: boolean;
};

const HardwareConfigurationTable: React.FC<HardwareConfigurationTableProps> = ({
  performanceArtifacts,
  isLoading = false,
}) => {
  const { filterData, setFilterData } = React.useContext(ModelCatalogContext);
  const { appliedHardwareTypes, setAppliedHardwareTypes, clearHardwareFilters } =
    useHardwareTypeFilterState();

  // Apply filters to the artifacts
  const filteredArtifacts = React.useMemo(
    () =>
      filterHardwareConfigurationArtifacts(performanceArtifacts, filterData, appliedHardwareTypes),
    [performanceArtifacts, filterData, appliedHardwareTypes],
  );

  if (isLoading) {
    return <Spinner size="lg" />;
  }

  const toolbarContent = (
    <HardwareConfigurationFilterToolbar
      performanceArtifacts={performanceArtifacts}
      appliedHardwareTypes={appliedHardwareTypes}
      onApplyHardwareFilters={setAppliedHardwareTypes}
      onResetHardwareFilters={clearHardwareFilters}
    />
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
            key={`${artifact.customProperties.hardware?.string_value} ${artifact.customProperties.hardware_count?.int_value}`}
            performanceArtifact={artifact}
          />
        )}
      />
    </OuterScrollContainer>
  );
};

export default HardwareConfigurationTable;
