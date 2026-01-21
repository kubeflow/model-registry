import { GenericObjectState } from 'mod-arch-core';
import useGenericObjectState from 'mod-arch-core/dist/utilities/useGenericObjectState';
import { ModelRegistryCustomProperties, ModelArtifact } from '~/app/types';

export enum ModelLocationType {
  ObjectStorage = 'Object storage',
  URI = 'URI',
}

export enum DestinationStorageType {
  S3 = 'S3',
  OCI = 'OCI',
}

export type RegistrationCommonFormData = {
  versionName: string;
  versionDescription: string;
  sourceModelFormat: string;
  sourceModelFormatVersion: string;
  modelLocationType: ModelLocationType;
  modelLocationEndpoint: string;
  modelLocationBucket: string;
  modelLocationRegion: string;
  modelLocationPath: string;
  modelLocationURI: string;
  destinationStorageType: DestinationStorageType;
  destinationS3AccessKeyId: string;
  destinationS3SecretAccessKey: string;
  destinationS3Endpoint: string;
  destinationS3Bucket: string;
  destinationS3Region: string;
  destinationS3Path: string;
  destinationOciRegistry: string;
  destinationOciUsername: string;
  destinationOciPassword: string;
  destinationOciUri: string;
  destinationOciEmail: string;
  versionCustomProperties?: ModelRegistryCustomProperties;
  modelCustomProperties?: ModelRegistryCustomProperties;
  additionalArtifactProperties?: Partial<ModelArtifact>;
};

export type RegisterModelFormData = RegistrationCommonFormData & {
  modelName: string;
  modelDescription: string;
};

export type RegisterVersionFormData = RegistrationCommonFormData & {
  registeredModelId: string;
};

export type RegisterCatalogModelFormData = RegisterModelFormData & {
  modelRegistry: string;
};

const registrationCommonFormDataDefaults: RegistrationCommonFormData = {
  versionName: '',
  versionDescription: '',
  sourceModelFormat: '',
  sourceModelFormatVersion: '',
  modelLocationType: ModelLocationType.ObjectStorage,
  modelLocationEndpoint: '',
  modelLocationBucket: '',
  modelLocationRegion: '',
  modelLocationPath: '',
  modelLocationURI: '',
  destinationStorageType: DestinationStorageType.S3,
  destinationS3AccessKeyId: '',
  destinationS3SecretAccessKey: '',
  destinationS3Endpoint: '',
  destinationS3Bucket: '',
  destinationS3Region: '',
  destinationS3Path: '',
  destinationOciRegistry: '',
  destinationOciUsername: '',
  destinationOciPassword: '',
  destinationOciUri: '',
  destinationOciEmail: '',
  modelCustomProperties: {},
  versionCustomProperties: {},
};

const registerModelFormDataDefaults: RegisterModelFormData = {
  ...registrationCommonFormDataDefaults,
  modelName: '',
  modelDescription: '',
};

const registerVersionFormDataDefaults: RegisterVersionFormData = {
  ...registrationCommonFormDataDefaults,
  registeredModelId: '',
};

const registerModelFormDataDefaultsForModelCatalog: RegisterCatalogModelFormData = {
  ...registerModelFormDataDefaults,
  modelRegistry: '',
};

export const useRegisterModelData = (): GenericObjectState<RegisterModelFormData> =>
  useGenericObjectState<RegisterModelFormData>(registerModelFormDataDefaults);

export const useRegisterVersionData = (
  registeredModelId?: string,
): GenericObjectState<RegisterVersionFormData> =>
  useGenericObjectState<RegisterVersionFormData>({
    ...registerVersionFormDataDefaults,
    registeredModelId: registeredModelId || '',
  });

export const useRegisterCatalogModelData = (
  initialData?: Partial<RegisterCatalogModelFormData>,
): GenericObjectState<RegisterCatalogModelFormData> =>
  useGenericObjectState<RegisterCatalogModelFormData>({
    ...registerModelFormDataDefaultsForModelCatalog,
    ...initialData,
  });
