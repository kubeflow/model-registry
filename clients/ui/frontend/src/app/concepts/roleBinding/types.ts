import { RoleBindingSubject } from '~/app/k8sTypes';

export enum RoleBindingPermissionsRBType {
  USER = 'User',
  GROUP = 'Group',
}

export enum RoleBindingPermissionsRoleType {
  EDIT = 'edit',
  ADMIN = 'admin',
  DEFAULT = 'default',
  CUSTOM = 'custom',
}

export type RoleBindingSubjectWithRole = RoleBindingSubject & {
  role: RoleBindingPermissionsRoleType;
  roleBindingName: string;
  roleBindingNamespace: string;
}; 