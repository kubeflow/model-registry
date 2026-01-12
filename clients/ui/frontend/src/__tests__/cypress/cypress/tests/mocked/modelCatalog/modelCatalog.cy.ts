/* eslint-disable camelcase */
import { modelCatalog } from '~/__tests__/cypress/cypress/pages/modelCatalog';
import {
  mockCatalogAccuracyMetricsArtifact,
  mockCatalogModel,
  mockCatalogModelArtifact,
  mockCatalogModelList,
  mockCatalogPerformanceMetricsArtifact,
  mockCatalogSource,
  mockCatalogSourceList,
} from '~/__mocks__';
import type { CatalogSource } from '~/app/modelCatalogTypes';
import { MODEL_CATALOG_API_VERSION } from '~/__tests__/cypress/cypress/support/commands/api';
import { mockCatalogFilterOptionsList } from '~/__mocks__/mockCatalogFilterOptionsList';
import { SourceLabel } from '~/app/modelCatalogTypes';
import { ModelRegistryMetadataType } from '~/app/types';

type HandlersProps = {
  sources?: CatalogSource[];
  modelsPerCategory?: number;
  hasValidatedModels?: boolean;
};
const initIntercepts = ({
  sources = [mockCatalogSource({}), mockCatalogSource({ id: 'source-2', name: 'source 2' })],
  modelsPerCategory = 4,
  hasValidatedModels = false,
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

  sources.forEach((source) => {
    source.labels.forEach((label) => {
      const models = hasValidatedModels
        ? [
            mockCatalogModel({
              name: 'validated-model',
              source_id: source.id,
              customProperties: {
                validated: { metadataType: ModelRegistryMetadataType.STRING, string_value: '' },
              },
            }),
            ...Array.from({ length: modelsPerCategory - 1 }, (_, i) =>
              mockCatalogModel({
                name: `${label.toLowerCase().replace(/\s+/g, '-')}-model-${i + 1}`,
                source_id: source.id,
              }),
            ),
          ]
        : Array.from({ length: modelsPerCategory }, (_, i) =>
            mockCatalogModel({
              name: `${label.toLowerCase().replace(/\s+/g, '-')}-model-${i + 1}`,
              source_id: source.id,
            }),
          );

      cy.interceptApi(
        `GET /api/:apiVersion/model_catalog/models`,
        {
          path: { apiVersion: MODEL_CATALOG_API_VERSION },
          query: {
            sourceLabel: label,
          },
        },
        mockCatalogModelList({ items: models }),
      );
    });
  });

  // Intercept requests for sources without labels if they exist
  const hasSourcesWithoutLabels = sources.some(
    (source) =>
      source.enabled !== false &&
      (source.labels.length === 0 || source.labels.every((label) => !label.trim())),
  );

  if (hasSourcesWithoutLabels) {
    cy.interceptApi(
      `GET /api/:apiVersion/model_catalog/models`,
      {
        path: { apiVersion: MODEL_CATALOG_API_VERSION },
        query: { sourceLabel: SourceLabel.other },
      },
      mockCatalogModelList({
        items: Array.from({ length: modelsPerCategory }, (_, i) =>
          mockCatalogModel({
            name: `custom-model-${i + 1}`,
            // eslint-disable-next-line camelcase
            source_id: sources.find((s) => s.labels.length === 0)?.id || 'custom-source',
          }),
        ),
      }),
    );
  }

  cy.interceptApi(
    `GET /api/:apiVersion/model_catalog/models/filter_options`,
    {
      path: { apiVersion: MODEL_CATALOG_API_VERSION },
      query: { namespace: 'kubeflow' },
    },
    mockCatalogFilterOptionsList(),
  );

  cy.intercept(
    {
      method: 'GET',
      url: new RegExp(
        `/api/${MODEL_CATALOG_API_VERSION}/model_catalog/sources/.*/performance_artifacts/.*`,
      ),
    },
    { items: [], size: 0, pageSize: 10, nextPageToken: '' },
  ).as('getPerformanceArtifacts');

  cy.interceptApi(
    `GET /api/:apiVersion/model_catalog/sources/:sourceId/artifacts/:modelName`,
    {
      path: {
        apiVersion: MODEL_CATALOG_API_VERSION,
        sourceId: 'sample-source',
        modelName: 'repo1/model1',
      },
    },
    {
      items: [
        mockCatalogPerformanceMetricsArtifact({}),
        mockCatalogAccuracyMetricsArtifact({}),
        mockCatalogModelArtifact({}),
      ],
    },
  );
};

describe('Model Catalog Page', () => {
  it('model catalog tab should be enabled', () => {
    initIntercepts({});
    modelCatalog.visit();
    modelCatalog.tabEnabled();
  });

  it('should show empty state when configmap has empty sources', () => {
    initIntercepts({ sources: [] });
    modelCatalog.visit();
    modelCatalog.findModelCatalogEmptyState().should('exist');
  });

  it('should display model catalog content when data is loaded', () => {
    initIntercepts({});
    modelCatalog.visit();
    modelCatalog.findLoadingState().should('not.exist');
    modelCatalog.findPageTitle().should('be.visible');
    modelCatalog.findPageDescription().should('be.visible');
    modelCatalog.findModelCatalogCards().should('have.length.at.least', 1);
  });

  it('should display model catalog filters', () => {
    initIntercepts({});
    modelCatalog.visit();
    modelCatalog.findFilter('Provider').should('be.visible');
    modelCatalog.findFilter('License').should('be.visible');
    modelCatalog.findFilter('Task').should('be.visible');
    modelCatalog.findFilter('Language').should('be.visible');
  });

  it('filters show more and show less button should work', () => {
    initIntercepts({});
    modelCatalog.visit();
    modelCatalog.findFilterShowMoreButton('Task').click({ scrollBehavior: false });
    modelCatalog.findFilterCheckbox('Task', 'text-generation').should('be.visible');
    modelCatalog.findFilterCheckbox('Task', 'text-to-text').should('be.visible');
    modelCatalog.findFilterCheckbox('Task', 'image-to-text').should('be.visible');
    modelCatalog.findFilterCheckbox('Task', 'image-text-to-text').should('be.visible');
    modelCatalog.findFilterCheckbox('Task', 'audio-to-text').should('be.visible');
    modelCatalog.findFilterCheckbox('Task', 'video-to-text').should('be.visible');
    modelCatalog.findFilterShowLessButton('Task').click({ scrollBehavior: false });
    modelCatalog.findFilterCheckbox('Task', 'audio-to-text').should('not.exist');
  });

  it('filters should be searchable', () => {
    initIntercepts({});
    modelCatalog.visit();
    modelCatalog.findFilterSearch('Task').type('audio-to-text');
    modelCatalog.findFilterCheckbox('Task', 'audio-to-text').should('be.visible');
    modelCatalog.findFilterCheckbox('Task', 'video-to-text').should('not.be.exist');
    modelCatalog.findFilterSearch('Task').type('test');
    modelCatalog.findFilterEmpty('Task').should('be.visible');
  });

  it('checkbox should work', () => {
    // Calculate expected category count based on sources
    const defaultSources = [
      mockCatalogSource({}),
      mockCatalogSource({ id: 'source-2', name: 'source 2' }),
    ];
    const uniqueLabels = new Set<string>();
    defaultSources.forEach((source) => {
      source.labels.forEach((label) => {
        if (label.trim()) {
          uniqueLabels.add(label.trim());
        }
      });
    });

    // Check if there are sources without labels
    const hasSourcesWithoutLabels = defaultSources.some(
      (source) =>
        source.enabled !== false &&
        (source.labels.length === 0 || source.labels.every((label) => !label.trim())),
    );

    // Expected count: unique labels + (1 if sources without labels exist)
    const expectedCategoryCount = uniqueLabels.size + (hasSourcesWithoutLabels ? 1 : 0);

    cy.interceptApi(
      `GET /api/:apiVersion/model_catalog/models`,
      {
        path: { apiVersion: MODEL_CATALOG_API_VERSION },
        query: { sourceLabel: '' },
      },
      mockCatalogModelList({
        items: [mockCatalogModel({})],
      }),
    ).as('getCatalogModelsBySource');

    initIntercepts({ sources: defaultSources });
    modelCatalog.visit();
    modelCatalog.findFilterCheckbox('Task', 'text-generation').click();
    modelCatalog.findFilterCheckbox('Task', 'text-to-text').click();
    modelCatalog.findFilterCheckbox('Provider', 'Google').click();

    // Wait for the expected number of API calls (one per category section when filters are applied)
    const waitCalls = Array.from(
      { length: expectedCategoryCount },
      () => '@getCatalogModelsBySource',
    );
    cy.wait(waitCalls).then((interceptions) => {
      const lastInterception = interceptions[interceptions.length - 1];
      expect(lastInterception.request.url).to.include(
        'tasks+IN+%28%27text-generation%27%2C%27text-to-text%27%29+AND+provider%3D%27Google%27',
      );
    });
  });
});

describe('Performance Empty State', () => {
  describe('Community & Custom Section', () => {
    it('should show performance empty state when toggle is ON', () => {
      initIntercepts({
        sources: [
          mockCatalogSource({ labels: ['Provider one'] }),
          mockCatalogSource({ id: 'custom-source', name: 'Custom Source', labels: [] }),
        ],
        hasValidatedModels: true,
      });
      modelCatalog.visit({ enableTempDevCatalogAdvancedFiltersFeature: true });

      modelCatalog.togglePerformanceView();
      modelCatalog.findPerformanceViewToggleValue().should('be.checked');

      modelCatalog.findCategoryToggle('no-labels').click();

      modelCatalog.findPerformanceEmptyState().should('be.visible');
      modelCatalog.findModelCatalogCards().should('not.exist');
    });

    it('should show models when toggle is OFF', () => {
      initIntercepts({
        sources: [
          mockCatalogSource({ labels: ['Provider one'] }),
          mockCatalogSource({ id: 'custom-source', name: 'Custom Source', labels: [] }),
        ],
      });
      modelCatalog.visit({ enableTempDevCatalogAdvancedFiltersFeature: true });

      modelCatalog.findCategoryToggle('no-labels').click();

      modelCatalog.findModelCatalogCards().should('have.length.at.least', 1);
      modelCatalog.findPerformanceEmptyState().should('not.exist');
    });
  });

  describe('Labeled Section Without Validated Models', () => {
    it('should show performance empty state when toggle is ON and no validated models', () => {
      initIntercepts({
        sources: [mockCatalogSource({ labels: ['Provider one'] })],
        hasValidatedModels: false,
      });
      modelCatalog.visit({ enableTempDevCatalogAdvancedFiltersFeature: true });

      modelCatalog.togglePerformanceView();
      modelCatalog.findCategoryToggle('label-Provider one').click();

      modelCatalog.findPerformanceEmptyState().should('be.visible');
    });

    it('should show models when toggle is OFF', () => {
      initIntercepts({
        sources: [mockCatalogSource({ labels: ['Provider one'] })],
        hasValidatedModels: false,
      });
      modelCatalog.visit({ enableTempDevCatalogAdvancedFiltersFeature: true });

      modelCatalog.findCategoryToggle('label-Provider one').click();

      modelCatalog.findModelCatalogCards().should('have.length.at.least', 1);
    });
  });

  describe('Labeled Section With Validated Models', () => {
    it('should show models when toggle is ON and section has validated models', () => {
      initIntercepts({
        sources: [mockCatalogSource({ labels: ['Provider one'] })],
        hasValidatedModels: true,
      });
      modelCatalog.visit({ enableTempDevCatalogAdvancedFiltersFeature: true });

      modelCatalog.togglePerformanceView();
      modelCatalog.findCategoryToggle('label-Provider one').click();

      modelCatalog.findModelCatalogCards().should('have.length.at.least', 1);
      modelCatalog.findPerformanceEmptyState().should('not.exist');
    });
  });

  describe('Empty State Actions', () => {
    it('should turn off toggle when clicking "set Explore model performance to off"', () => {
      initIntercepts({
        sources: [
          mockCatalogSource({ labels: ['Provider one'] }),
          mockCatalogSource({ id: 'custom-source', name: 'Custom Source', labels: [] }),
        ],
        hasValidatedModels: true,
      });
      modelCatalog.visit({ enableTempDevCatalogAdvancedFiltersFeature: true });

      modelCatalog.togglePerformanceView();
      modelCatalog.findCategoryToggle('no-labels').click();
      modelCatalog.findPerformanceEmptyState().should('be.visible');

      modelCatalog.findSetPerformanceOffLink().click();

      modelCatalog.findPerformanceViewToggleValue().should('not.be.checked');
      modelCatalog.findModelCatalogCards().should('have.length.at.least', 1);
    });

    it('should navigate to All models when clicking "Select the All models category"', () => {
      initIntercepts({
        sources: [
          mockCatalogSource({ labels: ['Provider one'] }),
          mockCatalogSource({ id: 'custom-source', name: 'Custom Source', labels: [] }),
        ],
        hasValidatedModels: true,
      });
      modelCatalog.visit({ enableTempDevCatalogAdvancedFiltersFeature: true });

      modelCatalog.togglePerformanceView();
      modelCatalog.findCategoryToggle('no-labels').click();
      modelCatalog.findPerformanceEmptyState().should('be.visible');

      modelCatalog.findSelectAllModelsCategoryButton().click();

      modelCatalog.findAllModelsToggle().find('button').should('have.attr', 'aria-pressed', 'true');
    });
  });

  it('should work correctly when toggling performance view', () => {
    initIntercepts({
      sources: [
        mockCatalogSource({ labels: ['Provider one'] }),
        mockCatalogSource({ id: 'custom-source', name: 'Custom Source', labels: [] }),
      ],
      hasValidatedModels: true,
    });
    modelCatalog.visit({ enableTempDevCatalogAdvancedFiltersFeature: true });

    modelCatalog.togglePerformanceView();
    modelCatalog.findPerformanceViewToggleValue().should('be.checked');
    modelCatalog.findCategoryToggle('no-labels').click();
    modelCatalog.findPerformanceEmptyState().should('be.visible');

    modelCatalog.findSetPerformanceOffLink().click();
    modelCatalog.findPerformanceViewToggleValue().should('not.be.checked');
    modelCatalog.findPerformanceEmptyState().should('not.exist');
    modelCatalog.findModelCatalogCards().should('have.length.at.least', 1);

    modelCatalog.togglePerformanceView();
    modelCatalog.findPerformanceViewToggleValue().should('be.checked');
    modelCatalog.findPerformanceEmptyState().should('be.visible');
    modelCatalog.findModelCatalogCards().should('not.exist');
  });
});
