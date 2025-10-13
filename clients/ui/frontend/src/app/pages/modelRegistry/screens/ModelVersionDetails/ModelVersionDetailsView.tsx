import * as React from 'react';
import {
  DescriptionList,
  Divider,
  ContentVariants,
  Title,
  Bullseye,
  Spinner,
  Alert,
  StackItem,
  Stack,
  Card,
  CardHeader,
  CardBody,
  Sidebar,
  SidebarPanel,
  SidebarContent,
} from '@patternfly/react-core';
import {
  EditableLabelsDescriptionListGroup,
  EditableTextDescriptionListGroup,
  DashboardDescriptionListGroup,
  InlineTruncatedClipboardCopy,
} from 'mod-arch-shared';
import { ModelVersion, ModelArtifactList, RegisteredModel } from '~/app/types';
import { ModelRegistryContext } from '~/app/context/ModelRegistryContext';
import { getLabels, mergeUpdatedLabels } from '~/app/pages/modelRegistry/screens/utils';
import ModelPropertiesDescriptionListGroup from '~/app/pages/modelRegistry/screens/ModelPropertiesDescriptionListGroup';
import ModelTimestamp from '~/app/pages/modelRegistry/screens/components/ModelTimestamp';
import { bumpBothTimestamps, bumpRegisteredModelTimestamp } from '~/app/api/updateTimestamps';
import { uriToStorageFields } from '~/app/utils';
import ModelDetailsCard from '~/app/pages/modelRegistry/screens/ModelVersions/ModelDetailsCard';
import ModelVersionRegisteredFromLink from '~/app/pages/modelRegistry/screens/components/ModelVersionRegisteredFromLink';

type ModelVersionDetailsViewProps = {
  registeredModel: RegisteredModel | null;
  modelVersion: ModelVersion;
  isArchiveVersion?: boolean;
  refresh: () => void;
  modelArtifacts: ModelArtifactList;
  modelArtifactsLoaded: boolean;
  modelArtifactsLoadError: Error | undefined;
};

