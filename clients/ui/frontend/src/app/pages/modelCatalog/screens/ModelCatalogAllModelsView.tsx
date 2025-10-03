import React from 'react';
import { Stack } from '@patternfly/react-core';

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

  const sourceLabels = React.useMemo(() => {
    const enabledSources = filterEnabledCatalogSources(catalogSources);
    return [...getUniqueSourceLabels(enabledSources), 'Other'];
  }, [catalogSources]);

  const handleShowMoreCategory = (categoryLabel: string) => {
    updateSelectedSourceLabel(categoryLabel);
  };

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
    </Stack>
  );
};

export default ModelCatalogAllModelsView;
