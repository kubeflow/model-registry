import type { CatalogFilterOptionsList } from '~/app/modelCatalogTypes';
import type { McpServer } from '~/app/mcpServerCatalogTypes';
import type { McpCatalogFiltersState } from '~/app/pages/mcpCatalog/types/mcpCatalogFilterOptions';

export type McpCatalogPaginationState = {
  page: number;
  pageSize: number;
  totalItems: number;
};

export type McpCatalogContextType = {
  filters: McpCatalogFiltersState;
  setFilters: (
    filters: McpCatalogFiltersState | ((prev: McpCatalogFiltersState) => McpCatalogFiltersState),
  ) => void;
  searchQuery: string;
  setSearchQuery: (query: string) => void;
  namedQuery: string | null;
  setNamedQuery: (query: string | null) => void;
  pagination: McpCatalogPaginationState;
  setPage: (page: number) => void;
  setPageSize: (pageSize: number) => void;
  setTotalItems: (totalItems: number) => void;
  selectedSourceLabel: string | undefined;
  setSelectedSourceLabel: (label: string | undefined) => void;
  clearAllFilters: () => void;
  sourceLabels: string[];
  sourceLabelNames: Record<string, string>;
  catalogSourcesLoaded: boolean;
  catalogSourcesLoadError: Error | undefined;
  mcpServers: { items: McpServer[] };
  mcpServersLoaded: boolean;
  mcpServersLoadError: Error | undefined;
  refreshMcpServers: () => void;
  filterOptions: CatalogFilterOptionsList | null;
  filterOptionsLoaded: boolean;
  filterOptionsLoadError: Error | undefined;
};
