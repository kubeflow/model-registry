import { getCatalogModelDetailsRoute } from '~/app/routes/modelCatalog/catalogModelDetails';

export const modelCatalogUrl = (sourceId = ''): string => `/model-catalog/${sourceId}`;

export const catalogModelDetailsFromModel = (catalogModelName = '', sourceId = ''): string =>
  getCatalogModelDetailsRoute({ sourceId, modelName: catalogModelName });
