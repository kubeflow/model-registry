import React from 'react';
import { ActionsColumn, Td, Tr } from '@patternfly/react-table';
import { Button, Tooltip } from '@patternfly/react-core';
import { useNavigate } from 'react-router';
import {
  ModelRegistryKind,
  PlatformMode,
  ResourceNameTooltip,
  useModularArchContext,
} from 'mod-arch-shared';
import { ModelRegistryTableRowStatus } from './ModelRegistryTableRowStatus';

type ModelRegistriesTableRowProps = {
  modelRegistry: ModelRegistryKind;
  // roleBindings: ContextResourceData<RoleBindingKind>; // TODO: [Midstream] Filter role bindings for this model registry
  onEditRegistry: (obj: ModelRegistryKind) => void;
  onDeleteRegistry: (obj: ModelRegistryKind) => void;
};

const ModelRegistriesTableRow: React.FC<ModelRegistriesTableRowProps> = ({
  modelRegistry: mr,
  // roleBindings, // TODO: [Midstream] Filter role bindings for this model registry
  onEditRegistry,
  onDeleteRegistry,
}) => {
  const navigate = useNavigate();
  const { platformMode } = useModularArchContext();
  const isPlatformKubeflow = platformMode === PlatformMode.Kubeflow;
  const filteredRoleBindings = []; // TODO: [Midstream] Filter role bindings for this model registry

  return (
    <Tr>
      <Td dataLabel="Model registry name">
        <ResourceNameTooltip resource={mr}>
          <strong>{mr.metadata.displayName || mr.metadata.name}</strong>
        </ResourceNameTooltip>
        {mr.metadata.description && <p>{mr.metadata.description}</p>}
      </Td>
      <Td dataLabel="Status">
        <ModelRegistryTableRowStatus
          conditions={[
            {
              type: 'Available',
              status: 'True',
              reason: 'Ready',
              message: 'Model registry is ready.',
            },
          ]}
        />
      </Td>
      {!isPlatformKubeflow && (
        <Td modifier="fitContent">
          {filteredRoleBindings.length === 0 ? (
            <Tooltip content="You can manage permissions when the model registry becomes available.">
              <Button isAriaDisabled variant="link">
                Manage permissions
              </Button>
            </Tooltip>
          ) : (
            <Button
              variant="link"
              onClick={() => navigate(`/model-registry-settings/permissions/${mr.metadata.name}`)}
            >
              Manage permissions
            </Button>
          )}
        </Td>
      )}
      {!isPlatformKubeflow && (
        <Td isActionCell>
          <ActionsColumn
            disabled={isPlatformKubeflow}
            items={[
              {
                title: 'Edit model registry',
                disabled: isPlatformKubeflow,
                onClick: () => {
                  onEditRegistry(mr);
                },
              },
              {
                title: 'Delete model registry',
                disabled: isPlatformKubeflow,
                onClick: () => {
                  onDeleteRegistry(mr);
                },
              },
            ]}
          />
        </Td>
      )}
    </Tr>
  );
};

export default ModelRegistriesTableRow;
