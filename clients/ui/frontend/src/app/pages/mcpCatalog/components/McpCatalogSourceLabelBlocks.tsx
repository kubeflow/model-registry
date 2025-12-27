import { ToggleGroup, ToggleGroupItem } from '@patternfly/react-core';
import React from 'react';
import { useMcpCatalog } from '~/app/context/mcpCatalog/McpCatalogContext';
import { McpCategoryName, McpSourceLabel } from '~/app/pages/mcpCatalog/types';
import {
  getUniqueMcpSourceLabels,
  filterEnabledMcpSources,
  hasMcpSourcesWithoutLabels,
} from '~/app/pages/mcpCatalog/utils/mcpCatalogUtils';

type SourceLabelBlock = {
  id: string;
  label: string;
  displayName: string;
};

const McpCatalogSourceLabelBlocks: React.FC = () => {
  const { mcpSources, updateSelectedSourceLabel, selectedSourceLabel } = useMcpCatalog();

  const blocks: SourceLabelBlock[] = React.useMemo(() => {
    if (!mcpSources) {
      return [];
    }

    const enabledSources = filterEnabledMcpSources(mcpSources);
    const uniqueLabels = getUniqueMcpSourceLabels(enabledSources);
    const hasNoLabels = hasMcpSourcesWithoutLabels(mcpSources);

    const allBlock: SourceLabelBlock = {
      id: 'all',
      label: McpCategoryName.allServers,
      displayName: McpCategoryName.allServers,
    };

    const labelBlocks: SourceLabelBlock[] = uniqueLabels.map((label) => ({
      id: `label-${label.toLowerCase().replace(/\s+/g, '-')}`,
      label,
      displayName: `${label} servers`,
    }));

    const blocksToReturn: SourceLabelBlock[] = [allBlock, ...labelBlocks];

    if (hasNoLabels) {
      const noLabelsBlock: SourceLabelBlock = {
        id: 'no-labels',
        label: McpSourceLabel.other,
        displayName: `${McpCategoryName.communityAndCustomServers} servers`,
      };
      blocksToReturn.push(noLabelsBlock);
    }

    return blocksToReturn;
  }, [mcpSources]);

  if (!mcpSources) {
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

export default McpCatalogSourceLabelBlocks;
