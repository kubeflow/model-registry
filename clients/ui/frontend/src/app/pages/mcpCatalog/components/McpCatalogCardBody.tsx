import * as React from 'react';
import { Stack, StackItem } from '@patternfly/react-core';
import McpSecurityIndicators from '~/app/pages/mcpCatalog/components/McpSecurityIndicators';
import { McpServer } from '~/app/pages/mcpCatalog/types';

type McpCatalogCardBodyProps = {
  server: McpServer;
};

const McpCatalogCardBody: React.FC<McpCatalogCardBodyProps> = ({ server }) => (
  <Stack hasGutter>
    <StackItem>
      <div
        data-testid="mcp-catalog-card-description"
        style={{
          overflow: 'hidden',
          textOverflow: 'ellipsis',
          WebkitLineClamp: 3,
          WebkitBoxOrient: 'vertical',
          display: '-webkit-box',
          minHeight: '4.5em',
        }}
      >
        {server.description}
      </div>
    </StackItem>
    {server.securityIndicators && (
      <StackItem>
        <McpSecurityIndicators indicators={server.securityIndicators} />
      </StackItem>
    )}
  </Stack>
);

export default McpCatalogCardBody;
