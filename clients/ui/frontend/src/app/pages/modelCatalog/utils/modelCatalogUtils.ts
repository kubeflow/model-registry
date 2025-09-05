import { CatalogModel, CatalogModelList } from '~/app/modelCatalogTypes';

export const extractVersionTag = (tags?: string[]): string | undefined =>
  tags?.find((tag) => /^\d+\.\d+\.\d+$/.test(tag));
export const filterNonVersionTags = (tags?: string[]): string[] | undefined => {
  const versionTag = extractVersionTag(tags);
  return tags?.filter((tag) => tag !== versionTag);
};

export const hasModelName = (baseName: string, targetName: string): boolean => {
  const name = baseName.split('/');
  if (name.length < 2) return false;
  return name[1].includes(targetName);
};

export const findCatalogModel = (
  catalogModels: CatalogModelList,
  sourceId: string,
  modelName: string,
): CatalogModel | null => {
  const model = catalogModels.items.find(
    (m) => m.sourceId === sourceId && hasModelName(m.name, modelName),
  );
  return model || null;
};

export const getModelName = (modelName: string): string => {
  const index = modelName.indexOf('/');
  if (index === -1) {
    return modelName;
  }
  return modelName.slice(index + 1);
};
