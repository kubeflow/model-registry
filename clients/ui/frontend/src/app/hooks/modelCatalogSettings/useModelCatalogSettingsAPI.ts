import React from 'react';
import { ModelCatalogSettingsContext } from '~/app/context/modelCatalogSettings/ModelCatalogSettingsContext';
import { ModelCatalogSettingsAPIState } from './useModelCatalogSettingsAPIState';

type UseModelCatalogSettingsAPI = ModelCatalogSettingsAPIState & {
  refreshAllAPI: () => void;
};

export const useModelCatalogSettingsAPI = (): UseModelCatalogSettingsAPI => {
  const { apiState, refreshAPIState: refreshAllAPI } = React.useContext(
    ModelCatalogSettingsContext,
  );

  return {
    refreshAllAPI,
    ...apiState,
  };
};
