import { waitFor } from '@testing-library/react';
import { useFetchState, POLL_INTERVAL } from 'mod-arch-core';
import useModelTransferJobs from '~/app/hooks/useModelTransferJobs';
import { useModelRegistryAPI } from '~/app/hooks/useModelRegistryAPI';
import { ModelRegistryAPIs, ModelTransferJobList, ModelTransferJobStatus } from '~/app/types';
import { mockModelTransferJob } from '~/__mocks__/mockModelTransferJob';
import { testHook } from '~/__tests__/unit/testUtils/hooks';

// Mock mod-arch-core to avoid React context issues and to inspect useFetchState calls
jest.mock('mod-arch-core', () => {
  const actual = jest.requireActual('mod-arch-core');
  return {
    ...actual,
    useFetchState: jest.fn(),
  };
});

// Mock the useModelRegistryAPI hook
jest.mock('~/app/hooks/useModelRegistryAPI', () => ({
  useModelRegistryAPI: jest.fn(),
}));

const mockUseModelRegistryAPI = jest.mocked(useModelRegistryAPI);
const mockUseFetchState = jest.mocked(useFetchState);

const mockModelRegistryAPIs: ModelRegistryAPIs = {
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
};

type FetchStateCapture = {
  getCapturedCallback: () => ((opts: unknown) => Promise<ModelTransferJobList>) | undefined;
  optionsCalls: Array<{ refreshRate?: number }>;
};

const setupFetchStateCapture = (): FetchStateCapture => {
  let capturedCallback: ((opts: unknown) => Promise<ModelTransferJobList>) | undefined;
  const optionsCalls: Array<{ refreshRate?: number }> = [];

  // The concrete types here are not important for the test; relax them to avoid clashing with useFetchState's generics.
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  mockUseFetchState.mockImplementation((cb: any, initialData: any, options: any) => {
    capturedCallback = cb as typeof capturedCallback;
    optionsCalls.push({ refreshRate: (options as { refreshRate?: number }).refreshRate });
    // Return a basic useFetchState tuple
    return [initialData, true, undefined, jest.fn()];
  });

  return {
    getCapturedCallback: () => capturedCallback,
    optionsCalls,
  };
};

describe('useModelTransferJobs', () => {
  beforeEach(() => {
    jest.clearAllMocks();

    mockUseModelRegistryAPI.mockReturnValue({
      api: mockModelRegistryAPIs,
      apiAvailable: true,
      refreshAllAPI: jest.fn(),
    });

    // Default mock for useFetchState: just return a basic empty state.
    mockUseFetchState.mockReturnValue([
      { items: [], size: 0, pageSize: 0, nextPageToken: '' },
      false,
      undefined,
      jest.fn(),
    ]);
  });

  it('sets refreshRate to POLL_INTERVAL when there are active (RUNNING/PENDING) jobs', async () => {
    // Arrange a response with one RUNNING job
    const activeJobsList: ModelTransferJobList = {
      items: [
        mockModelTransferJob({
          id: '1',
          name: 'job-running',
          jobDisplayName: 'job-running',
          description: 'Running job',
          namespace: 'kubeflow',
          status: ModelTransferJobStatus.RUNNING,
          createTimeSinceEpoch: '0',
          lastUpdateTimeSinceEpoch: '0',
        }),
      ],
      size: 1,
      pageSize: 10,
      nextPageToken: '',
    };

    const listModelTransferJobsMock =
      mockModelRegistryAPIs.listModelTransferJobs as jest.MockedFunction<
        ModelRegistryAPIs['listModelTransferJobs']
      >;

    // First call returns active jobs; subsequent calls can return the same
    listModelTransferJobsMock.mockResolvedValue(activeJobsList);

    // Capture the callback and options passed into useFetchState
    const { getCapturedCallback, optionsCalls } = setupFetchStateCapture();

    const renderResult = testHook(useModelTransferJobs)();

    // Wait for initial render
    await waitFor(() => {
      expect(renderResult.getUpdateCount()).toBeGreaterThan(0);
    });

    // Simulate the fetch callback resolving with active jobs, which should
    // toggle hasActiveJobs to true and cause a re-render with refreshRate = POLL_INTERVAL.
    const capturedCallback = getCapturedCallback();
    expect(capturedCallback).toBeDefined();
    await capturedCallback?.({});

    // Wait for the hook to re-render after state update
    await renderResult.waitForNextUpdate();

    // The latest call to useFetchState should have refreshRate set to POLL_INTERVAL
    const lastCallOptions = optionsCalls[optionsCalls.length - 1];
    expect(lastCallOptions.refreshRate).toBe(POLL_INTERVAL);
  });

  it('does not set refreshRate when all jobs are in terminal states', async () => {
    // Arrange a response with only COMPLETED and FAILED jobs
    const terminalJobsList: ModelTransferJobList = {
      items: [
        mockModelTransferJob({
          id: '1',
          name: 'job-completed',
          jobDisplayName: 'job-completed',
          description: 'Completed job',
          namespace: 'kubeflow',
          status: ModelTransferJobStatus.COMPLETED,
          createTimeSinceEpoch: '0',
          lastUpdateTimeSinceEpoch: '0',
        }),
        mockModelTransferJob({
          id: '2',
          name: 'job-failed',
          jobDisplayName: 'job-failed',
          description: 'Failed job',
          namespace: 'kubeflow',
          status: ModelTransferJobStatus.FAILED,
          createTimeSinceEpoch: '0',
          lastUpdateTimeSinceEpoch: '0',
        }),
      ],
      size: 2,
      pageSize: 10,
      nextPageToken: '',
    };

    const listModelTransferJobsMock =
      mockModelRegistryAPIs.listModelTransferJobs as jest.MockedFunction<
        ModelRegistryAPIs['listModelTransferJobs']
      >;

    listModelTransferJobsMock.mockResolvedValue(terminalJobsList);

    const { getCapturedCallback, optionsCalls } = setupFetchStateCapture();

    testHook(useModelTransferJobs)();

    // Simulate fetch resolving with only terminal jobs; hasActiveJobs should stay false,
    // so refreshRate should remain undefined (no need to wait for another update).
    const capturedCallback = getCapturedCallback();
    expect(capturedCallback).toBeDefined();
    await capturedCallback?.({});

    const lastCallOptions = optionsCalls[optionsCalls.length - 1];
    expect(lastCallOptions.refreshRate).toBeUndefined();
  });
});
