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
} from '@patternfly/react-core';
import DashboardDescriptionListGroup from '~/shared/components/DashboardDescriptionListGroup';
import EditableTextDescriptionListGroup from '~/shared/components/EditableTextDescriptionListGroup';
import { EditableLabelsDescriptionListGroup } from '~/shared/components/EditableLabelsDescriptionListGroup';
import { ModelVersion } from '~/app/types';
import useModelArtifactsByVersionId from '~/app/hooks/useModelArtifactsByVersionId';
import { ModelRegistryContext } from '~/app/context/ModelRegistryContext';
import InlineTruncatedClipboardCopy from '~/shared/components/InlineTruncatedClipboardCopy';
import { getLabels, mergeUpdatedLabels } from '~/app/pages/modelRegistry/screens/utils';
import { uriToObjectStorageFields } from '~/app/utils';
import ModelPropertiesDescriptionListGroup from '~/app/pages/modelRegistry/screens/ModelPropertiesDescriptionListGroup';
import ModelTimestamp from '~/app/pages/modelRegistry/screens/components/ModelTimestamp';
import { bumpBothTimestamps, bumpRegisteredModelTimestamp } from '~/app/utils/updateTimestamps';

type ModelVersionDetailsViewProps = {
  modelVersion: ModelVersion;
  isArchiveVersion?: boolean;
  refresh: () => void;
};

const ModelVersionDetailsView: React.FC<ModelVersionDetailsViewProps> = ({
  modelVersion: mv,
  isArchiveVersion,
  refresh,
}) => {
  const [modelArtifacts, modelArtifactsLoaded, modelArtifactsLoadError, refreshModelArtifacts] =
    useModelArtifactsByVersionId(mv.id);

  const modelArtifact = modelArtifacts.items.length ? modelArtifacts.items[0] : null;
  const { apiState } = React.useContext(ModelRegistryContext);
  const storageFields = uriToObjectStorageFields(modelArtifact?.uri || '');

  if (!modelArtifactsLoaded) {
    return (
      <Bullseye>
        <Spinner size="xl" />
      </Bullseye>
    );
  }
  const handleVersionUpdate = async (updatePromise: Promise<unknown>): Promise<void> => {
    await updatePromise;

    if (!mv.registeredModelId) {
      return;
    }

    await bumpRegisteredModelTimestamp(apiState.api, mv.registeredModelId);
    refresh();
  };

  const handleArtifactUpdate = async (updatePromise: Promise<unknown>): Promise<void> => {
    try {
      await updatePromise;
      await bumpBothTimestamps(apiState.api, mv.id, mv.registeredModelId);
      refreshModelArtifacts();
    } catch (error) {
      throw new Error(
        `Failed to update artifact: ${error instanceof Error ? error.message : String(error)}`,
      );
    }
  };

  return (
    <Flex
      direction={{ default: 'column', md: 'row' }}
      columnGap={{ default: 'columnGap4xl' }}
      rowGap={{ default: 'rowGapLg' }}
    >
      <FlexItem flex={{ default: 'flex_1' }}>
        <DescriptionList isFillColumns>
          <EditableTextDescriptionListGroup
            editableVariant="TextArea"
            baseTestId="model-version-description"
            isArchive={isArchiveVersion}
            title="Description"
            contentWhenEmpty="No description"
            value={mv.description || ''}
            saveEditedValue={(value) =>
              handleVersionUpdate(apiState.api.patchModelVersion({}, { description: value }, mv.id))
            }
          />
          <EditableLabelsDescriptionListGroup
            labels={getLabels(mv.customProperties)}
            isArchive={isArchiveVersion}
            allExistingKeys={Object.keys(mv.customProperties)}
            title="Labels"
            contentWhenEmpty="No labels"
            onLabelsChange={(editedLabels) =>
              handleVersionUpdate(
                apiState.api.patchModelVersion(
                  {},
                  { customProperties: mergeUpdatedLabels(mv.customProperties, editedLabels) },
                  mv.id,
                ),
              )
            }
            data-testid="model-version-labels"
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
      <FlexItem flex={{ default: 'flex_1' }}>
        <DescriptionList isFillColumns>
          <DashboardDescriptionListGroup
            title="Version ID"
            isEmpty={!mv.id}
            contentWhenEmpty="No model ID"
          >
            <InlineTruncatedClipboardCopy testId="model-version-id" textToCopy={mv.id} />
          </DashboardDescriptionListGroup>
        </DescriptionList>
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
              {storageFields && (
                <>
                  <DashboardDescriptionListGroup
                    title="Endpoint"
                    isEmpty={modelArtifacts.size === 0 || !storageFields.endpoint}
                    contentWhenEmpty="No endpoint"
                  >
                    <InlineTruncatedClipboardCopy
                      testId="storage-endpoint"
                      textToCopy={storageFields.endpoint}
                    />
                  </DashboardDescriptionListGroup>
                  <DashboardDescriptionListGroup
                    title="Region"
                    isEmpty={modelArtifacts.size === 0 || !storageFields.region}
                    contentWhenEmpty="No region"
                  >
                    <InlineTruncatedClipboardCopy
                      testId="storage-region"
                      textToCopy={storageFields.region || ''}
                    />
                  </DashboardDescriptionListGroup>
                  <DashboardDescriptionListGroup
                    title="Bucket"
                    isEmpty={modelArtifacts.size === 0 || !storageFields.bucket}
                    contentWhenEmpty="No bucket"
                  >
                    <InlineTruncatedClipboardCopy
                      testId="storage-bucket"
                      textToCopy={storageFields.bucket}
                    />
                  </DashboardDescriptionListGroup>
                  <DashboardDescriptionListGroup
                    title="Path"
                    isEmpty={modelArtifacts.size === 0 || !storageFields.path}
                    contentWhenEmpty="No path"
                  >
                    <InlineTruncatedClipboardCopy
                      testId="storage-path"
                      textToCopy={storageFields.path}
                    />
                  </DashboardDescriptionListGroup>
                </>
              )}
              {!storageFields && (
                <>
                  <DashboardDescriptionListGroup
                    title="URI"
                    isEmpty={modelArtifacts.size === 0 || !modelArtifact?.uri}
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
            <Title style={{ margin: '1em 0' }} headingLevel={ContentVariants.h3}>
              Source model format
            </Title>
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
                title="Model Format"
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
                title="Version"
                contentWhenEmpty="No source model format version"
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
            title="Last modified at"
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
  );
};
export default ModelVersionDetailsView;
