import React, { useState } from 'react';
import {
  Button,
  Card,
  CardBody,
  CardFooter,
  CardHeader,
  CardTitle,
  Divider,
  Flex,
  FlexItem,
  Label,
  Popover,
  Content,
  ContentVariants,
  Skeleton,
  Stack,
  StackItem,
  List,
  ListItem,
  Truncate,
} from '@patternfly/react-core';
import {
  MonitoringIcon,
  HelpIcon,
  ChevronLeftIcon,
  ChevronRightIcon,
} from '@patternfly/react-icons';
import { Link } from 'react-router-dom';
import {
  CatalogModel,
  CatalogSource,
  CatalogPerformanceMetricsArtifact,
  CatalogAccuracyMetricsArtifact,
} from '~/app/modelCatalogTypes';
import { catalogModelDetailsFromModel } from '~/app/routes/modelCatalog/catalogModel';
import { getLabels } from '~/app/pages/modelRegistry/screens/utils';
import { extractValidatedModelMetrics } from '~/app/pages/modelCatalog/utils/validatedModelUtils';
import ModelCatalogLabels from '~/app/pages/modelCatalog/components/ModelCatalogLabels';

type ValidatedModelCardProps = {
  model: CatalogModel;
  source: CatalogSource | undefined;
  performanceMetrics: CatalogPerformanceMetricsArtifact[];
  accuracyMetrics: CatalogAccuracyMetricsArtifact[];
  truncate?: boolean;
};

const ValidatedModelCard: React.FC<ValidatedModelCardProps> = ({
  model,
  source,
  performanceMetrics,
  accuracyMetrics,
  truncate = false,
}) => {
  const [currentPerformanceIndex, setCurrentPerformanceIndex] = useState(0);

  const allLabels = model.customProperties ? getLabels(model.customProperties) : [];
  const validatedLabels = allLabels.includes('validated') ? ['validated'] : [];

  const metrics = extractValidatedModelMetrics(
    performanceMetrics,
    accuracyMetrics,
    currentPerformanceIndex,
  );

  const handlePreviousBenchmark = () => {
    setCurrentPerformanceIndex((prev) => (prev > 0 ? prev - 1 : performanceMetrics.length - 1));
  };

  const handleNextBenchmark = () => {
    setCurrentPerformanceIndex((prev) => (prev < performanceMetrics.length - 1 ? prev + 1 : 0));
  };

  return (
    <Card
      isFullHeight
      data-testid="validated-model-catalog-card"
      key={`${model.name}/${model.source_id}`}
    >
      <CardHeader>
        <CardTitle>
          <Flex alignItems={{ default: 'alignItemsCenter' }} gap={{ default: 'gapSm' }}>
            {model.logo ? (
              <img src={model.logo} alt="model logo" style={{ height: '56px', width: '56px' }} />
            ) : (
              <Skeleton
                shape="square"
                width="56px"
                height="56px"
                screenreaderText="Brand image loading"
              />
            )}
            <FlexItem align={{ default: 'alignRight' }}>
              <Label color="purple">Validated</Label>
            </FlexItem>
          </Flex>
          <Link to={catalogModelDetailsFromModel(model.name, source?.id)}>
            <Button
              data-testid="model-catalog-detail-link"
              variant="link"
              tabIndex={-1}
              isInline
              style={{
                fontSize: 'var(--pf-t--global--font--size--body--default)',
                fontWeight: 'var(--pf-t--global--font--weight--body--bold)',
              }}
            >
              {truncate ? (
                <Truncate data-testid="validated-model-catalog-card-name" content={model.name} />
              ) : (
                <span data-testid="validated-model-catalog-card-name">{model.name}</span>
              )}
            </Button>
          </Link>
        </CardTitle>
      </CardHeader>
      <CardBody>
        <Stack hasGutter>
          <StackItem>
            <Flex
              alignItems={{ default: 'alignItemsCenter' }}
              justifyContent={{ default: 'justifyContentSpaceBetween' }}
            >
              <Flex
                alignItems={{ default: 'alignItemsCenter' }}
                gap={{ default: 'gapSm' }}
                wrap="nowrap"
              >
                <MonitoringIcon />
                <span data-testid="validated-model-accuracy" className="pf-v6-u-font-weight-bold">
                  {metrics.accuracy}%
                </span>
              </Flex>
              <Flex alignItems={{ default: 'alignItemsCenter' }} gap={{ default: 'gapXs' }}>
                <span style={{ fontSize: '14px', color: 'var(--pf-v5-global--Color--200)' }}>
                  Average accuracy
                </span>
                <Popover
                  headerContent="Average accuracy"
                  bodyContent="The weighted average of normalized scores from all benchmarks. Each benchmark is normalized to a 0-100 scale. All normalized benchmarks are then averaged together."
                >
                  <button
                    type="button"
                    aria-label="More info for average accuracy"
                    style={{
                      all: 'unset',
                      cursor: 'pointer',
                      display: 'inline-flex',
                    }}
                  >
                    <HelpIcon />
                  </button>
                </Popover>
              </Flex>
            </Flex>
          </StackItem>

          <Divider />

          <StackItem>
            <Flex justifyContent={{ default: 'justifyContentSpaceBetween' }}>
              <Flex direction={{ default: 'column' }}>
                <span className="pf-v6-u-font-weight-bold" data-testid="validated-model-hardware">
                  {metrics.hardwareCount}x{metrics.hardware}
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
                  {metrics.ttftMean} ms
                </span>
                <Flex alignItems={{ default: 'alignItemsCenter' }} gap={{ default: 'gapXs' }}>
                  <span style={{ fontSize: '14px', color: 'var(--pf-v5-global--Color--200)' }}>
                    TTFT
                  </span>
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
                    <button
                      type="button"
                      aria-label="More info for latency"
                      style={{
                        all: 'unset',
                        cursor: 'pointer',
                        display: 'inline-flex',
                      }}
                    >
                      <HelpIcon />
                    </button>
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
              <Content component={ContentVariants.p} data-testid="validated-model-benchmarks">
                {currentPerformanceIndex + 1} of {performanceMetrics.length}{' '}
                <Button variant="link" isInline style={{ padding: 0, fontSize: 'inherit' }}>
                  benchmarks
                </Button>
              </Content>
              <Flex gap={{ default: 'gapSm' }} alignItems={{ default: 'alignItemsCenter' }}>
                <Button
                  variant="plain"
                  icon={<ChevronLeftIcon />}
                  aria-label="Previous benchmark"
                  data-testid="validated-model-benchmark-prev"
                  onClick={handlePreviousBenchmark}
                  isDisabled={performanceMetrics.length <= 1}
                />
                <Button
                  variant="plain"
                  icon={<ChevronRightIcon />}
                  aria-label="Next benchmark"
                  data-testid="validated-model-benchmark-next"
                  onClick={handleNextBenchmark}
                  isDisabled={performanceMetrics.length <= 1}
                />
              </Flex>
            </Flex>
          </StackItem>
        </Stack>
      </CardBody>
      <CardFooter>
        <ModelCatalogLabels
          tasks={model.tasks ?? []}
          license={model.license}
          provider={model.provider}
          labels={validatedLabels}
        />
      </CardFooter>
    </Card>
  );
};

export default ValidatedModelCard;
