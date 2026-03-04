import { Button, Checkbox, Content, ContentVariants, SearchInput } from '@patternfly/react-core';
import * as React from 'react';
import { McpCatalogContext } from '~/app/context/mcpCatalog/McpCatalogContext';
import type {
  McpFilterCategoryKey,
  McpCatalogFilterStringOption,
} from '~/app/pages/mcpCatalog/types/mcpCatalogFilterOptions';
import { MCP_FILTER_MAX_VISIBLE } from '~/app/pages/mcpCatalog/constants/mcpCatalogFilterOptions';

function useMcpFilterState(filterKey: McpFilterCategoryKey) {
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
  return { isSelected, setSelected };
}

type McpCatalogStringFilterProps = {
  title: string;
  filterKey: McpFilterCategoryKey;
  filters: McpCatalogFilterStringOption | undefined;
  showSearch?: boolean;
};

const McpCatalogStringFilter: React.FC<McpCatalogStringFilterProps> = ({
  title,
  filterKey,
  filters,
  showSearch = false,
}) => {
  const [showMore, setShowMore] = React.useState(false);
  const [searchValue, setSearchValue] = React.useState('');
  const { isSelected, setSelected } = useMcpFilterState(filterKey);

  const filterValues = React.useMemo(() => filters?.values ?? [], [filters?.values]);

  const valuesMatchingSearch = React.useMemo(
    () =>
      filterValues.filter((value) => {
        const q = searchValue.trim().toLowerCase();
        return !q || value.toLowerCase().includes(q) || isSelected(value);
      }),
    [filterValues, searchValue, isSelected],
  );

  const visibleValues = showMore
    ? valuesMatchingSearch
    : valuesMatchingSearch.slice(0, MCP_FILTER_MAX_VISIBLE);
  const hasMoreThanMax = valuesMatchingSearch.length > MCP_FILTER_MAX_VISIBLE;

  return (
    <Content data-testid={`mcp-filter-${filterKey}`}>
      <Content component={ContentVariants.h6}>{title}</Content>
      {showSearch && (
        <SearchInput
          placeholder={`Search ${title.toLowerCase()}`}
          value={searchValue}
          onChange={(_event, value) => setSearchValue(value)}
          data-testid={`mcp-filter-${filterKey}-search`}
          className="pf-v6-u-mb-sm"
        />
      )}
      {visibleValues.length === 0 && (
        <div data-testid={`mcp-filter-${filterKey}-empty`}>No results found</div>
      )}
      {visibleValues.map((value) => (
        <Checkbox
          key={value}
          id={`${filterKey}-${value}`}
          label={value}
          isChecked={isSelected(value)}
          onChange={(_event, checked) => setSelected(value, checked)}
          data-testid={`mcp-filter-${filterKey}-${value}`}
        />
      ))}
      {!showMore && hasMoreThanMax && (
        <Button
          variant="link"
          onClick={() => setShowMore(true)}
          data-testid={`mcp-filter-${filterKey}-show-more`}
        >
          Show more
        </Button>
      )}
      {showMore && hasMoreThanMax && (
        <Button
          variant="link"
          onClick={() => setShowMore(false)}
          data-testid={`mcp-filter-${filterKey}-show-less`}
        >
          Show less
        </Button>
      )}
    </Content>
  );
};

export default McpCatalogStringFilter;
