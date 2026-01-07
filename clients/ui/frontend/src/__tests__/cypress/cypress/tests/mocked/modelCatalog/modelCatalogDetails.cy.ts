/* eslint-disable camelcase */
import {
  mockCatalogModel,
  mockCatalogModelArtifactList,
  mockCatalogModelList,
  mockCatalogSource,
  mockCatalogSourceList,
  mockCatalogModelArtifact,
} from '~/__mocks__';
import { modelCatalog } from '~/__tests__/cypress/cypress/pages/modelCatalog';
import { mockModelRegistry } from '~/__mocks__/mockModelRegistry';
import type { CatalogSource } from '~/app/modelCatalogTypes';
import { MODEL_CATALOG_API_VERSION } from '~/__tests__/cypress/cypress/support/commands/api';
import { mockCatalogFilterOptionsList } from '~/__mocks__/mockCatalogFilterOptionsList';
import { ModelRegistryMetadataType } from '~/app/types';

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
          query: { sourceLabel: label },
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

  cy.interceptApi(
    `GET /api/:apiVersion/model_catalog/sources/:sourceId/models/:modelName`,
    {
      path: {
        apiVersion: MODEL_CATALOG_API_VERSION,
        sourceId: 'source-2',
        modelName: 'sample%20category%201-model-1',
      },
    },
    mockCatalogModel({}),
  );

  cy.interceptApi(
    `GET /api/:apiVersion/model_catalog/sources/:sourceId/artifacts/:modelName`,
    {
      path: {
        apiVersion: MODEL_CATALOG_API_VERSION,
        sourceId: 'source-2',
        modelName: 'sample%20category%201-model-1',
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
};

describe('Model Catalog Details Page', () => {
  beforeEach(() => {
    // Mock model registries for register button functionality
    cy.intercept('GET', '/model-registry/api/v1/model_registry*', [
      mockModelRegistry({ name: 'modelregistry-sample' }),
    ]).as('getModelRegistries');

    initIntercepts({});
    modelCatalog.visit();
  });

  it('navigates to details and shows header, breadcrumb and description', () => {
    modelCatalog.findLoadingState().should('not.exist');
    modelCatalog.findModelCatalogDetailLink().first().click();
    modelCatalog.findBreadcrumb().should('exist');
    modelCatalog.findDetailsProviderText().should('be.visible');
    modelCatalog.findDetailsDescription().should('exist');
  });

  it('does not show architecture field when no architectures are available', () => {
    modelCatalog.findLoadingState().should('not.exist');
    modelCatalog.findModelCatalogDetailLink().first().click();
    // Architecture field should not exist when no valid architectures
    modelCatalog.findModelArchitecture().should('not.exist');
  });

  it('shows architecture field with valid architectures', () => {
    // Override the artifacts intercept with architecture data
    cy.interceptApi(
      `GET /api/:apiVersion/model_catalog/sources/:sourceId/artifacts/:modelName`,
      {
        path: {
          apiVersion: MODEL_CATALOG_API_VERSION,
          sourceId: 'source-2',
          modelName: 'sample%20category%201-model-1',
        },
      },
      mockCatalogModelArtifactList({
        items: [
          mockCatalogModelArtifact({
            customProperties: {
              architecture: {
                string_value: '["amd64", "arm64", "s390x"]',
                metadataType: ModelRegistryMetadataType.STRING,
              },
            },
          }),
        ],
      }),
    );

    modelCatalog.findLoadingState().should('not.exist');
    modelCatalog.findModelCatalogDetailLink().first().click();
    modelCatalog.findBreadcrumb().should('exist');

    // Architecture field should exist and show correct values
    modelCatalog.findModelArchitecture().should('be.visible');
    modelCatalog.findModelArchitecture().should('contain.text', 'amd64, arm64, s390x');
  });

  it('shows architecture field with uppercase architectures normalized to lowercase', () => {
    // Override the artifacts intercept with uppercase architecture data
    cy.interceptApi(
      `GET /api/:apiVersion/model_catalog/sources/:sourceId/artifacts/:modelName`,
      {
        path: {
          apiVersion: MODEL_CATALOG_API_VERSION,
          sourceId: 'source-2',
          modelName: 'sample%20category%201-model-1',
        },
      },
      mockCatalogModelArtifactList({
        items: [
          mockCatalogModelArtifact({
            customProperties: {
              architecture: {
                string_value: '["AMD64", "ARM64"]',
                metadataType: ModelRegistryMetadataType.STRING,
              },
            },
          }),
        ],
      }),
    );

    modelCatalog.findLoadingState().should('not.exist');
    modelCatalog.findModelCatalogDetailLink().first().click();

    // Architecture should be normalized to lowercase
    modelCatalog.findModelArchitecture().should('be.visible');
    modelCatalog.findModelArchitecture().should('contain.text', 'amd64, arm64');
  });

  it('does not show architecture field when only invalid architectures are present', () => {
    // Override the artifacts intercept with invalid architecture data
    cy.interceptApi(
      `GET /api/:apiVersion/model_catalog/sources/:sourceId/artifacts/:modelName`,
      {
        path: {
          apiVersion: MODEL_CATALOG_API_VERSION,
          sourceId: 'source-2',
          modelName: 'sample%20category%201-model-1',
        },
      },
      mockCatalogModelArtifactList({
        items: [
          mockCatalogModelArtifact({
            customProperties: {
              architecture: {
                string_value: '["invalid-arch", "unknown"]',
                metadataType: ModelRegistryMetadataType.STRING,
              },
            },
          }),
        ],
      }),
    );

    modelCatalog.findLoadingState().should('not.exist');
    modelCatalog.findModelCatalogDetailLink().first().click();

    // Architecture field should not exist when all architectures are invalid
    modelCatalog.findModelArchitecture().should('not.exist');
  });
});
