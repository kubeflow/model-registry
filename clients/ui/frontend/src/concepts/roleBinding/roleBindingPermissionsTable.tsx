import * as React from 'react';
import { K8sResourceCommon, K8sStatus } from '@openshift/dynamic-plugin-sdk-utils';
import { Table } from '~/components/table';
import { RoleBindingKind, RoleBindingRoleRef, RoleBindingSubject } from '~/k8sTypes';
import { generateRoleBindingPermissions } from '~/api';
import RoleBindingPermissionsTableRow from './RoleBindingPermissionsTableRow';
import { columnsRoleBindingPermissions } from './data';
import { RoleBindingPermissionsRoleType } from './types';
import { firstSubject } from './utils';
import RoleBindingPermissionsTableRowAdd from './RoleBindingPermissionsTableRowAdd';

type RoleBindingPermissionsTableProps = {
  ownerReference?: K8sResourceCommon;
  subjectKind: RoleBindingSubject['kind'];
  namespace: string;
  roleRefKind: RoleBindingRoleRef['kind'];
  roleRefName?: RoleBindingRoleRef['name'];
  labels?: { [key: string]: string };
  isProjectSubject?: boolean;
  defaultRoleBindingName?: string;
  permissions: RoleBindingKind[];
  permissionOptions: {
    type: RoleBindingPermissionsRoleType;
    description: string;
  }[];
  isAdding: boolean;
  typeAhead?: string[];
  createRoleBinding: (roleBinding: RoleBindingKind) => Promise<RoleBindingKind>;
  deleteRoleBinding: (name: string, namespace: string) => Promise<K8sStatus>;
  onDismissNewRow: () => void;
  onError: (error: Error) => void;
  refresh: () => void;
};

const RoleBindingPermissionsTable: React.FC<RoleBindingPermissionsTableProps> = ({
  ownerReference,
  subjectKind,
  namespace,
  roleRefKind,
  roleRefName,
  labels,
  defaultRoleBindingName,
  permissions,
  permissionOptions,
  typeAhead,
  isProjectSubject,
  isAdding,
  createRoleBinding,
  deleteRoleBinding,
  onDismissNewRow,
  onError,
  refresh,
}) => {
  const [editCell, setEditCell] = React.useState<string[]>([]);

  return (
    <Table
      variant="compact"
      data={permissions}
      data-testid={`role-binding-table ${subjectKind}`}
      columns={columnsRoleBindingPermissions}
      disableRowRenderSupport
      footerRow={() =>
        isAdding ? (
          <RoleBindingPermissionsTableRowAdd
            key="add-permission-row"
            subjectKind={subjectKind}
            isProjectSubject={isProjectSubject}
            permissionOptions={permissionOptions}
            typeAhead={typeAhead}
            onChange={(subjectName, rbRoleRefName) => {
              const newRBObject = generateRoleBindingPermissions(
                namespace,
                subjectKind,
                subjectName,
                roleRefName || rbRoleRefName,
                roleRefKind,
                labels,
                ownerReference,
              );
              createRoleBinding(newRBObject)
                .then(() => {
                  onDismissNewRow();
                  refresh();
                })
                .catch((e) => {
                  onError(e);
                });
            }}
            onCancel={onDismissNewRow}
          />
        ) : null
      }
      rowRenderer={(rb) => (
        <RoleBindingPermissionsTableRow
          isProjectSubject={isProjectSubject}
          defaultRoleBindingName={defaultRoleBindingName}
          key={rb.metadata.name || ''}
          permissionOptions={permissionOptions}
          roleBindingObject={rb}
          subjectKind={subjectKind}
          isEditing={
            firstSubject(rb, isProjectSubject) === '' || editCell.includes(rb.metadata.name)
          }
          typeAhead={typeAhead}
          onChange={(subjectName, rbRoleRefName) => {
            const newRBObject = generateRoleBindingPermissions(
              namespace,
              subjectKind,
              subjectName,
              roleRefName || rbRoleRefName,
              roleRefKind,
              labels,
              ownerReference,
            );
            createRoleBinding(newRBObject)
              .then(() =>
                deleteRoleBinding(rb.metadata.name, rb.metadata.namespace)
                  .then(() => refresh())
                  .catch((e) => {
                    onError(e);
                    setEditCell((prev) => prev.filter((cell) => cell !== rb.metadata.name));
                  }),
              )
              .catch((e) => {
                onError(e);
                setEditCell((prev) => prev.filter((cell) => cell !== rb.metadata.name));
              });
            refresh();
          }}
          onDelete={() => {
            deleteRoleBinding(rb.metadata.name, rb.metadata.namespace).then(() => refresh());
          }}
          onEdit={() => {
            setEditCell((prev) => [...prev, rb.metadata.name]);
          }}
          onCancel={() => {
            setEditCell((prev) => prev.filter((cell) => cell !== rb.metadata.name));
          }}
        />
      )}
    />
  );
};
export default RoleBindingPermissionsTable;
