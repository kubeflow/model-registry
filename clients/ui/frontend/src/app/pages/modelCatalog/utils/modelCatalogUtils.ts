import { isEnumMember } from 'mod-arch-core';
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
  MetricsType,
} from '~/app/modelCatalogTypes';
import { getLabels } from '~/app/pages/modelRegistry/screens/utils';
import {
  ModelCatalogStringFilterKey,
  ModelCatalogNumberFilterKey,
  ModelCatalogTask,
  ModelCatalogProvider,
  ModelCatalogLicense,
  AllLanguageCode,
} from '~/concepts/modelCatalog/const';

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

  const filteredItems = catalogSources.items.filter((source) => source.enabled !== false);

  return {
    ...catalogSources,
    items: filteredItems,
    size: filteredItems.length,
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

export const filterArtifactsByType = <T extends CatalogArtifacts>(
  artifacts: CatalogArtifacts[],
  artifactType: CatalogArtifactType,
  metricsType?: MetricsType,
): T[] =>
  artifacts.filter((artifact): artifact is T => {
    if (artifact.artifactType !== artifactType) {
      return false;
    }
    if (metricsType && 'metricsType' in artifact) {
      return artifact.metricsType === metricsType;
    }
    return true;
  });

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

// Define array-based filter keys (excluding USE_CASE which is single-selection)
type ArrayFilterKey =
  | ModelCatalogStringFilterKey.TASK
  | ModelCatalogStringFilterKey.PROVIDER
  | ModelCatalogStringFilterKey.LICENSE
  | ModelCatalogStringFilterKey.LANGUAGE
  | ModelCatalogStringFilterKey.HARDWARE_TYPE;

// Type mapping for array filter values
type ArrayFilterValueType = {
  [ModelCatalogStringFilterKey.TASK]: ModelCatalogTask;
  [ModelCatalogStringFilterKey.PROVIDER]: ModelCatalogProvider;
  [ModelCatalogStringFilterKey.LICENSE]: ModelCatalogLicense;
  [ModelCatalogStringFilterKey.LANGUAGE]: AllLanguageCode;
  [ModelCatalogStringFilterKey.HARDWARE_TYPE]: string;
};

// Type guard to check if a value is an array of the expected type
const isArrayOfValues = <T>(value: unknown): value is T[] => Array.isArray(value);

// Type guard to check if filter key is valid for array operations
const isArrayFilterKey = (filterKey: string): filterKey is ArrayFilterKey =>
  isEnumMember(filterKey, ModelCatalogStringFilterKey) &&
  filterKey !== ModelCatalogStringFilterKey.USE_CASE;

export const useCatalogStringFilterState = <K extends ArrayFilterKey>(
  filterKey: K,
): {
  isSelected: (value: ArrayFilterValueType[K]) => boolean;
  setSelected: (value: ArrayFilterValueType[K], selected: boolean) => void;
} => {
  type Value = ArrayFilterValueType[K];
  const { filterData, setFilterData } = React.useContext(ModelCatalogContext);
  const selections = filterData[filterKey];

  const isSelected = React.useCallback(
    (value: Value) => {
      if (!isArrayOfValues<Value>(selections)) {
        return false;
      }
      return selections.includes(value);
    },
    [selections],
  );

  const setSelected = (value: Value, selected: boolean) => {
    if (!isArrayOfValues<Value>(selections)) {
      return;
    }

    const nextState: Value[] = selected
      ? [...selections, value]
      : selections.filter((item) => item !== value);

    if (isArrayFilterKey(filterKey)) {
      // Type assertion is safe here because we've verified the key is an array filter
      // eslint-disable-next-line @typescript-eslint/consistent-type-assertions
      setFilterData(filterKey, nextState as ModelCatalogFilterStates[K]);
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

// TODO: Implement performance filters.
// type FilterId = keyof CatalogFilterOptionsList['filters'];
// const KNOWN_LESS_THAN_IDS: FilterId[] = [ModelCatalogNumberFilterKey.TTFT_MEAN]; // TODO: populate with filters that need to talk about "less" values
// const isKnownLessThanValue = (
//   filterOption: CatalogFilterOptions[keyof CatalogFilterOptions],
//   filterId: FilterId,
//   data: unknown,
// ): data is number =>
//   filterOption.type === 'number' &&
//   KNOWN_LESS_THAN_IDS.includes(filterId) &&
//   typeof data === 'number';

// const KNOWN_MORE_THAN_IDS: FilterId[] = [ModelCatalogNumberFilterKey.RPS_MEAN]; // TODO: populate with filters that need to talk about "more" values
// const isKnownMoreThanValue = (
//   filterOption: CatalogFilterOptions[keyof CatalogFilterOptions],
//   filterId: FilterId,
//   data: unknown,
// ): data is number =>
//   filterOption.type === 'number' &&
//   KNOWN_MORE_THAN_IDS.includes(filterId) &&
//   typeof data === 'number';

const isFilterIdInMap = (
  filterId: unknown,
  filters: CatalogFilterOptions,
): filterId is keyof CatalogFilterOptions => typeof filterId === 'string' && filterId in filters;

// TODO tech debt: different filterQuery syntax is needed depending on whether the API stores an array of values or a single string value.
//   the current filter_options API response does not indicate the difference between these two types of fields, so for now we hard-code them.
const KNOWN_ARRAY_FILTER_IDS: (keyof CatalogFilterOptions)[] = [
  ModelCatalogStringFilterKey.LANGUAGE,
  ModelCatalogStringFilterKey.TASK,
];

// If using LIKE on an array field, we need %" "% around value within the ' '
const wrapInQuotes = (v: string, isArrayLikeFilter = false): string =>
  isArrayLikeFilter ? `'%"${v}"%'` : `'${v}'`;

// LIKE works for any string filter but is only required for array fields
const likeFilter = (k: string, v: string, isArrayField: boolean): string =>
  `${k} LIKE ${wrapInQuotes(v, isArrayField)}`;

// = and IN only work for non-array fields
const eqFilter = (k: string, v: string) => `${k}=${wrapInQuotes(v)}`;
const inFilter = (k: string, values: string[]) =>
  `${k} IN (${values.map((v) => wrapInQuotes(v)).join(',')})`;

export const filtersToFilterQuery = (
  filterData: ModelCatalogFilterStates,
  options: CatalogFilterOptionsList,
): string => {
  const serializedFilters: string[] = Object.entries(filterData).map(([filterId, data]) => {
    if (!isFilterIdInMap(filterId, options.filters) || typeof data === 'undefined') {
      // Unhandled key or no data
      return '';
    }

    const filterOption = options.filters[filterId];
    const isArrayField = KNOWN_ARRAY_FILTER_IDS.includes(filterId);

    if (isArrayOfSelections(filterOption, data)) {
      switch (data.length) {
        case 0:
          return '';
        case 1:
          if (isArrayField) {
            return likeFilter(filterId, data[0], true);
          }
          return eqFilter(filterId, data[0]);
        default:
          // 2 or more
          if (isArrayField) {
            return `(${data.map((value) => likeFilter(filterId, value, true)).join(' OR ')})`;
          }
          return inFilter(filterId, data);
      }
    }

    // Handle single string values (like USE_CASE)
    if (filterOption?.type === 'string' && typeof data === 'string') {
      return `${filterId}=${wrapInQuotes(data)}`;
    }

    // TODO: Implement performance filters.
    // if (isKnownLessThanValue(filterOption, filterId, data)) {
    //   return `${filterId} < ${data}`;
    // }

    // if (isKnownMoreThanValue(filterOption, filterId, data)) {
    //   return `${filterId} > ${data}`;
    // }

    // TODO: Implement more data transforms
    // Shouldn't reach this far, but if it does, log & ignore the case
    // eslint-disable-next-line no-console
    console.warn('Unhandled option', filterId, data, filterOption);
    return '';
  });

  const nonEmptyFilters = serializedFilters.filter((v) => !!v);

  // eg. filterQuery=rps_mean >1 AND license IN ('mit','apache-2.0') AND ttft_mean < 10
  return nonEmptyFilters.length === 0 ? '' : nonEmptyFilters.join(' AND ');
};

export const getUniqueSourceLabels = (catalogSources: CatalogSourceList | null): string[] => {
  if (!catalogSources) {
    return [];
  }

  const allLabels = new Set<string>();

  catalogSources.items.forEach((source) => {
    if (source.enabled && source.labels.length > 0) {
      source.labels.forEach((label) => {
        if (label.trim()) {
          allLabels.add(label.trim());
        }
      });
    }
  });

  return Array.from(allLabels);
};

export const getSourceFromSourceId = (
  sourceId: string,
  catalogSources: CatalogSourceList | null,
): CatalogSource | undefined => {
  if (!catalogSources || !sourceId) {
    return undefined;
  }

  return catalogSources.items.find((source) => source.id === sourceId);
};

export const hasFiltersApplied = (filterData: ModelCatalogFilterStates): boolean =>
  Object.values(filterData).some((value) => {
    if (Array.isArray(value)) {
      return value.length > 0;
    }
    return value !== undefined;
  });
