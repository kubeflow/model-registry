import {
  OwnerReference,
  k8sCreateResource,
  k8sDeleteResource,
  k8sGetResource,
  k8sListResource,
  k8sPatchResource,
  K8sStatus,
  K8sResourceCommon,
} from 'mod-arch-shared';
import { RoleBindingPermissionsRoleType } from '~/app/pages/settings/roleBinding/types';
import { addOwnerReference } from '~/app/api/k8sUtils';
import {
  K8sAPIOptions,
  KnownLabels,
  RoleBindingKind,
  RoleBindingRoleRef,
  RoleBindingSubject,
} from '~/app/k8sTypes';
import { RoleBindingModel } from '~/app/api/models';
import { genRandomChars } from '~/app/utils/string';

export const generateRoleBindingServiceAccount = (
  name: string,
  serviceAccountName: string,
  roleRef: Omit<RoleBindingRoleRef, 'apiGroup'>,
  namespace: string,
): RoleBindingKind => {
  const roleBindingObject: RoleBindingKind = {
    apiVersion: 'rbac.authorization.k8s.io/v1',
    kind: 'RoleBinding',
    metadata: {
      name,
      namespace,
      labels: {
        [KnownLabels.DASHBOARD_RESOURCE]: 'true',
      },
    },
    roleRef: {
      apiGroup: 'rbac.authorization.k8s.io',
      ...roleRef,
    },
    subjects: [
      {
        kind: 'ServiceAccount',
        name: serviceAccountName,
      },
    ],
  };
  return roleBindingObject;
};

export const generateRoleBindingPermissions = (
  namespace: string,
  rbSubjectKind: RoleBindingSubject['kind'],
  rbSubjectName: RoleBindingSubject['name'],
  rbRoleRefName: RoleBindingPermissionsRoleType | string, //string because with MR this can include MR name
  rbRoleRefKind: RoleBindingRoleRef['kind'],
  rbLabels: { [key: string]: string } = {
    [KnownLabels.DASHBOARD_RESOURCE]: 'true',
    [KnownLabels.PROJECT_SHARING]: 'true',
  },
  ownerReference?: K8sResourceCommon,
): RoleBindingKind => {
  const roleBindingObject: RoleBindingKind = {
    apiVersion: 'rbac.authorization.k8s.io/v1',
    kind: 'RoleBinding',
    metadata: {
      name: `dashboard-permissions-${genRandomChars()}`,
      namespace,
      labels: rbLabels,
    },
    roleRef: {
      apiGroup: 'rbac.authorization.k8s.io',
      kind: rbRoleRefKind,
      name: rbRoleRefName,
    },
    subjects: [
      {
        apiGroup: 'rbac.authorization.k8s.io',
        kind: rbSubjectKind,
        name: rbSubjectName,
      },
    ],
  };
  return addOwnerReference(roleBindingObject, ownerReference);
};

export const listRoleBindings = (
  namespace?: string,
  labelSelector?: string,
): Promise<RoleBindingKind[]> => {
  const queryOptions = {
    ...(namespace && { ns: namespace }),
    ...(labelSelector && { queryParams: { labelSelector } }),
  };
  return k8sListResource<RoleBindingKind>({
    model: RoleBindingModel,
    queryOptions,
  }).then((listResource) => listResource.items);
};

export const getRoleBinding = (projectName: string, rbName: string): Promise<RoleBindingKind> =>
  k8sGetResource({
    model: RoleBindingModel,
    queryOptions: { name: rbName, ns: projectName },
  });

export const createRoleBinding = (
  data: RoleBindingKind,
  opts?: K8sAPIOptions,
): Promise<RoleBindingKind> =>
  k8sCreateResource<RoleBindingKind>({ model: RoleBindingModel, resource: data, ...opts });

export const deleteRoleBinding = (
  rbName: string,
  namespace: string,
  opts?: K8sAPIOptions,
): Promise<K8sStatus> =>
  k8sDeleteResource<RoleBindingKind, K8sStatus>({
    model: RoleBindingModel,
    queryOptions: { name: rbName, ns: namespace },
    ...opts,
  });

export const patchRoleBindingOwnerRef = (
  rbName: string,
  namespace: string,
  ownerReferences: OwnerReference[],
  opts?: K8sAPIOptions,
): Promise<RoleBindingKind> =>
  k8sPatchResource<RoleBindingKind>({
    model: RoleBindingModel,
    queryOptions: { name: rbName, ns: namespace },
    patches: [
      {
        op: 'replace',
        path: '/metadata/ownerReferences',
        value: ownerReferences,
      },
    ],
    ...opts,
  });

export const patchRoleBindingSubjects = (
  rbName: string,
  namespace: string,
  subjects: RoleBindingSubject[],
  opts?: K8sAPIOptions,
): Promise<RoleBindingKind> =>
  k8sPatchResource<RoleBindingKind>({
    model: RoleBindingModel,
    queryOptions: { name: rbName, ns: namespace },
    patches: [
      {
        op: 'replace',
        path: '/subjects',
        value: subjects,
      },
    ],
    ...opts,
  });
