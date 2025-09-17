import * as React from 'react';
import { PageSection } from '@patternfly/react-core';
import { ApplicationsPage, ProjectObjectType, TitleWithIcon } from 'mod-arch-shared';
import ScrollViewOnMount from '~/app/shared/components/ScrollViewOnMount';
import { modelCatalogUrl } from '~/app/routes/modelCatalog/catalogModel';
import ModelCatalogPage from './ModelCatalogPage';
import ModelCatalogSourceSelectorNavigator from './ModelCatalogSourceSelectorNavigator';

const ModelCatalog: React.FC = () => {
  const [searchTerm, setSearchTerm] = React.useState('');

  const handleSearch = React.useCallback((term: string) => {
    setSearchTerm(term);
  }, []);

  const handleClearSearch = React.useCallback(() => {
    setSearchTerm('');
  }, []);

  return (
    <>
      <ScrollViewOnMount shouldScroll />
      <ApplicationsPage
        title={<TitleWithIcon title="Model Catalog" objectType={ProjectObjectType.modelCatalog} />}
        description="Discover models that are available for your organization to register, deploy, and customize."
        empty={false}
        headerContent={
          <ModelCatalogSourceSelectorNavigator
            getRedirectPath={(sourceId: string) => modelCatalogUrl(sourceId)}
            searchTerm={searchTerm}
            onSearch={handleSearch}
            onClearSearch={handleClearSearch}
          />
        }
        loaded
        provideChildrenPadding
      >
        <PageSection isFilled>
          <ModelCatalogPage searchTerm={searchTerm} />
        </PageSection>
      </ApplicationsPage>
    </>
  );
};

export default ModelCatalog;
