import * as React from 'react';
import { Alert, Spinner, Stack } from '@patternfly/react-core';
import { McpCatalogContext } from '~/app/context/mcpCatalog/McpCatalogContext';
import DeploymentModeFilter from '~/app/pages/mcpCatalog/components/globalFilters/DeploymentModeFilter';
import SupportedTransportsFilter from '~/app/pages/mcpCatalog/components/globalFilters/SupportedTransportsFilter';
import McpLicenseFilter from '~/app/pages/mcpCatalog/components/globalFilters/McpLicenseFilter';
import LabelsFilter from '~/app/pages/mcpCatalog/components/globalFilters/LabelsFilter';
import SecurityVerificationFilter from '~/app/pages/mcpCatalog/components/globalFilters/SecurityVerificationFilter';

const McpCatalogFilters: React.FC = () => {
  const { filterOptions, filterOptionsLoaded, filterOptionsLoadError } =
    React.useContext(McpCatalogContext);

  if (!filterOptionsLoaded) {
    return <Spinner />;
  }

  if (filterOptionsLoadError) {
    return (
      <Alert variant="danger" title="Failed to load filter options" isInline>
        {filterOptionsLoadError.message}
      </Alert>
    );
  }

  const filters = filterOptions?.filters;

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
