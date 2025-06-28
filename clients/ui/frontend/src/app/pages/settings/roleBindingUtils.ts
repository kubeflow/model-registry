import { RoleBindingKind, K8sStatus } from 'mod-arch-shared';
import { createRoleBinding, deleteRoleBinding } from '~/app/api/k8s';

export const createModelRegistryRoleBindingWrapper = async (
  roleBinding: RoleBindingKind,
): Promise<RoleBindingKind> => {
  const hostPath = window.location.origin;
  return createRoleBinding(hostPath, {})({}, roleBinding);
};

export const deleteModelRegistryRoleBindingWrapper = async (name: string): Promise<K8sStatus> => {
  const hostPath = window.location.origin;
  await deleteRoleBinding(hostPath, {})({}, name);
  return {
    apiVersion: 'v1',
    kind: 'Status',
    status: 'Success' as const,
    code: 200,
    message: 'Role binding deleted successfully',
    reason: 'Deleted',
  };
};
