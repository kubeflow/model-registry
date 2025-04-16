import { capitalize } from '@patternfly/react-core';
import { ProjectKind, RoleBindingKind } from '~/k8sTypes';
import { namespaceToProjectDisplayName } from '~/concepts/projects/utils';
import { RoleBindingPermissionsRBType, RoleBindingPermissionsRoleType } from './types';

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
  return RoleBindingPermissionsRoleType.DEFAULT;
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