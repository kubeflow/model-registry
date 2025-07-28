import { RoleBindingKind, K8sStatus } from 'mod-arch-shared';
import { createRoleBinding, deleteRoleBinding } from '~/app/api/k8s';

export const createModelRegistryRoleBindingWrapper = async (
  roleBinding: RoleBindingKind,
): Promise<RoleBindingKind> => {
  const hostPath = window.location.origin;
  return createRoleBinding(hostPath, {})({}, roleBinding);
};

export const deleteModelRegistryRoleBindingWrapper = async (
  name: string,
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  namespace: string,
): Promise<K8sStatus> => {
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

export const createModelRegistryProjectRoleBinding = async (
  roleBinding: RoleBindingKind,
): Promise<RoleBindingKind> => {
  const hostPath = window.location.origin;
  // Add project-specific labels
  const projectRoleBinding = {
    ...roleBinding,
    metadata: {
      ...roleBinding.metadata,
      labels: {
        ...roleBinding.metadata.labels,
        'app.kubernetes.io/component': 'model-registry-project-rbac',
      },
    },
  };
  return createRoleBinding(hostPath, {})({}, projectRoleBinding);
};

export const deleteModelRegistryProjectRoleBinding = async (
  name: string,
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  namespace: string,
): Promise<K8sStatus> => {
  const hostPath = window.location.origin;
  await deleteRoleBinding(hostPath, {})({}, name);
  return {
    apiVersion: 'v1',
    kind: 'Status',
    status: 'Success' as const,
    code: 200,
    message: 'Project role binding deleted successfully',
    reason: 'Deleted',
  };
};
