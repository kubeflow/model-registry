import * as React from 'react';

// Simple hardware type filter state - separate from ModelCatalog system
export const useHardwareTypeFilterState = () => {
  const [appliedHardwareTypes, setAppliedHardwareTypes] = React.useState<string[]>([]);

  return {
    appliedHardwareTypes,
    setAppliedHardwareTypes,
    clearHardwareFilters: () => setAppliedHardwareTypes([]),
  };
};
