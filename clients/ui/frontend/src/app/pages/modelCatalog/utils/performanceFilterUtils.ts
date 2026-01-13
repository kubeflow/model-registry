import {
  ModelCatalogStringFilterKey,
  ModelCatalogNumberFilterKey,
  UseCaseOptionValue,
  isLatencyFilterKey,
  isPerformanceStringFilterKey,
  isPerformanceNumberFilterKey,
  DEFAULT_PERFORMANCE_FILTERS_QUERY_NAME,
} from '~/concepts/modelCatalog/const';
import {
  CatalogFilterOptionsList,
  FieldFilter,
  FilterOperator,
  ModelCatalogFilterStates,
  NamedQuery,
} from '~/app/modelCatalogTypes';
import { isUseCaseOptionValue } from './workloadTypeUtils';

/**
 * Type for a function that sets filter data.
 */
export type SetFilterDataFn = <K extends keyof ModelCatalogFilterStates>(
  key: K,
  value: ModelCatalogFilterStates[K],
) => void;

/**
 * Extracts and validates string array values from a filter value.
 * Returns a properly typed array based on the filter key.
 */
const extractStringArrayValue = (
  value: ModelCatalogFilterStates[keyof ModelCatalogFilterStates] | undefined,
): string[] => {
  if (!Array.isArray(value)) {
    return [];
  }
  return value.filter((v): v is string => typeof v === 'string');
};

/**
 * Extracts and validates USE_CASE filter values.
 */
const extractUseCaseValues = (
  value: ModelCatalogFilterStates[keyof ModelCatalogFilterStates] | undefined,
): UseCaseOptionValue[] => extractStringArrayValue(value).filter(isUseCaseOptionValue);

/**
 * Extracts a number value from a filter value.
 */
const extractNumberValue = (
  value: ModelCatalogFilterStates[keyof ModelCatalogFilterStates] | undefined,
): number | undefined => (typeof value === 'number' ? value : undefined);

/**
 * Applies a filter value to the filter state with proper type handling.
 * Uses centralized filter type categorization from const.ts.
 * This centralizes the type coercion logic for filter values.
 * Accepts string for filterKey to work with Object.entries().
 *
 * To add a new performance filter:
 * 1. Add to PERFORMANCE_STRING_FILTER_KEYS or PERFORMANCE_NUMBER_FILTER_KEYS in const.ts
 * 2. For string filters with special validation (like USE_CASE), add an extraction function
 * 3. Add a case in the appropriate switch/if statement below
 */
export const applyFilterValue = (
  setFilterData: SetFilterDataFn,
  filterKey: string,
  value: ModelCatalogFilterStates[keyof ModelCatalogFilterStates] | undefined,
): void => {
  // Handle performance string filters (arrays of strings)
  if (isPerformanceStringFilterKey(filterKey)) {
    // Each string filter may have different validation, handle explicitly
    if (filterKey === ModelCatalogStringFilterKey.USE_CASE) {
      setFilterData(ModelCatalogStringFilterKey.USE_CASE, extractUseCaseValues(value));
    } else if (filterKey === ModelCatalogStringFilterKey.HARDWARE_TYPE) {
      setFilterData(ModelCatalogStringFilterKey.HARDWARE_TYPE, extractStringArrayValue(value));
    }
    // Future string filters: add else-if cases above
    return;
  }

  // Handle performance number filters
  // Currently only MAX_RPS, but structured for future extensibility
  if (isPerformanceNumberFilterKey(filterKey)) {
    const numberValue = extractNumberValue(value);
    // Use explicit key to ensure type safety (MAX_RPS is currently the only performance number filter)
    setFilterData(ModelCatalogNumberFilterKey.MAX_RPS, numberValue);
    return;
  }

  // Handle latency filters (also numbers)
  if (isLatencyFilterKey(filterKey)) {
    const numberValue = extractNumberValue(value);
    setFilterData(filterKey, numberValue);
  }
};

/**
 * Resolves a filter value, handling special values like 'max' or 'min'
 * that should be looked up from the filter options range.
 */
export const resolveFilterValue = (
  filterOptions: CatalogFilterOptionsList | null,
  fieldName: string,
  value: string | number | boolean | (string | number)[],
): number | undefined => {
  if (value === 'max' || value === 'min') {
    // Look up the range from filterOptions using Object.entries to find the matching field
    const filters = filterOptions?.filters;
    if (filters) {
      const entries = Object.entries(filters);
      const matchingEntry = entries.find(([key]) => key === fieldName);
      if (matchingEntry) {
        const [, filterOption] = matchingEntry;
        if ('range' in filterOption && filterOption.range) {
          return value === 'max' ? filterOption.range.max : filterOption.range.min;
        }
      }
    }
    return undefined;
  }
  if (typeof value === 'number') {
    return value;
  }
  return undefined;
};

