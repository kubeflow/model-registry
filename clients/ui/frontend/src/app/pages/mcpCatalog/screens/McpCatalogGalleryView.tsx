import * as React from 'react';
import {
  Button,
  EmptyState,
  EmptyStateActions,
  EmptyStateBody,
  EmptyStateFooter,
  Grid,
  GridItem,
  Skeleton,
  Stack,
} from '@patternfly/react-core';
import { ExclamationCircleIcon, SearchIcon } from '@patternfly/react-icons';
import { McpCatalogContext } from '~/app/context/mcpCatalog/McpCatalogContext';
import { SourceLabel } from '~/app/modelCatalogTypes';
import { filterEnabledCatalogSources } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import { MCP_CATALOG_GALLERY, MCP_CATALOG_GRID_SPAN } from '~/app/pages/mcpCatalog/const';
import McpCatalogLoadMoreSection from '~/app/pages/mcpCatalog/components/McpCatalogLoadMoreSection';
import McpCatalogCategorySection from '~/app/pages/mcpCatalog/screens/McpCatalogCategorySection';

const hasActiveSearchOrFilters = (
  searchQuery: string,
  filters: Record<string, string[] | undefined>,
): boolean =>
  searchQuery.trim() !== '' ||
  Object.values(filters).some((arr) => Array.isArray(arr) && arr.length > 0);

