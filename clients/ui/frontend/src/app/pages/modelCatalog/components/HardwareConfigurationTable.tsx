/* eslint-disable @typescript-eslint/consistent-type-assertions */
import * as React from 'react';
import { DashboardEmptyTableView, Table } from 'mod-arch-shared';
import { Spinner } from '@patternfly/react-core';
import { OuterScrollContainer } from '@patternfly/react-table';
import { CatalogPerformanceMetricsArtifact } from '~/app/modelCatalogTypes';
import { hardwareConfigColumns } from './HardwareConfigurationTableColumns';
import HardwareConfigurationTableRow from './HardwareConfigurationTableRow';

type HardwareConfigurationTableProps = {
  performanceArtifacts: CatalogPerformanceMetricsArtifact[];
  isLoading?: boolean;
};

const HardwareConfigurationTable: React.FC<HardwareConfigurationTableProps> = ({
  performanceArtifacts,
  isLoading = false,
}) => {
  if (isLoading) {
    return <Spinner size="lg" />;
  }

  // TODO when we add filters - lift these out as props, reference what tables like RegisteredModelTable do
  const toolbarContent = <>TODO filter toolbar goes here</>;
  const clearFilters = () => {
    // TODO
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
        onClearFilters={clearFilters}
        defaultSortColumn={0}
        emptyTableView={<DashboardEmptyTableView onClearFilters={clearFilters} />}
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
