import * as React from 'react';
import { PageSection, Sidebar, SidebarContent, SidebarPanel } from '@patternfly/react-core';
import { ApplicationsPage, ProjectObjectType, TitleWithIcon } from 'mod-arch-shared';
import ScrollViewOnMount from '~/app/shared/components/ScrollViewOnMount';
import ModelCatalogFilters from '~/app/pages/modelCatalog/components/ModelCatalogFilters';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import ModelCatalogPage from './ModelCatalogPage';
import ModelCatalogSourceLabelSelectorNavigator from './ModelCatalogSourceLabelSelectorNavigator';
import ModelCatalogAllModelsView from './ModelCatalogAllModelsView';

const ModelCatalog: React.FC = () => {
  const [searchTerm, setSearchTerm] = React.useState('');
  const { selectedSourceLabel } = React.useContext(ModelCatalogContext);
  const isAllModelsView = selectedSourceLabel === 'All models' && !searchTerm;

  const handleSearch = React.useCallback((term: string) => {
    setSearchTerm(term);
  }, []);

  const handleClearSearch = React.useCallback(() => {
    setSearchTerm('');
  }, []);

  const handleFilterReset = React.useCallback(() => {
    setSearchTerm('');
  }, []);

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
            <PageSection isFilled style={{ paddingLeft: '0px', paddingTop: '25px' }}>
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
