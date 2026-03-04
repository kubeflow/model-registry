export type McpFilterCategoryKey =
  | 'deploymentMode'
  | 'supportedTransports'
  | 'license'
  | 'labels'
  | 'securityVerification';

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
