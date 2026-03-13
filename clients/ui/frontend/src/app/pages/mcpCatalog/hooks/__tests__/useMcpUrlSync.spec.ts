import '@testing-library/jest-dom';
import { renderHook, act } from '@testing-library/react';
import * as React from 'react';
import { MemoryRouter } from 'react-router-dom';
import { useMcpUrlSync } from '~/app/pages/mcpCatalog/hooks/useMcpUrlSync';

function createWrapper(initialEntry = '/') {
  const Wrapper = ({ children }: { children: React.ReactNode }) =>
    React.createElement(MemoryRouter, { initialEntries: [initialEntry] }, children);
  Wrapper.displayName = 'TestRouterWrapper';
  return Wrapper;
}

describe('useMcpUrlSync', () => {
  it('returns empty initial state when URL has no params', () => {
    const { result } = renderHook(() => useMcpUrlSync(), {
      wrapper: createWrapper('/'),
    });
    expect(result.current.initialState).toEqual({
      searchQuery: '',
      filters: {},
      selectedSourceLabel: undefined,
    });
  });

  it('reads search query from URL', () => {
    const { result } = renderHook(() => useMcpUrlSync(), {
      wrapper: createWrapper('/?q=test-search'),
    });
    expect(result.current.initialState.searchQuery).toBe('test-search');
  });

  it('reads source label from URL', () => {
    const { result } = renderHook(() => useMcpUrlSync(), {
      wrapper: createWrapper('/?source=my-source'),
    });
    expect(result.current.initialState.selectedSourceLabel).toBe('my-source');
  });

  it('reads filter params from URL', () => {
    const { result } = renderHook(() => useMcpUrlSync(), {
      wrapper: createWrapper('/?deploymentMode=Local,Remote&license=MIT'),
    });
    expect(result.current.initialState.filters).toEqual({
      deploymentMode: ['Local', 'Remote'],
      license: ['MIT'],
    });
  });

  it('reads combined params from URL', () => {
    const { result } = renderHook(() => useMcpUrlSync(), {
      wrapper: createWrapper('/?q=hello&source=src1&labels=gpu,cpu'),
    });
    expect(result.current.initialState).toEqual({
      searchQuery: 'hello',
      selectedSourceLabel: 'src1',
      filters: { labels: ['gpu', 'cpu'] },
    });
  });

  it('syncToUrl writes state to URL params', () => {
    const { result } = renderHook(() => useMcpUrlSync(), {
      wrapper: createWrapper('/'),
    });
    act(() => {
      result.current.syncToUrl({
        searchQuery: 'my-query',
        selectedSourceLabel: 'source-1',
        filters: { deploymentMode: ['Local'] },
      });
    });
    expect(result.current.initialState.searchQuery).toBe('');
  });

  it('syncToUrl removes empty params', () => {
    const { result } = renderHook(() => useMcpUrlSync(), {
      wrapper: createWrapper('/?q=old&source=old-src&deploymentMode=Local'),
    });
    act(() => {
      result.current.syncToUrl({
        searchQuery: '',
        selectedSourceLabel: undefined,
        filters: {},
      });
    });
    expect(result.current.initialState.searchQuery).toBe('old');
  });
});
