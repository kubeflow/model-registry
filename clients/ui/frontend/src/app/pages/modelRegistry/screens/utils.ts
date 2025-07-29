import { KeyValuePair } from 'mod-arch-shared';
import { SearchType } from 'mod-arch-shared/dist/components/DashboardSearchField';
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

export const filterModelVersions = (
  unfilteredModelVersions: ModelVersion[],
  search: string,
  searchType: SearchType,
): ModelVersion[] => {
  const searchLower = search.toLowerCase();

  return unfilteredModelVersions.filter((mv: ModelVersion) => {
    if (!search) {
      return true;
    }

    switch (searchType) {
      case SearchType.KEYWORD:
        return (
          mv.name.toLowerCase().includes(searchLower) ||
          (mv.description && mv.description.toLowerCase().includes(searchLower)) ||
          getLabels(mv.customProperties).some((label) => label.toLowerCase().includes(searchLower))
        );

      case SearchType.AUTHOR: {
        return mv.author && mv.author.toLowerCase().includes(searchLower);
      }

      default:
        return true;
    }
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
  search: string,
  searchType: SearchType,
): RegisteredModel[] => {
  const searchLower = search.toLowerCase();

  return unfilteredRegisteredModels.filter((rm: RegisteredModel) => {
    if (!search) {
      return true;
    }
    const modelVersions = unfilteredModelVersions.filter((mv) => mv.registeredModelId === rm.id);

    switch (searchType) {
      case SearchType.KEYWORD: {
        const matchesModel =
          rm.name.toLowerCase().includes(searchLower) ||
          (rm.description && rm.description.toLowerCase().includes(searchLower)) ||
          getLabels(rm.customProperties).some((label) => label.toLowerCase().includes(searchLower));

        const matchesVersion = modelVersions.some(
          (mv: ModelVersion) =>
            mv.name.toLowerCase().includes(searchLower) ||
            (mv.description && mv.description.toLowerCase().includes(searchLower)) ||
            getLabels(mv.customProperties).some((label) =>
              label.toLowerCase().includes(searchLower),
            ),
        );

        return matchesModel || matchesVersion;
      }
      case SearchType.OWNER: {
        return rm.owner && rm.owner.toLowerCase().includes(searchLower);
      }

      default:
        return true;
    }
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
