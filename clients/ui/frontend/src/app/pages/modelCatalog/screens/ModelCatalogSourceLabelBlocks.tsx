import { ToggleGroup, ToggleGroupItem } from '@patternfly/react-core';
import React from 'react';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import { CategoryName, SourceLabel } from '~/app/modelCatalogTypes';
import {
  getUniqueSourceLabels,
  filterEnabledCatalogSources,
  hasSourcesWithoutLabels,
} from '~/app/pages/modelCatalog/utils/modelCatalogUtils';

type SourceLabelBlock = {
  id: string;
  label: string;
  displayName: string;
};

const ModelCatalogSourceLabelBlocks: React.FC = () => {
  const { catalogSources, updateSelectedSourceLabel, selectedSourceLabel } =
    React.useContext(ModelCatalogContext);

  const blocks: SourceLabelBlock[] = React.useMemo(() => {
    if (!catalogSources) {
      return [];
    }

    const enabledSources = filterEnabledCatalogSources(catalogSources);
    const uniqueLabels = getUniqueSourceLabels(enabledSources);
    const hasNoLabels = hasSourcesWithoutLabels(enabledSources);

    const allBlock: SourceLabelBlock = {
      id: 'all',
      label: CategoryName.allModels,
      displayName: CategoryName.allModels,
    };

    const labelBlocks: SourceLabelBlock[] = uniqueLabels.map((label) => ({
      id: `label-${label}`,
      label,
      displayName: `${label} models`,
    }));

    const blocksToReturn: SourceLabelBlock[] = [allBlock, ...labelBlocks];

    if (hasNoLabels) {
      const noLabelsBlock: SourceLabelBlock = {
        id: 'no-labels',
        label: SourceLabel.other,
        displayName: CategoryName.otherModels,
      };
      blocksToReturn.push(noLabelsBlock);
    }

    return blocksToReturn;
  }, [catalogSources]);

  if (!catalogSources) {
    return null;
  }

  const handleToggleClick = (label: string) => {
    updateSelectedSourceLabel(label);
  };

  return (
    <ToggleGroup aria-label="Source label selection" className="pf-v6-u-pb-xl pf-v6-u-pt-xl">
      {blocks.map((block) => (
        <ToggleGroupItem
          buttonId={block.id}
          data-testid={block.id}
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
