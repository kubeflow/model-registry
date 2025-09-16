import { encodeParams } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';

export const getRegisterCatalogModelRoute = (id = '', name = '', repository = ''): string => {
  const {
    sourceId = '',
    repositoryName = '',
    modelName = '',
  } = encodeParams({
    sourceId: id,
    repositoryName: repository,
    modelName: name,
  });
  return `/model-catalog/${sourceId}/${repositoryName}/${modelName}/register` || '#';
};
