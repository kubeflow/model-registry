export const getRegisterCatalogModelRoute = (modelId: string): string =>
  `/model-catalog/${encodeURIComponent(modelId)}/register`;
