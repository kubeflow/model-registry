import React from 'react';
import { useParams } from 'react-router-dom';
import {
  ModelRegistryKind,
  RoleBindingKind,
  useQueryParamNamespaces,
  GroupKind,
  FetchStateObject,
} from 'mod-arch-shared';
import { useGroups } from '~/app/hooks/useGroups';
import { useModelRegistryCR } from '~/app/hooks/useModelRegistryCR';
import useModelRegistryRoleBindings from '~/app/pages/modelRegistrySettings/useModelRegistryRoleBindings';
import { RoleBindingPermissionsRoleType } from '~/app/pages/settings/roleBinding/types';
import {
  createModelRegistryRoleBindingWrapper,
  deleteModelRegistryRoleBindingWrapper,
  createModelRegistryProjectRoleBinding,
  deleteModelRegistryProjectRoleBinding,
} from '~/app/pages/settings/roleBindingUtils';

/**
 * Configuration object returned by useModelRegistryPermissionsLogic hook
 */
export interface ModelRegistryPermissionsConfig {
  /** Current active tab index (0 = Users, 1 = Projects) */
  activeTabKey: number;
  /** Function to set the active tab */
  setActiveTabKey: (key: number) => void;
  /** Model registry owner reference for role bindings */
  ownerReference: ModelRegistryKind | undefined;
  /** Available groups for user role bindings */
  groups: GroupKind[];
  /** Filtered role bindings for the current model registry (users) */
  filteredRoleBindings: RoleBindingKind[];
  /** Filtered role bindings for projects */
  filteredProjectRoleBindings: RoleBindingKind[];
  /** Query parameters including namespace */
  queryParams: ReturnType<typeof useQueryParamNamespaces>;
  /** Model registry name from URL params */
  mrName: string | undefined;
  /** Model registry namespace (currently hardcoded) */
  modelRegistryNamespace: string;
  /** Full role bindings fetch state object */
  roleBindings: FetchStateObject<RoleBindingKind[]>;

  // Permission options for both tabs
  /** Permission options for the Users tab */
  userPermissionOptions: Array<{
    type: RoleBindingPermissionsRoleType;
    description: string;
  }>;
  /** Permission options for the Projects tab */
  projectPermissionOptions: Array<{
    type: RoleBindingPermissionsRoleType;
    description: string;
  }>;

  // Role binding functions
  /** Function to create user role bindings */
  createUserRoleBinding: typeof createModelRegistryRoleBindingWrapper;
  /** Function to delete user role bindings */
  deleteUserRoleBinding: typeof deleteModelRegistryRoleBindingWrapper;
  /** Function to create project role bindings */
  createProjectRoleBinding: typeof createModelRegistryProjectRoleBinding;
  /** Function to delete project role bindings */
  deleteProjectRoleBinding: typeof deleteModelRegistryProjectRoleBinding;

  // Role ref configurations
  /** Role reference name for user permissions */
  userRoleRefName: string;
  /** Role reference name for project permissions */
  projectRoleRefName: string;

  // Error/redirect flags
  /** Whether to show error state due to missing namespace */
  shouldShowError: boolean;
  /** Whether to redirect due to no role bindings or model registry */
  shouldRedirect: boolean;
}

/**
 * Custom hook that encapsulates the shared logic for Model Registry permissions management.
 *
 * This hook provides all the necessary state, handlers, and configurations needed for both
 * PatternFly and Material-UI versions of the permissions management components.
 *
 * @returns {ModelRegistryPermissionsConfig} Configuration object with all necessary state and functions
 *
 * @example
 * ```tsx
 * const {
 *   activeTabKey,
 *   setActiveTabKey,
 *   ownerReference,
 *   groups,
 *   filteredRoleBindings,
 *   shouldShowError,
 *   shouldRedirect,
 * } = useModelRegistryPermissionsLogic();
 *
 * if (shouldShowError) {
 *   return <ErrorComponent />;
 * }
 *
 * if (shouldRedirect) {
 *   return <Navigate to="/modelRegistrySettings" replace />;
 * }
 * ```
 */
export const useModelRegistryPermissionsLogic = (): ModelRegistryPermissionsConfig => {
  const [activeTabKey, setActiveTabKey] = React.useState(0);
  const modelRegistryNamespace = 'model-registry'; // TODO: This is a placeholder
  const [ownerReference, setOwnerReference] = React.useState<ModelRegistryKind>();
  const queryParams = useQueryParamNamespaces();
  const [groups] = useGroups(queryParams);
  const roleBindings = useModelRegistryRoleBindings(queryParams);
  const { mrName } = useParams<{ mrName: string }>();
  const [modelRegistryCR, crLoaded] = useModelRegistryCR(modelRegistryNamespace, queryParams);

  // Filter role bindings based on model registry name
  const filteredRoleBindings = roleBindings.data.filter(
    (rb: RoleBindingKind) => rb.metadata.labels?.['app.kubernetes.io/name'] === mrName,
  );

  const filteredProjectRoleBindings = roleBindings.data.filter(
    (rb: RoleBindingKind) =>
      rb.metadata.labels?.['app.kubernetes.io/name'] === mrName &&
      rb.metadata.labels?.['app.kubernetes.io/component'] === 'model-registry-project-rbac',
  );

  // Update owner reference when model registry CR changes
  React.useEffect(() => {
    if (modelRegistryCR) {
      setOwnerReference(modelRegistryCR);
    } else {
      setOwnerReference(undefined);
    }
  }, [modelRegistryCR]);

  // Permission configurations
  const userPermissionOptions = [
    {
      type: RoleBindingPermissionsRoleType.DEFAULT,
      description: 'Default role for all users',
    },
  ];

  const projectPermissionOptions = [
    {
      type: RoleBindingPermissionsRoleType.DEFAULT,
      description: 'Default project access role',
    },
  ];

  // Role ref names
  const userRoleRefName = `registry-user-${mrName ?? ''}`;
  const projectRoleRefName = `registry-project-${mrName ?? ''}`;

  // Error handling flags
  const shouldShowError = !queryParams.namespace;
  const shouldRedirect =
    (roleBindings.loaded && filteredRoleBindings.length === 0) || (crLoaded && !modelRegistryCR);

  return {
    activeTabKey,
    setActiveTabKey,
    ownerReference,
    groups,
    filteredRoleBindings,
    filteredProjectRoleBindings,
    queryParams,
    mrName,
    modelRegistryNamespace,
    roleBindings,
    userPermissionOptions,
    projectPermissionOptions,
    createUserRoleBinding: createModelRegistryRoleBindingWrapper,
    deleteUserRoleBinding: deleteModelRegistryRoleBindingWrapper,
    createProjectRoleBinding: createModelRegistryProjectRoleBinding,
    deleteProjectRoleBinding: deleteModelRegistryProjectRoleBinding,
    userRoleRefName,
    projectRoleRefName,
    shouldShowError,
    shouldRedirect,
  };
};
