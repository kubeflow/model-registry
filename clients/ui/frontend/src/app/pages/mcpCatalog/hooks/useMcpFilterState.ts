import * as React from 'react';
import { McpCatalogContext } from '~/app/context/mcpCatalog/McpCatalogContext';
import type { McpFilterCategoryKey } from '~/app/pages/mcpCatalog/types/mcpCatalogFilterOptions';

export function useMcpFilterState(filterKey: McpFilterCategoryKey): {
  selectedValues: string[];
  setSelected: (value: string, checked: boolean) => void;
  isSelected: (value: string) => boolean;
} {
  const { filters, setFilters } = React.useContext(McpCatalogContext);
  const selected = React.useMemo(() => {
    const v = filters[filterKey];
    return Array.isArray(v) ? v : [];
  }, [filters, filterKey]);
  const setSelected = React.useCallback(
    (value: string, checked: boolean) => {
      setFilters((prev) => {
        const current = prev[filterKey];
        const arr = Array.isArray(current) ? current : [];
        if (checked) {
          return { ...prev, [filterKey]: [...arr, value] };
        }
        return { ...prev, [filterKey]: arr.filter((x) => x !== value) };
      });
    },
    [filterKey, setFilters],
  );
  const isSelected = React.useCallback((value: string) => selected.includes(value), [selected]);
  return { selectedValues: selected, setSelected, isSelected };
}
