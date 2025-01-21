/* eslint-disable camelcase */
import { mockRegisteredModelList } from '~/__mocks__/mockRegisteredModelsList';
import { mockModelVersionList } from '~/__mocks__/mockModelVersionList';
import { mockModelVersion } from '~/__mocks__/mockModelVersion';
import { mockRegisteredModel } from '~/__mocks__/mockRegisteredModel';
import { verifyRelativeURL } from '~/__tests__/cypress/cypress/utils/url';
import { labelModal, modelRegistry } from '~/__tests__/cypress/cypress/pages/modelRegistry';
import type { ModelRegistry, ModelVersion } from '~/app/types';
import { ModelRegistryMetadataType, ModelState } from '~/app/types';
import { mockModelRegistry } from '~/__mocks__/mockModelRegistry';
import { mockBFFResponse } from '~/__mocks__/utils';
import {
  archiveVersionModal,
  modelVersionArchive,
  restoreVersionModal,
} from '~/__tests__/cypress/cypress/pages/modelRegistryView/modelVersionArchive';
import { MODEL_REGISTRY_API_VERSION } from '~/__tests__/cypress/cypress/support/commands/api';
import { ToastNotification } from '~/__tests__/cypress/cypress/pages/components/Notification';

type HandlersProps = {
  registeredModelsSize?: number;
  modelVersions?: ModelVersion[];
  modelRegistries?: ModelRegistry[];
};

