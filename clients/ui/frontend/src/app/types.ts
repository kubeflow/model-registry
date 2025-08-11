import { APIOptions } from 'mod-arch-core';
import { ModelSourceProperties } from '~/concepts/modelRegistry/types';

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
  serverAddress?: string;
};

export type ModelRegistryPayload = {
  modelRegistry: {
    metadata: {
      name: string;
      annotations: {
        'openshift.io/display-name': string;
        'openshift.io/description': string;
      };
    };
    spec: {
      mysql: {
        host: string;
        port: number;
        username: string;
        database: string;
      };
    };
  };
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

export type ModelArtifact = ModelRegistryBase &
  ModelSourceProperties & {
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
  registeredModel: RegisteredModel,
  isFirstVersion?: boolean,
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

export type GetListModelVersions = (opts: APIOptions) => Promise<ModelVersionList>;

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

export type PatchModelArtifact = (
  opts: APIOptions,
  data: Partial<ModelArtifact>,
  modelartifactId: string,
) => Promise<ModelArtifact>;

export type ModelRegistryAPIs = {
  createRegisteredModel: CreateRegisteredModel;
  createModelVersionForRegisteredModel: CreateModelVersionForRegisteredModel;
  createModelArtifactForModelVersion: CreateModelArtifactForModelVersion;
  getRegisteredModel: GetRegisteredModel;
  getModelVersion: GetModelVersion;
  listModelVersions: GetListModelVersions;
  listRegisteredModels: GetListRegisteredModels;
  getModelVersionsByRegisteredModel: GetModelVersionsByRegisteredModel;
  getModelArtifactsByModelVersion: GetModelArtifactsByModelVersion;
  patchRegisteredModel: PatchRegisteredModel;
  patchModelVersion: PatchModelVersion;
  patchModelArtifact: PatchModelArtifact;
};
