import { mcpCatalog, mcpServerDetails } from '~/__tests__/cypress/cypress/pages/mcpCatalog';
import { mockMcpServers } from '~/app/pages/mcpCatalog/mocks/mockMcpServers';
import { initMcpCatalogIntercepts, initServerDetailIntercept } from './mcpCatalogTestUtils';

const kubernetesServer = mockMcpServers[0];
const customServer = mockMcpServers[4];

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

  describe('Server header and description', () => {
    beforeEach(() => {
      initServerDetailIntercept(kubernetesServer);
    });

    it('should display server name, deploy button, and description', () => {
      mcpServerDetails.visit(String(kubernetesServer.id));

      cy.findByTestId('app-page-title').should('contain.text', kubernetesServer.name);

      mcpServerDetails.findDeployButton().should('be.visible');
      mcpServerDetails.findDeployButton().should('contain.text', 'Deploy MCP Server');

      mcpServerDetails.findDescription().should('contain.text', kubernetesServer.description);
    });
  });

  describe('README card', () => {
    it('should render README with markdown elements', () => {
      initServerDetailIntercept(kubernetesServer);
      mcpServerDetails.visit(String(kubernetesServer.id));
      mcpServerDetails.findReadmeMarkdown().should('be.visible');
      mcpServerDetails.findReadmeMarkdown().should('contain.text', 'Kubernetes MCP Server');
      mcpServerDetails.findReadmeMarkdown().find('h3').should('exist');
      mcpServerDetails.findReadmeMarkdown().find('code').should('exist');
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

    it('should display labels, license, version, and deployment mode', () => {
      mcpServerDetails.visit(String(kubernetesServer.id));

      mcpServerDetails.findLabels().should('have.length.at.least', 1);

      mcpServerDetails.findLicenseLink().should('be.visible');
      mcpServerDetails.findLicenseLink().should('contain.text', kubernetesServer.license);

      mcpServerDetails.findVersion().should('contain.text', kubernetesServer.version);

      mcpServerDetails.findDeploymentMode().should('contain.text', 'Local to cluster');
    });

    it('should display artifacts, source code, provider, and transport type', () => {
      mcpServerDetails.visit(String(kubernetesServer.id));

      mcpServerDetails.findArtifactCopy().should('be.visible');
      mcpServerDetails
        .findArtifactCopy()
        .first()
        .find('input')
        .should('have.value', kubernetesServer.artifacts![0].uri);

      mcpServerDetails.findSourceCodeLink().should('be.visible');

      mcpServerDetails.findProvider().should('contain.text', kubernetesServer.provider);

      mcpServerDetails.findTransportType().should('contain.text', 'http-streaming');
    });
  });

  describe('Error handling', () => {
    it('should show not-found state for invalid server ID', () => {
      cy.visit('/mcp-catalog/999');
      cy.findByTestId('mcp-server-not-found').should('be.visible');
      cy.contains('MCP server not found').should('be.visible');
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
