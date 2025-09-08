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
  // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
  if (
    // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
    !registeredFrom ||
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

  // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
  if (
    // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
    !sourceModel ||
    // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
    !sourceModelVersion ||
    sourceModel.metadataType !== ModelRegistryMetadataType.STRING ||
    sourceModelVersion.metadataType !== ModelRegistryMetadataType.STRING
  ) {
    return null;
  }

  return {
    modelName: sourceModel.string_value,
    tag: sourceModelVersion.string_value,
    sourceName:
      // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
      provider && provider.metadataType === ModelRegistryMetadataType.STRING
        ? provider.string_value
        : undefined,
    repositoryName:
      // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
      sourceModelId && sourceModelId.metadataType === ModelRegistryMetadataType.STRING
        ? sourceModelId.string_value
        : undefined,
  };
};
