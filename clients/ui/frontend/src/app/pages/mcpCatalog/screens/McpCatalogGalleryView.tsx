import * as React from 'react';
import { EmptyState, EmptyStateBody, Stack } from '@patternfly/react-core';
import { SearchIcon } from '@patternfly/react-icons';
import { McpCatalogContext } from '~/app/context/mcpCatalog/McpCatalogContext';
import type { McpCatalogCategoryId } from '~/app/context/mcpCatalog/McpCatalogContext';
import type { McpServerMock } from '~/app/pages/mcpCatalog/types/mcpServer';
import { mockMcpServers } from '~/app/pages/mcpCatalog/mocks/mockMcpServers';
import McpCatalogCategorySection from './McpCatalogCategorySection';

type McpCatalogGalleryViewProps = {
  searchTerm: string;
};

const filterBySearch = (servers: McpServerMock[], query: string): McpServerMock[] => {
  if (!query.trim()) {
    return servers;
  }
  const q = query.trim().toLowerCase();
  return servers.filter(
    (s) => s.name.toLowerCase().includes(q) || s.description.toLowerCase().includes(q),
  );
};

const filterByCategory = (
  servers: McpServerMock[],
  category: McpCatalogCategoryId,
): McpServerMock[] => {
  if (category === 'all') {
    return servers;
  }
  return servers.filter((s) => s.category === category);
};

const McpCatalogGalleryView: React.FC<McpCatalogGalleryViewProps> = ({ searchTerm }) => {
  const { selectedCategory } = React.useContext(McpCatalogContext);

  const filteredServers = React.useMemo(() => {
    const bySearch = filterBySearch(mockMcpServers, searchTerm);
    return filterByCategory(bySearch, selectedCategory);
  }, [searchTerm, selectedCategory]);

  const sampleServers = React.useMemo(
    () => filteredServers.filter((s) => s.category === 'sample'),
    [filteredServers],
  );
  const otherServers = React.useMemo(
    () => filteredServers.filter((s) => s.category === 'other'),
    [filteredServers],
  );

  const showSample = selectedCategory === 'all' || selectedCategory === 'sample';
  const showOther = selectedCategory === 'all' || selectedCategory === 'other';

  if (filteredServers.length === 0) {
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

  return (
    <Stack hasGutter>
      {showSample && (
        <McpCatalogCategorySection
          title="Sample MCP servers"
          description="Sample of MCP cards category description"
          servers={sampleServers}
        />
      )}
      {showOther && (
        <McpCatalogCategorySection
          title="Other MCP servers"
          description="A broad collection of community and third-party MCP servers available for integration."
          servers={otherServers}
        />
      )}
    </Stack>
  );
};

export default McpCatalogGalleryView;
