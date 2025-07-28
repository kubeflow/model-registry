/* eslint-disable camelcase */
import { mockModArchResponse } from 'mod-arch-shared';
import { mockModelRegistry } from '~/__mocks__/mockModelRegistry';
import { mockRegisteredModel } from '~/__mocks__/mockRegisteredModel';
import { mockModelVersionList } from '~/__mocks__/mockModelVersionList';
import { mockModelVersion } from '~/__mocks__/mockModelVersion';
import { ModelRegistryMetadataType, ModelState, type ModelRegistry } from '~/app/types';
import { MODEL_REGISTRY_API_VERSION } from '~/__tests__/cypress/cypress/support/commands/api';
import { modelDetailsCard } from '~/__tests__/cypress/cypress/pages/modelRegistryView/modelDetailsCard';

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
    mockModelVersionList({
      items: [mockModelVersion({ name: 'Version 1', author: 'Author 1', registeredModelId: '1' })],
    }),
  );

  cy.interceptApi(
    `PATCH /api/:apiVersion/model_registry/:modelRegistryName/registered_models/:registeredModelId`,
    {
      path: {
        modelRegistryName: 'modelregistry-sample',
        apiVersion: MODEL_REGISTRY_API_VERSION,
        registeredModelId: 1,
      },
    },
    mockRegisteredModelWithData,
  ).as('patchRegisteredModel');
};

