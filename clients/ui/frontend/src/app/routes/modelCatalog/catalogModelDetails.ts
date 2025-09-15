import { CatalogModelDetailsParams } from '~/app/modelCatalogTypes';
import { encodeParams } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';

export const getCatalogModelDetailsRoute = (params: CatalogModelDetailsParams): string => {
  const {
    sourceId = '',
    repositoryName = '',
    modelName = '',
  } = encodeParams({
    sourceId: params.sourceId,
    repositoryName: params.repositoryName,
    modelName: params.modelName,
  });
  return `/model-catalog/${sourceId}/${repositoryName}/${modelName}` || '#';
};
