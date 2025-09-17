import React from 'react';
import { NotReadyError } from 'mod-arch-core';
import { CatalogModel } from '~/app/modelCatalogTypes';
import { useModelCatalogAPI } from './useModelCatalogAPI';

type PaginatedCatalogModelList = {
  items: CatalogModel[];
  size: number;
  pageSize: number;
  nextPageToken: string;
  loadMore: () => void;
  isLoadingMore: boolean;
  hasMore: boolean;
};

interface CatalogModelsState {
  items: CatalogModel[];
  nextPageToken: string;
  totalSize: number;
  isLoading: boolean;
  isLoadingMore: boolean;
  loaded: boolean;
  error: Error | undefined;
}

type CatalogModelList = [
  models: PaginatedCatalogModelList,
  catalogModelLoaded: boolean,
  catalogModelLoadError: Error | undefined,
  refresh: () => void,
];

export const useCatalogModelsBySources = (
  sourceId: string,
  pageSize = 10,
  searchQuery = '',
): CatalogModelList => {
  const { api, apiAvailable } = useModelCatalogAPI();

  const [state, setState] = React.useState<CatalogModelsState>({
    items: [],
    nextPageToken: '',
    totalSize: 0,
    isLoading: false,
    isLoadingMore: false,
    loaded: false,
    error: undefined,
  });

  const fetchModels = React.useCallback(
    async (nextPageToken?: string) => {
      const isLoadMore = Boolean(nextPageToken);

      setState((prev) => ({
        ...prev,
        isLoading: !isLoadMore,
        isLoadingMore: isLoadMore,
      }));

      try {
        if (!apiAvailable) {
          return await Promise.reject(new Error('API not yet available'));
        }
        if (!sourceId) {
          return await Promise.reject(new NotReadyError('No source id'));
        }
        const response = await api.getCatalogModelsBySource(
          {},
          sourceId,
          {
            pageSize: pageSize.toString(),
            ...(nextPageToken && { nextPageToken }),
          },
          searchQuery && searchQuery.trim() ? searchQuery.trim() : undefined,
        );

        setState((prev) => ({
          items: isLoadMore ? [...prev.items, ...response.items] : response.items,
          nextPageToken: response.nextPageToken,
          totalSize: response.size,
          isLoading: false,
          isLoadingMore: false,
          loaded: true,
          error: undefined,
        }));
      } catch (error) {
        setState((prev) => ({
          ...prev,
          isLoading: false,
          isLoadingMore: false,
          loaded: true,
          error: new Error(
            `Failed to load models ${error instanceof Error ? error.message : String(error)}`,
          ),
        }));
      }
    },
    [api, sourceId, pageSize, apiAvailable, searchQuery],
  );

  React.useEffect(() => {
    fetchModels();
  }, [fetchModels]);

  const loadMore = React.useCallback(() => {
    if (state.nextPageToken && !state.isLoadingMore) {
      fetchModels(state.nextPageToken);
    }
  }, [fetchModels, state.nextPageToken, state.isLoadingMore]);

  const refresh = React.useCallback(() => {
    setState((prev) => ({ ...prev, items: [], nextPageToken: '' }));
    fetchModels();
  }, [fetchModels]);

  return [
    {
      items: state.items,
      size: state.totalSize,
      pageSize: 10,
      nextPageToken: state.nextPageToken,
      loadMore,
      isLoadingMore: state.isLoadingMore,
      hasMore: Boolean(state.nextPageToken),
    },
    state.loaded,
    state.error,
    refresh,
  ];
};
