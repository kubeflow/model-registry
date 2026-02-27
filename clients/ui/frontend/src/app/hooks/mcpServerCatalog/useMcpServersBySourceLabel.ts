import React from 'react';
import { FetchStateCallbackPromise, useFetchState } from 'mod-arch-core';
import { McpServer, McpServerList, McpServerListParams } from '~/app/mcpServerCatalogTypes';
import { useModelCatalogAPI } from '~/app/hooks/modelCatalog/useModelCatalogAPI';

type PaginatedMcpServerList = {
  items: McpServer[];
  size: number;
  pageSize: number;
  nextPageToken: string;
  loadMore: () => void;
  isLoadingMore: boolean;
  hasMore: boolean;
  refresh: () => void;
};

type McpServers = {
  mcpServers: PaginatedMcpServerList;
  mcpServersLoaded: boolean;
  mcpServersLoadError: Error | undefined;
  refresh: () => void;
};

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
): McpServers => {
  const { api, apiAvailable } = useModelCatalogAPI();

  const [allItems, setAllItems] = React.useState<McpServer[]>([]);
  const [totalSize, setTotalSize] = React.useState(0);
  const [nextPageTokenValue, setNextPageTokenValue] = React.useState('');
  const [isLoadingMore, setIsLoadingMore] = React.useState(false);

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
    if (loaded && !error) {
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

    try {
      const response = await api.getMcpServerList({}, buildMcpServerListParams(nextPageTokenValue));

      setAllItems((prev) => [...prev, ...(response.items ?? [])]);
      setTotalSize(response.size);
      setNextPageTokenValue(response.nextPageToken);
    } catch (err) {
      throw new Error(
        `Failed to load more servers: ${err instanceof Error ? err.message : String(err)}`,
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
  };

  return {
    mcpServers: paginatedData,
    mcpServersLoaded: loaded,
    mcpServersLoadError: error,
    refresh,
  };
};
