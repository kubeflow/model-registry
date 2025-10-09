import {
  CatalogPerformanceMetricsArtifact,
  ModelCatalogFilterStates,
  PerformanceMetricsCustomProperties,
} from '~/app/modelCatalogTypes';
import { getDoubleValue, getStringValue } from '~/app/utils';
import {
  ModelCatalogStringFilterKey,
  ModelCatalogNumberFilterKey,
} from '~/concepts/modelCatalog/const';
import { getTotalRps } from './performanceMetricsUtils';

// Type for storing complex latency filter configuration
export type LatencyFilterConfig = {
  metric: 'E2E' | 'TTFT' | 'TPS' | 'ITL';
  percentile: 'Mean' | 'P90' | 'P95' | 'P99';
  value: number;
};

/**
 * Maps metric and percentile combination to the corresponding artifact field
 */
const getLatencyFieldName = (
  metric: string,
  percentile: string,
): keyof PerformanceMetricsCustomProperties => {
  const metricPrefix = metric.toLowerCase();
  const percentileSuffix = percentile === 'Mean' ? '_mean' : `_${percentile.toLowerCase()}`;
  const fieldName = `${metricPrefix}${percentileSuffix}`;

  // Validate that the field exists in PerformanceMetricsCustomProperties
  const validFields = [
    'ttft_mean',
    'ttft_p90',
    'ttft_p95',
    'ttft_p99',
    'e2e_mean',
    'e2e_p90',
    'e2e_p95',
    'e2e_p99',
    'tps_mean',
    'tps_p90',
    'tps_p95',
    'tps_p99',
    'itl_mean',
    'itl_p90',
    'itl_p95',
    'itl_p99',
  ];

  return validFields.includes(fieldName)
    ? // eslint-disable-next-line @typescript-eslint/consistent-type-assertions
      (fieldName as keyof PerformanceMetricsCustomProperties)
    : 'ttft_mean'; // Default fallback
};

/**
 * Extracts unique hardware types from performance artifacts
 */
export const getUniqueHardwareTypes = (
  artifacts: CatalogPerformanceMetricsArtifact[],
): string[] => {
  const hardwareTypes = artifacts
    .map((artifact) => getStringValue(artifact.customProperties, 'hardware'))
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
  appliedHardwareTypes?: string[],
): CatalogPerformanceMetricsArtifact[] =>
  artifacts.filter((artifact) => {
    // Hardware Type Filter (using dedicated hardware filter state)
    if (appliedHardwareTypes && appliedHardwareTypes.length > 0) {
      const hardwareType = getStringValue(artifact.customProperties, 'hardware');
      if (!hardwareType || !appliedHardwareTypes.includes(hardwareType)) {
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

    // Max Latency Filter (enhanced with metric and percentile support)
    const maxLatencyFilter = filterState[ModelCatalogNumberFilterKey.MAX_LATENCY];
    if (maxLatencyFilter !== undefined) {
      // TODO: For now using default TTFT Mean since we only store numeric value
      // In future iterations, we should store the full LatencyFilterConfig
      // and use applyMaxLatencyFilter(artifact, fullConfig)

      const defaultConfig: LatencyFilterConfig = {
        metric: 'TTFT',
        percentile: 'Mean',
        value: maxLatencyFilter,
      };

      if (!applyMaxLatencyFilter(artifact, defaultConfig)) {
        return false;
      }
    }

    // Workload Type Filter (Task-based, would need to be enhanced with real workload mapping)
    const workloadTypeFilter = filterState[ModelCatalogNumberFilterKey.WORKLOAD_TYPE];
    if (workloadTypeFilter !== undefined) {
      // For now, we'll assume all artifacts match workload type since we don't have
      // this property in the artifact data structure yet
      // This could be enhanced when workload type is available in the artifact data
    }

    return true;
  });

/**
 * Gets the count of active filters
 */
export const getActiveFilterCount = (filterState: ModelCatalogFilterStates): number => {
  let count = 0;

  // Count hardware type filters
  const hardwareTypeFilters = filterState[ModelCatalogStringFilterKey.PROVIDER];
  if (hardwareTypeFilters.length > 0) {
    count++;
  }

  // Count number filters
  if (filterState[ModelCatalogNumberFilterKey.MIN_RPS] !== undefined) {
    count++;
  }
  if (filterState[ModelCatalogNumberFilterKey.MAX_LATENCY] !== undefined) {
    count++;
  }
  if (filterState[ModelCatalogNumberFilterKey.WORKLOAD_TYPE] !== undefined) {
    count++;
  }

  return count;
};

/**
 * Clears all active filters
 */
export const clearAllFilters = (
  setFilterData: <K extends keyof ModelCatalogFilterStates>(
    key: K,
    value: ModelCatalogFilterStates[K],
  ) => void,
): void => {
  // Clear hardware filters
  setFilterData(ModelCatalogStringFilterKey.PROVIDER, []);

  // Clear number filters
  setFilterData(ModelCatalogNumberFilterKey.MIN_RPS, undefined);
  setFilterData(ModelCatalogNumberFilterKey.MAX_LATENCY, undefined);
  setFilterData(ModelCatalogNumberFilterKey.WORKLOAD_TYPE, undefined);
};
