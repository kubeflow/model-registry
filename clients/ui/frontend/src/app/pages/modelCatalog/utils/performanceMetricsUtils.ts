import { CatalogPerformanceMetricsArtifact } from '~/app/modelCatalogTypes';
import { getIntValue, getStringValue } from '~/app/utils';

export const getHardwareConfiguration = (artifact: CatalogPerformanceMetricsArtifact): string => {
  const count = getIntValue(artifact.customProperties, 'hardware_count');
  const hardware = getStringValue(artifact.customProperties, 'hardware_type');
  return `${count} x ${hardware}`;
};

export const formatLatency = (value: number): string => `${value.toFixed(2)} ms`;

export const formatTokenValue = (value: number): string => value.toFixed(0);
