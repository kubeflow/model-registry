/* eslint-disable camelcase */
import { ModelRegistryMetadataType } from '~/app/types';
import { UseCaseOptionValue } from '~/concepts/modelCatalog/const';
import { mockCatalogPerformanceMetricsArtifact } from '~/__mocks__';
import {
  getHardwareConfiguration,
  formatLatency,
  formatTokenValue,
  getWorkloadType,
  getSliderRange,
  FALLBACK_LATENCY_RANGE,
  FALLBACK_RPS_RANGE,
} from '~/app/pages/modelCatalog/utils/performanceMetricsUtils';
import {
  CatalogPerformanceMetricsArtifact,
  CatalogArtifactType,
  MetricsType,
} from '~/app/modelCatalogTypes';
import { getDoubleValue } from '~/app/utils';

describe('performanceMetricsUtils', () => {
  describe('getHardwareConfiguration', () => {
    it('should return formatted hardware configuration string', () => {
      const artifact = mockCatalogPerformanceMetricsArtifact();
      expect(getHardwareConfiguration(artifact)).toBe('2 x H100-80');
    });

    it('should handle single hardware count', () => {
      const artifact = mockCatalogPerformanceMetricsArtifact({
        customProperties: {
          hardware_type: {
            metadataType: ModelRegistryMetadataType.STRING,
            string_value: 'A100-40',
          },
          hardware_count: {
            metadataType: ModelRegistryMetadataType.INT,
            int_value: '1',
          },
          requests_per_second: {
            metadataType: ModelRegistryMetadataType.DOUBLE,
            double_value: 5,
          },
        },
      });
      expect(getHardwareConfiguration(artifact)).toBe('1 x A100-40');
    });

    it('should handle large hardware counts', () => {
      const artifact = mockCatalogPerformanceMetricsArtifact({
        customProperties: {
          hardware_type: {
            metadataType: ModelRegistryMetadataType.STRING,
            string_value: 'GPU-V100',
          },
          hardware_count: {
            metadataType: ModelRegistryMetadataType.INT,
            int_value: '8',
          },
          requests_per_second: {
            metadataType: ModelRegistryMetadataType.DOUBLE,
            double_value: 10,
          },
        },
      });
      expect(getHardwareConfiguration(artifact)).toBe('8 x GPU-V100');
    });

    it('should return "0 x " when hardware_count is missing', () => {
      const artifact = mockCatalogPerformanceMetricsArtifact();
      delete artifact.customProperties.hardware_count;
      expect(getHardwareConfiguration(artifact)).toBe('0 x H100-80');
    });

    it('should return count with "-" when hardware_type is missing', () => {
      const artifact = mockCatalogPerformanceMetricsArtifact();
      delete artifact.customProperties.hardware_type;
      expect(getHardwareConfiguration(artifact)).toBe('2 x -');
    });
  });

  describe('formatLatency', () => {
    it('should format latency value with 2 decimal places and "ms" suffix', () => {
      expect(formatLatency(35.48818160947744)).toBe('35.49 ms');
    });

    it('should format whole numbers with 2 decimal places', () => {
      expect(formatLatency(100)).toBe('100.00 ms');
    });

    it('should round to 2 decimal places', () => {
      expect(formatLatency(123.456789)).toBe('123.46 ms');
    });

    it('should handle 0', () => {
      expect(formatLatency(0)).toBe('0.00 ms');
    });

    it('should handle very small numbers', () => {
      expect(formatLatency(0.001)).toBe('0.00 ms');
    });

    it('should handle very large numbers', () => {
      expect(formatLatency(99999.999)).toBe('100000.00 ms');
    });

    it('should handle negative numbers', () => {
      expect(formatLatency(-50.5)).toBe('-50.50 ms');
    });
  });

  describe('formatTokenValue', () => {
    it('should format token value with 0 decimal places', () => {
      expect(formatTokenValue(123.456)).toBe('123');
    });

    it('should round to nearest integer', () => {
      expect(formatTokenValue(99.9)).toBe('100');
      expect(formatTokenValue(99.4)).toBe('99');
    });

    it('should handle whole numbers', () => {
      expect(formatTokenValue(500)).toBe('500');
    });

    it('should handle 0', () => {
      expect(formatTokenValue(0)).toBe('0');
    });

    it('should handle very large numbers', () => {
      expect(formatTokenValue(999999.99)).toBe('1000000');
    });

    it('should handle negative numbers', () => {
      expect(formatTokenValue(-42.7)).toBe('-43');
    });
  });

  describe('getWorkloadType', () => {
    it('should return pretty-printed label for valid use case', () => {
      const artifact = mockCatalogPerformanceMetricsArtifact({
        customProperties: {
          use_case: {
            metadataType: ModelRegistryMetadataType.STRING,
            string_value: UseCaseOptionValue.CODE_FIXING,
          },
        },
      });
      expect(getWorkloadType(artifact)).toBe('Code Fixing');
    });

    it('should handle chatbot use case', () => {
      const artifact = mockCatalogPerformanceMetricsArtifact({
        customProperties: {
          use_case: {
            metadataType: ModelRegistryMetadataType.STRING,
            string_value: UseCaseOptionValue.CHATBOT,
          },
        },
      });
      expect(getWorkloadType(artifact)).toBe('Chatbot');
    });

    it('should handle RAG use case', () => {
      const artifact = mockCatalogPerformanceMetricsArtifact({
        customProperties: {
          use_case: {
            metadataType: ModelRegistryMetadataType.STRING,
            string_value: UseCaseOptionValue.RAG,
          },
        },
      });
      expect(getWorkloadType(artifact)).toBe('RAG');
    });

    it('should handle Long RAG use case', () => {
      const artifact = mockCatalogPerformanceMetricsArtifact({
        customProperties: {
          use_case: {
            metadataType: ModelRegistryMetadataType.STRING,
            string_value: UseCaseOptionValue.LONG_RAG,
          },
        },
      });
      expect(getWorkloadType(artifact)).toBe('Long RAG');
    });

    it('should return "-" when use_case is missing', () => {
      const artifact = mockCatalogPerformanceMetricsArtifact();
      delete artifact.customProperties.use_case;
      expect(getWorkloadType(artifact)).toBe('-');
    });

    it('should return "-" when use_case is empty string', () => {
      const artifact = mockCatalogPerformanceMetricsArtifact({
        customProperties: {
          use_case: {
            metadataType: ModelRegistryMetadataType.STRING,
            string_value: '',
          },
        },
      });
      expect(getWorkloadType(artifact)).toBe('-');
    });

    it('should return "-" when use_case is not a valid enum value', () => {
      const artifact = mockCatalogPerformanceMetricsArtifact({
        customProperties: {
          use_case: {
            metadataType: ModelRegistryMetadataType.STRING,
            string_value: 'invalid-use-case',
          },
        },
      });
      expect(getWorkloadType(artifact)).toBe('-');
    });

    it('should handle code_fixing with underscores', () => {
      const artifact = mockCatalogPerformanceMetricsArtifact();
      expect(getWorkloadType(artifact)).toBe('Code Fixing');
    });

    it('should handle long_rag with underscores', () => {
      const artifact = mockCatalogPerformanceMetricsArtifact({
        customProperties: {
          use_case: {
            metadataType: ModelRegistryMetadataType.STRING,
            string_value: 'long_rag',
          },
        },
      });
      expect(getWorkloadType(artifact)).toBe('Long RAG');
    });

    it('should return "-" when customProperties is undefined', () => {
      const artifact = mockCatalogPerformanceMetricsArtifact();
      // @ts-expect-error - Testing undefined customProperties
      artifact.customProperties = undefined;
      expect(getWorkloadType(artifact)).toBe('-');
    });
  });
});

