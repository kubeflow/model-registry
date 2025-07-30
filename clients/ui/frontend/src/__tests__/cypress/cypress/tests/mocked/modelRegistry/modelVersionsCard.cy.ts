/* eslint-disable camelcase */
import { mockModArchResponse } from 'mod-arch-shared';
import { mockModelRegistry } from '~/__mocks__/mockModelRegistry';
import { mockRegisteredModel } from '~/__mocks__/mockRegisteredModel';
import { mockModelVersionList } from '~/__mocks__/mockModelVersionList';
import { mockModelVersion } from '~/__mocks__/mockModelVersion';
import { ModelRegistryMetadataType, type ModelRegistry } from '~/app/types';
import { MODEL_REGISTRY_API_VERSION } from '~/__tests__/cypress/cypress/support/commands/api';
import { modelVersionsCard } from '../../../pages/modelRegistryView/modelVersionsCard';

const mockRegisteredModelWithData = mockRegisteredModel({
  id: '1',
  name: 'Test Model',
  description: 'Test model description',
  owner: 'test-owner',
  customProperties: {
    label1: {
      metadataType: ModelRegistryMetadataType.STRING,
      string_value: '',
    },
    label2: {
      metadataType: ModelRegistryMetadataType.STRING,
      string_value: '',
    },
    property1: {
      metadataType: ModelRegistryMetadataType.STRING,
      string_value: 'value1',
    },
    property2: {
      metadataType: ModelRegistryMetadataType.STRING,
      string_value: 'value2',
    },
    'url-property': {
      metadataType: ModelRegistryMetadataType.STRING,
      string_value: 'https://example.com',
    },
  },
});

const mockModelVersionListWithData = mockModelVersionList({
  items: [
    mockModelVersion({
      id: '1',
      name: 'Version 1',
      author: 'Author 1',
      registeredModelId: '1',
      createTimeSinceEpoch: '1725282249921',
      customProperties: {
        label1: {
          metadataType: ModelRegistryMetadataType.STRING,
          string_value: '',
        },
        label2: {
          metadataType: ModelRegistryMetadataType.STRING,
          string_value: '',
        },
        property1: {
          metadataType: ModelRegistryMetadataType.STRING,
          string_value: 'value1',
        },
        property2: {
          metadataType: ModelRegistryMetadataType.STRING,
          string_value: 'value2',
        },
      },
    }),
    mockModelVersion({
      id: '2',
      name: 'Version 2',
      author: 'Author 2',
      createTimeSinceEpoch: '1725282249920',
      registeredModelId: '1',
    }),
    mockModelVersion({
      id: '3',
      name: 'Version 3',
      author: 'Author 3',
      createTimeSinceEpoch: '1725282249925',
      registeredModelId: '1',
    }),
    mockModelVersion({
      id: '4',
      name: 'Version 4',
      author: 'Author 4',
      registeredModelId: '1',
      createTimeSinceEpoch: '1725282349921',
    }),
  ],
});

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
    mockRegisteredModelWithData,
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
    mockModelVersionListWithData,
  );

  cy.interceptApi(
    `GET /api/:apiVersion/model_registry/:modelRegistryName/model_versions/:modelVersionId`,
    {
      path: {
        modelRegistryName: 'modelregistry-sample',
        apiVersion: MODEL_REGISTRY_API_VERSION,
        modelVersionId: 4,
      },
    },
    mockModelVersionListWithData.items[3],
  );
};

describe('Model Versions Card', () => {
  beforeEach(() => {
    initIntercepts({});
    mockModArchResponse({});
  });

  it('does not display model versions list if there are no model versions', () => {
    cy.visit('/model-registry/modelregistry-sample/registeredModels/1/overview');
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
        items: [],
      }),
    );

    modelVersionsCard.findNoVersionsText().should('be.visible');
  });

  it('should display model versions list correctly', () => {
    cy.visit('/model-registry/modelregistry-sample/registeredModels/1/overview');

    modelVersionsCard.findModelVersion('1').should('exist');

    modelVersionsCard.findModelVersionProperty('1', 'label1').should('exist');
    modelVersionsCard.findModelVersionProperty('1', 'label2').should('exist');
    modelVersionsCard.findModelVersionProperty('1', 'property1').should('exist');
    modelVersionsCard.findModelVersionLink('1').should('exist');

    modelVersionsCard.findModelVersion('2').should('not.exist');
    modelVersionsCard.findModelVersion('3').should('exist');
    modelVersionsCard.findModelVersion('4').should('exist');

    modelVersionsCard.findViewAllVersionsLink().click();
    cy.url().should('include', '/model-registry/modelregistry-sample/registeredModels/1/versions');
  });

  it('should have the correct link to the model version', () => {
    cy.visit('/model-registry/modelregistry-sample/registeredModels/1/overview');

    modelVersionsCard.findModelVersionLink('4').click();
    cy.url().should(
      'include',
      '/model-registry/modelregistry-sample/registeredModels/1/versions/4',
    );
    cy.contains('Version 4').should('be.visible');
  });
});
