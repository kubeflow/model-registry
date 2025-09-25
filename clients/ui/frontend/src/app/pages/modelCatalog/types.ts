export type ModelCatalogFilterResponseType = {
  filters: Record<string, ModelCatalogFilterCategoryType | undefined>;
};

export type ModelCatalogFilterCategoryType = ModelCatalogStringFilterType;

export type ModelCatalogFilterDataType = Record<string, ModelCatalogStringFilterStateType>;

export type ModelCatalogStringFilterType = {
  type: 'string';
  values: string[];
};

export type ModelCatalogStringFilterStateType = Record<string, boolean>;
