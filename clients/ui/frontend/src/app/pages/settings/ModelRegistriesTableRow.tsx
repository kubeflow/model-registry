import React from 'react';
import { ActionsColumn, Td, Tr } from '@patternfly/react-table';
import { Button, Tooltip } from '@patternfly/react-core';
import { useNavigate } from 'react-router';
import { ModelRegistry } from '~/app/types';
import ResourceNameTooltip from '~/shared/components/ResourceNameTooltip';
import { convertToK8sResourceCommon } from '~/app/utils';
import { isPlatformDefault } from '~/shared/utilities/const';
import { ModelRegistryTableRowStatus } from './ModelRegistryTableRowStatus';

type ModelRegistriesTableRowProps = {
  modelRegistry: ModelRegistry;
  // roleBindings: ContextResourceData<RoleBindingKind>; // TODO: [Midstream] Filter role bindings for this model registry
  onEditRegistry: (obj: ModelRegistry) => void;
  onDeleteRegistry: (obj: ModelRegistry) => void;
};

const ModelRegistriesTableRow: React.FC<ModelRegistriesTableRowProps> = ({
  modelRegistry: mr,
  // roleBindings, // TODO: [Midstream] Filter role bindings for this model registry
  onEditRegistry,
  onDeleteRegistry,
}) => {
  const navigate = useNavigate();
  const filteredRoleBindings = []; // TODO: [Midstream] Filter role bindings for this model registry

  return (
    <Tr>
      <Td dataLabel="Model registry name">
        <ResourceNameTooltip resource={convertToK8sResourceCommon(mr)}>
          <strong>{mr.displayName || mr.name}</strong>
        </ResourceNameTooltip>
        {mr.description && <p>{mr.description}</p>}
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
      {isPlatformDefault() && (
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
              onClick={() => navigate(`/model-registry-settings/permissions/${mr.name}`)}
            >
              Manage permissions
            </Button>
          )}
        </Td>
      )}
      {isPlatformDefault() && (
        <Td isActionCell>
          <ActionsColumn
            disabled={!isPlatformDefault()}
            items={[
              {
                title: 'Edit model registry',
                disabled: !isPlatformDefault(),
                onClick: () => {
                  onEditRegistry(mr);
                },
              },
              {
                title: 'Delete model registry',
                disabled: !isPlatformDefault(),
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
