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

const mockListModelTransferJobs = jest.fn();

jest.mock('~/app/api/service', () => ({
  getListModelTransferJobs: () => mockListModelTransferJobs,
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

describe('TransferJobNotificationsContext', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    jest.useFakeTimers();
  });

  afterEach(() => {
    jest.useRealTimers();
  });

  it('provides watchJob function', () => {
    let contextValue: { watchJob: (job: never) => void } | undefined;
    function Consumer() {
      contextValue = React.useContext(TransferJobNotificationsContext);
      return null;
    }
    render(
      <TransferJobNotificationsProvider>
        <Consumer />
      </TransferJobNotificationsProvider>,
    );
    expect(contextValue).toBeDefined();
    expect(typeof contextValue!.watchJob).toBe('function');
  });

  it('shows success toast when a watched job completes', async () => {
    mockListModelTransferJobs.mockResolvedValue({
      items: [
        {
          id: 'job-1',
          name: 'transfer-job-1',
          status: ModelTransferJobStatus.COMPLETED,
        },
      ],
      size: 1,
      pageSize: 10,
      nextPageToken: '',
    });

    function TestConsumer() {
      const { watchJob } = React.useContext(TransferJobNotificationsContext);
      React.useEffect(() => {
        watchJob({
          jobId: 'job-1',
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

    // The immediate poll should have fired
    await act(async () => {
      await Promise.resolve();
    });

    expect(mockNotification.success).toHaveBeenCalledWith(
      'Model transfer job succeeded',
      expect.anything(),
    );
  });

  it('shows error toast when a watched job fails', async () => {
    mockListModelTransferJobs.mockResolvedValue({
      items: [
        {
          id: 'job-2',
          name: 'transfer-job-2',
          status: ModelTransferJobStatus.FAILED,
        },
      ],
      size: 1,
      pageSize: 10,
      nextPageToken: '',
    });

    function TestConsumer() {
      const { watchJob } = React.useContext(TransferJobNotificationsContext);
      React.useEffect(() => {
        watchJob({
          jobId: 'job-2',
          registryName: 'mr-sample',
          displayParams: { versionModelName: 'My Model / v2', mrName: 'mr-sample' },
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

    expect(mockNotification.error).toHaveBeenCalledWith(
      'Model transfer job failed',
      expect.anything(),
    );
  });

  it('does not show toast for running jobs and keeps polling', async () => {
    mockListModelTransferJobs.mockResolvedValue({
      items: [
        {
          id: 'job-3',
          name: 'transfer-job-3',
          status: ModelTransferJobStatus.RUNNING,
        },
      ],
      size: 1,
      pageSize: 10,
      nextPageToken: '',
    });

    function TestConsumer() {
      const { watchJob } = React.useContext(TransferJobNotificationsContext);
      React.useEffect(() => {
        watchJob({
          jobId: 'job-3',
          registryName: 'mr-sample',
          displayParams: { versionModelName: 'My Model / v3', mrName: 'mr-sample' },
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

    expect(mockNotification.success).not.toHaveBeenCalled();
    expect(mockNotification.error).not.toHaveBeenCalled();

    // Simulate the job completing on next poll
    mockListModelTransferJobs.mockResolvedValue({
      items: [
        {
          id: 'job-3',
          name: 'transfer-job-3',
          status: ModelTransferJobStatus.COMPLETED,
        },
      ],
      size: 1,
      pageSize: 10,
      nextPageToken: '',
    });

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
    mockListModelTransferJobs.mockResolvedValue({
      items: [
        {
          id: 'job-4',
          name: 'transfer-job-4',
          status: ModelTransferJobStatus.CANCELLED,
        },
      ],
      size: 1,
      pageSize: 10,
      nextPageToken: '',
    });

    function TestConsumer() {
      const { watchJob } = React.useContext(TransferJobNotificationsContext);
      React.useEffect(() => {
        watchJob({
          jobId: 'job-4',
          registryName: 'mr-sample',
          displayParams: { versionModelName: 'My Model / v4', mrName: 'mr-sample' },
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

    expect(mockNotification.success).not.toHaveBeenCalled();
    expect(mockNotification.error).not.toHaveBeenCalled();
  });

  it('handles API errors gracefully and continues polling', async () => {
    mockListModelTransferJobs.mockRejectedValueOnce(new Error('Network error'));

    function TestConsumer() {
      const { watchJob } = React.useContext(TransferJobNotificationsContext);
      React.useEffect(() => {
        watchJob({
          jobId: 'job-5',
          registryName: 'mr-sample',
          displayParams: { versionModelName: 'My Model / v5', mrName: 'mr-sample' },
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

    // No toast on error
    expect(mockNotification.success).not.toHaveBeenCalled();
    expect(mockNotification.error).not.toHaveBeenCalled();

    // Should still poll on next interval
    mockListModelTransferJobs.mockResolvedValue({
      items: [
        {
          id: 'job-5',
          name: 'transfer-job-5',
          status: ModelTransferJobStatus.COMPLETED,
        },
      ],
      size: 1,
      pageSize: 10,
      nextPageToken: '',
    });

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
