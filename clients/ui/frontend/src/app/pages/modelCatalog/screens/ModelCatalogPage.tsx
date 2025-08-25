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
import { useNavigate } from 'react-router-dom';
import { useModelCatalogSources } from '~/app/hooks/modelCatalog/useModelCatalogSources';
import { ModelCatalogItem } from '~/app/modelCatalogTypes';
import ModelCatalogCard from '~/app/pages/modelCatalog/components/ModelCatalogCard';

const ModelCatalogPage: React.FC = () => {
  const { sources, loading, error, refreshSources } = useModelCatalogSources();
  const navigate = useNavigate();

  const handleModelSelect = (model: ModelCatalogItem) => {
    if (model.id) {
      navigate(`/model-catalog/${encodeURIComponent(model.id)}`);
    }
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

  if (sources.length === 0) {
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
      <Title headingLevel="h1" size="2xl" className="pf-v5-u-mb-md">
        Model Catalog
      </Title>
      <p className="pf-v5-u-mb-lg">
        Discover models that are available for your organization to register, deploy, and customize.
      </p>
      <Gallery hasGutter minWidths={{ default: '300px' }}>
        {sources.map((source) =>
          (source.models || []).map((model) => (
            <GalleryItem key={model.id}>
              <ModelCatalogCard
                model={model}
                source={source.displayName}
                onSelect={handleModelSelect}
              />
            </GalleryItem>
          )),
        )}
      </Gallery>
    </PageSection>
  );
};

export default ModelCatalogPage;
