import { CatalogModelDetailsParams } from '~/app/pages/modelRegistry/screens/types';
import { ModelRegistryCustomProperties, ModelRegistryMetadataType } from '~/app/types';
import { ModelSourceKind, ModelSourceProperties } from './types';

/**
 * Converts model source properties to catalog parameters
 * @param properties - The model source properties
 * @returns CatalogModelDetailsParams object or null if not a catalog source or if required properties are missing
 */
export const modelSourcePropertiesToCatalogParams = (
  properties: ModelSourceProperties,
): CatalogModelDetailsParams | null => {
  if (
    properties.modelSourceKind !== ModelSourceKind.CATALOG ||
    !properties.modelSourceClass ||
    !properties.modelSourceGroup ||
    !properties.modelSourceName ||
    !properties.modelSourceId
  ) {
    return null;
  }

  return {
    sourceName: properties.modelSourceClass,
    repositoryName: properties.modelSourceGroup,
    modelName: properties.modelSourceName,
    tag: properties.modelSourceId,
  };
};

export const catalogParamsToModelSourceProperties = (
  params: CatalogModelDetailsParams,
): ModelSourceProperties => ({
  modelSourceKind: ModelSourceKind.CATALOG,
  modelSourceClass: params.sourceName || '',
  modelSourceGroup: params.repositoryName || '',
  modelSourceName: params.modelName,
  modelSourceId: params.tag,
});

export const modelSourcePropertiesToCustomProperties = (
  properties: ModelSourceProperties,
): ModelRegistryCustomProperties => {
  const customProperties: ModelRegistryCustomProperties = {};

  if (properties.modelSourceKind === ModelSourceKind.CATALOG) {
    customProperties['Registered from'] = {
      // eslint-disable-next-line camelcase
      string_value: 'Model catalog',
      metadataType: ModelRegistryMetadataType.STRING,
    };

    if (properties.modelSourceName) {
      customProperties['Source model'] = {
        // eslint-disable-next-line camelcase
        string_value: properties.modelSourceName,
        metadataType: ModelRegistryMetadataType.STRING,
      };
    }

    if (properties.modelSourceId) {
      customProperties['Source model version'] = {
        // eslint-disable-next-line camelcase
        string_value: properties.modelSourceId,
        metadataType: ModelRegistryMetadataType.STRING,
      };
    }

    if (properties.modelSourceClass) {
      customProperties.Provider = {
        // eslint-disable-next-line camelcase
        string_value: properties.modelSourceClass,
        metadataType: ModelRegistryMetadataType.STRING,
      };
    }

    if (properties.modelSourceGroup) {
      customProperties['Source model id'] = {
        // eslint-disable-next-line camelcase
        string_value: properties.modelSourceGroup,
        metadataType: ModelRegistryMetadataType.STRING,
      };
    }
  }

  return customProperties;
};
