import * as React from 'react';
import {
  Card,
  CardBody,
  CardHeader,
  ClipboardCopy,
  Content,
  DescriptionList,
  DescriptionListDescription,
  DescriptionListGroup,
  DescriptionListTerm,
  Icon,
  Label,
  LabelGroup,
  PageSection,
  Sidebar,
  SidebarContent,
  SidebarPanel,
  Stack,
  StackItem,
  Title,
} from '@patternfly/react-core';
import { OutlinedClockIcon } from '@patternfly/react-icons';
import type { McpServer } from '~/app/mcpServerCatalogTypes';
import ExternalLink from '~/app/shared/components/ExternalLink';
import MarkdownComponent from '~/app/shared/markdown/MarkdownComponent';
import ModelTimestamp from '~/app/pages/modelRegistry/screens/components/ModelTimestamp';

type McpServerDetailsViewProps = {
  server: McpServer;
};

const VISIBLE_LABELS = 3;

const McpServerDetailsView: React.FC<McpServerDetailsViewProps> = ({ server }) => {
  const deploymentModeLabel = React.useMemo(() => {
    if (!server.deploymentMode) {
      return 'N/A';
    }
    return server.deploymentMode === 'local' ? 'Local to cluster' : 'Remote';
  }, [server.deploymentMode]);

  const transportTypeLabel = React.useMemo(() => {
    if (!server.transports || server.transports.length === 0) {
      return 'N/A';
    }
    return server.transports
      .map((t) => {
        switch (t) {
          case 'http':
            return 'http-streaming';
          case 'sse':
            return 'SSE';
          case 'stdio':
            return 'stdio';
          default:
            return t;
        }
      })
      .join(', ');
  }, [server.transports]);

  return (
    <PageSection hasBodyWrapper={false} isFilled padding={{ default: 'noPadding' }}>
      <Sidebar hasGutter isPanelRight>
        <SidebarContent style={{ minWidth: 0, overflow: 'hidden' }}>
          <Stack hasGutter>
            <StackItem>
              <Card>
                <CardHeader>
                  <Title headingLevel="h2" size="lg">
                    Description
                  </Title>
                </CardHeader>
                <CardBody>
                  <Content className="pf-v6-u-text-break-word">
                    <p data-testid="mcp-server-description">
                      {server.description || 'No description'}
                    </p>
                  </Content>
                </CardBody>
              </Card>
            </StackItem>
            <StackItem>
              <Card>
                <CardHeader>
                  <Title headingLevel="h2" size="lg">
                    README
                  </Title>
                </CardHeader>
                <CardBody>
                  {!server.readme && (
                    <Content component="p" data-testid="mcp-server-no-readme">
                      No README available
                    </Content>
                  )}
                  {server.readme && (
                    <MarkdownComponent
                      data={server.readme}
                      dataTestId="mcp-server-readme-markdown"
                      maxHeading={3}
                    />
                  )}
                </CardBody>
              </Card>
            </StackItem>
          </Stack>
        </SidebarContent>
        <SidebarPanel width={{ default: 'width_33' }}>
          <Card>
            <CardHeader>
              <Title headingLevel="h2" size="lg">
                Server details
              </Title>
            </CardHeader>
            <CardBody>
              <DescriptionList>
                {server.tags && server.tags.length > 0 && (
                  <DescriptionListGroup>
                    <DescriptionListTerm>Labels</DescriptionListTerm>
                    <DescriptionListDescription>
                      <LabelGroup numLabels={VISIBLE_LABELS} isCompact>
                        {server.tags.map((tag) => (
                          <Label key={tag} variant="outline" data-testid="mcp-server-detail-label">
                            {tag}
                          </Label>
                        ))}
                      </LabelGroup>
                    </DescriptionListDescription>
                  </DescriptionListGroup>
                )}
                <DescriptionListGroup>
                  <DescriptionListTerm>License</DescriptionListTerm>
                  <DescriptionListDescription>
                    {server.licenseLink ? (
                      <ExternalLink
                        text={server.license || 'Agreement'}
                        to={server.licenseLink}
                        testId="mcp-server-license-link"
                      />
                    ) : (
                      <span data-testid="mcp-server-license">
                        {server.license || 'N/A'}
                      </span>
                    )}
                  </DescriptionListDescription>
                </DescriptionListGroup>
                <DescriptionListGroup>
                  <DescriptionListTerm>Version</DescriptionListTerm>
                  <DescriptionListDescription data-testid="mcp-server-version">
                    {server.version || 'N/A'}
                  </DescriptionListDescription>
                </DescriptionListGroup>
                <DescriptionListGroup>
                  <DescriptionListTerm>Deployment mode</DescriptionListTerm>
                  <DescriptionListDescription data-testid="mcp-server-deployment-mode">
                    {deploymentModeLabel}
                  </DescriptionListDescription>
                </DescriptionListGroup>
                <DescriptionListGroup>
                  <DescriptionListTerm>Transport type</DescriptionListTerm>
                  <DescriptionListDescription data-testid="mcp-server-transport-type">
                    {transportTypeLabel}
                  </DescriptionListDescription>
                </DescriptionListGroup>
                {server.artifacts && server.artifacts.length > 0 && (
                  <DescriptionListGroup>
                    <DescriptionListTerm>Artifacts</DescriptionListTerm>
                    <DescriptionListDescription>
                      <Stack hasGutter>
                        {server.artifacts.map((artifact) => (
                          <StackItem key={artifact.uri}>
                            <ClipboardCopy
                              hoverTip="Copy"
                              clickTip="Copied"
                              isReadOnly
                              data-testid="mcp-server-artifact-copy"
                            >
                              {artifact.uri}
                            </ClipboardCopy>
                          </StackItem>
                        ))}
                      </Stack>
                    </DescriptionListDescription>
                  </DescriptionListGroup>
                )}
                {server.provider && (
                  <DescriptionListGroup>
                    <DescriptionListTerm>Provider</DescriptionListTerm>
                    <DescriptionListDescription data-testid="mcp-server-provider">
                      {server.provider}
                    </DescriptionListDescription>
                  </DescriptionListGroup>
                )}
                {(server.sourceCode || server.repositoryUrl) && (
                  <DescriptionListGroup>
                    <DescriptionListTerm>Source Code</DescriptionListTerm>
                    <DescriptionListDescription>
                      {server.repositoryUrl ? (
                        <ExternalLink
                          text={server.sourceCode || server.repositoryUrl}
                          to={server.repositoryUrl}
                          testId="mcp-server-source-code-link"
                        />
                      ) : (
                        <span data-testid="mcp-server-source-code">
                          {server.sourceCode}
                        </span>
                      )}
                    </DescriptionListDescription>
                  </DescriptionListGroup>
                )}
                {server.lastUpdated && (
                  <DescriptionListGroup>
                    <DescriptionListTerm>Last modified</DescriptionListTerm>
                    <DescriptionListDescription>
                      <Icon isInline style={{ marginRight: 4 }}>
                        <OutlinedClockIcon />
                      </Icon>
                      <ModelTimestamp timeSinceEpoch={server.lastUpdated} />
                    </DescriptionListDescription>
                  </DescriptionListGroup>
                )}
              </DescriptionList>
            </CardBody>
          </Card>
        </SidebarPanel>
      </Sidebar>
    </PageSection>
  );
};

export default McpServerDetailsView;
