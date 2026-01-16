import { CatalogPerformanceMetricsArtifact } from '~/app/modelCatalogTypes';
import { getStringValue } from '~/app/utils';
import {
  ModelCatalogStringFilterKey,
  ModelCatalogNumberFilterKey,
  LatencyMetricFieldName,
  ALL_LATENCY_FILTER_KEYS,
  PerformancePropertyKey,
} from '~/concepts/modelCatalog/const';

/**
 * Extracts unique hardware types from performance artifacts
 */
export const getUniqueHardwareTypes = (
  artifacts: CatalogPerformanceMetricsArtifact[],
): string[] => {
  // Use the short property key for accessing customProperties
  const hardwareTypes = artifacts
    .map((artifact) =>
      getStringValue(artifact.customProperties, PerformancePropertyKey.HARDWARE_TYPE),
    )
    .filter((hardware): hardware is string => !!hardware);

  return [...new Set(hardwareTypes)].toSorted();
};

/**
 * Gets all filter keys (string filters + number filters + latency filters)
 */
export const getAllFilterKeys = (): {
  stringFilterKeys: ModelCatalogStringFilterKey[];
  numberFilterKeys: ModelCatalogNumberFilterKey[];
  latencyFilterKeys: LatencyMetricFieldName[];
} => ({
  stringFilterKeys: Object.values(ModelCatalogStringFilterKey),
  numberFilterKeys: Object.values(ModelCatalogNumberFilterKey),
  latencyFilterKeys: ALL_LATENCY_FILTER_KEYS,
});
