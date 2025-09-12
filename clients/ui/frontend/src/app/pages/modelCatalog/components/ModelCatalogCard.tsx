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
import { CatalogModel, CatalogSource } from '~/app/modelCatalogTypes';
import { getModelName } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import { catalogModelDetailsFromModel } from '~/app/routes/modelCatalog/catalogModel';
import ModelCatalogLabels from './ModelCatalogLabels';

type ModelCatalogCardProps = {
  model: CatalogModel;
  source: CatalogSource | undefined;
  truncate?: boolean;
};

const ModelCatalogCard: React.FC<ModelCatalogCardProps> = ({ model, source, truncate = false }) => {
  const navigate = useNavigate();

  return (
    <Card isFullHeight data-testid="model-catalog-card" key={`${model.name}/${model.source_id}`}>
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
              {source && <Label>{source.name}</Label>}
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
                navigate(catalogModelDetailsFromModel(model.name, source?.id));
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
                {model.description}
              </div>
            ) : (
              model.description
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
