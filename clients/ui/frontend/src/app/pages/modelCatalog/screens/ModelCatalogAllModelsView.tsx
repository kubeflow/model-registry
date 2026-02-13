import React from 'react';
import { Stack } from '@patternfly/react-core';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import {
  filterEnabledCatalogSources,
  getUniqueSourceLabels,
  hasSourcesWithoutLabels,
  orderLabelsByPriority,
} from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import { CategoryName, SourceLabel } from '~/app/modelCatalogTypes';
import CatalogCategorySection from './CatalogCategorySection';

type ModelCatalogAllModelsViewProps = {
  searchTerm: string;
};

const ModelCatalogAllModelsView: React.FC<ModelCatalogAllModelsViewProps> = ({ searchTerm }) => {
  const { catalogSources, catalogLabels, updateSelectedSourceLabel } =
    React.useContext(ModelCatalogContext);

  const sourceLabels = React.useMemo(() => {
    const enabledSources = filterEnabledCatalogSources(catalogSources);
    const uniqueLabels = getUniqueSourceLabels(enabledSources);
    // Order labels according to catalogLabels priority
    return orderLabelsByPriority(uniqueLabels, catalogLabels);
  }, [catalogSources, catalogLabels]);

  const hasSourcesWithoutLabelsValue = React.useMemo(
    () => hasSourcesWithoutLabels(catalogSources),
    [catalogSources],
  );

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
      {hasSourcesWithoutLabelsValue && (
        <CatalogCategorySection
          key={CategoryName.otherModels}
          label={SourceLabel.other}
          searchTerm={searchTerm}
          pageSize={4}
          catalogSources={catalogSources}
          onShowMore={handleShowMoreCategory}
        />
      )}
    </Stack>
  );
};

export default ModelCatalogAllModelsView;
