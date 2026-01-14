import { CatalogSourceConfig } from '~/app/modelCatalogTypes';
import { EMPTY_CUSTOM_PROPERTY_VALUE } from '~/concepts/modelCatalog/const';
import { isHuggingFaceSource } from './const';

/**
 * Checks if a catalog source has filters applied
 * @param config - The catalog source configuration
 * @returns true if the source has included or excluded models
 */
export const hasSourceFilters = (config: CatalogSourceConfig): boolean => {
  const hasIncluded = config.includedModels && config.includedModels.length > 0;
  const hasExcluded = config.excludedModels && config.excludedModels.length > 0;
  return !!(hasIncluded || hasExcluded);
};

/**
 * Gets the organization display value for a catalog source
 * @param config - The catalog source configuration
 * @param isDefault - Whether this is a default source
 * @returns The organization name or '-' if not applicable
 */
export const getOrganizationDisplay = (config: CatalogSourceConfig, isDefault: boolean): string => {
  if (isDefault) {
    return EMPTY_CUSTOM_PROPERTY_VALUE;
  }

  if (isHuggingFaceSource(config)) {
    return config.allowedOrganization || EMPTY_CUSTOM_PROPERTY_VALUE;
  }

  return EMPTY_CUSTOM_PROPERTY_VALUE;
};