const ModelVersionDetailsView: React.FC<ModelVersionDetailsViewProps> = ({
  registeredModel,
  modelVersion: mv,
  isArchiveVersion,
  refresh,
  modelArtifacts,
  modelArtifactsLoaded,
  modelArtifactsLoadError,
}) => {
  const modelArtifact = modelArtifacts.items.length ? modelArtifacts.items[0] : null;
  const { apiState } = React.useContext(ModelRegistryContext);
  const storageFields = uriToStorageFields(modelArtifact?.uri || '');

  if (!modelArtifactsLoaded) {
    return (
      <Bullseye>
        <Spinner size="xl" />
      </Bullseye>
    );
  }
  const handleVersionUpdate = async (updatePromise: Promise<unknown>): Promise<void> => {
    await updatePromise;

    if (!mv.registeredModelId || !registeredModel) {
      return;
    }

    await bumpRegisteredModelTimestamp(apiState.api, registeredModel);
    refresh();
  };

  const handleArtifactUpdate = async (updatePromise: Promise<unknown>): Promise<void> => {
    try {
      await updatePromise;
      if (registeredModel) {
        await bumpBothTimestamps(apiState.api, registeredModel, mv);
        refresh();
      }
    } catch (error) {
      throw new Error(
        `Failed to update artifact: ${error instanceof Error ? error.message : String(error)}`,
      );
    }
  };

  return (
    <Stack hasGutter>
      {registeredModel && (
        <StackItem>
          <ModelDetailsCard
            registeredModel={registeredModel}
            refresh={refresh}
            isArchiveModel={isArchiveVersion}
            isExpandable
          />
        </StackItem>
      )}
      <StackItem>
        <Card>
          <CardHeader>
            <Title headingLevel="h2">Version details</Title>
          </CardHeader>
          <CardBody>
            <Sidebar hasBorder hasGutter isPanelRight>
              <SidebarContent>
                <DescriptionList>
                  <EditableLabelsDescriptionListGroup
                    labels={getLabels(mv.customProperties)}
                    isArchive={isArchiveVersion}
                    allExistingKeys={Object.keys(mv.customProperties)}
                    title="Labels"
                    contentWhenEmpty="No labels"
                    labelProps={{ variant: 'outline', color: 'grey' }}
                    onLabelsChange={(editedLabels) =>
                      handleVersionUpdate(
                        apiState.api.patchModelVersion(
                          {},
                          {
                            customProperties: mergeUpdatedLabels(mv.customProperties, editedLabels),
                          },
                          mv.id,
                        ),
                      )
                    }
                    data-testid="model-version-labels"
                  />
                  <EditableTextDescriptionListGroup
                    editableVariant="TextArea"
                    baseTestId="model-version-description"
                    isArchive={isArchiveVersion}
                    title="Description"
                    contentWhenEmpty="No description"
                    value={mv.description || ''}
                    saveEditedValue={(value) =>
                      handleVersionUpdate(
                        apiState.api.patchModelVersion({}, { description: value }, mv.id),
                      )
                    }
                  />
                  <ModelPropertiesDescriptionListGroup
                    isArchive={isArchiveVersion}
                    customProperties={mv.customProperties}
                    saveEditedCustomProperties={(editedProperties) =>
                      apiState.api
                        .patchModelVersion({}, { customProperties: editedProperties }, mv.id)
                        .then(refresh)
                    }
                  />
                </DescriptionList>
              </SidebarContent>
              <SidebarPanel width={{ default: 'width_33' }}>
                {modelArtifact && (
                  <ModelVersionRegisteredFromLink
                    modelArtifact={modelArtifact}
                    isModelCatalogAvailable
                  />
                )}
                <Title style={{ margin: '1em 0' }} headingLevel={ContentVariants.h3}>
                  Model location
                </Title>
                {modelArtifactsLoadError ? (
                  <Alert variant="danger" isInline title={modelArtifactsLoadError.name}>
                    {modelArtifactsLoadError.message}
                  </Alert>
                ) : (
                  <>
                    <DescriptionList>
                      {storageFields?.s3Fields && (
                        <>
                          <DashboardDescriptionListGroup
                            title="Endpoint"
                            isEmpty={!storageFields.s3Fields.endpoint}
                            contentWhenEmpty="No endpoint"
                          >
                            <InlineTruncatedClipboardCopy
                              testId="storage-endpoint"
                              textToCopy={storageFields.s3Fields.endpoint}
                            />
                          </DashboardDescriptionListGroup>
                          <DashboardDescriptionListGroup
                            title="Region"
                            isEmpty={!storageFields.s3Fields.region}
                            contentWhenEmpty="No region"
                          >
                            <InlineTruncatedClipboardCopy
                              testId="storage-region"
                              textToCopy={storageFields.s3Fields.region || ''}
                            />
                          </DashboardDescriptionListGroup>
                          <DashboardDescriptionListGroup
                            title="Bucket"
                            isEmpty={!storageFields.s3Fields.bucket}
                            contentWhenEmpty="No bucket"
                          >
                            <InlineTruncatedClipboardCopy
                              testId="storage-bucket"
                              textToCopy={storageFields.s3Fields.bucket}
                            />
                          </DashboardDescriptionListGroup>
                          <DashboardDescriptionListGroup
                            title="Path"
                            isEmpty={!storageFields.s3Fields.path}
                            contentWhenEmpty="No path"
                          >
                            <InlineTruncatedClipboardCopy
                              testId="storage-path"
                              textToCopy={storageFields.s3Fields.path}
                            />
                          </DashboardDescriptionListGroup>
                        </>
                      )}
                      {(storageFields?.uri || storageFields?.ociUri) && (
                        <>
                          <DashboardDescriptionListGroup
                            title="URI"
                            isEmpty={!modelArtifact?.uri}
                            contentWhenEmpty="No URI"
                          >
                            <InlineTruncatedClipboardCopy
                              testId="storage-uri"
                              textToCopy={modelArtifact?.uri || ''}
                            />
                          </DashboardDescriptionListGroup>
                        </>
                      )}
                    </DescriptionList>
                    <Divider style={{ marginTop: '1em' }} />
                    <DescriptionList>
                      <EditableTextDescriptionListGroup
                        editableVariant="TextInput"
                        baseTestId="source-model-format"
                        isArchive={isArchiveVersion}
                        value={modelArtifact?.modelFormatName || ''}
                        saveEditedValue={(value) =>
                          handleArtifactUpdate(
                            apiState.api.patchModelArtifact(
                              {},
                              { modelFormatName: value },
                              modelArtifact?.id || '',
                            ),
                          )
                        }
                        title="Model format"
                        contentWhenEmpty="No model format specified"
                      />
                      <EditableTextDescriptionListGroup
                        editableVariant="TextInput"
                        baseTestId="source-model-version"
                        value={modelArtifact?.modelFormatVersion || ''}
                        isArchive={isArchiveVersion}
                        saveEditedValue={(newVersion) =>
                          handleArtifactUpdate(
                            apiState.api.patchModelArtifact(
                              {},
                              { modelFormatVersion: newVersion },
                              modelArtifact?.id || '',
                            ),
                          )
                        }
                        title="Model format version"
                        contentWhenEmpty="No model format version"
                      />
                    </DescriptionList>
                  </>
                )}
                <Divider style={{ marginTop: '1em' }} />
                <DescriptionList isFillColumns style={{ marginTop: '1em' }}>
                  <DashboardDescriptionListGroup
                    title="Author"
                    popover="The author is the user who registered the model version."
                  >
                    {mv.author}
                  </DashboardDescriptionListGroup>
                  <DashboardDescriptionListGroup
                    title="Version ID"
                    isEmpty={!mv.id}
                    contentWhenEmpty="No model ID"
                  >
                    <InlineTruncatedClipboardCopy testId="model-version-id" textToCopy={mv.id} />
                  </DashboardDescriptionListGroup>
                  <DashboardDescriptionListGroup
                    title="Last modified"
                    isEmpty={!mv.lastUpdateTimeSinceEpoch}
                    contentWhenEmpty="Unknown"
                  >
                    <ModelTimestamp timeSinceEpoch={mv.lastUpdateTimeSinceEpoch} />
                  </DashboardDescriptionListGroup>
                  <DashboardDescriptionListGroup
                    title="Registered"
                    isEmpty={!mv.createTimeSinceEpoch}
                    contentWhenEmpty="Unknown"
                  >
                    <ModelTimestamp timeSinceEpoch={mv.createTimeSinceEpoch} />
                  </DashboardDescriptionListGroup>
                </DescriptionList>
              </SidebarPanel>
            </Sidebar>
          </CardBody>
        </Card>
      </StackItem>
    </Stack>
  );
};
export default ModelVersionDetailsView;
