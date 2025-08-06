import { modelCatalog } from '~/__tests__/cypress/cypress/pages/modelCatalog';

describe('ModelCatalogCard Component', () => {
  beforeEach(() => {
    cy.visit('/model-catalog');
    modelCatalog.findLoadingState().should('not.exist');
    modelCatalog.findModelCatalogCards().should('have.length.at.least', 1);
  });

  describe('Card Layout and Content', () => {
    it('should render all cards from the mock data', () => {
      modelCatalog.findModelCatalogCards().should('have.length.at.least', 1);
    });

    it('should display correct source labels', () => {
      modelCatalog.findFirstModelCatalogCard().within(() => {
        modelCatalog.findSourceLabel().should('contain.text', 'Red Hat');
      });
    });

    it('should handle cards with logos', () => {
      modelCatalog.findFirstModelCatalogCard().within(() => {
        modelCatalog
          .findModelLogo()
          .should('exist')
          .and('have.attr', 'src')
          .and('include', 'data:image/svg+xml;base64');
      });
    });
  });

  describe('Version Tag Display', () => {
    it('should extract and display version tags correctly', () => {
      modelCatalog.findFirstModelCatalogCard().within(() => {
        modelCatalog.findVersionIcon().should('exist');
        cy.contains('1.4.0').should('exist');
      });
    });
  });

  describe('Description Handling', () => {
    it('should display model descriptions', () => {
      modelCatalog.findFirstModelCatalogCard().within(() => {
        modelCatalog
          .findModelCatalogDescription()
          .should('contain.text', 'Base model for customizing and fine-tuning');
      });
    });
  });

  describe('Navigation and Interaction', () => {
    it('should show all model metadata correctly', () => {
      modelCatalog.findFirstModelCatalogCard().within(() => {
        modelCatalog.findModelCatalogDetailLink().should('contain.text', 'granite-7b-starter');

        modelCatalog.findFrameworkLabel().should('exist');
        modelCatalog.findTaskLabel().should('exist');
        modelCatalog.findLicenseLabel().should('exist');
        modelCatalog.findLabBaseLabel().should('exist');
      });
    });
  });
});
