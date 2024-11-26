class ModelVersionDetails {
  visit() {
    const preferredModelRegistry = 'modelregistry-sample';
    const rmId = '1';
    const mvId = '1';
    cy.visit(`/model-registry/${preferredModelRegistry}/registeredModels/${rmId}/versions/${mvId}`);
    this.wait();
  }

  private wait() {
    cy.findByTestId('app-page-title').should('exist');
    cy.testA11y();
  }

  findVersionId() {
    return cy.findByTestId('model-version-id');
  }

  findDescription() {
    return cy.findByTestId('model-version-description');
  }

  findMoreLabelsButton() {
    return cy.findByTestId('label-group').find('button');
  }

  findStorageURI() {
    return cy.findByTestId('storage-uri');
  }

  findStorageEndpoint() {
    return cy.findByTestId('storage-endpoint');
  }

  findStorageRegion() {
    return cy.findByTestId('storage-region');
  }

  findStorageBucket() {
    return cy.findByTestId('storage-bucket');
  }

  findStoragePath() {
    return cy.findByTestId('storage-path');
  }

  shouldContainsModalLabels(labels: string[]) {
    cy.findByTestId('label-group').within(() => labels.map((label) => cy.contains(label)));
    return this;
  }

  findModelVersionDropdownButton() {
    return cy.findByTestId('model-version-toggle-button');
  }

  findModelVersionDropdownSearch() {
    return cy.findByTestId('search-input');
  }

  findModelVersionDropdownItem(name: string) {
    return cy.findByTestId('model-version-selector-list').find('li').contains(name);
  }

  findDetailsTab() {
    return cy.findByTestId('model-versions-details-tab');
  }

  findRegisteredDeploymentsTab() {
    return cy.findByTestId('deployments-tab');
  }
}

export const modelVersionDetails = new ModelVersionDetails();
