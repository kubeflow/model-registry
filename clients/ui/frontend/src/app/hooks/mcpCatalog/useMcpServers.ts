import { FetchState, FetchStateCallbackPromise, useFetchState } from 'mod-arch-core';
import React from 'react';
import { McpServerList } from '~/app/pages/mcpCatalog/types';
import { McpCatalogAPIState } from './useMcpCatalogAPIState';

// Default page size for fetching MCP servers - set high to get all servers in one request
const DEFAULT_PAGE_SIZE = 100;

export type UseMcpServersOptions = {
  filterQuery?: string;
  searchTerm?: string;
};

export const useMcpServers = (
  apiState: McpCatalogAPIState,
  options?: UseMcpServersOptions,
): FetchState<McpServerList> => {
  const { filterQuery, searchTerm } = options ?? {};

  const call = React.useCallback<FetchStateCallbackPromise<McpServerList>>(
    (opts) => {
      if (!apiState.apiAvailable) {
        return Promise.reject(new Error('API not yet available'));
      }
      return apiState.api.getMcpServers(
        opts,
        undefined,
        DEFAULT_PAGE_SIZE,
        filterQuery,
        searchTerm,
      );
    },
    [apiState, filterQuery, searchTerm],
  );

  return useFetchState(
    call,
    { items: [], size: 0, pageSize: 0, nextPageToken: '' },
    { initialPromisePurity: true },
  );
};
