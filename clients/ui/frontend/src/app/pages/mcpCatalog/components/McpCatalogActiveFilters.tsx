import * as React from 'react';
import { ToolbarFilter, ToolbarLabel, ToolbarLabelGroup } from '@patternfly/react-core';
import { McpCatalogContext } from '~/app/context/mcpCatalog/McpCatalogContext';
import type { McpFilterCategoryKey } from '~/app/pages/mcpCatalog/types/mcpCatalogFilterOptions';
import {
  MCP_FILTER_KEYS,
  MCP_FILTER_CATEGORY_NAMES,
} from '~/app/pages/mcpCatalog/constants/mcpCatalogFilterChipNames';

const SEARCH_CHIP_CATEGORY = 'Search';

const McpCatalogActiveFilters: React.FC = () => {
  const { filters, setFilters, searchQuery, setSearchQuery } = React.useContext(McpCatalogContext);

  const hasSearchChip = searchQuery.trim().length > 0;

  const handleClearSearch = React.useCallback(() => {
    setSearchQuery('');
  }, [setSearchQuery]);

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
      {hasSearchChip && (
        <ToolbarFilter
          key="search"
          categoryName={{ key: 'search', name: SEARCH_CHIP_CATEGORY }}
          labels={[
            {
              key: searchQuery.trim(),
              node: <span data-testid="mcp-filter-chip-search">{searchQuery.trim()}</span>,
            },
          ]}
          deleteLabel={handleClearSearch}
          deleteLabelGroup={handleClearSearch}
          data-testid="mcp-filter-container-search"
        >
          {null}
        </ToolbarFilter>
      )}
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
