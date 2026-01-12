import React from 'react';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import {
  CatalogArtifacts,
  CatalogArtifactType,
  CatalogFilterOptions,
  CatalogFilterOptionsList,
  CatalogModel,
  CatalogModelDetailsParams,
  CatalogSource,
  CatalogSourceList,
  ModelCatalogFilterStates,
  ModelCatalogStringFilterValueType,
  MetricsType,
  ModelCatalogFilterKey,
} from '~/app/modelCatalogTypes';
import { getLabels } from '~/app/pages/modelRegistry/screens/utils';
import {
  ModelCatalogStringFilterKey,
  ModelCatalogNumberFilterKey,
  ALL_LATENCY_FIELD_NAMES,
  LatencyMetricFieldName,
  DEFAULT_PERFORMANCE_FILTERS_QUERY_NAME,
} from '~/concepts/modelCatalog/const';
import { CatalogSourceStatus } from '~/concepts/modelCatalogSettings/const';

/**
 * Prefix used by the backend for artifact-specific filter options.
 * Filter options with this prefix are applicable to the artifacts endpoint.
 */
export const ARTIFACTS_FILTER_PREFIX = 'artifacts.';

export const extractVersionTag = (tags?: string[]): string | undefined =>
  tags?.find((tag) => /^\d+\.\d+\.\d+$/.test(tag));
export const filterNonVersionTags = (tags?: string[]): string[] | undefined => {
  const versionTag = extractVersionTag(tags);
  return tags?.filter((tag) => tag !== versionTag);
};

export const getModelName = (modelName: string): string => {
  const index = modelName.indexOf('/');
  if (index === -1) {
    return modelName;
  }
  return modelName.slice(index + 1);
};

export const decodeParams = (
  params: Readonly<CatalogModelDetailsParams>,
): CatalogModelDetailsParams =>
  Object.fromEntries(
    Object.entries(params).map(([key, value]) => [key, decodeURIComponent(value)]),
  );

export const encodeParams = (params: CatalogModelDetailsParams): CatalogModelDetailsParams =>
  Object.fromEntries(
    Object.entries(params).map(([key, value]) => [
      key,
      encodeURIComponent(value).replace(/\./g, '%252E'),
    ]),
  );

export const filterEnabledCatalogSources = (
  catalogSources: CatalogSourceList | null,
): CatalogSourceList | null => {
  if (!catalogSources) {
    return null;
  }

  // Filter sources that are enabled AND have available models
  const filteredItems = catalogSources.items?.filter(
    (source) => source.enabled !== false && source.status === CatalogSourceStatus.AVAILABLE,
  );

  return {
    ...catalogSources,
    items: filteredItems || [],
    size: filteredItems?.length || 0,
  };
};

export const getModelArtifactUri = (artifacts: CatalogArtifacts[]): string => {
  const modelArtifact = artifacts.find(
    (artifact) => artifact.artifactType === CatalogArtifactType.modelArtifact,
  );

  if (modelArtifact) {
    return modelArtifact.uri || '';
  }

  return '';
};

export const hasModelArtifacts = (artifacts: CatalogArtifacts[]): boolean =>
  artifacts.some((artifact) => artifact.artifactType === CatalogArtifactType.modelArtifact);

export const hasPerformanceArtifacts = (artifacts: CatalogArtifacts[]): boolean =>
  artifacts.some(
    (artifact) =>
      artifact.artifactType === CatalogArtifactType.metricsArtifact &&
      'metricsType' in artifact &&
      artifact.metricsType === MetricsType.performanceMetrics,
  );

// Utility function to check if a model is validated
export const isModelValidated = (model: CatalogModel): boolean => {
  if (!model.customProperties) {
    return false;
  }
  const labels = getLabels(model.customProperties);
  return labels.includes('validated');
};

export const shouldShowValidatedInsights = (
  model: CatalogModel,
  artifacts: CatalogArtifacts[],
): boolean => isModelValidated(model) && hasPerformanceArtifacts(artifacts);

