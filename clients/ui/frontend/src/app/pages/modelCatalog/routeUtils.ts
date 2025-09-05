import { CatalogModel } from '~/app/modelCatalogTypes';

export const modelCatalogUrl = (sourceId = ''): string => `/model-catalog/${sourceId}`;

export const modelCatalogDetailsUrl = (model: CatalogModel): string =>
  `/model-catalog/${encodeURIComponent(model.sourceId)}/${model.name}` || '#';
