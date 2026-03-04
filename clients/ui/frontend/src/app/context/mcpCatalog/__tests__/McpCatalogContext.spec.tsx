import '@testing-library/jest-dom';
import React from 'react';
import { renderHook, act } from '@testing-library/react';
import {
  McpCatalogContextProvider,
  McpCatalogContext,
} from '~/app/context/mcpCatalog/McpCatalogContext';

describe('McpCatalogContext', () => {
  const wrapper = ({ children }: { children: React.ReactNode }) => (
    <McpCatalogContextProvider>{children}</McpCatalogContextProvider>
  );

  it('provides default filter state', () => {
    const { result } = renderHook(() => React.useContext(McpCatalogContext), { wrapper });
    expect(result.current.filters).toEqual({});
    expect(result.current.searchQuery).toBe('');
    expect(result.current.namedQuery).toBeNull();
    expect(result.current.selectedCategory).toBe('all');
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

  it('updates selectedCategory via setSelectedCategory', () => {
    const { result } = renderHook(() => React.useContext(McpCatalogContext), { wrapper });
    act(() => {
      result.current.setSelectedCategory('sample');
    });
    expect(result.current.selectedCategory).toBe('sample');
    act(() => {
      result.current.setSelectedCategory('other');
    });
    expect(result.current.selectedCategory).toBe('other');
  });

  it('clearAllFilters resets searchQuery and filters', () => {
    const { result } = renderHook(() => React.useContext(McpCatalogContext), { wrapper });
    act(() => {
      result.current.setSearchQuery('q');
      result.current.setFilters({ deploymentMode: ['Local'] });
    });
    expect(result.current.searchQuery).toBe('q');
    expect(result.current.filters).toEqual({ deploymentMode: ['Local'] });
    act(() => {
      result.current.clearAllFilters();
    });
    expect(result.current.searchQuery).toBe('');
    expect(result.current.filters).toEqual({});
  });
});
