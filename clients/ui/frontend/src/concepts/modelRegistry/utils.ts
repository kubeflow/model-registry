import { CatalogModelDetailsParams } from '~/app/pages/modelRegistry/screens/types';
import { ModelRegistryCustomProperties, ModelRegistryMetadataType } from '~/app/types';
import { ModelCatalogItem } from '~/app/modelCatalogTypes';
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

const EMPTY_CUSTOM_PROPERTY_STRING = {
  // eslint-disable-next-line camelcase
  string_value: '',
  metadataType: ModelRegistryMetadataType.STRING,
} as const;

/**
 * Creates custom properties from a catalog model
 * @param model - The catalog model item
 * @returns ModelRegistryCustomProperties object with labels and tasks
 */
export const createCustomPropertiesFromModel = (
  model: ModelCatalogItem,
): ModelRegistryCustomProperties => {
  const labels = (model.tags || []).reduce<ModelRegistryCustomProperties>((acc, cur) => {
    acc[cur] = EMPTY_CUSTOM_PROPERTY_STRING;
    return acc;
  }, {});

  // Add single Task property from model.task (like versionCustomProperties)
  const taskProperty: ModelRegistryCustomProperties = {};
  if (model.task) {
    taskProperty.Task = {
      // eslint-disable-next-line camelcase
      string_value: model.task,
      metadataType: ModelRegistryMetadataType.STRING,
    };
  }

  return { ...labels, ...taskProperty };
};
