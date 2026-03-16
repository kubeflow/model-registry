import '@testing-library/jest-dom';
import * as React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { McpCatalogContext } from '~/app/context/mcpCatalog/McpCatalogContext';
import type { McpCatalogContextType } from '~/app/pages/mcpCatalog/types/mcpCatalogContext';
import type { McpServer } from '~/app/mcpServerCatalogTypes';
import { MCP_CATALOG_GALLERY } from '~/app/pages/mcpCatalog/const';
import McpCatalogGalleryView from '~/app/pages/mcpCatalog/screens/McpCatalogGalleryView';

const buildServer = (id: number, sourceId: string): McpServer => ({
  id,
  name: `Server ${id}`,
  description: `Description ${id}`,
  deploymentMode: 'local',
  securityIndicators: {},
  source_id: sourceId, // eslint-disable-line camelcase
  toolCount: 0,
});

const defaultContext: McpCatalogContextType = {
  filters: {},
  setFilters: jest.fn(),
  searchQuery: '',
  setSearchQuery: jest.fn(),
  namedQuery: null,
  setNamedQuery: jest.fn(),
  pagination: { page: 1, pageSize: 10, totalItems: 0 },
  setPage: jest.fn(),
  setPageSize: jest.fn(),
  setTotalItems: jest.fn(),
  selectedSourceLabel: undefined,
  setSelectedSourceLabel: jest.fn(),
  clearAllFilters: jest.fn(),
  sourceLabels: [],
  sourceLabelNames: {},
  hasNoLabelSources: false,
  catalogSources: null,
  catalogSourcesLoaded: true,
  catalogSourcesLoadError: undefined,
  mcpServers: { items: [] },
  mcpServersLoaded: true,
  mcpServersLoadError: undefined,
  refreshMcpServers: jest.fn(),
  filterOptions: null,
  filterOptionsLoaded: true,
  filterOptionsLoadError: undefined,
};

const renderWithContext = (overrides: Partial<McpCatalogContextType> = {}) => {
  const ctx = { ...defaultContext, ...overrides };
  return render(
    <MemoryRouter>
      <McpCatalogContext.Provider value={ctx}>
        <McpCatalogGalleryView />
      </McpCatalogContext.Provider>
    </MemoryRouter>,
  );
};

