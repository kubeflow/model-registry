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

export const useCatalogModelsbySources = (sourceId: string, pageSize = 10): CatalogModelList => {
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
        const response = await api.getCatalogModelsBySource({}, sourceId, {
          pageSize: pageSize.toString(),
          ...(nextPageToken && { nextPageToken }),
        });

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
    [api, sourceId, pageSize, apiAvailable],
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

// export const useCatalogModelsbySources = (
//   sourceId: string,
//   initialPageSize: number = 10,
//   loadMoreSize: number = 10,
// ): FetchState<PaginatedCatalogModelList> => {
//   const { api, apiAvailable } = useModelCatalogAPI();

//   // ✅ Single state object with proper typing
//   const [state, setState] = React.useState<{
//     allItems: CatalogModel[];
//     nextPageToken: string;
//     totalSize: number;
//     isLoadingMore: boolean;
//   }>({
//     allItems: [],
//     nextPageToken: '',
//     totalSize: 0,
//     isLoadingMore: false,
//   });

//   // ✅ Keep current state in ref to avoid circular dependencies
//   const stateRef = React.useRef(state);
//   stateRef.current = state;

//   // ✅ loadMore with stable dependencies - no circular dependency
//   const loadMore = React.useCallback(async () => {
//     const current = stateRef.current;

//     // Check if we can load more
//     if (!current.nextPageToken || current.isLoadingMore || !apiAvailable) {
//       return;
//     }

//     // Set loading state
//     setState((prev) => ({ ...prev, isLoadingMore: true }));

//     try {
//       const moreData = await api.getCatalogModelsBySource({}, sourceId, {
//         pageSize: loadMoreSize.toString(),
//         nextPageToken: current.nextPageToken,
//       });

//       // Update state with new data
//       setState((prev) => ({
//         allItems: [...prev.allItems, ...(moreData.items || [])],
//         nextPageToken: moreData.nextPageToken || '',
//         totalSize: moreData.size || prev.totalSize,
//         isLoadingMore: false,
//       }));
//     } catch (error) {
//       console.error('Load more failed:', error);
//       setState((prev) => ({ ...prev, isLoadingMore: false }));
//     }
//   }, [api, sourceId, loadMoreSize, apiAvailable]); // ✅ Stable dependencies only

//   // ✅ Initial fetch function
//   const call = React.useCallback<FetchStateCallbackPromise<PaginatedCatalogModelList>>(
//     async (opts) => {
//       if (!apiAvailable) {
//         throw new NotReadyError('API not yet available');
//       }
//       if (!sourceId) {
//         throw new NotReadyError('No source id');
//       }

//       // Reset state for fresh load
//       setState({
//         allItems: [],
//         nextPageToken: '',
//         totalSize: 0,
//         isLoadingMore: false,
//       });

//       // Make initial API call
//       const response = await api.getCatalogModelsBySource(opts, sourceId, {
//         pageSize: initialPageSize.toString(),
//       });

//       // Update state with initial data
//       const initialData = {
//         allItems: response.items || [],
//         nextPageToken: response.nextPageToken || '',
//         totalSize: response.size || 0,
//         isLoadingMore: false,
//       };

//       setState(initialData);

//       // Return data for useFetchState
//       return {
//         items: initialData.allItems,
//         size: initialData.totalSize,
//         pageSize: response.pageSize || initialPageSize,
//         nextPageToken: initialData.nextPageToken,
//         loadMore, // ✅ Stable loadMore function
//         isLoadingMore: false,
//         hasMore: Boolean(initialData.nextPageToken),
//       };
//     },
//     [api, apiAvailable, sourceId, initialPageSize, loadMore],
//   );

//   // ✅ Use useFetchState for loading/error states
//   const [data, loaded, loadError, refresh] = useFetchState(call, {
//     items: [],
//     size: 0,
//     pageSize: initialPageSize,
//     nextPageToken: '',
//     loadMore: () => {},
//     isLoadingMore: false,
//     hasMore: false,
//   });

//   // ✅ Return final data with accumulated state
//   const finalData = React.useMemo(
//     () => ({
//       ...data,
//       items: state.allItems.length > 0 ? state.allItems : data.items,
//       isLoadingMore: state.isLoadingMore,
//       hasMore: Boolean(state.nextPageToken),
//       loadMore,
//     }),
//     [data, state, loadMore],
//   );

//   return [finalData, loaded, loadError, refresh];
// };
