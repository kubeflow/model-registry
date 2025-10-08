import {
  Alert,
  Bullseye,
  Button,
  EmptyState,
  Flex,
  Gallery,
  Spinner,
  Title,
} from '@patternfly/react-core';
import { SearchIcon } from '@patternfly/react-icons';
import React from 'react';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import { useCatalogModelsBySources } from '~/app/hooks/modelCatalog/useCatalogModelsBySource';
import { CatalogModel } from '~/app/modelCatalogTypes';
import ModelCatalogCard from '~/app/pages/modelCatalog/components/ModelCatalogCard';
import ValidatedModelCard from '~/app/pages/modelCatalog/components/ValidatedModelCard';
import { isModelValidated } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import { mockPerformanceMetricsArtifacts } from '~/app/pages/modelCatalog/mocks/hardwareConfigurationMock';
import { mockAccuracyMetricsArtifacts } from '~/app/pages/modelCatalog/mocks/accuracyMetricsMock';
import EmptyModelCatalogState from '~/app/pages/modelCatalog/EmptyModelCatalogState';

type ModelCatalogPageProps = {
  searchTerm: string;
};

const ModelCatalogPage: React.FC<ModelCatalogPageProps> = ({ searchTerm }) => {
  const { selectedSource } = React.useContext(ModelCatalogContext);
  const { catalogModels, catalogModelsLoaded, catalogModelsLoadError } = useCatalogModelsBySources(
    selectedSource?.id || '',
    10,
    searchTerm,
  );

  if (catalogModelsLoadError) {
    return (
      <Alert variant="danger" title="Failed to load model catalog" isInline>
        {catalogModelsLoadError.message}
      </Alert>
    );
  }

  if (!catalogModelsLoaded) {
    return (
      <EmptyState>
        <Spinner />
        <Title headingLevel="h4" size="lg">
          Loading model catalog...
        </Title>
      </EmptyState>
    );
  }

  if (catalogModels.items.length === 0) {
    return (
      <EmptyModelCatalogState
        testid="empty-model-catalog-state"
        title="No result found"
        headerIcon={SearchIcon}
        description={
          <>
            No models from the <b>{selectedSource?.name}</b> source match the search criteria.
            Adjust your search, or select a different source
          </>
        }
      />
    );
  }

  return (
    <>
      <Gallery hasGutter minWidths={{ default: '300px' }}>
        {catalogModels.items.map((model: CatalogModel) => {
          // Show ValidatedModelCard for validated models, ModelCatalogCard for others. This will be take care of by Pushpa's PR
          // that will implement sections for validated and non-validated models.
          if (isModelValidated(model)) {
            return (
              <ValidatedModelCard
                key={`${model.name}/${model.source_id}`}
                model={model}
                source={selectedSource}
                performanceMetrics={mockPerformanceMetricsArtifacts}
                accuracyMetrics={mockAccuracyMetricsArtifacts}
              />
            );
          }
          return (
            <ModelCatalogCard
              model={model}
              source={selectedSource}
              key={`${model.name}/${model.source_id}`}
            />
          );
        })}
      </Gallery>
      {catalogModels.hasMore && (
        <Flex direction={{ default: 'column' }} gap={{ default: 'gapLg' }}>
          <Bullseye>
            {catalogModels.isLoadingMore ? (
              <Flex
                direction={{ default: 'column' }}
                alignItems={{ default: 'alignItemsCenter' }}
                gap={{ default: 'gapMd' }}
              >
                <Spinner size="lg" />
                <Title size="lg" headingLevel="h5">
                  Loading more catalog models...
                </Title>
              </Flex>
            ) : (
              <Button variant="tertiary" onClick={catalogModels.loadMore} size="lg">
                Load more models
              </Button>
            )}
          </Bullseye>
        </Flex>
      )}
    </>
  );
};

export default ModelCatalogPage;
