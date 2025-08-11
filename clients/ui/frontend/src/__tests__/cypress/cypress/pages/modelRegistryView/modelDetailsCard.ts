class ModelDetailsCard {
  // Model Details Card selectors and methods
  findOwner() {
    return cy.findByTestId('registered-model-owner');
  }

  // Description section
  findDescriptionEditButton() {
    return cy.findByTestId('model-description-edit');
  }

  findDescriptionTextArea() {
    return cy.findByTestId('model-description-input');
  }

  findDescriptionSaveButton() {
    return cy.findByTestId('model-description-save');
  }

  // Properties section
  findPropertiesTable() {
    return cy.findByTestId('properties-table');
  }

  findAddPropertyButton() {
    return cy.findByTestId('add-property-button');
  }

  findAddPropertyKeyInput() {
    return cy.findByTestId('add-property-key-input');
  }

  findAddPropertyValueInput() {
    return cy.findByTestId('add-property-value-input');
  }

  findSavePropertyButton() {
    return cy.findByTestId('save-edit-button-property');
  }

  findExpandControlButton() {
    return cy.findByTestId('expand-control-button');
  }

  findToggleButton() {
    return cy.findByTestId('model-details-card-toggle-button');
  }
}

export const modelDetailsCard = new ModelDetailsCard();
