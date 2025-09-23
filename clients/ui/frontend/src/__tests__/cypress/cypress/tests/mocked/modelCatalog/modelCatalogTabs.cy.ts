import {
  mockCatalogModel,
  mockCatalogModelArtifactList,
  mockCatalogModelList,
  mockCatalogSource,
  mockCatalogSourceList,
} from '~/__mocks__';
import { modelCatalog } from '~/__tests__/cypress/cypress/pages/modelCatalog';
import { mockModelRegistry } from '~/__mocks__/mockModelRegistry';
import type { CatalogSource } from '~/app/modelCatalogTypes';
import { MODEL_CATALOG_API_VERSION } from '~/__tests__/cypress/cypress/support/commands/api';

type HandlersProps = {
  sources?: CatalogSource[];
};

const initIntercepts = ({
  sources = [mockCatalogSource({}), mockCatalogSource({ id: 'source-2', name: 'source 2' })],
}: HandlersProps) => {
  cy.interceptApi(
    `GET /api/:apiVersion/model_catalog/sources`,
    {
      path: { apiVersion: MODEL_CATALOG_API_VERSION },
    },
    mockCatalogSourceList({
      items: sources,
    }),
  );

  cy.interceptApi(
    `GET /api/:apiVersion/model_catalog/models`,
    {
      path: { apiVersion: MODEL_CATALOG_API_VERSION },
      query: { source: 'sample-source' },
    },
    mockCatalogModelList({
      items: [mockCatalogModel({})],
    }),
  );

  cy.interceptApi(
    `GET /api/:apiVersion/model_catalog/sources/:sourceId/models/:modelName`,
    {
      path: {
        apiVersion: MODEL_CATALOG_API_VERSION,
        sourceId: 'sample-source',
        modelName: 'repo1%2Fmodel1',
      },
    },
    mockCatalogModel({}),
  );

  cy.interceptApi(
    `GET /api/:apiVersion/model_catalog/sources/:sourceId/artifacts/:modelName`,
    {
      path: {
        apiVersion: MODEL_CATALOG_API_VERSION,
        sourceId: 'sample-source',
        modelName: 'repo1%2Fmodel1',
      },
    },
    mockCatalogModelArtifactList({}),
  );
};

describe('Model Catalog Details Tabs', () => {
  beforeEach(() => {
    // Mock model registries for register button functionality
    cy.intercept('GET', '/model-registry/api/v1/model_registry*', [
      mockModelRegistry({ name: 'modelregistry-sample' }),
    ]).as('getModelRegistries');

    initIntercepts({});
    modelCatalog.visit();
    modelCatalog.navigate();
  });

  describe('Tab Navigation', () => {
    it('should display tabs on model details page', () => {
      modelCatalog.findLoadingState().should('not.exist');
      modelCatalog.findModelCatalogDetailLink().first().click();

      // Verify tabs are present
      modelCatalog.findModelDetailsTabs().should('be.visible');
      modelCatalog.findOverviewTab().should('be.visible');
      modelCatalog.findPerformanceInsightsTab().should('be.visible');
    });

    it('should show Overview tab as active by default', () => {
      modelCatalog.findModelCatalogDetailLink().first().click();

      // Overview tab should be active and content should be visible
      modelCatalog.findOverviewTab().should('have.attr', 'aria-selected', 'true');
      modelCatalog.findOverviewTabContent().should('be.visible');
      modelCatalog.findDetailsDescription().should('be.visible');
    });

    it('should switch to Performance Insights tab when clicked', () => {
      modelCatalog.findModelCatalogDetailLink().first().click();

      // Click Performance Insights tab
      modelCatalog.clickPerformanceInsightsTab();

      // Verify tab switch
      modelCatalog.findPerformanceInsightsTab().should('have.attr', 'aria-selected', 'true');
      modelCatalog.findOverviewTab().should('have.attr', 'aria-selected', 'false');
      modelCatalog.findPerformanceInsightsTabContent().should('be.visible');
    });

    it('should switch back to Overview tab when clicked', () => {
      modelCatalog.findModelCatalogDetailLink().first().click();

      // First switch to Performance Insights
      modelCatalog.clickPerformanceInsightsTab();
      modelCatalog.findPerformanceInsightsTab().should('have.attr', 'aria-selected', 'true');

      // Then switch back to Overview
      modelCatalog.clickOverviewTab();
      modelCatalog.findOverviewTab().should('have.attr', 'aria-selected', 'true');
      modelCatalog.findPerformanceInsightsTab().should('have.attr', 'aria-selected', 'false');
      modelCatalog.findOverviewTabContent().should('be.visible');
    });
  });

  describe('Tab Content', () => {
    it('should display placeholder content in Performance Insights tab', () => {
      modelCatalog.findModelCatalogDetailLink().first().click();

      // Switch to Performance Insights tab
      modelCatalog.clickPerformanceInsightsTab();

      // Verify placeholder content
      modelCatalog.findPerformanceInsightsTabContent().should('be.visible');
      cy.contains('Performance Insights - Coming Soon').should('be.visible');
    });
  });

  describe('Accessibility', () => {
    it('should have proper ARIA attributes for tabs', () => {
      modelCatalog.findModelCatalogDetailLink().first().click();

      // Check tab container has proper role
      modelCatalog.findModelDetailsTabs().should('have.attr', 'role', 'region');
      modelCatalog
        .findModelDetailsTabs()
        .should('have.attr', 'aria-label', 'Model details page tabs');

      // Check individual tabs have proper attributes
      modelCatalog.findOverviewTab().should('have.attr', 'aria-label', 'Model overview tab');
      modelCatalog
        .findPerformanceInsightsTab()
        .should('have.attr', 'aria-label', 'Performance insights tab');
    });
  });

  describe('Tab State Management', () => {
    it('should maintain tab state when switching between tabs', () => {
      modelCatalog.findModelCatalogDetailLink().first().click();

      // Switch to Performance Insights
      modelCatalog.clickPerformanceInsightsTab();
      modelCatalog.findPerformanceInsightsTab().should('have.attr', 'aria-selected', 'true');

      // Switch back to Overview
      modelCatalog.clickOverviewTab();
      modelCatalog.findOverviewTab().should('have.attr', 'aria-selected', 'true');

      // Switch to Performance Insights again
      modelCatalog.clickPerformanceInsightsTab();
      modelCatalog.findPerformanceInsightsTab().should('have.attr', 'aria-selected', 'true');
    });
  });
});
