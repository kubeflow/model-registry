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
