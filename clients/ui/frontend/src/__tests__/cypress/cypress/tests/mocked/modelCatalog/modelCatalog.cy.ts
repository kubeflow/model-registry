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

type HandlersProps = {
  sources?: CatalogSource[];
  modelsPerCategory?: number;
};
const initIntercepts = ({
  sources = [mockCatalogSource({}), mockCatalogSource({ id: 'source-2', name: 'source 2' })],
  modelsPerCategory = 4,
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
      cy.interceptApi(
        `GET /api/:apiVersion/model_catalog/models`,
        {
          path: { apiVersion: MODEL_CATALOG_API_VERSION },
          query: {
            sourceLabel: label,
          },
        },
        mockCatalogModelList({
          items: Array.from({ length: modelsPerCategory }, (_, i) =>
            mockCatalogModel({
              name: `${label.toLowerCase()}-model-${i + 1}`,
              // eslint-disable-next-line camelcase
              source_id: source.id,
            }),
          ),
        }),
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
