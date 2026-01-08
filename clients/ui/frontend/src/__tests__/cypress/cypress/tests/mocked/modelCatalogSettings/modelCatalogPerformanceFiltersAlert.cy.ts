import { mockModArchResponse } from 'mod-arch-core';
import {
  mockCatalogModel,
  mockCatalogModelList,
  mockCatalogSource,
  mockCatalogSourceList,
  mockValidatedModel,
} from '~/__mocks__';
import { mockCatalogPerformanceMetricsArtifactList } from '~/__mocks__/mockCatalogModelArtifactList';
import { modelCatalog } from '~/__tests__/cypress/cypress/pages/modelCatalog';
import { mockModelRegistry } from '~/__mocks__/mockModelRegistry';
import type { CatalogSource } from '~/app/modelCatalogTypes';
import { MODEL_CATALOG_API_VERSION } from '~/__tests__/cypress/cypress/support/commands/api';
import { mockCatalogFilterOptionsList } from '~/__mocks__/mockCatalogFilterOptionsList';
import { ModelRegistryMetadataType } from '~/app/types';
import type { ModelRegistryCustomProperties } from '~/app/types';

type HandlersProps = {
  sources?: CatalogSource[];
  modelsPerCategory?: number;
};

const initIntercepts = ({
  sources = [mockCatalogSource({}), mockCatalogSource({ id: 'source-2', name: 'source 2' })],
  modelsPerCategory = 4,
}: HandlersProps) => {
  const testModel = mockValidatedModel;
  const testArtifacts = mockCatalogPerformanceMetricsArtifactList({});

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
      // Use regex intercept to handle all requests including those with filterQuery
      const encodedLabel = encodeURIComponent(label);
      const mockModels = mockCatalogModelList({
        items: Array.from({ length: modelsPerCategory }, (_, i) => {
          const customProperties =
            i === 0
              ? ({
                  validated: {
                    metadataType: ModelRegistryMetadataType.STRING,
                    // eslint-disable-next-line camelcase
                    string_value: '',
                  },
                } as ModelRegistryCustomProperties)
              : undefined;
          const name = i === 0 ? 'validated-model' : `${label.toLowerCase()}-model-${i + 1}`;
          return mockCatalogModel({
            name,
            // eslint-disable-next-line camelcase
            source_id: source.id,
            customProperties,
          });
        }),
      });
      cy.intercept(
        {
          method: 'GET',
          url: new RegExp(
            `/model-registry/api/${MODEL_CATALOG_API_VERSION}/model_catalog/models.*sourceLabel=${encodedLabel}`,
          ),
        },
        mockModArchResponse(mockModels),
      ).as(`getModels-${label}`);
    });
  });

  // When "All models" is selected and filters are applied (GalleryView), the request
  // may not include sourceLabel. Create mock models that include validated models.
  const allModelsResponse = mockCatalogModelList({
    items: Array.from({ length: modelsPerCategory }, (_, i) => {
      const customProperties =
        i === 0
          ? ({
              validated: {
                metadataType: ModelRegistryMetadataType.STRING,
                // eslint-disable-next-line camelcase
                string_value: '',
              },
            } as ModelRegistryCustomProperties)
          : undefined;
      const name = i === 0 ? 'validated-model' : `all-models-model-${i + 1}`;
      return mockCatalogModel({
        name,
        // eslint-disable-next-line camelcase
        source_id: 'sample-source',
        customProperties,
      });
    }),
  });

  // Intercept for GalleryView when filters are applied (no sourceLabel, but has filterQuery or pageSize)
  cy.intercept(
    {
      method: 'GET',
      url: new RegExp(
        `/model-registry/api/${MODEL_CATALOG_API_VERSION}/model_catalog/models\\?(?!.*sourceLabel=)`,
      ),
    },
    mockModArchResponse(allModelsResponse),
  ).as('getModelsFiltered');

  // Use regex to match any source's model details requests for validated-model
  cy.intercept(
    {
      method: 'GET',
      url: new RegExp(
        `/model-registry/api/${MODEL_CATALOG_API_VERSION}/model_catalog/sources/.*/models/validated-model`,
      ),
    },
    mockModArchResponse(testModel),
  ).as('getCatalogModel');

  // Use regex to match any source's artifacts requests for validated-model
  cy.intercept(
    {
      method: 'GET',
      url: new RegExp(
        `/model-registry/api/${MODEL_CATALOG_API_VERSION}/model_catalog/sources/.*/artifacts/validated-model`,
      ),
    },
    mockModArchResponse(testArtifacts),
  ).as('getCatalogModelArtifacts');

  // Use regex to match any source's validated-model performance artifacts requests
  cy.intercept(
    {
      method: 'GET',
      url: new RegExp(
        `/model-registry/api/${MODEL_CATALOG_API_VERSION}/model_catalog/sources/.*/performance_artifacts/validated-model`,
      ),
    },
    mockModArchResponse(testArtifacts),
  ).as('getCatalogSourceModelArtifacts');

  cy.interceptApi(
    `GET /api/:apiVersion/model_catalog/models/filter_options`,
    {
      path: { apiVersion: MODEL_CATALOG_API_VERSION },
      query: { namespace: 'kubeflow' },
    },
    mockCatalogFilterOptionsList(),
  );
};

