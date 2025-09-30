import { APIOptions } from 'mod-arch-core';
import { ModelRegistryCustomProperties, ModelRegistryCustomProperty } from './types';

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

type PerformanceMetricsCustomProperties = Record<string, ModelRegistryCustomProperty> & {
  config_id?: ModelRegistryCustomProperty;

  ttft_mean?: ModelRegistryCustomProperty;
  ttft_p90?: ModelRegistryCustomProperty;
  ttft_p95?: ModelRegistryCustomProperty;
  ttft_p99?: ModelRegistryCustomProperty;

  e2e_mean?: ModelRegistryCustomProperty;
  e2e_p90?: ModelRegistryCustomProperty;
  e2e_p95?: ModelRegistryCustomProperty;
  e2e_p99?: ModelRegistryCustomProperty;

  tps_mean?: ModelRegistryCustomProperty;
  tps_p90?: ModelRegistryCustomProperty;
  tps_p95?: ModelRegistryCustomProperty;
  tps_p99?: ModelRegistryCustomProperty;

  itl_mean?: ModelRegistryCustomProperty;
  itl_p90?: ModelRegistryCustomProperty;
  itl_p95?: ModelRegistryCustomProperty;
  itl_p99?: ModelRegistryCustomProperty;

  requests_per_second?: ModelRegistryCustomProperty;
  max_input_tokens?: ModelRegistryCustomProperty;
  max_output_tokens?: ModelRegistryCustomProperty;
  mean_input_tokens?: ModelRegistryCustomProperty;
  mean_output_tokens?: ModelRegistryCustomProperty;

  hardware?: ModelRegistryCustomProperty;
  hardware_count?: ModelRegistryCustomProperty;
  framework?: ModelRegistryCustomProperty;
  framework_version?: ModelRegistryCustomProperty;
  docker_image?: ModelRegistryCustomProperty;
  entrypoint?: ModelRegistryCustomProperty;
  inserted_at?: ModelRegistryCustomProperty;
  created_at?: ModelRegistryCustomProperty;
  updated_at?: ModelRegistryCustomProperty;
  model_hf_repo_name?: ModelRegistryCustomProperty;
};

export type CatalogMetricsArtifact = CatalogArtifactBase &
  (
    | {
        artifactType: CatalogArtifactType.metricsArtifact;
        metricsType?: MetricsType.performanceMetrics;
        customProperties?: PerformanceMetricsCustomProperties;
      }
    | {
        artifactType: CatalogArtifactType.metricsArtifact;
        metricsType?: string;
        customProperties?: ModelRegistryCustomProperties;
      }
  );

export type CatalogArtifactBase = {
  createTimeSinceEpoch: string;
  lastUpdateTimeSinceEpoch: string;
};

export type CatalogModelArtifact = CatalogArtifactBase & {
  artifactType: CatalogArtifactType.modelArtifact;
  uri: string;
  customProperties?: ModelRegistryCustomProperties;
};

export type CatalogArtifacts = CatalogModelArtifact | CatalogMetricsArtifact;

export type CatalogArtifactList = ModelCatalogListParams & { items: CatalogArtifacts[] };

export type CatalogFilterOption = {
  type: string;
  range?: {
    max?: number;
    min?: number;
  };
  values?: string[];
};

export type CatalogFilterOptionsList = {
  filters: CatalogFilterOption;
};

export type GetCatalogModelsBySource = (
  opts: APIOptions,
  sourceId: string,
  paginationParams?: {
    pageSize?: string;
    nextPageToken?: string;
    orderBy?: string;
    sortOrder?: string;
  },
  searchKeyword?: string,
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
