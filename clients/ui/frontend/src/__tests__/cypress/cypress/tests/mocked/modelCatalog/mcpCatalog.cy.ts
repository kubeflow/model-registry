import { mcpCatalog } from '~/__tests__/cypress/cypress/pages/mcpCatalog';

describe('MCP Catalog Page', () => {
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
    mcpCatalog.findMcpCatalogCards().should('have.length.at.least', 1);
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

  it('should display known mock server cards by id', () => {
    mcpCatalog.visit();
    mcpCatalog.findMcpCatalogCard('kubernetes').should('be.visible');
    mcpCatalog.findMcpCatalogCard('github').should('be.visible');
  });
});
