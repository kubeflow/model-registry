import * as React from 'react';
import {
  Card,
  CardTitle,
  CardBody,
  List,
  ListItem,
  Label,
  LabelGroup,
  Button,
  Flex,
  FlexItem,
  Divider,
} from '@patternfly/react-core';
import { Link } from 'react-router';
import { ArrowRightIcon } from '@patternfly/react-icons';
import { RegisteredModel } from '~/app/types';
import useModelVersionsByRegisteredModel from '~/app/hooks/useModelVersionsByRegisteredModel';
import { filterLiveVersions } from '~/app/utils';
import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';
import { modelVersionListUrl, modelVersionUrl } from '../routeUtils';

type ModelVersionsCardProps = {
  rm: RegisteredModel;
};

const ModelVersionsCard: React.FC<ModelVersionsCardProps> = ({ rm }) => {
  const [modelVersions] = useModelVersionsByRegisteredModel(rm.id);
  const liveModelVersions = filterLiveVersions(modelVersions.items);
  const { preferredModelRegistry } = React.useContext(ModelRegistrySelectorContext);
  const latestModelVersions = liveModelVersions
    .toSorted((a, b) => Number(b.createTimeSinceEpoch) - Number(a.createTimeSinceEpoch))
    .slice(0, 3);

  return (
    <Card>
      <CardTitle>Latest versions</CardTitle>
      <CardBody>
        <Divider />
        {modelVersions.items.length > 0 ? (
          <List isPlain isBordered>
            {latestModelVersions.map((mv) => (
              <ListItem
                key={mv.id}
                className="pf-v6-u-py-md"
                data-testid={`model-version-${mv.id}`}
              >
                <Flex spaceItems={{ default: 'spaceItemsXs' }} direction={{ default: 'column' }}>
                  <FlexItem>
                    <Link
                      to={modelVersionUrl(mv.id, rm.id, preferredModelRegistry?.name)}
                      data-testid={`model-version-${mv.id}-link`}
                    >
                      {mv.name}
                    </Link>
                  </FlexItem>
                  <FlexItem>
                    <>{mv.description}</>
                  </FlexItem>
                  <FlexItem>
                    <LabelGroup>
                      {Object.keys(mv.customProperties).map((key) => (
                        <Label
                          variant="outline"
                          key={key}
                          data-testid={`model-version-${mv.id}-property-${key}`}
                        >
                          {key}
                        </Label>
                      ))}
                    </LabelGroup>
                  </FlexItem>
                </Flex>
              </ListItem>
            ))}
            <ListItem className="pf-v6-u-pt-md">
              <Button
                component="a"
                isInline
                data-testid="versions-route-link"
                href={modelVersionListUrl(rm.id, preferredModelRegistry?.name)}
                variant="link"
                icon={<ArrowRightIcon />}
                iconPosition="right"
              >
                View all versions
              </Button>
            </ListItem>
          </List>
        ) : (
          <div className="pf-v6-u-pt-md" data-testid="no-versions-text">
            No versions
          </div>
        )}
      </CardBody>
    </Card>
  );
};

export default ModelVersionsCard;
