import { FetchState, FetchStateCallbackPromise, NotReadyError, useFetchState } from 'mod-arch-core';
import React from 'react';
import { CatalogModelList } from '~/app/modelCatalogTypes';
import { useModelCatalogAPI } from './useModelCatalogAPI';

export const useCatalogModelsbySources = (sourceId: string): FetchState<CatalogModelList> => {
  const { api, apiAvailable } = useModelCatalogAPI();

  const call = React.useCallback<FetchStateCallbackPromise<CatalogModelList>>(
    (opts) => {
      if (!apiAvailable) {
        return Promise.reject(new Error('API not yet available'));
      }
      if (!sourceId) {
        return Promise.reject(new NotReadyError('No source id'));
      }
      return api.getCatalogModelsBySource(opts, sourceId);
    },
    [api, apiAvailable, sourceId],
  );
  return useFetchState(
    call,
    { items: [], size: 0, pageSize: 0, nextPageToken: '' },
    {
      initialPromisePurity: true,
    },
  );
};
