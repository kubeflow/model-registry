import {
  CatalogArtifacts,
  CatalogArtifactType,
  CatalogModel,
  CatalogModelDetailsParams,
  CatalogSourceList,
} from '~/app/modelCatalogTypes';
import { getLabels } from '~/app/pages/modelRegistry/screens/utils';
import {
  ModelCatalogFilterCategoryResponseType,
  ModelCatalogFilterCategoryType,
  ModelCatalogFilterResponseType,
} from '~/app/pages/modelCatalog/types';

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
    license: {
      type: 'string',
      values: ['apache-2.0'],
    },
    provider: {
      type: 'string',
      values: ['provider1', 'provider2', 'provider3', 'Hugging Face', 'Admin model 1'],
    },
    language: {
      type: 'string',
      values: ['ar', 'cs', 'de', 'en', 'es', 'fr', 'it', 'ja', 'ko', 'nl', 'pt', 'zh'],
    },
  },
});

export const processModelCatalogFilters = (
  filters: Record<string, ModelCatalogFilterCategoryResponseType>,
): Record<string, ModelCatalogFilterCategoryType> =>
  Object.fromEntries(
    Object.entries(filters).map(([key, value]) => [
      key,
      {
        type: value.type,
        values: Object.fromEntries(value.values.map((val) => [val, false])),
      },
    ]),
  );
