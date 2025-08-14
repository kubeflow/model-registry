export enum ModelRegistryFilterOptions {
  keyword = 'Keyword',
  owner = 'Owner',
}

export const modelRegistryFilterOptions = {
  [ModelRegistryFilterOptions.keyword]: 'Keyword',
  [ModelRegistryFilterOptions.owner]: 'Owner',
};

export type ModelRegistryFilterDataType = Record<ModelRegistryFilterOptions, string | undefined>;

export const initialModelRegistryFilterData: ModelRegistryFilterDataType = {
  [ModelRegistryFilterOptions.keyword]: '',
  [ModelRegistryFilterOptions.owner]: '',
};

export enum ModelRegistryVersionsFilterOptions {
  keyword = 'Keyword',
  author = 'Author',
}

export const modelRegistryVersionsFilterOptions = {
  [ModelRegistryVersionsFilterOptions.keyword]: 'Keyword',
  [ModelRegistryVersionsFilterOptions.author]: 'Author',
};

export type ModelRegistryVersionsFilterDataType = Record<
  ModelRegistryVersionsFilterOptions,
  string | undefined
>;

export const initialModelRegistryVersionsFilterData: ModelRegistryVersionsFilterDataType = {
  [ModelRegistryVersionsFilterOptions.keyword]: '',
  [ModelRegistryVersionsFilterOptions.author]: '',
};
