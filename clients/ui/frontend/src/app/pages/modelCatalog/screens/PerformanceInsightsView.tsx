import * as React from 'react';
import { PageSection, Card, CardBody, Title, Flex, FlexItem, Alert } from '@patternfly/react-core';
import { useParams } from 'react-router-dom';
import HardwareConfigurationTable from '~/app/pages/modelCatalog/components/HardwareConfigurationTable';
import { CatalogModelDetailsParams } from '~/app/modelCatalogTypes';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import { useCatalogPerformanceArtifacts } from '~/app/hooks/modelCatalog/useCatalogPerformanceArtifacts';
import {
  ModelCatalogNumberFilterKey,
  ModelCatalogStringFilterKey,
  DEFAULT_PERFORMANCE_FILTERS_QUERY_NAME,
} from '~/concepts/modelCatalog/const';
import {
  decodeParams,
  getActiveLatencyFieldName,
} from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import { parseLatencyFieldName } from '~/app/pages/modelCatalog/utils/hardwareConfigurationFilterUtils';
import {
  applyFilterValue,
  getDefaultFiltersFromNamedQuery,
} from '~/app/pages/modelCatalog/utils/performanceFilterUtils';

const PerformanceInsightsView = (): React.JSX.Element => {
  const params = useParams<CatalogModelDetailsParams>();
  const decodedParams = decodeParams(params);
  const {
    filterData,
    filterOptions,
    filterOptionsLoaded,
    performanceViewEnabled,
    setPerformanceFiltersChangedOnDetailsPage,
    setFilterData,
  } = React.useContext(ModelCatalogContext);

  // Apply default performance filters on mount if they don't have values yet
  // Only apply defaults if performance view toggle is enabled
  React.useEffect(() => {
    if (!filterOptionsLoaded || !filterOptions?.namedQueries || !performanceViewEnabled) {
      return;
    }

    // Check if any performance filter already has a value
    const hasUseCaseValue = filterData[ModelCatalogStringFilterKey.USE_CASE].length > 0;
    const hasRpsValue = filterData[ModelCatalogNumberFilterKey.MAX_RPS] !== undefined;
    const hasLatencyValue = getActiveLatencyFieldName(filterData) !== undefined;

    // If no performance filters are set, apply defaults
    if (!hasUseCaseValue && !hasRpsValue && !hasLatencyValue) {
      const defaultQuery = filterOptions.namedQueries[DEFAULT_PERFORMANCE_FILTERS_QUERY_NAME];
      const defaults = getDefaultFiltersFromNamedQuery(filterOptions, defaultQuery);
      Object.entries(defaults).forEach(([filterKey, value]) => {
        applyFilterValue(setFilterData, filterKey, value);
      });
    }
    // Only run on mount when filterOptions become available
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [filterOptionsLoaded, performanceViewEnabled]);

  // Get performance-specific filter params for the /performance_artifacts endpoint
  const targetRPS = filterData[ModelCatalogNumberFilterKey.MAX_RPS];
  // Get full filter key and convert to short property key for the catalog API
  const latencyFieldName = getActiveLatencyFieldName(filterData);
  const latencyProperty = latencyFieldName
    ? parseLatencyFieldName(latencyFieldName).propertyKey
    : undefined;

  // Fetch performance artifacts from server with filtering/sorting/pagination
  const [performanceArtifactsList, performanceArtifactsLoaded, performanceArtifactsError] =
    useCatalogPerformanceArtifacts(
      decodedParams.sourceId || '',
      encodeURIComponent(`${decodedParams.modelName}`),
      {
        targetRPS,
        latencyProperty,
        recommendations: true,
      },
      filterData,
      filterOptions,
    );

  React.useEffect(() => {
    setPerformanceFiltersChangedOnDetailsPage(false);
  }, [setPerformanceFiltersChangedOnDetailsPage]);

  if (performanceArtifactsError) {
    return (
      <PageSection padding={{ default: 'noPadding' }}>
        <Alert variant="danger" isInline title="Error loading performance data">
          {performanceArtifactsError.message}
        </Alert>
      </PageSection>
    );
  }

  return (
    <PageSection padding={{ default: 'noPadding' }}>
      <Card>
        <CardBody>
          <Flex direction={{ default: 'column' }} gap={{ default: 'gapLg' }}>
            <FlexItem>
              <Flex direction={{ default: 'column' }} gap={{ default: 'gapSm' }}>
                <FlexItem>
                  <Title headingLevel="h2" size="lg">
                    Hardware Configuration
                  </Title>
                </FlexItem>
                <FlexItem>
                  <p>
                    Compare the performance metrics of hardware configuration to determine the most
                    suitable option for deployment.
                  </p>
                </FlexItem>
              </Flex>
            </FlexItem>
            <FlexItem>
              <HardwareConfigurationTable
                performanceArtifacts={performanceArtifactsList.items}
                isLoading={!performanceArtifactsLoaded}
              />
            </FlexItem>
          </Flex>
        </CardBody>
      </Card>
    </PageSection>
  );
};

export default PerformanceInsightsView;
