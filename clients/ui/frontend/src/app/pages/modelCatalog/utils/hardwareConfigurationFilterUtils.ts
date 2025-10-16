import {
  CatalogPerformanceMetricsArtifact,
  ModelCatalogFilterStates,
} from '~/app/modelCatalogTypes';
import { getDoubleValue, getStringValue } from '~/app/utils';
import {
  ModelCatalogStringFilterKey,
  ModelCatalogNumberFilterKey,
  LatencyMetricFieldName,
  LatencyPercentile,
  LatencyMetric,
} from '~/concepts/modelCatalog/const';
import { getTotalRps } from './performanceMetricsUtils';
import { LatencyFilterConfig as SharedLatencyFilterConfig } from './latencyFilterState';

// Type for storing complex latency filter configuration with value
export type LatencyFilterConfig = SharedLatencyFilterConfig & {
  value: number;
};

const isMetricLowercase = (str: string): str is Lowercase<LatencyMetric> =>
  Object.values(LatencyMetric)
    .map((value) => value.toLowerCase())
    .includes(str);
const isPercentileLowercase = (str: string): str is Lowercase<LatencyPercentile> =>
  Object.values(LatencyPercentile)
    .map((value) => value.toLowerCase())
    .includes(str);

/**
 * Maps metric and percentile combination to the corresponding artifact field
 */
export const getLatencyFieldName = (
  metric: LatencyMetric,
  percentile: LatencyPercentile,
): LatencyMetricFieldName => {
  const metricPrefix = metric.toLowerCase();
  const percentileSuffix = percentile.toLowerCase();
  if (!isMetricLowercase(metricPrefix) || !isPercentileLowercase(percentileSuffix)) {
    return 'ttft_mean'; // Default fallback
  }
  return `${metricPrefix}_${percentileSuffix}`;
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
 * Enhanced filter for Max Latency that supports metric and percentile selection
 */
export const applyMaxLatencyFilter = (
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
  latencyConfig?: SharedLatencyFilterConfig,
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
      const totalRps = getTotalRps(artifact.customProperties);
      if (totalRps < minRpsFilter) {
        return false;
      }
    }

    // Max Latency Filter - use provided config or fall back to default
    const maxLatencyFilter = filterState[ModelCatalogNumberFilterKey.MAX_LATENCY];
    if (maxLatencyFilter !== undefined) {
      const baseConfig = latencyConfig || {
        metric: LatencyMetric.TTFT,
        percentile: LatencyPercentile.Mean,
      };

      // Create the full config with the current filter value
      const fullConfig: LatencyFilterConfig = { ...baseConfig, value: maxLatencyFilter };

      if (!applyMaxLatencyFilter(artifact, fullConfig)) {
        return false;
      }
    }

    // Use Case Filter
    const useCaseFilter = filterState[ModelCatalogStringFilterKey.USE_CASE];

    if (useCaseFilter) {
      // Get the artifact's use case
      const artifactUseCase = getStringValue(artifact.customProperties, 'use_case');

      // Check if the artifact's use case matches the selected use case
      // Use includes() to handle potential comma-separated values or partial matches
      if (!artifactUseCase || !artifactUseCase.includes(useCaseFilter)) {
        return false;
      }
    }

    return true;
  });

/**
 * Clears all active filters
 */
export const clearAllFilters = (
  setFilterData: <K extends keyof ModelCatalogFilterStates>(
    key: K,
    value: ModelCatalogFilterStates[K],
  ) => void,
): void => {
  // Clear string filters (arrays)
  setFilterData(ModelCatalogStringFilterKey.TASK, []);
  setFilterData(ModelCatalogStringFilterKey.PROVIDER, []);
  setFilterData(ModelCatalogStringFilterKey.LICENSE, []);
  setFilterData(ModelCatalogStringFilterKey.LANGUAGE, []);
  setFilterData(ModelCatalogStringFilterKey.HARDWARE_TYPE, []);

  // Clear use case filter (single value)
  setFilterData(ModelCatalogStringFilterKey.USE_CASE, undefined);

  // Clear number filters
  Object.values(ModelCatalogNumberFilterKey).forEach((key) => {
    setFilterData(key, undefined);
  });
};
