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
import { HelpIcon, AngleLeftIcon, AngleRightIcon, ArrowRightIcon } from '@patternfly/react-icons';
import { TruncatedText } from 'mod-arch-shared';
import { CatalogModel, CatalogSource } from '~/app/modelCatalogTypes';
import {
  extractValidatedModelMetrics,
  getLatencyValue,
} from '~/app/pages/modelCatalog/utils/validatedModelUtils';
import { catalogModelDetailsTabFromModel } from '~/app/routes/modelCatalog/catalogModel';
import {
  ModelDetailsTab,
  ModelCatalogNumberFilterKey,
  LatencyMetric,
  parseLatencyFilterKey,
  SortOrder,
} from '~/concepts/modelCatalog/const';
import { useCatalogPerformanceArtifacts } from '~/app/hooks/modelCatalog/useCatalogPerformanceArtifacts';
import {
  getActiveLatencyFieldName,
  stripArtifactsPrefix,
} from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import { formatLatency } from '~/app/pages/modelCatalog/utils/performanceMetricsUtils';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';

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
  const { filterData, filterOptions, performanceViewEnabled } =
    React.useContext(ModelCatalogContext);

  const handlePreviousBenchmark = () => {
    setCurrentPerformanceIndex((prev) => (prev > 0 ? prev - 1 : performanceMetrics.length - 1));
  };

  const handleNextBenchmark = () => {
    setCurrentPerformanceIndex((prev) => (prev < performanceMetrics.length - 1 ? prev + 1 : 0));
  };

  // Get performance-specific filter params for the /performance_artifacts endpoint
  // Only apply performance filters when toggle is ON
  const targetRPS = performanceViewEnabled
    ? filterData[ModelCatalogNumberFilterKey.MAX_RPS]
    : undefined;
  // Get full filter key for display purposes
  const latencyFieldName = performanceViewEnabled
    ? getActiveLatencyFieldName(filterData)
    : undefined;
  // Use short property key (e.g., 'ttft_p90') for the catalog API, not the full filter key
  const latencyProperty = latencyFieldName
    ? parseLatencyFilterKey(latencyFieldName).propertyKey
    : undefined;

  // Fetch performance artifacts from the new endpoint with server-side filtering
  // When toggle is OFF, don't pass filterData so no perf filters are applied
  const [performanceArtifactsList, performanceArtifactsLoaded, performanceArtifactsError] =
    useCatalogPerformanceArtifacts(
      source?.id || '',
      model.name,
      {
        targetRPS,
        latencyProperty,
        recommendations: true,
        // TODO this is a temporary workaround to avoid capping performance artifacts with a default page size of 20.
        //      we need to implement proper cursor-based pagination as the user clicks through artifacts on a card.
        pageSize: '999',
        // If a latency filter is applied, sort artifacts on the card by lowest latency.
        ...(latencyFieldName && {
          orderBy: stripArtifactsPrefix(latencyFieldName),
          sortOrder: SortOrder.ASC,
        }),
      },
      performanceViewEnabled ? filterData : undefined,
      performanceViewEnabled ? filterOptions : undefined,
      isValidated, // Only fetch if validated
    );

  // Performance artifacts are already filtered by the server endpoint
  const performanceMetrics = performanceArtifactsList.items;

  // NOTE: Accuracy metrics are not currently returned by the /performance_artifacts endpoint.
  // This is kept as a placeholder for when accuracy metrics support is restored.
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const accuracyMetrics: any[] = [];

  const isLoading = isValidated && !performanceArtifactsLoaded;

  if (isLoading) {
    return <Spinner />;
  }

  if (performanceArtifactsError && isValidated) {
    return (
      <Alert variant="danger" isInline title={performanceArtifactsError.name}>
        {performanceArtifactsError.message}
      </Alert>
    );
  }

  if (isValidated && performanceMetrics.length > 0) {
    // When performance view toggle is OFF, show description with a link to benchmarks
    if (!performanceViewEnabled) {
      return (
        <Stack hasGutter>
          <StackItem>
            <TruncatedText
              content={model.description || ''}
              maxLines={4}
              data-testid="model-catalog-card-description"
            />
          </StackItem>
          <StackItem>
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
                icon={<ArrowRightIcon />}
                iconPosition="end"
                style={{ padding: 0, fontSize: 'inherit' }}
                data-testid="validated-model-benchmark-link"
              >
                View {performanceMetrics.length} benchmark
                {performanceMetrics.length !== 1 ? 's' : ''}
              </Button>
            </Link>
          </StackItem>
        </Stack>
      );
    }

    // When performance view toggle is ON, show hardware, latency and replicas data
    const metrics = extractValidatedModelMetrics(
      performanceMetrics,
      accuracyMetrics,
      currentPerformanceIndex,
    );

    // Get the selected latency metric from filters, or default to TTFT
    const activeLatencyField = latencyFieldName;
    const latencyValue =
      getLatencyValue(metrics.latencyMetrics, activeLatencyField) ?? metrics.ttftMean;
    const latencyLabel = activeLatencyField
      ? parseLatencyFilterKey(activeLatencyField).metric
      : LatencyMetric.TTFT;

    return (
      <Stack hasGutter>
        <StackItem>
          <Flex justifyContent={{ default: 'justifyContentSpaceBetween' }}>
            <Flex direction={{ default: 'column' }}>
              <span className="pf-v6-u-font-weight-bold" data-testid="validated-model-hardware">
                {metrics.hardwareConfiguration}
              </span>
              <Content component={ContentVariants.small}>Hardware</Content>
            </Flex>
            <Flex direction={{ default: 'column' }}>
              <span className="pf-v6-u-font-weight-bold" data-testid="validated-model-replicas">
                {metrics.replicas !== undefined ? metrics.replicas : metrics.rpsPerReplica}
              </span>
              <Content component={ContentVariants.small}>
                {metrics.replicas !== undefined ? 'Replicas' : 'RPS/rep.'}
              </Content>
            </Flex>
            <Flex direction={{ default: 'column' }}>
              <span className="pf-v6-u-font-weight-bold" data-testid="validated-model-latency">
                {formatLatency(latencyValue)}
              </span>
              <Flex alignItems={{ default: 'alignItemsBaseline' }} gap={{ default: 'gapXs' }}>
                <Content component={ContentVariants.small}>{latencyLabel}</Content>
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
    <TruncatedText
      content={model.description || ''}
      maxLines={4}
      data-testid="model-catalog-card-description"
    />
  );
};

export default ModelCatalogCardBody;
