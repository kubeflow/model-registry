import {
  APIOptions,
  assembleModArchBody,
  handleRestFailures,
  isModArchResponse,
  restCREATE,
  restDELETE,
  restGET,
  restPATCH,
} from 'mod-arch-core';
import {
  CatalogSourceConfig,
  CatalogSourceConfigList,
  CatalogSourceConfigPayload,
  CatalogSourcePreviewRequest,
  CatalogSourcePreviewResult,
} from '~/app/modelCatalogTypes';

export const getCatalogSourceConfigs =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions): Promise<CatalogSourceConfigList> =>
    handleRestFailures(restGET(hostPath, '/source_configs', queryParams, opts)).then((response) => {
      if (isModArchResponse<CatalogSourceConfigList>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const createCatalogSourceConfig =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions, data: CatalogSourceConfigPayload): Promise<CatalogSourceConfig> =>
    handleRestFailures(
      restCREATE(hostPath, '/source_configs', assembleModArchBody(data), queryParams, opts),
    ).then((response) => {
      if (isModArchResponse<CatalogSourceConfig>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const getCatalogSourceConfig =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions, sourceId: string): Promise<CatalogSourceConfig> =>
    handleRestFailures(restGET(hostPath, `/source_configs/${sourceId}`, queryParams, opts)).then(
      (response) => {
        if (isModArchResponse<CatalogSourceConfig>(response)) {
          return response.data;
        }
        throw new Error('Invalid response format');
      },
    );

export const updateCatalogSourceConfig =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (
    opts: APIOptions,
    sourceId: string,
    data: Partial<CatalogSourceConfigPayload>,
  ): Promise<CatalogSourceConfig> =>
    handleRestFailures(
      restPATCH(
        hostPath,
        `/source_configs/${sourceId}`,
        assembleModArchBody(data),
        queryParams,
        opts,
      ),
    ).then((response) => {
      if (isModArchResponse<CatalogSourceConfig>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

export const deleteCatalogSourceConfig =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions, sourceId: string): Promise<void> =>
    handleRestFailures(restDELETE(hostPath, `/source_configs/${sourceId}`, {}, queryParams, opts));

export const previewCatalogSource =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions, data: CatalogSourcePreviewRequest): Promise<CatalogSourcePreviewResult> =>
    handleRestFailures(
      restCREATE(hostPath, '/source_preview', assembleModArchBody(data), queryParams, opts),
    ).then((response) => {
      if (isModArchResponse<CatalogSourcePreviewResult>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });
