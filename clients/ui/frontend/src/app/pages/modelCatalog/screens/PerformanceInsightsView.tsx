import * as React from 'react';
import {
  PageSection,
  Card,
  CardBody,
  Title,
  Flex,
  FlexItem,
  Alert,
  Stack,
  StackItem,
} from '@patternfly/react-core';
import { useParams } from 'react-router-dom';
import HardwareConfigurationTable from '~/app/pages/modelCatalog/components/HardwareConfigurationTable';
import { CatalogModel, CatalogModelDetailsParams } from '~/app/modelCatalogTypes';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import { useCatalogPerformanceArtifacts } from '~/app/hooks/modelCatalog/useCatalogPerformanceArtifacts';
import { ModelCatalogNumberFilterKey } from '~/concepts/modelCatalog/const';
import {
  decodeParams,
  getActiveLatencyFieldName,
} from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import TensorTypeComparisonCard from './TensorTypeComparisonCard';

type PerformanceInsightsViewProps = {
  model: CatalogModel;
};

const PerformanceInsightsView: React.FC<PerformanceInsightsViewProps> = ({ model }) => {
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
    <Stack hasGutter>
      <StackItem>
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
                      Compare the performance metrics of hardware configuration to determine the
                      most suitable option for deployment.
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
      </StackItem>
      <StackItem>
        <TensorTypeComparisonCard model={model} />
      </StackItem>
    </Stack>
  );
};

export default PerformanceInsightsView;
