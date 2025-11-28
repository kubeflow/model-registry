import React, { useState } from 'react';
import {
  Alert,
  Button,
  Content,
  ContentVariants,
  Flex,
  List,
  ListItem,
  Popover,
  Spinner,
  Stack,
  StackItem,
} from '@patternfly/react-core';
import { Link } from 'react-router-dom';
import { HelpIcon, AngleLeftIcon, AngleRightIcon } from '@patternfly/react-icons';
import {
  CatalogModel,
  CatalogSource,
  CatalogArtifactType,
  MetricsType,
  CatalogPerformanceMetricsArtifact,
  CatalogAccuracyMetricsArtifact,
} from '~/app/modelCatalogTypes';
import { extractValidatedModelMetrics } from '~/app/pages/modelCatalog/utils/validatedModelUtils';
import { catalogModelDetailsTabFromModel } from '~/app/routes/modelCatalog/catalogModel';
import { ModelDetailsTab } from '~/app/pages/modelCatalog/screens/ModelDetailsTabs';
import { useCatalogModelArtifacts } from '~/app/hooks/modelCatalog/useCatalogModelArtifacts';
import { filterArtifactsByType } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import { formatLatency } from '~/app/pages/modelCatalog/utils/performanceMetricsUtils';

type ModelCatalogCardBodyProps = {
  model: CatalogModel;
  isValidated: boolean;
  source: CatalogSource | undefined;
};

const ModelCatalogCardBody: React.FC<ModelCatalogCardBodyProps> = ({
  model,
  isValidated,
  source,
}) => {
  const [currentPerformanceIndex, setCurrentPerformanceIndex] = useState(0);

  const handlePreviousBenchmark = () => {
    setCurrentPerformanceIndex((prev) => (prev > 0 ? prev - 1 : performanceMetrics.length - 1));
  };

  const handleNextBenchmark = () => {
    setCurrentPerformanceIndex((prev) => (prev < performanceMetrics.length - 1 ? prev + 1 : 0));
  };

  const [artifacts, artifactsLoaded, artifactsLoadError] = useCatalogModelArtifacts(
    source?.id || '',
    model.name,
    isValidated,
    true,
  );

  const performanceMetrics = filterArtifactsByType<CatalogPerformanceMetricsArtifact>(
    artifacts.items,
    CatalogArtifactType.metricsArtifact,
    MetricsType.performanceMetrics,
  );

  const accuracyMetrics = filterArtifactsByType<CatalogAccuracyMetricsArtifact>(
    artifacts.items,
    CatalogArtifactType.metricsArtifact,
    MetricsType.accuracyMetrics,
  );

  if (!artifactsLoaded && isValidated) {
    return <Spinner />;
  }

  if (artifactsLoadError && isValidated) {
    return (
      <Alert variant="danger" isInline title={artifactsLoadError.name}>
        {artifactsLoadError.message}
      </Alert>
    );
  }

  if (isValidated && performanceMetrics.length > 0 && accuracyMetrics.length > 0) {
    const metrics = extractValidatedModelMetrics(
      performanceMetrics,
      accuracyMetrics,
      currentPerformanceIndex,
    );

    return (
      <Stack hasGutter>
        <StackItem>
          <Flex justifyContent={{ default: 'justifyContentSpaceBetween' }}>
            <Flex direction={{ default: 'column' }}>
              <span className="pf-v6-u-font-weight-bold" data-testid="validated-model-hardware">
                {metrics.hardwareCount}x{metrics.hardwareType}
              </span>
              <Content component={ContentVariants.small}>Hardware</Content>
            </Flex>
            <Flex direction={{ default: 'column' }}>
              <span className="pf-v6-u-font-weight-bold" data-testid="validated-model-rps">
                {metrics.rpsPerReplica}
              </span>
              <Content component={ContentVariants.small}>RPS/rep.</Content>
            </Flex>
            <Flex direction={{ default: 'column' }}>
              <span className="pf-v6-u-font-weight-bold" data-testid="validated-model-ttft">
                {formatLatency(metrics.ttftMean)}
              </span>
              <Flex alignItems={{ default: 'alignItemsBaseline' }} gap={{ default: 'gapXs' }}>
                <Content component={ContentVariants.small}>TTFT</Content>
                <Popover
                  headerContent="Latency"
                  bodyContent={
                    <div>
                      <p>
                        The delay (in milliseconds) between sending a request and receiving the
                        first response.
                      </p>
                      <List>
                        <ListItem>
                          <strong>TTFT (Time to First Token)</strong> - The time between when a
                          request is sent to a model and when the model begins streaming its first
                          token in the response.
                        </ListItem>
                        <ListItem>
                          <strong>ITL (Inter-Token Latency)</strong> - The average time between
                          successive output tokens after the model has started generating.
                        </ListItem>
                        <ListItem>
                          <strong>E2E (End-to-End latency)</strong> - The total time from when the
                          request is sent until the last token is received.
                        </ListItem>
                      </List>
                    </div>
                  }
                >
                  <Button
                    icon={<HelpIcon />}
                    hasNoPadding
                    aria-label="More info for latency"
                    variant="plain"
                  />
                </Popover>
              </Flex>
            </Flex>
          </Flex>
        </StackItem>

        <StackItem>
          <Flex
            alignItems={{ default: 'alignItemsCenter' }}
            justifyContent={{ default: 'justifyContentSpaceBetween' }}
          >
            <span data-testid="validated-model-benchmarks">
              {currentPerformanceIndex + 1} of {performanceMetrics.length}{' '}
              <Link
                to={catalogModelDetailsTabFromModel(
                  ModelDetailsTab.PERFORMANCE_INSIGHTS,
                  model.name,
                  source?.id,
                )}
              >
                <Button
                  variant="link"
                  isInline
                  tabIndex={-1}
                  style={{ padding: 0, fontSize: 'inherit' }}
                  data-testid="validated-model-benchmark-link"
                >
                  benchmarks
                </Button>
              </Link>
            </span>
            <Flex gap={{ default: 'gapSm' }} alignItems={{ default: 'alignItemsCenter' }}>
              <Button
                variant="plain"
                icon={<AngleLeftIcon />}
                aria-label="Previous benchmark"
                data-testid="validated-model-benchmark-prev"
                onClick={handlePreviousBenchmark}
                isDisabled={performanceMetrics.length <= 1}
              />
              <Button
                variant="plain"
                icon={<AngleRightIcon />}
                aria-label="Next benchmark"
                data-testid="validated-model-benchmark-next"
                onClick={handleNextBenchmark}
                isDisabled={performanceMetrics.length <= 1}
              />
            </Flex>
          </Flex>
        </StackItem>
      </Stack>
    );
  }

  // Standard card body for non-validated models
  return (
    <div
      data-testid="model-catalog-card-description"
      style={{
        overflow: 'hidden',
        textOverflow: 'ellipsis',
        WebkitLineClamp: 4,
        WebkitBoxOrient: 'vertical',
        display: '-webkit-box',
      }}
    >
      {model.description}
    </div>
  );
};

export default ModelCatalogCardBody;
