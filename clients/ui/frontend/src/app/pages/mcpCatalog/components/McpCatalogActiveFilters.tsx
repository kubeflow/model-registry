import * as React from 'react';
import { ToolbarFilter, ToolbarLabel, ToolbarLabelGroup } from '@patternfly/react-core';
import { McpCatalogContext } from '~/app/context/mcpCatalog/McpCatalogContext';
import type { McpFilterCategoryKey } from '~/app/pages/mcpCatalog/types/mcpCatalogFilterOptions';
import {
  MCP_FILTER_KEYS,
  MCP_FILTER_CATEGORY_NAMES,
} from '~/app/pages/mcpCatalog/constants/mcpCatalogFilterChipNames';

const McpCatalogActiveFilters: React.FC = () => {
  const { filters, setFilters } = React.useContext(McpCatalogContext);

  const handleRemoveFilter = React.useCallback(
    (categoryKey: McpFilterCategoryKey, valueKey: string) => {
      setFilters((prev) => {
        const current = prev[categoryKey];
        const arr = Array.isArray(current) ? current : [];
        const newValues = arr.filter((v) => v !== valueKey);
        return { ...prev, [categoryKey]: newValues };
      });
    },
    [setFilters],
  );

  const handleClearCategory = React.useCallback(
    (categoryKey: McpFilterCategoryKey) => {
      setFilters((prev) => ({ ...prev, [categoryKey]: [] }));
    },
    [setFilters],
  );

  return (
    <>
      {MCP_FILTER_KEYS.map((filterKey) => {
        const filterValue = filters[filterKey];
        const values = Array.isArray(filterValue) ? filterValue : [];
        const hasValue = values.length > 0;

        const labels: ToolbarLabel[] = hasValue
          ? values.map((value) => ({
              key: value,
              node: <span data-testid={`mcp-filter-chip-${filterKey}-${value}`}>{value}</span>,
            }))
          : [];

        const categoryLabelGroup: ToolbarLabelGroup = {
          key: filterKey,
          name: MCP_FILTER_CATEGORY_NAMES[filterKey],
        };

        return (
          <ToolbarFilter
            key={filterKey}
            categoryName={categoryLabelGroup}
            labels={labels}
            deleteLabel={(_, label) => {
              const labelKey = typeof label === 'string' ? label : label.key;
              handleRemoveFilter(filterKey, labelKey);
            }}
            deleteLabelGroup={() => handleClearCategory(filterKey)}
            data-testid={`mcp-filter-container-${filterKey}`}
          >
            {null}
          </ToolbarFilter>
        );
      })}
    </>
  );
};

export default McpCatalogActiveFilters;
