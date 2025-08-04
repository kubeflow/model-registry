describe('ModelCatalogCard Component', () => {
  beforeEach(() => {
    cy.visit('/model-catalog');
    // Wait for cards to be rendered
    cy.get('[data-testid="model-catalog-card"]', { timeout: 10000 }).should(
      'have.length.at.least',
      1,
    );
  });

  describe('Card Layout and Content', () => {
    it('should render all cards from the mock data', () => {
      cy.get('[data-testid="model-catalog-card"]').should('have.length.at.least', 1);
    });

    it('should display correct source labels', () => {
      cy.get('[data-testid="model-catalog-card"]')
        .first()
        .find('.pf-v6-c-label')
        .should('contain.text', 'Red Hat');
    });

    it('should handle cards with logos', () => {
      cy.get('[data-testid="model-catalog-card"]')
        .first()
        .find('img[alt="model logo"]')
        .should('exist')
        .and('have.attr', 'src')
        .and('include', 'data:image/svg+xml;base64');
    });
  });

  describe('Version Tag Display', () => {
    it('should extract and display version tags correctly', () => {
      cy.get('[data-testid="model-catalog-card"]')
        .first()
        .within(() => {
          cy.get('.pf-v6-c-icon').should('exist');
          cy.contains('1.4.0');
        });
    });
  });

  describe('Description Handling', () => {
    it('should display model descriptions', () => {
      cy.get('[data-testid="model-catalog-card"]')
        .first()
        .find('[data-testid="model-catalog-card-description"]')
        .should('contain.text', 'Base model for customizing and fine-tuning');
    });
  });

  describe('Navigation and Interaction', () => {
    it('should show all model metadata correctly', () => {
      cy.get('[data-testid="model-catalog-card"]')
        .first()
        .within(() => {
          // Model name (using displayName from mock data)
          cy.get('[data-testid="model-catalog-detail-link"]').should(
            'contain.text',
            'granite-7b-starter',
          );

          // Framework and task labels
          cy.contains('PyTorch');
          cy.contains('text-generation');
          cy.contains('apache-2.0');
          cy.contains('lab-base');
        });
    });
  });
});
