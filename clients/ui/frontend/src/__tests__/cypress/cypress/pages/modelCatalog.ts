class ModelCatalog {
  visit() {
    cy.visit('/model-catalog');
  }

  findModelCatalogCards() {
    return cy.findAllByTestId('model-catalog-card');
  }

  findFirstModelCatalogCard() {
    return this.findModelCatalogCards().first().should('be.visible');
  }

  findModelCatalogDetailLink() {
    return cy.findAllByTestId('model-catalog-detail-link');
  }

  findModelCatalogDescription() {
    return cy.findByTestId('model-catalog-card-description');
  }

  findSourceLabel() {
    return cy.get('.pf-v6-c-label');
  }

  findModelLogo() {
    return cy.get('img[alt="model logo"]');
  }

  findVersionIcon() {
    return cy.get('.pf-v6-c-icon');
  }

  findFrameworkLabel() {
    return cy.contains('PyTorch');
  }

  findTaskLabel() {
    return cy.contains('text-generation');
  }

  findLicenseLabel() {
    return cy.contains('apache-2.0');
  }

  findLabBaseLabel() {
    return cy.contains('lab-base');
  }

  findLoadingState() {
    return cy.contains('Loading model catalog...');
  }

  findPageTitle() {
    return cy.contains('Model Catalog');
  }

  findPageDescription() {
    return cy.contains('Discover models that are available for your organization');
  }

  // Details page helpers
  findBreadcrumb() {
    return cy.contains('Model catalog');
  }

  findDetailsProviderText() {
    return cy.contains('Provided by');
  }

  findDetailsDescription() {
    return cy.findByTestId('model-long-description');
  }
}

export const modelCatalog = new ModelCatalog();
