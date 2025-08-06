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
import { TruncatedText } from 'mod-arch-shared';
import { ArrowRightIcon } from '@patternfly/react-icons';
import { RegisteredModel } from '~/app/types';
import useModelVersionsByRegisteredModel from '~/app/hooks/useModelVersionsByRegisteredModel';
import { filterLiveVersions } from '~/app/utils';
import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';
import {
  archiveModelVersionDetailsUrl,
  modelVersionUrl,
} from '~/app/pages/modelRegistry/screens/routeUtils';
import { getLabels } from '~/app/pages/modelRegistry/screens/utils';
import ViewAllVersionsButton from '~/app/pages/modelRegistry/screens/components/ViewAllVersionsButton';

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
                        isArchiveModel
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
              <ViewAllVersionsButton
                rmId={rm.id}
                totalVersions={filteredVersions.length}
                isArchiveModel={isArchiveModel}
                preferredModelRegistry={preferredModelRegistry}
                icon={<ArrowRightIcon />}
              />
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