export const useCatalogStringFilterState = <K extends ModelCatalogStringFilterKey>(
  filterKey: K,
): {
  isSelected: (value: ModelCatalogStringFilterValueType[K]) => boolean;
  setSelected: (value: ModelCatalogStringFilterValueType[K], selected: boolean) => void;
} => {
  type Value = ModelCatalogStringFilterValueType[K];
  const { filterData, setFilterData } = React.useContext(ModelCatalogContext);
  const selections: string[] = filterData[filterKey];
  const isValidStringState = (state: string[]): state is ModelCatalogFilterStates[K] =>
    Object.values(ModelCatalogStringFilterKey).includes(filterKey);
  const isSelected = React.useCallback((value: Value) => selections.includes(value), [selections]);
  const setSelected = (value: Value, selected: boolean) => {
    const nextState = selected
      ? [...selections, value]
      : selections.filter((item) => item !== value);
    if (isValidStringState(nextState)) {
      setFilterData(filterKey, nextState);
    }
  };

  return { isSelected, setSelected };
};

export const useCatalogNumberFilterState = (
  filterKey: ModelCatalogNumberFilterKey,
): {
  value: number | undefined;
  setValue: (value: number | undefined) => void;
} => {
  const { filterData, setFilterData } = React.useContext(ModelCatalogContext);
  const value = filterData[filterKey];
  const setValue = React.useCallback(
    (newValue: number | undefined) => {
      setFilterData(filterKey, newValue);
    },
    [filterKey, setFilterData],
  );
  return { value, setValue };
};

const isArrayOfSelections = (
  filterOption: CatalogFilterOptions[keyof CatalogFilterOptions],
  data: unknown,
): data is string[] =>
  filterOption?.type === 'string' && Array.isArray(filterOption.values) && Array.isArray(data);

/**
 * Filter IDs that use numeric comparison (latency filters + Max RPS).
 */
const KNOWN_NUMERIC_FILTER_IDS: string[] = [
  ...ALL_LATENCY_FIELD_NAMES,
  ModelCatalogNumberFilterKey.MAX_RPS,
];

/**
 * Type guard to check if a filter is a known numeric filter with a number value.
 */
const isKnownNumericFilter = (
  filterOption: CatalogFilterOptions[keyof CatalogFilterOptions],
  filterId: string,
  data: unknown,
): data is number =>
  filterOption?.type === 'number' &&
  KNOWN_NUMERIC_FILTER_IDS.includes(filterId) &&
  typeof data === 'number';

/**
 * Gets the comparison operator for a numeric filter from namedQueries.
 * Looks up the operator in the default performance filters namedQuery,
 * falls back to '<' if not found.
 */
const getNumericFilterOperator = (options: CatalogFilterOptionsList, filterId: string): string => {
  const defaultQuery = options.namedQueries?.[DEFAULT_PERFORMANCE_FILTERS_QUERY_NAME];
  if (defaultQuery && filterId in defaultQuery) {
    const fieldFilter = defaultQuery[filterId];
    // Return the operator from the namedQuery (e.g., '<=', '<', '>')
    return fieldFilter.operator;
  }
  // Fall back to '<' if this filter isn't in the namedQuery
  return '<';
};

const isFilterIdInMap = (
  filterId: unknown,
  filters: CatalogFilterOptions,
): filterId is keyof CatalogFilterOptions => typeof filterId === 'string' && filterId in filters;

/**
 * Gets the active latency field name from the filter state (if any)
 */
export const getActiveLatencyFieldName = (
  filterData: ModelCatalogFilterStates,
): LatencyMetricFieldName | undefined => {
  for (const fieldName of ALL_LATENCY_FIELD_NAMES) {
    const value = filterData[fieldName];
    if (value !== undefined && typeof value === 'number') {
      return fieldName;
    }
  }
  return undefined;
};

