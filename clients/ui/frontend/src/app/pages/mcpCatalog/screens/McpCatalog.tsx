import * as React from 'react';
import {
  Button,
  Flex,
  PageSection,
  Sidebar,
  SidebarContent,
  SidebarPanel,
  Stack,
  StackItem,
  Toolbar,
  ToolbarContent,
  ToolbarItem,
  ToolbarToggleGroup,
  ToolbarGroup,
  Spinner,
  Bullseye,
  EmptyState,
  EmptyStateBody,
  EmptyStateVariant,
} from '@patternfly/react-core';
import { FilterIcon, ExclamationCircleIcon, ArrowRightIcon } from '@patternfly/react-icons';
import { ApplicationsPage, ProjectObjectType, TitleWithIcon } from 'mod-arch-shared';
import { useThemeContext } from 'mod-arch-kubeflow';
import ScrollViewOnMount from '~/app/shared/components/ScrollViewOnMount';
import McpCatalogFilters from '~/app/pages/mcpCatalog/components/McpCatalogFilters';
import McpCatalogSourceLabelBlocks from '~/app/pages/mcpCatalog/components/McpCatalogSourceLabelBlocks';
import { useMcpCatalog } from '~/app/context/mcpCatalog/McpCatalogContext';
import { McpCategoryName } from '~/app/pages/mcpCatalog/types';
import { hasMcpFiltersActive } from '~/app/pages/mcpCatalog/utils/mcpCatalogUtils';
import ThemeAwareSearchInput from '~/app/pages/modelRegistry/screens/components/ThemeAwareSearchInput';
import McpCatalogAllServersView from './McpCatalogAllServersView';
import McpCatalogGalleryView from './McpCatalogGalleryView';

