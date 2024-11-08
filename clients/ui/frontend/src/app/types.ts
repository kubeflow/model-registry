import { AlertVariant } from '@patternfly/react-core';
import { APIOptions } from '~/shared/api/types';

export enum ModelState {
  LIVE = 'LIVE',
  ARCHIVED = 'ARCHIVED',
}

export enum ModelArtifactState {
  UNKNOWN = 'UNKNOWN',
  PENDING = 'PENDING',
  LIVE = 'LIVE',
  MARKED_FOR_DELETION = 'MARKED_FOR_DELETION',
  DELETED = 'DELETED',
  ABANDONED = 'ABANDONED',
  REFERENCE = 'REFERENCE',
}

export type ModelRegistry = {
  name: string;
  displayName: string;
  description: string;
};

export type ModelRegistryBody<T> = {
  data: T;
  metadata?: Record<string, unknown>;
};

export enum ModelRegistryMetadataType {
  INT = 'MetadataIntValue',
  DOUBLE = 'MetadataDoubleValue',
  STRING = 'MetadataStringValue',
  STRUCT = 'MetadataStructValue',
  PROTO = 'MetadataProtoValue',
  BOOL = 'MetadataBoolValue',
}

export type ModelRegistryCustomPropertyInt = {
  metadataType: ModelRegistryMetadataType.INT;
  int_value: string; // int64-formatted string
};

export type ModelRegistryCustomPropertyDouble = {
  metadataType: ModelRegistryMetadataType.DOUBLE;
  double_value: number;
};

export type ModelRegistryCustomPropertyString = {
  metadataType: ModelRegistryMetadataType.STRING;
  string_value: string;
};

export type ModelRegistryCustomPropertyStruct = {
  metadataType: ModelRegistryMetadataType.STRUCT;
  struct_value: string; // Base64 encoded bytes for struct value
};

export type ModelRegistryCustomPropertyProto = {
  metadataType: ModelRegistryMetadataType.PROTO;
  type: string; // url describing proto value
  proto_value: string; // Base64 encoded bytes for proto value
};

export type ModelRegistryCustomPropertyBool = {
  metadataType: ModelRegistryMetadataType.BOOL;
  bool_value: boolean;
};

export type ModelRegistryCustomProperty =
  | ModelRegistryCustomPropertyInt
  | ModelRegistryCustomPropertyDouble
  | ModelRegistryCustomPropertyString
  | ModelRegistryCustomPropertyStruct
  | ModelRegistryCustomPropertyProto
  | ModelRegistryCustomPropertyBool;

export type ModelRegistryCustomProperties = Record<string, ModelRegistryCustomProperty>;
export type ModelRegistryStringCustomProperties = Record<string, ModelRegistryCustomPropertyString>;

export type ModelRegistryBase = {
  id: string;
  name: string;
  externalID?: string;
  description?: string;
  createTimeSinceEpoch: string;
  lastUpdateTimeSinceEpoch: string;
  customProperties: ModelRegistryCustomProperties;
};

export type ModelArtifact = ModelRegistryBase & {
  uri?: string;
  state?: ModelArtifactState;
  author?: string;
  modelFormatName?: string;
  storageKey?: string;
  storagePath?: string;
  modelFormatVersion?: string;
  serviceAccountName?: string;
  artifactType: string;
};

export type ModelVersion = ModelRegistryBase & {
  state?: ModelState;
  author?: string;
  registeredModelId: string;
};

export type RegisteredModel = ModelRegistryBase & {
  state?: ModelState;
  owner?: string;
};

export type CreateRegisteredModelData = Omit<
  RegisteredModel,
  'lastUpdateTimeSinceEpoch' | 'createTimeSinceEpoch' | 'id'
>;

export type CreateModelVersionData = Omit<
  ModelVersion,
  'lastUpdateTimeSinceEpoch' | 'createTimeSinceEpoch' | 'id'
>;

export type CreateModelArtifactData = Omit<
  ModelArtifact,
  'lastUpdateTimeSinceEpoch' | 'createTimeSinceEpoch' | 'id'
>;

export type ModelRegistryListParams = {
  size: number;
  pageSize: number;
  nextPageToken: string;
};

export type RegisteredModelList = ModelRegistryListParams & { items: RegisteredModel[] };

export type ModelVersionList = ModelRegistryListParams & { items: ModelVersion[] };

export type ModelArtifactList = ModelRegistryListParams & { items: ModelArtifact[] };

export type CreateRegisteredModel = (
  opts: APIOptions,
  data: CreateRegisteredModelData,
) => Promise<RegisteredModel>;

export type CreateModelVersionForRegisteredModel = (
  opts: APIOptions,
  registeredModelId: string,
  data: CreateModelVersionData,
) => Promise<ModelVersion>;

export type CreateModelArtifactForModelVersion = (
  opts: APIOptions,
  modelVersionId: string,
  data: CreateModelArtifactData,
) => Promise<ModelArtifact>;

export type GetRegisteredModel = (
  opts: APIOptions,
  registeredModelId: string,
) => Promise<RegisteredModel>;

export type GetModelVersion = (opts: APIOptions, modelversionId: string) => Promise<ModelVersion>;

export type GetListRegisteredModels = (opts: APIOptions) => Promise<RegisteredModelList>;

export type GetModelVersionsByRegisteredModel = (
  opts: APIOptions,
  registeredmodelId: string,
) => Promise<ModelVersionList>;

export type GetModelArtifactsByModelVersion = (
  opts: APIOptions,
  modelVersionId: string,
) => Promise<ModelArtifactList>;

export type PatchRegisteredModel = (
  opts: APIOptions,
  data: Partial<RegisteredModel>,
  registeredModelId: string,
) => Promise<RegisteredModel>;

export type PatchModelVersion = (
  opts: APIOptions,
  data: Partial<ModelVersion>,
  modelversionId: string,
) => Promise<ModelVersion>;

export type ModelRegistryAPIs = {
  createRegisteredModel: CreateRegisteredModel;
  createModelVersionForRegisteredModel: CreateModelVersionForRegisteredModel;
  createModelArtifactForModelVersion: CreateModelArtifactForModelVersion;
  getRegisteredModel: GetRegisteredModel;
  getModelVersion: GetModelVersion;
  listRegisteredModels: GetListRegisteredModels;
  getModelVersionsByRegisteredModel: GetModelVersionsByRegisteredModel;
  getModelArtifactsByModelVersion: GetModelArtifactsByModelVersion;
  patchRegisteredModel: PatchRegisteredModel;
  patchModelVersion: PatchModelVersion;
};

export type Notification = {
  id?: number;
  status: AlertVariant;
  title: string;
  message?: React.ReactNode;
  hidden?: boolean;
  read?: boolean;
  timestamp: Date;
};

export enum NotificationActionTypes {
  ADD_NOTIFICATION = 'add_notification',
  DELETE_NOTIFICATION = 'delete_notification',
}

export type NotificationAction =
  | {
      type: NotificationActionTypes.ADD_NOTIFICATION;
      payload: Notification;
    }
  | {
      type: NotificationActionTypes.DELETE_NOTIFICATION;
      payload: { id: Notification['id'] };
    };
