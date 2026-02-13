/* eslint-disable camelcase */
import type { Namespace } from 'mod-arch-core';
import { mockNamespace } from '~/__mocks__/mockNamespace';
import { mockModelRegistry } from '~/__mocks__/mockModelRegistry';
import { registerAndStoreFields } from '~/__tests__/cypress/cypress/pages/modelRegistryView/registerAndStoreFields';
import { MODEL_REGISTRY_API_VERSION } from '~/__tests__/cypress/cypress/support/commands/api';
import type { ModelRegistry, ModelTransferJob, RegisteredModel } from '~/app/types';
import { mockRegisteredModelList } from '~/__mocks__/mockRegisteredModelsList';
import { mockModelTransferJob } from '~/__mocks__/mockModelTransferJob';

type HandlersProps = {
  modelRegistries?: ModelRegistry[];
  namespaces?: Namespace[];
  registeredModels?: RegisteredModel[];
};

const initIntercepts = ({
  modelRegistries = [mockModelRegistry({ name: 'modelregistry-sample' })],
  namespaces = [
    mockNamespace({ name: 'namespace-1' }),
    mockNamespace({ name: 'namespace-2' }),
    mockNamespace({ name: 'namespace-3' }),
  ],
  registeredModels = [],
}: HandlersProps = {}) => {
  cy.interceptApi(
    'GET /api/:apiVersion/namespaces',
    {
      path: { apiVersion: MODEL_REGISTRY_API_VERSION },
    },
    namespaces,
  );

  cy.interceptApi(
    `GET /api/:apiVersion/model_registry`,
    {
      path: { apiVersion: MODEL_REGISTRY_API_VERSION },
    },
    modelRegistries,
  );

  cy.interceptApi(
    'GET /api/:apiVersion/model_registry/:modelRegistryName/registered_models',
    {
      path: {
        apiVersion: MODEL_REGISTRY_API_VERSION,
        modelRegistryName: 'modelregistry-sample',
      },
    },
    mockRegisteredModelList({ items: registeredModels }),
  );
};

describe('Register and Store Fields - Toggle Behavior', () => {
  beforeEach(() => {
    initIntercepts({});
    registerAndStoreFields.visit();
  });

  it('Should display registration mode toggle when feature is enabled', () => {
    registerAndStoreFields.shouldHaveRegistrationModeToggle();
  });

  it('Should have "Register" mode selected by default', () => {
    registerAndStoreFields.shouldHaveRegisterModeSelected();
  });

  it('Should switch to "Register and store" mode', () => {
    registerAndStoreFields.selectRegisterAndStoreMode();
    registerAndStoreFields.shouldHaveRegisterAndStoreModeSelected();
  });

  it('Should show namespace selector only in "Register and store" mode', () => {
    registerAndStoreFields.findNamespaceFormGroup().should('not.exist');
    registerAndStoreFields.selectRegisterAndStoreMode();
    registerAndStoreFields.findNamespaceSelector().should('exist');
  });

  it('Should switch back to "Register" mode and hide namespace selector', () => {
    registerAndStoreFields.selectRegisterAndStoreMode();
    registerAndStoreFields.findNamespaceSelector().should('exist');
    registerAndStoreFields.selectRegisterMode();
    registerAndStoreFields.shouldHaveRegisterModeSelected();
    registerAndStoreFields.findNamespaceSelector().should('not.exist');
  });

  it('Should reset namespace selection when switching modes', () => {
    registerAndStoreFields.selectRegisterAndStoreMode();
    registerAndStoreFields.selectNamespace('namespace-1');
    registerAndStoreFields.shouldShowSelectedNamespace('namespace-1');
    registerAndStoreFields.selectRegisterMode();
    registerAndStoreFields.selectRegisterAndStoreMode();
    registerAndStoreFields.shouldShowPlaceholder('Select a namespace');
  });
});

describe('Register and Store Fields - NamespaceSelector', () => {
  beforeEach(() => {
    initIntercepts({});
    registerAndStoreFields.visit();
    registerAndStoreFields.selectRegisterAndStoreMode();
  });

  it('Should show placeholder text instead of auto-selecting', () => {
    registerAndStoreFields.shouldShowPlaceholder('Select a namespace');
  });

  it('Should display all available namespaces in dropdown', () => {
    registerAndStoreFields.shouldHaveNamespaceOptions([
      'namespace-1',
      'namespace-2',
      'namespace-3',
    ]);
  });

  it('Should hide form sections until namespace is selected', () => {
    registerAndStoreFields.shouldHideOriginLocationSection().shouldHideDestinationLocationSection();
  });

  it('Should show form sections after namespace selection', () => {
    registerAndStoreFields.selectNamespace('namespace-1');

    registerAndStoreFields.shouldShowOriginLocationSection();
    registerAndStoreFields.shouldShowDestinationLocationSection();
  });

  it('Should update selected namespace in dropdown', () => {
    registerAndStoreFields.selectNamespace('namespace-2');
    registerAndStoreFields.shouldShowSelectedNamespace('namespace-2');
  });

  it('Should handle empty namespace list gracefully', () => {
    initIntercepts({ namespaces: [] });
    registerAndStoreFields.visit();
    registerAndStoreFields.selectRegisterAndStoreMode();

    registerAndStoreFields.findNamespaceSelector().should('exist');
    registerAndStoreFields.findNamespaceSelector().should('be.disabled');

    registerAndStoreFields.shouldShowPlaceholder('Select a namespace');
  });
});

