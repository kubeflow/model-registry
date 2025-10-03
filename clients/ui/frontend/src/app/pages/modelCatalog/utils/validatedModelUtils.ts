import {
  CatalogPerformanceMetricsArtifact,
  CatalogAccuracyMetricsArtifact,
} from '~/app/modelCatalogTypes';

export type ValidatedModelMetrics = {
  accuracy: number;
  hardware: string;
  hardwareCount: string;
  rpsPerReplica: number;
  ttftMean: number;
};

export const extractValidatedModelMetrics = (
  performanceMetrics?: CatalogPerformanceMetricsArtifact,
  accuracyMetrics?: CatalogAccuracyMetricsArtifact,
): ValidatedModelMetrics => ({
  accuracy: accuracyMetrics?.customProperties.overall_average?.double_value || 53.9,
  hardware: performanceMetrics?.customProperties.hardware?.string_value || '8xH100-80',
  hardwareCount: performanceMetrics?.customProperties.hardware_count?.int_value || '8',
  rpsPerReplica: performanceMetrics?.customProperties.requests_per_second?.double_value || 1,
  ttftMean: performanceMetrics?.customProperties.ttft_mean?.double_value || 1428,
});
