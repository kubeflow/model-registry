import * as React from 'react';
import {
  Button,
  Card,
  CardBody,
  CardHeader,
  Content,
  DescriptionList,
  DescriptionListDescription,
  DescriptionListGroup,
  DescriptionListTerm,
  Icon,
  PageSection,
  Sidebar,
  SidebarContent,
  SidebarPanel,
  Stack,
  StackItem,
  Title,
} from '@patternfly/react-core';
import { GithubIcon, OutlinedClockIcon } from '@patternfly/react-icons';
import { InlineTruncatedClipboardCopy } from 'mod-arch-shared';
import text from '@patternfly/react-styles/css/utilities/Text/text';
import { McpServer } from '~/app/pages/mcpCatalog/types';
import {
  formatDeploymentMode,
  formatTransports,
  isRemoteMcpServer,
} from '~/app/pages/mcpCatalog/utils/mcpCatalogUtils';
import McpToolsList from '~/app/pages/mcpCatalog/components/McpToolsList';
import McpCatalogLabels from '~/app/pages/mcpCatalog/components/McpCatalogLabels';
import ExternalLink from '~/app/shared/components/ExternalLink';
import MarkdownComponent from '~/app/shared/markdown/MarkdownComponent';

type McpServerDetailsViewProps = {
  server: McpServer;
};

