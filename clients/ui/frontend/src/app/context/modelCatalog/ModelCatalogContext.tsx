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
  NamedQuery,
} from '~/app/modelCatalogTypes';
import {
  ModelDetailsTab,
  ModelCatalogStringFilterKey,
  ModelCatalogNumberFilterKey,
  ALL_LATENCY_FIELD_NAMES,
  DEFAULT_PERFORMANCE_FILTERS_QUERY_NAME,
} from '~/concepts/modelCatalog/const';
import {
  getSingleFilterDefault,
  applyFilterValue,
  getDefaultFiltersFromNamedQuery,
} from '~/app/pages/modelCatalog/utils/performanceFilterUtils';
import { BFF_API_VERSION, URL_PREFIX } from '~/app/utilities/const';

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
  clearAllFilters: () => void;
  clearBasicFiltersAndResetPerformanceToDefaults: () => void;
  resetPerformanceFiltersToDefaults: () => void;
  resetSinglePerformanceFilterToDefault: (filterKey: keyof ModelCatalogFilterStates) => void;
  getPerformanceFilterDefaultValue: (
    filterKey: keyof ModelCatalogFilterStates,
  ) => string | number | string[] | undefined;
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
    [ModelCatalogNumberFilterKey.MAX_RPS]: undefined,
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
  clearAllFilters: () => undefined,
  clearBasicFiltersAndResetPerformanceToDefaults: () => undefined,
  resetPerformanceFiltersToDefaults: () => undefined,
  resetSinglePerformanceFilterToDefault: () => undefined,
  getPerformanceFilterDefaultValue: () => undefined,
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
  const [filterData, baseSetFilterData] = useGenericObjectState<ModelCatalogFilterStates>({
    [ModelCatalogStringFilterKey.TASK]: [],
    [ModelCatalogStringFilterKey.PROVIDER]: [],
    [ModelCatalogStringFilterKey.LICENSE]: [],
    [ModelCatalogStringFilterKey.LANGUAGE]: [],
    [ModelCatalogStringFilterKey.HARDWARE_TYPE]: [],
    [ModelCatalogStringFilterKey.USE_CASE]: [],
    [ModelCatalogNumberFilterKey.MAX_RPS]: undefined,
  });
  const [filterOptions, filterOptionsLoaded, filterOptionsLoadError] =
    useCatalogFilterOptionList(apiState);
  const [selectedSourceLabel, setSelectedSourceLabel] = React.useState<
    ModelCatalogContextType['selectedSourceLabel']
  >(CategoryName.allModels);
  const [basePerformanceViewEnabled, setBasePerformanceViewEnabled] = React.useState(false);
  const [performanceFiltersChangedOnDetailsPage, setPerformanceFiltersChangedOnDetailsPage] =
    React.useState(false);

  const location = useLocation();
  const isOnDetailsPage = location.pathname.includes(ModelDetailsTab.PERFORMANCE_INSIGHTS);

  /**
   * Applies filter values from a named query to the filter state.
   * Uses getDefaultFiltersFromNamedQuery to parse the namedQuery and applyFilterValue to set each filter.
   */
  const applyNamedQueryDefaults = React.useCallback(
    (namedQuery: NamedQuery) => {
      const defaults = getDefaultFiltersFromNamedQuery(filterOptions, namedQuery);
      Object.entries(defaults).forEach(([filterKey, value]) => {
        applyFilterValue(baseSetFilterData, filterKey, value);
      });
    },
    [baseSetFilterData, filterOptions],
  );

  /**
   * Clears all filters to their empty/undefined state.
   * Used on the landing page "Reset all filters" button.
   */
  const clearAllFilters = React.useCallback(() => {
    // Clear all string filters (arrays set to empty)
    Object.values(ModelCatalogStringFilterKey).forEach((filterKey) => {
      baseSetFilterData(filterKey, []);
    });

    // Clear all number filters (set to undefined)
    Object.values(ModelCatalogNumberFilterKey).forEach((filterKey) => {
      baseSetFilterData(filterKey, undefined);
    });

    // Clear all latency filters (set to undefined)
    ALL_LATENCY_FIELD_NAMES.forEach((fieldName) => {
      baseSetFilterData(fieldName, undefined);
    });
  }, [baseSetFilterData]);

  /**
   * Clears only performance filters to their empty/undefined state.
   * Used on the details page where we don't want to affect basic filters.
   */
  const clearPerformanceFilters = React.useCallback(() => {
    // Clear performance string filters
    baseSetFilterData(ModelCatalogStringFilterKey.USE_CASE, []);
    baseSetFilterData(ModelCatalogStringFilterKey.HARDWARE_TYPE, []);

    // Clear performance number filters
    baseSetFilterData(ModelCatalogNumberFilterKey.MAX_RPS, undefined);

    // Clear all latency filters
    ALL_LATENCY_FIELD_NAMES.forEach((fieldName) => {
      baseSetFilterData(fieldName, undefined);
    });
  }, [baseSetFilterData]);

  const setPerformanceViewEnabled = React.useCallback(
    (enabled: boolean) => {
      setBasePerformanceViewEnabled(enabled);
      if (enabled) {
        // Apply default performance filters from namedQueries if available
        const defaultQuery = filterOptions?.namedQueries?.[DEFAULT_PERFORMANCE_FILTERS_QUERY_NAME];
        if (defaultQuery) {
          applyNamedQueryDefaults(defaultQuery);
        }
      } else {
        // Clear performance-related filters when toggle is disabled
        clearPerformanceFilters();
        setPerformanceFiltersChangedOnDetailsPage(false);
      }
    },
    [filterOptions?.namedQueries, applyNamedQueryDefaults, clearPerformanceFilters],
  );

  /**
   * Resets performance filters then applies defaults from namedQueries.
   * This is used by "Reset all filters" button on the details page.
   */
  const resetPerformanceFiltersToDefaults = React.useCallback(() => {
    clearPerformanceFilters();

    // Then apply performance defaults from namedQueries if available
    const defaultQuery = filterOptions?.namedQueries?.[DEFAULT_PERFORMANCE_FILTERS_QUERY_NAME];
    if (defaultQuery) {
      applyNamedQueryDefaults(defaultQuery);
    }
  }, [clearPerformanceFilters, filterOptions?.namedQueries, applyNamedQueryDefaults]);

  /**
   * Clears basic filters (Task, Provider, License, Language) completely,
   * and resets performance filters to their default values.
   * This is used by "Reset all filters" on the landing page when performance view is enabled.
   */
  const clearBasicFiltersAndResetPerformanceToDefaults = React.useCallback(() => {
    // Clear basic filters (non-performance string filters)
    baseSetFilterData(ModelCatalogStringFilterKey.TASK, []);
    baseSetFilterData(ModelCatalogStringFilterKey.PROVIDER, []);
    baseSetFilterData(ModelCatalogStringFilterKey.LICENSE, []);
    baseSetFilterData(ModelCatalogStringFilterKey.LANGUAGE, []);

    // Reset performance filters to defaults (reuse existing function)
    resetPerformanceFiltersToDefaults();
  }, [baseSetFilterData, resetPerformanceFiltersToDefaults]);

  /**
   * Resets a single performance filter to its default value from namedQueries.
   * Used when clicking the undo button on individual performance filter chips.
   */
  const resetSinglePerformanceFilterToDefault = React.useCallback(
    (filterKey: keyof ModelCatalogFilterStates) => {
      const { value } = getSingleFilterDefault(filterOptions, filterKey);
      applyFilterValue(baseSetFilterData, filterKey, value);
    },
    [filterOptions, baseSetFilterData],
  );

  /**
   * Gets the default value for a performance filter from namedQueries.
   * Wrapper around the utility function that provides filterOptions from context.
   */
  const getDefaultValueForPerformanceFilter = React.useCallback(
    (filterKey: keyof ModelCatalogFilterStates): string | number | string[] | undefined => {
      const { value } = getSingleFilterDefault(filterOptions, filterKey);
      // Return value - the type is already compatible
      if (Array.isArray(value) || typeof value === 'string' || typeof value === 'number') {
        return value;
      }
      return undefined;
    },
    [filterOptions],
  );

  const setFilterData = React.useCallback(
    <K extends keyof ModelCatalogFilterStates>(key: K, value: ModelCatalogFilterStates[K]) => {
      baseSetFilterData(key, value);
      if (isOnDetailsPage) {
        setPerformanceFiltersChangedOnDetailsPage(true);
      } else {
        setPerformanceFiltersChangedOnDetailsPage(false);
      }
    },
    [baseSetFilterData, isOnDetailsPage],
  );

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
      performanceViewEnabled: basePerformanceViewEnabled,
      setPerformanceViewEnabled,
      performanceFiltersChangedOnDetailsPage,
      setPerformanceFiltersChangedOnDetailsPage,
      clearAllFilters,
      clearBasicFiltersAndResetPerformanceToDefaults,
      resetPerformanceFiltersToDefaults,
      resetSinglePerformanceFilterToDefault,
      getPerformanceFilterDefaultValue: getDefaultValueForPerformanceFilter,
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
      basePerformanceViewEnabled,
      setPerformanceViewEnabled,
      performanceFiltersChangedOnDetailsPage,
      setPerformanceFiltersChangedOnDetailsPage,
      clearAllFilters,
      clearBasicFiltersAndResetPerformanceToDefaults,
      resetPerformanceFiltersToDefaults,
      resetSinglePerformanceFilterToDefault,
      getDefaultValueForPerformanceFilter,
    ],
  );

  return (
    <ModelCatalogContext.Provider value={contextValue}>{children}</ModelCatalogContext.Provider>
  );
};
