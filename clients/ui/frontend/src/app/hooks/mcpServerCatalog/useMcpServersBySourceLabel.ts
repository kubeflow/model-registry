import { McpServer, McpServerList } from '~/app/mcpServerCatalogTypes';
import { useModelCatalogAPI } from '../modelCatalog/useModelCatalogAPI';
import React from 'react';
import { FetchStateCallbackPromise, useFetchState } from 'mod-arch-core';

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

export const useMcpServersBySourceLabel = (sourceLabel?: string): McpServers => {
  const { api, apiAvailable } = useModelCatalogAPI();

  const [allItems, setAllItems] = React.useState<McpServer[]>([]);
  const [totalSize, setTotalSize] = React.useState(0);
  const [nextPageToken, setNextPageToken] = React.useState('');
  const [isLoadingMore, setIsLoadingMore] = React.useState(false);

  const fetchMcpServers = React.useCallback<FetchStateCallbackPromise<McpServerList>>(
    (opts) => {
      if (!apiAvailable) {
        return Promise.reject(new Error('API not yet available'));
      }

      return api.getMcpServerList(opts, sourceLabel);
    },
    [api, apiAvailable, sourceLabel],
  );

  const [firstPageData, loaded, error, refetch] = useFetchState(
    fetchMcpServers,
    { items: [], size: 0, pageSize: 10, nextPageToken: '' },
    { initialPromisePurity: true },
  );

  React.useEffect(() => {
    if (loaded && !error && firstPageData.items.length > 0) {
      setAllItems(firstPageData.items);
      setTotalSize(firstPageData.size);
      setNextPageToken(firstPageData.nextPageToken);
    }
  }, [firstPageData, loaded, error]);

  const loadMore = React.useCallback(async () => {
    if (!nextPageToken || isLoadingMore || !apiAvailable) {
      return;
    }

    setIsLoadingMore(true);

    try {
      const response = await api.getMcpServerList({}, sourceLabel);

      setAllItems((prev) => [...prev, ...response.items]);
      setTotalSize(response.size);
      setNextPageToken(response.nextPageToken);
    } catch (err) {
      throw new Error(
        `Failed to load more servers: ${err instanceof Error ? err.message : String(err)}`,
      );
    } finally {
      setIsLoadingMore(false);
    }
  }, [api, apiAvailable, nextPageToken, isLoadingMore, sourceLabel]);

  React.useEffect(() => {
    setAllItems([]);
    setTotalSize(0);
    setNextPageToken('');
    setIsLoadingMore(false);
  }, [sourceLabel]);

  const refresh = React.useCallback(() => {
    setAllItems([]);
    setTotalSize(0);
    setNextPageToken('');
    setIsLoadingMore(false);
    refetch();
  }, [refetch]);

  const paginatedData: PaginatedMcpServerList = {
    items: allItems,
    size: totalSize,
    pageSize: firstPageData.pageSize,
    nextPageToken,
    loadMore,
    isLoadingMore,
    hasMore: Boolean(nextPageToken),
    refresh,
  };

  return {
    mcpServers: paginatedData,
    mcpServersLoaded: loaded,
    mcpServersLoadError: error,
    refresh,
  };
};
