/* eslint-disable camelcase */
/**
 * Shared intercept helpers for Model Catalog tests.
 * These functions consolidate common API intercept patterns used across multiple test files.
 */
import { mockModArchResponse } from 'mod-arch-core';
import {
  mockCatalogModel,
  mockCatalogModelList,
  mockCatalogPerformanceMetricsArtifact,
  mockCatalogSource,
  mockCatalogSourceList,
  mockNonValidatedModel,
  mockValidatedModel,
} from '~/__mocks__';
import { mockCatalogPerformanceMetricsArtifactList } from '~/__mocks__/mockCatalogModelArtifactList';
import { mockCatalogFilterOptionsList } from '~/__mocks__/mockCatalogFilterOptionsList';
import { mockModelRegistry } from '~/__mocks__/mockModelRegistry';
import type { CatalogModel, CatalogSource } from '~/app/modelCatalogTypes';
import type { ModelRegistryCustomProperties } from '~/app/types';
import { ModelRegistryMetadataType } from '~/app/types';
import { MODEL_CATALOG_API_VERSION } from '~/__tests__/cypress/cypress/support/commands/api';

/**
 * Options for setting up model catalog intercepts
 */
export type ModelCatalogInterceptOptions = {
  /** Catalog sources to mock. Defaults to two sample sources. */
  sources?: CatalogSource[];
  /** Number of models per category/label. Defaults to 4. */
  modelsPerCategory?: number;
  /** Whether to include a validated model (with performance data). Defaults to false. */
  useValidatedModel?: boolean;
  /** Whether to include performance artifacts. Defaults to false. */
  includePerformanceArtifacts?: boolean;
  /** Whether to include filter options. Defaults to true. */
  includeFilterOptions?: boolean;
  /** Whether to include model registry intercept. Defaults to false. */
  includeModelRegistry?: boolean;
  /** Custom validated model to use. Defaults to mockValidatedModel. */
  customValidatedModel?: CatalogModel;
  /** Custom non-validated model to use. Defaults to mockNonValidatedModel. */
  customNonValidatedModel?: CatalogModel;
};

/**
 * Default sources used when none are provided
 */
export const defaultSources = (): CatalogSource[] => [
  mockCatalogSource({}),
  mockCatalogSource({ id: 'source-2', name: 'source 2' }),
];

/**
 * Intercepts the GET /model_catalog/sources endpoint
 */
export const interceptSources = (sources: CatalogSource[]): void => {
  cy.interceptApi(
    `GET /api/:apiVersion/model_catalog/sources`,
    {
      path: { apiVersion: MODEL_CATALOG_API_VERSION },
    },
    mockCatalogSourceList({
      items: sources,
    }),
  );
};

/**
 * Creates mock models for a given label with optional validated model support
 */
export const createMockModelsForLabel = (
  source: CatalogSource,
  label: string,
  modelsPerCategory: number,
  useValidatedModel: boolean,
): ReturnType<typeof mockCatalogModelList> =>
  mockCatalogModelList({
    items: Array.from({ length: modelsPerCategory }, (_, i) => {
      const customProperties =
        i === 0 && useValidatedModel
          ? ({
              validated: {
                metadataType: ModelRegistryMetadataType.STRING,
                string_value: '',
              },
            } as ModelRegistryCustomProperties)
          : undefined;
      const name =
        i === 0 && useValidatedModel ? 'validated-model' : `${label.toLowerCase()}-model-${i + 1}`;

      return mockCatalogModel({
        name,
        source_id: source.id,
        customProperties,
      });
    }),
  });

/**
 * Intercepts models by label using regex to handle filterQuery parameters
 */
