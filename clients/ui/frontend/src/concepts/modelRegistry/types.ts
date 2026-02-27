export enum ModelSourceKind {
  CATALOG = 'catalog',
  KFP = 'kfp',
  TRANSFER_JOB = 'transfer_job',
}

export type TransferJobParams = {
  jobNamespace: string;
  jobName: string;
};

export type ModelSourceProperties = {
  modelSourceKind?: ModelSourceKind;
  modelSourceClass?: string;
  modelSourceGroup?: string;
  modelSourceName?: string;
  modelSourceId?: string;
};
