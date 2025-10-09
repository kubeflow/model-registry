import * as React from 'react';
import { PageSection, Sidebar, SidebarContent, SidebarPanel } from '@patternfly/react-core';
import { ApplicationsPage, ProjectObjectType, TitleWithIcon } from 'mod-arch-shared';
import ScrollViewOnMount from '~/app/shared/components/ScrollViewOnMount';
import ModelCatalogFilters from '~/app/pages/modelCatalog/components/ModelCatalogFilters';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import { hasFiltersApplied } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import { ModelCatalogStringFilterKey } from '~/concepts/modelCatalog/const';
import ModelCatalogPage from './ModelCatalogPage';
import ModelCatalogSourceLabelSelectorNavigator from './ModelCatalogSourceLabelSelectorNavigator';
import ModelCatalogAllModelsView from './ModelCatalogAllModelsView';

const ModelCatalog: React.FC = () => {
  const [searchTerm, setSearchTerm] = React.useState('');
  const { selectedSourceLabel, filterData, setFilterData } = React.useContext(ModelCatalogContext);
  const filtersApplied = hasFiltersApplied(filterData);
  const isAllModelsView = selectedSourceLabel === 'All models' && !searchTerm && !filtersApplied;

  const handleSearch = React.useCallback((term: string) => {
    setSearchTerm(term);
  }, []);

  const handleClearSearch = React.useCallback(() => {
    setSearchTerm('');
  }, []);

  const resetAllFilters = React.useCallback(() => {
    setFilterData(ModelCatalogStringFilterKey.TASK, []);
    setFilterData(ModelCatalogStringFilterKey.PROVIDER, []);
    setFilterData(ModelCatalogStringFilterKey.LICENSE, []);
    setFilterData(ModelCatalogStringFilterKey.LANGUAGE, []);
  }, [setFilterData]);

  const handleFilterReset = React.useCallback(() => {
    setSearchTerm('');
    resetAllFilters();
  }, [resetAllFilters]);

  return (
    <>
      <ScrollViewOnMount shouldScroll />
      <ApplicationsPage
        title={<TitleWithIcon title="Model Catalog" objectType={ProjectObjectType.modelCatalog} />}
        description="Discover models that are available for your organization to register, deploy, and customize."
        empty={false}
        loaded
        provideChildrenPadding
      >
        <Sidebar hasBorder hasGutter>
          <SidebarPanel>
            <ModelCatalogFilters />
          </SidebarPanel>
          <SidebarContent>
            <ModelCatalogSourceLabelSelectorNavigator
              searchTerm={searchTerm}
              onSearch={handleSearch}
              onClearSearch={handleClearSearch}
            />
            <PageSection isFilled padding={{ default: 'noPadding' }}>
              {isAllModelsView ? (
                <ModelCatalogAllModelsView searchTerm={searchTerm} />
              ) : (
                <ModelCatalogPage searchTerm={searchTerm} handleFilterReset={handleFilterReset} />
              )}
            </PageSection>
          </SidebarContent>
        </Sidebar>
      </ApplicationsPage>
    </>
  );
};

export default ModelCatalog;
