import { APIOptions, handleRestFailures, isModArchResponse, restGET } from 'mod-arch-core';
import { McpServer, McpServerList, McpToolList } from '~/app/mcpServerCatalogTypes';
import { CatalogFilterOptionsList } from '~/app/modelCatalogTypes';

export const getMcpServerList =
  (hostPath: string, queryParams: Record<string, unknown> = {}) =>
  (opts: APIOptions): Promise<McpServerList> =>
    handleRestFailures(restGET(hostPath, '/mcp_servers', queryParams, opts)).then((response) => {
      if (isModArchResponse<McpServerList>(response)) {
        return response.data;
      }
      throw new Error('Invalid response format');
    });

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
