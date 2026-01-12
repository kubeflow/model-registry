import { asEnumMember } from 'mod-arch-core';
import { CatalogPerformanceMetricsArtifact } from '~/app/modelCatalogTypes';
import { getIntValue, getStringValue } from '~/app/utils';
import { UseCaseOptionValue, PerformancePropertyKey } from '~/concepts/modelCatalog/const';
import { getUseCaseOption } from './workloadTypeUtils';

export type SliderRange = {
  minValue: number;
  maxValue: number;
  isSliderDisabled: boolean;
};

export const FALLBACK_RPS_RANGE: SliderRange = {
  minValue: 1,
  maxValue: 300,
  isSliderDisabled: false,
};

export const FALLBACK_LATENCY_RANGE: SliderRange = {
  minValue: 20,
  maxValue: 893,
  isSliderDisabled: false,
};

type CalculateSliderRangeOptions = {
  performanceArtifacts: CatalogPerformanceMetricsArtifact[];
  getArtifactFilterValue: (artifact: CatalogPerformanceMetricsArtifact) => number;
  fallbackRange: SliderRange;
  shouldRound?: boolean;
};

export const getHardwareConfiguration = (artifact: CatalogPerformanceMetricsArtifact): string => {
  const count = getIntValue(artifact.customProperties, 'hardware_count');
  const hardware = getStringValue(artifact.customProperties, PerformancePropertyKey.HARDWARE_TYPE);
  return `${count} x ${hardware}`;
};

export const formatLatency = (value: number): string => `${value.toFixed(2)} ms`;

export const formatTokenValue = (value: number): string => value.toFixed(0);

export const getWorkloadType = (artifact: CatalogPerformanceMetricsArtifact): string => {
  const useCaseValue = getStringValue(artifact.customProperties, PerformancePropertyKey.USE_CASE);
  if (!useCaseValue) {
    return '-';
  }
  const useCaseEnum = asEnumMember(useCaseValue, UseCaseOptionValue);
  if (!useCaseEnum) {
    return '-';
  }
  return getUseCaseOption(useCaseEnum)?.label || '-';
};

export const getSliderRange = ({
  performanceArtifacts,
  getArtifactFilterValue,
  fallbackRange,
  shouldRound = false,
}: CalculateSliderRangeOptions): SliderRange => {
  if (performanceArtifacts.length === 0) {
    return fallbackRange;
  }

  const values = performanceArtifacts.map(getArtifactFilterValue).filter((value) => value > 0);

  if (values.length === 0) {
    return fallbackRange;
  }

  const minValue = Math.min(...values);
  const maxValue = Math.max(...values);

  const calculatedMin = shouldRound ? Math.round(minValue) : minValue;
  const calculatedMax = shouldRound ? Math.round(maxValue) : maxValue;
  const hasIdenticalValues = calculatedMin === calculatedMax;

  return {
    minValue: calculatedMin,
    maxValue: hasIdenticalValues ? calculatedMin + 1 : calculatedMax,
    isSliderDisabled: hasIdenticalValues,
  };
};
