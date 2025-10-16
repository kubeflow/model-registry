import { APIOptions } from 'mod-arch-core';
import {
  ModelCatalogTask,
  ModelCatalogProvider,
  ModelCatalogLicense,
  AllLanguageCode,
  ModelCatalogStringFilterKey,
  ModelCatalogNumberFilterKey,
} from '~/concepts/modelCatalog/const';
import {
  ModelRegistryCustomProperties,
  ModelRegistryCustomPropertyString,
  ModelRegistryCustomPropertyInt,
  ModelRegistryCustomPropertyDouble,
} from './types';

export type CatalogSource = {
  id: string;
  name: string;
  labels: string[];
  enabled?: boolean;
};

export type CatalogSourceList = ModelCatalogListParams & { items: CatalogSource[] };

export type CatalogModel = {
  source_id?: string;
  name: string;
  provider?: string;
  readme?: string;
  maturity?: string;
  language?: string[];
  logo?: string;
  tasks?: string[];
  libraryName?: string;
  license?: string;
  licenseLink?: string;
  description?: string;
  createTimeSinceEpoch?: string;
  lastUpdateTimeSinceEpoch?: string;
  customProperties?: ModelRegistryCustomProperties;
};

export type ModelCatalogListParams = {
  size: number;
  pageSize: number;
  nextPageToken: string;
};

export type CatalogModelList = ModelCatalogListParams & { items: CatalogModel[] };

export enum CatalogArtifactType {
  modelArtifact = 'model-artifact',
  metricsArtifact = 'metrics-artifact',
}

export enum MetricsType {
  accuracyMetrics = 'accuracy-metrics',
  performanceMetrics = 'performance-metrics',
}

export enum CategoryName {
  allModels = 'All models',
  communityAndCustomModels = 'Community and custom',
}

export enum SourceLabel {
  other = 'null',
}

export type CatalogArtifactBase = {
  createTimeSinceEpoch: string;
  lastUpdateTimeSinceEpoch: string;
  customProperties: ModelRegistryCustomProperties;
};

export type CatalogModelArtifact = CatalogArtifactBase & {
  artifactType: CatalogArtifactType.modelArtifact;
  uri: string;
};

export type PerformanceMetricsCustomProperties = {
  config_id?: ModelRegistryCustomPropertyString;
  hardware_type?: ModelRegistryCustomPropertyString;
  hardware_count?: ModelRegistryCustomPropertyInt;
  requests_per_second?: ModelRegistryCustomPropertyDouble;
  // TTFT (Time To First Token) latency metrics
  ttft_mean?: ModelRegistryCustomPropertyDouble;
  ttft_p90?: ModelRegistryCustomPropertyDouble;
  ttft_p95?: ModelRegistryCustomPropertyDouble;
  ttft_p99?: ModelRegistryCustomPropertyDouble;
  // E2E (End-to-End) latency metrics
  e2e_mean?: ModelRegistryCustomPropertyDouble;
  e2e_p90?: ModelRegistryCustomPropertyDouble;
  e2e_p95?: ModelRegistryCustomPropertyDouble;
  e2e_p99?: ModelRegistryCustomPropertyDouble;
  // TPS (Tokens Per Second) latency metrics
  tps_mean?: ModelRegistryCustomPropertyDouble;
  tps_p90?: ModelRegistryCustomPropertyDouble;
  tps_p95?: ModelRegistryCustomPropertyDouble;
  tps_p99?: ModelRegistryCustomPropertyDouble;
  // ITL (Inter-Token Latency) metrics
  itl_mean?: ModelRegistryCustomPropertyDouble;
  itl_p90?: ModelRegistryCustomPropertyDouble;
  itl_p95?: ModelRegistryCustomPropertyDouble;
  itl_p99?: ModelRegistryCustomPropertyDouble;
  // Token metrics
  max_input_tokens?: ModelRegistryCustomPropertyDouble;
  max_output_tokens?: ModelRegistryCustomPropertyDouble;
  mean_input_tokens?: ModelRegistryCustomPropertyDouble;
  mean_output_tokens?: ModelRegistryCustomPropertyDouble;
  // Framework information
  framework?: ModelRegistryCustomPropertyString;
  framework_version?: ModelRegistryCustomPropertyString;
  // Additional fields from ADR (excluded from display per requirements)
  docker_image?: ModelRegistryCustomPropertyString;
  entrypoint?: ModelRegistryCustomPropertyString;
  inserted_at?: ModelRegistryCustomPropertyString;
  created_at?: ModelRegistryCustomPropertyString;
  updated_at?: ModelRegistryCustomPropertyString;
  model_hf_repo_name?: ModelRegistryCustomPropertyString;
  scenario_id?: ModelRegistryCustomPropertyString;
};

