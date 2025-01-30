export const modelRegistryUrl = (preferredModelRegistry = ''): string =>
  `/model-registry/${preferredModelRegistry}`;

export const registeredModelsUrl = (preferredModelRegistry?: string): string =>
  `${modelRegistryUrl(preferredModelRegistry)}/registeredModels`;

export const registeredModelUrl = (rmId = '', preferredModelRegistry?: string): string =>
  `${registeredModelsUrl(preferredModelRegistry)}/${rmId}`;

export const registeredModelArchiveUrl = (preferredModelRegistry?: string): string =>
  `${registeredModelsUrl(preferredModelRegistry)}/archive`;

export const registeredModelArchiveDetailsUrl = (
  rmId = '',
  preferredModelRegistry?: string,
): string => `${registeredModelArchiveUrl(preferredModelRegistry)}/${rmId}`;

export const modelVersionListUrl = (rmId?: string, preferredModelRegistry?: string): string =>
  `${registeredModelUrl(rmId, preferredModelRegistry)}/versions`;

export const archiveModelVersionListUrl = (
  rmId?: string,
  preferredModelRegistry?: string,
): string => `${registeredModelArchiveDetailsUrl(rmId, preferredModelRegistry)}/versions`;

export const modelVersionUrl = (
  mvId: string,
  rmId?: string,
  preferredModelRegistry?: string,
): string => `${modelVersionListUrl(rmId, preferredModelRegistry)}/${mvId}`;

export const modelVersionArchiveUrl = (rmId?: string, preferredModelRegistry?: string): string =>
  `${modelVersionListUrl(rmId, preferredModelRegistry)}/archive`;

export const archiveModelVersionDetailsUrl = (
  mvId: string,
  rmId?: string,
  preferredModelRegistry?: string,
): string => `${archiveModelVersionListUrl(rmId, preferredModelRegistry)}/${mvId}`;

export const modelVersionArchiveDetailsUrl = (
  mvId: string,
  rmId?: string,
  preferredModelRegistry?: string,
): string => `${modelVersionArchiveUrl(rmId, preferredModelRegistry)}/${mvId}`;

export const registerModelUrl = (preferredModelRegistry?: string): string =>
  `${modelRegistryUrl(preferredModelRegistry)}/registerModel`;

export const registerVersionUrl = (preferredModelRegistry?: string): string =>
  `${modelRegistryUrl(preferredModelRegistry)}/registerVersion`;

export const registerVersionForModelUrl = (
  rmId?: string,
  preferredModelRegistry?: string,
): string => `${registeredModelUrl(rmId, preferredModelRegistry)}/registerVersion`;

export const modelVersionDeploymentsUrl = (
  mvId: string,
  rmId?: string,
  preferredModelRegistry?: string,
): string => `${modelVersionUrl(mvId, rmId, preferredModelRegistry)}/deployments`;
