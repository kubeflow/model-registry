import * as React from 'react';
import { render, screen, fireEvent, within } from '@testing-library/react';
import '@testing-library/jest-dom';
import { MemoryRouter } from 'react-router-dom';
import { Table, Tbody } from '@patternfly/react-table';
import { mockModelTransferJob } from '~/__mocks__/mockModelTransferJob';
import { ModelTransferJobStatus } from '~/app/types';
import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';
import ModelTransferJobTableRow from '~/app/pages/modelRegistry/screens/ModelTransferJobs/ModelTransferJobTableRow';

jest.mock('~/app/hooks/useModelRegistryAPI', () => ({
  useModelRegistryAPI: jest.fn().mockReturnValue({
    api: {
      createRegisteredModel: jest.fn(),
      createModelVersionForRegisteredModel: jest.fn(),
      createModelArtifactForModelVersion: jest.fn(),
      getRegisteredModel: jest.fn(),
      getModelVersion: jest.fn(),
      listModelVersions: jest.fn(),
      listRegisteredModels: jest.fn(),
      getModelVersionsByRegisteredModel: jest.fn(),
      getModelArtifactsByModelVersion: jest.fn(),
      patchRegisteredModel: jest.fn(),
      patchModelVersion: jest.fn(),
      patchModelArtifact: jest.fn(),
      listModelTransferJobs: jest.fn(),
      getModelTransferJobByName: jest.fn(),
      createModelTransferJob: jest.fn(),
      updateModelTransferJob: jest.fn(),
      deleteModelTransferJob: jest.fn(),
      getModelTransferJobEvents: jest.fn(),
    },
    apiAvailable: true,
    refreshAllAPI: jest.fn(),
  }),
}));

jest.mock('mod-arch-core', () => {
  const actual = jest.requireActual('mod-arch-core');
  return {
    ...actual,
    useFetchState: jest.fn().mockReturnValue([[], true, undefined, jest.fn()]),
  };
});

const mockContextValue = {
  modelRegistriesLoaded: true,
  modelRegistriesLoadError: undefined,
  modelRegistries: [],
  preferredModelRegistry: undefined,
  updatePreferredModelRegistry: jest.fn(),
};

const renderRow = (props: React.ComponentProps<typeof ModelTransferJobTableRow>) =>
  render(
    <MemoryRouter>
      <ModelRegistrySelectorContext.Provider value={mockContextValue}>
        <Table>
          <Tbody>
            <ModelTransferJobTableRow {...props} />
          </Tbody>
        </Table>
      </ModelRegistrySelectorContext.Provider>
    </MemoryRouter>,
  );

describe('ModelTransferJobTableRow', () => {
  it('shows truncated error message below Failed label when errorMessage exists', () => {
    const job = mockModelTransferJob({
      status: ModelTransferJobStatus.FAILED,
      errorMessage: 'Connection timeout while uploading to destination bucket',
    });

    renderRow({ job });

    expect(screen.getByTestId('job-status')).toHaveTextContent('Failed');
    expect(screen.getByTestId('job-error-message')).toBeInTheDocument();
    expect(screen.getByTestId('job-error-message')).toHaveTextContent(
      'Connection timeout while uploading to destination bucket',
    );
  });

  it('does not show error message when job is FAILED but errorMessage is undefined', () => {
    const job = mockModelTransferJob({
      status: ModelTransferJobStatus.FAILED,
      errorMessage: undefined,
    });

    renderRow({ job });

    expect(screen.getByTestId('job-status')).toHaveTextContent('Failed');
    expect(screen.queryByTestId('job-error-message')).not.toBeInTheDocument();
  });

  it('does not show error message for non-FAILED statuses', () => {
    const job = mockModelTransferJob({
      status: ModelTransferJobStatus.COMPLETED,
    });

    renderRow({ job });

    expect(screen.getByTestId('job-status')).toHaveTextContent('Complete');
    expect(screen.queryByTestId('job-error-message')).not.toBeInTheDocument();
  });

  it('opens status modal when clicking the truncated error message', () => {
    const job = mockModelTransferJob({
      status: ModelTransferJobStatus.FAILED,
      errorMessage: 'Upload failed: insufficient storage',
    });

    renderRow({ job });

    expect(screen.queryByTestId('transfer-job-status-modal')).not.toBeInTheDocument();

    fireEvent.click(screen.getByTestId('job-error-message'));

    expect(screen.getByTestId('transfer-job-status-modal')).toBeInTheDocument();
  });

  it('opens status modal when clicking the Failed label', () => {
    const job = mockModelTransferJob({
      status: ModelTransferJobStatus.FAILED,
      errorMessage: 'Some error',
    });

    renderRow({ job });

    expect(screen.queryByTestId('transfer-job-status-modal')).not.toBeInTheDocument();

    const labelWrapper = screen.getByTestId('job-status');
    const labelButton = within(labelWrapper).getByRole('button');
    fireEvent.click(labelButton);

    expect(screen.getByTestId('transfer-job-status-modal')).toBeInTheDocument();
  });

  it('shows retry button for failed jobs when onRequestRetry is provided', () => {
    const mockRetry = jest.fn();
    const job = mockModelTransferJob({
      status: ModelTransferJobStatus.FAILED,
      errorMessage: 'Connection error',
    });

    renderRow({ job, onRequestRetry: mockRetry });

    expect(screen.getByTestId('job-retry-button')).toBeInTheDocument();
    fireEvent.click(screen.getByTestId('job-retry-button'));
    expect(mockRetry).toHaveBeenCalledWith(job);
  });

  it('applies dotted underline style to the error message', () => {
    const job = mockModelTransferJob({
      status: ModelTransferJobStatus.FAILED,
      errorMessage: 'Some error message',
    });

    renderRow({ job });

    const errorButton = screen.getByTestId('job-error-message');
    expect(errorButton).toHaveStyle({ textDecoration: 'underline dotted' });
  });

  it('does not show error message for empty string errorMessage', () => {
    const job = mockModelTransferJob({
      status: ModelTransferJobStatus.FAILED,
      errorMessage: '',
    });

    renderRow({ job });

    expect(screen.getByTestId('job-status')).toHaveTextContent('Failed');
    expect(screen.queryByTestId('job-error-message')).not.toBeInTheDocument();
  });
});
