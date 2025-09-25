import { FetchState, FetchStateCallbackPromise, useFetchState } from 'mod-arch-core';
import React from 'react';
import { CatalogFilterOptionsList } from '~/app/modelCatalogTypes';
import { ModelCatalogAPIState } from './useModelCatalogAPIState';

type State = CatalogFilterOptionsList | null;

export const useCatalogFilterOptionList = (apiState: ModelCatalogAPIState): FetchState<State> => {
  const call = React.useCallback<FetchStateCallbackPromise<State>>(
    (opts) => {
      if (!apiState.apiAvailable) {
        return Promise.reject(new Error('API not yet available'));
      }

      return apiState.api.getCatalogFilterOptionList(opts);
    },
    [apiState],
  );
  return useFetchState(call, null, { initialPromisePurity: true });
};
