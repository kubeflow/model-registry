import * as React from 'react';
import { useQueryParamNamespaces } from 'mod-arch-core';
import { BFF_API_VERSION, URL_PREFIX } from '~/app/utilities/const';
import useModelCatalogAPIState from '~/app/hooks/modelCatalog/useModelCatalogAPIState';
import { useCatalogSources } from '~/app/hooks/modelCatalog/useCatalogSources';
import { useMcpServersBySourceLabelWithAPI } from '~/app/hooks/mcpServerCatalog/useMcpServersBySourceLabel';
import { useMcpServerFilterOptionListWithAPI } from '~/app/hooks/mcpServerCatalog/useMcpServerFilterOptionList';
import {
  filterEnabledCatalogSources,
  getUniqueSourceLabels,
  hasSourcesWithoutLabels,
} from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import type {
  McpCatalogContextType,
  McpCatalogPaginationState,
} from '~/app/pages/mcpCatalog/types/mcpCatalogContext';
import type { McpCatalogFiltersState } from '~/app/pages/mcpCatalog/types/mcpCatalogFilterOptions';
import {
  filterMcpServersByFilters,
  filterMcpServersBySearchQuery,
} from '~/app/pages/mcpCatalog/utils/mcpCatalogUtils';
import { useMcpUrlSync } from '~/app/pages/mcpCatalog/hooks/useMcpUrlSync';

export type {
  McpCatalogContextType,
  McpCatalogPaginationState,
} from '~/app/pages/mcpCatalog/types/mcpCatalogContext';
export type { McpCatalogFiltersState } from '~/app/pages/mcpCatalog/types/mcpCatalogFilterOptions';

type McpCatalogContextProviderProps = {
  children: React.ReactNode;
};

const defaultPagination: McpCatalogPaginationState = {
  page: 1,
  pageSize: 10,
  totalItems: 0,
};

export const McpCatalogContext = React.createContext<McpCatalogContextType>({
  filters: {},
  setFilters: () => undefined,
  searchQuery: '',
  setSearchQuery: () => undefined,
  namedQuery: null,
  setNamedQuery: () => undefined,
  pagination: defaultPagination,
  setPage: () => undefined,
  setPageSize: () => undefined,
  setTotalItems: () => undefined,
  selectedSourceLabel: undefined,
  setSelectedSourceLabel: () => undefined,
  clearAllFilters: () => undefined,
  sourceLabels: [],
  sourceLabelNames: {},
  hasNoLabelSources: false,
  catalogSources: null,
  catalogSourcesLoaded: false,
  catalogSourcesLoadError: undefined,
  mcpServers: { items: [] },
  mcpServersLoaded: false,
  mcpServersLoadError: undefined,
  refreshMcpServers: () => undefined,
  filterOptions: null,
  filterOptionsLoaded: false,
  filterOptionsLoadError: undefined,
});

const MODEL_CATALOG_PATH = `${URL_PREFIX}/api/${BFF_API_VERSION}/model_catalog`;
const MCP_CATALOG_PATH = `${URL_PREFIX}/api/${BFF_API_VERSION}/mcp_catalog`;

