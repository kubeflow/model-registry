import {
  modelCatalogSettings,
  manageSourcePage,
} from '~/__tests__/cypress/cypress/pages/modelCatalogSettings';
import { MODEL_CATALOG_API_VERSION } from '~/__tests__/cypress/cypress/support/commands/api';
import {
  mockCatalogSource,
  mockCatalogSourceList,
  mockCatalogSourceConfigList,
  mockYamlCatalogSourceConfig,
  mockHuggingFaceCatalogSourceConfig,
} from '~/__mocks__';
import type { CatalogSource, CatalogSourceConfig } from '~/app/modelCatalogTypes';

const NAMESPACE = 'kubeflow';
const userMock = {
  data: {
    userId: 'user@example.com',
    clusterAdmin: true,
  },
};

const setupMocks = (sources: CatalogSource[] = [], sourceConfigs: CatalogSourceConfig[] = []) => {
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
  cy.intercept(
    'GET',
    `/model-registry/api/${MODEL_CATALOG_API_VERSION}/settings/model_catalog/source_configs*`,
    {
      statusCode: 200,
      body: {
        data: mockCatalogSourceConfigList({
          catalogs: sourceConfigs,
        }),
      },
    },
  ).as('getCatalogSourceConfigs');
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
  const defaultYamlSource = mockYamlCatalogSourceConfig({
    id: 'default-yaml',
    name: 'Default Catalog',
    isDefault: true,
  });

  beforeEach(() => {
    setupMocks([], [defaultYamlSource]);
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

describe('Catalog Source Configs Table', () => {
  const defaultYamlSource = mockYamlCatalogSourceConfig({
    id: 'default-yaml',
    name: 'Default Catalog',
    isDefault: true,
    enabled: true,
    includedModels: [],
    excludedModels: [],
  });

  const huggingFaceSource = mockHuggingFaceCatalogSourceConfig({
    id: 'hf-google',
    name: 'HuggingFace Google',
    isDefault: false,
    enabled: true,
    allowedOrganization: 'Google',
    includedModels: ['model1', 'model2'],
  });

  const customYamlSource = mockYamlCatalogSourceConfig({
    id: 'custom-yaml',
    name: 'Custom YAML',
    isDefault: false,
    enabled: false,
    excludedModels: ['excluded-model'],
  });

  beforeEach(() => {
    setupMocks([], [defaultYamlSource, huggingFaceSource, customYamlSource]);
  });

  it('should display empty state when no source configs exist', () => {
    setupMocks([], []);
    modelCatalogSettings.visit();
    modelCatalogSettings.shouldBeEmpty();
    modelCatalogSettings.findEmptyState().should('contain', 'No catalog sources');
  });

  it('should display table with source configs', () => {
    modelCatalogSettings.visit();
    modelCatalogSettings.shouldHaveSourceConfigs();
    modelCatalogSettings.findRows().should('have.length', 3);
  });

  it('should render table column headers correctly', () => {
    modelCatalogSettings.visit();
    modelCatalogSettings.findTable().should('be.visible');
    modelCatalogSettings.findTable().contains('th', 'Name').should('be.visible');
    modelCatalogSettings.findTable().contains('th', 'Organization').should('be.visible');
    modelCatalogSettings.findTable().contains('th', 'Model visibility').should('be.visible');
    modelCatalogSettings.findTable().contains('th', 'Source type').should('be.visible');
    modelCatalogSettings.findTable().contains('th', 'Enable').should('be.visible');
    modelCatalogSettings.findTable().contains('th', 'Validation status').should('be.visible');
  });

  describe('Table row rendering', () => {
    it('should render default YAML source correctly', () => {
      modelCatalogSettings.visit();
      const row = modelCatalogSettings.getRow('Default Catalog');
      row.findName().should('be.visible').and('contain', 'Default Catalog');
      row.shouldHaveOrganization('-');
      row.shouldHaveModelVisibility('Unfiltered');
      row.shouldHaveSourceType('YAML file');
      row.shouldHaveEnableToggle(false); // Default sources don't have toggle
    });

    it('should render Hugging Face source correctly', () => {
      modelCatalogSettings.visit();
      const row = modelCatalogSettings.getRow('HuggingFace Google');
      row.findName().should('be.visible').and('contain', 'HuggingFace Google');
      row.shouldHaveOrganization('Google');
      row.shouldHaveModelVisibility('Filtered');
      row.shouldHaveSourceType('Hugging Face');
      row.shouldHaveEnableToggle(true);
      row.shouldHaveEnableState(true);
    });

    it('should render custom YAML source correctly', () => {
      modelCatalogSettings.visit();
      const row = modelCatalogSettings.getRow('Custom YAML');
      row.findName().should('be.visible').and('contain', 'Custom YAML');
      row.shouldHaveOrganization('-');
      row.shouldHaveModelVisibility('Filtered');
      row.shouldHaveSourceType('YAML file');
      row.shouldHaveEnableToggle(true);
      row.shouldHaveEnableState(false);
    });
  });

  describe('Enable toggle functionality', () => {
    it('should show alert when enable toggle is clicked', () => {
      modelCatalogSettings.visit();
      const row = modelCatalogSettings.getRow('HuggingFace Google');
      row.findName().should('be.visible');
      row.findEnableToggle().should('exist').and('be.checked');

      cy.window().then((win) => {
        cy.stub(win, 'alert').as('windowAlert');
      });

      row.toggleEnable();

      cy.get('@windowAlert').should(
        'have.been.calledWith',
        'Toggle clicked! "HuggingFace Google" will be disabled when functionality is implemented.',
      );
    });

    it('should not show toggle for default sources', () => {
      modelCatalogSettings.visit();
      const row = modelCatalogSettings.getRow('Default Catalog');
      row.findName().should('be.visible');
      row.shouldHaveEnableToggle(false);
    });
  });

  describe('Manage source button', () => {
    it('should navigate to manage source page when button is clicked', () => {
      modelCatalogSettings.visit();
      const row = modelCatalogSettings.getRow('HuggingFace Google');
      row.findName().should('be.visible');
      row.findManageSourceButton().should('be.visible').click();
      cy.url().should('include', '/model-catalog-settings/manage-source/hf-google');
      manageSourcePage.findManageSourceTitle();
    });

    it('should navigate to correct manage source page for each row', () => {
      modelCatalogSettings.visit();
      const customRow = modelCatalogSettings.getRow('Custom YAML');
      customRow.findName().should('be.visible');
      customRow.findManageSourceButton().should('be.visible').click();
      cy.url().should('include', '/model-catalog-settings/manage-source/custom-yaml');
    });
  });

  describe('Kebab menu actions', () => {
    it('should show delete action for non-default sources', () => {
      modelCatalogSettings.visit();
      const row = modelCatalogSettings.getRow('HuggingFace Google');
      row.findName().should('be.visible');
      row.findKebab().should('be.visible').click();
      cy.findByRole('menuitem', { name: 'Delete source' })
        .should('be.visible')
        .and('not.be.disabled');
    });

    it('should disable delete action for default sources', () => {
      modelCatalogSettings.visit();
      const row = modelCatalogSettings.getRow('Default Catalog');
      row.findName().should('be.visible');
      row.findKebab().should('be.visible').click();
      cy.findByRole('menuitem', { name: 'Delete source' }).should('be.visible').and('be.disabled');
    });
  });

  describe('Model visibility badges', () => {
    it('should show "Filtered" badge when source has included models', () => {
      modelCatalogSettings.visit();
      const row = modelCatalogSettings.getRow('HuggingFace Google');
      row.findName().should('be.visible');
      row.findModelVisibility().should('be.visible').and('contain', 'Filtered');
      row
        .findModelVisibility()
        .find('[data-testid*="model-visibility-filtered"]')
        .should('be.visible');
    });

    it('should show "Filtered" badge when source has excluded models', () => {
      modelCatalogSettings.visit();
      const row = modelCatalogSettings.getRow('Custom YAML');
      row.findName().should('be.visible');
      row.findModelVisibility().should('be.visible').and('contain', 'Filtered');
    });

    it('should show "Unfiltered" badge when source has no filters', () => {
      modelCatalogSettings.visit();
      const row = modelCatalogSettings.getRow('Default Catalog');
      row.findName().should('be.visible');
      row.findModelVisibility().should('be.visible').and('contain', 'Unfiltered');
      row
        .findModelVisibility()
        .find('[data-testid*="model-visibility-unfiltered"]')
        .should('be.visible');
    });
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
