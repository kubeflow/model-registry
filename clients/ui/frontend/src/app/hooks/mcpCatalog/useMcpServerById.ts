import { FetchState, FetchStateCallbackPromise, useFetchState } from 'mod-arch-core';
import React from 'react';
import { McpServer } from '~/app/pages/mcpCatalog/types';
import { McpCatalogAPIState } from './useMcpCatalogAPIState';

const DEFAULT_MCP_SERVER: McpServer = {
  id: '',
  name: '',
  description: '',
};

export const useMcpServerById = (
  apiState: McpCatalogAPIState,
  serverId: string | undefined,
): FetchState<McpServer> => {
  const call = React.useCallback<FetchStateCallbackPromise<McpServer>>(
    (opts) => {
      if (!apiState.apiAvailable) {
        return Promise.reject(new Error('API not yet available'));
      }
      if (!serverId) {
        return Promise.reject(new Error('Server ID is required'));
      }
      return apiState.api.getMcpServer(opts, serverId);
    },
    [apiState, serverId],
  );

  return useFetchState(call, DEFAULT_MCP_SERVER, { initialPromisePurity: true });
};
