import * as React from 'react';

// Simple hardware type filter state - separate from ModelCatalog system
export const useHardwareTypeFilterState = (): {
  appliedHardwareTypes: string[];
  setAppliedHardwareTypes: React.Dispatch<React.SetStateAction<string[]>>;
  clearHardwareFilters: () => void;
} => {
  const [appliedHardwareTypes, setAppliedHardwareTypes] = React.useState<string[]>([]);

  return {
    appliedHardwareTypes,
    setAppliedHardwareTypes,
    clearHardwareFilters: () => setAppliedHardwareTypes([]),
  };
};
