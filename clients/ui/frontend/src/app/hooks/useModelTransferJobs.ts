import * as React from 'react';
import { useFetchState, FetchState, FetchStateCallbackPromise, POLL_INTERVAL } from 'mod-arch-core';
import { ModelTransferJobList, ModelTransferJobStatus } from '~/app/types';
import { useModelRegistryAPI } from '~/app/hooks/useModelRegistryAPI';

const useModelTransferJobs = (): FetchState<ModelTransferJobList> => {
  const { api, apiAvailable } = useModelRegistryAPI();
  const [hasActiveJobs, setHasActiveJobs] = React.useState(false);

  const callback = React.useCallback<FetchStateCallbackPromise<ModelTransferJobList>>(
    (opts) => {
      if (!apiAvailable) {
        return Promise.reject(new Error('API not yet available'));
      }
      return api.listModelTransferJobs(opts).then((result) => {
        const active = result.items.some(
          (job) =>
            job.status === ModelTransferJobStatus.RUNNING ||
            job.status === ModelTransferJobStatus.PENDING,
        );
        setHasActiveJobs(active);
        return result;
      });
    },
    [api, apiAvailable],
  );

  return useFetchState(
    callback,
    { items: [], size: 0, pageSize: 0, nextPageToken: '' },
    {
      initialPromisePurity: true,
      refreshRate: hasActiveJobs ? POLL_INTERVAL : undefined,
    },
  );
};

export default useModelTransferJobs;
