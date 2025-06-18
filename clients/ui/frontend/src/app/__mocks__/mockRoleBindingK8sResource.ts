import { genUID } from '~/app/__mocks__/mockUtils';
import { KnownLabels, RoleBindingKind, RoleBindingSubject } from '~/app/k8sTypes';

type MockResourceConfigType = {
  name?: string;
  namespace?: string;
  subjects?: RoleBindingSubject[];
  roleRefName?: string;
  uid?: string;
  modelRegistryName?: string;
  isProjectSubject?: boolean;
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
  isProjectSubject = false,
  modelRegistryName = '',
}: MockResourceConfigType): RoleBindingKind => {
  let labels;
  if (modelRegistryName) {
    labels = {
      'app.kubernetes.io/name': modelRegistryName,
      app: modelRegistryName,
      'app.kubernetes.io/component': 'model-registry',
      'app.kubernetes.io/part-of': 'model-registry',
      [KnownLabels.DASHBOARD_RESOURCE]: 'true',
      component: 'model-registry',
      ...(isProjectSubject && { [KnownLabels.PROJECT_SUBJECT]: 'true' }),
    };
  } else {
    labels = {
      [KnownLabels.DASHBOARD_RESOURCE]: 'true',
    };
  }
  return {
    kind: 'RoleBinding',
    apiVersion: 'rbac.authorization.k8s.io/v1',
    metadata: {
      name,
      namespace,
      uid,
      creationTimestamp: '2023-02-14T21:43:59Z',
      labels,
    },
    subjects,
    roleRef: {
      apiGroup: 'rbac.authorization.k8s.io',
      kind: modelRegistryName ? 'Role' : 'ClusterRole',
      name: roleRefName,
    },
  };
};
