import { TableRow } from '~/__tests__/cypress/cypress/pages/components/table';
import { Modal } from '~/__tests__/cypress/cypress/pages/components/Modal';

class ExpandedModelDetailsCardPropertyRow extends TableRow {
  findSaveButton() {
    return cy.findByTestId('save-edit-button-property');
  }
}

class DeletePropertyModal extends Modal {
  constructor() {
    super('Delete property from all model versions?');
  }

  find() {
    return cy.findByTestId('delete-property-modal');
  }

  findConfirmButton() {
    return this.findFooter().findByTestId('delete-property-modal-confirm');
  }
}

class ModelDetailsExpandedCard {
  findExpandedButton() {
    return cy.findByTestId('model-details-card-toggle-button');
  }

  find() {
    return cy.findByTestId('model-details-card-expandable-content');
  }

  findLabelEditButton() {
    return this.find().findByTestId('editable-labels-group-edit');
  }

  findLabelSaveButton() {
    return this.find().findByTestId('editable-labels-group-save');
  }

  findDescriptionEditButton() {
    return this.find().findByTestId('model-description-edit');
  }

  findDescriptionSaveButton() {
    return this.find().findByTestId('model-description-save');
  }

  findAlert() {
    return cy.findByTestId('edit-alert');
  }

  findAddPropertyButton() {
    return this.find().findByTestId('add-property-button');
  }

  findTable() {
    return this.find().findByTestId('properties-table');
  }

  findPropertiesExpandableButton() {
    return this.find().findByTestId('properties-expandable-section').findByRole('button');
  }

  getRow(name: string) {
    return new ExpandedModelDetailsCardPropertyRow(() =>
      this.findTable().find(`[data-label=Key]`).contains(name).parents('tr'),
    );
  }
}

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
export const modelDetailsExpandedCard = new ModelDetailsExpandedCard();
export const deletePropertyModal = new DeletePropertyModal();
