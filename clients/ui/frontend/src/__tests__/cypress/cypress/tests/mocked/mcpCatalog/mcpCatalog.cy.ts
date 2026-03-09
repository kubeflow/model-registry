import { mockModArchResponse } from 'mod-arch-core';
import { mcpCatalog } from '~/__tests__/cypress/cypress/pages/mcpCatalog';
import { mockCatalogSource, mockCatalogSourceList } from '~/__mocks__';
import { MODEL_CATALOG_API_VERSION } from '~/__tests__/cypress/cypress/support/commands/api';

const testMcpServers = [
  {
    id: 1,
    name: 'Kubernetes',
    description: 'Control and inspect Kubernetes clusters.',
    deploymentMode: 'local',
    securityIndicators: { verifiedSource: true, sast: true },
    source_id: 'sample', // eslint-disable-line camelcase
    toolCount: 0,
  },
  {
    id: 2,
    name: 'GitHub',
    description: 'Integrate with GitHub repositories.',
    deploymentMode: 'remote',
    securityIndicators: { verifiedSource: true, secureEndpoint: true },
    source_id: 'sample', // eslint-disable-line camelcase
    toolCount: 0,
  },
];

const testFilterOptions = {
  filters: {
    deploymentMode: { type: 'string', values: ['Remote', 'Local'] },
    supportedTransports: { type: 'string', values: ['SSE', 'http-streaming'] },
    license: { type: 'string', values: ['MIT', 'Apache-2.0'] },
    labels: {
      type: 'string',
      values: ['kubernetes', 'github', 'database', 'monitoring', 'security', 'automation'],
    },
    securityVerification: {
      type: 'string',
      values: ['Verified source', 'Secure endpoint', 'SAST', 'Read only tools'],
    },
  },
};

const MCP_SERVERS_RESPONSE = {
  items: testMcpServers,
  size: testMcpServers.length,
  pageSize: 10,
  nextPageToken: '',
};

const MCP_SERVERS_PATH = `/model-registry/api/${MODEL_CATALOG_API_VERSION}/model_catalog/mcp_servers`;

const MCP_FILTER_OPTIONS_PATH = `/model-registry/api/${MODEL_CATALOG_API_VERSION}/model_catalog/mcp_servers_filter_options`;

const initMcpCatalogIntercepts = () => {
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
    mockModArchResponse(testFilterOptions),
  );
};

describe('MCP Catalog Page', () => {
  beforeEach(() => {
    initMcpCatalogIntercepts();
  });

  it('MCP Catalog tab should be enabled in nav', () => {
    mcpCatalog.visit();
    mcpCatalog.tabEnabled();
  });

  it('should display page title and description', () => {
    mcpCatalog.visit();
    mcpCatalog.findPageTitle().should('be.visible');
    mcpCatalog.findPageDescription().should('be.visible');
  });

  it('should display MCP server cards', () => {
    mcpCatalog.visit();
    cy.get('[data-testid^="mcp-catalog-card-"]', { timeout: 15000 }).should(
      'have.length.at.least',
      1,
    );
  });

  it('should display sidebar filters', () => {
    mcpCatalog.visit();
    mcpCatalog.findFilter('deploymentMode').should('be.visible');
    mcpCatalog.findFilter('supportedTransports').should('be.visible');
    mcpCatalog.findFilter('license').should('be.visible');
    mcpCatalog.findFilter('labels').should('be.visible');
    mcpCatalog.findFilter('securityVerification').should('be.visible');
  });

  it('should display Deployment mode filter with Local and Remote options', () => {
    mcpCatalog.visit();
    mcpCatalog.findFilterCheckbox('deploymentMode', 'Local').should('be.visible');
    mcpCatalog.findFilterCheckbox('deploymentMode', 'Remote').should('be.visible');
  });

  it('filter Show more should expand labels list', () => {
    mcpCatalog.visit();
    mcpCatalog.findFilterShowMore('labels').scrollIntoView();
    mcpCatalog.findFilterShowMore('labels').click();
    cy.findByTestId('mcp-filter-labels-show-less').scrollIntoView();
    cy.findByTestId('mcp-filter-labels-show-less').should('be.visible');
  });

  it('should display known mock server cards', () => {
    mcpCatalog.visit();
    cy.get('[data-testid="mcp-catalog-card-1"]', { timeout: 15000 }).should('be.visible');
    cy.get('[data-testid="mcp-catalog-card-2"]', { timeout: 15000 }).should('be.visible');
  });
});
