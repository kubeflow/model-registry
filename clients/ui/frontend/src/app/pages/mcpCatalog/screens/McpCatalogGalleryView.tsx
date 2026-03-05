import * as React from 'react';
import { EmptyState, EmptyStateBody, Grid, GridItem, Stack } from '@patternfly/react-core';
import { SearchIcon } from '@patternfly/react-icons';
import { McpCatalogContext } from '~/app/context/mcpCatalog/McpCatalogContext';
import McpCatalogCard from '~/app/pages/mcpCatalog/components/McpCatalogCard';

type McpCatalogGalleryViewProps = {
  searchTerm: string;
};

const McpCatalogGalleryView: React.FC<McpCatalogGalleryViewProps> = () => {
  const { mcpServers, mcpServersLoaded, mcpServersLoadError } = React.useContext(McpCatalogContext);
  const { items } = mcpServers;

  if (mcpServersLoadError) {
    return (
      <EmptyState
        icon={SearchIcon}
        headingLevel="h3"
        titleText="Unable to load MCP servers"
        data-testid="mcp-catalog-load-error"
      >
        <EmptyStateBody>{mcpServersLoadError.message}</EmptyStateBody>
      </EmptyState>
    );
  }

  if (mcpServersLoaded && items.length === 0) {
    return (
      <EmptyState
        icon={SearchIcon}
        headingLevel="h3"
        titleText="No result found"
        data-testid="mcp-catalog-empty-search"
      >
        <EmptyStateBody>Adjust your filters and try again.</EmptyStateBody>
      </EmptyState>
    );
  }

  if (!mcpServersLoaded && items.length === 0) {
    return null;
  }

  return (
    <Stack hasGutter>
      <Grid hasGutter>
        {items.map((server) => (
          <GridItem key={String(server.id)} sm={12} md={6} lg={4} xl2={4}>
            <McpCatalogCard server={server} />
          </GridItem>
        ))}
      </Grid>
    </Stack>
  );
};

export default McpCatalogGalleryView;
