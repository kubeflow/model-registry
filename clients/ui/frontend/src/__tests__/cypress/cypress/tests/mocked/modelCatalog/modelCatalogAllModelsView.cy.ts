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
  modelsPerCategory?: number;
  isEmpty?: boolean;
};

const initIntercepts = ({
  sources = [
    mockCatalogSource({ id: 'huggingface', name: 'Hugging Face', labels: ['Hugging Face'] }),
    mockCatalogSource({ id: 'openvino', name: 'OpenVINO', labels: ['OpenVINO'] }),
    mockCatalogSource({ id: 'community', name: 'Community', labels: ['Community'] }),
  ],
  modelsPerCategory = 4,
  isEmpty = false,
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
    `GET /api/:apiVersion/model_catalog/models/filter_options`,
    {
      path: { apiVersion: MODEL_CATALOG_API_VERSION },
      query: { namespace: 'kubeflow' },
    },
    mockCatalogFilterOptionsList(),
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
          items: isEmpty
            ? []
            : Array.from({ length: modelsPerCategory }, (_, i) =>
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
};

describe('Model Catalog All Models View', () => {
  beforeEach(() => {
    initIntercepts({});
    modelCatalog.visit();
  });

  describe('Category Sections', () => {
    it('should display all category sections', () => {
      modelCatalog.findAllModelsToggle().should('be.visible');
      modelCatalog.findCategoryToggle('label-Hugging Face').should('be.visible');
      modelCatalog.findCategoryToggle('label-OpenVINO').should('be.visible');
      modelCatalog.findCategoryToggle('label-Community').should('be.visible');
      modelCatalog.findCategoryToggle('no-labels').should('be.visible');
    });

    it('should show category titles', () => {
      modelCatalog.findCategoryTitle('OpenVINO').should('contain.text', 'OpenVINO models');
      cy.findByTestId('title Hugging Face').should('contain.text', 'Hugging Face models');
      modelCatalog.findCategoryTitle('Community').should('contain.text', 'Community models');
      modelCatalog.findCategoryTitle('null').should('contain.text', 'Community and custom models');
    });
  });

  describe('Show More Functionality', () => {
    it('should display show more button when category has 4 or more models', () => {
      modelCatalog.findShowMoreModelsLink('hugging-face').should('be.visible');
      modelCatalog.findShowMoreModelsLink('hugging-face').click();
      modelCatalog.findAllModelsToggle().click();
      modelCatalog.findShowMoreModelsLink('openvino').should('be.visible');
      modelCatalog.findShowMoreModelsLink('openvino').click();
      modelCatalog.findAllModelsToggle().click();
      modelCatalog.findShowMoreModelsLink('community').should('be.visible');
      modelCatalog.findAllModelsToggle().click();
      modelCatalog.findShowMoreModelsLink('community').click();
    });
  });

  describe('Error Handling', () => {
    it('should display error message when category fails to load', () => {
      modelCatalog.findErrorState('null').should('be.visible');
      modelCatalog
        .findErrorState('null')
        .should('contain.text', 'Failed to load Community and custom models');
    });
  });

  describe('Empty States', () => {
    it('should show empty state when category has no models', () => {
      initIntercepts({ isEmpty: true });
      modelCatalog.visit();

      modelCatalog.findEmptyState('OpenVINO').should('be.visible');
      modelCatalog
        .findEmptyState('OpenVINO')
        .should('contain.text', 'No result foundAdjust your filters and try again.');
    });
  });
});
