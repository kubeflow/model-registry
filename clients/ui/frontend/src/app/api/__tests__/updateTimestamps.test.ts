import { mockModelVersion, mockRegisteredModel } from '~/__mocks__';
import {
  ModelRegistryAPIs,
  ModelState,
  ModelRegistryMetadataType,
  ModelVersion,
  RegisteredModel,
} from '~/app/types';
import {
  bumpModelVersionTimestamp,
  bumpRegisteredModelTimestamp,
  bumpBothTimestamps,
} from '~/app/api/updateTimestamps';

describe('updateTimestamps', () => {
  const mockApi = jest.mocked<ModelRegistryAPIs>({
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
  });
  const fakeModelVersionId = 'test-model-version-id';
  const fakeRegisteredModelId = 'test-registered-model-id';

  beforeEach(() => {
    jest.spyOn(Date.prototype, 'toISOString').mockReturnValue('2024-01-01T00:00:00.000Z');
  });

  describe('bumpModelVersionTimestamp', () => {
    it('should successfully update model version timestamp', async () => {
      await bumpModelVersionTimestamp(mockApi, mockModelVersion({ id: fakeModelVersionId }));

      expect(mockApi.patchModelVersion).toHaveBeenCalledWith(
        {},
        {
          state: ModelState.LIVE,
          customProperties: {
            _lastModified: {
              metadataType: ModelRegistryMetadataType.STRING,
              // eslint-disable-next-line camelcase
              string_value: '2024-01-01T00:00:00.000Z',
            },
          },
        },
        fakeModelVersionId,
      );
    });

    it('should successfully update model version timestamp with customProperties of version', async () => {
      await bumpModelVersionTimestamp(
        mockApi,
        mockModelVersion({
          id: fakeModelVersionId,
          customProperties: {
            'Registered from': {
              // eslint-disable-next-line camelcase
              string_value: 'Model catalog',
              metadataType: ModelRegistryMetadataType.STRING,
            },
            'Source model': {
              // eslint-disable-next-line camelcase
              string_value: 'test',
              metadataType: ModelRegistryMetadataType.STRING,
            },
          },
        }),
      );

      expect(mockApi.patchModelVersion).toHaveBeenCalledWith(
        {},
        {
          state: ModelState.LIVE,
          customProperties: {
            'Registered from': {
              // eslint-disable-next-line camelcase
              string_value: 'Model catalog',
              metadataType: ModelRegistryMetadataType.STRING,
            },
            'Source model': {
              // eslint-disable-next-line camelcase
              string_value: 'test',
              metadataType: ModelRegistryMetadataType.STRING,
            },
            _lastModified: {
              metadataType: ModelRegistryMetadataType.STRING,
              // eslint-disable-next-line camelcase
              string_value: '2024-01-01T00:00:00.000Z',
            },
          },
        },
        fakeModelVersionId,
      );
    });

    it('should throw error if modelVersionId is empty', async () => {
      await expect(
        bumpModelVersionTimestamp(mockApi, mockModelVersion({ id: '' })),
      ).rejects.toThrow('Model version ID is required');
    });

    it('should handle API errors appropriately', async () => {
      const errorMessage = 'API Error';
      // Use proper type for mock function
      const mockFn = mockApi.patchModelVersion;
      mockFn.mockRejectedValue(new Error(errorMessage));

      await expect(
        bumpModelVersionTimestamp(mockApi, mockModelVersion({ id: fakeModelVersionId })),
      ).rejects.toThrow(`Failed to update model version timestamp: ${errorMessage}`);
    });
  });

  describe('bumpRegisteredModelTimestamp', () => {
    it('should successfully update registered model timestamp with model customProperties', async () => {
      await bumpRegisteredModelTimestamp(
        mockApi,
        mockRegisteredModel({
          id: fakeRegisteredModelId,
          customProperties: {
            'Registered from': {
              // eslint-disable-next-line camelcase
              string_value: 'Model catalog',
              metadataType: ModelRegistryMetadataType.STRING,
            },
            'Source model': {
              // eslint-disable-next-line camelcase
              string_value: 'test',
              metadataType: ModelRegistryMetadataType.STRING,
            },
          },
        }),
      );

      expect(mockApi.patchRegisteredModel).toHaveBeenCalledWith(
        {},
        {
          state: ModelState.LIVE,
          customProperties: {
            'Registered from': {
              // eslint-disable-next-line camelcase
              string_value: 'Model catalog',
              metadataType: ModelRegistryMetadataType.STRING,
            },
            'Source model': {
              // eslint-disable-next-line camelcase
              string_value: 'test',
              metadataType: ModelRegistryMetadataType.STRING,
            },
            _lastModified: {
              metadataType: ModelRegistryMetadataType.STRING,
              // eslint-disable-next-line camelcase
              string_value: '2024-01-01T00:00:00.000Z',
            },
          },
        },
        fakeRegisteredModelId,
      );
    });

    it('should successfully update registered model timestamp', async () => {
      await bumpRegisteredModelTimestamp(
        mockApi,
        mockRegisteredModel({
          id: fakeRegisteredModelId,
        }),
      );

      expect(mockApi.patchRegisteredModel).toHaveBeenCalledWith(
        {},
        {
          state: ModelState.LIVE,
          customProperties: {
            _lastModified: {
              metadataType: ModelRegistryMetadataType.STRING,
              // eslint-disable-next-line camelcase
              string_value: '2024-01-01T00:00:00.000Z',
            },
          },
        },
        fakeRegisteredModelId,
      );
    });

    it('should throw error if registeredModelId is empty', async () => {
      await expect(
        bumpRegisteredModelTimestamp(
          mockApi,
          mockRegisteredModel({
            id: '',
          }),
        ),
      ).rejects.toThrow('Registered model ID is required');
    });

    it('should handle API errors appropriately', async () => {
      const errorMessage = 'API Error';
      // Use proper type for mock function
      const mockFn = mockApi.patchRegisteredModel;
      mockFn.mockRejectedValue(new Error(errorMessage));

      await expect(
        bumpRegisteredModelTimestamp(
          mockApi,
          mockRegisteredModel({
            id: fakeRegisteredModelId,
          }),
        ),
      ).rejects.toThrow(`Failed to update registered model timestamp: ${errorMessage}`);
    });
  });

  describe('bumpBothTimestamps', () => {
    it('should update both timestamps successfully', async () => {
      mockApi.patchModelVersion.mockResolvedValue({} as ModelVersion);
      mockApi.patchRegisteredModel.mockResolvedValue({} as RegisteredModel);

      await bumpBothTimestamps(
        mockApi,
        mockRegisteredModel({
          id: fakeRegisteredModelId,
        }),
        mockModelVersion({ id: fakeModelVersionId }),
      );

      expect(mockApi.patchModelVersion).toHaveBeenCalled();
      expect(mockApi.patchRegisteredModel).toHaveBeenCalled();
    });

    it('should update both timestamps successfully with both model and version customProperties', async () => {
      mockApi.patchModelVersion.mockResolvedValue({} as ModelVersion);
      mockApi.patchRegisteredModel.mockResolvedValue({} as RegisteredModel);

      await bumpBothTimestamps(
        mockApi,

        mockRegisteredModel({
          id: fakeRegisteredModelId,
          customProperties: {
            'Registered from': {
              // eslint-disable-next-line camelcase
              string_value: 'Model catalog',
              metadataType: ModelRegistryMetadataType.STRING,
            },
            'Source model': {
              // eslint-disable-next-line camelcase
              string_value: 'test',
              metadataType: ModelRegistryMetadataType.STRING,
            },
          },
        }),
        mockModelVersion({
          id: fakeModelVersionId,
          customProperties: {
            'Registered from': {
              // eslint-disable-next-line camelcase
              string_value: 'Model catalog',
              metadataType: ModelRegistryMetadataType.STRING,
            },
            'Source model': {
              // eslint-disable-next-line camelcase
              string_value: 'test-version',
              metadataType: ModelRegistryMetadataType.STRING,
            },
          },
        }),
      );

      expect(mockApi.patchModelVersion).toHaveBeenCalled();
      expect(mockApi.patchRegisteredModel).toHaveBeenCalled();
    });

    it('should handle errors from either update', async () => {
      const errorMessage = 'API Error';
      mockApi.patchModelVersion.mockRejectedValue(new Error(errorMessage));

      await expect(
        bumpBothTimestamps(
          mockApi,
          mockRegisteredModel({
            id: fakeRegisteredModelId,
          }),
          mockModelVersion({ id: fakeModelVersionId }),
        ),
      ).rejects.toThrow();
    });
  });
});
