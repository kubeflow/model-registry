/* eslint-disable camelcase */
import { verifyRelativeURL } from '~/__tests__/cypress/cypress/utils/url';
import { mockModelRegistry } from '~/__mocks__/mockModelRegistry';
import { mockRegisteredModel } from '~/__mocks__/mockRegisteredModel';
import { mockModelVersionList } from '~/__mocks__/mockModelVersionList';
import { mockModelVersion } from '~/__mocks__/mockModelVersion';
import { mockModelArtifactList } from '~/__mocks__/mockModelArtifactList';
import { mockModelArtifact } from '~/__mocks__/mockModelArtifact';
import type { ModelRegistry } from '~/app/types';
import { MODEL_REGISTRY_API_VERSION } from '~/__tests__/cypress/cypress/support/commands/api';
import { modelVersionDetails } from '~/__tests__/cypress/cypress/pages/modelRegistryView/modelVersionDetails';

type HandlersProps = {
  modelRegistries?: ModelRegistry[];
};

const initIntercepts = ({
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
    `GET /api/:apiVersion/model_registry/:modelRegistryName/registered_models/:registeredModelId/versions`,
    {
      path: {
        modelRegistryName: 'modelregistry-sample',
        apiVersion: MODEL_REGISTRY_API_VERSION,
        registeredModelId: 1,
      },
    },
    mockModelVersionList({
      items: [
        mockModelVersion({ name: 'Version 1', author: 'Author 1', registeredModelId: '1' }),
        mockModelVersion({
          author: 'Author 2',
          registeredModelId: '1',
          id: '2',
          name: 'Version 2',
        }),
      ],
    }),
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
    mockModelVersion({
      id: '1',
      name: 'Version 1',
      labels: [
        'Testing label',
        'Financial data',
        'Fraud detection',
        'Long label data to be truncated abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc',
        'Machine learning',
        'Next data to be overflow',
        'Label x',
        'Label y',
        'Label z',
      ],
    }),
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
    mockModelVersion({ id: '2', name: 'Version 2' }),
  );

  cy.interceptApi(
    `GET /api/:apiVersion/model_registry/:modelRegistryName/model_versions/:modelVersionId/artifacts`,
    {
      path: {
        modelRegistryName: 'modelregistry-sample',
        apiVersion: MODEL_REGISTRY_API_VERSION,
        modelVersionId: 1,
      },
    },
    mockModelArtifactList({
      items: [
        mockModelArtifact({}),
        mockModelArtifact({
          author: 'Author 2',
          id: '2',
          name: 'Artifact 2',
        }),
      ],
    }),
  );
};

describe('Model version details', () => {
  describe('Details tab', () => {
    beforeEach(() => {
      initIntercepts({});
      modelVersionDetails.visit();
    });

    it('Model version details page header', () => {
      verifyRelativeURL(
        '/model-registry/modelregistry-sample/registeredModels/1/versions/1/details',
      );
      cy.findByTestId('app-page-title').should('have.text', 'Version 1');
      cy.findByTestId('breadcrumb-version-name').should('have.text', 'Version 1');
    });

    it('Model version details tab', () => {
      modelVersionDetails.findVersionId().contains('1');
      modelVersionDetails.findDescription().should('have.text', 'Description of model version');
      modelVersionDetails.findMoreLabelsButton().contains('6 more');
      modelVersionDetails.findMoreLabelsButton().click();
      modelVersionDetails.shouldContainsModalLabels([
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
      modelVersionDetails.findStorageEndpoint().contains('test-endpoint');
      modelVersionDetails.findStorageRegion().contains('test-region');
      modelVersionDetails.findStorageBucket().contains('test-bucket');
      modelVersionDetails.findStoragePath().contains('demo-models/test-path');
    });

    it('Switching model versions', () => {
      modelVersionDetails.findVersionId().contains('1');
      modelVersionDetails.findModelVersionDropdownButton().click();
      modelVersionDetails.findModelVersionDropdownSearch().fill('Version 2');
      modelVersionDetails.findModelVersionDropdownItem('Version 2').click();
      modelVersionDetails.findVersionId().contains('2');
    });
  });
});
