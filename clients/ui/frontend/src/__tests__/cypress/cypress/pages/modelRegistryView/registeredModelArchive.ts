import { TableRow } from '~/__tests__/cypress/cypress/pages/components/table';
import { Modal } from '~/__tests__/cypress/cypress/pages/components/Modal';

class ArchiveModelTableRow extends TableRow {
  findName() {
    return this.find().findByTestId('model-name');
  }

  findDescription() {
    return this.find().findByTestId('description');
  }

  findLabelPopoverText() {
    return this.find().findByTestId('popover-label-text');
  }

  findLabelModalText() {
    return this.find().findByTestId('modal-label-text');
  }

  shouldContainsPopoverLabels(labels: string[]) {
    cy.findByTestId('popover-label-group').within(() => labels.map((label) => cy.contains(label)));
    return this;
  }
}

class RestoreModelModal extends Modal {
  constructor() {
    super('Restore model?');
  }

  findRestoreButton() {
    return cy.findByTestId('modal-submit-button');
  }
}

class ArchiveModelModal extends Modal {
  constructor() {
    super('Archive model?');
  }

  findArchiveButton() {
    return cy.findByTestId('modal-submit-button');
  }

  findModalTextInput() {
    return cy.findByTestId('confirm-archive-input');
  }
}

class ModelArchive {
  private wait() {
    cy.findByTestId('app-page-title').should('exist');
    cy.testA11y();
  }

  visit() {
    const preferredModelRegistry = 'modelregistry-sample';
    cy.visit(`/model-registry/${preferredModelRegistry}/registeredModels/archive`);
    this.wait();
  }

  visitArchiveModelDetail() {
    const rmId = '2';
    const preferredModelRegistry = 'modelregistry-sample';
    cy.visit(`/model-registry/${preferredModelRegistry}/registeredModels/archive/${rmId}`);
  }

  visitArchiveModelVersionList() {
    const rmId = '2';
    const preferredModelRegistry = 'modelregistry-sample';
    cy.visit(`/model-registry/${preferredModelRegistry}/registeredModels/archive/${rmId}/versions`);
  }

  visitModelList() {
    cy.visit('/model-registry/modelregistry-sample');
    this.wait();
  }

  visitModelDetails() {
    const rmId = '2';
    const preferredModelRegistry = 'modelregistry-sample';
    cy.visit(`/model-registry/${preferredModelRegistry}/registeredModels/${rmId}`);
    this.wait();
  }

  findTableKebabMenu() {
    return cy.findByTestId('registered-models-table-kebab-action');
  }

  shouldArchiveVersionsEmpty() {
    cy.findByTestId('empty-archive-model-state').should('exist');
  }

  findArchiveModelBreadcrumbItem() {
    return cy.findByTestId('archive-model-page-breadcrumb');
  }

  findRegisteredModelsArchiveTableHeaderButton(name: string) {
    return this.findArchiveModelTable().find('thead').findByRole('button', { name });
  }

  findArchiveModelTable() {
    return cy.findByTestId('registered-models-archive-table');
  }

  findArchiveModelsTableRows() {
    return this.findArchiveModelTable().find('tbody tr');
  }

  findRestoreButton() {
    return cy.findByTestId('restore-button');
  }

  getRow(name: string) {
    return new ArchiveModelTableRow(() =>
      this.findArchiveModelTable().find(`[data-label="Model name"]`).contains(name).parents('tr'),
    );
  }

  findModelVersionsDetailsHeaderAction() {
    return cy.findByTestId('model-version-action-toggle');
  }
}

export const registeredModelArchive = new ModelArchive();
export const restoreModelModal = new RestoreModelModal();
export const archiveModelModal = new ArchiveModelModal();
