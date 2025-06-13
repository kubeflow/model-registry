import { capitalize } from '@patternfly/react-core';
import { ProjectKind, RoleBindingKind } from '~/app/k8sTypes';
import { namespaceToProjectDisplayName } from '~/app/concepts/projects/utils';
import { patchRoleBindingSubjects } from '~/app/api/k8s/roleBindings';
import { RoleBindingPermissionsRBType, RoleBindingPermissionsRoleType } from '~/app/pages/settings/roleBinding/types';

export const filterRoleBindingSubjects = (
  roleBindings: RoleBindingKind[],
  type: RoleBindingPermissionsRBType,
  isProjectSubject?: boolean,
): RoleBindingKind[] =>
  roleBindings.filter(
    (roles) =>
      roles.subjects[0]?.kind === type &&
      (isProjectSubject
        ? roles.metadata.labels?.['opendatahub.io/rb-project-subject'] === 'true'
        : !(roles.metadata.labels?.['opendatahub.io/rb-project-subject'] === 'true')),
  );

export const castRoleBindingPermissionsRoleType = (
  role: string,
): RoleBindingPermissionsRoleType => {
  if (role === RoleBindingPermissionsRoleType.ADMIN) {
    return RoleBindingPermissionsRoleType.ADMIN;
  }
  if (role === RoleBindingPermissionsRoleType.EDIT) {
    return RoleBindingPermissionsRoleType.EDIT;
  }
  if (role.includes('registry-user')) {
    return RoleBindingPermissionsRoleType.DEFAULT;
  }
  return RoleBindingPermissionsRoleType.CUSTOM;
};

export const firstSubject = (
  roleBinding: RoleBindingKind,
  isProjectSubject?: boolean,
  project?: ProjectKind[],
): string =>
  (isProjectSubject && project
    ? namespaceToProjectDisplayName(
        roleBinding.subjects[0]?.name.replace(/^system:serviceaccounts:/, ''),
        project,
      )
    : roleBinding.subjects[0]?.name) || '';

export const roleLabel = (value: RoleBindingPermissionsRoleType): string => {
  if (value === RoleBindingPermissionsRoleType.EDIT) {
    return 'Contributor';
  }
  return capitalize(value);
};

export const removePrefix = (roleBindings: RoleBindingKind[]): string[] =>
  roleBindings.map((rb) => rb.subjects[0]?.name.replace(/^system:serviceaccounts:/, ''));

export const isCurrentUserChanging = (
  roleBinding: RoleBindingKind | undefined,
  currentUsername: string,
): boolean => {
  if (!roleBinding) {
    return false;
  }
  return currentUsername === roleBinding.subjects[0].name;
};

export const tryPatchRoleBinding = async (
  oldRBObject: RoleBindingKind,
  newRBObject: RoleBindingKind,
): Promise<boolean> => {
  // Trying to patch roleRef will always fail
  if (oldRBObject.roleRef.name !== newRBObject.roleRef.name) {
    return false;
  }
  try {
    await patchRoleBindingSubjects(
      oldRBObject.metadata.name,
      oldRBObject.metadata.namespace,
      newRBObject.subjects,
      { dryRun: true },
    );
  } catch (e) {
    return false;
  }
  try {
    await patchRoleBindingSubjects(
      oldRBObject.metadata.name,
      oldRBObject.metadata.namespace,
      newRBObject.subjects,
      { dryRun: false },
    );
    return true;
  } catch {
    return false;
  }
}; 