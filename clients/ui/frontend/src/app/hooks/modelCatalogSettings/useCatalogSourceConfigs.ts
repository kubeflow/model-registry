import { FetchState, FetchStateCallbackPromise, useFetchState } from 'mod-arch-core';
import React from 'react';
import { CatalogSourceConfigList } from '~/app/modelCatalogTypes';
import { ModelCatalogSettingsAPIState } from './useModelCatalogSettingsAPIState';

export const useCatalogSourceConfigs = (
  apiState: ModelCatalogSettingsAPIState,
): FetchState<CatalogSourceConfigList> => {
  const call = React.useCallback<FetchStateCallbackPromise<CatalogSourceConfigList>>(
    (opts) => {
      if (!apiState.apiAvailable) {
        return Promise.reject(new Error('API not yet available'));
      }

      return apiState.api.getCatalogSourceConfigs(opts);
    },
    [apiState],
  );
  return useFetchState(call, { catalogs: [] }, { initialPromisePurity: true });
};
