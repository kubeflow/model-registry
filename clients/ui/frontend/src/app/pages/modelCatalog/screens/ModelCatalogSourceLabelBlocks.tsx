import { ToggleGroup, ToggleGroupItem } from '@patternfly/react-core';
import React from 'react';
import { useSearchParams } from 'react-router-dom';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import {
  getUniqueSourceLabels,
  filterEnabledCatalogSources,
} from '~/app/pages/modelCatalog/utils/modelCatalogUtils';

type SourceLabelBlock = {
  id: string;
  label: string;
  displayName: string;
};

const ModelCatalogSourceLabelBlocks: React.FC = () => {
  const { catalogSources, updateSelectedSourceLabel, selectedSourceLabel } =
    React.useContext(ModelCatalogContext);
  const [searchParams, setSearchParams] = useSearchParams();

  React.useEffect(() => {
    const categoryFromURL = searchParams.get('category');
    if (categoryFromURL && categoryFromURL !== selectedSourceLabel) {
      updateSelectedSourceLabel(categoryFromURL);
    } else if (!categoryFromURL && selectedSourceLabel !== 'All models') {
      updateSelectedSourceLabel('All models');
    }
  }, [searchParams, selectedSourceLabel, updateSelectedSourceLabel]);

  const blocks: SourceLabelBlock[] = React.useMemo(() => {
    if (!catalogSources) {
      return [];
    }

    const enabledSources = filterEnabledCatalogSources(catalogSources);
    const uniqueLabels = getUniqueSourceLabels(enabledSources);

    const allBlock: SourceLabelBlock = {
      id: 'all',
      label: 'All models',
      displayName: 'All models',
    };

    const labelBlocks: SourceLabelBlock[] = uniqueLabels.map((label) => ({
      id: `label-${label}`,
      label,
      displayName: label,
    }));

    const noLabelsBlock: SourceLabelBlock = {
      id: 'no-labels',
      label: 'Other',
      displayName: 'Community and custom models',
    };

    return [allBlock, ...labelBlocks, noLabelsBlock];
  }, [catalogSources]);

  if (!catalogSources) {
    return null;
  }

  const handleToggleClick = (label: string) => {
    updateSelectedSourceLabel(label);
    const newSearchParams = new URLSearchParams(searchParams);
    newSearchParams.set('category', label);
    setSearchParams(newSearchParams);
  };

  return (
    <ToggleGroup aria-label="Source label selection" className="pf-v6-u-pb-xl pf-v6-u-pt-xl">
      {blocks.map((block) => (
        <ToggleGroupItem
          key={block.id}
          text={block.displayName}
          isSelected={block.label === selectedSourceLabel}
          onChange={() => {
            handleToggleClick(block.label);
          }}
        />
      ))}
    </ToggleGroup>
  );
};

export default ModelCatalogSourceLabelBlocks;
