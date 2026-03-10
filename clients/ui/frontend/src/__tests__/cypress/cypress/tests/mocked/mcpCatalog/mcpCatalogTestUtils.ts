import { mockModArchResponse } from 'mod-arch-core';
import {
  mockCatalogSource,
  mockCatalogSourceList,
  testFilterOptions,
  testMcpServers,
} from '~/__mocks__';
import { MODEL_CATALOG_API_VERSION } from '~/__tests__/cypress/cypress/support/commands/api';

export { testMcpServers, testFilterOptions };

export const MCP_SERVERS_RESPONSE = {
  items: testMcpServers,
  size: testMcpServers.length,
  pageSize: 10,
  nextPageToken: '',
};

export const MCP_SERVERS_PATH = `/model-registry/api/${MODEL_CATALOG_API_VERSION}/model_catalog/mcp_servers`;

export const MCP_FILTER_OPTIONS_PATH = `/model-registry/api/${MODEL_CATALOG_API_VERSION}/model_catalog/mcp_servers_filter_options`;

export const initMcpCatalogIntercepts = (): void => {
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
    mockModArchResponse(testFilterOptions),
  );
};

export const initServerDetailIntercept = (server: (typeof testMcpServers)[number]): void => {
  cy.intercept(
    {
      method: 'GET',
      pathname: `${MCP_SERVERS_PATH}/${server.id}`,
    },
    mockModArchResponse(server),
  );
};
