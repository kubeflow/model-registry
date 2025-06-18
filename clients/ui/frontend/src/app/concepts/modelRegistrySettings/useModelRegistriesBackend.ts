import * as React from 'react';
import { ModelRegistryKind } from '~/app/k8sTypes';
import useFetch, { FetchState } from '~/app/utils/useFetch';
import { listModelRegistriesBackend } from '~/app/services/modelRegistrySettingsService';
import { POLL_INTERVAL } from '~/app/utils/const';

const useModelRegistriesBackend = (): FetchState<ModelRegistryKind[]> => {
  const getModelRegistries = React.useCallback(() => listModelRegistriesBackend(), []);
  return useFetch<ModelRegistryKind[]>(getModelRegistries, [], { refreshRate: POLL_INTERVAL });
};

export default useModelRegistriesBackend;
