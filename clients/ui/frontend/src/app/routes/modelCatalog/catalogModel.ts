import { getCatalogModelDetailsRoute } from '~/app/routes/modelCatalog/catalogModelDetails';

export const modelCatalogUrl = (): string => '/model-catalog';

export const catalogModelDetailsFromModel = (catalogModelName = '', sourceId = ''): string =>
  getCatalogModelDetailsRoute({ sourceId, modelName: catalogModelName });
