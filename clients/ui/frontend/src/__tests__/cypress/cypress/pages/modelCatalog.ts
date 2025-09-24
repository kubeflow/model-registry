import { appChrome } from './appChrome';

class ModelCatalog {
  visit() {
    cy.visit('/model-catalog');
    this.wait();
  }

  navigate() {
    appChrome.findNavItem('Model Catalog').click();
    this.wait();
  }

  private wait() {
    cy.findByTestId('app-page-title').should('exist');
    cy.findByTestId('app-page-title').contains('Model Catalog');
    cy.testA11y();
  }

  tabEnabled() {
    appChrome.findNavItem('Model Catalog').should('exist');
    return this;
  }

  findModelCatalogEmptyState() {
    return cy.findByTestId('empty-model-catalog-state');
  }

  findModelCatalogCards() {
    return cy.findAllByTestId('model-catalog-card');
  }

  findFirstModelCatalogCard() {
    return this.findModelCatalogCards().first().should('be.visible');
  }

  findModelCatalogDetailLink() {
    return cy.findAllByTestId('model-catalog-detail-link');
  }

  findModelCatalogDescription() {
    return cy.findByTestId('model-catalog-card-description');
  }

  findSourceLabel() {
    return cy.get('.pf-v6-c-label');
  }

  findModelLogo() {
    return cy.get('img[alt="model logo"]');
  }

  findVersionIcon() {
    return cy.get('.pf-v6-c-icon');
  }

  findFrameworkLabel() {
    return cy.contains('PyTorch');
  }

  findTaskLabel() {
    return cy.contains('text-generation');
  }

  findLicenseLabel() {
    return cy.contains('apache-2.0');
  }

  findProviderLabel() {
    return cy.contains('provider1');
  }

  findLoadingState() {
    return cy.contains('Loading model catalog...');
  }

  findPageTitle() {
    return cy.contains('Model Catalog');
  }

  findPageDescription() {
    return cy.contains('Discover models that are available for your organization');
  }

  // Details page helpers
  findBreadcrumb() {
    return cy.contains('Model catalog');
  }

  findDetailsProviderText() {
    return cy.contains('Provided by');
  }

  findDetailsDescription() {
    return cy.findByTestId('model-long-description');
  }

  // Tabs functionality
  findModelDetailsTabs() {
    return cy.findByTestId('model-details-page-tabs');
  }

  findOverviewTab() {
    return cy.findByTestId('model-overview-tab');
  }

  findPerformanceInsightsTab() {
    return cy.findByTestId('performance-insights-tab');
  }

  findOverviewTabContent() {
    return cy.findByTestId('model-overview-tab-content');
  }

  findPerformanceInsightsTabContent() {
    return cy.findByTestId('performance-insights-tab-content');
  }

  clickOverviewTab() {
    this.findOverviewTab().click();
    return this;
  }

  clickPerformanceInsightsTab() {
    this.findPerformanceInsightsTab().click();
    return this;
  }
}

export const modelCatalog = new ModelCatalog();