describe('performanceMetricsUtils', () => {
  const createMockPerformanceArtifact = (
    rps: number,
    latency: number,
  ): CatalogPerformanceMetricsArtifact => ({
    artifactType: CatalogArtifactType.metricsArtifact,
    metricsType: MetricsType.performanceMetrics,
    createTimeSinceEpoch: '1739210683000',
    lastUpdateTimeSinceEpoch: '1739210683000',
    customProperties: {
      requests_per_second: {
        metadataType: ModelRegistryMetadataType.DOUBLE,
        double_value: rps,
      },
      ttft_mean: {
        metadataType: ModelRegistryMetadataType.DOUBLE,
        double_value: latency,
      },
    },
  });

  describe('getSliderRange', () => {
    describe('with empty performance artifacts', () => {
      it('should return fallback range', () => {
        const result = getSliderRange({
          performanceArtifacts: [],
          getArtifactFilterValue: (artifact) =>
            getDoubleValue(artifact.customProperties, 'requests_per_second'),
          fallbackRange: FALLBACK_RPS_RANGE,
        });

        expect(result).toEqual(FALLBACK_RPS_RANGE);
      });
    });

    describe('with no valid values', () => {
      it('should return fallback range when all values are invalid (zero or negative)', () => {
        const artifacts = [
          createMockPerformanceArtifact(0, 100),
          createMockPerformanceArtifact(-5, 200),
          createMockPerformanceArtifact(0, 300),
        ];

        const result = getSliderRange({
          performanceArtifacts: artifacts,
          getArtifactFilterValue: (artifact) =>
            getDoubleValue(artifact.customProperties, 'requests_per_second'),
          fallbackRange: FALLBACK_RPS_RANGE,
        });

        expect(result).toEqual(FALLBACK_RPS_RANGE);
      });
    });

    describe('with valid values', () => {
      it('should calculate min and max from RPS values', () => {
        const artifacts = [
          createMockPerformanceArtifact(10, 100),
          createMockPerformanceArtifact(50, 200),
          createMockPerformanceArtifact(25, 150),
        ];

        const result = getSliderRange({
          performanceArtifacts: artifacts,
          getArtifactFilterValue: (artifact) =>
            getDoubleValue(artifact.customProperties, 'requests_per_second'),
          fallbackRange: FALLBACK_RPS_RANGE,
        });

        expect(result).toEqual({
          minValue: 10,
          maxValue: 50,
          isSliderDisabled: false,
        });
      });

      it('should calculate min and max from latency values', () => {
        const artifacts = [
          createMockPerformanceArtifact(10, 150.5),
          createMockPerformanceArtifact(20, 200.8),
          createMockPerformanceArtifact(30, 100.2),
        ];

        const result = getSliderRange({
          performanceArtifacts: artifacts,
          getArtifactFilterValue: (artifact) =>
            getDoubleValue(artifact.customProperties, 'ttft_mean'),
          fallbackRange: FALLBACK_LATENCY_RANGE,
        });

        expect(result).toEqual({
          minValue: 100.2,
          maxValue: 200.8,
          isSliderDisabled: false,
        });
      });

      it('should filter out zero and negative values', () => {
        const artifacts = [
          createMockPerformanceArtifact(0, 100),
          createMockPerformanceArtifact(10, 200),
          createMockPerformanceArtifact(-5, 300),
          createMockPerformanceArtifact(30, 400),
        ];

        const result = getSliderRange({
          performanceArtifacts: artifacts,
          getArtifactFilterValue: (artifact) =>
            getDoubleValue(artifact.customProperties, 'requests_per_second'),
          fallbackRange: FALLBACK_RPS_RANGE,
        });

        expect(result).toEqual({
          minValue: 10,
          maxValue: 30,
          isSliderDisabled: false,
        });
      });
    });

    describe('with identical values', () => {
      it('should disable slider and add 1 to max when all values are identical', () => {
        const artifacts = [
          createMockPerformanceArtifact(25, 100),
          createMockPerformanceArtifact(25, 200),
          createMockPerformanceArtifact(25, 300),
        ];

        const result = getSliderRange({
          performanceArtifacts: artifacts,
          getArtifactFilterValue: (artifact) =>
            getDoubleValue(artifact.customProperties, 'requests_per_second'),
          fallbackRange: FALLBACK_RPS_RANGE,
        });

        expect(result).toEqual({
          minValue: 25,
          maxValue: 26, // 25 + 1
          isSliderDisabled: true,
        });
      });
    });

    describe('with shouldRound flag', () => {
      it('should round values when shouldRound is true', () => {
        const artifacts = [
          createMockPerformanceArtifact(10, 150.4),
          createMockPerformanceArtifact(20, 200.6),
          createMockPerformanceArtifact(30, 100.9),
        ];

        const result = getSliderRange({
          performanceArtifacts: artifacts,
          getArtifactFilterValue: (artifact) =>
            getDoubleValue(artifact.customProperties, 'ttft_mean'),
          fallbackRange: FALLBACK_LATENCY_RANGE,
          shouldRound: true,
        });

        expect(result).toEqual({
          minValue: 101, // Math.round(100.9)
          maxValue: 201, // Math.round(200.6)
          isSliderDisabled: false,
        });
      });

      it('should not round values when shouldRound is false', () => {
        const artifacts = [
          createMockPerformanceArtifact(10, 150.4),
          createMockPerformanceArtifact(20, 200.6),
          createMockPerformanceArtifact(30, 100.9),
        ];

        const result = getSliderRange({
          performanceArtifacts: artifacts,
          getArtifactFilterValue: (artifact) =>
            getDoubleValue(artifact.customProperties, 'ttft_mean'),
          fallbackRange: FALLBACK_LATENCY_RANGE,
          shouldRound: false,
        });

        expect(result).toEqual({
          minValue: 100.9,
          maxValue: 200.6,
          isSliderDisabled: false,
        });
      });

      it('should not round by default when shouldRound is omitted', () => {
        const artifacts = [
          createMockPerformanceArtifact(10.5, 100),
          createMockPerformanceArtifact(50.7, 200),
        ];

        const result = getSliderRange({
          performanceArtifacts: artifacts,
          getArtifactFilterValue: (artifact) =>
            getDoubleValue(artifact.customProperties, 'requests_per_second'),
          fallbackRange: FALLBACK_RPS_RANGE,
        });

        expect(result).toEqual({
          minValue: 10.5,
          maxValue: 50.7,
          isSliderDisabled: false,
        });
      });

      it('should handle identical values after rounding', () => {
        const artifacts = [
          createMockPerformanceArtifact(10, 100.3),
          createMockPerformanceArtifact(20, 100.4),
        ];

        const result = getSliderRange({
          performanceArtifacts: artifacts,
          getArtifactFilterValue: (artifact) =>
            getDoubleValue(artifact.customProperties, 'ttft_mean'),
          fallbackRange: FALLBACK_LATENCY_RANGE,
          shouldRound: true,
        });

        expect(result).toEqual({
          minValue: 100, // Math.round(100.3) = 100, Math.round(100.4) = 100
          maxValue: 101, // Since identical, adds 1
          isSliderDisabled: true,
        });
      });
    });
  });
});
