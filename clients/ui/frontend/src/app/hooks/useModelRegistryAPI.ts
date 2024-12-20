import * as React from 'react';
import { ModelRegistryAPIState } from '~/app/hooks/useModelRegistryAPIState';
import { ModelRegistryContext } from '~/app/context/ModelRegistryContext';

type UseModelRegistryAPI = ModelRegistryAPIState & {
  refreshAllAPI: () => void;
};

export const useModelRegistryAPI = (): UseModelRegistryAPI => {
  const { apiState, refreshAPIState: refreshAllAPI } = React.useContext(ModelRegistryContext);

  return {
    refreshAllAPI,
    ...apiState,
  };
};
