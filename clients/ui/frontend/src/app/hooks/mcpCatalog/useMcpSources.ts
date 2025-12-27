import { FetchState, FetchStateCallbackPromise, useFetchState } from 'mod-arch-core';
import React from 'react';
import { McpCatalogSourceList } from '~/app/pages/mcpCatalog/types';
import { McpCatalogAPIState } from './useMcpCatalogAPIState';

export const useMcpSources = (apiState: McpCatalogAPIState): FetchState<McpCatalogSourceList> => {
  const call = React.useCallback<FetchStateCallbackPromise<McpCatalogSourceList>>(
    (opts) => {
      if (!apiState.apiAvailable) {
        return Promise.reject(new Error('API not yet available'));
      }
      return apiState.api.getMcpSources(opts);
    },
    [apiState],
  );

  return useFetchState(
    call,
    { items: [], size: 0, pageSize: 0, nextPageToken: '' },
    { initialPromisePurity: true },
  );
};
