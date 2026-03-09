import { mockModArchResponse } from 'mod-arch-core';
import { mcpCatalog, mcpServerDetails } from '~/__tests__/cypress/cypress/pages/mcpCatalog';
import { mockMcpServers } from '~/app/pages/mcpCatalog/mocks/mockMcpServers';
import { mockMcpCatalogFilterOptions } from '~/app/pages/mcpCatalog/mocks/mockMcpCatalogFilterOptions';
import { mockCatalogSource, mockCatalogSourceList } from '~/__mocks__';
import { MODEL_CATALOG_API_VERSION } from '~/__tests__/cypress/cypress/support/commands/api';

const MCP_SERVERS_RESPONSE = {
  items: mockMcpServers,
  size: mockMcpServers.length,
  pageSize: 10,
  nextPageToken: '',
};

const kubernetesServer = mockMcpServers[0];
const customServer = mockMcpServers[4];

const initMcpCatalogIntercepts = () => {
  cy.intercept('GET', '*mcp_servers*', mockModArchResponse(MCP_SERVERS_RESPONSE));
  cy.intercept(
    'GET',
    `**/api/${MODEL_CATALOG_API_VERSION}/model_catalog/mcp_servers*`,
    mockModArchResponse(MCP_SERVERS_RESPONSE),
  );
  cy.intercept(
    {
      method: 'GET',
      pathname: `/model-registry/api/${MODEL_CATALOG_API_VERSION}/model_catalog/mcp_servers`,
    },
    mockModArchResponse(MCP_SERVERS_RESPONSE),
  );
  cy.interceptApi(
    `GET /api/:apiVersion/model_catalog/sources`,
    { path: { apiVersion: MODEL_CATALOG_API_VERSION } },
    mockCatalogSourceList({ items: [mockCatalogSource({})] }),
  );
  cy.intercept(
    {
      method: 'GET',
      pathname: `/model-registry/api/${MODEL_CATALOG_API_VERSION}/model_catalog/mcp_servers_filter_options`,
    },
    mockModArchResponse(mockMcpCatalogFilterOptions),
  );
};

const initServerDetailIntercept = (server: (typeof mockMcpServers)[number]) => {
  cy.intercept(
    {
      method: 'GET',
      pathname: `/model-registry/api/${MODEL_CATALOG_API_VERSION}/model_catalog/mcp_servers/${server.id}`,
    },
    mockModArchResponse(server),
  );
  cy.intercept('GET', `**/mcp_servers/${server.id}`, mockModArchResponse(server));
};

