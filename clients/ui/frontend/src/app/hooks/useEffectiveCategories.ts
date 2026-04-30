import * as React from 'react';
import { CatalogLabelList, CatalogSourceList } from '~/app/modelCatalogTypes';
import { getActiveSourceLabels } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';

type UseEffectiveCategoriesResult = {
  effectiveActiveCategories: string[];
  isSingleCategory: boolean;
  hasNoCategories: boolean;
};

const useEffectiveCategories = (
  catalogSources: CatalogSourceList | null,
  catalogLabels: CatalogLabelList | null,
  emptyCategoryLabels: Set<string>,
  catalogSourcesLoaded: boolean,
  updateSelectedSourceLabel: (label: string | undefined) => void,
): UseEffectiveCategoriesResult => {
  const activeCategories = React.useMemo(
    () => getActiveSourceLabels(catalogSources, catalogLabels),
    [catalogSources, catalogLabels],
  );

  const effectiveActiveCategories = React.useMemo(
    () => activeCategories.filter((c) => !emptyCategoryLabels.has(c)),
    [activeCategories, emptyCategoryLabels],
  );

  const isSingleCategory = effectiveActiveCategories.length === 1;
  const hasNoCategories = effectiveActiveCategories.length === 0;

  React.useEffect(() => {
    if (catalogSourcesLoaded && isSingleCategory) {
      updateSelectedSourceLabel(effectiveActiveCategories[0]);
    }
  }, [
    catalogSourcesLoaded,
    isSingleCategory,
    effectiveActiveCategories,
    updateSelectedSourceLabel,
  ]);

  return { effectiveActiveCategories, isSingleCategory, hasNoCategories };
};

export default useEffectiveCategories;
