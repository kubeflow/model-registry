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

// Type for storing complex latency filter configuration
export type LatencyFilterConfig = {
  metric: LatencyMetric;
  percentile: LatencyPercentile;
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

    // Max Latency Filter (using MAX_LATENCY key for now)
    // TODO: This currently uses TTFT Mean as default - should be enhanced to use
    // the specific metric/percentile selected by the user
    const maxLatencyFilter = filterState[ModelCatalogNumberFilterKey.MAX_LATENCY];
    if (maxLatencyFilter !== undefined) {
      const defaultConfig: LatencyFilterConfig = {
        metric: LatencyMetric.TTFT,
        percentile: LatencyPercentile.Mean,
        value: maxLatencyFilter,
      };

      if (!applyMaxLatencyFilter(artifact, defaultConfig)) {
        return false;
      }
    }

    // Workload Type Filter (based on max input/output tokens as minimum thresholds)
    const maxInputTokensFilter = filterState[ModelCatalogNumberFilterKey.MAX_INPUT_TOKENS];
    const maxOutputTokensFilter = filterState[ModelCatalogNumberFilterKey.MAX_OUTPUT_TOKENS];

    if (maxInputTokensFilter !== undefined && maxOutputTokensFilter !== undefined) {
      // Get the artifact's max input/output token capabilities
      const artifactMaxInputTokens = getDoubleValue(artifact.customProperties, 'max_input_tokens');
      const artifactMaxOutputTokens = getDoubleValue(
        artifact.customProperties,
        'max_output_tokens',
      );

      // Apply minimum threshold logic: artifact must support AT LEAST the selected workload requirements
      if (
        artifactMaxInputTokens < maxInputTokensFilter ||
        artifactMaxOutputTokens < maxOutputTokensFilter
      ) {
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
  // Clear string filters
  setFilterData(ModelCatalogStringFilterKey.HARDWARE_TYPE, []);

  // Clear number filters
  Object.values(ModelCatalogNumberFilterKey).forEach((key) => {
    setFilterData(key, undefined);
  });
};
