import type { ModelVersion, RegisteredModel, ModelArtifactList, ModelRegistry } from '~/app/types';
import { getServerAddress } from '~/app/pages/modelRegistry/screens/utils';

// Utility function to get model version tuning data without hook dependencies
export const getModelVersionTuningData = (
  modelVersionId: string | null,
  modelVersion: ModelVersion | null,
  registeredModel: RegisteredModel | null,
  artifacts: ModelArtifactList,
  preferredModelRegistry: ModelRegistry | undefined,
  modelRegistries: ModelRegistry[],
): {
  tuningData: {
    modelRegistryName: string;
    modelRegistryDisplayName: string;
    registeredModelId: string;
    registeredModelName: string;
    modelVersionId: string;
    modelVersionName: string;
    inputModelLocationUri: string;
    outputModelRegistryApiUrl: string;
  } | null;
} => {
  const registryService = modelRegistries.find((s) => s.name === preferredModelRegistry?.name);
  const inputModelLocationUri = artifacts.items[0]?.uri;
  const modelRegistryDisplayName = registryService ? registryService.displayName : '';
  const outputModelRegistryApiUrl = registryService ? getServerAddress(registryService) : '';

  return {
    tuningData:
      modelVersionId && modelVersion && inputModelLocationUri && registryService && registeredModel
        ? {
            modelRegistryName: registryService.name,
            modelRegistryDisplayName,
            registeredModelId: modelVersion.registeredModelId,
            registeredModelName: registeredModel.name,
            modelVersionId: modelVersion.id,
            modelVersionName: modelVersion.name,
            inputModelLocationUri,
            outputModelRegistryApiUrl,
          }
        : null,
  };
};
