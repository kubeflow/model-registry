// eslint-disable-next-line @typescript-eslint/no-unused-vars
import * as React from 'react';
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

const captureCallback = (): ((opts: unknown) => Promise<unknown>) => {
  mockUseFetchState.mockReturnValue([null, false, undefined, jest.fn()]);
  return mockUseFetchState.mock.calls[0][0] as (opts: unknown) => Promise<unknown>;
};

describe('useModelVersionById', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('should reject with an error when the API is not available', async () => {
    mockUseModelRegistryAPI.mockReturnValue({
      api: mockModelRegistryAPIs,
      apiAvailable: false,
      refreshAllAPI: jest.fn(),
    });

    testHook(useModelVersionById)('version-1');
    const callback = captureCallback();

    await expect(callback({})).rejects.toThrow('API not yet available');
  });

  it('should reject with NotReadyError when modelVersionId is not provided', async () => {
    mockUseModelRegistryAPI.mockReturnValue({
      api: mockModelRegistryAPIs,
      apiAvailable: true,
      refreshAllAPI: jest.fn(),
    });

    testHook(useModelVersionById)();
    const callback = captureCallback();

    await expect(callback({})).rejects.toThrow('No model version id');
    await expect(callback({})).rejects.toMatchObject({ name: 'NotReadyError' });
  });

  it('should call api.getModelVersion with the correct modelVersionId', async () => {
    const getModelVersionMock = jest.fn().mockResolvedValue(mockModelVersion({ id: 'v-42' }));

    mockUseModelRegistryAPI.mockReturnValue({
      api: { ...mockModelRegistryAPIs, getModelVersion: getModelVersionMock },
      apiAvailable: true,
      refreshAllAPI: jest.fn(),
    });

    testHook(useModelVersionById)('v-42');
    const callback = captureCallback();

    const result = await callback({});

    expect(getModelVersionMock).toHaveBeenCalledWith({}, 'v-42');
    expect(result).toEqual(mockModelVersion({ id: 'v-42' }));
  });
});
