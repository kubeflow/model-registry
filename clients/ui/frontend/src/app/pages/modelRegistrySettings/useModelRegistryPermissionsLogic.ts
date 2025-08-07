import React from 'react';
import { useParams } from 'react-router-dom';
import { GroupKind, RoleBindingKind, FetchStateObject, ModelRegistryKind } from 'mod-arch-shared';
import { useQueryParamNamespaces } from 'mod-arch-core';
import { useGroups } from '~/app/hooks/useGroups';
import useModelRegistryRoleBindings from '~/app/pages/modelRegistrySettings/useModelRegistryRoleBindings';
import { useModelRegistryCR } from '~/app/hooks/useModelRegistryCR';
import { RoleBindingPermissionsRoleType } from '~/app/pages/settings/roleBinding/types';
import {
  createModelRegistryRoleBindingWrapper,
  deleteModelRegistryRoleBindingWrapper,
  createModelRegistryNamespaceRoleBinding,
  deleteModelRegistryNamespaceRoleBinding,
} from '~/app/pages/settings/roleBindingUtils';

export interface ModelRegistryPermissionsConfig {
  activeTabKey: number;
  setActiveTabKey: (key: number) => void;
  ownerReference: ModelRegistryKind | undefined;
  groups: GroupKind[];
  filteredRoleBindings: RoleBindingKind[];
  filteredNamespaceRoleBindings: RoleBindingKind[];
  queryParams: ReturnType<typeof useQueryParamNamespaces>;
  mrName: string | undefined;
  modelRegistryNamespace: string;
  roleBindings: FetchStateObject<RoleBindingKind[]>;
  userPermissionOptions: Array<{
    type: RoleBindingPermissionsRoleType;
    description: string;
  }>;
  namespacePermissionOptions: Array<{
    type: RoleBindingPermissionsRoleType;
    description: string;
  }>;
  createUserRoleBinding: typeof createModelRegistryRoleBindingWrapper;
  deleteUserRoleBinding: typeof deleteModelRegistryRoleBindingWrapper;
  createNamespaceRoleBinding: typeof createModelRegistryNamespaceRoleBinding;
  deleteNamespaceRoleBinding: typeof deleteModelRegistryNamespaceRoleBinding;
  userRoleRefName: string;
  namespaceRoleRefName: string;
  shouldShowError: boolean;
  shouldRedirect: boolean;
}

export const useModelRegistryPermissionsLogic = (): ModelRegistryPermissionsConfig => {
  const [activeTabKey, setActiveTabKey] = React.useState(0);
  const modelRegistryNamespace = 'model-registry'; // TODO: This is a placeholder
  const [ownerReference, setOwnerReference] = React.useState<ModelRegistryKind>();
  const queryParams = useQueryParamNamespaces();
  const [groups] = useGroups(queryParams);
  const roleBindings = useModelRegistryRoleBindings(queryParams);
  const { mrName } = useParams<{ mrName: string }>();
  const [modelRegistryCR, crLoaded] = useModelRegistryCR(modelRegistryNamespace, queryParams);

  const filteredRoleBindings = roleBindings.data.filter(
    (rb: RoleBindingKind) => rb.metadata.labels?.['app.kubernetes.io/name'] === mrName,
  );

  const filteredNamespaceRoleBindings = roleBindings.data.filter(
    (rb: RoleBindingKind) =>
      rb.metadata.labels?.['app.kubernetes.io/name'] === mrName &&
      rb.metadata.labels?.['app.kubernetes.io/component'] === 'model-registry-namespace-rbac',
  );

  React.useEffect(() => {
    if (modelRegistryCR) {
      setOwnerReference(modelRegistryCR);
    } else {
      setOwnerReference(undefined);
    }
  }, [modelRegistryCR]);

  const userPermissionOptions = [
    {
      type: RoleBindingPermissionsRoleType.DEFAULT,
      description: 'Default role for all users',
    },
  ];

  const namespacePermissionOptions = [
    {
      type: RoleBindingPermissionsRoleType.DEFAULT,
      description: 'Default namespace access role',
    },
  ];

  const userRoleRefName = `registry-user-${mrName ?? ''}`;
  const namespaceRoleRefName = `registry-namespace-${mrName ?? ''}`;

  const shouldShowError = !queryParams.namespace;
  const shouldRedirect =
    (roleBindings.loaded && filteredRoleBindings.length === 0) || (crLoaded && !modelRegistryCR);

  return {
    activeTabKey,
    setActiveTabKey,
    ownerReference,
    groups,
    filteredRoleBindings,
    filteredNamespaceRoleBindings,
    queryParams,
    mrName,
    modelRegistryNamespace,
    roleBindings,
    userPermissionOptions,
    namespacePermissionOptions,
    createUserRoleBinding: createModelRegistryRoleBindingWrapper,
    deleteUserRoleBinding: deleteModelRegistryRoleBindingWrapper,
    createNamespaceRoleBinding: createModelRegistryNamespaceRoleBinding,
    deleteNamespaceRoleBinding: deleteModelRegistryNamespaceRoleBinding,
    userRoleRefName,
    namespaceRoleRefName,
    shouldShowError,
    shouldRedirect,
  };
};
