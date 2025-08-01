import { capitalize } from '@patternfly/react-core';
import { RoleBindingKind } from 'mod-arch-shared';
import { patchRoleBinding } from '~/app/api/k8s';
import { getDisplayNameFromK8sResource } from '~/app/shared/components/utils';
import { ProjectKind, NamespaceKind } from '~/app/shared/components/types';
import { RoleBindingPermissionsRBType, RoleBindingPermissionsRoleType } from './types';

export const filterRoleBindingSubjects = (
  roleBindings: RoleBindingKind[],
  type: RoleBindingPermissionsRBType,
): RoleBindingKind[] =>
  roleBindings.filter(
    (roles) =>
      roles.subjects[0]?.kind === type &&
      !(roles.metadata.labels?.['opendatahub.io/rb-project-subject'] === 'true'),
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
    await patchRoleBinding('', { namespace: oldRBObject.metadata.namespace, dryRun: true })(
      {},
      newRBObject,
      oldRBObject.metadata.name,
    );
  } catch {
    return false;
  }
  try {
    // Actual patch
    await patchRoleBinding('', { namespace: oldRBObject.metadata.namespace, dryRun: false })(
      {},
      newRBObject,
      oldRBObject.metadata.name,
    );
    return true;
  } catch {
    return false;
  }
};

export const namespaceToProjectDisplayName = (
  namespace: string,
  projects: ProjectKind[],
): string => {
  const project = projects.find((p) => p.metadata.name === namespace);
  return project ? getDisplayNameFromK8sResource(project) : namespace;
};

export const projectDisplayNameToNamespace = (
  displayName: string,
  projects: ProjectKind[],
): string => {
  const project = projects.find(
    (p) => p.metadata.annotations?.['openshift.io/display-name'] === displayName,
  );
  return project?.metadata.name || displayName;
};

// New utility functions for NamespaceKind
/**
 * Get the display name for a namespace.
 * @param namespaceName The name of the namespace
 * @param namespaces Array of NamespaceKind objects
 * @returns The display name or namespace name if not found
 */
export const namespaceToDisplayName = (
  namespaceName: string,
  namespaces: NamespaceKind[],
): string => namespaces.find((ns) => ns.name === namespaceName)?.displayName || namespaceName;

/**
 * Find a namespace by its display name.
 * @param displayName The display name to search for
 * @param namespaces Array of NamespaceKind objects
 * @returns The namespace name or the display name if not found
 */
export const displayNameToNamespace = (displayName: string, namespaces: NamespaceKind[]): string =>
  namespaces.find((ns) => ns.displayName === displayName)?.name || displayName;
