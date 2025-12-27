import * as React from 'react';
import { Link } from 'react-router-dom';
import {
  Card,
  CardBody,
  CardFooter,
  CardHeader,
  CardTitle,
  Flex,
  FlexItem,
  Label,
  Skeleton,
  Truncate,
} from '@patternfly/react-core';
import McpCatalogCardBody from '~/app/pages/mcpCatalog/components/McpCatalogCardBody';
import McpCatalogLabels from '~/app/pages/mcpCatalog/components/McpCatalogLabels';
import { getMcpServerDetailsRoute } from '~/app/routes/mcpCatalog/mcpServerDetails';
import { McpServer, McpDeploymentMode } from '~/app/pages/mcpCatalog/types';

type McpCatalogCardProps = {
  server: McpServer;
};

const McpCatalogCard: React.FC<McpCatalogCardProps> = ({ server }) => {
  const isRemote = server.deploymentMode === McpDeploymentMode.REMOTE;

  return (
    <Card isFullHeight data-testid="mcp-catalog-card" key={server.id}>
      <CardHeader>
        <CardTitle>
          <Flex alignItems={{ default: 'alignItemsFlexStart' }} className="pf-v6-u-mb-md">
            {server.logo ? (
              <img src={server.logo} alt="MCP Server" style={{ height: '56px', width: '56px' }} />
            ) : (
              <Skeleton
                shape="square"
                width="56px"
                height="56px"
                screenreaderText="MCP server icon loading"
              />
            )}
            {isRemote && (
              <FlexItem align={{ default: 'alignRight' }}>
                <Label color="purple">Remote</Label>
              </FlexItem>
            )}
          </Flex>
          <Link
            to={getMcpServerDetailsRoute(server.id)}
            data-testid="mcp-catalog-detail-link"
            style={{
              fontSize: 'var(--pf-t--global--font--size--body--default)',
              fontWeight: 'var(--pf-t--global--font--weight--body--bold)',
              textDecoration: 'none',
            }}
          >
            <Truncate
              data-testid="mcp-catalog-card-name"
              content={server.name}
              position="middle"
              tooltipPosition="top"
              style={{ textDecoration: 'underline' }}
            />
          </Link>
        </CardTitle>
      </CardHeader>
      <CardBody>
        <McpCatalogCardBody server={server} />
      </CardBody>
      <CardFooter>
        <McpCatalogLabels
          tags={server.tags}
          provider={server.provider}
          numLabels={isRemote ? 2 : 3}
        />
      </CardFooter>
    </Card>
  );
};

export default McpCatalogCard;
