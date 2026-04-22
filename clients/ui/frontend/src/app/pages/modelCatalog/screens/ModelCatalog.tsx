import * as React from 'react';
import { PageSection, Sidebar, SidebarContent, SidebarPanel } from '@patternfly/react-core';
import { ApplicationsPage, ProjectObjectType, TitleWithIcon } from 'mod-arch-shared';
import { SearchIcon } from '@patternfly/react-icons';
import ScrollViewOnMount from '~/app/shared/components/ScrollViewOnMount';
import ModelCatalogFilters from '~/app/pages/modelCatalog/components/ModelCatalogFilters';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import { CategoryName } from '~/app/modelCatalogTypes';
import { useHasVisibleFiltersApplied } from '~/app/hooks/modelCatalog/useHasVisibleFiltersApplied';
import useEffectiveCategories from '~/app/hooks/useEffectiveCategories';
import EmptyModelCatalogState from '~/app/pages/modelCatalog/EmptyModelCatalogState';
import ModelCatalogSourceLabelSelectorNavigator from './ModelCatalogSourceLabelSelectorNavigator';
import ModelCatalogAllModelsView from './ModelCatalogAllModelsView';
import ModelCatalogGalleryView from './ModelCatalogGalleryView';

const ModelCatalog: React.FC = () => {
  const [searchTerm, setSearchTerm] = React.useState('');
  const {
    selectedSourceLabel,
    updateSelectedSourceLabel,
    clearAllFilters,
    catalogSources,
    catalogLabels,
    catalogSourcesLoaded,
    emptyCategoryLabels,
  } = React.useContext(ModelCatalogContext);
  const filtersApplied = useHasVisibleFiltersApplied();

  const { effectiveActiveCategories, isSingleCategory, hasNoCategories } = useEffectiveCategories(
    catalogSources,
    catalogLabels,
    emptyCategoryLabels,
    catalogSourcesLoaded,
    updateSelectedSourceLabel,
  );

  const isAllModelsView =
    selectedSourceLabel === CategoryName.allModels && !searchTerm && !filtersApplied;

  const handleSearch = React.useCallback((term: string) => {
    setSearchTerm(term);
  }, []);

  const handleClearSearch = React.useCallback(() => {
    setSearchTerm('');
  }, []);

  const handleFilterReset = React.useCallback(() => {
    setSearchTerm('');
    // clearAllFilters clears basic filters to empty and resets performance filters to defaults
    clearAllFilters();
  }, [clearAllFilters]);

  return (
    <>
      <ScrollViewOnMount shouldScroll scrollToTop />
      <ApplicationsPage
        title={<TitleWithIcon title="Model Catalog" objectType={ProjectObjectType.modelCatalog} />}
        description="Discover models that are available for your organization to register, deploy, and customize."
        empty={false}
        loaded
        provideChildrenPadding
      >
        {catalogSourcesLoaded && hasNoCategories ? (
          <EmptyModelCatalogState
            testid="empty-model-catalog-no-categories"
            title="No models available"
            headerIcon={SearchIcon}
            description="There are no model categories available. Configure model sources in settings to get started."
          />
        ) : (
          <Sidebar hasBorder hasGutter>
            <SidebarPanel variant="sticky">
              <ModelCatalogFilters />
            </SidebarPanel>
            <SidebarContent>
              <ModelCatalogSourceLabelSelectorNavigator
                searchTerm={searchTerm}
                onSearch={handleSearch}
                onClearSearch={handleClearSearch}
                onResetAllFilters={handleFilterReset}
              />
              <PageSection isFilled padding={{ default: 'noPadding' }}>
                {isAllModelsView && !isSingleCategory ? (
                  <ModelCatalogAllModelsView searchTerm={searchTerm} />
                ) : (
                  <ModelCatalogGalleryView
                    searchTerm={searchTerm}
                    handleFilterReset={handleFilterReset}
                    isSingleCategory={isSingleCategory}
                    singleCategoryLabel={
                      isSingleCategory ? effectiveActiveCategories[0] : undefined
                    }
                  />
                )}
              </PageSection>
            </SidebarContent>
          </Sidebar>
        )}
      </ApplicationsPage>
    </>
  );
};

export default ModelCatalog;
