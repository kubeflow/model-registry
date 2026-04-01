import React from 'react';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import { useFetchState } from 'mod-arch-core';
import { mockModelTransferJob } from '~/__mocks__/mockModelTransferJob';
import { ModelTransferJobStatus } from '~/app/types';
import ModelTransferJobStatusModal from '~/app/pages/modelRegistry/screens/ModelTransferJobs/ModelTransferJobStatusModal';
import { useModelRegistryAPI } from '~/app/hooks/useModelRegistryAPI';

// Mock mod-arch-core's useFetchState so we can drive the modal's loading/error states directly.
jest.mock('mod-arch-core', () => {
  const actual = jest.requireActual('mod-arch-core');
  return {
    ...actual,
    useFetchState: jest.fn(),
  };
});

// Mock the useModelRegistryAPI hook to avoid needing real context/API wiring.
jest.mock('~/app/hooks/useModelRegistryAPI', () => ({
  useModelRegistryAPI: jest.fn(),
}));

const mockUseFetchState = jest.mocked(useFetchState);
const mockUseModelRegistryAPI = jest.mocked(useModelRegistryAPI);

describe('ModelTransferJobStatusModal', () => {
  beforeEach(() => {
    jest.clearAllMocks();

    mockUseModelRegistryAPI.mockReturnValue({
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
    });
  });

  it('renders the status label in the modal title with outline variant', () => {
    const job = mockModelTransferJob({
      status: ModelTransferJobStatus.COMPLETED,
      namespace: 'kubeflow',
    });

    mockUseFetchState.mockReturnValue([[], true, undefined, jest.fn()]);

    render(<ModelTransferJobStatusModal job={job} isOpen onClose={jest.fn()} />);

    const label = screen.getByText('Complete');
    expect(label).toBeVisible();

    const labelWrapper = label.closest('span.pf-v6-c-label');
    expect(labelWrapper).not.toBeNull();
    expect(labelWrapper!.className).toMatch(/outline/);
    expect(labelWrapper!.className).not.toMatch(/filled/);
  });

  it('shows unknown failure reason and danger alert when events fail to load', () => {
    // Arrange job without an explicit errorMessage so the fallback text is used.
    const job = mockModelTransferJob({
      status: ModelTransferJobStatus.FAILED,
      errorMessage: undefined,
      namespace: 'kubeflow',
    });

    // useFetchState returns: no events, loaded=true, and a load error.
    mockUseFetchState.mockReturnValue([[], true, new Error('Events API failed'), jest.fn()]);

    render(<ModelTransferJobStatusModal job={job} isOpen onClose={jest.fn()} />);

    // Failure alert uses the fallback message when errorMessage is missing.
    expect(screen.getByTestId('transfer-job-failure-alert')).toBeInTheDocument();
    expect(screen.getByText('Failure reason (unknown)')).toBeInTheDocument();

    // Events tab shows the "Failed to load events" danger alert with the error message.
    expect(screen.getByText('Failed to load events')).toBeInTheDocument();
    expect(screen.getByText('Events API failed')).toBeInTheDocument();
  });
});
