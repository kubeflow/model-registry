import { CatalogModelArtifact, CatalogModelArtifactList } from '~/app/modelCatalogTypes';

export const mockCatalogModelArtifact = (
  partial?: Partial<CatalogModelArtifact>,
): CatalogModelArtifact => ({
  createTimeSinceEpoch: '1739210683000',
  lastUpdateTimeSinceEpoch: '1739210683000',
  uri: '',
  customProperties: {},
  ...partial,
});

export const mockCatalogModelArtifactList = (
  partial?: Partial<CatalogModelArtifactList>,
): CatalogModelArtifactList => ({
  items: [mockCatalogModelArtifact({})],
  pageSize: 10,
  size: 15,
  nextPageToken: '',
  ...partial,
});
