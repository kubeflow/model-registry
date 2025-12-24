import {
  mockCatalogModel,
  mockCatalogModelArtifactList,
  mockCatalogModelList,
  mockCatalogSource,
  mockCatalogSourceList,
  mockNonValidatedModel,
  mockValidatedModel,
} from '~/__mocks__';
import { mockCatalogPerformanceMetricsArtifactList } from '~/__mocks__/mockCatalogModelArtifactList';
import { modelCatalog } from '~/__tests__/cypress/cypress/pages/modelCatalog';
import { mockModelRegistry } from '~/__mocks__/mockModelRegistry';
import type { CatalogSource } from '~/app/modelCatalogTypes';
import { MODEL_CATALOG_API_VERSION } from '~/__tests__/cypress/cypress/support/commands/api';
import type { ModelRegistryCustomProperties } from '~/app/types';
import { ModelRegistryMetadataType } from '~/app/types';
import { mockCatalogFilterOptionsList } from '~/__mocks__/mockCatalogFilterOptionsList';

type HandlersProps = {
  sources?: CatalogSource[];
  useValidatedModel?: boolean;
  modelsPerCategory?: number;
  hasPerformanceArtifacts?: boolean;
};

const initIntercepts = ({
  sources = [mockCatalogSource({}), mockCatalogSource({ id: 'source-2', name: 'source 2' })],
  useValidatedModel = true,
  modelsPerCategory = 4,
  hasPerformanceArtifacts = true,
}: HandlersProps) => {
  const testModel = useValidatedModel ? mockValidatedModel : mockNonValidatedModel;

  const testArtifacts = hasPerformanceArtifacts
    ? mockCatalogPerformanceMetricsArtifactList({})
    : mockCatalogModelArtifactList({});

  cy.interceptApi(
    `GET /api/:apiVersion/model_catalog/sources`,
    {
      path: { apiVersion: MODEL_CATALOG_API_VERSION },
    },
    mockCatalogSourceList({
      items: sources,
    }),
  );

  sources.forEach((source) => {
    source.labels.forEach((label) => {
      cy.interceptApi(
        `GET /api/:apiVersion/model_catalog/models`,
        {
          path: { apiVersion: MODEL_CATALOG_API_VERSION },
          query: { sourceLabel: label },
        },
        mockCatalogModelList({
          items: Array.from({ length: modelsPerCategory }, (_, i) => {
            const customProperties =
              i === 0 && useValidatedModel
                ? ({
                    validated: {
                      metadataType: ModelRegistryMetadataType.STRING,
                      // eslint-disable-next-line camelcase
                      string_value: '',
                    },
                  } as ModelRegistryCustomProperties)
                : undefined;
            const name =
              i === 0 && useValidatedModel
                ? 'validated-model'
                : `${label.toLowerCase()}-model-${i + 1}`;

            return mockCatalogModel({
              name,
              // eslint-disable-next-line camelcase
              source_id: source.id,
              customProperties,
            });
          }),
        }),
      );
    });
  });

  cy.interceptApi(
    `GET /api/:apiVersion/model_catalog/sources/:sourceId/models/:modelName`,
    {
      path: {
        apiVersion: MODEL_CATALOG_API_VERSION,
        sourceId: 'source-2',
        modelName: testModel.name.replace('/', '%2F'),
      },
    },
    testModel,
  );

  // Intercept for /artifacts/ - used to determine if tabs should show
  cy.interceptApi(
    `GET /api/:apiVersion/model_catalog/sources/:sourceId/artifacts/:modelName`,
    {
      path: {
        apiVersion: MODEL_CATALOG_API_VERSION,
        sourceId: 'source-2',
        modelName: testModel.name.replace('/', '%2F'),
      },
    },
    testArtifacts,
  );

  // Intercept for /performance_artifacts/ - used for server-side filtered performance data
  cy.interceptApi(
    `GET /api/:apiVersion/model_catalog/sources/:sourceId/performance_artifacts/:modelName`,
    {
      path: {
        apiVersion: MODEL_CATALOG_API_VERSION,
        sourceId: 'source-2',
        modelName: testModel.name.replace('/', '%2F'),
      },
    },
    testArtifacts,
  );

  cy.interceptApi(
    `GET /api/:apiVersion/model_catalog/models/filter_options`,
    {
      path: { apiVersion: MODEL_CATALOG_API_VERSION },
      query: { namespace: 'kubeflow' },
    },
    mockCatalogFilterOptionsList(),
  );
};

