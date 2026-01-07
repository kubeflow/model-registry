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
 * Filter IDs that should use "less than" comparison (latency filters).
 * All latency field names use less-than comparison.
 */
const KNOWN_LESS_THAN_IDS: string[] = ALL_LATENCY_FIELD_NAMES;

const isKnownLessThanValue = (
  filterOption: CatalogFilterOptions[keyof CatalogFilterOptions],
  filterId: string,
  data: unknown,
): data is number =>
  filterOption?.type === 'number' &&
  KNOWN_LESS_THAN_IDS.includes(filterId) &&
  typeof data === 'number';

/**
 * Filter IDs that should use "greater than" comparison (RPS filter).
 */
const KNOWN_GREATER_THAN_IDS: string[] = [ModelCatalogNumberFilterKey.MIN_RPS];

const isKnownGreaterThanValue = (
  filterOption: CatalogFilterOptions[keyof CatalogFilterOptions],
  filterId: string,
  data: unknown,
): data is number =>
  filterOption?.type === 'number' &&
  KNOWN_GREATER_THAN_IDS.includes(filterId) &&
  typeof data === 'number';

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
 * Supports string filters (equality/IN), and numeric filters (greater than for RPS, less than for latency).
 */
export const filtersToFilterQuery = (
  filterData: ModelCatalogFilterStates,
  options: CatalogFilterOptionsList,
): string => {
  const serializedFilters: string[] = Object.entries(filterData).map(([filterId, data]) => {
    if (
      !options.filters ||
      !isFilterIdInMap(filterId, options.filters) ||
      typeof data === 'undefined'
    ) {
      // Unhandled key or no data
      return '';
    }

    const filterOption = options.filters[filterId];

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

    // Numeric filters: less-than for latency, greater-than for RPS
    if (isKnownLessThanValue(filterOption, filterId, data)) {
      return `${filterId} < ${data}`;
    }

    if (isKnownGreaterThanValue(filterOption, filterId, data)) {
      return `${filterId} > ${data}`;
    }

    // Shouldn't reach this far, but if it does, log & ignore the case
    // eslint-disable-next-line no-console
    console.warn('Unhandled option', filterId, data, filterOption);
    return '';
  });

  const nonEmptyFilters = serializedFilters.filter((v) => !!v);

  // eg. filterQuery=rps_mean > 1 AND license IN ('mit','apache-2.0') AND ttft_mean < 10
  return nonEmptyFilters.length === 0 ? '' : nonEmptyFilters.join(' AND ');
};

/**
 * Converts filter data into a filter query string for the /artifacts/performance endpoint.
 * Only includes filters that have the 'artifacts.' prefix and strips that prefix in the output.
 * RPS is NOT included in filterQuery for artifacts - it's passed as targetRPS param instead.
 */
export const filtersToArtifactsFilterQuery = (
  filterData: ModelCatalogFilterStates,
  options: CatalogFilterOptionsList,
): string => {
  const isLatencyFieldName = (id: string): boolean =>
    ALL_LATENCY_FIELD_NAMES.some((name) => name === id);

  const serializedFilters: string[] = Object.entries(filterData)
    .filter(([filterId]) => {
      // Only include artifact-specific filters (those with artifacts. prefix in filter options)
      // OR performance-related filters like latency and hardware_type/use_case
      // But NOT rps_mean - that goes to targetRPS param
      if (filterId === ModelCatalogNumberFilterKey.MIN_RPS) {
        return false; // RPS is passed as targetRPS param, not in filterQuery
      }
      // Include latency filters, hardware_type, and use_case for artifacts filtering
      if (isLatencyFieldName(filterId)) {
        return true;
      }
      if (
        filterId === ModelCatalogStringFilterKey.HARDWARE_TYPE ||
        filterId === ModelCatalogStringFilterKey.USE_CASE
      ) {
        return true;
      }
      return false;
    })
    .map(([filterId, data]) => {
      if (typeof data === 'undefined') {
        return '';
      }

      // For artifacts endpoint, we use the filter ID directly (no prefix stripping needed
      // since our local state doesn't have the prefix - the backend filter_options have it)
      const filterOption =
        options.filters && isFilterIdInMap(filterId, options.filters)
          ? options.filters[filterId]
          : undefined;

      if (isArrayOfSelections(filterOption, data)) {
        switch (data.length) {
          case 0:
            return '';
          case 1:
            return eqFilter(filterId, data[0]);
          default:
            return inFilter(filterId, data);
        }
      }

      // Numeric filters for artifacts: latency uses less-than
      if (isKnownLessThanValue(filterOption, filterId, data)) {
        return `${filterId} < ${data}`;
      }

      return '';
    });

  const nonEmptyFilters = serializedFilters.filter((v) => !!v);
  return nonEmptyFilters.length === 0 ? '' : nonEmptyFilters.join(' AND ');
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
