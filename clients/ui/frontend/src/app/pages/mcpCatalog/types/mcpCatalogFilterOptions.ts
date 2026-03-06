export type McpFilterCategoryKey =
  | 'deploymentMode'
  | 'supportedTransports'
  | 'license'
  | 'labels'
  | 'securityVerification';

export type McpCatalogFiltersState = {
  [K in McpFilterCategoryKey]?: string[];
};

export type McpCatalogFilterStringOption = {
  type: 'string';
  values?: string[];
};

export type McpCatalogFilterOptions = {
  [key in McpFilterCategoryKey]?: McpCatalogFilterStringOption;
};

export type McpCatalogFilterOptionsList = {
  filters?: McpCatalogFilterOptions;
};
