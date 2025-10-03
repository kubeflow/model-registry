import { CatalogSource, CatalogSourceList } from '~/app/modelCatalogTypes';

export const mockCatalogSource = (partial?: Partial<CatalogSource>): CatalogSource => ({
  id: 'sample-source',
  name: 'sample source',
  enabled: true,
  labels: ['Sample category 1', 'Sample category 2', 'Community'],
  ...partial,
});

export const mockCatalogSourceList = (partial?: Partial<CatalogSourceList>): CatalogSourceList => ({
  items: [mockCatalogSource({})],
  pageSize: 10,
  size: 25,
  nextPageToken: '',
  ...partial,
});
