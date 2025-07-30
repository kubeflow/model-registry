import * as React from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { ActionsColumn, IAction, Td, Tr } from '@patternfly/react-table';
import { Content, ContentVariants, FlexItem, Truncate } from '@patternfly/react-core';
import { ModelRegistryContext } from '~/app/context/ModelRegistryContext';
import { ModelState, RegisteredModel } from '~/app/types';
import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';
import ModelLabels from '~/app/pages/modelRegistry/screens/components/ModelLabels';
import ModelTimestamp from '~/app/pages/modelRegistry/screens/components/ModelTimestamp';
import { ArchiveRegisteredModelModal } from '~/app/pages/modelRegistry/screens/components/ArchiveRegisteredModelModal';
import { RestoreRegisteredModelModal } from '~/app/pages/modelRegistry/screens/components/RestoreRegisteredModel';
import {
  registeredModelArchiveDetailsUrl,
  registeredModelArchiveUrl,
  registeredModelUrl,
} from '~/app/pages/modelRegistry/screens/routeUtils';
import { ModelVersionsTab } from '~/app/pages/modelRegistry/screens/ModelVersions/const';

type RegisteredModelTableRowProps = {
  registeredModel: RegisteredModel;
  isArchiveRow?: boolean;
  hasDeploys?: boolean;
  refresh: () => void;
};

const RegisteredModelTableRow: React.FC<RegisteredModelTableRowProps> = ({
  registeredModel: rm,
  isArchiveRow,
  hasDeploys = false,
  refresh,
}) => {
  const { apiState } = React.useContext(ModelRegistryContext);
  const navigate = useNavigate();
  const { preferredModelRegistry } = React.useContext(ModelRegistrySelectorContext);
  const [isArchiveModalOpen, setIsArchiveModalOpen] = React.useState(false);
  const [isRestoreModalOpen, setIsRestoreModalOpen] = React.useState(false);
  const rmUrl = registeredModelUrl(rm.id, preferredModelRegistry?.name);

  const actions: IAction[] = [
    {
      title: 'Overview',
      onClick: () => {
        navigate(
          isArchiveRow
            ? `${registeredModelArchiveUrl(preferredModelRegistry?.name)}/${rm.id}/${
                ModelVersionsTab.OVERVIEW
              }`
            : `${rmUrl}/${ModelVersionsTab.OVERVIEW}`,
        );
      },
    },
    ...(isArchiveRow
      ? [
          {
            title: 'Restore model',
            onClick: () => setIsRestoreModalOpen(true),
          },
        ]
      : [
          { isSeparator: true },
          {
            title: 'Archive model',
            onClick: () => setIsArchiveModalOpen(true),
            isAriaDisabled: hasDeploys,
            tooltipProps: hasDeploys
              ? { content: 'Models with deployed versions cannot be archived.' }
              : undefined,
          },
        ]),
  ];

  return (
    <Tr>
      <Td dataLabel="Model name">
        <div id="model-name" data-testid="model-name">
          <FlexItem>
            <Link
              to={
                isArchiveRow
                  ? registeredModelArchiveDetailsUrl(rm.id, preferredModelRegistry?.name)
                  : rmUrl
              }
            >
              <Truncate content={rm.name} />
            </Link>
          </FlexItem>
        </div>
        {rm.description && (
          <Content data-testid="description" component={ContentVariants.small}>
            <Truncate content={rm.description} />
          </Content>
        )}
      </Td>
      <Td dataLabel="Labels">
        <ModelLabels customProperties={rm.customProperties} name={rm.name} />
      </Td>
      <Td dataLabel="Last modified">
        <ModelTimestamp timeSinceEpoch={rm.lastUpdateTimeSinceEpoch} />
      </Td>
      <Td dataLabel="Owner">
        <Content component="p" data-testid="registered-model-owner">
          {rm.owner || '-'}
        </Content>
      </Td>
      <Td isActionCell>
        <ActionsColumn items={actions} />
        {isArchiveModalOpen ? (
          <ArchiveRegisteredModelModal
            onCancel={() => setIsArchiveModalOpen(false)}
            onSubmit={() =>
              apiState.api
                .patchRegisteredModel(
                  {},
                  {
                    state: ModelState.ARCHIVED,
                  },
                  rm.id,
                )
                .then(refresh)
            }
            registeredModelName={rm.name}
          />
        ) : null}
        {isRestoreModalOpen ? (
          <RestoreRegisteredModelModal
            onCancel={() => setIsRestoreModalOpen(false)}
            onSubmit={() =>
              apiState.api
                .patchRegisteredModel(
                  {},
                  {
                    state: ModelState.LIVE,
                  },
                  rm.id,
                )
                .then(() => navigate(registeredModelUrl(rm.id, preferredModelRegistry?.name)))
            }
            registeredModelName={rm.name}
          />
        ) : null}
      </Td>
    </Tr>
  );
};

export default RegisteredModelTableRow;
