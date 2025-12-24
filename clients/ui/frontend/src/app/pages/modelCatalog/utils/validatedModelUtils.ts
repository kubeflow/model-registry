import {
  CatalogPerformanceMetricsArtifact,
  CatalogAccuracyMetricsArtifact,
} from '~/app/modelCatalogTypes';

export type ValidatedModelMetrics = {
  // accuracy: number; // NOTE: overall_average is currently omitted from the API and will be restored
  hardwareType: string;
  hardwareCount: string;
  rpsPerReplica: number;
  ttftMean: number;
  replicas: number | undefined;
  totalRequestsPerSecond: number | undefined;
};

export type PerformanceMetrics = {
  hardwareType: string;
  hardwareCount: string;
  rpsPerReplica: number;
  ttftMean: number;
  replicas: number | undefined;
  totalRequestsPerSecond: number | undefined;
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
});

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
      };

  return {
    // accuracy: calculateAverageAccuracy(accuracyMetrics), // NOTE: overall_average is currently omitted from the API and will be restored
    ...performance,
  };
};
