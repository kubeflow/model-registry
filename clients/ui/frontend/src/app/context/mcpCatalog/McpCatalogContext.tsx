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
} from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import type { CatalogFilterOptionsList } from '~/app/modelCatalogTypes';
import { mockMcpServers } from '~/app/pages/mcpCatalog/mocks/mockMcpServers';

export type McpCatalogPaginationState = {
  page: number;
  pageSize: number;
  totalItems: number;
};

export type McpCatalogFiltersState = Record<string, unknown>;

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
  catalogSourcesLoaded: boolean;
  catalogSourcesLoadError: Error | undefined;
  mcpServers: { items: import('~/app/mcpServerCatalogTypes').McpServer[] };
  mcpServersLoaded: boolean;
  mcpServersLoadError: Error | undefined;
  filterOptions: CatalogFilterOptionsList | null;
  filterOptionsLoaded: boolean;
  filterOptionsLoadError: Error | undefined;
};

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
  catalogSourcesLoaded: false,
  catalogSourcesLoadError: undefined,
  mcpServers: { items: [] },
  mcpServersLoaded: false,
  mcpServersLoadError: undefined,
  filterOptions: null,
  filterOptionsLoaded: false,
  filterOptionsLoadError: undefined,
});

export const McpCatalogContextProvider: React.FC<McpCatalogContextProviderProps> = ({
  children,
}) => {
  const hostPath = `${URL_PREFIX}/api/${BFF_API_VERSION}/model_catalog`;
  const queryParams = useQueryParamNamespaces();
  const [apiState] = useModelCatalogAPIState(hostPath, queryParams);

  const mcpListParams = React.useMemo(() => ({ assetType: 'mcp_servers' as const }), []);
  const [catalogSources, catalogSourcesLoaded, catalogSourcesLoadError] = useCatalogSources(
    apiState,
    mcpListParams,
  );
  const [filterOptions, filterOptionsLoaded, filterOptionsLoadError] =
    useMcpServerFilterOptionListWithAPI(apiState);

  const [filters, setFilters] = React.useState<McpCatalogFiltersState>({});
  const [searchQuery, setSearchQuery] = React.useState('');
  const [namedQuery, setNamedQuery] = React.useState<string | null>(null);
  const [pagination, setPaginationState] =
    React.useState<McpCatalogPaginationState>(defaultPagination);
  const [selectedSourceLabel, setSelectedSourceLabel] = React.useState<string | undefined>(
    undefined,
  );

  const sourceLabelsFromApi = React.useMemo(() => {
    const enabled = filterEnabledCatalogSources(catalogSources);
    return getUniqueSourceLabels(enabled);
  }, [catalogSources]);

  const useMockLabels = !apiState.apiAvailable || !catalogSourcesLoaded;
  const sourceLabels = useMockLabels
    ? Array.from(
        new Set(mockMcpServers.map((s) => s.source_id).filter((id): id is string => Boolean(id))),
      )
    : sourceLabelsFromApi;

  const mcpServersResult = useMcpServersBySourceLabelWithAPI(apiState, {
    sourceLabel: selectedSourceLabel,
    pageSize: pagination.pageSize,
    searchQuery,
  });

  const useMockData = !apiState.apiAvailable;
  const filteredMockItems = mockMcpServers.filter(
    (s) =>
      selectedSourceLabel === undefined || (s.source_id && s.source_id === selectedSourceLabel),
  );
  const apiHasItems = mcpServersResult.mcpServers.items.length > 0;
  const mcpServers = useMockData
    ? { items: filteredMockItems }
    : apiHasItems
      ? mcpServersResult.mcpServers
      : { items: filteredMockItems };
  const mcpServersLoaded = useMockData ? true : mcpServersResult.mcpServersLoaded;
  const mcpServersLoadError = useMockData ? undefined : mcpServersResult.mcpServersLoadError;

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
      catalogSourcesLoaded,
      catalogSourcesLoadError,
      mcpServers: { items: mcpServers.items },
      mcpServersLoaded,
      mcpServersLoadError,
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
      catalogSourcesLoaded,
      catalogSourcesLoadError,
      mcpServers.items,
      mcpServersLoaded,
      mcpServersLoadError,
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
