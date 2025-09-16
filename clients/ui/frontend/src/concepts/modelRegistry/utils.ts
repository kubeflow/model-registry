import {
  ModelRegistryCustomProperties,
  ModelRegistryCustomProperty,
  ModelRegistryCustomPropertyString,
  ModelRegistryMetadataType,
} from '~/app/types';
import { CatalogModel, CatalogModelDetailsParams } from '~/app/modelCatalogTypes';
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
    !properties.modelSourceName
  ) {
    return null;
  }

  return {
    sourceId: properties.modelSourceClass,
    repositoryName: properties.modelSourceGroup,
    modelName: properties.modelSourceName,
  };
};

export const catalogParamsToModelSourceProperties = (
  params: CatalogModelDetailsParams,
): ModelSourceProperties => ({
  modelSourceKind: ModelSourceKind.CATALOG,
  modelSourceClass: params.sourceId,
  modelSourceGroup: params.repositoryName,
  modelSourceName: params.modelName,
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
export const getLabelsFromModelTasks = (
  model: CatalogModel | null,
): ModelRegistryCustomProperties => {
  const tasks = model?.tasks?.reduce<ModelRegistryCustomProperties>((acc, cur) => {
    acc[cur] = EMPTY_CUSTOM_PROPERTY_STRING;
    return acc;
  }, {});

  return { ...tasks };
};

const isStringProperty = (
  prop: ModelRegistryCustomProperty,
): prop is ModelRegistryCustomPropertyString =>
  prop.metadataType === ModelRegistryMetadataType.STRING && prop.string_value === '';

export const getLabelsFromCustomProperties = (
  customProperties?: ModelRegistryCustomProperties,
): Record<string, ModelRegistryCustomPropertyString> => {
  const filteredProperties: Record<string, ModelRegistryCustomPropertyString> = {};

  if (!customProperties) {
    return filteredProperties;
  }

  Object.keys(customProperties).forEach((key) => {
    const prop = customProperties[key];
    if (isStringProperty(prop)) {
      filteredProperties[key] = prop;
    }
  });

  return filteredProperties;
};
