import { ModelArtifact, ModelRegistryMetadataType } from '~/app/types';

export interface CatalogModelParams {
  modelName: string;
  tag: string;
  sourceName?: string;
  repositoryName?: string;
}

export const modelSourcePropertiesToCatalogParams = (
  modelArtifact: ModelArtifact,
): CatalogModelParams | null => {
  const { customProperties } = modelArtifact;

  // Check if this was registered from Model catalog
  const registeredFrom = customProperties['Registered from'];
  if (
    registeredFrom.metadataType !== ModelRegistryMetadataType.STRING ||
    registeredFrom.string_value !== 'Model catalog'
  ) {
    return null;
  }

  // Extract catalog model information
  const sourceModel = customProperties['Source model'];
  const sourceModelVersion = customProperties['Source model version'];
  const sourceModelId = customProperties['Source model id'];
  const provider = customProperties.Provider;

  if (
    sourceModel.metadataType !== ModelRegistryMetadataType.STRING ||
    sourceModelVersion.metadataType !== ModelRegistryMetadataType.STRING
  ) {
    return null;
  }

  return {
    modelName: sourceModel.string_value,
    tag: sourceModelVersion.string_value,
    sourceName:
      provider.metadataType === ModelRegistryMetadataType.STRING
        ? provider.string_value
        : undefined,
    repositoryName:
      sourceModelId.metadataType === ModelRegistryMetadataType.STRING
        ? sourceModelId.string_value
        : undefined,
  };
};
