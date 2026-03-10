import * as React from 'react';
import { Content, Flex, FlexItem, Grid, GridItem, StackItem, Title } from '@patternfly/react-core';
import type { McpServer } from '~/app/mcpServerCatalogTypes';
import McpCatalogCard from '~/app/pages/mcpCatalog/components/McpCatalogCard';

type McpCatalogCategorySectionProps = {
  title: string;
  description?: string;
  servers: McpServer[];
};

const McpCatalogCategorySection: React.FC<McpCatalogCategorySectionProps> = ({
  title,
  description,
  servers,
}) => {
  if (servers.length === 0) {
    return null;
  }

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
      </Flex>
      <Grid hasGutter>
        {servers.map((server) => (
          <GridItem key={String(server.id)} sm={12} md={6} lg={4} xl2={4}>
            <McpCatalogCard server={server} />
          </GridItem>
        ))}
      </Grid>
    </StackItem>
  );
};

export default McpCatalogCategorySection;
