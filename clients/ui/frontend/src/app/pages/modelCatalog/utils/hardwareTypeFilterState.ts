import * as React from 'react';
import { ModelCatalogStringFilterKey } from '~/concepts/modelCatalog/const';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';

export const useHardwareTypeFilterState = (): {
  appliedHardwareTypes: string[];
  setAppliedHardwareTypes: (hardwareTypes: string[]) => void;
  clearHardwareFilters: () => void;
} => {
  const { filterData, setFilterData } = React.useContext(ModelCatalogContext);
  const appliedHardwareTypes = filterData[ModelCatalogStringFilterKey.HARDWARE_TYPE];

  const setAppliedHardwareTypes = React.useCallback(
    (hardwareTypes: string[]) => {
      setFilterData(ModelCatalogStringFilterKey.HARDWARE_TYPE, hardwareTypes);
    },
    [setFilterData],
  );

  const clearHardwareFilters = React.useCallback(() => {
    setFilterData(ModelCatalogStringFilterKey.HARDWARE_TYPE, []);
  }, [setFilterData]);

  return {
    appliedHardwareTypes,
    setAppliedHardwareTypes,
    clearHardwareFilters,
  };
};