export const interceptModelsByLabel = (
  sources: CatalogSource[],
  modelsPerCategory: number,
  useValidatedModel: boolean,
): void => {
  sources.forEach((source) => {
    source.labels.forEach((label) => {
      const mockModels = createMockModelsForLabel(
        source,
        label,
        modelsPerCategory,
        useValidatedModel,
      );
      const encodedLabel = encodeURIComponent(label);

      // Use regex-based intercept to match requests with this sourceLabel
      // This handles both basic requests and requests with filterQuery
      cy.intercept(
        {
          method: 'GET',
          url: new RegExp(
            `/model-registry/api/${MODEL_CATALOG_API_VERSION}/model_catalog/models.*sourceLabel=${encodedLabel}`,
          ),
        },
        mockModArchResponse(mockModels),
      ).as(`getModels-${label}`);
    });
  });
};

/**
 * Intercepts models requests without sourceLabel (for "All models" / GalleryView)
 */
export const interceptAllModels = (modelsPerCategory: number, useValidatedModel: boolean): void => {
  const allModelsResponse = mockCatalogModelList({
    items: Array.from({ length: modelsPerCategory }, (_, i) => {
      const customProperties =
        i === 0 && useValidatedModel
          ? ({
              validated: {
                metadataType: ModelRegistryMetadataType.STRING,
                string_value: '',
              },
            } as ModelRegistryCustomProperties)
          : undefined;
      const name = i === 0 && useValidatedModel ? 'validated-model' : `all-models-model-${i + 1}`;
      return mockCatalogModel({
        name,
        source_id: 'sample-source',
        customProperties,
      });
    }),
  });

  // Intercept for GalleryView when filters are applied (no sourceLabel, but has filterQuery or pageSize)
  cy.intercept(
    {
      method: 'GET',
      url: new RegExp(
        `/model-registry/api/${MODEL_CATALOG_API_VERSION}/model_catalog/models\\?(?!.*sourceLabel=)`,
      ),
    },
    mockModArchResponse(allModelsResponse),
  ).as('getModelsFiltered');
};

/**
 * Intercepts the filter_options endpoint
 */
export const interceptFilterOptions = (): void => {
  cy.interceptApi(
    `GET /api/:apiVersion/model_catalog/models/filter_options`,
    {
      path: { apiVersion: MODEL_CATALOG_API_VERSION },
      query: { namespace: 'kubeflow' },
    },
    mockCatalogFilterOptionsList(),
  );
};

/**
 * Intercepts a single model detail request
 */
export const interceptSingleModel = (sourceId: string, model: CatalogModel): void => {
  cy.interceptApi(
    `GET /api/:apiVersion/model_catalog/sources/:sourceId/models/:modelName`,
    {
      path: {
        apiVersion: MODEL_CATALOG_API_VERSION,
        sourceId,
        modelName: model.name.replace('/', '%2F'),
      },
    },
    model,
  );
};

/**
 * Intercepts single model using regex to match any source/model combination
 */
export const interceptSingleModelRegex = (model: CatalogModel): void => {
  cy.intercept(
    {
      method: 'GET',
      url: new RegExp(
        `/model-registry/api/${MODEL_CATALOG_API_VERSION}/model_catalog/sources/[^/]+/models/[^/]+$`,
      ),
    },
    mockModArchResponse(model),
  ).as('getSingleModel');
};

/**
 * Intercepts performance artifacts for a specific model using regex
 */
export const interceptPerformanceArtifacts = (modelName = 'validated-model'): void => {
  const performanceArtifactsResponse = {
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
          ttft_p90: {
            metadataType: ModelRegistryMetadataType.DOUBLE,
            double_value: 82.34,
          },
          use_case: {
            metadataType: ModelRegistryMetadataType.STRING,
            string_value: 'chatbot',
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
          ttft_p90: {
            metadataType: ModelRegistryMetadataType.DOUBLE,
            double_value: 58.45,
          },
          use_case: {
            metadataType: ModelRegistryMetadataType.STRING,
            string_value: 'chatbot',
          },
        },
      }),
    ],
    pageSize: 10,
    size: 3,
    nextPageToken: '',
  };

  // Use regex to match any source's model performance artifacts requests
  cy.intercept(
    {
      method: 'GET',
      url: new RegExp(
        `/model-registry/api/${MODEL_CATALOG_API_VERSION}/model_catalog/sources/.*/performance_artifacts/${modelName}`,
      ),
    },
    mockModArchResponse(performanceArtifactsResponse),
  ).as('getCatalogSourceModelArtifacts');
};

