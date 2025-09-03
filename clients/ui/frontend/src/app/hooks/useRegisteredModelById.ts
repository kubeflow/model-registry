import * as React from 'react';
import { useFetchState, FetchState, FetchStateCallbackPromise, NotReadyError } from 'mod-arch-core';
import { RegisteredModel } from '~/app/types';
import { useModelRegistryAPI } from '~/app/hooks/useModelRegistryAPI';

const useRegisteredModelById = (registeredModel?: string): FetchState<RegisteredModel | null> => {
  const { api, apiAvailable } = useModelRegistryAPI();

  const call = React.useCallback<FetchStateCallbackPromise<RegisteredModel | null>>(
    (opts) => {
      if (!apiAvailable) {
        return Promise.reject(new Error('API not yet available'));
      }
      if (!registeredModel) {
        return Promise.reject(new NotReadyError('No registered model id'));
      }

      return api.getRegisteredModel(opts, registeredModel);
    },
    [api, apiAvailable, registeredModel],
  );

  return useFetchState(call, null);
};

export default useRegisteredModelById;
