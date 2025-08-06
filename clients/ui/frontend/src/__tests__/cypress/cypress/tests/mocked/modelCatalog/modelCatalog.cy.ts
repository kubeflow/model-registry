import { modelCatalog } from '~/__tests__/cypress/cypress/pages/modelCatalog';

describe('Model Catalog Page', () => {
  beforeEach(() => {
    cy.visit('/model-catalog');
  });

  it('should display loading state initially', () => {
    modelCatalog.findLoadingState().should('be.visible');
  });

  it('should display model catalog content when data is loaded', () => {
    modelCatalog.findLoadingState().should('not.exist');
    modelCatalog.findPageTitle().should('be.visible');
    modelCatalog.findPageDescription().should('be.visible');
    modelCatalog.findModelCatalogCards().should('have.length.at.least', 1);
  });
});
