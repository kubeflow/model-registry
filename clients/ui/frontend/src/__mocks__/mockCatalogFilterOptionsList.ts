/* eslint-disable camelcase */
import { CatalogFilterOptionsList, NamedQuery, FilterOperator } from '~/app/modelCatalogTypes';
import {
  ModelCatalogStringFilterKey,
  ModelCatalogNumberFilterKey,
  ModelCatalogLicense,
  ModelCatalogProvider,
  ModelCatalogTask,
  AllLanguageCode,
  UseCaseOptionValue,
} from '~/concepts/modelCatalog/const';

export const mockNamedQueries: Record<string, NamedQuery> = {
  high_performance_gpu: {
    'hardware_type.string_value': { operator: FilterOperator.IN, value: ['H100-80', 'A100-80'] },
    'requests_per_second.double_value': {
      operator: FilterOperator.GREATER_THAN_OR_EQUAL,
      value: 50,
    },
  },
  low_latency: {
    'ttft_p90.double_value': { operator: FilterOperator.LESS_THAN, value: 100 },
    'e2e_p90.double_value': { operator: FilterOperator.LESS_THAN, value: 500 },
  },
  chatbot_optimized: {
    'use_case.string_value': { operator: FilterOperator.EQUALS, value: UseCaseOptionValue.CHATBOT },
  },
  rag_optimized: {
    'use_case.string_value': {
      operator: FilterOperator.IN,
      value: [UseCaseOptionValue.RAG, UseCaseOptionValue.LONG_RAG],
    },
  },
  cost_effective: {
    'hardware_count.int_value': { operator: FilterOperator.LESS_THAN_OR_EQUAL, value: 2 },
  },
};

export const mockCatalogFilterOptionsList = (
  partial?: Partial<CatalogFilterOptionsList>,
): CatalogFilterOptionsList => ({
  filters: {
    [ModelCatalogStringFilterKey.PROVIDER]: {
      type: 'string',
      values: [ModelCatalogProvider.RED_HAT, ModelCatalogProvider.IBM, ModelCatalogProvider.GOOGLE],
    },
    [ModelCatalogStringFilterKey.LICENSE]: {
      type: 'string',
      values: [ModelCatalogLicense.APACHE_2_0, ModelCatalogLicense.MIT],
    },
    [ModelCatalogStringFilterKey.TASK]: {
      type: 'string',
      values: [
        ModelCatalogTask.TEXT_GENERATION,
        ModelCatalogTask.TEXT_TO_TEXT,
        ModelCatalogTask.IMAGE_TO_TEXT,
        ModelCatalogTask.IMAGE_TEXT_TO_TEXT,
        ModelCatalogTask.VIDEO_TO_TEXT,
        ModelCatalogTask.AUDIO_TO_TEXT,
      ],
    },
    [ModelCatalogStringFilterKey.LANGUAGE]: {
      type: 'string',
      values: [
        AllLanguageCode.AR,
        AllLanguageCode.CS,
        AllLanguageCode.DE,
        AllLanguageCode.EN,
        AllLanguageCode.ES,
        AllLanguageCode.FR,
        AllLanguageCode.IT,
        AllLanguageCode.JA,
        AllLanguageCode.KO,
        AllLanguageCode.NL,
        AllLanguageCode.PT,
        AllLanguageCode.ZH,
      ],
    },
    [ModelCatalogStringFilterKey.HARDWARE_TYPE]: {
      type: 'string',
      values: ['GPU', 'CPU', 'TPU', 'FPGA'],
    },
    [ModelCatalogStringFilterKey.USE_CASE]: {
      type: 'string',
      values: [
        UseCaseOptionValue.CHATBOT,
        UseCaseOptionValue.CODE_FIXING,
        UseCaseOptionValue.LONG_RAG,
        UseCaseOptionValue.RAG,
      ],
    },
    [ModelCatalogNumberFilterKey.MIN_RPS]: {
      type: 'number',
      range: {
        min: 1,
        max: 300,
      },
    },
    // All latency metric combinations for dropdown options
    ttft_mean: {
      type: 'number' as const,
      range: { min: 20, max: 893 },
    },
    ttft_p90: {
      type: 'number' as const,
      range: { min: 25, max: 600 },
    },
    ttft_p95: {
      type: 'number' as const,
      range: { min: 30, max: 700 },
    },
    ttft_p99: {
      type: 'number' as const,
      range: { min: 40, max: 893 },
    },
    e2e_mean: {
      type: 'number' as const,
      range: { min: 50, max: 800 },
    },
    e2e_p90: {
      type: 'number' as const,
      range: { min: 60, max: 900 },
    },
    e2e_p95: {
      type: 'number' as const,
      range: { min: 70, max: 1000 },
    },
    e2e_p99: {
      type: 'number' as const,
      range: { min: 80, max: 1200 },
    },
    tps_mean: {
      type: 'number' as const,
      range: { min: 10, max: 300 },
    },
    tps_p90: {
      type: 'number' as const,
      range: { min: 15, max: 350 },
    },
    tps_p95: {
      type: 'number' as const,
      range: { min: 20, max: 400 },
    },
    tps_p99: {
      type: 'number' as const,
      range: { min: 25, max: 500 },
    },
    itl_mean: {
      type: 'number' as const,
      range: { min: 5, max: 100 },
    },
    itl_p90: {
      type: 'number' as const,
      range: { min: 8, max: 120 },
    },
    itl_p95: {
      type: 'number' as const,
      range: { min: 10, max: 150 },
    },
    itl_p99: {
      type: 'number' as const,
      range: { min: 15, max: 200 },
    },
  },
  namedQueries: mockNamedQueries,
  ...partial,
});

// Mock for artifact-specific filter options (performance artifacts endpoint)
// This is a subset of filters relevant to performance artifacts
export const mockArtifactFilterOptionsList = (
  partial?: Partial<CatalogFilterOptionsList>,
): CatalogFilterOptionsList => {
  const base = mockCatalogFilterOptionsList();
  return {
    ...base,
    filters: {
      ...base.filters,
      [ModelCatalogStringFilterKey.HARDWARE_TYPE]: {
        type: 'string',
        values: ['H100-80', 'A100-80', 'L40S', 'MI300X'],
      },
    },
    namedQueries: mockNamedQueries,
    ...partial,
  };
};
