export type ModelCatalogFilterResponseType = {
  filters: Record<string, ModelCatalogFilterCategoryResponseType>;
};

export type ModelCatalogFiltersType = Record<string, ModelCatalogFilterCategoryType>;

export type ModelCatalogFilterCategoryResponseType = {
  type: string;
  values: string[];
};

export type ModelCatalogFilterCategoryType = {
  type: string;
  values: Record<string, boolean>;
};