export type AccuracyMetricsCustomProperties = {
  // overall_average?: ModelRegistryCustomPropertyDouble; // NOTE: overall_average is currently omitted from the API and will be restored
  arc_v1?: ModelRegistryCustomPropertyDouble;
} & Record<string, ModelRegistryCustomPropertyDouble>;

export type CatalogPerformanceMetricsArtifact = Omit<CatalogArtifactBase, 'customProperties'> & {
  artifactType: CatalogArtifactType.metricsArtifact;
  metricsType: MetricsType.performanceMetrics;
  customProperties: PerformanceMetricsCustomProperties;
};

export type CatalogAccuracyMetricsArtifact = Omit<CatalogArtifactBase, 'customProperties'> & {
  artifactType: CatalogArtifactType.metricsArtifact;
  metricsType: MetricsType.accuracyMetrics;
  customProperties: AccuracyMetricsCustomProperties;
};

export type CatalogMetricsArtifact =
  | CatalogPerformanceMetricsArtifact
  | CatalogAccuracyMetricsArtifact;

export type CatalogArtifacts = CatalogModelArtifact | CatalogMetricsArtifact;

export type CatalogArtifactList = ModelCatalogListParams & { items: CatalogArtifacts[] };

export type CatalogFilterNumberOption = {
  type: 'number';
  range: {
    max: number;
    min: number;
  };
};

export type CatalogFilterStringOption<T extends string> = {
  type: 'string';
  values: T[];
};

export type GetCatalogModelsBySource = (
  opts: APIOptions,
  sourceId?: string,
  sourceLabel?: string,
  paginationParams?: {
    pageSize?: string;
    nextPageToken?: string;
    orderBy?: string;
    sortOrder?: string;
  },
  searchKeyword?: string,
  filterData?: ModelCatalogFilterStates,
  filterOptions?: CatalogFilterOptionsList | null,
) => Promise<CatalogModelList>;

export type GetListSources = (opts: APIOptions) => Promise<CatalogSourceList>;

export type GetCatalogModel = (
  opts: APIOptions,
  sourceId: string,
  modelName: string,
) => Promise<CatalogModel>;

export type GetListCatalogModelArtifacts = (
  opts: APIOptions,
  sourceId: string,
  modelName: string,
) => Promise<CatalogArtifactList>;

export type GetCatalogFilterOptionList = (opts: APIOptions) => Promise<CatalogFilterOptionsList>;

export type ModelCatalogAPIs = {
  getCatalogModelsBySource: GetCatalogModelsBySource;
  getListSources: GetListSources;
  getCatalogModel: GetCatalogModel;
  getListCatalogModelArtifacts: GetListCatalogModelArtifacts;
  getCatalogFilterOptionList: GetCatalogFilterOptionList;
};

export type CatalogModelDetailsParams = {
  sourceId?: string;
  modelName?: string;
};

export type ModelCatalogFilterKey = ModelCatalogStringFilterKey | ModelCatalogNumberFilterKey;

// Not used for a run time value, just for mapping other types
export type ModelCatalogStringFilterValueType = {
  [ModelCatalogStringFilterKey.TASK]: ModelCatalogTask;
  [ModelCatalogStringFilterKey.PROVIDER]: ModelCatalogProvider;
  [ModelCatalogStringFilterKey.LICENSE]: ModelCatalogLicense;
  [ModelCatalogStringFilterKey.LANGUAGE]: AllLanguageCode;
  [ModelCatalogStringFilterKey.HARDWARE_TYPE]: string;
};

export type ModelCatalogStringFilterOptions = {
  [key in ModelCatalogStringFilterKey]: CatalogFilterStringOption<
    ModelCatalogStringFilterValueType[key]
  >;
};

export type CatalogFilterOptions = ModelCatalogStringFilterOptions & {
  [key in ModelCatalogNumberFilterKey]: CatalogFilterNumberOption;
};

export type CatalogFilterOptionsList = {
  filters: CatalogFilterOptions;
};

export type ModelCatalogFilterStates = {
  [key in ModelCatalogStringFilterKey]: ModelCatalogStringFilterValueType[key][];
} & {
  [key in ModelCatalogNumberFilterKey]: number | undefined;
};
