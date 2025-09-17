import {
  Alert,
  Bullseye,
  Button,
  EmptyState,
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
import EmptyModelCatalogState from '~/app/pages/modelCatalog/EmptyModelCatalogState';

type ModelCatalogPageProps = {
  searchTerm: string;
};

const ModelCatalogPage: React.FC<ModelCatalogPageProps> = ({ searchTerm }) => {
  const { selectedSource } = React.useContext(ModelCatalogContext);
  const [catalogModels, catalogModelsLoaded, catalogModelsLoadError, refresh] =
    useCatalogModelsBySources(selectedSource?.id || '', 10, searchTerm);

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

  if (catalogModelsLoadError) {
    return (
      <Alert variant="danger" title="Failed to load model catalog" isInline>
        {catalogModelsLoadError.message}
        <Button variant="link" onClick={refresh}>
          Try again
        </Button>
      </Alert>
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
            Adjust your seach, or select a differenct source
          </>
        }
      />
    );
  }

  return (
    <>
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
    </>
  );
};

export default ModelCatalogPage;
