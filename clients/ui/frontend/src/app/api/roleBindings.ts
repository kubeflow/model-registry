import {
  genRandomChars,
  K8sResourceCommon,
  RoleBindingKind,
  RoleBindingRoleRef,
  RoleBindingSubject,
} from 'mod-arch-shared';
import { RoleBindingPermissionsRoleType } from '~/app/pages/settings/roleBinding/types';

export const addOwnerReference = <R extends K8sResourceCommon>(
  resource: R,
  owner?: K8sResourceCommon,
  blockOwnerDeletion = false,
): R => {
  if (!owner) {
    return resource;
  }
  const ownerReferences = resource.metadata?.ownerReferences || [];
  if (
    owner.metadata?.uid &&
    owner.metadata.name &&
    !ownerReferences.find((r) => r.uid === owner.metadata?.uid)
  ) {
    ownerReferences.push({
      uid: owner.metadata.uid,
      name: owner.metadata.name,
      apiVersion: owner.apiVersion,
      kind: owner.kind,
      blockOwnerDeletion,
    });
  }
  return {
    ...resource,
    metadata: {
      ...resource.metadata,
      ownerReferences,
    },
  };
};

export const generateRoleBindingPermissions = (
  namespace: string,
  rbSubjectKind: RoleBindingSubject['kind'],
  rbSubjectName: RoleBindingSubject['name'],
  rbRoleRefName: RoleBindingPermissionsRoleType | string, //string because with MR this can include MR name
  rbRoleRefKind: RoleBindingRoleRef['kind'],
  rbLabels?: { [key: string]: string },
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
