import * as React from 'react';
import { PageSection, Sidebar, SidebarContent, SidebarPanel, Stack } from '@patternfly/react-core';
import { ApplicationsPage, ProjectObjectType, TitleWithIcon } from 'mod-arch-shared';
import { SearchIcon } from '@patternfly/react-icons';
import ScrollViewOnMount from '~/app/shared/components/ScrollViewOnMount';
import { McpCatalogContext } from '~/app/context/mcpCatalog/McpCatalogContext';
import { hasMcpFiltersApplied } from '~/app/pages/mcpCatalog/utils/mcpCatalogUtils';
import McpCatalogFilters from '~/app/pages/mcpCatalog/components/McpCatalogFilters';
import { MCP_CATALOG_TITLE, MCP_CATALOG_DESCRIPTION } from '~/app/pages/mcpCatalog/const';
import useEffectiveCategories from '~/app/hooks/useEffectiveCategories';
import EmptyModelCatalogState from '~/app/pages/modelCatalog/EmptyModelCatalogState';
import McpCatalogSourceLabelSelector from './McpCatalogSourceLabelSelector';
import McpCatalogAllServersView from './McpCatalogAllServersView';
import McpCatalogGalleryView from './McpCatalogGalleryView';

const McpCatalog: React.FC = () => {
  const {
    searchQuery,
    setSearchQuery,
    clearAllFilters,
    selectedSourceLabel,
    setSelectedSourceLabel,
    filters,
    catalogSources,
    catalogLabels,
    catalogSourcesLoaded,
    emptyCategoryLabels,
  } = React.useContext(McpCatalogContext);

  const filtersApplied = hasMcpFiltersApplied(filters, searchQuery);
  const isAllServersView = selectedSourceLabel === undefined && !filtersApplied;

  const { effectiveActiveCategories, isSingleCategory, hasNoCategories } = useEffectiveCategories(
    catalogSources,
    catalogLabels,
    emptyCategoryLabels,
    catalogSourcesLoaded,
    setSelectedSourceLabel,
  );

  const handleSearch = React.useCallback(
    (term: string) => {
      setSearchQuery(term);
    },
    [setSearchQuery],
  );

  const handleClearSearch = React.useCallback(() => {
    setSearchQuery('');
  }, [setSearchQuery]);

  const handleResetAllFilters = React.useCallback(() => {
    clearAllFilters();
  }, [clearAllFilters]);

  return (
    <>
      <ScrollViewOnMount shouldScroll scrollToTop />
      <ApplicationsPage
        title={
          <TitleWithIcon title={MCP_CATALOG_TITLE} objectType={ProjectObjectType.mcpCatalog} />
        }
        description={MCP_CATALOG_DESCRIPTION}
        empty={false}
        loaded
        provideChildrenPadding
      >
        {catalogSourcesLoaded && hasNoCategories ? (
          <EmptyModelCatalogState
            testid="empty-mcp-catalog-no-categories"
            title="No MCP servers available"
            headerIcon={SearchIcon}
            description="There are no MCP server categories available. Configure sources in settings to get started."
          />
        ) : (
          <Sidebar hasBorder hasGutter>
            <SidebarPanel variant="sticky">
              <McpCatalogFilters />
            </SidebarPanel>
            <SidebarContent>
              <Stack hasGutter>
                <McpCatalogSourceLabelSelector
                  searchTerm={searchQuery}
                  onSearch={handleSearch}
                  onClearSearch={handleClearSearch}
                  onResetAllFilters={handleResetAllFilters}
                />
                <PageSection isFilled padding={{ default: 'noPadding' }}>
                  {isAllServersView && !isSingleCategory ? (
                    <McpCatalogAllServersView searchTerm={searchQuery} />
                  ) : (
                    <McpCatalogGalleryView
                      handleFilterReset={handleResetAllFilters}
                      isSingleCategory={isSingleCategory}
                      singleCategoryLabel={
                        isSingleCategory ? effectiveActiveCategories[0] : undefined
                      }
                    />
                  )}
                </PageSection>
              </Stack>
            </SidebarContent>
          </Sidebar>
        )}
      </ApplicationsPage>
    </>
  );
};

export default McpCatalog;
