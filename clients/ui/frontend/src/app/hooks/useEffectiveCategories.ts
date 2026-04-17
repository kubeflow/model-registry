import * as React from 'react';
import { CatalogLabelList, CatalogSourceList } from '~/app/modelCatalogTypes';
import { getActiveSourceLabels } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';

type UseEffectiveCategoriesResult = {
  activeCategories: string[];
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

  const effectiveCategoriesKey = effectiveActiveCategories.join(',');

  React.useEffect(() => {
    if (catalogSourcesLoaded && isSingleCategory) {
      updateSelectedSourceLabel(effectiveActiveCategories[0]);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [catalogSourcesLoaded, isSingleCategory, effectiveCategoriesKey, updateSelectedSourceLabel]);

  return { activeCategories, effectiveActiveCategories, isSingleCategory, hasNoCategories };
};

export default useEffectiveCategories;
