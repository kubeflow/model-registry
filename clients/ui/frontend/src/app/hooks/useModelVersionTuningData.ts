import { useContext } from 'react';
import type { ModelVersion, RegisteredModel } from '~/app/types';
import useModelArtifactsByVersionId from '~/app/hooks/useModelArtifactsByVersionId';
import { getServerAddress } from '~/app/pages/modelRegistry/screens/utils';
import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';

export const useModelVersionTuningData = (
  modelVersionId: string | null,
  modelVersion: ModelVersion | null,
  registeredModel: RegisteredModel | null,
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
  loaded: boolean;
  loadError: Error | null;
} => {
  const { preferredModelRegistry, modelRegistries } = useContext(ModelRegistrySelectorContext);
  const registryService = modelRegistries.find((s) => s.name === preferredModelRegistry?.name);

  const [artifacts, loaded, loadError] = useModelArtifactsByVersionId(modelVersionId || undefined);

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
    loaded,
    loadError: loadError || null,
  };
};
