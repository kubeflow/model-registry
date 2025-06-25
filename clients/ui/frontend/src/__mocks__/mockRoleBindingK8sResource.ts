import { RoleBindingKind, RoleBindingSubject } from 'mod-arch-shared';
import { genUID } from './mockUtils';

export const mockRoleBindingK8sResource = ({
  name = 'test-name-view',
  namespace = 'test-project',
  subjects = [
    {
      kind: 'ServiceAccount',
      apiGroup: 'rbac.authorization.k8s.io',
      name: 'test-name-sa',
    },
  ],
  roleRefName = 'view',
  uid = genUID('rolebinding'),
}: {
  name?: string;
  namespace?: string;
  subjects?: RoleBindingSubject[];
  roleRefName?: string;
  uid?: string;
}): RoleBindingKind => ({
  kind: 'RoleBinding',
  apiVersion: 'rbac.authorization.k8s.io/v1',
  metadata: {
    name,
    namespace,
    uid,
    creationTimestamp: '2023-02-14T21:43:59Z',
    // No dashboard or known labels, but keep a generic label for testing if needed
    labels: { 'mock-label': 'true' },
  },
  subjects,
  roleRef: {
    apiGroup: 'rbac.authorization.k8s.io',
    kind: 'ClusterRole',
    name: roleRefName,
  },
});
