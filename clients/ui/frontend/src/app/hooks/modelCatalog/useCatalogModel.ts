import { FetchState, FetchStateCallbackPromise, NotReadyError, useFetchState } from 'mod-arch-core';
import React from 'react';
import { CatalogModel } from '~/app/modelCatalogTypes';
import { useModelCatalogAPI } from './useModelCatalogAPI';

type State = CatalogModel | null;

export const useCatalogModel = (sourceId: string, modelName: string): FetchState<State> => {
  const { api, apiAvailable } = useModelCatalogAPI();

  const call = React.useCallback<FetchStateCallbackPromise<State>>(
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
      return api.getCatalogModel(opts, sourceId, modelName);
    },
    [api, apiAvailable, sourceId, modelName],
  );
  return useFetchState(call, null, {
    initialPromisePurity: true,
  });
};
