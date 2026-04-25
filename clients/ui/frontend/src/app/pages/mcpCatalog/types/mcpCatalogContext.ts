import type { CatalogLabelList, CatalogSourceList } from '~/app/modelCatalogTypes';
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
  mcpApiState: ModelCatalogAPIState;
  catalogSources: CatalogSourceList | null;
  catalogSourcesLoaded: boolean;
  catalogSourcesLoadError: Error | undefined;
  catalogLabels: CatalogLabelList | null;
  catalogLabelsLoaded: boolean;
  catalogLabelsLoadError: Error | undefined;
  filterOptions: McpCatalogFilterOptionsList | null;
  filterOptionsLoaded: boolean;
  filterOptionsLoadError: Error | undefined;
  emptyCategoryLabels: Set<string>;
  reportCategoryEmpty: (label: string, isEmpty: boolean) => void;
};
