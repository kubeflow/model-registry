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
  // Intercept the catalog source configs endpoint
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
    modelCatalogSettings.findTable().contains('th', 'Name');
    modelCatalogSettings.findTable().contains('th', 'Organization');
    modelCatalogSettings.findTable().contains('th', 'Model visibility');
    modelCatalogSettings.findTable().contains('th', 'Source type');
    modelCatalogSettings.findTable().contains('th', 'Enable');
    modelCatalogSettings.findTable().contains('th', 'Validation status');
  });

  describe('Table row rendering', () => {
    it('should render default YAML source correctly', () => {
      modelCatalogSettings.visit();
      const row = modelCatalogSettings.getRow('Default Catalog');
      row.findName().should('contain', 'Default Catalog');
      row.shouldHaveOrganization('-');
      row.shouldHaveModelVisibility('Unfiltered');
      row.shouldHaveSourceType('YAML file');
      row.shouldHaveEnableToggle(false); // Default sources don't have toggle
    });

    it('should render Hugging Face source correctly', () => {
      modelCatalogSettings.visit();
      const row = modelCatalogSettings.getRow('HuggingFace Google');
      row.findName().should('contain', 'HuggingFace Google');
      row.shouldHaveOrganization('Google');
      row.shouldHaveModelVisibility('Filtered');
      row.shouldHaveSourceType('Hugging Face');
      row.shouldHaveEnableToggle(true);
      row.shouldHaveEnableState(true);
    });

    it('should render custom YAML source correctly', () => {
      modelCatalogSettings.visit();
      const row = modelCatalogSettings.getRow('Custom YAML');
      row.findName().should('contain', 'Custom YAML');
      row.shouldHaveOrganization('-');
      row.shouldHaveModelVisibility('Filtered');
      row.shouldHaveSourceType('YAML file');
      row.shouldHaveEnableToggle(true);
      row.shouldHaveEnableState(false);
    });
  });

  describe('Enable toggle functionality', () => {
    it('should show notification when enable toggle is clicked', () => {
      modelCatalogSettings.visit();
      const row = modelCatalogSettings.getRow('HuggingFace Google');
      row.toggleEnable();
      // Check for notification toast
      cy.get('.pf-v5-c-alert-group').should('exist');
      cy.get('.pf-v5-c-alert').should('contain', 'Toggle disabled');
    });

    it('should not show toggle for default sources', () => {
      modelCatalogSettings.visit();
      const row = modelCatalogSettings.getRow('Default Catalog');
      row.shouldHaveEnableToggle(false);
    });
  });

  describe('Manage source button', () => {
    it('should navigate to manage source page when button is clicked', () => {
      modelCatalogSettings.visit();
      const row = modelCatalogSettings.getRow('HuggingFace Google');
      row.findManageSourceButton().click();
      cy.url().should('include', '/model-catalog-settings/manage-source/hf-google');
      manageSourcePage.findManageSourceTitle();
    });

    it('should navigate to correct manage source page for each row', () => {
      modelCatalogSettings.visit();
      const customRow = modelCatalogSettings.getRow('Custom YAML');
      customRow.findManageSourceButton().click();
      cy.url().should('include', '/model-catalog-settings/manage-source/custom-yaml');
    });
  });

  describe('Kebab menu actions', () => {
    it('should show delete action for non-default sources', () => {
      modelCatalogSettings.visit();
      const row = modelCatalogSettings.getRow('HuggingFace Google');
      row.findKebab().click();
      cy.findByRole('menuitem', { name: 'Delete source' }).should('exist');
      cy.findByRole('menuitem', { name: 'Delete source' }).should('not.be.disabled');
    });

    it('should disable delete action for default sources', () => {
      modelCatalogSettings.visit();
      const row = modelCatalogSettings.getRow('Default Catalog');
      row.findKebab().click();
      cy.findByRole('menuitem', { name: 'Delete source' }).should('be.disabled');
    });
  });

  describe('Model visibility badges', () => {
    it('should show "Filtered" badge when source has included models', () => {
      modelCatalogSettings.visit();
      const row = modelCatalogSettings.getRow('HuggingFace Google');
      row.findModelVisibility().should('contain', 'Filtered');
      row.findModelVisibility().find('.pf-v5-c-label').should('have.class', 'pf-m-blue');
    });

    it('should show "Filtered" badge when source has excluded models', () => {
      modelCatalogSettings.visit();
      const row = modelCatalogSettings.getRow('Custom YAML');
      row.findModelVisibility().should('contain', 'Filtered');
    });

    it('should show "Unfiltered" badge when source has no filters', () => {
      modelCatalogSettings.visit();
      const row = modelCatalogSettings.getRow('Default Catalog');
      row.findModelVisibility().should('contain', 'Unfiltered');
      row.findModelVisibility().find('.pf-v5-c-label').should('have.class', 'pf-m-grey');
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
  });
});
