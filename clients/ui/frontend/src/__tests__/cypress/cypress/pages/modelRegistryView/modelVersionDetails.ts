import { TableRow } from '~/__tests__/cypress/cypress/pages/components/table';

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

  findSourceModelFormat(subComponent: 'group' | 'edit' | 'save' | 'cancel') {
    return cy.findByTestId(`source-model-format-${subComponent}`);
  }

  findSourceModelVersion(subComponent: 'group' | 'edit' | 'save' | 'cancel') {
    return cy.findByTestId(`source-model-version-${subComponent}`);
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

  findAddPropertyButton() {
    return cy.findByTestId('add-property-button');
  }

  findAddKeyInput() {
    return cy.findByTestId('add-property-key-input');
  }

  findAddValueInput() {
    return cy.findByTestId('add-property-value-input');
  }

  findKeyEditInput(key: string) {
    return cy.findByTestId(['edit-property-key-input', key]);
  }

  findValueEditInput(value: string) {
    return cy.findByTestId(['edit-property-value-input', value]);
  }

  findSaveButton() {
    return cy.findByTestId('save-edit-button-property');
  }

  findCancelButton() {
    return cy.findByTestId('discard-edit-button-property');
  }

  findExpandControlButton() {
    return cy.findByTestId('expand-control-button');
  }

  private findTable() {
    return cy.findByTestId('properties-table');
  }

  findPropertiesTableRows() {
    return this.findTable().find('tbody tr');
  }

  getRow(name: string) {
    return new PropertyRow(() =>
      this.findTable().find(`[data-label=Key]`).contains(name).parents('tr'),
    );
  }

  findEditLabelsButton() {
    return cy.findByTestId('editable-labels-group-edit');
  }

  findAddLabelButton() {
    return cy.findByTestId('add-label-button');
  }

  findLabelInput(label: string) {
    return cy.findByTestId(`edit-label-input-${label}`);
  }

  findLabel(label: string) {
    return cy.findByTestId(`editable-label-${label}`);
  }

  findLabelErrorAlert() {
    return cy.findByTestId('label-error-alert');
  }

  findSaveLabelsButton() {
    return cy.findByTestId('editable-labels-group-save');
  }
}

class PropertyRow extends TableRow {}

export const modelVersionDetails = new ModelVersionDetails();
