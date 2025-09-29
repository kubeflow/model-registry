import {
  CatalogArtifacts,
  CatalogArtifactType,
  CatalogModel,
  CatalogModelDetailsParams,
  CatalogSourceList,
} from '~/app/modelCatalogTypes';
import { getLabels } from '~/app/pages/modelRegistry/screens/utils';
import {
  ModelCatalogFilterDataType,
  ModelCatalogFilterResponseType,
} from '~/app/pages/modelCatalog/types';
import {
  ModelCatalogFilterKeys,
  ModelCatalogTasks,
  ModelCatalogLicenses,
  ModelCatalogProviders,
  AllLanguageCodes,
} from '~/concepts/modelCatalog/const';

export const extractVersionTag = (tags?: string[]): string | undefined =>
  tags?.find((tag) => /^\d+\.\d+\.\d+$/.test(tag));
export const filterNonVersionTags = (tags?: string[]): string[] | undefined => {
  const versionTag = extractVersionTag(tags);
  return tags?.filter((tag) => tag !== versionTag);
};

export const getModelName = (modelName: string): string => {
  const index = modelName.indexOf('/');
  if (index === -1) {
    return modelName;
  }
  return modelName.slice(index + 1);
};

export const decodeParams = (
  params: Readonly<CatalogModelDetailsParams>,
): CatalogModelDetailsParams =>
  Object.fromEntries(
    Object.entries(params).map(([key, value]) => [key, decodeURIComponent(value)]),
  );

export const encodeParams = (params: CatalogModelDetailsParams): CatalogModelDetailsParams =>
  Object.fromEntries(
    Object.entries(params).map(([key, value]) => [
      key,
      encodeURIComponent(value).replace(/\./g, '%252E'),
    ]),
  );

export const filterEnabledCatalogSources = (
  catalogSources: CatalogSourceList | null,
): CatalogSourceList | null => {
  if (!catalogSources) {
    return null;
  }

  const filteredItems = catalogSources.items.filter((source) => source.enabled !== false);

  return {
    ...catalogSources,
    items: filteredItems,
    size: filteredItems.length,
  };
};

export const getModelArtifactUri = (artifacts: CatalogArtifacts[]): string => {
  const modelArtifact = artifacts.find(
    (artifact) => artifact.artifactType === CatalogArtifactType.modelArtifact,
  );

  if (modelArtifact) {
    return modelArtifact.uri || '';
  }

  return '';
};

export const hasModelArtifacts = (artifacts: CatalogArtifacts[]): boolean =>
  artifacts.some((artifact) => artifact.artifactType === CatalogArtifactType.modelArtifact);

// Utility function to check if a model is validated
export const isModelValidated = (model: CatalogModel): boolean => {
  if (!model.customProperties) {
    return false;
  }
  const labels = getLabels(model.customProperties);
  return labels.includes('validated');
};

export const getModelCatalogFilters = (): ModelCatalogFilterResponseType => ({
  filters: {
    [ModelCatalogFilterKeys.TASK]: {
      type: 'string',
      values: [
        ModelCatalogTasks.TEXT_GENERATION,
        ModelCatalogTasks.TEXT_TO_TEXT,
        ModelCatalogTasks.IMAGE_TO_TEXT,
      ],
    },
    [ModelCatalogFilterKeys.LICENSE]: {
      type: 'string',
      values: [ModelCatalogLicenses.APACHE_2_0, ModelCatalogLicenses.MIT],
    },
    [ModelCatalogFilterKeys.PROVIDER]: {
      type: 'string',
      values: [ModelCatalogProviders.GOOGLE, ModelCatalogProviders.META],
    },
    [ModelCatalogFilterKeys.LANGUAGE]: {
      type: 'string',
      values: [AllLanguageCodes.EN, AllLanguageCodes.ES, AllLanguageCodes.FR],
    },
  },
});

export const filterModelCatalogModels = (
  models: CatalogModel[],
  filterData: ModelCatalogFilterDataType,
): CatalogModel[] =>
  models.filter((model) =>
    Object.entries(filterData).every(([filterKey, filterState]) => {
      const activeFilters = Object.entries(filterState).filter(([, isActive]) => isActive);

      if (activeFilters.length === 0) {
        return true;
      }

      let modelValue: string | string[] | undefined;
      switch (filterKey) {
        case 'task':
          modelValue = model.tasks;
          break;
        case 'language':
          modelValue = model.language;
          break;
        case 'provider':
          modelValue = model.provider;
          break;
        case 'license':
          modelValue = model.license;
          break;
        default:
          return true;
      }

      if (!modelValue) {
        return false;
      }

      return activeFilters.every(([filterValue]) => {
        if (Array.isArray(modelValue)) {
          return modelValue.includes(filterValue);
        }
        return modelValue === filterValue;
      });
    }),
  );
