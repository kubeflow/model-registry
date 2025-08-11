import * as React from 'react';
import { useFetchState, FetchState, FetchStateCallbackPromise } from 'mod-arch-core';
import { ModelVersionList } from '~/app/types';
import { useModelRegistryAPI } from '~/app/hooks/useModelRegistryAPI';

const useModelVersions = (): FetchState<ModelVersionList> => {
  const { api, apiAvailable } = useModelRegistryAPI();
  const callback = React.useCallback<FetchStateCallbackPromise<ModelVersionList>>(
    (opts) => {
      if (!apiAvailable) {
        return Promise.reject(new Error('API not yet available'));
      }
      return api.listModelVersions(opts).then((r) => r);
    },
    [api, apiAvailable],
  );
  return useFetchState(
    callback,
    { items: [], size: 0, pageSize: 0, nextPageToken: '' },
    { initialPromisePurity: true },
  );
};

export default useModelVersions;
