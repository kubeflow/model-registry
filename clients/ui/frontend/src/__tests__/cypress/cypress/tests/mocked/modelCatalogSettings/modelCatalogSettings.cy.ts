import {
  modelCatalogSettings,
  manageSourcePage,
} from '~/__tests__/cypress/cypress/pages/modelCatalogSettings';
import { MODEL_CATALOG_API_VERSION } from '~/__tests__/cypress/cypress/support/commands/api';
import { mockCatalogSource, mockCatalogSourceList } from '~/__mocks__';
import type { CatalogSource } from '~/app/modelCatalogTypes';

const NAMESPACE = 'kubeflow';
const userMock = {
  data: {
    userId: 'user@example.com',
    clusterAdmin: true,
  },
};

const setupMocks = (sources: CatalogSource[] = []) => {
  cy.intercept('GET', '/model-registry/api/v1/namespaces', {
    data: [{ metadata: { name: NAMESPACE } }],
  });
  cy.intercept('GET', '/model-registry/api/v1/user', userMock);
  cy.interceptApi(
    `GET /api/:apiVersion/model_catalog/sources`,
    {
      path: { apiVersion: MODEL_CATALOG_API_VERSION },
    },
    mockCatalogSourceList({
      items: sources,
    }),
  );
};

function selectNamespaceIfPresent() {
  cy.get('body').then(($body) => {
    if ($body.find('[data-testid="namespace-select"]').length) {
      cy.get('[data-testid="namespace-select"]').click();
      cy.findByText(NAMESPACE).click();
    }
  });
}

describe('Model Catalog Settings', () => {
  beforeEach(() => {
    setupMocks();
  });

  it('should display the settings page', () => {
    modelCatalogSettings.visit();
    modelCatalogSettings.findHeading();
    modelCatalogSettings.findDescription();
  });

  it('should navigate to settings page from nav', () => {
    selectNamespaceIfPresent();
    cy.get('body').then(($body) => {
      if ($body.find('#page-sidebar').length > 0) {
        modelCatalogSettings.navigate();
        modelCatalogSettings.findHeading();
      } else {
        cy.log('Sidebar not available, skipping nav test');
      }
    });
  });

  it('should display add source button', () => {
    modelCatalogSettings.visit();
    modelCatalogSettings.findAddSourceButton().should('be.visible');
    modelCatalogSettings.findAddSourceButton().should('contain', 'Add a source');
  });

  it('should navigate to add source page when button is clicked', () => {
    modelCatalogSettings.visit();
    modelCatalogSettings.findAddSourceButton().click();
    manageSourcePage.findAddSourceTitle();
    manageSourcePage.findAddSourceDescription();
    manageSourcePage.findBreadcrumb().should('exist');
    manageSourcePage.findBreadcrumbAction().should('contain', 'Add a source');
  });
});

