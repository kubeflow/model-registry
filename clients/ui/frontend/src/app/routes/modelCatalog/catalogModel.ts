import { getCatalogModelDetailsRoute } from '~/app/routes/modelCatalog/catalogModelDetails';

export const modelCatalogUrl = (sourceId = ''): string => `/model-catalog/${sourceId}`;

export const catalogModelDetailsFromModel = (catalogModelName = '', sourceId = ''): string => {
  const parts = catalogModelName.split('/');
  return getCatalogModelDetailsRoute({ sourceId, modelName: parts[1], repositoryName: parts[0] });
};
