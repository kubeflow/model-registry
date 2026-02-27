import {
  ModelRegistryCustomProperties,
  ModelRegistryMetadataType,
  ModelState,
  ModelVersion,
  ModelTransferJob,
  ModelTransferJobSourceType,
  ModelTransferJobDestinationType,
  RegisteredModel,
} from '~/app/types';
import { EMPTY_CUSTOM_PROPERTY_VALUE } from '~/concepts/modelCatalog/const';

export type StorageType = ModelTransferJobSourceType | ModelTransferJobDestinationType;

export type ObjectStorageFields = {
  endpoint: string;
  bucket: string;
  region?: string;
  path: string;
};

export type RegisteredModelLocation = {
  s3Fields: ObjectStorageFields | null;
  uri: string | null;
  ociUri: string | null;
} | null;

export const objectStorageFieldsToUri = (fields: ObjectStorageFields): string | null => {
  const { endpoint, bucket, region, path } = fields;
  if (!endpoint || !bucket || !path) {
    return null;
  }
  const searchParams = new URLSearchParams();
  searchParams.set('endpoint', endpoint);
  if (region) {
    searchParams.set('defaultRegion', region);
  }
  return `s3://${bucket}/${path}?${searchParams.toString()}`;
};

export const uriToStorageFields = (uri: string): RegisteredModelLocation => {
  try {
    const urlObj = new URL(uri);
    if (urlObj.toString().startsWith('s3:')) {
      // Some environments include the first token after the protocol (our bucket) in the pathname and some have it as the hostname
      const [bucket, ...pathSplit] = [urlObj.hostname, ...urlObj.pathname.split('/')].filter(
        Boolean,
      );
      const path = pathSplit.join('/');
      const searchParams = new URLSearchParams(urlObj.search);
      const endpoint = searchParams.get('endpoint');
      const region = searchParams.get('defaultRegion');
      if (endpoint && bucket && path) {
        return {
          s3Fields: { endpoint, bucket, region: region || undefined, path },
          uri: null,
          ociUri: null,
        };
      }
      return null;
    }
    if (uri.startsWith('oci:')) {
      return { s3Fields: null, uri: null, ociUri: uri };
    }
    return { s3Fields: null, uri, ociUri: null };
  } catch {
    return null;
  }
};

export const getLastCreatedItem = <T extends { createTimeSinceEpoch?: string }>(
  items?: T[],
): T | undefined =>
  items?.toSorted(
    ({ createTimeSinceEpoch: createTimeA }, { createTimeSinceEpoch: createTimeB }) => {
      if (!createTimeA || !createTimeB) {
        return 0;
      }
      return Number(createTimeB) - Number(createTimeA);
    },
  )[0];

export const filterArchiveVersions = (modelVersions: ModelVersion[]): ModelVersion[] =>
  modelVersions.filter((mv) => mv.state === ModelState.ARCHIVED);

export const filterLiveVersions = (modelVersions: ModelVersion[]): ModelVersion[] =>
  modelVersions.filter((mv) => mv.state === ModelState.LIVE);

export const filterArchiveModels = (registeredModels: RegisteredModel[]): RegisteredModel[] =>
  registeredModels.filter((rm) => rm.state === ModelState.ARCHIVED);

export const filterLiveModels = (registeredModels: RegisteredModel[]): RegisteredModel[] =>
  registeredModels.filter((rm) => rm.state === ModelState.LIVE);

export const getStringValue = <T extends ModelRegistryCustomProperties>(
  customProperties: T | undefined,
  key: keyof T,
): string => {
  const prop = customProperties?.[key];
  if (prop && prop.metadataType === ModelRegistryMetadataType.STRING) {
    return prop.string_value;
  }
  return EMPTY_CUSTOM_PROPERTY_VALUE;
};
export const getIntValue = <T extends ModelRegistryCustomProperties>(
  customProperties: T | undefined,
  key: keyof T,
): number => {
  const prop = customProperties?.[key];
  if (prop && prop.metadataType === ModelRegistryMetadataType.INT) {
    const value = prop.int_value;
    return value ? parseInt(value, 10) : 0;
  }
  return 0;
};
export const getDoubleValue = <T extends ModelRegistryCustomProperties>(
  customProperties: T | undefined,
  key: keyof T,
): number => {
  const prop = customProperties?.[key];
  if (prop && prop.metadataType === ModelRegistryMetadataType.DOUBLE) {
    return prop.double_value;
  }
  return 0;
};

export const getStorageTypeLabel = (type: StorageType | string): string => {
  switch (type) {
    case ModelTransferJobSourceType.S3:
    case ModelTransferJobDestinationType.S3:
      return 'S3';
    case ModelTransferJobSourceType.OCI:
    case ModelTransferJobDestinationType.OCI:
      return 'OCI';
    case ModelTransferJobSourceType.URI:
      return 'URI';
    default:
      return type;
  }
};

export const getModelUriPopoverContent = (destType: StorageType | string): string => {
  switch (destType) {
    case ModelTransferJobDestinationType.OCI:
      return 'The URI of the OCI connection that is currently being used as the model storage location.';
    case ModelTransferJobDestinationType.S3:
      return 'The path of the S3-compatible object storage connection that is currently being used as the model storage location.';
    default:
      return 'The URI of the URI connection that is currently being used as the model storage location.';
  }
};

export const getModelUriLabel = (destType: StorageType | string): string => {
  switch (destType) {
    case ModelTransferJobDestinationType.S3:
      return 'Path';
    default:
      return 'Model URI';
  }
};

export const getDestinationUri = (job: ModelTransferJob): string => {
  const { destination } = job;
  if ('uri' in destination && destination.uri) {
    return destination.uri;
  }
  if ('key' in destination && destination.key) {
    return destination.key;
  }
  return '';
};

export const getSourcePath = (job: ModelTransferJob): string => {
  const { source } = job;
  if (source.type === ModelTransferJobSourceType.S3 && 'key' in source) {
    return source.key || '';
  }
  if ('uri' in source && source.uri) {
    return source.uri;
  }
  return '';
};
