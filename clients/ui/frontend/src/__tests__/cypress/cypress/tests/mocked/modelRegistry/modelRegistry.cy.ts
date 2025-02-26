/* eslint-disable camelcase */
import { mockModelRegistry } from '~/__mocks__/mockModelRegistry';
import { mockModelVersion } from '~/__mocks__/mockModelVersion';
import { mockModelVersionList } from '~/__mocks__/mockModelVersionList';
import { mockRegisteredModel } from '~/__mocks__/mockRegisteredModel';
import { mockRegisteredModelList } from '~/__mocks__/mockRegisteredModelsList';
import { labelModal, modelRegistry } from '~/__tests__/cypress/cypress/pages/modelRegistry';
import {
  ModelRegistryMetadataType,
  type ModelRegistry,
  type ModelVersion,
  type RegisteredModel,
} from '~/app/types';
import { be } from '~/__tests__/cypress/cypress/utils/should';
import { MODEL_REGISTRY_API_VERSION } from '~/__tests__/cypress/cypress/support/commands/api';

type HandlersProps = {
  modelRegistries?: ModelRegistry[];
  registeredModels?: RegisteredModel[];
  modelVersions?: ModelVersion[];
};

const initIntercepts = ({
  modelRegistries = [
    mockModelRegistry({
      name: 'modelregistry-sample',
    }),
    mockModelRegistry({
      name: 'modelregistry-sample-2',
      description: '',
      displayName: 'modelregistry-sample-2',
    }),
  ],
  registeredModels = [
    mockRegisteredModel({
      name: 'Fraud detection model',
      description:
        'A machine learning model trained to detect fraudulent transactions in financial data',
      customProperties: {
        'Financial data': {
          metadataType: ModelRegistryMetadataType.STRING,
          string_value: '',
        },
        'Fraud detection': {
          metadataType: ModelRegistryMetadataType.STRING,
          string_value: '',
        },
        'Test label': {
          metadataType: ModelRegistryMetadataType.STRING,
          string_value: '',
        },
        'Machine learning': {
          metadataType: ModelRegistryMetadataType.STRING,
          string_value: '',
        },
        'Next data to be overflow': {
          metadataType: ModelRegistryMetadataType.STRING,
          string_value: '',
        },
      },
    }),
    mockRegisteredModel({
      name: 'Label modal',
      description:
        'A machine learning model trained to detect fraudulent transactions in financial data',
      customProperties: {
        'Testing label': {
          metadataType: ModelRegistryMetadataType.STRING,
          string_value: '',
        },
        'Financial data': {
          metadataType: ModelRegistryMetadataType.STRING,
          string_value: '',
        },
        'Fraud detection': {
          metadataType: ModelRegistryMetadataType.STRING,
          string_value: '',
        },
        'Long label data to be truncated abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc':
          {
            metadataType: ModelRegistryMetadataType.STRING,
            string_value: '',
          },
        'Machine learning': {
          metadataType: ModelRegistryMetadataType.STRING,
          string_value: '',
        },
        'Next data to be overflow': {
          metadataType: ModelRegistryMetadataType.STRING,
          string_value: '',
        },
        'Label x': {
          metadataType: ModelRegistryMetadataType.STRING,
          string_value: '',
        },
        'Label y': {
          metadataType: ModelRegistryMetadataType.STRING,
          string_value: '',
        },
        'Label z': {
          metadataType: ModelRegistryMetadataType.STRING,
          string_value: '',
        },
      },
    }),
  ],
  modelVersions = [
    mockModelVersion({ author: 'Author 1' }),
    mockModelVersion({ name: 'model version' }),
  ],
}: HandlersProps) => {
  cy.interceptApi(
    `GET /api/:apiVersion/model_registry`,
    {
      path: { apiVersion: MODEL_REGISTRY_API_VERSION },
    },
    modelRegistries,
  );

  cy.interceptApi(
    `GET /api/:apiVersion/model_registry/:modelRegistryName/model_versions`,
    {
      path: { modelRegistryName: 'modelregistry-sample', apiVersion: MODEL_REGISTRY_API_VERSION },
    },
    mockModelVersionList({ items: modelVersions }),
  );

  cy.interceptApi(
    `GET /api/:apiVersion/model_registry/:modelRegistryName/registered_models`,
    {
      path: { modelRegistryName: 'modelregistry-sample', apiVersion: MODEL_REGISTRY_API_VERSION },
    },
    mockRegisteredModelList({ items: registeredModels }),
  );

  cy.interceptApi(
    `GET /api/:apiVersion/model_registry/:modelRegistryName/registered_models/:registeredModelId/versions`,
    {
      path: {
        modelRegistryName: 'modelregistry-sample',
        apiVersion: MODEL_REGISTRY_API_VERSION,
        registeredModelId: 1,
      },
    },
    mockModelVersionList({ items: modelVersions }),
  );
};

