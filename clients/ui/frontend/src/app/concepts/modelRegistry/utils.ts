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
  if (
    !('Registered from' in customProperties) ||
    customProperties['Registered from'].metadataType !== ModelRegistryMetadataType.STRING ||
    customProperties['Registered from'].string_value !== 'Model catalog'
  ) {
    return null;
  }

  // Extract catalog model information
  if (
    !('Source model' in customProperties) ||
    !('Source model version' in customProperties) ||
    customProperties['Source model'].metadataType !== ModelRegistryMetadataType.STRING ||
    customProperties['Source model version'].metadataType !== ModelRegistryMetadataType.STRING
  ) {
    return null;
  }

  const sourceModel = customProperties['Source model'];
  const sourceModelVersion = customProperties['Source model version'];
  const sourceModelId = customProperties['Source model id'];
  const provider = customProperties.Provider;

  return {
    modelName: sourceModel.string_value,
    tag: sourceModelVersion.string_value,
    sourceName:
      'Provider' in customProperties && provider.metadataType === ModelRegistryMetadataType.STRING
        ? provider.string_value
        : undefined,
    repositoryName:
      'Source model id' in customProperties &&
      sourceModelId.metadataType === ModelRegistryMetadataType.STRING
        ? sourceModelId.string_value
        : undefined,
  };
};
