import * as React from 'react';
import {
  Card,
  CardHeader,
  CardTitle,
  CardBody,
  CardFooter,
  Flex,
  FlexItem,
  Button,
} from '@patternfly/react-core';
import { ExternalLinkAltIcon } from '@patternfly/react-icons';
import { ModelCatalogItem } from '~/app/modelCatalogTypes';
import ModelCatalogLabels from './ModelCatalogLabels';

type ModelCatalogCardProps = {
  model: ModelCatalogItem;
  onSelect?: (model: ModelCatalogItem) => void;
};

const ModelCatalogCard: React.FC<ModelCatalogCardProps> = ({ model, onSelect }) => {
  const { displayName, description, provider, url, tags, framework, task, license, metrics } =
    model;

  return (
    <Card isCompact>
      <CardHeader>
        <CardTitle>{displayName}</CardTitle>
        {provider && <div className="pf-v5-u-color-200">{provider}</div>}
      </CardHeader>
      <CardBody>
        {description && <p>{description}</p>}
        <ModelCatalogLabels tags={tags} framework={framework} task={task} license={license} />
        {metrics && Object.keys(metrics).length > 0 && (
          <div className="pf-v5-u-font-size-sm">
            Metrics:{' '}
            {Object.entries(metrics)
              .map(([key, value]) => `${key}: ${value}`)
              .join(', ')}
          </div>
        )}
      </CardBody>
      <CardFooter>
        <Flex>
          <FlexItem>
            <Button
              variant="primary"
              onClick={() => onSelect?.(model)}
              data-testid="select-model-button"
            >
              Select model
            </Button>
          </FlexItem>
          {url && (
            <FlexItem>
              <Button
                variant="link"
                component="a"
                href={url}
                target="_blank"
                rel="noopener noreferrer"
                icon={<ExternalLinkAltIcon />}
                iconPosition="right"
                data-testid="view-model-link"
              >
                View model
              </Button>
            </FlexItem>
          )}
        </Flex>
      </CardFooter>
    </Card>
  );
};

export default ModelCatalogCard;
