import * as React from 'react';
import { Alert, Divider, Spinner, Stack, StackItem } from '@patternfly/react-core';
import { McpCatalogContext } from '~/app/context/mcpCatalog/McpCatalogContext';
import McpCatalogStringFilter from '~/app/pages/mcpCatalog/components/McpCatalogStringFilter';
import type { McpFilterCategoryKey } from '~/app/pages/mcpCatalog/types/mcpCatalogFilterOptions';

type FilterConfig = {
  filterKey: McpFilterCategoryKey;
  title: string;
  hasDivider: boolean;
};

const FILTER_CONFIGS: FilterConfig[] = [
  { filterKey: 'deploymentMode', title: 'Deployment mode', hasDivider: true },
  { filterKey: 'supportedTransports', title: 'Supported transports', hasDivider: true },
  { filterKey: 'license', title: 'License', hasDivider: true },
  { filterKey: 'labels', title: 'Labels', hasDivider: true },
  { filterKey: 'securityVerification', title: 'Security & Verification', hasDivider: false },
];

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
      {FILTER_CONFIGS.map(({ filterKey, title, hasDivider }) => {
        const value = filters?.[filterKey];
        if (!value) {
          return null;
        }
        return (
          <React.Fragment key={filterKey}>
            <StackItem>
              <McpCatalogStringFilter title={title} filterKey={filterKey} filters={value} />
            </StackItem>
            {hasDivider && <Divider />}
          </React.Fragment>
        );
      })}
    </Stack>
  );
};

export default McpCatalogFilters;
