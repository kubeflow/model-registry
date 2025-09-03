/* eslint-disable camelcase */
import { mockModArchResponse } from 'mod-arch-core';
import { mockRegisteredModelList } from '~/__mocks__/mockRegisteredModelsList';
import { mockModelVersion } from '~/__mocks__/mockModelVersion';
import { mockRegisteredModel } from '~/__mocks__/mockRegisteredModel';
import { verifyRelativeURL } from '~/__tests__/cypress/cypress/utils/url';
import { labelModal, modelRegistry } from '~/__tests__/cypress/cypress/pages/modelRegistry';
import { mockModelVersionList } from '~/__mocks__/mockModelVersionList';
import { be } from '~/__tests__/cypress/cypress/utils/should';
import type { ModelRegistry, ModelVersion, RegisteredModel } from '~/app/types';
import { ModelRegistryMetadataType, ModelState } from '~/app/types';
import { mockModelRegistry } from '~/__mocks__/mockModelRegistry';
import { MODEL_REGISTRY_API_VERSION } from '~/__tests__/cypress/cypress/support/commands/api';
import {
  archiveModelModal,
  registeredModelArchive,
  restoreModelModal,
} from '~/__tests__/cypress/cypress/pages/modelRegistryView/registeredModelArchive';
import { ToastNotification } from '~/__tests__/cypress/cypress/pages/components/Notification';

type HandlersProps = {
  registeredModels?: RegisteredModel[];
  modelVersions?: ModelVersion[];
  modelRegistries?: ModelRegistry[];
};

const initIntercepts = ({
  registeredModels = [
    mockRegisteredModel({
      name: 'model 1',
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
      state: ModelState.ARCHIVED,
    }),
    mockRegisteredModel({
      id: '2',
      owner: 'Author 2',
      name: 'model 2',
      state: ModelState.ARCHIVED,
    }),
    mockRegisteredModel({ id: '3', name: 'model 3' }),
    mockRegisteredModel({ id: '4', name: 'model 4' }),
  ],
  modelVersions = [
    mockModelVersion({ author: 'Author 1', registeredModelId: '2' }),
    mockModelVersion({ name: 'model version' }),
  ],
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
    `GET /api/:apiVersion/model_registry/:modelRegistryName/model_versions/:modelVersionId`,
    {
      path: {
        modelRegistryName: 'modelregistry-sample',
        apiVersion: MODEL_REGISTRY_API_VERSION,
        modelVersionId: 1,
      },
    },
    mockModelVersion({ id: '1', name: 'Version 2' }),
  );

  cy.interceptApi(
    `GET /api/:apiVersion/model_registry/:modelRegistryName/registered_models/:registeredModelId/versions`,
    {
      path: {
        modelRegistryName: 'modelregistry-sample',
        apiVersion: MODEL_REGISTRY_API_VERSION,
        registeredModelId: 2,
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
        registeredModelId: 2,
      },
    },
    mockRegisteredModel({ id: '2', name: 'model 2', state: ModelState.ARCHIVED }),
  );

  cy.interceptApi(
    `GET /api/:apiVersion/model_registry/:modelRegistryName/registered_models/:registeredModelId`,
    {
      path: {
        modelRegistryName: 'modelregistry-sample',
        apiVersion: MODEL_REGISTRY_API_VERSION,
        registeredModelId: 3,
      },
    },
    mockRegisteredModel({ id: '3', name: 'model 3' }),
  );
};

