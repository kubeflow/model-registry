import { FetchState, FetchStateCallbackPromise, NotReadyError, useFetchState } from 'mod-arch-core';
import React from 'react';
import { useModelCatalogAPI } from '../modelCatalog/useModelCatalogAPI';
import { McpServer } from '~/app/mcpServerCatalogTypes';

type State = McpServer | null;
export const useMcpServer = (serverId: string): FetchState<State> => {
  const { api, apiAvailable } = useModelCatalogAPI();

  const call = React.useCallback<FetchStateCallbackPromise<State>>(
    (opts) => {
      if (!apiAvailable) {
        return Promise.reject(new Error('API not yet available'));
      }
      if (!serverId) {
        return Promise.reject(new NotReadyError('No server id'));
      }
      return api.getMcpServer(opts, serverId);
    },
    [api, apiAvailable, serverId],
  );
  return useFetchState(call, null, {
    initialPromisePurity: true,
  });
};