const wrapInQuotes = (v: string): string => `'${v}'`;

const eqFilter = (k: string, v: string) => `${k}=${wrapInQuotes(v)}`;
const inFilter = (k: string, values: string[]) =>
  `${k} IN (${values.map((v) => wrapInQuotes(v)).join(',')})`;

/**
 * Converts filter data into a filter query string for the /models endpoint.
 * Supports string filters (equality/IN) and numeric filters (less-than-or-equal for latency).
 * Filter keys are used directly as they already match backend format.
 * Note: RPS is NOT included in filterQuery - it's passed as targetRPS param instead.
 */
export const filtersToFilterQuery = (
  filterData: ModelCatalogFilterStates,
  options: CatalogFilterOptionsList,
): string => {
  const serializedFilters: string[] = Object.entries(filterData)
    .filter(([filterId]) => {
      // RPS is passed as targetRPS param, not in filterQuery
      if (filterId === ModelCatalogNumberFilterKey.MAX_RPS) {
        return false;
      }
      return true;
    })
    .map(([filterId, data]) => {
      if (typeof data === 'undefined') {
        return '';
      }

      // Check if this filter exists in the options
      const filterOption =
        options.filters && isFilterIdInMap(filterId, options.filters)
          ? options.filters[filterId]
          : undefined;

      if (!filterOption) {
        // Filter not found in options
        return '';
      }

      // Handle string array filters (multi-select)
      if (isArrayOfSelections(filterOption, data)) {
        switch (data.length) {
          case 0:
            return '';
          case 1:
            return eqFilter(filterId, data[0]);
          default:
            // 2 or more
            return inFilter(filterId, data);
        }
      }

      // Handle numeric filters - look up operator from namedQueries, fallback to '<'
      if (isKnownNumericFilter(filterOption, filterId, data)) {
        const operator = getNumericFilterOperator(options, filterId);
        return `${filterId} ${operator} ${data}`;
      }

      // Shouldn't reach this far, but if it does, log & ignore the case
      // eslint-disable-next-line no-console
      console.warn('Unhandled option', filterId, data, filterOption);
      return '';
    });

  const nonEmptyFilters = serializedFilters.filter((v) => !!v);

  // eg. filterQuery=license IN ('mit','apache-2.0') AND artifacts.hardware_type.string_value='H100'
  return nonEmptyFilters.length === 0 ? '' : nonEmptyFilters.join(' AND ');
};

/**
 * Find the server filter key for a given filter ID.
 * Handles the case where local state uses short keys but server uses fully qualified keys.
 */
const findServerFilterKey = (
  filterId: string,
  filters: CatalogFilterOptions | undefined,
): string | undefined => {
  if (!filters) {
    return undefined;
  }

  if (filterId in filters) {
    return filterId;
  }

  for (const suffix of ['.string_value', '.double_value', '.int_value', '.array_value']) {
    const key = `artifacts.${filterId}${suffix}`;
    if (key in filters) {
      return key;
    }
  }
  return undefined;
};

/**
 * Convert a server filter key to the format expected by the filterQuery.
 * Strips the 'artifacts.' prefix if present.
 */
const convertToFilterQueryKey = (serverFilterKey: string): string => {
  if (serverFilterKey.startsWith(ARTIFACTS_FILTER_PREFIX)) {
    return serverFilterKey.substring(ARTIFACTS_FILTER_PREFIX.length);
  }
  return serverFilterKey;
};

/**
 * Check if a filter ID is a latency field name.
 */
const isLatencyFieldName = (id: string): id is LatencyMetricFieldName =>
  ALL_LATENCY_FIELD_NAMES.some((name) => name === id);

/**
 * Converts filter data into a filter query string for the /artifacts/performance endpoint.
 * Only includes filters that have the 'artifacts.' prefix and strips that prefix in the output.
 * RPS is NOT included in filterQuery - it's passed as targetRPS param instead.
 */
