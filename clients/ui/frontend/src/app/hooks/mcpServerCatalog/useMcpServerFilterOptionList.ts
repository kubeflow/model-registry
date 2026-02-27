import { FetchState, FetchStateCallbackPromise, useFetchState } from 'mod-arch-core';
import React from 'react';
import { CatalogFilterOptionsList } from '~/app/modelCatalogTypes';
import { useModelCatalogAPI } from '~/app/hooks/modelCatalog/useModelCatalogAPI';

type State = CatalogFilterOptionsList | null;

export const useMcpServerFilterOptionList = (): FetchState<State> => {
  const { api, apiAvailable } = useModelCatalogAPI();
  const call = React.useCallback<FetchStateCallbackPromise<State>>(
    (opts) => {
      if (!apiAvailable) {
        return Promise.reject(new Error('API not yet available'));
      }

      return api.getMcpServerFilterOptionList(opts);
    },
    [api, apiAvailable],
  );
  return useFetchState(call, null, { initialPromisePurity: true });
};
