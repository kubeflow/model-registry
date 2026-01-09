import { isEnumMember } from 'mod-arch-core';
import {
  CatalogPerformanceMetricsArtifact,
  ModelCatalogFilterStates,
  ModelCatalogFilterKey,
} from '~/app/modelCatalogTypes';
import { getStringValue } from '~/app/utils';
import {
  ModelCatalogStringFilterKey,
  ModelCatalogNumberFilterKey,
  LatencyMetricFieldName,
  ALL_LATENCY_FIELD_NAMES,
  parseLatencyFilterKey,
  PerformancePropertyKey,
} from '~/concepts/modelCatalog/const';

// Re-export parseLatencyFilterKey as parseLatencyFieldName for backward compatibility
export const parseLatencyFieldName = parseLatencyFilterKey;

/**
 * Extracts unique hardware types from performance artifacts
 */
export const getUniqueHardwareTypes = (
  artifacts: CatalogPerformanceMetricsArtifact[],
): string[] => {
  // Use the short property key for accessing customProperties
  const hardwareTypes = artifacts
    .map((artifact) =>
      getStringValue(artifact.customProperties, PerformancePropertyKey.HARDWARE_TYPE),
    )
    .filter((hardware): hardware is string => !!hardware);

  return [...new Set(hardwareTypes)].toSorted();
};

/**
 * Gets all filter keys (string filters + number filters + latency filters)
 */
export const getAllFilterKeys = (): {
  stringFilterKeys: ModelCatalogStringFilterKey[];
  numberFilterKeys: ModelCatalogNumberFilterKey[];
  latencyFilterKeys: LatencyMetricFieldName[];
} => ({
  stringFilterKeys: Object.values(ModelCatalogStringFilterKey),
  numberFilterKeys: Object.values(ModelCatalogNumberFilterKey),
  latencyFilterKeys: ALL_LATENCY_FIELD_NAMES,
});

/**
 * Clears filters. If filterKeys is provided, only clears those specific filters.
 * Otherwise clears all filters.
 */
export const clearAllFilters = (
  setFilterData: <K extends keyof ModelCatalogFilterStates>(
    key: K,
    value: ModelCatalogFilterStates[K],
  ) => void,
  filterKeys?: ModelCatalogFilterKey[],
): void => {
  const { stringFilterKeys, numberFilterKeys, latencyFilterKeys } = getAllFilterKeys();

  // If specific filter keys are provided, only clear those
  if (filterKeys) {
    filterKeys.forEach((key) => {
      if (isEnumMember(key, ModelCatalogStringFilterKey)) {
        setFilterData(key, []);
      } else {
        setFilterData(key, undefined);
      }
    });
    return;
  }

  // Clear all string filters (arrays)
  stringFilterKeys.forEach((key) => {
    setFilterData(key, []);
  });

  // Clear all number filters
  numberFilterKeys.forEach((key) => {
    setFilterData(key, undefined);
  });

  // Clear all latency metric filters (e.g., ttft_mean, ttft_p90, etc.)
  latencyFilterKeys.forEach((fieldName) => {
    setFilterData(fieldName, undefined);
  });
};
