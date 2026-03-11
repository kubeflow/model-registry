import { mockModArchResponse } from 'mod-arch-core';
import {
  mockCatalogSource,
  mockCatalogSourceList,
  mockMcpCatalogFilterOptions,
  mockMcpServers,
} from '~/__mocks__';
import { MODEL_CATALOG_API_VERSION } from '~/__tests__/cypress/cypress/support/commands/api';

export { mockMcpServers, mockMcpCatalogFilterOptions };

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
