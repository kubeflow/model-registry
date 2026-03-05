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

  const getFilterProps = () => filters;

  return (
    <Stack hasGutter>
      <DeploymentModeFilter filters={getFilterProps()} />
      <SupportedTransportsFilter filters={getFilterProps()} />
      <McpLicenseFilter filters={getFilterProps()} />
      <LabelsFilter filters={getFilterProps()} />
      <SecurityVerificationFilter filters={getFilterProps()} />
    </Stack>
  );
};

export default McpCatalogFilters;
