import {
  ModelRegistryAPIs,
  ModelState,
  ModelRegistryMetadataType,
  ModelVersion,
  RegisteredModel,
} from '~/app/types';

type MinimalModelRegistryAPI = Pick<ModelRegistryAPIs, 'patchRegisteredModel'>;

export const bumpModelVersionTimestamp = async (
  api: ModelRegistryAPIs,
  modelVersion: ModelVersion,
): Promise<void> => {
  if (!modelVersion.id) {
    throw new Error('Model version ID is required');
  }

  try {
    const currentTime = new Date().toISOString();
    await api.patchModelVersion(
      {},
      {
        // This is a workaround to update the timestamp on the backend. There is a bug opened for model registry team
        // to fix this issue.
        state: ModelState.LIVE,
        customProperties: {
          ...modelVersion.customProperties,
          _lastModified: {
            metadataType: ModelRegistryMetadataType.STRING,
            // eslint-disable-next-line camelcase
            string_value: currentTime,
          },
        },
      },
      modelVersion.id,
    );
  } catch (error) {
    throw new Error(
      `Failed to update model version timestamp: ${
        error instanceof Error ? error.message : String(error)
      }`,
    );
  }
};

export const bumpRegisteredModelTimestamp = async (
  api: MinimalModelRegistryAPI,
  registeredModel: RegisteredModel,
): Promise<void> => {
  if (!registeredModel.id) {
    throw new Error('Registered model ID is required');
  }

  try {
    const currentTime = new Date().toISOString();
    await api.patchRegisteredModel(
      {},
      {
        state: ModelState.LIVE,
        customProperties: {
          ...registeredModel.customProperties,
          // This is a workaround to update the timestamp on the backend. There is a bug opened for model registry team
          // to fix this issue.
          _lastModified: {
            metadataType: ModelRegistryMetadataType.STRING,
            // eslint-disable-next-line camelcase
            string_value: currentTime,
          },
        },
      },
      registeredModel.id,
    );
  } catch (error) {
    throw new Error(
      `Failed to update registered model timestamp: ${
        error instanceof Error ? error.message : String(error)
      }`,
    );
  }
};

export const bumpBothTimestamps = async (
  api: ModelRegistryAPIs,
  registeredModel: RegisteredModel,
  modelVersion: ModelVersion,
): Promise<void> => {
  await Promise.all([
    bumpModelVersionTimestamp(api, modelVersion),
    bumpRegisteredModelTimestamp(api, registeredModel),
  ]);
};
