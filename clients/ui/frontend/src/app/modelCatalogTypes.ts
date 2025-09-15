import { APIOptions } from 'mod-arch-core';
import { ModelRegistryCustomProperties } from './types';

export type CatalogSource = {
  id: string;
  name: string;
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

export type CatalogModelArtifact = {
  createTimeSinceEpoch: string;
  lastUpdateTimeSinceEpoch: string;
  uri: string;
  customProperties: ModelRegistryCustomProperties;
};

export type CatalogModelArtifactList = ModelCatalogListParams & { items: CatalogModelArtifact[] };

export type GetCatalogModelsBySource = (
  opts: APIOptions,
  sourceId: string,
  paginationParams?: {
    pageSize?: string;
    nextPageToken?: string;
    orderBy?: string;
    sortOrder?: string;
  },
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
) => Promise<CatalogModelArtifactList>;

export type ModelCatalogAPIs = {
  getCatalogModelsBySource: GetCatalogModelsBySource;
  getListSources: GetListSources;
  getCatalogModel: GetCatalogModel;
  getListCatalogModelArtifacts: GetListCatalogModelArtifacts;
};

export type CatalogModelDetailsParams = {
  sourceId?: string;
  repositoryName?: string;
  modelName?: string;
};
