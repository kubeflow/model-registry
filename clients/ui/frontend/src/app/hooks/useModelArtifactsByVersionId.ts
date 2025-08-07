import * as React from 'react';
import { useFetchState, FetchState, FetchStateCallbackPromise, NotReadyError } from 'mod-arch-core';
import { ModelArtifactList } from '~/app/types';
import { useModelRegistryAPI } from '~/app/hooks/useModelRegistryAPI';

const useModelArtifactsByVersionId = (modelVersionId?: string): FetchState<ModelArtifactList> => {
  const { api, apiAvailable } = useModelRegistryAPI();
  const callback = React.useCallback<FetchStateCallbackPromise<ModelArtifactList>>(
    (opts) => {
      if (!apiAvailable) {
        return Promise.reject(new Error('API not yet available'));
      }
      if (!modelVersionId) {
        return Promise.reject(new NotReadyError('No model registeredModel id'));
      }
      return api.getModelArtifactsByModelVersion(opts, modelVersionId);
    },
    [api, apiAvailable, modelVersionId],
  );
  return useFetchState(
    callback,
    { items: [], size: 0, pageSize: 0, nextPageToken: '' },
    { initialPromisePurity: true },
  );
};

export default useModelArtifactsByVersionId;
