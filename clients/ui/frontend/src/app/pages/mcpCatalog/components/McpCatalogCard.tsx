import * as React from 'react';
import {
  Button,
  Card,
  CardBody,
  CardHeader,
  CardTitle,
  Flex,
  FlexItem,
  Label,
  Truncate,
} from '@patternfly/react-core';
import { TruncatedText } from 'mod-arch-shared';
import { ApplicationsIcon } from '@patternfly/react-icons';
import type { McpServer } from '~/app/mcpServerCatalogTypes';
import { getSecurityIndicatorLabels } from '~/app/pages/mcpCatalog/utils/mcpCatalogUtils';
import {
  McpCardIconType,
  McpCardIconByLabel,
  getMcpCardIconConfig,
} from '~/app/pages/mcpCatalog/constants/mcpCatalogCardIcons';

type McpCatalogCardProps = {
  server: McpServer;
};

const SecurityTag: React.FC<{ label: string }> = ({ label }) => (
  <FlexItem>
    <Flex alignItems={{ default: 'alignItemsCenter' }} gap={{ default: 'gapXs' }}>
      <McpCardIconByLabel label={label} />
      <FlexItem>{label}</FlexItem>
    </Flex>
  </FlexItem>
);

const McpCatalogCard: React.FC<McpCatalogCardProps> = React.memo(({ server }) => {
  const deploymentType =
    server.deploymentMode === 'local' ? McpCardIconType.LOCAL_TO_CLUSTER : McpCardIconType.REMOTE;
  const deploymentConfig = getMcpCardIconConfig(deploymentType);
  const securityLabels = getSecurityIndicatorLabels(server.securityIndicators);
  const serverId = String(server.id);

  return (
    <Card isFullHeight data-testid={`mcp-catalog-card-${serverId}`}>
      <CardHeader>
        <Flex
          alignItems={{ default: 'alignItemsFlexStart' }}
          justifyContent={{ default: 'justifyContentSpaceBetween' }}
          gap={{ default: 'gapXs' }}
          className="pf-v6-u-mb-md"
        >
          <FlexItem>
            <span
              className="pf-v6-u-display-inline-block pf-v6-u-font-size-2xl pf-v6-u-color-200"
              aria-hidden
            >
              <ApplicationsIcon />
            </span>
          </FlexItem>
          {server.deploymentMode === 'remote' && (
            <FlexItem>
              <Label data-testid={`mcp-catalog-card-deployment-${serverId}`}>
                {deploymentConfig.label}
              </Label>
            </FlexItem>
          )}
        </Flex>
        <CardTitle>
          <Button
            variant="link"
            isInline
            component="a"
            href="#"
            data-testid={`mcp-catalog-card-detail-link-${serverId}`}
            style={{
              fontSize: 'var(--pf-t--global--font--size--body--default)',
              fontWeight: 'var(--pf-t--global--font--weight--body--bold)',
            }}
          >
            <Truncate
              content={server.name}
              position="middle"
              tooltipPosition="top"
              data-testid={`mcp-catalog-card-name-${serverId}`}
            />
          </Button>
        </CardTitle>
      </CardHeader>
      <CardBody>
        <TruncatedText
          content={server.description ?? ''}
          maxLines={4}
          data-testid={`mcp-catalog-card-description-${serverId}`}
        />
        {securityLabels.length > 0 && (
          <Flex
            direction={{ default: 'column' }}
            gap={{ default: 'gapSm' }}
            className="pf-v6-u-mt-lg"
          >
            {securityLabels.map((tag) => (
              <SecurityTag key={tag} label={tag} />
            ))}
          </Flex>
        )}
      </CardBody>
    </Card>
  );
});
McpCatalogCard.displayName = 'McpCatalogCard';

export default McpCatalogCard;
