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
import { isModelValidated } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import { mockPerformanceMetricsArtifacts } from '~/app/pages/modelCatalog/mocks/hardwareConfigurationMock';
import { mockAccuracyMetricsArtifacts } from '~/app/pages/modelCatalog/mocks/accuracyMetricsMock';
import EmptyModelCatalogState from '~/app/pages/modelCatalog/EmptyModelCatalogState';
import { getSourceFromSourceId } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import ModelCatalogAllModelsView from './ModelCatalogAllModelsView';

type ModelCatalogPageProps = {
  searchTerm: string;
  handleFilterReset?: () => void;
};

const ModelCatalogPage: React.FC<ModelCatalogPageProps> = ({ searchTerm, handleFilterReset }) => {
  const {
    selectedSourceLabel,
    filterData,
    filterOptions,
    filterOptionsLoaded,
    filterOptionsLoadError,
    catalogSources,
  } = React.useContext(ModelCatalogContext);
  const { catalogModels, catalogModelsLoaded, catalogModelsLoadError } = useCatalogModelsBySources(
    '',
    selectedSourceLabel,
    10,
    searchTerm,
    filterData,
    filterOptions,
  );

  const isAllModelsView = selectedSourceLabel === 'All models';

  if (isAllModelsView) {
    return <ModelCatalogAllModelsView searchTerm={searchTerm} />;
  }

  const loaded = catalogModelsLoaded && filterOptionsLoaded;
  const loadError = catalogModelsLoadError || filterOptionsLoadError;

  if (loadError) {
    return (
      <Alert variant="danger" title="Failed to load model catalog" isInline>
        {loadError.message}
      </Alert>
    );
  }

  if (!loaded) {
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
        description={<>Adjust your filters and try again</>}
        customAction={<Button onClick={handleFilterReset}>Reset filters</Button>}
      />
    );
  }

  return (
    <>
      <Gallery hasGutter minWidths={{ default: '300px' }}>
        {catalogModels.items.map((model: CatalogModel) => (
          <ModelCatalogCard
            key={`${model.name}/${model.source_id}`}
            model={model}
            source={getSourceFromSourceId(model.source_id || '', catalogSources)}
            performanceMetrics={
              isModelValidated(model) ? mockPerformanceMetricsArtifacts : undefined
            }
            accuracyMetrics={isModelValidated(model) ? mockAccuracyMetricsArtifacts : undefined}
          />
        ))}
      </Gallery>
      {catalogModels.hasMore && (
        <Bullseye className="pf-v6-u-mt-lg">
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
      )}
    </>
  );
};

export default ModelCatalogPage;