export const McpCatalogContextProvider: React.FC<McpCatalogContextProviderProps> = ({
  children,
}) => {
  const queryParams = useQueryParamNamespaces();
  const [apiStateModelCatalog] = useModelCatalogAPIState(MODEL_CATALOG_PATH, queryParams);
  const [apiStateMcpCatalog] = useModelCatalogAPIState(MCP_CATALOG_PATH, queryParams);

  const mcpListParams = React.useMemo(() => ({ assetType: 'mcp_servers' as const }), []);
  const [catalogSources, catalogSourcesLoaded, catalogSourcesLoadError] = useCatalogSources(
    apiStateModelCatalog,
    mcpListParams,
  );
  const [filterOptions, filterOptionsLoaded, filterOptionsLoadError] =
    useMcpServerFilterOptionListWithAPI(apiStateMcpCatalog);

  const { initialState, syncToUrl } = useMcpUrlSync();

  const [filters, setFilters] = React.useState<McpCatalogFiltersState>(initialState.filters);
  const [searchQuery, setSearchQuery] = React.useState(initialState.searchQuery);
  const [namedQuery, setNamedQuery] = React.useState<string | null>(null);
  const [pagination, setPaginationState] =
    React.useState<McpCatalogPaginationState>(defaultPagination);
  const [selectedSourceLabel, setSelectedSourceLabel] = React.useState<string | undefined>(
    initialState.selectedSourceLabel,
  );

  React.useEffect(() => {
    syncToUrl({ searchQuery, filters, selectedSourceLabel });
  }, [searchQuery, filters, selectedSourceLabel, syncToUrl]);

  const { sourceLabels, sourceLabelNames, hasNoLabelSources } = React.useMemo(() => {
    const enabled = filterEnabledCatalogSources(catalogSources);
    const labels = getUniqueSourceLabels(enabled);
    const nameMap: Record<string, string> = {};
    if (enabled?.items) {
      for (const source of enabled.items) {
        for (const label of source.labels) {
          const trimmed = label.trim();
          if (trimmed && !nameMap[trimmed]) {
            nameMap[trimmed] = source.name;
          }
        }
      }
    }
    return {
      sourceLabels: labels,
      sourceLabelNames: nameMap,
      hasNoLabelSources: hasSourcesWithoutLabels(enabled),
    };
  }, [catalogSources]);

  const mcpServersResult = useMcpServersBySourceLabelWithAPI(apiStateMcpCatalog, {
    sourceLabel: selectedSourceLabel,
    pageSize: pagination.pageSize,
    searchQuery,
  });

  const mcpServers = React.useMemo(() => {
    const { items } = mcpServersResult.mcpServers;
    let filtered = items;
    if (searchQuery.trim().length > 0) {
      filtered = filterMcpServersBySearchQuery(filtered, searchQuery);
    }
    return { items: filterMcpServersByFilters(filtered, filters) };
    // mcpServersResult.mcpServers is new each render; .items gives stable deps to avoid cascading re-renders
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedSourceLabel, searchQuery, filters, mcpServersResult.mcpServers.items]);

  const { mcpServersLoaded, mcpServersLoadError, refresh: refreshMcpServers } = mcpServersResult;

  const setPage = React.useCallback((page: number) => {
    setPaginationState((prev) => ({ ...prev, page }));
  }, []);

  const setPageSize = React.useCallback((pageSize: number) => {
    setPaginationState((prev) => ({ ...prev, pageSize, page: 1 }));
  }, []);

  const setTotalItems = React.useCallback((totalItems: number) => {
    setPaginationState((prev) => ({ ...prev, totalItems }));
  }, []);

  const clearAllFilters = React.useCallback(() => {
    setSearchQuery('');
    setFilters({});
    setSelectedSourceLabel(undefined);
    setNamedQuery(null);
  }, []);

  const value = React.useMemo<McpCatalogContextType>(
    () => ({
      filters,
      setFilters,
      searchQuery,
      setSearchQuery,
      namedQuery,
      setNamedQuery,
      pagination,
      setPage,
      setPageSize,
      setTotalItems,
      selectedSourceLabel,
      setSelectedSourceLabel,
      clearAllFilters,
      sourceLabels,
      sourceLabelNames,
      hasNoLabelSources,
      catalogSources,
      catalogSourcesLoaded,
      catalogSourcesLoadError,
      mcpServers,
      mcpServersLoaded,
      mcpServersLoadError,
      refreshMcpServers,
      filterOptions,
      filterOptionsLoaded,
      filterOptionsLoadError,
    }),
    [
      filters,
      searchQuery,
      namedQuery,
      pagination,
      selectedSourceLabel,
      sourceLabels,
      sourceLabelNames,
      hasNoLabelSources,
      catalogSources,
      catalogSourcesLoaded,
      catalogSourcesLoadError,
      mcpServers,
      mcpServersLoaded,
      mcpServersLoadError,
      refreshMcpServers,
      filterOptions,
      filterOptionsLoaded,
      filterOptionsLoadError,
      setPage,
      setPageSize,
      setTotalItems,
      clearAllFilters,
    ],
  );

  return <McpCatalogContext.Provider value={value}>{children}</McpCatalogContext.Provider>;
};
