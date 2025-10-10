import React from 'react';
import { Stack } from '@patternfly/react-core';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import {
  filterEnabledCatalogSources,
  getUniqueSourceLabels,
} from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import { CategoryName, SourceLabel } from '~/app/modelCatalogTypes';
import CatalogCategorySection from './CatalogCategorySection';

type ModelCatalogAllModelsViewProps = {
  searchTerm: string;
};

const ModelCatalogAllModelsView: React.FC<ModelCatalogAllModelsViewProps> = ({ searchTerm }) => {
  const { catalogSources, updateSelectedSourceLabel } = React.useContext(ModelCatalogContext);

  const sourceLabels = React.useMemo(() => {
    const enabledSources = filterEnabledCatalogSources(catalogSources);
    return getUniqueSourceLabels(enabledSources);
  }, [catalogSources]);

  const handleShowMoreCategory = React.useCallback(
    (categoryLabel: string) => {
      updateSelectedSourceLabel(categoryLabel);
    },
    [updateSelectedSourceLabel],
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
        key={CategoryName.communityAndCustomModels}
        label={SourceLabel.other}
        searchTerm={searchTerm}
        pageSize={4}
        catalogSources={catalogSources}
        onShowMore={handleShowMoreCategory}
        displayName={CategoryName.communityAndCustomModels}
      />
    </Stack>
  );
};

export default ModelCatalogAllModelsView;
