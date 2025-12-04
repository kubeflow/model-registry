import { useQueryParamNamespaces } from 'mod-arch-core';
import useGenericObjectState from 'mod-arch-core/dist/utilities/useGenericObjectState';
import * as React from 'react';
import { useLocation } from 'react-router-dom';
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
  ModelCatalogFilterStates,
} from '~/app/modelCatalogTypes';
import { BFF_API_VERSION, URL_PREFIX } from '~/app/utilities/const';
import {
  ModelCatalogStringFilterKey,
  ModelCatalogNumberFilterKey,
} from '~/concepts/modelCatalog/const';
import { isModelDetailsPage } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';

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
  setFilterData: <K extends keyof ModelCatalogFilterStates>(
    key: K,
    value: ModelCatalogFilterStates[K],
  ) => void;
  filterOptions: CatalogFilterOptionsList | null;
  filterOptionsLoaded: boolean;
  filterOptionsLoadError?: Error;
  performanceViewEnabled: boolean;
  setPerformanceViewEnabled: (enabled: boolean) => void;
  performanceFiltersChangedOnDetailsPage: boolean;
  setPerformanceFiltersChangedOnDetailsPage: (changed: boolean) => void;
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
    [ModelCatalogStringFilterKey.USE_CASE]: [],
    [ModelCatalogNumberFilterKey.MIN_RPS]: undefined,
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
  performanceViewEnabled: false,
  setPerformanceViewEnabled: () => undefined,
  performanceFiltersChangedOnDetailsPage: false,
  setPerformanceFiltersChangedOnDetailsPage: () => undefined,
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
    [ModelCatalogStringFilterKey.USE_CASE]: [],
    [ModelCatalogNumberFilterKey.MIN_RPS]: undefined,
  });
  const [filterOptions, filterOptionsLoaded, filterOptionsLoadError] =
    useCatalogFilterOptionList(apiState);
  const [selectedSourceLabel, setSelectedSourceLabel] = React.useState<
    ModelCatalogContextType['selectedSourceLabel']
  >(CategoryName.allModels);
  const [performanceViewEnabled, setPerformanceViewEnabled] = React.useState(false);
  const [performanceFiltersChangedOnDetailsPage, setPerformanceFiltersChangedOnDetailsPage] =
    React.useState(false);

  const location = useLocation();
  const prevLocationRef = React.useRef(location.pathname);

  React.useEffect(() => {
    const prevPath = prevLocationRef.current;
    const currentPath = location.pathname;

    const wasOnDetailsPath = isModelDetailsPage(prevPath);
    const isNowOnDetailsPath = isModelDetailsPage(currentPath);

    if (!wasOnDetailsPath && isNowOnDetailsPath) {
      setPerformanceFiltersChangedOnDetailsPage(false);
    }

    prevLocationRef.current = currentPath;
  }, [location.pathname]);

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
      performanceViewEnabled,
      setPerformanceViewEnabled,
      performanceFiltersChangedOnDetailsPage,
      setPerformanceFiltersChangedOnDetailsPage,
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
      performanceViewEnabled,
      performanceFiltersChangedOnDetailsPage,
    ],
  );

  return (
    <ModelCatalogContext.Provider value={contextValue}>{children}</ModelCatalogContext.Provider>
  );
};
