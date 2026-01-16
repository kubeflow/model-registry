import { ModelCatalogStringFilterKey } from '~/concepts/modelCatalog/const';
import { createStringArrayFilterStateHook } from './useFilterState';

/**
 * Hook for managing hardware configuration filter state.
 * Uses the generic filter state hook factory to eliminate duplication.
 */
export const useHardwareConfigurationFilterState = createStringArrayFilterStateHook(
  ModelCatalogStringFilterKey.HARDWARE_CONFIGURATION,
);
