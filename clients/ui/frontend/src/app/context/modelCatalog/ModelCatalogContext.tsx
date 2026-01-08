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
  FilterOperator,
  ModelCatalogFilterStates,
  NamedQuery,
} from '~/app/modelCatalogTypes';
import {
  ModelDetailsTab,
  ModelCatalogStringFilterKey,
  ModelCatalogNumberFilterKey,
  UseCaseOptionValue,
  ALL_LATENCY_FIELD_NAMES,
  isLatencyMetricFieldName,
} from '~/concepts/modelCatalog/const';
import { isUseCaseOptionValue } from '~/app/pages/modelCatalog/utils/workloadTypeUtils';
import {
  getPerformanceFilterDefaultValue,
  resolveFilterValue,
  getLatencyFieldKey,
  getSingleFilterDefault,
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
   * Maps backend field names (e.g., 'artifacts.use_case.string_value') to frontend filter keys.
   * Handles special values like 'max' by resolving them from filter ranges.
   */
  const applyNamedQueryDefaults = React.useCallback(
    (namedQuery: NamedQuery) => {
      Object.entries(namedQuery).forEach(([fieldName, fieldFilter]) => {
        // Handle artifacts.* prefix - check both with and without prefix
        const isUseCase =
          fieldName === 'artifacts.use_case.string_value' || fieldName === 'use_case.string_value';
        const isHardwareType =
          fieldName === 'artifacts.hardware_type.string_value' ||
          fieldName === 'hardware_type.string_value';
        const isRps =
          fieldName === 'artifacts.requests_per_second.double_value' ||
          fieldName === 'requests_per_second.double_value';

        // Check if it's a latency field
        const latencyFieldKey = getLatencyFieldKey(fieldName);

        if (isUseCase) {
          // Extract string values and filter to only valid UseCaseOptionValue entries
          const rawValues =
            fieldFilter.operator === FilterOperator.IN && Array.isArray(fieldFilter.value)
              ? fieldFilter.value.filter((v): v is string => typeof v === 'string')
              : typeof fieldFilter.value === 'string'
                ? [fieldFilter.value]
                : [];
          // Filter to only valid UseCaseOptionValue entries
          const validValues: UseCaseOptionValue[] = rawValues.filter(isUseCaseOptionValue);
          baseSetFilterData(ModelCatalogStringFilterKey.USE_CASE, validValues);
        } else if (isHardwareType) {
          const values =
            fieldFilter.operator === FilterOperator.IN && Array.isArray(fieldFilter.value)
              ? fieldFilter.value.filter((v): v is string => typeof v === 'string')
              : typeof fieldFilter.value === 'string'
                ? [fieldFilter.value]
                : [];
          baseSetFilterData(ModelCatalogStringFilterKey.HARDWARE_TYPE, values);
        } else if (isRps) {
          const resolvedValue = resolveFilterValue(filterOptions, fieldName, fieldFilter.value);
          if (resolvedValue !== undefined) {
            baseSetFilterData(ModelCatalogNumberFilterKey.MAX_RPS, resolvedValue);
          }
        } else if (latencyFieldKey) {
          // Apply latency filter using the resolved field name
          const resolvedValue = resolveFilterValue(filterOptions, fieldName, fieldFilter.value);
          if (resolvedValue !== undefined) {
            baseSetFilterData(latencyFieldKey, resolvedValue);
          }
        }
      });
    },
    [baseSetFilterData, filterOptions],
  );

  const setPerformanceViewEnabled = React.useCallback(
    (enabled: boolean) => {
      setBasePerformanceViewEnabled(enabled);
      if (enabled) {
        // Apply default performance filters from namedQueries if available
        const defaultQuery = filterOptions?.namedQueries?.['default-performance-filters'];
        if (defaultQuery) {
          applyNamedQueryDefaults(defaultQuery);
        }
      } else {
        // Clear performance-related filters when toggle is disabled
        baseSetFilterData(ModelCatalogStringFilterKey.USE_CASE, []);
        baseSetFilterData(ModelCatalogStringFilterKey.HARDWARE_TYPE, []);
        baseSetFilterData(ModelCatalogNumberFilterKey.MAX_RPS, undefined);
        // Clear all latency filters
        ALL_LATENCY_FIELD_NAMES.forEach((fieldName) => {
          baseSetFilterData(fieldName, undefined);
        });
        setPerformanceFiltersChangedOnDetailsPage(false);
      }
    },
    [filterOptions?.namedQueries, applyNamedQueryDefaults, baseSetFilterData],
  );

  /**
   * Resets all filters when performance view is enabled:
   * - Clears basic filters (Task, Provider, License, Language)
   * - Resets performance filters to default values from namedQueries
   * This is used by "Reset all filters" button in the performance toolbar.
   */
  const resetPerformanceFiltersToDefaults = React.useCallback(() => {
    // Clear basic filters
    baseSetFilterData(ModelCatalogStringFilterKey.TASK, []);
    baseSetFilterData(ModelCatalogStringFilterKey.PROVIDER, []);
    baseSetFilterData(ModelCatalogStringFilterKey.LICENSE, []);
    baseSetFilterData(ModelCatalogStringFilterKey.LANGUAGE, []);

    // Clear all performance filters
    baseSetFilterData(ModelCatalogStringFilterKey.USE_CASE, []);
    baseSetFilterData(ModelCatalogStringFilterKey.HARDWARE_TYPE, []);
    baseSetFilterData(ModelCatalogNumberFilterKey.MAX_RPS, undefined);
    ALL_LATENCY_FIELD_NAMES.forEach((fieldName) => {
      baseSetFilterData(fieldName, undefined);
    });

    // Then apply performance defaults from namedQueries if available
    const defaultQuery = filterOptions?.namedQueries?.['default-performance-filters'];
    if (defaultQuery) {
      applyNamedQueryDefaults(defaultQuery);
    }
  }, [filterOptions?.namedQueries, applyNamedQueryDefaults, baseSetFilterData]);

  /**
   * Resets a single performance filter to its default value from namedQueries.
   * Used when clicking the undo button on individual performance filter chips.
   */
  const resetSinglePerformanceFilterToDefault = React.useCallback(
    (filterKey: keyof ModelCatalogFilterStates) => {
      const { hasDefault, value } = getSingleFilterDefault(filterOptions, filterKey);
      const filterKeyStr = String(filterKey);

      if (filterKey === ModelCatalogStringFilterKey.USE_CASE) {
        if (hasDefault && Array.isArray(value)) {
          // Filter to only valid UseCaseOptionValue entries
          const validValues: UseCaseOptionValue[] = value
            .filter((v): v is string => typeof v === 'string')
            .filter(isUseCaseOptionValue);
          baseSetFilterData(ModelCatalogStringFilterKey.USE_CASE, validValues);
        } else {
          // No default found, clear the filter
          baseSetFilterData(ModelCatalogStringFilterKey.USE_CASE, []);
        }
      } else if (filterKey === ModelCatalogStringFilterKey.HARDWARE_TYPE) {
        if (hasDefault && Array.isArray(value)) {
          baseSetFilterData(
            ModelCatalogStringFilterKey.HARDWARE_TYPE,
            value.filter((v): v is string => typeof v === 'string'),
          );
        } else {
          baseSetFilterData(ModelCatalogStringFilterKey.HARDWARE_TYPE, []);
        }
      } else if (filterKey === ModelCatalogNumberFilterKey.MAX_RPS) {
        if (hasDefault && typeof value === 'number') {
          baseSetFilterData(ModelCatalogNumberFilterKey.MAX_RPS, value);
        } else {
          baseSetFilterData(ModelCatalogNumberFilterKey.MAX_RPS, undefined);
        }
      } else if (isLatencyMetricFieldName(filterKeyStr)) {
        if (hasDefault && typeof value === 'number') {
          baseSetFilterData(filterKeyStr, value);
        } else {
          baseSetFilterData(filterKeyStr, undefined);
        }
      }
    },
    [filterOptions, baseSetFilterData],
  );

  /**
   * Gets the default value for a performance filter from namedQueries.
   * Wrapper around the utility function that provides filterOptions from context.
   */
  const getDefaultValueForPerformanceFilter = React.useCallback(
    (filterKey: keyof ModelCatalogFilterStates): string | number | string[] | undefined =>
      getPerformanceFilterDefaultValue(filterOptions, filterKey),
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
      resetPerformanceFiltersToDefaults,
      resetSinglePerformanceFilterToDefault,
      getDefaultValueForPerformanceFilter,
    ],
  );

  return (
    <ModelCatalogContext.Provider value={contextValue}>{children}</ModelCatalogContext.Provider>
  );
};
