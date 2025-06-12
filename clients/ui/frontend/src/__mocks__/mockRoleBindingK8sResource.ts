import { RoleBindingKind, RoleBindingSubject } from 'mod-arch-shared';
import { genUID } from './mockUtils';

type MockResourceConfigType = {
  name?: string;
  namespace?: string;
  subjects?: RoleBindingSubject[];
  roleRefName?: string;
  uid?: string;
  modelRegistryName?: string;
};

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
  modelRegistryName = '',
}: MockResourceConfigType): RoleBindingKind => ({
  kind: 'RoleBinding',
  apiVersion: 'rbac.authorization.k8s.io/v1',
  metadata: {
    name,
    namespace,
    uid,
    creationTimestamp: '2023-02-14T21:43:59Z',
    labels: {
      'app.kubernetes.io/name': modelRegistryName,
      app: modelRegistryName,
      'app.kubernetes.io/component': 'model-registry',
      'app.kubernetes.io/part-of': 'model-registry',
      component: 'model-registry',
    },
  },
  subjects,
  roleRef: {
    apiGroup: 'rbac.authorization.k8s.io',
    kind: modelRegistryName ? 'Role' : 'ClusterRole',
    name: roleRefName,
  },
});
