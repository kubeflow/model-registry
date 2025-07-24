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
import { useModelCatalogSources } from '~/app/hooks/modelCatalog/useModelCatalogSources';
import { ModelCatalogItem } from '~/app/modelCatalogTypes';
import ModelCatalogCard from '~/app/pages/modelCatalog/components/ModelCatalogCard';

const ModelCatalogPage: React.FC = () => {
  const { sources, loading, error, refreshSources } = useModelCatalogSources();

  const handleModelSelect = (model: ModelCatalogItem) => {
    // TODO: Implement model selection logic
    // eslint-disable-next-line no-console
    console.log('Selected model:', model);
  };

  if (loading) {
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

  if (error) {
    return (
      <PageSection>
        <Alert variant="danger" title="Failed to load model catalog" isInline>
          {error.message}
          <Button variant="link" onClick={refreshSources}>
            Try again
          </Button>
        </Alert>
      </PageSection>
    );
  }

  const allModels = sources.flatMap((source) => source.models || []);

  if (allModels.length === 0) {
    return (
      <PageSection>
        <EmptyState>
          <CubesIcon />
          <Title headingLevel="h4" size="lg">
            No models available
          </Title>
          <EmptyStateBody>
            There are no models available in the catalog. Try refreshing or contact your
            administrator.
          </EmptyStateBody>
          <Button variant="primary" onClick={refreshSources}>
            Refresh
          </Button>
        </EmptyState>
      </PageSection>
    );
  }

  return (
    <PageSection>
      <Title headingLevel="h1" size="2xl" className="pf-v5-u-mb-lg">
        Model Catalog
      </Title>
      <Gallery hasGutter minWidths={{ default: '300px' }}>
        {allModels.map((model) => (
          <GalleryItem key={model.id}>
            <ModelCatalogCard model={model} onSelect={handleModelSelect} />
          </GalleryItem>
        ))}
      </Gallery>
    </PageSection>
  );
};

export default ModelCatalogPage;
