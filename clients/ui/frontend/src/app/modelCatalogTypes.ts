import { APIOptions } from 'mod-arch-core';
import { ModelRegistryCustomProperties } from './types';

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

export type CatalogArtifactBase = {
  createTimeSinceEpoch: string;
  lastUpdateTimeSinceEpoch: string;
  customProperties: ModelRegistryCustomProperties;
};

export type CatalogModelArtifact = CatalogArtifactBase & {
  artifactType: CatalogArtifactType.modelArtifact;
  uri: string;
};

export type CatalogMetricsArtifact = CatalogArtifactBase & {
  artifactType: CatalogArtifactType.metricsArtifact;
  metricsType?: string;
};

export type CatalogArtifacts = CatalogModelArtifact | CatalogMetricsArtifact;

export type CatalogArtifactList = ModelCatalogListParams & { items: CatalogArtifacts[] };

export type CatalogFilterOption = {
  type: string;
  range?: {
    max?: number;
    min?: number;
  };
  Values?: string[];
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
