/* eslint-disable camelcase */
import { mockModArchResponse } from 'mod-arch-core';
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

type FilteredModelsInterceptConfig = {
  returnModelsForFilters?: boolean;
  modelsToReturn?: ReturnType<typeof mockCatalogModel>[];
};

const setupFilteredModelsIntercept = ({
  returnModelsForFilters = false,
  modelsToReturn = [],
}: FilteredModelsInterceptConfig = {}) => {
  cy.intercept(
    {
      method: 'GET',
      url: '**/model_catalog/models*filterQuery=*',
    },
    (req) => {
      const items = returnModelsForFilters ? modelsToReturn : [];
      req.reply(mockModArchResponse(mockCatalogModelList({ items })));
    },
  ).as('getFilteredModels');
};

const generateMockModels = (
  count: number,
  hasValidated: boolean,
  namePrefix: string,
  sourceId?: string,
): ReturnType<typeof mockCatalogModel>[] => {
  if (hasValidated) {
    return [
      mockCatalogModel({
        name: 'validated-model',
        source_id: sourceId,
        customProperties: {
          validated: { metadataType: ModelRegistryMetadataType.STRING, string_value: '' },
        },
      }),
      ...Array.from({ length: count - 1 }, (_, i) =>
        mockCatalogModel({ name: `${namePrefix}-${i + 1}`, source_id: sourceId }),
      ),
    ];
  }
  return Array.from({ length: count }, (_, i) =>
    mockCatalogModel({ name: `${namePrefix}-${i + 1}`, source_id: sourceId }),
  );
};

type HandlersProps = {
  sources?: CatalogSource[];
  modelsPerCategory?: number;
  hasValidatedModels?: boolean;
  includeAllModelsIntercept?: boolean;
};

