import React from 'react';
import { ActionsColumn, Td, Tr } from '@patternfly/react-table';
import { useNavigate } from 'react-router-dom';
import { Button, Tooltip } from '@patternfly/react-core';
import { FetchStateObject } from 'mod-arch-shared/dist/types/common';
import { ModelRegistryKind, ResourceNameTooltip, RoleBindingKind } from 'mod-arch-shared';
import { DeploymentMode, useModularArchContext } from 'mod-arch-core';
import { ModelRegistryTableRowStatus } from './ModelRegistryTableRowStatus';

type ModelRegistriesTableRowProps = {
  modelRegistry: ModelRegistryKind;
  roleBindings: FetchStateObject<RoleBindingKind[]>;
  onEditRegistry: (obj: ModelRegistryKind) => void;
  onDeleteRegistry: (obj: ModelRegistryKind) => void;
};

const ModelRegistriesTableRow: React.FC<ModelRegistriesTableRowProps> = ({
  modelRegistry: mr,
  roleBindings,
  onEditRegistry,
  onDeleteRegistry,
}) => {
  const navigate = useNavigate();
  const { config } = useModularArchContext();
  const { deploymentMode } = config;
  const isDeploymentKubeflow = deploymentMode === DeploymentMode.Kubeflow;
  const filteredRoleBindings = roleBindings.data.filter(
    (rb) =>
      rb.metadata.labels?.['app.kubernetes.io/name'] ===
      (mr.metadata.name || mr.metadata.annotations?.['openshift.io/display-name']),
  );

  return (
    <Tr>
      <Td dataLabel="Model registry name">
        <ResourceNameTooltip resource={mr}>
          <strong>
            {mr.metadata.annotations?.['openshift.io/display-name'] || mr.metadata.name}
          </strong>
        </ResourceNameTooltip>
        {mr.metadata.annotations?.['openshift.io/description'] && (
          <p>{mr.metadata.annotations['openshift.io/description']}</p>
        )}
      </Td>
      <Td dataLabel="Status">
        <ModelRegistryTableRowStatus conditions={mr.status?.conditions} />
      </Td>
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
      <Td isActionCell>
        <ActionsColumn
          items={[
            {
              title: 'Edit model registry',
              onClick: () => {
                onEditRegistry(mr);
              },
            },
            {
              title: 'Delete model registry',
              disabled: isDeploymentKubeflow,
              onClick: () => {
                onDeleteRegistry(mr);
              },
            },
          ]}
        />
      </Td>
    </Tr>
  );
};

export default ModelRegistriesTableRow;
