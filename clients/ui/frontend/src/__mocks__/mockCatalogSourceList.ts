import { CatalogSource, CatalogSourceList } from '~/app/modelCatalogTypes';

export const mockCatalogSource = (partial?: Partial<CatalogSource>): CatalogSource => ({
  id: 'sample-source',
  name: 'sample source',
  enabled: true,
  labels: ['Sample category 1', 'Sample category 2', 'Community'],
  status: 'available',
  ...partial,
});

// Mock source with no status (Starting state - no status field)
export const mockCatalogSourceStarting = (): CatalogSource => ({
  id: 'starting-source',
  name: 'Starting Source',
  enabled: true,
  labels: ['Community'],
  // No status field - represents "Starting" state
});

// Mock source with error status and error message (invalid credential)
export const mockCatalogSourceFailedCredential = (): CatalogSource => ({
  id: 'failed-credential-source',
  name: 'Failed Credential Source',
  enabled: true,
  labels: ['Enterprise'],
  status: 'error',
  error: 'The provided API key is invalid or has expired. Please update your credentials.',
});

// Mock source with error status and error message (invalid organization)
export const mockCatalogSourceFailedOrg = (): CatalogSource => ({
  id: 'failed-org-source',
  name: 'Failed Organization Source',
  enabled: true,
  labels: ['Enterprise'],
  status: 'error',
  error: "The specified organization 'invalid-org' does not exist or you don't have access to it.",
});

// Mock source with disabled status
export const mockCatalogSourceDisabled = (): CatalogSource => ({
  id: 'disabled-source',
  name: 'Disabled Source',
  enabled: false,
  labels: ['Community'],
  status: 'disabled',
});

// Mock source with available status
export const mockCatalogSourceActive = (): CatalogSource => ({
  id: 'active-source',
  name: 'Active Source',
  enabled: true,
  labels: ['Community', 'Enterprise'],
  status: 'available',
});

export const mockCatalogSourceList = (partial?: Partial<CatalogSourceList>): CatalogSourceList => ({
  items: [
    mockCatalogSourceActive(),
    mockCatalogSourceStarting(),
    mockCatalogSourceFailedCredential(),
    mockCatalogSourceFailedOrg(),
    mockCatalogSourceDisabled(),
  ],
  pageSize: 10,
  size: 25,
  nextPageToken: '',
  ...partial,
});