describe('Model Details Card', () => {
  beforeEach(() => {
    initIntercepts({});
    mockModArchResponse({});
  });

  it('displays model details correctly', () => {
    cy.visit('/model-registry/modelregistry-sample/registeredModels/1/overview');

    // Verify card title
    cy.contains('Model details').should('be.visible');

    // Verify description
    cy.findByText('Test model description').should('be.visible');

    // Verify owner
    modelDetailsCard.findOwner().should('contain', 'test-owner');

    // Verify model ID clipboard copy
    cy.contains('Model ID').should('be.visible');
    cy.findByTestId('registered-model-id-clipboard-copy').should('exist');

    // Verify timestamps sections exist
    cy.contains('Last modified at').should('be.visible');
    cy.contains('Created at').should('be.visible');
  });

  it('displays labels section correctly', () => {
    cy.visit('/model-registry/modelregistry-sample/registeredModels/1/overview');

    // Verify labels are displayed
    cy.contains('Labels').should('be.visible');
    cy.contains('label1').should('be.visible');
    cy.contains('label2').should('be.visible');

    // Verify labels section has edit capabilities
    cy.contains('Labels').parent().find('button').should('exist');
  });

  it('displays properties in expandable section', () => {
    cy.visit('/model-registry/modelregistry-sample/registeredModels/1/overview');

    // Check properties section exists with correct badge count
    cy.contains('Properties').should('be.visible');
    cy.contains('Properties').parent().find('.pf-v6-c-badge').should('contain', '3'); // property1, property2, url-property

    // Expand properties section
    cy.contains('Properties').click();

    // Verify properties table is visible
    modelDetailsCard.findPropertiesTable().should('be.visible');

    // Verify properties are displayed
    cy.contains('property1').should('be.visible');
    cy.contains('value1').should('be.visible');
    cy.contains('property2').should('be.visible');
    cy.contains('value2').should('be.visible');

    // Verify URL property is displayed as link
    cy.contains('url-property').should('be.visible');
    cy.get('a[href="https://example.com"]').should('exist');
  });

  it('shows add property button and validates input', () => {
    cy.visit('/model-registry/modelregistry-sample/registeredModels/1/overview');

    // Expand properties section
    cy.contains('Properties').click();

    // Verify add property button exists
    modelDetailsCard.findAddPropertyButton().should('be.visible');

    // Click add property button to test interaction
    modelDetailsCard.findAddPropertyButton().click();

    // Verify input fields appear
    modelDetailsCard.findAddPropertyKeyInput().should('be.visible');
    modelDetailsCard.findAddPropertyValueInput().should('be.visible');

    // Test validation
    modelDetailsCard.findAddPropertyKeyInput().type('property1'); // Already exists
    modelDetailsCard.findAddPropertyValueInput().type('someValue');

    // Verify validation error appears
    cy.contains('Key must not match an existing property key or label').should('be.visible');

    // Verify save button is disabled
    modelDetailsCard.findSavePropertyButton().should('be.disabled');
  });

  it('validates property key length correctly', () => {
    cy.visit('/model-registry/modelregistry-sample/registeredModels/1/overview');

    // Expand properties section
    cy.contains('Properties').click();

    // Click add property button
    modelDetailsCard.findAddPropertyButton().click();

    // Test key length validation
    modelDetailsCard.findAddPropertyKeyInput().type('a'.repeat(64)); // Too long
    modelDetailsCard.findAddPropertyValueInput().type('someValue');

    cy.contains("Key text can't exceed 63 characters").should('be.visible');
    modelDetailsCard.findSavePropertyButton().should('be.disabled');
  });

  it('handles expand/collapse for many properties', () => {
    // Create mock with many properties
    const customPropsEntries = Array.from({ length: 10 }, (_, i) => [
      `property${i + 1}`,
      {
        metadataType: ModelRegistryMetadataType.STRING as const,
        string_value: `value${i + 1}`,
      },
    ]);
    const customProps = Object.fromEntries(customPropsEntries);

    const manyPropertiesModel = mockRegisteredModel({
      id: '1',
      customProperties: customProps,
    });

    cy.interceptApi(
      `GET /api/:apiVersion/model_registry/:modelRegistryName/registered_models/:registeredModelId`,
      {
        path: {
          modelRegistryName: 'modelregistry-sample',
          apiVersion: MODEL_REGISTRY_API_VERSION,
          registeredModelId: 1,
        },
      },
      manyPropertiesModel,
    );

    cy.visit('/model-registry/modelregistry-sample/registeredModels/1/overview');

    // Expand properties section
    cy.contains('Properties').click();

    // Verify only first 5 properties are shown initially
    cy.contains('property1').should('be.visible');
    cy.contains('property5').should('be.visible');
    cy.contains('property6').should('not.exist');

    // Click expand control button
    modelDetailsCard.findExpandControlButton().click();

    // Verify more properties are now visible
    cy.contains('property6').should('be.visible');
    cy.contains('property10').should('be.visible');

    // Verify button text changed
    modelDetailsCard.findExpandControlButton().should('contain', 'Show fewer properties');
  });

  it('handles archived model state correctly', () => {
    // Test with archived model - editing should be disabled
    const archivedModel = mockRegisteredModel({
      id: '1',
      state: ModelState.ARCHIVED,
      customProperties: {
        property1: {
          metadataType: ModelRegistryMetadataType.STRING,
          string_value: 'value1',
        },
      },
    });

    cy.interceptApi(
      `GET /api/:apiVersion/model_registry/:modelRegistryName/registered_models/:registeredModelId`,
      {
        path: {
          modelRegistryName: 'modelregistry-sample',
          apiVersion: MODEL_REGISTRY_API_VERSION,
          registeredModelId: 1,
        },
      },
      archivedModel,
    );

    cy.visit('/model-registry/modelregistry-sample/registeredModels/1/overview');

    // Expand properties section
    cy.contains('Properties').click();

    // Verify add property button is not present for archived model
    modelDetailsCard.findAddPropertyButton().should('not.exist');
  });

  it('shows the correct tab structure and navigation', () => {
    cy.visit('/model-registry/modelregistry-sample/registeredModels/1/overview');

    // Verify tab structure exists
    cy.findByTestId('model-versions-page-tabs').should('exist');
    cy.findByTestId('model-overview-tab').should('exist');

    // Verify overview tab is active
    cy.findByTestId('model-overview-tab').should('have.attr', 'aria-selected', 'true');

    // Verify model details content is displayed
    cy.findByTestId('model-details-tab-content').should('be.visible');
  });

  // TODO: Pending tests for complex interactions with mod-arch-shared components
  // These tests need investigation of exact DOM structure of mod-arch-shared components

  it('allows editing model description', () => {
    cy.visit('/model-registry/modelregistry-sample/registeredModels/1/overview');

    // Verify initial description is visible
    cy.findByText('Test model description').should('be.visible');

    // Find and click the edit button near the description section
    modelDetailsCard.findDescriptionEditButton().click();

    // Verify edit mode is active (should show input fields)
    modelDetailsCard.findDescriptionTextArea().should('be.visible');

    // Edit the description
    modelDetailsCard.findDescriptionTextArea().clear();
    modelDetailsCard.findDescriptionTextArea().type('Updated model description for testing');

    // Save the changes
    modelDetailsCard.findDescriptionSaveButton().click();

    // Verify API call was made
    cy.wait('@patchRegisteredModel').then((interception) => {
      expect(interception.request.body.data.description).to.equal(
        'Updated model description for testing',
      );
    });
  });
});
