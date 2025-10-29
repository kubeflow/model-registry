import { asEnumMember } from 'mod-arch-core';
import {
  CatalogPerformanceMetricsArtifact,
  PerformanceMetricsCustomProperties,
} from '~/app/modelCatalogTypes';
import { getDoubleValue, getIntValue, getStringValue } from '~/app/utils';
import { UseCaseOptionValue } from '~/concepts/modelCatalog/const';
import { getUseCaseOption } from './workloadTypeUtils';

export const getHardwareConfiguration = (artifact: CatalogPerformanceMetricsArtifact): string => {
  const count = getIntValue(artifact.customProperties, 'hardware_count');
  const hardware = getStringValue(artifact.customProperties, 'hardware_type');
  return `${count} x ${hardware}`;
};

export const getTotalRps = (
  customProperties: PerformanceMetricsCustomProperties | undefined,
): number =>
  getIntValue(customProperties, 'hardware_count') *
  getDoubleValue(customProperties, 'requests_per_second');

export const formatLatency = (value: number): string => `${value.toFixed(2)} ms`;

export const formatTokenValue = (value: number): string => value.toFixed(0);

export const getWorkloadType = (artifact: CatalogPerformanceMetricsArtifact): string => {
  const useCaseValue = getStringValue(artifact.customProperties, 'use_case');
  if (!useCaseValue) {
    return '-';
  }
  const useCaseEnum = asEnumMember(useCaseValue, UseCaseOptionValue);
  if (!useCaseEnum) {
    return '-';
  }
  return getUseCaseOption(useCaseEnum)?.label || '-';
};
