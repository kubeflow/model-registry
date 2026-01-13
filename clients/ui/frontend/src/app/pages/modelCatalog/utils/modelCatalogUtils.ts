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
  ALL_LATENCY_FILTER_KEYS,
  LatencyMetricFieldName,
  DEFAULT_PERFORMANCE_FILTERS_QUERY_NAME,
  isPerformanceStringFilterKey,
  PERFORMANCE_FILTER_KEYS,
} from '~/concepts/modelCatalog/const';
import { CatalogSourceStatus } from '~/concepts/modelCatalogSettings/const';

/**
 * Prefix used by the backend for artifact-specific filter options.
 * Filter options with this prefix are applicable to the artifacts endpoint.
 */
const ARTIFACTS_FILTER_PREFIX = 'artifacts.';

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
  ...ALL_LATENCY_FILTER_KEYS,
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
  for (const fieldName of ALL_LATENCY_FILTER_KEYS) {
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
 * Check if a filter key has the artifacts.* prefix.
 * Filter keys with this prefix are artifact-specific filters.
 */
const hasArtifactsPrefix = (filterId: string): boolean =>
  filterId.startsWith(ARTIFACTS_FILTER_PREFIX);

/**
 * Strips the 'artifacts.' prefix from a filter key if present.
 * Used when constructing filterQuery for artifacts endpoint.
 * Example: 'artifacts.use_case.string_value' -> 'use_case.string_value'
 */
const stripArtifactsPrefix = (filterId: string): string => {
  if (hasArtifactsPrefix(filterId)) {
    return filterId.substring(ARTIFACTS_FILTER_PREFIX.length);
  }
  return filterId;
};

/**
 * Target endpoint type for filter query construction.
 * - 'models': Include all filters (except RPS), use filter keys directly
 * - 'artifacts': Only include artifact-prefixed filters, strip the prefix in output
 */
export type FilterQueryTarget = 'models' | 'artifacts';

/**
 * Determines if a filter should be included based on the target endpoint.
 * - For models: Include all filters except RPS (which is passed as a separate param)
 * - For artifacts: Only include filters that have the artifacts.* prefix
 */
const shouldIncludeFilter = (filterId: string, target: FilterQueryTarget): boolean => {
  // RPS is always passed as a separate param, not in filterQuery
  if (filterId === ModelCatalogNumberFilterKey.MAX_RPS) {
    return false;
  }

  if (target === 'models') {
    // For models, include all filters (except RPS which is already excluded)
    return true;
  }

  // For artifacts, only include filters with the artifacts.* prefix
  return hasArtifactsPrefix(filterId);
};

/**
 * Gets the query key to use in the filterQuery string.
 * - For models: Use the filter ID directly (it already includes artifacts.* prefix if needed)
 * - For artifacts: Strip the artifacts.* prefix (the endpoint doesn't need it)
 */
const getQueryKey = (filterId: string, target: FilterQueryTarget): string => {
  if (target === 'artifacts') {
    return stripArtifactsPrefix(filterId);
  }
  return filterId;
};

/**
 * Serializes a single filter entry to a filter query clause.
 * Handles string arrays (IN/equality) and numeric filters (comparison operators).
 */
const serializeFilterEntry = (
  filterId: string,
  data: ModelCatalogFilterStates[keyof ModelCatalogFilterStates],
  options: CatalogFilterOptionsList,
  target: FilterQueryTarget,
): string => {
  if (typeof data === 'undefined') {
    return '';
  }

  // Get the filter option from the options map
  const filterOption =
    options.filters && isFilterIdInMap(filterId, options.filters)
      ? options.filters[filterId]
      : undefined;

  if (!filterOption) {
    return '';
  }

  const queryKey = getQueryKey(filterId, target);

  // Handle string array filters (multi-select)
  if (isArrayOfSelections(filterOption, data)) {
    switch (data.length) {
      case 0:
        return '';
      case 1:
        return eqFilter(queryKey, data[0]);
      default:
        return inFilter(queryKey, data);
    }
  }

  // Handle numeric filters
  if (isKnownNumericFilter(filterOption, filterId, data)) {
    const operator = getNumericFilterOperator(options, filterId);
    return `${queryKey} ${operator} ${data}`;
  }

  return '';
};

/**
 * Converts filter data into a filter query string.
 *
 * @param filterData - The current filter state
 * @param options - Filter options from the server (includes namedQueries for operators)
 * @param target - The target endpoint:
 *   - 'models': Include all filters (except RPS), use filter keys directly
 *   - 'artifacts': Only include artifact-prefixed filters, strip the prefix in output
 *
 * Note: RPS is NOT included in filterQuery for either target - it's passed as targetRPS param.
 */
export const filtersToFilterQuery = (
  filterData: ModelCatalogFilterStates,
  options: CatalogFilterOptionsList,
  target: FilterQueryTarget = 'models',
): string => {
  const serializedFilters: string[] = Object.entries(filterData)
    .filter(([filterId]) => shouldIncludeFilter(filterId, target))
    .map(([filterId, data]) => serializeFilterEntry(filterId, data, options, target));

  const nonEmptyFilters = serializedFilters.filter((v) => !!v);
  return nonEmptyFilters.length === 0 ? '' : nonEmptyFilters.join(' AND ');
};

/**
 * Returns a copy of filterData with only basic (non-performance) filters.
 * Used when performance view is disabled to exclude performance filters from API queries.
 * Performance filters are cleared to their empty state ([] for arrays, undefined for numbers).
 */
export const getBasicFiltersOnly = (
  filterData: ModelCatalogFilterStates,
): ModelCatalogFilterStates => {
  // Start with a copy of filterData
  const result: ModelCatalogFilterStates = { ...filterData };

  // Clear all performance filters using the centralized list
  PERFORMANCE_FILTER_KEYS.forEach((perfKey) => {
    if (isPerformanceStringFilterKey(perfKey)) {
      // String filters clear to empty array
      result[perfKey] = [];
    } else {
      // Number filters (MAX_RPS and latency) clear to undefined
      result[perfKey] = undefined;
    }
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
