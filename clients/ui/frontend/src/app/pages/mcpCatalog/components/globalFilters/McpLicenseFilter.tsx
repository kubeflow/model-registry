import * as React from 'react';
import { Divider, StackItem } from '@patternfly/react-core';
import McpCatalogStringFilter from '~/app/pages/mcpCatalog/components/McpCatalogStringFilter';
import type { McpCatalogFilterOptions } from '~/app/pages/mcpCatalog/types/mcpCatalogFilterOptions';

const filterKey = 'license';

type McpLicenseFilterProps = {
  filters?: McpCatalogFilterOptions;
};

const McpLicenseFilter: React.FC<McpLicenseFilterProps> = ({ filters }) => {
  const value = filters?.[filterKey];
  if (!value) {
    return null;
  }
  return (
    <>
      <StackItem>
        <McpCatalogStringFilter title="License" filterKey={filterKey} filters={value} />
      </StackItem>
      <Divider />
    </>
  );
};

export default McpLicenseFilter;
