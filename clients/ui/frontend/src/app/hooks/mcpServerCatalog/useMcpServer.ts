import { FetchState, FetchStateCallbackPromise, NotReadyError, useFetchState } from 'mod-arch-core';
import React from 'react';
import { McpServer } from '~/app/mcpServerCatalogTypes';
import { useModelCatalogAPI } from '~/app/hooks/modelCatalog/useModelCatalogAPI';
import type { ModelCatalogAPIState } from '~/app/hooks/modelCatalog/useModelCatalogAPIState';

type State = McpServer | null;

export const useMcpServerWithAPI = (
  apiState: ModelCatalogAPIState,
  serverId: string,
): FetchState<State> => {
  const { api, apiAvailable } = apiState;

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

export const useMcpServer = (serverId: string): FetchState<State> => {
  const { api, apiAvailable } = useModelCatalogAPI();
  return useMcpServerWithAPI({ api, apiAvailable }, serverId);
};
