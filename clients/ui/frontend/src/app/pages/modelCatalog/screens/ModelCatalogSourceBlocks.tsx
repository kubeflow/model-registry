import { ToggleGroup, ToggleGroupItem } from '@patternfly/react-core';
import React from 'react';
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

const ModelCatalogSourceBlocks: React.FC = () => {
  const { catalogSources, updateSelectedSourceLabel, selectedSourceLabel } =
    React.useContext(ModelCatalogContext);

  const blocks: SourceLabelBlock[] = React.useMemo(() => {
    if (!catalogSources) {
      return [];
    }

    const enabledSources = filterEnabledCatalogSources(catalogSources);
    const uniqueLabels = getUniqueSourceLabels(enabledSources);

    const allBlock: SourceLabelBlock = {
      id: 'all',
      label: '',
      displayName: 'All models',
    };

    const labelBlocks: SourceLabelBlock[] = uniqueLabels.map((label) => ({
      id: `label-${label}`,
      label,
      displayName: label,
    }));

    const noLabelsBlock: SourceLabelBlock = {
      id: 'no-labels',
      label: '',
      displayName: 'Community and custom models',
    };

    return [allBlock, ...labelBlocks, noLabelsBlock];
  }, [catalogSources]);

  if (!catalogSources) {
    return null;
  }

  return (
    <ToggleGroup aria-label="Source label selection">
      {blocks.map((block) => (
        <ToggleGroupItem
          key={block.id}
          text={block.displayName}
          isSelected={block.displayName === selectedSourceLabel}
          onChange={() => {
            updateSelectedSourceLabel(block.displayName);
          }}
        />
      ))}
    </ToggleGroup>
  );
};

export default ModelCatalogSourceBlocks;
