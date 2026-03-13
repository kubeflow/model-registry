import * as React from 'react';
import { useSearchParams } from 'react-router-dom';
import { MCP_FILTER_KEYS } from '~/app/pages/mcpCatalog/const';
import type { McpCatalogFiltersState } from '~/app/pages/mcpCatalog/types/mcpCatalogFilterOptions';

const SEARCH_PARAM = 'q';
const SOURCE_PARAM = 'source';

function filtersFromParams(params: URLSearchParams): McpCatalogFiltersState {
  const filters: McpCatalogFiltersState = {};
  for (const key of MCP_FILTER_KEYS) {
    const raw = params.get(key);
    if (raw) {
      filters[key] = raw.split(',');
    }
  }
  return filters;
}

function filtersToParams(filters: McpCatalogFiltersState, params: URLSearchParams): void {
  for (const key of MCP_FILTER_KEYS) {
    const values = filters[key];
    if (values && values.length > 0) {
      params.set(key, values.join(','));
    } else {
      params.delete(key);
    }
  }
}

type UrlState = {
  searchQuery: string;
  filters: McpCatalogFiltersState;
  selectedSourceLabel: string | undefined;
};

type UseMcpUrlSyncReturn = {
  initialState: UrlState;
  syncToUrl: (state: UrlState) => void;
};

export function useMcpUrlSync(): UseMcpUrlSyncReturn {
  const [searchParams, setSearchParams] = useSearchParams();

  const initialState = React.useMemo<UrlState>(() => {
    const q = searchParams.get(SEARCH_PARAM) || '';
    const source = searchParams.get(SOURCE_PARAM) || undefined;
    const filters = filtersFromParams(searchParams);
    return { searchQuery: q, filters, selectedSourceLabel: source };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const syncToUrl = React.useCallback(
    (state: UrlState) => {
      setSearchParams(
        (prev) => {
          const next = new URLSearchParams(prev);
          if (state.searchQuery) {
            next.set(SEARCH_PARAM, state.searchQuery);
          } else {
            next.delete(SEARCH_PARAM);
          }
          if (state.selectedSourceLabel) {
            next.set(SOURCE_PARAM, state.selectedSourceLabel);
          } else {
            next.delete(SOURCE_PARAM);
          }
          filtersToParams(state.filters, next);
          return next;
        },
        { replace: true },
      );
    },
    [setSearchParams],
  );

  return { initialState, syncToUrl };
}
