import * as React from 'react';
import { ActionsColumn, Td, Tr } from '@patternfly/react-table';
import { Content, ContentVariants, Truncate, FlexItem } from '@patternfly/react-core';
import { Link, useNavigate } from 'react-router-dom';
import { ModelState, ModelVersion } from '~/app/types';
import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';
import { ModelRegistryContext } from '~/app/context/ModelRegistryContext';
import {
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
  refresh: () => void;
};

const ModelVersionsTableRow: React.FC<ModelVersionsTableRowProps> = ({
  modelVersion: mv,
  isArchiveRow,
  refresh,
}) => {
  const navigate = useNavigate();
  const { preferredModelRegistry } = React.useContext(ModelRegistrySelectorContext);
  const [isArchiveModalOpen, setIsArchiveModalOpen] = React.useState(false);
  const [isRestoreModalOpen, setIsRestoreModalOpen] = React.useState(false);
  const { apiState } = React.useContext(ModelRegistryContext);

  const actions = isArchiveRow
    ? [
        {
          title: 'Restore version',
          onClick: () => setIsRestoreModalOpen(true),
        },
      ]
    : [
        {
          title: 'Deploy',
          onClick: () => setIsDeployModalOpen(true),
        },
        {
          title: 'Archive model version',
          onClick: () => setIsArchiveModalOpen(true),
        },
      ];

  return (
    <Tr>
      <Td dataLabel="Version name">
        <div id="model-version-name" data-testid="model-version-name">
          <FlexItem>
            <Link
              to={
                isArchiveRow
                  ? modelVersionArchiveDetailsUrl(
                      mv.id,
                      mv.registeredModelId,
                      preferredModelRegistry?.name,
                    )
                  : modelVersionUrl(mv.id, mv.registeredModelId, preferredModelRegistry?.name)
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
      <Td isActionCell>
        <ActionsColumn items={actions} />
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
          isOpen={isArchiveModalOpen}
          modelVersionName={mv.name}
        />
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
                  modelVersionUrl(mv.id, mv.registeredModelId, preferredModelRegistry?.name),
                ),
              )
          }
          isOpen={isRestoreModalOpen}
          modelVersionName={mv.name}
        />
      </Td>
    </Tr>
  );
};

export default ModelVersionsTableRow;
