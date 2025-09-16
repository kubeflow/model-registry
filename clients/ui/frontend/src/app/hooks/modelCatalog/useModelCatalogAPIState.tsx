import { APIState, useAPIState } from 'mod-arch-core';
import React from 'react';
import {
  getCatalogModel,
  getCatalogModelsBySource,
  getListCatalogModelArtifacts,
  getListSources,
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
    }),
    [queryParameters],
  );

  return useAPIState(hostPath, createAPI);
};

export default useModelCatalogAPIState;
