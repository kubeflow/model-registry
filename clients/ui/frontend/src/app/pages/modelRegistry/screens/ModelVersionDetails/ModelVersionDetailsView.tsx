import * as React from 'react';
import {
  DescriptionList,
  Divider,
  Flex,
  FlexItem,
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
} from '@patternfly/react-core';
import {
  EditableLabelsDescriptionListGroup,
  EditableTextDescriptionListGroup,
  DashboardDescriptionListGroup,
  InlineTruncatedClipboardCopy,
} from 'mod-arch-shared';
import { ModelVersion, ModelArtifactList } from '~/app/types';
import { ModelRegistryContext } from '~/app/context/ModelRegistryContext';
import { getLabels, mergeUpdatedLabels } from '~/app/pages/modelRegistry/screens/utils';
import ModelPropertiesDescriptionListGroup from '~/app/pages/modelRegistry/screens/ModelPropertiesDescriptionListGroup';
import ModelTimestamp from '~/app/pages/modelRegistry/screens/components/ModelTimestamp';
import { bumpBothTimestamps, bumpRegisteredModelTimestamp } from '~/app/api/updateTimestamps';
import { uriToStorageFields } from '~/app/utils';
import useRegisteredModelById from '~/app/hooks/useRegisteredModelById';
import ModelDetailsCard from '~/app/pages/modelRegistry/screens/ModelVersions/ModelDetailsCard';

type ModelVersionDetailsViewProps = {
  modelVersion: ModelVersion;
  isArchiveVersion?: boolean;
  refresh: () => void;
  modelArtifacts: ModelArtifactList;
};

const ModelVersionDetailsView: React.FC<ModelVersionDetailsViewProps> = ({
  modelVersion: mv,
  isArchiveVersion,
  refresh,
  modelArtifacts,
}) => {
  const modelArtifact = modelArtifacts.items.length ? modelArtifacts.items[0] : null;
  const { apiState } = React.useContext(ModelRegistryContext);
  const storageFields = uriToStorageFields(modelArtifact?.uri || '');
  const [registeredModel, registeredModelLoaded, registeredModelLoadError] = useRegisteredModelById(
    mv.registeredModelId,
  );

  const loaded = registeredModelLoaded;
  const loadError = registeredModelLoadError;

  if (!loaded) {
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
            <Flex
              direction={{ default: 'column', md: 'row' }}
              columnGap={{ default: 'columnGap4xl' }}
              rowGap={{ default: 'rowGapLg' }}
            >
              <FlexItem flex={{ default: 'flex_1' }}>
                <DescriptionList isFillColumns>
                  <EditableLabelsDescriptionListGroup
                    labels={getLabels(mv.customProperties)}
                    isArchive={isArchiveVersion}
                    allExistingKeys={Object.keys(mv.customProperties)}
                    title="Labels"
                    contentWhenEmpty="No labels"
                    labelProps={{ variant: 'outline' }}
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
              </FlexItem>
              <Divider orientation={{ default: 'vertical' }} />
              <FlexItem flex={{ default: 'flex_1' }}>
                <Title style={{ margin: '1em 0' }} headingLevel={ContentVariants.h3}>
                  Model location
                </Title>
                {loadError ? (
                  <Alert variant="danger" isInline title={loadError.name}>
                    {loadError.message}
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
              </FlexItem>
            </Flex>
          </CardBody>
        </Card>
      </StackItem>
    </Stack>
  );
};
export default ModelVersionDetailsView;
