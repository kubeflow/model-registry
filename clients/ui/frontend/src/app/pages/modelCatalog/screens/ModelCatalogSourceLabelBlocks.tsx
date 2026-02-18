import { ToggleGroup, ToggleGroupItem } from '@patternfly/react-core';
import React from 'react';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import { CategoryName, SourceLabel } from '~/app/modelCatalogTypes';
import {
  getUniqueSourceLabels,
  filterEnabledCatalogSources,
  hasSourcesWithoutLabels,
  orderLabelsByPriority,
  getLabelDisplayName,
} from '~/app/pages/modelCatalog/utils/modelCatalogUtils';

type SourceLabelBlock = {
  id: string;
  label: string;
  displayName: string;
};

const ModelCatalogSourceLabelBlocks: React.FC = () => {
  const { catalogSources, catalogLabels, updateSelectedSourceLabel, selectedSourceLabel } =
    React.useContext(ModelCatalogContext);

  const blocks: SourceLabelBlock[] = React.useMemo(() => {
    if (!catalogSources) {
      return [];
    }

    const enabledSources = filterEnabledCatalogSources(catalogSources);
    const uniqueLabels = getUniqueSourceLabels(enabledSources);
    const hasNoLabels = hasSourcesWithoutLabels(enabledSources);

    // Order labels according to catalogLabels priority
    const orderedLabels = orderLabelsByPriority(uniqueLabels, catalogLabels);

    const allBlock: SourceLabelBlock = {
      id: 'all',
      label: CategoryName.allModels,
      displayName: CategoryName.allModels,
    };

    const labelBlocks: SourceLabelBlock[] = orderedLabels.map((label) => ({
      id: `label-${label}`,
      label,
      displayName: getLabelDisplayName(label, catalogLabels),
    }));

    const blocksToReturn: SourceLabelBlock[] = [allBlock, ...labelBlocks];

    if (hasNoLabels) {
      const noLabelsBlock: SourceLabelBlock = {
        id: 'no-labels',
        label: SourceLabel.other,
        displayName: getLabelDisplayName(SourceLabel.other, catalogLabels),
      };
      blocksToReturn.push(noLabelsBlock);
    }

    return blocksToReturn;
  }, [catalogSources, catalogLabels]);

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
