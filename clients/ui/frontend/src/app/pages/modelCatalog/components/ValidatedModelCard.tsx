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
  Stack,
  StackItem,
  Tooltip,
  Content,
  ContentVariants,
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
import ModelCatalogLabels from './ModelCatalogLabels';

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
  // Extract labels from customProperties and check for validated label
  const allLabels = model.customProperties ? getLabels(model.customProperties) : [];
  const validatedLabels = allLabels.includes('validated') ? ['validated'] : [];

  // Extract performance metrics data using utility function
  const metrics = extractValidatedModelMetrics(performanceMetrics, accuracyMetrics);

  return (
    <Card
      isFullHeight
      data-testid="validated-model-catalog-card"
      key={`${model.name}/${model.source_id}`}
    >
      <CardHeader>
        <CardTitle>
          <Flex alignItems={{ default: 'alignItemsCenter' }}>
            {/* Model icon placeholder - 3x3 grid with X marks */}
            <div
              style={{
                display: 'grid',
                gridTemplateColumns: 'repeat(3, 1fr)',
                gap: '2px',
                width: '24px',
                height: '24px',
                marginRight: '8px',
              }}
            >
              {Array.from({ length: 9 }, (_, i) => (
                <div
                  key={i}
                  style={{
                    width: '6px',
                    height: '6px',
                    backgroundColor: [0, 2, 4, 6, 8].includes(i) ? '#c9190b' : '#d2d2d2',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    fontSize: '4px',
                    color: 'white',
                  }}
                >
                  {[0, 2, 4, 6, 8].includes(i) && '✕'}
                </div>
              ))}
            </div>
            <FlexItem align={{ default: 'alignRight' }}>
              <Badge>Validated</Badge>
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
              <Button data-testid="model-catalog-detail-link" variant="link" tabIndex={-1} isInline>
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
        <Stack hasGutter>
          {/* Performance Metrics */}
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
                <Tooltip content="Placeholder tooltip text for accuracy">
                  <HelpIcon style={{ cursor: 'help' }} />
                </Tooltip>
              </Flex>
            </Flex>
          </StackItem>

          {/* Hardware Configuration and Latency */}
          <StackItem>
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
                  <Tooltip content="Placeholder tooltip text for TTFT">
                    <HelpIcon />
                  </Tooltip>
                </Flex>
              </div>
            </Flex>
          </StackItem>

          {/* Benchmarks */}
          <StackItem>
            <Flex
              alignItems={{ default: 'alignItemsCenter' }}
              justifyContent={{ default: 'justifyContentSpaceBetween' }}
            >
              <Content component={ContentVariants.small} data-testid="validated-model-benchmarks">
                1 of 3 benchmarks
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
