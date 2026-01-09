import {
  CatalogPerformanceMetricsArtifact,
  CatalogAccuracyMetricsArtifact,
} from '~/app/modelCatalogTypes';
import {
  LatencyFilterKey,
  LatencyMetric,
  LatencyPercentile,
  getLatencyPropertyKey,
  getLatencyFilterKey,
} from '~/concepts/modelCatalog/const';

/**
 * Type for latency metrics - uses LatencyFilterKey from const.ts
 * to dynamically define all possible latency field keys
 */
export type LatencyMetricsMap = Partial<Record<LatencyFilterKey, number>>;

export type ValidatedModelMetrics = {
  // accuracy: number; // NOTE: overall_average is currently omitted from the API and will be restored
  hardwareType: string;
  hardwareCount: string;
  rpsPerReplica: number;
  ttftMean: number;
  replicas: number | undefined;
  totalRequestsPerSecond: number | undefined;
  latencyMetrics: LatencyMetricsMap;
};

export type PerformanceMetrics = {
  hardwareType: string;
  hardwareCount: string;
  rpsPerReplica: number;
  ttftMean: number;
  replicas: number | undefined;
  totalRequestsPerSecond: number | undefined;
  latencyMetrics: LatencyMetricsMap;
};

/**
 * Extracts all latency metrics from artifact custom properties.
 * Loops over LatencyMetric and LatencyPercentile enums to build the map.
 * Uses the full filter key format (e.g., 'artifacts.ttft_p90.double_value') as keys.
 */
const extractLatencyMetrics = (
  customProperties: CatalogPerformanceMetricsArtifact['customProperties'],
): LatencyMetricsMap => {
  const result: LatencyMetricsMap = {};
  Object.values(LatencyMetric).forEach((metric) => {
    Object.values(LatencyPercentile).forEach((percentile) => {
      const propertyKey = getLatencyPropertyKey(metric, percentile);
      const filterKey = getLatencyFilterKey(metric, percentile);
      const value = customProperties?.[propertyKey]?.double_value;
      if (value !== undefined) {
        result[filterKey] = value;
      }
    });
  });
  return result;
};

export const extractPerformanceMetrics = (
  performanceMetrics: CatalogPerformanceMetricsArtifact,
): PerformanceMetrics => {
  const ttftMeanKey = getLatencyPropertyKey(LatencyMetric.TTFT, LatencyPercentile.Mean);
  return {
    hardwareType: performanceMetrics.customProperties?.hardware_type?.string_value || 'H100-80',
    hardwareCount: performanceMetrics.customProperties?.hardware_count?.int_value || '1',
    rpsPerReplica: performanceMetrics.customProperties?.requests_per_second?.double_value || 1,
    ttftMean: performanceMetrics.customProperties?.[ttftMeanKey]?.double_value || 1428,
    replicas: performanceMetrics.customProperties?.replicas?.int_value
      ? Number(performanceMetrics.customProperties.replicas.int_value)
      : undefined,
    totalRequestsPerSecond:
      performanceMetrics.customProperties?.total_requests_per_second?.double_value,
    latencyMetrics: extractLatencyMetrics(performanceMetrics.customProperties),
  };
};

/**
 * Gets the latency value for a specific filter key from the latency metrics.
 * The filterKey should be in the full format (e.g., 'artifacts.ttft_p90.double_value').
 */
export const getLatencyValue = (
  latencyMetrics: ValidatedModelMetrics['latencyMetrics'],
  filterKey: LatencyFilterKey | undefined,
): number | undefined => {
  if (!filterKey) {
    // Default to ttft_mean if no field specified
    const defaultKey = getLatencyFilterKey(LatencyMetric.TTFT, LatencyPercentile.Mean);
    return latencyMetrics[defaultKey];
  }
  return latencyMetrics[filterKey];
};

// NOTE: overall_average is currently omitted from the API and will be restored
// export const calculateAverageAccuracy = (
//   accuracyMetrics: CatalogAccuracyMetricsArtifact[],
// ): number => {
//   if (accuracyMetrics.length === 0) {
//     return 53.9; // Default fallback
//   }
//
//   const totalAccuracy = accuracyMetrics.reduce((sum, artifact) => {
//     const accuracy = artifact.customProperties.overall_average?.double_value || 0;
//     return sum + accuracy;
//   }, 0);
//
//   return Math.round((totalAccuracy / accuracyMetrics.length) * 10) / 10; // Round to 1 decimal place
// };

export const extractValidatedModelMetrics = (
  performanceMetrics: CatalogPerformanceMetricsArtifact[],
  _accuracyMetrics: CatalogAccuracyMetricsArtifact[],
  currentPerformanceIndex = 0,
): ValidatedModelMetrics => {
  const currentPerformance = performanceMetrics[currentPerformanceIndex];
  // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
  const performance = currentPerformance
    ? extractPerformanceMetrics(currentPerformance)
    : {
        hardwareType: 'H100-80',
        hardwareCount: '1',
        rpsPerReplica: 1,
        ttftMean: 1428,
        replicas: undefined,
        totalRequestsPerSecond: undefined,
        latencyMetrics: {},
      };

  return {
    // accuracy: calculateAverageAccuracy(accuracyMetrics), // NOTE: overall_average is currently omitted from the API and will be restored
    ...performance,
  };
};
