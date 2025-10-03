import React from 'react';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import { CatalogModelList } from '~/app/modelCatalogTypes';
import {
  filterEnabledCatalogSources,
  getUniqueSourceLabels,
} from '~/app/pages/modelCatalog/utils/modelCatalogUtils';

export type CategoryModels = {
  models: CatalogModelList | null;
  loading: boolean;
  loaded: boolean;
  error: Error | null;
};

type CatalogAllModels = {
  categoriesData: Record<string, CategoryModels>;
  allCategoriesLoaded: boolean;
  isAnyLoading: boolean;
};

export const useCatalogAllCategoriesModels = (
  searchTerm: string,
  initialPageSize = 4,
): CatalogAllModels => {
  const { catalogSources, apiState } = React.useContext(ModelCatalogContext);
  const sourceLabels = React.useMemo(() => {
    const enabledSources = filterEnabledCatalogSources(catalogSources);
    return [...getUniqueSourceLabels(enabledSources), 'Other'];
  }, [catalogSources]);

  const [categoriesData, setCategoriesData] = React.useState<Record<string, CategoryModels>>(() => {
    if (sourceLabels.length === 0) {
      return {};
    }

    return sourceLabels.reduce(
      (acc, label) => {
        acc[label] = {
          models: null,
          loading: false,
          loaded: false,
          error: null,
        };
        return acc;
      },
      {} as Record<string, CategoryModels>,
    );
  });

  React.useEffect(() => {
    if (!apiState.apiAvailable || sourceLabels.length === 0) {
      return;
    }

    const loadAllCategories = async () => {
      setCategoriesData((prev) => {
        const newData = { ...prev };
        sourceLabels.forEach((label) => {
          newData[label] = { ...newData[label], loading: true };
        });
        return newData;
      });

      const promises = sourceLabels.map(async (label) => {
        try {
          const data = await apiState.api.getCatalogModelsBySource(
            {},
            undefined,
            label,
            {
              pageSize: initialPageSize.toString(),
              nextPageToken: '',
            },
            searchTerm || undefined,
          );

          return { label, data, error: null };
        } catch (error) {
          return {
            label,
            data: null,
            error: error instanceof Error ? error : new Error('fetching models failed'),
          };
        }
      });

      const results = await Promise.allSettled(promises);

      setCategoriesData((prev) => {
        const newData = { ...prev };

        results.forEach((result) => {
          if (result.status === 'fulfilled') {
            const { label, data, error } = result.value;
            newData[label] = {
              models: data,
              loading: false,
              loaded: !error,
              error,
            };
          }
        });

        return newData;
      });
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
