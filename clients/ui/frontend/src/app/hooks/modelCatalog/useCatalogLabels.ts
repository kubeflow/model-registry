import { FetchState, FetchStateCallbackPromise, useFetchState } from 'mod-arch-core';
import React from 'react';
import { CatalogLabelList } from '~/app/modelCatalogTypes';
import { ModelCatalogAPIState } from './useModelCatalogAPIState';

export const useCatalogLabels = (apiState: ModelCatalogAPIState): FetchState<CatalogLabelList> => {
  const call = React.useCallback<FetchStateCallbackPromise<CatalogLabelList>>(
    (opts) => {
      if (!apiState.apiAvailable) {
        return Promise.reject(new Error('API not yet available'));
      }
      return apiState.api.getCatalogLabels(opts);
    },
    [apiState],
  );
  return useFetchState(
    call,
    { items: [], size: 0, pageSize: 0, nextPageToken: '' },
    { initialPromisePurity: true },
  );
};
