/* eslint-disable camelcase */
import { mockModArchResponse } from 'mod-arch-core';
import { verifyRelativeURL } from '~/__tests__/cypress/cypress/utils/url';
import { mockModelRegistry } from '~/__mocks__/mockModelRegistry';
import { mockRegisteredModel } from '~/__mocks__/mockRegisteredModel';
import { mockModelVersionList } from '~/__mocks__/mockModelVersionList';
import { mockModelVersion } from '~/__mocks__/mockModelVersion';
import { mockModelArtifactList } from '~/__mocks__/mockModelArtifactList';
import { ModelRegistryMetadataType, ModelState, type ModelRegistry } from '~/app/types';
import { MODEL_REGISTRY_API_VERSION } from '~/__tests__/cypress/cypress/support/commands/api';
import { modelVersionDetails } from '~/__tests__/cypress/cypress/pages/modelRegistryView/modelVersionDetails';

const mockModelVersions = mockModelVersion({
  id: '1',
  name: 'Version 1',
  customProperties: {
    a1: {
      metadataType: ModelRegistryMetadataType.STRING,
      string_value: 'v1',
    },
    a2: {
      metadataType: ModelRegistryMetadataType.STRING,
      string_value: 'v2',
    },
    a3: {
      metadataType: ModelRegistryMetadataType.STRING,
      string_value: 'v3',
    },
    a4: {
      metadataType: ModelRegistryMetadataType.STRING,
      string_value: 'v4',
    },
    a5: {
      metadataType: ModelRegistryMetadataType.STRING,
      string_value: 'v5',
    },
    a6: {
      metadataType: ModelRegistryMetadataType.STRING,
      string_value: 'v1',
    },
    a7: {
      metadataType: ModelRegistryMetadataType.STRING,
      string_value: 'v7',
    },
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
        mockModelVersion({
          name: 'Version 1',
          author: 'Author 1',
          registeredModelId: '1',
          createTimeSinceEpoch: '1712234877000', // Older
        }),
        mockModelVersion({
          author: 'Author 2',
          registeredModelId: '1',
          id: '2',
          name: 'Version 2',
          createTimeSinceEpoch: '1712234879000', // Latest
        }),
        mockModelVersion({
          author: 'Author 3',
          registeredModelId: '1',
          id: '3',
          name: 'Version 3',
          state: ModelState.ARCHIVED,
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
    mockModelVersions,
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
    `PATCH /api/:apiVersion/model_registry/:modelRegistryName/model_versions/:modelVersionId`,
    {
      path: {
        modelRegistryName: 'modelregistry-sample',
        apiVersion: MODEL_REGISTRY_API_VERSION,
        modelVersionId: 1,
      },
    },
    mockModelVersions,
  ).as('UpdatePropertyRow');

  cy.interceptApi(
    `GET /api/:apiVersion/model_registry/:modelRegistryName/model_versions/:modelVersionId/artifacts`,
    {
      path: {
        modelRegistryName: 'modelregistry-sample',
        apiVersion: MODEL_REGISTRY_API_VERSION,
        modelVersionId: 1,
      },
    },
    mockModelArtifactList({}),
  );
};

describe('Model version details', () => {
  describe('Overview tab', () => {
    beforeEach(() => {
      initIntercepts({});
      modelVersionDetails.visit();
    });

    it('Model version details page header', () => {
      verifyRelativeURL(
        '/model-registry/modelregistry-sample/registeredModels/1/versions/1/details',
      );
      cy.findByTestId('app-page-title').should('contain.text', 'Version 1');
      cy.findByTestId('breadcrumb-version-name').should('have.text', 'Version 1');
      cy.findByTestId('breadcrumb-model-version').should('contain.text', 'test');
    });

    it('should add a property', () => {
      modelVersionDetails.findAddPropertyButton().click();
      modelVersionDetails.findAddKeyInput().type('new_key');
      modelVersionDetails.findAddValueInput().type('new_value');
      modelVersionDetails.findCancelButton().click();

      modelVersionDetails.findAddPropertyButton().click();
      modelVersionDetails.findAddKeyInput().type('new_key');
      modelVersionDetails.findAddValueInput().type('new_value');
      modelVersionDetails.findSaveButton().click();
      cy.wait('@UpdatePropertyRow');
    });

    it('should edit a property row', () => {
      modelVersionDetails.findExpandControlButton().should('have.text', 'Show 2 more properties');
      modelVersionDetails.findExpandControlButton().click();
      const propertyRow = modelVersionDetails.getRow('a6');
      propertyRow.find().findKebabAction('Edit').click();
      modelVersionDetails.findKeyEditInput('a6').clear().type('edit_key');
      modelVersionDetails.findValueEditInput('v1').clear().type('edit_value');

      modelVersionDetails.findCancelButton().click();
      propertyRow.find().findKebabAction('Edit').click();
      modelVersionDetails.findKeyEditInput('a6').clear().type('edit_key');
      modelVersionDetails.findValueEditInput('v1').clear().type('edit_value');
      modelVersionDetails.findSaveButton().click();
      cy.wait('@UpdatePropertyRow').then((interception) => {
        expect(interception.request.body).to.eql(
          mockModArchResponse({
            customProperties: {
              a1: { metadataType: 'MetadataStringValue', string_value: 'v1' },
              a2: { metadataType: 'MetadataStringValue', string_value: 'v2' },
              a3: { metadataType: 'MetadataStringValue', string_value: 'v3' },
              a4: { metadataType: 'MetadataStringValue', string_value: 'v4' },
              a5: { metadataType: 'MetadataStringValue', string_value: 'v5' },
              a7: { metadataType: 'MetadataStringValue', string_value: 'v7' },
              'Testing label': { metadataType: 'MetadataStringValue', string_value: '' },
              'Financial data': { metadataType: 'MetadataStringValue', string_value: '' },
              'Fraud detection': { metadataType: 'MetadataStringValue', string_value: '' },
              'Long label data to be truncated abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc':
                { metadataType: 'MetadataStringValue', string_value: '' },
              'Machine learning': { metadataType: 'MetadataStringValue', string_value: '' },
              'Next data to be overflow': { metadataType: 'MetadataStringValue', string_value: '' },
              'Label x': { metadataType: 'MetadataStringValue', string_value: '' },
              'Label y': { metadataType: 'MetadataStringValue', string_value: '' },
              'Label z': { metadataType: 'MetadataStringValue', string_value: '' },
              edit_key: { string_value: 'edit_value', metadataType: 'MetadataStringValue' },
            },
          }),
        );
      });
    });

    it('should delete a property row', () => {
      modelVersionDetails.findExpandControlButton().should('have.text', 'Show 2 more properties');
      modelVersionDetails.findExpandControlButton().click();
      const propertyRow = modelVersionDetails.getRow('a6');
      modelVersionDetails.findPropertiesTableRows().should('have.length', 7);
      propertyRow.find().findKebabAction('Delete').click();
      cy.wait('@UpdatePropertyRow').then((interception) => {
        expect(interception.request.body).to.eql(
          mockModArchResponse({
            customProperties: {
              a1: { metadataType: 'MetadataStringValue', string_value: 'v1' },
              a2: { metadataType: 'MetadataStringValue', string_value: 'v2' },
              a3: { metadataType: 'MetadataStringValue', string_value: 'v3' },
              a4: { metadataType: 'MetadataStringValue', string_value: 'v4' },
              a5: { metadataType: 'MetadataStringValue', string_value: 'v5' },
              a7: { metadataType: 'MetadataStringValue', string_value: 'v7' },
              'Testing label': { metadataType: 'MetadataStringValue', string_value: '' },
              'Financial data': { metadataType: 'MetadataStringValue', string_value: '' },
              'Fraud detection': { metadataType: 'MetadataStringValue', string_value: '' },
              'Long label data to be truncated abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc':
                { metadataType: 'MetadataStringValue', string_value: '' },
              'Machine learning': { metadataType: 'MetadataStringValue', string_value: '' },
              'Next data to be overflow': { metadataType: 'MetadataStringValue', string_value: '' },
              'Label x': { metadataType: 'MetadataStringValue', string_value: '' },
              'Label y': { metadataType: 'MetadataStringValue', string_value: '' },
              'Label z': { metadataType: 'MetadataStringValue', string_value: '' },
            },
          }),
        );
      });
    });

    it('Switching model versions', () => {
      cy.interceptApi(
        `GET /api/:apiVersion/model_registry/:modelRegistryName/model_versions/:modelVersionId/artifacts`,
        {
          path: {
            modelRegistryName: 'modelregistry-sample',
            apiVersion: MODEL_REGISTRY_API_VERSION,
            modelVersionId: 2,
          },
        },
        mockModelArtifactList({}),
      );
      modelVersionDetails.findVersionId().contains('1');
      modelVersionDetails.findModelVersionDropdownButton().click();
      modelVersionDetails.findModelVersionDropdownItem('Version 3').should('not.exist');
      modelVersionDetails.findModelVersionDropdownSearch().fill('Version 2');
      modelVersionDetails.findModelVersionDropdownItem('Version 2').click();
      modelVersionDetails.findVersionId().contains('2');
    });

    it('should show "Latest" badge on the most recent version', () => {
      modelVersionDetails.findModelVersionDropdownButton().click();

      // Check that the "Latest" badge exists and is associated with Version 2
      cy.findByTestId('model-version-selector-list')
        .should('contain.text', 'Latest')
        .and('contain.text', 'Version 2');

      // Verify that Version 1 exists but doesn't have its own "Latest" badge
      cy.findByTestId('model-version-selector-list').should('contain.text', 'Version 1');

      cy.findByTestId('model-version-selector-list')
        .find('.pf-v6-c-badge')
        .should('have.length', 1);
    });

    it('should handle label editing', () => {
      modelVersionDetails.findEditLabelsButton().click();

      modelVersionDetails.findAddLabelButton().click();
      cy.findByTestId('editable-label-group')
        .should('exist')
        .within(() => {
          cy.contains('New Label').should('exist').click();
          cy.focused().type('First Label{enter}');
        });

      modelVersionDetails.findAddLabelButton().click();
      cy.findByTestId('editable-label-group')
        .should('exist')
        .within(() => {
          cy.contains('New Label').should('exist').click();
          cy.focused().type('Second Label{enter}');
        });

      cy.findByTestId('editable-label-group').within(() => {
        cy.contains('First Label').should('exist').click();
        cy.focused().type('Updated First Label{enter}');
      });

      cy.findByTestId('editable-label-group').within(() => {
        cy.contains('Second Label').parent().find('[data-testid^="remove-label-"]').click();
      });

      modelVersionDetails.findSaveLabelsButton().should('exist').click();
    });

    it('should validate label length', () => {
      modelVersionDetails.findEditLabelsButton().click();

      const longLabel = 'a'.repeat(64);
      modelVersionDetails.findAddLabelButton().click();
      cy.findByTestId('editable-label-group')
        .should('exist')
        .within(() => {
          cy.contains('New Label').should('exist').click();
          cy.focused().type(`${longLabel}{enter}`);
        });

      cy.findByTestId('label-error-alert')
        .should('be.visible')
        .within(() => {
          cy.contains(`can't exceed 63 characters`).should('exist');
        });
    });

    it('should validate duplicate labels', () => {
      modelVersionDetails.findEditLabelsButton().click();

      modelVersionDetails.findAddLabelButton().click();
      cy.findByTestId('editable-label-group')
        .should('exist')
        .within(() => {
          cy.get('[data-testid^="editable-label-"]').last().click();
          cy.focused().type('{selectall}{backspace}Testing label{enter}');
        });

      modelVersionDetails.findAddLabelButton().click();
      cy.findByTestId('editable-label-group')
        .should('exist')
        .within(() => {
          cy.get('[data-testid^="editable-label-"]').last().click();
          cy.focused().type('{selectall}{backspace}Testing label{enter}');
        });

      cy.findByTestId('label-error-alert')
        .should('be.visible')
        .within(() => {
          cy.contains('Testing label already exists').should('exist');
        });
    });

    it('should navigate to versions list when clicking ViewAllVersionsButton', () => {
      modelVersionDetails.visit();
      modelVersionDetails.findModelVersionDropdownButton().click();

      cy.findByTestId('versions-route-link')
        .should('exist')
        .and('contain.text', 'View all')
        .and('contain.text', 'versions');

      // Click the link and verify navigation
      cy.findByTestId('versions-route-link').click();

      // Verify we navigated to the versions list page
      cy.url().should(
        'include',
        '/model-registry/modelregistry-sample/registeredModels/1/versions',
      );
      cy.findByTestId('model-versions-tab-content').should('exist');
    });
  });
});
