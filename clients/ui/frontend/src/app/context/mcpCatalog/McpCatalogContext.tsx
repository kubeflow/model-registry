import { useQueryParamNamespaces } from 'mod-arch-core';
import * as React from 'react';
import useMcpCatalogAPIState, {
  McpCatalogAPIState,
} from '~/app/hooks/mcpCatalog/useMcpCatalogAPIState';
import { useMcpFilterOptions } from '~/app/hooks/mcpCatalog/useMcpFilterOptions';
import { useMcpServers } from '~/app/hooks/mcpCatalog/useMcpServers';
import { useMcpSources } from '~/app/hooks/mcpCatalog/useMcpSources';
import {
  McpCategoryName,
  McpCatalogSourceList,
  McpFilterOptionsList,
  McpServerList,
} from '~/app/pages/mcpCatalog/types';
import {
  McpServerFilterState,
  mcpFiltersToFilterQuery,
} from '~/app/pages/mcpCatalog/utils/mcpCatalogUtils';
import { BFF_API_VERSION, URL_PREFIX } from '~/app/utilities/const';

/**
 * Filter options derived from backend filter_options endpoint.
 * These remain stable regardless of applied filters.
 */
export type McpFilterOptions = {
  providers: string[];
  licenses: string[];
  tags: string[];
  transports: string[];
  deploymentModes: string[];
};

const EMPTY_FILTER_OPTIONS: McpFilterOptions = {
  providers: [],
  licenses: [],
  tags: [],
  transports: [],
  deploymentModes: [],
};

/**
 * Convert backend filter options response to the frontend McpFilterOptions type.
 * The backend returns a map of field names to filter option objects with values arrays.
 */
const convertBackendFilterOptions = (backendOptions: McpFilterOptionsList): McpFilterOptions => {
  const filters = backendOptions.filters ?? {};

  const getStringValues = (key: string): string[] => {
    if (!Object.prototype.hasOwnProperty.call(filters, key)) {
      return [];
    }
    const option = filters[key];
    if (!option.values) {
      return [];
    }
    return option.values.filter((v): v is string => typeof v === 'string').toSorted();
  };

  return {
    providers: getStringValues('provider'),
    licenses: getStringValues('license'),
    tags: getStringValues('tags'),
    transports: getStringValues('transports'),
    deploymentModes: getStringValues('deploymentMode'),
  };
};

// Initial empty filter state
const EMPTY_FILTER_STATE: McpServerFilterState = {
  selectedProviders: [],
  selectedLicenses: [],
  selectedTags: [],
  selectedTransports: [],
  selectedDeploymentModes: [],
};

export type McpCatalogContextType = {
  mcpServersLoaded: boolean;
  mcpServersLoadError?: Error;
  mcpServers: McpServerList | null;
  mcpSources: McpCatalogSourceList | null;
  mcpSourcesLoaded: boolean;
  mcpSourcesLoadError?: Error;
  selectedSourceLabel: string;
  updateSelectedSourceLabel: (label: string) => void;
  // Filter state management
  filters: McpServerFilterState;
  searchTerm: string;
  updateFilters: (filters: McpServerFilterState) => void;
  updateSearchTerm: (term: string) => void;
  resetFilters: () => void;
  // Filter options from backend filter_options endpoint
  filterOptions: McpFilterOptions;
  filterOptionsLoaded: boolean;
  filterOptionsLoadError?: Error;
  apiState: McpCatalogAPIState;
  refreshAPIState: () => void;
  refreshMcpServers: () => void;
  refreshMcpSources: () => void;
};

type McpCatalogContextProviderProps = {
  children: React.ReactNode;
};

