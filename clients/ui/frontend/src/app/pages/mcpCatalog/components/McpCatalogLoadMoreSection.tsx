import * as React from 'react';
import { Button, Flex, StackItem } from '@patternfly/react-core';

type McpCatalogLoadMoreSectionProps = {
  onLoadMore: () => void;
};

const McpCatalogLoadMoreSection: React.FC<McpCatalogLoadMoreSectionProps> = ({ onLoadMore }) => (
  <StackItem>
    <Flex justifyContent={{ default: 'justifyContentCenter' }}>
      <Button variant="secondary" data-testid="mcp-load-more-button" onClick={onLoadMore}>
        Load more MCP servers
      </Button>
    </Flex>
  </StackItem>
);

export default McpCatalogLoadMoreSection;
