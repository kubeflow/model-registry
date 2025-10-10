import { encodeParams } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import { modelCatalogUrl } from './catalogModel';

export const getRegisterCatalogModelRoute = (id = '', name = ''): string => {
  const { sourceId = '', modelName = '' } = encodeParams({
    sourceId: id,
    modelName: name,
  });
  return `${modelCatalogUrl(sourceId)}/${modelName}/register` || '#';
};
