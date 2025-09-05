import { APIOptions, handleRestFailures, isModArchResponse, restGET } from 'mod-arch-core';
import {
  CatalogModel,
  CatalogModelArtifactList,
  CatalogModelList,
  CatalogSourceList,
} from '~/app/modelCatalogTypes';

export const getCatalogModelsBySource =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions, sourceId: string): Promise<CatalogModelList> =>
    handleRestFailures(
      restGET(hostPath, '/models', { source: sourceId, ...queryParams }, opts),
    ).then((response) => {
      if (isModArchResponse<CatalogModelList>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const getListSources =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions): Promise<CatalogSourceList> =>
    handleRestFailures(restGET(hostPath, '/sources', queryParams, opts)).then((response) => {
      if (isModArchResponse<CatalogSourceList>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const getCatalogModel =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions, sourceId: string, modelName: string): Promise<CatalogModel> =>
    handleRestFailures(
      restGET(hostPath, `/sources/${sourceId}/models/${modelName}`, queryParams, opts),
    ).then((response) => {
      if (isModArchResponse<CatalogModel>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const getListCatalogModelArtifacts =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions, sourceId: string, modelName: string): Promise<CatalogModelArtifactList> =>
    handleRestFailures(
      restGET(hostPath, `/sources/${sourceId}/models/${modelName}/artifacts`, queryParams, opts),
    ).then((response) => {
      if (isModArchResponse<CatalogModelArtifactList>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });
