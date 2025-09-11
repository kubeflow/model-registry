/* eslint-disable camelcase */
import { encodeParams } from './utils/modelCatalogUtils';

export const modelCatalogUrl = (sourceId = ''): string => `/model-catalog/${sourceId}`;

export const modelCatalogDetailsUrl = (catalogModelName: string, source_id = ''): string => {
  const parts = catalogModelName.split('/');
  const name = parts[1];
  const repository = parts[0];
  const {
    sourceId = '',
    repositoryName = '',
    modelName = '',
  } = encodeParams({
    sourceId: source_id,
    repositoryName: repository,
    modelName: name,
  });
  return `/model-catalog/${sourceId}/${repositoryName}/${modelName}` || '#';
};
