import { modelCatalog } from '~/__tests__/cypress/cypress/pages/modelCatalog';

describe('Model Catalog Details Page', () => {
  beforeEach(() => {
    cy.visit('/model-catalog');
  });

  it('navigates to details and shows header, breadcrumb and description', () => {
    modelCatalog.findLoadingState().should('not.exist');
    modelCatalog.findModelCatalogDetailLink().first().click();
    modelCatalog.findBreadcrumb().should('be.visible');
    modelCatalog.findDetailsProviderText().should('be.visible');
    modelCatalog.findDetailsDescription().should('exist');
  });
});
