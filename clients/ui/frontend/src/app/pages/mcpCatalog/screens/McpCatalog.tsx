import * as React from 'react';
import { PageSection, Sidebar, SidebarContent, SidebarPanel, Stack } from '@patternfly/react-core';
import { ApplicationsPage } from 'mod-arch-shared';
import ScrollViewOnMount from '~/app/shared/components/ScrollViewOnMount';
import { McpCatalogContext } from '~/app/context/mcpCatalog/McpCatalogContext';
import McpCatalogFilters from '~/app/pages/mcpCatalog/components/McpCatalogFilters';
import McpCatalogSourceLabelSelector from './McpCatalogSourceLabelSelector';
import McpCatalogGalleryView from './McpCatalogGalleryView';

const MCP_CATALOG_TITLE = 'MCP Catalog';
const MCP_CATALOG_SUBTITLE =
  'Discover and manage MCP servers and tools available for your organization.';

const McpCatalog: React.FC = () => {
  const { searchQuery, setSearchQuery, clearAllFilters } = React.useContext(McpCatalogContext);

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
        title={MCP_CATALOG_TITLE}
        description={MCP_CATALOG_SUBTITLE}
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
                <McpCatalogGalleryView />
              </PageSection>
            </Stack>
          </SidebarContent>
        </Sidebar>
      </ApplicationsPage>
    </>
  );
};

export default McpCatalog;
