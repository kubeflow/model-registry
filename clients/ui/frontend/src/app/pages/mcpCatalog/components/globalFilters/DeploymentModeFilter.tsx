import * as React from 'react';
import { Divider, StackItem } from '@patternfly/react-core';
import McpCatalogStringFilter from '~/app/pages/mcpCatalog/components/McpCatalogStringFilter';
import type { McpCatalogFilterOptions } from '~/app/pages/mcpCatalog/types/mcpCatalogFilterOptions';

const filterKey = 'deploymentMode';

type DeploymentModeFilterProps = {
  filters?: McpCatalogFilterOptions;
};

const DeploymentModeFilter: React.FC<DeploymentModeFilterProps> = ({ filters }) => {
  const value = filters?.[filterKey];
  if (!value) {
    return null;
  }
  return (
    <>
      <StackItem>
        <McpCatalogStringFilter title="Deployment mode" filterKey={filterKey} filters={value} />
      </StackItem>
      <Divider />
    </>
  );
};

export default DeploymentModeFilter;
