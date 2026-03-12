import type { CatalogSourceList } from '~/app/modelCatalogTypes';
import type { McpServer } from '~/app/mcpServerCatalogTypes';
import type { ModelCatalogAPIState } from '~/app/hooks/modelCatalog/useModelCatalogAPIState';
import type {
  McpCatalogFilterOptionsList,
  McpCatalogFiltersState,
} from '~/app/pages/mcpCatalog/types/mcpCatalogFilterOptions';

export type McpCatalogPaginationState = {
  page: number;
  pageSize: number;
  totalItems: number;
};

export type McpCatalogContextType = {
  apiState: ModelCatalogAPIState;
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
  hasNoLabelSources: boolean;
  catalogSources: CatalogSourceList | null;
  catalogSourcesLoaded: boolean;
  catalogSourcesLoadError: Error | undefined;
  mcpServers: { items: McpServer[] };
  mcpServersLoaded: boolean;
  mcpServersLoadError: Error | undefined;
  refreshMcpServers: () => void;
  filterOptions: McpCatalogFilterOptionsList | null;
  filterOptionsLoaded: boolean;
  filterOptionsLoadError: Error | undefined;
};
