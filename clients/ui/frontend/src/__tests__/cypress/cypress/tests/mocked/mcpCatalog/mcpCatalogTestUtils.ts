import { mockModArchResponse } from 'mod-arch-core';
import {
  mockCatalogSource,
  mockCatalogSourceList,
  mockMcpCatalogFilterOptions,
  mockMcpServers,
  mockMcpToolWithServer,
  mockMcpToolList,
} from '~/__mocks__';
import type { McpToolList } from '~/app/mcpServerCatalogTypes';
import { MODEL_CATALOG_API_VERSION } from '~/__tests__/cypress/cypress/support/commands/api';

export { mockMcpServers, mockMcpCatalogFilterOptions, mockMcpToolWithServer, mockMcpToolList };

export const MCP_SERVERS_RESPONSE = {
  items: mockMcpServers,
  size: mockMcpServers.length,
  pageSize: 10,
  nextPageToken: '',
};

export const MCP_SERVERS_PATH = `/model-registry/api/${MODEL_CATALOG_API_VERSION}/mcp_catalog/mcp_servers`;

export const MCP_FILTER_OPTIONS_PATH = `/model-registry/api/${MODEL_CATALOG_API_VERSION}/mcp_catalog/mcp_servers_filter_options`;

export const initMcpCatalogIntercepts = (): void => {
  cy.intercept('GET', '*mcp_servers*', mockModArchResponse(MCP_SERVERS_RESPONSE));
  cy.intercept(
    'GET',
    `**/api/${MODEL_CATALOG_API_VERSION}/mcp_catalog/mcp_servers*`,
    mockModArchResponse(MCP_SERVERS_RESPONSE),
  );
  cy.intercept(
    { method: 'GET', pathname: MCP_SERVERS_PATH },
    mockModArchResponse(MCP_SERVERS_RESPONSE),
  );
  cy.interceptApi(
    `GET /api/:apiVersion/model_catalog/sources`,
    { path: { apiVersion: MODEL_CATALOG_API_VERSION } },
    mockCatalogSourceList({ items: [mockCatalogSource({})] }),
  );
  cy.intercept(
    { method: 'GET', pathname: MCP_FILTER_OPTIONS_PATH },
    mockModArchResponse(mockMcpCatalogFilterOptions()),
  );
};

export const initServerDetailIntercept = (server: (typeof mockMcpServers)[number]): void => {
  cy.intercept(
    {
      method: 'GET',
      pathname: `${MCP_SERVERS_PATH}/${server.id}`,
    },
    mockModArchResponse(server),
  );
};

export const initServerToolsIntercept = (serverId: string, toolList: McpToolList): void => {
  cy.intercept(
    {
      method: 'GET',
      pathname: `${MCP_SERVERS_PATH}/${serverId}/tools`,
    },
    mockModArchResponse(toolList),
  );
};

export const initServerToolsErrorIntercept = (serverId: string): void => {
  cy.intercept(
    {
      method: 'GET',
      pathname: `${MCP_SERVERS_PATH}/${serverId}/tools`,
    },
    { statusCode: 500, body: { error: 'Internal server error' } },
  );
};
