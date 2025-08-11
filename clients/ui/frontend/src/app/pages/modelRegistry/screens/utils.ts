import { KeyValuePair } from 'mod-arch-core';
import {
  ModelRegistry,
  ModelRegistryCustomProperties,
  ModelRegistryCustomProperty,
  ModelRegistryMetadataType,
  ModelRegistryStringCustomProperties,
  ModelVersion,
  RegisteredModel,
} from '~/app/types';
import { COMPANY_URI } from '~/app/utilities/const';
import { getLastCreatedItem } from '~/app/utils';
import {
  ModelRegistryFilterDataType,
  ModelRegistryVersionsFilterDataType,
} from '~/app/pages/modelRegistry/screens/const';

export type ObjectStorageFields = {
  endpoint: string;
  bucket: string;
  region?: string;
  path: string;
};

// Retrieves the labels from customProperties that have non-empty string_value.
export const getLabels = <T extends ModelRegistryCustomProperties>(customProperties: T): string[] =>
  Object.keys(customProperties).filter((key) => {
    const prop = customProperties[key];
    return prop.metadataType === ModelRegistryMetadataType.STRING && prop.string_value === '';
  });

// Returns the customProperties object with an updated set of labels (non-empty string_value) without affecting other properties.
export const mergeUpdatedLabels = (
  customProperties: ModelRegistryCustomProperties,
  updatedLabels: string[],
): ModelRegistryCustomProperties => {
  const existingLabels = getLabels(customProperties);
  const addedLabels = updatedLabels.filter((label) => !existingLabels.includes(label));
  const removedLabels = existingLabels.filter((label) => !updatedLabels.includes(label));
  const customPropertiesCopy = { ...customProperties };
  removedLabels.forEach((label) => {
    delete customPropertiesCopy[label];
  });
  addedLabels.forEach((label) => {
    customPropertiesCopy[label] = {
      // eslint-disable-next-line camelcase
      string_value: '',
      metadataType: ModelRegistryMetadataType.STRING,
    };
  });
  return customPropertiesCopy;
};

// Retrieves the customProperties that are not special (_registeredFrom) or labels (they have a defined string_value).
export const getProperties = <T extends ModelRegistryCustomProperties>(
  customProperties: T,
): ModelRegistryStringCustomProperties => {
  const initial: ModelRegistryStringCustomProperties = {};
  return Object.keys(customProperties).reduce((acc, key) => {
    // _lastModified is a property that is required to update the timestamp on the backend and we have a workaround for it. It should be resolved by
    // backend team
    if (key === '_lastModified' || /^_registeredFrom/.test(key)) {
      return acc;
    }

    const prop = customProperties[key];
    if (prop.metadataType === ModelRegistryMetadataType.STRING && prop.string_value !== '') {
      return { ...acc, [key]: prop };
    }
    return acc;
  }, initial);
};

// Returns the customProperties object with a single string property added, updated or deleted
export const mergeUpdatedProperty = (
  args: { customProperties: ModelRegistryCustomProperties } & (
    | { op: 'create'; newPair: KeyValuePair }
    | { op: 'update'; oldKey: string; newPair: KeyValuePair }
    | { op: 'delete'; oldKey: string }
  ),
): ModelRegistryCustomProperties => {
  const { op } = args;
  const customPropertiesCopy = { ...args.customProperties };
  if (op === 'delete' || (op === 'update' && args.oldKey !== args.newPair.key)) {
    delete customPropertiesCopy[args.oldKey];
  }
  if (op === 'create' || op === 'update') {
    const { key, value } = args.newPair;
    customPropertiesCopy[key] = {
      // eslint-disable-next-line camelcase
      string_value: value,
      metadataType: ModelRegistryMetadataType.STRING,
    };
  }
  return customPropertiesCopy;
};

export const getCustomPropString = <
  T extends Record<string, ModelRegistryCustomProperty | undefined>,
