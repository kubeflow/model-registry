export const CATALOG_SETTINGS_PAGE_TITLE = 'Model catalog sources';
export const CATALOG_SETTINGS_DESCRIPTION =
  'Add and manage model sources that populate the model catalog for users in your organization.';

export const ADD_SOURCE_TITLE = 'Add a source';
export const ADD_SOURCE_DESCRIPTION = 'Add a new model catalog source to your organization.';

export const MANAGE_SOURCE_TITLE = 'Manage source';
export const MANAGE_SOURCE_DESCRIPTION = 'Manage the selected model catalog source.';

export const catalogSettingsUrl = (): string => '/model-catalog-settings';

export const addSourceUrl = (): string => `${catalogSettingsUrl()}/add-source`;

export const manageSourceUrl = (catalogSourceId: string): string =>
  `${catalogSettingsUrl()}/manage-source/${encodeURIComponent(catalogSourceId)}`;