const initIntercepts = ({
  sources = [mockCatalogSource({}), mockCatalogSource({ id: 'source-2', name: 'source 2' })],
  modelsPerCategory = 4,
  hasValidatedModels = false,
  includeAllModelsIntercept = true,
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
      const models = generateMockModels(
        modelsPerCategory,
        hasValidatedModels,
        `${label.toLowerCase().replace(/\s+/g, '-')}-model`,
        source.id,
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
            source_id: sources.find((s) => s.labels.length === 0)?.id || 'custom-source',
          }),
        ),
      }),
    );
  }
  if (includeAllModelsIntercept) {
    cy.intercept(
      {
        method: 'GET',
        url: new RegExp(
          `/api/${MODEL_CATALOG_API_VERSION}/model_catalog/models\\?(?!.*sourceLabel=)`,
        ),
      },
      (req) => {
        const models = generateMockModels(modelsPerCategory, hasValidatedModels, 'all-model');
        req.reply(mockModArchResponse(mockCatalogModelList({ items: models })));
      },
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
  );

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

    initIntercepts({ sources: defaultSources, includeAllModelsIntercept: false });

    setupFilteredModelsIntercept({
      returnModelsForFilters: true,
      modelsToReturn: [mockCatalogModel({})],
    });

    modelCatalog.visit();
    modelCatalog.findFilterShowMoreButton('Task').click();
    modelCatalog.findFilterCheckbox('Task', 'text-generation').click();
    modelCatalog.findFilterCheckbox('Task', 'text-to-text').click();
    modelCatalog.findFilterCheckbox('Provider', 'Google').click();

    // Wait for the expected number of API calls (one per category section when filters are applied)
    const waitCalls = Array.from({ length: expectedCategoryCount }, () => '@getFilteredModels');
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
      setupFilteredModelsIntercept({ returnModelsForFilters: false });
      modelCatalog.visit();

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
      modelCatalog.visit();

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
      setupFilteredModelsIntercept({ returnModelsForFilters: false });
      modelCatalog.visit();

      modelCatalog.togglePerformanceView();
      modelCatalog.findCategoryToggle('label-Provider one').click();

      modelCatalog.findPerformanceEmptyState().should('be.visible');
    });

    it('should show models when toggle is OFF', () => {
      initIntercepts({
        sources: [mockCatalogSource({ labels: ['Provider one'] })],
        hasValidatedModels: false,
      });
      modelCatalog.visit();

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
      modelCatalog.visit();

      modelCatalog.togglePerformanceView();
      modelCatalog.findCategoryToggle('label-Provider one').click();

      modelCatalog.findModelCatalogCards().should('have.length.at.least', 1);
      modelCatalog.findPerformanceEmptyState().should('not.exist');
    });
  });

  describe('Empty State Actions', () => {
    it('should turn off toggle when clicking "Turn Model performance view off"', () => {
      initIntercepts({
        sources: [
          mockCatalogSource({ labels: ['Provider one'] }),
          mockCatalogSource({ id: 'custom-source', name: 'Custom Source', labels: [] }),
        ],
        hasValidatedModels: true,
      });
      setupFilteredModelsIntercept({ returnModelsForFilters: false });
      modelCatalog.visit();

      modelCatalog.togglePerformanceView();
      modelCatalog.findCategoryToggle('no-labels').click();
      modelCatalog.findPerformanceEmptyState().should('be.visible');

      modelCatalog.findSetPerformanceOffLink().click();

      modelCatalog.findPerformanceViewToggleValue().should('not.be.checked');
      modelCatalog.findModelCatalogCards().should('have.length.at.least', 1);
    });

    it('should navigate to All models when clicking "View all models with performance data"', () => {
      initIntercepts({
        sources: [
          mockCatalogSource({ labels: ['Provider one'] }),
          mockCatalogSource({ id: 'custom-source', name: 'Custom Source', labels: [] }),
        ],
        hasValidatedModels: true,
      });
      setupFilteredModelsIntercept({ returnModelsForFilters: false });
      modelCatalog.visit();

      modelCatalog.togglePerformanceView();
      modelCatalog.findCategoryToggle('no-labels').click();
      modelCatalog.findPerformanceEmptyState().should('be.visible');

      modelCatalog.findSelectAllModelsCategoryButton().click();

      modelCatalog.findAllModelsToggle().find('button').should('have.attr', 'aria-pressed', 'true');
    });

    it('should show performance empty state after clicking Reset filters when toggle is ON', () => {
      initIntercepts({
        sources: [mockCatalogSource({ labels: ['Provider one'] })],
        hasValidatedModels: false,
      });
      setupFilteredModelsIntercept({ returnModelsForFilters: false });
      modelCatalog.visit();

      modelCatalog.togglePerformanceView();
      modelCatalog.findCategoryToggle('label-Provider one').click();
      modelCatalog.findPerformanceEmptyState().should('be.visible');

      modelCatalog.findFilterShowMoreButton('Task').click();
      modelCatalog.findFilterCheckbox('Task', 'audio-to-text').click();

      modelCatalog.findModelCatalogEmptyState().should('contain.text', 'No result found');

      modelCatalog.findEmptyStateResetFiltersButton().click();
      modelCatalog.findPerformanceEmptyState().should('be.visible');
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
    setupFilteredModelsIntercept({ returnModelsForFilters: false });
    modelCatalog.visit();

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

  it('should show "No result found" when toggle is ON and user applies filter that returns 0 results', () => {
    initIntercepts({
      sources: [mockCatalogSource({ labels: ['Provider one'] })],
      hasValidatedModels: true,
    });
    setupFilteredModelsIntercept({ returnModelsForFilters: false });

    modelCatalog.visit();
    modelCatalog.togglePerformanceView();
    modelCatalog.findCategoryToggle('label-Provider one').click();
    modelCatalog.findFilterShowMoreButton('Task').click();
    modelCatalog.findFilterCheckbox('Task', 'audio-to-text').click();
    modelCatalog.findPerformanceEmptyState().should('not.exist');
    modelCatalog.findModelCatalogEmptyState().should('contain.text', 'No result found');
    modelCatalog.findAllModelsToggle().click();
    modelCatalog.findPerformanceEmptyState().should('not.exist');
    modelCatalog.findModelCatalogEmptyState().should('contain.text', 'No result found');
  });
});

describe('All Models Section', () => {
  it('should show models in All models section even when toggle is ON', () => {
    initIntercepts({
      sources: [mockCatalogSource({ labels: ['Provider one'] })],
      hasValidatedModels: true,
    });
    modelCatalog.visit();

    modelCatalog.togglePerformanceView();
    modelCatalog.findPerformanceViewToggleValue().should('be.checked');

    modelCatalog.findAllModelsToggle().click();
    modelCatalog.findModelCatalogCards().should('have.length.at.least', 1);
    modelCatalog.findPerformanceEmptyState().should('not.exist');
  });
});
