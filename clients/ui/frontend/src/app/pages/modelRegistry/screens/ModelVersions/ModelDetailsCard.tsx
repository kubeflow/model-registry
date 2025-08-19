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
  CardHeader,
  CardExpandableContent,
  Sidebar,
  SidebarPanel,
  SidebarContent,
} from '@patternfly/react-core';
import {
  EditableTextDescriptionListGroup,
  DashboardDescriptionListGroup,
  EditableLabelsDescriptionListGroup,
} from 'mod-arch-shared';
import { RegisteredModel } from '~/app/types';
import ModelTimestamp from '~/app/pages/modelRegistry/screens/components/ModelTimestamp';
import ModelPropertiesExpandableSection from '~/app/pages/modelRegistry/screens/components/ModelPropertiesExpandableSection';
import { ModelRegistryContext } from '~/app/context/ModelRegistryContext';
import { getLabels, mergeUpdatedLabels } from '~/app/pages/modelRegistry/screens/utils';

type ModelDetailsCardProps = {
  registeredModel: RegisteredModel;
  refresh: () => void;
  isArchiveModel?: boolean;
  isExpandable?: boolean;
};

const ModelDetailsCard: React.FC<ModelDetailsCardProps> = ({
  registeredModel: rm,
  refresh,
  isArchiveModel,
  isExpandable,
}) => {
  const { apiState } = React.useContext(ModelRegistryContext);
  const [isExpanded, setIsExpanded] = React.useState(false);

  const labelsSection = (
    <EditableLabelsDescriptionListGroup
      labels={getLabels(rm.customProperties)}
      isArchive={isArchiveModel}
      allExistingKeys={Object.keys(rm.customProperties)}
      title="Labels"
      contentWhenEmpty="No labels"
      onLabelsChange={(editedLabels) =>
        apiState.api
          .patchRegisteredModel(
            {},
            {
              customProperties: mergeUpdatedLabels(rm.customProperties, editedLabels),
            },
            rm.id,
          )
          .then(refresh)
      }
      isCollapsible={false}
      labelProps={{ variant: 'outline' }}
    />
  );

  const descriptionSection = (
    <EditableTextDescriptionListGroup
      truncateMaxLines={3}
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
  );

  const infoSection = (
    <>
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
    </>
  );

  const propertiesSection = (
    <ModelPropertiesExpandableSection
      isArchive={isArchiveModel}
      customProperties={rm.customProperties}
      saveEditedCustomProperties={(editedProperties) =>
        apiState.api
          .patchRegisteredModel({}, { customProperties: editedProperties }, rm.id)
          .then(refresh)
      }
    />
  );

  const cardBody = (
    <CardBody>
      {isExpandable ? (
        <Sidebar hasBorder hasGutter isPanelRight>
          <SidebarContent>
            <DescriptionList>
              {labelsSection}
              {descriptionSection}
              {propertiesSection}
            </DescriptionList>
            {/* TODO: Add model card markdown here  */}
          </SidebarContent>
          <SidebarPanel width={{ default: 'width_33' }}>
            <DescriptionList>{infoSection}</DescriptionList>
          </SidebarPanel>
        </Sidebar>
      ) : (
        <Stack hasGutter>
          <StackItem>
            <DescriptionList>
              {labelsSection}
              {descriptionSection}
            </DescriptionList>
          </StackItem>
          <StackItem>
            <DescriptionList columnModifier={{ default: '1Col', md: '2Col' }}>
              {infoSection}
            </DescriptionList>
          </StackItem>
          <StackItem>{propertiesSection}</StackItem>
          {/* TODO: Add model card markdown here  */}
        </Stack>
      )}
    </CardBody>
  );

  return (
    <Card isExpanded={isExpanded} style={{ overflow: 'visible' }}>
      {isExpandable ? (
        <>
          <CardHeader
            onExpand={() => setIsExpanded(!isExpanded)}
            toggleButtonProps={{
              id: 'toggle-button1',
              'data-testid': 'model-details-card-toggle-button',
              'aria-label': 'Details',
              'aria-expanded': isExpanded,
            }}
          >
            <CardTitle>Model details</CardTitle>
          </CardHeader>
          <CardExpandableContent>{cardBody}</CardExpandableContent>
        </>
      ) : (
        <>
          <CardTitle>Model details</CardTitle>
          {cardBody}
        </>
      )}
    </Card>
  );
};

export default ModelDetailsCard;
