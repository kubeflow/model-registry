import React from 'react';
import {
  Button,
  Card,
  CardBody,
  CardFooter,
  CardHeader,
  CardTitle,
  Flex,
  FlexItem,
  Label,
  Skeleton,
  Truncate,
} from '@patternfly/react-core';
import { Link } from 'react-router-dom';
import {
  CatalogModel,
  CatalogSource,
  CatalogPerformanceMetricsArtifact,
  CatalogAccuracyMetricsArtifact,
} from '~/app/modelCatalogTypes';
import { catalogModelDetailsFromModel } from '~/app/routes/modelCatalog/catalogModel';
import { getLabels } from '~/app/pages/modelRegistry/screens/utils';
import { isModelValidated, getModelName } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import ModelCatalogLabels from './ModelCatalogLabels';
import ModelCatalogCardBody from './ModelCatalogCardBody';

type ModelCatalogCardProps = {
  model: CatalogModel;
  source: CatalogSource | undefined;
  truncate?: boolean;
  // TODO: Later these will be fetched based on the model, for now using props
  performanceMetrics?: CatalogPerformanceMetricsArtifact[];
  accuracyMetrics?: CatalogAccuracyMetricsArtifact[];
};

const ModelCatalogCard: React.FC<ModelCatalogCardProps> = ({
  model,
  source,
  truncate = false,
  performanceMetrics,
  accuracyMetrics,
}) => {
  // Extract labels from customProperties and check for validated label
  const allLabels = model.customProperties ? getLabels(model.customProperties) : [];
  const validatedLabels = allLabels.includes('validated') ? ['validated'] : [];
  const isValidated = isModelValidated(model);

  return (
    <Card isFullHeight data-testid="model-catalog-card" key={`${model.name}/${model.source_id}`}>
      <CardHeader>
        <CardTitle>
          <Flex alignItems={{ default: 'alignItemsFlexStart' }} className="pf-v6-u-mb-md">
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
              {isValidated ? (
                <Label color="purple">Validated</Label>
              ) : (
                source && <Label>{source.name}</Label>
              )}
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
                <Truncate
                  data-testid="model-catalog-card-name"
                  content={getModelName(model.name)}
                  position="middle"
                  tooltipPosition="top"
                  style={{ textDecoration: 'underline' }}
                />
              ) : (
                <span data-testid="model-catalog-card-name">{getModelName(model.name)}</span>
              )}
            </Button>
          </Link>
        </CardTitle>
      </CardHeader>
      <CardBody>
        <ModelCatalogCardBody
          model={model}
          isValidated={isValidated}
          performanceMetrics={performanceMetrics}
          accuracyMetrics={accuracyMetrics}
        />
      </CardBody>
      <CardFooter>
        <ModelCatalogLabels
          tasks={model.tasks ?? []}
          license={model.license}
          provider={model.provider}
          labels={validatedLabels}
          numLabels={isValidated ? 2 : 3}
        />
      </CardFooter>
    </Card>
  );
};

export default ModelCatalogCard;
