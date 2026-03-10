import { mockModArchResponse } from 'mod-arch-core';
import { mockCatalogSource, mockCatalogSourceList } from '~/__mocks__';
import { mcpCatalog } from '~/__tests__/cypress/cypress/pages/mcpCatalog';
import { MODEL_CATALOG_API_VERSION } from '~/__tests__/cypress/cypress/support/commands/api';
import {
  initMcpCatalogIntercepts,
  MCP_FILTER_OPTIONS_PATH,
  MCP_SERVERS_PATH,
  testFilterOptions,
} from './mcpCatalogTestUtils';

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
    mcpCatalog.findMcpCatalogCards().should('have.length.at.least', 1, { timeout: 15000 });
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
    mcpCatalog.findMcpCatalogCard('1').should('be.visible', { timeout: 15000 });
    mcpCatalog.findMcpCatalogCard('2').should('be.visible', { timeout: 15000 });
  });

  it('card name should be a clickable link to details page', () => {
    mcpCatalog.visit();
    mcpCatalog
      .findCardDetailLink('1')
      .should('be.visible', { timeout: 15000 })
      .invoke('attr', 'href')
      .should('include', '/mcp-catalog/');
  });

  it('card description should be truncated', () => {
    mcpCatalog.visit();
    mcpCatalog
      .findCardDescription('1')
      .should('be.visible', { timeout: 15000 })
      .and('have.css', '-webkit-line-clamp', '4');
  });
});

describe('MCP Catalog Empty State', () => {
  it('should show empty state with Reset filters when no results', () => {
    cy.intercept(
      { method: 'GET', pathname: MCP_SERVERS_PATH },
      mockModArchResponse({ items: [], size: 0, pageSize: 10, nextPageToken: '' }),
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
    mcpCatalog.visit();
    cy.findByTestId('mcp-catalog-empty-search', { timeout: 15000 }).should('be.visible');
    cy.contains('No servers found').should('be.visible');
    cy.findByTestId('mcp-catalog-reset-filters')
      .should('be.visible')
      .and('contain', 'Reset filters');
  });
});

describe('MCP Catalog Error State', () => {
  it('should show error state with Retry button on load failure', () => {
    cy.intercept(
      { method: 'GET', pathname: MCP_SERVERS_PATH },
      { statusCode: 500, body: 'Internal Server Error' },
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
    mcpCatalog.visit();
    cy.findByTestId('mcp-catalog-load-error', { timeout: 15000 }).should('be.visible');
    cy.contains('Unable to load MCP servers').should('be.visible');
    cy.findByTestId('mcp-catalog-retry').should('be.visible').and('contain', 'Retry');
  });
});
