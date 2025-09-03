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
  Icon,
  Label,
  Skeleton,
  Split,
  SplitItem,
  Stack,
  StackItem,
  Truncate,
} from '@patternfly/react-core';
import { TagIcon } from '@patternfly/react-icons';
import { useNavigate } from 'react-router-dom';
import { ModelCatalogItem } from '~/app/modelCatalogTypes';
import {
  extractVersionTag,
  filterNonVersionTags,
} from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import ModelCatalogLabels from './ModelCatalogLabels';

type ModelCatalogCardProps = {
  model: ModelCatalogItem;
  source: string;
  truncate?: boolean;
  onSelect?: (model: ModelCatalogItem) => void;
};

const ModelCatalogCard: React.FC<ModelCatalogCardProps> = ({
  model,
  source,
  truncate = false,
  onSelect,
}) => {
  const navigate = useNavigate();

  const versionTag = extractVersionTag(model.tags);
  const nonVersionTags = filterNonVersionTags(model.tags);

  return (
    <Card isFullHeight data-testid="model-catalog-card">
      <CardHeader>
        <CardTitle>
          <Flex alignItems={{ default: 'alignItemsCenter' }}>
            {model.logo ? (
              <img src={model.logo} alt="model logo" style={{ height: '36px', width: '36px' }} />
            ) : (
              <Skeleton
                shape="square"
                width="36px"
                height="36px"
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
                if (onSelect) {
                  onSelect(model);
                } else {
                  navigate(`/model-catalog/${encodeURIComponent(model.id)}` || '#');
                }
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
                <span>{model.name}</span>
              )}
            </Button>
            <Split hasGutter>
              <SplitItem>
                <Icon isInline>
                  <TagIcon />
                </Icon>
                <span style={{ marginLeft: 'var(--pf-t--global--spacer--sm)' }}>
                  {versionTag || 'No version'}
                </span>
              </SplitItem>
            </Split>
          </StackItem>
          <StackItem isFilled data-testid="model-catalog-card-description">
            {truncate ? (
              <div
                style={{
                  overflow: 'hidden',
                  textOverflow: 'ellipsis',
                  WebkitLineClamp: 2,
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
          tags={nonVersionTags}
          framework={model.framework}
          task={model.task}
          license={model.license}
        />
      </CardFooter>
    </Card>
  );
};

export default ModelCatalogCard;
