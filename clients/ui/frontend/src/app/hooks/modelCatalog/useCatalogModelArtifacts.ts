import { FetchState, FetchStateCallbackPromise, NotReadyError, useFetchState } from 'mod-arch-core';
import React from 'react';
import { CatalogArtifactList } from '~/app/modelCatalogTypes';
import { useModelCatalogAPI } from './useModelCatalogAPI';

export const useCatalogModelArtifacts = (
  sourceId: string,
  modelName: string,
): FetchState<CatalogArtifactList> => {
  const { api, apiAvailable } = useModelCatalogAPI();

  const call = React.useCallback<FetchStateCallbackPromise<CatalogArtifactList>>(
    (opts) => {
      if (!apiAvailable) {
        return Promise.reject(new Error('API not yet available'));
      }
      if (!sourceId) {
        return Promise.reject(new NotReadyError('No source id'));
      }
      if (!modelName) {
        return Promise.reject(new NotReadyError('No model name'));
      }
      return api.getListCatalogModelArtifacts(opts, sourceId, modelName);
    },
    [api, apiAvailable, sourceId, modelName],
  );
  return useFetchState(
    call,
    { items: [], size: 0, pageSize: 0, nextPageToken: '' },
    {
      initialPromisePurity: true,
    },
  );
};
