import * as React from 'react';
import {
  K8sResourceCommon,
  K8sStatus,
  RoleBindingKind,
  RoleBindingRoleRef,
  RoleBindingSubject,
  Table,
} from 'mod-arch-shared';
import { generateRoleBindingPermissions } from '~/app/api/k8s';
import RoleBindingPermissionsTableRow from './RoleBindingPermissionsTableRow';
import { columnsRoleBindingPermissions } from './data';
import { RoleBindingPermissionsRoleType } from './types';
import { firstSubject, tryPatchRoleBinding } from './utils';

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
  onError: (error: React.ReactNode) => void;
  refresh: () => void;
};

const RoleBindingPermissionsTable: React.FC<RoleBindingPermissionsTableProps> = ({
  ownerReference,
  subjectKind,
  namespace,
  roleRefKind,
  roleRefName,
  labels,
  isProjectSubject,
  defaultRoleBindingName,
  permissions,
  permissionOptions,
  isAdding,
  typeAhead,
  createRoleBinding,
  deleteRoleBinding,
  onDismissNewRow,
  onError,
  refresh,
}) => {
  const [editCell, setEditCell] = React.useState<string[]>([]);

  const createProjectRoleBinding = async (
    newRBObject: RoleBindingKind,
    oldRBObject?: RoleBindingKind,
  ) => {
    if (isAdding) {
      // Add new role binding
      createRoleBinding(newRBObject)
        .then(() => {
          onDismissNewRow();
          refresh();
        })
        .catch((e) => {
          onError(<>{e}</>);
        });
    } else if (oldRBObject) {
      const patchSucceeded = await tryPatchRoleBinding(oldRBObject, newRBObject);
      if (patchSucceeded) {
        setEditCell((prev) => prev.filter((cell) => cell !== oldRBObject.metadata.name));
        onDismissNewRow();
        refresh();
      } else {
        createRoleBinding(newRBObject)
          .then(() => {
            deleteRoleBinding(oldRBObject.metadata.name, oldRBObject.metadata.namespace)
              .then(() => refresh())
              .catch((e) => {
                onError(<>{e}</>);
                setEditCell((prev) => prev.filter((cell) => cell !== oldRBObject.metadata.name));
              });
          })
          .then(() => {
            onDismissNewRow();
            refresh();
          })
          .catch((e) => {
            onError(<>{e}</>);
            setEditCell((prev) => prev.filter((cell) => cell !== oldRBObject.metadata.name));
          });
      }
    }
  };
  return (
    <Table
      variant="compact"
      data={permissions}
      data-testid={`role-binding-table ${subjectKind}`}
      columns={columnsRoleBindingPermissions}
      disableRowRenderSupport
      footerRow={() =>
        isAdding ? (
          <RoleBindingPermissionsTableRow
            key="add-permissions-row"
            subjectKind={subjectKind}
            permissionOptions={permissionOptions}
            isProjectSubject={isProjectSubject}
            typeAhead={typeAhead}
            isEditing={false}
            isAdding
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
                .catch((e) => onError(e?.message || e || 'Unknown error'));
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
          isAdding={false}
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
            createProjectRoleBinding(newRBObject, rb);
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
            onDismissNewRow();
          }}
        />
      )}
    />
  );
};
export default RoleBindingPermissionsTable;
