import { ModelRegistryAPIs, ModelState, ModelRegistryMetadataType } from '~/app/types';

type MinimalModelRegistryAPI = Pick<ModelRegistryAPIs, 'patchRegisteredModel'>;

export const bumpModelVersionTimestamp = async (
  api: ModelRegistryAPIs,
  modelVersionId: string,
): Promise<void> => {
  if (!modelVersionId) {
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
          _lastModified: {
            metadataType: ModelRegistryMetadataType.STRING,
            // eslint-disable-next-line camelcase
            string_value: currentTime,
          },
        },
      },
      modelVersionId,
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
  registeredModelId: string,
): Promise<void> => {
  if (!registeredModelId) {
    throw new Error('Registered model ID is required');
  }

  try {
    const currentTime = new Date().toISOString();
    await api.patchRegisteredModel(
      {},
      {
        state: ModelState.LIVE,
        customProperties: {
          // This is a workaround to update the timestamp on the backend. There is a bug opened for model registry team
          // to fix this issue.
          _lastModified: {
            metadataType: ModelRegistryMetadataType.STRING,
            // eslint-disable-next-line camelcase
            string_value: currentTime,
          },
        },
      },
      registeredModelId,
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
  modelVersionId: string,
  registeredModelId: string,
): Promise<void> => {
  await Promise.all([
    bumpModelVersionTimestamp(api, modelVersionId),
    bumpRegisteredModelTimestamp(api, registeredModelId),
  ]);
};