describe('Model Catalog Performance Filters Alert', () => {
  beforeEach(() => {
    cy.intercept('GET', '/model-registry/api/v1/model_registry*', [
      mockModelRegistry({ name: 'modelregistry-sample' }),
    ]).as('getModelRegistries');

    initIntercepts({});
    modelCatalog.visit({ enableTempDevCatalogAdvancedFiltersFeature: true });
  });

  describe('Alert Display Logic', () => {
    it('should not show alert initially', () => {
      modelCatalog.findLoadingState().should('not.exist');
      modelCatalog.findPerformanceFiltersUpdatedAlert().should('not.exist');
    });

    it('should not show alert when performance view is disabled', () => {
      modelCatalog.findModelCatalogDetailLink().first().click();
      modelCatalog.clickPerformanceInsightsTab();

      modelCatalog.findWorkloadTypeFilter().click();
      modelCatalog.selectWorkloadType('code_fixing');

      cy.go('back');
      cy.go('back');
      modelCatalog.findLoadingState().should('not.exist');

      modelCatalog.findPerformanceFiltersUpdatedAlert().should('not.exist');
    });

    it('should show alert when returning from details page after changing performance filters', () => {
      modelCatalog.togglePerformanceView();
      modelCatalog.findPerformanceViewToggleValue().should('be.checked');
      // Wait for the models to reload after toggle applies default filters
      modelCatalog.findLoadingState().should('not.exist');

      modelCatalog.findModelCatalogDetailLink().first().click();
      modelCatalog.clickPerformanceInsightsTab();

      modelCatalog.findWorkloadTypeFilter().click();
      modelCatalog.selectWorkloadType('code_fixing');

      cy.go('back');
      cy.go('back');
      modelCatalog.findLoadingState().should('not.exist');

      modelCatalog.findPerformanceFiltersUpdatedAlert().should('be.visible');
      modelCatalog
        .findPerformanceFiltersUpdatedAlert()
        .should(
          'contain.text',
          'The results list has been updated to match the latest performance criteria set on the details page.',
        );
    });

    it('should not show alert when no filters were changed on details page', () => {
      modelCatalog.togglePerformanceView();
      modelCatalog.findLoadingState().should('not.exist');

      modelCatalog.findModelCatalogDetailLink().first().click();
      modelCatalog.clickPerformanceInsightsTab();

      cy.go('back');
      cy.go('back');
      modelCatalog.findLoadingState().should('not.exist');

      modelCatalog.findPerformanceFiltersUpdatedAlert().should('not.exist');
    });
  });

  describe('Alert Dismissal', () => {
    it('should dismiss alert when close button is clicked', () => {
      modelCatalog.togglePerformanceView();
      modelCatalog.findLoadingState().should('not.exist');

      modelCatalog.findModelCatalogDetailLink().first().click();
      modelCatalog.clickPerformanceInsightsTab();

      modelCatalog.findWorkloadTypeFilter().click();
      modelCatalog.selectWorkloadType('code_fixing');

      cy.go('back');
      cy.go('back');
      modelCatalog.findLoadingState().should('not.exist');

      modelCatalog.findPerformanceFiltersUpdatedAlert().should('be.visible');

      modelCatalog.dismissPerformanceFiltersUpdatedAlert();

      modelCatalog.findPerformanceFiltersUpdatedAlert().should('not.exist');
    });
  });

  describe('Alert Hidden Scenarios', () => {
    it('should hide alert when performance toggle is turned OFF', () => {
      modelCatalog.togglePerformanceView();
      modelCatalog.findLoadingState().should('not.exist');

      modelCatalog.findModelCatalogDetailLink().first().click();
      modelCatalog.clickPerformanceInsightsTab();

      modelCatalog.findWorkloadTypeFilter().click();
      modelCatalog.selectWorkloadType('code_fixing');

      cy.go('back');
      cy.go('back');
      modelCatalog.findLoadingState().should('not.exist');

      modelCatalog.findPerformanceFiltersUpdatedAlert().should('be.visible');

      modelCatalog.togglePerformanceView();

      modelCatalog.findPerformanceFiltersUpdatedAlert().should('not.exist');
    });

    it('should hide alert when filters change on catalog page', () => {
      modelCatalog.togglePerformanceView();
      modelCatalog.findLoadingState().should('not.exist');

      modelCatalog.findModelCatalogDetailLink().first().click();
      modelCatalog.clickPerformanceInsightsTab();

      modelCatalog.findWorkloadTypeFilter().click();
      modelCatalog.selectWorkloadType('code_fixing');

      cy.go('back');
      cy.go('back');
      modelCatalog.findLoadingState().should('not.exist');

      modelCatalog.findPerformanceFiltersUpdatedAlert().should('be.visible');

      modelCatalog.findFilter('Task').should('be.visible');
      modelCatalog.findFilterCheckbox('Task', 'text-generation').click();

      modelCatalog.findPerformanceFiltersUpdatedAlert().should('not.exist');
    });
  });

  describe('Multiple Filter Changes', () => {
    it('should show alert after changing multiple performance filters', () => {
      modelCatalog.togglePerformanceView();
      modelCatalog.findLoadingState().should('not.exist');

      modelCatalog.findModelCatalogDetailLink().first().click();
      modelCatalog.clickPerformanceInsightsTab();

      modelCatalog.findWorkloadTypeFilter().click();
      modelCatalog.selectWorkloadType('code_fixing');
      // Re-open dropdown to select second option
      modelCatalog.findWorkloadTypeFilter().click();
      modelCatalog.selectWorkloadType('chatbot');

      cy.go('back');
      cy.go('back');
      modelCatalog.findLoadingState().should('not.exist');

      modelCatalog.findPerformanceFiltersUpdatedAlert().should('be.visible');
    });
  });
});
