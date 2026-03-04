import type { McpCatalogFiltersState } from '~/app/context/mcpCatalog/McpCatalogContext';
import { MCP_FILTER_KEYS } from '~/app/pages/mcpCatalog/constants/mcpCatalogFilterChipNames';

export const hasMcpFiltersApplied = (
  filters: McpCatalogFiltersState,
  searchQuery: string,
): boolean => {
  if (searchQuery && searchQuery.trim().length > 0) {
    return true;
  }
  for (const key of MCP_FILTER_KEYS) {
    const value = filters[key];
    if (Array.isArray(value) && value.length > 0) {
      return true;
    }
  }
  return false;
};
