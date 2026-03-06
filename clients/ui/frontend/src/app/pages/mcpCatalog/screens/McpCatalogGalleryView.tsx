import * as React from 'react';
import { EmptyState, EmptyStateBody, Stack } from '@patternfly/react-core';
import { SearchIcon } from '@patternfly/react-icons';
import { McpCatalogContext } from '~/app/context/mcpCatalog/McpCatalogContext';
import McpCatalogCategorySection from '~/app/pages/mcpCatalog/screens/McpCatalogCategorySection';

function getCategoryDisplayName(sourceLabel: string): string {
  if (!sourceLabel) {
    return 'Other';
  }
  return sourceLabel.charAt(0).toUpperCase() + sourceLabel.slice(1).toLowerCase();
}

type McpCatalogGalleryViewProps = {
  searchTerm: string;
};

const McpCatalogGalleryView: React.FC<McpCatalogGalleryViewProps> = () => {
  const { mcpServers, mcpServersLoaded, mcpServersLoadError, selectedSourceLabel, sourceLabels } =
    React.useContext(McpCatalogContext);
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

  if (selectedSourceLabel !== undefined) {
    return (
      <Stack hasGutter>
        <McpCatalogCategorySection
          title={getCategoryDisplayName(selectedSourceLabel)}
          servers={items}
        />
      </Stack>
    );
  }

  const knownLabels = new Set(sourceLabels);
  const uncategorized = items.filter((s) => !s.source_id || !knownLabels.has(s.source_id));
  const hasUncategorized = uncategorized.length > 0;

  if (sourceLabels.length === 0) {
    return (
      <Stack hasGutter>
        <McpCatalogCategorySection title="Servers" servers={items} />
      </Stack>
    );
  }

  return (
    <Stack hasGutter>
      {sourceLabels.map((label) => {
        const sectionItems = items.filter((s) => s.source_id === label);
        return (
          <McpCatalogCategorySection
            key={label}
            title={getCategoryDisplayName(label)}
            servers={sectionItems}
          />
        );
      })}
      {hasUncategorized && (
        <McpCatalogCategorySection key="other" title="Other" servers={uncategorized} />
      )}
    </Stack>
  );
};

export default McpCatalogGalleryView;
