/**
 * Temporary development hook for toggling the catalog advanced filters feature.
 *
 * This hook provides a browser storage-backed feature flag that can be toggled
 * via the browser console using window.setTempDevCatalogAdvancedFiltersFeatureAvailable(true/false).
 * The state persists across page reloads using browser storage.
 *
 * This should be removed once the advanced filters feature is ready to be rolled out.
 *
 * @returns {boolean} Whether the catalog advanced filters feature is enabled
 */

import * as React from 'react';
import { useBrowserStorage } from 'mod-arch-core';

declare global {
  interface Window {
    setTempDevCatalogAdvancedFiltersFeatureAvailable?: (enabled: boolean) => void;
  }
}

export const TEMP_DEV_CATALOG_ADVANCED_FILTERS_FEATURE_KEY =
  'tempDevCatalogAdvancedFiltersFeatureAvailable';

export const useTempDevCatalogAdvancedFiltersFeatureAvailable = (): boolean => {
  const [isAvailable, setIsAvailable] = useBrowserStorage(
    TEMP_DEV_CATALOG_ADVANCED_FILTERS_FEATURE_KEY,
    false,
  );

  // Expose setter to window for easy toggling via browser console
  React.useEffect(() => {
    window.setTempDevCatalogAdvancedFiltersFeatureAvailable = setIsAvailable;
  }, [setIsAvailable]);

  return isAvailable;
};
