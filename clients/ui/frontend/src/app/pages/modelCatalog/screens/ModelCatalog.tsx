import * as React from 'react';
import { Divider, Flex, FlexItem, PageSection } from '@patternfly/react-core';
import { ApplicationsPage, ProjectObjectType, TitleWithIcon } from 'mod-arch-shared';
import ScrollViewOnMount from '~/app/shared/components/ScrollViewOnMount';
import { modelCatalogUrl } from '~/app/routes/modelCatalog/catalogModel';
import ModelCatalogPage from './ModelCatalogPage';
import ModelCatalogSourceSelectorNavigator from './ModelCatalogSourceSelectorNavigator';
import ModelCatalogFilters from '../components/ModelCatalogFilters';

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
        loaded
        provideChildrenPadding
      >
        <Flex flexWrap={{ default: 'nowrap' }}>
          <FlexItem style={{ minWidth: '280px' }}>
            <ModelCatalogFilters />
          </FlexItem>
          <Divider orientation={{ default: 'vertical' }} />
          <FlexItem>
            <ModelCatalogSourceSelectorNavigator
              getRedirectPath={(sourceId: string) => modelCatalogUrl(sourceId)}
              searchTerm={searchTerm}
              onSearch={handleSearch}
              onClearSearch={handleClearSearch}
            />
            <PageSection isFilled>
              <ModelCatalogPage searchTerm={searchTerm} />
            </PageSection>
          </FlexItem>
        </Flex>
      </ApplicationsPage>
    </>
  );
};

export default ModelCatalog;
