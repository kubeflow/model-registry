import {
  ModelCatalogStringFilterKey,
  ModelCatalogNumberFilterKey,
  UseCaseOptionValue,
  isLatencyMetricFieldName,
  ALL_LATENCY_FIELD_NAMES,
  LatencyMetricFieldName,
} from '~/concepts/modelCatalog/const';
import {
  CatalogFilterOptionsList,
  FilterOperator,
  ModelCatalogFilterStates,
  NamedQuery,
} from '~/app/modelCatalogTypes';
import { isUseCaseOptionValue } from './workloadTypeUtils';

/**
 * Maps a frontend filter key to its corresponding backend field name in namedQueries.
 */
export const getBackendFieldName = (filterKey: keyof ModelCatalogFilterStates): string | null => {
  const filterKeyStr = String(filterKey);

  if (filterKey === ModelCatalogStringFilterKey.USE_CASE) {
    return 'artifacts.use_case.string_value';
  }
  if (filterKey === ModelCatalogStringFilterKey.HARDWARE_TYPE) {
    return 'artifacts.hardware_type.string_value';
  }
  if (filterKey === ModelCatalogNumberFilterKey.MAX_RPS) {
    return 'artifacts.requests_per_second.double_value';
  }
  if (isLatencyMetricFieldName(filterKeyStr)) {
    return `artifacts.${filterKeyStr}.double_value`;
  }
  return null;
};

/**
 * Checks if a field name corresponds to a latency metric field.
 * Handles both prefixed (artifacts.ttft_p90.double_value) and unprefixed (ttft_p90) formats.
 */
export const getLatencyFieldKey = (fieldName: string): LatencyMetricFieldName | null => {
  // Check if it's a direct latency field name
  const directMatch = ALL_LATENCY_FIELD_NAMES.find((name) => name === fieldName);
  if (directMatch) {
    return directMatch;
  }
  // Check if it's an artifacts.* prefixed latency field
  const artifactMatch = ALL_LATENCY_FIELD_NAMES.find(
    (name) => fieldName === `artifacts.${name}.double_value`,
  );
  if (artifactMatch) {
    return artifactMatch;
  }
  return null;
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
    // Look up the range from filterOptions
    const filterOption = filterOptions?.filters?.[fieldName];
    if (filterOption && 'range' in filterOption && filterOption.range) {
      return value === 'max' ? filterOption.range.max : filterOption.range.min;
    }
    return undefined;
  }
  if (typeof value === 'number') {
    return value;
  }
  return undefined;
};

/**
 * Gets the default value for a performance filter from namedQueries.
 * Used to determine if a filter chip should be shown (only when value differs from default).
 */
export const getPerformanceFilterDefaultValue = (
  filterOptions: CatalogFilterOptionsList | null,
  filterKey: keyof ModelCatalogFilterStates,
): string | number | string[] | undefined => {
  const defaultQuery = filterOptions?.namedQueries?.['default-performance-filters'];
  if (!defaultQuery) {
    return undefined;
  }

  // Check if it's a latency field first (before type narrowing)
  const filterKeyStr = String(filterKey);
  const isLatencyField = isLatencyMetricFieldName(filterKeyStr);

  // Map frontend filter keys to their corresponding backend field names in namedQueries
  let backendFieldName: string | null = null;
  if (filterKey === ModelCatalogStringFilterKey.USE_CASE) {
    backendFieldName = 'artifacts.use_case.string_value';
  } else if (filterKey === ModelCatalogStringFilterKey.HARDWARE_TYPE) {
    backendFieldName = 'artifacts.hardware_type.string_value';
  } else if (filterKey === ModelCatalogNumberFilterKey.MAX_RPS) {
    backendFieldName = 'artifacts.requests_per_second.double_value';
  } else if (isLatencyField) {
    backendFieldName = `artifacts.${filterKeyStr}.double_value`;
  }

  // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
  if (!backendFieldName) {
    return undefined;
  }

  const fieldFilter = defaultQuery[backendFieldName];
  // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
  if (!fieldFilter) {
    return undefined;
  }

  if (filterKey === ModelCatalogStringFilterKey.USE_CASE) {
    const rawValues =
      fieldFilter.operator === FilterOperator.IN && Array.isArray(fieldFilter.value)
        ? fieldFilter.value.filter((v): v is string => typeof v === 'string')
        : typeof fieldFilter.value === 'string'
          ? [fieldFilter.value]
          : [];
    const validValues: UseCaseOptionValue[] = rawValues.filter(isUseCaseOptionValue);
    return validValues.length === 1 ? validValues[0] : validValues;
  }

  if (filterKey === ModelCatalogStringFilterKey.HARDWARE_TYPE) {
    const values =
      fieldFilter.operator === FilterOperator.IN && Array.isArray(fieldFilter.value)
        ? fieldFilter.value.filter((v): v is string => typeof v === 'string')
        : typeof fieldFilter.value === 'string'
          ? [fieldFilter.value]
          : [];
    return values;
  }

  if (filterKey === ModelCatalogNumberFilterKey.MAX_RPS || isLatencyField) {
    return resolveFilterValue(filterOptions, backendFieldName, fieldFilter.value);
  }

  return undefined;
};