describe('MCP Server Details Page', () => {
  beforeEach(() => {
    initMcpCatalogIntercepts();
  });

  describe('Navigation from catalog', () => {
    it('should navigate to details page when clicking server card link', () => {
      mcpCatalog.visit();
      cy.get(`[data-testid="mcp-catalog-card-detail-link-1"]`, { timeout: 15000 }).should(
        'be.visible',
      );
      mcpCatalog.findCardDetailLink('1').click();
      cy.url().should('include', '/mcp-catalog/1');
    });
  });

  describe('Breadcrumb navigation', () => {
    beforeEach(() => {
      initServerDetailIntercept(kubernetesServer);
    });

    it('should display breadcrumb with MCP Catalog link and server name', () => {
      mcpServerDetails.visit(String(kubernetesServer.id));
      mcpServerDetails.findBreadcrumbCatalogLink().should('be.visible');
      mcpServerDetails.findBreadcrumbServerName().should('contain.text', kubernetesServer.name);
    });

    it('should navigate back to catalog when clicking breadcrumb link', () => {
      mcpServerDetails.visit(String(kubernetesServer.id));
      mcpServerDetails.findBreadcrumbCatalogLink().click();
      cy.url().should('include', '/mcp-catalog');
      cy.url().should('not.include', `/${kubernetesServer.id}`);
    });
  });

  describe('Server header', () => {
    beforeEach(() => {
      initServerDetailIntercept(kubernetesServer);
    });

    it('should display server name in the page title', () => {
      mcpServerDetails.visit(String(kubernetesServer.id));
      cy.findByTestId('app-page-title').should('contain.text', kubernetesServer.name);
    });

    it('should display Deploy MCP Server button', () => {
      mcpServerDetails.visit(String(kubernetesServer.id));
      mcpServerDetails.findDeployButton().should('be.visible');
      mcpServerDetails.findDeployButton().should('contain.text', 'Deploy MCP Server');
    });
  });

  describe('Description card', () => {
    beforeEach(() => {
      initServerDetailIntercept(kubernetesServer);
    });

    it('should display the server description', () => {
      mcpServerDetails.visit(String(kubernetesServer.id));
      mcpServerDetails.findDescription().should('contain.text', kubernetesServer.description);
    });
  });

  describe('README card', () => {
    it('should render README markdown content', () => {
      initServerDetailIntercept(kubernetesServer);
      mcpServerDetails.visit(String(kubernetesServer.id));
      mcpServerDetails.findReadmeMarkdown().should('be.visible');
      mcpServerDetails.findReadmeMarkdown().should('contain.text', 'Kubernetes MCP Server');
    });

    it('should display empty state when no README is available', () => {
      initServerDetailIntercept(customServer);
      mcpServerDetails.visit(String(customServer.id));
      mcpServerDetails.findNoReadme().should('be.visible');
      mcpServerDetails.findNoReadme().should('contain.text', 'No README available');
    });
  });

  describe('Server details sidebar', () => {
    beforeEach(() => {
      initServerDetailIntercept(kubernetesServer);
    });

    it('should display labels with truncation', () => {
      mcpServerDetails.visit(String(kubernetesServer.id));
      mcpServerDetails.findLabels().should('have.length.at.least', 1);
    });

    it('should display license as external link', () => {
      mcpServerDetails.visit(String(kubernetesServer.id));
      mcpServerDetails.findLicenseLink().should('be.visible');
      mcpServerDetails.findLicenseLink().should('contain.text', kubernetesServer.license);
    });

    it('should display version', () => {
      mcpServerDetails.visit(String(kubernetesServer.id));
      mcpServerDetails.findVersion().should('contain.text', kubernetesServer.version);
    });

    it('should display deployment mode', () => {
      mcpServerDetails.visit(String(kubernetesServer.id));
      mcpServerDetails.findDeploymentMode().should('contain.text', 'Local to cluster');
    });

    it('should display transport type', () => {
      mcpServerDetails.visit(String(kubernetesServer.id));
      mcpServerDetails.findTransportType().should('contain.text', 'http-streaming');
    });

    it('should display provider', () => {
      mcpServerDetails.visit(String(kubernetesServer.id));
      mcpServerDetails.findProvider().should('contain.text', kubernetesServer.provider);
    });

    it('should display source code link', () => {
      mcpServerDetails.visit(String(kubernetesServer.id));
      mcpServerDetails.findSourceCodeLink().should('be.visible');
    });
  });

  describe('Error handling', () => {
    it('should show error state for invalid server ID', () => {
      cy.intercept(
        {
          method: 'GET',
          url: '**/mcp_servers/999*',
        },
        { statusCode: 404, body: { error: 'Not found' } },
      );
      cy.visit('/mcp-catalog/999');
      cy.contains('Details not found').should('be.visible');
    });
  });

  describe('Browser navigation', () => {
    beforeEach(() => {
      initServerDetailIntercept(kubernetesServer);
    });

    it('should support browser back navigation', () => {
      mcpCatalog.visit();
      cy.get(`[data-testid="mcp-catalog-card-detail-link-1"]`, { timeout: 15000 }).should(
        'be.visible',
      );
      mcpCatalog.findCardDetailLink('1').click();
      cy.url().should('include', '/mcp-catalog/1');
      cy.go('back');
      cy.url().should('eq', `${Cypress.config().baseUrl}/mcp-catalog`);
    });
  });
});
