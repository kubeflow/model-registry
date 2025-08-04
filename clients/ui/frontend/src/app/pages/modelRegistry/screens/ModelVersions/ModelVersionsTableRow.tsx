import * as React from 'react';
import { ActionsColumn, IAction, Td, Tr } from '@patternfly/react-table';
import { Content, ContentVariants, Truncate, FlexItem } from '@patternfly/react-core';
import { Link, useNavigate } from 'react-router-dom';
import { ModelState, ModelVersion } from '~/app/types';
import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';
import { ModelRegistryContext } from '~/app/context/ModelRegistryContext';
import {
  archiveModelVersionDetailsUrl,
  modelVersionArchiveDetailsUrl,
  modelVersionUrl,
} from '~/app/pages/modelRegistry/screens/routeUtils';
import ModelTimestamp from '~/app/pages/modelRegistry/screens/components/ModelTimestamp';
import ModelLabels from '~/app/pages/modelRegistry/screens/components/ModelLabels';
import { ArchiveModelVersionModal } from '~/app/pages/modelRegistry/screens/components/ArchiveModelVersionModal';
import { RestoreModelVersionModal } from '~/app/pages/modelRegistry/screens/components/RestoreModelVersionModal';

type ModelVersionsTableRowProps = {
  modelVersion: ModelVersion;
  isArchiveRow?: boolean;
  isArchiveModel?: boolean;
  hasDeployment?: boolean;
  refresh: () => void;
};

const ModelVersionsTableRow: React.FC<ModelVersionsTableRowProps> = ({
  modelVersion: mv,
  isArchiveRow,
  isArchiveModel,
  hasDeployment = false,
  refresh,
}) => {
  const navigate = useNavigate();
  const { preferredModelRegistry } = React.useContext(ModelRegistrySelectorContext);
  const { apiState } = React.useContext(ModelRegistryContext);

  // TODO: Fetch model artifacts for when deploy functionality is enabled
  // const [modelArtifacts, modelArtifactsLoaded, modelArtifactsLoadError] =
  //   useModelArtifactsByVersionId(mv.id);

  const [isArchiveModalOpen, setIsArchiveModalOpen] = React.useState(false);
  const [isRestoreModalOpen, setIsRestoreModalOpen] = React.useState(false);

  if (!preferredModelRegistry) {
    return null;
  }

  const actions: IAction[] = isArchiveRow
    ? [
        {
          title: 'Restore model version',
          onClick: () => setIsRestoreModalOpen(true),
        },
      ]
    : [
        { isSeparator: true },
        {
          title: 'Archive model version',
          onClick: () => setIsArchiveModalOpen(true),
          isAriaDisabled: hasDeployment,
          tooltipProps: hasDeployment
            ? { content: 'Deployed versions cannot be archived' }
            : undefined,
        },
      ];

  return (
    <Tr>
      <Td dataLabel="Version name">
        <div id="model-version-name" data-testid="model-version-name">
          <FlexItem>
            <Link
              to={
                isArchiveModel
                  ? archiveModelVersionDetailsUrl(
                      mv.id,
                      mv.registeredModelId,
                      preferredModelRegistry.name,
                    )
                  : isArchiveRow
                    ? modelVersionArchiveDetailsUrl(
                        mv.id,
                        mv.registeredModelId,
                        preferredModelRegistry.name,
                      )
                    : modelVersionUrl(mv.id, mv.registeredModelId, preferredModelRegistry.name)
              }
            >
              <Truncate content={mv.name} />
            </Link>
          </FlexItem>
        </div>
        {mv.description && (
          <Content data-testid="model-version-description" component={ContentVariants.small}>
            <Truncate content={mv.description} />
          </Content>
        )}
      </Td>
      <Td dataLabel="Last modified">
        <ModelTimestamp timeSinceEpoch={mv.lastUpdateTimeSinceEpoch} />
      </Td>
      <Td dataLabel="Author">{mv.author}</Td>
      <Td dataLabel="Labels">
        <ModelLabels customProperties={mv.customProperties} name={mv.name} />
      </Td>
      {!isArchiveModel && (
        <Td isActionCell>
          <ActionsColumn items={actions} />
          {isArchiveModalOpen ? (
            <ArchiveModelVersionModal
              onCancel={() => setIsArchiveModalOpen(false)}
              onSubmit={() =>
                apiState.api
                  .patchModelVersion(
                    {},
                    {
                      state: ModelState.ARCHIVED,
                    },
                    mv.id,
                  )
                  .then(refresh)
              }
              modelVersionName={mv.name}
            />
          ) : null}
          {/* TODO: [Model Serving] Uncomment when model serving is available */}
          {/* NOTE: When uncommenting, pass modelArtifacts prop to avoid duplicate fetching */}
          {/* {isDeployModalOpen ? (
            <DeployRegisteredModelModal
              onSubmit={() => {
                navigate(
                  modelVersionDeploymentsUrl(
                    mv.id,
                    mv.registeredModelId,
                    preferredModelRegistry.metadata.name,
                  ),
                );
              }}
              onCancel={() => setIsDeployModalOpen(false)}
              modelVersion={mv}
            />
          ) : null} */}
          {isRestoreModalOpen ? (
            <RestoreModelVersionModal
              onCancel={() => setIsRestoreModalOpen(false)}
              onSubmit={() =>
                apiState.api
                  .patchModelVersion(
                    {},
                    {
                      state: ModelState.LIVE,
                    },
                    mv.id,
                  )
                  .then(() =>
                    navigate(
                      modelVersionUrl(mv.id, mv.registeredModelId, preferredModelRegistry.name),
                    ),
                  )
              }
              modelVersionName={mv.name}
            />
          ) : null}
        </Td>
      )}
    </Tr>
  );
};

export default ModelVersionsTableRow;
