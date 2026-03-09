import * as React from 'react';
import { ToggleGroup, ToggleGroupItem } from '@patternfly/react-core';
import { McpCatalogContext } from '~/app/context/mcpCatalog/McpCatalogContext';

const ALL_SERVERS_LABEL = 'All Servers';

type SourceLabelBlock = { id: string; label?: string; displayName: string };

const McpCatalogSourceLabelBlocks: React.FC = () => {
  const { sourceLabels, sourceLabelNames, selectedSourceLabel, setSelectedSourceLabel } =
    React.useContext(McpCatalogContext);

  const blocks = React.useMemo((): SourceLabelBlock[] => {
    const allBlock: SourceLabelBlock = { id: 'all', displayName: ALL_SERVERS_LABEL };
    const labelBlocks: SourceLabelBlock[] = sourceLabels.map((label) => ({
      id: `label-${label}`,
      label,
      displayName: sourceLabelNames[label] || label,
    }));
    return [allBlock, ...labelBlocks];
  }, [sourceLabels, sourceLabelNames]);

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
