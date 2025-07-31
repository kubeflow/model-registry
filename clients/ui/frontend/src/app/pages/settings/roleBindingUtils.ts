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

export const createModelRegistryNamespaceRoleBinding = async (
  roleBinding: RoleBindingKind,
): Promise<RoleBindingKind> => {
  const hostPath = window.location.origin;
  // Add namespace-specific labels
  const namespaceRoleBinding = {
    ...roleBinding,
    metadata: {
      ...roleBinding.metadata,
      labels: {
        ...roleBinding.metadata.labels,
        'app.kubernetes.io/component': 'model-registry-namespace-rbac',
      },
    },
  };
  return createRoleBinding(hostPath, {})({}, namespaceRoleBinding);
};

export const deleteModelRegistryNamespaceRoleBinding = async (
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
    message: 'Namespace role binding deleted successfully',
    reason: 'Deleted',
  };
};