/**
 * Applies filter values from a named query to the filter state.
 * Returns an object with the filter values to apply.
 */
export const getDefaultFiltersFromNamedQuery = (
  filterOptions: CatalogFilterOptionsList | null,
  namedQuery: NamedQuery,
): Partial<ModelCatalogFilterStates> => {
  const result: Partial<ModelCatalogFilterStates> = {};

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

    // Check if it's a latency field by extracting the field name from artifacts.*.double_value
    const latencyMatch = fieldName.match(/^artifacts\.([a-z0-9_]+)\.double_value$/);
    const potentialLatencyField = latencyMatch?.[1];
    const isLatencyField = potentialLatencyField && isLatencyMetricFieldName(potentialLatencyField);

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
      // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
    } else if (isLatencyField && potentialLatencyField) {
      // Apply latency filter using the resolved field name
      const resolvedValue = resolveFilterValue(filterOptions, fieldName, fieldFilter.value);
      if (resolvedValue !== undefined) {
        result[potentialLatencyField] = resolvedValue;
      }
    }
  });

  return result;
};

/**
 * Type representing the result of getting a single filter's default value.
 * Includes the value and whether a default was found.
 */
export type SingleFilterDefaultResult = {
  hasDefault: boolean;
  value: UseCaseOptionValue[] | string[] | number | undefined;
};

/**
 * Gets the default value for a single performance filter from namedQueries.
 * Returns both the value and whether a default was found.
 * This is used when resetting individual filter chips.
 */
export const getSingleFilterDefault = (
  filterOptions: CatalogFilterOptionsList | null,
  filterKey: keyof ModelCatalogFilterStates,
): SingleFilterDefaultResult => {
  const defaultQuery = filterOptions?.namedQueries?.['default-performance-filters'];
  const backendFieldName = getBackendFieldName(filterKey);

  if (!defaultQuery || !backendFieldName) {
    return { hasDefault: false, value: undefined };
  }

  const fieldFilter = defaultQuery[backendFieldName];
  // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
  if (!fieldFilter) {
    return { hasDefault: false, value: undefined };
  }

  const filterKeyStr = String(filterKey);

  if (filterKey === ModelCatalogStringFilterKey.USE_CASE) {
    const rawValues =
      fieldFilter.operator === FilterOperator.IN && Array.isArray(fieldFilter.value)
        ? fieldFilter.value.filter((v): v is string => typeof v === 'string')
        : typeof fieldFilter.value === 'string'
          ? [fieldFilter.value]
          : [];
    const validValues: UseCaseOptionValue[] = rawValues.filter(isUseCaseOptionValue);
    return { hasDefault: true, value: validValues };
  }

  if (filterKey === ModelCatalogStringFilterKey.HARDWARE_TYPE) {
    const values =
      fieldFilter.operator === FilterOperator.IN && Array.isArray(fieldFilter.value)
        ? fieldFilter.value.filter((v): v is string => typeof v === 'string')
        : typeof fieldFilter.value === 'string'
          ? [fieldFilter.value]
          : [];
    return { hasDefault: true, value: values };
  }

  if (filterKey === ModelCatalogNumberFilterKey.MAX_RPS || isLatencyMetricFieldName(filterKeyStr)) {
    const resolvedValue = resolveFilterValue(filterOptions, backendFieldName, fieldFilter.value);
    return { hasDefault: true, value: resolvedValue };
  }

  return { hasDefault: false, value: undefined };
};
