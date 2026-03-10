import '@testing-library/jest-dom';
import * as React from 'react';
import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { McpCatalogContext } from '~/app/context/mcpCatalog/McpCatalogContext';
import type { McpCatalogContextType } from '~/app/pages/mcpCatalog/types/mcpCatalogContext';
import McpCatalogCoreLoader from '~/app/pages/mcpCatalog/McpCatalogCoreLoader';

jest.mock('mod-arch-kubeflow', () => ({
  ...jest.requireActual('mod-arch-kubeflow'),
  useThemeContext: jest.fn(() => ({ isMUITheme: false })),
}));

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  Outlet: () => <div data-testid="mcp-catalog-outlet">Outlet</div>,
}));

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
  const value = { ...defaultContext, ...overrides };
  return render(
    <MemoryRouter>
      <McpCatalogContext.Provider value={value}>
        <McpCatalogCoreLoader />
      </McpCatalogContext.Provider>
    </MemoryRouter>,
  );
};

describe('McpCatalogCoreLoader', () => {
  it('shows error state when catalogSourcesLoadError is set', () => {
    renderWithContext({
      catalogSourcesLoadError: new Error('Failed to load sources'),
      catalogSourcesLoaded: true,
    });
    expect(screen.getByText('MCP catalog source load error')).toBeInTheDocument();
    expect(screen.getByText('Failed to load sources')).toBeInTheDocument();
    expect(screen.queryByTestId('mcp-catalog-outlet')).not.toBeInTheDocument();
  });

  it('shows loading state when catalog sources are not loaded', () => {
    renderWithContext({ catalogSourcesLoaded: false });
    expect(screen.getByText('Loading')).toBeInTheDocument();
    expect(screen.queryByTestId('mcp-catalog-outlet')).not.toBeInTheDocument();
  });

  it('shows empty state when sourceLabels is empty', () => {
    renderWithContext({ sourceLabels: [], catalogSourcesLoaded: true });
    expect(screen.getByTestId('empty-mcp-catalog-state')).toBeInTheDocument();
    expect(screen.getByText('MCP catalog configuration required')).toBeInTheDocument();
    expect(
      screen.getByText(
        'There are no MCP sources to display. Request that your administrator configure model sources for the catalog.',
      ),
    ).toBeInTheDocument();
    expect(screen.queryByTestId('mcp-catalog-outlet')).not.toBeInTheDocument();
  });

  it('shows MUI empty state title when isMUITheme is true', () => {
    const { useThemeContext } = jest.requireMock('mod-arch-kubeflow');
    useThemeContext.mockReturnValue({ isMUITheme: true });
    renderWithContext({ sourceLabels: [], catalogSourcesLoaded: true });
    expect(screen.getByText('Deploy a model catalog')).toBeInTheDocument();
    expect(
      screen.getByText(
        'To deploy model catalog and discover MCP servers, follow the instructions in the docs below.',
      ),
    ).toBeInTheDocument();
    useThemeContext.mockReturnValue({ isMUITheme: false });
  });

  it('renders Outlet when sourceLabels has items', () => {
    renderWithContext({ sourceLabels: ['community_mcp_servers'], catalogSourcesLoaded: true });
    expect(screen.getByTestId('mcp-catalog-outlet')).toBeInTheDocument();
    expect(screen.getByText('Outlet')).toBeInTheDocument();
    expect(screen.queryByTestId('empty-mcp-catalog-state')).not.toBeInTheDocument();
  });
});
