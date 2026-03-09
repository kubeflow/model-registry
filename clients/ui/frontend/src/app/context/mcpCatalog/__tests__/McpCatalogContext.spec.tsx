import '@testing-library/jest-dom';
import * as React from 'react';
import { renderHook, act } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import {
  McpCatalogContextProvider,
  McpCatalogContext,
} from '~/app/context/mcpCatalog/McpCatalogContext';

jest.mock('mod-arch-core', () => ({
  useQueryParamNamespaces: jest.fn(() => ({})),
  asEnumMember: jest.fn((val: unknown) => val),
  DeploymentMode: {},
}));

jest.mock('~/app/utilities/const', () => ({
  BFF_API_VERSION: 'v1',
  URL_PREFIX: '/model-registry',
}));

jest.mock('~/app/hooks/modelCatalog/useModelCatalogAPIState', () => ({
  __esModule: true,
  default: jest.fn(() => [
    {
      apiAvailable: false,
      api: {
        getMcpServerList: jest.fn(),
        getMcpServerFilterOptionList: jest.fn(),
      },
    },
  ]),
}));

jest.mock('~/app/hooks/modelCatalog/useCatalogSources', () => ({
  useCatalogSources: jest.fn(() => [
    { items: [], size: 0, pageSize: 0, nextPageToken: '' },
    true,
    undefined,
  ]),
}));

jest.mock('~/app/hooks/mcpServerCatalog/useMcpServersBySourceLabel', () => ({
  useMcpServersBySourceLabelWithAPI: jest.fn(() => ({
    mcpServers: { items: [] },
    mcpServersLoaded: true,
    mcpServersLoadError: undefined,
  })),
}));

jest.mock('~/app/hooks/mcpServerCatalog/useMcpServerFilterOptionList', () => ({
  useMcpServerFilterOptionListWithAPI: jest.fn(() => [null, true, undefined]),
}));

describe('McpCatalogContext', () => {
  const wrapper = ({ children }: { children: React.ReactNode }) => (
    <MemoryRouter>
      <McpCatalogContextProvider>{children}</McpCatalogContextProvider>
    </MemoryRouter>
  );

  it('provides default filter state', () => {
    const { result } = renderHook(() => React.useContext(McpCatalogContext), { wrapper });
    expect(result.current.filters).toEqual({});
    expect(result.current.searchQuery).toBe('');
    expect(result.current.namedQuery).toBeNull();
    expect(result.current.selectedSourceLabel).toBeUndefined();
    expect(result.current.pagination).toEqual({
      page: 1,
      pageSize: 10,
      totalItems: 0,
    });
  });

  it('updates searchQuery via setSearchQuery', () => {
    const { result } = renderHook(() => React.useContext(McpCatalogContext), { wrapper });
    act(() => {
      result.current.setSearchQuery('test');
    });
    expect(result.current.searchQuery).toBe('test');
  });

  it('updates namedQuery via setNamedQuery', () => {
    const { result } = renderHook(() => React.useContext(McpCatalogContext), { wrapper });
    act(() => {
      result.current.setNamedQuery('named');
    });
    expect(result.current.namedQuery).toBe('named');
    act(() => {
      result.current.setNamedQuery(null);
    });
    expect(result.current.namedQuery).toBeNull();
  });

  it('updates filters via setFilters', () => {
    const { result } = renderHook(() => React.useContext(McpCatalogContext), { wrapper });
    act(() => {
      result.current.setFilters({ deploymentMode: ['Local'] });
    });
    expect(result.current.filters).toEqual({ deploymentMode: ['Local'] });
    act(() => {
      result.current.setFilters((prev) => ({ ...prev, license: ['MIT'] }));
    });
    expect(result.current.filters).toEqual({ deploymentMode: ['Local'], license: ['MIT'] });
  });

  it('updates pagination via setPage, setPageSize, setTotalItems', () => {
    const { result } = renderHook(() => React.useContext(McpCatalogContext), { wrapper });
    act(() => {
      result.current.setPage(2);
    });
    expect(result.current.pagination.page).toBe(2);
    act(() => {
      result.current.setPageSize(20);
    });
    expect(result.current.pagination.pageSize).toBe(20);
    expect(result.current.pagination.page).toBe(1);
    act(() => {
      result.current.setTotalItems(50);
    });
    expect(result.current.pagination.totalItems).toBe(50);
  });

  it('updates selectedSourceLabel via setSelectedSourceLabel', () => {
    const { result } = renderHook(() => React.useContext(McpCatalogContext), { wrapper });
    act(() => {
      result.current.setSelectedSourceLabel('sample');
    });
    expect(result.current.selectedSourceLabel).toBe('sample');
    act(() => {
      result.current.setSelectedSourceLabel('other');
    });
    expect(result.current.selectedSourceLabel).toBe('other');
    act(() => {
      result.current.setSelectedSourceLabel(undefined);
    });
    expect(result.current.selectedSourceLabel).toBeUndefined();
  });

  it('clearAllFilters resets searchQuery, filters, selectedSourceLabel and namedQuery', () => {
    const { result } = renderHook(() => React.useContext(McpCatalogContext), { wrapper });
    act(() => {
      result.current.setSearchQuery('q');
      result.current.setFilters({ deploymentMode: ['Local'] });
      result.current.setSelectedSourceLabel('sample');
      result.current.setNamedQuery('named');
    });
    expect(result.current.searchQuery).toBe('q');
    expect(result.current.filters).toEqual({ deploymentMode: ['Local'] });
    expect(result.current.selectedSourceLabel).toBe('sample');
    expect(result.current.namedQuery).toBe('named');
    act(() => {
      result.current.clearAllFilters();
    });
    expect(result.current.searchQuery).toBe('');
    expect(result.current.filters).toEqual({});
    expect(result.current.selectedSourceLabel).toBeUndefined();
    expect(result.current.namedQuery).toBeNull();
  });
});
