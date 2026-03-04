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
      'Discover and manage MCP servers and tools available for your organization.',
    );
  }

  findFilter(filterKey: string) {
    return cy.findByTestId(`mcp-filter-${filterKey}`);
  }

  findFilterShowMore(filterKey: string) {
    return cy.findByTestId(`mcp-filter-${filterKey}-show-more`);
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
}

export const mcpCatalog = new McpCatalog();
