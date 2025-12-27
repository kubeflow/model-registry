import React from 'react';
import { Grid, GridItem } from '@patternfly/react-core';
import { SearchIcon } from '@patternfly/react-icons';
import { useMcpCatalog } from '~/app/context/mcpCatalog/McpCatalogContext';
import { McpCategoryName, McpSourceLabel, McpServer } from '~/app/pages/mcpCatalog/types';
import { hasMcpFiltersActive } from '~/app/pages/mcpCatalog/utils/mcpCatalogUtils';
import McpCatalogCard from '~/app/pages/mcpCatalog/components/McpCatalogCard';
import EmptyMcpCatalogState from '~/app/pages/mcpCatalog/EmptyMcpCatalogState';

/**
 * Gallery view for MCP servers.
 * Displays servers in a grid layout with filtering applied server-side.
 * Client-side filtering is only used for source label filtering (not sent to backend).
 */
const McpCatalogGalleryView: React.FC = () => {
  const { mcpServers, mcpSources, selectedSourceLabel, searchTerm, filters } = useMcpCatalog();

  const servers = React.useMemo(() => mcpServers?.items ?? [], [mcpServers?.items]);

  // Filter servers by selected source label (client-side)
  // This is the only client-side filtering needed - source label is UI-only
  const filteredServers = React.useMemo((): McpServer[] => {
    if (selectedSourceLabel === McpCategoryName.allServers) {
      return servers;
    }

    if (!mcpSources?.items) {
      return [];
    }

    // Handle "other" label (sources without labels)
    if (selectedSourceLabel === McpSourceLabel.other) {
      const sourcesWithoutLabels = mcpSources.items.filter(
        (source) => source.labels.length === 0 || source.labels.every((label) => !label.trim()),
      );
      const sourceIds = sourcesWithoutLabels.map((s) => s.id);
      return servers.filter((server) => server.source_id && sourceIds.includes(server.source_id));
    }

    // Find sources with the selected label
    const matchingSources = mcpSources.items.filter((source) =>
      source.labels.includes(selectedSourceLabel),
    );
    const matchingSourceIds = matchingSources.map((s) => s.id);

    return servers.filter(
      (server) => server.source_id && matchingSourceIds.includes(server.source_id),
    );
  }, [servers, mcpSources, selectedSourceLabel]);

  // Show empty state if no servers match filters
  if (filteredServers.length === 0) {
    const hasActiveFilters =
      searchTerm.length > 0 ||
      hasMcpFiltersActive(filters) ||
      selectedSourceLabel !== McpCategoryName.allServers;
    return (
      <EmptyMcpCatalogState
        testid="empty-mcp-catalog-state"
        title="No MCP servers found"
        headerIcon={SearchIcon}
        description={
          hasActiveFilters
            ? 'Adjust your search or filters and try again.'
            : 'No MCP servers are currently available.'
        }
      />
    );
  }

  return (
    <Grid hasGutter>
      {filteredServers.map((server) => (
        <GridItem key={`${server.name}/${server.source_id}`} sm={6} md={6} lg={6} xl={6} xl2={3}>
          <McpCatalogCard server={server} />
        </GridItem>
      ))}
    </Grid>
  );
};

export default McpCatalogGalleryView;
