import useFetch, { FetchState } from '../../utils/useFetch';
import { listModelRegistriesBackend } from '../../services/modelRegistrySettingsService';
import { POLL_INTERVAL } from '../../utils/const';
import * as React from 'react';
import { ModelRegistryKind } from '../../k8sTypes';

const useModelRegistriesBackend = (): FetchState<ModelRegistryKind[]> => {
  const getModelRegistries = React.useCallback(() => listModelRegistriesBackend(), []);
  return useFetch<ModelRegistryKind[]>(getModelRegistries, [], { refreshRate: POLL_INTERVAL });
};

export default useModelRegistriesBackend; 