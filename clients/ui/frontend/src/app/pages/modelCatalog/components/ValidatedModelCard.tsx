import React from 'react';
import {
  Badge,
  Button,
  Card,
  CardBody,
  CardFooter,
  CardHeader,
  CardTitle,
  Flex,
  FlexItem,
  Tooltip,
  Content,
  ContentVariants,
  Skeleton,
} from '@patternfly/react-core';
import text from '@patternfly/react-styles/css/utilities/Text/text';
import { ChartBarIcon, HelpIcon, ChevronLeftIcon, ChevronRightIcon } from '@patternfly/react-icons';
import { Link } from 'react-router-dom';
import {
  CatalogModel,
  CatalogSource,
  CatalogPerformanceMetricsArtifact,
  CatalogAccuracyMetricsArtifact,
} from '~/app/modelCatalogTypes';
import { getModelName } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import { catalogModelDetailsFromModel } from '~/app/routes/modelCatalog/catalogModel';
import { getLabels } from '~/app/pages/modelRegistry/screens/utils';
import { extractValidatedModelMetrics } from '~/app/pages/modelCatalog/utils/validatedModelUtils';
import ModelCatalogLabels from '~/app/pages/modelCatalog/components/ModelCatalogLabels';

type ValidatedModelCardProps = {
  model: CatalogModel;
  source: CatalogSource | undefined;
  performanceMetrics?: CatalogPerformanceMetricsArtifact;
  accuracyMetrics?: CatalogAccuracyMetricsArtifact;
  truncate?: boolean;
};

