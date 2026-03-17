import * as React from 'react';
import {
  Alert,
  Bullseye,
  Button,
  EmptyState,
  Flex,
  Grid,
  GridItem,
  Spinner,
  Title,
} from '@patternfly/react-core';
import { SearchIcon } from '@patternfly/react-icons';
import { McpCatalogContext } from '~/app/context/mcpCatalog/McpCatalogContext';
import { useMcpServersBySourceLabelWithAPI } from '~/app/hooks/mcpServerCatalog/useMcpServersBySourceLabel';
import { MCP_CATALOG_GRID_SPAN } from '~/app/pages/mcpCatalog/const';
import { filterMcpServersByFilters } from '~/app/pages/mcpCatalog/utils/mcpCatalogUtils';
import EmptyModelCatalogState from '~/app/pages/modelCatalog/EmptyModelCatalogState';
import ScrollViewOnMount from '~/app/shared/components/ScrollViewOnMount';
import McpCatalogCard from '~/app/pages/mcpCatalog/components/McpCatalogCard';

const PAGE_SIZE = 10;

type McpCatalogGalleryViewProps = {
  handleFilterReset: () => void;
};

const McpCatalogGalleryView: React.FC<McpCatalogGalleryViewProps> = ({ handleFilterReset }) => {
  const { mcpApiState, selectedSourceLabel, searchQuery, filters, catalogLabelsLoaded } =
    React.useContext(McpCatalogContext);

  const { mcpServers, mcpServersLoaded, mcpServersLoadError } = useMcpServersBySourceLabelWithAPI(
    mcpApiState,
    {
      sourceLabel: selectedSourceLabel,
      pageSize: PAGE_SIZE,
      searchQuery,
    },
  );

  const items = React.useMemo(
    () => filterMcpServersByFilters(mcpServers.items, filters),
    [mcpServers.items, filters],
  );

  const loaded = mcpServersLoaded && catalogLabelsLoaded;

  if (mcpServersLoadError) {
    return (
      <Alert variant="danger" title="Failed to load MCP servers" isInline>
        {mcpServersLoadError.message}
      </Alert>
    );
  }

  if (!loaded) {
    return (
      <EmptyState>
        <Spinner />
        <Title headingLevel="h4" size="lg">
          Loading MCP servers...
        </Title>
      </EmptyState>
    );
  }

  if (items.length === 0) {
    return (
      <EmptyModelCatalogState
        testid="empty-mcp-catalog-state"
        title="No results found"
        headerIcon={SearchIcon}
        description="Adjust your filters and try again."
        primaryAction={
          <Button variant="link" onClick={handleFilterReset}>
            Reset filters
          </Button>
        }
      />
    );
  }

  return (
    <>
      <ScrollViewOnMount shouldScroll scrollToTop />
      <Grid hasGutter>
        {items.map((server) => (
          <GridItem
            key={String(server.id)}
            sm={MCP_CATALOG_GRID_SPAN.sm}
            md={MCP_CATALOG_GRID_SPAN.md}
            lg={MCP_CATALOG_GRID_SPAN.lg}
            xl2={MCP_CATALOG_GRID_SPAN.xl2}
          >
            <McpCatalogCard server={server} />
          </GridItem>
        ))}
      </Grid>
      {mcpServers.hasMore && items.length >= PAGE_SIZE && (
        <Bullseye className="pf-v6-u-mt-lg">
          {mcpServers.isLoadingMore ? (
            <Flex
              direction={{ default: 'column' }}
              alignItems={{ default: 'alignItemsCenter' }}
              gap={{ default: 'gapMd' }}
            >
              <Spinner size="lg" />
              <Title size="lg" headingLevel="h5">
                Loading more MCP servers...
              </Title>
            </Flex>
          ) : (
            <Button variant="tertiary" onClick={mcpServers.loadMore} size="lg">
              Load more servers
            </Button>
          )}
        </Bullseye>
      )}
    </>
  );
};

export default McpCatalogGalleryView;