describe('Model Registry core', () => {
  it('Model Registry Enabled in the cluster', () => {
    initIntercepts({
      registeredModels: [],
    });

    modelRegistry.visit();
    modelRegistry.navigate();

    modelRegistry.tabEnabled();
  });

  it('Renders empty state with no model registries', () => {
    initIntercepts({
      modelRegistries: [],
    });

    modelRegistry.visit();
    modelRegistry.navigate();
    modelRegistry.findModelRegistryEmptyState().should('exist');
  });
  it('No registered models in the selected Model Registry', () => {
    initIntercepts({
      registeredModels: [],
    });

    modelRegistry.visit();
    modelRegistry.navigate();
    modelRegistry.shouldModelRegistrySelectorExist();
    modelRegistry.shouldregisteredModelsEmpty();

    modelRegistry.findViewDetailsButton().click();
    modelRegistry.findDetailsPopover().should('exist');
    modelRegistry.findDetailsPopover().findByText('Model registry description').should('exist');

    // Model registry with no description
    modelRegistry.findModelRegistry().findSelectOption('modelregistry-sample-2').click();
    modelRegistry.findViewDetailsButton().click();
    modelRegistry.findDetailsPopover().should('exist');
    modelRegistry.findDetailsPopover().findByText('No description').should('exist');

    //  Model registry help content
    modelRegistry.findHelpContentButton().should('exist');
    // modelRegistry.findHelpContentButton().click(); TODO: Get conditional render depnding on the style
    // modelRegistry.findHelpContentPopover().should('exist');
    // modelRegistry
    //   .findHelpContentPopover()
    //   .findByText(
    //     'To request access to a new or existing model registry, contact your administrator.',
    //   )
    //   .should('exist');
  });

  describe('Registered model table', () => {
    beforeEach(() => {
      initIntercepts({});
      modelRegistry.visit();
    });

    it('Renders row contents', () => {
      const registeredModelRow = modelRegistry.getRow('Fraud detection model');
      registeredModelRow.findName().contains('Fraud detection model');
      registeredModelRow
        .findDescription()
        .contains(
          'A machine learning model trained to detect fraudulent transactions in financial data',
        );
      registeredModelRow.findOwner().contains('Author 1');

      // Label popover
      registeredModelRow.findLabelPopoverText().contains('2 more');
      registeredModelRow.findLabelPopoverText().click();
      registeredModelRow.shouldContainsPopoverLabels([
        'Machine learning',
        'Next data to be overflow',
      ]);
    });

    it('Renders labels in modal', () => {
      const registeredModelRow2 = modelRegistry.getRow('Label modal');
      registeredModelRow2.findLabelModalText().contains('6 more');
      registeredModelRow2.findLabelModalText().click();
      labelModal.shouldContainsModalLabels([
        'Testing label',
        'Financial',
        'Financial data',
        'Fraud detection',
        'Machine learning',
        'Next data to be overflow',
        'Label x',
        'Label y',
        'Label z',
      ]);
      labelModal.findModalSearchInput().type('Financial');
      labelModal.shouldContainsModalLabels(['Financial', 'Financial data']);
      labelModal.findCloseModal().click();
    });

    it('Sort by Model name', () => {
      modelRegistry.findRegisteredModelTableHeaderButton('Model name').click();
      modelRegistry.findRegisteredModelTableHeaderButton('Model name').should(be.sortAscending);
      modelRegistry.findRegisteredModelTableHeaderButton('Model name').click();
      modelRegistry.findRegisteredModelTableHeaderButton('Model name').should(be.sortDescending);
    });

    it('Sort by Last modified', () => {
      modelRegistry.findRegisteredModelTableHeaderButton('Last modified').should(be.sortAscending);
      modelRegistry.findRegisteredModelTableHeaderButton('Last modified').click();
      modelRegistry.findRegisteredModelTableHeaderButton('Last modified').should(be.sortDescending);
    });

    it('Filter by keyword', () => {
      modelRegistry.findTableSearch().type('Fraud detection model');
      modelRegistry.findTableRows().should('have.length', 1);
      modelRegistry.findTableRows().contains('Fraud detection model');
    });
  });
});

describe('Register Model button', () => {
  it('Navigates to register page from empty state', () => {
    initIntercepts({ registeredModels: [] });
    modelRegistry.visit();
    modelRegistry.findRegisterModelButton().click();
    cy.findByTestId('app-page-title').should('exist');
    cy.findByTestId('app-page-title').contains('Register model');
    cy.findByText('Model registry - modelregistry-sample').should('exist');
  });

  it('Navigates to register page from table toolbar', () => {
    initIntercepts({ registeredModels: [] });
    modelRegistry.visit();
    modelRegistry.findRegisterModelButton().click();
    cy.findByTestId('app-page-title').should('exist');
    cy.findByTestId('app-page-title').contains('Register model');
    cy.findByText('Model registry - modelregistry-sample').should('exist');
  });
});