export const filtersToArtifactsFilterQuery = (
  filterData: ModelCatalogFilterStates,
  options: CatalogFilterOptionsList,
): string => {
  const serializedFilters: string[] = Object.entries(filterData)
    .filter(
      ([filterId]) =>
        filterId === ModelCatalogStringFilterKey.HARDWARE_TYPE ||
        filterId === ModelCatalogStringFilterKey.USE_CASE ||
        isLatencyFieldName(filterId),
    )
    .map(([filterId, data]) => {
      if (typeof data === 'undefined') {
        return '';
      }

      const serverFilterKey = findServerFilterKey(filterId, options.filters);
      const queryKey = serverFilterKey ? convertToFilterQueryKey(serverFilterKey) : filterId;

      const filterOption =
        serverFilterKey && options.filters && isFilterIdInMap(serverFilterKey, options.filters)
          ? options.filters[serverFilterKey]
          : undefined;

      if (filterOption?.type === 'string' && Array.isArray(data)) {
        switch (data.length) {
          case 0:
            return '';
          case 1:
            return eqFilter(queryKey, data[0]);
          default:
            return inFilter(queryKey, data);
        }
      }

      if (
        filterOption?.type === 'number' &&
        isLatencyFieldName(filterId) &&
        typeof data === 'number'
      ) {
        return `${queryKey}<${data}`; // e.g., ttft_p90.double_value<60
      }

      return '';
    });

  const nonEmptyFilters = serializedFilters.filter((v) => !!v);
  return nonEmptyFilters.length === 0 ? '' : nonEmptyFilters.join(' AND ');
};

/**
 * Returns a copy of filterData with only basic (non-performance) filters.
 * Used when performance view is disabled to exclude performance filters from API queries.
 */
export const getBasicFiltersOnly = (
  filterData: ModelCatalogFilterStates,
): ModelCatalogFilterStates => {
  const result: ModelCatalogFilterStates = {
    ...filterData,
    // Clear performance string filters
    [ModelCatalogStringFilterKey.USE_CASE]: [],
    [ModelCatalogStringFilterKey.HARDWARE_TYPE]: [],
    // Clear performance number filter
    [ModelCatalogNumberFilterKey.MAX_RPS]: undefined,
  };

  // Clear all latency fields
  ALL_LATENCY_FIELD_NAMES.forEach((latencyKey) => {
    result[latencyKey] = undefined;
  });

  return result;
};

export const getUniqueSourceLabels = (catalogSources: CatalogSourceList | null): string[] => {
  if (!catalogSources || !catalogSources.items) {
    return [];
  }

  const allLabels = new Set<string>();

  catalogSources.items.forEach((source) => {
    // Only include labels from sources that are enabled AND have available models
    if (
      source.enabled &&
      source.status === CatalogSourceStatus.AVAILABLE &&
      source.labels.length > 0
    ) {
      source.labels.forEach((label) => {
        if (label.trim()) {
          allLabels.add(label.trim());
        }
      });
    }
  });

  return Array.from(allLabels);
};

export const hasSourcesWithoutLabels = (catalogSources: CatalogSourceList | null): boolean => {
  if (!catalogSources || !catalogSources.items) {
    return false;
  }

  return catalogSources.items.some((source) => {
    // Only consider sources that are enabled AND have available models
    if (source.enabled !== false && source.status === CatalogSourceStatus.AVAILABLE) {
      // Check if source has no labels or only empty/whitespace labels
      return source.labels.length === 0 || source.labels.every((label) => !label.trim());
    }
    return false;
  });
};

export const getSourceFromSourceId = (
  sourceId: string,
  catalogSources: CatalogSourceList | null,
): CatalogSource | undefined => {
  if (!catalogSources || !sourceId || !catalogSources.items) {
    return undefined;
  }

  return catalogSources.items.find((source) => source.id === sourceId);
};

/**
 * Checks if any filters are applied. If filterKeys is provided, only checks those specific filters.
 * Otherwise checks all filters.
 */
