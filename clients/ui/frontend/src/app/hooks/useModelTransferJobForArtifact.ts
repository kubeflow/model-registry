import * as React from 'react';
import { useFetchState, FetchState, FetchStateCallbackPromise } from 'mod-arch-core';
import { ModelArtifact, ModelTransferJob } from '~/app/types';
import { useModelRegistryAPI } from '~/app/hooks/useModelRegistryAPI';
import { modelSourcePropertiesToTransferJobParams } from '~/concepts/modelRegistry/utils';

const useModelTransferJobForArtifact = (
  modelArtifact: ModelArtifact | null,
): FetchState<ModelTransferJob | null> => {
  const { api, apiAvailable } = useModelRegistryAPI();
  const transferJobParams = modelArtifact
    ? modelSourcePropertiesToTransferJobParams(modelArtifact)
    : null;
  const jobNamespace = transferJobParams?.jobNamespace;
  const jobName = transferJobParams?.jobName;

  const callback = React.useCallback<FetchStateCallbackPromise<ModelTransferJob | null>>(
    (opts) => {
      if (!apiAvailable) {
        return Promise.reject(new Error('API not yet available'));
      }
      if (!jobNamespace || !jobName) {
        return Promise.resolve(null);
      }
      return api.getModelTransferJobByName(opts, jobNamespace, jobName);
    },
    [api, apiAvailable, jobNamespace, jobName],
  );

  return useFetchState(callback, null, { initialPromisePurity: true });
};

export default useModelTransferJobForArtifact;
