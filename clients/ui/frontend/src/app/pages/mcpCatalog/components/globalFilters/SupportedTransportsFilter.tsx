import * as React from 'react';
import { Divider, StackItem } from '@patternfly/react-core';
import McpCatalogStringFilter from '~/app/pages/mcpCatalog/components/McpCatalogStringFilter';
import type { McpCatalogFilterOptions } from '~/app/pages/mcpCatalog/types/mcpCatalogFilterOptions';

const filterKey = 'supportedTransports';

type SupportedTransportsFilterProps = {
  filters?: McpCatalogFilterOptions;
};

const SupportedTransportsFilter: React.FC<SupportedTransportsFilterProps> = ({ filters }) => {
  const value = filters?.[filterKey];
  if (!value) {
    return null;
  }
  return (
    <>
      <StackItem>
        <McpCatalogStringFilter
          title="Supported transports"
          filterKey={filterKey}
          filters={value}
        />
      </StackItem>
      <Divider />
    </>
  );
};

export default SupportedTransportsFilter;
