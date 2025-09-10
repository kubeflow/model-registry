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
  Stack,
  StackItem,
  Truncate,
} from '@patternfly/react-core';
import { useNavigate } from 'react-router-dom';
import { CatalogModel } from '~/app/modelCatalogTypes';
import { getModelName } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import { modelCatalogDetailsUrl } from '~/app/pages/modelCatalog/routeUtils';
import ModelCatalogLabels from './ModelCatalogLabels';

type ModelCatalogCardProps = {
  model: CatalogModel;
  source: string;
  truncate?: boolean;
};

const ModelCatalogCard: React.FC<ModelCatalogCardProps> = ({ model, source, truncate = false }) => {
  const navigate = useNavigate();

  return (
    <Card isFullHeight data-testid="model-catalog-card" key={`${model.name}/${model.sourceId}`}>
      <CardHeader>
        <CardTitle>
          <Flex alignItems={{ default: 'alignItemsCenter' }}>
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
              <Label>{source}</Label>
            </FlexItem>
          </Flex>
        </CardTitle>
      </CardHeader>
      <CardBody>
        <Stack hasGutter>
          <StackItem isFilled>
            <Button
              data-testid="model-catalog-detail-link"
              variant="link"
              isInline
              component="a"
              onClick={() => {
                navigate(modelCatalogDetailsUrl(model));
              }}
              style={{
                fontSize: 'var(--pf-t--global--font--size--body--default)',
                fontWeight: 'var(--pf-t--global--font--weight--body--bold)',
              }}
            >
              {truncate ? (
                <Truncate
                  data-testid="model-catalog-card-name"
                  content={model.name}
                  position="middle"
                  tooltipPosition="top"
                  style={{ textDecoration: 'underline' }}
                />
              ) : (
                <span>{getModelName(model.name)}</span>
              )}
            </Button>
          </StackItem>
          <StackItem isFilled data-testid="model-catalog-card-description">
            {truncate ? (
              <div
                style={{
                  overflow: 'hidden',
                  textOverflow: 'ellipsis',
                  WebkitLineClamp: 4,
                  WebkitBoxOrient: 'vertical',
                  display: '-webkit-box',
                }}
              >
                {model.longDescription}
              </div>
            ) : (
              model.longDescription
            )}
          </StackItem>
        </Stack>
      </CardBody>
      <CardFooter>
        <ModelCatalogLabels
          tasks={model.tasks ?? []}
          license={model.license}
          provider={model.provider}
        />
      </CardFooter>
    </Card>
  );
};

export default ModelCatalogCard;
