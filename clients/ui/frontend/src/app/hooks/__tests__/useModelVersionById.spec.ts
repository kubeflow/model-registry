// eslint-disable-next-line @typescript-eslint/no-unused-vars
import * as React from 'react';
import { waitFor } from '@testing-library/react';
import { useFetchState } from 'mod-arch-core';
import useModelVersionById from '~/app/hooks/useModelVersionById';
import { useModelRegistryAPI } from '~/app/hooks/useModelRegistryAPI';
import { ModelRegistryAPIs } from '~/app/types';
import { mockModelVersion } from '~/__mocks__/mockModelVersion';
import { testHook } from '~/__tests__/unit/testUtils/hooks';

jest.mock('mod-arch-core', () => ({
  useFetchState: jest.fn(),
  NotReadyError: class NotReadyError extends Error {
    constructor(message: string) {
      super(message);
      this.name = 'NotReadyError';
    }
  },
}));

global.fetch = jest.fn();

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

describe('useModelVersionById', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('should return an error when the API is not available', async () => {
    mockUseModelRegistryAPI.mockReturnValue({
      api: mockModelRegistryAPIs,
      apiAvailable: false,
      refreshAllAPI: jest.fn(),
    });

    const mockError = new Error('API not yet available');
    mockUseFetchState.mockReturnValue([null, false, mockError, jest.fn()]);

    const { result } = testHook(useModelVersionById)('version-1');

    await waitFor(() => {
      const [, , error] = result.current;
      expect(error).toBeInstanceOf(Error);
      expect(error?.message).toBe('API not yet available');
    });
  });

  it('should silently fail when modelVersionId is not provided', async () => {
    mockUseModelRegistryAPI.mockReturnValue({
      api: mockModelRegistryAPIs,
      apiAvailable: true,
      refreshAllAPI: jest.fn(),
    });

    mockUseFetchState.mockReturnValue([null, false, undefined, jest.fn()]);

    const { result } = testHook(useModelVersionById)();

    await waitFor(() => {
      const [data, , error] = result.current;
      expect(data).toBeNull();
      expect(error).toBeUndefined();
    });
  });

  it('should fetch the model version when API is available and id is provided', async () => {
    const mockedVersion = mockModelVersion({ id: 'version-1', name: 'my-version' });

    mockUseModelRegistryAPI.mockReturnValue({
      api: {
        ...mockModelRegistryAPIs,
        getModelVersion: jest.fn().mockResolvedValue(mockedVersion),
      },
      apiAvailable: true,
      refreshAllAPI: jest.fn(),
    });

    mockUseFetchState.mockReturnValue([mockedVersion, true, undefined, jest.fn()]);

    const { result } = testHook(useModelVersionById)('version-1');

    await waitFor(() => {
      const [data, loaded] = result.current;
      expect(loaded).toBe(true);
      expect(data).toEqual(mockedVersion);
    });
  });

  it('should pass the correct modelVersionId to api.getModelVersion', async () => {
    const getModelVersionMock = jest.fn().mockResolvedValue(mockModelVersion({ id: 'v-42' }));

    mockUseModelRegistryAPI.mockReturnValue({
      api: { ...mockModelRegistryAPIs, getModelVersion: getModelVersionMock },
      apiAvailable: true,
      refreshAllAPI: jest.fn(),
    });

    // Capture the callback that the hook passes to useFetchState
    mockUseFetchState.mockReturnValue([null, false, undefined, jest.fn()]);

    testHook(useModelVersionById)('v-42');

    // useFetchState receives the callback as its first argument
    const capturedCallback = mockUseFetchState.mock.calls[0][0] as (
      opts: unknown,
    ) => Promise<unknown>;

    await capturedCallback({});

    expect(getModelVersionMock).toHaveBeenCalledWith({}, 'v-42');
  });
});
