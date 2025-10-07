import { CatalogModelDetailsParams } from '~/app/modelCatalogTypes';
import { encodeParams } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import { modelCatalogUrl } from './catalogModel';

export const getCatalogModelDetailsRoute = (params: CatalogModelDetailsParams): string => {
  const { sourceId = '', modelName = '' } = encodeParams({
    sourceId: params.sourceId,
    modelName: params.modelName,
  });
  return `${modelCatalogUrl(sourceId)}/${modelName}` || '#';
};
