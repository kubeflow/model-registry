import * as React from 'react';

export type McpCatalogPaginationState = {
  page: number;
  pageSize: number;
  totalItems: number;
};

export type McpCatalogFiltersState = Record<string, unknown>;

export type McpCatalogCategoryId = 'all' | 'sample' | 'other';

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
  selectedCategory: McpCatalogCategoryId;
  setSelectedCategory: (category: McpCatalogCategoryId) => void;
  clearAllFilters: () => void;
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
  selectedCategory: 'all',
  setSelectedCategory: () => undefined,
  clearAllFilters: () => undefined,
});

export const McpCatalogContextProvider: React.FC<McpCatalogContextProviderProps> = ({
  children,
}) => {
  const [filters, setFilters] = React.useState<McpCatalogFiltersState>({});
  const [searchQuery, setSearchQuery] = React.useState('');
  const [namedQuery, setNamedQuery] = React.useState<string | null>(null);
  const [pagination, setPaginationState] =
    React.useState<McpCatalogPaginationState>(defaultPagination);
  const [selectedCategory, setSelectedCategory] =
    React.useState<McpCatalogContextType['selectedCategory']>('all');

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
      selectedCategory,
      setSelectedCategory,
      clearAllFilters,
    }),
    [
      filters,
      searchQuery,
      namedQuery,
      pagination,
      selectedCategory,
      setPage,
      setPageSize,
      setTotalItems,
      clearAllFilters,
    ],
  );

  return <McpCatalogContext.Provider value={value}>{children}</McpCatalogContext.Provider>;
};
