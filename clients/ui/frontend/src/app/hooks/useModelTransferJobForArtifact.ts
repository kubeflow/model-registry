import * as React from 'react';
import { useFetchState, FetchState, FetchStateCallbackPromise, NotReadyError } from 'mod-arch-core';
import { ModelTransferJob } from '~/app/types';
import { useModelRegistryAPI } from '~/app/hooks/useModelRegistryAPI';
import { TransferJobParams } from '~/concepts/modelRegistry/types';

const useModelTransferJobForArtifact = (
  transferJobParams: TransferJobParams | null,
): FetchState<ModelTransferJob | null> => {
  const { api, apiAvailable } = useModelRegistryAPI();
  const jobNamespace = transferJobParams?.jobNamespace;
  const jobName = transferJobParams?.jobName;

  const callback = React.useCallback<FetchStateCallbackPromise<ModelTransferJob | null>>(
    (opts) => {
      if (!apiAvailable) {
        return Promise.reject(new NotReadyError('API not yet available'));
      }
      if (!jobNamespace || !jobName) {
        return Promise.reject(new NotReadyError('No jobName or jobNamespace'));
      }
      return api.getModelTransferJobByName(opts, jobNamespace, jobName);
    },
    [api, apiAvailable, jobNamespace, jobName],
  );

  return useFetchState(callback, null, { initialPromisePurity: true });
};

export default useModelTransferJobForArtifact;
