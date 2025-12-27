import { APIState, useAPIState } from 'mod-arch-core';
import React from 'react';
import {
  getMcpFilterOptions,
  getMcpServer,
  getMcpServers,
  getMcpSources,
} from '~/app/api/mcpCatalog/service';
import { McpCatalogAPIs } from '~/app/pages/mcpCatalog/types';

export type McpCatalogAPIState = APIState<McpCatalogAPIs>;

const useMcpCatalogAPIState = (
  hostPath: string | null,
  queryParameters?: Record<string, unknown>,
): [apiState: McpCatalogAPIState, refreshAPIState: () => void] => {
  const createAPI = React.useCallback(
    (path: string) => ({
      getMcpServers: getMcpServers(path, queryParameters),
      getMcpServer: getMcpServer(path, queryParameters),
      getMcpSources: getMcpSources(path, queryParameters),
      getMcpFilterOptions: getMcpFilterOptions(path, queryParameters),
    }),
    [queryParameters],
  );

  return useAPIState(hostPath, createAPI);
};

export default useMcpCatalogAPIState;
