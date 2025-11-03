/* eslint-disable camelcase */
import { ModelRegistryMetadataType } from '~/app/types';
import { UseCaseOptionValue } from '~/concepts/modelCatalog/const';
import { mockCatalogPerformanceMetricsArtifact } from '~/__mocks__';
import {
  getHardwareConfiguration,
  formatLatency,
  formatTokenValue,
  getWorkloadType,
} from '~/app/pages/modelCatalog/utils/performanceMetricsUtils';

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
