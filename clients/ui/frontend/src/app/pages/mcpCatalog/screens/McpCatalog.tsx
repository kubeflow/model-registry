import * as React from 'react';
import { PageSection, Sidebar, SidebarContent, SidebarPanel, Stack } from '@patternfly/react-core';
import { ApplicationsPage, ProjectObjectType, TitleWithIcon } from 'mod-arch-shared';
import ScrollViewOnMount from '~/app/shared/components/ScrollViewOnMount';
import { McpCatalogContext } from '~/app/context/mcpCatalog/McpCatalogContext';
import { hasMcpFiltersApplied } from '~/app/pages/mcpCatalog/utils/mcpCatalogUtils';
import McpCatalogFilters from '~/app/pages/mcpCatalog/components/McpCatalogFilters';
import { MCP_CATALOG_TITLE, MCP_CATALOG_DESCRIPTION } from '~/app/pages/mcpCatalog/const';
import McpCatalogSourceLabelSelector from './McpCatalogSourceLabelSelector';
import McpCatalogAllServersView from './McpCatalogAllServersView';
import McpCatalogGalleryView from './McpCatalogGalleryView';

const McpCatalog: React.FC = () => {
  const { searchQuery, setSearchQuery, clearAllFilters, selectedSourceLabel, filters } =
    React.useContext(McpCatalogContext);

  const filtersApplied = hasMcpFiltersApplied(filters, searchQuery);
  const isAllServersView = selectedSourceLabel === undefined && !filtersApplied;

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
        <Sidebar hasBorder hasGutter>
          <SidebarPanel>
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
                {isAllServersView ? (
                  <McpCatalogAllServersView searchTerm={searchQuery} />
                ) : (
                  <McpCatalogGalleryView handleFilterReset={handleResetAllFilters} />
                )}
              </PageSection>
            </Stack>
          </SidebarContent>
        </Sidebar>
      </ApplicationsPage>
    </>
  );
};

export default McpCatalog;
