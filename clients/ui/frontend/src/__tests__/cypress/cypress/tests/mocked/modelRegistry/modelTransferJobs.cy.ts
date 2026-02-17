import { mockModArchResponse } from 'mod-arch-core';
import { mockModelTransferJob, mockModelTransferJobList } from '~/__mocks__/mockModelTransferJob';
import { mockModelRegistry } from '~/__mocks__/mockModelRegistry';
import { modelTransferJobsPage } from '~/__tests__/cypress/cypress/pages/modelRegistryView/modelTransferJobs';
import { MODEL_REGISTRY_API_VERSION } from '~/__tests__/cypress/cypress/support/commands/api';

const modelRegistryName = 'modelregistry-sample';

const jobList = mockModelTransferJobList({
  items: [
    mockModelTransferJob({ id: 'job-to-delete', name: 'job-to-delete' }),
    mockModelTransferJob({ id: 'job-to-keep', name: 'job-to-keep' }),
  ],
  size: 2,
  pageSize: 10,
  nextPageToken: '',
});

function visitAndWaitForJobsTable() {
  modelTransferJobsPage.visit(modelRegistryName);
  modelTransferJobsPage.findTableRows().should('have.length', 2);
}

const setupIntercepts = () => {
  cy.intercept(
    'GET',
    `**/api/${MODEL_REGISTRY_API_VERSION}/model_registry*`,
    mockModArchResponse([mockModelRegistry({ name: modelRegistryName })]),
  );

  cy.intercept(
    'GET',
    `**/api/${MODEL_REGISTRY_API_VERSION}/model_registry/${modelRegistryName}/model_transfer_jobs*`,
    mockModArchResponse(jobList),
  );

  cy.intercept(
    'DELETE',
    `**/api/${MODEL_REGISTRY_API_VERSION}/model_registry/${modelRegistryName}/model_transfer_jobs/job-to-delete*`,
    { statusCode: 200, body: { status: 'deleted' } },
  );
};

describe('Model transfer jobs', () => {
  beforeEach(() => {
    setupIntercepts();
  });

  it('should delete a model transfer job via kebab action and confirmation modal', () => {
    visitAndWaitForJobsTable();

    modelTransferJobsPage.getRow('job-to-delete').findKebabAction('Delete').click();

    modelTransferJobsPage.findDeleteModal().should('be.visible');
    modelTransferJobsPage.findDeleteModal().contains('Delete model transfer job?').should('exist');
    modelTransferJobsPage.findDeleteModal().contains('job-to-delete').should('exist');

    modelTransferJobsPage.findDeleteModalSubmitButton().should('be.disabled');
    modelTransferJobsPage.findDeleteModalInput().type('job-to-delete');
    modelTransferJobsPage.findDeleteModalSubmitButton().should('be.enabled').click();

    modelTransferJobsPage.findDeleteModal().should('not.exist');
  });

  it('should close delete modal on Cancel without deleting', () => {
    visitAndWaitForJobsTable();

    modelTransferJobsPage.getRow('job-to-delete').findKebabAction('Delete').click();

    modelTransferJobsPage.findDeleteModal().should('be.visible');
    modelTransferJobsPage.findDeleteModalCancelButton().click();

    modelTransferJobsPage.findDeleteModal().should('not.exist');
    modelTransferJobsPage.findTableRows().should('have.length', 2);
  });
});
