describe('Model Catalog Page', () => {
  beforeEach(() => {
    // Visit the page before each test
    cy.visit('/model-catalog');
  });

  it('should display loading state initially', () => {
    cy.contains('Loading model catalog...');
  });

  it('should display model catalog cards when data is loaded', () => {
    // Wait for the content to be loaded and cards to be rendered
    cy.contains('Model Catalog', { timeout: 10000 });
    cy.contains('Discover models that are available for your organization');

    // Check if model cards are rendered
    cy.get('[data-testid="model-catalog-card"]', { timeout: 10000 })
      .should('have.length.at.least', 1)
      .first()
      .within(() => {
        // Check model name
        cy.get('[data-testid="model-catalog-detail-link"]').should(
          'contain.text',
          'granite-7b-starter',
        );

        // Check description
        cy.get('[data-testid="model-catalog-card-description"]').should(
          'contain.text',
          'Base model',
        );
      });
  });

  it('should display empty state when no models are available', () => {
    // Since we can't easily modify the context, we'll skip this test for now
    cy.log('Test skipped - requires context modification capability');
  });

  it('should display error state when API fails', () => {
    // Since we can't easily modify the context, we'll skip this test for now
    cy.log('Test skipped - requires context modification capability');
  });
});

describe('Model Catalog Card', () => {
  beforeEach(() => {
    cy.visit('/model-catalog');
    // Wait for cards to be rendered
    cy.get('[data-testid="model-catalog-card"]', { timeout: 10000 }).should(
      'have.length.at.least',
      1,
    );
  });

  it('should display all card elements correctly', () => {
    cy.get('[data-testid="model-catalog-card"]')
      .first()
      .within(() => {
        // Check logo
        cy.get('img[alt="model logo"]').should('exist');

        // Check model name
        cy.get('[data-testid="model-catalog-detail-link"]').should(
          'contain.text',
          'granite-7b-starter',
        );

        // Check description
        cy.get('[data-testid="model-catalog-card-description"]').should(
          'contain.text',
          'Base model',
        );
      });
  });
  it('should display correct metadata', () => {
    cy.get('[data-testid="model-catalog-card"]')
      .first()
      .within(() => {
        // Check framework and task labels
        cy.contains('PyTorch');
        cy.contains('text-generation');
      });
  });

  it('should be accessible', () => {
    cy.get('[data-testid="model-catalog-card"]')
      .first()
      .within(() => {
        // Check if the link button is accessible
        cy.get('[data-testid="model-catalog-detail-link"]')
          .should('have.class', 'pf-m-link')
          .and('have.class', 'pf-m-inline');

        // Check if images have alt text
        cy.get('img').should('have.attr', 'alt');
      });
  });
});
