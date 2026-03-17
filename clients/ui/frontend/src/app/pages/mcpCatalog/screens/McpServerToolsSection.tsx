import React from 'react';
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionToggle,
  Bullseye,
  Button,
  Card,
  CardBody,
  CardHeader,
  Content,
  Flex,
  FlexItem,
  Icon,
  Label,
  SearchInput,
  Spinner,
  Stack,
  StackItem,
  Title,
} from '@patternfly/react-core';
import { AngleLeftIcon, AngleRightIcon, WrenchIcon } from '@patternfly/react-icons';
import type { McpTool, McpToolAccessType } from '~/app/mcpServerCatalogTypes';
import { useMcpServerToolList } from '~/app/hooks/mcpServerCatalog/useMcpServerToolList';

const TOOLS_PAGE_SIZE = 5;

const getAccessTypeConfig = (
  accessType: McpToolAccessType,
): { text: string; color: 'blue' | 'orange' | 'purple' } => {
  switch (accessType) {
    case 'read_only':
      return { text: 'read-only', color: 'blue' };
    case 'read_write':
      return { text: 'read/write', color: 'orange' };
    case 'execute':
      return { text: 'execute', color: 'purple' };
  }
};

type McpServerToolsSectionProps = {
  serverId: string;
};

const McpServerToolsSection: React.FC<McpServerToolsSectionProps> = ({ serverId }) => {
  const [toolList, toolsLoaded, toolsLoadError] = useMcpServerToolList(serverId);

  const tools: McpTool[] = React.useMemo(
    () => (toolList.items ?? []).map((t) => t.tool),
    [toolList.items],
  );

  const [expanded, setExpanded] = React.useState<string[]>([]);
  const [filterText, setFilterText] = React.useState('');
  const [page, setPage] = React.useState(1);

  const filteredTools = React.useMemo(() => {
    if (!filterText) {
      return tools;
    }
    const lower = filterText.toLowerCase();
    return tools.filter(
      (t) =>
        t.name.toLowerCase().includes(lower) ||
        (t.description?.toLowerCase().includes(lower) ?? false),
    );
  }, [tools, filterText]);

  const totalPages = Math.max(1, Math.ceil(filteredTools.length / TOOLS_PAGE_SIZE));
  const paginatedTools = filteredTools.slice((page - 1) * TOOLS_PAGE_SIZE, page * TOOLS_PAGE_SIZE);

  React.useEffect(() => {
    setPage(1);
  }, [filterText]);

  const toggle = (id: string) => {
    setExpanded((prev) => (prev.includes(id) ? prev.filter((x) => x !== id) : [...prev, id]));
  };

  const toolsCardShell = (children: React.ReactNode) => (
    <Card data-testid="mcp-server-tools" style={{ overflow: 'visible' }}>
      <CardHeader>
        <Title headingLevel="h2" size="lg">
          <Icon isInline className="pf-v6-u-mr-sm">
            <WrenchIcon />
          </Icon>
          Tools
        </Title>
      </CardHeader>
      <CardBody>{children}</CardBody>
    </Card>
  );

  if (toolsLoadError) {
    return toolsCardShell(
      <Bullseye>
        <Content component="p" data-testid="mcp-server-tools-error">
          Unable to load tools.
        </Content>
      </Bullseye>,
    );
  }

  if (!toolsLoaded) {
    return toolsCardShell(
      <Bullseye>
        <Spinner size="lg" data-testid="mcp-server-tools-loading" />
      </Bullseye>,
    );
  }

  if (tools.length === 0) {
    return null;
  }

  return (
    <Card data-testid="mcp-server-tools" style={{ overflow: 'visible' }}>
      <CardHeader>
        <Flex
          justifyContent={{ default: 'justifyContentSpaceBetween' }}
          alignItems={{ default: 'alignItemsCenter' }}
        >
          <FlexItem>
            <Title headingLevel="h2" size="lg">
              <Icon isInline className="pf-v6-u-mr-sm">
                <WrenchIcon />
              </Icon>
              Tools
            </Title>
          </FlexItem>
          {totalPages > 1 && (
            <FlexItem>
              <Flex alignItems={{ default: 'alignItemsCenter' }} gap={{ default: 'gapSm' }}>
                <Button
                  variant="plain"
                  isDisabled={page <= 1}
                  onClick={() => setPage((p) => p - 1)}
                  aria-label="Previous tools page"
                  data-testid="mcp-tools-page-prev"
                >
                  <AngleLeftIcon />
                </Button>
                <span data-testid="mcp-tools-page-indicator">
                  {page} / {totalPages}
                </span>
                <Button
                  variant="plain"
                  isDisabled={page >= totalPages}
                  onClick={() => setPage((p) => p + 1)}
                  aria-label="Next tools page"
                  data-testid="mcp-tools-page-next"
                >
                  <AngleRightIcon />
                </Button>
              </Flex>
            </FlexItem>
          )}
        </Flex>
      </CardHeader>
      <CardBody>
        <Stack hasGutter>
          <StackItem>
            <SearchInput
              placeholder="Filter by name or description"
              value={filterText}
              onChange={(_event, value) => setFilterText(value)}
              onClear={() => setFilterText('')}
              data-testid="mcp-tools-filter"
            />
          </StackItem>
          <StackItem>
            {paginatedTools.length === 0 ? (
              <Content component="p" data-testid="mcp-tools-empty-filter">
                No tools match the filter criteria.
              </Content>
            ) : (
              <Accordion isBordered asDefinitionList={false}>
                {paginatedTools.map((tool) => {
                  const toggleId = `tool-toggle-${tool.name}`;
                  const isExpanded = expanded.includes(toggleId);
                  const accessConfig = getAccessTypeConfig(tool.accessType);

                  return (
                    <AccordionItem key={tool.name} isExpanded={isExpanded}>
                      <AccordionToggle
                        onClick={() => toggle(toggleId)}
                        id={toggleId}
                        data-testid={`mcp-tool-toggle-${tool.name}`}
                      >
                        <Flex direction={{ default: 'column' }} gap={{ default: 'gapXs' }}>
                          <Flex
                            gap={{ default: 'gapSm' }}
                            alignItems={{ default: 'alignItemsCenter' }}
                          >
                            <FlexItem>
                              <strong>{tool.name}</strong>
                            </FlexItem>
                            <FlexItem>
                              <Label color={accessConfig.color} isCompact>
                                {accessConfig.text}
                              </Label>
                            </FlexItem>
                            {tool.revoked && (
                              <FlexItem>
                                <Label
                                  color="red"
                                  isCompact
                                  data-testid={`mcp-tool-revoked-${tool.name}`}
                                >
                                  revoked
                                </Label>
                              </FlexItem>
                            )}
                          </Flex>
                          {tool.description && (
                            <FlexItem>
                              <Content component="small">{tool.description}</Content>
                            </FlexItem>
                          )}
                        </Flex>
                      </AccordionToggle>
                      <AccordionContent id={`tool-content-${tool.name}`}>
                        <Stack hasGutter>
                          {tool.revoked && tool.revokedReason && (
                            <StackItem>
                              <Content
                                component="p"
                                data-testid={`mcp-tool-revoked-reason-${tool.name}`}
                              >
                                <strong>Revoked:</strong> {tool.revokedReason}
                              </Content>
                            </StackItem>
                          )}
                          {tool.description && (
                            <StackItem>
                              <Content component="p">{tool.description}</Content>
                            </StackItem>
                          )}
                          {tool.parameters && tool.parameters.length > 0 && (
                            <StackItem>
                              <Content component="p">
                                <strong>Input Parameters:</strong>
                              </Content>
                              <Stack hasGutter>
                                {tool.parameters.map((param) => (
                                  <StackItem key={param.name}>
                                    <Flex
                                      direction={{ default: 'column' }}
                                      gap={{ default: 'gapXs' }}
                                      className="pf-v6-u-pl-md pf-v6-u-pb-sm"
                                    >
                                      <Flex
                                        gap={{ default: 'gapSm' }}
                                        alignItems={{ default: 'alignItemsCenter' }}
                                      >
                                        <FlexItem>
                                          <strong>{param.name}</strong>
                                        </FlexItem>
                                        <FlexItem>
                                          <Label color="grey" isCompact>
                                            {param.type}
                                          </Label>
                                        </FlexItem>
                                        <FlexItem>
                                          <Label color={param.required ? 'red' : 'grey'} isCompact>
                                            {param.required ? 'required' : 'optional'}
                                          </Label>
                                        </FlexItem>
                                      </Flex>
                                      {param.description && (
                                        <FlexItem>
                                          <Content component="small">{param.description}</Content>
                                        </FlexItem>
                                      )}
                                    </Flex>
                                  </StackItem>
                                ))}
                              </Stack>
                            </StackItem>
                          )}
                        </Stack>
                      </AccordionContent>
                    </AccordionItem>
                  );
                })}
              </Accordion>
            )}
          </StackItem>
        </Stack>
      </CardBody>
    </Card>
  );
};

export default McpServerToolsSection;
