/* eslint-disable camelcase */
import {
  CatalogArtifacts,
  CatalogArtifactList,
  CatalogModelArtifact,
  CatalogArtifactType,
  MetricsType,
  CatalogPerformanceMetricsArtifact,
  CatalogAccuracyMetricsArtifact,
} from '~/app/modelCatalogTypes';
import { ModelRegistryMetadataType } from '~/app/types';
import { UseCaseOptionValue } from '~/concepts/modelCatalog/const';

export const mockCatalogModelArtifact = (
  partial?: Partial<CatalogModelArtifact>,
): CatalogArtifacts => ({
  artifactType: CatalogArtifactType.modelArtifact,
  createTimeSinceEpoch: '1739210683000',
  lastUpdateTimeSinceEpoch: '1739210683000',
  uri: '',
  customProperties: {},
  ...partial,
});

export const mockCatalogAccuracyMetricsArtifact = (
  partial?: Partial<CatalogAccuracyMetricsArtifact>,
): CatalogAccuracyMetricsArtifact => ({
  artifactType: CatalogArtifactType.metricsArtifact,
  metricsType: MetricsType.accuracyMetrics,
  createTimeSinceEpoch: '1739210683000',
  lastUpdateTimeSinceEpoch: '1739210683000',
  customProperties: {
    overall_average: {
      metadataType: ModelRegistryMetadataType.DOUBLE,
      double_value: 0.582439,
    },
    arc_v1: {
      metadataType: ModelRegistryMetadataType.DOUBLE,
      double_value: 0.659556,
    },
  },
  ...partial,
});

export const mockCatalogPerformanceMetricsArtifact = (
  partial?: Partial<CatalogPerformanceMetricsArtifact>,
): CatalogPerformanceMetricsArtifact => ({
  artifactType: CatalogArtifactType.metricsArtifact,
  metricsType: MetricsType.performanceMetrics,
  createTimeSinceEpoch: '1739210683000',
  lastUpdateTimeSinceEpoch: '1739210683000',
  customProperties: {
    config_id: {
      metadataType: ModelRegistryMetadataType.STRING,
      string_value: '0055d94f6a542f6932cac5dfa5ffdd38',
    },
    hardware_count: {
      metadataType: ModelRegistryMetadataType.INT,
      int_value: '2',
    },
    hardware_type: {
      metadataType: ModelRegistryMetadataType.STRING,
      string_value: 'H100-80',
    },
    requests_per_second: {
      metadataType: ModelRegistryMetadataType.DOUBLE,
      double_value: 7,
    },
    ttft_mean: {
      metadataType: ModelRegistryMetadataType.DOUBLE,
      double_value: 35.48818160947744,
    },
    ttft_p90: {
      metadataType: ModelRegistryMetadataType.DOUBLE,
      double_value: 51.55777931213379,
    },
    ttft_p95: {
      metadataType: ModelRegistryMetadataType.DOUBLE,
      double_value: 61.26761436462402,
    },
    ttft_p99: {
      metadataType: ModelRegistryMetadataType.DOUBLE,
      double_value: 72.95823097229004,
    },
    e2e_mean: {
      metadataType: ModelRegistryMetadataType.DOUBLE,
      double_value: 1994.480013381083,
    },
    use_case: {
      metadataType: ModelRegistryMetadataType.STRING,
      string_value: UseCaseOptionValue.CODE_FIXING,
    },
  },
  ...partial,
});

export const mockCatalogModelArtifactList = (
  partial?: Partial<CatalogModelArtifact>,
): CatalogArtifactList => ({
  items: [mockCatalogModelArtifact({})],
  pageSize: 10,
  size: 15,
  nextPageToken: '',
  ...partial,
});

export const mockCatalogPerformanceMetricsArtifactList = (
  partial?: Partial<CatalogPerformanceMetricsArtifact>,
): CatalogArtifactList => ({
  items: [mockCatalogPerformanceMetricsArtifact({}), mockCatalogModelArtifact({})],
  pageSize: 10,
  size: 15,
  nextPageToken: '',
  ...partial,
});

export const mockCatalogAccuracyMetricsArtifactList = (
  partial?: Partial<CatalogAccuracyMetricsArtifact>,
): CatalogArtifactList => ({
  items: [mockCatalogAccuracyMetricsArtifact({}), mockCatalogModelArtifact({})],
  pageSize: 10,
  size: 15,
  nextPageToken: '',
  ...partial,
});

// Performance artifact with computed properties (when targetRPS is provided)
export const mockCatalogPerformanceMetricsArtifactWithRPS = (
  targetRPS: number,
  partial?: Partial<CatalogPerformanceMetricsArtifact>,
): CatalogPerformanceMetricsArtifact => {
  const baseArtifact = mockCatalogPerformanceMetricsArtifact(partial);
  const rps = baseArtifact.customProperties?.requests_per_second?.double_value || 7;
  const replicas = Math.ceil(targetRPS / rps);

  return {
    ...baseArtifact,
    customProperties: {
      ...baseArtifact.customProperties,
      replicas: {
        metadataType: ModelRegistryMetadataType.INT,
        int_value: String(replicas),
      },
      total_requests_per_second: {
        metadataType: ModelRegistryMetadataType.DOUBLE,
        double_value: replicas * rps,
      },
    },
  };
};

// Mock for Pareto-filtered (recommendations=true) performance artifacts
export const mockParetoFilteredPerformanceArtifactList = (
  targetRPS?: number,
): CatalogArtifactList => ({
  items: [
    mockCatalogPerformanceMetricsArtifactWithRPS(targetRPS || 100, {
      customProperties: {
        config_id: {
          metadataType: ModelRegistryMetadataType.STRING,
          string_value: 'pareto-optimal-1',
        },
        hardware_count: { metadataType: ModelRegistryMetadataType.INT, int_value: '1' },
        hardware_type: { metadataType: ModelRegistryMetadataType.STRING, string_value: 'H100-80' },
        requests_per_second: { metadataType: ModelRegistryMetadataType.DOUBLE, double_value: 50 },
        ttft_p90: { metadataType: ModelRegistryMetadataType.DOUBLE, double_value: 35 },
        use_case: {
          metadataType: ModelRegistryMetadataType.STRING,
          string_value: UseCaseOptionValue.CHATBOT,
        },
      },
    }),
    mockCatalogPerformanceMetricsArtifactWithRPS(targetRPS || 100, {
      customProperties: {
        config_id: {
          metadataType: ModelRegistryMetadataType.STRING,
          string_value: 'pareto-optimal-2',
        },
        hardware_count: { metadataType: ModelRegistryMetadataType.INT, int_value: '2' },
        hardware_type: { metadataType: ModelRegistryMetadataType.STRING, string_value: 'A100-80' },
        requests_per_second: { metadataType: ModelRegistryMetadataType.DOUBLE, double_value: 30 },
        ttft_p90: { metadataType: ModelRegistryMetadataType.DOUBLE, double_value: 28 },
        use_case: {
          metadataType: ModelRegistryMetadataType.STRING,
          string_value: UseCaseOptionValue.RAG,
        },
      },
    }),
  ],
  pageSize: 10,
  size: 2,
  nextPageToken: '',
});
