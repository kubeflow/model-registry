import * as React from 'react';
import { Stack } from '@patternfly/react-core';
import { McpCatalogContext } from '~/app/context/mcpCatalog/McpCatalogContext';
import { mockMcpCatalogFilterOptions } from '~/app/pages/mcpCatalog/mocks/mockMcpCatalogFilterOptions';
import DeploymentModeFilter from '~/app/pages/mcpCatalog/components/globalFilters/DeploymentModeFilter';
import SupportedTransportsFilter from '~/app/pages/mcpCatalog/components/globalFilters/SupportedTransportsFilter';
import McpLicenseFilter from '~/app/pages/mcpCatalog/components/globalFilters/McpLicenseFilter';
import LabelsFilter from '~/app/pages/mcpCatalog/components/globalFilters/LabelsFilter';
import SecurityVerificationFilter from '~/app/pages/mcpCatalog/components/globalFilters/SecurityVerificationFilter';

const McpCatalogFilters: React.FC = () => {
  const { filterOptions } = React.useContext(McpCatalogContext);
  const filters = filterOptions?.filters ?? mockMcpCatalogFilterOptions.filters;

  return (
    <Stack hasGutter>
      <DeploymentModeFilter filters={filters} />
      <SupportedTransportsFilter filters={filters} />
      <McpLicenseFilter filters={filters} />
      <LabelsFilter filters={filters} />
      <SecurityVerificationFilter filters={filters} />
    </Stack>
  );
};

export default McpCatalogFilters;