const McpServerDetailsView: React.FC<McpServerDetailsViewProps> = ({ server }) => {
  const isRemote = isRemoteMcpServer(server.deploymentMode);

  return (
    <PageSection hasBodyWrapper={false} isFilled padding={{ default: 'noPadding' }}>
      <Sidebar hasGutter isPanelRight>
        <SidebarContent style={{ minWidth: 0, overflow: 'hidden' }}>
          <Stack hasGutter>
            {/* Description Section */}
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

            {/* MCP Server Card Section */}
            <StackItem>
              <Card>
                <CardHeader>
                  <Title headingLevel="h2" size="lg">
                    MCP server card
                  </Title>
                </CardHeader>
                <CardBody>
                  {!server.readme && <p className={text.textColorDisabled}>No MCP server card</p>}
                  {server.readme && (
                    <MarkdownComponent
                      data={server.readme}
                      dataTestId="mcp-card-markdown"
                      maxHeading={3}
                    />
                  )}
                </CardBody>
              </Card>
            </StackItem>
          </Stack>
        </SidebarContent>

        {/* Details Panel */}
        <SidebarPanel width={{ default: 'width_33' }}>
          <Stack hasGutter>
            <StackItem>
              <Card>
                <CardHeader>
                  <Title headingLevel="h2" size="lg">
                    MCP Server Details
                  </Title>
                </CardHeader>
                <CardBody>
                  <DescriptionList>
                    {/* Tags - matching Model's Labels pattern */}
                    {server.tags && server.tags.length > 0 && (
                      <DescriptionListGroup>
                        <DescriptionListTerm>Tags</DescriptionListTerm>
                        <DescriptionListDescription>
                          <McpCatalogLabels tags={server.tags} numLabels={5} />
                        </DescriptionListDescription>
                      </DescriptionListGroup>
                    )}
                    {/* License with link - matching Model's License pattern */}
                    {server.license_link && (
                      <DescriptionListGroup>
                        <DescriptionListTerm>License</DescriptionListTerm>
                        <ExternalLink
                          text="Agreement"
                          to={server.license_link}
                          testId="mcp-license-link"
                        />
                      </DescriptionListGroup>
                    )}
                    <DescriptionListGroup>
                      <DescriptionListTerm>Provider</DescriptionListTerm>
                      <DescriptionListDescription>
                        {server.provider || 'N/A'}
                      </DescriptionListDescription>
                    </DescriptionListGroup>
                    {!isRemote && server.version && (
                      <DescriptionListGroup>
                        <DescriptionListTerm>Version</DescriptionListTerm>
                        <DescriptionListDescription>{server.version}</DescriptionListDescription>
                      </DescriptionListGroup>
                    )}
                    <DescriptionListGroup>
                      <DescriptionListTerm>Transport(s)</DescriptionListTerm>
                      <DescriptionListDescription>
                        <strong>{formatTransports(server.transports)}</strong>
                      </DescriptionListDescription>
                    </DescriptionListGroup>
                    <DescriptionListGroup>
                      <DescriptionListTerm>Deployment Mode</DescriptionListTerm>
                      <DescriptionListDescription>
                        {formatDeploymentMode(server.deploymentMode)}
                      </DescriptionListDescription>
                    </DescriptionListGroup>
                    {isRemote && server.endpoints && (
                      <>
                        {server.endpoints.http && (
                          <DescriptionListGroup>
                            <DescriptionListTerm>HTTP Endpoint</DescriptionListTerm>
                            <DescriptionListDescription>
                              <InlineTruncatedClipboardCopy
                                testId="mcp-server-http-endpoint"
                                textToCopy={server.endpoints.http}
                              />
                            </DescriptionListDescription>
                          </DescriptionListGroup>
                        )}
                        {server.endpoints.sse && (
                          <DescriptionListGroup>
                            <DescriptionListTerm>SSE Endpoint</DescriptionListTerm>
                            <DescriptionListDescription>
                              <InlineTruncatedClipboardCopy
                                testId="mcp-server-sse-endpoint"
                                textToCopy={server.endpoints.sse}
                              />
                            </DescriptionListDescription>
                          </DescriptionListGroup>
                        )}
                      </>
                    )}
                    {!isRemote && (
                      <>
                        <DescriptionListGroup>
                          <DescriptionListTerm>MCP server location</DescriptionListTerm>
                          <DescriptionListDescription>
                            {server.artifacts &&
                            server.artifacts.length > 0 &&
                            server.artifacts[0].uri ? (
                              <InlineTruncatedClipboardCopy
                                testId="mcp-server-location"
                                textToCopy={server.artifacts[0].uri}
                              />
                            ) : (
                              'N/A'
                            )}
                          </DescriptionListDescription>
                        </DescriptionListGroup>
                        <DescriptionListGroup>
                          <DescriptionListTerm>Source Code</DescriptionListTerm>
                          <DescriptionListDescription>
                            {server.sourceCode ? (
                              <Button
                                variant="link"
                                isInline
                                component="a"
                                href={
                                  server.repositoryUrl || `https://github.com/${server.sourceCode}`
                                }
                                target="_blank"
                                rel="noopener noreferrer"
                                icon={<GithubIcon />}
                              >
                                {server.sourceCode}
                              </Button>
                            ) : (
                              'N/A'
                            )}
                          </DescriptionListDescription>
                        </DescriptionListGroup>
                        <DescriptionListGroup>
                          <DescriptionListTerm>Last Modified</DescriptionListTerm>
                          <DescriptionListDescription>
                            <Icon isInline style={{ marginRight: 4 }}>
                              <OutlinedClockIcon />
                            </Icon>
                            {server.lastUpdated
                              ? new Date(server.lastUpdated).toLocaleDateString('en-US', {
                                  month: '2-digit',
                                  day: '2-digit',
                                  year: 'numeric',
                                })
                              : 'N/A'}
                          </DescriptionListDescription>
                        </DescriptionListGroup>
                        <DescriptionListGroup>
                          <DescriptionListTerm>Published</DescriptionListTerm>
                          <DescriptionListDescription>
                            <Icon isInline style={{ marginRight: 4 }}>
                              <OutlinedClockIcon />
                            </Icon>
                            {server.publishedDate
                              ? new Date(server.publishedDate).toLocaleDateString('en-US', {
                                  month: '2-digit',
                                  day: '2-digit',
                                  year: 'numeric',
                                })
                              : 'N/A'}
                          </DescriptionListDescription>
                        </DescriptionListGroup>
                      </>
                    )}
                  </DescriptionList>
                </CardBody>
              </Card>
            </StackItem>

            {/* Tools Section */}
            {server.tools && server.tools.length > 0 && (
              <StackItem>
                <McpToolsList tools={server.tools} />
              </StackItem>
            )}
          </Stack>
        </SidebarPanel>
      </Sidebar>
    </PageSection>
  );
};

export default McpServerDetailsView;
