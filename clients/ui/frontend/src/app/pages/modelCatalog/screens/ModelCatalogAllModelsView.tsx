import React from 'react';
import { Stack } from '@patternfly/react-core';

import { useSearchParams } from 'react-router-dom';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';

import {
  filterEnabledCatalogSources,
  getUniqueSourceLabels,
} from '~/app/pages/modelCatalog/utils/modelCatalogUtils';

import CatalogCategorySection from './CatalogCategorySection';

type ModelCatalogAllModelsViewProps = {
  searchTerm: string;
};

const ModelCatalogAllModelsView: React.FC<ModelCatalogAllModelsViewProps> = ({ searchTerm }) => {
  const { catalogSources, updateSelectedSourceLabel } = React.useContext(ModelCatalogContext);
  const [searchParams, setSearchParams] = useSearchParams();

  const sourceLabels = React.useMemo(() => {
    const enabledSources = filterEnabledCatalogSources(catalogSources);
    return getUniqueSourceLabels(enabledSources);
  }, [catalogSources]);

  const handleShowMoreCategory = React.useCallback(
    (categoryLabel: string) => {
      updateSelectedSourceLabel(categoryLabel);
      const newSearchParams = new URLSearchParams(searchParams);
      newSearchParams.set('category', categoryLabel);
      setSearchParams(newSearchParams);
    },
    [searchParams, setSearchParams, updateSelectedSourceLabel],
  );

  return (
    <Stack hasGutter>
      {sourceLabels.map((label) => (
        <CatalogCategorySection
          key={label}
          label={label}
          searchTerm={searchTerm}
          pageSize={4}
          catalogSources={catalogSources}
          onShowMore={handleShowMoreCategory}
        />
      ))}
      <CatalogCategorySection
        key="Other"
        label="Other"
        searchTerm={searchTerm}
        pageSize={4}
        catalogSources={catalogSources}
        onShowMore={handleShowMoreCategory}
        displayName="Community and custom models"
      />
    </Stack>
  );
};

export default ModelCatalogAllModelsView;
