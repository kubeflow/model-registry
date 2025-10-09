import { modelCatalog } from '~/__tests__/cypress/cypress/pages/modelCatalog';
import {
  mockCatalogModel,
  mockCatalogModelList,
  mockCatalogSource,
  mockCatalogSourceList,
} from '~/__mocks__';
import type { CatalogSource } from '~/app/modelCatalogTypes';
import { MODEL_CATALOG_API_VERSION } from '~/__tests__/cypress/cypress/support/commands/api';
import { mockCatalogFilterOptionsList } from '~/__mocks__/mockCatalogFilterOptionsList';

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
  ).as('getCatalogModelsBySource');

  cy.interceptApi(
    `GET /api/:apiVersion/model_catalog/models/filter_options`,
    {
      path: { apiVersion: MODEL_CATALOG_API_VERSION },
      query: { namespace: 'kubeflow' },
    },
    mockCatalogFilterOptionsList(),
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

describe('Model Catalog Page', () => {
  it('model catalog tab should be enabled', () => {
    initIntercepts({});
    modelCatalog.visit();
    modelCatalog.navigate();
    modelCatalog.tabEnabled();
  });

  it('should show empty state when configmap has empty sources', () => {
    initIntercepts({ sources: [] });
    modelCatalog.visit();
    modelCatalog.navigate();
    modelCatalog.visit();
    modelCatalog.findModelCatalogEmptyState().should('exist');
  });

  it('should display model catalog content when data is loaded', () => {
    initIntercepts({});
    modelCatalog.visit();
    modelCatalog.navigate();
    modelCatalog.findLoadingState().should('not.exist');
    modelCatalog.findPageTitle().should('be.visible');
    modelCatalog.findPageDescription().should('be.visible');
    modelCatalog.findModelCatalogCards().should('have.length.at.least', 1);
  });

  it('should display model catalog filters', () => {
    initIntercepts({});
    modelCatalog.visit();
    modelCatalog.navigate();
    modelCatalog.findFilter('Provider').should('be.visible');
    modelCatalog.findFilter('License').should('be.visible');
    modelCatalog.findFilter('Task').should('be.visible');
    modelCatalog.findFilter('Language').should('be.visible');
  });

  it('filters show more and show less button should work', () => {
    initIntercepts({});
    modelCatalog.visit();
    modelCatalog.navigate();
    modelCatalog.findFilterShowMoreButton('Task').click();
    modelCatalog.findFilterCheckbox('Task', 'text-generation').should('be.visible');
    modelCatalog.findFilterCheckbox('Task', 'text-to-text').should('be.visible');
    modelCatalog.findFilterCheckbox('Task', 'image-to-text').should('be.visible');
    modelCatalog.findFilterCheckbox('Task', 'image-text-to-text').should('be.visible');
    modelCatalog.findFilterCheckbox('Task', 'audio-to-text').should('be.visible');
    modelCatalog.findFilterCheckbox('Task', 'video-to-text').should('be.visible');
    modelCatalog.findFilterShowLessButton('Task').click();
    modelCatalog.findFilterCheckbox('Task', 'audio-to-text').should('not.exist');
  });

  it('filters should be searchable', () => {
    initIntercepts({});
    modelCatalog.visit();
    modelCatalog.navigate();
    modelCatalog.findFilterSearch('Task').type('audio-to-text');
    modelCatalog.findFilterCheckbox('Task', 'audio-to-text').should('be.visible');
    modelCatalog.findFilterCheckbox('Task', 'video-to-text').should('not.be.exist');
    modelCatalog.findFilterSearch('Task').type('test');
    modelCatalog.findFilterEmpty('Task').should('be.visible');
  });

  it('checkbox should work', () => {
    initIntercepts({});
    modelCatalog.visit();
    modelCatalog.navigate();
    modelCatalog.findFilterCheckbox('Task', 'text-generation').click();
    modelCatalog.findFilterCheckbox('Task', 'text-to-text').click();
    modelCatalog.findFilterCheckbox('Provider', 'Google').click();
    cy.wait([
      '@getCatalogModelsBySource',
      '@getCatalogModelsBySource',
      '@getCatalogModelsBySource',
      '@getCatalogModelsBySource',
      '@getCatalogModelsBySource',
    ]).then((interceptions) => {
      const lastInterception = interceptions[interceptions.length - 1];
      expect(lastInterception.request.url).to.include(
        'tasks%2BIN%2B%28%27text-generation%27%2C%27text-to-text%27%29%2BAND%2Bprovider%2B%3D%2B%27Google%27',
      );
    });
  });
});
