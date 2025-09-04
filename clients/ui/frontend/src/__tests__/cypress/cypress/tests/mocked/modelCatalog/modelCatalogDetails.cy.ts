import { modelCatalog } from '~/__tests__/cypress/cypress/pages/modelCatalog';
import { mockModelRegistry } from '~/__mocks__/mockModelRegistry';

describe('Model Catalog Details Page', () => {
  beforeEach(() => {
    // Mock model registries for register button functionality
    cy.intercept('GET /api/model-registries', [
      mockModelRegistry({ name: 'modelregistry-sample' }),
    ]);

    cy.visit('/model-catalog');
  });

  it('navigates to details and shows header, breadcrumb and description', () => {
    modelCatalog.findLoadingState().should('not.exist');
    modelCatalog.findModelCatalogDetailLink().first().click();
    modelCatalog.findBreadcrumb().should('be.visible');
    modelCatalog.findDetailsProviderText().should('be.visible');
    modelCatalog.findDetailsDescription().should('exist');
  });

  it('shows register button when model registries are available', () => {
    modelCatalog.findLoadingState().should('not.exist');
    modelCatalog.findModelCatalogDetailLink().first().click();

    // Should show register button in header actions
    cy.findByTestId('register-model-button').should('be.visible');
    cy.findByTestId('register-model-button').should('contain', 'Register model');
  });

  it('navigates to register form when clicking register button', () => {
    modelCatalog.findLoadingState().should('not.exist');
    modelCatalog.findModelCatalogDetailLink().first().click();

    cy.findByTestId('register-model-button').click();

    // Should navigate to register page
    cy.url().should('include', '/register');
  });
});
