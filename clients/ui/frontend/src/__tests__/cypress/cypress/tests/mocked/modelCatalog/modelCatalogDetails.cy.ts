import { modelCatalog } from '~/__tests__/cypress/cypress/pages/modelCatalog';
import { mockModelRegistry } from '~/__mocks__/mockModelRegistry';
import {
  setupModelCatalogIntercepts,
  setupValidatedModelIntercepts,
  interceptArtifactsList,
  interceptPerformanceArtifactsList,
} from '~/__tests__/cypress/cypress/support/interceptHelpers/modelCatalog';

describe('Model Catalog Details Page', () => {
  beforeEach(() => {
    // Mock model registries for register button functionality
    cy.intercept('GET', '/model-registry/api/v1/model_registry*', [
      mockModelRegistry({ name: 'modelregistry-sample' }),
    ]).as('getModelRegistries');

    setupModelCatalogIntercepts({});
    modelCatalog.visit();
  });

  it('navigates to details and shows header, breadcrumb and description', () => {
    modelCatalog.findLoadingState().should('not.exist');
    modelCatalog.findModelCatalogDetailLink().first().click();
    modelCatalog.findBreadcrumb().should('exist');
    modelCatalog.findDetailsProviderText().should('be.visible');
    modelCatalog.findDetailsDescription().should('exist');
  });
});

/**
 * NOTE: Performance Insights Tab Navigation, Hardware Configuration Table,
 * Workload Type Filter, and Latency Filter tests are covered in modelCatalogTabs.cy.ts.
 * This file focuses on filter state management across pages.
 */

describe('Model Catalog Details Page - Filter State Management', () => {
  beforeEach(() => {
    cy.intercept('GET', '/model-registry/api/v1/model_registry*', [
      mockModelRegistry({ name: 'modelregistry-sample' }),
    ]).as('getModelRegistries');

    // Use validated model intercepts which include performance artifacts
    setupValidatedModelIntercepts({});
    interceptArtifactsList();
    interceptPerformanceArtifactsList();
  });

  it('should persist filter state when navigating between Overview and Performance Insights tabs', () => {
    modelCatalog.visit({ enableTempDevCatalogAdvancedFiltersFeature: true });
    modelCatalog.findLoadingState().should('not.exist');
    modelCatalog.togglePerformanceView();
    modelCatalog.findLoadingState().should('not.exist');

    modelCatalog.findModelCatalogDetailLink().first().click();
    modelCatalog.clickPerformanceInsightsTab();

    // Change a filter
    modelCatalog.findWorkloadTypeFilter().click();
    modelCatalog.selectWorkloadType('code_fixing');

    // Navigate to Overview tab
    modelCatalog.clickOverviewTab();
    modelCatalog.findOverviewTabContent().should('be.visible');

    // Navigate back to Performance Insights
    modelCatalog.clickPerformanceInsightsTab();

    // Filter should still show the selected value (capitalized as "Code Fixing")
    modelCatalog.findWorkloadTypeFilter().should('contain.text', 'Code Fixing');
  });

  it('should sync filter changes back to catalog page', () => {
    modelCatalog.visit({ enableTempDevCatalogAdvancedFiltersFeature: true });
    modelCatalog.findLoadingState().should('not.exist');
    modelCatalog.togglePerformanceView();
    modelCatalog.findLoadingState().should('not.exist');

    modelCatalog.findModelCatalogDetailLink().first().click();
    modelCatalog.clickPerformanceInsightsTab();

    // Change a filter on details page
    modelCatalog.findWorkloadTypeFilter().click();
    modelCatalog.selectWorkloadType('rag');

    // Go back to catalog page
    cy.go('back');
    cy.go('back');
    modelCatalog.findLoadingState().should('not.exist');

    // The alert should show indicating filters were updated
    modelCatalog.findPerformanceFiltersUpdatedAlert().should('be.visible');
  });
});
