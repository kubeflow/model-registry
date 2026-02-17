import { TableRow } from '~/__tests__/cypress/cypress/pages/components/table';
import { TempDevFeature } from '~/app/hooks/useTempDevFeatureAvailable';

class ModelTransferJobsTableRow extends TableRow {
  findJobName() {
    return this.find().findByTestId('job-name');
  }
}

class ModelTransferJobsPage {
  visit(modelRegistryName = 'modelregistry-sample') {
    cy.visit(`/model-registry/${modelRegistryName}/model-transfer-jobs`, {
      onBeforeLoad: (win) => {
        win.localStorage.setItem(TempDevFeature.RegistryStorage, 'true');
      },
    });
    this.wait();
  }

  private wait() {
    cy.findByTestId('app-page-title').should('exist');
    cy.findByTestId('app-page-title').contains('Model transfer jobs');
  }

  findTable() {
    return cy.findByTestId('model-transfer-jobs-table');
  }

  findTableRows() {
    return this.findTable().find('tbody tr');
  }

  getRow(jobName: string) {
    return new ModelTransferJobsTableRow(
      () =>
        this.findTable()
          .find('tbody tr')
          .filter((_, row) => {
            const cell = Cypress.$(row).find('[data-testid="job-name"]');
            return cell.length > 0 && cell.text().trim() === jobName;
          })
          .first() as unknown as Cypress.Chainable<JQuery<HTMLTableRowElement>>,
    );
  }

  findDeleteModal() {
    return cy.findByTestId('delete-model-transfer-job-modal');
  }

  findDeleteModalInput() {
    return this.findDeleteModal().findByTestId('delete-modal-input');
  }

  findDeleteModalSubmitButton() {
    return this.findDeleteModal().find('button.pf-m-danger');
  }

  findDeleteModalCancelButton() {
    return this.findDeleteModal().findByRole('button', { name: 'Cancel' });
  }

  findEmptyState() {
    return cy.findByTestId('empty-model-transfer-jobs');
  }
}

export const modelTransferJobsPage = new ModelTransferJobsPage();
