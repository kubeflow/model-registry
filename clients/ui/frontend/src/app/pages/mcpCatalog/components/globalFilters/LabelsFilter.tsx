import * as React from 'react';
import { Divider, StackItem } from '@patternfly/react-core';
import McpCatalogStringFilter from '~/app/pages/mcpCatalog/components/McpCatalogStringFilter';
import type { McpCatalogFilterOptions } from '~/app/pages/mcpCatalog/types/mcpCatalogFilterOptions';

const filterKey = 'labels';

type LabelsFilterProps = {
  filters?: McpCatalogFilterOptions;
};

const LabelsFilter: React.FC<LabelsFilterProps> = ({ filters }) => {
  const value = filters?.[filterKey];
  if (!value) {
    return null;
  }
  return (
    <>
      <StackItem>
        <McpCatalogStringFilter title="Labels" filterKey={filterKey} filters={value} showSearch />
      </StackItem>
      <Divider />
    </>
  );
};

export default LabelsFilter;