describe('Model Catalog Details Tabs', () => {
  describe('Validated Models with performance artifacts (with tabs)', () => {
    beforeEach(() => {
      // Mock model registries for register button functionality
      cy.intercept('GET', '/model-registry/api/v1/model_registry*', [
        mockModelRegistry({ name: 'modelregistry-sample' }),
      ]).as('getModelRegistries');

      initIntercepts({ useValidatedModel: true, hasPerformanceArtifacts: true });
      modelCatalog.visit();
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
        cy.url().should('include', '/model-catalog/source-2/validated-model/overview');
      });

      it('should switch to Performance Insights tab when clicked', () => {
        modelCatalog.findModelCatalogDetailLink().first().click();

        cy.url().should('include', '/model-catalog/source-2/validated-model/overview');

        // Click Performance Insights tab
        modelCatalog.clickPerformanceInsightsTab();

        // Verify tab switch
        modelCatalog.findPerformanceInsightsTab().should('have.attr', 'aria-selected', 'true');
        modelCatalog.findOverviewTab().should('have.attr', 'aria-selected', 'false');
        modelCatalog.findPerformanceInsightsTabContent().should('be.visible');
        cy.url().should('include', '/model-catalog/source-2/validated-model/performance-insights');
      });

      it('should switch back to Overview tab when clicked', () => {
        modelCatalog.findModelCatalogDetailLink().first().click();

        cy.url().should('include', '/model-catalog/source-2/validated-model/overview');

        // First switch to Performance Insights
        modelCatalog.clickPerformanceInsightsTab();
        modelCatalog.findPerformanceInsightsTab().should('have.attr', 'aria-selected', 'true');

        // Then switch back to Overview
        modelCatalog.clickOverviewTab();
        modelCatalog.findOverviewTab().should('have.attr', 'aria-selected', 'true');
        modelCatalog.findPerformanceInsightsTab().should('have.attr', 'aria-selected', 'false');
        modelCatalog.findOverviewTabContent().should('be.visible');
        cy.url().should('include', '/model-catalog/source-2/validated-model/overview');
      });
    });

    describe('Tab Content', () => {
      it('should display Hardware Configuration content in Performance Insights tab', () => {
        modelCatalog.findModelCatalogDetailLink().first().click();

        // Switch to Performance Insights tab
        modelCatalog.clickPerformanceInsightsTab();

        // Verify Hardware Configuration content is displayed
        modelCatalog.findPerformanceInsightsTabContent().should('be.visible');
        modelCatalog.findHardwareConfigurationTitle().should('be.visible');
        modelCatalog.findHardwareConfigurationDescription().should('be.visible');
        modelCatalog.findHardwareConfigurationTable().should('be.visible');
      });

      it('should display Workload type column as the second column in hardware configuration table', () => {
        modelCatalog.findModelCatalogDetailLink().first().click();
        modelCatalog.clickPerformanceInsightsTab();

        modelCatalog
          .findHardwareConfigurationTableHeaders()
          .eq(1)
          .should('contain.text', 'Workload type');
        modelCatalog
          .findHardwareConfigurationColumn('Workload type')
          .first()
          .should('contain.text', 'Code Fixing')
          .should('not.contain.text', 'code_fixing');
      });
    });

    describe('Workload Type Filter', () => {
      it('should display workload type filter in the toolbar', () => {
        modelCatalog.findModelCatalogDetailLink().first().click();
        modelCatalog.clickPerformanceInsightsTab();

        modelCatalog
          .findWorkloadTypeFilter()
          .should('be.visible')
          .should('contain.text', 'Workload type');
      });

      it('should show workload type options when clicked', () => {
        modelCatalog.findModelCatalogDetailLink().first().click();
        modelCatalog.clickPerformanceInsightsTab();
        modelCatalog.findWorkloadTypeFilter().click();
        modelCatalog.findWorkloadTypeOption('chatbot').should('be.visible');
        modelCatalog.findWorkloadTypeOption('code_fixing').should('be.visible');
        modelCatalog.findWorkloadTypeOption('long_rag').should('be.visible');
        modelCatalog.findWorkloadTypeOption('rag').should('be.visible');
      });

      it('should update toggle text when workload type is selected', () => {
        modelCatalog.findModelCatalogDetailLink().first().click();
        modelCatalog.clickPerformanceInsightsTab();
        modelCatalog.findWorkloadTypeFilter().click();
        modelCatalog.selectWorkloadType('code_fixing');
        modelCatalog
          .findWorkloadTypeFilter()
          .should('contain.text', 'Workload type')
          .should('contain.text', '1 selected');
      });

      it('should filter hardware configuration table by selected workload type', () => {
        modelCatalog.findModelCatalogDetailLink().first().click();
        modelCatalog.clickPerformanceInsightsTab();
        modelCatalog.findHardwareConfigurationTableRows().should('have.length.at.least', 1);
        modelCatalog.findWorkloadTypeFilter().click();
        modelCatalog.selectWorkloadType('code_fixing');
        modelCatalog.findHardwareConfigurationTableRows().should('exist');
        modelCatalog.findHardwareConfigurationColumn('Workload type').each(($el) => {
          cy.wrap($el).should('contain.text', 'Code Fixing');
        });
      });

      it('should clear workload type filter when clicking selected option again', () => {
        modelCatalog.findModelCatalogDetailLink().first().click();
        modelCatalog.clickPerformanceInsightsTab();
        modelCatalog.findWorkloadTypeFilter().click();
        modelCatalog.selectWorkloadType('code_fixing');
        modelCatalog
          .findWorkloadTypeFilter()
          .should('contain.text', 'Workload type')
          .should('contain.text', '1 selected');

        // Re-open dropdown before deselecting
        modelCatalog.findWorkloadTypeFilter().click();
        modelCatalog.selectWorkloadType('code_fixing');
        modelCatalog.findWorkloadTypeFilter().should('contain.text', 'Workload type');
        modelCatalog.findWorkloadTypeFilter().should('not.contain.text', '1 selected');
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

        cy.url().should('include', '/model-catalog/source-2/validated-model/overview');

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

  describe('Validated Models without performance artifacts (without tabs)', () => {
    beforeEach(() => {
      cy.intercept('GET', '/model-registry/api/v1/model_registry*', [
        mockModelRegistry({ name: 'modelregistry-sample' }),
      ]).as('getModelRegistries');

      initIntercepts({ useValidatedModel: true, hasPerformanceArtifacts: false });
      modelCatalog.visit();
    });

    it('should not display tabs for validated models without performance artifacts', () => {
      modelCatalog.findLoadingState().should('not.exist');
      modelCatalog.findModelCatalogDetailLink().first().click();
      modelCatalog.findModelDetailsTabs().should('not.exist');
      modelCatalog.findOverviewTab().should('not.exist');
      modelCatalog.findPerformanceInsightsTab().should('not.exist');

      modelCatalog.findOverviewTabContent().should('be.visible');
      modelCatalog.findDetailsDescription().should('be.visible');
    });
  });

  describe('Non-Validated Models (without tabs)', () => {
    beforeEach(() => {
      // Mock model registries for register button functionality
      cy.intercept('GET', '/model-registry/api/v1/model_registry*', [
        mockModelRegistry({ name: 'modelregistry-sample' }),
      ]).as('getModelRegistries');

      initIntercepts({ useValidatedModel: false, hasPerformanceArtifacts: false });
      modelCatalog.visit();
    });

    it('should not display tabs for non-validated models', () => {
      modelCatalog.findLoadingState().should('not.exist');
      modelCatalog.findModelCatalogDetailLink().first().click();

      // Tabs should not be present
      modelCatalog.findModelDetailsTabs().should('not.exist');
      modelCatalog.findOverviewTab().should('not.exist');
      modelCatalog.findPerformanceInsightsTab().should('not.exist');

      // But overview content should still be visible
      modelCatalog.findOverviewTabContent().should('be.visible');
      modelCatalog.findDetailsDescription().should('be.visible');
    });
  });
});
