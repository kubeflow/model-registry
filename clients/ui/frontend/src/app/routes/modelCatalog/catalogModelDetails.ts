export interface CatalogModelDetailsParams {
  modelName: string;
  tag: string;
  sourceName?: string;
  repositoryName?: string;
}

export const getCatalogModelDetailsRoute = (params: CatalogModelDetailsParams): string =>
  // For now, we'll use a simple modelName-based route
  // This should be updated to match the actual routing structure
  `/model-catalog/${encodeURIComponent(params.modelName)}`;
