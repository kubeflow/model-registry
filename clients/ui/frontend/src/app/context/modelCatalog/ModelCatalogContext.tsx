import { useQueryParamNamespaces } from 'mod-arch-core';
import * as React from 'react';
import { useCatalogSources } from '~/app/hooks/modelCatalog/useCatalogSources';
import useModelCatalogAPIState, {
  ModelCatalogAPIState,
} from '~/app/hooks/modelCatalog/useModelCatalogAPIState';
import { CatalogSource, CatalogSourceList } from '~/app/modelCatalogTypes';
import { BFF_API_VERSION } from '~/app/utilities/const';

export type ModelCatalogContextType = {
  catalogSourcesLoaded: boolean;
  catalogSourcesLoadError?: Error;
  catalogSources: CatalogSourceList | null;
  selectedSource: CatalogSource | undefined;
  updateSelectedSource: (modelRegistry: CatalogSource | undefined) => void;
  apiState: ModelCatalogAPIState;
  refreshAPIState: () => void;
};

type ModelCatalogContextProviderProps = {
  children: React.ReactNode;
};

export const ModelCatalogContext = React.createContext<ModelCatalogContextType>({
  catalogSourcesLoaded: false,
  catalogSourcesLoadError: undefined,
  catalogSources: null,
  selectedSource: undefined,
  updateSelectedSource: () => undefined,
  // eslint-disable-next-line @typescript-eslint/consistent-type-assertions
  apiState: { apiAvailable: false, api: null as unknown as ModelCatalogAPIState['api'] },
  refreshAPIState: () => undefined,
});

export const ModelCatalogContextProvider: React.FC<ModelCatalogContextProviderProps> = ({
  children,
}) => {
  const hostPath = `/api/${BFF_API_VERSION}/model_catalog`;
  const queryParams = useQueryParamNamespaces();
  const [apiState, refreshAPIState] = useModelCatalogAPIState(hostPath, queryParams);
  const [catalogSources, isLoaded, error] = useCatalogSources(apiState);
  const [selectedSource, setSelectedSource] =
    React.useState<ModelCatalogContextType['selectedSource']>(undefined);

  const contextValue = React.useMemo(
    () => ({
      catalogSourcesLoaded: isLoaded,
      catalogSourcesLoadError: error,
      catalogSources,
      selectedSource: selectedSource ?? undefined,
      updateSelectedSource: setSelectedSource,
      apiState,
      refreshAPIState,
    }),
    [isLoaded, error, catalogSources, selectedSource, apiState, refreshAPIState],
  );

  return (
    <ModelCatalogContext.Provider value={contextValue}>{children}</ModelCatalogContext.Provider>
  );
};
