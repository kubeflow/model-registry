import { FetchState, FetchStateCallbackPromise, NotReadyError, useFetchState } from 'mod-arch-core';
import React from 'react';
import { McpToolList } from '~/app/mcpServerCatalogTypes';
import { useModelCatalogAPI } from '~/app/hooks/modelCatalog/useModelCatalogAPI';

export const useMcpServerToolList = (serverId: string): FetchState<McpToolList> => {
  const { api, apiAvailable } = useModelCatalogAPI();
  const call = React.useCallback<FetchStateCallbackPromise<McpToolList>>(
    (opts) => {
      if (!apiAvailable) {
        return Promise.reject(new Error('API not yet available'));
      }
      if (!serverId) {
        return Promise.reject(new NotReadyError('No server id'));
      }
      return api.getMcpServerToolList(opts, serverId).then((data) => ({
        ...data,
        items: data.items ?? [],
      }));
    },
    [api, apiAvailable, serverId],
  );
  return useFetchState(
    call,
    { items: [], size: 0, pageSize: 0, nextPageToken: '' },
    { initialPromisePurity: true },
  );
};
