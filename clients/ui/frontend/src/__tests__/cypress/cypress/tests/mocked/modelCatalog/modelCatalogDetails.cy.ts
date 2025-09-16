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

describe('Model Catalog Details Page', () => {
  beforeEach(() => {
    // Mock model registries for register button functionality
    cy.intercept('GET', '/model-registry/api/v1/model_registry*', [
      mockModelRegistry({ name: 'modelregistry-sample' }),
    ]).as('getModelRegistries');

    initIntercepts({});
    modelCatalog.visit();
    modelCatalog.navigate();
  });

  it('navigates to details and shows header, breadcrumb and description', () => {
    modelCatalog.findLoadingState().should('not.exist');
    modelCatalog.findModelCatalogDetailLink().first().click();
    modelCatalog.findBreadcrumb().should('be.visible');
    modelCatalog.findDetailsProviderText().should('be.visible');
    modelCatalog.findDetailsDescription().should('exist');
  });
});
