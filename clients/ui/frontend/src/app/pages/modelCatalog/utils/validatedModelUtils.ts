import {
  CatalogPerformanceMetricsArtifact,
  CatalogAccuracyMetricsArtifact,
} from '~/app/modelCatalogTypes';
import {
  LatencyMetric,
  LatencyMetricFieldName,
  ALL_LATENCY_FIELD_NAMES,
} from '~/concepts/modelCatalog/const';

/**
 * Type for latency metrics - uses LatencyMetricFieldName from const.ts
 * to dynamically define all possible latency field keys
 */
export type LatencyMetricsMap = Partial<Record<LatencyMetricFieldName, number>>;

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
 * Extracts all latency metrics from artifact custom properties
 * using ALL_LATENCY_FIELD_NAMES from const.ts
 */
const extractLatencyMetrics = (
  customProperties: CatalogPerformanceMetricsArtifact['customProperties'],
): LatencyMetricsMap => {
  const result: LatencyMetricsMap = {};
  ALL_LATENCY_FIELD_NAMES.forEach((fieldName) => {
    const value = customProperties?.[fieldName]?.double_value;
    if (value !== undefined) {
      result[fieldName] = value;
    }
  });
  return result;
};

export const extractPerformanceMetrics = (
  performanceMetrics: CatalogPerformanceMetricsArtifact,
): PerformanceMetrics => ({
  hardwareType: performanceMetrics.customProperties?.hardware_type?.string_value || 'H100-80',
  hardwareCount: performanceMetrics.customProperties?.hardware_count?.int_value || '1',
  rpsPerReplica: performanceMetrics.customProperties?.requests_per_second?.double_value || 1,
  ttftMean: performanceMetrics.customProperties?.ttft_mean?.double_value || 1428,
  replicas: performanceMetrics.customProperties?.replicas?.int_value
    ? Number(performanceMetrics.customProperties.replicas.int_value)
    : undefined,
  totalRequestsPerSecond:
    performanceMetrics.customProperties?.total_requests_per_second?.double_value,
  latencyMetrics: extractLatencyMetrics(performanceMetrics.customProperties),
});

/**
 * Gets the latency value for a specific field name from the latency metrics
 */
export const getLatencyValue = (
  latencyMetrics: ValidatedModelMetrics['latencyMetrics'],
  fieldName: LatencyMetricFieldName | undefined,
): number | undefined => {
  if (!fieldName) {
    // Default to ttft_mean if no field specified
    return latencyMetrics.ttft_mean;
  }
  return latencyMetrics[fieldName];
};

/**
 * Gets the display label for a latency metric (e.g., "TTFT", "ITL", "E2E")
 */
export const getLatencyLabel = (fieldName: LatencyMetricFieldName | undefined): string => {
  if (!fieldName) {
    return LatencyMetric.TTFT;
  }
  const [prefix] = fieldName.split('_');
  const metric = Object.values(LatencyMetric).find((m) => m.toLowerCase() === prefix);
  return metric || LatencyMetric.TTFT;
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
