import * as React from 'react';
import { ToggleGroup, ToggleGroupItem } from '@patternfly/react-core';
import { McpCatalogContext } from '~/app/context/mcpCatalog/McpCatalogContext';
import type { McpCatalogCategoryId } from '~/app/context/mcpCatalog/McpCatalogContext';

const CATEGORY_BLOCKS: { id: McpCatalogCategoryId; displayName: string }[] = [
  { id: 'all', displayName: 'All Servers' },
  { id: 'sample', displayName: 'Sample MCP servers' },
  { id: 'other', displayName: 'Other MCP servers' },
];

const McpCatalogSourceLabelBlocks: React.FC = () => {
  const { selectedCategory, setSelectedCategory } = React.useContext(McpCatalogContext);

  return (
    <ToggleGroup aria-label="MCP category selection" data-testid="mcp-catalog-category-toggle">
      {CATEGORY_BLOCKS.map((block) => (
        <ToggleGroupItem
          buttonId={block.id}
          data-testid={`mcp-category-${block.id}`}
          key={block.id}
          text={block.displayName}
          isSelected={selectedCategory === block.id}
          onChange={() => setSelectedCategory(block.id)}
        />
      ))}
    </ToggleGroup>
  );
};

export default McpCatalogSourceLabelBlocks;
