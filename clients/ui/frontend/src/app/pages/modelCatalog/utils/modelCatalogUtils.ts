import { isEnumMember } from 'mod-arch-core';
import React from 'react';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import {
  CatalogArtifacts,
  CatalogArtifactType,
  CatalogModel,
  CatalogModelDetailsParams,
  CatalogSourceList,
  ModelCatalogFilterStates,
  ModelCatalogStringFilterValueType,
} from '~/app/modelCatalogTypes';
import { getLabels } from '~/app/pages/modelRegistry/screens/utils';
import { ModelCatalogStringFilterKey } from '~/concepts/modelCatalog/const';

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

const isStringFilterValid = <K extends ModelCatalogStringFilterKey>(
  filterKey: K,
  value: ModelCatalogStringFilterValueType[ModelCatalogStringFilterKey][],
): value is ModelCatalogFilterStates[K] => isEnumMember(filterKey, ModelCatalogStringFilterKey);

export const useCatalogStringFilterState = (
  filterKey: ModelCatalogStringFilterKey,
): {
  isSelected: (value: ModelCatalogStringFilterValueType[ModelCatalogStringFilterKey]) => boolean;
  setSelected: (
    value: ModelCatalogStringFilterValueType[ModelCatalogStringFilterKey],
    selected: boolean,
  ) => void;
} => {
  type Value = ModelCatalogStringFilterValueType[ModelCatalogStringFilterKey];
  const { filterData, setFilterData } = React.useContext(ModelCatalogContext);
  const selections: Value[] = filterData[filterKey];
  const isSelected = React.useCallback((value: Value) => selections.includes(value), [selections]);
  const setSelected = (value: Value, selected: boolean) => {
    const nextState: Value[] = selected
      ? [...selections, value]
      : selections.filter((item) => item !== value);
    if (isStringFilterValid(filterKey, nextState)) {
      setFilterData(filterKey, nextState);
    }
  };
  return { isSelected, setSelected };
};
