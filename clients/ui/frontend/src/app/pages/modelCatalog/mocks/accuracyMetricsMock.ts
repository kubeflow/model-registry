/* eslint-disable camelcase */
import {
  CatalogAccuracyMetricsArtifact,
  CatalogArtifactType,
  MetricsType,
} from '~/app/modelCatalogTypes';
import { ModelRegistryMetadataType } from '~/app/types';

const MOCK_TIMESTAMP = '1739210683000';

const createAccuracyMetricsArtifact = (
  id: string,
  // overallAverage: number, // NOTE: overall_average is currently omitted from the API and will be restored
  arcV1: number,
): CatalogAccuracyMetricsArtifact => ({
  artifactType: CatalogArtifactType.metricsArtifact,
  metricsType: MetricsType.accuracyMetrics,
  createTimeSinceEpoch: MOCK_TIMESTAMP,
  lastUpdateTimeSinceEpoch: MOCK_TIMESTAMP,
  customProperties: {
    // overall_average: { // NOTE: overall_average is currently omitted from the API and will be restored
    //   metadataType: ModelRegistryMetadataType.DOUBLE,
    //   double_value: overallAverage,
    // },
    arc_v1: {
      metadataType: ModelRegistryMetadataType.DOUBLE,
      double_value: arcV1,
    },
  },
});

export const mockAccuracyMetricsArtifacts: CatalogAccuracyMetricsArtifact[] = [
  createAccuracyMetricsArtifact('1', /* 53.9, */ 45.2),
  createAccuracyMetricsArtifact('2', /* 67.3, */ 58.1),
  createAccuracyMetricsArtifact('3', /* 42.1, */ 38.7),
  createAccuracyMetricsArtifact('4', /* 78.5, */ 72.3),
  createAccuracyMetricsArtifact('5', /* 61.2, */ 55.8),
  createAccuracyMetricsArtifact('6', /* 49.7, */ 43.9),
];
