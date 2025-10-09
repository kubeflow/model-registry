import { APIOptions, handleRestFailures, isModArchResponse, restGET } from 'mod-arch-core';
import {
  CatalogArtifactList,
  CatalogFilterOptionsList,
  CatalogModel,
  CatalogModelList,
  CatalogSourceList,
  ModelCatalogFilterStates,
} from '~/app/modelCatalogTypes';
import { filtersToFilterQuery } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';

export const getCatalogModelsBySource =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (
    opts: APIOptions,
    sourceId: string,
    paginationParams?: {
      pageSize?: string;
      nextPageToken?: string;
      orderBy?: string;
      sortOrder?: string;
    },
    searchKeyword?: string,
    filterData?: ModelCatalogFilterStates,
    filterOptions?: CatalogFilterOptionsList | null,
  ): Promise<CatalogModelList> => {
    const allParams = {
      source: sourceId,
      ...paginationParams,
      ...(searchKeyword && { q: searchKeyword }),
      ...queryParams,
      ...(filterData &&
        filterOptions && { filterQuery: filtersToFilterQuery(filterData, filterOptions) }),
    };
    return handleRestFailures(restGET(hostPath, '/models', allParams, opts)).then((response) => {
      if (isModArchResponse<CatalogModelList>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });
  };

export const getCatalogFilterOptionList =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions): Promise<CatalogFilterOptionsList> =>
    handleRestFailures(restGET(hostPath, '/models/filter_options', queryParams, opts)).then(
      (response) => {
        if (isModArchResponse<CatalogFilterOptionsList>(response)) {
          return response.data;
        }
        throw new Error('Invalid response format');
      },
    );

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
  (opts: APIOptions, sourceId: string, modelName: string): Promise<CatalogArtifactList> =>
    handleRestFailures(
      restGET(hostPath, `/sources/${sourceId}/artifacts/${modelName}`, queryParams, opts),
    ).then((response) => {
      if (isModArchResponse<CatalogArtifactList>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });
