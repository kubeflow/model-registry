import { SearchType } from '~/app/components/DashboardSearchField';
import {
  ModelRegistryCustomProperties,
  ModelRegistryMetadataType,
  ModelRegistryStringCustomProperties,
  ModelState,
  ModelVersion,
  RegisteredModel,
} from '~/app/types';
import { KeyValuePair } from '~/types';

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

// Retrives the customProperties that are not labels (they have a defined string_value).
export const getProperties = <T extends ModelRegistryCustomProperties>(
  customProperties: T,
): ModelRegistryStringCustomProperties => {
  const initial: ModelRegistryStringCustomProperties = {};
  return Object.keys(customProperties).reduce((acc, key) => {
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

export const filterModelVersions = (
  unfilteredModelVersions: ModelVersion[],
  search: string,
  searchType: SearchType,
): ModelVersion[] =>
  unfilteredModelVersions.filter((mv: ModelVersion) => {
    if (!search) {
      return true;
    }

    switch (searchType) {
      case SearchType.KEYWORD:
        return (
          mv.name.toLowerCase().includes(search.toLowerCase()) ||
          (mv.description && mv.description.toLowerCase().includes(search.toLowerCase()))
        );

      case SearchType.AUTHOR:
        return (
          mv.author &&
          (mv.author.toLowerCase().includes(search.toLowerCase()) ||
            (mv.author && mv.author.toLowerCase().includes(search.toLowerCase())))
        );

      default:
        return true;
    }
  });

export const sortModelVersionsByCreateTime = (registeredModels: ModelVersion[]): ModelVersion[] =>
  registeredModels.toSorted((a, b) => {
    const first = parseInt(a.createTimeSinceEpoch);
    const second = parseInt(b.createTimeSinceEpoch);
    return new Date(second).getTime() - new Date(first).getTime();
  });

export const filterRegisteredModels = (
  unfilteredRegisteredModels: RegisteredModel[],
  search: string,
  searchType: SearchType,
): RegisteredModel[] =>
  unfilteredRegisteredModels.filter((rm: RegisteredModel) => {
    if (!search) {
      return true;
    }

    switch (searchType) {
      case SearchType.KEYWORD:
        return (
          rm.name.toLowerCase().includes(search.toLowerCase()) ||
          (rm.description && rm.description.toLowerCase().includes(search.toLowerCase()))
        );

      case SearchType.OWNER:
        return rm.owner && rm.owner.toLowerCase().includes(search.toLowerCase());

      default:
        return true;
    }
  });

export const filterArchiveVersions = (modelVersions: ModelVersion[]): ModelVersion[] =>
  modelVersions.filter((mv) => mv.state === ModelState.ARCHIVED);

export const filterLiveVersions = (modelVersions: ModelVersion[]): ModelVersion[] =>
  modelVersions.filter((mv) => mv.state === ModelState.LIVE);

export const filterArchiveModels = (registeredModels: RegisteredModel[]): RegisteredModel[] =>
  registeredModels.filter((rm) => rm.state === ModelState.ARCHIVED);

export const filterLiveModels = (registeredModels: RegisteredModel[]): RegisteredModel[] =>
  registeredModels.filter((rm) => rm.state === ModelState.LIVE);
