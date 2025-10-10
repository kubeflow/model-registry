import { ModelDetailsTab } from '~/app/pages/modelCatalog/screens/ModelDetailsTabs';
import { getCatalogModelDetailsRoute } from '~/app/routes/modelCatalog/catalogModelDetails';

export const modelCatalogUrl = (sourceId?: string): string =>
  `/model-catalog${sourceId ? `/${sourceId}` : ''}`;

export const catalogModelDetailsFromModel = (catalogModelName = '', sourceId = ''): string =>
  getCatalogModelDetailsRoute({ sourceId, modelName: catalogModelName });

export const catalogModelDetailsTabFromModel = (
  tab: ModelDetailsTab,
  catalogModelName = '',
  sourceId = '',
): string =>
  `${getCatalogModelDetailsRoute({ sourceId, modelName: catalogModelName })}/${encodeURIComponent(
    tab,
  )}`;