describe('Model archive list', () => {
  it('No archive models in the selected model registry', () => {
    initIntercepts({
      registeredModels: [],
    });
    registeredModelArchive.visit();
    verifyRelativeURL('/model-registry/modelregistry-sample/registeredModels/archive');
    registeredModelArchive.shouldArchiveVersionsEmpty();
  });

  it('Archived model details browser back button should lead to archived models table', () => {
    initIntercepts({});
    registeredModelArchive.visit();
    verifyRelativeURL('/model-registry/modelregistry-sample/registeredModels/archive');
    registeredModelArchive.findArchiveModelBreadcrumbItem().contains('Archived models');
    const archiveModelRow = registeredModelArchive.getRow('model 2');
    archiveModelRow.findName().contains('model 2').click();
    verifyRelativeURL('/model-registry/modelregistry-sample/registeredModels/archive/2/overview');
    cy.findByTestId('app-page-title').should('have.text', 'model 2Archived');
    cy.go('back');
    verifyRelativeURL('/model-registry/modelregistry-sample/registeredModels/archive');
    registeredModelArchive.findArchiveModelTable().should('be.visible');
  });

  it('Archived model with no versions', () => {
    initIntercepts({ modelVersions: [] });
    registeredModelArchive.visit();
    verifyRelativeURL('/model-registry/modelregistry-sample/registeredModels/archive');
    registeredModelArchive.findArchiveModelBreadcrumbItem().contains('Archived models');
    const archiveModelRow = registeredModelArchive.getRow('model 2');
    archiveModelRow.findName().contains('model 2').click();
    modelRegistry.shouldArchiveModelVersionsEmpty();
  });

  it('Archived model flow', () => {
    initIntercepts({});
    registeredModelArchive.visitArchiveModelVersionList();
    verifyRelativeURL('/model-registry/modelregistry-sample/registeredModels/archive/2/versions');

    modelRegistry.findModelVersionsTable().should('be.visible');
    modelRegistry.findModelVersionsTableRows().should('have.length', 2);
    const version = modelRegistry.getModelVersionRow('model version');
    version.findModelVersionName().contains('model version').click();
    verifyRelativeURL(
      '/model-registry/modelregistry-sample/registeredModels/archive/2/versions/1/details',
    );
    cy.go('back');
    verifyRelativeURL('/model-registry/modelregistry-sample/registeredModels/archive/2/versions');
  });

  it('Archive models list', () => {
    initIntercepts({});
    registeredModelArchive.visit();
    verifyRelativeURL('/model-registry/modelregistry-sample/registeredModels/archive');

    //breadcrumb
    registeredModelArchive.findArchiveModelBreadcrumbItem().contains('Archived models');

    // name, last modified, owner, labels modal
    registeredModelArchive.findArchiveModelTable().should('be.visible');
    registeredModelArchive.findTableSearch().type('model 1');
    registeredModelArchive.findArchiveModelsTableRows().should('have.length', 1);
    registeredModelArchive
      .findArchiveModelsTableToolbar()
      .findByRole('button', { name: 'Clear all filters' })
      .click();
    registeredModelArchive.findArchiveModelsTableRows().should('have.length', 2);

    const archiveModelRow = registeredModelArchive.getRow('model 1');

    archiveModelRow.findLabelModalText().contains('5 more');
    archiveModelRow.findLabelModalText().click();
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
    labelModal.findCloseModal().click();

    // sort by Last modified
    registeredModelArchive.findRegisteredModelsArchiveTableHeaderButton('Last modified').click();
    registeredModelArchive
      .findRegisteredModelsArchiveTableHeaderButton('Last modified')
      .should(be.sortAscending);
    registeredModelArchive.findRegisteredModelsArchiveTableHeaderButton('Last modified').click();
    registeredModelArchive
      .findRegisteredModelsArchiveTableHeaderButton('Last modified')
      .should(be.sortDescending);

    // sort by Model name
    registeredModelArchive.findRegisteredModelsArchiveTableHeaderButton('Model name').click();
    registeredModelArchive
      .findRegisteredModelsArchiveTableHeaderButton('Model name')
      .should(be.sortAscending);
    registeredModelArchive.findRegisteredModelsArchiveTableHeaderButton('Model name').click();
    registeredModelArchive
      .findRegisteredModelsArchiveTableHeaderButton('Model name')
      .should(be.sortDescending);
  });

  it('Filter by keyword then both', () => {
    initIntercepts({});
    registeredModelArchive.visit();
    registeredModelArchive.findTableSearch().type('model 1');
    registeredModelArchive.findArchiveModelsTableRows().should('have.length', 1);
    registeredModelArchive.findFilterDropdownItem('Owner').click();
    registeredModelArchive.findTableSearch().type('Author 1');
    registeredModelArchive.findArchiveModelsTableRows().should('have.length', 1);
    registeredModelArchive.findArchiveModelsTableRows().contains('model 1');
    registeredModelArchive.findTableSearch().type('2');
    registeredModelArchive.findArchiveModelsTableRows().should('have.length', 0);
  });

  it('Filter by owner then both', () => {
    initIntercepts({});
    registeredModelArchive.visit();
    registeredModelArchive.findFilterDropdownItem('Owner').click();
    registeredModelArchive.findTableSearch().type('Author 2');
    registeredModelArchive.findArchiveModelsTableRows().should('have.length', 1);
    registeredModelArchive.findFilterDropdownItem('Keyword').click();
    registeredModelArchive.findTableSearch().type('model 2');
    registeredModelArchive.findArchiveModelsTableRows().should('have.length', 1);
    registeredModelArchive.findTableSearch().type('.');
    registeredModelArchive.findArchiveModelsTableRows().should('have.length', 0);
  });

  it('latest version column', () => {
    initIntercepts({});
    registeredModelArchive.visit();
    const archiveModelRow = registeredModelArchive.getRow('model 2');
    archiveModelRow.findLatestVersion().contains('new model version');
    archiveModelRow.findLatestVersion().click();
    verifyRelativeURL(
      `/model-registry/modelregistry-sample/registeredModels/archive/2/versions/1/details`,
    );
  });

  it('Opens the detail page when we select "Overview" from action menu and verison list when we select "Versions', () => {
    initIntercepts({});
    registeredModelArchive.visit();
    const archiveModelRow = registeredModelArchive.getRow('model 2');
    archiveModelRow.findKebabAction('Overview').click();
    cy.location('pathname').should(
      'be.equals',
      '/model-registry/modelregistry-sample/registeredModels/archive/2/overview',
    );
    cy.go('back');
    archiveModelRow.findKebabAction('Versions').click();
    verifyRelativeURL(`/model-registry/modelregistry-sample/registeredModels/archive/2/versions`);
  });
});

