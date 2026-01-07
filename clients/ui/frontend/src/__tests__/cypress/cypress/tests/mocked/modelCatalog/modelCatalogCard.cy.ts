/* eslint-disable camelcase */
import { modelCatalog } from '~/__tests__/cypress/cypress/pages/modelCatalog';
import {
  mockCatalogModel,
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
    `GET /api/:apiVersion/model_catalog/sources/:sourceId/performance_artifacts/:modelName`,
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

  // The /performance_artifacts endpoint only returns performance metrics artifacts
  // (no accuracy or model artifacts - those are filtered server-side)
  cy.interceptApi(
    `GET /api/:apiVersion/model_catalog/sources/:sourceId/performance_artifacts/:modelName`,
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
              double_value: 67.15,
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
              double_value: 42.12,
            },
          },
        }),
      ],
      pageSize: 10,
      size: 3,
      nextPageToken: '',
    },
  ).as('getCatalogSourceModelArtifacts');
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
        modelCatalog.findSourceLabel().should('contain.text', 'source 2text-generationprovider1');
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
    describe('Toggle OFF (default)', () => {
      beforeEach(() => {
        initIntercepts({ useValidatedModel: true });
        modelCatalog.visit();
      });

      it('should show description with View benchmarks link when toggle is OFF', () => {
        cy.wait('@getCatalogSourceModelArtifacts');
        modelCatalog.findFirstModelCatalogCard().within(() => {
          // Should show description
          modelCatalog.findModelCatalogDescription().should('be.visible');

          // Should show "View X benchmarks" link
          modelCatalog.findValidatedModelBenchmarkLink().should('be.visible');
          modelCatalog
            .findValidatedModelBenchmarkLink()
            .should('contain.text', 'View 3 benchmarks');

          // Should NOT show hardware, replicas, TTFT metrics when toggle is OFF
          modelCatalog.findValidatedModelHardware().should('not.exist');
          modelCatalog.findValidatedModelReplicas().should('not.exist');
          modelCatalog.findValidatedModelTtft().should('not.exist');
        });
      });

      it('should navigate to Performance Insights tab when clicking View benchmarks link', () => {
        cy.wait('@getCatalogSourceModelArtifacts');
        modelCatalog.findFirstModelCatalogCard().within(() => {
          modelCatalog.findValidatedModelBenchmarkLink().click();
        });
        cy.url().should('include', 'performance-insights');
      });
    });

    describe('Toggle ON', () => {
      beforeEach(() => {
        initIntercepts({ useValidatedModel: true });
        // Enable feature flag and visit
        modelCatalog.visit({ enableTempDevCatalogAdvancedFiltersFeature: true });
        cy.wait('@getCatalogSourceModelArtifacts');
        // Turn the toggle ON before each test in this block
        modelCatalog.togglePerformanceView();
      });

      it('should show validated model metrics correctly when toggle is ON', () => {
        modelCatalog.findFirstModelCatalogCard().within(() => {
          // Should show hardware, replicas, TTFT metrics
          modelCatalog.findValidatedModelHardware().should('contain.text', '2xH100-80');
          modelCatalog.findValidatedModelReplicas().should('contain.text', '7');
          modelCatalog.findValidatedModelTtft().should('contain.text', '35.49');

          // Should NOT show description when toggle is ON
          modelCatalog.findModelCatalogDescription().should('not.exist');

          // Navigate through benchmarks
          modelCatalog.findValidatedModelBenchmarkNext().click();
          modelCatalog.findValidatedModelHardware().should('contain.text', '33xRTX 4090');
          modelCatalog.findValidatedModelReplicas().should('contain.text', '10');
          modelCatalog.findValidatedModelTtft().should('contain.text', '67.15');

          modelCatalog.findValidatedModelBenchmarkNext().click();
          modelCatalog.findValidatedModelHardware().should('contain.text', '40xA100');
          modelCatalog.findValidatedModelReplicas().should('contain.text', '15');
          modelCatalog.findValidatedModelTtft().should('contain.text', '42.12');

          modelCatalog.findValidatedModelBenchmarkPrev().click();
          modelCatalog.findValidatedModelHardware().should('contain.text', '33xRTX 4090');
          modelCatalog.findValidatedModelReplicas().should('contain.text', '10');
          modelCatalog.findValidatedModelTtft().should('contain.text', '67.15');

          // Click benchmark link to navigate to Performance Insights
          modelCatalog.findValidatedModelBenchmarkLink().click();
        });
        cy.url().should('include', 'performance-insights');
      });

      it('should navigate through benchmarks correctly', () => {
        modelCatalog.findFirstModelCatalogCard().within(() => {
          // Initial state - first benchmark
          modelCatalog.findValidatedModelHardware().should('contain.text', '2xH100-80');

          // Navigate to next benchmark
          modelCatalog.findValidatedModelBenchmarkNext().click();
          modelCatalog.findValidatedModelHardware().should('contain.text', '33xRTX 4090');

          // Navigate to next benchmark
          modelCatalog.findValidatedModelBenchmarkNext().click();
          modelCatalog.findValidatedModelHardware().should('contain.text', '40xA100');

          // Navigate back
          modelCatalog.findValidatedModelBenchmarkPrev().click();
          modelCatalog.findValidatedModelHardware().should('contain.text', '33xRTX 4090');
        });
      });
    });
  });
});
