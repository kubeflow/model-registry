import {
  Alert,
  Button,
  Flex,
  FlexItem,
  Grid,
  GridItem,
  Skeleton,
  StackItem,
  Title,
} from '@patternfly/react-core';
import React from 'react';
import { ArrowRightIcon, SearchIcon } from '@patternfly/react-icons';
import EmptyMcpCatalogState from '~/app/pages/mcpCatalog/EmptyMcpCatalogState';
import { McpCatalogSourceList, McpServer, McpSourceLabel } from '~/app/pages/mcpCatalog/types';
import McpCatalogCard from './McpCatalogCard';

type McpCatalogCategorySectionProps = {
  label: string;
  pageSize: number;
  servers: McpServer[];
  sources: McpCatalogSourceList | null;
  loaded: boolean;
  loadError?: Error;
  onShowMore: (label: string) => void;
  displayName?: string;
};

/**
 * Category section for the All Servers view.
 * Groups servers by source label. Server-side filtering is already applied via context.
 * Only source label filtering is done client-side.
 */
const McpCatalogCategorySection: React.FC<McpCatalogCategorySectionProps> = ({
  label,
  pageSize,
  servers,
  sources,
  loaded,
  loadError,
  onShowMore,
  displayName,
}) => {
  // Filter servers by source label (client-side - source labels are UI grouping only)
  const filteredServers = React.useMemo(() => {
    if (!sources?.items) {
      return [];
    }

    // Find sources matching this label
    // For McpSourceLabel.other ("null"), find sources with empty labels (catch-all category)
    const isOtherCategory = label === McpSourceLabel.other;
    const matchingSources = sources.items.filter((source) => {
      if (isOtherCategory) {
        // Match sources with no labels or only empty/whitespace labels
        return source.labels.length === 0 || source.labels.every((l) => !l.trim());
      }
      return source.labels.includes(label);
    });
    const matchingSourceIds = matchingSources.map((s) => s.id);

    // Filter servers by source_id
    return servers.filter((server) => {
      if (!server.source_id) {
        return false;
      }
      return matchingSourceIds.includes(server.source_id);
    });
  }, [servers, sources, label]);

  const itemsToDisplay = filteredServers.slice(0, pageSize);

  return (
    <StackItem className="pf-v6-u-pb-xl">
      <Flex
        alignItems={{ default: 'alignItemsCenter' }}
        justifyContent={{ default: 'justifyContentSpaceBetween' }}
        className="pf-v6-u-mb-md"
      >
        <FlexItem>
          <Title headingLevel="h3" size="lg" data-testid={`title ${label}`}>
            {`${displayName ?? label} servers`}
          </Title>
        </FlexItem>

        {filteredServers.length >= pageSize && (
          <FlexItem>
            <Button
              variant="link"
              size="sm"
              isInline
              icon={<ArrowRightIcon />}
              iconPosition="end"
              data-testid={`show-more-button ${label.toLowerCase().replace(/\s+/g, '-')}`}
              onClick={() => onShowMore(label)}
            >
              Show all {displayName ?? label} servers
            </Button>
          </FlexItem>
        )}
      </Flex>

      {loadError ? (
        <Alert
          variant="danger"
          title={`Failed to load ${displayName ?? label} servers`}
          data-testid={`error-state ${label}`}
        >
          {loadError.message}
        </Alert>
      ) : !loaded ? (
        <Grid hasGutter>
          {Array.from({ length: 4 }).map((_, index) => (
            <GridItem key={index} sm={6} md={6} lg={6} xl={6} xl2={3}>
              <Skeleton
                height="280px"
                width="100%"
                screenreaderText={`Loading ${label} servers`}
                data-testid={`category-skeleton-${label.toLowerCase().replace(/\s+/g, '-')}-${index}`}
              />
            </GridItem>
          ))}
        </Grid>
      ) : filteredServers.length === 0 ? (
        <EmptyMcpCatalogState
          testid={`empty-mcp-catalog-state ${label}`}
          title="No result found"
          headerIcon={SearchIcon}
          description="Adjust your filters and try again."
        />
      ) : (
        <Grid hasGutter>
          {itemsToDisplay.map((server) => (
            <GridItem
              key={`${server.name}/${server.source_id}`}
              sm={6}
              md={6}
              lg={6}
              xl={6}
              xl2={3}
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
