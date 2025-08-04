import { CatalogModelDetailsParams } from '~/app/pages/modelRegistry/screens/types';
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
