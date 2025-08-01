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
  Truncate,
} from '@patternfly/react-core';
import { ArrowRightIcon } from '@patternfly/react-icons';
import { TruncatedText } from 'mod-arch-shared';
import { ModelState, RegisteredModel } from '~/app/types';
import useModelVersionsByRegisteredModel from '~/app/hooks/useModelVersionsByRegisteredModel';
import { filterLiveVersions } from '~/app/utils';
import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';
import {
  archiveModelVersionDetailsUrl,
  archiveModelVersionListUrl,
  modelVersionListUrl,
  modelVersionUrl,
} from '~/app/pages/modelRegistry/screens/routeUtils';
import { getLabels } from '~/app/pages/modelRegistry/screens/utils';

type ModelVersionsCardProps = {
  rm: RegisteredModel;
  isArchiveModel?: boolean;
};

const ModelVersionsCard: React.FC<ModelVersionsCardProps> = ({ rm, isArchiveModel }) => {
  const [modelVersions] = useModelVersionsByRegisteredModel(rm.id);
  const { preferredModelRegistry } = React.useContext(ModelRegistrySelectorContext);
  const filteredVersions = isArchiveModel
    ? modelVersions.items
    : filterLiveVersions(modelVersions.items);
  const latestModelVersions = filteredVersions
    .toSorted((a, b) => Number(b.createTimeSinceEpoch) - Number(a.createTimeSinceEpoch))
    .slice(0, 3);

  return (
    <Card>
      <CardTitle>Latest versions</CardTitle>
      <CardBody>
        <Divider />
        {latestModelVersions.length > 0 ? (
          <List isPlain isBordered>
            {latestModelVersions.map((mv) => (
              <ListItem
                key={mv.id}
                className="pf-v6-u-py-md"
                data-testid={`model-version-${mv.id}`}
              >
                <Flex spaceItems={{ default: 'spaceItemsXs' }} direction={{ default: 'column' }}>
                  <FlexItem>
                    <Button
                      component="a"
                      isInline
                      data-testid={`model-version-${mv.id}-link`}
                      href={
                        mv.state === ModelState.ARCHIVED
                          ? archiveModelVersionDetailsUrl(
                              mv.id,
                              rm.id,
                              preferredModelRegistry?.name,
                            )
                          : modelVersionUrl(mv.id, rm.id, preferredModelRegistry?.name)
                      }
                      variant="link"
                    >
                      <Truncate content={mv.name} />
                    </Button>
                  </FlexItem>
                  <FlexItem>
                    <TruncatedText content={mv.description} maxLines={1} />
                  </FlexItem>
                  <FlexItem>
                    <LabelGroup>
                      {getLabels(mv.customProperties).map((label) => (
                        <Label
                          variant="outline"
                          key={label}
                          data-testid={`model-version-${mv.id}-property-${label}`}
                        >
                          {label}
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
                href={
                  isArchiveModel
                    ? archiveModelVersionListUrl(rm.id, preferredModelRegistry?.name)
                    : modelVersionListUrl(rm.id, preferredModelRegistry?.name)
                }
                variant="link"
                icon={<ArrowRightIcon />}
                iconPosition="right"
              >
                {`View all ${filteredVersions.length} versions`}
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
