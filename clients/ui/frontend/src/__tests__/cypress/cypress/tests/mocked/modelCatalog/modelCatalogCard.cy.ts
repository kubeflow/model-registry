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

describe('ModelCatalogCard Component', () => {
  beforeEach(() => {
    initIntercepts({});
    modelCatalog.visit();
    modelCatalog.navigate();
  });
  describe('Card Layout and Content', () => {
    it('should render all cards from the mock data', () => {
      modelCatalog.findModelCatalogCards().should('have.length.at.least', 1);
    });

    it('should display correct source labels', () => {
      modelCatalog.findFirstModelCatalogCard().within(() => {
        modelCatalog.findSourceLabel().should('contain.text', 'sample source');
      });
    });

    it('should handle cards with logos', () => {
      modelCatalog.findFirstModelCatalogCard().within(() => {
        modelCatalog
          .findModelLogo()
          .should('exist')
          .and('have.attr', 'src')
          .and('include', 'data:image/svg+xml;base64');
      });
    });
  });

  describe('Description Handling', () => {
    it('should display model descriptions', () => {
      modelCatalog.findFirstModelCatalogCard().within(() => {
        modelCatalog
          .findModelCatalogDescription()
          .should(
            'contain.text',
            'Granite-8B-Code-Instruct is a 8B parameter model fine tuned from\nGranite-8B-Code-Base on a combination of permissively licensed instruction\ndata to enhance instruction following capabilities including logical\nreasoning and problem-solving skills.',
          );
      });
    });
  });

  describe('Navigation and Interaction', () => {
    it('should show all model metadata correctly', () => {
      modelCatalog.findFirstModelCatalogCard().within(() => {
        modelCatalog.findModelCatalogDetailLink().should('contain.text', 'model1');
        modelCatalog.findTaskLabel().should('exist');
        modelCatalog.findLicenseLabel().should('exist');
        modelCatalog.findProviderLabel().should('exist');
      });
    });
  });
});
