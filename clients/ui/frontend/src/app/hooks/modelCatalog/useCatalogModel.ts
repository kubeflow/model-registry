import { FetchState, FetchStateCallbackPromise, NotReadyError, useFetchState } from 'mod-arch-core';
import { useModelCatalogAPI } from './useModelCatalogAPI';
import { CatalogModel } from '~/app/modelCatalogTypes';
import React from 'react';

type state = CatalogModel | null;

export const useCatalogModel = (sourceId: string, modelName: string): FetchState<state> | null => {
  const { api, apiAvailable } = useModelCatalogAPI();

  const call = React.useCallback<FetchStateCallbackPromise<state>>(
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
