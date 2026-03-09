import '@testing-library/jest-dom';
import * as React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { McpCatalogContext } from '~/app/context/mcpCatalog/McpCatalogContext';
import type { McpCatalogContextType } from '~/app/pages/mcpCatalog/types/mcpCatalogContext';
import type { McpServer } from '~/app/mcpServerCatalogTypes';
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
        <McpCatalogGalleryView searchTerm="" />
      </McpCatalogContext.Provider>
    </MemoryRouter>,
  );
};

describe('McpCatalogGalleryView', () => {
  it('renders 6 skeleton cards when loading', () => {
    renderWithContext({ mcpServersLoaded: false, mcpServers: { items: [] } });
    for (let i = 0; i < 6; i++) {
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

  it('renders Load more button when category has more than 6 items', () => {
    const servers = Array.from({ length: 8 }, (_, i) => buildServer(i + 1, 'cat-a'));
    renderWithContext({
      mcpServersLoaded: true,
      mcpServers: { items: servers },
      selectedSourceLabel: 'cat-a',
      sourceLabels: ['cat-a'],
      sourceLabelNames: { 'cat-a': 'Category A' },
    });
    expect(screen.getByTestId('mcp-load-more-button')).toBeInTheDocument();
    expect(screen.getByText('Load more MCP servers')).toBeInTheDocument();
    expect(screen.getAllByTestId(/^mcp-catalog-card-\d+$/)).toHaveLength(6);
  });

  it('Load more button reveals all items when clicked', () => {
    const servers = Array.from({ length: 8 }, (_, i) => buildServer(i + 1, 'cat-a'));
    renderWithContext({
      mcpServersLoaded: true,
      mcpServers: { items: servers },
      selectedSourceLabel: 'cat-a',
      sourceLabels: ['cat-a'],
      sourceLabelNames: { 'cat-a': 'Category A' },
    });
    fireEvent.click(screen.getByTestId('mcp-load-more-button'));
    expect(screen.getAllByTestId(/^mcp-catalog-card-\d+$/)).toHaveLength(8);
    expect(screen.queryByTestId('mcp-load-more-button')).not.toBeInTheDocument();
  });

  it('does not show Load more when category has 6 or fewer items', () => {
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

  it('renders Show all link in All Servers view with max 3 items per category', () => {
    const servers = [
      ...Array.from({ length: 5 }, (_, i) => buildServer(i + 1, 'cat-a')),
      ...Array.from({ length: 4 }, (_, i) => buildServer(10 + i, 'cat-b')),
    ];
    renderWithContext({
      mcpServersLoaded: true,
      mcpServers: { items: servers },
      selectedSourceLabel: undefined,
      sourceLabels: ['cat-a', 'cat-b'],
      sourceLabelNames: { 'cat-a': 'Category A', 'cat-b': 'Category B' },
    });
    expect(screen.getByText('Show all Category A')).toBeInTheDocument();
    expect(screen.getByText('Show all Category B')).toBeInTheDocument();
  });
});
