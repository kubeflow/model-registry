import * as React from 'react';
import { Button, EmptyState, EmptyStateBody, Flex, Stack, StackItem } from '@patternfly/react-core';
import { SearchIcon } from '@patternfly/react-icons';
import { McpCatalogContext } from '~/app/context/mcpCatalog/McpCatalogContext';
import McpCatalogCategorySection from '~/app/pages/mcpCatalog/screens/McpCatalogCategorySection';

const CATEGORY_PAGE_SIZE = 6;

type McpCatalogGalleryViewProps = {
  searchTerm: string;
};

const McpCatalogGalleryView: React.FC<McpCatalogGalleryViewProps> = () => {
  const {
    mcpServers,
    mcpServersLoaded,
    mcpServersLoadError,
    selectedSourceLabel,
    setSelectedSourceLabel,
    sourceLabels,
    sourceLabelNames,
  } = React.useContext(McpCatalogContext);

  const getDisplayName = React.useCallback(
    (label: string): string => sourceLabelNames[label] || label,
    [sourceLabelNames],
  );
  const isAllServersView = selectedSourceLabel === undefined;
  const { items } = mcpServers;

  const [visibleCount, setVisibleCount] = React.useState(CATEGORY_PAGE_SIZE);

  React.useEffect(() => {
    setVisibleCount(CATEGORY_PAGE_SIZE);
  }, [selectedSourceLabel]);

  const { itemsByLabel, uncategorizedItems } = React.useMemo(() => {
    if (!isAllServersView) {
      return { itemsByLabel: new Map(), uncategorizedItems: [] };
    }
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
  }, [items, sourceLabels, isAllServersView]);

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
    const visibleItems = items.slice(0, visibleCount);
    const hasMore = items.length > visibleCount;

    return (
      <Stack hasGutter>
        <McpCatalogCategorySection
          title={getDisplayName(selectedSourceLabel)}
          servers={visibleItems}
        />
        {hasMore && (
          <StackItem>
            <Flex justifyContent={{ default: 'justifyContentCenter' }}>
              <Button
                variant="secondary"
                data-testid="mcp-load-more-button"
                onClick={() => setVisibleCount(items.length)}
              >
                Load more MCP servers
              </Button>
            </Flex>
          </StackItem>
        )}
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
          title={getDisplayName(label)}
          servers={itemsByLabel.get(label) ?? []}
          maxItems={isAllServersView ? 3 : undefined}
          onShowAll={isAllServersView ? () => setSelectedSourceLabel(label) : undefined}
        />
      ))}
      {uncategorizedItems.length > 0 && (
        <McpCatalogCategorySection key="other" title="Other" servers={uncategorizedItems} />
      )}
    </Stack>
  );
};

export default McpCatalogGalleryView;
