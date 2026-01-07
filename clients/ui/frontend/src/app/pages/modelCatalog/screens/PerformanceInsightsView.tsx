import * as React from 'react';
import { PageSection, Card, CardBody, Title, Flex, FlexItem, Alert } from '@patternfly/react-core';
import { useParams } from 'react-router-dom';
import HardwareConfigurationTable from '~/app/pages/modelCatalog/components/HardwareConfigurationTable';
import { CatalogModelDetailsParams } from '~/app/modelCatalogTypes';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import { useCatalogPerformanceArtifacts } from '~/app/hooks/modelCatalog/useCatalogPerformanceArtifacts';
import { ModelCatalogNumberFilterKey } from '~/concepts/modelCatalog/const';
import {
  decodeParams,
  getActiveLatencyFieldName,
} from '~/app/pages/modelCatalog/utils/modelCatalogUtils';

const PerformanceInsightsView = (): React.JSX.Element => {
  const params = useParams<CatalogModelDetailsParams>();
  const decodedParams = decodeParams(params);
  const { filterData, filterOptions, setPerformanceFiltersChangedOnDetailsPage } =
    React.useContext(ModelCatalogContext);

  // Get performance-specific filter params for the /performance_artifacts endpoint
  const targetRPS = filterData[ModelCatalogNumberFilterKey.MIN_RPS];
  const latencyProperty = getActiveLatencyFieldName(filterData);

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
