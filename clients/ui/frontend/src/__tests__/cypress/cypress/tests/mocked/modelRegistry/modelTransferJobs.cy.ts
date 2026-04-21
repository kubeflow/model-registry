import { mockModArchResponse } from 'mod-arch-core';
import { ModelTransferJobStatus, ModelTransferJobUploadIntent } from '~/app/types';
import { mockModelTransferJob, mockModelTransferJobList } from '~/__mocks__/mockModelTransferJob';
import { mockModelRegistry } from '~/__mocks__/mockModelRegistry';
import { modelTransferJobsPage } from '~/__tests__/cypress/cypress/pages/modelRegistryView/modelTransferJobs';
import { MODEL_REGISTRY_API_VERSION } from '~/__tests__/cypress/cypress/support/commands/api';

const modelRegistryName = 'modelregistry-sample';

const jobList = mockModelTransferJobList({
  items: [
    mockModelTransferJob({
      id: 'job-to-delete',
      name: 'job-to-delete',
      jobDisplayName: 'job-to-delete',
    }),
    mockModelTransferJob({
      id: 'job-to-keep',
      name: 'job-to-keep',
      jobDisplayName: 'job-to-keep',
    }),
  ],
  size: 2,
  pageSize: 10,
  nextPageToken: '',
});

const interceptWithSingleFailedJob = (id: string, errorMessage: string) => {
  const failedJobList = mockModelTransferJobList({
    items: [
      mockModelTransferJob({
        id,
        name: id,
        jobDisplayName: id,
        status: ModelTransferJobStatus.FAILED,
        namespace: 'kubeflow',
        errorMessage,
      }),
    ],
    size: 1,
    pageSize: 10,
    nextPageToken: '',
  });

  cy.intercept(
    'GET',
    `**/api/${MODEL_REGISTRY_API_VERSION}/model_registry/${modelRegistryName}/model_transfer_jobs*`,
    mockModArchResponse(failedJobList),
  );
};

