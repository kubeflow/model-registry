import * as React from 'react';
import { ModelRegistryKind } from '~/app/k8sTypes';
import { useFetchState, FetchStateObject } from 'mod-arch-shared';
import { listModelRegistriesBackend } from '~/app/services/modelRegistrySettingsService';
import { POLL_INTERVAL } from 'mod-arch-shared';

const useModelRegistriesBackend = (): FetchStateObject<ModelRegistryKind[]> => {
  const getModelRegistries = React.useCallback(() => listModelRegistriesBackend(), []);
  const [data, loaded, error, refresh] = useFetchState<ModelRegistryKind[]>(getModelRegistries, [], {
    refreshRate: POLL_INTERVAL,
  });
  return { data, loaded, error, refresh };
};

export default useModelRegistriesBackend;
