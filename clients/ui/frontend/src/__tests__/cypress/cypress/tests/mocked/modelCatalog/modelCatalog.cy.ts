import { modelCatalog } from '~/__tests__/cypress/cypress/pages/modelCatalog';
import {
  mockCatalogModel,
  mockCatalogModelList,
  mockCatalogSource,
  mockCatalogSourceList,
} from '~/__mocks__';
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
});