/**
 * Intercepts performance artifacts using the standard list mock
 * @param overrideArtifacts - Optional override for the artifacts response (e.g., empty list)
 */
export const interceptPerformanceArtifactsList = (overrideArtifacts?: {
  items: unknown[];
  size: number;
  pageSize: number;
  nextPageToken: string;
}): void => {
  const testArtifacts = overrideArtifacts ?? mockCatalogPerformanceMetricsArtifactList({});

  cy.intercept(
    {
      method: 'GET',
      url: new RegExp(
        `/model-registry/api/${MODEL_CATALOG_API_VERSION}/model_catalog/sources/.*/performance_artifacts/.*`,
      ),
    },
    mockModArchResponse(testArtifacts),
  ).as('getCatalogSourceModelArtifacts');
};

/**
 * Intercepts the /artifacts/ endpoint (used to determine if tabs should show)
 * @param overrideArtifacts - Optional override for the artifacts response (e.g., empty list)
 */
export const interceptArtifactsList = (overrideArtifacts?: {
  items: unknown[];
  size: number;
  pageSize: number;
  nextPageToken: string;
}): void => {
  const testArtifacts = overrideArtifacts ?? mockCatalogPerformanceMetricsArtifactList({});

  cy.intercept(
    {
      method: 'GET',
      url: new RegExp(
        `/model-registry/api/${MODEL_CATALOG_API_VERSION}/model_catalog/sources/.*/artifacts/.*`,
      ),
    },
    mockModArchResponse(testArtifacts),
  ).as('getArtifacts');
};

/**
 * Intercepts model registry endpoint
 */
export const interceptModelRegistry = (): void => {
  cy.interceptApi(
    `GET /api/:apiVersion/model_registry`,
    {
      path: { apiVersion: 'v1' },
    },
    [mockModelRegistry({})],
  );
};

/**
 * Sets up all common model catalog intercepts with sensible defaults.
 * This is the main function to use for setting up test intercepts.
 */
export const setupModelCatalogIntercepts = (options: ModelCatalogInterceptOptions = {}): void => {
  const {
    sources = defaultSources(),
    modelsPerCategory = 4,
    useValidatedModel = false,
    includePerformanceArtifacts = false,
    includeFilterOptions = true,
    includeModelRegistry = false,
    customValidatedModel,
    customNonValidatedModel,
  } = options;

  // Always intercept sources
  interceptSources(sources);

  // Intercept models by label
  interceptModelsByLabel(sources, modelsPerCategory, useValidatedModel);

  // Intercept all models (for GalleryView without sourceLabel)
  interceptAllModels(modelsPerCategory, useValidatedModel);

  // Intercept single model detail
  const testModel = useValidatedModel
    ? (customValidatedModel ?? mockValidatedModel)
    : (customNonValidatedModel ?? mockNonValidatedModel);
  interceptSingleModel('source-2', testModel);
  interceptSingleModelRegex(testModel);

  // Optional intercepts
  if (includeFilterOptions) {
    interceptFilterOptions();
  }

  if (includePerformanceArtifacts) {
    interceptPerformanceArtifacts();
  }

  if (includeModelRegistry) {
    interceptModelRegistry();
  }
};

/**
 * Sets up intercepts specifically for tests that need validated models with performance data
 */
export const setupValidatedModelIntercepts = (
  options: Omit<
    ModelCatalogInterceptOptions,
    'useValidatedModel' | 'includePerformanceArtifacts'
  > = {},
): void => {
  setupModelCatalogIntercepts({
    ...options,
    useValidatedModel: true,
    includePerformanceArtifacts: true,
  });
};

/**
 * Sets up intercepts for model details page tests
 */
export const setupModelDetailsIntercepts = (options: ModelCatalogInterceptOptions = {}): void => {
  setupModelCatalogIntercepts({
    ...options,
    includeModelRegistry: true,
    includePerformanceArtifacts: true,
  });
};