const McpCatalog: React.FC = () => {
  const {
    mcpServersLoaded,
    mcpServersLoadError,
    selectedSourceLabel,
    filters,
    searchTerm,
    updateFilters,
    updateSearchTerm,
    resetFilters,
    filterOptions,
    filterOptionsLoaded,
  } = useMcpCatalog();

  const [inputValue, setInputValue] = React.useState('');
  const { isMUITheme } = useThemeContext();

  const isAllServersView = selectedSourceLabel === McpCategoryName.allServers && !searchTerm;

  const handleSearch = React.useCallback(
    (value: string) => {
      updateSearchTerm(value);
    },
    [updateSearchTerm],
  );

  const handleClearSearch = React.useCallback(() => {
    updateSearchTerm('');
    setInputValue('');
  }, [updateSearchTerm]);

  const handleSearchInputChange = React.useCallback((value: string) => {
    setInputValue(value);
  }, []);

  const handleSearchInputSearch = React.useCallback(
    (_: React.SyntheticEvent<HTMLButtonElement>, value: string) => {
      handleSearch(value.trim());
    },
    [handleSearch],
  );

  const handleSearchButtonClick = React.useCallback(() => {
    if (inputValue.trim() !== searchTerm) {
      handleSearch(inputValue.trim());
    }
  }, [inputValue, searchTerm, handleSearch]);

  // Handle filter changes - update context filters which triggers server-side refetch
  const handleProviderChange = React.useCallback(
    (provider: string, checked: boolean) => {
      const newProviders = checked
        ? [...filters.selectedProviders, provider]
        : filters.selectedProviders.filter((p) => p !== provider);
      updateFilters({ ...filters, selectedProviders: newProviders });
    },
    [filters, updateFilters],
  );

  const handleLicenseChange = React.useCallback(
    (license: string, checked: boolean) => {
      const newLicenses = checked
        ? [...filters.selectedLicenses, license]
        : filters.selectedLicenses.filter((l) => l !== license);
      updateFilters({ ...filters, selectedLicenses: newLicenses });
    },
    [filters, updateFilters],
  );

  const handleTagChange = React.useCallback(
    (tag: string, checked: boolean) => {
      const newTags = checked
        ? [...filters.selectedTags, tag]
        : filters.selectedTags.filter((t) => t !== tag);
      updateFilters({ ...filters, selectedTags: newTags });
    },
    [filters, updateFilters],
  );

  const handleTransportChange = React.useCallback(
    (transport: string, checked: boolean) => {
      const newTransports = checked
        ? [...filters.selectedTransports, transport]
        : filters.selectedTransports.filter((t) => t !== transport);
      updateFilters({ ...filters, selectedTransports: newTransports });
    },
    [filters, updateFilters],
  );

  const handleDeploymentModeChange = React.useCallback(
    (mode: string, checked: boolean) => {
      const newModes = checked
        ? [...filters.selectedDeploymentModes, mode]
        : filters.selectedDeploymentModes.filter((m) => m !== mode);
      updateFilters({ ...filters, selectedDeploymentModes: newModes });
    },
    [filters, updateFilters],
  );

  // Check if any filters are active
  const hasActiveFilters = searchTerm.length > 0 || hasMcpFiltersActive(filters);

  const handleResetAllFilters = React.useCallback(() => {
    resetFilters();
    setInputValue('');
  }, [resetFilters]);

  // Show loading state (wait for both servers and filter options to load)
  if (!mcpServersLoaded || !filterOptionsLoaded) {
    return (
      <ApplicationsPage
        title={<TitleWithIcon title="MCP catalog" objectType={ProjectObjectType.modelCatalog} />}
        description="Loading MCP servers..."
        empty={false}
        loaded={false}
        provideChildrenPadding
      >
        <Bullseye>
          <Spinner size="xl" />
        </Bullseye>
      </ApplicationsPage>
    );
  }

  // Show error state
  if (mcpServersLoadError) {
    return (
      <ApplicationsPage
        title={<TitleWithIcon title="MCP catalog" objectType={ProjectObjectType.modelCatalog} />}
        description="Error loading MCP servers"
        empty={false}
        loaded
        provideChildrenPadding
      >
        <Bullseye>
          <EmptyState
            headingLevel="h4"
            icon={ExclamationCircleIcon}
            titleText="Error loading MCP servers"
            variant={EmptyStateVariant.lg}
          >
            <EmptyStateBody>
              {mcpServersLoadError.message ||
                'An unexpected error occurred while loading MCP servers.'}
            </EmptyStateBody>
          </EmptyState>
        </Bullseye>
      </ApplicationsPage>
    );
  }

  return (
    <>
      <ScrollViewOnMount shouldScroll scrollToTop />
      <ApplicationsPage
        title={<TitleWithIcon title="MCP catalog" objectType={ProjectObjectType.modelCatalog} />}
        description="Discover MCP servers that are available for your organization to integrate with AI agents."
        empty={false}
        loaded
        provideChildrenPadding
      >
        <Sidebar hasBorder hasGutter>
          <SidebarPanel>
            <McpCatalogFilters
              allProviders={filterOptions.providers}
              allLicenses={filterOptions.licenses}
              allTags={filterOptions.tags}
              allTransports={filterOptions.transports}
              allDeploymentModes={filterOptions.deploymentModes}
              selectedProviders={filters.selectedProviders}
              selectedLicenses={filters.selectedLicenses}
              selectedTags={filters.selectedTags}
              selectedTransports={filters.selectedTransports}
              selectedDeploymentModes={filters.selectedDeploymentModes}
              onProviderChange={handleProviderChange}
              onLicenseChange={handleLicenseChange}
              onTagChange={handleTagChange}
              onTransportChange={handleTransportChange}
              onDeploymentModeChange={handleDeploymentModeChange}
            />
          </SidebarPanel>
          <SidebarContent>
            <Stack>
              <StackItem>
                <Toolbar
                  className="pf-v6-u-pb-0"
                  {...(hasActiveFilters
                    ? {
                        clearAllFilters: handleResetAllFilters,
                        clearFiltersButtonText: 'Reset all filters',
                      }
                    : {})}
                >
                  <ToolbarContent>
                    <Flex>
                      <ToolbarToggleGroup breakpoint="md" toggleIcon={<FilterIcon />}>
                        <ToolbarGroup
                          variant="filter-group"
                          gap={{ default: 'gapMd' }}
                          alignItems="center"
                        >
                          <ToolbarItem>
                            <ThemeAwareSearchInput
                              data-testid="mcp-search-input"
                              fieldLabel="Filter by name, description and provider"
                              aria-label="Search MCP servers"
                              className="toolbar-fieldset-wrapper"
                              placeholder="Filter by name, description and provider"
                              value={inputValue}
                              style={{ minWidth: '400px' }}
                              onChange={handleSearchInputChange}
                              onSearch={handleSearchInputSearch}
                              onClear={handleClearSearch}
                            />
                          </ToolbarItem>
                          <ToolbarItem>
                            {isMUITheme && (
                              <Button
                                isInline
                                aria-label="Search"
                                data-testid="mcp-search-button"
                                variant="link"
                                icon={<ArrowRightIcon />}
                                iconPosition="end"
                                onClick={handleSearchButtonClick}
                              />
                            )}
                          </ToolbarItem>
                        </ToolbarGroup>
                      </ToolbarToggleGroup>
                    </Flex>
                  </ToolbarContent>
                </Toolbar>
              </StackItem>
              <StackItem>
                <McpCatalogSourceLabelBlocks />
              </StackItem>
              <StackItem>
                <PageSection isFilled padding={{ default: 'noPadding' }}>
                  {isAllServersView ? <McpCatalogAllServersView /> : <McpCatalogGalleryView />}
                </PageSection>
              </StackItem>
            </Stack>
          </SidebarContent>
        </Sidebar>
      </ApplicationsPage>
    </>
  );
};

export default McpCatalog;
