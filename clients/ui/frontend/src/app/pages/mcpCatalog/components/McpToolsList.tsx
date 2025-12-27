import * as React from 'react';
import {
  Card,
  CardBody,
  CardExpandableContent,
  CardHeader,
  CardTitle,
  DescriptionList,
  DescriptionListDescription,
  DescriptionListGroup,
  DescriptionListTerm,
  Flex,
  FlexItem,
  Label,
  Pagination,
  PaginationVariant,
  Title,
  Tooltip,
} from '@patternfly/react-core';
import { BanIcon } from '@patternfly/react-icons';
import { McpTool } from '~/app/pages/mcpCatalog/types';

type McpToolsListProps = {
  tools: McpTool[];
};

const TOOLS_PER_PAGE = 5;

const McpToolsList: React.FC<McpToolsListProps> = ({ tools }) => {
  const [expandedTools, setExpandedTools] = React.useState<Set<string>>(new Set());
  const [page, setPage] = React.useState(1);

  const totalTools = tools.length;
  const startIndex = (page - 1) * TOOLS_PER_PAGE;
  const endIndex = Math.min(startIndex + TOOLS_PER_PAGE, totalTools);
  const paginatedTools = tools.slice(startIndex, endIndex);

  const toggleExpanded = (toolName: string) => {
    setExpandedTools((prev) => {
      const next = new Set(prev);
      if (next.has(toolName)) {
        next.delete(toolName);
      } else {
        next.add(toolName);
      }
      return next;
    });
  };

  return (
    <Card>
      <CardHeader>
        <Flex justifyContent={{ default: 'justifyContentSpaceBetween' }} style={{ width: '100%' }}>
          <FlexItem>
            <Title headingLevel="h2" size="lg">
              Tools
            </Title>
          </FlexItem>
          <FlexItem>
            <Pagination
              itemCount={totalTools}
              perPage={TOOLS_PER_PAGE}
              page={page}
              onSetPage={(_event, newPage) => setPage(newPage)}
              variant={PaginationVariant.top}
              isCompact
            />
          </FlexItem>
        </Flex>
      </CardHeader>
      <CardBody>
        {paginatedTools.map((tool) => {
          const isExpanded = expandedTools.has(tool.name);
          const isRevoked = tool.revoked === true;
          return (
            <Card
              key={tool.name}
              isExpanded={isExpanded}
              isCompact
              className="pf-v6-u-mb-sm"
              data-testid={`mcp-tool-card-${tool.name}`}
              style={
                isRevoked
                  ? {
                      opacity: 0.6,
                      backgroundColor: 'var(--pf-t--global--background--color--disabled--default)',
                    }
                  : undefined
              }
            >
              <CardHeader
                onExpand={() => toggleExpanded(tool.name)}
                isToggleRightAligned
                toggleButtonProps={{
                  'aria-label': `Toggle ${tool.name} details`,
                  'aria-expanded': isExpanded,
                }}
              >
                <Flex alignItems={{ default: 'alignItemsCenter' }} gap={{ default: 'gapSm' }}>
                  <FlexItem>
                    <CardTitle>{tool.name}</CardTitle>
                  </FlexItem>
                  {isRevoked && (
                    <FlexItem>
                      <Tooltip
                        content={
                          tool.revokedReason || 'This tool has been revoked and should not be used.'
                        }
                      >
                        <Label color="grey" icon={<BanIcon />} isCompact>
                          Revoked
                        </Label>
                      </Tooltip>
                    </FlexItem>
                  )}
                </Flex>
              </CardHeader>
              <CardExpandableContent>
                <CardBody>
                  <DescriptionList isCompact>
                    <DescriptionListGroup>
                      <DescriptionListDescription>{tool.description}</DescriptionListDescription>
                    </DescriptionListGroup>
                    {tool.parameters && tool.parameters.length > 0 && (
                      <DescriptionListGroup>
                        <DescriptionListTerm>Input Parameters:</DescriptionListTerm>
                        <DescriptionListDescription>
                          {tool.parameters.map((param) => (
                            <Card key={param.name} isCompact className="pf-v6-u-mb-sm">
                              <CardBody>
                                <Flex
                                  alignItems={{ default: 'alignItemsCenter' }}
                                  gap={{ default: 'gapSm' }}
                                >
                                  <FlexItem>
                                    <strong>{param.name}</strong>
                                  </FlexItem>
                                  <FlexItem>
                                    <Label isCompact color="grey">
                                      {param.type}
                                    </Label>
                                  </FlexItem>
                                  <FlexItem>
                                    <Label isCompact color={param.required ? 'red' : 'grey'}>
                                      {param.required ? 'required' : 'optional'}
                                    </Label>
                                  </FlexItem>
                                </Flex>
                                <div className="pf-v6-u-mt-sm">{param.description}</div>
                              </CardBody>
                            </Card>
                          ))}
                        </DescriptionListDescription>
                      </DescriptionListGroup>
                    )}
                  </DescriptionList>
                </CardBody>
              </CardExpandableContent>
            </Card>
          );
        })}
      </CardBody>
    </Card>
  );
};

export default McpToolsList;
