/* eslint-disable camelcase */
import { modelCatalog } from '~/__tests__/cypress/cypress/pages/modelCatalog';
import {
  mockCatalogAccuracyMetricsArtifact,
  mockCatalogModel,
  mockCatalogModelArtifact,
  mockCatalogModelArtifactList,
  mockCatalogModelList,
  mockCatalogPerformanceMetricsArtifact,
  mockCatalogSource,
  mockCatalogSourceList,
  mockNonValidatedModel,
  mockValidatedModel,
} from '~/__mocks__';
import type { CatalogSource } from '~/app/modelCatalogTypes';
import { MODEL_CATALOG_API_VERSION } from '~/__tests__/cypress/cypress/support/commands/api';
import { mockCatalogFilterOptionsList } from '~/__mocks__/mockCatalogFilterOptionsList';
import type { ModelRegistryCustomProperties } from '~/app/types';
import { ModelRegistryMetadataType } from '~/app/types';

type HandlersProps = {
  sources?: CatalogSource[];
  modelsPerCategory?: number;
  useValidatedModel?: boolean;
};

const initIntercepts = ({
  sources = [mockCatalogSource({}), mockCatalogSource({ id: 'source-2', name: 'source 2' })],
  modelsPerCategory = 4,
  useValidatedModel = false,
}: HandlersProps) => {
  const testModel = useValidatedModel ? mockValidatedModel : mockNonValidatedModel;

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

  cy.interceptApi(
    `GET /api/:apiVersion/model_catalog/sources/:sourceId/artifacts/:modelName`,
    {
      path: {
        apiVersion: MODEL_CATALOG_API_VERSION,
        sourceId: 'source-2',
        modelName: testModel.name.replace('/', '%2F'),
      },
    },
    mockCatalogModelArtifactList({}),
  );

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
        sourceId: 'source-2',
        modelName: 'validated-model',
      },
    },
    {
      items: [
        mockCatalogPerformanceMetricsArtifact({}),
        mockCatalogPerformanceMetricsArtifact({
          customProperties: {
            hardware_type: {
              metadataType: ModelRegistryMetadataType.STRING,
              string_value: 'RTX 4090',
            },
            hardware_count: {
              metadataType: ModelRegistryMetadataType.INT,
              int_value: '33',
            },
            requests_per_second: {
              metadataType: ModelRegistryMetadataType.DOUBLE,
              double_value: 10,
            },
            ttft_mean: {
              metadataType: ModelRegistryMetadataType.DOUBLE,
              double_value: 67.14892749816,
            },
          },
        }),
        mockCatalogPerformanceMetricsArtifact({
          customProperties: {
            hardware_type: {
              metadataType: ModelRegistryMetadataType.STRING,
              string_value: 'A100',
            },
            hardware_count: {
              metadataType: ModelRegistryMetadataType.INT,
              int_value: '40',
            },
            requests_per_second: {
              metadataType: ModelRegistryMetadataType.DOUBLE,
              double_value: 15,
            },
            ttft_mean: {
              metadataType: ModelRegistryMetadataType.DOUBLE,
              double_value: 42.123791232,
            },
          },
        }),
        mockCatalogAccuracyMetricsArtifact({}),
        mockCatalogModelArtifact({}),
      ],
    },
  ).as('getCatalogModelArtifacts');
};

describe('ModelCatalogCard Component', () => {
  beforeEach(() => {
    initIntercepts({});
    modelCatalog.visit();
  });
  describe('Card Layout and Content', () => {
    it('should render all cards from the mock data', () => {
      modelCatalog.findModelCatalogCards().should('have.length.at.least', 1);
    });

    it('should display correct source labels', () => {
      modelCatalog.findFirstModelCatalogCard().within(() => {
        modelCatalog
          .findSourceLabel()
          .should('contain.text', 'source 2text-generationprovider1apache-2.0');
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
        modelCatalog
          .findModelCatalogDetailLink()
          .should('contain.text', 'sample category 1-model-1');
        modelCatalog.findTaskLabel().should('exist');
        modelCatalog.findProviderLabel().should('exist');
      });
    });
  });

  describe('Validated Model', () => {
    beforeEach(() => {
      initIntercepts({ useValidatedModel: true });
      modelCatalog.visit();
    });
    it('should show validated model correctly', () => {
      cy.wait('@getCatalogModelArtifacts');
      modelCatalog.findFirstModelCatalogCard().within(() => {
        modelCatalog.findValidatedModelHardware().should('contain.text', '2xH100-80');
        modelCatalog.findValidatedModelRps().should('contain.text', '7');
        modelCatalog.findValidatedModelTtft().should('contain.text', '35.48818160947744');
        modelCatalog.findValidatedModelBenchmarkNext().click();
        modelCatalog.findValidatedModelHardware().should('contain.text', '33xRTX 4090');
        modelCatalog.findValidatedModelRps().should('contain.text', '10');
        modelCatalog.findValidatedModelTtft().should('contain.text', '67.14892749816');
        modelCatalog.findValidatedModelBenchmarkNext().click();
        modelCatalog.findValidatedModelHardware().should('contain.text', '40xA100');
        modelCatalog.findValidatedModelRps().should('contain.text', '15');
        modelCatalog.findValidatedModelTtft().should('contain.text', '42.123791232');
        modelCatalog.findValidatedModelBenchmarkPrev().click();
        modelCatalog.findValidatedModelHardware().should('contain.text', '33xRTX 4090');
        modelCatalog.findValidatedModelRps().should('contain.text', '10');
        modelCatalog.findValidatedModelTtft().should('contain.text', '67.14892749816');
        modelCatalog.findValidatedModelBenchmarkLink().click();
        cy.url().should('include', 'performance-insights');
      });
    });
  });
});
