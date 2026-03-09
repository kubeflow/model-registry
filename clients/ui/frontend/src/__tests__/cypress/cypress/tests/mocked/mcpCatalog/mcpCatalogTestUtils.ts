import { mockModArchResponse } from 'mod-arch-core';
import { mockMcpServers } from '~/app/pages/mcpCatalog/mocks/mockMcpServers';
import { mockMcpCatalogFilterOptions } from '~/app/pages/mcpCatalog/mocks/mockMcpCatalogFilterOptions';
import { mockCatalogSource, mockCatalogSourceList } from '~/__mocks__';
import { MODEL_CATALOG_API_VERSION } from '~/__tests__/cypress/cypress/support/commands/api';

export const MCP_SERVERS_RESPONSE = {
  items: mockMcpServers,
  size: mockMcpServers.length,
  pageSize: 10,
  nextPageToken: '',
};

export const MCP_SERVERS_PATH = `/model-registry/api/${MODEL_CATALOG_API_VERSION}/model_catalog/mcp_servers`;

export const MCP_FILTER_OPTIONS_PATH = `/model-registry/api/${MODEL_CATALOG_API_VERSION}/model_catalog/mcp_servers_filter_options`;

export const initMcpCatalogIntercepts = (): void => {
  cy.intercept('GET', '*mcp_servers*', mockModArchResponse(MCP_SERVERS_RESPONSE));
  cy.intercept(
    'GET',
    `**/api/${MODEL_CATALOG_API_VERSION}/model_catalog/mcp_servers*`,
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
    mockModArchResponse(mockMcpCatalogFilterOptions),
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
  cy.intercept('GET', `**/mcp_servers/${server.id}`, mockModArchResponse(server));
};
