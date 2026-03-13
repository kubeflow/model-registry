import '@testing-library/jest-dom';
import React from 'react';
import { render, act } from '@testing-library/react';
import {
  TransferJobNotificationsProvider,
  TransferJobNotificationsContext,
} from '~/app/context/TransferJobNotificationsContext';
import { ModelTransferJobStatus } from '~/app/types';

const mockNotification = {
  info: jest.fn(),
  success: jest.fn(),
  error: jest.fn(),
  warning: jest.fn(),
  remove: jest.fn(),
};

jest.mock('~/app/hooks/useNotification', () => ({
  useNotification: () => mockNotification,
}));

jest.mock('mod-arch-core', () => ({
  useQueryParamNamespaces: () => ({}),
}));

const mockGetModelTransferJobByName = jest.fn();

jest.mock('~/app/api/service', () => ({
  getModelTransferJobByName: () => mockGetModelTransferJobByName,
}));

jest.mock('~/app/utilities/const', () => ({
  URL_PREFIX: '/model-registry',
  BFF_API_VERSION: 'v1',
  POLL_INTERVAL: 1000,
  REGISTRATION_TOAST_TITLES: {
    REGISTER_AND_STORE_STARTED: 'Model transfer job started',
    REGISTER_AND_STORE_SUCCEEDED: 'Model transfer job succeeded',
    REGISTER_AND_STORE_ERROR: 'Model transfer job failed',
  },
}));

jest.mock('~/app/pages/modelRegistry/screens/routeUtils', () => ({
  modelTransferJobsUrl: jest.fn((mrName: string) => `/model-registry/${mrName}/jobs`),
}));

const mockJobResponse = (name: string, status: ModelTransferJobStatus) => ({
  id: `id-${name}`,
  name,
  status,
});

const renderWithWatcher = async (jobName: string) => {
  function TestConsumer() {
    const { watchJob } = React.useContext(TransferJobNotificationsContext);
    React.useEffect(() => {
      watchJob({
        jobName,
        jobNamespace: 'test-namespace',
        registryName: 'mr-sample',
        displayParams: { versionModelName: 'My Model / v1', mrName: 'mr-sample' },
      });
    }, [watchJob]);
    return null;
  }

  await act(async () => {
    render(
      <TransferJobNotificationsProvider>
        <TestConsumer />
      </TransferJobNotificationsProvider>,
    );
  });

  await act(async () => {
    await Promise.resolve();
  });
};

describe('TransferJobNotificationsContext', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    jest.useFakeTimers();
  });

  afterEach(() => {
    jest.useRealTimers();
  });

  it('shows success toast when a watched job completes', async () => {
    mockGetModelTransferJobByName.mockResolvedValue(
      mockJobResponse('job-1', ModelTransferJobStatus.COMPLETED),
    );
    await renderWithWatcher('job-1');

    expect(mockNotification.success).toHaveBeenCalledWith(
      'Model transfer job succeeded',
      expect.anything(),
    );
  });

  it('shows error toast when a watched job fails', async () => {
    mockGetModelTransferJobByName.mockResolvedValue(
      mockJobResponse('job-2', ModelTransferJobStatus.FAILED),
    );
    await renderWithWatcher('job-2');

    expect(mockNotification.error).toHaveBeenCalledWith(
      'Model transfer job failed',
      expect.anything(),
    );
  });

  it('does not show toast for running jobs and keeps polling until completion', async () => {
    mockGetModelTransferJobByName.mockResolvedValue(
      mockJobResponse('job-3', ModelTransferJobStatus.RUNNING),
    );
    await renderWithWatcher('job-3');

    expect(mockNotification.success).not.toHaveBeenCalled();
    expect(mockNotification.error).not.toHaveBeenCalled();

    mockGetModelTransferJobByName.mockResolvedValue(
      mockJobResponse('job-3', ModelTransferJobStatus.COMPLETED),
    );

    await act(async () => {
      jest.advanceTimersByTime(1000);
      await Promise.resolve();
    });

    expect(mockNotification.success).toHaveBeenCalledWith(
      'Model transfer job succeeded',
      expect.anything(),
    );
  });

  it('silently removes cancelled jobs without showing toast', async () => {
    mockGetModelTransferJobByName.mockResolvedValue(
      mockJobResponse('job-4', ModelTransferJobStatus.CANCELLED),
    );
    await renderWithWatcher('job-4');

    expect(mockNotification.success).not.toHaveBeenCalled();
    expect(mockNotification.error).not.toHaveBeenCalled();
  });

  it('handles API errors gracefully and continues polling', async () => {
    mockGetModelTransferJobByName.mockRejectedValueOnce(new Error('Network error'));
    await renderWithWatcher('job-5');

    expect(mockNotification.success).not.toHaveBeenCalled();
    expect(mockNotification.error).not.toHaveBeenCalled();

    mockGetModelTransferJobByName.mockResolvedValue(
      mockJobResponse('job-5', ModelTransferJobStatus.COMPLETED),
    );

    await act(async () => {
      jest.advanceTimersByTime(1000);
      await Promise.resolve();
    });

    expect(mockNotification.success).toHaveBeenCalledWith(
      'Model transfer job succeeded',
      expect.anything(),
    );
  });
});
