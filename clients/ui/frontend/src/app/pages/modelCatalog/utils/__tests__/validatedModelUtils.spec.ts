/* eslint-disable camelcase */
import {
  CatalogPerformanceMetricsArtifact,
  CatalogAccuracyMetricsArtifact,
  CatalogArtifactType,
  MetricsType,
} from '~/app/modelCatalogTypes';
import { ModelRegistryMetadataType } from '~/app/types';
import {
  extractPerformanceMetrics,
  calculateAverageAccuracy,
  extractValidatedModelMetrics,
} from '~/app/pages/modelCatalog/utils/validatedModelUtils';

describe('validatedModelUtils', () => {
  const createMockPerformanceArtifact = (
    hardware: string,
    hardwareCount: number,
    rpsPerReplica: number,
    ttftMean: number,
  ): CatalogPerformanceMetricsArtifact => ({
    artifactType: CatalogArtifactType.metricsArtifact,
    metricsType: MetricsType.performanceMetrics,
    createTimeSinceEpoch: '1739210683000',
    lastUpdateTimeSinceEpoch: '1739210683000',
    customProperties: {
      hardware: {
        metadataType: ModelRegistryMetadataType.STRING,
        string_value: hardware,
      },
      hardware_count: {
        metadataType: ModelRegistryMetadataType.INT,
        int_value: hardwareCount.toString(),
      },
      requests_per_second: {
        metadataType: ModelRegistryMetadataType.DOUBLE,
        double_value: rpsPerReplica,
      },
      ttft_mean: {
        metadataType: ModelRegistryMetadataType.DOUBLE,
        double_value: ttftMean,
      },
    },
  });

  const createMockAccuracyArtifact = (overallAverage: number): CatalogAccuracyMetricsArtifact => ({
    artifactType: CatalogArtifactType.metricsArtifact,
    metricsType: MetricsType.accuracyMetrics,
    createTimeSinceEpoch: '1739210683000',
    lastUpdateTimeSinceEpoch: '1739210683000',
    customProperties: {
      overall_average: {
        metadataType: ModelRegistryMetadataType.DOUBLE,
        double_value: overallAverage,
      },
    },
  });

  describe('extractPerformanceMetrics', () => {
    it('should extract performance metrics from a single artifact', () => {
      const artifact = createMockPerformanceArtifact('H100-80', 2, 3.5, 1200);

      const result = extractPerformanceMetrics(artifact);

      expect(result).toEqual({
        hardware: 'H100-80',
        hardwareCount: '2',
        rpsPerReplica: 3.5,
        ttftMean: 1200,
      });
    });

    it('should handle missing properties with default values', () => {
      const artifact: CatalogPerformanceMetricsArtifact = {
        artifactType: CatalogArtifactType.metricsArtifact,
        metricsType: MetricsType.performanceMetrics,
        createTimeSinceEpoch: '1739210683000',
        lastUpdateTimeSinceEpoch: '1739210683000',
        customProperties: {},
      };

      const result = extractPerformanceMetrics(artifact);

      expect(result).toEqual({
        hardware: 'H100-80',
        hardwareCount: '1',
        rpsPerReplica: 1,
        ttftMean: 1428,
      });
    });
  });

  describe('calculateAverageAccuracy', () => {
    it('should calculate average accuracy from multiple artifacts', () => {
      const artifacts = [
        createMockAccuracyArtifact(50.0),
        createMockAccuracyArtifact(60.0),
        createMockAccuracyArtifact(70.0),
      ];

      const result = calculateAverageAccuracy(artifacts);

      expect(result).toBe(60.0);
    });

    it('should handle single artifact', () => {
      const artifacts = [createMockAccuracyArtifact(75.5)];

      const result = calculateAverageAccuracy(artifacts);

      expect(result).toBe(75.5);
    });

    it('should handle empty array with default fallback', () => {
      const result = calculateAverageAccuracy([]);

      expect(result).toBe(53.9);
    });

    it('should round to 1 decimal place', () => {
      const artifacts = [
        createMockAccuracyArtifact(50.0),
        createMockAccuracyArtifact(60.0),
        createMockAccuracyArtifact(70.0),
        createMockAccuracyArtifact(80.0),
      ];

      const result = calculateAverageAccuracy(artifacts);

      expect(result).toBe(65.0);
    });

    it('should handle artifacts with missing accuracy values', () => {
      const artifacts = [
        createMockAccuracyArtifact(50.0),
        {
          ...createMockAccuracyArtifact(0),
          customProperties: {},
        },
        createMockAccuracyArtifact(70.0),
      ];

      const result = calculateAverageAccuracy(artifacts);

      expect(result).toBe(40.0); // (50 + 0 + 70) / 3
    });
  });

  describe('extractValidatedModelMetrics', () => {
    it('should extract metrics from arrays of artifacts with specific performance index', () => {
      const performanceArtifacts = [
        createMockPerformanceArtifact('A100-80', 1, 2.0, 1000),
        createMockPerformanceArtifact('H100-80', 2, 3.5, 1200),
        createMockPerformanceArtifact('A100-40', 4, 5.0, 1500),
      ];

      const accuracyArtifacts = [
        createMockAccuracyArtifact(50.0),
        createMockAccuracyArtifact(60.0),
        createMockAccuracyArtifact(70.0),
      ];

      const result = extractValidatedModelMetrics(performanceArtifacts, accuracyArtifacts, 1);

      expect(result).toEqual({
        accuracy: 60.0,
        hardware: 'H100-80',
        hardwareCount: '2',
        rpsPerReplica: 3.5,
        ttftMean: 1200,
      });
    });

    it('should use first performance artifact by default', () => {
      const performanceArtifacts = [
        createMockPerformanceArtifact('A100-80', 1, 2.0, 1000),
        createMockPerformanceArtifact('H100-80', 2, 3.5, 1200),
      ];

      const accuracyArtifacts = [createMockAccuracyArtifact(75.0)];

      const result = extractValidatedModelMetrics(performanceArtifacts, accuracyArtifacts);

      expect(result).toEqual({
        accuracy: 75.0,
        hardware: 'A100-80',
        hardwareCount: '1',
        rpsPerReplica: 2.0,
        ttftMean: 1000,
      });
    });

    it('should handle empty performance artifacts array with defaults', () => {
      const performanceArtifacts: CatalogPerformanceMetricsArtifact[] = [];
      const accuracyArtifacts = [createMockAccuracyArtifact(80.0)];

      const result = extractValidatedModelMetrics(performanceArtifacts, accuracyArtifacts);

      expect(result).toEqual({
        accuracy: 80.0,
        hardware: 'H100-80',
        hardwareCount: '1',
        rpsPerReplica: 1,
        ttftMean: 1428,
      });
    });

    it('should handle invalid performance index with defaults', () => {
      const performanceArtifacts = [createMockPerformanceArtifact('A100-80', 1, 2.0, 1000)];
      const accuracyArtifacts = [createMockAccuracyArtifact(90.0)];

      const result = extractValidatedModelMetrics(performanceArtifacts, accuracyArtifacts, 5);

      expect(result).toEqual({
        accuracy: 90.0,
        hardware: 'H100-80',
        hardwareCount: '1',
        rpsPerReplica: 1,
        ttftMean: 1428,
      });
    });
  });
});