export const hasFiltersApplied = (
  filterData: ModelCatalogFilterStates,
  filterKeys?: ModelCatalogFilterKey[],
): boolean =>
  Object.entries(filterData).some(([key, value]) => {
    if (filterKeys && !filterKeys.some((k) => k === key)) {
      return false;
    }
    if (Array.isArray(value)) {
      return value.length > 0;
    }
    return value !== undefined;
  });

/**
 * Checks if a filter value differs from its default value.
 * Used to determine if a filter chip should be visible.
 */
export const isValueDifferentFromDefault = (
  currentValue: string | number | string[] | undefined,
  defaultValue: string | number | string[] | undefined,
): boolean => {
  if (defaultValue === undefined) {
    // No default defined, show chip if value exists
    if (Array.isArray(currentValue)) {
      return currentValue.length > 0;
    }
    return currentValue !== undefined;
  }

  if (currentValue === undefined) {
    return false;
  }

  // Compare arrays
  if (Array.isArray(currentValue) && Array.isArray(defaultValue)) {
    if (currentValue.length !== defaultValue.length) {
      return true;
    }
    return !currentValue.every((v) => defaultValue.includes(String(v)));
  }

  // Compare single value with array
  if (Array.isArray(currentValue) && !Array.isArray(defaultValue)) {
    if (currentValue.length !== 1) {
      return true;
    }
    return currentValue[0] !== defaultValue;
  }

  // Compare single values
  return currentValue !== defaultValue;
};

/**
 * Checks if there are any visible filter chips to display.
 * For performance filters, this accounts for default values - chips are only
 * visible when the value differs from the default.
 * For basic filters, any non-empty value is visible.
 *
 * @param filterData - Current filter state
 * @param filterKeys - Filter keys to check for visibility
 * @param getDefaultValue - Function to get default value for a filter key
 * @param performanceViewEnabled - Whether performance view is enabled
 * @param isPerformanceFilter - Function to check if a filter is a performance filter
 */
export const hasVisibleFilterChips = (
  filterData: ModelCatalogFilterStates,
  filterKeys: ModelCatalogFilterKey[],
  getDefaultValue: (key: ModelCatalogFilterKey) => string | number | string[] | undefined,
  performanceViewEnabled: boolean,
  isPerformanceFilter: (key: ModelCatalogFilterKey) => boolean,
): boolean =>
  filterKeys.some((key) => {
    const value = filterData[key];

    // Skip if no value
    if (value === undefined) {
      return false;
    }

    // For array values, skip if empty
    if (Array.isArray(value) && value.length === 0) {
      return false;
    }

    // For performance filters when toggle is on, check if differs from default
    if (performanceViewEnabled && isPerformanceFilter(key)) {
      const defaultValue = getDefaultValue(key);
      return isValueDifferentFromDefault(value, defaultValue);
    }

    // For basic filters, any value means visible
    return true;
  });

/**
 * Filters catalog sources to only include those with available models.
 * A source has models if its status is AVAILABLE.
 * This is used to filter out disabled sources or sources with errors from the switcher.
 */
export const filterSourcesWithModels = (
  catalogSources: CatalogSourceList | null,
): CatalogSourceList | null => {
  if (!catalogSources) {
    return null;
  }

  const filteredItems = catalogSources.items?.filter(
    (source) => source.status === CatalogSourceStatus.AVAILABLE,
  );

  return {
    ...catalogSources,
    items: filteredItems || [],
    size: filteredItems?.length || 0,
  };
};

/**
 * Checks if there are any catalog sources that have models available.
 * Returns true if at least one source has status === AVAILABLE.
 */
export const hasSourcesWithModels = (catalogSources: CatalogSourceList | null): boolean => {
  if (!catalogSources?.items) {
    return false;
  }

  return catalogSources.items.some((source) => source.status === CatalogSourceStatus.AVAILABLE);
};
