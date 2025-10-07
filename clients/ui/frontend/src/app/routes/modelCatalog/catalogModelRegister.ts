import { encodeParams } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';

export const getRegisterCatalogModelRoute = (id = '', name = ''): string => {
  const { sourceId = '', modelName = '' } = encodeParams({
    sourceId: id,
    modelName: name,
  });
  return `/model-catalog/${sourceId}/${modelName}/register` || '#';
};
