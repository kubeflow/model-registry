import { APIState, useAPIState } from 'mod-arch-core';
import React from 'react';
import {
  getMcpServerFilterOptionList,
  getMcpServer,
  getMcpServerList,
  getMcpServerToolList,
} from '~/app/api/mcpServerCatalog/service';
import {
  getCatalogFilterOptionList,
  getCatalogLabels,
  getCatalogModel,
  getCatalogModelsBySource,
  getListCatalogModelArtifacts,
  getListSources,
  getPerformanceArtifacts,
} from '~/app/api/modelCatalog/service';
import { ModelCatalogAPIs } from '~/app/modelCatalogTypes';

export type ModelCatalogAPIState = APIState<ModelCatalogAPIs>;

const useModelCatalogAPIState = (
  hostPath: string | null,
  queryParameters?: Record<string, unknown>,
): [apiState: ModelCatalogAPIState, refreshAPIState: () => void] => {
  const createAPI = React.useCallback(
    (path: string) => ({
      getCatalogModelsBySource: getCatalogModelsBySource(path, queryParameters),
      getListSources: getListSources(path, queryParameters),
      getCatalogModel: getCatalogModel(path, queryParameters),
      getListCatalogModelArtifacts: getListCatalogModelArtifacts(path, queryParameters),
      getCatalogFilterOptionList: getCatalogFilterOptionList(path, queryParameters),
      getPerformanceArtifacts: getPerformanceArtifacts(path, queryParameters),
      getCatalogLabels: getCatalogLabels(path, queryParameters),
      getMcpServerList: getMcpServerList(path, queryParameters),
      getMcpServerFilterOptionList: getMcpServerFilterOptionList(path, queryParameters),
      getMcpServer: getMcpServer(path, queryParameters),
      getMcpServerToolList: getMcpServerToolList(path, queryParameters),
    }),
    [queryParameters],
  );

  return useAPIState(hostPath, createAPI);
};

export default useModelCatalogAPIState;
