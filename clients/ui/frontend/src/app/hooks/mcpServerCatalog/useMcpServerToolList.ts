import { FetchState, FetchStateCallbackPromise, NotReadyError, useFetchState } from 'mod-arch-core';
import React from 'react';
import { McpToolList } from '~/app/mcpServerCatalogTypes';
import { McpCatalogContext } from '~/app/context/mcpCatalog/McpCatalogContext';

export const useMcpServerToolList = (
  serverId: string,
  serverName?: string,
): FetchState<McpToolList> => {
  const { mcpApiState } = React.useContext(McpCatalogContext);
  const { api, apiAvailable } = mcpApiState;
  const call = React.useCallback<FetchStateCallbackPromise<McpToolList>>(
    (opts) => {
      if (!apiAvailable) {
        return Promise.reject(new Error('API not yet available'));
      }
      if (!serverId) {
        return Promise.reject(new NotReadyError('No server id'));
      }
      return api.getMcpServerToolList(opts, serverId, serverName).then((data) => ({
        ...data,
        items: data.items ?? [],
      }));
    },
    [api, apiAvailable, serverId, serverName],
  );
  return useFetchState(
    call,
    { items: [], size: 0, pageSize: 0, nextPageToken: '' },
    { initialPromisePurity: true },
  );
};