const McpCatalogGalleryView: React.FC = () => {
  const {
    mcpServers,
    mcpServersLoaded,
    mcpServersLoadError,
    refreshMcpServers,
    selectedSourceLabel,
    setSelectedSourceLabel,
    clearAllFilters,
    searchQuery,
    filters,
    catalogSources,
    sourceLabels,
    sourceLabelNames,
    hasNoLabelSources,
  } = React.useContext(McpCatalogContext);

  const getDisplayName = React.useCallback(
    (label: string): string =>
      label === SourceLabel.other ? 'Other MCP Servers' : sourceLabelNames[label] || label,
    [sourceLabelNames],
  );
  const isAllServersView = selectedSourceLabel === undefined;
  const { items } = mcpServers;
  const hasSearchOrFilters = React.useMemo(
    () => hasActiveSearchOrFilters(searchQuery, filters),
    [searchQuery, filters],
  );

  const [visibleCount, setVisibleCount] = React.useState<number>(MCP_CATALOG_GALLERY.PAGE_SIZE);

  React.useEffect(() => {
    if (selectedSourceLabel !== undefined || hasSearchOrFilters) {
      setVisibleCount(MCP_CATALOG_GALLERY.PAGE_SIZE);
    }
  }, [selectedSourceLabel, hasSearchOrFilters]);

  const { itemsByLabel, uncategorizedItems } = React.useMemo(() => {
    if (!isAllServersView) {
      return { itemsByLabel: new Map(), uncategorizedItems: [] };
    }
    const enabled = filterEnabledCatalogSources(catalogSources);
    const sourceIdToLabel = new Map<string, string>();
    if (enabled?.items) {
      for (const source of enabled.items) {
        const label =
          source.labels.length > 0
            ? source.labels.map((l) => l.trim()).find(Boolean)
            : SourceLabel.other;
        if (label) {
          sourceIdToLabel.set(source.id, label);
        }
      }
    }
    const knownLabels = new Set(sourceLabels);
    if (hasNoLabelSources) {
      knownLabels.add(SourceLabel.other);
    }
    const byLabel = new Map<string, typeof items>();
    const uncategorized: typeof items = [];

    for (const item of items) {
      const resolvedLabel = item.source_id ? sourceIdToLabel.get(item.source_id) : undefined;
      if (resolvedLabel && knownLabels.has(resolvedLabel)) {
        const group = byLabel.get(resolvedLabel);
        if (group) {
          group.push(item);
        } else {
          byLabel.set(resolvedLabel, [item]);
        }
      } else {
        uncategorized.push(item);
      }
    }

    return { itemsByLabel: byLabel, uncategorizedItems: uncategorized };
  }, [items, catalogSources, sourceLabels, hasNoLabelSources, isAllServersView]);

  if (mcpServersLoadError) {
    return (
      <EmptyState
        icon={ExclamationCircleIcon}
        headingLevel="h3"
        titleText="Unable to load MCP servers"
        data-testid="mcp-catalog-load-error"
      >
        <EmptyStateBody>{mcpServersLoadError.message}</EmptyStateBody>
        <EmptyStateFooter>
          <EmptyStateActions>
            <Button variant="link" data-testid="mcp-catalog-retry" onClick={refreshMcpServers}>
              Retry
            </Button>
          </EmptyStateActions>
        </EmptyStateFooter>
      </EmptyState>
    );
  }

  if (mcpServersLoaded && items.length === 0) {
    return (
      <EmptyState
        icon={SearchIcon}
        headingLevel="h3"
        titleText="No servers found"
        data-testid="mcp-catalog-empty-search"
      >
        <EmptyStateBody>Adjust your filters and try again.</EmptyStateBody>
        <EmptyStateFooter>
          <EmptyStateActions>
            <Button
              variant="link"
              data-testid="mcp-catalog-reset-filters"
              onClick={clearAllFilters}
            >
              Reset filters
            </Button>
          </EmptyStateActions>
        </EmptyStateFooter>
      </EmptyState>
    );
  }

  if (!mcpServersLoaded && items.length === 0) {
    return (
      <Grid hasGutter>
        {Array.from({
          length: MCP_CATALOG_GALLERY.CARDS_PER_ROW * 2,
        }).map((_, index) => (
          <GridItem
            key={index}
            sm={MCP_CATALOG_GRID_SPAN.sm}
            md={MCP_CATALOG_GRID_SPAN.md}
            lg={MCP_CATALOG_GRID_SPAN.lg}
            xl2={MCP_CATALOG_GRID_SPAN.xl2}
          >
            <Skeleton
              height="280px"
              width="100%"
              screenreaderText="Loading MCP servers"
              data-testid={`mcp-catalog-skeleton-${index}`}
            />
          </GridItem>
        ))}
      </Grid>
    );
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
          <McpCatalogLoadMoreSection
            onLoadMore={() =>
              setVisibleCount((prev) =>
                Math.min(prev + MCP_CATALOG_GALLERY.PAGE_SIZE, items.length),
              )
            }
          />
        )}
      </Stack>
    );
  }

  if (isAllServersView && hasSearchOrFilters) {
    const visibleItems = items.slice(0, visibleCount);
    const hasMore = items.length > MCP_CATALOG_GALLERY.PAGE_SIZE && items.length > visibleCount;

    return (
      <Stack hasGutter>
        <McpCatalogCategorySection
          title={MCP_CATALOG_GALLERY.SECTION_TITLE}
          servers={visibleItems}
        />
        {hasMore && (
          <McpCatalogLoadMoreSection
            onLoadMore={() =>
              setVisibleCount((prev) =>
                Math.min(prev + MCP_CATALOG_GALLERY.PAGE_SIZE, items.length),
              )
            }
          />
        )}
      </Stack>
    );
  }

  if (sourceLabels.length === 0) {
    return (
      <Stack hasGutter>
        <McpCatalogCategorySection title={MCP_CATALOG_GALLERY.SECTION_TITLE} servers={items} />
      </Stack>
    );
  }

  const otherSectionServers = hasNoLabelSources
    ? (itemsByLabel.get(SourceLabel.other) ?? []).concat(uncategorizedItems)
    : uncategorizedItems;
  const showOtherSection = otherSectionServers.length > 0;

  return (
    <Stack hasGutter>
      {sourceLabels.map((label) => (
        <McpCatalogCategorySection
          key={label}
          title={getDisplayName(label)}
          servers={itemsByLabel.get(label) ?? []}
          maxItems={isAllServersView ? MCP_CATALOG_GALLERY.CARDS_PER_ROW : undefined}
          onShowAll={isAllServersView ? () => setSelectedSourceLabel(label) : undefined}
        />
      ))}
      {showOtherSection && (
        <McpCatalogCategorySection
          key="other"
          title={hasNoLabelSources ? 'Other MCP Servers' : 'Other'}
          servers={otherSectionServers}
          maxItems={isAllServersView ? MCP_CATALOG_GALLERY.CARDS_PER_ROW : undefined}
          onShowAll={
            isAllServersView && hasNoLabelSources
              ? () => setSelectedSourceLabel(SourceLabel.other)
              : undefined
          }
        />
      )}
    </Stack>
  );
};

export default McpCatalogGalleryView;