describe('McpCatalogGalleryView', () => {
  it('renders skeleton cards when loading', () => {
    const skeletonCount = MCP_CATALOG_GALLERY.CARDS_PER_ROW * 2;
    expect(skeletonCount).toBe(8);
    renderWithContext({ mcpServersLoaded: false, mcpServers: { items: [] } });
    for (let i = 0; i < skeletonCount; i++) {
      expect(screen.getByTestId(`mcp-catalog-skeleton-${i}`)).toBeInTheDocument();
    }
  });

  it('renders error state with Retry button', () => {
    const refreshFn = jest.fn();
    renderWithContext({
      mcpServersLoadError: new Error('Network failure'),
      refreshMcpServers: refreshFn,
    });
    expect(screen.getByTestId('mcp-catalog-load-error')).toBeInTheDocument();
    expect(screen.getByText('Unable to load MCP servers')).toBeInTheDocument();
    expect(screen.getByText('Network failure')).toBeInTheDocument();
    const retryBtn = screen.getByTestId('mcp-catalog-retry');
    expect(retryBtn).toBeInTheDocument();
    fireEvent.click(retryBtn);
    expect(refreshFn).toHaveBeenCalledTimes(1);
  });

  it('renders empty state with Reset filters button', () => {
    const clearFn = jest.fn();
    renderWithContext({
      mcpServersLoaded: true,
      mcpServers: { items: [] },
      clearAllFilters: clearFn,
    });
    expect(screen.getByTestId('mcp-catalog-empty-search')).toBeInTheDocument();
    expect(screen.getByText('No servers found')).toBeInTheDocument();
    expect(screen.getByText('Adjust your filters and try again.')).toBeInTheDocument();
    const resetBtn = screen.getByTestId('mcp-catalog-reset-filters');
    expect(resetBtn).toBeInTheDocument();
    fireEvent.click(resetBtn);
    expect(clearFn).toHaveBeenCalledTimes(1);
  });

  it('renders Load more button when category has more than PAGE_SIZE items', () => {
    const totalServers = MCP_CATALOG_GALLERY.PAGE_SIZE + 2;
    const servers = Array.from({ length: totalServers }, (_, i) => buildServer(i + 1, 'cat-a'));
    renderWithContext({
      mcpServersLoaded: true,
      mcpServers: { items: servers },
      selectedSourceLabel: 'cat-a',
      sourceLabels: ['cat-a'],
      sourceLabelNames: { 'cat-a': 'Category A' },
    });
    expect(screen.getByTestId('mcp-load-more-button')).toBeInTheDocument();
    expect(screen.getByText('Load more MCP servers')).toBeInTheDocument();
    expect(screen.getAllByTestId(/^mcp-catalog-card-\d+$/)).toHaveLength(10);
  });

  it('Load more button pages through batches and reveals all items after multiple clicks', () => {
    const totalServers = MCP_CATALOG_GALLERY.PAGE_SIZE * 2 + 1;
    const servers = Array.from({ length: totalServers }, (_, i) => buildServer(i + 1, 'cat-a'));
    renderWithContext({
      mcpServersLoaded: true,
      mcpServers: { items: servers },
      selectedSourceLabel: 'cat-a',
      sourceLabels: ['cat-a'],
      sourceLabelNames: { 'cat-a': 'Category A' },
    });
    expect(screen.getAllByTestId(/^mcp-catalog-card-\d+$/)).toHaveLength(10);
    fireEvent.click(screen.getByTestId('mcp-load-more-button'));
    expect(screen.getAllByTestId(/^mcp-catalog-card-\d+$/)).toHaveLength(20);
    expect(screen.getByTestId('mcp-load-more-button')).toBeInTheDocument();
    fireEvent.click(screen.getByTestId('mcp-load-more-button'));
    expect(screen.getAllByTestId(/^mcp-catalog-card-\d+$/)).toHaveLength(21);
    expect(screen.queryByTestId('mcp-load-more-button')).not.toBeInTheDocument();
  });

  it('does not show Load more when category has PAGE_SIZE or fewer items', () => {
    const servers = Array.from({ length: 5 }, (_, i) => buildServer(i + 1, 'cat-a'));
    renderWithContext({
      mcpServersLoaded: true,
      mcpServers: { items: servers },
      selectedSourceLabel: 'cat-a',
      sourceLabels: ['cat-a'],
      sourceLabelNames: { 'cat-a': 'Category A' },
    });
    expect(screen.queryByTestId('mcp-load-more-button')).not.toBeInTheDocument();
    expect(screen.getAllByTestId(/^mcp-catalog-card-\d+$/)).toHaveLength(5);
  });

  it('renders Show all link in All Servers view with max CARDS_PER_ROW items per category', () => {
    const minItemsForShowAll = MCP_CATALOG_GALLERY.CARDS_PER_ROW + 1;
    expect(minItemsForShowAll).toBe(5);
    const servers = [
      ...Array.from({ length: minItemsForShowAll }, (_, i) => buildServer(i + 1, 'cat-a')),
      ...Array.from({ length: minItemsForShowAll }, (_, i) => buildServer(10 + i, 'cat-b')),
    ];
    const catalogSources = {
      items: [
        {
          id: 'cat-a',
          name: 'Category A',
          labels: ['cat-a'],
          enabled: true,
          status: 'available' as const,
        },
        {
          id: 'cat-b',
          name: 'Category B',
          labels: ['cat-b'],
          enabled: true,
          status: 'available' as const,
        },
      ],
      size: 2,
      pageSize: 10,
      nextPageToken: '',
    };
    renderWithContext({
      mcpServersLoaded: true,
      mcpServers: { items: servers },
      selectedSourceLabel: undefined,
      catalogSources,
      sourceLabels: ['cat-a', 'cat-b'],
      sourceLabelNames: { 'cat-a': 'Category A', 'cat-b': 'Category B' },
    });
    expect(screen.getByText('Show all Category A')).toBeInTheDocument();
    expect(screen.getByText('Show all Category B')).toBeInTheDocument();
  });
});
