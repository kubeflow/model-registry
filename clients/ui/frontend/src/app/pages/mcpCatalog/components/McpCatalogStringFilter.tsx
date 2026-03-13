import * as React from 'react';
import CatalogStringFilter from '~/app/shared/components/catalog/CatalogStringFilter';
import { useMcpFilterState } from '~/app/pages/mcpCatalog/hooks/useMcpFilterState';
import type {
  McpFilterCategoryKey,
  McpCatalogFilterStringOption,
} from '~/app/pages/mcpCatalog/types/mcpCatalogFilterOptions';

type McpCatalogStringFilterProps = {
  title: string;
  filterKey: McpFilterCategoryKey;
  filters: McpCatalogFilterStringOption | undefined;
};

const McpCatalogStringFilter: React.FC<McpCatalogStringFilterProps> = ({
  title,
  filterKey,
  filters,
}) => {
  const { selectedValues, setSelected } = useMcpFilterState(filterKey);
  const filterValues = React.useMemo(() => filters?.values ?? [], [filters?.values]);

  return (
    <CatalogStringFilter
      title={title}
      filterValues={filterValues}
      selectedValues={selectedValues}
      onToggle={setSelected}
      testIdBase={`mcp-filter-${filterKey}`}
    />
  );
};

export default McpCatalogStringFilter;