/**
 * Extracts string array values from a FieldFilter for namedQuery parsing.
 * Handles both IN operator (array) and single string values.
 */
const extractStringArrayFromFieldFilter = (fieldFilter: FieldFilter): string[] => {
  const rawValues =
    fieldFilter.operator === FilterOperator.IN && Array.isArray(fieldFilter.value)
      ? fieldFilter.value.filter((v): v is string => typeof v === 'string')
      : typeof fieldFilter.value === 'string'
        ? [fieldFilter.value]
        : [];
  return rawValues;
};

/**
 * Extracts USE_CASE values from a FieldFilter with validation.
 */
const extractUseCaseValuesFromFieldFilter = (fieldFilter: FieldFilter): UseCaseOptionValue[] =>
  extractStringArrayFromFieldFilter(fieldFilter).filter(isUseCaseOptionValue);

/**
 * Parses filter values from a named query and returns them as a partial filter state.
 * Uses centralized filter type categorization from const.ts.
 * Filter keys in the frontend now match backend field names directly.
 *
 * To add a new performance filter:
 * 1. Add to PERFORMANCE_STRING_FILTER_KEYS or PERFORMANCE_NUMBER_FILTER_KEYS in const.ts
 * 2. For string filters with special validation (like USE_CASE), add an extraction function
 * 3. Add a case in the appropriate if/else statement below
 */
export const getDefaultFiltersFromNamedQuery = (
  filterOptions: CatalogFilterOptionsList | null,
  namedQuery: NamedQuery,
): Partial<ModelCatalogFilterStates> => {
  const result: Partial<ModelCatalogFilterStates> = {};

  Object.entries(namedQuery).forEach(([fieldName, fieldFilter]) => {
    // Handle performance string filters (arrays of strings)
    if (isPerformanceStringFilterKey(fieldName)) {
      // Each string filter may have different validation, handle explicitly
      if (fieldName === ModelCatalogStringFilterKey.USE_CASE) {
        result[ModelCatalogStringFilterKey.USE_CASE] =
          extractUseCaseValuesFromFieldFilter(fieldFilter);
      } else if (fieldName === ModelCatalogStringFilterKey.HARDWARE_TYPE) {
        result[ModelCatalogStringFilterKey.HARDWARE_TYPE] =
          extractStringArrayFromFieldFilter(fieldFilter);
      }
      // Future string filters: add else-if cases above
      return;
    }

    // Handle performance number filters
    // Currently only MAX_RPS, but structured for future extensibility
    if (isPerformanceNumberFilterKey(fieldName)) {
      const resolvedValue = resolveFilterValue(filterOptions, fieldName, fieldFilter.value);
      if (resolvedValue !== undefined) {
        // Use explicit key to ensure type safety (MAX_RPS is currently the only performance number filter)
        result[ModelCatalogNumberFilterKey.MAX_RPS] = resolvedValue;
      }
      return;
    }

    // Handle latency filters (also numbers)
    if (isLatencyFilterKey(fieldName)) {
      const resolvedValue = resolveFilterValue(filterOptions, fieldName, fieldFilter.value);
      if (resolvedValue !== undefined) {
        result[fieldName] = resolvedValue;
      }
    }
  });

  return result;
};

/**
 * Gets all default performance filter values from namedQueries.
 * Returns a partial filter state with all default values.
 */
export const getDefaultPerformanceFilters = (
  filterOptions: CatalogFilterOptionsList | null,
): Partial<ModelCatalogFilterStates> => {
  const defaultQuery = filterOptions?.namedQueries?.[DEFAULT_PERFORMANCE_FILTERS_QUERY_NAME];
  if (!defaultQuery) {
    return {};
  }
  return getDefaultFiltersFromNamedQuery(filterOptions, defaultQuery);
};

/**
 * Gets the default value for a single performance filter from namedQueries.
 * Returns the value and whether a default was found.
 */
export const getSingleFilterDefault = (
  filterOptions: CatalogFilterOptionsList | null,
  filterKey: keyof ModelCatalogFilterStates,
): { hasDefault: boolean; value: ModelCatalogFilterStates[keyof ModelCatalogFilterStates] } => {
  const defaults = getDefaultPerformanceFilters(filterOptions);
  const value = defaults[filterKey];
  return {
    hasDefault: value !== undefined,
    value,
  };
};
