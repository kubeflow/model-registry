import { FetchState, FetchStateCallbackPromise, useFetchState } from 'mod-arch-core';
import React from 'react';
import { CatalogSourceList } from '~/app/modelCatalogTypes';
import { ModelCatalogAPIState } from './useModelCatalogAPIState';

export const useCatalogSources = (
  apiState: ModelCatalogAPIState,
): FetchState<CatalogSourceList> => {
  const call = React.useCallback<FetchStateCallbackPromise<CatalogSourceList>>(
    (opts) => {
      if (!apiState.apiAvailable) {
        return Promise.reject(new Error('API not yet available'));
      }

      return apiState.api.getListSources(opts);
    },
    [apiState],
  );
  return useFetchState(
    call,
    { items: [], size: 0, pageSize: 0, nextPageToken: '' },
    { initialPromisePurity: true },
  );
};
