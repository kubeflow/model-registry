import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import { ModelCatalogAPIState } from './useModelCatalogAPIState';
import React from 'react';

type UseModelRegistryAPI = ModelCatalogAPIState & {
  refreshAllAPI: () => void;
};

export const useModelCatalogAPI = (): UseModelRegistryAPI => {
  const { apiState, refreshAPIState: refreshAllAPI } = React.useContext(ModelCatalogContext);

  return {
    refreshAllAPI,
    ...apiState,
  };
};
