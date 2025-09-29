import {
  ModelCatalogTasks,
  ModelCatalogProviders,
  ModelCatalogLicenses,
  AllLanguageCodes,
  ModelCatalogFilterKeys,
} from '~/concepts/modelCatalog/const';

export interface ModelCatalogTasksFilterType {
  type: 'string';
  values: ModelCatalogTasks[];
}

export interface ModelCatalogProvidersFilterType {
  type: 'string';
  values: ModelCatalogProviders[];
}

export interface ModelCatalogLicensesFilterType {
  type: 'string';
  values: ModelCatalogLicenses[];
}

export interface ModelCatalogLanguagesFilterType {
  type: 'string';
  values: AllLanguageCodes[];
}

export type ModelCatalogFilterTypesByKey = {
  [ModelCatalogFilterKeys.TASK]: ModelCatalogTasksFilterType;
  [ModelCatalogFilterKeys.PROVIDER]: ModelCatalogProvidersFilterType;
  [ModelCatalogFilterKeys.LICENSE]: ModelCatalogLicensesFilterType;
  [ModelCatalogFilterKeys.LANGUAGE]: ModelCatalogLanguagesFilterType;
};

export type ModelCatalogFilterState<K extends ModelCatalogFilterKeys> = Partial<
  Record<ModelCatalogFilterTypesByKey[K]['values'][number], boolean>
>;

export type ModelCatalogTasksFilterStateType = ModelCatalogFilterState<ModelCatalogFilterKeys.TASK>;

export type ModelCatalogProvidersFilterStateType =
  ModelCatalogFilterState<ModelCatalogFilterKeys.PROVIDER>;

export type ModelCatalogLicensesFilterStateType =
  ModelCatalogFilterState<ModelCatalogFilterKeys.LICENSE>;

export type ModelCatalogLanguagesFilterStateType =
  ModelCatalogFilterState<ModelCatalogFilterKeys.LANGUAGE>;

export type ModelCatalogFilterStatesByKey = {
  [K in ModelCatalogFilterKeys]: ModelCatalogFilterState<K>;
};

export type ModelCatalogFilterResponseType = {
  filters: Partial<ModelCatalogFilterTypesByKey>;
};

export type ModelCatalogFilterDataType = Partial<ModelCatalogFilterStatesByKey>;