describe('Register and Store Fields - Credential Validation', () => {
  beforeEach(() => {
    initIntercepts({});
    registerAndStoreFields.visit();
    registerAndStoreFields.selectRegisterAndStoreMode();
    registerAndStoreFields.selectNamespace('namespace-1');
  });

  it('Should have submit button disabled when S3 access key ID is missing', () => {
    // Fill all fields except S3 access key ID
    registerAndStoreFields.fillModelName('test-model');
    registerAndStoreFields.fillVersionName('v1.0.0');
    registerAndStoreFields.fillJobName('my-transfer-job');
    registerAndStoreFields.fillSourceEndpoint('https://s3.amazonaws.com');
    registerAndStoreFields.fillSourceBucket('test-bucket');
    registerAndStoreFields.fillSourcePath('models/test');
    // Skip: fillSourceS3AccessKeyId
    registerAndStoreFields.fillSourceS3SecretAccessKey('wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY');
    registerAndStoreFields.fillDestinationOciRegistry('quay.io');
    registerAndStoreFields.fillDestinationOciUri('quay.io/my-org/my-model:v1');
    registerAndStoreFields.fillDestinationOciUsername('testuser');
    registerAndStoreFields.fillDestinationOciPassword('testpassword123');

    registerAndStoreFields.findSubmitButton().should('be.disabled');
  });

  it('Should have submit button disabled when S3 secret access key is missing', () => {
    // Fill all fields except S3 secret access key
    registerAndStoreFields.fillModelName('test-model');
    registerAndStoreFields.fillVersionName('v1.0.0');
    registerAndStoreFields.fillJobName('my-transfer-job');
    registerAndStoreFields.fillSourceEndpoint('https://s3.amazonaws.com');
    registerAndStoreFields.fillSourceBucket('test-bucket');
    registerAndStoreFields.fillSourcePath('models/test');
    registerAndStoreFields.fillSourceS3AccessKeyId('AKIAIOSFODNN7EXAMPLE');
    // Skip: fillSourceS3SecretAccessKey
    registerAndStoreFields.fillDestinationOciRegistry('quay.io');
    registerAndStoreFields.fillDestinationOciUri('quay.io/my-org/my-model:v1');
    registerAndStoreFields.fillDestinationOciUsername('testuser');
    registerAndStoreFields.fillDestinationOciPassword('testpassword123');

    registerAndStoreFields.findSubmitButton().should('be.disabled');
  });

  it('Should have submit button disabled when OCI username is missing', () => {
    // Fill all fields except OCI username
    registerAndStoreFields.fillModelName('test-model');
    registerAndStoreFields.fillVersionName('v1.0.0');
    registerAndStoreFields.fillJobName('my-transfer-job');
    registerAndStoreFields.fillSourceEndpoint('https://s3.amazonaws.com');
    registerAndStoreFields.fillSourceBucket('test-bucket');
    registerAndStoreFields.fillSourcePath('models/test');
    registerAndStoreFields.fillSourceS3AccessKeyId('AKIAIOSFODNN7EXAMPLE');
    registerAndStoreFields.fillSourceS3SecretAccessKey('wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY');
    registerAndStoreFields.fillDestinationOciRegistry('quay.io');
    registerAndStoreFields.fillDestinationOciUri('quay.io/my-org/my-model:v1');
    // Skip: fillDestinationOciUsername
    registerAndStoreFields.fillDestinationOciPassword('testpassword123');

    registerAndStoreFields.findSubmitButton().should('be.disabled');
  });

  it('Should have submit button disabled when OCI password is missing', () => {
    // Fill all fields except OCI password
    registerAndStoreFields.fillModelName('test-model');
    registerAndStoreFields.fillVersionName('v1.0.0');
    registerAndStoreFields.fillJobName('my-transfer-job');
    registerAndStoreFields.fillSourceEndpoint('https://s3.amazonaws.com');
    registerAndStoreFields.fillSourceBucket('test-bucket');
    registerAndStoreFields.fillSourcePath('models/test');
    registerAndStoreFields.fillSourceS3AccessKeyId('AKIAIOSFODNN7EXAMPLE');
    registerAndStoreFields.fillSourceS3SecretAccessKey('wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY');
    registerAndStoreFields.fillDestinationOciRegistry('quay.io');
    registerAndStoreFields.fillDestinationOciUri('quay.io/my-org/my-model:v1');
    registerAndStoreFields.fillDestinationOciUsername('testuser');
    // Skip: fillDestinationOciPassword

    registerAndStoreFields.findSubmitButton().should('be.disabled');
  });

  it('Should enable submit button when all credentials are provided', () => {
    registerAndStoreFields.fillAllRequiredFields();
    registerAndStoreFields.findSubmitButton().should('not.be.disabled');
  });
});

