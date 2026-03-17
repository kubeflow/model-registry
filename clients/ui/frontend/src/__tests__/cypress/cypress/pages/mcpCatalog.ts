import { appChrome } from './appChrome';

class McpCatalog {
  visit() {
    cy.visit('/mcp-catalog');
    this.wait();
  }

  private wait() {
    cy.contains('MCP Catalog').should('exist');
    cy.testA11y();
  }

  tabEnabled() {
    appChrome.findNavItem('MCP Catalog').should('exist');
    return this;
  }

  findPageTitle() {
    return cy.contains('MCP Catalog');
  }

  findPageDescription() {
    return cy.contains(
      'Browse and deploy MCP servers provided by Red Hat partners and other providers.',
    );
  }

  findFilter(filterKey: string) {
    return cy.findByTestId(`mcp-filter-${filterKey}`);
  }

  findFilterShowMore(filterKey: string) {
    return cy.findByTestId(`mcp-filter-${filterKey}-show-more`);
  }

  findFilterShowLess(filterKey: string) {
    return cy.findByTestId(`mcp-filter-${filterKey}-show-less`);
  }

  findFilterCheckbox(filterKey: string, value: string) {
    return cy.findByTestId(`mcp-filter-${filterKey}-${value}`);
  }

  findMcpCatalogCards() {
    return cy.get('[data-testid^="mcp-catalog-card-"]');
  }

  findMcpCatalogCard(serverId: string) {
    return cy.findByTestId(`mcp-catalog-card-${serverId}`);
  }

  findCardDetailLink(serverId: string) {
    return cy.findByTestId(`mcp-catalog-card-detail-link-${serverId}`);
  }

  findCardDescription(serverId: string) {
    return cy.findByTestId(`mcp-catalog-card-description-${serverId}`);
  }

  findEmptyState() {
    return cy.findByTestId('mcp-catalog-empty-search');
  }

  findResetFilters() {
    return cy.findByTestId('mcp-catalog-reset-filters');
  }

  findLoadError() {
    return cy.findByTestId('mcp-catalog-load-error');
  }

  findRetry() {
    return cy.findByTestId('mcp-catalog-retry');
  }
}

class McpServerDetails {
  visit(serverId: string) {
    cy.visit(`/mcp-catalog/${serverId}`);
    this.wait();
  }

  private wait() {
    cy.findByTestId('app-page-title').should('exist');
    cy.testA11y();
  }

  findBreadcrumbCatalogLink() {
    return cy.get('.pf-v6-c-breadcrumb').contains('MCP Catalog');
  }

  findBreadcrumbServerName() {
    return cy.findByTestId('breadcrumb-server-name');
  }

  findDeployButton() {
    return cy.findByTestId('deploy-mcp-server-button');
  }

  findDescription() {
    return cy.findByTestId('mcp-server-description');
  }

  findReadmeMarkdown() {
    return cy.findByTestId('mcp-server-readme-markdown');
  }

  findNoReadme() {
    return cy.findByTestId('mcp-server-no-readme');
  }

  findVersion() {
    return cy.findByTestId('mcp-server-version');
  }

  findDeploymentMode() {
    return cy.findByTestId('mcp-server-deployment-mode');
  }

  findTransportType() {
    return cy.findByTestId('mcp-server-transport-type');
  }

  findProvider() {
    return cy.findByTestId('mcp-server-provider');
  }

  findLicense() {
    return cy.findByTestId('mcp-server-license');
  }

  findLicenseLink() {
    return cy.findByTestId('mcp-server-license-link');
  }

  findLabels() {
    return cy.get('[data-testid="mcp-server-detail-label"]');
  }

  findArtifactCopy() {
    return cy.get('[data-testid="mcp-server-artifact-copy"]');
  }

  findSourceCodeLink() {
    return cy.findByTestId('mcp-server-source-code-link');
  }
}

export const mcpCatalog = new McpCatalog();
export const mcpServerDetails = new McpServerDetails();