const initIntercepts = ({
  registeredModelsSize = 4,
  modelVersions = [
    mockModelVersion({
      name: 'model version 1',
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
      state: ModelState.ARCHIVED,
    }),
    mockModelVersion({ id: '2', name: 'model version 2', state: ModelState.ARCHIVED }),
    mockModelVersion({ id: '3', name: 'model version 3' }),
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
    mockModelVersionList({
      items: modelVersions,
    }),
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
    mockRegisteredModel({ name: 'test-1' }),
  );

  cy.interceptApi(
    `GET /api/:apiVersion/model_registry/:modelRegistryName/model_versions/:modelVersionId`,
    {
      path: {
        modelRegistryName: 'modelregistry-sample',
        apiVersion: MODEL_REGISTRY_API_VERSION,
        modelVersionId: 2,
      },
    },
    mockModelVersion({ id: '2', name: 'model version 2', state: ModelState.ARCHIVED }),
  );

  cy.interceptApi(
    `GET /api/:apiVersion/model_registry/:modelRegistryName/model_versions/:modelVersionId`,
    {
      path: {
        modelRegistryName: 'modelregistry-sample',
        apiVersion: MODEL_REGISTRY_API_VERSION,
        modelVersionId: 3,
      },
    },
    mockModelVersion({ id: '3', name: 'model version 3', state: ModelState.LIVE }),
  );
};

describe('Model version archive list', () => {
  it('No archive versions in the selected registered model', () => {
    initIntercepts({ modelVersions: [mockModelVersion({ id: '3', name: 'model version 2' })] });
    modelVersionArchive.visitModelVersionList();
    verifyRelativeURL('/model-registry/modelregistry-sample/registeredModels/1/versions');
    modelVersionArchive
      .findModelVersionsTableKebab()
      .findDropdownItem('View archived versions')
      .click();
    modelVersionArchive.shouldArchiveVersionsEmpty();
  });

  it('Archived version details browser back button should lead to archived versions table', () => {
    initIntercepts({});
    modelVersionArchive.visit();
    verifyRelativeURL('/model-registry/modelregistry-sample/registeredModels/1/versions/archive');
    modelVersionArchive.findArchiveVersionBreadcrumbItem().contains('Archived version');
    const archiveVersionRow = modelVersionArchive.getRow('model version 2');
    archiveVersionRow.findName().contains('model version 2').click();
    verifyRelativeURL(
      '/model-registry/modelregistry-sample/registeredModels/1/versions/archive/2/details',
    );
    cy.go('back');
    verifyRelativeURL('/model-registry/modelregistry-sample/registeredModels/1/versions/archive');
    modelVersionArchive.findArchiveVersionBreadcrumbItem().contains('Archived version');
    archiveVersionRow.findName().contains('model version 2').should('exist');
  });

  it('Archive version list', () => {
    initIntercepts({});
    modelVersionArchive.visit();
    verifyRelativeURL('/model-registry/modelregistry-sample/registeredModels/1/versions/archive');

    //breadcrumb
    modelVersionArchive.findArchiveVersionBreadcrumbItem().contains('Archived version');

    // name, last modified, owner, labels modal
    modelVersionArchive.findArchiveVersionTable().should('be.visible');
    modelVersionArchive.findArchiveVersionsTableRows().should('have.length', 2);

    const archiveVersionRow = modelVersionArchive.getRow('model version 1');

    archiveVersionRow.findLabelModalText().contains('5 more');
    archiveVersionRow.findLabelModalText().click();
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
  });
});

describe('Restoring archive version', () => {
  it('Restore from archive table', () => {
    cy.interceptApi(
      'PATCH /api/:apiVersion/model_registry/:modelRegistryName/model_versions/:modelVersionId',
      {
        path: {
          modelRegistryName: 'modelregistry-sample',
          apiVersion: MODEL_REGISTRY_API_VERSION,
          modelVersionId: 2,
        },
      },
      mockModelVersion({}),
    ).as('versionRestored');

    initIntercepts({});
    modelVersionArchive.visit();

    const archiveVersionRow = modelVersionArchive.getRow('model version 2');
    archiveVersionRow.findKebabAction('Restore model version').click();

    restoreVersionModal.findRestoreButton().click();

    const notification = new ToastNotification('model version 2 restored.');
    notification.find();

    cy.wait('@versionRestored').then((interception) => {
      expect(interception.request.body).to.eql(mockBFFResponse({ state: 'LIVE' }));
    });
  });

  it('Restore from archive version details', () => {
    cy.interceptApi(
      'PATCH /api/:apiVersion/model_registry/:modelRegistryName/model_versions/:modelVersionId',
      {
        path: {
          modelRegistryName: 'modelregistry-sample',
          apiVersion: MODEL_REGISTRY_API_VERSION,
          modelVersionId: 2,
        },
      },
      mockModelVersion({}),
    ).as('versionRestored');

    initIntercepts({});
    modelVersionArchive.visitArchiveVersionDetail();

    modelVersionArchive.findRestoreButton().click();
    restoreVersionModal.findRestoreButton().click();

    const notification = new ToastNotification('model version 2 restored.');
    notification.find();

    cy.wait('@versionRestored').then((interception) => {
      expect(interception.request.body).to.eql(mockBFFResponse({ state: 'LIVE' }));
    });
  });
});

describe('Archiving version', () => {
  it('Archive version from versions table', () => {
    cy.interceptApi(
      'PATCH /api/:apiVersion/model_registry/:modelRegistryName/model_versions/:modelVersionId',
      {
        path: {
          modelRegistryName: 'modelregistry-sample',
          apiVersion: MODEL_REGISTRY_API_VERSION,
          modelVersionId: 3,
        },
      },
      mockModelVersion({}),
    ).as('versionArchived');

    initIntercepts({});
    modelVersionArchive.visitModelVersionList();

    const modelVersionRow = modelRegistry.getModelVersionRow('model version 3');
    modelVersionRow.findKebabAction('Archive model version').click();
    archiveVersionModal.findArchiveButton().should('be.disabled');
    archiveVersionModal.findModalTextInput().fill('model version 3');
    archiveVersionModal.findArchiveButton().should('be.enabled').click();

    const notification = new ToastNotification('model version 3 archived.');
    notification.find();

    cy.wait('@versionArchived').then((interception) => {
      expect(interception.request.body).to.eql(mockBFFResponse({ state: 'ARCHIVED' }));
    });
  });

  it('Archived version details page does not have the Deployments tab', () => {
    initIntercepts({});
    modelVersionArchive.visitArchiveVersionDetail();
    modelVersionArchive.findVersionDetailsTab().should('exist');
    modelVersionArchive.findVersionDeploymentTab().should('not.exist');
  });

  it('Archive version from versions details', () => {
    cy.interceptApi(
      'PATCH /api/:apiVersion/model_registry/:modelRegistryName/model_versions/:modelVersionId',
      {
        path: {
          modelRegistryName: 'modelregistry-sample',
          apiVersion: MODEL_REGISTRY_API_VERSION,
          modelVersionId: 3,
        },
      },
      mockModelVersion({}),
    ).as('versionArchived');

    initIntercepts({});
    modelVersionArchive.visitModelVersionDetails();
    modelVersionArchive
      .findModelVersionsDetailsHeaderAction()
      .findDropdownItem('Archive model version')
      .click();

    archiveVersionModal.findArchiveButton().should('be.disabled');
    archiveVersionModal.findModalTextInput().fill('model version 3');
    archiveVersionModal.findArchiveButton().should('be.enabled').click();

    const notification = new ToastNotification('model version 3 archived.');
    notification.find();

    cy.wait('@versionArchived').then((interception) => {
      expect(interception.request.body).to.eql(mockBFFResponse({ state: 'ARCHIVED' }));
    });
  });
});
