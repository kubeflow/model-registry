import {
  CatalogPerformanceMetricsArtifact,
  PerformanceMetricsCustomProperties,
} from '~/app/modelCatalogTypes';
import { getDoubleValue, getIntValue, getStringValue } from '~/app/utils';

export const getHardwareConfiguration = (artifact: CatalogPerformanceMetricsArtifact): string => {
  const count = getIntValue(artifact.customProperties, 'hardware_count');
  const hardware = getStringValue(artifact.customProperties, 'hardware');
  return `${count} x ${hardware}`;
};

export const getTotalRps = (
  customProperties: PerformanceMetricsCustomProperties | undefined,
): number =>
  getIntValue(customProperties, 'hardware_count') *
  getDoubleValue(customProperties, 'requests_per_second');

export const formatLatency = (value: number): string => `${value.toFixed(2)} ms`;

export const formatTokenValue = (value: number): string => value.toFixed(0);