describe('Restoring archive model', () => {
  it('Restore from archive models table', () => {
    cy.interceptApi(
      'PATCH /api/:apiVersion/model_registry/:modelRegistryName/registered_models/:registeredModelId',
      {
        path: {
          modelRegistryName: 'modelregistry-sample',
          apiVersion: MODEL_REGISTRY_API_VERSION,
          registeredModelId: 2,
        },
      },
      mockRegisteredModel({ id: '2', name: 'model 2', state: ModelState.LIVE }),
    ).as('modelRestored');

    initIntercepts({});
    registeredModelArchive.visit();

    const archiveModelRow = registeredModelArchive.getRow('model 2');
    archiveModelRow.findKebabAction('Restore model').click();

    restoreModelModal.findRestoreButton().click();

    const notification = new ToastNotification(`model 2 and all its versions restored.`);
    notification.find();

    cy.wait('@modelRestored').then((interception) => {
      expect(interception.request.body).to.eql(mockModArchResponse({ state: 'LIVE' }));
    });
  });

  it('Restore from archive model details', () => {
    cy.interceptApi(
      'PATCH /api/:apiVersion/model_registry/:modelRegistryName/registered_models/:registeredModelId',
      {
        path: {
          modelRegistryName: 'modelregistry-sample',
          apiVersion: MODEL_REGISTRY_API_VERSION,
          registeredModelId: 2,
        },
      },
      mockRegisteredModel({ id: '2', name: 'model 2', state: ModelState.LIVE }),
    ).as('modelRestored');

    initIntercepts({});
    registeredModelArchive.visitArchiveModelDetail();

    registeredModelArchive.findRestoreButton().click();
    restoreModelModal.findRestoreButton().click();

    const notification = new ToastNotification(`model 2 and all its versions restored.`);
    notification.find();

    cy.wait('@modelRestored').then((interception) => {
      expect(interception.request.body).to.eql(mockModArchResponse({ state: 'LIVE' }));
    });
  });
});

describe('Archiving model', () => {
  it('Archive model from registered models table', () => {
    cy.interceptApi(
      'PATCH /api/:apiVersion/model_registry/:modelRegistryName/registered_models/:registeredModelId',
      {
        path: {
          modelRegistryName: 'modelregistry-sample',
          apiVersion: MODEL_REGISTRY_API_VERSION,
          registeredModelId: 3,
        },
      },
      mockRegisteredModel({ id: '3', name: 'model 3', state: ModelState.ARCHIVED }),
    ).as('modelArchived');

    initIntercepts({});
    registeredModelArchive.visitModelList();

    const modelRow = modelRegistry.getRow('model 3');
    modelRow.findKebabAction('Archive model').click();
    archiveModelModal.findArchiveButton().should('be.disabled');
    archiveModelModal.findModalTextInput().fill('model 3');
    archiveModelModal.findArchiveButton().should('be.enabled').click();

    const notification = new ToastNotification('model 3 and all its versions archived.');
    notification.find();

    cy.wait('@modelArchived').then((interception) => {
      expect(interception.request.body).to.eql(mockModArchResponse({ state: 'ARCHIVED' }));
    });
  });

  it('Archive model from model details', () => {
    cy.interceptApi(
      'PATCH /api/:apiVersion/model_registry/:modelRegistryName/registered_models/:registeredModelId',
      {
        path: {
          modelRegistryName: 'modelregistry-sample',
          apiVersion: MODEL_REGISTRY_API_VERSION,
          registeredModelId: 3,
        },
      },
      mockRegisteredModel({ id: '3', name: 'model 3', state: ModelState.ARCHIVED }),
    ).as('modelArchived');

    initIntercepts({});
    registeredModelArchive.visitModelList();

    const modelRow = modelRegistry.getRow('model 3');
    modelRow.findName().contains('model 3').click();
    registeredModelArchive
      .findModelVersionsDetailsHeaderAction()
      .findDropdownItem('Archive model')
      .click();

    archiveModelModal.findArchiveButton().should('be.disabled');
    archiveModelModal.findModalTextInput().fill('model 3');
    archiveModelModal.findArchiveButton().should('be.enabled').click();

    const notification = new ToastNotification('model 3 and all its versions archived.');
    notification.find();

    cy.wait('@modelArchived').then((interception) => {
      expect(interception.request.body).to.eql(mockModArchResponse({ state: 'ARCHIVED' }));
    });
  });
});
