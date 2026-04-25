import * as React from 'react';
import {
  Alert,
  Button,
  Content,
  Flex,
  FlexItem,
  Grid,
  GridItem,
  Skeleton,
  StackItem,
  Title,
} from '@patternfly/react-core';
import { ArrowRightIcon, SearchIcon } from '@patternfly/react-icons';
import { useMcpServersBySourceLabelWithAPI } from '~/app/hooks/mcpServerCatalog/useMcpServersBySourceLabel';
import useReportCategoryEmpty from '~/app/hooks/useReportCategoryEmpty';
import {
  getLabelDescription,
  getLabelDisplayName,
} from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import EmptyModelCatalogState from '~/app/pages/modelCatalog/EmptyModelCatalogState';
import { McpCatalogContext } from '~/app/context/mcpCatalog/McpCatalogContext';
import {
  MCP_CATALOG_GRID_SPAN,
  OTHER_MCP_SERVERS_DISPLAY_NAME,
} from '~/app/pages/mcpCatalog/const';
import McpCatalogCard from '~/app/pages/mcpCatalog/components/McpCatalogCard';

type McpCatalogCategorySectionProps = {
  label: string;
  searchTerm: string;
  pageSize: number;
  onShowMore: (label: string) => void;
};

const McpCatalogCategorySection: React.FC<McpCatalogCategorySectionProps> = ({
  label,
  searchTerm,
  pageSize,
  onShowMore,
}) => {
  const { mcpApiState, catalogLabels, reportCategoryEmpty } = React.useContext(McpCatalogContext);
  const { mcpServers, mcpServersLoaded, mcpServersLoadError } = useMcpServersBySourceLabelWithAPI(
    mcpApiState,
    {
      sourceLabel: label,
      pageSize,
      searchQuery: searchTerm,
    },
  );

  const itemsToDisplay = mcpServers.items.slice(0, pageSize);

  const categoryTitle = getLabelDisplayName(
    label,
    catalogLabels,
    OTHER_MCP_SERVERS_DISPLAY_NAME,
    'servers',
  );
  const description = getLabelDescription(label, catalogLabels);

  useReportCategoryEmpty(
    reportCategoryEmpty,
    label,
    mcpServersLoaded,
    mcpServers.items.length,
    searchTerm,
  );

  if (mcpServersLoaded && mcpServers.items.length === 0 && !searchTerm) {
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
          <Title headingLevel="h3" size="lg" data-testid={`mcp-category-title-${label}`}>
            {categoryTitle}
          </Title>
          {description && (
            <Content component="p" className="pf-v6-u-color-200 pf-v6-u-mt-sm">
              {description}
            </Content>
          )}
        </FlexItem>

        {mcpServers.items.length >= pageSize && (
          <FlexItem>
            <Button
              variant="link"
              size="sm"
              isInline
              icon={<ArrowRightIcon />}
              iconPosition="right"
              data-testid={`mcp-show-all-${label.toLowerCase().replace(/\s+/g, '-')}`}
              onClick={() => onShowMore(label)}
            >
              Show all {categoryTitle}
            </Button>
          </FlexItem>
        )}
      </Flex>

      {mcpServersLoadError ? (
        <Alert
          variant="danger"
          title={`Failed to load ${categoryTitle}`}
          data-testid={`mcp-error-state-${label}`}
        >
          {mcpServersLoadError.message}
        </Alert>
      ) : !mcpServersLoaded ? (
        <Grid hasGutter>
          {Array.from({ length: pageSize }).map((_, index) => (
            <GridItem
              key={index}
              sm={MCP_CATALOG_GRID_SPAN.sm}
              md={MCP_CATALOG_GRID_SPAN.md}
              lg={MCP_CATALOG_GRID_SPAN.lg}
              xl2={MCP_CATALOG_GRID_SPAN.xl2}
            >
              <Skeleton
                height="280px"
                width="100%"
                screenreaderText={`Loading ${label} servers`}
                data-testid={`mcp-category-skeleton-${label.toLowerCase().replace(/\s+/g, '-')}-${index}`}
              />
            </GridItem>
          ))}
        </Grid>
      ) : mcpServers.items.length === 0 ? (
        <EmptyModelCatalogState
          testid={`empty-mcp-catalog-state-${label}`}
          title="No result found"
          headerIcon={SearchIcon}
          description="Adjust your filters and try again."
        />
      ) : (
        <Grid hasGutter>
          {itemsToDisplay.map((server) => (
            <GridItem
              key={server.id}
              sm={MCP_CATALOG_GRID_SPAN.sm}
              md={MCP_CATALOG_GRID_SPAN.md}
              lg={MCP_CATALOG_GRID_SPAN.lg}
              xl2={MCP_CATALOG_GRID_SPAN.xl2}
            >
              <McpCatalogCard server={server} />
            </GridItem>
          ))}
        </Grid>
      )}
    </StackItem>
  );
};
export default McpCatalogCategorySection;
