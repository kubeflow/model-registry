import { CatalogSource, CatalogSourceList } from '~/app/modelCatalogTypes';

export const mockCatalogSource = (partial?: Partial<CatalogSource>): CatalogSource => ({
  id: 'sample-source',
  name: 'sample source',
  enabled: true,
  labels: ['Sample category 1', 'Sample category 2', 'Community'],
  status: {
    state: 'Active',
  },
  ...partial,
});

// Mock source with no status (Starting state)
export const mockCatalogSourceStarting = (): CatalogSource => ({
  id: 'starting-source',
  name: 'Starting Source',
  enabled: true,
  labels: ['Community'],
  // No status field - represents "Starting" state
});

// Mock source with Failed status and invalid credential reason
export const mockCatalogSourceFailedCredential = (): CatalogSource => ({
  id: 'failed-credential-source',
  name: 'Failed Credential Source',
  enabled: true,
  labels: ['Enterprise'],
  status: {
    state: 'Failed',
    reason: 'InvalidCredential',
    message: 'The provided API key is invalid or has expired. Please update your credentials.',
  },
});

// Mock source with Failed status and invalid organization reason
export const mockCatalogSourceFailedOrg = (): CatalogSource => ({
  id: 'failed-org-source',
  name: 'Failed Organization Source',
  enabled: true,
  labels: ['Enterprise'],
  status: {
    state: 'Failed',
    reason: 'InvalidOrganization',
    message:
      "The specified organization 'invalid-org' does not exist or you don't have access to it.",
  },
});

// Mock source with Disabled status
export const mockCatalogSourceDisabled = (): CatalogSource => ({
  id: 'disabled-source',
  name: 'Disabled Source',
  enabled: false,
  labels: ['Community'],
  status: {
    state: 'Disabled',
  },
});

// Mock source with Active status
export const mockCatalogSourceActive = (): CatalogSource => ({
  id: 'active-source',
  name: 'Active Source',
  enabled: true,
  labels: ['Community', 'Enterprise'],
  status: {
    state: 'Active',
  },
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
