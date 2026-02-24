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
      // TODO: Replace with a GET single job endpoint when the BFF supports it,
      // to avoid fetching the full job list and filtering client-side.
      return api.getModelTransferJobsByNamespace(opts, jobNamespace).then((jobList) => {
        const job = jobList.items.find((j) => j.name === jobName);
        return job ?? null;
      });
    },
    [api, apiAvailable, jobNamespace, jobName],
  );

  return useFetchState(callback, null, { initialPromisePurity: true });
};

export default useModelTransferJobForArtifact;
