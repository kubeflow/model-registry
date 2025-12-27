import { APIOptions, handleRestFailures, isModArchResponse, restGET } from 'mod-arch-core';
import {
  McpCatalogSourceList,
  McpFilterOptionsList,
  McpServer,
  McpServerList,
} from '~/app/pages/mcpCatalog/types';

/**
 * Get all MCP servers from the catalog
 * @param hostPath - The base URL for the MCP catalog API
 * @param queryParams - Additional query parameters to include in the request
 * @returns A function that fetches MCP servers with optional filters
 */
export const getMcpServers =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (
    opts: APIOptions,
    sourceLabel?: string,
    pageSize?: number,
    filterQuery?: string,
    searchTerm?: string,
  ): Promise<McpServerList> => {
    const params = { ...queryParams };
    if (sourceLabel) {
      params.sourceLabel = sourceLabel;
    }
    if (pageSize) {
      params.pageSize = pageSize;
    }
    if (filterQuery) {
      params.filterQuery = filterQuery;
    }
    if (searchTerm) {
      params.q = searchTerm;
    }
    return handleRestFailures(restGET(hostPath, '/mcp_servers', params, opts)).then((response) => {
      if (isModArchResponse<McpServerList>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });
  };

/**
 * Get a specific MCP server by ID
 */
export const getMcpServer =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions, serverId: string): Promise<McpServer> =>
    handleRestFailures(
      restGET(hostPath, `/mcp_servers/server/${serverId}`, queryParams, opts),
    ).then((response) => {
      if (isModArchResponse<McpServer>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

/**
 * Get all MCP catalog sources
 * Filters by assetType=mcp_servers to only return MCP server sources
 */
export const getMcpSources =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions): Promise<McpCatalogSourceList> =>
    handleRestFailures(
      restGET(hostPath, '/sources', { ...queryParams, assetType: 'mcp_servers' }, opts),
    ).then((response) => {
      if (isModArchResponse<McpCatalogSourceList>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

/**
 * Get filter options for MCP servers
 * Returns available values for each filterable field
 */
export const getMcpFilterOptions =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions): Promise<McpFilterOptionsList> =>
    handleRestFailures(restGET(hostPath, '/mcp_servers/filter_options', queryParams, opts)).then(
      (response) => {
        if (isModArchResponse<McpFilterOptionsList>(response)) {
          return response.data;
        }
        throw new Error('Invalid response format');
      },
    );
