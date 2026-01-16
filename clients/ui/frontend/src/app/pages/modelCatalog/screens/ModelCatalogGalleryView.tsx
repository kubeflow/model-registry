import {
  Alert,
  Bullseye,
  Button,
  EmptyState,
  EmptyStateVariant,
  Flex,
  Gallery,
  Spinner,
  Title,
} from '@patternfly/react-core';
import { ChartBarIcon, SearchIcon } from '@patternfly/react-icons';
import React from 'react';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import { useCatalogModelsBySources } from '~/app/hooks/modelCatalog/useCatalogModelsBySource';
import { CatalogModel, CategoryName, SourceLabel } from '~/app/modelCatalogTypes';
import ModelCatalogCard from '~/app/pages/modelCatalog/components/ModelCatalogCard';
import {
  getSourceFromSourceId,
  hasFiltersApplied,
  getBasicFiltersOnly,
  getActiveLatencyFieldName,
  getSortParams,
  generateCategoryName,
} from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import EmptyModelCatalogState from '~/app/pages/modelCatalog/EmptyModelCatalogState';
import ScrollViewOnMount from '~/app/shared/components/ScrollViewOnMount';
import {
  BASIC_FILTER_KEYS,
  ModelCatalogNumberFilterKey,
  parseLatencyFilterKey,
} from '~/concepts/modelCatalog/const';

type ModelCatalogPageProps = {
  searchTerm: string;
  handleFilterReset: () => void;
};

const ModelCatalogGalleryView: React.FC<ModelCatalogPageProps> = ({
  searchTerm,
  handleFilterReset,
}) => {
  const {
    selectedSourceLabel,
    filterData,
    filterOptions,
    filterOptionsLoaded,
    filterOptionsLoadError,
    catalogSources,
    setPerformanceViewEnabled,
    updateSelectedSourceLabel,
    performanceViewEnabled,
    sortBy,
  } = React.useContext(ModelCatalogContext);
  const filtersApplied = hasFiltersApplied(filterData);

  // When performance view is disabled, exclude performance filters from API queries
  // Memoize to prevent infinite re-fetching
  const effectiveFilterData = React.useMemo(
    () => (performanceViewEnabled ? filterData : getBasicFiltersOnly(filterData)),
    [performanceViewEnabled, filterData],
  );

  // Optimize: Only track the active latency field instead of entire filterData
  // This prevents unnecessary recalculations when non-latency filters change
  const activeLatencyField = React.useMemo(
    () => getActiveLatencyFieldName(filterData),
    [filterData],
  );

  const sortParams = React.useMemo(
    () => getSortParams(sortBy, performanceViewEnabled, activeLatencyField),
    [sortBy, performanceViewEnabled, activeLatencyField],
  );

  // Derive performance params to pass to the models API when performance view is enabled
  const performanceParams = React.useMemo(() => {
    if (!performanceViewEnabled) {
      return undefined;
    }

    const targetRPS = filterData[ModelCatalogNumberFilterKey.MAX_RPS];
    const latencyProperty = activeLatencyField
      ? parseLatencyFilterKey(activeLatencyField).propertyKey
      : undefined;

    return {
      targetRPS,
      latencyProperty,
      recommendations: true,
    };
  }, [performanceViewEnabled, filterData, activeLatencyField]);

  const { catalogModels, catalogModelsLoaded, catalogModelsLoadError } = useCatalogModelsBySources(
    '',
    selectedSourceLabel === CategoryName.allModels ? undefined : selectedSourceLabel,
    10,
    searchTerm,
    effectiveFilterData,
    filterOptions,
    undefined, // filterQuery - will be computed from filterData and filterOptions
    sortParams.orderBy,
    sortParams.sortOrder,
    performanceParams,
  );

  const loaded = catalogModelsLoaded && filterOptionsLoaded;
  const loadError = catalogModelsLoadError || filterOptionsLoadError;

  const isNoLabelsSection = selectedSourceLabel === SourceLabel.other;

  const areOnlyDefaultFiltersApplied = React.useMemo(
    () => !hasFiltersApplied(filterData, BASIC_FILTER_KEYS),
    [filterData],
  );

  const noUserFiltersOrSearch = areOnlyDefaultFiltersApplied && !searchTerm;

  const shouldShowPerformanceEmptyState = React.useMemo(() => {
    const isEmptyResult = catalogModels.items.length === 0;
    const isNotAllModelsCategory = selectedSourceLabel !== CategoryName.allModels;
    const isPerformanceExcludedSection = isNoLabelsSection || noUserFiltersOrSearch;

    return (
      performanceViewEnabled &&
      isEmptyResult &&
      isNotAllModelsCategory &&
      isPerformanceExcludedSection
    );
  }, [
    performanceViewEnabled,
    catalogModels.items.length,
    selectedSourceLabel,
    isNoLabelsSection,
    noUserFiltersOrSearch,
  ]);

  const handleDisablePerformanceView = () => {
    setPerformanceViewEnabled(false);
  };

  const handleSelectAllModels = () => {
    updateSelectedSourceLabel(CategoryName.allModels);
  };

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

  if (shouldShowPerformanceEmptyState) {
    return (
      <EmptyModelCatalogState
        testid="performance-empty-state"
        title="No performance data available in selected category"
        headerIcon={ChartBarIcon}
        variant={EmptyStateVariant.lg}
        description={
          <>
            No models in the{' '}
            <strong>
              {selectedSourceLabel === 'null'
                ? CategoryName.otherModels
                : generateCategoryName(selectedSourceLabel || '')}
            </strong>{' '}
            category have performance data. Select another model category, or turn off model
            performance view to see models in the selected category.
          </>
        }
        primaryAction={
          <Button variant="primary" onClick={handleSelectAllModels}>
            View all models with performance data
          </Button>
        }
        secondaryAction={
          <Button variant="link" onClick={handleDisablePerformanceView}>
            Turn off model performance view
          </Button>
        }
      />
    );
  }

  if (catalogModels.items.length === 0 && !searchTerm && !filtersApplied) {
    return (
      <EmptyModelCatalogState
        testid="empty-model-catalog-state"
        title="No models available"
        headerIcon={SearchIcon}
        description="No models are available in this category."
      />
    );
  }

  if (catalogModels.items.length === 0 && (searchTerm || filtersApplied)) {
    return (
      <EmptyModelCatalogState
        testid="empty-model-catalog-state"
        title="No result found"
        headerIcon={SearchIcon}
        description="Adjust your filters and try again."
        primaryAction={<Button onClick={handleFilterReset}>Reset filters</Button>}
      />
    );
  }

  return (
    <>
      <ScrollViewOnMount shouldScroll scrollToTop />
      <Gallery hasGutter minWidths={{ default: '300px' }}>
        {catalogModels.items.map((model: CatalogModel) => (
          <ModelCatalogCard
            key={`${model.name}/${model.source_id}`}
            model={model}
            source={getSourceFromSourceId(model.source_id || '', catalogSources)}
            truncate
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

export default ModelCatalogGalleryView;
