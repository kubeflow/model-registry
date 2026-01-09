import {
  ModelCatalogStringFilterKey,
  ModelCatalogNumberFilterKey,
  UseCaseOptionValue,
  isLatencyMetricFieldName,
  DEFAULT_PERFORMANCE_FILTERS_QUERY_NAME,
} from '~/concepts/modelCatalog/const';
import {
  CatalogFilterOptionsList,
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
 * Applies a filter value to the filter state with proper type handling.
 * Handles string arrays (USE_CASE, HARDWARE_TYPE), numbers (MAX_RPS, latency), and clears.
 * This centralizes the type coercion logic for filter values.
 * Accepts string for filterKey to work with Object.entries().
 */
export const applyFilterValue = (
  setFilterData: SetFilterDataFn,
  filterKey: string,
  value: ModelCatalogFilterStates[keyof ModelCatalogFilterStates] | undefined,
): void => {
  if (filterKey === ModelCatalogStringFilterKey.USE_CASE) {
    if (Array.isArray(value)) {
      const validValues: UseCaseOptionValue[] = value
        .filter((v): v is string => typeof v === 'string')
        .filter(isUseCaseOptionValue);
      setFilterData(ModelCatalogStringFilterKey.USE_CASE, validValues);
    } else {
      setFilterData(ModelCatalogStringFilterKey.USE_CASE, []);
    }
  } else if (filterKey === ModelCatalogStringFilterKey.HARDWARE_TYPE) {
    if (Array.isArray(value)) {
      const validValues = value.filter((v): v is string => typeof v === 'string');
      setFilterData(ModelCatalogStringFilterKey.HARDWARE_TYPE, validValues);
    } else {
      setFilterData(ModelCatalogStringFilterKey.HARDWARE_TYPE, []);
    }
  } else if (filterKey === ModelCatalogNumberFilterKey.MAX_RPS) {
    setFilterData(
      ModelCatalogNumberFilterKey.MAX_RPS,
      typeof value === 'number' ? value : undefined,
    );
  } else if (isLatencyMetricFieldName(filterKey)) {
    setFilterData(filterKey, typeof value === 'number' ? value : undefined);
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
 * Parses filter values from a named query and returns them as a partial filter state.
 * This is the core function for extracting default filter values from namedQueries.
 * Filter keys in the frontend now match backend field names directly.
 */
export const getDefaultFiltersFromNamedQuery = (
  filterOptions: CatalogFilterOptionsList | null,
  namedQuery: NamedQuery,
): Partial<ModelCatalogFilterStates> => {
  const result: Partial<ModelCatalogFilterStates> = {};

  Object.entries(namedQuery).forEach(([fieldName, fieldFilter]) => {
    // Check which filter type this is using the enum values (which now match backend keys)
    const isUseCase = fieldName === ModelCatalogStringFilterKey.USE_CASE;
    const isHardwareType = fieldName === ModelCatalogStringFilterKey.HARDWARE_TYPE;
    const isRps = fieldName === ModelCatalogNumberFilterKey.MAX_RPS;
    const isLatencyField = isLatencyMetricFieldName(fieldName);

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
      result[ModelCatalogStringFilterKey.USE_CASE] = validValues;
    } else if (isHardwareType) {
      const values =
        fieldFilter.operator === FilterOperator.IN && Array.isArray(fieldFilter.value)
          ? fieldFilter.value.filter((v): v is string => typeof v === 'string')
          : typeof fieldFilter.value === 'string'
            ? [fieldFilter.value]
            : [];
      result[ModelCatalogStringFilterKey.HARDWARE_TYPE] = values;
    } else if (isRps) {
      const resolvedValue = resolveFilterValue(filterOptions, fieldName, fieldFilter.value);
      if (resolvedValue !== undefined) {
        result[ModelCatalogNumberFilterKey.MAX_RPS] = resolvedValue;
      }
    } else if (isLatencyField) {
      // Apply latency filter using the full filter key (e.g., 'artifacts.ttft_p90.double_value')
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
