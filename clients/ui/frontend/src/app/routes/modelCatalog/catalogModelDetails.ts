import { CatalogModelDetailsParams } from '~/app/modelCatalogTypes';
import { encodeParams } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';

export const getCatalogModelDetailsRoute = (params: CatalogModelDetailsParams): string => {
  const { sourceId = '', modelName = '' } = encodeParams({
    sourceId: params.sourceId,
    modelName: params.modelName,
  });
  return `/model-catalog/${sourceId}/${modelName}` || '#';
};