const buildFailedJobWithIntent = (
  id: string,
  intent: ModelTransferJobUploadIntent,
  errorMessage: string,
) =>
  mockModelTransferJob({
    id,
    name: id,
    jobDisplayName: id,
    status: ModelTransferJobStatus.FAILED,
    uploadIntent: intent,
    errorMessage,
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

  it('should render registered model and version as links only for completed jobs with IDs', () => {
    const jobsWithVariants = mockModelTransferJobList({
      items: [
        mockModelTransferJob({
          id: 'completed-with-ids',
          name: 'completed-with-ids',
          jobDisplayName: 'completed-with-ids',
          status: ModelTransferJobStatus.COMPLETED,
          registeredModelId: 'rm-1',
          registeredModelName: 'Completed model',
          modelVersionId: 'mv-1',
          modelVersionName: 'v1.0.0',
        }),
        mockModelTransferJob({
          id: 'completed-no-ids',
          name: 'completed-no-ids',
          jobDisplayName: 'completed-no-ids',
          status: ModelTransferJobStatus.COMPLETED,
          registeredModelId: undefined,
          modelVersionId: undefined,
          registeredModelName: 'Completed no IDs',
          modelVersionName: 'v-no-ids',
        }),
      ],
      size: 2,
      pageSize: 10,
      nextPageToken: '',
    });

    cy.intercept(
      'GET',
      `**/api/${MODEL_REGISTRY_API_VERSION}/model_registry/${modelRegistryName}/model_transfer_jobs*`,
      mockModArchResponse(jobsWithVariants),
    );

    modelTransferJobsPage.visit(modelRegistryName);
    modelTransferJobsPage.findTableRows().should('have.length', 2);

    // COMPLETED with IDs: model and version names are clickable and navigate to the correct routes
    const completedWithIdsRow = modelTransferJobsPage.getRow('completed-with-ids');

    completedWithIdsRow.find().findByRole('button', { name: 'Completed model' }).click();
    cy.location('pathname').should(
      'eq',
      `/model-registry/${modelRegistryName}/registered-models/rm-1/overview`,
    );

    cy.go('back');

    modelTransferJobsPage
      .getRow('completed-with-ids')
      .find()
      .findByRole('button', { name: 'v1.0.0' })
      .click();
    cy.location('pathname').should(
      'eq',
      `/model-registry/${modelRegistryName}/registered-models/rm-1/versions/mv-1/details`,
    );

    cy.go('back');

    // COMPLETED without IDs: names render as plain text (no buttons/links)
    const completedNoIdsRow = modelTransferJobsPage.getRow('completed-no-ids');
    completedNoIdsRow.find().contains('td', 'Completed no IDs').should('exist');
    completedNoIdsRow.find().findByRole('button', { name: 'Completed no IDs' }).should('not.exist');
    completedNoIdsRow.find().findByRole('button', { name: 'v-no-ids' }).should('not.exist');
  });

  it('should render status labels for all statuses and open/close the status modal', () => {
    const jobsWithStatuses = mockModelTransferJobList({
      items: [
        mockModelTransferJob({
          id: 'job-completed',
          name: 'job-completed',
          jobDisplayName: 'job-completed',
          status: ModelTransferJobStatus.COMPLETED,
        }),
        mockModelTransferJob({
          id: 'job-running',
          name: 'job-running',
          jobDisplayName: 'job-running',
          status: ModelTransferJobStatus.RUNNING,
        }),
        mockModelTransferJob({
          id: 'job-pending',
          name: 'job-pending',
          jobDisplayName: 'job-pending',
          status: ModelTransferJobStatus.PENDING,
        }),
        mockModelTransferJob({
          id: 'job-failed',
          name: 'job-failed',
          jobDisplayName: 'job-failed',
          status: ModelTransferJobStatus.FAILED,
        }),
        mockModelTransferJob({
          id: 'job-cancelled',
          name: 'job-cancelled',
          jobDisplayName: 'job-cancelled',
          status: ModelTransferJobStatus.CANCELLED,
        }),
      ],
      size: 5,
      pageSize: 10,
      nextPageToken: '',
    });

    cy.intercept(
      'GET',
      `**/api/${MODEL_REGISTRY_API_VERSION}/model_registry/${modelRegistryName}/model_transfer_jobs*`,
      mockModArchResponse(jobsWithStatuses),
    );

    modelTransferJobsPage.visit(modelRegistryName);
    modelTransferJobsPage.findTableRows().should('have.length', 5);

    const completedRow = modelTransferJobsPage.getRow('job-completed');
    completedRow.find().findByTestId('job-status').should('contain.text', 'Complete');

    const runningRow = modelTransferJobsPage.getRow('job-running');
    runningRow.find().findByTestId('job-status').should('contain.text', 'Running');

    const pendingRow = modelTransferJobsPage.getRow('job-pending');
    pendingRow.find().findByTestId('job-status').should('contain.text', 'Pending');

    const failedRow = modelTransferJobsPage.getRow('job-failed');
    failedRow.find().findByTestId('job-status').should('contain.text', 'Failed');

    const cancelledRow = modelTransferJobsPage.getRow('job-cancelled');
    cancelledRow.find().findByTestId('job-status').should('contain.text', 'Canceled');

    completedRow.find().findByTestId('job-status').click();
    cy.findByTestId('transfer-job-status-modal').should('be.visible');

    cy.findByLabelText('Close').click();
    cy.findByTestId('transfer-job-status-modal').should('not.exist');
  });

  it('should render the correct modal title for each upload intent', () => {
    const jobsWithIntents = mockModelTransferJobList({
      items: [
        buildFailedJobWithIntent(
          'job-create-model',
          ModelTransferJobUploadIntent.CREATE_MODEL,
          'Create model failed',
        ),
        buildFailedJobWithIntent(
          'job-create-version',
          ModelTransferJobUploadIntent.CREATE_VERSION,
          'Create version failed',
        ),
        buildFailedJobWithIntent(
          'job-update-artifact',
          ModelTransferJobUploadIntent.UPDATE_ARTIFACT,
          'Update artifact failed',
        ),
      ],
      size: 3,
      pageSize: 10,
      nextPageToken: '',
    });

    cy.intercept(
      'GET',
      `**/api/${MODEL_REGISTRY_API_VERSION}/model_registry/${modelRegistryName}/model_transfer_jobs*`,
      mockModArchResponse(jobsWithIntents),
    );

    modelTransferJobsPage.visit(modelRegistryName);
    modelTransferJobsPage.findTableRows().should('have.length', 3);

    const assertModalTitleForJob = (jobName: string, expectedTitle: string) => {
      const row = modelTransferJobsPage.getRow(jobName);
      row.find().findByTestId('job-status').click();
      cy.findByTestId('transfer-job-status-modal').contains(expectedTitle).should('be.visible');
      cy.findByLabelText('Close').click();
      cy.findByTestId('transfer-job-status-modal').should('not.exist');
    };

    assertModalTitleForJob('job-create-model', 'Model creation status');
    assertModalTitleForJob('job-create-version', 'Model version status');
    assertModalTitleForJob('job-update-artifact', 'Transfer job status');
  });

  it('should show failure alert and render event log entries for a failed job', () => {
    interceptWithSingleFailedJob('job-failed-events', 'Connection timeout while uploading');

    const events = [
      {
        timestamp: '2025-01-02T00:00:00Z',
        reason: 'BackOff',
        type: 'Warning',
        message: 'Back-off pulling image "quay.io/example/image:latest"',
      },
      {
        timestamp: '2025-01-01T00:00:00Z',
        reason: 'Pulling',
        type: 'Normal',
        message: 'Pulling image "quay.io/example/image:latest"',
      },
    ];

    cy.intercept(
      'GET',
      '**/model_transfer_jobs/job-failed-events/events*',
      mockModArchResponse({ events }),
    ).as('getFailedJobEvents');

    modelTransferJobsPage.visit(modelRegistryName);
    modelTransferJobsPage.findTableRows().should('have.length', 1);

    const failedRow = modelTransferJobsPage.getRow('job-failed-events');
    failedRow.find().findByTestId('job-status').should('contain.text', 'Failed').click();

    cy.findByTestId('transfer-job-status-modal').should('be.visible');

    // Failure alert shows the job error message
    cy.findByTestId('transfer-job-failure-alert')
      .should('be.visible')
      .and('contain.text', 'Connection timeout while uploading');

    cy.wait('@getFailedJobEvents');

    cy.findByTestId('transfer-job-event-log')
      .find('[data-testid="transfer-job-event-log-entry"]')
      .should('have.length', 2)
      .then(($items) => {
        const firstText = $items.eq(0).text();
        const secondText = $items.eq(1).text();
        expect(firstText).to.contain('2025-01-02T00:00:00Z [BackOff] [Warning]');
        expect(secondText).to.contain('2025-01-01T00:00:00Z [Pulling] [Normal]');
      });
  });

  it('should fall back to unknown failure reason when errorMessage is missing', () => {
    const failedJobList = mockModelTransferJobList({
      items: [
        mockModelTransferJob({
          id: 'job-failed-no-message',
          name: 'job-failed-no-message',
          jobDisplayName: 'job-failed-no-message',
          status: ModelTransferJobStatus.FAILED,
          namespace: 'kubeflow',
          errorMessage: undefined,
        }),
      ],
      size: 1,
      pageSize: 10,
      nextPageToken: '',
    });

    cy.intercept(
      'GET',
      `**/api/${MODEL_REGISTRY_API_VERSION}/model_registry/${modelRegistryName}/model_transfer_jobs*`,
      mockModArchResponse(failedJobList),
    );

    modelTransferJobsPage.visit(modelRegistryName);

    const failedRow = modelTransferJobsPage.getRow('job-failed-no-message');
    failedRow.find().findByTestId('job-status').click();

    cy.findByTestId('transfer-job-status-modal').should('be.visible');
    cy.findByTestId('transfer-job-failure-alert')
      .should('be.visible')
      .and('contain.text', 'Failure reason (unknown)');
  });

  it('should show empty message when there are no events', () => {
    interceptWithSingleFailedJob('job-failed-no-events', 'Some failure');

    cy.intercept(
      'GET',
      '**/model_transfer_jobs/job-failed-no-events/events*',
      mockModArchResponse({ events: [] }),
    ).as('getNoEvents');

    modelTransferJobsPage.visit(modelRegistryName);

    const failedRow = modelTransferJobsPage.getRow('job-failed-no-events');
    failedRow.find().findByTestId('job-status').click();

    cy.findByTestId('transfer-job-status-modal').should('be.visible');

    cy.wait('@getNoEvents');

    // When there are no events, EventLog shows its empty message
    cy.findByTestId('transfer-job-status-modal')
      .contains('There are no recent events.')
      .should('be.visible');
  });

  it('should keep showing spinner and not render event entries when the events API fails', () => {
    interceptWithSingleFailedJob('job-failed-error-events', 'Unexpected failure');

    cy.intercept('GET', '**/model_transfer_jobs/job-failed-error-events/events*', {
      forceNetworkError: true,
    }).as('getEventsError');

    modelTransferJobsPage.visit(modelRegistryName);

    const failedRow = modelTransferJobsPage.getRow('job-failed-error-events');
    failedRow.find().findByTestId('job-status').click();

    cy.findByTestId('transfer-job-status-modal').should('be.visible');

    cy.wait('@getEventsError');

    // When the events API fails with a network error, the spinner should remain
    // and no event log entries should render (the hook does not surface an error alert).
    cy.findByTestId('transfer-job-status-modal').findByLabelText('Contents').should('be.visible');
    cy.findByTestId('transfer-job-status-modal')
      .find('[data-testid="transfer-job-event-log-entry"]')
      .should('not.exist');
  });

  it('should render all key columns when the transfer jobs list is not empty', () => {
    const jobsWithDetails = mockModelTransferJobList({
      items: [
        mockModelTransferJob({
          id: 'detailed-job',
          name: 'detailed-job',
          jobDisplayName: 'Detailed transfer job',
          registeredModelName: 'My registered model',
          modelVersionName: 'v9.9.9',
          namespace: 'kubeflow',
          author: 'Sherlock Holmes',
          status: ModelTransferJobStatus.RUNNING,
        }),
      ],
      size: 1,
      pageSize: 10,
      nextPageToken: '',
    });

    cy.intercept(
      'GET',
      `**/api/${MODEL_REGISTRY_API_VERSION}/model_registry/${modelRegistryName}/model_transfer_jobs*`,
      mockModArchResponse(jobsWithDetails),
    );

    modelTransferJobsPage.visit(modelRegistryName);

    // Smoke check: table renders and contains the expected row and columns
    cy.findByTestId('model-transfer-jobs-table').within(() => {
      cy.findAllByRole('row').should('have.length.greaterThan', 1); // header + at least one data row

      cy.findByTestId('job-name').should('contain.text', 'Detailed transfer job');
      cy.contains('td', 'My registered model').should('exist');
      cy.contains('td', 'v9.9.9').should('exist');
      cy.findByTestId('job-namespace').should('contain.text', 'kubeflow');
      cy.findByTestId('job-author').should('contain.text', 'Sherlock Holmes');
      cy.findByTestId('job-status').should('contain.text', 'Running');
    });
  });

  it('should show the empty state when there are no transfer jobs', () => {
    const emptyList = mockModelTransferJobList({
      items: [],
      size: 0,
      pageSize: 10,
      nextPageToken: '',
    });

    cy.intercept(
      'GET',
      `**/api/${MODEL_REGISTRY_API_VERSION}/model_registry/${modelRegistryName}/model_transfer_jobs*`,
      mockModArchResponse(emptyList),
    );

    modelTransferJobsPage.visit(modelRegistryName);

    // When there are no jobs, the empty state should be visible and the table should not render
    modelTransferJobsPage.findEmptyState().should('be.visible');
    cy.findByTestId('model-transfer-jobs-table').should('not.exist');
  });
});
