/* eslint-disable camelcase */
import { modelCatalog } from '~/__tests__/cypress/cypress/pages/modelCatalog';
import { mockModelRegistry } from '~/__mocks__/mockModelRegistry';
import {
  setupModelCatalogIntercepts,
  setupValidatedModelIntercepts,
  interceptArtifactsList,
  interceptPerformanceArtifactsList,
} from '~/__tests__/cypress/cypress/support/interceptHelpers/modelCatalog';
import { mockCatalogModelArtifact } from '~/__mocks__';
import { ModelRegistryMetadataType } from '~/app/types';

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

  it('does not show architecture field when no architectures are available', () => {
    modelCatalog.findLoadingState().should('not.exist');
    modelCatalog.findModelCatalogDetailLink().first().click();
    // Architecture field should not exist when no valid architectures
    modelCatalog.findModelArchitecture().should('not.exist');
  });
});

describe('Model Catalog Details Page - Architecture Field', () => {
  beforeEach(() => {
    // Mock model registries for register button functionality
    cy.intercept('GET', '/model-registry/api/v1/model_registry*', [
      mockModelRegistry({ name: 'modelregistry-sample' }),
    ]).as('getModelRegistries');

    setupModelCatalogIntercepts({});
  });

  it('shows architecture field with valid architectures', () => {
    // Set up intercept with architecture data before navigation
    interceptArtifactsList({
      items: [
        mockCatalogModelArtifact({
          customProperties: {
            architecture: {
              string_value: '["amd64", "arm64", "s390x"]',
              metadataType: ModelRegistryMetadataType.STRING,
            },
          },
        }),
      ],
      size: 1,
      pageSize: 10,
      nextPageToken: '',
    });

    modelCatalog.visit();
    modelCatalog.findLoadingState().should('not.exist');
    modelCatalog.findModelCatalogDetailLink().first().click();
    modelCatalog.findBreadcrumb().should('exist');

    // Architecture field should exist and show correct values
    modelCatalog.findModelArchitecture().should('be.visible');
    modelCatalog.findModelArchitecture().should('contain.text', 'amd64, arm64, s390x');
  });

  it('shows architecture field with uppercase architectures normalized to lowercase', () => {
    // Set up intercept with uppercase architecture data before navigation
    interceptArtifactsList({
      items: [
        mockCatalogModelArtifact({
          customProperties: {
            architecture: {
              string_value: '["AMD64", "ARM64"]',
              metadataType: ModelRegistryMetadataType.STRING,
            },
          },
        }),
      ],
      size: 1,
      pageSize: 10,
      nextPageToken: '',
    });

    modelCatalog.visit();
    modelCatalog.findLoadingState().should('not.exist');
    modelCatalog.findModelCatalogDetailLink().first().click();

    // Architecture should be normalized to lowercase
    modelCatalog.findModelArchitecture().should('be.visible');
    modelCatalog.findModelArchitecture().should('contain.text', 'amd64, arm64');
  });

  it('shows architecture field with custom architecture values', () => {
    // Set up intercept with custom architecture data before navigation
    interceptArtifactsList({
      items: [
        mockCatalogModelArtifact({
          customProperties: {
            architecture: {
              string_value: '["custom-arch", "unknown"]',
              metadataType: ModelRegistryMetadataType.STRING,
            },
          },
        }),
      ],
      size: 1,
      pageSize: 10,
      nextPageToken: '',
    });

    modelCatalog.visit();
    modelCatalog.findLoadingState().should('not.exist');
    modelCatalog.findModelCatalogDetailLink().first().click();

    // Architecture field should display all architecture values without validation
    modelCatalog.findModelArchitecture().should('be.visible');
    modelCatalog.findModelArchitecture().should('contain.text', 'custom-arch, unknown');
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
    modelCatalog.visit();
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
    modelCatalog.visit();
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
