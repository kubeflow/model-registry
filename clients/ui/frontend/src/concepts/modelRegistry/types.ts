export enum ModelSourceKind {
  CATALOG = 'catalog',
  KFP = 'kfp',
}

export type ModelSourceProperties = {
  modelSourceKind?: ModelSourceKind;
  modelSourceClass?: string;
  modelSourceGroup?: string;
  modelSourceName?: string;
  modelSourceId?: string;
};
