import React from 'react';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import { CatalogModelList } from '~/app/modelCatalogTypes';
import {
  filterEnabledCatalogSources,
  getUniqueSourceLabels,
} from '~/app/pages/modelCatalog/utils/modelCatalogUtils';

export type CategoryData = {
  models: CatalogModelList | null;
  loading: boolean;
  loaded: boolean;
  error: Error | undefined;
};

type CatalogAllModels = {
  categoriesData: Record<string, CategoryData>;
  allCategoriesLoaded: boolean;
  isAnyLoading: boolean;
};

export const useCatalogAllCategoriesModels = (
  searchTerm: string,
  initialPageSize = 4,
): CatalogAllModels => {
  const { catalogSources, apiState } = React.useContext(ModelCatalogContext);
  const [categoriesData, setCategoriesData] = React.useState<Record<string, CategoryData>>({});

  const sourceLabels = React.useMemo(() => {
    const enabledSources = filterEnabledCatalogSources(catalogSources);
    return getUniqueSourceLabels(enabledSources);
  }, [catalogSources]);

  React.useEffect(() => {
    if (sourceLabels.length > 0) {
      const initialData: Record<string, CategoryData> = {};
      sourceLabels.forEach((label) => {
        initialData[label] = {
          models: null,
          loading: false,
          loaded: false,
          error: undefined,
        };
      });
      setCategoriesData(initialData);
    }
  }, [sourceLabels]);

  React.useEffect(() => {
    if (!apiState.apiAvailable || sourceLabels.length === 0) {
      return;
    }

    const loadAllCategories = async () => {
      const promises = sourceLabels.map(async (label) => {
        setCategoriesData((prev) => ({
          ...prev,
          [label]: { ...prev[label], loading: true },
        }));

        try {
          const data = await apiState.api.getCatalogModelsBySource(
            { signal: new AbortController().signal },
            undefined,
            label,
            {
              pageSize: initialPageSize.toString(),
              nextPageToken: '',
            },
            searchTerm || undefined,
          );

          setCategoriesData((prev) => ({
            ...prev,
            [label]: {
              models: data,
              loading: false,
              loaded: true,
              error: undefined,
            },
          }));
        } catch (error) {
          setCategoriesData((prev) => ({
            ...prev,
            [label]: {
              ...prev[label],
              loading: false,
              loaded: true,
              error: error instanceof Error ? error : new Error('fetching models failed'),
            },
          }));
        }
      });

      await Promise.allSettled(promises);
    };

    loadAllCategories();
  }, [apiState, sourceLabels, searchTerm, initialPageSize]);

  const allCategoriesLoaded = React.useMemo(
    () => Object.values(categoriesData).every((data) => data.loaded),
    [categoriesData],
  );

  const isAnyLoading = React.useMemo(
    () => Object.values(categoriesData).some((data) => data.loading),
    [categoriesData],
  );

  return {
    categoriesData,
    allCategoriesLoaded,
    isAnyLoading,
  };
};