const ValidatedModelCard: React.FC<ValidatedModelCardProps> = ({
  model,
  source,
  performanceMetrics,
  accuracyMetrics,
  truncate = false,
}) => {
  const allLabels = model.customProperties ? getLabels(model.customProperties) : [];
  const validatedLabels = allLabels.includes('validated') ? ['validated'] : [];

  const metrics = extractValidatedModelMetrics(performanceMetrics, accuracyMetrics);

  return (
    <Card
      isFullHeight
      data-testid="validated-model-catalog-card"
      key={`${model.name}/${model.source_id}`}
    >
      <CardHeader>
        <CardTitle>
          <Flex alignItems={{ default: 'alignItemsCenter' }} style={{ paddingLeft: '8px' }}>
            {model.logo ? (
              <img
                src={model.logo}
                alt="model logo"
                style={{ height: '56px', width: '56px', marginRight: '8px' }}
              />
            ) : (
              <Skeleton
                shape="square"
                width="56px"
                height="56px"
                screenreaderText="Brand image loading"
                style={{ marginRight: '8px' }}
              />
            )}
            <FlexItem align={{ default: 'alignRight' }}>
              <Badge className="pf-m-purple">Validated</Badge>
            </FlexItem>
          </Flex>
          <div style={{ marginTop: '8px' }}>
            <Flex alignItems={{ default: 'alignItemsCenter' }} gap={{ default: 'gapSm' }}>
              <div
                data-testid="validated-model-status-indicator"
                className={text.textColorStatusSuccess}
                style={{
                  width: '8px',
                  height: '8px',
                  backgroundColor: 'var(--pf-t--global--palette--green-400)',
                  borderRadius: '50%',
                }}
              />
              <Content component={ContentVariants.small} className={text.textColorStatusSuccess}>
                Success
              </Content>
            </Flex>
            <Link to={catalogModelDetailsFromModel(model.name, source?.id)}>
              <Button
                data-testid="model-catalog-detail-link"
                variant="link"
                tabIndex={-1}
                isInline
                style={{ fontSize: '1.25rem', fontWeight: 'bold' }}
              >
                {truncate ? (
                  <span>{getModelName(model.name)}</span>
                ) : (
                  <span>{getModelName(model.name)}</span>
                )}
              </Button>
            </Link>
          </div>
        </CardTitle>
      </CardHeader>
      <CardBody>
        <div>
          <div style={{ paddingLeft: '8px', marginBottom: '16px' }}>
            <Flex
              alignItems={{ default: 'alignItemsCenter' }}
              justifyContent={{ default: 'justifyContentSpaceBetween' }}
            >
              <Flex
                alignItems={{ default: 'alignItemsCenter' }}
                gap={{ default: 'gapSm' }}
                wrap="nowrap"
              >
                <ChartBarIcon />
                <span
                  data-testid="validated-model-accuracy"
                  style={{ fontSize: '1.25rem', fontWeight: 'bold' }}
                >
                  {metrics.accuracy}%
                </span>
              </Flex>
              <Flex
                alignItems={{ default: 'alignItemsCenter' }}
                gap={{ default: 'gapSm' }}
                wrap="nowrap"
              >
                <span className="pf-v5-c-content pf-m-small">Average accuracy</span>
                <Tooltip
                  content="Average accuracy

The weighted average of normalized scores from all benchmarks. Each benchmark is normalized to a 0-100 scale. All normalized benchmarks are then averaged together."
                  position="top"
                  className="pf-v6-c-tooltip"
                >
                  <HelpIcon style={{ cursor: 'help' }} />
                </Tooltip>
              </Flex>
            </Flex>
          </div>

          {/* Hardware Configuration and Latency */}
          <div style={{ paddingLeft: '8px', marginBottom: '16px' }}>
            <Flex justifyContent={{ default: 'justifyContentSpaceBetween' }}>
              <div>
                <Content component={ContentVariants.small} style={{ color: '#6a6e73' }}>
                  Hardware
                </Content>
                <Content component={ContentVariants.p} data-testid="validated-model-hardware">
                  {metrics.hardwareCount}x{metrics.hardware}
                </Content>
              </div>
              <div>
                <Content component={ContentVariants.small} style={{ color: '#6a6e73' }}>
                  RPS/rep.
                </Content>
                <Content component={ContentVariants.p} data-testid="validated-model-rps">
                  {metrics.rpsPerReplica}
                </Content>
              </div>
              <div>
                <Content component={ContentVariants.small} style={{ color: '#6a6e73' }}>
                  TTFT
                </Content>
                <Flex
                  alignItems={{ default: 'alignItemsCenter' }}
                  gap={{ default: 'gapSm' }}
                  wrap="nowrap"
                >
                  <span data-testid="validated-model-ttft">{metrics.ttftMean} ms</span>
                  <Tooltip
                    content="Latency

The delay (in milliseconds) between sending a request and receiving the first response.

*TTFT (Time to First Token)
The time between when a request is sent to a model and when the model begins streaming its first token in the response.

*ITL (Inter-Token Latency)
The average time between successive output tokens after the model has started generating.

*E2E (End-to-End latency)
The total time from when the request is sent until the last token is received."
                    position="top"
                    className="pf-v6-c-tooltip"
                  >
                    <HelpIcon />
                  </Tooltip>
                </Flex>
              </div>
            </Flex>
          </div>

          {/* Benchmarks */}
          <div style={{ paddingLeft: '8px' }}>
            <Flex
              alignItems={{ default: 'alignItemsCenter' }}
              justifyContent={{ default: 'justifyContentSpaceBetween' }}
            >
              <Content component={ContentVariants.p} data-testid="validated-model-benchmarks">
                1 of 3{' '}
                <Button variant="link" isInline style={{ padding: 0, fontSize: 'inherit' }}>
                  benchmarks
                </Button>
              </Content>
              <Flex gap={{ default: 'gapSm' }}>
                <Button
                  variant="plain"
                  icon={<ChevronLeftIcon />}
                  aria-label="Previous benchmark"
                  data-testid="validated-model-benchmark-prev"
                />
                <Button
                  variant="plain"
                  icon={<ChevronRightIcon />}
                  aria-label="Next benchmark"
                  data-testid="validated-model-benchmark-next"
                />
              </Flex>
            </Flex>
          </div>
        </div>
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