describe('Manage Source Page', () => {
  beforeEach(() => {
    setupMocks();
  });

  describe('Add Source Mode', () => {
    it('should display add source page', () => {
      manageSourcePage.visitAddSource();
      manageSourcePage.findAddSourceTitle();
      manageSourcePage.findAddSourceDescription();
    });

    it('should display correct breadcrumb for add source', () => {
      manageSourcePage.visitAddSource();
      manageSourcePage.findBreadcrumb().should('exist');
      manageSourcePage.findBreadcrumbAction().should('contain', 'Add a source');
    });

    it('should navigate back to settings from breadcrumb', () => {
      manageSourcePage.visitAddSource();
      manageSourcePage.findBreadcrumb().click({ force: true });
      modelCatalogSettings.findHeading();
    });

    it('should navigate back to settings from cancel button', () => {
      manageSourcePage.visitAddSource();
      manageSourcePage.findCancelButton().click();
      modelCatalogSettings.findHeading();
    });

    it('should display form fields', () => {
      manageSourcePage.visitAddSource();
      manageSourcePage.findNameInput().should('exist');
      manageSourcePage.findSourceTypeHuggingFace().should('exist');
      manageSourcePage.findSourceTypeYaml().should('exist');
      manageSourcePage.findEnableSourceCheckbox().should('exist');
      manageSourcePage.findSubmitButton().should('exist');
      manageSourcePage.findPreviewButton().should('exist');
      manageSourcePage.findCancelButton().should('exist');
    });

    it('should show Hugging Face fields by default', () => {
      manageSourcePage.visitAddSource();
      manageSourcePage.findSourceTypeHuggingFace().should('be.checked');
      manageSourcePage.findCredentialsSection().should('exist');
      manageSourcePage.findAccessTokenInput().should('exist');
      manageSourcePage.findOrganizationInput().should('exist');
      manageSourcePage.findYamlSection().should('not.exist');
    });

    it('should show YAML fields when YAML type is selected', () => {
      manageSourcePage.visitAddSource();
      manageSourcePage.selectSourceType('yaml');
      manageSourcePage.findSourceTypeYaml().should('be.checked');
      manageSourcePage.findYamlSection().should('exist');
      manageSourcePage.findYamlContentInput().should('exist');
      manageSourcePage.findCredentialsSection().should('not.exist');
    });

    it('should have Add button disabled by default', () => {
      manageSourcePage.visitAddSource();
      manageSourcePage.findSubmitButton().should('be.disabled');
    });

    it('should have Preview button disabled by default', () => {
      manageSourcePage.visitAddSource();
      manageSourcePage.findPreviewButton().should('be.disabled');
    });

    it('should show validation error when name field is empty and touched', () => {
      manageSourcePage.visitAddSource();
      manageSourcePage.findNameInput().focus().blur();
      manageSourcePage.findNameError().should('exist');
      manageSourcePage.findNameError().should('contain', 'Name is required');
    });

    it('should enable Add button when all required HF fields are filled', () => {
      manageSourcePage.visitAddSource();
      manageSourcePage.fillSourceName('Test Source');
      manageSourcePage.fillAccessToken('test-token-123');
      manageSourcePage.fillOrganization('Google');
      manageSourcePage.findSubmitButton().should('not.be.disabled');
    });

    it('should enable Preview button when HF credentials are filled', () => {
      manageSourcePage.visitAddSource();
      manageSourcePage.fillAccessToken('test-token-123');
      manageSourcePage.fillOrganization('Google');
      manageSourcePage.findPreviewButton().should('not.be.disabled');
    });

    it('should enable Add button when all required YAML fields are filled', () => {
      manageSourcePage.visitAddSource();
      manageSourcePage.selectSourceType('yaml');
      manageSourcePage.fillSourceName('Test Source');
      manageSourcePage.fillYamlContent('test: yaml\ncontent: here');
      manageSourcePage.findSubmitButton().should('not.be.disabled');
    });

    it('should enable Preview button when YAML content is filled', () => {
      manageSourcePage.visitAddSource();
      manageSourcePage.selectSourceType('yaml');
      manageSourcePage.fillYamlContent('test: yaml\ncontent: here');
      manageSourcePage.findPreviewButton().should('not.be.disabled');
    });

    it('should show validation errors for HF fields when touched', () => {
      manageSourcePage.visitAddSource();
      manageSourcePage.findAccessTokenInput().focus().blur();
      manageSourcePage.findAccessTokenError().should('exist');
      manageSourcePage.findAccessTokenError().should('contain', 'Access token is required');

      manageSourcePage.findOrganizationInput().focus().blur();
      manageSourcePage.findOrganizationError().should('exist');
      manageSourcePage.findOrganizationError().should('contain', 'Organization is required');
    });

    it('should show validation error for YAML content when touched', () => {
      manageSourcePage.visitAddSource();
      manageSourcePage.selectSourceType('yaml');
      manageSourcePage.findYamlContentInput().focus().blur();
      manageSourcePage.findYamlContentError().should('exist');
      manageSourcePage.findYamlContentError().should('contain', 'YAML content is required');
    });

    it('should expand and collapse model visibility section', () => {
      manageSourcePage.visitAddSource();
      manageSourcePage.findAllowedModelsInput().should('not.exist');
      manageSourcePage.findExcludedModelsInput().should('not.exist');

      manageSourcePage.toggleModelVisibility();
      manageSourcePage.findAllowedModelsInput().should('exist');
      manageSourcePage.findExcludedModelsInput().should('exist');

      manageSourcePage.toggleModelVisibility();
      manageSourcePage.findAllowedModelsInput().should('not.exist');
      manageSourcePage.findExcludedModelsInput().should('not.exist');
    });

    it('should allow entering filter values', () => {
      manageSourcePage.visitAddSource();
      manageSourcePage.toggleModelVisibility();

      const allowedModels = 'model-1\nmodel-2*\nmodel-3';
      const excludedModels = 'test-model*\ndemo-model';

      manageSourcePage.fillAllowedModels(allowedModels);
      manageSourcePage.fillExcludedModels(excludedModels);

      manageSourcePage.findAllowedModelsInput().should('have.value', allowedModels);
      manageSourcePage.findExcludedModelsInput().should('have.value', excludedModels);
    });

    it('should have enable source checkbox unchecked by default', () => {
      manageSourcePage.visitAddSource();
      manageSourcePage.findEnableSourceCheckbox().should('not.be.checked');
    });

    it('should allow toggling enable source checkbox', () => {
      manageSourcePage.visitAddSource();
      manageSourcePage.findEnableSourceCheckbox().should('not.be.checked');
      manageSourcePage.toggleEnableSource();
      manageSourcePage.findEnableSourceCheckbox().should('be.checked');
      manageSourcePage.toggleEnableSource();
      manageSourcePage.findEnableSourceCheckbox().should('not.be.checked');
    });

    it('should clear validation errors when fields are filled', () => {
      manageSourcePage.visitAddSource();

      // Trigger validation errors
      manageSourcePage.findNameInput().focus().blur();
      manageSourcePage.findAccessTokenInput().focus().blur();
      manageSourcePage.findOrganizationInput().focus().blur();

      manageSourcePage.findNameError().should('exist');
      manageSourcePage.findAccessTokenError().should('exist');
      manageSourcePage.findOrganizationError().should('exist');

      // Fill fields
      manageSourcePage.fillSourceName('Test Source');
      manageSourcePage.fillAccessToken('test-token');
      manageSourcePage.fillOrganization('Google');

      // Errors should be cleared
      manageSourcePage.findNameError().should('not.exist');
      manageSourcePage.findAccessTokenError().should('not.exist');
      manageSourcePage.findOrganizationError().should('not.exist');
    });

    it('should maintain form state when switching between source types', () => {
      manageSourcePage.visitAddSource();

      // Fill name and HF fields
      manageSourcePage.fillSourceName('Test Source');
      manageSourcePage.fillAccessToken('test-token');
      manageSourcePage.fillOrganization('Google');

      // Switch to YAML
      manageSourcePage.selectSourceType('yaml');
      manageSourcePage.findYamlSection().should('exist');

      // Name should be maintained
      manageSourcePage.findNameInput().should('have.value', 'Test Source');

      // Fill YAML
      manageSourcePage.fillYamlContent('test: yaml');

      // Switch back to HF
      manageSourcePage.selectSourceType('huggingface');
      manageSourcePage.findCredentialsSection().should('exist');

      // All values should be maintained
      manageSourcePage.findNameInput().should('have.value', 'Test Source');
      manageSourcePage.findAccessTokenInput().should('have.value', 'test-token');
      manageSourcePage.findOrganizationInput().should('have.value', 'Google');
    });

    it('should dynamically update filter descriptions with organization name', () => {
      manageSourcePage.visitAddSource();
      manageSourcePage.toggleModelVisibility();

      // Before entering organization, should show generic text
      cy.contains(
        'Optionally filter which models from your source appear in the model catalog',
      ).should('exist');

      // Fill organization name
      manageSourcePage.fillOrganization('Google');

      // After entering organization, should show organization-specific text
      cy.contains(
        'Optionally filter which Google models from your source appear in the model catalog',
      ).should('exist');
      cy.contains('all Google models from the source will be visible').should('exist');

      // Change organization name
      manageSourcePage.fillOrganization('Meta');

      // Text should update to new organization
      cy.contains(
        'Optionally filter which Meta models from your source appear in the model catalog',
      ).should('exist');
      cy.contains('all Meta models from the source will be visible').should('exist');
    });

    it('should display model catalog preview panel', () => {
      manageSourcePage.visitAddSource();
      manageSourcePage.findPreviewPanel().should('exist');
      manageSourcePage.findPreviewPanelTitle().should('be.visible');
      manageSourcePage.findPreviewPanelEmptyMessage().should('be.visible');
    });

    it('should have three preview buttons', () => {
      manageSourcePage.visitAddSource();
      // One in the action group (bottom left)
      manageSourcePage.findPreviewButton().should('exist');
      // One in the preview panel header (top right)
      manageSourcePage.findPreviewButtonHeader().should('exist');
      // One in the preview panel body (center)
      manageSourcePage.findPreviewButtonPanel().should('exist');
    });

    it('should have all three preview buttons disabled by default', () => {
      manageSourcePage.visitAddSource();
      manageSourcePage.findPreviewButton().should('be.disabled');
      manageSourcePage.findPreviewButtonHeader().should('be.disabled');
      manageSourcePage.findPreviewButtonPanel().should('be.disabled');
    });

    it('should enable all three preview buttons when credentials are filled', () => {
      manageSourcePage.visitAddSource();
      manageSourcePage.fillAccessToken('test-token');
      manageSourcePage.fillOrganization('Google');
      manageSourcePage.findPreviewButton().should('not.be.disabled');
      manageSourcePage.findPreviewButtonHeader().should('not.be.disabled');
      manageSourcePage.findPreviewButtonPanel().should('not.be.disabled');
    });
  });

  describe('Manage Source Mode', () => {
    const catalogSourceId = 'test-source-id';
    const catalogSource = mockCatalogSource({
      id: catalogSourceId,
      name: 'Test Source',
    });

    beforeEach(() => {
      setupMocks([catalogSource]);
    });

    it('should display manage source page', () => {
      manageSourcePage.visitManageSource(catalogSourceId);
      manageSourcePage.findManageSourceTitle();
      manageSourcePage.findManageSourceDescription();
    });

    it('should display correct breadcrumb for manage source', () => {
      manageSourcePage.visitManageSource(catalogSourceId);
      manageSourcePage.findBreadcrumb().should('exist');
      manageSourcePage.findBreadcrumbAction().should('contain', 'Manage source');
    });

    it('should navigate back to settings from breadcrumb', () => {
      manageSourcePage.visitManageSource(catalogSourceId);
      manageSourcePage.findBreadcrumb().click({ force: true });
      modelCatalogSettings.findHeading();
    });

    it('should show Save button instead of Add button', () => {
      manageSourcePage.visitManageSource(catalogSourceId);
      manageSourcePage.findSubmitButton().should('exist');
      manageSourcePage.findSubmitButton().should('contain', 'Save');
    });
  });
});
