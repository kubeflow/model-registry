class ModelVersionsCard {
  findNoVersionsText() {
    return cy.findByTestId('no-versions-text');
  }

  findModelVersion(id: string) {
    return cy.findByTestId(`model-version-${id}`);
  }

  findModelVersionProperty(id: string, property: string) {
    return cy.findByTestId(`model-version-${id}-property-${property}`);
  }

  findViewAllVersionsLink() {
    return cy.findByTestId('versions-route-link');
  }

  findModelVersionLink(id: string) {
    return cy.findByTestId(`model-version-${id}-link`);
  }
}

export const modelVersionsCard = new ModelVersionsCard();
