import '@testing-library/jest-dom';
import { renderHook, act } from '@testing-library/react';
import * as React from 'react';
import { MemoryRouter, useSearchParams } from 'react-router-dom';
import { useMcpUrlSync } from '~/app/pages/mcpCatalog/hooks/useMcpUrlSync';

const searchParamsRef = { current: new URLSearchParams() };

function SearchParamsCapture({ children }: { children: React.ReactNode }) {
  const [params] = useSearchParams();
  searchParamsRef.current = params;
  return React.createElement(React.Fragment, null, children);
}

function createWrapper(initialEntry = '/') {
  const Wrapper = ({ children }: { children: React.ReactNode }) =>
    React.createElement(
      MemoryRouter,
      { initialEntries: [initialEntry] },
      React.createElement(SearchParamsCapture, null, children),
    );
  Wrapper.displayName = 'TestRouterWrapper';
  return Wrapper;
}

function getSearchParams(): URLSearchParams {
  return searchParamsRef.current;
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
    const params = getSearchParams();
    expect(params.get('q')).toBe('my-query');
    expect(params.get('source')).toBe('source-1');
    expect(params.get('deploymentMode')).toBe('Local');
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
    const params = getSearchParams();
    expect(params.get('q')).toBeNull();
    expect(params.get('source')).toBeNull();
    expect(params.get('deploymentMode')).toBeNull();
  });
});
