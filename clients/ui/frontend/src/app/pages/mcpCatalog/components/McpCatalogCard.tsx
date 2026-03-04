import * as React from 'react';
import {
  Card,
  CardBody,
  CardHeader,
  CardTitle,
  Content,
  Flex,
  FlexItem,
  Label,
  Truncate,
} from '@patternfly/react-core';
import { ApplicationsIcon } from '@patternfly/react-icons';
import { Link } from 'react-router-dom';
import type { McpServerMock } from '~/app/pages/mcpCatalog/types/mcpServer';
import {
  McpCardIconType,
  McpCardIconByLabel,
  getMcpCardIconConfig,
} from '~/app/pages/mcpCatalog/constants/mcpCatalogCardIcons';

const MCP_DETAILS_LINK_PLACEHOLDER = '#';

type McpCatalogCardProps = {
  server: McpServerMock;
};

const SecurityTag: React.FC<{ label: string }> = ({ label }) => (
  <FlexItem>
    <span className="pf-v5-u-display-inline-flex pf-v5-u-align-items-center">
      <McpCardIconByLabel label={label} className="pf-v5-u-mr-xs" />
      <span>{label}</span>
    </span>
  </FlexItem>
);

const McpCatalogCard: React.FC<McpCatalogCardProps> = ({ server }) => {
  const deploymentType =
    server.deploymentMode === 'Local' ? McpCardIconType.LOCAL_TO_CLUSTER : McpCardIconType.REMOTE;
  const deploymentConfig = getMcpCardIconConfig(deploymentType);

  return (
    <Card
      isFullHeight
      style={{ minHeight: '296.58px' }}
      data-testid={`mcp-catalog-card-${server.id}`}
    >
      <CardHeader>
        <Flex
          alignItems={{ default: 'alignItemsFlexStart' }}
          justifyContent={{ default: 'justifyContentSpaceBetween' }}
          style={{ gap: '4px' }}
          className="pf-v6-u-mb-md"
        >
          <FlexItem>
            <span
              className="pf-v5-u-display-inline-block"
              style={{ fontSize: '2rem', color: 'var(--pf-v5-global--default-color--200)' }}
              aria-hidden
            >
              <ApplicationsIcon />
            </span>
          </FlexItem>
          <FlexItem>
            <Label variant="outline" data-testid={`mcp-catalog-card-deployment-${server.id}`}>
              {deploymentConfig.label}
            </Label>
          </FlexItem>
        </Flex>
        <CardTitle>
          <Link to={MCP_DETAILS_LINK_PLACEHOLDER}>
            <Truncate
              content={server.name}
              position="middle"
              tooltipPosition="top"
              data-testid={`mcp-catalog-card-name-${server.id}`}
              style={{ textDecoration: 'underline', color: 'var(--pf-v5-global--link--Color)' }}
            />
          </Link>
        </CardTitle>
      </CardHeader>
      <CardBody>
        <Content
          component="p"
          className="pf-v5-u-mb-md"
          data-testid={`mcp-catalog-card-description-${server.id}`}
        >
          {server.description}
        </Content>
        {server.securityVerification.length > 0 && (
          <Flex
            direction={{ default: 'column' }}
            gap={{ default: 'gapSm' }}
            className="pf-v5-u-mt-lg"
          >
            {server.securityVerification.map((tag) => (
              <SecurityTag key={tag} label={tag} />
            ))}
          </Flex>
        )}
      </CardBody>
    </Card>
  );
};

export default McpCatalogCard;
