import * as React from 'react';
import {
  PageSection,
  Title,
  Gallery,
  EmptyState,
  Button,
  Spinner,
  Alert,
  Bullseye,
} from '@patternfly/react-core';
import { ApplicationsPage, ProjectObjectType, TitleWithIcon } from 'mod-arch-shared';
import ModelCatalogCard from '~/app/pages/modelCatalog/components/ModelCatalogCard';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import ScrollViewOnMount from '~/app/shared/components/ScrollViewOnMount';
import EmptyModelCatalogState from '~/app/pages/modelCatalog/EmptyModelCatalogState';
import { CatalogModel } from '~/app/modelCatalogTypes';
import { useCatalogModelsBySource } from '~/app/hooks/modelCatalog/useCatalogModelsBySource';

const ModelCatalogPage: React.FC = () => {
  const { selectedSource } = React.useContext(ModelCatalogContext);
  const [catalogModels, catalogModelsLoaded, catalogModelsLoadError, refresh] =
    useCatalogModelsBySource(selectedSource?.id || '', 10);

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
        title={<TitleWithIcon title="Model Catalog" objectType={ProjectObjectType.modelCatalog} />}
        description="Discover models that are available for your organization to register, deploy, and customize."
        empty={catalogModels.items.length === 0}
        emptyStatePage={
          <EmptyModelCatalogState
            testid="empty-model-catalog-state"
            title="No models available"
            description="There are no models available in the catalog. Try refreshing or to request access to model catalog, contact your administrator."
          >
            <Button variant="primary" onClick={refresh}>
              Refresh
            </Button>
          </EmptyModelCatalogState>
        }
        headerContent={null}
        loaded
        errorMessage="Unable to load model catalog"
        provideChildrenPadding
      >
        <PageSection isFilled>
          <Gallery hasGutter minWidths={{ default: '300px' }}>
            {catalogModels.items.map((model: CatalogModel) => (
              <ModelCatalogCard
                model={model}
                source={selectedSource}
                key={`${model.name}/${model.source_id}`}
              />
            ))}
          </Gallery>
          {catalogModels.hasMore && (
            <div style={{ marginTop: '2rem' }}>
              <Bullseye>
                {catalogModels.isLoadingMore ? (
                  <>
                    <Spinner size="lg" className="pf-v5-u-mb-md" />
                    <Title size="lg" headingLevel="h5">
                      Loading more catalog models...
                    </Title>
                  </>
                ) : (
                  <Button variant="tertiary" onClick={catalogModels.loadMore} size="lg">
                    Load more models
                  </Button>
                )}
              </Bullseye>
            </div>
          )}
        </PageSection>
      </ApplicationsPage>
    </>
  );
};

export default ModelCatalogPage;