export const McpCatalogContext = React.createContext<McpCatalogContextType>({
  mcpServersLoaded: false,
  mcpServersLoadError: undefined,
  mcpServers: null,
  mcpSources: null,
  mcpSourcesLoaded: false,
  mcpSourcesLoadError: undefined,
  selectedSourceLabel: McpCategoryName.allServers,
  updateSelectedSourceLabel: () => undefined,
  filters: EMPTY_FILTER_STATE,
  searchTerm: '',
  updateFilters: () => undefined,
  updateSearchTerm: () => undefined,
  resetFilters: () => undefined,
  filterOptions: EMPTY_FILTER_OPTIONS,
  filterOptionsLoaded: false,
  filterOptionsLoadError: undefined,
  // eslint-disable-next-line @typescript-eslint/consistent-type-assertions
  apiState: { apiAvailable: false, api: null as unknown as McpCatalogAPIState['api'] },
  refreshAPIState: () => undefined,
  refreshMcpServers: () => undefined,
  refreshMcpSources: () => undefined,
});

export const McpCatalogContextProvider: React.FC<McpCatalogContextProviderProps> = ({
  children,
}) => {
  const hostPath = `${URL_PREFIX}/api/${BFF_API_VERSION}/mcp_catalog`;
  const queryParams = useQueryParamNamespaces();
  const [apiState, refreshAPIState] = useMcpCatalogAPIState(hostPath, queryParams);
  const [selectedSourceLabel, setSelectedSourceLabel] = React.useState<string>(
    McpCategoryName.allServers,
  );

  // Filter and search state
  const [filters, setFilters] = React.useState<McpServerFilterState>(EMPTY_FILTER_STATE);
  const [searchTerm, setSearchTerm] = React.useState<string>('');

  // Build filterQuery from current filter state
  const filterQuery = React.useMemo(() => mcpFiltersToFilterQuery(filters), [filters]);

  // Fetch filter options from the dedicated backend endpoint
  const [backendFilterOptions, filterOptionsLoaded, filterOptionsLoadError] =
    useMcpFilterOptions(apiState);

  // Pass filterQuery and searchTerm to useMcpServers for the filtered results
  const [mcpServers, mcpServersLoaded, mcpServersLoadError, refreshMcpServers] = useMcpServers(
    apiState,
    { filterQuery: filterQuery || undefined, searchTerm: searchTerm || undefined },
  );
  const [mcpSources, mcpSourcesLoaded, mcpSourcesLoadError, refreshMcpSources] =
    useMcpSources(apiState);

  // Convert backend filter options to frontend format
  const filterOptions = React.useMemo(
    () => convertBackendFilterOptions(backendFilterOptions),
    [backendFilterOptions],
  );

  const updateSelectedSourceLabel = React.useCallback((label: string) => {
    setSelectedSourceLabel(label);
  }, []);

  const updateFilters = React.useCallback((newFilters: McpServerFilterState) => {
    setFilters(newFilters);
  }, []);

  const updateSearchTerm = React.useCallback((term: string) => {
    setSearchTerm(term);
  }, []);

  const resetFilters = React.useCallback(() => {
    setFilters(EMPTY_FILTER_STATE);
    setSearchTerm('');
  }, []);

  const contextValue = React.useMemo(
    () => ({
      mcpServersLoaded,
      mcpServersLoadError,
      mcpServers,
      mcpSources,
      mcpSourcesLoaded,
      mcpSourcesLoadError,
      selectedSourceLabel,
      updateSelectedSourceLabel,
      filters,
      searchTerm,
      updateFilters,
      updateSearchTerm,
      resetFilters,
      filterOptions,
      filterOptionsLoaded,
      filterOptionsLoadError,
      apiState,
      refreshAPIState,
      refreshMcpServers,
      refreshMcpSources,
    }),
    [
      mcpServersLoaded,
      mcpServersLoadError,
      mcpServers,
      mcpSources,
      mcpSourcesLoaded,
      mcpSourcesLoadError,
      selectedSourceLabel,
      updateSelectedSourceLabel,
      filters,
      searchTerm,
      updateFilters,
      updateSearchTerm,
      resetFilters,
      filterOptions,
      filterOptionsLoaded,
      filterOptionsLoadError,
      apiState,
      refreshAPIState,
      refreshMcpServers,
      refreshMcpSources,
    ],
  );

  return <McpCatalogContext.Provider value={contextValue}>{children}</McpCatalogContext.Provider>;
};

export const useMcpCatalog = (): McpCatalogContextType => React.useContext(McpCatalogContext);
