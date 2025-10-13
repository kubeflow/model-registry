import { ToggleGroup, ToggleGroupItem } from '@patternfly/react-core';
import React from 'react';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import { CategoryName, SourceLabel } from '~/app/modelCatalogTypes';
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

  const blocks: SourceLabelBlock[] = React.useMemo(() => {
    if (!catalogSources) {
      return [];
    }

    const enabledSources = filterEnabledCatalogSources(catalogSources);
    const uniqueLabels = getUniqueSourceLabels(enabledSources);

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

    const noLabelsBlock: SourceLabelBlock = {
      id: 'no-labels',
      label: SourceLabel.other,
      displayName: `${CategoryName.communityAndCustomModels} models`,
    };

    return [allBlock, ...labelBlocks, noLabelsBlock];
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
