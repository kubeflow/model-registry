import * as React from 'react';
import { useFetchState, FetchState, FetchStateCallbackPromise } from 'mod-arch-core';
import { ModelTransferJobList } from '~/app/types';
import { useModelRegistryAPI } from '~/app/hooks/useModelRegistryAPI';

const useModelTransferJobs = (): FetchState<ModelTransferJobList> => {
  const { api, apiAvailable } = useModelRegistryAPI();

  const callback = React.useCallback<FetchStateCallbackPromise<ModelTransferJobList>>(
    (opts) => {
      if (!apiAvailable) {
        return Promise.reject(new Error('API not yet available'));
      }
      return api.listModelTransferJobs(opts);
    },
    [api, apiAvailable],
  );

  return useFetchState(
    callback,
    { items: [], size: 0, pageSize: 0, nextPageToken: '' },
    { initialPromisePurity: true },
  );
};

export default useModelTransferJobs;
