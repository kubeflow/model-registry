import * as React from 'react';
import { K8sResourceCommon, K8sStatus } from '@openshift/dynamic-plugin-sdk-utils';
import {
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
} from '@mui/material';
import { RoleBindingKind, RoleBindingRoleRef, RoleBindingSubject } from '~/app/k8sTypes';
import { generateRoleBindingPermissions } from '~/app/api/k8s/roleBindings';
import RoleBindingPermissionsTableRow from './RoleBindingPermissionsTableRow';
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
    <TableContainer component={Paper}>
      <Table data-testid={`role-binding-table-${subjectKind}`}>
        <TableHead>
          <TableRow>
            <TableCell>User</TableCell>
            <TableCell>Permission</TableCell>
            <TableCell>Date Added</TableCell>
            <TableCell />
          </TableRow>
        </TableHead>
        <TableBody>
          {permissions.map((rb) => (
            <RoleBindingPermissionsTableRow
              key={rb.metadata.name || ''}
              permissionOptions={permissionOptions}
              roleBindingObject={rb}
              subjectKind={subjectKind}
              isEditing={
                firstSubject(rb) === '' || editCell.includes(rb.metadata.name)
              }
              isAdding={false}
              typeAhead={typeAhead}
              onChange={(subjectName: string, rbRoleRefName: RoleBindingPermissionsRoleType) => {
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
          ))}
          {isAdding && (
            <RoleBindingPermissionsTableRow
              subjectKind={subjectKind}
              permissionOptions={permissionOptions}
              isAdding
              onChange={(subjectName: string, rbRoleRefName: RoleBindingPermissionsRoleType) => {
                const newRBObject = generateRoleBindingPermissions(
                  namespace,
                  subjectKind,
                  subjectName,
                  roleRefName || rbRoleRefName,
                  roleRefKind,
                  labels,
                  ownerReference,
                );
                createProjectRoleBinding(newRBObject);
              }}
              onCancel={onDismissNewRow}
            />
          )}
        </TableBody>
      </Table>
    </TableContainer>
  );
};
export default RoleBindingPermissionsTable;
