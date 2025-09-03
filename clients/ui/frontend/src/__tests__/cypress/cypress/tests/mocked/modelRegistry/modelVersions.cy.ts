/* eslint-disable camelcase */
import { mockModelVersionList } from '~/__mocks__/mockModelVersionList';
import { mockRegisteredModelList } from '~/__mocks__/mockRegisteredModelsList';
import { labelModal, modelRegistry } from '~/__tests__/cypress/cypress/pages/modelRegistry';
import { be } from '~/__tests__/cypress/cypress/utils/should';
import { mockRegisteredModel } from '~/__mocks__/mockRegisteredModel';
import { ModelRegistryMetadataType, type ModelRegistry, type ModelVersion } from '~/app/types';
import { verifyRelativeURL } from '~/__tests__/cypress/cypress/utils/url';
import { mockModelRegistry } from '~/__mocks__/mockModelRegistry';
import { mockModelVersion } from '~/__mocks__/mockModelVersion';
import { MODEL_REGISTRY_API_VERSION } from '~/__tests__/cypress/cypress/support/commands/api';

type HandlersProps = {
  registeredModelsSize?: number;
  modelVersions?: ModelVersion[];
  modelRegistries?: ModelRegistry[];
};

const initIntercepts = ({
  registeredModelsSize = 4,
  modelRegistries = [
    mockModelRegistry({
      name: 'modelregistry-sample',
      description: 'New model registry',
      displayName: 'Model Registry Sample',
    }),
    mockModelRegistry({
      name: 'modelregistry-sample-2',
      description: 'New model registry 2',
      displayName: 'Model Registry Sample 2',
    }),
  ],
  modelVersions = [
    mockModelVersion({
      author: 'Author 1',
      id: '1',
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
        'Test label x': {
          metadataType: ModelRegistryMetadataType.STRING,
          string_value: '',
        },
        'Test label y': {
          metadataType: ModelRegistryMetadataType.STRING,
          string_value: '',
        },
        'Test label z': {
          metadataType: ModelRegistryMetadataType.STRING,
          string_value: '',
        },
      },
    }),
    mockModelVersion({ id: '2', name: 'model version' }),
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
    mockRegisteredModelList({ size: registeredModelsSize }),
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

  cy.interceptApi(
    `GET /api/:apiVersion/model_registry/:modelRegistryName/registered_models/:registeredModelId`,
    {
      path: {
        modelRegistryName: 'modelregistry-sample',
        apiVersion: MODEL_REGISTRY_API_VERSION,
        registeredModelId: 1,
      },
    },
    mockRegisteredModel({}),
  );

  cy.interceptApi(
    `GET /api/:apiVersion/model_registry/:modelRegistryName/model_versions/:modelVersionId`,
    {
      path: {
        modelRegistryName: 'modelregistry-sample',
        apiVersion: MODEL_REGISTRY_API_VERSION,
        modelVersionId: 1,
      },
    },
    mockModelVersion({ id: '1', name: 'model version' }),
  );
};

describe('Model Versions', () => {
  it('No model versions in the selected registered model', () => {
    initIntercepts({
      modelVersions: [],
    });

    modelRegistry.visit();
    const registeredModelRow = modelRegistry.getRow('Fraud detection model');
    registeredModelRow.findName().contains('Fraud detection model').click();
    verifyRelativeURL(`/model-registry/modelregistry-sample/registeredModels/1/overview`);

    modelRegistry.shouldModelVersionsEmpty();
  });

  it('Model versions table browser back button should lead to Registered models table', () => {
    initIntercepts({
      modelVersions: [],
    });

    modelRegistry.visit();
    const registeredModelRow = modelRegistry.getRow('Fraud detection model');
    registeredModelRow.findName().contains('Fraud detection model').click();
    verifyRelativeURL(`/model-registry/modelregistry-sample/registeredModels/1/overview`);
    cy.go('back');
    verifyRelativeURL(`/model-registry/modelregistry-sample`);
    registeredModelRow.findName().contains('Fraud detection model').should('exist');
  });

  it('Model versions table', () => {
    // TODO: [Testing] Uncomment when we fix finding listbox items

    initIntercepts({
      modelRegistries: [
        // mockModelRegistry({ name: 'modelRegistry-1', displayName: 'modelRegistry-1' }),
        mockModelRegistry({}),
      ],
    });

    modelRegistry.visit();
    //modelRegistry.findModelRegistry().findSelectOption('Model Registry Sample').click();
    //cy.reload();
    const registeredModelRow = modelRegistry.getRow('Fraud detection model');
    registeredModelRow.findName().contains('Fraud detection model').click();
    verifyRelativeURL(`/model-registry/modelregistry-sample/registeredModels/1/overview`);
    modelRegistry.findModelBreadcrumbItem().contains('test');

    // Navigate to versions tab
    cy.findByTestId('model-versions-tab').click();
    //modelRegistry.findModelVersionsTableKebab().findDropdownItem('View archived versions');
    //modelRegistry.findModelVersionsHeaderAction().findDropdownItem('Archive model');
    modelRegistry.findModelVersionsTable().should('be.visible');
    modelRegistry.findModelVersionsTableRows().should('have.length', 2);

    // Label modal
    const modelVersionRow = modelRegistry.getModelVersionRow('new model version');

    modelVersionRow.findLabelModalText().contains('5 more');
    modelVersionRow.findLabelModalText().click();
    labelModal.shouldContainsModalLabels([
      'Financial',
      'Financial data',
      'Fraud detection',
      'Test label',
      'Machine learning',
      'Next data to be overflow',
      'Test label x',
      'Test label y',
      'Test label y',
    ]);
    labelModal.findModalSearchInput().type('Financial');
    labelModal.shouldContainsModalLabels(['Financial', 'Financial data']);
    labelModal.findCloseModal().click();

    // sort by model version name
    modelRegistry.findModelVersionsTableHeaderButton('Version name').click();
    modelRegistry.findModelVersionsTableHeaderButton('Version name').should(be.sortAscending);
    modelRegistry.findModelVersionsTableHeaderButton('Version name').click();
    modelRegistry.findModelVersionsTableHeaderButton('Version name').should(be.sortDescending);

    // sort by Last modified
    modelRegistry.findModelVersionsTableHeaderButton('Last modified').click();
    modelRegistry.findModelVersionsTableHeaderButton('Last modified').should(be.sortAscending);
    modelRegistry.findModelVersionsTableHeaderButton('Last modified').click();
    modelRegistry.findModelVersionsTableHeaderButton('Last modified').should(be.sortDescending);

    // sort by model version author
    modelRegistry.findModelVersionsTableHeaderButton('Author').click();
    modelRegistry.findModelVersionsTableHeaderButton('Author').should(be.sortAscending);
    modelRegistry.findModelVersionsTableHeaderButton('Author').click();
    modelRegistry.findModelVersionsTableHeaderButton('Author').should(be.sortDescending);

    // filtering by label then both
    modelRegistry.findTableSearch().type('Financial');
    modelRegistry.findModelVersionsTableRows().should('have.length', 1);
    modelRegistry.findModelVersionsTableRows().contains('new model version');
    modelRegistry.findModelVersionsTableFilterOption('Author').click();
    modelRegistry.findTableSearch().type('Author 1');
    modelRegistry.findModelVersionsTableRows().should('have.length', 1);
    modelRegistry.findModelVersionsTableRows().contains('new model version');
    modelRegistry.findTableSearch().type('2');
    modelRegistry.findModelVersionsTableRows().should('have.length', 0);
    modelRegistry.findTableSearch().focused().clear();
    modelRegistry.findModelVersionsTableFilterOption('Keyword').click();
    modelRegistry.findTableSearch().click();
    modelRegistry.findTableSearch().focused().clear();

    // filtering by model version author then both
    modelRegistry.findModelVersionsTableFilterOption('Author').click();
    modelRegistry.findTableSearch().type('Test author');
    modelRegistry.findModelVersionsTableRows().should('have.length', 1);
    modelRegistry.findModelVersionsTableRows().contains('Test author');
    modelRegistry.findModelVersionsTableFilterOption('Keyword').click();
    modelRegistry.findTableSearch().type('model version');
    modelRegistry.findModelVersionsTableRows().should('have.length', 1);
    modelRegistry.findModelVersionsTableRows().contains('model version');
    modelRegistry.findTableSearch().type('2');

    // searching with no matches shows no results screen
    modelRegistry.findModelVersionsTableRows().should('have.length', 0);
  });

  it('Model version details back button should lead to versions table', () => {
    initIntercepts({});

    modelRegistry.visit();
    const registeredModelRow = modelRegistry.getRow('Fraud detection model');
    registeredModelRow.findName().contains('Fraud detection model').click();
    // Navigate to versions tab
    cy.findByTestId('model-versions-tab').click();
    const modelVersionRow = modelRegistry.getModelVersionRow('model version');
    modelVersionRow.findModelVersionName().contains('model version').click();
    verifyRelativeURL('/model-registry/modelregistry-sample/registeredModels/1/versions/1/details');
    cy.findByTestId('app-page-title').should('contain.text', 'test');
    cy.findByTestId('breadcrumb-version-name').should('have.text', 'model version');
    cy.go('back');
    verifyRelativeURL('/model-registry/modelregistry-sample/registeredModels/1/versions');
  });
});
