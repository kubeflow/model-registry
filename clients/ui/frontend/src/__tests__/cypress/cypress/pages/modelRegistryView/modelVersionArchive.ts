import { TableRow } from '~/__tests__/cypress/cypress/pages/components/table';
import { Modal } from '~/__tests__/cypress/cypress/pages/components/Modal';

class ArchiveVersionTableRow extends TableRow {
  findName() {
    return this.find().findByTestId('model-version-name');
  }

  findDescription() {
    return this.find().findByTestId('model-version-description');
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

class RestoreVersionModal extends Modal {
  constructor() {
    super('Restore version?');
  }

  findRestoreButton() {
    return cy.findByTestId('modal-submit-button');
  }
}

class ArchiveVersionModal extends Modal {
  constructor() {
    super('Archive version?');
  }

  findArchiveButton() {
    return cy.findByTestId('modal-submit-button');
  }

  findModalTextInput() {
    return cy.findByTestId('confirm-archive-input');
  }
}

class ModelVersionArchive {
  private wait() {
    cy.findByTestId('app-page-title').should('exist');
    cy.testA11y();
  }

  visit() {
    const rmId = '1';
    const preferredModelRegistry = 'modelregistry-sample';
    cy.visit(`/model-registry/${preferredModelRegistry}/registeredModels/${rmId}/versions/archive`);
    this.wait();
  }

  visitArchiveVersionDetail() {
    const mvId = '2';
    const rmId = '1';
    const preferredModelRegistry = 'modelregistry-sample';
    cy.visit(
      `/model-registry/${preferredModelRegistry}/registeredModels/${rmId}/versions/archive/${mvId}`,
    );
  }

  visitModelVersionList() {
    const rmId = '1';
    const preferredModelRegistry = 'modelregistry-sample';
    cy.visit(`/model-registry/${preferredModelRegistry}/registeredModels/${rmId}/versions`);
    this.wait();
  }

  visitModelVersionDetails() {
    const mvId = '3';
    const rmId = '1';
    const preferredModelRegistry = 'modelregistry-sample';
    cy.visit(`/model-registry/${preferredModelRegistry}/registeredModels/${rmId}/versions/${mvId}`);
    this.wait();
  }

  findModelVersionsTableKebab() {
    return cy.findByTestId('model-versions-table-kebab-action');
  }

  shouldArchiveVersionsEmpty() {
    cy.findByTestId('empty-archive-state').should('exist');
  }

  findArchiveVersionBreadcrumbItem() {
    return cy.findByTestId('archive-version-page-breadcrumb');
  }

  findVersionDetailsTab() {
    return cy.findByTestId('model-versions-details-tab');
  }

  findVersionDeploymentTab() {
    return cy.findByTestId('deployments-tab');
  }

  findArchiveVersionTable() {
    return cy.findByTestId('model-versions-archive-table');
  }

  findArchiveVersionsTableRows() {
    return this.findArchiveVersionTable().find('tbody tr');
  }

  findRestoreButton() {
    return cy.findByTestId('restore-button');
  }

  getRow(name: string) {
    return new ArchiveVersionTableRow(() =>
      this.findArchiveVersionTable()
        .find(`[data-label="Version name"]`)
        .contains(name)
        .parents('tr'),
    );
  }

  findModelVersionsDetailsHeaderAction() {
    return cy.findByTestId('model-version-details-action-button');
  }
}

export const modelVersionArchive = new ModelVersionArchive();
export const restoreVersionModal = new RestoreVersionModal();
export const archiveVersionModal = new ArchiveVersionModal();
