import { isEnumMember } from 'mod-arch-core';
import {
  CatalogPerformanceMetricsArtifact,
  ModelCatalogFilterStates,
  ModelCatalogFilterKey,
} from '~/app/modelCatalogTypes';
import { getDoubleValue, getStringValue } from '~/app/utils';
import {
  ModelCatalogStringFilterKey,
  ModelCatalogNumberFilterKey,
  LatencyMetricFieldName,
  LatencyPercentile,
  LatencyMetric,
  ALL_LATENCY_FIELD_NAMES,
  getLatencyFieldName,
} from '~/concepts/modelCatalog/const';

// Type for storing complex latency filter configuration with value
export type LatencyFilterConfig = {
  metric: LatencyMetric;
  percentile: LatencyPercentile;
  value: number;
};

/**
 * Inverse of getLatencyFieldName
 */
export const parseLatencyFieldName = (
  fieldName: LatencyMetricFieldName,
): { metric: LatencyMetric; percentile: LatencyPercentile } | null => {
  const [prefix, suffix] = fieldName.split('_');
  const metric = Object.values(LatencyMetric).find((m) => m.toLowerCase() === prefix);
  const percentile = Object.values(LatencyPercentile).find((p) => p.toLowerCase() === suffix);
  return metric && percentile ? { metric, percentile } : null;
};

/**
 * Extracts unique hardware types from performance artifacts
 */
export const getUniqueHardwareTypes = (
  artifacts: CatalogPerformanceMetricsArtifact[],
): string[] => {
  const hardwareTypes = artifacts
    .map((artifact) =>
      getStringValue(artifact.customProperties, ModelCatalogStringFilterKey.HARDWARE_TYPE),
    )
    .filter((hardware): hardware is string => !!hardware);

  return [...new Set(hardwareTypes)].toSorted();
};

/**
 * Enhanced filter for Latency that supports metric and percentile selection
 */
export const applyLatencyFilter = (
  artifact: CatalogPerformanceMetricsArtifact,
  config: LatencyFilterConfig,
): boolean => {
  const fieldName = getLatencyFieldName(config.metric, config.percentile);
  const latencyValue = getDoubleValue(artifact.customProperties, fieldName);
  return latencyValue <= config.value;
};

/**
 * Filters hardware configuration artifacts based on current filter state
 */
export const filterHardwareConfigurationArtifacts = (
  artifacts: CatalogPerformanceMetricsArtifact[],
  filterState: ModelCatalogFilterStates,
): CatalogPerformanceMetricsArtifact[] =>
  artifacts.filter((artifact) => {
    // Hardware Type Filter (using central filter state)
    const hardwareTypeFilters = filterState[ModelCatalogStringFilterKey.HARDWARE_TYPE];
    if (hardwareTypeFilters.length > 0) {
      const hardwareType = getStringValue(
        artifact.customProperties,
        ModelCatalogStringFilterKey.HARDWARE_TYPE,
      );
      if (!hardwareType || !hardwareTypeFilters.includes(hardwareType)) {
        return false;
      }
    }

    // Min RPS Filter
    const minRpsFilter = filterState[ModelCatalogNumberFilterKey.MIN_RPS];
    if (minRpsFilter !== undefined) {
      const rpsPerReplica = getDoubleValue(artifact.customProperties, 'requests_per_second');
      if (rpsPerReplica < minRpsFilter) {
        return false;
      }
    }

    // Max Latency Filter - check for any active latency field
    for (const metric of Object.values(LatencyMetric)) {
      for (const percentile of Object.values(LatencyPercentile)) {
        const fieldName = getLatencyFieldName(metric, percentile);
        const filterValue = filterState[fieldName];
        if (filterValue !== undefined && typeof filterValue === 'number') {
          const latencyValue = getDoubleValue(artifact.customProperties, fieldName);
          if (latencyValue > filterValue) {
            return false;
          }
        }
      }
    }

    // Use Case Filter
    const useCaseFilters = filterState[ModelCatalogStringFilterKey.USE_CASE];

    if (useCaseFilters.length > 0) {
      // Get the artifact's use case
      const artifactUseCase = getStringValue(artifact.customProperties, 'use_case');

      // Check if the artifact's use case matches any of the selected use cases (exact match)
      if (!artifactUseCase || !useCaseFilters.some((filter) => filter === artifactUseCase)) {
        return false;
      }
    }

    return true;
  });

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
