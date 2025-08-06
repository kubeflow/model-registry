import * as React from 'react';
import { useFetchState, FetchState, FetchStateCallbackPromise, NotReadyError } from 'mod-arch-core';
import { ModelVersion } from '~/app/types';
import { useModelRegistryAPI } from '~/app/hooks/useModelRegistryAPI';

const useModelVersionById = (modelVersionId?: string): FetchState<ModelVersion | null> => {
  const { api, apiAvailable } = useModelRegistryAPI();

  const call = React.useCallback<FetchStateCallbackPromise<ModelVersion | null>>(
    (opts) => {
      if (!apiAvailable) {
        return Promise.reject(new Error('API not yet available'));
      }
      if (!modelVersionId) {
        return Promise.reject(new NotReadyError('No model version id'));
      }

      return api.getModelVersion(opts, modelVersionId);
    },
    [api, apiAvailable, modelVersionId],
  );

  return useFetchState(call, null);
};

export default useModelVersionById;
