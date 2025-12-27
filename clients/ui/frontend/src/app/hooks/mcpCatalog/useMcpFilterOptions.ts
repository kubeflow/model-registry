import { FetchState, FetchStateCallbackPromise, useFetchState } from 'mod-arch-core';
import React from 'react';
import { McpFilterOptionsList } from '~/app/pages/mcpCatalog/types';
import { McpCatalogAPIState } from './useMcpCatalogAPIState';

const EMPTY_FILTER_OPTIONS: McpFilterOptionsList = { filters: {} };

/**
 * Hook to fetch filter options for MCP servers.
 * Returns available values for each filterable field from the backend.
 */
export const useMcpFilterOptions = (
  apiState: McpCatalogAPIState,
): FetchState<McpFilterOptionsList> => {
  const call = React.useCallback<FetchStateCallbackPromise<McpFilterOptionsList>>(
    (opts) => {
      if (!apiState.apiAvailable) {
        return Promise.reject(new Error('API not yet available'));
      }
      return apiState.api.getMcpFilterOptions(opts);
    },
    [apiState],
  );

  return useFetchState(call, EMPTY_FILTER_OPTIONS, { initialPromisePurity: true });
};
