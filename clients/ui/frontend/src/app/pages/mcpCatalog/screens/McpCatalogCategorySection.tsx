import * as React from 'react';
import {
  Button,
  Content,
  Flex,
  FlexItem,
  Grid,
  GridItem,
  StackItem,
  Title,
} from '@patternfly/react-core';
import { ArrowRightIcon } from '@patternfly/react-icons';
import type { McpServer } from '~/app/mcpServerCatalogTypes';
import McpCatalogCard from '~/app/pages/mcpCatalog/components/McpCatalogCard';

type McpCatalogCategorySectionProps = {
  title: string;
  description?: string;
  servers: McpServer[];
  maxItems?: number;
  onShowAll?: () => void;
};

const McpCatalogCategorySection: React.FC<McpCatalogCategorySectionProps> = React.memo(
  ({ title, description, servers, maxItems, onShowAll }) => {
    if (servers.length === 0) {
      return null;
    }

    const displayItems = maxItems ? servers.slice(0, maxItems) : servers;
    const hasMore = maxItems !== undefined && servers.length > maxItems;

    return (
      <StackItem className="pf-v6-u-pb-xl">
        <Flex
          alignItems={{ default: 'alignItemsCenter' }}
          justifyContent={{ default: 'justifyContentSpaceBetween' }}
          className="pf-v6-u-mb-md"
        >
          <FlexItem>
            <Title headingLevel="h3" size="lg" data-testid={`mcp-category-title-${title}`}>
              {title}
            </Title>
            {description && (
              <Content component="p" className="pf-v6-u-color-200 pf-v6-u-mt-sm">
                {description}
              </Content>
            )}
          </FlexItem>
          {hasMore && onShowAll && (
            <FlexItem>
              <Button
                variant="link"
                size="sm"
                isInline
                icon={<ArrowRightIcon />}
                iconPosition="right"
                data-testid={`mcp-show-all-${title.toLowerCase().replace(/\s+/g, '-')}`}
                onClick={onShowAll}
              >
                Show all {title}
              </Button>
            </FlexItem>
          )}
        </Flex>
        <Grid hasGutter>
          {displayItems.map((server) => (
            <GridItem key={String(server.id)} sm={12} md={6} lg={4} xl2={4}>
              <McpCatalogCard server={server} />
            </GridItem>
          ))}
        </Grid>
      </StackItem>
    );
  },
);
McpCatalogCategorySection.displayName = 'McpCatalogCategorySection';

export default McpCatalogCategorySection;
