import * as React from 'react';
import { DescriptionList, Flex, FlexItem, ContentVariants, Title } from '@patternfly/react-core';
import DashboardDescriptionListGroup from '~/shared/components/DashboardDescriptionListGroup';
import EditableTextDescriptionListGroup from '~/shared/components/EditableTextDescriptionListGroup';
import EditableLabelsDescriptionListGroup from '~/shared/components/EditableLabelsDescriptionListGroup';
import { ModelVersion } from '~/app/types';
import useModelArtifactsByVersionId from '~/app/hooks/useModelArtifactsByVersionId';
import { ModelRegistryContext } from '~/app/context/ModelRegistryContext';
import InlineTruncatedClipboardCopy from '~/shared/components/InlineTruncatedClipboardCopy';
import DashboardHelpTooltip from '~/shared/components/DashboardHelpTooltip';
import {
  getLabels,
  mergeUpdatedLabels,
  uriToObjectStorageFields,
} from '~/app/pages/modelRegistry/screens/utils';
import ModelPropertiesDescriptionListGroup from '~/app/pages/modelRegistry/screens/ModelPropertiesDescriptionListGroup';
import ModelTimestamp from '~/app/pages/modelRegistry/screens/components/ModelTimestamp';

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
  const [modelArtifact] = useModelArtifactsByVersionId(mv.id);
  const { apiState } = React.useContext(ModelRegistryContext);
  const storageFields = uriToObjectStorageFields(modelArtifact.items[0]?.uri || '');

  return (
    <Flex
      direction={{ default: 'column', md: 'row' }}
      columnGap={{ default: 'columnGap4xl' }}
      rowGap={{ default: 'rowGapLg' }}
    >
      <FlexItem flex={{ default: 'flex_1' }}>
        <DescriptionList isFillColumns>
          <EditableTextDescriptionListGroup
            testid="model-version-description"
            isArchive={isArchiveVersion}
            title="Description"
            contentWhenEmpty="No description"
            value={mv.description || ''}
            saveEditedValue={(value) =>
              apiState.api
                .patchModelVersion(
                  {},
                  {
                    description: value,
                  },
                  mv.id,
                )
                .then(refresh)
            }
          />
          <EditableLabelsDescriptionListGroup
            labels={getLabels(mv.customProperties)}
            isArchive={isArchiveVersion}
            allExistingKeys={Object.keys(mv.customProperties)}
            saveEditedLabels={(editedLabels) =>
              apiState.api
                .patchModelVersion(
                  {},
                  {
                    customProperties: mergeUpdatedLabels(mv.customProperties, editedLabels),
                  },
                  mv.id,
                )
                .then(refresh)
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
        <Title style={{ marginTop: '1em' }} headingLevel={ContentVariants.h3}>
          Model location
        </Title>
        <DescriptionList isFillColumns>
          {storageFields && (
            <>
              <DashboardDescriptionListGroup
                title="Endpoint"
                isEmpty={modelArtifact.size === 0 || !storageFields.endpoint}
                contentWhenEmpty="No endpoint"
              >
                <InlineTruncatedClipboardCopy
                  testId="storage-endpoint"
                  textToCopy={storageFields.endpoint}
                />
              </DashboardDescriptionListGroup>
              <DashboardDescriptionListGroup
                title="Region"
                isEmpty={modelArtifact.size === 0 || !storageFields.region}
                contentWhenEmpty="No region"
              >
                <InlineTruncatedClipboardCopy
                  testId="storage-region"
                  textToCopy={storageFields.region || ''}
                />
              </DashboardDescriptionListGroup>
              <DashboardDescriptionListGroup
                title="Bucket"
                isEmpty={modelArtifact.size === 0 || !storageFields.bucket}
                contentWhenEmpty="No bucket"
              >
                <InlineTruncatedClipboardCopy
                  testId="storage-bucket"
                  textToCopy={storageFields.bucket}
                />
              </DashboardDescriptionListGroup>
              <DashboardDescriptionListGroup
                title="Path"
                isEmpty={modelArtifact.size === 0 || !storageFields.path}
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
                isEmpty={modelArtifact.size === 0 || !modelArtifact.items[0].uri}
                contentWhenEmpty="No URI"
              >
                <InlineTruncatedClipboardCopy
                  testId="storage-uri"
                  textToCopy={modelArtifact.items[0]?.uri || ''}
                />
              </DashboardDescriptionListGroup>
            </>
          )}
          <DashboardDescriptionListGroup
            title="Source model format"
            isEmpty={modelArtifact.size === 0 || !modelArtifact.items[0].modelFormatName}
            contentWhenEmpty="No source model format"
          >
            {modelArtifact.items[0]?.modelFormatName}
          </DashboardDescriptionListGroup>
          <DashboardDescriptionListGroup
            title="Author"
            tooltip={
              <DashboardHelpTooltip content="The author is the user who registered the model version." />
            }
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