>(
  customProperties: T,
  key: string,
): string => {
  const prop = customProperties[key];

  if (prop?.metadataType === 'MetadataStringValue') {
    return prop.string_value;
  }
  return '';
};

const isMatchVersionKeyword = (mv: ModelVersion, keywordFilter: string): boolean =>
  mv.name.toLowerCase().includes(keywordFilter) ||
  (mv.description && mv.description.toLowerCase().includes(keywordFilter)) ||
  getLabels(mv.customProperties).some((label) => label.toLowerCase().includes(keywordFilter));

export const filterModelVersions = (
  unfilteredModelVersions: ModelVersion[],
  filterData: ModelRegistryVersionsFilterDataType,
): ModelVersion[] => {
  const keywordFilter = filterData.Keyword?.toLowerCase();
  const authorFilter = filterData.Author?.toLowerCase();

  return unfilteredModelVersions.filter((mv: ModelVersion) => {
    if (!keywordFilter && !authorFilter) {
      return true;
    }

    const doesNotMatchVersion = keywordFilter && !isMatchVersionKeyword(mv, keywordFilter);

    if (doesNotMatchVersion) {
      return false;
    }

    return !authorFilter || mv.author?.toLowerCase().includes(authorFilter);
  });
};

export const sortModelVersionsByCreateTime = (registeredModels: ModelVersion[]): ModelVersion[] =>
  registeredModels.toSorted((a, b) => {
    const first = parseInt(a.createTimeSinceEpoch);
    const second = parseInt(b.createTimeSinceEpoch);
    return new Date(second).getTime() - new Date(first).getTime();
  });

export const filterRegisteredModels = (
  unfilteredRegisteredModels: RegisteredModel[],
  unfilteredModelVersions: ModelVersion[],
  filterData: ModelRegistryFilterDataType,
): RegisteredModel[] => {
  const keywordFilter = filterData.Keyword?.toLowerCase();
  const ownerFilter = filterData.Owner?.toLowerCase();

  return unfilteredRegisteredModels.filter((rm: RegisteredModel) => {
    if (!keywordFilter && !ownerFilter) {
      return true;
    }
    const modelVersions = unfilteredModelVersions.filter((mv) => mv.registeredModelId === rm.id);
    const doesNotMatchModel =
      keywordFilter &&
      !(
        rm.name.toLowerCase().includes(keywordFilter) ||
        (rm.description && rm.description.toLowerCase().includes(keywordFilter)) ||
        getLabels(rm.customProperties).some((label) => label.toLowerCase().includes(keywordFilter))
      );

    const doesNotMatchVersions =
      keywordFilter &&
      !modelVersions.some((mv: ModelVersion) => isMatchVersionKeyword(mv, keywordFilter));

    if (doesNotMatchModel && doesNotMatchVersions) {
      return false;
    }

    return !ownerFilter || rm.owner?.toLowerCase().includes(ownerFilter);
  });
};

export const getServerAddress = (resource: ModelRegistry): string => resource.serverAddress || '';

export const isValidHttpUrl = (value: string): boolean => {
  try {
    const url = new URL(value);
    const isHttp = url.protocol === 'http:' || url.protocol === 'https:';
    // Domain validation
    const domainPattern = /^(?!-)[A-Za-z0-9-]+(\.[A-Za-z0-9-]+)*\.[A-Za-z]{2,}$/;

    return isHttp && domainPattern.test(url.hostname);
  } catch {
    return false;
  }
};

export const isCompanyUri = (uri: string): boolean => uri.startsWith(`${COMPANY_URI}/`);

export const getLatestVersionForRegisteredModel = (
  modelVersions: ModelVersion[],
  rmId: string,
): ModelVersion | undefined => {
  const filteredVersions = modelVersions.filter((mv) => mv.registeredModelId === rmId);
  const latestVersion = getLastCreatedItem(filteredVersions);
  return latestVersion;
};
