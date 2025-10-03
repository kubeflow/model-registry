import {
  CatalogPerformanceMetricsArtifact,
  PerformanceMetricsCustomProperties,
} from '~/app/modelCatalogTypes';
import {
  ModelRegistryCustomPropertyString,
  ModelRegistryCustomPropertyInt,
  ModelRegistryCustomPropertyDouble,
  ModelRegistryMetadataType,
} from '~/app/types';

export const getStringValue = (
  customProperties: PerformanceMetricsCustomProperties | undefined,
  key: keyof PerformanceMetricsCustomProperties,
): string => {
  const prop = customProperties?.[key];
  if (prop && prop.metadataType === ModelRegistryMetadataType.STRING) {
    // eslint-disable-next-line @typescript-eslint/consistent-type-assertions
    return (prop as ModelRegistryCustomPropertyString).string_value;
  }
  return '-';
};

export const getIntValue = (
  customProperties: PerformanceMetricsCustomProperties | undefined,
  key: keyof PerformanceMetricsCustomProperties,
): number => {
  const prop = customProperties?.[key];
  if (prop && prop.metadataType === ModelRegistryMetadataType.INT) {
    // eslint-disable-next-line @typescript-eslint/consistent-type-assertions
    const value = (prop as ModelRegistryCustomPropertyInt).int_value;
    return value ? parseInt(value, 10) : 0;
  }
  return 0;
};

export const getDoubleValue = (
  customProperties: PerformanceMetricsCustomProperties | undefined,
  key: keyof PerformanceMetricsCustomProperties,
): number => {
  const prop = customProperties?.[key];
  if (prop && prop.metadataType === ModelRegistryMetadataType.DOUBLE) {
    // eslint-disable-next-line @typescript-eslint/consistent-type-assertions
    return (prop as ModelRegistryCustomPropertyDouble).double_value;
  }
  return 0;
};

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
