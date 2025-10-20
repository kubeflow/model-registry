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
