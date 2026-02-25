import { APIOptions, handleRestFailures, isModArchResponse, restGET } from 'mod-arch-core';
import {
  McpServer,
  McpServerList,
  McpServerListParams,
  McpToolList,
} from '~/app/mcpServerCatalogTypes';
import { CatalogFilterOptionsList } from '~/app/modelCatalogTypes';

export const getMcpServerList =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions, listParams?: McpServerListParams): Promise<McpServerList> => {
    const pageSize = listParams?.pageSize !== undefined ? String(listParams.pageSize) : undefined;
    const allParams = {
      ...queryParams,
      ...(listParams?.sourceLabel !== undefined && { sourceLabel: listParams.sourceLabel }),
      ...(listParams?.nextPageToken !== undefined && { nextPageToken: listParams.nextPageToken }),
      ...(pageSize !== undefined && { pageSize }),
      ...(listParams?.filterQuery !== undefined &&
        listParams.filterQuery !== '' && { filterQuery: listParams.filterQuery }),
      ...(listParams?.namedQuery !== undefined &&
        listParams.namedQuery !== '' && { namedQuery: listParams.namedQuery }),
      ...(listParams?.includeTools !== undefined && { includeTools: listParams.includeTools }),
      ...(listParams?.toolLimit !== undefined && { toolLimit: listParams.toolLimit }),
      ...(listParams?.orderBy !== undefined &&
        listParams.orderBy !== '' && { orderBy: listParams.orderBy }),
      ...(listParams?.sortOrder !== undefined &&
        listParams.sortOrder !== '' && { sortOrder: listParams.sortOrder }),
      ...(listParams?.name !== undefined && listParams.name !== '' && { name: listParams.name }),
      ...(listParams?.q !== undefined && listParams.q !== '' && { q: listParams.q }),
    };
    return handleRestFailures(restGET(hostPath, '/mcp_servers', allParams, opts)).then(
      (response) => {
        if (isModArchResponse<McpServerList>(response)) {
          return response.data;
        }
        throw new Error('Invalid response format');
      },
    );
  };

export const getMcpServerFilterOptionList =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions): Promise<CatalogFilterOptionsList> =>
    handleRestFailures(restGET(hostPath, '/mcp_servers_filter_options', queryParams, opts)).then(
      (response) => {
        if (isModArchResponse<CatalogFilterOptionsList>(response)) {
          return response.data;
        }
        throw new Error('Invalid response format');
      },
    );

export const getMcpServer =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions, serverId: string): Promise<McpServer> =>
    handleRestFailures(restGET(hostPath, `/mcp_servers/${serverId}`, queryParams, opts)).then(
      (response) => {
        if (isModArchResponse<McpServer>(response)) {
          return response.data;
        }
        throw new Error('Invalid response format');
      },
    );

export const getMcpServerToolList =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions, serverId: string): Promise<McpToolList> =>
    handleRestFailures(restGET(hostPath, `/mcp_servers/${serverId}/tools`, queryParams, opts)).then(
      (response) => {
        if (isModArchResponse<McpToolList>(response)) {
          return response.data;
        }
        throw new Error('Invalid response format');
      },
    );
