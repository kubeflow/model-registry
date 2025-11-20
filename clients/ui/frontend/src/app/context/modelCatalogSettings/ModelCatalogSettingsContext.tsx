import * as React from 'react';
import { useQueryParamNamespaces } from 'mod-arch-core';
import useModelCatalogSettingsAPIState, {
  ModelCatalogSettingsAPIState,
} from '~/app/hooks/modelCatalogSettings/useModelCatalogSettingsAPIState';
import { useCatalogSourceConfigs } from '~/app/hooks/modelCatalogSettings/useCatalogSourceConfigs';
import { CatalogSourceConfigList } from '~/app/modelCatalogTypes';
import { BFF_API_VERSION, URL_PREFIX } from '~/app/utilities/const';

export type ModelCatalogSettingsContextType = {
  apiState: ModelCatalogSettingsAPIState;
  refreshAPIState: () => void;
  catalogSourceConfigs: CatalogSourceConfigList | null;
  catalogSourceConfigsLoaded: boolean;
  catalogSourceConfigsLoadError?: Error;
  refreshCatalogSourceConfigs: () => void;
};

type ModelCatalogSettingsContextProviderProps = {
  children: React.ReactNode;
};

export const ModelCatalogSettingsContext = React.createContext<ModelCatalogSettingsContextType>({
  // eslint-disable-next-line @typescript-eslint/consistent-type-assertions
  apiState: { apiAvailable: false, api: null as unknown as ModelCatalogSettingsAPIState['api'] },
  refreshAPIState: () => undefined,
  catalogSourceConfigs: null,
  catalogSourceConfigsLoaded: false,
  catalogSourceConfigsLoadError: undefined,
  refreshCatalogSourceConfigs: () => undefined,
});

export const ModelCatalogSettingsContextProvider: React.FC<
  ModelCatalogSettingsContextProviderProps
> = ({ children }) => {
  const hostPath = `${URL_PREFIX}/api/${BFF_API_VERSION}/settings/model_catalog`;
  const queryParams = useQueryParamNamespaces();
  const [apiState, refreshAPIState] = useModelCatalogSettingsAPIState(hostPath, queryParams);
  const [
    catalogSourceConfigs,
    catalogSourceConfigsLoaded,
    catalogSourceConfigsLoadError,
    refreshCatalogSourceConfigs,
  ] = useCatalogSourceConfigs(apiState);

  const contextValue = React.useMemo(
    () => ({
      apiState,
      refreshAPIState,
      catalogSourceConfigs,
      catalogSourceConfigsLoaded,
      catalogSourceConfigsLoadError,
      refreshCatalogSourceConfigs,
    }),
    [
      apiState,
      refreshAPIState,
      catalogSourceConfigs,
      catalogSourceConfigsLoaded,
      catalogSourceConfigsLoadError,
      refreshCatalogSourceConfigs,
    ],
  );

  return (
    <ModelCatalogSettingsContext.Provider value={contextValue}>
      {children}
    </ModelCatalogSettingsContext.Provider>
  );
};
