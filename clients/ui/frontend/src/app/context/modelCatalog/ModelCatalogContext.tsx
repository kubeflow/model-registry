import { useQueryParamNamespaces } from 'mod-arch-core';
import useGenericObjectState from 'mod-arch-core/dist/utilities/useGenericObjectState';
import * as React from 'react';
import { useCatalogFilterOptionList } from '~/app/hooks/modelCatalog/useCatalogFilterOptionList';
import { useCatalogSources } from '~/app/hooks/modelCatalog/useCatalogSources';
import useModelCatalogAPIState, {
  ModelCatalogAPIState,
} from '~/app/hooks/modelCatalog/useModelCatalogAPIState';
import {
  CatalogFilterOptionsList,
  CatalogSource,
  CatalogSourceList,
  CategoryName,
  ModelCatalogFilterKey,
  ModelCatalogFilterStates,
} from '~/app/modelCatalogTypes';
import { BFF_API_VERSION, URL_PREFIX } from '~/app/utilities/const';
import {
  ModelCatalogStringFilterKey,
  ModelCatalogNumberFilterKey,
} from '~/concepts/modelCatalog/const';

export type ModelCatalogContextType = {
  catalogSourcesLoaded: boolean;
  catalogSourcesLoadError?: Error;
  catalogSources: CatalogSourceList | null;
  selectedSource: CatalogSource | undefined;
  updateSelectedSource: (source: CatalogSource | undefined) => void;
  selectedSourceLabel: string | undefined;
  updateSelectedSourceLabel: (sourceLabel: string | undefined) => void;
  apiState: ModelCatalogAPIState;
  refreshAPIState: () => void;
  filterData: ModelCatalogFilterStates;
  setFilterData: <K extends ModelCatalogFilterKey>(
    key: K,
    value: ModelCatalogFilterStates[K],
  ) => void;
  filterOptions: CatalogFilterOptionsList | null;
  filterOptionsLoaded: boolean;
  filterOptionsLoadError?: Error;
};

type ModelCatalogContextProviderProps = {
  children: React.ReactNode;
};

export const ModelCatalogContext = React.createContext<ModelCatalogContextType>({
  catalogSourcesLoaded: false,
  catalogSourcesLoadError: undefined,
  catalogSources: null,
  selectedSource: undefined,
  filterData: {
    [ModelCatalogStringFilterKey.TASK]: [],
    [ModelCatalogStringFilterKey.PROVIDER]: [],
    [ModelCatalogStringFilterKey.LICENSE]: [],
    [ModelCatalogStringFilterKey.LANGUAGE]: [],
    [ModelCatalogStringFilterKey.HARDWARE_TYPE]: [],
    [ModelCatalogNumberFilterKey.MIN_RPS]: undefined,
    [ModelCatalogNumberFilterKey.MAX_LATENCY]: undefined,
    [ModelCatalogNumberFilterKey.MAX_INPUT_TOKENS]: undefined,
    [ModelCatalogNumberFilterKey.MAX_OUTPUT_TOKENS]: undefined,
  },
  updateSelectedSource: () => undefined,
  selectedSourceLabel: undefined,
  updateSelectedSourceLabel: () => undefined,
  // eslint-disable-next-line @typescript-eslint/consistent-type-assertions
  apiState: { apiAvailable: false, api: null as unknown as ModelCatalogAPIState['api'] },
  refreshAPIState: () => undefined,
  setFilterData: () => undefined,
  filterOptions: null,
  filterOptionsLoaded: false,
  filterOptionsLoadError: undefined,
});

export const ModelCatalogContextProvider: React.FC<ModelCatalogContextProviderProps> = ({
  children,
}) => {
  const hostPath = `${URL_PREFIX}/api/${BFF_API_VERSION}/model_catalog`;
  const queryParams = useQueryParamNamespaces();
  const [apiState, refreshAPIState] = useModelCatalogAPIState(hostPath, queryParams);
  const [catalogSources, catalogSourcesLoaded, catalogSourcesLoadError] =
    useCatalogSources(apiState);
  const [selectedSource, setSelectedSource] =
    React.useState<ModelCatalogContextType['selectedSource']>(undefined);
  const [filterData, setFilterData] = useGenericObjectState<ModelCatalogFilterStates>({
    [ModelCatalogStringFilterKey.TASK]: [],
    [ModelCatalogStringFilterKey.PROVIDER]: [],
    [ModelCatalogStringFilterKey.LICENSE]: [],
    [ModelCatalogStringFilterKey.LANGUAGE]: [],
    [ModelCatalogStringFilterKey.HARDWARE_TYPE]: [],
    [ModelCatalogNumberFilterKey.MIN_RPS]: undefined,
    [ModelCatalogNumberFilterKey.MAX_LATENCY]: undefined,
    [ModelCatalogNumberFilterKey.MAX_INPUT_TOKENS]: undefined,
    [ModelCatalogNumberFilterKey.MAX_OUTPUT_TOKENS]: undefined,
  });
  const [filterOptions, filterOptionsLoaded, filterOptionsLoadError] =
    useCatalogFilterOptionList(apiState);
  const [selectedSourceLabel, setSelectedSourceLabel] = React.useState<
    ModelCatalogContextType['selectedSourceLabel']
  >(CategoryName.allModels);

  const contextValue = React.useMemo(
    () => ({
      catalogSourcesLoaded,
      catalogSourcesLoadError,
      catalogSources,
      selectedSource: selectedSource ?? undefined,
      updateSelectedSource: setSelectedSource,
      selectedSourceLabel: selectedSourceLabel ?? undefined,
      updateSelectedSourceLabel: setSelectedSourceLabel,
      apiState,
      refreshAPIState,
      filterData,
      setFilterData,
      filterOptions,
      filterOptionsLoaded,
      filterOptionsLoadError,
    }),
    [
      catalogSourcesLoaded,
      catalogSourcesLoadError,
      catalogSources,
      selectedSource,
      apiState,
      refreshAPIState,
      filterData,
      setFilterData,
      filterOptions,
      filterOptionsLoaded,
      filterOptionsLoadError,
      selectedSourceLabel,
    ],
  );

  return (
    <ModelCatalogContext.Provider value={contextValue}>{children}</ModelCatalogContext.Provider>
  );
};
