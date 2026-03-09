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

  const { itemsByLabel, uncategorizedItems } = React.useMemo(() => {
    const knownLabels = new Set(sourceLabels);
    const byLabel = new Map<string, typeof items>();
    const uncategorized: typeof items = [];

    for (const item of items) {
      if (!item.source_id || !knownLabels.has(item.source_id)) {
        uncategorized.push(item);
      } else {
        const group = byLabel.get(item.source_id);
        if (group) {
          group.push(item);
        } else {
          byLabel.set(item.source_id, [item]);
        }
      }
    }

    return { itemsByLabel: byLabel, uncategorizedItems: uncategorized };
  }, [items, sourceLabels]);

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

  if (sourceLabels.length === 0) {
    return (
      <Stack hasGutter>
        <McpCatalogCategorySection title="Servers" servers={items} />
      </Stack>
    );
  }

  return (
    <Stack hasGutter>
      {sourceLabels.map((label) => (
        <McpCatalogCategorySection
          key={label}
          title={getCategoryDisplayName(label)}
          servers={itemsByLabel.get(label) ?? []}
        />
      ))}
      {uncategorizedItems.length > 0 && (
        <McpCatalogCategorySection key="other" title="Other" servers={uncategorizedItems} />
      )}
    </Stack>
  );
};

export default McpCatalogGalleryView;
