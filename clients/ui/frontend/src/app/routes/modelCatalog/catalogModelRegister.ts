import { encodeParams } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';

export const getRegisterCatalogModelRoute = (id = '', name = ''): string => {
  const { sourceId = '', modelName = '' } = encodeParams({
    sourceId: id,
    modelName: name,
  });
  return `/ai-hub/catalog/${sourceId}/${modelName}/register` || '#';
};
