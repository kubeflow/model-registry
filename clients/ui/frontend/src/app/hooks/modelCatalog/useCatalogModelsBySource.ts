import { useFetchState, FetchStateCallbackPromise } from 'mod-arch-core';
import React from 'react';
import {
  CatalogFilterOptionsList,
  CatalogModel,
  CatalogModelList,
  ModelCatalogFilterStates,
} from '~/app/modelCatalogTypes';
import { useModelCatalogAPI } from './useModelCatalogAPI';

type PaginatedCatalogModelList = {
  items: CatalogModel[];
  size: number;
  pageSize: number;
  nextPageToken: string;
  loadMore: () => void;
  isLoadingMore: boolean;
  hasMore: boolean;
  refresh: () => void;
};

type ModelList = {
  catalogModels: PaginatedCatalogModelList;
  catalogModelsLoaded: boolean;
  catalogModelsLoadError: Error | undefined;
  refresh: () => void;
};

export const useCatalogModelsBySources = (
  sourceId?: string,
  sourceLabel?: string,
  pageSize = 10,
  searchQuery = '',
  filterData?: ModelCatalogFilterStates,
  filterOptions?: CatalogFilterOptionsList | null,
): ModelList => {
  const { api, apiAvailable } = useModelCatalogAPI();

  const [allItems, setAllItems] = React.useState<CatalogModel[]>([]);
  const [totalSize, setTotalSize] = React.useState(0);
  const [nextPageToken, setNextPageToken] = React.useState('');
  const [isLoadingMore, setIsLoadingMore] = React.useState(false);

  const fetchModels = React.useCallback<FetchStateCallbackPromise<CatalogModelList>>(
    (opts) => {
      if (!apiAvailable) {
        return Promise.reject(new Error('API not yet available'));
      }

      return api.getCatalogModelsBySource(
        opts,
        sourceId,
        sourceLabel,
        { pageSize: pageSize.toString() },
        searchQuery.trim() || undefined,
        filterData,
        filterOptions,
      );
    },
    [api, apiAvailable, sourceId, pageSize, searchQuery, filterData, filterOptions, sourceLabel],
  );

  const [firstPageData, loaded, error, refetch] = useFetchState(
    fetchModels,
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
      const response = await api.getCatalogModelsBySource(
        {},
        sourceId,
        sourceLabel,
        {
          pageSize: pageSize.toString(),
          nextPageToken,
        },
        searchQuery.trim() || undefined,
        filterData,
        filterOptions,
      );

      setAllItems((prev) => [...prev, ...response.items]);
      setTotalSize(response.size);
      setNextPageToken(response.nextPageToken);
    } catch (err) {
      throw new Error(
        `Failed to load more models: ${err instanceof Error ? err.message : String(err)}`,
      );
    } finally {
      setIsLoadingMore(false);
    }
  }, [
    api,
    apiAvailable,
    sourceId,
    pageSize,
    searchQuery,
    nextPageToken,
    isLoadingMore,
    sourceLabel,
    filterData,
    filterOptions,
  ]);

  React.useEffect(() => {
    setAllItems([]);
    setTotalSize(0);
    setNextPageToken('');
    setIsLoadingMore(false);
  }, [sourceId, searchQuery, sourceLabel, filterData, filterOptions]);

  const refresh = React.useCallback(() => {
    setAllItems([]);
    setTotalSize(0);
    setNextPageToken('');
    setIsLoadingMore(false);
    refetch();
  }, [refetch]);

  const paginatedData: PaginatedCatalogModelList = {
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
    catalogModels: paginatedData,
    catalogModelsLoaded: loaded,
    catalogModelsLoadError: error,
    refresh,
  };
};
