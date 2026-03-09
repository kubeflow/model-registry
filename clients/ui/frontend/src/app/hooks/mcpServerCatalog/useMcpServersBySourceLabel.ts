import React from 'react';
import { FetchStateCallbackPromise, useFetchState } from 'mod-arch-core';
import { McpServer, McpServerList, McpServerListParams } from '~/app/mcpServerCatalogTypes';
import { useModelCatalogAPI } from '~/app/hooks/modelCatalog/useModelCatalogAPI';
import type { ModelCatalogAPIState } from '~/app/hooks/modelCatalog/useModelCatalogAPIState';

type PaginatedMcpServerList = {
  items: McpServer[];
  size: number;
  pageSize: number;
  nextPageToken: string;
  loadMore: () => void;
  isLoadingMore: boolean;
  hasMore: boolean;
  refresh: () => void;
  loadMoreError?: Error;
};

export type McpServersResult = {
  mcpServers: PaginatedMcpServerList;
  mcpServersLoaded: boolean;
  mcpServersLoadError: Error | undefined;
  refresh: () => void;
};

type UseMcpServersBySourceLabelParams = {
  sourceLabel?: string;
  pageSize?: number;
  searchQuery?: string;
  filterQuery?: string;
  namedQuery?: string;
  includeTools?: boolean;
  toolLimit?: number;
  sortBy?: string | null;
  sortOrder?: string;
};

export function useMcpServersBySourceLabelWithAPI(
  apiState: ModelCatalogAPIState,
  params: UseMcpServersBySourceLabelParams,
): McpServersResult {
  const {
    sourceLabel,
    pageSize = 10,
    searchQuery = '',
    filterQuery,
    namedQuery,
    includeTools,
    toolLimit,
    sortBy,
    sortOrder,
  } = params;
  const { api, apiAvailable } = apiState;

  const [allItems, setAllItems] = React.useState<McpServer[]>([]);
  const [totalSize, setTotalSize] = React.useState(0);
  const [nextPageTokenValue, setNextPageTokenValue] = React.useState('');
  const [isLoadingMore, setIsLoadingMore] = React.useState(false);
  const [loadMoreError, setLoadMoreError] = React.useState<Error | undefined>();

  const buildMcpServerListParams = React.useCallback(
    (nextPageToken?: string): McpServerListParams => ({
      sourceLabel,
      pageSize: pageSize.toString(),
      ...(nextPageToken !== undefined && nextPageToken !== '' && { nextPageToken }),
      q: searchQuery.trim() || undefined,
      filterQuery,
      namedQuery,
      includeTools,
      toolLimit,
      orderBy: sortBy ?? undefined,
      sortOrder,
    }),
    [
      sourceLabel,
      pageSize,
      searchQuery,
      filterQuery,
      namedQuery,
      includeTools,
      toolLimit,
      sortBy,
      sortOrder,
    ],
  );

  const fetchMcpServers = React.useCallback<FetchStateCallbackPromise<McpServerList>>(
    (opts) => {
      if (!apiAvailable) {
        return Promise.reject(new Error('API not yet available'));
      }

      return api.getMcpServerList(opts, buildMcpServerListParams());
    },
    [api, apiAvailable, buildMcpServerListParams],
  );

  const [firstPageData, loaded, error, refetch] = useFetchState(
    fetchMcpServers,
    { items: [], size: 0, pageSize: 10, nextPageToken: '' },
    { initialPromisePurity: true },
  );

  React.useEffect(() => {
    if (loaded && !error && (firstPageData.items?.length ?? 0) > 0) {
      setAllItems(firstPageData.items ?? []);
      setTotalSize(firstPageData.size);
      setNextPageTokenValue(firstPageData.nextPageToken);
    }
  }, [firstPageData, loaded, error]);

  const loadMore = React.useCallback(async () => {
    if (!nextPageTokenValue || isLoadingMore || !apiAvailable) {
      return;
    }

    setIsLoadingMore(true);
    setLoadMoreError(undefined);

    try {
      const response = await api.getMcpServerList({}, buildMcpServerListParams(nextPageTokenValue));

      setAllItems((prev) => [...prev, ...(response.items ?? [])]);
      setTotalSize(response.size);
      setNextPageTokenValue(response.nextPageToken);
      setLoadMoreError(undefined);
    } catch (err) {
      setLoadMoreError(
        new Error(
          `Failed to load more servers: ${err instanceof Error ? err.message : String(err)}`,
        ),
      );
    } finally {
      setIsLoadingMore(false);
    }
  }, [api, apiAvailable, buildMcpServerListParams, isLoadingMore, nextPageTokenValue]);

  React.useEffect(() => {
    setAllItems([]);
    setTotalSize(0);
    setNextPageTokenValue('');
    setIsLoadingMore(false);
    setLoadMoreError(undefined);
  }, [
    sourceLabel,
    pageSize,
    searchQuery,
    filterQuery,
    namedQuery,
    includeTools,
    toolLimit,
    sortBy,
    sortOrder,
  ]);

  const refresh = React.useCallback(() => {
    setAllItems([]);
    setTotalSize(0);
    setNextPageTokenValue('');
    setIsLoadingMore(false);
    setLoadMoreError(undefined);
    refetch();
  }, [refetch]);

  const paginatedData: PaginatedMcpServerList = {
    items: allItems,
    size: totalSize,
    pageSize: firstPageData.pageSize,
    nextPageToken: nextPageTokenValue,
    loadMore,
    isLoadingMore,
    hasMore: Boolean(nextPageTokenValue),
    refresh,
    loadMoreError,
  };

  return {
    mcpServers: paginatedData,
    mcpServersLoaded: loaded,
    mcpServersLoadError: error,
    refresh,
  };
}

export const useMcpServersBySourceLabel = (
  sourceLabel?: string,
  pageSize = 10,
  searchQuery = '',
  filterQuery?: string,
  namedQuery?: string,
  includeTools?: boolean,
  toolLimit?: number,
  sortBy?: string | null,
  sortOrder?: string,
): McpServersResult => {
  const { api, apiAvailable } = useModelCatalogAPI();
  return useMcpServersBySourceLabelWithAPI(
    { api, apiAvailable },
    {
      sourceLabel,
      pageSize,
      searchQuery,
      filterQuery,
      namedQuery,
      includeTools,
      toolLimit,
      sortBy,
      sortOrder,
    },
  );
};
