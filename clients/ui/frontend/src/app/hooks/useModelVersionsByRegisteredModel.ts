import * as React from 'react';
import { useFetchState, FetchState, FetchStateCallbackPromise, NotReadyError } from 'mod-arch-core';
import { ModelVersionList } from '~/app/types';
import { useModelRegistryAPI } from '~/app/hooks/useModelRegistryAPI';

const useModelVersionsByRegisteredModel = (
  registeredModelId?: string,
): FetchState<ModelVersionList> => {
  const { api, apiAvailable } = useModelRegistryAPI();

  const call = React.useCallback<FetchStateCallbackPromise<ModelVersionList>>(
    (opts) => {
      if (!apiAvailable) {
        return Promise.reject(new Error('API not yet available'));
      }
      if (!registeredModelId) {
        return Promise.reject(new NotReadyError('No model registeredModel id'));
      }

      return api.getModelVersionsByRegisteredModel(opts, registeredModelId);
    },
    [api, apiAvailable, registeredModelId],
  );

  return useFetchState(
    call,
    { items: [], size: 0, pageSize: 0, nextPageToken: '' },
    { initialPromisePurity: true },
  );
};

export default useModelVersionsByRegisteredModel;