describe('Register and Store Fields - Form Submission', () => {
  beforeEach(() => {
    initIntercepts({});
    registerAndStoreFields.visit();
  });

  it('Should have submit button disabled when required fields are empty', () => {
    registerAndStoreFields.selectRegisterAndStoreMode();
    registerAndStoreFields.selectNamespace('namespace-1');
    registerAndStoreFields.findSubmitButton().should('be.disabled');
  });

  it('Should enable submit button when all required fields are filled', () => {
    registerAndStoreFields.selectRegisterAndStoreMode();
    registerAndStoreFields.selectNamespace('namespace-1');
    registerAndStoreFields.fillAllRequiredFields();
    registerAndStoreFields.findSubmitButton().should('not.be.disabled');
  });

  it('Should create transfer job and navigate to model list on success', () => {
    const mockJob = mockModelTransferJob({ id: 'new-job-id' });

    cy.interceptApi(
      'POST /api/:apiVersion/model_registry/:modelRegistryName/model_transfer_jobs',
      {
        path: {
          apiVersion: MODEL_REGISTRY_API_VERSION,
          modelRegistryName: 'modelregistry-sample',
        },
      },
      mockJob,
    ).as('createTransferJob');

    registerAndStoreFields.selectRegisterAndStoreMode();
    registerAndStoreFields.selectNamespace('namespace-1');
    // Verify namespace is selected before filling other fields
    registerAndStoreFields.shouldShowSelectedNamespace('namespace-1');
    registerAndStoreFields.fillAllRequiredFields();
    registerAndStoreFields.findSubmitButton().click();

    cy.wait('@createTransferJob').then((interception) => {
      // Body might be a string if Content-Type isn't detected correctly
      const rawBody =
        typeof interception.request.body === 'string'
          ? JSON.parse(interception.request.body)
          : interception.request.body;
      // assembleModArchBody wraps the payload in { data: ... }
      const body: ModelTransferJob = rawBody.data || rawBody;
      expect(body.namespace).to.equal('namespace-1');
      expect(body.destination.uri).to.equal('quay.io/my-org/my-model:v1');
    });

    // Should navigate to model list (not version page)
    cy.url().should('include', '/model-registry/modelregistry-sample');
    cy.url().should('not.include', '/register');
  });

  it('Should show error when transfer job creation fails', () => {
    // Use raw cy.intercept for error responses to avoid mockModArchResponse wrapper
    cy.intercept(
      {
        method: 'POST',
        pathname: `/model-registry/api/${MODEL_REGISTRY_API_VERSION}/model_registry/modelregistry-sample/model_transfer_jobs`,
      },
      { statusCode: 500, body: { error: 'Failed to create transfer job' } },
    ).as('createTransferJobError');

    registerAndStoreFields.selectRegisterAndStoreMode();
    registerAndStoreFields.selectNamespace('namespace-1');
    registerAndStoreFields.fillAllRequiredFields();
    registerAndStoreFields.findSubmitButton().click();

    cy.wait('@createTransferJobError');
    cy.url().should('include', '/register');
  });

  it('Should NOT call registerModel API in Register and Store mode', () => {
    const mockJob = mockModelTransferJob({ id: 'new-job-id' });

    cy.interceptApi(
      'POST /api/:apiVersion/model_registry/:modelRegistryName/model_transfer_jobs',
      {
        path: {
          apiVersion: MODEL_REGISTRY_API_VERSION,
          modelRegistryName: 'modelregistry-sample',
        },
      },
      mockJob,
    ).as('createTransferJob');

    cy.interceptApi(
      'POST /api/:apiVersion/model_registry/:modelRegistryName/registered_models',
      {
        path: {
          apiVersion: MODEL_REGISTRY_API_VERSION,
          modelRegistryName: 'modelregistry-sample',
        },
      },
      {} as RegisteredModel,
    ).as('createRegisteredModel');

    registerAndStoreFields.selectRegisterAndStoreMode();
    registerAndStoreFields.selectNamespace('namespace-1');
    registerAndStoreFields.fillAllRequiredFields();
    registerAndStoreFields.findSubmitButton().click();

    cy.wait('@createTransferJob');
    cy.get('@createRegisteredModel.all').should('have.length', 0);
  });
});
