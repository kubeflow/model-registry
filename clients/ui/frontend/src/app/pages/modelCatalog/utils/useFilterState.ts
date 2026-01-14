import * as React from 'react';
import { ModelCatalogStringFilterKey } from '~/concepts/modelCatalog/const';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import { ModelCatalogFilterStates } from '~/app/modelCatalogTypes';

/**
 * Type for string array filter keys (filters that store string[] values).
 * This hook factory is only intended for filters that use plain string arrays,
 * not filters with specific types like ModelCatalogTask[] or UseCaseOptionValue[].
 */
type StringArrayFilterKey =
  | ModelCatalogStringFilterKey.HARDWARE_TYPE
  | ModelCatalogStringFilterKey.HARDWARE_CONFIGURATION;

/**
 * Creates a generic hook factory for string array filter state.
 * This eliminates duplication across filter state hooks that follow the same pattern.
 *
 * @param filterKey - The filter key to manage state for (must be a string array filter)
 * @returns A hook that provides appliedValues, setAppliedValues, and clearFilters
 */
export const createStringArrayFilterStateHook =
  <K extends StringArrayFilterKey>(filterKey: K) =>
  (): {
    appliedValues: string[];
    setAppliedValues: (values: string[]) => void;
    clearFilters: () => void;
  } => {
    const { filterData, setFilterData } = React.useContext(ModelCatalogContext);
    // filterData[filterKey] is always defined as an array (initialized in ModelCatalogContext)
    const appliedValues: string[] = filterData[filterKey];

    const setAppliedValues = React.useCallback(
      (values: string[]) => {
        // Safe to call setFilterData because K is constrained to StringArrayFilterKey
        // which only includes keys that map to string[] in ModelCatalogFilterStates
        // Using 'as unknown as' to satisfy TypeScript's strict type checking
        // eslint-disable-next-line @typescript-eslint/consistent-type-assertions
        setFilterData(filterKey, values as unknown as ModelCatalogFilterStates[K]);
      },
      // filterKey is a constant parameter from the outer scope, not a reactive dependency
      // eslint-disable-next-line react-hooks/exhaustive-deps
      [setFilterData],
    );

    const clearFilters = React.useCallback(() => {
      // Safe to call setFilterData because K is constrained to StringArrayFilterKey
      // Using 'as unknown as' to satisfy TypeScript's strict type checking
      // eslint-disable-next-line @typescript-eslint/consistent-type-assertions
      setFilterData(filterKey, [] as unknown as ModelCatalogFilterStates[K]);
      // filterKey is a constant parameter from the outer scope, not a reactive dependency
      // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [setFilterData]);

    return { appliedValues, setAppliedValues, clearFilters };
  };
