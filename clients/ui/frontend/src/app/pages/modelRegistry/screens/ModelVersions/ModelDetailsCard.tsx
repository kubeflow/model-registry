import React from 'react';
import {
  Card,
  DescriptionList,
  StackItem,
  Stack,
  CardBody,
  CardTitle,
  ClipboardCopy,
  Content,
} from '@patternfly/react-core';
import { EditableTextDescriptionListGroup, DashboardDescriptionListGroup } from 'mod-arch-shared';
import ModelEditableLabelsDescriptionListGroup from '~/app/pages/modelRegistry/screens/components/ModelEditableLabelsDescriptionListGroup';
import { RegisteredModel } from '~/app/types';
import ModelTimestamp from '~/app/pages/modelRegistry/screens/components/ModelTimestamp';
import ModelPropertiesExpandableSection from '~/app/pages/modelRegistry/screens/components/ModelPropertiesExpandableSection';
import { ModelRegistryContext } from '~/app/context/ModelRegistryContext';

type ModelDetailsCardProps = {
  registeredModel: RegisteredModel;
  refresh: () => void;
  isArchiveModel?: boolean;
};

const ModelDetailsCard: React.FC<ModelDetailsCardProps> = ({
  registeredModel: rm,
  refresh,
  isArchiveModel,
}) => {
  const { apiState } = React.useContext(ModelRegistryContext);

  return (
    <Card>
      <CardTitle>Model details</CardTitle>
      <CardBody>
        <Stack hasGutter>
          <StackItem>
            <DescriptionList>
              <ModelEditableLabelsDescriptionListGroup
                isArchiveModel={isArchiveModel}
                rm={rm}
                refresh={refresh}
              />
              <EditableTextDescriptionListGroup
                editableVariant="TextArea"
                baseTestId="model-description"
                title="Description"
                isArchive={isArchiveModel}
                contentWhenEmpty="No description"
                value={rm.description || ''}
                saveEditedValue={(value) =>
                  apiState.api
                    .patchRegisteredModel(
                      {},
                      {
                        description: value,
                      },
                      rm.id,
                    )
                    .then(refresh)
                }
              />
            </DescriptionList>
          </StackItem>
          <StackItem>
            <DescriptionList columnModifier={{ default: '1Col', md: '2Col' }}>
              <DashboardDescriptionListGroup
                title="Owner"
                popover="The owner is the user who registered the model."
              >
                <Content component="p" data-testid="registered-model-owner">
                  {rm.owner || '-'}
                </Content>
              </DashboardDescriptionListGroup>
              <DashboardDescriptionListGroup title="Model ID">
                <ClipboardCopy
                  hoverTip="Copy"
                  clickTip="Copied"
                  variant="inline-compact"
                  data-testid="registered-model-id-clipboard-copy"
                >
                  {rm.id}
                </ClipboardCopy>
              </DashboardDescriptionListGroup>
              <DashboardDescriptionListGroup
                isEmpty={!rm.lastUpdateTimeSinceEpoch}
                contentWhenEmpty="Unknown"
                title="Last modified"
              >
                <ModelTimestamp timeSinceEpoch={rm.lastUpdateTimeSinceEpoch} />
              </DashboardDescriptionListGroup>
              <DashboardDescriptionListGroup
                isEmpty={!rm.createTimeSinceEpoch}
                contentWhenEmpty="Unknown"
                title="Created"
              >
                <ModelTimestamp timeSinceEpoch={rm.createTimeSinceEpoch} />
              </DashboardDescriptionListGroup>
            </DescriptionList>
          </StackItem>
          <StackItem>
            <ModelPropertiesExpandableSection
              isArchive={isArchiveModel}
              customProperties={rm.customProperties}
              saveEditedCustomProperties={(editedProperties) =>
                apiState.api
                  .patchRegisteredModel({}, { customProperties: editedProperties }, rm.id)
                  .then(refresh)
              }
            />
          </StackItem>
          {/* TODO: Add model card markdown here  */}
        </Stack>
      </CardBody>
    </Card>
  );
};

export default ModelDetailsCard;
