import * as React from 'react';
import {
  PageSection,
  Title,
  Gallery,
  GalleryItem,
  EmptyState,
  EmptyStateBody,
  Button,
  Spinner,
  Alert,
} from '@patternfly/react-core';
import { CubesIcon } from '@patternfly/react-icons';
import { useCatalogModelsbySources } from '~/app/hooks/modelCatalog/useCatalogModelsbySources';
import ModelCatalogCard from '~/app/pages/modelCatalog/components/ModelCatalogCard';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import ScrollViewOnMount from '~/app/shared/components/ScrollViewOnMount';
import { ApplicationsPage, ProjectObjectType, TitleWithIcon } from 'mod-arch-shared';
import EmptyModelCatalogState from '../EmptyModelCatalogState';

const ModelCatalogPage: React.FC = () => {
  const { selectedSource } = React.useContext(ModelCatalogContext);
  const [catalogModels, catalogModelsLoaded, catalogModelsLoadError, refresh] =
    useCatalogModelsbySources(selectedSource?.id || '');

  if (!catalogModelsLoaded) {
    return (
      <PageSection>
        <EmptyState>
          <Spinner />
          <Title headingLevel="h4" size="lg">
            Loading model catalog...
          </Title>
        </EmptyState>
      </PageSection>
    );
  }

  if (catalogModelsLoadError) {
    return (
      <PageSection>
        <Alert variant="danger" title="Failed to load model catalog" isInline>
          {catalogModelsLoadError.message}
          <Button variant="link" onClick={refresh}>
            Try again
          </Button>
        </Alert>
      </PageSection>
    );
  }

  return (
    <>
      <ScrollViewOnMount shouldScroll />
      <ApplicationsPage
        title={<TitleWithIcon title="Model catalog" objectType={ProjectObjectType.singleModel} />}
        description="Discover models that are available for your organization to register, deploy, and customize."
        empty={catalogModels.items.length === 0}
        emptyStatePage={
          <EmptyModelCatalogState
            testid="empty-model-catalog-state"
            title="No models available"
            description="There are no models available in the catalog. Try refreshing or to request access to model catalog, contact your administrator."
            children={
              <Button variant="primary" onClick={refresh}>
                Refresh
              </Button>
            }
          />
        }
        headerContent={null}
        loaded
        errorMessage="Unable to load model catalog"
        provideChildrenPadding
      >
        <PageSection isFilled>
          <Gallery hasGutter minWidths={{ default: '300px' }}>
            {catalogModels.items.map((model) => (
              <ModelCatalogCard model={model} source={selectedSource?.name || ''} />
            ))}
          </Gallery>
        </PageSection>
      </ApplicationsPage>
    </>
  );
};

export default ModelCatalogPage;
