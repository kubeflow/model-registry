import * as React from 'react';
import { ToggleGroup, ToggleGroupItem } from '@patternfly/react-core';
import { SourceLabel } from '~/app/modelCatalogTypes';
import {
  filterEnabledCatalogSources,
  getLabelDisplayName,
  getUniqueSourceLabels,
  hasSourcesWithoutLabels,
  orderLabelsByPriority,
} from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import { McpCatalogContext } from '~/app/context/mcpCatalog/McpCatalogContext';
import { OTHER_MCP_SERVERS_DISPLAY_NAME } from '~/app/pages/mcpCatalog/const';

const ALL_SERVERS_LABEL = 'All MCP servers';

type SourceLabelBlock = { id: string; label?: string; displayName: string };

const McpCatalogSourceLabelBlocks: React.FC = () => {
  const {
    catalogSources,
    catalogLabels,
    selectedSourceLabel,
    setSelectedSourceLabel,
    emptyCategoryLabels,
  } = React.useContext(McpCatalogContext);

  const blocks: SourceLabelBlock[] = React.useMemo(() => {
    if (!catalogSources) {
      return [];
    }

    const enabledSources = filterEnabledCatalogSources(catalogSources);
    const uniqueLabels = getUniqueSourceLabels(enabledSources);
    const hasNoLabels = hasSourcesWithoutLabels(enabledSources);
    const orderedLabels = orderLabelsByPriority(uniqueLabels, catalogLabels);

    const allBlock: SourceLabelBlock = { id: 'all', displayName: ALL_SERVERS_LABEL };

    const labelBlocks: SourceLabelBlock[] = orderedLabels
      .filter((label) => !emptyCategoryLabels.has(label))
      .map((label) => ({
        id: `label-${label}`,
        label,
        displayName: getLabelDisplayName(
          label,
          catalogLabels,
          OTHER_MCP_SERVERS_DISPLAY_NAME,
          'servers',
        ),
      }));

    const result: SourceLabelBlock[] = [allBlock, ...labelBlocks];

    if (hasNoLabels && !emptyCategoryLabels.has(SourceLabel.other)) {
      result.push({
        id: 'no-labels',
        label: SourceLabel.other,
        displayName: getLabelDisplayName(
          SourceLabel.other,
          catalogLabels,
          OTHER_MCP_SERVERS_DISPLAY_NAME,
          'servers',
        ),
      });
    }

    return result;
  }, [catalogSources, catalogLabels, emptyCategoryLabels]);

  if (!catalogSources) {
    return null;
  }

  const activeCategoryCount = blocks.length - 1;
  if (activeCategoryCount <= 1) {
    return null;
  }

  const isSelected = (block: SourceLabelBlock) =>
    block.label === undefined
      ? selectedSourceLabel === undefined
      : selectedSourceLabel === block.label;

  return (
    <ToggleGroup aria-label="MCP category selection" data-testid="mcp-catalog-category-toggle">
      {blocks.map((block) => (
        <ToggleGroupItem
          buttonId={block.id}
          data-testid={`mcp-category-${block.id}`}
          key={block.id}
          text={block.displayName}
          isSelected={isSelected(block)}
          onChange={() => setSelectedSourceLabel(block.label)}
        />
      ))}
    </ToggleGroup>
  );
};

export default McpCatalogSourceLabelBlocks;
