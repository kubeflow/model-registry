import * as React from 'react';
import { ModelCatalogStringFilterKey } from '~/concepts/modelCatalog/const';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';

export const useHardwareConfigurationFilterState = (): {
  appliedHardwareConfigurations: string[];
  setAppliedHardwareConfigurations: (hardwareConfigurations: string[]) => void;
  clearHardwareFilters: () => void;
} => {
  const { filterData, setFilterData } = React.useContext(ModelCatalogContext);
  const appliedHardwareConfigurations =
    filterData[ModelCatalogStringFilterKey.HARDWARE_CONFIGURATION];

  const setAppliedHardwareConfigurations = React.useCallback(
    (hardwareConfigurations: string[]) => {
      setFilterData(ModelCatalogStringFilterKey.HARDWARE_CONFIGURATION, hardwareConfigurations);
    },
    [setFilterData],
  );

  const clearHardwareFilters = React.useCallback(() => {
    setFilterData(ModelCatalogStringFilterKey.HARDWARE_CONFIGURATION, []);
  }, [setFilterData]);

  return {
    appliedHardwareConfigurations,
    setAppliedHardwareConfigurations,
    clearHardwareFilters,
  };
};
