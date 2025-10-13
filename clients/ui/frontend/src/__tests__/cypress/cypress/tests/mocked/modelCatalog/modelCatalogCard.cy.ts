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
import { ModelRegistryCustomProperties, ModelRegistryMetadataType } from '~/app/types';

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

  sources.forEach((source, sourceIndex) => {
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
              sourceIndex === 1
                ? ({
                    validated: {
                      metadataType: ModelRegistryMetadataType.STRING,
                      // eslint-disable-next-line camelcase
                      string_value: '',
                    },
                  } as ModelRegistryCustomProperties)
                : undefined;

            return mockCatalogModel({
              name: `${label.toLowerCase()}-model-${i + 1}`,
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
    `GET /api/:apiVersion/model_catalog/models/filter_options`,
    {
      path: { apiVersion: MODEL_CATALOG_API_VERSION },
      query: { namespace: 'kubeflow' },
    },
    mockCatalogFilterOptionsList(),
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
        modelCatalog
          .findSourceLabel()
          .should('contain.text', 'source 2text-generationapache-2.0provider1');
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
    it('should show model metadata correctly', () => {
      modelCatalog.findFirstModelCatalogCard().within(() => {
        modelCatalog.findModelCatalogDetailLink().should('contain.text', 'community-model-1');
        modelCatalog.findTaskLabel().should('exist');
        modelCatalog.findProviderLabel().should('exist');
      });
    });
  });

  describe('Validated Model', () => {
    it('should show validated model correctly', () => {
      modelCatalog.findLastModelCatalogCard().within(() => {
        modelCatalog.findValidatedModelBenchmarkLink().click();
        cy.url().should('include', 'performance-insights');
      });
    });
  });
});
